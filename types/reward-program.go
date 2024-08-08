package types

type OChainRewardProgramType uint64

const (
	OChainValidatorRewardProgramType OChainRewardProgramType = 0
	OChainUniverseRewardProgramType  OChainRewardProgramType = 1
)

type OChainRewardProgram struct {
	Id                   uint64                  `cbor:"1,keyasint"`
	Type                 OChainRewardProgramType `cbor:"2,keyasint"`
	UniverseId           uint64                  `cbor:"3,keyasint"`
	TopPlayerRewarded    uint64                  `cbor:"4,keyasint"`
	DistributionPerEpoch uint64                  `cbor:"5,keyasint"`
	StartingAtEpoch      uint64                  `cbor:"6,keyasint"`
	EndingAtEpoch        uint64                  `cbor:"7,keyasint"`
}
