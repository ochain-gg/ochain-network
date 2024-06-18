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
	case transactions.OChainPortalInteraction:

		transaction, err := transactions.ParseNewOChainPortalInteraction(tx)
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

		if transaction.Data.Type == transactions.NewValidatorPortalInteractionType {
			formatedTx, err := transactions.ParseNewValidatorTransaction(transaction)
			if err != nil {
				return &abcitypes.ExecTxResult{Code: types.ParsingTransactionDataError, GasWanted: 0, GasUsed: 0}, valUpdates
			}

			pubkeyBytes, err := hex.DecodeString(formatedTx.Data.Arguments.PublicKey)
			if err != nil {
				return &abcitypes.ExecTxResult{Code: types.ParsingTransactionDataError, GasWanted: 0, GasUsed: 0}, valUpdates
			}

			valUpdates = append(valUpdates, abcitypes.UpdateValidator(pubkeyBytes, 10000, "ed25519"))
		} else if transaction.Data.Type == transactions.RemoveValidatorPortalInteractionType {
			formatedTx, err := transactions.ParseRemoveValidatorTransaction(transaction)
			if err != nil {
				return &abcitypes.ExecTxResult{Code: types.ParsingTransactionDataError, GasWanted: 0, GasUsed: 0}, valUpdates
			}

			validator, err := ctx.Db.Validators.GetById(formatedTx.Data.Arguments.ValidatorId)
			if err != nil {
				return &abcitypes.ExecTxResult{Code: types.ExecuteTransactionFailure, GasWanted: 0, GasUsed: 0}, valUpdates
			}

			pubkeyBytes, err := hex.DecodeString(validator.PublicKey)
			if err != nil {
				return &abcitypes.ExecTxResult{Code: types.ParsingTransactionDataError, GasWanted: 0, GasUsed: 0}, valUpdates
			}

			valUpdates = append(valUpdates, abcitypes.UpdateValidator(pubkeyBytes, 0, "ed25519"))
		}
	}

	return &abcitypes.ExecTxResult{Code: types.NoError, GasWanted: 0, GasUsed: 0}, valUpdates
}
