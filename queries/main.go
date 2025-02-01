package queries

import (
	"errors"

	abcitypes "github.com/cometbft/cometbft/abci/types"
	"github.com/ochain-gg/ochain-network/database"
)

const (
	GetUniversesPath = "ochain_getUniverses"
	GetUniversePath  = "ochain_getUniverse"

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

	case GetAccountPath:
		value, err := ResolveGetAccountQuery(req.Data, db)
		if err != nil {
			return []byte{}, err
		}
		return value, nil
	}

	return []byte{}, errors.New("unknown query path")
}
