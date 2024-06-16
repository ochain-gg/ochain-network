package types

type OChainReward struct {
	TokenEarned uint64 `cbor:"1,keyasint"`
	USDEarned   uint64 `cbor:"2,keyasint"`
}

type OChainRewardsDistribution struct {
	StackingRewards  float64 `cbor:"4,keyasint"`
	PlayerRewards    float64 `cbor:"4,keyasint"`
	ValidatorRewards float64 `cbor:"4,keyasint"`
}

type OChainEpoch struct {
	Id uint64 `cbor:"1,keyasint"`

	StartedAt int64 `cbor:"2,keyasint"`
	EndedAt   int64 `cbor:"3,keyasint"`

	TokenEarned uint64 `cbor:"4,keyasint"`
	USDEarned   uint64 `cbor:"5,keyasint"`
}

type OChainEpochStackingRewards struct {
	Id         string `cbor:"1,keyasint"`
	UniverseId string `cbor:"2,keyasint"`
	EpochId    uint64 `cbor:"3,keyasint"`
	Account    string `cbor:"4,keyasint"` //0x address which will claim the rewards

	TokenEarned uint64 `cbor:"5,keyasint"`
	USDEarned   uint64 `cbor:"6,keyasint"`
}

type OChainEpochPlayerRewards struct {
	Id         string `cbor:"1,keyasint"`
	UniverseId string `cbor:"2,keyasint"`
	EpochId    uint64 `cbor:"3,keyasint"`
	Rank       uint64 `cbor:"4,keyasint"`
	Account    string `cbor:"6,keyasint"` //0x address which will claim the rewards

	TokenEarned uint64 `cbor:"7,keyasint"`
	USDEarned   uint64 `cbor:"8,keyasint"`
}

type OChainEpochValidatorRewards struct {
	Id        string `cbor:"1,keyasint"`
	EpochId   uint64 `cbor:"2,keyasint"`
	Validator string `cbor:"3,keyasint"` //0x address which will claim the rewards

	TokenEarned uint64 `cbor:"4,keyasint"`
	USDEarned   uint64 `cbor:"5,keyasint"`
}
