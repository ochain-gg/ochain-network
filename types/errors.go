package types

const (
	NoError uint32 = 0

	ParsingTransactionError     uint32 = 1
	ParsingTransactionDataError uint32 = 2
	InvalidTransactionError     uint32 = 3
	InvalidTransactionSignature uint32 = 4
	InvalidNonce                uint32 = 5
	CheckTransactionFailure     uint32 = 6
	ExecuteTransactionFailure   uint32 = 7
	NotImplemented              uint32 = 8
	GasCostHigherThanBalance    uint32 = 9
)

type OChainError struct {
	Code    uint
	Message string
}
