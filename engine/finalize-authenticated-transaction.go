package engine

import (
	abcitypes "github.com/cometbft/cometbft/abci/types"
	"github.com/ochain-gg/ochain-network/engine/transactions"
	account_transactions "github.com/ochain-gg/ochain-network/engine/transactions/account"
	planet_transactions "github.com/ochain-gg/ochain-network/engine/transactions/game/planet"
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

		account.Nonce += 1

		if tx.Type != transactions.ClaimFaucet {
			_, err = account.ApplyGasCost(uint64(ctx.Date.Unix()))
			if err != nil {
				return &abcitypes.ExecTxResult{
					Code: types.GasCostHigherThanBalance,
				}, valUpdates
			}
		}

		err = ctx.Db.GlobalsAccounts.Update(account)
		if err != nil {
			return &abcitypes.ExecTxResult{
				Code: types.InvalidTransactionError,
			}, valUpdates
		}

	} else {
		if signer != tx.From || tx.Nonce != 0 {
			return &abcitypes.ExecTxResult{Code: types.InvalidTransactionSignature}, nil
		}
	}

	switch tx.Type {
	case transactions.RegisterAccount:

		transaction, err := account_transactions.ParseRegisterAccountTransaction(tx)
		if err != nil {
			return &abcitypes.ExecTxResult{Code: types.ParsingTransactionDataError}, valUpdates
		}

		checkResult := transaction.Check(ctx)
		if checkResult.Code != types.NoError {
			return &abcitypes.ExecTxResult{Code: checkResult.Code, GasWanted: 0, GasUsed: 0}, valUpdates
		}

		return transaction.Execute(ctx), valUpdates

	case transactions.RegisterUniverseAccount:

		transaction, err := account_transactions.ParseRegisterUniverseAccountTransaction(tx)
		if err != nil {
			return &abcitypes.ExecTxResult{Code: types.ParsingTransactionDataError}, valUpdates
		}

		checkResult := transaction.Check(ctx)
		if checkResult.Code != types.NoError {
			return &abcitypes.ExecTxResult{Code: checkResult.Code, GasWanted: 0, GasUsed: 0}, valUpdates
		}

		return transaction.Execute(ctx), valUpdates

	case transactions.ClaimFaucet:

		transaction, err := account_transactions.ParseClaimFaucetTransaction(tx)
		if err != nil {
			return &abcitypes.ExecTxResult{Code: types.ParsingTransactionDataError}, valUpdates
		}

		checkResult := transaction.Check(ctx)
		if checkResult.Code != types.NoError {
			return &abcitypes.ExecTxResult{Code: checkResult.Code, GasWanted: 0, GasUsed: 0}, valUpdates
		}

		return transaction.Execute(ctx), valUpdates

	case transactions.StartBuildingUpgrade:

		transaction, err := planet_transactions.ParseUpgradeBuildingTransaction(tx)
		if err != nil {
			return &abcitypes.ExecTxResult{Code: types.ParsingTransactionDataError}, valUpdates
		}

		checkResult := transaction.Check(ctx)
		if checkResult.Code != types.NoError {
			return &abcitypes.ExecTxResult{Code: checkResult.Code, GasWanted: 0, GasUsed: 0}, valUpdates
		}

		return transaction.Execute(ctx), valUpdates

	case transactions.StartTechnologyUpgrade:

		transaction, err := planet_transactions.ParseUpgradeTechnologyTransaction(tx)
		if err != nil {
			return &abcitypes.ExecTxResult{Code: types.ParsingTransactionDataError}, valUpdates
		}

		checkResult := transaction.Check(ctx)
		if checkResult.Code != types.NoError {
			return &abcitypes.ExecTxResult{Code: checkResult.Code, GasWanted: 0, GasUsed: 0}, valUpdates
		}

		return transaction.Execute(ctx), valUpdates

	case transactions.Build:

		transaction, err := planet_transactions.ParseBuildTransaction(tx)
		if err != nil {
			return &abcitypes.ExecTxResult{Code: types.ParsingTransactionDataError}, valUpdates
		}

		checkResult := transaction.Check(ctx)
		if checkResult.Code != types.NoError {
			return &abcitypes.ExecTxResult{Code: checkResult.Code, GasWanted: 0, GasUsed: 0}, valUpdates
		}

		return transaction.Execute(ctx), valUpdates

	}

	return &abcitypes.ExecTxResult{Code: types.NotImplemented, GasWanted: 0}, valUpdates
}
