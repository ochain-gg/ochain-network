package queries

import (
	"github.com/fxamacker/cbor/v2"
	"github.com/ochain-gg/ochain-network/database"
	"github.com/ochain-gg/ochain-network/types"
)

type GetAccountsQueryParameters struct {
	Wallet string `cbor:"wallet"`
}

type GetAccountsQueryResponseEntity struct {
	Account          types.OChainGlobalAccount     `cbor:"account"`
	Deleguated       bool                          `cbor:"deleguated"`
	UniverseAccounts []types.OChainUniverseAccount `cbor:"universeAccounts"`
}

func ResolveGetAccountsQuery(q []byte, db *database.OChainDatabase) ([]byte, error) {

	var queryParams GetAccountsQueryParameters
	err := cbor.Unmarshal(q, &queryParams)
	if err != nil {
		return []byte(""), err
	}

	var results []GetAccountsQueryResponseEntity

	relatedGlobalAccounts, err := db.GlobalsAccounts.GetRelatedToAddress(queryParams.Wallet)
	if err != nil {
		return []byte(""), err
	}

	for i := 0; i < len(relatedGlobalAccounts); i++ {

		accounts, err := db.UniverseAccounts.GetByAddress(relatedGlobalAccounts[i].Address)
		if err != nil {
			return []byte(""), err
		}

		result := GetAccountsQueryResponseEntity{
			Account:          relatedGlobalAccounts[i],
			Deleguated:       relatedGlobalAccounts[i].Address == queryParams.Wallet,
			UniverseAccounts: accounts,
		}

		results = append(results, result)
	}

	response, err := cbor.Marshal(results)
	if err != nil {
		return []byte(""), err
	}

	return response, nil
}
