package governance_transactions

import (
	abcitypes "github.com/cometbft/cometbft/abci/types"
	"github.com/fxamacker/cbor/v2"
	t "github.com/ochain-gg/ochain-network/transactions"
	"github.com/ochain-gg/ochain-network/types"
)

type ExecuteOChainTokenUnstakeTransactionData struct {
	Amount string `cbor:"1,keyasint"`
}

type ExecuteOChainTokenUnstakeTransaction struct {
	Type      t.TransactionType                        `cbor:"1,keyasint"`
	From      string                                   `cbor:"2,keyasint"`
	Nonce     uint64                                   `cbor:"3,keyasint"`
	Data      ExecuteOChainTokenUnstakeTransactionData `cbor:"4,keyasint"`
	Signature []byte                                   `cbor:"5,keyasint"`
}

func (tx *ExecuteOChainTokenUnstakeTransaction) Check(ctx t.TransactionContext) *abcitypes.ResponseCheckTx {
	globalAccount, err := ctx.Db.GlobalsAccounts.Get(tx.From)
	if err != nil {
		return &abcitypes.ResponseCheckTx{
			Code: types.InvalidTransactionError,
		}
	}

	if globalAccount.StackedBalance < tx.Data.Amount {
		return &abcitypes.ResponseCheckTx{
			Code: types.InvalidTransactionError,
		}
	}

	return &abcitypes.ResponseCheckTx{
		Code: types.NoError,
	}
}

func (tx *ExecuteOChainTokenUnstakeTransaction) Execute(ctx t.TransactionContext) *abcitypes.ExecTxResult {
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

	globalAccount.TokenBalance -= tx.Data.Amount
	globalAccount.StackedBalance += tx.Data.Amount

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

	events := []abcitypes.Event{
		{
			Type: "OChainTokenStacked",
			Attributes: []abcitypes.EventAttribute{
				{Key: "account", Value: tx.From, Index: true},
				{Key: "amount", Value: string(tx.Data.Amount), Index: true},
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

func (tx *ExecuteOChainTokenUnstakeTransaction) Transaction() (t.Transaction, error) {
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

func ParseExecuteOChainTokenUnstakeTransaction(tx t.Transaction) (ExecuteOChainTokenUnstakeTransaction, error) {
	var txData ExecuteOChainTokenUnstakeTransactionData
	err := cbor.Unmarshal(tx.Data, &txData)

	if err != nil {
		return ExecuteOChainTokenUnstakeTransaction{}, err
	}

	return ExecuteOChainTokenUnstakeTransaction{
		Type:      tx.Type,
		From:      tx.From,
		Nonce:     tx.Nonce,
		Data:      txData,
		Signature: tx.Signature,
	}, nil
}