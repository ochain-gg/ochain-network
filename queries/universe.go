package queries

import (
	"log"

	"github.com/fxamacker/cbor/v2"
	"github.com/ochain-gg/ochain-network/database"
	"github.com/ochain-gg/ochain-network/types"
)

type GetUniverseQueryParameters struct {
	Id string `cbor:"id"`
}

type GetUniverseAccountsQueryParameters struct {
	Address string `cbor:"address"`
}

func ResolveGetUniversesQuery(q []byte, db *database.OChainDatabase) ([]byte, error) {
	universes, err := db.Universes.GetAll()
	if err != nil {
		log.Println(err)
		return []byte(""), err
	}

	var result []types.OChainUniverseWithAttributes
	for index := range universes {
		result = append(result, universes[index].WithAttributes())
	}

	r, err := cbor.Marshal(result)
	if err != nil {
		return []byte(""), err
	}

	return r, nil
}

func ResolveGetUniverseQuery(q []byte, db *database.OChainDatabase) ([]byte, error) {
	log.Println(string(q))

	var parameters GetUniverseQueryParameters
	err := cbor.Unmarshal(q, &parameters)
	if err != nil {
		return []byte(""), err
	}

	universes, err := db.Universes.Get(parameters.Id)
	if err != nil {
		return []byte(""), err
	}

	result, err := cbor.Marshal(universes.WithAttributes())
	if err != nil {
		return []byte(""), err
	}

	return result, nil
}
