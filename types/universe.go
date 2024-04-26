package types

type OChainResources struct {
	OCT       uint
	Metal     uint
	Crystal   uint
	Deutereum uint
}

type OChainCost struct {
	Duration  uint
	Resources OChainResources
}

type OChainUniverseConfiguration struct {
	//Global configuration
	Speed                   uint
	MaxGalaxy               uint
	MaxSolarSystemPerGalaxy uint
	MaxPlanetPerSolarSystem uint

	Spaceships []OChainSpaceship
	Defenses   []OChainDefense
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

type OChainUniverse struct {
	Id            uint `badgerhold:"key"`
	Name          string
	Speed         uint
	Configuration OChainUniverseConfiguration
	CreatedAt     uint
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
