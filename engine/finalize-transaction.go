package engine

import (
	"encoding/hex"
	"log"

	abcitypes "github.com/cometbft/cometbft/abci/types"
	"github.com/ochain-gg/ochain-network/engine/transactions"
	"github.com/ochain-gg/ochain-network/types"
)

func FinalizeTx(ctx transactions.TransactionContext, tx []byte) (*abcitypes.ExecTxResult, []abcitypes.ValidatorUpdate) {

	log.Println("FINALIZE TX DATA RECEIVED: ", string(tx))

	txBytes, err := hex.DecodeString(string(tx))
	transac, err := transactions.ParseTransaction(txBytes)
	if err != nil {
		log.Println("finalize tx failed: " + err.Error())
		return &abcitypes.ExecTxResult{Code: types.ParsingTransactionError}, nil
	}

	log.Println("finalize tx decoded: ", transac)

	isSystemTx := uint64(transac.Type) <= 6
	if isSystemTx {
		return FinalizeSystemTx(ctx, transac)
	} else {
		return FinalizeAuthenticatedTx(ctx, transac)
	}
}
