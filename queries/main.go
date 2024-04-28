package queries

type QueryType string

const (
	GetUniverses QueryType = "GetUniverses"
	GetUniverse  QueryType = "GetUniverse"

	GetAccounts    QueryType = "GetAccounts"
	GetAccount     QueryType = "GetAccount"
	GetLeaderboard QueryType = "GetLeaderboard"
)

type OChainQuery struct {
}

func ResolveQuery(query []byte) ([]byte, error) {
	return []byte(""), nil
}
