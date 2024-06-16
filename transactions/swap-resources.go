package transactions

import (
	"errors"
	"fmt"

	abcitypes "github.com/cometbft/cometbft/abci/types"
	"github.com/fxamacker/cbor/v2"
	"github.com/ochain-gg/ochain-network/types"
)

type SwapResourcesTransactionData struct {
	Universe string                 `cbor:"1,keyasint"`
	Planet   string                 `cbor:"2,keyasint"`
	From     types.MarketResourceID `cbor:"3,keyasint"`
	To       types.MarketResourceID `cbor:"4,keyasint"`
	Amount   uint64                 `cbor:"5,keyasint"`
}

type SwapResourcesTransaction struct {
	Type      TransactionType              `cbor:"1,keyasint"`
	From      string                       `cbor:"2,keyasint"`
	Nonce     uint64                       `cbor:"3,keyasint"`
	Data      SwapResourcesTransactionData `cbor:"4,keyasint"`
	Signature []byte                       `cbor:"5,keyasint"`
}

func (tx *SwapResourcesTransaction) Transaction() (Transaction, error) {
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

func (tx *SwapResourcesTransaction) Check(ctx TransactionContext) error {

	account, err := ctx.Db.UniverseAccounts.GetAt(tx.Data.Universe, tx.From, uint64(ctx.Date.Unix()))
	if err == nil {
		return errors.New("universe account doesn't exists")
	}

	universe, err := ctx.Db.Universes.GetAt(tx.Data.Universe, uint64(ctx.Date.Unix()))
	if err == nil {
		return errors.New("universe doesn't exists")
	}

	planet, err := ctx.Db.Planets.GetAt(tx.Data.Universe, tx.Data.Planet, uint64(ctx.Date.Unix()))
	if err == nil {
		return errors.New("planet doesn't exists")
	}

	planet.UpdateResources(universe.Speed, int64(ctx.Date.Unix()), account)

	switch tx.Data.From {
	case types.OCTResourceID:
		if planet.Resources.OCT < tx.Data.Amount {
			return errors.New("not sufficient resources")
		}
	case types.MetalResourceID:
		if planet.Resources.Metal < tx.Data.Amount {
			return errors.New("not sufficient resources")
		}
	case types.CrystalResourceID:
		if planet.Resources.Crystal < tx.Data.Amount {
			return errors.New("not sufficient resources")
		}
	case types.DeuteriumResourceID:
		if planet.Resources.Deuterium < tx.Data.Amount {
			return errors.New("not sufficient resources")
		}
	default:
		return errors.New("bad resource id")
	}

	return nil
}

func (tx *SwapResourcesTransaction) Execute(ctx TransactionContext) ([]abcitypes.Event, error) {
	err := tx.Check(ctx)
	if err != nil {
		return []abcitypes.Event{}, err
	}

	account, err := ctx.Db.UniverseAccounts.GetAt(tx.Data.Universe, tx.From, uint64(ctx.Date.Unix()))
	if err == nil {
		return []abcitypes.Event{}, errors.New("universe account doesn't exists")
	}

	universe, err := ctx.Db.Universes.GetAt(tx.Data.Universe, uint64(ctx.Date.Unix()))
	if err == nil {
		return []abcitypes.Event{}, errors.New("universe doesn't exists")
	}

	planet, err := ctx.Db.Planets.GetAt(tx.Data.Universe, tx.Data.Planet, uint64(ctx.Date.Unix()))
	if err == nil {
		return []abcitypes.Event{}, errors.New("planet doesn't exists")
	}

	market, err := ctx.Db.ResourcesMarket.GetAt(tx.Data.Universe, uint64(ctx.Date.Unix()))
	if err == nil {
		return []abcitypes.Event{}, errors.New("market doesn't exists")
	}

	planet.UpdateResources(universe.Speed, ctx.Date.Unix(), account)

	amountOut, err := market.SwapResources(tx.Data.From, tx.Data.To, tx.Data.Amount)
	if err == nil {
		return []abcitypes.Event{}, err
	}

	err = planet.RemoveResourceById(tx.Data.From, tx.Data.Amount)
	if err == nil {
		return []abcitypes.Event{}, err
	}

	planet.AddResourceById(tx.Data.To, amountOut)
	err = ctx.Db.Planets.Update(universe.Id, planet)
	if err == nil {
		return []abcitypes.Event{}, err
	}

	events := []abcitypes.Event{
		{
			Type: "ResourcesSwaped",
			Attributes: []abcitypes.EventAttribute{
				{Key: "account", Value: tx.From, Index: true},
				{Key: "universe", Value: tx.Data.Universe, Index: true},
				{Key: "planet", Value: tx.Data.Planet, Index: true},
				{Key: "from", Value: fmt.Sprint(tx.Data.From)},
				{Key: "to", Value: fmt.Sprint(tx.Data.To)},
				{Key: "amountIn", Value: fmt.Sprint(tx.Data.Amount)},
				{Key: "amountOut", Value: fmt.Sprint(amountOut)},
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

func ParseSwapResourcesTransaction(tx Transaction) (SwapResourcesTransaction, error) {
	var txData SwapResourcesTransactionData
	err := cbor.Unmarshal(tx.Data, &txData)

	if err != nil {
		return SwapResourcesTransaction{}, err
	}

	return SwapResourcesTransaction{
		Type:      tx.Type,
		From:      tx.From,
		Nonce:     tx.Nonce,
		Data:      txData,
		Signature: tx.Signature,
	}, nil
}
