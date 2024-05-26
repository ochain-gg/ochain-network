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
	Address string `cbor:"1,keyasint"`

	UniverseId string `cbor:"2,keyasint"`
	Points     uint64 `cbor:"3,keyasint"`
	CreatedAt  int64  `cbor:"4,keyasint"`
}

func (acc *OChainGlobalAccount) getAllowedSigners() []string {
	var addressList []string

	addressList = append(addressList, acc.Address)
	for i := 0; i < len(acc.IAM.DeleguatedTo); i++ {
		addressList = append(addressList, acc.IAM.DeleguatedTo[i])
	}

	return addressList
}

func (acc *OChainGlobalAccount) IsAllowedSigner(address string, deleguationAuthorized bool) bool {

	if !deleguationAuthorized {
		return address == acc.Address
	} else {
		addressList := acc.getAllowedSigners()
		for i := 0; i < len(addressList); i++ {
			if addressList[i] == address {
				return true
			}
		}

		return false
	}
}
