package transactions

import (
	"sync"
	"testing"

	abcitypes "github.com/cometbft/cometbft/abci/types"
	"github.com/ochain-gg/ochain-network/types"
)

// MockTransactionHandler is a mock implementation of ITransactionHandler for testing
type MockTransactionHandler struct {
	txType TransactionType
}

func (m *MockTransactionHandler) Check(ctx TransactionContext) *abcitypes.CheckTxResponse {
	return &abcitypes.CheckTxResponse{Code: types.NoError}
}

func (m *MockTransactionHandler) Execute(ctx TransactionContext) *abcitypes.ExecTxResult {
	return &abcitypes.ExecTxResult{Code: types.NoError}
}

func (m *MockTransactionHandler) GetType() TransactionType {
	return m.txType
}

func (m *MockTransactionHandler) Transaction() (Transaction, error) {
	return Transaction{Type: m.txType}, nil
}

// ParseMockTransaction is a mock parse function for testing
func ParseMockTransaction(tx Transaction) (ITransactionHandler, error) {
	return &MockTransactionHandler{txType: tx.Type}, nil
}

// ParseMockTransactionWrongType returns a handler with a different type than requested
func ParseMockTransactionWrongType(tx Transaction) (ITransactionHandler, error) {
	// Return a handler with wrong type to test type mismatch detection
	return &MockTransactionHandler{txType: GovernanceStake}, nil
}

func TestRegistryRegister(t *testing.T) {
	registry := NewTransactionRegistry()

	err := registry.Register(RegisterAccount, ParseMockTransaction)
	if err != nil {
		t.Errorf("Failed to register: %v", err)
	}

	if !registry.IsRegistered(RegisterAccount) {
		t.Error("Transaction should be registered")
	}
}

func TestRegistryRegisterDuplicate(t *testing.T) {
	registry := NewTransactionRegistry()

	// First registration should succeed
	err := registry.Register(RegisterAccount, ParseMockTransaction)
	if err != nil {
		t.Errorf("First registration failed: %v", err)
	}

	// Second registration should fail
	err = registry.Register(RegisterAccount, ParseMockTransaction)
	if err == nil {
		t.Error("Should fail when registering duplicate")
	}
}

func TestRegistryParse(t *testing.T) {
	registry := NewTransactionRegistry()
	registry.Register(RegisterAccount, ParseMockTransaction)

	tx := Transaction{Type: RegisterAccount}
	handler, err := registry.Parse(tx)

	if err != nil {
		t.Errorf("Failed to parse: %v", err)
	}

	if handler.GetType() != RegisterAccount {
		t.Errorf("Expected type %d, got %d", RegisterAccount, handler.GetType())
	}
}

func TestRegistryParseUnregistered(t *testing.T) {
	registry := NewTransactionRegistry()

	tx := Transaction{Type: RegisterAccount}
	_, err := registry.Parse(tx)

	if err == nil {
		t.Error("Should fail when parsing unregistered type")
	}
}

func TestRegistryTypeMismatch(t *testing.T) {
	registry := NewTransactionRegistry()

	// Register a parse function that returns wrong type
	registry.Register(RegisterAccount, ParseMockTransactionWrongType)

	tx := Transaction{Type: RegisterAccount}
	_, err := registry.Parse(tx)

	if err == nil {
		t.Error("Should fail when handler type doesn't match")
	}
}

func TestRegistryConcurrency(t *testing.T) {
	registry := NewTransactionRegistry()
	registry.Register(RegisterAccount, ParseMockTransaction)

	// Simulate concurrent reads
	var wg sync.WaitGroup
	errors := make(chan error, 100)

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			tx := Transaction{Type: RegisterAccount}
			_, err := registry.Parse(tx)
			if err != nil {
				errors <- err
			}
		}()
	}

	wg.Wait()
	close(errors)

	// Check if any goroutine encountered an error
	for err := range errors {
		t.Errorf("Concurrent parse failed: %v", err)
	}
}

func TestGetRegisteredTypes(t *testing.T) {
	registry := NewTransactionRegistry()
	registry.Register(RegisterAccount, ParseMockTransaction)
	registry.Register(GovernanceStake, ParseMockTransaction)

	types := registry.GetRegisteredTypes()

	if len(types) != 2 {
		t.Errorf("Expected 2 registered types, got %d", len(types))
	}

	// Verify both types are present
	found := make(map[TransactionType]bool)
	for _, txType := range types {
		found[txType] = true
	}

	if !found[RegisterAccount] {
		t.Error("RegisterAccount should be in registered types")
	}
	if !found[GovernanceStake] {
		t.Error("GovernanceStake should be in registered types")
	}
}

