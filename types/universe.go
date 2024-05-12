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
	Id string `cbor:"id"`

	Speed                   uint `cbor:"speed"`
	MaxGalaxy               uint `cbor:"maxGalaxy"`
	MaxSolarSystemPerGalaxy uint `cbor:"maxSolarSystemPerGalaxy"`
	MaxPlanetPerSolarSystem uint `cbor:"maxPlanetPerSolarSystem"`

	Spaceships []OChainSpaceship `cbor:"spaceships"`
	Defenses   []OChainDefense   `cbor:"defenses"`
}

type OChainUniverse struct {
	Id        string `cbor:"id"`
	Name      string `cbor:"name"`
	CreatedAt uint   `cbor:"createdAt"`
}

type OChainFighterStats struct {
	Armor  uint
	Shield uint
	Attack uint
}

type OChainSpaceshipStats struct {
	Capacity   uint
	Speed      uint
	Consumtion uint
}

type OChainSpaceship struct {
	Id           uint `badgerhold:"key"`
	Name         string
	FighterStats OChainFighterStats
	Stats        OChainSpaceshipStats
	Cost         OChainResources
}

type OChainDependencyType uint

const (
	OChainBuildingDependency   OChainDependencyType = 0
	OChainTechnologyDependency OChainDependencyType = 1
)

type OChainDependency struct {
	DependencyType OChainDependencyType
	DependencyId   uint
	Level          uint
}

type OChainDefense struct {
	Id           uint `badgerhold:"key"`
	Name         string
	FighterStats OChainFighterStats
	Cost         OChainResources
}

type OChainAccountTechnologies struct {
	Computer  uint
	Weapon    uint
	Shielding uint
	Armor     uint
	Energy    uint

	CombustionDrive uint
	ImpulseDrive    uint
	HyperspaceDrive uint

	Hyperspace uint
	Laser      uint
	Ion        uint
	Plasma     uint

	IntergalacticResearchNetwork uint
	Astrophysics                 uint
	Graviton                     uint
}

type OChainPlanetBuildings struct {
	MetalMine       uint
	CrystalMine     uint
	DeutereumMine   uint
	SolarPowerPlant uint

	RoboticFactory   uint
	NaniteFactory    uint
	SpaceshipFactory uint

	IntergalacticPortal uint
	ResearchLaboratory  uint
	ShieldDome          uint
}

type OChainFleetSpaceships struct {
	SmallCargo   uint
	LargeCargo   uint
	LightFighter uint
	HeavyFighter uint

	Cruiser       uint
	Battleship    uint
	Battlecruiser uint

	Bomber    uint
	Destroyer uint
	Deathstar uint
	Reaper    uint
	Recycler  uint
}

type OChainPlanetDefences struct {
	RocketLauncher  uint
	LightLaser      uint
	HeavyLaser      uint
	IonCannon       uint
	GaussCannon     uint
	PlasmaTurret    uint
	DarkMatterCanon uint
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
