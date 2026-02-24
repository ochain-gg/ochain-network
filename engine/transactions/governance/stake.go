package governance_transactions

import (
	"fmt"

	abcitypes "github.com/cometbft/cometbft/abci/types"
	"github.com/fxamacker/cbor/v2"

	t "github.com/ochain-gg/ochain-network/engine/transactions"
	"github.com/ochain-gg/ochain-network/types"
)

type StakeTransactionData struct {
	Amount uint64 `cbor:"1,keyasint"`
}

type StakeTransaction struct {
	Type      t.TransactionType    `cbor:"1,keyasint"`
	From      string               `cbor:"2,keyasint"`
	Nonce     uint64               `cbor:"3,keyasint"`
	Data      StakeTransactionData `cbor:"4,keyasint"`
	Signature []byte               `cbor:"5,keyasint"`
}

func (tx *StakeTransaction) Transaction() (t.Transaction, error) {
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

func (tx *StakeTransaction) Check(ctx t.TransactionContext) *abcitypes.CheckTxResponse {
	// Check if account exists
	account, err := ctx.Db.GlobalsAccounts.Get(tx.From)
	if err != nil {
		return &abcitypes.CheckTxResponse{
			Code: types.GlobalAccountNotFoundError,
		}
	}

	// Check if amount is valid
	if tx.Data.Amount == 0 {
		return &abcitypes.CheckTxResponse{
			Code: types.InvalidAmountError,
		}
	}

	// Check if account has sufficient token balance
	if account.TokenBalance < tx.Data.Amount {
		return &abcitypes.CheckTxResponse{
			Code: types.InsufficientBalanceError,
		}
	}

	return &abcitypes.CheckTxResponse{
		Code: types.NoError,
	}
}

func (tx *StakeTransaction) Execute(ctx t.TransactionContext) *abcitypes.ExecTxResult {
	response := tx.Check(ctx)
	if response.Code != types.NoError {
		return &abcitypes.ExecTxResult{
			Code: response.Code,
		}
	}

	// Get account
	account, err := ctx.Db.GlobalsAccounts.Get(tx.From)
	if err != nil {
		return &abcitypes.ExecTxResult{
			Code: types.GlobalAccountNotFoundError,
		}
	}

	// Calculate gas cost
	gasCost, err := account.ApplyGasCost(uint64(ctx.Date.Unix()))
	if err != nil {
		return &abcitypes.ExecTxResult{
			Code: types.GasCostHigherThanBalance,
		}
	}

	// Update account balances
	account.TokenBalance -= tx.Data.Amount
	account.StackedBalance += tx.Data.Amount

	// Update account in database
	err = ctx.Db.GlobalsAccounts.Update(account)
	if err != nil {
		return &abcitypes.ExecTxResult{
			Code: types.DatabaseOperationError,
		}
	}

	events := []abcitypes.Event{
		{
			Type: "TokensStaked",
			Attributes: []abcitypes.EventAttribute{
				{Key: "address", Value: tx.From, Index: true},
				{Key: "amount", Value: fmt.Sprint(tx.Data.Amount)},
				{Key: "newStakedBalance", Value: fmt.Sprint(account.StackedBalance)},
			},
		},
	}

	receipt := t.TransactionReceipt{
		GasCost: gasCost,
	}

	receiptBytes, err := cbor.Marshal(receipt)
	if err != nil {
		return &abcitypes.ExecTxResult{
			Code: types.SerializationError,
		}
	}

	return &abcitypes.ExecTxResult{
		Code:      types.NoError,
		Events:    events,
		GasUsed:   int64(gasCost),
		GasWanted: int64(gasCost),
		Data:      receiptBytes,
	}
}

func ParseStakeTransaction(tx t.Transaction) (StakeTransaction, error) {
	var txData StakeTransactionData
	err := cbor.Unmarshal(tx.Data, &txData)

	if err != nil {
		return StakeTransaction{}, err
	}

	return StakeTransaction{
		Type:      tx.Type,
		From:      tx.From,
		Nonce:     tx.Nonce,
		Data:      txData,
		Signature: tx.Signature,
	}, nil
}
