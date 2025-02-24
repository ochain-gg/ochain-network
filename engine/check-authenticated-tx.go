package engine

import (
	abcitypes "github.com/cometbft/cometbft/abci/types"
	"github.com/ochain-gg/ochain-network/engine/transactions"
	account_transactions "github.com/ochain-gg/ochain-network/engine/transactions/account"
	planet_transactions "github.com/ochain-gg/ochain-network/engine/transactions/game/planet"
	"github.com/ochain-gg/ochain-network/types"
)

const (
	AppVersion uint64 = 1
)

type OChainEngine struct {
}

func IsDeleguatedAuthorized(txType transactions.TransactionType) bool {
	switch txType {
	case transactions.RegisterAccount:
		return false
	case transactions.UniverseOChainTokenWithdraw:
		return false
	case transactions.OChainTokenWithdrawal:
		return false
	default:
		return true
	}
}

func CheckAuthenticatedTx(ctx transactions.TransactionContext, req *abcitypes.CheckTxRequest, tx transactions.Transaction) *abcitypes.CheckTxResponse {

	signer, err := tx.GetSigner()
	if err != nil {
		return &abcitypes.CheckTxResponse{Code: types.InvalidTransactionSignature, GasWanted: 0, GasUsed: 0}
	}

	if tx.Type != transactions.RegisterAccount {
		account, err := ctx.Db.GlobalsAccounts.Get(tx.From)
		if err != nil {
			return &abcitypes.CheckTxResponse{Code: types.InvalidTransactionSignature, GasWanted: 0, GasUsed: 0}
		}

		if !account.IsAllowedSigner(signer, IsDeleguatedAuthorized(tx.Type)) {
			return &abcitypes.CheckTxResponse{Code: types.InvalidTransactionSignature, GasWanted: 0, GasUsed: 0}
		}

		//verify nonce
		if tx.Nonce != account.Nonce {
			return &abcitypes.CheckTxResponse{Code: types.InvalidNonce, GasWanted: 0, GasUsed: 0}
		}
	} else {
		if signer != tx.From || tx.Nonce != 0 {
			return &abcitypes.CheckTxResponse{Code: types.InvalidTransactionSignature, GasWanted: 0, GasUsed: 0}
		}
	}

	switch tx.Type {
	case transactions.RegisterAccount:

		transaction, err := account_transactions.ParseRegisterAccountTransaction(tx)
		if err != nil {
			return &abcitypes.CheckTxResponse{Code: types.ParsingTransactionDataError, GasWanted: 0, GasUsed: 0}
		}

		return transaction.Check(ctx)

	case transactions.RegisterUniverseAccount:

		transaction, err := account_transactions.ParseRegisterUniverseAccountTransaction(tx)
		if err != nil {
			return &abcitypes.CheckTxResponse{Code: types.ParsingTransactionDataError, GasWanted: 0, GasUsed: 0}
		}

		return transaction.Check(ctx)

	case transactions.ClaimFaucet:

		transaction, err := account_transactions.ParseClaimFaucetTransaction(tx)
		if err != nil {
			return &abcitypes.CheckTxResponse{Code: types.ParsingTransactionDataError, GasWanted: 0, GasUsed: 0}
		}

		return transaction.Check(ctx)

	case transactions.StartBuildingUpgrade:

		transaction, err := planet_transactions.ParseUpgradeBuildingTransaction(tx)
		if err != nil {
			return &abcitypes.CheckTxResponse{Code: types.ParsingTransactionDataError, GasWanted: 0, GasUsed: 0}
		}

		return transaction.Check(ctx)

	case transactions.StartTechnologyUpgrade:

		transaction, err := planet_transactions.ParseUpgradeTechnologyTransaction(tx)
		if err != nil {
			return &abcitypes.CheckTxResponse{Code: types.ParsingTransactionDataError, GasWanted: 0, GasUsed: 0}
		}

		return transaction.Check(ctx)
	case transactions.Build:

		transaction, err := planet_transactions.ParseBuildTransaction(tx)
		if err != nil {
			return &abcitypes.CheckTxResponse{Code: types.ParsingTransactionDataError, GasWanted: 0, GasUsed: 0}
		}

		return transaction.Check(ctx)
	}

	return &abcitypes.CheckTxResponse{Code: types.NotImplemented, GasWanted: 0, GasUsed: 0}
}
