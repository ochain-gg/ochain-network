package types

import (
	"errors"
	"math"
)

type OChainGlobalAccountIAM struct {
	GuardianQuorum uint64   `cbor:"1,keyasint"`
	Guardians      []string `cbor:"2,keyasint"`
	DeleguatedTo   []string `cbor:"3,keyasint"`
}

type OChainGlobalAccount struct {
	Address                 string                 `cbor:"1,keyasint"`
	IAM                     OChainGlobalAccountIAM `cbor:"2,keyasint"`
	Nonce                   uint64                 `cbor:"3,keyasint"`
	TokenBalance            uint64                 `cbor:"4,keyasint"`
	StackedBalance          uint64                 `cbor:"5,keyasint"`
	VotingPowerDeleguated   bool                   `cbor:"5,keyasint"`
	VotingPowerDeleguatedTo string                 `cbor:"5,keyasint"`
	CreditBalance           uint64                 `cbor:"6,keyasint"`
	LastDailyClaim          int64                  `cbor:"7,keyasint"`
	CreatedAt               int64                  `cbor:"8,keyasint"`

	LastTransactionHour  uint64 `cbor:"9,keyasint"`
	LastTransactionCount uint64 `cbor:"10,keyasint"`
}

func (acc *OChainGlobalAccount) GetGasCost(timestamp uint64) uint64 {
	defaultGasCost := 100_000
	currentTransactionHour := timestamp / 60

	if currentTransactionHour > acc.LastTransactionHour {
		return uint64(defaultGasCost)
	}

	lastHourCount := acc.LastTransactionCount + 1
	return uint64(float64(defaultGasCost) * math.Pow(float64(lastHourCount), 3))
}

func (acc *OChainGlobalAccount) ApplyGasCost(timestamp uint64) (uint64, error) {
	currentTransactionHour := timestamp / 60

	if currentTransactionHour > acc.LastTransactionHour {
		acc.LastTransactionHour = currentTransactionHour
		acc.LastTransactionCount = 1
	} else {
		acc.LastTransactionCount += 1
	}

	gasCost := acc.GetGasCost(timestamp)
	if gasCost > acc.TokenBalance {
		return 0, errors.New("account balance don't cover gas cost")
	}

	acc.TokenBalance -= gasCost
	return gasCost, nil
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

	LastOCTWithdrawalAt int64 `cbor:"8,keyasint"`
}

type OChainUniverseAccountWeeklyUsage struct {
	Address    string `cbor:"1,keyasint"`
	UniverseId string `cbor:"2,keyasint"`

	Year int `cbor:"3,keyasint"`
	Week int `cbor:"4,keyasint"`

	WithdrawalsExecuted uint64 `cbor:"5,keyasint"`
	DepositedAmount     uint64 `cbor:"6,keyasint"`
}

func (acc *OChainUniverseAccount) TechnologyLevel(id OChainTechnologyID) uint64 {
	switch id {
	case ComputerID:
		return acc.Technologies.Computer
	case WeaponID:
		return acc.Technologies.Weapon
	case ShieldingID:
		return acc.Technologies.Shielding
	case ArmorID:
		return acc.Technologies.Armor
	case EnergyID:
		return acc.Technologies.Energy
	case CombustionDriveID:
		return acc.Technologies.CombustionDrive
	case ImpulseDriveID:
		return acc.Technologies.ImpulseDrive
	case HyperspaceDriveID:
		return acc.Technologies.HyperspaceDrive
	case HyperspaceID:
		return acc.Technologies.Hyperspace
	case LaserID:
		return acc.Technologies.Laser
	case IonID:
		return acc.Technologies.Ion
	case PlasmaID:
		return acc.Technologies.Plasma
	case IntergalacticResearchNetworkID:
		return acc.Technologies.IntergalacticResearchNetwork
	case AstrophysicsID:
		return acc.Technologies.Astrophysics
	case GravitonID:
		return acc.Technologies.Graviton
	}
	return 0
}

func (acc *OChainUniverseAccount) SetTechnologyLevel(id OChainTechnologyID, level uint64) {
	switch id {
	case ComputerID:
		acc.Technologies.Computer = level
	case WeaponID:
		acc.Technologies.Weapon = level
	case ShieldingID:
		acc.Technologies.Shielding = level
	case ArmorID:
		acc.Technologies.Armor = level
	case EnergyID:
		acc.Technologies.Energy = level
	case CombustionDriveID:
		acc.Technologies.CombustionDrive = level
	case ImpulseDriveID:
		acc.Technologies.ImpulseDrive = level
	case HyperspaceDriveID:
		acc.Technologies.HyperspaceDrive = level
	case HyperspaceID:
		acc.Technologies.Hyperspace = level
	case LaserID:
		acc.Technologies.Laser = level
	case IonID:
		acc.Technologies.Ion = level
	case PlasmaID:
		acc.Technologies.Plasma = level
	case IntergalacticResearchNetworkID:
		acc.Technologies.IntergalacticResearchNetwork = level
	case AstrophysicsID:
		acc.Technologies.Astrophysics = level
	case GravitonID:
		acc.Technologies.Graviton = level
	}
}

func (acc *OChainUniverseAccount) SubscribeToCommander(timestamp int64, commander OChainCommanderID) {
	for i := 0; i < len(acc.Commanders); i++ {
		if acc.Commanders[i].CommanderId == EngineerID {
			if acc.Commanders[i].EndedAt < timestamp {
				acc.Commanders[i].EndedAt = timestamp + 1209600 //2weeks
			}
			return
		} else {
			acc.Commanders[i].EndedAt += 1209600
		}
	}

	acc.Commanders = append(acc.Commanders, OChainCommanderBonus{
		CommanderId: commander,
		EndedAt:     timestamp + 1209600,
	})

	return
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
