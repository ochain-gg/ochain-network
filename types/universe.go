package types

import "math"

type OChainResources struct {
	OCT       uint64 `cbor:"1,keyasint"`
	Metal     uint64 `cbor:"2,keyasint"`
	Crystal   uint64 `cbor:"3,keyasint"`
	Deuterium uint64 `cbor:"4,keyasint"`
}

func (resource *OChainResources) Add(r OChainResources) {
	resource.OCT += r.OCT
	resource.Metal += r.Metal
	resource.Crystal += r.Crystal
	resource.Deuterium += r.Deuterium
}

func (resource *OChainResources) Sub(r OChainResources) {
	resource.OCT -= r.OCT
	resource.Metal -= r.Metal
	resource.Crystal -= r.Crystal
	resource.Deuterium -= r.Deuterium
}

func (resource *OChainResources) Mul(factor uint64) {
	resource.OCT *= factor
	resource.Metal *= factor
	resource.Crystal *= factor
	resource.Deuterium *= factor
}

type OChainCost struct {
	Duration  uint64
	Resources OChainResources
}

type OChainUniverse struct {
	Id                     string `cbor:"1,keyasint"`
	Name                   string `cbor:"2,keyasint"`
	Speed                  uint64 `cbor:"3,keyasint"`
	ResourcesMarketEnabled bool   `cbor:"4,keyasint"`
	IsStretchable          bool   `cbor:"5,keyasint"` // if true, max galaxy increase in function of colinized planets

	MaxGalaxy               uint64 `cbor:"6,keyasint"`
	MaxSolarSystemPerGalaxy uint64 `cbor:"7,keyasint"`
	MaxPlanetPerSolarSystem uint64 `cbor:"8,keyasint"`

	Accounts         uint64 `cbor:"9,keyasint"`
	ColonizedPlanets uint64 `cbor:"10,keyasint"`

	CreatedAt uint64 `cbor:"11,keyasint"`
	EndingAt  uint64 `cbor:"12,keyasint"`
}

type OChainFighterStats struct {
	Armor  uint64 `cbor:"1,keyasint"`
	Shield uint64 `cbor:"2,keyasint"`
	Attack uint64 `cbor:"3,keyasint"`
}

type OChainSpaceshipStats struct {
	Capacity   uint64 `cbor:"1,keyasint"`
	Speed      uint64 `cbor:"2,keyasint"`
	Consumtion uint64 `cbor:"3,keyasint"`
}

type OChainSpaceship struct {
	Id           OChainSpaceshipID    `cbor:"1,keyasint"`
	Name         string               `cbor:"2,keyasint"`
	FighterStats OChainFighterStats   `cbor:"3,keyasint"`
	Stats        OChainSpaceshipStats `cbor:"4,keyasint"`
	Cost         OChainResources      `cbor:"5,keyasint"`
	Dependencies []OChainDependency   `cbor:"6,keyasint"`
}

func (spaceship *OChainSpaceship) MeetRequirements(planet OChainPlanet, acc OChainUniverseAccount) bool {

	for i := 0; i < len(spaceship.Dependencies); i++ {
		dep := spaceship.Dependencies[i]
		if dep.DependencyType == OChainBuildingDependency {
			if planet.BuildingLevel(OChainBuildingID(dep.DependencyId)) < dep.Level {
				return false
			}
		}

		if dep.DependencyType == OChainTechnologyDependency {
			if acc.TechnologyLevel(OChainTechnologyID(dep.DependencyId)) < dep.Level {
				return false
			}
		}
	}

	return true
}

type OChainDependencyType uint64

const (
	OChainBuildingDependency   OChainDependencyType = 0
	OChainTechnologyDependency OChainDependencyType = 1
)

type OChainDependency struct {
	DependencyType OChainDependencyType `cbor:"1,keyasint"`
	DependencyId   string               `cbor:"2,keyasint"`
	Level          uint64               `cbor:"3,keyasint"`
}

type OChainDefense struct {
	Id           OChainDefenseID    `cbor:"1,keyasint"`
	Name         string             `cbor:"2,keyasint"`
	FighterStats OChainFighterStats `cbor:"3,keyasint"`
	Cost         OChainResources    `cbor:"4,keyasint"`
	Dependencies []OChainDependency `cbor:"5,keyasint"`
}

