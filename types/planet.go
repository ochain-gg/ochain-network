package types

import (
	"errors"
	"fmt"
	"math"
)

type OChainUpgradeType uint64

const (
	OChainBuildingUpgrade   OChainUpgradeType = 0
	OChainTechnologyUpgrade OChainUpgradeType = 1
)

type OChainBuildType uint64

const (
	OChainSpaceshipBuild OChainBuildType = 0
	OChainDefenseBuild   OChainBuildType = 1
)

type OChainUpgrade struct {
	UniverseId         string            `cbor:"1,keyasint"`
	PlanetCoordinateId string            `cbor:"2,keyasint"`
	UpgradeType        OChainUpgradeType `cbor:"3,keyasint"`
	UpgradeId          string            `cbor:"4,keyasint"`
	Level              uint64            `cbor:"5,keyasint"`
	StartedAt          int64             `cbor:"6,keyasint"`
	EndedAt            int64             `cbor:"7,keyasint"`
	Executed           bool              `cbor:"8,keyasint"`
}

type OChainBuild struct {
	UniverseId         string          `cbor:"1,keyasint"`
	PlanetCoordinateId string          `cbor:"2,keyasint"`
	BuildType          OChainBuildType `cbor:"3,keyasint"`
	BuildId            string          `cbor:"4,keyasint"`
	StartedAt          int64           `cbor:"5,keyasint"`
	EndedAt            int64           `cbor:"6,keyasint"`
}

type OChainPlanetBuildings struct {
	MetalMine       uint64 `cbor:"1,keyasint"`
	CrystalMine     uint64 `cbor:"2,keyasint"`
	DeuteriumMine   uint64 `cbor:"3,keyasint"`
	SolarPowerPlant uint64 `cbor:"4,keyasint"`

	RoboticFactory   uint64 `cbor:"5,keyasint"`
	NaniteFactory    uint64 `cbor:"6,keyasint"`
	SpaceshipFactory uint64 `cbor:"7,keyasint"`

	IntergalacticPortal uint64 `cbor:"8,keyasint"`
	ResearchLaboratory  uint64 `cbor:"9,keyasint"`
	ShieldDome          uint64 `cbor:"10,keyasint"`
}

type OChainFleetSpaceships struct {
	SmallCargo   uint64 `cbor:"1,keyasint"`
	LargeCargo   uint64 `cbor:"2,keyasint"`
	LightFighter uint64 `cbor:"3,keyasint"`
	HeavyFighter uint64 `cbor:"4,keyasint"`

	Cruiser       uint64 `cbor:"5,keyasint"`
	Battleship    uint64 `cbor:"6,keyasint"`
	Battlecruiser uint64 `cbor:"7,keyasint"`

	Bomber    uint64 `cbor:"8,keyasint"`
	Destroyer uint64 `cbor:"9,keyasint"`
	Deathstar uint64 `cbor:"10,keyasint"`
	Reaper    uint64 `cbor:"11,keyasint"`
	Recycler  uint64 `cbor:"12,keyasint"`
}

type OChainPlanetDefences struct {
	RocketLauncher  uint64 `cbor:"1,keyasint"`
	LightLaser      uint64 `cbor:"2,keyasint"`
	HeavyLaser      uint64 `cbor:"3,keyasint"`
	IonCannon       uint64 `cbor:"4,keyasint"`
	GaussCannon     uint64 `cbor:"5,keyasint"`
	PlasmaTurret    uint64 `cbor:"6,keyasint"`
	DarkMatterCanon uint64 `cbor:"7,keyasint"`
}

type OChainBuild struct {
	BuildType OChainBuildType `cbor:"1,keyasint"`
	BuildId   string          `cbor:"2,keyasint"`
	Count     uint64          `cbor:"3,keyasint"`
}

type OChainBuildQueueItem struct {
	BuildType OChainBuildType `cbor:"1,keyasint"`
	BuildId   string          `cbor:"2,keyasint"`
	Count     uint64          `cbor:"3,keyasint"`
	StartedAt uint64          `cbor:"4,keyasint"`
	FinishAt  uint64          `cbor:"5,keyasint"`
}

type OChainPlanet struct {
	Owner       string `cbor:"1,keyasint"`
	Universe    string `cbor:"2,keyasint"`
	Galaxy      uint64 `cbor:"3,keyasint"`
	SolarSystem uint64 `cbor:"4,keyasint"`
	Planet      uint64 `cbor:"5,keyasint"`

	Buildings  OChainPlanetBuildings `cbor:"6,keyasint"`
	Spaceships OChainFleetSpaceships `cbor:"7,keyasint"`
	Defenses   OChainPlanetDefences  `cbor:"8,keyasint"`
	Resources  OChainResources       `cbor:"9,keyasint"`

	BuildQueue []OChainBuildQueueItem `cbor:"10,keyasint"`

	LastResourceUpdate int64 `cbor:"11,keyasint"`
}

