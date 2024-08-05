package engine

import (
	abcitypes "github.com/cometbft/cometbft/abci/types"
	"github.com/ochain-gg/ochain-network/transactions"
	"github.com/ochain-gg/ochain-network/types"
)

func CheckSystemTx(ctx transactions.TransactionContext, req *abcitypes.RequestCheckTx, tx transactions.Transaction) *abcitypes.ResponseCheckTx {

	tx, err := transactions.ParseTransaction(req.Tx)
	if err != nil {
		return &abcitypes.ResponseCheckTx{Code: types.ParsingTransactionError, GasWanted: 0, GasUsed: 0}
	}

	err = tx.IsValid()
	if err != nil {
		return &abcitypes.ResponseCheckTx{Code: types.InvalidTransactionError, GasWanted: 0, GasUsed: 0}
	}

	switch tx.Type {
	case transactions.NewValidator:
		transaction, err := transactions.ParseNewValidatorTransaction(tx)
		if err != nil {
			return &abcitypes.ResponseCheckTx{Code: types.ParsingTransactionDataError, GasWanted: 0, GasUsed: 0}
		}

		return transaction.Check(ctx)

	case transactions.RemoveValidator:
		transaction, err := transactions.ParseRemoveValidatorTransaction(tx)
		if err != nil {
			return &abcitypes.ResponseCheckTx{Code: types.ParsingTransactionDataError, GasWanted: 0, GasUsed: 0}
		}

		return transaction.Check(ctx)

	case transactions.OChainTokenDeposit:
		transaction, err := transactions.ParseNewOChainTokenDepositTransaction(tx)
		if err != nil {
			return &abcitypes.ResponseCheckTx{Code: types.ParsingTransactionDataError, GasWanted: 0, GasUsed: 0}
		}

		return transaction.Check(ctx)

	case transactions.OChainCreditDeposit:
		transaction, err := transactions.ParseNewOChainCreditDepositTransaction(tx)
		if err != nil {
			return &abcitypes.ResponseCheckTx{Code: types.ParsingTransactionDataError, GasWanted: 0, GasUsed: 0}
		}

		return transaction.Check(ctx)

	}

	return &abcitypes.ResponseCheckTx{Code: types.NotImplemented, GasWanted: 0, GasUsed: 0}
}
