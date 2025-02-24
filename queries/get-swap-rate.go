package queries

import (
	"github.com/fxamacker/cbor/v2"
	"github.com/ochain-gg/ochain-network/engine/database"
	"github.com/ochain-gg/ochain-network/types"
)

type GetSwapRateQueryParameters struct {
	UniverseId string                 `cbor:"universeId"`
	From       types.MarketResourceID `cbor:"from"`
	To         types.MarketResourceID `cbor:"to"`
	Amount     uint64                 `cbor:"amount"`
}

type GetSwapRateQueryResponse struct {
	UniverseId string                 `cbor:"universeId"`
	From       types.MarketResourceID `cbor:"from"`
	To         types.MarketResourceID `cbor:"to"`
	Amount     uint64                 `cbor:"amount"`
	AmountOut  uint64                 `cbor:"amountOut"`
}

func ResolveGetSwapRateQuery(q []byte, db *database.OChainDatabase) ([]byte, error) {
	var parameters GetSwapRateQueryParameters
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

	amountOut, err := market.GetSwapAmountOut(parameters.From, parameters.To, parameters.Amount)

	result, err := cbor.Marshal(GetSwapRateQueryResponse{
		UniverseId: parameters.UniverseId,
		From:       parameters.From,
		To:         parameters.To,
		Amount:     parameters.Amount,
		AmountOut:  amountOut,
	})
	if err != nil {
		return []byte(""), err
	}

	return result, nil
}