func (planet *OChainPlanet) CoordinateId() string {
	return fmt.Sprint(planet.Galaxy) + "_" + fmt.Sprint(planet.SolarSystem) + "_" + fmt.Sprint(planet.Planet)
}

func (planet *OChainPlanet) BuildingLevel(id OChainBuildingID) uint64 {
	switch id {
	case MetalMineID:
		return planet.Buildings.MetalMine
	case CrystalMineID:
		return planet.Buildings.CrystalMine
	case DeuteriumMineID:
		return planet.Buildings.DeuteriumMine
	case SolarPowerPlantID:
		return planet.Buildings.SolarPowerPlant
	case RoboticFactoryID:
		return planet.Buildings.RoboticFactory
	case NaniteFactoryID:
		return planet.Buildings.NaniteFactory
	case SpaceshipFactoryID:
		return planet.Buildings.SpaceshipFactory
	case IntergalacticPortalID:
		return planet.Buildings.IntergalacticPortal
	case ResearchLaboratoryID:
		return planet.Buildings.ResearchLaboratory
	case ShieldDomeID:
		return planet.Buildings.ShieldDome
	}
	return 0
}

func (planet *OChainPlanet) SetBuildingLevel(id OChainBuildingID, level uint64) {
	switch id {
	case MetalMineID:
		planet.Buildings.MetalMine = level
	case CrystalMineID:
		planet.Buildings.CrystalMine = level
	case DeuteriumMineID:
		planet.Buildings.DeuteriumMine = level
	case SolarPowerPlantID:
		planet.Buildings.SolarPowerPlant = level
	case RoboticFactoryID:
		planet.Buildings.RoboticFactory = level
	case NaniteFactoryID:
		planet.Buildings.NaniteFactory = level
	case SpaceshipFactoryID:
		planet.Buildings.SpaceshipFactory = level
	case IntergalacticPortalID:
		planet.Buildings.IntergalacticPortal = level
	case ResearchLaboratoryID:
		planet.Buildings.ResearchLaboratory = level
	case ShieldDomeID:
		planet.Buildings.ShieldDome = level
	}

}

func (planet *OChainPlanet) AddDefenses(id OChainDefenseID, count uint64) {
	switch id {
	case RocketLauncherID:
		planet.Defenses.RocketLauncher += count
	case LightLaserID:
		planet.Defenses.LightLaser += count
	case HeavyLaserID:
		planet.Defenses.HeavyLaser += count
	case IonCannonID:
		planet.Defenses.IonCannon += count
	case GaussCannonID:
		planet.Defenses.GaussCannon += count
	case PlasmaTurretID:
		planet.Defenses.PlasmaTurret += count
	case DarkMatterCanonID:
		planet.Defenses.DarkMatterCanon += count
	}
}

func (planet *OChainPlanet) RemoveDefenses(id OChainDefenseID, count uint64) {
	switch id {
	case RocketLauncherID:
		planet.Defenses.RocketLauncher -= count
	case LightLaserID:
		planet.Defenses.LightLaser -= count
	case HeavyLaserID:
		planet.Defenses.HeavyLaser -= count
	case IonCannonID:
		planet.Defenses.IonCannon -= count
	case GaussCannonID:
		planet.Defenses.GaussCannon -= count
	case PlasmaTurretID:
		planet.Defenses.PlasmaTurret -= count
	case DarkMatterCanonID:
		planet.Defenses.DarkMatterCanon -= count
	}
}

func (planet *OChainPlanet) AddSpaceships(id OChainSpaceshipID, count uint64) {
	switch id {
	case SmallCargoID:
		planet.Spaceships.SmallCargo += count
	case LargeCargoID:
		planet.Spaceships.LargeCargo += count
	case LightFighterID:
		planet.Spaceships.LightFighter += count
	case HeavyFighterID:
		planet.Spaceships.HeavyFighter += count
	case CruiserID:
		planet.Spaceships.Cruiser += count
	case BattleshipID:
		planet.Spaceships.Battleship += count
	case BattlecruiserID:
		planet.Spaceships.Battlecruiser += count
	case BomberID:
		planet.Spaceships.Bomber += count
	case DestroyerID:
		planet.Spaceships.Destroyer += count
	case DeathstarID:
		planet.Spaceships.Deathstar += count
	case ReaperID:
		planet.Spaceships.Reaper += count
	case RecyclerID:
		planet.Spaceships.Recycler += count
	}
}