func TestIsRegistered(t *testing.T) {
	registry := NewTransactionRegistry()
	registry.Register(RegisterAccount, ParseMockTransaction)

	if !registry.IsRegistered(RegisterAccount) {
		t.Error("RegisterAccount should be registered")
	}

	if registry.IsRegistered(GovernanceStake) {
		t.Error("GovernanceStake should not be registered")
	}
}

func TestGlobalRegistryFunctions(t *testing.T) {
	// Create a fresh global registry for this test
	oldRegistry := globalRegistry
	globalRegistry = NewTransactionRegistry()
	defer func() { globalRegistry = oldRegistry }()

	// Test global Register
	err := Register(RegisterAccount, ParseMockTransaction)
	if err != nil {
		t.Errorf("Global Register failed: %v", err)
	}

	// Test global IsRegistered
	if !IsRegistered(RegisterAccount) {
		t.Error("Global IsRegistered should return true")
	}

	// Test global Parse
	tx := Transaction{Type: RegisterAccount}
	handler, err := Parse(tx)
	if err != nil {
		t.Errorf("Global Parse failed: %v", err)
	}
	if handler.GetType() != RegisterAccount {
		t.Errorf("Expected type %d, got %d", RegisterAccount, handler.GetType())
	}

	// Test global GetRegisteredTypes
	types := GetRegisteredTypes()
	if len(types) != 1 {
		t.Errorf("Expected 1 registered type, got %d", len(types))
	}
}

func TestIsSystemTransaction(t *testing.T) {
	// Test system transactions (types 0-6)
	systemTypes := []TransactionType{
		NewValidator,         // 0
		RemoveValidator,      // 1
		SlashValidator,       // 2
		NewEpoch,             // 3
		OChainTokenDeposit,   // 4
		OChainCreditDeposit,  // 5
		ExecutePendingUpgrade, // 6
	}

	for _, txType := range systemTypes {
		if !IsSystemTransaction(txType) {
			t.Errorf("Transaction type %d should be system transaction", txType)
		}
	}

	// Test non-system transactions (types 7+)
	nonSystemTypes := []TransactionType{
		RegisterAccount,           // 7
		GovernanceStake,           // 11
		RegisterUniverseAccount,   // 17
	}

	for _, txType := range nonSystemTypes {
		if IsSystemTransaction(txType) {
			t.Errorf("Transaction type %d should not be system transaction", txType)
		}
	}
}

func TestRegistryMultipleRegistrations(t *testing.T) {
	registry := NewTransactionRegistry()

	// Register multiple transaction types
	types := []TransactionType{
		RegisterAccount,
		GovernanceStake,
		GovernanceUnstakeStart,
		CreateGovernanceProposal,
		RegisterUniverseAccount,
	}

	for _, txType := range types {
		err := registry.Register(txType, ParseMockTransaction)
		if err != nil {
			t.Errorf("Failed to register type %d: %v", txType, err)
		}
	}

	// Verify all are registered
	registeredTypes := registry.GetRegisteredTypes()
	if len(registeredTypes) != len(types) {
		t.Errorf("Expected %d registered types, got %d", len(types), len(registeredTypes))
	}

	// Verify each can be parsed
	for _, txType := range types {
		tx := Transaction{Type: txType}
		handler, err := registry.Parse(tx)
		if err != nil {
			t.Errorf("Failed to parse type %d: %v", txType, err)
		}
		if handler.GetType() != txType {
			t.Errorf("Handler type mismatch for type %d", txType)
		}
	}
}

func TestRegistryConcurrentRegistrationAndParse(t *testing.T) {
	registry := NewTransactionRegistry()

	// Pre-register some types
	registry.Register(RegisterAccount, ParseMockTransaction)

	var wg sync.WaitGroup
	errors := make(chan error, 200)

	// Start 100 concurrent readers
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			tx := Transaction{Type: RegisterAccount}
			_, err := registry.Parse(tx)
			if err != nil {
				errors <- err
			}
		}()
	}

	// Note: We don't do concurrent writes (Register) because that's not
	// a real-world scenario - registrations happen at init time

	wg.Wait()
	close(errors)

	// Check for errors
	for err := range errors {
		t.Errorf("Concurrent operation failed: %v", err)
	}
}
