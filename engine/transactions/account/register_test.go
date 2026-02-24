package account_transactions

import (
	"testing"

	tx "github.com/ochain-gg/ochain-network/engine/transactions"
)

func TestRegisterAccountTransactionGetType(t *testing.T) {
	registerTx := &RegisterAccountTransaction{
		Type: tx.RegisterAccount,
	}

	if registerTx.GetType() != tx.RegisterAccount {
		t.Errorf("Expected type %d, got %d", tx.RegisterAccount, registerTx.GetType())
	}
}

func TestRegisterAccountTransactionRegistered(t *testing.T) {
	if !tx.IsRegistered(tx.RegisterAccount) {
		t.Error("RegisterAccount should be registered")
	}
}

func TestRegisterUniverseAccountTransactionGetType(t *testing.T) {
	registerTx := &RegisterUniverseAccountTransaction{
		Type: tx.RegisterUniverseAccount,
	}

	if registerTx.GetType() != tx.RegisterUniverseAccount {
		t.Errorf("Expected type %d, got %d", tx.RegisterUniverseAccount, registerTx.GetType())
	}
}

func TestRegisterUniverseAccountTransactionRegistered(t *testing.T) {
	if !tx.IsRegistered(tx.RegisterUniverseAccount) {
		t.Error("RegisterUniverseAccount should be registered")
	}
}

func TestUniverseAccountWithdrawTransactionGetType(t *testing.T) {
	withdrawTx := &UniverseAccountWithdrawTransaction{
		Type: tx.UniverseOChainTokenWithdraw,
	}

	if withdrawTx.GetType() != tx.UniverseOChainTokenWithdraw {
		t.Errorf("Expected type %d, got %d", tx.UniverseOChainTokenWithdraw, withdrawTx.GetType())
	}
}

func TestUniverseAccountWithdrawTransactionRegistered(t *testing.T) {
	if !tx.IsRegistered(tx.UniverseOChainTokenWithdraw) {
		t.Error("UniverseOChainTokenWithdraw should be registered")
	}
}

func TestClaimFaucetTransactionGetType(t *testing.T) {
	faucetTx := &ClaimFaucetTransaction{
		Type: tx.ClaimFaucet,
	}

	if faucetTx.GetType() != tx.ClaimFaucet {
		t.Errorf("Expected type %d, got %d", tx.ClaimFaucet, faucetTx.GetType())
	}
}

func TestClaimFaucetTransactionRegistered(t *testing.T) {
	if !tx.IsRegistered(tx.ClaimFaucet) {
		t.Error("ClaimFaucet should be registered")
	}
}

func TestAllAccountTransactionsRegistered(t *testing.T) {
	accountTypes := []tx.TransactionType{
		tx.RegisterAccount,
		tx.RegisterUniverseAccount,
		tx.UniverseOChainTokenWithdraw,
		tx.ClaimFaucet,
	}

	for _, txType := range accountTypes {
		if !tx.IsRegistered(txType) {
			t.Errorf("Account transaction type %d should be registered", txType)
		}
	}
}
