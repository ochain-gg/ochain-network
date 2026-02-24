package governance_transactions

import (
	"fmt"

	abcitypes "github.com/cometbft/cometbft/abci/types"
	"github.com/fxamacker/cbor/v2"

	t "github.com/ochain-gg/ochain-network/engine/transactions"
	"github.com/ochain-gg/ochain-network/types"
)

type CreateProposalTransactionData struct {
	ProposalType types.GovernanceProposalType `cbor:"1,keyasint"`
	Payload      []byte                       `cbor:"2,keyasint"`
}

type CreateProposalTransaction struct {
	Type      t.TransactionType              `cbor:"1,keyasint"`
	From      string                         `cbor:"2,keyasint"`
	Nonce     uint64                         `cbor:"3,keyasint"`
	Data      CreateProposalTransactionData  `cbor:"4,keyasint"`
	Signature []byte                         `cbor:"5,keyasint"`
}

func (tx *CreateProposalTransaction) GetType() t.TransactionType {
	return t.CreateGovernanceProposal
}

func (tx *CreateProposalTransaction) Transaction() (t.Transaction, error) {
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

func (tx *CreateProposalTransaction) Check(ctx t.TransactionContext) *abcitypes.CheckTxResponse {
	// Check if account exists
	account, err := ctx.Db.GlobalsAccounts.Get(tx.From)
	if err != nil {
		return &abcitypes.CheckTxResponse{
			Code: types.GlobalAccountNotFoundError,
		}
	}

	// Check if user has enough voting power (staked balance) to create proposal
	if account.StackedBalance < types.VotingPowerRequiredForProposalCreation {
		return &abcitypes.CheckTxResponse{
			Code: types.InsufficientVotingPowerError,
		}
	}

	// Validate proposal payload based on type
	switch tx.Data.ProposalType {
	case types.CreateUniverseProposal:
		var payload types.CreateUniverseProposalPayload
		if err := cbor.Unmarshal(tx.Data.Payload, &payload); err != nil {
			return &abcitypes.CheckTxResponse{
				Code: types.InvalidProposalPayloadError,
			}
		}
		// Check if universe already exists
		exists, err := ctx.Db.Universes.ExistsAt(payload.Id, ctx.Txn.ReadTs())
		if err != nil {
			return &abcitypes.CheckTxResponse{
				Code: types.DatabaseOperationError,
			}
		}
		if exists {
			return &abcitypes.CheckTxResponse{
				Code: types.UniverseAlreadyExistsError,
			}
		}

	case types.PauseUniverseProposal:
		var payload types.PauseUniverseProposalPayload
		if err := cbor.Unmarshal(tx.Data.Payload, &payload); err != nil {
			return &abcitypes.CheckTxResponse{
				Code: types.InvalidProposalPayloadError,
			}
		}
		// Check if universe exists
		universe, err := ctx.Db.Universes.Get(payload.UniverseId)
		if err != nil {
			return &abcitypes.CheckTxResponse{
				Code: types.UniverseNotFoundError,
			}
		}
		// Check if universe is already paused
		if universe.Paused {
			return &abcitypes.CheckTxResponse{
				Code: types.UniverseAlreadyPausedError,
			}
		}

	case types.ResumeUniverseProposal:
		var payload types.ResumeUniverseProposalPayload
		if err := cbor.Unmarshal(tx.Data.Payload, &payload); err != nil {
			return &abcitypes.CheckTxResponse{
				Code: types.InvalidProposalPayloadError,
			}
		}
		// Check if universe exists
		universe, err := ctx.Db.Universes.Get(payload.UniverseId)
		if err != nil {
			return &abcitypes.CheckTxResponse{
				Code: types.UniverseNotFoundError,
			}
		}
		// Check if universe is paused
		if !universe.Paused {
			return &abcitypes.CheckTxResponse{
				Code: types.UniverseNotPausedError,
			}
		}

	case types.UpdateUniverseConfigProposal:
		var payload types.UpdateUniverseConfigProposalPayload
		if err := cbor.Unmarshal(tx.Data.Payload, &payload); err != nil {
			return &abcitypes.CheckTxResponse{
				Code: types.InvalidProposalPayloadError,
			}
		}
		// Check if universe exists
		_, err := ctx.Db.Universes.Get(payload.UniverseId)
		if err != nil {
			return &abcitypes.CheckTxResponse{
				Code: types.UniverseNotFoundError,
			}
		}

	case types.UpdateSpaceshipStatsProposal:
		var payload types.UpdateSpaceshipStatsProposalPayload
		if err := cbor.Unmarshal(tx.Data.Payload, &payload); err != nil {
			return &abcitypes.CheckTxResponse{
				Code: types.InvalidProposalPayloadError,
			}
		}
		// Validate payload has required fields
		if payload.SpaceshipId == "" {
			return &abcitypes.CheckTxResponse{
				Code: types.InvalidProposalPayloadError,
			}
		}

	case types.UpdateDefenseStatsProposal:
		var payload types.UpdateDefenseStatsProposalPayload
		if err := cbor.Unmarshal(tx.Data.Payload, &payload); err != nil {
			return &abcitypes.CheckTxResponse{
				Code: types.InvalidProposalPayloadError,
			}
		}
		// Validate payload has required fields
		if payload.DefenseId == "" {
			return &abcitypes.CheckTxResponse{
				Code: types.InvalidProposalPayloadError,
			}
		}

	default:
		return &abcitypes.CheckTxResponse{
			Code: types.InvalidProposalTypeError,
		}
	}

	return &abcitypes.CheckTxResponse{
		Code: types.NoError,
	}
}

func (tx *CreateProposalTransaction) Execute(ctx t.TransactionContext) *abcitypes.ExecTxResult {
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

	// Get current proposal count from state to generate new proposal ID
	proposalId := ctx.State.TotalGovernanceProposalsCreated + 1
	ctx.State.TotalGovernanceProposalsCreated = proposalId

	// Create proposal
	votingPeriodEnd := uint64(ctx.Date.Unix()) + types.ProposalVotingPeriod
	proposal := types.OChainGovernanceProposal{
		Id:               proposalId,
		Type:             tx.Data.ProposalType,
		Payload:          tx.Data.Payload,
		TotalVote:        0,
		TotalForVote:     0,
		TotalAbstainVote: 0,
		TotalAgainstVote: 0,
		CreatedAt:        uint64(ctx.Date.Unix()),
		EndingAt:         votingPeriodEnd,
		Executed:         false,
	}

	err = ctx.Db.GovernanceProposals.Insert(proposal)
	if err != nil {
		return &abcitypes.ExecTxResult{
			Code: types.DatabaseOperationError,
		}
	}

	events := []abcitypes.Event{
		{
			Type: "ProposalCreated",
			Attributes: []abcitypes.EventAttribute{
				{Key: "proposalId", Value: fmt.Sprint(proposalId), Index: true},
				{Key: "proposalType", Value: fmt.Sprint(tx.Data.ProposalType), Index: true},
				{Key: "creator", Value: tx.From, Index: true},
				{Key: "endingAt", Value: fmt.Sprint(votingPeriodEnd)},
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

func ParseCreateProposalTransaction(tx t.Transaction) (*CreateProposalTransaction, error) {
	var txData CreateProposalTransactionData
	err := cbor.Unmarshal(tx.Data, &txData)

	if err != nil {
		return nil, err
	}

	return &CreateProposalTransaction{
		Type:      tx.Type,
		From:      tx.From,
		Nonce:     tx.Nonce,
		Data:      txData,
		Signature: tx.Signature,
	}, nil
}

// ParseCreateProposalTransactionHandler is a wrapper that returns ITransactionHandler for registry
func ParseCreateProposalTransactionHandler(tx t.Transaction) (t.ITransactionHandler, error) {
	return ParseCreateProposalTransaction(tx)
}
