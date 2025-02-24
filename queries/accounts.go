package queries

import (
	"errors"

	"github.com/dgraph-io/badger/v4"
	"github.com/fxamacker/cbor/v2"
	"github.com/ochain-gg/ochain-network/engine/database"
	"github.com/ochain-gg/ochain-network/types"
)

type GetAccountQueryParameters struct {
	Address string `cbor:"address"`
}

type GetAccountQueryResponse struct {
	Account          types.OChainGlobalAccountWithAttributes     `cbor:"account"`
	UniverseAccounts []types.OChainUniverseAccountWithAttributes `cbor:"universeAccounts"`
}

func ResolveGetAccountQuery(q []byte, db *database.OChainDatabase) ([]byte, error) {

	var queryParams GetAccountQueryParameters
	err := cbor.Unmarshal(q, &queryParams)
	if err != nil {
		return []byte(""), err
	}

	var result GetAccountQueryResponse

	relatedGlobalAccount, err := db.GlobalsAccounts.Get(queryParams.Address)
	if err != nil {
		if err == badger.ErrKeyNotFound {
			return []byte(""), errors.New("global account doesn't exists")
		}
	}

	result.Account = relatedGlobalAccount.WithAttributes()

	universeAccounts, err := db.UniverseAccounts.GetByAddress(result.Account.Address)
	if err != nil {
		return []byte(""), err
	}

	for i := 0; i < len(universeAccounts); i++ {
		result.UniverseAccounts = append(result.UniverseAccounts, universeAccounts[i].WithAttribute())
	}

	response, err := cbor.Marshal(result)
	if err != nil {
		return []byte(""), err
	}

	return response, nil
}
