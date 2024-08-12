package engine

import (
	abcitypes "github.com/cometbft/cometbft/abci/types"
	"github.com/ochain-gg/ochain-network/types"
)

func CheckSystemTx(ctx types.TransactionContext, req *abcitypes.RequestCheckTx, tx types.Transaction) *abcitypes.ResponseCheckTx {

	tx, err := types.ParseTransaction(req.Tx)
	if err != nil {
		return &abcitypes.ResponseCheckTx{Code: types.ParsingTransactionError, GasWanted: 0, GasUsed: 0}
	}

	switch tx.Type {
	case types.NewValidator:
		transaction, err := types.ParseOChainBridgeNewValidatorTransaction(tx)
		if err != nil {
			return &abcitypes.ResponseCheckTx{Code: types.ParsingTransactionDataError, GasWanted: 0, GasUsed: 0}
		}

		return transaction.Check(ctx)

	case types.RemoveValidator:
		transaction, err := types.ParseOChainBridgeRemoveValidatorTransaction(tx)
		if err != nil {
			return &abcitypes.ResponseCheckTx{Code: types.ParsingTransactionDataError, GasWanted: 0, GasUsed: 0}
		}

		return transaction.Check(ctx)

	case types.OChainTokenDeposit:
		transaction, err := types.ParseOChainBridgeTokenDepositTransaction(tx)
		if err != nil {
			return &abcitypes.ResponseCheckTx{Code: types.ParsingTransactionDataError, GasWanted: 0, GasUsed: 0}
		}

		return transaction.Check(ctx)

	case types.OChainCreditDeposit:
		transaction, err := types.ParseOChainBridgeCreditDepositTransaction(tx)
		if err != nil {
			return &abcitypes.ResponseCheckTx{Code: types.ParsingTransactionDataError, GasWanted: 0, GasUsed: 0}
		}

		return transaction.Check(ctx)

	}

	return &abcitypes.ResponseCheckTx{Code: types.NotImplemented, GasWanted: 0, GasUsed: 0}
}
