package transactions

import (
	"errors"
	"fmt"

	abcitypes "github.com/cometbft/cometbft/abci/types"
	"github.com/fxamacker/cbor/v2"
	"github.com/ochain-gg/ochain-network/types"
)

type UniverseAccountWithdrawTransactionData struct {
	Universe string `cbor:"1,keyasint"`
	Planet   string `cbor:"2,keyasint"`
	Amount   uint64 `cbor:"3,keyasint"`
}

type UniverseAccountWithdrawTransaction struct {
	Type      TransactionType                        `cbor:"1,keyasint"`
	From      string                                 `cbor:"2,keyasint"`
	Nonce     uint64                                 `cbor:"3,keyasint"`
	Data      UniverseAccountWithdrawTransactionData `cbor:"4,keyasint"`
	Signature []byte                                 `cbor:"5,keyasint"`
}

func (tx *UniverseAccountWithdrawTransaction) Transaction() (Transaction, error) {
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

func (tx *UniverseAccountWithdrawTransaction) Check(ctx TransactionContext) error {
	_, err := ctx.Db.GlobalsAccounts.GetAt(tx.From, uint64(ctx.Date.Unix()))
	if err == nil {
		return errors.New("account doesn't exists")
	}

	_, err = ctx.Db.UniverseAccounts.GetAt(tx.Data.Universe, tx.From, uint64(ctx.Date.Unix()))
	if err == nil {
		return errors.New("universe account doesn't exists")
	}

	planet, err := ctx.Db.Planets.GetAt(tx.Data.Universe, tx.Data.Planet, uint64(ctx.Date.Unix()))
	if err == nil {
		return errors.New("planet doesn't exists")
	}

	year, week := ctx.Date.ISOWeek()
	weeklyUsage, err := ctx.Db.UniverseAccountWeeklyUsage.GetAt(tx.Data.Universe, tx.From, year, week, uint64(ctx.Date.Unix()))
	if err == nil {
		return errors.New("error on retrieve limits")
	}

	if weeklyUsage.WithdrawalsExecuted >= types.MaxWeeklyOCTWithdrawals {
		return errors.New("withdrawal limits exceeded")
	}

	if planet.Resources.OCT < tx.Data.Amount {
		return errors.New("unsuficient planet oct balance")
	}

	return nil
}

func (tx *UniverseAccountWithdrawTransaction) Execute(ctx TransactionContext) ([]abcitypes.Event, error) {
	err := tx.Check(ctx)
	if err != nil {
		return []abcitypes.Event{}, err
	}

	currentDate := uint64(ctx.Date.Unix())

	globalAccount, err := ctx.Db.GlobalsAccounts.GetAt(tx.From, currentDate)
	if err == nil {
		return []abcitypes.Event{}, errors.New("account doesn't exists")
	}

	universe, err := ctx.Db.Universes.GetAt(tx.Data.Universe, currentDate)
	if err == nil {
		return []abcitypes.Event{}, errors.New("universe doesn't exists")
	}

	account, err := ctx.Db.UniverseAccounts.GetAt(universe.Id, tx.From, currentDate)
	if err == nil {
		return []abcitypes.Event{}, errors.New("account doesn't exists")
	}

	planet, err := ctx.Db.Planets.GetAt(tx.Data.Universe, tx.Data.Planet, currentDate)
	if err == nil {
		return []abcitypes.Event{}, errors.New("planet doesn't exists")
	}

	year, week := ctx.Date.ISOWeek()
	weeklyUsage, err := ctx.Db.UniverseAccountWeeklyUsage.GetAt(tx.Data.Universe, tx.From, year, week, uint64(ctx.Date.Unix()))
	if err == nil {
		return []abcitypes.Event{}, errors.New("error on retrieve limits")
	}

	weeklyUsage.WithdrawalsExecuted += 1

	globalAccount.CreditBalance += tx.Data.Amount
	planet.Resources.OCT -= tx.Data.Amount
	planet.UpdateResources(universe.Speed, ctx.Date.Unix(), account)

	err = ctx.Db.GlobalsAccounts.Update(globalAccount)
	if err != nil {
		return []abcitypes.Event{}, err
	}

	err = ctx.Db.Planets.Update(tx.Data.Universe, planet)
	if err != nil {
		return []abcitypes.Event{}, err
	}

	err = ctx.Db.UniverseAccountWeeklyUsage.Upsert(weeklyUsage)
	if err != nil {
		return []abcitypes.Event{}, err
	}

	events := []abcitypes.Event{
		{
			Type: "UniverseAccountWithdraw",
			Attributes: []abcitypes.EventAttribute{
				{Key: "account", Value: tx.From, Index: true},
				{Key: "universe", Value: tx.Data.Universe, Index: true},
				{Key: "planet", Value: tx.Data.Planet, Index: true},
				{Key: "amount", Value: fmt.Sprint(tx.Data.Amount)},
			},
		},
		{
			Type: "PlanetResourcesUpdated",
			Attributes: []abcitypes.EventAttribute{
				{Key: "account", Value: tx.From, Index: true},
				{Key: "universe", Value: tx.Data.Universe, Index: true},
				{Key: "planet", Value: tx.Data.Planet, Index: true},
				{Key: "oct", Value: fmt.Sprint(planet.Resources.OCT)},
				{Key: "metal", Value: fmt.Sprint(planet.Resources.Metal)},
				{Key: "crystal", Value: fmt.Sprint(planet.Resources.Crystal)},
				{Key: "deuterium", Value: fmt.Sprint(planet.Resources.Deuterium)},
			},
		},
	}

	return events, nil
}

func ParseUniverseAccountWithdrawTransaction(tx Transaction) (UniverseAccountWithdrawTransaction, error) {
	var txData UniverseAccountWithdrawTransactionData
	err := cbor.Unmarshal(tx.Data, &txData)

	if err != nil {
		return UniverseAccountWithdrawTransaction{}, err
	}

	return UniverseAccountWithdrawTransaction{
		Type:      tx.Type,
		From:      tx.From,
		Nonce:     tx.Nonce,
		Data:      txData,
		Signature: tx.Signature,
	}, nil
}
