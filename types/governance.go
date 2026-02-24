package types

type GovernanceProposalType uint64

type VoteType uint64

const (
	CreateUniverseProposal        GovernanceProposalType = 0
	PauseUniverseProposal         GovernanceProposalType = 1
	ResumeUniverseProposal        GovernanceProposalType = 2
	UpdateUniverseConfigProposal  GovernanceProposalType = 3
	UpdateSpaceshipStatsProposal  GovernanceProposalType = 4
	UpdateDefenseStatsProposal    GovernanceProposalType = 5
	NewRewardProgramProposal      GovernanceProposalType = 6
	GameConfigurationChangeProposal GovernanceProposalType = 7

	ForProposalVote     VoteType = 0
	AbstainProposalVote VoteType = 1
	AgainstProposalVote VoteType = 2

	VotingPowerRequiredForProposalCreation uint64 = 100_000_000_000_000 // 100K staked OCT required to create a proposal
	VoteQuorumForProposalExecution         uint64 = 100_000_000_000_000 // 100K quorum
	ProposalVotingPeriod                   uint64 = 7 * 24 * 60 * 60    // 7 days in seconds
)

type OChainGovernanceProposal struct {
	Id      uint64                 `cbor:"1,keyasint"`
	Type    GovernanceProposalType `cbor:"2,keyasint"`
	Payload []byte                 `cbor:"3,keyasint"`

	TotalVote        uint64 `cbor:"4,keyasint"`
	TotalForVote     uint64 `cbor:"5,keyasint"`
	TotalAbstainVote uint64 `cbor:"6,keyasint"`
	TotalAgainstVote uint64 `cbor:"7,keyasint"`

	CreatedAt uint64 `cbor:"8,keyasint"`
	EndingAt  uint64 `cbor:"9,keyasint"`
	Executed  bool   `cbor:"10,keyasint"`
}

type OChainGovernanceProposalVote struct {
	Id         string   `cbor:"1,keyasint"`
	ProposalId uint64   `cbor:"2,keyasint"`
	Voter      string   `cbor:"3,keyasint"`
	Vote       VoteType `cbor:"4,keyasint"`

	CreatedAt uint64 `cbor:"5,keyasint"`
	UpdatedAt uint64 `cbor:"6,keyasint"`
}

type OChainPendingUnstake struct {
	Address         string `cbor:"1,keyasint"`
	Amount          uint64 `cbor:"2,keyasint"`
	StartedAt       int64  `cbor:"3,keyasint"`
	WithdrawableAt  int64  `cbor:"4,keyasint"`
}

type OChainPendingUnstakeWithAttributes struct {
	Address         string `cbor:"address"`
	Amount          uint64 `cbor:"amount"`
	StartedAt       int64  `cbor:"startedAt"`
	WithdrawableAt  int64  `cbor:"withdrawableAt"`
}

func (pu *OChainPendingUnstake) WithAttributes() OChainPendingUnstakeWithAttributes {
	return OChainPendingUnstakeWithAttributes{
		Address:        pu.Address,
		Amount:         pu.Amount,
		StartedAt:      pu.StartedAt,
		WithdrawableAt: pu.WithdrawableAt,
	}
}

// Proposal payload types

type CreateUniverseProposalPayload struct {
	Id                      string `cbor:"1,keyasint"`
	Name                    string `cbor:"2,keyasint"`
	Speed                   uint64 `cbor:"3,keyasint"`
	ResourcesMarketEnabled  bool   `cbor:"4,keyasint"`
	IsStretchable           bool   `cbor:"5,keyasint"`
	MaxGalaxy               uint64 `cbor:"6,keyasint"`
	MaxSolarSystemPerGalaxy uint64 `cbor:"7,keyasint"`
	MaxPlanetPerSolarSystem uint64 `cbor:"8,keyasint"`
}

type PauseUniverseProposalPayload struct {
	UniverseId string `cbor:"1,keyasint"`
}

type ResumeUniverseProposalPayload struct {
	UniverseId string `cbor:"1,keyasint"`
}

type UpdateUniverseConfigProposalPayload struct {
	UniverseId              string `cbor:"1,keyasint"`
	MaxGalaxy               uint64 `cbor:"2,keyasint"`
	MaxSolarSystemPerGalaxy uint64 `cbor:"3,keyasint"`
	MaxPlanetPerSolarSystem uint64 `cbor:"4,keyasint"`
	Speed                   uint64 `cbor:"5,keyasint"`
	ResourcesMarketEnabled  bool   `cbor:"6,keyasint"`
	IsStretchable           bool   `cbor:"7,keyasint"`
}

type UpdateSpaceshipStatsProposalPayload struct {
	SpaceshipId  OChainSpaceshipID    `cbor:"1,keyasint"`
	FighterStats OChainFighterStats   `cbor:"2,keyasint"`
	Stats        OChainSpaceshipStats `cbor:"3,keyasint"`
}

type UpdateDefenseStatsProposalPayload struct {
	DefenseId    OChainDefenseID    `cbor:"1,keyasint"`
	FighterStats OChainFighterStats `cbor:"2,keyasint"`
}
