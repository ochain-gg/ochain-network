package engine

import (
	"encoding/hex"

	abcitypes "github.com/cometbft/cometbft/abci/types"
	"github.com/ochain-gg/ochain-network/engine/transactions"
	"github.com/ochain-gg/ochain-network/types"
)

func CheckTx(ctx transactions.TransactionContext, req *abcitypes.CheckTxRequest) *abcitypes.CheckTxResponse {
	txBytes, err := hex.DecodeString(string(req.Tx))
	if err != nil {
		return &abcitypes.CheckTxResponse{Code: types.ParsingTransactionError}
	}

	tx, err := transactions.ParseTransaction(txBytes)
	if err != nil {
		return &abcitypes.CheckTxResponse{Code: types.ParsingTransactionError}
	}

	isSystemTx := uint64(tx.Type) <= 6
	if isSystemTx {
		return CheckSystemTx(ctx, req, tx)
	} else {
		return CheckAuthenticatedTx(ctx, req, tx)
	}
}