func (planet *OChainPlanet) RemoveSpaceships(id OChainSpaceshipID, count uint64) {
	switch id {
	case SmallCargoID:
		planet.Spaceships.SmallCargo -= count
	case LargeCargoID:
		planet.Spaceships.LargeCargo -= count
	case LightFighterID:
		planet.Spaceships.LightFighter -= count
	case HeavyFighterID:
		planet.Spaceships.HeavyFighter -= count
	case CruiserID:
		planet.Spaceships.Cruiser -= count
	case BattleshipID:
		planet.Spaceships.Battleship -= count
	case BattlecruiserID:
		planet.Spaceships.Battlecruiser -= count
	case BomberID:
		planet.Spaceships.Bomber -= count
	case DestroyerID:
		planet.Spaceships.Destroyer -= count
	case DeathstarID:
		planet.Spaceships.Deathstar -= count
	case ReaperID:
		planet.Spaceships.Reaper -= count
	case RecyclerID:
		planet.Spaceships.Recycler -= count
	}
}

func (planet *OChainPlanet) AddResourceById(id MarketResourceID, amount uint64) {
	switch id {
	case OCTResourceID:
		planet.Resources.OCT += amount
	case MetalResourceID:
		planet.Resources.Metal += amount
	case CrystalResourceID:
		planet.Resources.Crystal += amount
	case DeuteriumResourceID:
		planet.Resources.Deuterium += amount
	}
}

func (planet *OChainPlanet) RemoveResourceById(id MarketResourceID, amount uint64) error {
	switch id {
	case OCTResourceID:
		if planet.Resources.OCT < amount {
			return errors.New("remove resource overflow")
		}
		planet.Resources.OCT -= amount
	case MetalResourceID:
		if planet.Resources.Metal < amount {
			return errors.New("remove resource overflow")
		}
		planet.Resources.Metal -= amount
	case CrystalResourceID:
		if planet.Resources.Crystal < amount {
			return errors.New("remove resource overflow")
		}
		planet.Resources.Crystal -= amount
	case DeuteriumResourceID:
		if planet.Resources.Deuterium < amount {
			return errors.New("remove resource overflow")
		}
		planet.Resources.Deuterium -= amount
	}

	return nil
}

func (planet *OChainPlanet) CanPay(cost OChainResources) bool {

	if planet.Resources.OCT < cost.OCT {
		return false
	}

	if planet.Resources.Metal < cost.Metal {
		return false
	}

	if planet.Resources.Crystal < cost.Crystal {
		return false
	}

	if planet.Resources.Deuterium < cost.Deuterium {
		return false
	}

	return true
}

func (planet *OChainPlanet) Pay(cost OChainResources) error {
	payable := planet.CanPay(cost)
	if !payable {
		return errors.New("insuficient resources")
	}

	planet.Resources.OCT -= cost.OCT
	planet.Resources.Metal -= cost.Metal
	planet.Resources.Crystal -= cost.Crystal
	planet.Resources.Deuterium -= cost.Deuterium
	return nil
}

func (planet *OChainPlanet) UpdateBuildQueue(timestamp uint64) {
	newState := []OChainBuildQueueItem{}
	for i := 0; i < len(planet.BuildQueue); i++ {
		item := planet.BuildQueue[i]

		//if currently processing or ended
		if item.StartedAt > timestamp {

			//if item in queue is finished
			if item.FinishAt <= timestamp {
				if item.BuildType == OChainDefenseBuild {
					planet.AddDefenses(OChainDefenseID(item.BuildId), item.Count)
				}
				if item.BuildType == OChainSpaceshipBuild {
					planet.AddSpaceships(OChainSpaceshipID(item.BuildId), item.Count)
				}
			} else {
				//if item isn't finished yet
				itemDuration := (item.FinishAt - item.StartedAt) / item.Count
				effectiveTime := timestamp - item.StartedAt

				itemsFinished := effectiveTime / itemDuration
				if item.BuildType == OChainDefenseBuild {
					planet.AddDefenses(OChainDefenseID(item.BuildId), itemsFinished)
				}
				if item.BuildType == OChainSpaceshipBuild {
					planet.AddSpaceships(OChainSpaceshipID(item.BuildId), itemsFinished)
				}

				item.Count -= itemsFinished
				item.StartedAt = item.StartedAt + (itemsFinished * itemDuration)

				newState = append(newState, item)
			}
		}

	}
	planet.BuildQueue = newState
}

func (planet *OChainPlanet) UpdateResources(speed uint64, timestamp int64, account OChainUniverseAccount) {
	var totalConsumption int64 = planet.computeEnergyConsumtion()
	var totalProducedEnergy int64 = planet.computeTotalEnergy(account)

	if account.HasEngineerCommander(timestamp) {
		totalProducedEnergy = totalProducedEnergy * 110 / 100
	}

	var energyRate float64 = float64(totalProducedEnergy) / float64(totalConsumption)
	if energyRate > 1 {
		energyRate = 1
	}

	planet.Resources.Metal = planet.computeMetalProduction(energyRate, timestamp, speed, account)
	planet.Resources.Crystal = planet.computeCrystalProduction(energyRate, timestamp, speed, account)
	planet.Resources.Deuterium = planet.computeDeuteriumProduction(energyRate, timestamp, speed, account)
	planet.LastResourceUpdate = timestamp
}

