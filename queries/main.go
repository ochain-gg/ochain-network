package queries

import (
	"errors"

	abcitypes "github.com/cometbft/cometbft/abci/types"
	"github.com/ochain-gg/ochain-network/database"
)

const (
	GetUniversesPath    = "ochain_getUniverses"
	GetUniversePath     = "ochain_getUniverse"
	GetPlanetPath       = "ochain_getPlanet"
	GetGameEntitiesPath = "ochain_getGameEntities"

	GetAccountPath     = "ochain_getAccount"
	GetLeaderboardPath = "ochain_getLeaderboard"
)

func GetQueryResponse(req *abcitypes.QueryRequest, db *database.OChainDatabase) ([]byte, error) {
	switch req.Path {
	case GetUniversesPath:
		value, err := ResolveGetUniversesQuery(req.Data, db)
		if err != nil {
			return []byte{}, err
		}
		return value, nil

	case GetUniversePath:
		value, err := ResolveGetUniverseQuery(req.Data, db)
		if err != nil {
			return []byte{}, err
		}
		return value, nil

	case GetAccountPath:
		value, err := ResolveGetAccountQuery(req.Data, db)
		if err != nil {
			return []byte{}, err
		}
		return value, nil

	case GetPlanetPath:
		value, err := ResolveGetPlanetQuery(req.Data, db)
		if err != nil {
			return []byte{}, err
		}
		return value, nil

	case GetGameEntitiesPath:
		value, err := ResolveGetGameEntitiesQuery(req.Data, db)
		if err != nil {
			return []byte{}, err
		}
		return value, nil
	}

	return []byte{}, errors.New("unknown query path")
}
