package governance_transactions

import (
	"log"

	t "github.com/ochain-gg/ochain-network/engine/transactions"
)

// init registers all governance transaction handlers with the global registry
func init() {
	// Register governance transactions
	if err := t.Register(t.GovernanceStake, ParseStakeTransactionHandler); err != nil {
		log.Fatalf("Failed to register GovernanceStake: %v", err)
	}
	if err := t.Register(t.GovernanceUnstakeStart, ParseUnstakeStartTransactionHandler); err != nil {
		log.Fatalf("Failed to register GovernanceUnstakeStart: %v", err)
	}
	if err := t.Register(t.GovernanceUnstakeExecute, ParseUnstakeExecuteTransactionHandler); err != nil {
		log.Fatalf("Failed to register GovernanceUnstakeExecute: %v", err)
	}
	if err := t.Register(t.CreateGovernanceProposal, ParseCreateProposalTransactionHandler); err != nil {
		log.Fatalf("Failed to register CreateGovernanceProposal: %v", err)
	}
	if err := t.Register(t.VoteOnProposal, ParseVoteOnProposalTransactionHandler); err != nil {
		log.Fatalf("Failed to register VoteOnProposal: %v", err)
	}
	if err := t.Register(t.ExecuteProposal, ParseExecuteProposalTransactionHandler); err != nil {
		log.Fatalf("Failed to register ExecuteProposal: %v", err)
	}

	log.Println("Governance transactions handler loaded")
}
