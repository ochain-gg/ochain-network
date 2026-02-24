package transactions

import (
	"fmt"
	"sync"
)

// TransactionRegistry manages the mapping from transaction types to parse functions.
// It provides thread-safe registration and lookup of transaction handlers.
// The registry uses a RWMutex to allow concurrent reads (Parse calls) while
// preventing concurrent writes (Register calls).
type TransactionRegistry struct {
	mu       sync.RWMutex
	handlers map[TransactionType]ParseFunc
}

// Global registry instance used by all transaction types.
// Transaction handlers register themselves in their package's init() function.
var globalRegistry = NewTransactionRegistry()

// NewTransactionRegistry creates a new transaction registry.
// This is primarily used for testing; production code uses the global registry.
func NewTransactionRegistry() *TransactionRegistry {
	return &TransactionRegistry{
		handlers: make(map[TransactionType]ParseFunc),
	}
}

// Register adds a transaction handler parse function to the registry.
// Returns an error if the transaction type is already registered.
// This function is typically called from init() functions in transaction handler packages.
func (r *TransactionRegistry) Register(txType TransactionType, parseFunc ParseFunc) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.handlers[txType]; exists {
		return fmt.Errorf("transaction type %d already registered", txType)
	}

	r.handlers[txType] = parseFunc
	return nil
}

// Parse parses a generic Transaction into a specific handler using the registered parse function.
// Returns an error if the transaction type is not registered or if parsing fails.
// After parsing, verifies that the handler's GetType() matches the transaction type.
func (r *TransactionRegistry) Parse(tx Transaction) (ITransactionHandler, error) {
	r.mu.RLock()
	parseFunc, exists := r.handlers[tx.Type]
	r.mu.RUnlock()

	if !exists {
		return nil, fmt.Errorf("unknown transaction type: %d", tx.Type)
	}

	handler, err := parseFunc(tx)
	if err != nil {
		return nil, fmt.Errorf("failed to parse transaction type %d: %w", tx.Type, err)
	}

	// Verify handler returns correct type (safety check)
	if handler.GetType() != tx.Type {
		return nil, fmt.Errorf("handler type mismatch: expected %d, got %d", tx.Type, handler.GetType())
	}

	return handler, nil
}

// IsRegistered checks if a transaction type is registered in the registry.
// This is useful for distinguishing between "not implemented" and "parse error" cases.
func (r *TransactionRegistry) IsRegistered(txType TransactionType) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	_, exists := r.handlers[txType]
	return exists
}

// GetRegisteredTypes returns a slice of all registered transaction types.
// This is primarily useful for testing and debugging.
// The order of types in the returned slice is not guaranteed.
func (r *TransactionRegistry) GetRegisteredTypes() []TransactionType {
	r.mu.RLock()
	defer r.mu.RUnlock()

	types := make([]TransactionType, 0, len(r.handlers))
	for txType := range r.handlers {
		types = append(types, txType)
	}
	return types
}

// Register adds a transaction handler to the global registry.
// This is a convenience function that calls globalRegistry.Register().
// Returns an error if the transaction type is already registered.
func Register(txType TransactionType, parseFunc ParseFunc) error {
	return globalRegistry.Register(txType, parseFunc)
}

// Parse parses a transaction using the global registry.
// This is a convenience function that calls globalRegistry.Parse().
func Parse(tx Transaction) (ITransactionHandler, error) {
	return globalRegistry.Parse(tx)
}

// IsRegistered checks if a transaction type is registered in the global registry.
// This is a convenience function that calls globalRegistry.IsRegistered().
func IsRegistered(txType TransactionType) bool {
	return globalRegistry.IsRegistered(txType)
}

// GetRegisteredTypes returns all registered types from the global registry.
// This is a convenience function that calls globalRegistry.GetRegisteredTypes().
func GetRegisteredTypes() []TransactionType {
	return globalRegistry.GetRegisteredTypes()
}
