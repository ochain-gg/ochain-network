package types

type OChainTokenTransferType uint64

const (
	UniverseDepositTransfer  OChainTokenTransferType = 0
	UniverseWithdrawTransfer OChainTokenTransferType = 1
)

type OChainTokenTransfer struct {
	Id              string                  `cbor:"1,keyasint"`
	Type            OChainTokenTransferType `cbor:"2,keyasint"`
	TransactionHash string                  `cbor:"3,keyasint"`
	Account         string                  `cbor:"4,keyasint"`
	Amount          uint64                  `cbor:"5,keyasint"`
	ExecutedAt      bool                    `cbor:"6,keyasint"`
}
