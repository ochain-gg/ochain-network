package types

import (
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

	LastResourceUpdate int64 `cbor:"10,keyasint"`
}

func (planet *OChainPlanet) CoordinateId() string {
	return fmt.Sprint(planet.Galaxy) + "_" + fmt.Sprint(planet.SolarSystem) + "_" + fmt.Sprint(planet.Planet)
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
