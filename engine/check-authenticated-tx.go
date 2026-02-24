package engine

import (
	"log"

	abcitypes "github.com/cometbft/cometbft/abci/types"
	"github.com/ochain-gg/ochain-network/engine/transactions"
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

	// Use registry to parse and check transaction
	if !transactions.IsRegistered(tx.Type) {
		log.Printf("[CheckAuthenticatedTx] Transaction type %d is not registered", tx.Type)
		return &abcitypes.CheckTxResponse{Code: types.NotImplemented, GasWanted: 0, GasUsed: 0}
	}

	handler, err := transactions.Parse(tx)
	if err != nil {
		log.Printf("[CheckAuthenticatedTx] Failed to parse transaction type %d: %v", tx.Type, err)
		return &abcitypes.CheckTxResponse{Code: types.ParsingTransactionDataError, GasWanted: 0, GasUsed: 0}
	}

	log.Printf("[CheckAuthenticatedTx] Successfully parsed transaction type %d", tx.Type)
	return handler.Check(ctx)
}
