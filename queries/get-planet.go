package queries

import (
	"time"

	"github.com/fxamacker/cbor/v2"
	"github.com/ochain-gg/ochain-network/database"
	"github.com/ochain-gg/ochain-network/types"
)

type GetPlanetQueryParameters struct {
	UniverseId   string `cbor:"universeId"`
	CoordinateId string `cbor:"coordinateId"`
}

type PlanetStatistics struct {
	TotalEnergy uint64 `cbor:"universeId"`
}

type GetPlanetQueryResponse struct {
	Planet           types.OChainPlanetWithAttributes    `cbor:"planet"`
	Stats            types.OChainPlanetStatistics        `cbor:"stats"`
	BuildingUpgrades []types.OChainUpgradeWithAttributes `cbor:"pendingBuildingUpgrades"`
}

func ResolveGetPlanetQuery(q []byte, db *database.OChainDatabase) ([]byte, error) {
	var parameters GetPlanetQueryParameters
	err := cbor.Unmarshal(q, &parameters)
	if err != nil {
		return []byte(""), err
	}

	universe, err := db.Universes.Get(parameters.UniverseId)
	if err != nil {
		return []byte(""), err
	}

	planet, err := db.Planets.Get(parameters.UniverseId, parameters.CoordinateId)
	if err != nil {
		return []byte(""), err
	}

	pendingUpgrades, err := db.Upgrades.GetPendingBuildingUpgradesByPlanet(parameters.UniverseId, parameters.CoordinateId)
	if err != nil {
		return []byte(""), err
	}

	var upgrades []types.OChainUpgradeWithAttributes
	for i := range pendingUpgrades {
		upgrades = append(upgrades, types.OChainUpgradeWithAttributes(pendingUpgrades[i].WithAttributes()))
	}

	account, err := db.UniverseAccounts.Get(parameters.UniverseId, planet.Owner)
	if err != nil {
		return []byte(""), err
	}

	planet.UpdateResources(universe.Speed, time.Now().Unix(), account)

	result, err := cbor.Marshal(GetPlanetQueryResponse{
		Planet:           planet.WithAttributes(),
		Stats:            planet.PlanetStatistics(universe.Speed, time.Now().Unix(), account),
		BuildingUpgrades: upgrades,
	})

	if err != nil {
		return []byte(""), err
	}

	return result, nil
}
