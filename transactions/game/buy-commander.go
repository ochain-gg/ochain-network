package game_transactions

import (
	abcitypes "github.com/cometbft/cometbft/abci/types"
	"github.com/fxamacker/cbor/v2"

	t "github.com/ochain-gg/ochain-network/transactions"
	"github.com/ochain-gg/ochain-network/types"
)

type BuyCommanderTransactionData struct {
	Universe  string                  `cbor:"2,keyasint"`
	Commander types.OChainCommanderID `cbor:"3,keyasint"`
}

type BuyCommanderTransaction struct {
	Type      t.TransactionType           `cbor:"1,keyasint"`
	From      string                      `cbor:"2,keyasint"`
	Nonce     uint64                      `cbor:"3,keyasint"`
	Data      BuyCommanderTransactionData `cbor:"4,keyasint"`
	Signature []byte                      `cbor:"5,keyasint"`
}

func (tx *BuyCommanderTransaction) Transaction() (t.Transaction, error) {
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

func (tx *BuyCommanderTransaction) Check(ctx t.TransactionContext) *abcitypes.ResponseCheckTx {
	globalAccount, err := ctx.Db.GlobalsAccounts.Get(tx.From)
	if err != nil {
		return &abcitypes.ResponseCheckTx{
			Code: types.InvalidTransactionError,
		}
	}

	txGasCost := globalAccount.GetGasCost(uint64(ctx.Date.Unix()))
	if txGasCost > globalAccount.TokenBalance {
		return &abcitypes.ResponseCheckTx{
			Code: types.GasCostHigherThanBalance,
		}
	}

	_, err = ctx.Db.UniverseAccounts.Get(tx.Data.Universe, tx.From)
	if err == nil {
		return &abcitypes.ResponseCheckTx{
			Code: types.InvalidTransactionError,
		}
	}

	//if account balance >= TwoWeekCommanderPrice
	if globalAccount.CreditBalance < types.TwoWeekCommanderPrice {
		return &abcitypes.ResponseCheckTx{
			Code: types.InvalidTransactionError,
		}
	}

	return nil
}

func (tx *BuyCommanderTransaction) Execute(ctx t.TransactionContext) *abcitypes.ExecTxResult {
	response := tx.Check(ctx)
	if response.Code != types.NoError {
		return &abcitypes.ExecTxResult{
			Code: response.Code,
		}
	}

	globalAccount, err := ctx.Db.GlobalsAccounts.Get(tx.From)
	if err == nil {
		return &abcitypes.ExecTxResult{
			Code: types.InvalidTransactionError,
		}
	}

	account, err := ctx.Db.UniverseAccounts.Get(tx.Data.Universe, tx.From)
	if err == nil {
		return &abcitypes.ExecTxResult{
			Code: types.InvalidTransactionError,
		}
	}

	globalAccount.CreditBalance -= types.TwoWeekCommanderPrice
	account.SubscribeToCommander(ctx.Date.Unix(), tx.Data.Commander)

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

	err = ctx.Db.UniverseAccounts.Update(account)
	if err != nil {
		return &abcitypes.ExecTxResult{
			Code: types.InvalidTransactionError,
		}
	}

	events := []abcitypes.Event{
		{
			Type: "CommanderSubscribed",
			Attributes: []abcitypes.EventAttribute{
				{Key: "account", Value: tx.From, Index: true},
				{Key: "universe", Value: tx.Data.Universe, Index: true},
				{Key: "commander", Value: string(tx.Data.Commander), Index: true},
			},
		},
	}

	receipt := t.TransactionReceipt{
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

func ParseBuyCommanderTransaction(tx t.Transaction) (BuyCommanderTransaction, error) {
	var txData BuyCommanderTransactionData
	err := cbor.Unmarshal(tx.Data, &txData)

	if err != nil {
		return BuyCommanderTransaction{}, err
	}

	return BuyCommanderTransaction{
		Type:      tx.Type,
		From:      tx.From,
		Nonce:     tx.Nonce,
		Data:      txData,
		Signature: tx.Signature,
	}, nil
}
