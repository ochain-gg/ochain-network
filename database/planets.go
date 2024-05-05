package database

import (
	"errors"

	"github.com/dgraph-io/badger/v4"
	"github.com/ochain-gg/ochain-network/types"
	"github.com/timshannon/badgerhold/v4"
)

type OChainPlanetTable struct {
	db *badgerhold.Store
}

type OChainPlanet struct {
	Id           uint64 `badgerhold:"key"`
	Owner        uint64 `badgerhold:"index"`
	Universe     uint64 `badgerhold:"index"`
	Galaxy       uint64
	SolarSystem  uint64
	Planet       uint64
	CoordinateId string `badgerhold:"unique"`

	Buildings  types.OChainPlanetBuildings
	Spaceships types.OChainFleetSpaceships
	Defenses   types.OChainPlanetDefences
	Resources  types.OChainResources

	LastResourceUpdate uint
}

func (table *OChainPlanetTable) Get(id uint) (OChainPlanet, error) {
	var result []OChainPlanet
	err := table.db.Find(&result, badgerhold.Where("Id").Eq(id))

	if err != nil {
		return OChainPlanet{}, err
	}

	if len(result) == 0 {
		return OChainPlanet{}, errors.New("planet not found")
	}

	return result[0], nil
}

func (table *OChainPlanetTable) GetAll(id uint) []OChainPlanet {
	var result []OChainPlanet
	table.db.Find(&result, badgerhold.Where("Id").Eq(id))

	return result
}

func (table *OChainPlanetTable) Insert(planet OChainPlanet, tx *badger.Txn) error {
	err := table.db.TxInsert(tx, badgerhold.NextSequence(), &planet)
	return err
}

func (table *OChainPlanetTable) Save(planet OChainPlanet, tx *badger.Txn) error {
	err := table.db.TxUpdate(tx, planet.Id, planet)
	return err
}

func (table *OChainPlanetTable) Delete(planet OChainPlanet, tx *badger.Txn) error {
	err := table.db.TxDelete(tx, planet.Id, planet)
	return err
}

func NewOChainPlanetTable(db *badgerhold.Store) *OChainPlanetTable {
	return &OChainPlanetTable{
		db: db,
	}
}
