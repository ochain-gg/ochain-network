package engine

import (
	"encoding/hex"
	"log"

	abcitypes "github.com/cometbft/cometbft/abci/types"
	"github.com/ochain-gg/ochain-network/transactions"
	"github.com/ochain-gg/ochain-network/types"
)

func FinalizeTx(ctx transactions.TransactionContext, tx []byte) (*abcitypes.ExecTxResult, []abcitypes.ValidatorUpdate) {

	log.Printf("Check tx: %s", hex.EncodeToString(tx))

	transaction, err := transactions.ParseTransaction(tx)
	if err != nil {
		return &abcitypes.ExecTxResult{Code: types.ParsingTransactionError}, []abcitypes.ValidatorUpdate{}
	}

	err = transaction.IsValid()
	if err != nil {
		return &abcitypes.ExecTxResult{Code: types.InvalidTransactionError}, []abcitypes.ValidatorUpdate{}
	}

	isSystemTx := transaction.Type == transactions.OChainPortalInteraction || transaction.Type == transactions.ExecutePendingUpdate
	if isSystemTx {
		return FinalizeSystemTx(ctx, transaction)
	} else {
		return FinalizeAuthenticatedTx(ctx, transaction)
	}
}
