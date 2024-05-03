package transactions

import (
	"errors"
)

type ChangeAccountIAMTransactionData struct {
	GuardianQuorum uint64   `json:"guardianQuorum"`
	Guardians      []string `json:"guardians"`
	DeleguatedTo   []string `json:"deleguatedTo"`
}

type ChangeAccountIAMTransaction struct {
	Type      TransactionType `json:"type"`
	From      string          `json:"from"`
	Nonce     uint64          `json:"nonce"`
	Data      string          `json:"data"`
	Signature string          `json:"signature"`
}

func (tx *ChangeAccountIAMTransaction) Check(ctx TransactionContext) (uint64, error) {
	_, err := ctx.Db.GlobalsAccounts.GetByOwnerAddress(string(tx.From))
	if err != nil {
		return 1, errors.New("account don't exists")
	}

	//TODO: verify authrorizer signature
	return 0, nil
}

func (tx *ChangeAccountIAMTransaction) Execute(ctx TransactionContext) (uint64, error) {
	_, err := ctx.Db.GlobalsAccounts.GetByOwnerAddress(tx.From)
	if err != nil {
		return 1, errors.New("account don't exists")
	}

	//TODO: verify authrorizer signature
	return 0, nil
}
