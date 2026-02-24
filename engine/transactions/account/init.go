package account_transactions

import (
	"log"

	t "github.com/ochain-gg/ochain-network/engine/transactions"
)

// init registers all account transaction handlers with the global registry
func init() {
	// Register account transactions
	if err := t.Register(t.RegisterAccount, ParseRegisterAccountTransactionHandler); err != nil {
		log.Fatalf("Failed to register RegisterAccount: %v", err)
	}
	if err := t.Register(t.RegisterUniverseAccount, ParseRegisterUniverseAccountTransactionHandler); err != nil {
		log.Fatalf("Failed to register RegisterUniverseAccount: %v", err)
	}
	if err := t.Register(t.UniverseOChainTokenWithdraw, ParseUniverseAccountWithdrawTransactionHandler); err != nil {
		log.Fatalf("Failed to register UniverseOChainTokenWithdraw: %v", err)
	}
	if err := t.Register(t.ClaimFaucet, ParseClaimFaucetTransactionHandler); err != nil {
		log.Fatalf("Failed to register ClaimFaucet: %v", err)
	}

	log.Println("Account transactions handler loaded")
}
