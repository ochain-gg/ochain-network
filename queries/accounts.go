package queries

import (
	"github.com/fxamacker/cbor/v2"
	"github.com/ochain-gg/ochain-network/database"
	"github.com/ochain-gg/ochain-network/types"
)

type GetAccountQueryParameters struct {
	Address string `cbor:"address"`
}

type GetAccountQueryResponse struct {
	Account          types.OChainGlobalAccount     `cbor:"account"`
	UniverseAccounts []types.OChainUniverseAccount `cbor:"universeAccounts"`
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
		return []byte(""), err
	}

	result.Account = relatedGlobalAccount

	universeAccounts, err := db.UniverseAccounts.GetByAddress(result.Account.Address)
	if err != nil {
		return []byte(""), err
	}

	for i := 0; i < len(universeAccounts); i++ {
		result.UniverseAccounts = append(result.UniverseAccounts, universeAccounts[i])
	}

	response, err := cbor.Marshal(result)
	if err != nil {
		return []byte(""), err
	}

	return response, nil
}
