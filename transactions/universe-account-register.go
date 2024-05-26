package transactions

import (
	"errors"

	abcitypes "github.com/cometbft/cometbft/abci/types"
	"github.com/fxamacker/cbor/v2"

	"github.com/ochain-gg/ochain-network/types"
)

type RegisterUniverseAccountTransactionData struct {
	UniverseId string `cbor:"1,keyasint"`
}

type RegisterUniverseAccountTransaction struct {
	Type      TransactionType                        `cbor:"1,keyasint"`
	From      string                                 `cbor:"2,keyasint"`
	Nonce     uint64                                 `cbor:"3,keyasint"`
	Data      RegisterUniverseAccountTransactionData `cbor:"4,keyasint"`
	Signature []byte                                 `cbor:"5,keyasint"`
}

func (tx *RegisterUniverseAccountTransaction) Transaction() (Transaction, error) {
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

func (tx *RegisterUniverseAccountTransaction) Check(ctx TransactionContext) error {
	_, err := ctx.Db.GlobalsAccounts.Get(tx.From)
	if err == nil {
		return errors.New("account doesn't exists")
	}

	universeExists, err := ctx.Db.Universes.Exists(tx.Data.UniverseId)
	if err != nil {
		return errors.New("error on universe account fetch")
	}
	if !universeExists {
		return errors.New("universe doesn't exists")
	}

	exists, err := ctx.Db.UniverseAccounts.Exists(tx.Data.UniverseId, tx.From)
	if err != nil {
		return errors.New("error on universe account fetch")
	}

	if exists {
		return errors.New("universe account already exists")
	}

	return nil
}

func (tx *RegisterUniverseAccountTransaction) Execute(ctx TransactionContext) ([]abcitypes.Event, error) {
	err := tx.Check(ctx)
	if err != nil {
		return []abcitypes.Event{}, err
	}

	universe, err := ctx.Db.Universes.Get(tx.Data.UniverseId)
	if err != nil {
		return []abcitypes.Event{}, errors.New("error on universe account fetch")
	}

	account := types.OChainUniverseAccount{
		Address:    tx.From,
		UniverseId: tx.Data.UniverseId,
		Points:     0,
		CreatedAt:  ctx.Date.Unix(),
	}

	err = ctx.Db.UniverseAccounts.Insert(account)
	if err != nil {
		return []abcitypes.Event{}, err
	}

	universe.Accounts += 1
	err = ctx.Db.Universes.Update(universe)
	if err != nil {
		return []abcitypes.Event{}, err
	}

	events := []abcitypes.Event{
		{
			Type: "UniverseAccountRegistered",
			Attributes: []abcitypes.EventAttribute{
				{Key: "address", Value: tx.From, Index: true},
				{Key: "universeId", Value: tx.Data.UniverseId, Index: true},
			},
		},
	}

	return events, nil
}

func ParseRegisterUniverseAccountTransaction(tx Transaction) (RegisterUniverseAccountTransaction, error) {
	var txData RegisterUniverseAccountTransactionData
	err := cbor.Unmarshal(tx.Data, &txData)

	if err != nil {
		return RegisterUniverseAccountTransaction{}, err
	}

	return RegisterUniverseAccountTransaction{
		Type:      tx.Type,
		From:      tx.From,
		Nonce:     tx.Nonce,
		Data:      txData,
		Signature: tx.Signature,
	}, nil
}