func (planet *OChainPlanet) computeEnergyConsumtion() int64 {
	var metalFactor float64 = math.Pow(float64(1.1), float64(planet.Buildings.MetalMine))
	var metalEnergy float64 = 10 * float64(planet.Buildings.MetalMine) * metalFactor

	var crystalFactor float64 = math.Pow(float64(1.6), float64(planet.Buildings.CrystalMine))
	var crystalEnergy float64 = 10 * float64(planet.Buildings.CrystalMine) * crystalFactor

	var deuteriumFactor float64 = math.Pow(float64(1.1), float64(planet.Buildings.DeuteriumMine))
	var deuteriumEnergy float64 = 20 * float64(planet.Buildings.DeuteriumMine) * deuteriumFactor

	return int64(metalEnergy) + int64(crystalEnergy) + int64(deuteriumEnergy)
}

func (planet *OChainPlanet) computeTotalEnergy(account OChainUniverseAccount) int64 {
	var solarPlantEnergy float64 = 20 * float64(planet.Buildings.SolarPowerPlant) * math.Pow(float64(1.1), float64(planet.Buildings.SolarPowerPlant))
	return int64(solarPlantEnergy * math.Pow(1.05, float64(account.Technologies.Energy)))
}

func (planet *OChainPlanet) computeMetalProduction(energyRate float64, timestamp int64, speed uint64, account OChainUniverseAccount) uint64 {
	baseProductionPerHour := 30 * speed
	var factor float64 = math.Pow(float64(1.1), float64(planet.Buildings.MetalMine-1))
	var mineProductionPerHour float64 = 30 * float64(speed) * float64(planet.Buildings.MetalMine) * factor * (1 + (float64(account.Technologies.Plasma) / 100))

	secondsSinceLastUpdate := timestamp - planet.LastResourceUpdate
	metalEarnedSinceLastUpdate := (baseProductionPerHour + uint64(mineProductionPerHour)) * uint64(secondsSinceLastUpdate) / 60

	currentTotalMetal := planet.Resources.Metal + metalEarnedSinceLastUpdate

	if account.HasGeologistCommander(timestamp) {
		currentTotalMetal = currentTotalMetal * 110 / 100
	}

	return uint64(float64(currentTotalMetal) * energyRate)
}

func (planet *OChainPlanet) computeCrystalProduction(energyRate float64, timestamp int64, speed uint64, account OChainUniverseAccount) uint64 {
	baseProductionPerHour := 15 * speed
	var factor float64 = math.Pow(float64(1.1), float64(planet.Buildings.CrystalMine))
	var mineProductionPerHour float64 = 20 * float64(speed) * float64(planet.Buildings.CrystalMine) * factor * (1 + (float64(account.Technologies.Plasma) * 0.0066))

	secondsSinceLastUpdate := timestamp - planet.LastResourceUpdate
	crystalEarnedSinceLastUpdate := (baseProductionPerHour + uint64(mineProductionPerHour)) * uint64(secondsSinceLastUpdate) / 60

	currentTotalCrystal := planet.Resources.Metal + crystalEarnedSinceLastUpdate

	if account.HasGeologistCommander(timestamp) {
		currentTotalCrystal = currentTotalCrystal * 110 / 100
	}

	return uint64(float64(currentTotalCrystal) * energyRate)
}

func (planet *OChainPlanet) computeDeuteriumProduction(energyRate float64, timestamp int64, speed uint64, account OChainUniverseAccount) uint64 {
	var factor float64 = math.Pow(float64(1.1), float64(planet.Buildings.DeuteriumMine))
	var mineProductionPerHour float64 = 20 * float64(speed) * float64(planet.Buildings.DeuteriumMine) * factor

	secondsSinceLastUpdate := timestamp - planet.LastResourceUpdate
	deuteriumEarnedSinceLastUpdate := uint64(mineProductionPerHour) * uint64(secondsSinceLastUpdate) / 60

	if account.HasGeologistCommander(timestamp) {
		deuteriumEarnedSinceLastUpdate = deuteriumEarnedSinceLastUpdate * 110 / 100
	}

	return uint64(float64(deuteriumEarnedSinceLastUpdate) * energyRate)
}

func CoordinateId(galaxy uint64, solarSystem uint64, planet uint64) string {
	return fmt.Sprint(galaxy) + "_" + fmt.Sprint(solarSystem) + "_" + fmt.Sprint(planet)
}
