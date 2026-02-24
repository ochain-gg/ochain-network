package governance_transactions

import (
	"fmt"

	abcitypes "github.com/cometbft/cometbft/abci/types"
	"github.com/fxamacker/cbor/v2"

	t "github.com/ochain-gg/ochain-network/engine/transactions"
	"github.com/ochain-gg/ochain-network/types"
)

type VoteOnProposalTransactionData struct {
	ProposalId uint64          `cbor:"1,keyasint"`
	Vote       types.VoteType `cbor:"2,keyasint"`
}

type VoteOnProposalTransaction struct {
	Type      t.TransactionType              `cbor:"1,keyasint"`
	From      string                         `cbor:"2,keyasint"`
	Nonce     uint64                         `cbor:"3,keyasint"`
	Data      VoteOnProposalTransactionData  `cbor:"4,keyasint"`
	Signature []byte                         `cbor:"5,keyasint"`
}

func (tx *VoteOnProposalTransaction) GetType() t.TransactionType {
	return t.VoteOnProposal
}

func (tx *VoteOnProposalTransaction) Transaction() (t.Transaction, error) {
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

func (tx *VoteOnProposalTransaction) Check(ctx t.TransactionContext) *abcitypes.CheckTxResponse {
	// Check if account exists
	account, err := ctx.Db.GlobalsAccounts.Get(tx.From)
	if err != nil {
		return &abcitypes.CheckTxResponse{
			Code: types.GlobalAccountNotFoundError,
		}
	}

	// Check if proposal exists
	proposal, err := ctx.Db.GovernanceProposals.Get(fmt.Sprint(tx.Data.ProposalId))
	if err != nil {
		return &abcitypes.CheckTxResponse{
			Code: types.ProposalNotFoundError,
		}
	}

	// Check if voting period has ended
	if uint64(ctx.Date.Unix()) >= proposal.EndingAt {
		return &abcitypes.CheckTxResponse{
			Code: types.VotingPeriodExpiredError,
		}
	}

	// Check if user has already voted
	voteId := fmt.Sprintf("%d_%s", tx.Data.ProposalId, tx.From)
	exists, err := ctx.Db.GovernanceVotes.ExistsAt(voteId, ctx.Txn.ReadTs())
	if err != nil {
		return &abcitypes.CheckTxResponse{
			Code: types.DatabaseOperationError,
		}
	}
	if exists {
		return &abcitypes.CheckTxResponse{
			Code: types.AlreadyVotedError,
		}
	}

	// User must have staked tokens to vote
	if account.StackedBalance == 0 {
		return &abcitypes.CheckTxResponse{
			Code: types.InsufficientVotingPowerError,
		}
	}

	return &abcitypes.CheckTxResponse{
		Code: types.NoError,
	}
}

func (tx *VoteOnProposalTransaction) Execute(ctx t.TransactionContext) *abcitypes.ExecTxResult {
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

	// Update account in database
	err = ctx.Db.GlobalsAccounts.Update(account)
	if err != nil {
		return &abcitypes.ExecTxResult{
			Code: types.DatabaseOperationError,
		}
	}

	// Get proposal
	proposal, err := ctx.Db.GovernanceProposals.Get(fmt.Sprint(tx.Data.ProposalId))
	if err != nil {
		return &abcitypes.ExecTxResult{
			Code: types.ProposalNotFoundError,
		}
	}

	// Voting power is equal to staked balance (1 OCT = 1 vote)
	votingPower := account.StackedBalance

	// Update proposal vote counts
	proposal.TotalVote += votingPower
	switch tx.Data.Vote {
	case types.ForProposalVote:
		proposal.TotalForVote += votingPower
	case types.AbstainProposalVote:
		proposal.TotalAbstainVote += votingPower
	case types.AgainstProposalVote:
		proposal.TotalAgainstVote += votingPower
	}

	err = ctx.Db.GovernanceProposals.Update(proposal)
	if err != nil {
		return &abcitypes.ExecTxResult{
			Code: types.DatabaseOperationError,
		}
	}

	// Record the vote
	voteId := fmt.Sprintf("%d_%s", tx.Data.ProposalId, tx.From)
	vote := types.OChainGovernanceProposalVote{
		Id:         voteId,
		ProposalId: tx.Data.ProposalId,
		Voter:      tx.From,
		Vote:       tx.Data.Vote,
		CreatedAt:  uint64(ctx.Date.Unix()),
		UpdatedAt:  uint64(ctx.Date.Unix()),
	}

	err = ctx.Db.GovernanceVotes.Insert(vote)
	if err != nil {
		return &abcitypes.ExecTxResult{
			Code: types.DatabaseOperationError,
		}
	}

	events := []abcitypes.Event{
		{
			Type: "ProposalVoteCast",
			Attributes: []abcitypes.EventAttribute{
				{Key: "proposalId", Value: fmt.Sprint(tx.Data.ProposalId), Index: true},
				{Key: "voter", Value: tx.From, Index: true},
				{Key: "vote", Value: fmt.Sprint(tx.Data.Vote)},
				{Key: "votingPower", Value: fmt.Sprint(votingPower)},
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

func ParseVoteOnProposalTransaction(tx t.Transaction) (*VoteOnProposalTransaction, error) {
	var txData VoteOnProposalTransactionData
	err := cbor.Unmarshal(tx.Data, &txData)

	if err != nil {
		return nil, err
	}

	return &VoteOnProposalTransaction{
		Type:      tx.Type,
		From:      tx.From,
		Nonce:     tx.Nonce,
		Data:      txData,
		Signature: tx.Signature,
	}, nil
}

// ParseVoteOnProposalTransactionHandler is a wrapper that returns ITransactionHandler for registry
func ParseVoteOnProposalTransactionHandler(tx t.Transaction) (t.ITransactionHandler, error) {
	return ParseVoteOnProposalTransaction(tx)
}
