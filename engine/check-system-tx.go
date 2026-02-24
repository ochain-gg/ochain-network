package engine

import (
	abcitypes "github.com/cometbft/cometbft/abci/types"
	"github.com/ochain-gg/ochain-network/engine/transactions"
	"github.com/ochain-gg/ochain-network/types"
)

func CheckSystemTx(ctx transactions.TransactionContext, req *abcitypes.CheckTxRequest, tx transactions.Transaction) *abcitypes.CheckTxResponse {
	// Use registry to parse and check system transaction
	if !transactions.IsRegistered(tx.Type) {
		return &abcitypes.CheckTxResponse{Code: types.NotImplemented, GasWanted: 0, GasUsed: 0}
	}

	handler, err := transactions.Parse(tx)
	if err != nil {
		return &abcitypes.CheckTxResponse{Code: types.ParsingTransactionDataError, GasWanted: 0, GasUsed: 0}
	}

	return handler.Check(ctx)
}
