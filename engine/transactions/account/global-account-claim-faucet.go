package account_transactions

import (
	"fmt"
	"log"

	abcitypes "github.com/cometbft/cometbft/abci/types"
	"github.com/fxamacker/cbor/v2"

	t "github.com/ochain-gg/ochain-network/engine/transactions"
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

func (tx *ClaimFaucetTransaction) GetType() t.TransactionType {
	return t.ClaimFaucet
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

const (
	secondsPerDay int64 = 86400 // 24 * 60 * 60
)

// getDailyLoginReward returns the reward amount based on the streak day (1-7)
// Day 1-2: 1 OCT, Day 3-4: 2 OCT, Day 5-6: 5 OCT, Day 7: 10 OCT
func getDailyLoginReward(streakDay uint64) uint64 {
	switch streakDay {
	case 1, 2:
		return 1 * types.OneOCT
	case 3, 4:
		return 2 * types.OneOCT
	case 5, 6:
		return 5 * types.OneOCT
	case 7:
		return 10 * types.OneOCT
	default:
		return 1 * types.OneOCT
	}
}

func (tx *ClaimFaucetTransaction) Check(ctx t.TransactionContext) *abcitypes.CheckTxResponse {

	log.Println("[ClaimFaucetTransaction] Check")

	globalAccount, err := ctx.Db.GlobalsAccounts.Get(tx.From)
	if err != nil {
		return &abcitypes.CheckTxResponse{
			Code: types.InvalidTransactionError,
		}
	}

	// Check if already claimed today
	currentDay := ctx.Date.Unix() / secondsPerDay
	lastClaimDay := globalAccount.LastDailyClaim / secondsPerDay

	if currentDay <= lastClaimDay {
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

	globalAccount, err := ctx.Db.GlobalsAccounts.Get(tx.From)
	if err != nil {
		return &abcitypes.ExecTxResult{
			Code: types.InvalidTransactionError,
		}
	}

	// Calculate streak
	currentDay := ctx.Date.Unix() / secondsPerDay
	lastClaimDay := globalAccount.LastDailyClaim / secondsPerDay

	var newStreak uint64
	if currentDay == lastClaimDay+1 {
		// Consecutive day - continue streak
		newStreak = globalAccount.DailyLoginStreak + 1
		if newStreak > 7 {
			newStreak = 1 // Reset to day 1 after completing 7-day cycle
		}
	} else {
		// Missed a day or first claim - reset streak to day 1
		newStreak = 1
	}

	// Calculate reward based on streak day
	rewardAmount := getDailyLoginReward(newStreak)

	// Update account
	globalAccount.TokenBalance += rewardAmount
	globalAccount.LastDailyClaim = ctx.Date.Unix()
	globalAccount.DailyLoginStreak = newStreak

	err = ctx.Db.GlobalsAccounts.Update(globalAccount)
	if err != nil {
		return &abcitypes.ExecTxResult{
			Code: types.InvalidTransactionError,
		}
	}

	ctx.State.AddInGameCirculatingSupply(rewardAmount)

	events := []abcitypes.Event{
		{
			Type: "DailyLoginRewardClaimed",
			Attributes: []abcitypes.EventAttribute{
				{Key: "address", Value: tx.From, Index: true},
				{Key: "amount", Value: fmt.Sprint(rewardAmount)},
				{Key: "streakDay", Value: fmt.Sprint(newStreak)},
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

func ParseClaimFaucetTransaction(tx t.Transaction) (*ClaimFaucetTransaction, error) {
	return &ClaimFaucetTransaction{
		Type:      tx.Type,
		From:      tx.From,
		Nonce:     tx.Nonce,
		Data:      ClaimFaucetTransactionData{},
		Signature: tx.Signature,
	}, nil
}

// ParseClaimFaucetTransactionHandler is a wrapper that returns ITransactionHandler for registry
func ParseClaimFaucetTransactionHandler(tx t.Transaction) (t.ITransactionHandler, error) {
	return ParseClaimFaucetTransaction(tx)
}