func (defense *OChainDefense) MeetRequirements(planet OChainPlanet, acc OChainUniverseAccount) bool {

	for i := 0; i < len(defense.Dependencies); i++ {
		dep := defense.Dependencies[i]
		if dep.DependencyType == OChainBuildingDependency {
			if planet.BuildingLevel(OChainBuildingID(dep.DependencyId)) < dep.Level {
				return false
			}
		}

		if dep.DependencyType == OChainTechnologyDependency {
			if acc.TechnologyLevel(OChainTechnologyID(dep.DependencyId)) < dep.Level {
				return false
			}
		}
	}

	return true
}

type OChainAccountTechnologies struct {
	Computer  uint64 `cbor:"1,keyasint"`
	Weapon    uint64 `cbor:"2,keyasint"`
	Shielding uint64 `cbor:"3,keyasint"`
	Armor     uint64 `cbor:"4,keyasint"`
	Energy    uint64 `cbor:"5,keyasint"`

	CombustionDrive uint64 `cbor:"6,keyasint"`
	ImpulseDrive    uint64 `cbor:"7,keyasint"`
	HyperspaceDrive uint64 `cbor:"8,keyasint"`

	Hyperspace uint64 `cbor:"9,keyasint"`
	Laser      uint64 `cbor:"10,keyasint"`
	Ion        uint64 `cbor:"11,keyasint"`
	Plasma     uint64 `cbor:"12,keyasint"`

	IntergalacticResearchNetwork uint64 `cbor:"13,keyasint"`
	Astrophysics                 uint64 `cbor:"14,keyasint"`
	Graviton                     uint64 `cbor:"15,keyasint"`
}

type OChainBuildingID string

const (
	MetalMineID           OChainBuildingID = "METAL_MINE"
	CrystalMineID         OChainBuildingID = "CRYSTAL_MINE"
	DeuteriumMineID       OChainBuildingID = "DEUTERIUM_MINE"
	SolarPowerPlantID     OChainBuildingID = "SOLAR_POWER_PLANT"
	RoboticFactoryID      OChainBuildingID = "ROBOTIC_FACTORY"
	NaniteFactoryID       OChainBuildingID = "NANITE_FACTORY"
	SpaceshipFactoryID    OChainBuildingID = "SPACESHIP_FACTORY"
	IntergalacticPortalID OChainBuildingID = "INTERGALACTIC_PORTAL"
	ResearchLaboratoryID  OChainBuildingID = "RESEARCH_LABORATORY"
	ShieldDomeID          OChainBuildingID = "SHIELD_DOME"
)

type OChainTechnologyID string

const (
	ComputerID                     OChainTechnologyID = "COMPUTER"
	WeaponID                       OChainTechnologyID = "WEAPON"
	ShieldingID                    OChainTechnologyID = "SHIELDING"
	ArmorID                        OChainTechnologyID = "ARMOR"
	EnergyID                       OChainTechnologyID = "ENERGY"
	CombustionDriveID              OChainTechnologyID = "COMBUSTION_DRIVE"
	ImpulseDriveID                 OChainTechnologyID = "IMPULSE_DRIVE"
	HyperspaceDriveID              OChainTechnologyID = "HYPERSPACE_DRIVE"
	HyperspaceID                   OChainTechnologyID = "HYPERSPACE"
	LaserID                        OChainTechnologyID = "LASER"
	IonID                          OChainTechnologyID = "ION"
	PlasmaID                       OChainTechnologyID = "PLASMA"
	IntergalacticResearchNetworkID OChainTechnologyID = "INTERGALACTIC_RESEARCH_NETWORK"
	AstrophysicsID                 OChainTechnologyID = "ASTROPHYSICS"
	GravitonID                     OChainTechnologyID = "GRAVITON"
)

type OChainSpaceshipID string

