package types

type OChainBridgeTransactionType uint64

const (
	OChainBridgeTokenDepositTransaction  OChainBridgeTransactionType = 0
	OChainBridgeTokenWithdrawTransaction OChainBridgeTransactionType = 1
	OChainBridgeCreditDepositTransaction OChainBridgeTransactionType = 2
)

type OChainBridgeTransaction struct {
	Type     OChainBridgeTransactionType `cbor:"1,keyasint"`
	Hash     string                      `cbor:"2,keyasint"`
	Account  string                      `cbor:"3,keyasint"`
	Amount   uint64                      `cbor:"4,keyasint"`
	Executed bool                        `cbor:"5,keyasint"`
	Canceled bool                        `cbor:"6,keyasint"`
}
