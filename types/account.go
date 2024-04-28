package types

type OChainGlobalAccountIAM struct {
	OwnerAddress string

	GuardianQuorum uint64
	Guardians      []string

	DeleguatedTo []string
}

type OChainGlobalAccount struct {
	Id           uint `badgerhold:"key"`
	IAM          OChainGlobalAccountIAM
	Nonce        uint64
	TokenBalance string

	CreatedAt int64
}

type OChainUniverseAccountIAM struct {
}

type OChainUniverseAccount struct {
	Id           uint   `badgerhold:"key"`
	UniverseId   uint64 `badgerhold:"index"`
	OwnerAddress string
	Points       uint64
	IAM          OChainUniverseAccountIAM

	CreatedAt int64
}
