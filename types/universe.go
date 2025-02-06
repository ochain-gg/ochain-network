package types

import "math"

type OChainResourcesWithAttributes struct {
	OCT       uint64 `cbor:"OCT"`
	Metal     uint64 `cbor:"metal"`
	Crystal   uint64 `cbor:"crystal"`
	Deuterium uint64 `cbor:"deuterium"`
}

type OChainResources struct {
	OCT       uint64 `cbor:"1,keyasint"`
	Metal     uint64 `cbor:"2,keyasint"`
	Crystal   uint64 `cbor:"3,keyasint"`
	Deuterium uint64 `cbor:"4,keyasint"`
}

func (resource *OChainResources) WithAttributes() OChainResourcesWithAttributes {
	return OChainResourcesWithAttributes{
		OCT:       resource.OCT,
		Metal:     resource.Metal,
		Crystal:   resource.Crystal,
		Deuterium: resource.Deuterium,
	}
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

type OChainUniverseWithAttributes struct {
	Id                      string `cbor:"id"`
	Name                    string `cbor:"name"`
	Speed                   uint64 `cbor:"speed"`
	ResourcesMarketEnabled  bool   `cbor:"resourcesMarketEnabled"`
	IsStretchable           bool   `cbor:"isStretchable"` // if true, max galaxy increase in function of colinized planets
	MaxGalaxy               uint64 `cbor:"maxGalaxy"`
	MaxSolarSystemPerGalaxy uint64 `cbor:"maxSolarSystemPerGalaxy"`
	MaxPlanetPerSolarSystem uint64 `cbor:"maxPlanetPerSolarSystem"`
	Accounts                uint64 `cbor:"accounts"`
	ColonizedPlanets        uint64 `cbor:"colonizedPlanets"`
	CreatedAt               uint64 `cbor:"createdAt"`
	EndingAt                uint64 `cbor:"endingAt"`
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

func (universe *OChainUniverse) WithAttributes() OChainUniverseWithAttributes {
	return OChainUniverseWithAttributes{
		Id:                      universe.Id,
		Name:                    universe.Name,
		Speed:                   universe.Speed,
		ResourcesMarketEnabled:  universe.ResourcesMarketEnabled,
		IsStretchable:           universe.IsStretchable,
		MaxGalaxy:               universe.MaxGalaxy,
		MaxSolarSystemPerGalaxy: universe.MaxSolarSystemPerGalaxy,
		MaxPlanetPerSolarSystem: universe.MaxPlanetPerSolarSystem,
		Accounts:                universe.Accounts,
		ColonizedPlanets:        universe.ColonizedPlanets,
		CreatedAt:               universe.CreatedAt,
		EndingAt:                universe.EndingAt,
	}
}

type OChainFighterStatsWithAttributes struct {
	Armor  uint64 `cbor:"armor"`
	Shield uint64 `cbor:"shield"`
	Attack uint64 `cbor:"attack"`
}

type OChainSpaceshipStatsWithAttributes struct {
	Capacity   uint64 `cbor:"capacity"`
	Speed      uint64 `cbor:"speed"`
	Consumtion uint64 `cbor:"consumtion"`
}

type OChainSpaceshipWithAttributes struct {
	Id           OChainSpaceshipID                `cbor:"id"`
	Name         string                           `cbor:"name"`
	FighterStats OChainFighterStats               `cbor:"fighterStats"`
	Stats        OChainSpaceshipStats             `cbor:"stats"`
	Cost         OChainResourcesWithAttributes    `cbor:"cost"`
	Dependencies []OChainDependencyWithAttributes `cbor:"dependencies"`
}

type OChainFighterStats struct {
	Armor  uint64 `cbor:"1,keyasint"`
	Shield uint64 `cbor:"2,keyasint"`
	Attack uint64 `cbor:"3,keyasint"`
}

func (stat *OChainFighterStats) WithAttributes() OChainFighterStatsWithAttributes {
	return OChainFighterStatsWithAttributes{
		Armor:  stat.Armor,
		Shield: stat.Shield,
		Attack: stat.Attack,
	}
}

type OChainSpaceshipStats struct {
	Capacity   uint64 `cbor:"1,keyasint"`
	Speed      uint64 `cbor:"2,keyasint"`
	Consumtion uint64 `cbor:"3,keyasint"`
}

func (stat *OChainSpaceshipStats) WithAttributes() OChainSpaceshipStatsWithAttributes {
	return OChainSpaceshipStatsWithAttributes{
		Capacity:   stat.Capacity,
		Speed:      stat.Speed,
		Consumtion: stat.Consumtion,
	}
}

type OChainSpaceship struct {
	Id           OChainSpaceshipID    `cbor:"1,keyasint"`
	Name         string               `cbor:"2,keyasint"`
	FighterStats OChainFighterStats   `cbor:"3,keyasint"`
	Stats        OChainSpaceshipStats `cbor:"4,keyasint"`
	Cost         OChainResources      `cbor:"5,keyasint"`
	Dependencies []OChainDependency   `cbor:"6,keyasint"`
}

func (spaceship *OChainSpaceship) WithAttributes() OChainSpaceshipWithAttributes {
	var dependencies []OChainDependencyWithAttributes
	for i := range spaceship.Dependencies {
		dependencies = append(dependencies, spaceship.Dependencies[i].WithAttributes())
	}

	return OChainSpaceshipWithAttributes{
		Id:           spaceship.Id,
		Name:         spaceship.Name,
		FighterStats: spaceship.FighterStats,
		Stats:        spaceship.Stats,
		Cost:         spaceship.Cost.WithAttributes(),
		Dependencies: dependencies,
	}
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

type OChainDependencyWithAttributes struct {
	DependencyType OChainDependencyType `cbor:"dependencyType"`
	DependencyId   string               `cbor:"dependencyId"`
	Level          uint64               `cbor:"level"`
}

type OChainDependency struct {
	DependencyType OChainDependencyType `cbor:"1,keyasint"`
	DependencyId   string               `cbor:"2,keyasint"`
	Level          uint64               `cbor:"3,keyasint"`
}

func (dep *OChainDependency) WithAttributes() OChainDependencyWithAttributes {
	return OChainDependencyWithAttributes{
		DependencyType: dep.DependencyType,
		DependencyId:   dep.DependencyId,
		Level:          dep.Level,
	}
}

type OChainDefenseWithAttributes struct {
	Id           OChainDefenseID                  `cbor:"id"`
	Name         string                           `cbor:"name"`
	FighterStats OChainFighterStats               `cbor:"fighterStats"`
	Cost         OChainResourcesWithAttributes    `cbor:"cost"`
	Dependencies []OChainDependencyWithAttributes `cbor:"dependencies"`
}

type OChainDefense struct {
	Id           OChainDefenseID    `cbor:"1,keyasint"`
	Name         string             `cbor:"2,keyasint"`
	FighterStats OChainFighterStats `cbor:"3,keyasint"`
	Cost         OChainResources    `cbor:"4,keyasint"`
	Dependencies []OChainDependency `cbor:"5,keyasint"`
}

func (def *OChainDefense) WithAttributes() OChainDefenseWithAttributes {
	var dependencies []OChainDependencyWithAttributes
	for i := range def.Dependencies {
		dependencies = append(dependencies, def.Dependencies[i].WithAttributes())
	}

	return OChainDefenseWithAttributes{
		Id:           def.Id,
		Name:         def.Name,
		FighterStats: def.FighterStats,
		Cost:         def.Cost.WithAttributes(),
		Dependencies: dependencies,
	}
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

type OChainAccountTechnologiesWithAttributes struct {
	Computer                     uint64 `cbor:"computer"`
	Weapon                       uint64 `cbor:"weapon"`
	Shielding                    uint64 `cbor:"shielding"`
	Armor                        uint64 `cbor:"armor"`
	Energy                       uint64 `cbor:"energy"`
	CombustionDrive              uint64 `cbor:"combustionDrive"`
	ImpulseDrive                 uint64 `cbor:"impulseDrive"`
	HyperspaceDrive              uint64 `cbor:"hyperspaceDrive"`
	Hyperspace                   uint64 `cbor:"hyperspace"`
	Laser                        uint64 `cbor:"laser"`
	Ion                          uint64 `cbor:"ion"`
	Plasma                       uint64 `cbor:"plasma"`
	IntergalacticResearchNetwork uint64 `cbor:"intergalacticResearchNetwork"`
	Astrophysics                 uint64 `cbor:"astrophysics"`
	Graviton                     uint64 `cbor:"graviton"`
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

func (tech *OChainAccountTechnologies) WithAttributes() OChainAccountTechnologiesWithAttributes {
	return OChainAccountTechnologiesWithAttributes{
		Computer:                     tech.Computer,
		Weapon:                       tech.Weapon,
		Shielding:                    tech.Shielding,
		Armor:                        tech.Armor,
		Energy:                       tech.Energy,
		CombustionDrive:              tech.CombustionDrive,
		ImpulseDrive:                 tech.ImpulseDrive,
		HyperspaceDrive:              tech.HyperspaceDrive,
		Hyperspace:                   tech.Hyperspace,
		Laser:                        tech.Laser,
		Ion:                          tech.Ion,
		Plasma:                       tech.Plasma,
		IntergalacticResearchNetwork: tech.IntergalacticResearchNetwork,
		Astrophysics:                 tech.Astrophysics,
		Graviton:                     tech.Graviton,
	}
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
	MetalStorageID        OChainBuildingID = "METAL_STORAGE"
	CrystalStorageID      OChainBuildingID = "CRYSTAL_STORAGE"
	DeuteriumStorageID    OChainBuildingID = "DEUTERIUM_STORAGE"
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

type OChainBuildingWithAttributes struct {
	Id           OChainBuildingID                 `cbor:"id"`
	Name         string                           `cbor:"name"`
	BaseCost     OChainResourcesWithAttributes    `cbor:"baseCost"`
	Dependencies []OChainDependencyWithAttributes `cbor:"dependencies"`
}

type OChainBuilding struct {
	Id           OChainBuildingID   `cbor:"1,keyasint"`
	Name         string             `cbor:"2,keyasint"`
	BaseCost     OChainResources    `cbor:"3,keyasint"`
	Dependencies []OChainDependency `cbor:"4,keyasint"`
}

func (def *OChainBuilding) WithAttributes() OChainBuildingWithAttributes {
	var dependencies []OChainDependencyWithAttributes
	for i := range def.Dependencies {
		dependencies = append(dependencies, def.Dependencies[i].WithAttributes())
	}

	return OChainBuildingWithAttributes{
		Id:           def.Id,
		Name:         def.Name,
		BaseCost:     def.BaseCost.WithAttributes(),
		Dependencies: dependencies,
	}
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

type OChainTechnologyWithAttributes struct {
	Id           OChainTechnologyID               `cbor:"id"`
	Name         string                           `cbor:"name"`
	BaseCost     OChainResourcesWithAttributes    `cbor:"baseCost"`
	Dependencies []OChainDependencyWithAttributes `cbor:"dependencies"`
}

type OChainTechnology struct {
	Id           OChainTechnologyID `cbor:"1,keyasint"`
	Name         string             `cbor:"2,keyasint"`
	BaseCost     OChainResources    `cbor:"3,keyasint"`
	Dependencies []OChainDependency `cbor:"4,keyasint"`
}

func (tech *OChainTechnology) WithAttributes() OChainTechnologyWithAttributes {
	var dependencies []OChainDependencyWithAttributes
	for i := range tech.Dependencies {
		dependencies = append(dependencies, tech.Dependencies[i].WithAttributes())
	}

	return OChainTechnologyWithAttributes{
		Id:           tech.Id,
		Name:         tech.Name,
		BaseCost:     tech.BaseCost.WithAttributes(),
		Dependencies: dependencies,
	}
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
