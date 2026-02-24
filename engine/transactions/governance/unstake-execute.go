package governance_transactions

import (
	"fmt"

	abcitypes "github.com/cometbft/cometbft/abci/types"
	"github.com/fxamacker/cbor/v2"

	t "github.com/ochain-gg/ochain-network/engine/transactions"
	"github.com/ochain-gg/ochain-network/types"
)

type UnstakeExecuteTransactionData struct {
	// No additional data needed, just executes the pending unstake
}

type UnstakeExecuteTransaction struct {
	Type      t.TransactionType              `cbor:"1,keyasint"`
	From      string                         `cbor:"2,keyasint"`
	Nonce     uint64                         `cbor:"3,keyasint"`
	Data      UnstakeExecuteTransactionData  `cbor:"4,keyasint"`
	Signature []byte                         `cbor:"5,keyasint"`
}

func (tx *UnstakeExecuteTransaction) GetType() t.TransactionType {
	return t.GovernanceUnstakeExecute
}

func (tx *UnstakeExecuteTransaction) Transaction() (t.Transaction, error) {
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

func (tx *UnstakeExecuteTransaction) Check(ctx t.TransactionContext) *abcitypes.CheckTxResponse {
	// Check if account exists
	_, err := ctx.Db.GlobalsAccounts.Get(tx.From)
	if err != nil {
		return &abcitypes.CheckTxResponse{
			Code: types.GlobalAccountNotFoundError,
		}
	}

	// Check if there's a pending unstake
	pendingUnstake, err := ctx.Db.PendingUnstakes.GetAt(tx.From, ctx.Txn.ReadTs())
	if err != nil {
		return &abcitypes.CheckTxResponse{
			Code: types.NoPendingUnstakeFound,
		}
	}

	// Check if cooldown period has passed
	if ctx.Date.Unix() < pendingUnstake.WithdrawableAt {
		return &abcitypes.CheckTxResponse{
			Code: types.UnstakeCooldownNotCompleted,
		}
	}

	return &abcitypes.CheckTxResponse{
		Code: types.NoError,
	}
}

func (tx *UnstakeExecuteTransaction) Execute(ctx t.TransactionContext) *abcitypes.ExecTxResult {
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

	// Get pending unstake
	pendingUnstake, err := ctx.Db.PendingUnstakes.Get(tx.From)
	if err != nil {
		return &abcitypes.ExecTxResult{
			Code: types.NoPendingUnstakeFound,
		}
	}

	// Calculate gas cost
	gasCost, err := account.ApplyGasCost(uint64(ctx.Date.Unix()))
	if err != nil {
		return &abcitypes.ExecTxResult{
			Code: types.GasCostHigherThanBalance,
		}
	}

	// Return tokens to account balance
	account.TokenBalance += pendingUnstake.Amount

	// Update account in database
	err = ctx.Db.GlobalsAccounts.Update(account)
	if err != nil {
		return &abcitypes.ExecTxResult{
			Code: types.DatabaseOperationError,
		}
	}

	// Delete pending unstake
	err = ctx.Db.PendingUnstakes.Delete(tx.From)
	if err != nil {
		return &abcitypes.ExecTxResult{
			Code: types.DatabaseOperationError,
		}
	}

	events := []abcitypes.Event{
		{
			Type: "UnstakeExecuted",
			Attributes: []abcitypes.EventAttribute{
				{Key: "address", Value: tx.From, Index: true},
				{Key: "amount", Value: fmt.Sprint(pendingUnstake.Amount)},
				{Key: "newTokenBalance", Value: fmt.Sprint(account.TokenBalance)},
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

func ParseUnstakeExecuteTransaction(tx t.Transaction) (*UnstakeExecuteTransaction, error) {
	var txData UnstakeExecuteTransactionData
	err := cbor.Unmarshal(tx.Data, &txData)

	if err != nil {
		return nil, err
	}

	return &UnstakeExecuteTransaction{
		Type:      tx.Type,
		From:      tx.From,
		Nonce:     tx.Nonce,
		Data:      txData,
		Signature: tx.Signature,
	}, nil
}

// ParseUnstakeExecuteTransactionHandler is a wrapper that returns ITransactionHandler for registry
func ParseUnstakeExecuteTransactionHandler(tx t.Transaction) (t.ITransactionHandler, error) {
	return ParseUnstakeExecuteTransaction(tx)
}
