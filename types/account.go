package types

type OChainGlobalAccountIAM struct {
	GuardianQuorum uint64   `cbor:"1,keyasint"`
	Guardians      []string `cbor:"2,keyasint"`
	DeleguatedTo   []string `cbor:"3,keyasint"`
}

type OChainGlobalAccount struct {
	Address      string                 `cbor:"1,keyasint"`
	IAM          OChainGlobalAccountIAM `cbor:"2,keyasint"`
	Nonce        uint64                 `cbor:"3,keyasint"`
	TokenBalance string                 `cbor:"4,keyasint"`
	CreatedAt    int64                  `cbor:"5,keyasint"`
}

type OChainUniverseAccount struct {
	Address      string `cbor:"1,keyasint"`
	UniverseId   string `cbor:"2,keyasint"`
	OwnerAddress string `cbor:"3,keyasint"`
	Points       uint64 `cbor:"4,keyasint"`
	CreatedAt    int64  `cbor:"5,keyasint"`
}
