package governance_transactions

import (
	"fmt"
	"time"

	abcitypes "github.com/cometbft/cometbft/abci/types"
	"github.com/fxamacker/cbor/v2"

	t "github.com/ochain-gg/ochain-network/engine/transactions"
	"github.com/ochain-gg/ochain-network/types"
)

type ExecuteProposalTransactionData struct {
	ProposalId uint64 `cbor:"1,keyasint"`
}

type ExecuteProposalTransaction struct {
	Type      t.TransactionType               `cbor:"1,keyasint"`
	From      string                          `cbor:"2,keyasint"`
	Nonce     uint64                          `cbor:"3,keyasint"`
	Data      ExecuteProposalTransactionData  `cbor:"4,keyasint"`
	Signature []byte                          `cbor:"5,keyasint"`
}

func (tx *ExecuteProposalTransaction) Transaction() (t.Transaction, error) {
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

func (tx *ExecuteProposalTransaction) Check(ctx t.TransactionContext) *abcitypes.CheckTxResponse {
	// Check if account exists
	_, err := ctx.Db.GlobalsAccounts.Get(tx.From)
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

	// Check if proposal has already been executed
	if proposal.Executed {
		return &abcitypes.CheckTxResponse{
			Code: types.ProposalAlreadyExecutedError,
		}
	}

	// Check if voting period has ended
	if uint64(ctx.Date.Unix()) < proposal.EndingAt {
		return &abcitypes.CheckTxResponse{
			Code: types.VotingPeriodNotExpiredError,
		}
	}

	// Check if quorum was reached
	if proposal.TotalVote < types.VoteQuorumForProposalExecution {
		return &abcitypes.CheckTxResponse{
			Code: types.ProposalNotPassedError,
		}
	}

	// Check if proposal passed (more "for" than "against")
	if proposal.TotalForVote <= proposal.TotalAgainstVote {
		return &abcitypes.CheckTxResponse{
			Code: types.ProposalNotPassedError,
		}
	}

	return &abcitypes.CheckTxResponse{
		Code: types.NoError,
	}
}

func (tx *ExecuteProposalTransaction) Execute(ctx t.TransactionContext) *abcitypes.ExecTxResult {
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

	var events []abcitypes.Event

	// Execute proposal based on type
	switch proposal.Type {
	case types.CreateUniverseProposal:
		var payload types.CreateUniverseProposalPayload
		if err := cbor.Unmarshal(proposal.Payload, &payload); err != nil {
			return &abcitypes.ExecTxResult{
				Code: types.InvalidProposalPayloadError,
			}
		}

		// Create new universe
		universe := types.OChainUniverse{
			Id:                      payload.Id,
			Name:                    payload.Name,
			Speed:                   payload.Speed,
			ResourcesMarketEnabled:  payload.ResourcesMarketEnabled,
			IsStretchable:           payload.IsStretchable,
			Paused:                  false,
			MaxGalaxy:               payload.MaxGalaxy,
			MaxSolarSystemPerGalaxy: payload.MaxSolarSystemPerGalaxy,
			MaxPlanetPerSolarSystem: payload.MaxPlanetPerSolarSystem,
			Accounts:                0,
			ColonizedPlanets:        0,
			OCTCirculatingSupply:    0,
			CreatedAt:               uint64(time.Now().Unix()),
			EndingAt:                0,
		}

		err = ctx.Db.Universes.Insert(universe)
		if err != nil {
			return &abcitypes.ExecTxResult{
				Code: types.DatabaseOperationError,
			}
		}

		events = append(events, abcitypes.Event{
			Type: "UniverseCreated",
			Attributes: []abcitypes.EventAttribute{
				{Key: "universeId", Value: payload.Id, Index: true},
				{Key: "name", Value: payload.Name},
				{Key: "maxGalaxy", Value: fmt.Sprint(payload.MaxGalaxy)},
			},
		})

	case types.PauseUniverseProposal:
		var payload types.PauseUniverseProposalPayload
		if err := cbor.Unmarshal(proposal.Payload, &payload); err != nil {
			return &abcitypes.ExecTxResult{
				Code: types.InvalidProposalPayloadError,
			}
		}

		// Get universe
		universe, err := ctx.Db.Universes.Get(payload.UniverseId)
		if err != nil {
			return &abcitypes.ExecTxResult{
				Code: types.UniverseNotFoundError,
			}
		}

		// Pause universe
		universe.Paused = true
		err = ctx.Db.Universes.Update(universe)
		if err != nil {
			return &abcitypes.ExecTxResult{
				Code: types.DatabaseOperationError,
			}
		}

		events = append(events, abcitypes.Event{
			Type: "UniversePaused",
			Attributes: []abcitypes.EventAttribute{
				{Key: "universeId", Value: payload.UniverseId, Index: true},
			},
		})

	case types.ResumeUniverseProposal:
		var payload types.ResumeUniverseProposalPayload
		if err := cbor.Unmarshal(proposal.Payload, &payload); err != nil {
			return &abcitypes.ExecTxResult{
				Code: types.InvalidProposalPayloadError,
			}
		}

		// Get universe
		universe, err := ctx.Db.Universes.Get(payload.UniverseId)
		if err != nil {
			return &abcitypes.ExecTxResult{
				Code: types.UniverseNotFoundError,
			}
		}

		// Resume universe
		universe.Paused = false
		err = ctx.Db.Universes.Update(universe)
		if err != nil {
			return &abcitypes.ExecTxResult{
				Code: types.DatabaseOperationError,
			}
		}

		events = append(events, abcitypes.Event{
			Type: "UniverseResumed",
			Attributes: []abcitypes.EventAttribute{
				{Key: "universeId", Value: payload.UniverseId, Index: true},
			},
		})

	case types.UpdateUniverseConfigProposal:
		var payload types.UpdateUniverseConfigProposalPayload
		if err := cbor.Unmarshal(proposal.Payload, &payload); err != nil {
			return &abcitypes.ExecTxResult{
				Code: types.InvalidProposalPayloadError,
			}
		}

		// Get universe
		universe, err := ctx.Db.Universes.Get(payload.UniverseId)
		if err != nil {
			return &abcitypes.ExecTxResult{
				Code: types.UniverseNotFoundError,
			}
		}

		// Update universe configuration
		universe.MaxGalaxy = payload.MaxGalaxy
		universe.MaxSolarSystemPerGalaxy = payload.MaxSolarSystemPerGalaxy
		universe.MaxPlanetPerSolarSystem = payload.MaxPlanetPerSolarSystem
		universe.Speed = payload.Speed
		universe.ResourcesMarketEnabled = payload.ResourcesMarketEnabled
		universe.IsStretchable = payload.IsStretchable

		err = ctx.Db.Universes.Update(universe)
		if err != nil {
			return &abcitypes.ExecTxResult{
				Code: types.DatabaseOperationError,
			}
		}

		events = append(events, abcitypes.Event{
			Type: "UniverseConfigUpdated",
			Attributes: []abcitypes.EventAttribute{
				{Key: "universeId", Value: payload.UniverseId, Index: true},
				{Key: "maxGalaxy", Value: fmt.Sprint(payload.MaxGalaxy)},
				{Key: "speed", Value: fmt.Sprint(payload.Speed)},
			},
		})

	case types.UpdateSpaceshipStatsProposal:
		var payload types.UpdateSpaceshipStatsProposalPayload
		if err := cbor.Unmarshal(proposal.Payload, &payload); err != nil {
			return &abcitypes.ExecTxResult{
				Code: types.InvalidProposalPayloadError,
			}
		}

		// Get spaceship from database
		spaceship, err := ctx.Db.Spaceships.Get(string(payload.SpaceshipId))
		if err != nil {
			return &abcitypes.ExecTxResult{
				Code: types.SpaceshipNotFoundError,
			}
		}

		// Update spaceship stats (global change affecting all universes)
		spaceship.FighterStats = payload.FighterStats
		spaceship.Stats = payload.Stats

		err = ctx.Db.Spaceships.Update(spaceship)
		if err != nil {
			return &abcitypes.ExecTxResult{
				Code: types.DatabaseOperationError,
			}
		}

		events = append(events, abcitypes.Event{
			Type: "SpaceshipStatsUpdated",
			Attributes: []abcitypes.EventAttribute{
				{Key: "spaceshipId", Value: string(payload.SpaceshipId), Index: true},
				{Key: "armor", Value: fmt.Sprint(payload.FighterStats.Armor)},
				{Key: "shield", Value: fmt.Sprint(payload.FighterStats.Shield)},
				{Key: "attack", Value: fmt.Sprint(payload.FighterStats.Attack)},
			},
		})

	case types.UpdateDefenseStatsProposal:
		var payload types.UpdateDefenseStatsProposalPayload
		if err := cbor.Unmarshal(proposal.Payload, &payload); err != nil {
			return &abcitypes.ExecTxResult{
				Code: types.InvalidProposalPayloadError,
			}
		}

		// Get defense from database
		defense, err := ctx.Db.Defenses.Get(string(payload.DefenseId))
		if err != nil {
			return &abcitypes.ExecTxResult{
				Code: types.DefenseNotFoundError,
			}
		}

		// Update defense stats (global change affecting all universes)
		defense.FighterStats = payload.FighterStats

		err = ctx.Db.Defenses.Update(defense)
		if err != nil {
			return &abcitypes.ExecTxResult{
				Code: types.DatabaseOperationError,
			}
		}

		events = append(events, abcitypes.Event{
			Type: "DefenseStatsUpdated",
			Attributes: []abcitypes.EventAttribute{
				{Key: "defenseId", Value: string(payload.DefenseId), Index: true},
				{Key: "armor", Value: fmt.Sprint(payload.FighterStats.Armor)},
				{Key: "shield", Value: fmt.Sprint(payload.FighterStats.Shield)},
				{Key: "attack", Value: fmt.Sprint(payload.FighterStats.Attack)},
			},
		})

	default:
		return &abcitypes.ExecTxResult{
			Code: types.InvalidProposalTypeError,
		}
	}

	// Mark proposal as executed
	proposal.Executed = true
	err = ctx.Db.GovernanceProposals.Update(proposal)
	if err != nil {
		return &abcitypes.ExecTxResult{
			Code: types.DatabaseOperationError,
		}
	}

	// Add proposal executed event
	events = append(events, abcitypes.Event{
		Type: "ProposalExecuted",
		Attributes: []abcitypes.EventAttribute{
			{Key: "proposalId", Value: fmt.Sprint(tx.Data.ProposalId), Index: true},
			{Key: "proposalType", Value: fmt.Sprint(proposal.Type)},
			{Key: "executor", Value: tx.From, Index: true},
		},
	})

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

func ParseExecuteProposalTransaction(tx t.Transaction) (ExecuteProposalTransaction, error) {
	var txData ExecuteProposalTransactionData
	err := cbor.Unmarshal(tx.Data, &txData)

	if err != nil {
		return ExecuteProposalTransaction{}, err
	}

	return ExecuteProposalTransaction{
		Type:      tx.Type,
		From:      tx.From,
		Nonce:     tx.Nonce,
		Data:      txData,
		Signature: tx.Signature,
	}, nil
}
