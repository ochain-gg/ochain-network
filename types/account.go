package types

type OChainGlobalAccountIAM struct {
	GuardianQuorum uint64   `cbor:"1,keyasint"`
	Guardians      []string `cbor:"2,keyasint"`
	DeleguatedTo   []string `cbor:"3,keyasint"`
}

type OChainGlobalAccount struct {
	Address        string                 `cbor:"1,keyasint"`
	IAM            OChainGlobalAccountIAM `cbor:"2,keyasint"`
	Nonce          uint64                 `cbor:"3,keyasint"`
	TokenBalance   string                 `cbor:"4,keyasint"`
	USDBalance     string                 `cbor:"5,keyasint"`
	LastDailyClaim int64                  `cbor:"6,keyasint"`
	CreatedAt      int64                  `cbor:"7,keyasint"`
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

type OChainCommanderID string

const (
	EngineerID     OChainCommanderID = "ENGINEER"
	GeologistID    OChainCommanderID = "GEOLOGIST"
	TechnocratID   OChainCommanderID = "TECHNOCRAT"
	FleetAdmiralID OChainCommanderID = "FLEET_ADMIRAL"
)

type OChainCommanderBonus struct {
	CommanderId OChainCommanderID `cbor:"1,keyasint"`
	EndedAt     int64             `cbor:"2,keyasint"`
}

type OChainUniverseAccount struct {
	Address    string `cbor:"1,keyasint"`
	UniverseId string `cbor:"2,keyasint"`

	Points             uint64                    `cbor:"3,keyasint"`
	PlanetsCoordinates []string                  `cbor:"4,keyasint"`
	Technologies       OChainAccountTechnologies `cbor:"5,keyasint"`
	Commanders         []OChainCommanderBonus    `cbor:"6,keyasint"`
	CreatedAt          int64                     `cbor:"7,keyasint"`
}

func (acc *OChainUniverseAccount) HasEngineerCommander(timestamp int64) bool {
	for i := 0; i < len(acc.Commanders); i++ {
		if acc.Commanders[i].CommanderId == EngineerID {
			return acc.Commanders[i].EndedAt > timestamp
		}
	}

	return false
}

func (acc *OChainUniverseAccount) HasGeologistCommander(timestamp int64) bool {
	for i := 0; i < len(acc.Commanders); i++ {
		if acc.Commanders[i].CommanderId == GeologistID {
			return acc.Commanders[i].EndedAt > timestamp
		}
	}

	return false
}

func (acc *OChainUniverseAccount) HasTechnocratCommander(timestamp int64) bool {
	for i := 0; i < len(acc.Commanders); i++ {
		if acc.Commanders[i].CommanderId == TechnocratID {
			return acc.Commanders[i].EndedAt > timestamp
		}
	}

	return false
}

func (acc *OChainUniverseAccount) HasFleetAdmiralCommander(timestamp int64) bool {
	for i := 0; i < len(acc.Commanders); i++ {
		if acc.Commanders[i].CommanderId == FleetAdmiralID {
			return acc.Commanders[i].EndedAt > timestamp
		}
	}

	return false
}
