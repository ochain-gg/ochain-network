package database

import (
	"errors"
	"math"

	"github.com/dgraph-io/badger/v4"
	"github.com/fxamacker/cbor/v2"
	"github.com/ochain-gg/ochain-network/types"
)

const (
	OChainBuildingPrefix string = "building_"
)

type OChainBuildingTable struct {
	bdb        *badger.DB
	currentTxn *badger.Txn
}

func (db *OChainBuildingTable) SetCurrentTxn(tx *badger.Txn) {
	db.currentTxn = tx
}

func (db *OChainBuildingTable) Exists(id string) (bool, error) {
	var at uint64 = math.MaxUint64
	return db.ExistsAt(id, at)
}

func (db *OChainBuildingTable) ExistsAt(id string, at uint64) (bool, error) {
	key := []byte(OChainBuildingPrefix + id)
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

func (db *OChainBuildingTable) Get(id string) (types.OChainBuilding, error) {
	var at uint64 = math.MaxUint64
	return db.GetAt(id, at)
}

func (db *OChainBuildingTable) GetAt(id string, at uint64) (types.OChainBuilding, error) {
	var building types.OChainBuilding
	key := []byte(OChainBuildingPrefix + id)
	txn := db.bdb.NewTransactionAt(at, false)

	item, err := txn.Get([]byte(key))
	if err != nil {
		return types.OChainBuilding{}, err
	}

	value, err := item.ValueCopy(nil)
	if err != nil {
		return types.OChainBuilding{}, err
	}

	err = cbor.Unmarshal(value, &building)
	if err != nil {
		return types.OChainBuilding{}, err
	}

	return building, nil
}

func (db *OChainBuildingTable) Insert(building types.OChainBuilding) error {
	key := []byte(OChainBuildingPrefix + string(building.Id))

	exists, err := db.ExistsAt(string(building.Id), db.currentTxn.ReadTs())
	if err != nil {
		return err
	}

	if exists {
		return errors.New("global building already exists")
	}

	value, err := cbor.Marshal(building)
	if err != nil {
		return err
	}

	return db.currentTxn.Set(key, value)
}

func (db *OChainBuildingTable) Update(building types.OChainBuilding) error {
	key := []byte(OChainBuildingPrefix + string(building.Id))

	exists, err := db.ExistsAt(string(building.Id), db.currentTxn.ReadTs())
	if err != nil {
		return err
	}

	if !exists {
		return errors.New("building doesn't exists")
	}

	value, err := cbor.Marshal(building)
	if err != nil {
		return err
	}

	return db.currentTxn.Set(key, value)
}

func (db *OChainBuildingTable) Upsert(building types.OChainBuilding) error {
	key := []byte(OChainBuildingPrefix + string(building.Id))
	value, err := cbor.Marshal(building)
	if err != nil {
		return err
	}

	return db.currentTxn.Set(key, value)
}

func (db *OChainBuildingTable) Delete(id string) error {
	key := []byte(OChainBuildingPrefix + id)
	return db.currentTxn.Delete(key)
}

func (db *OChainBuildingTable) GetAll() ([]types.OChainBuilding, error) {
	var at uint64 = math.MaxUint64
	return db.GetAllAt(at)
}

func (db *OChainBuildingTable) GetAllAt(at uint64) ([]types.OChainBuilding, error) {
	var buildings []types.OChainBuilding

	txn := db.bdb.NewTransactionAt(at, false)
	it := txn.NewIterator(badger.DefaultIteratorOptions)
	defer it.Close()

	prefix := []byte(OChainBuildingPrefix)

	for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
		item := it.Item()

		var building types.OChainBuilding
		value, err := item.ValueCopy(nil)
		if err != nil {
			return []types.OChainBuilding{}, err
		}

		err = cbor.Unmarshal(value, &building)
		if err != nil {
			return []types.OChainBuilding{}, err
		}

		buildings = append(buildings, building)
	}

	return buildings, nil
}

func NewOChainBuildingTable(db *badger.DB) *OChainBuildingTable {
	return &OChainBuildingTable{
		bdb: db,
	}
}
