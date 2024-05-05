package database

import (
	"errors"

	"github.com/dgraph-io/badger/v4"
	"github.com/ochain-gg/ochain-network/types"
	"github.com/timshannon/badgerhold/v4"
)

type OChainUniverseTable struct {
	db *badgerhold.Store
}

func (table *OChainUniverseTable) Get(id uint) (types.OChainUniverse, error) {
	var result []types.OChainUniverse
	err := table.db.Find(&result, badgerhold.Where("Id").Eq(id))

	if err != nil {
		return types.OChainUniverse{}, err
	}

	if len(result) == 0 {
		return types.OChainUniverse{}, errors.New("universe not found")
	}

	return result[0], nil
}

func (table *OChainUniverseTable) GetAll() []types.OChainUniverse {
	var result []types.OChainUniverse
	q := &badgerhold.Query{}

	table.db.Find(&result, q)
	return result
}

func (table *OChainUniverseTable) Insert(universe types.OChainUniverse, tx *badger.Txn) error {
	err := table.db.TxInsert(tx, badgerhold.NextSequence(), &universe)
	return err
}

func (table *OChainUniverseTable) Save(universe types.OChainUniverse, tx *badger.Txn) error {
	err := table.db.TxUpdate(tx, universe.Id, universe)
	return err
}

func (table *OChainUniverseTable) Delete(universe types.OChainUniverse, tx *badger.Txn) error {
	err := table.db.TxDelete(tx, universe.Id, universe)
	return err
}

func NewOChainUniverseTable(db *badgerhold.Store) *OChainUniverseTable {
	return &OChainUniverseTable{
		db: db,
	}
}

type OChainUniverseConfigurationTable struct {
	db *badgerhold.Store
}

func (table *OChainUniverseConfigurationTable) Get(id uint) (types.OChainUniverseConfiguration, error) {
	var result []types.OChainUniverseConfiguration
	err := table.db.Find(&result, badgerhold.Where("Id").Eq(id))

	if err != nil {
		return types.OChainUniverseConfiguration{}, err
	}

	if len(result) == 0 {
		return types.OChainUniverseConfiguration{}, errors.New("universe not found")
	}

	return result[0], nil
}

func (table *OChainUniverseConfigurationTable) GetAll() []types.OChainUniverseConfiguration {
	var result []types.OChainUniverseConfiguration
	q := &badgerhold.Query{}

	table.db.Find(&result, q)
	return result
}

func (table *OChainUniverseConfigurationTable) Insert(universe types.OChainUniverseConfiguration, tx *badger.Txn) error {
	err := table.db.TxInsert(tx, badgerhold.NextSequence(), &universe)
	return err
}

func (table *OChainUniverseConfigurationTable) Save(universe types.OChainUniverseConfiguration, tx *badger.Txn) error {
	err := table.db.TxUpdate(tx, universe.Id, universe)
	return err
}

func (table *OChainUniverseConfigurationTable) Delete(universe types.OChainUniverseConfiguration, tx *badger.Txn) error {
	err := table.db.TxDelete(tx, universe.Id, universe)
	return err
}

func NewOChainUniverseConfigurationTable(db *badgerhold.Store) *OChainUniverseConfigurationTable {
	return &OChainUniverseConfigurationTable{
		db: db,
	}
}
