package database

import (
	"testing"

	"github.com/dgraph-io/badger/v4"
	"github.com/ochain-gg/ochain-network/types"
)

func TestOChainStateTable(t *testing.T) {
	// Create in-memory BadgerDB
	opts := badger.DefaultOptions("").WithInMemory(true).WithLogger(nil)
	bdb, err := badger.OpenManaged(opts)
	if err != nil {
		t.Fatalf("Failed to create in-memory database: %v", err)
	}
	defer bdb.Close()
	
	// Create database instance with in-memory DB
	db := &OChainDatabase{
		DB: bdb,
	}
	db.LoadTables()

	// Test initial state (should not exist)
	exists, err := db.State.Exists()
	if err != nil {
		t.Fatalf("Failed to check if state exists: %v", err)
	}
	if exists {
		t.Error("State should not exist initially")
	}

	// Get initial state (should return empty state)
	state, err := db.State.Get()
	if err != nil {
		t.Fatalf("Failed to get initial state: %v", err)
	}
	if state.Height != 0 {
		t.Errorf("Expected initial height 0, got %d", state.Height)
	}

	// Test upsert state
	txn := db.DB.NewTransactionAt(1, true)
	defer txn.Discard()
	
	db.State.SetCurrentTxn(txn)
	
	newState := types.OChainState{
		Size:               1,
		Height:             100,
		Hash:               []byte("test hash"),
		LatestPortalUpdate: 123456,
		OCTTotalSupply:     1000000,
	}

	err = db.State.Upsert(newState)
	if err != nil {
		t.Fatalf("Failed to upsert state: %v", err)
	}

	err = txn.CommitAt(txn.ReadTs(), nil)
	if err != nil {
		t.Fatalf("Failed to commit transaction: %v", err)
	}

	// Verify state was saved
	exists, err = db.State.Exists()
	if err != nil {
		t.Fatalf("Failed to check if state exists after upsert: %v", err)
	}
	if !exists {
		t.Error("State should exist after upsert")
	}

	// Get saved state
	savedState, err := db.State.Get()
	if err != nil {
		t.Fatalf("Failed to get saved state: %v", err)
	}

	// Verify all fields
	if savedState.Size != newState.Size {
		t.Errorf("Expected size %d, got %d", newState.Size, savedState.Size)
	}
	if savedState.Height != newState.Height {
		t.Errorf("Expected height %d, got %d", newState.Height, savedState.Height)
	}
	if string(savedState.Hash) != string(newState.Hash) {
		t.Errorf("Expected hash %s, got %s", string(newState.Hash), string(savedState.Hash))
	}
	if savedState.LatestPortalUpdate != newState.LatestPortalUpdate {
		t.Errorf("Expected portal update %d, got %d", newState.LatestPortalUpdate, savedState.LatestPortalUpdate)
	}
	if savedState.OCTTotalSupply != newState.OCTTotalSupply {
		t.Errorf("Expected OCT supply %d, got %d", newState.OCTTotalSupply, savedState.OCTTotalSupply)
	}
}

func TestOChainStateHelperMethods(t *testing.T) {
	state := &types.OChainState{
		Height:                     0,
		Size:                       0,
		OCTInGameCirculatingSupply: 1000,
	}

	// Test SetHeight
	state.SetHeight(100)
	if state.Height != 100 {
		t.Errorf("Expected height 100, got %d", state.Height)
	}

	// Test IncSize
	originalSize := state.Size
	state.IncSize()
	if state.Size != originalSize+1 {
		t.Errorf("Expected size %d, got %d", originalSize+1, state.Size)
	}

	// Test AddInGameCirculatingSupply
	originalSupply := state.OCTInGameCirculatingSupply
	state.AddInGameCirculatingSupply(500)
	expectedSupply := originalSupply + 500
	if state.OCTInGameCirculatingSupply != expectedSupply {
		t.Errorf("Expected supply %d, got %d", expectedSupply, state.OCTInGameCirculatingSupply)
	}

	// Test RemoveInGameCirculatingSupply
	state.RemoveInGameCirculatingSupply(200)
	expectedSupply = originalSupply + 500 - 200
	if state.OCTInGameCirculatingSupply != expectedSupply {
		t.Errorf("Expected supply %d, got %d", expectedSupply, state.OCTInGameCirculatingSupply)
	}

	// Test SetLatestPortalUpdate
	state.SetLatestPortalUpdate(789)
	if state.LatestPortalUpdate != 789 {
		t.Errorf("Expected portal update 789, got %d", state.LatestPortalUpdate)
	}
}

func TestOChainDatabaseCreation(t *testing.T) {
	// Create in-memory BadgerDB
	opts := badger.DefaultOptions("").WithInMemory(true).WithLogger(nil)
	bdb, err := badger.OpenManaged(opts)
	if err != nil {
		t.Fatalf("Failed to create in-memory database: %v", err)
	}
	defer bdb.Close()
	
	// Create database instance with in-memory DB
	db := &OChainDatabase{
		DB: bdb,
	}
	db.LoadTables()

	// Test that all table instances are created
	if db.State == nil {
		t.Error("State table should not be nil")
	}
	if db.GlobalsAccounts == nil {
		t.Error("GlobalsAccounts table should not be nil")
	}
	if db.UniverseAccounts == nil {
		t.Error("UniverseAccounts table should not be nil")
	}
	if db.Planets == nil {
		t.Error("Planets table should not be nil")
	}
	if db.Fleets == nil {
		t.Error("Fleets table should not be nil")
	}

	// Test that database connection works
	txn := db.DB.NewTransactionAt(1, false)
	defer txn.Discard()
	
	// Should be able to create a read transaction without errors
	if txn == nil {
		t.Error("Should be able to create read transaction")
	}
}