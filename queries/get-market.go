package queries

import (
	"github.com/fxamacker/cbor/v2"
	"github.com/ochain-gg/ochain-network/database"
	"github.com/ochain-gg/ochain-network/types"
)

type GetMarketQueryParameters struct {
	UniverseId string                 `cbor:"universeId"`
	From       types.MarketResourceID `cbor:"from"`
	To         types.MarketResourceID `cbor:"to"`
	Amount     uint64                 `cbor:"amount"`
}

func ResolveGetMarketQuery(q []byte, db *database.OChainDatabase) ([]byte, error) {
	var parameters GetMarketQueryParameters
	err := cbor.Unmarshal(q, &parameters)
	if err != nil {
		return []byte(""), err
	}

	universe, err := db.Universes.Get(parameters.UniverseId)
	if err != nil {
		return []byte(""), err
	}

	market, err := db.ResourcesMarket.Get(universe.Id)
	if err != nil {
		return []byte(""), err
	}

	result, err := cbor.Marshal(market.WithAttributes())
	if err != nil {
		return []byte(""), err
	}

	return result, nil
}
