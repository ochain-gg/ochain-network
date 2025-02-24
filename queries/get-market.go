package queries

import (
	"encoding/hex"
	"log"

	"github.com/fxamacker/cbor/v2"
	"github.com/ochain-gg/ochain-network/engine/database"
	"github.com/ochain-gg/ochain-network/engine/transactions"
)

type GetMarketQueryParameters struct {
	UniverseId string `cbor:"universeId"`
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

	tx := transactions.Transaction{
		From:      "a",
		Type:      transactions.ChangeAccountIAM,
		Data:      []byte("data"),
		Signature: []byte("sig"),
		Nonce:     1,
	}

	txm, err := cbor.Marshal(tx)
	log.Println("tx representation: ", hex.EncodeToString(txm))

	result, err := cbor.Marshal(market.WithAttributes())
	if err != nil {
		return []byte(""), err
	}

	return result, nil
}
