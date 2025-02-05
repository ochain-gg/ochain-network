package types

import (
	"errors"
	"fmt"
	"math"
	"time"
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

type OChainUpgradeWithAttributes struct {
	UniverseId         string            `cbor:"universeId"`
	PlanetCoordinateId string            `cbor:"planetCoordinateId"`
	UpgradeType        OChainUpgradeType `cbor:"upgradeType"`
	UpgradeId          string            `cbor:"upgradeId"`
	Level              uint64            `cbor:"level"`
	StartedAt          int64             `cbor:"startedAt"`
	EndedAt            int64             `cbor:"endedAt"`
	Executed           bool              `cbor:"executed"`
}

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

func (upgrade *OChainUpgrade) WithAttributes() OChainUpgradeWithAttributes {
	return OChainUpgradeWithAttributes{
		UniverseId:         upgrade.UniverseId,
		PlanetCoordinateId: upgrade.PlanetCoordinateId,
		UpgradeType:        upgrade.UpgradeType,
		UpgradeId:          upgrade.UpgradeId,
		Level:              upgrade.Level,
		StartedAt:          upgrade.StartedAt,
		EndedAt:            upgrade.EndedAt,
		Executed:           upgrade.Executed,
	}
}

type OChainBuildWithAttributes struct {
	BuildType OChainBuildType `cbor:"1,keyasint"`
	BuildId   string          `cbor:"2,keyasint"`
	Count     uint64          `cbor:"3,keyasint"`
}

type OChainBuild struct {
	BuildType OChainBuildType `cbor:"1,keyasint"`
	BuildId   string          `cbor:"2,keyasint"`
	Count     uint64          `cbor:"3,keyasint"`
}

func (build *OChainBuild) WithAttributes() OChainBuildWithAttributes {
	return OChainBuildWithAttributes{
		BuildType: build.BuildType,
		BuildId:   build.BuildId,
		Count:     build.Count,
	}
}

type OChainPlanetBuildingsWithAttributes struct {
	MetalMine           uint64 `cbor:"metalMine"`
	CrystalMine         uint64 `cbor:"crystalMine"`
	DeuteriumMine       uint64 `cbor:"deuteriumMine"`
	SolarPowerPlant     uint64 `cbor:"solarPowerPlant"`
	RoboticFactory      uint64 `cbor:"roboticFactory"`
	NaniteFactory       uint64 `cbor:"naniteFactory"`
	SpaceshipFactory    uint64 `cbor:"spaceshipFactory"`
	IntergalacticPortal uint64 `cbor:"intergalacticPortal"`
	ResearchLaboratory  uint64 `cbor:"researchLaboratory"`
	ShieldDome          uint64 `cbor:"shieldDome"`
	MetalStorage        uint64 `cbor:"metalStorage"`
	CrystalStorage      uint64 `cbor:"crystalStorage"`
	DeuteriumStorage    uint64 `cbor:"deuteriumStorage"`
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

	MetalStorage     uint64 `cbor:"11,keyasint"`
	CrystalStorage   uint64 `cbor:"12,keyasint"`
	DeuteriumStorage uint64 `cbor:"13,keyasint"`
}

func (builds *OChainPlanetBuildings) WithAttributes() OChainPlanetBuildingsWithAttributes {
	return OChainPlanetBuildingsWithAttributes{
		MetalMine:           builds.MetalMine,
		CrystalMine:         builds.CrystalMine,
		DeuteriumMine:       builds.DeuteriumMine,
		SolarPowerPlant:     builds.SolarPowerPlant,
		RoboticFactory:      builds.RoboticFactory,
		NaniteFactory:       builds.NaniteFactory,
		SpaceshipFactory:    builds.SpaceshipFactory,
		IntergalacticPortal: builds.IntergalacticPortal,
		ResearchLaboratory:  builds.ResearchLaboratory,
		ShieldDome:          builds.ShieldDome,
		MetalStorage:        builds.MetalStorage,
		CrystalStorage:      builds.CrystalStorage,
		DeuteriumStorage:    builds.DeuteriumStorage,
	}
}

type OChainFleetSpaceshipsWithAttributes struct {
	SmallCargo    uint64 `cbor:"smallCargo"`
	LargeCargo    uint64 `cbor:"largeCargo"`
	LightFighter  uint64 `cbor:"lightFighter"`
	HeavyFighter  uint64 `cbor:"heavyFighter"`
	Cruiser       uint64 `cbor:"cruiser"`
	Battleship    uint64 `cbor:"battleship"`
	Battlecruiser uint64 `cbor:"battlecruiser"`
	Bomber        uint64 `cbor:"bomber"`
	Destroyer     uint64 `cbor:"destroyer"`
	Deathstar     uint64 `cbor:"deathstar"`
	Reaper        uint64 `cbor:"reaper"`
	Recycler      uint64 `cbor:"recycler"`
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

func (fleet *OChainFleetSpaceships) WithAttributes() OChainFleetSpaceshipsWithAttributes {
	return OChainFleetSpaceshipsWithAttributes{
		SmallCargo:    fleet.SmallCargo,
		LargeCargo:    fleet.LargeCargo,
		LightFighter:  fleet.LightFighter,
		HeavyFighter:  fleet.HeavyFighter,
		Cruiser:       fleet.Cruiser,
		Battleship:    fleet.Battleship,
		Battlecruiser: fleet.Battlecruiser,
		Bomber:        fleet.Bomber,
		Destroyer:     fleet.Destroyer,
		Deathstar:     fleet.Deathstar,
		Reaper:        fleet.Reaper,
		Recycler:      fleet.Recycler,
	}
}

type OChainPlanetDefencesWithAttributes struct {
	RocketLauncher  uint64 `cbor:"rocketLauncher"`
	LightLaser      uint64 `cbor:"lightLaser"`
	HeavyLaser      uint64 `cbor:"heavyLaser"`
	IonCannon       uint64 `cbor:"ionCannon"`
	GaussCannon     uint64 `cbor:"gaussCannon"`
	PlasmaTurret    uint64 `cbor:"plasmaTurret"`
	DarkMatterCanon uint64 `cbor:"darkMatterCanon"`
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

func (def *OChainPlanetDefences) WithAttributes() OChainPlanetDefencesWithAttributes {
	return OChainPlanetDefencesWithAttributes{
		RocketLauncher:  def.RocketLauncher,
		LightLaser:      def.LightLaser,
		HeavyLaser:      def.HeavyLaser,
		IonCannon:       def.IonCannon,
		GaussCannon:     def.GaussCannon,
		PlasmaTurret:    def.PlasmaTurret,
		DarkMatterCanon: def.DarkMatterCanon,
	}
}

type OChainBuildQueueItemWithAttributes struct {
	BuildType OChainBuildType `cbor:"buildType"`
	BuildId   string          `cbor:"buildId"`
	Count     uint64          `cbor:"count"`
	StartAt   uint64          `cbor:"startAt"`
	FinishAt  uint64          `cbor:"finishAt"`
}

type OChainBuildQueueItem struct {
	BuildType OChainBuildType `cbor:"1,keyasint"`
	BuildId   string          `cbor:"2,keyasint"`
	Count     uint64          `cbor:"3,keyasint"`
	StartAt   uint64          `cbor:"4,keyasint"`
	FinishAt  uint64          `cbor:"5,keyasint"`
}

func (item *OChainBuildQueueItem) WithAttributes() OChainBuildQueueItemWithAttributes {
	return OChainBuildQueueItemWithAttributes{
		BuildType: item.BuildType,
		BuildId:   item.BuildId,
		Count:     item.Count,
		StartAt:   item.StartAt,
		FinishAt:  item.FinishAt,
	}
}

type OChainPlanetStatistics struct {
	MetalHourlyProduction     float64 `cbor:"metalHourlyProduction"`
	CrystalHourlyProduction   float64 `cbor:"crystalHourlyProduction"`
	DeutereumHourlyProduction float64 `cbor:"deutereumHourlyProduction"`

	MetalMaxStorageCapacity     uint64 `cbor:"metalMaxStorageCapacity"`
	CrystalMaxStorageCapacity   uint64 `cbor:"crystalMaxStorageCapacity"`
	DeutereumMaxStorageCapacity uint64 `cbor:"deutereumMaxStorageCapacity"`

	EnergyProduction int64 `cbor:"energyProduction"`
	EnergyConsumtion int64 `cbor:"energyConsumtion"`
}

type OChainPlanetWithAttributes struct {
	Owner       string `cbor:"owner"`
	Universe    string `cbor:"universe"`
	Galaxy      uint64 `cbor:"galaxy"`
	SolarSystem uint64 `cbor:"solarSystem"`
	Planet      uint64 `cbor:"planet"`

	Buildings  OChainPlanetBuildingsWithAttributes `cbor:"buildings"`
	Spaceships OChainFleetSpaceshipsWithAttributes `cbor:"spaceships"`
	Defenses   OChainPlanetDefencesWithAttributes  `cbor:"defenses"`
	Resources  OChainResourcesWithAttributes       `cbor:"resources"`

	BuildQueue         []OChainBuildQueueItemWithAttributes `cbor:"buildQueue"`
	LastResourceUpdate int64                                `cbor:"lastResourceUpdate"`
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

func (planet *OChainPlanet) WithAttributes() OChainPlanetWithAttributes {
	var queue []OChainBuildQueueItemWithAttributes
	for i := range planet.BuildQueue {
		queue = append(queue, planet.BuildQueue[i].WithAttributes())
	}

	return OChainPlanetWithAttributes{
		Owner:              planet.Owner,
		Universe:           planet.Universe,
		Galaxy:             planet.Galaxy,
		SolarSystem:        planet.SolarSystem,
		Planet:             planet.Planet,
		Buildings:          planet.Buildings.WithAttributes(),
		Spaceships:         planet.Spaceships.WithAttributes(),
		Defenses:           planet.Defenses.WithAttributes(),
		Resources:          planet.Resources.WithAttributes(),
		BuildQueue:         queue,
		LastResourceUpdate: planet.LastResourceUpdate,
	}
}

func (planet *OChainPlanet) CoordinateId() string {
	return fmt.Sprint(planet.Galaxy) + "_" + fmt.Sprint(planet.SolarSystem) + "_" + fmt.Sprint(planet.Planet)
}

func (planet *OChainPlanet) PlanetStatistics(speed uint64, timestamp int64, account OChainUniverseAccount) OChainPlanetStatistics {

	var totalConsumption int64 = planet.computeEnergyConsumtion()
	var totalProducedEnergy int64 = planet.computeTotalEnergy(account)

	var energyRate float64 = float64(totalProducedEnergy) / float64(totalConsumption)
	if energyRate > 1 {
		energyRate = 1
	}

	return OChainPlanetStatistics{
		MetalHourlyProduction:     planet.getMetalProduction(energyRate, timestamp, speed, account),
		CrystalHourlyProduction:   planet.getCrystalProduction(energyRate, timestamp, speed, account),
		DeutereumHourlyProduction: planet.getDeuteriumProduction(energyRate, timestamp, speed, account),

		MetalMaxStorageCapacity:     planet.getMetalStorageCapacity(),
		CrystalMaxStorageCapacity:   planet.getCrystalStorageCapacity(),
		DeutereumMaxStorageCapacity: planet.getDeuteriumStorageCapacity(),

		EnergyProduction: totalProducedEnergy,
		EnergyConsumtion: totalConsumption,
	}
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

func (planet *OChainPlanet) AddItemToBuildQueue(item OChainBuild, duration uint64) OChainBuildQueueItem {
	newState := []OChainBuildQueueItem{}
	queueFinishingAt := uint64(time.Now().Unix())
	for i := 0; i < len(planet.BuildQueue); i++ {
		queueItem := planet.BuildQueue[i]
		if queueFinishingAt < queueItem.FinishAt {
			queueFinishingAt = queueItem.FinishAt
		}

		newState = append(newState, queueItem)
	}

	newItem := OChainBuildQueueItem{
		BuildType: item.BuildType,
		BuildId:   item.BuildId,
		Count:     item.Count,
		StartAt:   queueFinishingAt,
		FinishAt:  queueFinishingAt + duration,
	}

	newState = append(newState, newItem)

	planet.BuildQueue = newState

	return newItem
}

func (planet *OChainPlanet) UpdateBuildQueue(timestamp uint64) {
	newState := []OChainBuildQueueItem{}
	for i := 0; i < len(planet.BuildQueue); i++ {
		item := planet.BuildQueue[i]

		//if currently processing or ended
		if item.StartAt > timestamp {

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
				itemDuration := (item.FinishAt - item.StartAt) / item.Count
				effectiveTime := timestamp - item.StartAt

				itemsFinished := effectiveTime / itemDuration
				if item.BuildType == OChainDefenseBuild {
					planet.AddDefenses(OChainDefenseID(item.BuildId), itemsFinished)
				}
				if item.BuildType == OChainSpaceshipBuild {
					planet.AddSpaceships(OChainSpaceshipID(item.BuildId), itemsFinished)
				}

				item.Count -= itemsFinished
				item.StartAt = item.StartAt + (itemsFinished * itemDuration)

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

	if planet.Resources.Metal < planet.getMetalStorageCapacity() {
		planet.Resources.Metal += planet.computeMetalProductionFromLastUpdate(energyRate, timestamp, speed, account)
		if planet.Resources.Metal > planet.getMetalStorageCapacity() {
			planet.Resources.Metal = planet.getMetalStorageCapacity()
		}
	}

	if planet.Resources.Crystal < planet.getCrystalStorageCapacity() {
		planet.Resources.Crystal += planet.computeCrystalProductionFromLastUpdate(energyRate, timestamp, speed, account)
		if planet.Resources.Crystal > planet.getCrystalStorageCapacity() {
			planet.Resources.Crystal = planet.getCrystalStorageCapacity()
		}
	}
	if planet.Resources.Deuterium < planet.getDeuteriumStorageCapacity() {
		planet.Resources.Deuterium += planet.computeDeuteriumProductionFromLastUpdate(energyRate, timestamp, speed, account)
		if planet.Resources.Deuterium > planet.getDeuteriumStorageCapacity() {
			planet.Resources.Deuterium = planet.getDeuteriumStorageCapacity()
		}
	}

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
	var solarPlantEnergy float64 = 20 + (20 * float64(planet.Buildings.SolarPowerPlant) * math.Pow(float64(1.1), float64(planet.Buildings.SolarPowerPlant)))
	return int64(solarPlantEnergy * math.Pow(1.05, float64(account.Technologies.Energy)))
}

func (planet *OChainPlanet) getMetalStorageCapacity() uint64 {
	return uint64(float64(5000) * float64(2.5) * math.Pow(float64(2.71828), float64(20*planet.Buildings.MetalStorage/33)))
}

func (planet *OChainPlanet) getMetalProduction(energyRate float64, timestamp int64, speed uint64, account OChainUniverseAccount) float64 {
	baseProductionPerHour := 30 * speed

	var mineProductionPerHour float64 = 0
	if planet.Buildings.MetalMine > 0 {
		var factor float64 = math.Pow(float64(1.1), float64(planet.Buildings.MetalMine-1))
		mineProductionPerHour = 30 * float64(speed) * float64(planet.Buildings.MetalMine) * factor * (1 + (float64(account.Technologies.Plasma) / 100))
	}

	productionPerHour := float64(baseProductionPerHour) + mineProductionPerHour

	if account.HasGeologistCommander(timestamp) {
		productionPerHour = productionPerHour * 110 / 100
	}

	return float64(productionPerHour) * energyRate
}

func (planet *OChainPlanet) computeMetalProductionFromLastUpdate(energyRate float64, timestamp int64, speed uint64, account OChainUniverseAccount) uint64 {

	hourlyProduction := planet.getMetalProduction(energyRate, timestamp, speed, account)

	secondsSinceLastUpdate := timestamp - planet.LastResourceUpdate
	metalEarnedSinceLastUpdate := hourlyProduction * float64(secondsSinceLastUpdate) / 3600

	return uint64(metalEarnedSinceLastUpdate)
}

func (planet *OChainPlanet) getCrystalStorageCapacity() uint64 {
	return uint64(float64(5000) * float64(2.5) * math.Pow(float64(2.71828), float64(20*planet.Buildings.CrystalStorage/33)))
}

func (planet *OChainPlanet) getCrystalProduction(energyRate float64, timestamp int64, speed uint64, account OChainUniverseAccount) float64 {
	baseProductionPerHour := 15 * speed
	var factor float64 = math.Pow(float64(1.1), float64(planet.Buildings.CrystalMine))
	var mineProductionPerHour float64 = 20 * float64(speed) * float64(planet.Buildings.CrystalMine) * factor * (1 + (float64(account.Technologies.Plasma) * 0.0066))

	productionPerHour := float64(baseProductionPerHour) + mineProductionPerHour

	if account.HasGeologistCommander(timestamp) {
		productionPerHour = productionPerHour * 110 / 100
	}

	return float64(productionPerHour) * energyRate
}

func (planet *OChainPlanet) computeCrystalProductionFromLastUpdate(energyRate float64, timestamp int64, speed uint64, account OChainUniverseAccount) uint64 {
	hourlyProduction := planet.getCrystalProduction(energyRate, timestamp, speed, account)

	secondsSinceLastUpdate := timestamp - planet.LastResourceUpdate
	crystalEarnedSinceLastUpdate := hourlyProduction * float64(secondsSinceLastUpdate) / 3600

	if account.HasGeologistCommander(timestamp) {
		crystalEarnedSinceLastUpdate = crystalEarnedSinceLastUpdate * 110 / 100
	}

	return uint64(crystalEarnedSinceLastUpdate)
}

func (planet *OChainPlanet) getDeuteriumStorageCapacity() uint64 {
	return uint64(float64(5000) * float64(2.5) * math.Pow(float64(2.71828), float64(20*planet.Buildings.DeuteriumStorage/33)))
}

func (planet *OChainPlanet) getDeuteriumProduction(energyRate float64, timestamp int64, speed uint64, account OChainUniverseAccount) float64 {
	var factor float64 = math.Pow(float64(1.1), float64(planet.Buildings.DeuteriumMine))
	var mineProductionPerHour float64 = 20 * float64(speed) * float64(planet.Buildings.DeuteriumMine) * factor

	if account.HasGeologistCommander(timestamp) {
		mineProductionPerHour = mineProductionPerHour * 110 / 100
	}

	return float64(mineProductionPerHour) * energyRate
}

func (planet *OChainPlanet) computeDeuteriumProductionFromLastUpdate(energyRate float64, timestamp int64, speed uint64, account OChainUniverseAccount) uint64 {
	var productionPerHour float64 = planet.getDeuteriumProduction(energyRate, timestamp, speed, account)

	secondsSinceLastUpdate := timestamp - planet.LastResourceUpdate
	deuteriumEarnedSinceLastUpdate := uint64(productionPerHour) * uint64(secondsSinceLastUpdate) / 3600

	return uint64(deuteriumEarnedSinceLastUpdate)
}

func CoordinateId(galaxy uint64, solarSystem uint64, planet uint64) string {
	return fmt.Sprint(galaxy) + "_" + fmt.Sprint(solarSystem) + "_" + fmt.Sprint(planet)
}
