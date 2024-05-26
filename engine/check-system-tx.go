package engine

import (
	"log"

	abcitypes "github.com/cometbft/cometbft/abci/types"
	"github.com/ochain-gg/ochain-network/transactions"
	"github.com/ochain-gg/ochain-network/types"
)

func CheckSystemTx(ctx transactions.TransactionContext, req *abcitypes.RequestCheckTx, tx transactions.Transaction) (*abcitypes.ResponseCheckTx, error) {

	tx, err := transactions.ParseTransaction(req.Tx)
	if err != nil {
		return &abcitypes.ResponseCheckTx{Code: types.ParsingTransactionError}, nil
	}

	err = tx.IsValid()
	if err != nil {
		return &abcitypes.ResponseCheckTx{Code: types.InvalidTransactionError}, nil
	}

	switch tx.Type {
	case transactions.OChainPortalInteraction:
		transaction, err := transactions.ParseNewOChainPortalInteraction(tx)
		if err != nil {
			return &abcitypes.ResponseCheckTx{Code: types.ParsingTransactionDataError}, nil
		}

		err = transaction.Check(ctx)
		if err != nil {
			log.Println(err)
			return &abcitypes.ResponseCheckTx{Code: types.CheckTransactionFailure}, nil
		}

		return &abcitypes.ResponseCheckTx{Code: types.NoError}, nil
	}

	return &abcitypes.ResponseCheckTx{Code: types.NotImplemented, GasWanted: 0}, nil
}
