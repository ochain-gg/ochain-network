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
		return &abcitypes.ExecTxResult{Code: types.ParsingTransactionError, GasWanted: 0, GasUsed: 0}, []abcitypes.ValidatorUpdate{}
	}

	isSystemTx := uint64(transaction.Type) <= 5
	if isSystemTx {
		return FinalizeSystemTx(ctx, transaction)
	} else {
		return FinalizeAuthenticatedTx(ctx, transaction)
	}
}
