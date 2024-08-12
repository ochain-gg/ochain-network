package types

type GovernanceProposalType uint64

type VoteType uint64

const (
	NewUniverseProposal             GovernanceProposalType = 0
	NewRewardProgramProposal        GovernanceProposalType = 1
	GameConfigurationChangeProposal GovernanceProposalType = 2

	ForProposalVote     VoteType = 0
	AbstainProposalVote VoteType = 1
	AgainstProposalVote VoteType = 2

	VotingPowerRequiredForProposalCreation uint64 = 100_000_000_000_000 // 100K staked OCT required to create a proposal
	VoteQuorumForProposalExecution         uint64 = 100_000_000_000_000 // 100K quorum
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