const (
	SmallCargoID    OChainSpaceshipID = "SMALL_CARGO"
	LargeCargoID    OChainSpaceshipID = "LARGE_CARGO"
	LightFighterID  OChainSpaceshipID = "LIGHT_FIGHTER"
	HeavyFighterID  OChainSpaceshipID = "HEAVY_FIGHTER"
	CruiserID       OChainSpaceshipID = "CRUISER"
	BattleshipID    OChainSpaceshipID = "BATTLESHIP"
	BattlecruiserID OChainSpaceshipID = "BATTLECRUISER"
	BomberID        OChainSpaceshipID = "BOMBER"
	DestroyerID     OChainSpaceshipID = "DESTROYER"
	DeathstarID     OChainSpaceshipID = "DEATHSTAR"
	ReaperID        OChainSpaceshipID = "REAPER"
	RecyclerID      OChainSpaceshipID = "RECYCLER"
)

type OChainDefenseID string

const (
	RocketLauncherID  OChainDefenseID = "ROCKET_LAUNCHER"
	LightLaserID      OChainDefenseID = "LIGHT_LASER"
	HeavyLaserID      OChainDefenseID = "HEAVY_LASER"
	IonCannonID       OChainDefenseID = "ION_CANNON"
	GaussCannonID     OChainDefenseID = "GAUSS_CANNON"
	PlasmaTurretID    OChainDefenseID = "PLASMA_TURRET"
	DarkMatterCanonID OChainDefenseID = "DARK_MATTER_CANON"
)

type OChainBuilding struct {
	Id           OChainBuildingID   `cbor:"1,keyasint"`
	Name         string             `cbor:"2,keyasint"`
	BaseCost     OChainResources    `cbor:"3,keyasint"`
	Dependencies []OChainDependency `cbor:"4,keyasint"`
}

func (building *OChainBuilding) GetUpgradeCost(level uint64) OChainResources {
	return OChainResources{
		OCT:       0,
		Metal:     uint64(float64(building.BaseCost.Metal) * math.Pow(1.75, float64(level-1))),
		Crystal:   uint64(float64(building.BaseCost.Crystal) * math.Pow(1.75, float64(level-1))),
		Deuterium: uint64(float64(building.BaseCost.Deuterium) * math.Pow(1.75, float64(level-1))),
	}
}

func (building *OChainBuilding) MeetRequirements(planet OChainPlanet, acc OChainUniverseAccount) bool {

	for i := 0; i < len(building.Dependencies); i++ {
		dep := building.Dependencies[i]
		if dep.DependencyType == OChainBuildingDependency {
			if planet.BuildingLevel(OChainBuildingID(dep.DependencyId)) < dep.Level {
				return false
			}
		}

		if dep.DependencyType == OChainTechnologyDependency {
			if acc.TechnologyLevel(OChainTechnologyID(dep.DependencyId)) < dep.Level {
				return false
			}
		}
	}

	return true
}

type OChainTechnology struct {
	Id           OChainTechnologyID `cbor:"1,keyasint"`
	Name         string             `cbor:"2,keyasint"`
	BaseCost     OChainResources    `cbor:"3,keyasint"`
	Dependencies []OChainDependency `cbor:"4,keyasint"`
}

func (technology *OChainTechnology) GetUpgradeCost(level uint64) OChainResources {
	return OChainResources{
		OCT:       0,
		Metal:     uint64(float64(technology.BaseCost.Metal) * math.Pow(2, float64(level-1))),
		Crystal:   uint64(float64(technology.BaseCost.Crystal) * math.Pow(2, float64(level-1))),
		Deuterium: uint64(float64(technology.BaseCost.Deuterium) * math.Pow(2, float64(level-1))),
	}
}

func (technology *OChainTechnology) MeetRequirements(planet OChainPlanet, acc OChainUniverseAccount) bool {

	for i := 0; i < len(technology.Dependencies); i++ {
		dep := technology.Dependencies[i]
		if dep.DependencyType == OChainBuildingDependency {
			if planet.BuildingLevel(OChainBuildingID(dep.DependencyId)) < dep.Level {
				return false
			}
		}

		if dep.DependencyType == OChainTechnologyDependency {
			if acc.TechnologyLevel(OChainTechnologyID(dep.DependencyId)) < dep.Level {
				return false
			}
		}
	}

	return true
}
