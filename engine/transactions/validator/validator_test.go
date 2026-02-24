package validator_transactions

import (
	"testing"

	tx "github.com/ochain-gg/ochain-network/engine/transactions"
)

func TestNewValidatorTransactionGetType(t *testing.T) {
	newValidatorTx := &OChainNewValidatorTransaction{
		Type: tx.NewValidator,
	}

	if newValidatorTx.GetType() != tx.NewValidator {
		t.Errorf("Expected type %d, got %d", tx.NewValidator, newValidatorTx.GetType())
	}
}

func TestNewValidatorTransactionRegistered(t *testing.T) {
	if !tx.IsRegistered(tx.NewValidator) {
		t.Error("NewValidator should be registered")
	}
}

func TestRemoveValidatorTransactionGetType(t *testing.T) {
	removeValidatorTx := &OChainRemoveValidatorTransaction{
		Type: tx.RemoveValidator,
	}

	if removeValidatorTx.GetType() != tx.RemoveValidator {
		t.Errorf("Expected type %d, got %d", tx.RemoveValidator, removeValidatorTx.GetType())
	}
}

func TestRemoveValidatorTransactionRegistered(t *testing.T) {
	if !tx.IsRegistered(tx.RemoveValidator) {
		t.Error("RemoveValidator should be registered")
	}
}

func TestOChainTokenDepositTransactionGetType(t *testing.T) {
	tokenDepositTx := &OChainTokenDepositTransaction{
		Type: tx.OChainTokenDeposit,
	}

	if tokenDepositTx.GetType() != tx.OChainTokenDeposit {
		t.Errorf("Expected type %d, got %d", tx.OChainTokenDeposit, tokenDepositTx.GetType())
	}
}

func TestOChainTokenDepositTransactionRegistered(t *testing.T) {
	if !tx.IsRegistered(tx.OChainTokenDeposit) {
		t.Error("OChainTokenDeposit should be registered")
	}
}

func TestOChainCreditDepositTransactionGetType(t *testing.T) {
	creditDepositTx := &OChainCreditDepositTransaction{
		Type: tx.OChainCreditDeposit,
	}

	if creditDepositTx.GetType() != tx.OChainCreditDeposit {
		t.Errorf("Expected type %d, got %d", tx.OChainCreditDeposit, creditDepositTx.GetType())
	}
}

func TestOChainCreditDepositTransactionRegistered(t *testing.T) {
	if !tx.IsRegistered(tx.OChainCreditDeposit) {
		t.Error("OChainCreditDeposit should be registered")
	}
}

func TestAllValidatorTransactionsRegistered(t *testing.T) {
	validatorTypes := []tx.TransactionType{
		tx.NewValidator,
		tx.RemoveValidator,
		tx.OChainTokenDeposit,
		tx.OChainCreditDeposit,
	}

	for _, txType := range validatorTypes {
		if !tx.IsRegistered(txType) {
			t.Errorf("Validator transaction type %d should be registered", txType)
		}
	}
}
