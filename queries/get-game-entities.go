package queries

import (
	"github.com/fxamacker/cbor/v2"
	"github.com/ochain-gg/ochain-network/engine/database"
	"github.com/ochain-gg/ochain-network/types"
)

type GetGameEntitiesQueryParameters struct {
}

type GetGameEntitiesQueryResponse struct {
	Buildings    []types.OChainBuildingWithAttributes   `cbor:"buildings"`
	Technologies []types.OChainTechnologyWithAttributes `cbor:"technologies"`
	Spaceships   []types.OChainSpaceshipWithAttributes  `cbor:"spaceships"`
	Defenses     []types.OChainDefenseWithAttributes    `cbor:"defenses"`
}

func ResolveGetGameEntitiesQuery(q []byte, db *database.OChainDatabase) ([]byte, error) {

	buildings, err := db.Buildings.GetAll()
	if err != nil {
		return []byte(""), err
	}

	var buildingsWithAttributes []types.OChainBuildingWithAttributes
	for i := range buildings {
		buildingsWithAttributes = append(buildingsWithAttributes, buildings[i].WithAttributes())
	}

	spaceships, err := db.Spaceships.GetAll()
	if err != nil {
		return []byte(""), err
	}

	var spaceshipsWithAttributes []types.OChainSpaceshipWithAttributes
	for i := range spaceships {
		spaceshipsWithAttributes = append(spaceshipsWithAttributes, spaceships[i].WithAttributes())
	}

	defenses, err := db.Defenses.GetAll()
	if err != nil {
		return []byte(""), err
	}

	var defensesWithAttributes []types.OChainDefenseWithAttributes
	for i := range defenses {
		defensesWithAttributes = append(defensesWithAttributes, defenses[i].WithAttributes())
	}

	technologies, err := db.Technologies.GetAll()
	if err != nil {
		return []byte(""), err
	}

	var technologiesWithAttributes []types.OChainTechnologyWithAttributes
	for i := range technologies {
		technologiesWithAttributes = append(technologiesWithAttributes, technologies[i].WithAttributes())
	}

	result, err := cbor.Marshal(GetGameEntitiesQueryResponse{
		Buildings:    buildingsWithAttributes,
		Spaceships:   spaceshipsWithAttributes,
		Defenses:     defensesWithAttributes,
		Technologies: technologiesWithAttributes,
	})

	if err != nil {
		return []byte(""), err
	}

	return result, nil
}
