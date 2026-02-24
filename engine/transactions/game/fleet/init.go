package fleet_transactions

import (
	"log"

	t "github.com/ochain-gg/ochain-network/engine/transactions"
)

// init registers all fleet transaction handlers with the global registry
func init() {
	// Register fleet transactions
	if err := t.Register(t.CreateFleet, ParseCreateFleetTransactionHandler); err != nil {
		log.Fatalf("Failed to register CreateFleet: %v", err)
	}
	if err := t.Register(t.FillCargo, ParseFeelFleetCargoTransactionHandler); err != nil {
		log.Fatalf("Failed to register FillCargo: %v", err)
	}
	if err := t.Register(t.UnfillCargo, ParseUnfeelFleetCargoTransactionHandler); err != nil {
		log.Fatalf("Failed to register UnfillCargo: %v", err)
	}
	if err := t.Register(t.LandingFleetInOrbit, ParseLendFleetCargoTransactionHandler); err != nil {
		log.Fatalf("Failed to register LandingFleetInOrbit: %v", err)
	}
	if err := t.Register(t.MergeFleets, ParseMergeFleetTransactionHandler); err != nil {
		log.Fatalf("Failed to register MergeFleets: %v", err)
	}
	if err := t.Register(t.SplitFleet, ParseSplitFleetTransactionHandler); err != nil {
		log.Fatalf("Failed to register SplitFleet: %v", err)
	}
	if err := t.Register(t.Fight, ParseFleetAttackTransactionHandler); err != nil {
		log.Fatalf("Failed to register Fight: %v", err)
	}

	log.Println("Fleet transactions handler loaded")
}
