package types

import "fmt"

type OChainResources struct {
	OCT       uint `cbor:"1,keyasint"`
	Metal     uint `cbor:"2,keyasint"`
	Crystal   uint `cbor:"3,keyasint"`
	Deutereum uint `cbor:"4,keyasint"`
}

type OChainCost struct {
	Duration  uint
	Resources OChainResources
}

type OChainUniverseConfiguration struct {
	Speed                   uint `cbor:"1,keyasint"`
	MaxGalaxy               uint `cbor:"2,keyasint"`
	MaxSolarSystemPerGalaxy uint `cbor:"3,keyasint"`
	MaxPlanetPerSolarSystem uint `cbor:"4,keyasint"`

	Spaceships []OChainSpaceship `cbor:"5,keyasint"`
	Defenses   []OChainDefense   `cbor:"6,keyasint"`
}

type OChainUniverse struct {
	Id            string                      `cbor:"1,keyasint"`
	Name          string                      `cbor:"2,keyasint"`
	Configuration OChainUniverseConfiguration `cbor:"3,keyasint"`
	Accounts      uint64                      `cbor:"4,keyasint"`
	CreatedAt     uint64                      `cbor:"5,keyasint"`
}

type OChainFighterStats struct {
	Armor  uint `cbor:"1,keyasint"`
	Shield uint `cbor:"2,keyasint"`
	Attack uint `cbor:"3,keyasint"`
}

type OChainSpaceshipStats struct {
	Capacity   uint `cbor:"1,keyasint"`
	Speed      uint `cbor:"2,keyasint"`
	Consumtion uint `cbor:"3,keyasint"`
}

type OChainSpaceship struct {
	Id           uint                 `cbor:"1,keyasint"`
	Name         string               `cbor:"2,keyasint"`
	FighterStats OChainFighterStats   `cbor:"3,keyasint"`
	Stats        OChainSpaceshipStats `cbor:"4,keyasint"`
	Cost         OChainResources      `cbor:"5,keyasint"`
}

type OChainDependencyType uint

const (
	OChainBuildingDependency   OChainDependencyType = 0
	OChainTechnologyDependency OChainDependencyType = 1
)

type OChainDependency struct {
	DependencyType OChainDependencyType `cbor:"1,keyasint"`
	DependencyId   uint                 `cbor:"2,keyasint"`
	Level          uint                 `cbor:"3,keyasint"`
}

type OChainDefense struct {
	Id           uint               `cbor:"1,keyasint"`
	Name         string             `cbor:"2,keyasint"`
	FighterStats OChainFighterStats `cbor:"3,keyasint"`
	Cost         OChainResources    `cbor:"4,keyasint"`
}

type OChainAccountTechnologies struct {
	Computer  uint `cbor:"1,keyasint"`
	Weapon    uint `cbor:"2,keyasint"`
	Shielding uint `cbor:"3,keyasint"`
	Armor     uint `cbor:"4,keyasint"`
	Energy    uint `cbor:"5,keyasint"`

	CombustionDrive uint `cbor:"6,keyasint"`
	ImpulseDrive    uint `cbor:"7,keyasint"`
	HyperspaceDrive uint `cbor:"8,keyasint"`

	Hyperspace uint `cbor:"9,keyasint"`
	Laser      uint `cbor:"10,keyasint"`
	Ion        uint `cbor:"11,keyasint"`
	Plasma     uint `cbor:"12,keyasint"`

	IntergalacticResearchNetwork uint `cbor:"13,keyasint"`
	Astrophysics                 uint `cbor:"14,keyasint"`
	Graviton                     uint `cbor:"15,keyasint"`
}

type OChainPlanetBuildings struct {
	MetalMine       uint `cbor:"1,keyasint"`
	CrystalMine     uint `cbor:"2,keyasint"`
	DeutereumMine   uint `cbor:"3,keyasint"`
	SolarPowerPlant uint `cbor:"4,keyasint"`

	RoboticFactory   uint `cbor:"5,keyasint"`
	NaniteFactory    uint `cbor:"6,keyasint"`
	SpaceshipFactory uint `cbor:"7,keyasint"`

	IntergalacticPortal uint `cbor:"8,keyasint"`
	ResearchLaboratory  uint `cbor:"9,keyasint"`
	ShieldDome          uint `cbor:"10,keyasint"`
}

type OChainFleetSpaceships struct {
	SmallCargo   uint `cbor:"1,keyasint"`
	LargeCargo   uint `cbor:"2,keyasint"`
	LightFighter uint `cbor:"3,keyasint"`
	HeavyFighter uint `cbor:"4,keyasint"`

	Cruiser       uint `cbor:"5,keyasint"`
	Battleship    uint `cbor:"6,keyasint"`
	Battlecruiser uint `cbor:"7,keyasint"`

	Bomber    uint `cbor:"8,keyasint"`
	Destroyer uint `cbor:"9,keyasint"`
	Deathstar uint `cbor:"10,keyasint"`
	Reaper    uint `cbor:"11,keyasint"`
	Recycler  uint `cbor:"12,keyasint"`
}

type OChainPlanetDefences struct {
	RocketLauncher  uint `cbor:"1,keyasint"`
	LightLaser      uint `cbor:"2,keyasint"`
	HeavyLaser      uint `cbor:"3,keyasint"`
	IonCannon       uint `cbor:"4,keyasint"`
	GaussCannon     uint `cbor:"5,keyasint"`
	PlasmaTurret    uint `cbor:"6,keyasint"`
	DarkMatterCanon uint `cbor:"7,keyasint"`
}

type OChainPlanet struct {
	Owner       uint64 `cbor:"1,keyasint"`
	Universe    uint64 `cbor:"2,keyasint"`
	Galaxy      uint64 `cbor:"3,keyasint"`
	SolarSystem uint64 `cbor:"4,keyasint"`
	Planet      uint64 `cbor:"5,keyasint"`

	Buildings  OChainPlanetBuildings `cbor:"6,keyasint"`
	Spaceships OChainFleetSpaceships `cbor:"7,keyasint"`
	Defenses   OChainPlanetDefences  `cbor:"8,keyasint"`
	Resources  OChainResources       `cbor:"9,keyasint"`

	LastResourceUpdate uint `cbor:"10,keyasint"`
}

func (planet *OChainPlanet) CoordinateId() string {
	return fmt.Sprint(planet.Galaxy) + "_" + fmt.Sprint(planet.SolarSystem) + fmt.Sprint(planet.Planet)
}
