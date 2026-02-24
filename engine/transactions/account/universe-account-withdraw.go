package account_transactions

import (
	"fmt"

	abcitypes "github.com/cometbft/cometbft/abci/types"
	"github.com/fxamacker/cbor/v2"
	t "github.com/ochain-gg/ochain-network/engine/transactions"
	"github.com/ochain-gg/ochain-network/types"
)

type UniverseAccountWithdrawTransactionData struct {
	Universe string `cbor:"1,keyasint"`
	Planet   string `cbor:"2,keyasint"`
	Amount   uint64 `cbor:"3,keyasint"`
}

type UniverseAccountWithdrawTransaction struct {
	Type      t.TransactionType                      `cbor:"1,keyasint"`
	From      string                                 `cbor:"2,keyasint"`
	Nonce     uint64                                 `cbor:"3,keyasint"`
	Data      UniverseAccountWithdrawTransactionData `cbor:"4,keyasint"`
	Signature []byte                                 `cbor:"5,keyasint"`
}

func (tx *UniverseAccountWithdrawTransaction) GetType() t.TransactionType {
	return t.UniverseOChainTokenWithdraw
}

func (tx *UniverseAccountWithdrawTransaction) Transaction() (t.Transaction, error) {
	txData, err := cbor.Marshal(tx.Data)
	if err != nil {
		return t.Transaction{}, err
	}

	return t.Transaction{
		Type:      tx.Type,
		From:      tx.From,
		Nonce:     tx.Nonce,
		Data:      txData,
		Signature: tx.Signature,
	}, nil
}

func (tx *UniverseAccountWithdrawTransaction) Check(ctx t.TransactionContext) *abcitypes.CheckTxResponse {
	_, err := ctx.Db.GlobalsAccounts.GetAt(tx.From, uint64(ctx.Date.Unix()))
	if err != nil {
		return &abcitypes.CheckTxResponse{
			Code: types.GlobalAccountNotFoundError,
		}
	}

	_, err = ctx.Db.UniverseAccounts.GetAt(tx.Data.Universe, tx.From, uint64(ctx.Date.Unix()))
	if err != nil {
		return &abcitypes.CheckTxResponse{
			Code: types.UniverseAccountNotFoundError,
		}
	}

	planet, err := ctx.Db.Planets.GetAt(tx.Data.Universe, tx.Data.Planet, uint64(ctx.Date.Unix()))
	if err != nil {
		return &abcitypes.CheckTxResponse{
			Code: types.PlanetNotFoundError,
		}
	}

	year, week := ctx.Date.ISOWeek()
	weeklyUsage, err := ctx.Db.UniverseAccountWeeklyUsage.GetAt(tx.Data.Universe, tx.From, year, week, uint64(ctx.Date.Unix()))
	if err != nil {
		return &abcitypes.CheckTxResponse{
			Code: types.DatabaseOperationError,
		}
	}

	if weeklyUsage.WithdrawalsExecuted >= types.MaxWeeklyOCTWithdrawals {
		return &abcitypes.CheckTxResponse{
			Code: types.RateLimitExceededError,
		}
	}

	if planet.Resources.OCT < tx.Data.Amount {
		return &abcitypes.CheckTxResponse{
			Code: types.InsufficientBalanceError,
		}
	}

	return &abcitypes.CheckTxResponse{
		Code: types.NoError,
	}
}

