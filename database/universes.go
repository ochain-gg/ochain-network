package database

import (
	"errors"

	"github.com/dgraph-io/badger/v4"
	"github.com/ochain.gg/ochain-network/types"
	"github.com/timshannon/badgerhold/v4"
)

type OChainUniverseTable struct {
	db *badgerhold.Store
}

type OChainUniverseConfiguration struct {
	//Global configuration
	Speed                   uint
	MaxGalaxy               uint
	MaxSolarSystemPerGalaxy uint
	MaxPlanetPerSolarSystem uint

	Spaceships []types.OChainSpaceship
	Defenses   []types.OChainDefense
}

type OChainUniverse struct {
	Id            uint `badgerhold:"key"`
	Name          string
	Configuration OChainUniverseConfiguration
	CreatedAt     uint
}

func (table *OChainUniverseTable) Get(id uint) (OChainUniverse, error) {
	var result []OChainUniverse
	err := table.db.Find(&result, badgerhold.Where("Id").Eq(id))

	if err != nil {
		return OChainUniverse{}, err
	}

	if len(result) == 0 {
		return OChainUniverse{}, errors.New("universe not found")
	}

	return result[0], nil
}

func (table *OChainUniverseTable) GetAll(id uint) []OChainUniverse {
	var result []OChainUniverse
	table.db.Find(&result, badgerhold.Where("Id").Eq(id))

	return result
}

func (table *OChainUniverseTable) Insert(universe OChainUniverse, tx *badger.Txn) error {
	err := table.db.TxInsert(tx, badgerhold.NextSequence(), &universe)
	return err
}

func (table *OChainUniverseTable) Save(universe OChainUniverse, tx *badger.Txn) error {
	err := table.db.TxUpdate(tx, universe.Id, universe)
	return err
}

func (table *OChainUniverseTable) Delete(universe OChainUniverse, tx *badger.Txn) error {
	err := table.db.TxDelete(tx, universe.Id, universe)
	return err
}

func NewOChainUniverseTable(db *badgerhold.Store) *OChainUniverseTable {
	return &OChainUniverseTable{
		db: db,
	}
}
