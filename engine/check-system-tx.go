package engine

import (
	abcitypes "github.com/cometbft/cometbft/abci/types"
	"github.com/ochain-gg/ochain-network/transactions"
	game_transactions "github.com/ochain-gg/ochain-network/transactions/game"
	validator_transactions "github.com/ochain-gg/ochain-network/transactions/validator"
	"github.com/ochain-gg/ochain-network/types"
)

func CheckSystemTx(ctx transactions.TransactionContext, req *abcitypes.CheckTxRequest, tx transactions.Transaction) *abcitypes.CheckTxResponse {

	switch tx.Type {
	case transactions.NewValidator:
		transaction, err := validator_transactions.ParseOChainNewValidatorTransaction(tx)
		if err != nil {
			return &abcitypes.CheckTxResponse{Code: types.ParsingTransactionDataError, GasWanted: 0, GasUsed: 0}
		}

		return transaction.Check(ctx)

	case transactions.RemoveValidator:
		transaction, err := validator_transactions.ParseOChainRemoveValidatorTransaction(tx)
		if err != nil {
			return &abcitypes.CheckTxResponse{Code: types.ParsingTransactionDataError, GasWanted: 0, GasUsed: 0}
		}

		return transaction.Check(ctx)

	case transactions.OChainTokenDeposit:
		transaction, err := validator_transactions.ParseOChainTokenDepositTransaction(tx)
		if err != nil {
			return &abcitypes.CheckTxResponse{Code: types.ParsingTransactionDataError, GasWanted: 0, GasUsed: 0}
		}

		return transaction.Check(ctx)

	case transactions.OChainCreditDeposit:
		transaction, err := validator_transactions.ParseOChainCreditDepositTransaction(tx)
		if err != nil {
			return &abcitypes.CheckTxResponse{Code: types.ParsingTransactionDataError, GasWanted: 0, GasUsed: 0}
		}

		return transaction.Check(ctx)

	case transactions.ExecutePendingUpgrade:
		transaction, err := game_transactions.ParseExecuteUpgradeTransaction(tx)
		if err != nil {
			return &abcitypes.CheckTxResponse{Code: types.ParsingTransactionDataError, GasWanted: 0, GasUsed: 0}
		}

		return transaction.Check(ctx)

	}

	return &abcitypes.CheckTxResponse{Code: types.NotImplemented, GasWanted: 0, GasUsed: 0}
}
