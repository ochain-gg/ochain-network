package transactions

import (
	"errors"

	abcitypes "github.com/cometbft/cometbft/abci/types"
	"github.com/fxamacker/cbor/v2"

	"github.com/ochain-gg/ochain-network/types"
)

type RegisterAccountTransactionData struct {
	OwnerAddress        string   `cbor:"1,keyasint"`
	GuardianQuorum      uint64   `cbor:"2,keyasint"`
	Guardians           []string `cbor:"3,keyasint"`
	DeleguatedTo        []string `cbor:"4,keyasint"`
	AuthorizerSignature string   `cbor:"5,keyasint"`
}

type RegisterAccountTransaction struct {
	Type      TransactionType
	From      string
	Nonce     uint64
	Data      RegisterAccountTransactionData
	Signature []byte
}

func (tx *RegisterAccountTransaction) Transaction() (Transaction, error) {
	txData, err := cbor.Marshal(tx)
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
