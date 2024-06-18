package transactions

import (
	"fmt"

	abcitypes "github.com/cometbft/cometbft/abci/types"
	"github.com/fxamacker/cbor/v2"
	"github.com/ochain-gg/ochain-network/types"
)

type SwapResourcesTransactionData struct {
	Universe string                 `cbor:"1,keyasint"`
	Planet   string                 `cbor:"2,keyasint"`
	From     types.MarketResourceID `cbor:"3,keyasint"`
	To       types.MarketResourceID `cbor:"4,keyasint"`
	Amount   uint64                 `cbor:"5,keyasint"`
}

type SwapResourcesTransaction struct {
	Type      TransactionType              `cbor:"1,keyasint"`
	From      string                       `cbor:"2,keyasint"`
	Nonce     uint64                       `cbor:"3,keyasint"`
	Data      SwapResourcesTransactionData `cbor:"4,keyasint"`
	Signature []byte                       `cbor:"5,keyasint"`
}

func (tx *SwapResourcesTransaction) Transaction() (Transaction, error) {
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

func (tx *SwapResourcesTransaction) Check(ctx TransactionContext) *abcitypes.ResponseCheckTx {

	account, err := ctx.Db.UniverseAccounts.GetAt(tx.Data.Universe, tx.From, uint64(ctx.Date.Unix()))
	if err != nil {
		return &abcitypes.ResponseCheckTx{
			Code: types.InvalidTransactionError,
		}
	}

	universe, err := ctx.Db.Universes.GetAt(tx.Data.Universe, uint64(ctx.Date.Unix()))
	if err != nil {
		return &abcitypes.ResponseCheckTx{
			Code: types.InvalidTransactionError,
		}
	}

	planet, err := ctx.Db.Planets.GetAt(tx.Data.Universe, tx.Data.Planet, uint64(ctx.Date.Unix()))
	if err != nil {
		return &abcitypes.ResponseCheckTx{
			Code: types.InvalidTransactionError,
		}
	}

	planet.UpdateResources(universe.Speed, int64(ctx.Date.Unix()), account)

	switch tx.Data.From {
	case types.OCTResourceID:
		if planet.Resources.OCT < tx.Data.Amount {
			return &abcitypes.ResponseCheckTx{
				Code: types.InvalidTransactionError,
			}
		}
	case types.MetalResourceID:
		if planet.Resources.Metal < tx.Data.Amount {
			return &abcitypes.ResponseCheckTx{
				Code: types.InvalidTransactionError,
			}
		}
	case types.CrystalResourceID:
		if planet.Resources.Crystal < tx.Data.Amount {
			return &abcitypes.ResponseCheckTx{
				Code: types.InvalidTransactionError,
			}
		}
	case types.DeuteriumResourceID:
		if planet.Resources.Deuterium < tx.Data.Amount {
			return &abcitypes.ResponseCheckTx{
				Code: types.InvalidTransactionError,
			}
		}
	default:
		return &abcitypes.ResponseCheckTx{
			Code: types.InvalidTransactionError,
		}
	}

	return nil
}

func (tx *SwapResourcesTransaction) Execute(ctx TransactionContext) *abcitypes.ExecTxResult {
	result := tx.Check(ctx)
	if result.Code != types.NoError {
		return &abcitypes.ExecTxResult{
			Code: result.GetCode(),
		}
	}

	globalAccount, err := ctx.Db.GlobalsAccounts.GetAt(tx.From, uint64(ctx.Date.Unix()))
	if err != nil {
		return &abcitypes.ExecTxResult{
			Code: types.InvalidTransactionError,
		}
	}

	account, err := ctx.Db.UniverseAccounts.GetAt(tx.Data.Universe, tx.From, uint64(ctx.Date.Unix()))
	if err != nil {
		return &abcitypes.ExecTxResult{
			Code: types.InvalidTransactionError,
		}
	}

	universe, err := ctx.Db.Universes.GetAt(tx.Data.Universe, uint64(ctx.Date.Unix()))
	if err != nil {
		return &abcitypes.ExecTxResult{
			Code: types.InvalidTransactionError,
		}
	}

	planet, err := ctx.Db.Planets.GetAt(tx.Data.Universe, tx.Data.Planet, uint64(ctx.Date.Unix()))
	if err != nil {
		return &abcitypes.ExecTxResult{
			Code: types.InvalidTransactionError,
		}
	}

	market, err := ctx.Db.ResourcesMarket.GetAt(tx.Data.Universe, uint64(ctx.Date.Unix()))
	if err != nil {
		return &abcitypes.ExecTxResult{
			Code: types.InvalidTransactionError,
		}
	}

	planet.UpdateResources(universe.Speed, ctx.Date.Unix(), account)

	amountOut, err := market.SwapResources(tx.Data.From, tx.Data.To, tx.Data.Amount)
	if err != nil {
		return &abcitypes.ExecTxResult{
			Code: types.InvalidTransactionError,
		}
	}

	err = planet.RemoveResourceById(tx.Data.From, tx.Data.Amount)
	if err != nil {
		return &abcitypes.ExecTxResult{
			Code: types.InvalidTransactionError,
		}
	}

	planet.AddResourceById(tx.Data.To, amountOut)

	txGasCost, err := globalAccount.ApplyGasCost(uint64(ctx.Date.Unix()))
	if err != nil {
		return &abcitypes.ExecTxResult{
			Code: types.GasCostHigherThanBalance,
		}
	}

	err = ctx.Db.GlobalsAccounts.Update(globalAccount)
	if err != nil {
		return &abcitypes.ExecTxResult{
			Code: types.InvalidTransactionError,
		}
	}

	err = ctx.Db.Planets.Update(universe.Id, planet)
	if err != nil {
		return &abcitypes.ExecTxResult{
			Code: types.InvalidTransactionError,
		}
	}

	events := []abcitypes.Event{
		{
			Type: "ResourcesSwaped",
			Attributes: []abcitypes.EventAttribute{
				{Key: "account", Value: tx.From, Index: true},
				{Key: "universe", Value: tx.Data.Universe, Index: true},
				{Key: "planet", Value: tx.Data.Planet, Index: true},
				{Key: "from", Value: fmt.Sprint(tx.Data.From)},
				{Key: "to", Value: fmt.Sprint(tx.Data.To)},
				{Key: "amountIn", Value: fmt.Sprint(tx.Data.Amount)},
				{Key: "amountOut", Value: fmt.Sprint(amountOut)},
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

	receipt := TransactionReceipt{
		GasCost: txGasCost,
	}

	return &abcitypes.ExecTxResult{
		Code:      types.NoError,
		Events:    events,
		GasUsed:   100,
		GasWanted: 100,
		Data:      receipt.Bytes(),
	}
}

func ParseSwapResourcesTransaction(tx Transaction) (SwapResourcesTransaction, error) {
	var txData SwapResourcesTransactionData
	err := cbor.Unmarshal(tx.Data, &txData)

	if err != nil {
		return SwapResourcesTransaction{}, err
	}

	return SwapResourcesTransaction{
		Type:      tx.Type,
		From:      tx.From,
		Nonce:     tx.Nonce,
		Data:      txData,
		Signature: tx.Signature,
	}, nil
}
