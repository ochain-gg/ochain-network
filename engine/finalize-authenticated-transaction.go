package engine

import (
	abcitypes "github.com/cometbft/cometbft/abci/types"
	"github.com/ochain-gg/ochain-network/transactions"
	"github.com/ochain-gg/ochain-network/types"
)

func FinalizeAuthenticatedTx(ctx transactions.TransactionContext, tx transactions.Transaction) (*abcitypes.ExecTxResult, []abcitypes.ValidatorUpdate) {
	var valUpdates []abcitypes.ValidatorUpdate

	signer, err := tx.GetSigner()
	if err != nil {
		return &abcitypes.ExecTxResult{Code: types.InvalidTransactionSignature}, valUpdates
	}

	if tx.Type != transactions.RegisterAccount {
		account, err := ctx.Db.GlobalsAccounts.Get(tx.From)
		if err != nil {
			return &abcitypes.ExecTxResult{Code: types.InvalidTransactionSignature}, nil
		}

		if !account.IsAllowedSigner(signer, IsDeleguatedAuthorized(tx.Type)) {
			return &abcitypes.ExecTxResult{Code: types.InvalidTransactionSignature}, nil
		}

		//verify nonce
		if tx.Nonce != account.Nonce {
			return &abcitypes.ExecTxResult{Code: types.InvalidTransactionSignature}, nil
		}
	} else {
		if signer != tx.From || tx.Nonce != 0 {
			return &abcitypes.ExecTxResult{Code: types.InvalidTransactionSignature}, nil
		}
	}

	switch tx.Type {
	case transactions.RegisterAccount:

		transaction, err := transactions.ParseRegisterAccountTransaction(tx)
		if err != nil {
			return &abcitypes.ExecTxResult{Code: types.ParsingTransactionDataError}, valUpdates
		}

		err = transaction.Check(ctx)
		if err != nil {
			return &abcitypes.ExecTxResult{Code: types.CheckTransactionFailure}, valUpdates
		}

		events, err := transaction.Execute(ctx)
		if err != nil {
			return &abcitypes.ExecTxResult{Code: types.ExecuteTransactionFailure}, valUpdates
		}

		return &abcitypes.ExecTxResult{
			Code:    types.NoError,
			Events:  events,
			GasUsed: 0,
		}, valUpdates

	case transactions.RegisterUniverseAccount:

		transaction, err := transactions.ParseRegisterUniverseAccountTransaction(tx)
		if err != nil {
			return &abcitypes.ExecTxResult{Code: types.ParsingTransactionDataError}, valUpdates
		}

		err = transaction.Check(ctx)
		if err != nil {
			return &abcitypes.ExecTxResult{Code: types.CheckTransactionFailure}, valUpdates
		}

		events, err := transaction.Execute(ctx)
		if err != nil {
			return &abcitypes.ExecTxResult{Code: types.ExecuteTransactionFailure}, valUpdates
		}

		return &abcitypes.ExecTxResult{
			Code:    types.NoError,
			Events:  events,
			GasUsed: 0,
		}, valUpdates
	}

	return &abcitypes.ExecTxResult{Code: types.NotImplemented, GasWanted: 0}, valUpdates
}
