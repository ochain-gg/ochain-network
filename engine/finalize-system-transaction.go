package engine

import (
	"encoding/hex"

	abcitypes "github.com/cometbft/cometbft/abci/types"
	"github.com/ochain-gg/ochain-network/transactions"
	"github.com/ochain-gg/ochain-network/types"
)

func FinalizeSystemTx(ctx transactions.TransactionContext, tx transactions.Transaction) (*abcitypes.ExecTxResult, []abcitypes.ValidatorUpdate) {
	var valUpdates []abcitypes.ValidatorUpdate

	switch tx.Type {
	case transactions.NewValidator:

		transaction, err := transactions.ParseOChainBridgeNewValidatorTransaction(tx)
		if err != nil {
			return &abcitypes.ExecTxResult{Code: types.ParsingTransactionDataError, GasWanted: 0, GasUsed: 0}, valUpdates
		}

		checkResult := transaction.Check(ctx)
		if checkResult.Code != types.NoError {
			return &abcitypes.ExecTxResult{Code: checkResult.Code, GasWanted: 0, GasUsed: 0}, valUpdates
		}

		executeResult := transaction.Execute(ctx)
		if executeResult.Code != types.NoError {
			return executeResult, valUpdates
		}

		pubkeyBytes, err := hex.DecodeString(transaction.Data.PublicKey)
		if err != nil {
			return &abcitypes.ExecTxResult{Code: types.ParsingTransactionDataError, GasWanted: 0, GasUsed: 0}, valUpdates
		}

		valUpdates = append(valUpdates, abcitypes.UpdateValidator(pubkeyBytes, 10000, "ed25519"))

	case transactions.RemoveValidator:

		transaction, err := transactions.ParseOChainBridgeRemoveValidatorTransaction(tx)
		if err != nil {
			return &abcitypes.ExecTxResult{Code: types.ParsingTransactionDataError, GasWanted: 0, GasUsed: 0}, valUpdates
		}

		checkResult := transaction.Check(ctx)
		if checkResult.Code != types.NoError {
			return &abcitypes.ExecTxResult{Code: checkResult.Code, GasWanted: 0, GasUsed: 0}, valUpdates
		}

		executeResult := transaction.Execute(ctx)
		if executeResult.Code != types.NoError {
			return executeResult, valUpdates
		}

		validator, err := ctx.Db.Validators.GetById(transaction.Data.ValidatorId)
		if err != nil {
			return &abcitypes.ExecTxResult{Code: types.ExecuteTransactionFailure, GasWanted: 0, GasUsed: 0}, valUpdates
		}

		pubkeyBytes, err := hex.DecodeString(validator.PublicKey)
		if err != nil {
			return &abcitypes.ExecTxResult{Code: types.ParsingTransactionDataError, GasWanted: 0, GasUsed: 0}, valUpdates
		}

		valUpdates = append(valUpdates, abcitypes.UpdateValidator(pubkeyBytes, 0, "ed25519"))
	}

	return &abcitypes.ExecTxResult{Code: types.NoError, GasWanted: 0, GasUsed: 0}, valUpdates
}
