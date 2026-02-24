package planet_transactions

import (
	"log"

	t "github.com/ochain-gg/ochain-network/engine/transactions"
)

// init registers all planet transaction handlers with the global registry
func init() {
	// Register planet transactions
	if err := t.Register(t.StartBuildingUpgrade, ParseUpgradeBuildingTransactionHandler); err != nil {
		log.Fatalf("Failed to register StartBuildingUpgrade: %v", err)
	}
	if err := t.Register(t.StartTechnologyUpgrade, ParseUpgradeTechnologyTransactionHandler); err != nil {
		log.Fatalf("Failed to register StartTechnologyUpgrade: %v", err)
	}
	if err := t.Register(t.Build, ParseBuildTransactionHandler); err != nil {
		log.Fatalf("Failed to register Build: %v", err)
	}
	if err := t.Register(t.SwapResources, ParseSwapResourcesTransactionHandler); err != nil {
		log.Fatalf("Failed to register SwapResources: %v", err)
	}
	if err := t.Register(t.ExecutePendingUpgrade, ParseExecuteUpgradeTransactionHandler); err != nil {
		log.Fatalf("Failed to register ExecutePendingUpgrade: %v", err)
	}

	log.Println("Planet transactions handler loaded")
}
