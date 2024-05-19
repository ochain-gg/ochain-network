package transactions

import (
	"errors"

	abcitypes "github.com/cometbft/cometbft/abci/types"
	"github.com/fxamacker/cbor/v2"

	"github.com/ochain-gg/ochain-network/types"
)

type RegisterAccountTransactionData struct {
	Address             string   `cbor:"1,keyasint"`
	GuardianQuorum      uint64   `cbor:"2,keyasint"`
	Guardians           []string `cbor:"3,keyasint"`
	DeleguatedTo        []string `cbor:"4,keyasint"`
	AuthorizerSignature string   `cbor:"5,keyasint"`
}

type RegisterAccountTransaction struct {
	Type      TransactionType                `cbor:"1,keyasint"`
	From      string                         `cbor:"2,keyasint"`
	Nonce     uint64                         `cbor:"3,keyasint"`
	Data      RegisterAccountTransactionData `cbor:"4,keyasint"`
	Signature []byte                         `cbor:"5,keyasint"`
}

func (tx *RegisterAccountTransaction) Transaction() (Transaction, error) {
	txData, err := cbor.Marshal(tx.Data)
	if err != nil {
		return Transaction{}, err
	}

	return Transaction{
		Type:      tx.Type,
		From:      tx.From,
		Nonce:     tx.Nonce,
		Data:      txData,
		Signature: tx.Signature,
	}, nil
}

func (tx *RegisterAccountTransaction) Check(ctx TransactionContext) error {
	_, err := ctx.Db.GlobalsAccounts.Get(tx.Data.Address)
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
		Address:      tx.From,
		Nonce:        1,
		TokenBalance: "0",
		IAM: types.OChainGlobalAccountIAM{
			GuardianQuorum: tx.Data.GuardianQuorum,
			Guardians:      tx.Data.Guardians,
			DeleguatedTo:   tx.Data.DeleguatedTo,
		},
	}

	err = ctx.Db.GlobalsAccounts.Insert(account)
	if err != nil {
		return []abcitypes.Event{}, err
	}

	events := []abcitypes.Event{
		{
			Type: "GlobalAccountRegistered",
			Attributes: []abcitypes.EventAttribute{
				{Key: "address", Value: tx.From, Index: true},
			},
		},
	}

	return events, nil
}

func ParseRegisterAccountTransaction(tx Transaction) (RegisterAccountTransaction, error) {
	var txData RegisterAccountTransactionData
	err := cbor.Unmarshal(tx.Data, &txData)

	if err != nil {
		return RegisterAccountTransaction{}, err
	}

	return RegisterAccountTransaction{
		Type:      tx.Type,
		From:      tx.From,
		Nonce:     tx.Nonce,
		Data:      txData,
		Signature: tx.Signature,
	}, nil
}
