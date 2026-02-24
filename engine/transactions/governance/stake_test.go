package governance_transactions

import (
	"testing"

	tx "github.com/ochain-gg/ochain-network/engine/transactions"
)

func TestStakeTransactionGetType(t *testing.T) {
	stakeTx := &StakeTransaction{
		Type: tx.GovernanceStake,
	}

	if stakeTx.GetType() != tx.GovernanceStake {
		t.Errorf("Expected type %d, got %d", tx.GovernanceStake, stakeTx.GetType())
	}
}

func TestStakeTransactionParsing(t *testing.T) {
	// Create a stake transaction
	stakeTx := &StakeTransaction{
		Type:  tx.GovernanceStake,
		From:  "0xtest",
		Nonce: 1,
		Data: StakeTransactionData{
			Amount: 100000,
		},
		Signature: []byte("test-signature"),
	}

	// Convert to generic transaction
	genericTx, err := stakeTx.Transaction()
	if err != nil {
		t.Fatalf("Failed to convert to generic tx: %v", err)
	}

	// Parse back using the handler
	handler, err := tx.Parse(genericTx)
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	if handler.GetType() != tx.GovernanceStake {
		t.Errorf("Expected type %d, got %d", tx.GovernanceStake, handler.GetType())
	}

	// Verify it's the correct type
	parsedTx, ok := handler.(*StakeTransaction)
	if !ok {
		t.Fatal("Handler is not a StakeTransaction")
	}

	if parsedTx.From != stakeTx.From {
		t.Errorf("Expected From %s, got %s", stakeTx.From, parsedTx.From)
	}

	if parsedTx.Data.Amount != stakeTx.Data.Amount {
		t.Errorf("Expected Amount %d, got %d", stakeTx.Data.Amount, parsedTx.Data.Amount)
	}
}

func TestStakeTransactionRegistered(t *testing.T) {
	if !tx.IsRegistered(tx.GovernanceStake) {
		t.Error("GovernanceStake should be registered")
	}
}
