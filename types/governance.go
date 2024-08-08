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
	VoteQuorumForProposalExecution         uint64 = 100_000_000_000_000 // 100K staked OCT required to create a proposal
)

type OChainGovernanceProposal struct {
	Id   uint64                 `cbor:"1,keyasint"`
	Type GovernanceProposalType `cbor:"2,keyasint"`
}

type OChainGovernanceProposalVote struct {
	Id         string `cbor:"1,keyasint"`
	ProposalId uint64 `cbor:"2,keyasint"`
	Voter      string `cbor:"3,keyasint"`
}
