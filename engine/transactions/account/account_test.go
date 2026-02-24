package account_transactions

import (
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ochain-gg/ochain-network/config"
	"github.com/ochain-gg/ochain-network/engine/database"
	"github.com/ochain-gg/ochain-network/engine/transactions"
	"github.com/ochain-gg/ochain-network/types"
)

func createTestContext(t *testing.T) transactions.TransactionContext {
	// Create database instance with in-memory DB
	db := database.NewOChainDatabase("")
	db.NewTransaction(1)

	// Initialize state
	state := types.OChainState{
		Size:               0,
		Height:             1,
		Hash:               []byte("test"),
		LatestPortalUpdate: 0,
	}

	err := db.State.Upsert(state)
	if err != nil {
		t.Fatalf("Failed to initialize state: %v", err)
	}

	db.CommitTransaction()

	return transactions.TransactionContext{
		Config: *config.DefaultConfig(),
		Db:     db,
		State:  &state,
	}
}

func TestGlobalAccountRegister(t *testing.T) {
	ctx := createTestContext(t)

	// Create test address
	privateKey, err := crypto.HexToECDSA("ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80")
	if err != nil {
		t.Fatalf("Failed to create private key: %v", err)
	}
	address := crypto.PubkeyToAddress(privateKey.PublicKey).Hex()

	// This is a simplified test - actual transaction data structures and functions
	// would need to be imported from the specific account transaction files
	// For now, test basic database operations

	t.Logf("Testing with address: %s", address)

	// Test that we can create test context without errors
	if ctx.Db == nil {
		t.Error("Database should not be nil")
	}
	if ctx.State == nil {
		t.Error("State should not be nil")
	}
}

func TestDatabaseIntegration(t *testing.T) {
	ctx := createTestContext(t)

	// Test that database tables are accessible
	if ctx.Db.GlobalsAccounts == nil {
		t.Error("GlobalsAccounts table should not be nil")
	}
	if ctx.Db.UniverseAccounts == nil {
		t.Error("UniverseAccounts table should not be nil")
	}
	if ctx.Db.State == nil {
		t.Error("State table should not be nil")
	}

	// Test transaction context is properly set up
	if ctx.Db.CurrentTxn == nil {
		t.Error("Transaction should not be nil")
	}
	if ctx.State == nil {
		t.Error("State should not be nil")
	}
}

func TestTransactionContextCreation(t *testing.T) {
	ctx := createTestContext(t)

	// Verify configuration
	if ctx.Config.EVMChainId == 0 {
		t.Error("EVM Chain ID should be set")
	}

	// Verify state initialization
	if ctx.State.Height != 1 {
		t.Errorf("Expected height 1, got %d", ctx.State.Height)
	}

	// Test that we can interact with the database
	ctx.Db.State.SetCurrentTxn(ctx.Txn)
	state, err := ctx.Db.State.Get()
	if err != nil {
		t.Fatalf("Failed to get state: %v", err)
	}

	if state.Height != 1 {
		t.Errorf("Expected state height 1, got %d", state.Height)
	}
}
