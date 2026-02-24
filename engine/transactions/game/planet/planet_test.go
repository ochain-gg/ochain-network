package planet_transactions

import (
	"testing"

	tx "github.com/ochain-gg/ochain-network/engine/transactions"
)

func TestUpgradeBuildingTransactionGetType(t *testing.T) {
	upgradeTx := &UpgradeBuildingTransaction{
		Type: tx.StartBuildingUpgrade,
	}

	if upgradeTx.GetType() != tx.StartBuildingUpgrade {
		t.Errorf("Expected type %d, got %d", tx.StartBuildingUpgrade, upgradeTx.GetType())
	}
}

func TestUpgradeBuildingTransactionRegistered(t *testing.T) {
	if !tx.IsRegistered(tx.StartBuildingUpgrade) {
		t.Error("StartBuildingUpgrade should be registered")
	}
}

func TestUpgradeTechnologyTransactionGetType(t *testing.T) {
	upgradeTx := &UpgradeTechnologyTransaction{
		Type: tx.StartTechnologyUpgrade,
	}

	if upgradeTx.GetType() != tx.StartTechnologyUpgrade {
		t.Errorf("Expected type %d, got %d", tx.StartTechnologyUpgrade, upgradeTx.GetType())
	}
}

func TestUpgradeTechnologyTransactionRegistered(t *testing.T) {
	if !tx.IsRegistered(tx.StartTechnologyUpgrade) {
		t.Error("StartTechnologyUpgrade should be registered")
	}
}

func TestBuildTransactionGetType(t *testing.T) {
	buildTx := &BuildTransaction{
		Type: tx.Build,
	}

	if buildTx.GetType() != tx.Build {
		t.Errorf("Expected type %d, got %d", tx.Build, buildTx.GetType())
	}
}

func TestBuildTransactionRegistered(t *testing.T) {
	if !tx.IsRegistered(tx.Build) {
		t.Error("Build should be registered")
	}
}

func TestSwapResourcesTransactionGetType(t *testing.T) {
	swapTx := &SwapResourcesTransaction{
		Type: tx.SwapResources,
	}

	if swapTx.GetType() != tx.SwapResources {
		t.Errorf("Expected type %d, got %d", tx.SwapResources, swapTx.GetType())
	}
}

func TestSwapResourcesTransactionRegistered(t *testing.T) {
	if !tx.IsRegistered(tx.SwapResources) {
		t.Error("SwapResources should be registered")
	}
}

func TestExecuteUpgradeTransactionGetType(t *testing.T) {
	executeUpgradeTx := &ExecuteUpgradeTransaction{
		Type: tx.ExecutePendingUpgrade,
	}

	if executeUpgradeTx.GetType() != tx.ExecutePendingUpgrade {
		t.Errorf("Expected type %d, got %d", tx.ExecutePendingUpgrade, executeUpgradeTx.GetType())
	}
}

func TestExecuteUpgradeTransactionRegistered(t *testing.T) {
	if !tx.IsRegistered(tx.ExecutePendingUpgrade) {
		t.Error("ExecutePendingUpgrade should be registered")
	}
}

func TestAllPlanetTransactionsRegistered(t *testing.T) {
	planetTypes := []tx.TransactionType{
		tx.StartBuildingUpgrade,
		tx.StartTechnologyUpgrade,
		tx.Build,
		tx.SwapResources,
		tx.ExecutePendingUpgrade,
	}

	for _, txType := range planetTypes {
		if !tx.IsRegistered(txType) {
			t.Errorf("Planet transaction type %d should be registered", txType)
		}
	}
}
