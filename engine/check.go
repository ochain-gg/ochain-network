package engine

import (
	"encoding/hex"
	"log"

	abcitypes "github.com/cometbft/cometbft/abci/types"
	"github.com/ochain-gg/ochain-network/transactions"
	"github.com/ochain-gg/ochain-network/types"
)

func CheckTx(ctx transactions.TransactionContext, req *abcitypes.RequestCheckTx) (*abcitypes.ResponseCheckTx, error) {

	log.Printf("Check tx: %s", hex.EncodeToString(req.Tx))

	tx, err := transactions.ParseTransaction(req.Tx)
	if err != nil {
		return &abcitypes.ResponseCheckTx{Code: types.ParsingTransactionError}, nil
	}

	err = tx.IsValid()
	if err != nil {
		return &abcitypes.ResponseCheckTx{Code: types.InvalidTransactionError}, nil
	}

	isSystemTx := tx.Type == transactions.OChainPortalInteraction || tx.Type == transactions.ExecutePendingUpgrade
	if isSystemTx {
		return CheckSystemTx(ctx, req, tx)
	} else {
		return CheckAuthenticatedTx(ctx, req, tx)
	}
}
