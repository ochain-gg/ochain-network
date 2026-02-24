package engine

import (
	"encoding/hex"

	abcitypes "github.com/cometbft/cometbft/abci/types"
	"github.com/ochain-gg/ochain-network/engine/transactions"
	validator_transactions "github.com/ochain-gg/ochain-network/engine/transactions/validator"
	"github.com/ochain-gg/ochain-network/types"
)

func FinalizeSystemTx(ctx transactions.TransactionContext, tx transactions.Transaction) (*abcitypes.ExecTxResult, []abcitypes.ValidatorUpdate) {
	var valUpdates []abcitypes.ValidatorUpdate

	// Use registry to parse, check, and execute system transaction
	if !transactions.IsRegistered(tx.Type) {
		return &abcitypes.ExecTxResult{Code: types.NotImplemented, GasWanted: 0, GasUsed: 0}, valUpdates
	}

	handler, err := transactions.Parse(tx)
	if err != nil {
		return &abcitypes.ExecTxResult{Code: types.ParsingTransactionDataError, GasWanted: 0, GasUsed: 0}, valUpdates
	}

	checkResult := handler.Check(ctx)
	if checkResult.Code != types.NoError {
		return &abcitypes.ExecTxResult{Code: checkResult.Code, GasWanted: 0, GasUsed: 0}, valUpdates
	}

	executeResult := handler.Execute(ctx)
	if executeResult.Code != types.NoError {
		return executeResult, valUpdates
	}

	// Special handling for validator updates
	switch tx.Type {
	case transactions.NewValidator:
		transaction := handler.(*validator_transactions.OChainNewValidatorTransaction)
		pubkeyBytes, err := hex.DecodeString(transaction.Data.PublicKey)
		if err != nil {
			return &abcitypes.ExecTxResult{Code: types.ParsingTransactionDataError, GasWanted: 0, GasUsed: 0}, valUpdates
		}
		valUpdates = append(valUpdates, abcitypes.ValidatorUpdate{
			PubKeyType:  "tendermint/PubKeyEd25519",
			PubKeyBytes: pubkeyBytes,
			Power:       10000,
		})

	case transactions.RemoveValidator:
		transaction := handler.(*validator_transactions.OChainRemoveValidatorTransaction)
		validator, err := ctx.Db.Validators.GetById(transaction.Data.ValidatorId)
		if err != nil {
			return &abcitypes.ExecTxResult{Code: types.ExecuteTransactionFailure, GasWanted: 0, GasUsed: 0}, valUpdates
		}
		pubkeyBytes, err := hex.DecodeString(validator.PublicKey)
		if err != nil {
			return &abcitypes.ExecTxResult{Code: types.ParsingTransactionDataError, GasWanted: 0, GasUsed: 0}, valUpdates
		}
		valUpdates = append(valUpdates, abcitypes.ValidatorUpdate{
			PubKeyType:  "tendermint/PubKeyEd25519",
			PubKeyBytes: pubkeyBytes,
			Power:       0,
		})
	}

	return executeResult, valUpdates
}
