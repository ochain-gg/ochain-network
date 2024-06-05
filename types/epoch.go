package types

type OChainEpoch struct {
	Id uint64 `cbor:"1,keyasint"`

	StartedAt int64 `cbor:"2,keyasint"`
	EndedAt   int64 `cbor:"3,keyasint"`

	TotalTokenEarned uint64 `cbor:"4,keyasint"`
	TotalUSDEarned   uint64 `cbor:"5,keyasint"`
}

type OChainEpochRewards struct {
	Id        string `cbor:"1,keyasint"`
	EpochId   uint64 `cbor:"2,keyasint"`
	Validator string `cbor:"3,keyasint"` //0x address which will claim the rewards

	TokenEarned uint64 `cbor:"4,keyasint"`
	USDEarned   uint64 `cbor:"5,keyasint"`
}
