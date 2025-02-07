package queries

import (
	"fmt"
	"log"

	"github.com/fxamacker/cbor/v2"
	"github.com/ochain-gg/ochain-network/database"
)

type GetUniverseQueryParameters struct {
	Id string `cbor:"id"`
}

type GetUniverseAccountsQueryParameters struct {
	Address string `cbor:"address"`
}

func ResolveGetUniversesQuery(q []byte, db *database.OChainDatabase) ([]byte, error) {
	log.Println("ResolveGetUniversesQuery ")

	universes, err := db.Universes.GetAll()
	if err != nil {
		log.Println(err)
		return []byte(""), err
	}

	log.Println("Universes founds: " + fmt.Sprint(len(universes)))

	result, err := cbor.Marshal(universes)
	if err != nil {
		return []byte(""), err
	}

	return result, nil
}

func ResolveGetUniverseQuery(q []byte, db *database.OChainDatabase) ([]byte, error) {
	var parameters GetUniverseQueryParameters
	cbor.Unmarshal(q, &parameters)

	universes, err := db.Universes.Get(parameters.Id)
	if err != nil {
		return []byte(""), err
	}

	result, err := cbor.Marshal(universes)
	if err != nil {
		return []byte(""), err
	}

	return result, nil
}
