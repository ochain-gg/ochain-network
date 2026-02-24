package governance_transactions

import (
	"fmt"

	abcitypes "github.com/cometbft/cometbft/abci/types"
	"github.com/fxamacker/cbor/v2"

	t "github.com/ochain-gg/ochain-network/engine/transactions"
	"github.com/ochain-gg/ochain-network/types"
)

const (
	UnstakeCooldownPeriod = 14 * 24 * 60 * 60 // 14 days in seconds
)

type UnstakeStartTransactionData struct {
	Amount uint64 `cbor:"1,keyasint"`
}

type UnstakeStartTransaction struct {
	Type      t.TransactionType           `cbor:"1,keyasint"`
	From      string                      `cbor:"2,keyasint"`
	Nonce     uint64                      `cbor:"3,keyasint"`
	Data      UnstakeStartTransactionData `cbor:"4,keyasint"`
	Signature []byte                      `cbor:"5,keyasint"`
}

func (tx *UnstakeStartTransaction) GetType() t.TransactionType {
	return t.GovernanceUnstakeStart
}

func (tx *UnstakeStartTransaction) Transaction() (t.Transaction, error) {
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

func (tx *UnstakeStartTransaction) Check(ctx t.TransactionContext) *abcitypes.CheckTxResponse {
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

	// Check if account has sufficient staked balance
	if account.StackedBalance < tx.Data.Amount {
		return &abcitypes.CheckTxResponse{
			Code: types.InsufficientStakedBalance,
		}
	}

	// Check if there's already a pending unstake
	exists, err := ctx.Db.PendingUnstakes.ExistsAt(tx.From, ctx.Txn.ReadTs())
	if err != nil {
		return &abcitypes.CheckTxResponse{
			Code: types.DatabaseOperationError,
		}
	}

	if exists {
		return &abcitypes.CheckTxResponse{
			Code: types.PendingUnstakeAlreadyExists,
		}
	}

	return &abcitypes.CheckTxResponse{
		Code: types.NoError,
	}
}

func (tx *UnstakeStartTransaction) Execute(ctx t.TransactionContext) *abcitypes.ExecTxResult {
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

	// Update account staked balance
	account.StackedBalance -= tx.Data.Amount

	// Update account in database
	err = ctx.Db.GlobalsAccounts.Update(account)
	if err != nil {
		return &abcitypes.ExecTxResult{
			Code: types.DatabaseOperationError,
		}
	}

	// Create pending unstake entry
	withdrawableAt := ctx.Date.Unix() + UnstakeCooldownPeriod
	pendingUnstake := types.OChainPendingUnstake{
		Address:        tx.From,
		Amount:         tx.Data.Amount,
		StartedAt:      ctx.Date.Unix(),
		WithdrawableAt: withdrawableAt,
	}

	err = ctx.Db.PendingUnstakes.Insert(pendingUnstake)
	if err != nil {
		return &abcitypes.ExecTxResult{
			Code: types.DatabaseOperationError,
		}
	}

	events := []abcitypes.Event{
		{
			Type: "UnstakeStarted",
			Attributes: []abcitypes.EventAttribute{
				{Key: "address", Value: tx.From, Index: true},
				{Key: "amount", Value: fmt.Sprint(tx.Data.Amount)},
				{Key: "withdrawableAt", Value: fmt.Sprint(withdrawableAt)},
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

func ParseUnstakeStartTransaction(tx t.Transaction) (*UnstakeStartTransaction, error) {
	var txData UnstakeStartTransactionData
	err := cbor.Unmarshal(tx.Data, &txData)

	if err != nil {
		return nil, err
	}

	return &UnstakeStartTransaction{
		Type:      tx.Type,
		From:      tx.From,
		Nonce:     tx.Nonce,
		Data:      txData,
		Signature: tx.Signature,
	}, nil
}

// ParseUnstakeStartTransactionHandler is a wrapper that returns ITransactionHandler for registry
func ParseUnstakeStartTransactionHandler(tx t.Transaction) (t.ITransactionHandler, error) {
	return ParseUnstakeStartTransaction(tx)
}
