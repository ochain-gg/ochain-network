package transactions

import (
	"fmt"

	abcitypes "github.com/cometbft/cometbft/abci/types"
	"github.com/fxamacker/cbor/v2"
	"github.com/ochain-gg/ochain-network/types"
)

type UniverseAccountDepositTransactionData struct {
	Universe string `cbor:"1,keyasint"`
	Planet   string `cbor:"2,keyasint"`
	Amount   uint64 `cbor:"3,keyasint"`
}

type UniverseAccountDepositTransaction struct {
	Type      TransactionType                       `cbor:"1,keyasint"`
	From      string                                `cbor:"2,keyasint"`
	Nonce     uint64                                `cbor:"3,keyasint"`
	Data      UniverseAccountDepositTransactionData `cbor:"4,keyasint"`
	Signature []byte                                `cbor:"5,keyasint"`
}

func (tx *UniverseAccountDepositTransaction) Transaction() (Transaction, error) {
	txData, err := cbor.Marshal(tx.Data)
	if err != nil {
		return Transaction{}, err
	}

	return Transaction{
		Type:      tx.Type,
		From:      tx.From,
		Nonce:     tx.Nonce,
		Data:      txData,
		Signature: tx.Signature,
	}, nil
}

func (tx *UniverseAccountDepositTransaction) Check(ctx TransactionContext) *abcitypes.ResponseCheckTx {
	globalAccount, err := ctx.Db.GlobalsAccounts.GetAt(tx.From, uint64(ctx.Date.Unix()))
	if err != nil {
		return &abcitypes.ResponseCheckTx{
			Code: types.InvalidTransactionError,
		}
	}

	_, err = ctx.Db.UniverseAccounts.GetAt(tx.Data.Universe, tx.From, uint64(ctx.Date.Unix()))
	if err != nil {
		return &abcitypes.ResponseCheckTx{
			Code: types.InvalidTransactionError,
		}
	}

	_, err = ctx.Db.Planets.GetAt(tx.Data.Universe, tx.Data.Planet, uint64(ctx.Date.Unix()))
	if err != nil {
		return &abcitypes.ResponseCheckTx{
			Code: types.InvalidTransactionError,
		}
	}

	if globalAccount.TokenBalance < tx.Data.Amount {
		return &abcitypes.ResponseCheckTx{
			Code: types.InvalidTransactionError,
		}
	}

	year, week := ctx.Date.ISOWeek()
	weeklyUsage, err := ctx.Db.UniverseAccountWeeklyUsage.GetAt(tx.Data.Universe, tx.From, year, week, uint64(ctx.Date.Unix()))
	if err != nil {
		return &abcitypes.ResponseCheckTx{
			Code: types.InvalidTransactionError,
		}
	}

	newDepositWeeklyUsage := weeklyUsage.DepositedAmount + tx.Data.Amount

	//if account balance >= TwoWeekCommanderPrice
	if newDepositWeeklyUsage > types.MaxWeeklyOCTDeposit {
		return &abcitypes.ResponseCheckTx{
			Code: types.InvalidTransactionError,
		}
	}

	return nil
}

func (tx *UniverseAccountDepositTransaction) Execute(ctx TransactionContext) *abcitypes.ExecTxResult {
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
			Code: types.InvalidTransactionError,
		}
	}

	universe, err := ctx.Db.Universes.GetAt(tx.Data.Universe, currentDate)
	if err != nil {
		return &abcitypes.ExecTxResult{
			Code: types.InvalidTransactionError,
		}
	}

	account, err := ctx.Db.UniverseAccounts.GetAt(universe.Id, tx.From, currentDate)
	if err != nil {
		return &abcitypes.ExecTxResult{
			Code: types.InvalidTransactionError,
		}
	}

	planet, err := ctx.Db.Planets.GetAt(tx.Data.Universe, tx.Data.Planet, currentDate)
	if err != nil {
		return &abcitypes.ExecTxResult{
			Code: types.InvalidTransactionError,
		}
	}
	year, week := ctx.Date.ISOWeek()
	weeklyUsage, err := ctx.Db.UniverseAccountWeeklyUsage.GetAt(tx.Data.Universe, tx.From, year, week, uint64(ctx.Date.Unix()))
	if err != nil {
		return &abcitypes.ExecTxResult{
			Code: types.InvalidTransactionError,
		}
	}

	newDepositWeeklyUsage := weeklyUsage.DepositedAmount + tx.Data.Amount
	if newDepositWeeklyUsage > types.MaxWeeklyOCTDeposit {
		return &abcitypes.ExecTxResult{
			Code: types.InvalidTransactionError,
		}
	}

	weeklyUsage.DepositedAmount = newDepositWeeklyUsage
	globalAccount.CreditBalance -= tx.Data.Amount

	planet.UpdateResources(universe.Speed, ctx.Date.Unix(), account)
	planet.Resources.OCT += tx.Data.Amount

	err = ctx.Db.GlobalsAccounts.Update(globalAccount)
	if err != nil {
		return &abcitypes.ExecTxResult{
			Code: types.InvalidTransactionError,
		}
	}

	err = ctx.Db.Planets.Update(tx.Data.Universe, planet)
	if err != nil {
		return &abcitypes.ExecTxResult{
			Code: types.InvalidTransactionError,
		}
	}

	err = ctx.Db.UniverseAccountWeeklyUsage.Upsert(weeklyUsage)
	if err != nil {
		return &abcitypes.ExecTxResult{
			Code: types.InvalidTransactionError,
		}
	}

	txGasCost, err := globalAccount.ApplyGasCost(uint64(ctx.Date.Unix()))
	if err != nil {
		return &abcitypes.ExecTxResult{
			Code: types.GasCostHigherThanBalance,
		}
	}

	receipt := TransactionReceipt{
		GasCost: txGasCost,
	}

	events := []abcitypes.Event{
		{
			Type: "UniverseAccountDeposit",
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

func ParseUniverseAccountDepositTransaction(tx Transaction) (UniverseAccountDepositTransaction, error) {
	var txData UniverseAccountDepositTransactionData
	err := cbor.Unmarshal(tx.Data, &txData)

	if err != nil {
		return UniverseAccountDepositTransaction{}, err
	}

	return UniverseAccountDepositTransaction{
		Type:      tx.Type,
		From:      tx.From,
		Nonce:     tx.Nonce,
		Data:      txData,
		Signature: tx.Signature,
	}, nil
}
