package alliance_transactions

import (
	"log"

	t "github.com/ochain-gg/ochain-network/engine/transactions"
)

// init registers all alliance transaction handlers with the global registry
func init() {
	// Register alliance transactions
	if err := t.Register(t.CreateAlliance, ParseCreateAllianceTransactionHandler); err != nil {
		log.Fatalf("Failed to register CreateAlliance: %v", err)
	}
	if err := t.Register(t.EditAlliance, ParseEditAllianceTransactionHandler); err != nil {
		log.Fatalf("Failed to register EditAlliance: %v", err)
	}
	if err := t.Register(t.AssignAllianceRole, ParseAssignAllianceRoleTransactionHandler); err != nil {
		log.Fatalf("Failed to register AssignAllianceRole: %v", err)
	}
	if err := t.Register(t.UpgradeAlliance, ParseUpgradeAllianceTransactionHandler); err != nil {
		log.Fatalf("Failed to register UpgradeAlliance: %v", err)
	}
	if err := t.Register(t.JoinAllianceRequest, ParseJoinAllianceRequestTransactionHandler); err != nil {
		log.Fatalf("Failed to register JoinAllianceRequest: %v", err)
	}
	if err := t.Register(t.AnswerJoinAllianceRequest, ParseAnswerJoinAllianceRequestTransactionHandler); err != nil {
		log.Fatalf("Failed to register AnswerJoinAllianceRequest: %v", err)
	}
	if err := t.Register(t.CancelJoinAllianceRequest, ParseCancelJoinAllianceRequestTransactionHandler); err != nil {
		log.Fatalf("Failed to register CancelJoinAllianceRequest: %v", err)
	}

	log.Println("Alliance transactions handler loaded")
}
