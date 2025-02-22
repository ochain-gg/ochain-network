package types

import (
	"errors"
	"fmt"
	"math"
)

type OChainFleetPositionWithAttributes struct {
	Galaxy      uint64 `cbor:"galaxy,keyasint"`
	SolarSystem uint64 `cbor:"solarSystem,keyasint"`
	Planet      uint64 `cbor:"planet,keyasint"`
}

type OChainFleetPosition struct {
	Galaxy      uint64 `cbor:"3,keyasint"`
	SolarSystem uint64 `cbor:"4,keyasint"`
	Planet      uint64 `cbor:"5,keyasint"`
}

func (position *OChainFleetPosition) CoordinateId() string {
	return fmt.Sprint(position.Galaxy) + "_" + fmt.Sprint(position.SolarSystem) + "_" + fmt.Sprint(position.Planet)
}

func (position *OChainFleetPosition) WithAttributes() OChainFleetPositionWithAttributes {
	return OChainFleetPositionWithAttributes{
		Galaxy:      position.Galaxy,
		SolarSystem: position.SolarSystem,
		Planet:      position.Planet,
	}
}

type OChainFleetWithAttributes struct {
	Id   string `cbor:"id,keyasint"`
	Name string `cbor:"name,keyasint"`

	Owner      string                              `cbor:"owner,keyasint"`
	Universe   string                              `cbor:"universe,keyasint"`
	Spaceships OChainFleetSpaceshipsWithAttributes `cbor:"spaceships,keyasint"`
	Cargo      OChainResourcesWithAttributes       `cbor:"cargo,keyasint"`

	Departure   OChainFleetPositionWithAttributes `cbor:"departure,keyasint"`
	Destination OChainFleetPositionWithAttributes `cbor:"destination,keyasint"`

	BeginTravelAt  int64 `cbor:"beginTravelAt,keyasint"`
	TravelDuration int64 `cbor:"travelDuration,keyasint"`
}

type OChainFleet struct {
	Id   string `cbor:"1,keyasint"`
	Name string `cbor:"2,keyasint"`

	Owner      string                `cbor:"3,keyasint"`
	Universe   string                `cbor:"4,keyasint"`
	Spaceships OChainFleetSpaceships `cbor:"5,keyasint"`
	Cargo      OChainResources       `cbor:"6,keyasint"`

	Departure   OChainFleetPosition `cbor:"7,keyasint"`
	Destination OChainFleetPosition `cbor:"8,keyasint"`

	BeginTravelAt  int64 `cbor:"9,keyasint"`
	TravelDuration int64 `cbor:"10,keyasint"`
}

func (fleet *OChainFleet) WithAttributes() OChainFleetWithAttributes {
	return OChainFleetWithAttributes{
		Id:             fleet.Id,
		Name:           fleet.Name,
		Owner:          fleet.Owner,
		Universe:       fleet.Universe,
		Spaceships:     fleet.Spaceships.WithAttributes(),
		Cargo:          fleet.Cargo.WithAttributes(),
		Departure:      fleet.Departure.WithAttributes(),
		Destination:    fleet.Destination.WithAttributes(),
		BeginTravelAt:  fleet.BeginTravelAt,
		TravelDuration: fleet.TravelDuration,
	}
}

func (fleet *OChainFleet) CanPay(cost OChainResources) bool {

	if fleet.Cargo.OCT < cost.OCT {
		return false
	}

	if fleet.Cargo.Metal < cost.Metal {
		return false
	}

	if fleet.Cargo.Crystal < cost.Crystal {
		return false
	}

	if fleet.Cargo.Deuterium < cost.Deuterium {
		return false
	}

	return true
}

func (fleet *OChainFleet) IsTraveling(timestamp int64) bool {
	return fleet.BeginTravelAt+fleet.TravelDuration > timestamp
}

func (fleet *OChainFleet) MaxSpeed(timestamp int64) uint64 {
	var maxSpeed uint64 = math.MaxUint64
	spaceshipIds := GetSpaceshipIds()
	for _, spaceshipId := range spaceshipIds {
		spaceshipNumber := fleet.GetSpaceships(spaceshipId)
		if spaceshipNumber > 0 {

		}
	}

	return maxSpeed
}

func (fleet *OChainFleet) GetSpaceships(id OChainSpaceshipID) uint64 {
	switch id {
	case SmallCargoID:
		return fleet.Spaceships.SmallCargo
	case LargeCargoID:
		return fleet.Spaceships.LargeCargo
	case LightFighterID:
		return fleet.Spaceships.LightFighter
	case HeavyFighterID:
		return fleet.Spaceships.HeavyFighter
	case CruiserID:
		return fleet.Spaceships.Cruiser
	case BattleshipID:
		return fleet.Spaceships.Battleship
	case BattlecruiserID:
		return fleet.Spaceships.Battlecruiser
	case BomberID:
		return fleet.Spaceships.Bomber
	case DestroyerID:
		return fleet.Spaceships.Destroyer
	case DeathstarID:
		return fleet.Spaceships.Deathstar
	case ReaperID:
		return fleet.Spaceships.Reaper
	case RecyclerID:
		return fleet.Spaceships.Recycler
	}

	return 0
}