func (tx *UniverseAccountWithdrawTransaction) Execute(ctx t.TransactionContext) *abcitypes.ExecTxResult {
	result := tx.Check(ctx)
	if result.Code != types.NoError {
		return &abcitypes.ExecTxResult{
			Code: result.GetCode(),
		}
	}

	currentDate := uint64(ctx.Date.Unix())

	globalAccount, err := ctx.Db.GlobalsAccounts.GetAt(tx.From, currentDate)
	if err != nil {
		return &abcitypes.ExecTxResult{
			Code: types.GlobalAccountNotFoundError,
		}
	}

	universe, err := ctx.Db.Universes.GetAt(tx.Data.Universe, currentDate)
	if err != nil {
		return &abcitypes.ExecTxResult{
			Code: types.UniverseNotFoundError,
		}
	}

	account, err := ctx.Db.UniverseAccounts.GetAt(universe.Id, tx.From, currentDate)
	if err != nil {
		return &abcitypes.ExecTxResult{
			Code: types.UniverseAccountNotFoundError,
		}
	}

	planet, err := ctx.Db.Planets.GetAt(tx.Data.Universe, tx.Data.Planet, currentDate)
	if err != nil {
		return &abcitypes.ExecTxResult{
			Code: types.PlanetNotFoundError,
		}
	}

	year, week := ctx.Date.ISOWeek()
	weeklyUsage, err := ctx.Db.UniverseAccountWeeklyUsage.GetAt(tx.Data.Universe, tx.From, year, week, uint64(ctx.Date.Unix()))
	if err != nil {
		return &abcitypes.ExecTxResult{
			Code: types.DatabaseOperationError,
		}
	}

	weeklyUsage.WithdrawalsExecuted += 1

	globalAccount.CreditBalance += tx.Data.Amount
	planet.Resources.OCT -= tx.Data.Amount
	planet.UpdateResources(universe.Speed, ctx.Date.Unix(), account)

	err = ctx.Db.GlobalsAccounts.Update(globalAccount)
	if err != nil {
		return &abcitypes.ExecTxResult{
			Code: types.DatabaseOperationError,
		}
	}

	err = ctx.Db.Planets.Update(tx.Data.Universe, planet)
	if err != nil {
		return &abcitypes.ExecTxResult{
			Code: types.DatabaseOperationError,
		}
	}

	err = ctx.Db.UniverseAccountWeeklyUsage.Upsert(weeklyUsage)
	if err != nil {
		return &abcitypes.ExecTxResult{
			Code: types.DatabaseOperationError,
		}
	}

	txGasCost, err := globalAccount.ApplyGasCost(uint64(ctx.Date.Unix()))
	if err != nil {
		return &abcitypes.ExecTxResult{
			Code: types.GasCostHigherThanBalance,
		}
	}

	receipt := t.TransactionReceipt{
		GasCost: txGasCost,
	}

	events := []abcitypes.Event{
		{
			Type: "UniverseAccountWithdraw",
			Attributes: []abcitypes.EventAttribute{
				{Key: "account", Value: tx.From, Index: true},
				{Key: "universe", Value: tx.Data.Universe, Index: true},
				{Key: "planet", Value: tx.Data.Planet, Index: true},
				{Key: "amount", Value: fmt.Sprint(tx.Data.Amount)},
			},
		},
		{
			Type: "PlanetResourcesUpdated",
			Attributes: []abcitypes.EventAttribute{
				{Key: "account", Value: tx.From, Index: true},
				{Key: "universe", Value: tx.Data.Universe, Index: true},
				{Key: "planet", Value: tx.Data.Planet, Index: true},
				{Key: "oct", Value: fmt.Sprint(planet.Resources.OCT)},
				{Key: "metal", Value: fmt.Sprint(planet.Resources.Metal)},
				{Key: "crystal", Value: fmt.Sprint(planet.Resources.Crystal)},
				{Key: "deuterium", Value: fmt.Sprint(planet.Resources.Deuterium)},
			},
		},
	}

	return &abcitypes.ExecTxResult{
		Code:      types.NoError,
		Events:    events,
		GasUsed:   100,
		GasWanted: 100,
		Data:      receipt.Bytes(),
	}
}

func ParseUniverseAccountWithdrawTransaction(tx t.Transaction) (*UniverseAccountWithdrawTransaction, error) {
	var txData UniverseAccountWithdrawTransactionData
	err := cbor.Unmarshal(tx.Data, &txData)

	if err != nil {
		return nil, err
	}

	return &UniverseAccountWithdrawTransaction{
		Type:      tx.Type,
		From:      tx.From,
		Nonce:     tx.Nonce,
		Data:      txData,
		Signature: tx.Signature,
	}, nil
}

// ParseUniverseAccountWithdrawTransactionHandler is a wrapper that returns ITransactionHandler for registry
func ParseUniverseAccountWithdrawTransactionHandler(tx t.Transaction) (t.ITransactionHandler, error) {
	return ParseUniverseAccountWithdrawTransaction(tx)
}
