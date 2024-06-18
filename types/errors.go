package types

const (
	NoError uint32 = 0

	ParsingTransactionError     uint32 = 1
	ParsingTransactionDataError uint32 = 2
	InvalidTransactionError     uint32 = 3
	InvalidTransactionSignature uint32 = 4
	CheckTransactionFailure     uint32 = 5
	ExecuteTransactionFailure   uint32 = 6
	NotImplemented              uint32 = 7
	GasCostHigherThanBalance    uint32 = 8
)

type OChainError struct {
	Code    uint
	Message string
}
