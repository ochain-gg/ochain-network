package transactions

import (
	"encoding/json"
	"errors"

	abcitypes "github.com/cometbft/cometbft/abci/types"

	"github.com/ochain.gg/ochain-network-validator/types"
)

type RegisterAccountTransactionData struct {
	OwnerAddress string `json:"owner"`

	GuardianQuorum uint64   `json:"guardianQuorum"`
	Guardians      []string `json:"guardians"`
	DeleguatedTo   []string `json:"deleguatedTo"`

	AuthorizerSignature string `json:"authorizerSignature"`
}

type RegisterAccountTransaction struct {
	Hash      string                         `json:"hash"`
	Type      TransactionType                `json:"type"`
	From      string                         `json:"from"`
	Nonce     uint64                         `json:"nonce"`
	Data      RegisterAccountTransactionData `json:"data"`
	Signature string                         `json:"signature"`
}

func (tx *RegisterAccountTransaction) Check(ctx TransactionContext) error {
	_, err := ctx.Db.GlobalsAccounts.GetByOwnerAddress(tx.Data.OwnerAddress)
	if err == nil {
		return errors.New("account aleady exists")
	}
	return nil
}

func (tx *RegisterAccountTransaction) Execute(ctx TransactionContext) ([]abcitypes.Event, error) {
	err := tx.Check(ctx)
	if err != nil {
		return []abcitypes.Event{}, err
	}

	account := types.OChainGlobalAccount{
		Nonce:        1,
		TokenBalance: "0",
		IAM: types.OChainGlobalAccountIAM{
			OwnerAddress: tx.From,

			GuardianQuorum: tx.Data.GuardianQuorum,
			Guardians:      tx.Data.Guardians,
			DeleguatedTo:   tx.Data.DeleguatedTo,
		},
	}

	ctx.Db.GlobalsAccounts.Insert(account, ctx.Txn)

	events := []abcitypes.Event{
		{
			Type: "GlobalAccountRegistered",
			Attributes: []abcitypes.EventAttribute{
				{Key: "adress", Value: tx.From, Index: true},
			},
		},
	}

	return events, nil
}

func ParseRegisterAccountTransaction(tx AuthenticatedTransaction) (RegisterAccountTransaction, error) {
	var txData RegisterAccountTransactionData
	err := json.Unmarshal([]byte(tx.Data), &txData)

	if err != nil {
		return RegisterAccountTransaction{}, err
	}

	return RegisterAccountTransaction{
		Hash:      tx.Hash,
		Type:      tx.Type,
		From:      tx.From,
		Nonce:     tx.Nonce,
		Data:      txData,
		Signature: tx.Signature,
	}, nil
}

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
	if err == nil {
		return 1, errors.New("account aleady exists")
	}

	//TODO: verify authrorizer signature
	return 0, nil
}
