package types

type OChainPortalTransactionType uint64

const (
	OChainBridgeDepositTransaction   OChainPortalTransactionType = 0
	OChainBridgeWithdrawTransaction  OChainPortalTransactionType = 1
	OChainBridgeSubscribeTransaction OChainPortalTransactionType = 2
)

type OChainBridgeTransaction struct {
	Type     OChainPortalTransactionType `cbor:"1,keyasint"`
	Hash     string                      `cbor:"2,keyasint"`
	Account  string                      `cbor:"3,keyasint"`
	Amount   uint64                      `cbor:"4,keyasint"`
	Executed bool                        `cbor:"5,keyasint"`
	Canceled bool                        `cbor:"6,keyasint"`
}
