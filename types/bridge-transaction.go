package types

type OChainTransactionType uint64

const (
	OChainTokenDepositTransaction  OChainTransactionType = 0
	OChainTokenWithdrawTransaction OChainTransactionType = 1
	OChainCreditDepositTransaction OChainTransactionType = 2
)

type OChainTransaction struct {
	Type     OChainTransactionType `cbor:"1,keyasint"`
	Hash     string                `cbor:"2,keyasint"`
	Account  string                `cbor:"3,keyasint"`
	Amount   uint64                `cbor:"4,keyasint"`
	Executed bool                  `cbor:"5,keyasint"`
	Canceled bool                  `cbor:"6,keyasint"`
}
