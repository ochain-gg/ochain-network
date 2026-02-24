package validator_transactions

import (
	"log"

	t "github.com/ochain-gg/ochain-network/engine/transactions"
)

// init registers all validator/system transaction handlers with the global registry
func init() {
	// Register validator transactions
	if err := t.Register(t.NewValidator, ParseOChainNewValidatorTransactionHandler); err != nil {
		log.Fatalf("Failed to register NewValidator: %v", err)
	}
	if err := t.Register(t.RemoveValidator, ParseOChainRemoveValidatorTransactionHandler); err != nil {
		log.Fatalf("Failed to register RemoveValidator: %v", err)
	}
	if err := t.Register(t.OChainTokenDeposit, ParseOChainTokenDepositTransactionHandler); err != nil {
		log.Fatalf("Failed to register OChainTokenDeposit: %v", err)
	}
	if err := t.Register(t.OChainCreditDeposit, ParseOChainCreditDepositTransactionHandler); err != nil {
		log.Fatalf("Failed to register OChainCreditDeposit: %v", err)
	}

	log.Println("Validator transactions handler loaded")
}
