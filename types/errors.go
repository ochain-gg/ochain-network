package types

type OChainErrorCode uint64

const (
	NoError OChainErrorCode = 0

	InvalidTransactionNonce     OChainErrorCode = 1
	InvalidTransactionSignature OChainErrorCode = 2
	InvalidTransactionType      OChainErrorCode = 4
)

type OChainError struct {
	Code    uint
	Message string
}
