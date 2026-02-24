package transactions

import (
	"encoding/hex"
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
)

func TestTransactionSigning(t *testing.T) {
	// Create test private key
	privateKey, err := crypto.HexToECDSA("ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80")
	if err != nil {
		t.Fatalf("Failed to create private key: %v", err)
	}
	privateKeyBytes := crypto.FromECDSA(privateKey)

	// Create test transaction
	tx := Transaction{
		Type:  RegisterAccount,
		From:  "0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266",
		Nonce: 1,
		Data:  []byte("test data"),
	}

	// Sign transaction
	err = tx.Sign(privateKeyBytes)
	if err != nil {
		t.Fatalf("Failed to sign transaction: %v", err)
	}

	// Verify signature exists
	if len(tx.Signature) == 0 {
		t.Error("Transaction signature is empty")
	}

	// Verify signer can be recovered
	signer, err := tx.GetSigner()
	if err != nil {
		t.Fatalf("Failed to recover signer: %v", err)
	}

	expectedSigner := "0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266"
	if signer != expectedSigner {
		t.Errorf("Expected signer %s, got %s", expectedSigner, signer)
	}
}

func TestTransactionParsing(t *testing.T) {
	// Create test transaction
	originalTx := Transaction{
		Type:      RegisterAccount,
		From:      "0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266",
		Nonce:     1,
		Data:      []byte("test data"),
		Signature: []byte("test signature"),
	}

	// Serialize transaction
	txBytes, err := originalTx.Bytes()
	if err != nil {
		t.Fatalf("Failed to serialize transaction: %v", err)
	}

	// Parse transaction back
	parsedTx, err := ParseTransaction(txBytes)
	if err != nil {
		t.Fatalf("Failed to parse transaction: %v", err)
	}

	// Verify fields match
	if parsedTx.Type != originalTx.Type {
		t.Errorf("Expected type %d, got %d", originalTx.Type, parsedTx.Type)
	}
	if parsedTx.From != originalTx.From {
		t.Errorf("Expected from %s, got %s", originalTx.From, parsedTx.From)
	}
	if parsedTx.Nonce != originalTx.Nonce {
		t.Errorf("Expected nonce %d, got %d", originalTx.Nonce, parsedTx.Nonce)
	}
}

func TestTransactionUniqueID(t *testing.T) {
	// Create two identical transactions
	tx1 := Transaction{
		Type:  RegisterAccount,
		From:  "0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266",
		Nonce: 1,
		Data:  []byte("test data"),
	}
	tx2 := Transaction{
		Type:  RegisterAccount,
		From:  "0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266",
		Nonce: 1,
		Data:  []byte("test data"),
	}

	// Get unique IDs
	id1, err := tx1.UniqueID()
	if err != nil {
		t.Fatalf("Failed to get unique ID for tx1: %v", err)
	}
	id2, err := tx2.UniqueID()
	if err != nil {
		t.Fatalf("Failed to get unique ID for tx2: %v", err)
	}

	// Verify IDs are identical for identical transactions
	if hex.EncodeToString(id1) != hex.EncodeToString(id2) {
		t.Error("Identical transactions should have identical unique IDs")
	}

	// Create different transaction
	tx3 := Transaction{
		Type:  RegisterAccount,
		From:  "0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266",
		Nonce: 2, // Different nonce
		Data:  []byte("test data"),
	}

	id3, err := tx3.UniqueID()
	if err != nil {
		t.Fatalf("Failed to get unique ID for tx3: %v", err)
	}

	// Verify different transactions have different IDs
	if hex.EncodeToString(id1) == hex.EncodeToString(id3) {
		t.Error("Different transactions should have different unique IDs")
	}
}

func TestInvalidTransactionType(t *testing.T) {
	// Create transaction with invalid type
	tx := Transaction{
		Type:  TransactionType(MaxTransactionType + 1),
		From:  "0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266",
		Nonce: 1,
		Data:  []byte("test data"),
	}

	// Should fail when getting unique ID
	_, err := tx.UniqueID()
	if err == nil {
		t.Error("Expected error for invalid transaction type")
	}
}

func TestTransactionTypedData(t *testing.T) {
	tx := Transaction{
		Type:  RegisterAccount,
		From:  "0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266",
		Nonce: 1,
		Data:  []byte("test data"),
	}

	// Get typed data hash
	typedDataHash, err := tx.GetTypedDataHash()
	if err != nil {
		t.Fatalf("Failed to get typed data hash: %v", err)
	}

	// Verify hash is not empty
	if len(typedDataHash) == 0 {
		t.Error("Typed data hash should not be empty")
	}

	// Verify hash is 32 bytes (Keccak256)
	if len(typedDataHash) != 32 {
		t.Errorf("Expected hash length 32, got %d", len(typedDataHash))
	}
}

func TestTransactionReceipt(t *testing.T) {
	receipt := TransactionReceipt{
		GasCost: 21000,
	}

	// Serialize receipt
	receiptBytes := receipt.Bytes()
	if len(receiptBytes) == 0 {
		t.Error("Receipt bytes should not be empty")
	}
}