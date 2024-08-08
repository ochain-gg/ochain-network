package transactions

import (
	abcitypes "github.com/cometbft/cometbft/abci/types"
	"github.com/fxamacker/cbor/v2"
	"github.com/ochain-gg/ochain-network/types"
)

type OChainBridgeExecuteTransactionData struct {
	RemoteTransactionHash string `cbor:"1,keyasint"`
}

type OChainBridgeExecuteTransaction struct {
	Type      TransactionType                    `cbor:"1,keyasint"`
	From      string                             `cbor:"2,keyasint"`
	Nonce     uint64                             `cbor:"3,keyasint"`
	Data      OChainBridgeExecuteTransactionData `cbor:"4,keyasint"`
	Signature []byte                             `cbor:"5,keyasint"`
}

func (tx OChainBridgeExecuteTransaction) Check(ctx TransactionContext) *abcitypes.ResponseCheckTx {

	transaction, err := ctx.Db.BridgeTransactions.GetAt(tx.Data.RemoteTransactionHash, uint64(ctx.Date.Unix()))
	if err != nil || transaction.Executed || transaction.Canceled {
		return &abcitypes.ResponseCheckTx{
			Code: types.InvalidTransactionError,
		}
	}

	_, err = ctx.Db.GlobalsAccounts.GetAt(tx.From, uint64(ctx.Date.Unix()))
	if err != nil {
		return &abcitypes.ResponseCheckTx{
			Code: types.InvalidTransactionError,
		}
	}

	return &abcitypes.ResponseCheckTx{
		Code: types.NoError,
	}

}

func (tx OChainBridgeExecuteTransaction) Execute(ctx TransactionContext) *abcitypes.ExecTxResult {
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

	case types.OChainBridgeTokenDepositTransaction:
		globalAccount.TokenBalance += transaction.Amount

	case types.OChainBridgeCreditDepositTransaction:
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

func (tx OChainBridgeExecuteTransaction) Transaction() (Transaction, error) {

	txData, err := cbor.Marshal(tx.Data)
	if err != nil {
		return Transaction{}, err
	}

	return Transaction{
		Type: tx.Type,
		Data: txData,
	}, nil
}

func ParseOChainBridgeExecuteTransaction(tx Transaction) (OChainBridgeExecuteTransaction, error) {
	var txData OChainBridgeExecuteTransactionData
	err := cbor.Unmarshal(tx.Data, &txData)

	if err != nil {
		return OChainBridgeExecuteTransaction{}, err
	}

	return OChainBridgeExecuteTransaction{
		Type: tx.Type,
		Data: txData,
	}, nil
}
