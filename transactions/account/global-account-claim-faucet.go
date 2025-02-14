package account_transactions

import (
	"fmt"
	"log"

	abcitypes "github.com/cometbft/cometbft/abci/types"
	"github.com/fxamacker/cbor/v2"

	t "github.com/ochain-gg/ochain-network/transactions"
	"github.com/ochain-gg/ochain-network/types"
)

type ClaimFaucetTransactionData struct {
}

type ClaimFaucetTransaction struct {
	Type      t.TransactionType          `cbor:"1,keyasint"`
	From      string                     `cbor:"2,keyasint"`
	Nonce     uint64                     `cbor:"3,keyasint"`
	Data      ClaimFaucetTransactionData `cbor:"4,keyasint"`
	Signature []byte                     `cbor:"5,keyasint"`
}

func (tx *ClaimFaucetTransaction) Transaction() (t.Transaction, error) {
	txData, err := cbor.Marshal(tx.Data)
	if err != nil {
		return t.Transaction{}, err
	}

	return t.Transaction{
		Type:      tx.Type,
		From:      tx.From,
		Nonce:     tx.Nonce,
		Data:      txData,
		Signature: tx.Signature,
	}, nil
}

func (tx *ClaimFaucetTransaction) Check(ctx t.TransactionContext) *abcitypes.CheckTxResponse {

	log.Println("[ClaimFaucetTransaction] Check")
	// var faucetAmount uint64 = 10_000_000

	globalAccount, err := ctx.Db.GlobalsAccounts.Get(tx.From)
	if err != nil {
		return &abcitypes.CheckTxResponse{
			Code: types.InvalidTransactionError,
		}
	}

	// if ctx.State.AvailableTokensInTreasury < ctx.State.OCTInGameCirculatingSupply+faucetAmount {
	// 	return &abcitypes.CheckTxResponse{
	// 		Code: types.InvalidTransactionError,
	// 	}
	// }

	if globalAccount.LastDailyClaim > ctx.Date.Unix() /*-3600*24*/ {
		return &abcitypes.CheckTxResponse{
			Code: types.InvalidTransactionError,
		}
	}

	return &abcitypes.CheckTxResponse{
		Code: types.NoError,
	}
}

func (tx *ClaimFaucetTransaction) Execute(ctx t.TransactionContext) *abcitypes.ExecTxResult {
	response := tx.Check(ctx)
	if response.Code != types.NoError {
		return &abcitypes.ExecTxResult{
			Code: response.Code,
		}
	}

	var faucetAmount uint64 = 10_000_000

	globalAccount, err := ctx.Db.GlobalsAccounts.Get(tx.From)
	if err != nil {
		return &abcitypes.ExecTxResult{
			Code: types.InvalidTransactionError,
		}
	}

	globalAccount.TokenBalance += faucetAmount
	err = ctx.Db.GlobalsAccounts.Update(globalAccount)
	if err != nil {
		return &abcitypes.ExecTxResult{
			Code: types.InvalidTransactionError,
		}
	}

	ctx.State.AddInGameCirculatingSupply(faucetAmount)

	events := []abcitypes.Event{
		{
			Type: "FaucetClaimed",
			Attributes: []abcitypes.EventAttribute{
				{Key: "address", Value: tx.From, Index: true},
				{Key: "amount", Value: fmt.Sprint(faucetAmount)},
			},
		},
	}

	return &abcitypes.ExecTxResult{
		Code:      types.NoError,
		Events:    events,
		GasUsed:   0,
		GasWanted: 0,
		Data:      []byte{},
	}
}

func ParseClaimFaucetTransaction(tx t.Transaction) (ClaimFaucetTransaction, error) {
	return ClaimFaucetTransaction{
		Type:      tx.Type,
		From:      tx.From,
		Nonce:     tx.Nonce,
		Data:      ClaimFaucetTransactionData{},
		Signature: tx.Signature,
	}, nil
}
