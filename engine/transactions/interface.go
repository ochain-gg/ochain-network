package transactions

import (
	abcitypes "github.com/cometbft/cometbft/abci/types"
)

// ITransactionHandler defines the interface that all transaction handlers must implement.
// This interface is implemented by both authenticated and system transactions.
// For system transactions, the base Transaction struct will have empty From, zero Nonce,
// and nil Signature fields, but the interface remains the same.
type ITransactionHandler interface {
	// Check validates the transaction without modifying state.
	// Returns a CheckTxResponse with appropriate error code if validation fails.
	Check(ctx TransactionContext) *abcitypes.CheckTxResponse

	// Execute performs the transaction and modifies state.
	// Should call Check() first for defensive validation.
	// Returns an ExecTxResult with events and receipt data.
	Execute(ctx TransactionContext) *abcitypes.ExecTxResult

	// GetType returns the transaction type constant.
	// Used by the registry to verify handlers match their registered type.
	GetType() TransactionType

	// Transaction converts the typed transaction back to the generic Transaction format.
	// Used for serialization and network transmission.
	Transaction() (Transaction, error)
}

// ParseFunc is a function that parses a generic Transaction into a typed transaction handler.
// Each transaction type must provide a ParseFunc that deserializes the transaction data
// and returns an implementation of ITransactionHandler.
type ParseFunc func(tx Transaction) (ITransactionHandler, error)

// IsSystemTransaction returns true if the transaction type is a system transaction (types 0-6).
// System transactions are unauthenticated and don't have From, Nonce, or Signature fields.
func IsSystemTransaction(txType TransactionType) bool {
	return uint64(txType) <= 6
}
