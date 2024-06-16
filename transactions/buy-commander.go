package transactions

import (
	"errors"

	abcitypes "github.com/cometbft/cometbft/abci/types"
	"github.com/fxamacker/cbor/v2"

	"github.com/ochain-gg/ochain-network/types"
)

type BuyCommanderTransactionData struct {
	Universe  string                  `cbor:"2,keyasint"`
	Commander types.OChainCommanderID `cbor:"3,keyasint"`
}

type BuyCommanderTransaction struct {
	Type      TransactionType             `cbor:"1,keyasint"`
	From      string                      `cbor:"2,keyasint"`
	Nonce     uint64                      `cbor:"3,keyasint"`
	Data      BuyCommanderTransactionData `cbor:"4,keyasint"`
	Signature []byte                      `cbor:"5,keyasint"`
}

func (tx *BuyCommanderTransaction) Transaction() (Transaction, error) {
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

func (tx *BuyCommanderTransaction) Check(ctx TransactionContext) error {
	globalAccount, err := ctx.Db.GlobalsAccounts.Get(tx.From)
	if err == nil {
		return errors.New("account doesn't exists")
	}

	_, err = ctx.Db.UniverseAccounts.Get(tx.Data.Universe, tx.From)
	if err == nil {
		return errors.New("universe account doesn't exists")
	}

	//if account balance >= TwoWeekCommanderPrice
	if globalAccount.CreditBalance < types.TwoWeekCommanderPrice {
		return errors.New("unsuficient usd balance")
	}

	return nil
}

func (tx *BuyCommanderTransaction) Execute(ctx TransactionContext) ([]abcitypes.Event, error) {
	err := tx.Check(ctx)
	if err != nil {
		return []abcitypes.Event{}, err
	}

	globalAccount, err := ctx.Db.GlobalsAccounts.Get(tx.From)
	if err == nil {
		return []abcitypes.Event{}, errors.New("account doesn't exists")
	}

	account, err := ctx.Db.UniverseAccounts.Get(tx.Data.Universe, tx.From)
	if err == nil {
		return []abcitypes.Event{}, errors.New("universe account doesn't exists")
	}

	globalAccount.CreditBalance -= types.TwoWeekCommanderPrice
	account.SubscribeToCommander(ctx.Date.Unix(), tx.Data.Commander)

	err = ctx.Db.GlobalsAccounts.Update(globalAccount)
	if err != nil {
		return []abcitypes.Event{}, err
	}
	err = ctx.Db.UniverseAccounts.Update(account)
	if err != nil {
		return []abcitypes.Event{}, err
	}

	events := []abcitypes.Event{
		{
			Type: "CommanderSubscribed",
			Attributes: []abcitypes.EventAttribute{
				{Key: "account", Value: tx.From, Index: true},
				{Key: "universe", Value: tx.Data.Universe, Index: true},
				{Key: "commander", Value: string(tx.Data.Commander), Index: true},
			},
		},
	}

	return events, nil
}

func ParseBuyCommanderTransaction(tx Transaction) (BuyCommanderTransaction, error) {
	var txData BuyCommanderTransactionData
	err := cbor.Unmarshal(tx.Data, &txData)

	if err != nil {
		return BuyCommanderTransaction{}, err
	}

	return BuyCommanderTransaction{
		Type:      tx.Type,
		From:      tx.From,
		Nonce:     tx.Nonce,
		Data:      txData,
		Signature: tx.Signature,
	}, nil
}
