package queries

import (
	"github.com/fxamacker/cbor/v2"
	"github.com/ochain-gg/ochain-network/database"
	"github.com/ochain-gg/ochain-network/types"
)

type GetGameEntitiesQueryParameters struct {
}

type GetGameEntitiesQueryResponse struct {
	Buildings    []types.OChainBuilding   `cbor:"Buildings"`
	Technologies []types.OChainTechnology `cbor:"technologies"`
	Spaceships   []types.OChainSpaceship  `cbor:"spaceships"`
	Defenses     []types.OChainDefense    `cbor:"defenses"`
}

func ResolveGetGameEntitiesQuery(q []byte, db *database.OChainDatabase) ([]byte, error) {

	buildings, err := db.Buildings.GetAll()
	if err != nil {
		return []byte(""), err
	}

	spaceships, err := db.Spaceships.GetAll()
	if err != nil {
		return []byte(""), err
	}

	defenses, err := db.Defenses.GetAll()
	if err != nil {
		return []byte(""), err
	}

	technologies, err := db.Technologies.GetAll()
	if err != nil {
		return []byte(""), err
	}

	result, err := cbor.Marshal(GetGameEntitiesQueryResponse{
		Buildings:    buildings,
		Spaceships:   spaceships,
		Defenses:     defenses,
		Technologies: technologies,
	})

	if err != nil {
		return []byte(""), err
	}

	return result, nil
}
