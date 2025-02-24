package database

import (
	"errors"
	"math"

	"github.com/dgraph-io/badger/v4"
	"github.com/fxamacker/cbor/v2"
	"github.com/ochain-gg/ochain-network/types"
)

const (
	OChainSpaceshipPrefix string = "spaceship_"
)

type OChainSpaceshipTable struct {
	bdb        *badger.DB
	currentTxn *badger.Txn
}

func (db *OChainSpaceshipTable) SetCurrentTxn(tx *badger.Txn) {
	db.currentTxn = tx
}

func (db *OChainSpaceshipTable) Exists(id string) (bool, error) {
	var at uint64 = math.MaxUint64
	return db.ExistsAt(id, at)
}

func (db *OChainSpaceshipTable) ExistsAt(id string, at uint64) (bool, error) {
	key := []byte(OChainSpaceshipPrefix + id)
	txn := db.bdb.NewTransactionAt(at, false)
	if _, err := txn.Get([]byte(key)); err != nil {
		if errors.Is(err, badger.ErrKeyNotFound) {
			return false, nil
		} else {
			return false, err
		}
	} else {
		return true, nil
	}
}

func (db *OChainSpaceshipTable) Get(id string) (types.OChainSpaceship, error) {
	var at uint64 = math.MaxUint64
	return db.GetAt(id, at)
}

func (db *OChainSpaceshipTable) GetAt(id string, at uint64) (types.OChainSpaceship, error) {
	var spaceship types.OChainSpaceship
	key := []byte(OChainSpaceshipPrefix + id)
	txn := db.bdb.NewTransactionAt(at, false)

	item, err := txn.Get([]byte(key))
	if err != nil {
		return types.OChainSpaceship{}, err
	}

	value, err := item.ValueCopy(nil)
	if err != nil {
		return types.OChainSpaceship{}, err
	}

	err = cbor.Unmarshal(value, &spaceship)
	if err != nil {
		return types.OChainSpaceship{}, err
	}

	return spaceship, nil
}

func (db *OChainSpaceshipTable) Insert(spaceship types.OChainSpaceship) error {
	key := []byte(OChainSpaceshipPrefix + string(spaceship.Id))

	exists, err := db.ExistsAt(string(spaceship.Id), db.currentTxn.ReadTs())
	if err != nil {
		return err
	}

	if exists {
		return errors.New("spaceship already exists")
	}

	value, err := cbor.Marshal(spaceship)
	if err != nil {
		return err
	}

	return db.currentTxn.Set(key, value)
}

func (db *OChainSpaceshipTable) Update(spaceship types.OChainSpaceship) error {
	key := []byte(OChainSpaceshipPrefix + string(spaceship.Id))

	exists, err := db.ExistsAt(string(spaceship.Id), db.currentTxn.ReadTs())
	if err != nil {
		return err
	}

	if !exists {
		return errors.New("spaceship doesn't exists")
	}

	value, err := cbor.Marshal(spaceship)
	if err != nil {
		return err
	}

	return db.currentTxn.Set(key, value)
}

func (db *OChainSpaceshipTable) Upsert(spaceship types.OChainSpaceship) error {
	key := []byte(OChainSpaceshipPrefix + string(spaceship.Id))
	value, err := cbor.Marshal(spaceship)
	if err != nil {
		return err
	}

	return db.currentTxn.Set(key, value)
}

func (db *OChainSpaceshipTable) Delete(id string) error {
	key := []byte(OChainSpaceshipPrefix + id)
	return db.currentTxn.Delete(key)
}

func (db *OChainSpaceshipTable) GetCargoOf(fleet types.OChainFleet) (uint64, error) {
	var at uint64 = math.MaxUint64
	return db.GetCargoOfAt(fleet, at)
}

func (db *OChainSpaceshipTable) GetCargoOfAt(fleet types.OChainFleet, at uint64) (uint64, error) {

	spaceships, err := db.GetAllAt(at)
	if err != nil {
		return 0, err
	}

	var cargo uint64 = 0
	for _, spaceship := range spaceships {
		cargo += spaceship.Stats.Capacity * fleet.GetSpaceships(spaceship.Id)
	}

	return cargo, nil
}

func (db *OChainSpaceshipTable) GetAll() ([]types.OChainSpaceship, error) {
	var at uint64 = math.MaxUint64
	return db.GetAllAt(at)
}

func (db *OChainSpaceshipTable) GetAllAt(at uint64) ([]types.OChainSpaceship, error) {
	var spaceships []types.OChainSpaceship

	txn := db.bdb.NewTransactionAt(at, false)
	it := txn.NewIterator(badger.DefaultIteratorOptions)
	defer it.Close()

	prefix := []byte(OChainSpaceshipPrefix)

	for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
		item := it.Item()

		var spaceship types.OChainSpaceship
		value, err := item.ValueCopy(nil)
		if err != nil {
			return []types.OChainSpaceship{}, err
		}

		err = cbor.Unmarshal(value, &spaceship)
		if err != nil {
			return []types.OChainSpaceship{}, err
		}

		spaceships = append(spaceships, spaceship)
	}

	return spaceships, nil
}

func NewOChainSpaceshipTable(db *badger.DB) *OChainSpaceshipTable {
	return &OChainSpaceshipTable{
		bdb: db,
	}
}
