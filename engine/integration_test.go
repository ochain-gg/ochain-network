package engine

import (
	"testing"

	"github.com/ochain-gg/ochain-network/engine/transactions"
)

// TestGovernanceTransactionsRegistered verifies all governance transactions are registered in the registry
func TestGovernanceTransactionsRegistered(t *testing.T) {
	governanceTypes := []transactions.TransactionType{
		transactions.GovernanceStake,
		transactions.GovernanceUnstakeStart,
		transactions.GovernanceUnstakeExecute,
		transactions.CreateGovernanceProposal,
		transactions.VoteOnProposal,
		transactions.ExecuteProposal,
	}

	for _, txType := range governanceTypes {
		if !transactions.IsRegistered(txType) {
			t.Errorf("Governance transaction type %d should be registered", txType)
		}
	}
}

// TestRegistryContainsGovernanceTransactions verifies the registry has the correct number of types
func TestRegistryContainsGovernanceTransactions(t *testing.T) {
	registeredTypes := transactions.GetRegisteredTypes()

	// Should have at least 6 governance transactions registered
	if len(registeredTypes) < 6 {
		t.Errorf("Expected at least 6 registered types, got %d", len(registeredTypes))
	}

	// Verify governance types are in the list
	found := make(map[transactions.TransactionType]bool)
	for _, txType := range registeredTypes {
		found[txType] = true
	}

	governanceTypes := []transactions.TransactionType{
		transactions.GovernanceStake,
		transactions.GovernanceUnstakeStart,
		transactions.GovernanceUnstakeExecute,
		transactions.CreateGovernanceProposal,
		transactions.VoteOnProposal,
		transactions.ExecuteProposal,
	}

	for _, txType := range governanceTypes {
		if !found[txType] {
			t.Errorf("Governance transaction type %d not found in registered types", txType)
		}
	}
}

// TestAllTransactionTypesRegistered verifies all 19 transaction types are registered
func TestAllTransactionTypesRegistered(t *testing.T) {
	allTransactionTypes := []transactions.TransactionType{
		// Governance (6)
		transactions.GovernanceStake,
		transactions.GovernanceUnstakeStart,
		transactions.GovernanceUnstakeExecute,
		transactions.CreateGovernanceProposal,
		transactions.VoteOnProposal,
		transactions.ExecuteProposal,
		// Account (4)
		transactions.RegisterAccount,
		transactions.RegisterUniverseAccount,
		transactions.UniverseOChainTokenWithdraw,
		transactions.ClaimFaucet,
		// Planet (5)
		transactions.StartBuildingUpgrade,
		transactions.StartTechnologyUpgrade,
		transactions.Build,
		transactions.SwapResources,
		transactions.ExecutePendingUpgrade,
		// Validator/System (4)
		transactions.NewValidator,
		transactions.RemoveValidator,
		transactions.OChainTokenDeposit,
		transactions.OChainCreditDeposit,
	}

	registeredTypes := transactions.GetRegisteredTypes()
	expectedCount := len(allTransactionTypes)

	if len(registeredTypes) != expectedCount {
		t.Errorf("Expected exactly %d registered types, got %d", expectedCount, len(registeredTypes))
	}

	for _, txType := range allTransactionTypes {
		if !transactions.IsRegistered(txType) {
			t.Errorf("Transaction type %d should be registered", txType)
		}
	}

	t.Logf("✓ All %d transaction types successfully registered", expectedCount)
}
