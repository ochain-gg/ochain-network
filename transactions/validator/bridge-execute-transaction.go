package validator_transactions

import (
	abcitypes "github.com/cometbft/cometbft/abci/types"
	"github.com/fxamacker/cbor/v2"

	t "github.com/ochain-gg/ochain-network/transactions"
	"github.com/ochain-gg/ochain-network/types"
)

type OChainExecuteTransactionData struct {
	RemoteTransactionHash string `cbor:"1,keyasint"`
}

type OChainExecuteTransaction struct {
	Type      t.TransactionType            `cbor:"1,keyasint"`
	From      string                       `cbor:"2,keyasint"`
	Nonce     uint64                       `cbor:"3,keyasint"`
	Data      OChainExecuteTransactionData `cbor:"4,keyasint"`
	Signature []byte                       `cbor:"5,keyasint"`
}

func (tx OChainExecuteTransaction) Check(ctx t.TransactionContext) *abcitypes.CheckTxResponse {

	transaction, err := ctx.Db.BridgeTransactions.GetAt(tx.Data.RemoteTransactionHash, uint64(ctx.Date.Unix()))
	if err != nil || transaction.Executed || transaction.Canceled {
		return &abcitypes.CheckTxResponse{
			Code: types.InvalidTransactionError,
		}
	}

	_, err = ctx.Db.GlobalsAccounts.GetAt(tx.From, uint64(ctx.Date.Unix()))
	if err != nil {
		return &abcitypes.CheckTxResponse{
			Code: types.InvalidTransactionError,
		}
	}

	return &abcitypes.CheckTxResponse{
		Code: types.NoError,
	}

}

func (tx OChainExecuteTransaction) Execute(ctx t.TransactionContext) *abcitypes.ExecTxResult {
	result := tx.Check(ctx)
	if result.Code != types.NoError {
		return &abcitypes.ExecTxResult{
			Code: result.GetCode(),
		}
	}

	transaction, err := ctx.Db.BridgeTransactions.GetAt(tx.Data.RemoteTransactionHash, uint64(ctx.Date.Unix()))
	if err != nil || transaction.Executed || transaction.Canceled {
		return &abcitypes.ExecTxResult{
			Code: types.InvalidTransactionError,
		}
	}

	globalAccount, err := ctx.Db.GlobalsAccounts.GetAt(tx.From, uint64(ctx.Date.Unix()))
	if err != nil {
		return &abcitypes.ExecTxResult{
			Code: types.InvalidTransactionError,
		}
	}

	switch transaction.Type {

	case types.OChainTokenDepositTransaction:
		globalAccount.TokenBalance += transaction.Amount

	case types.OChainCreditDepositTransaction:
		globalAccount.CreditBalance += transaction.Amount

	default:
		return &abcitypes.ExecTxResult{
			Code: types.InvalidTransactionError,
		}
	}

	transaction.Executed = true

	err = ctx.Db.BridgeTransactions.Update(transaction)
	if err != nil {
		return &abcitypes.ExecTxResult{
			Code: types.InvalidTransactionError,
		}
	}

	err = ctx.Db.GlobalsAccounts.Update(globalAccount)
	if err != nil {
		return &abcitypes.ExecTxResult{
			Code: types.InvalidTransactionError,
		}
	}

	return &abcitypes.ExecTxResult{
		Code: types.NoError,
	}
}

func (tx OChainExecuteTransaction) Transaction() (t.Transaction, error) {

	txData, err := cbor.Marshal(tx.Data)
	if err != nil {
		return t.Transaction{}, err
	}

	return t.Transaction{
		Type: tx.Type,
		Data: txData,
	}, nil
}

func ParseOChainExecuteTransaction(tx t.Transaction) (OChainExecuteTransaction, error) {
	var txData OChainExecuteTransactionData
	err := cbor.Unmarshal(tx.Data, &txData)

	if err != nil {
		return OChainExecuteTransaction{}, err
	}

	return OChainExecuteTransaction{
		Type: tx.Type,
		Data: txData,
	}, nil
}