func (fleet *OChainFleet) AddSpaceships(id OChainSpaceshipID, count uint64) {
	switch id {
	case SmallCargoID:
		fleet.Spaceships.SmallCargo += count
	case LargeCargoID:
		fleet.Spaceships.LargeCargo += count
	case LightFighterID:
		fleet.Spaceships.LightFighter += count
	case HeavyFighterID:
		fleet.Spaceships.HeavyFighter += count
	case CruiserID:
		fleet.Spaceships.Cruiser += count
	case BattleshipID:
		fleet.Spaceships.Battleship += count
	case BattlecruiserID:
		fleet.Spaceships.Battlecruiser += count
	case BomberID:
		fleet.Spaceships.Bomber += count
	case DestroyerID:
		fleet.Spaceships.Destroyer += count
	case DeathstarID:
		fleet.Spaceships.Deathstar += count
	case ReaperID:
		fleet.Spaceships.Reaper += count
	case RecyclerID:
		fleet.Spaceships.Recycler += count
	}
}

func (fleet *OChainFleet) RemoveSpaceships(id OChainSpaceshipID, count uint64) error {
	switch id {
	case SmallCargoID:
		if fleet.Spaceships.SmallCargo < count {
			return errors.New("no sufficient smallCargo")
		}
		fleet.Spaceships.SmallCargo -= count
		return nil
	case LargeCargoID:
		if fleet.Spaceships.LargeCargo < count {
			return errors.New("no sufficient largeCargo")
		}
		fleet.Spaceships.LargeCargo -= count
		return nil
	case LightFighterID:
		if fleet.Spaceships.LightFighter < count {
			return errors.New("no sufficient lightFighter")
		}
		fleet.Spaceships.LightFighter -= count
		return nil
	case HeavyFighterID:
		if fleet.Spaceships.HeavyFighter < count {
			return errors.New("no sufficient heavyFighter")
		}
		fleet.Spaceships.HeavyFighter -= count
		return nil
	case CruiserID:
		if fleet.Spaceships.Cruiser < count {
			return errors.New("no sufficient cruiser")
		}
		fleet.Spaceships.Cruiser -= count
		return nil
	case BattleshipID:
		if fleet.Spaceships.Battleship < count {
			return errors.New("no sufficient battleship")
		}
		fleet.Spaceships.Battleship -= count
		return nil
	case BattlecruiserID:
		if fleet.Spaceships.Battlecruiser < count {
			return errors.New("no sufficient battlecruiser")
		}
		fleet.Spaceships.Battlecruiser -= count
		return nil
	case BomberID:
		if fleet.Spaceships.Bomber < count {
			return errors.New("no sufficient bomber")
		}
		fleet.Spaceships.Bomber -= count
		return nil
	case DestroyerID:
		if fleet.Spaceships.Destroyer < count {
			return errors.New("no sufficient destroyer")
		}
		fleet.Spaceships.Destroyer -= count
		return nil
	case DeathstarID:
		if fleet.Spaceships.Deathstar < count {
			return errors.New("no sufficient deathstar")
		}
		fleet.Spaceships.Deathstar -= count
		return nil
	case ReaperID:
		if fleet.Spaceships.Reaper < count {
			return errors.New("no sufficient reaper")
		}
		fleet.Spaceships.Reaper -= count
		return nil
	case RecyclerID:
		if fleet.Spaceships.Recycler < count {
			return errors.New("no sufficient recycler")
		}
		fleet.Spaceships.Recycler -= count
		return nil
	}

	return errors.New("spaceship not found")
}

func (fleet *OChainFleet) AddResourceById(id MarketResourceID, amount uint64) {
	switch id {
	case OCTResourceID:
		fleet.Cargo.OCT += amount
	case MetalResourceID:
		fleet.Cargo.Metal += amount
	case CrystalResourceID:
		fleet.Cargo.Crystal += amount
	case DeuteriumResourceID:
		fleet.Cargo.Deuterium += amount
	}
}

func (fleet *OChainFleet) RemoveResourceById(id MarketResourceID, amount uint64) error {
	switch id {
	case OCTResourceID:
		if fleet.Cargo.OCT < amount {
			return errors.New("remove resource overflow")
		}
		fleet.Cargo.OCT -= amount
	case MetalResourceID:
		if fleet.Cargo.Metal < amount {
			return errors.New("remove resource overflow")
		}
		fleet.Cargo.Metal -= amount
	case CrystalResourceID:
		if fleet.Cargo.Crystal < amount {
			return errors.New("remove resource overflow")
		}
		fleet.Cargo.Crystal -= amount
	case DeuteriumResourceID:
		if fleet.Cargo.Deuterium < amount {
			return errors.New("remove resource overflow")
		}
		fleet.Cargo.Deuterium -= amount
	}

	return nil
}
