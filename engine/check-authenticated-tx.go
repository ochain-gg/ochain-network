package engine

import (
	abcitypes "github.com/cometbft/cometbft/abci/types"
	"github.com/ochain-gg/ochain-network/transactions"
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
	case transactions.UniverseOCTDeposit:
		return false
	case transactions.UniverseOCTWithdraw:
		return false
	case transactions.AccountOCTWithdraw:
		return false
	default:
		return true
	}
}

func CheckAuthenticatedTx(ctx transactions.TransactionContext, req *abcitypes.RequestCheckTx, tx transactions.Transaction) (*abcitypes.ResponseCheckTx, error) {

	signer, err := tx.GetSigner()
	if err != nil {
		return &abcitypes.ResponseCheckTx{Code: types.InvalidTransactionSignature}, nil
	}

	if tx.Type != transactions.RegisterAccount {
		account, err := ctx.Db.GlobalsAccounts.Get(tx.From)
		if err != nil {
			return &abcitypes.ResponseCheckTx{Code: types.InvalidTransactionSignature}, nil
		}

		if !account.IsAllowedSigner(signer, IsDeleguatedAuthorized(tx.Type)) {
			return &abcitypes.ResponseCheckTx{Code: types.InvalidTransactionSignature}, nil
		}

		//verify nonce
		if tx.Nonce != account.Nonce {
			return &abcitypes.ResponseCheckTx{Code: types.InvalidTransactionSignature}, nil
		}
	} else {
		if signer != tx.From || tx.Nonce != 0 {
			return &abcitypes.ResponseCheckTx{Code: types.InvalidTransactionSignature}, nil
		}
	}

	switch tx.Type {
	case transactions.RegisterAccount:

		transaction, err := transactions.ParseRegisterAccountTransaction(tx)
		if err != nil {
			return &abcitypes.ResponseCheckTx{Code: types.ParsingTransactionDataError}, nil
		}

		err = transaction.Check(ctx)
		if err != nil {
			return &abcitypes.ResponseCheckTx{Code: types.CheckTransactionFailure}, nil
		}

		return &abcitypes.ResponseCheckTx{Code: types.NoError}, nil

	case transactions.RegisterUniverseAccount:

		transaction, err := transactions.ParseRegisterUniverseAccountTransaction(tx)
		if err != nil {
			return &abcitypes.ResponseCheckTx{Code: types.ParsingTransactionDataError}, nil
		}

		err = transaction.Check(ctx)
		if err != nil {
			return &abcitypes.ResponseCheckTx{Code: types.CheckTransactionFailure}, nil
		}

		return &abcitypes.ResponseCheckTx{Code: types.NoError}, nil
	}

	return &abcitypes.ResponseCheckTx{Code: types.NotImplemented, GasWanted: 0}, nil
}
