package database

import (
	"errors"
	"math"

	"github.com/dgraph-io/badger/v4"
	"github.com/fxamacker/cbor/v2"
	"github.com/ochain-gg/ochain-network/types"
)

const (
	OChainFleetPrefix string = "fleet_"
)

type OChainFleetTable struct {
	bdb        *badger.DB
	currentTxn *badger.Txn
}

func (db *OChainFleetTable) SetCurrentTxn(tx *badger.Txn) {
	db.currentTxn = tx
}

func (db *OChainFleetTable) Exists(universeId string, id string) (bool, error) {
	var at uint64 = math.MaxUint64
	return db.ExistsAt(universeId, id, at)
}

func (db *OChainFleetTable) ExistsAt(universeId string, id string, at uint64) (bool, error) {
	key := []byte(OChainFleetPrefix + universeId + "_" + id)

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

func (db *OChainFleetTable) Get(universeId string, account types.OChainUniverseAccount, id string) (types.OChainFleet, error) {
	var at uint64 = math.MaxUint64
	return db.GetAt(universeId, account, id, at)
}

func (db *OChainFleetTable) GetAt(universeId string, account types.OChainUniverseAccount, id string, at uint64) (types.OChainFleet, error) {
	var planet types.OChainFleet
	key := []byte(OChainFleetPrefix + universeId + "_" + account.Address + "_" + id)
	txn := db.bdb.NewTransactionAt(at, false)

	item, err := txn.Get([]byte(key))
	if err != nil {
		return types.OChainFleet{}, err
	}

	value, err := item.ValueCopy(nil)
	if err != nil {
		return types.OChainFleet{}, err
	}

	err = cbor.Unmarshal(value, &planet)
	if err != nil {
		return types.OChainFleet{}, err
	}

	return planet, nil
}

func (db *OChainFleetTable) Insert(universeId string, account types.OChainUniverseAccount, fleet types.OChainFleet) error {
	key := []byte(OChainFleetPrefix + universeId + "_" + account.Address + "_" + fleet.Id)

	exists, err := db.ExistsAt(universeId, fleet.Id, db.currentTxn.ReadTs())
	if err != nil {
		return err
	}

	if exists {
		return errors.New("planet already exists")
	}

	value, err := cbor.Marshal(fleet)
	if err != nil {
		return err
	}

	return db.currentTxn.Set(key, value)
}

func (db *OChainFleetTable) Update(universeId string, account types.OChainUniverseAccount, fleet types.OChainFleet) error {
	key := []byte(OChainFleetPrefix + universeId + "_" + account.Address + "_" + fleet.Id)

	exists, err := db.ExistsAt(universeId, fleet.Id, db.currentTxn.ReadTs())
	if err != nil {
		return err
	}

	if !exists {
		return errors.New("fleet doesn't exists")
	}

	value, err := cbor.Marshal(fleet)
	if err != nil {
		return err
	}

	return db.currentTxn.Set(key, value)
}

func (db *OChainFleetTable) Upsert(universeId string, account types.OChainUniverseAccount, fleet types.OChainFleet) error {
	key := []byte(OChainFleetPrefix + universeId + "_" + account.Address + "_" + fleet.Id)
	value, err := cbor.Marshal(fleet)
	if err != nil {
		return err
	}

	return db.currentTxn.Set(key, value)
}

func (db *OChainFleetTable) Delete(universeId string, account types.OChainUniverseAccount, id string) error {

	key := []byte(OChainFleetPrefix + universeId + "_" + account.Address + "_" + id)
	return db.currentTxn.Delete(key)
}

func (db *OChainFleetTable) GetAccountFleet(universeId string, account types.OChainUniverseAccount) ([]types.OChainFleet, error) {
	var at uint64 = math.MaxUint64
	return db.GetAccountFleetAt(universeId, account, at)
}

func (db *OChainFleetTable) GetAccountFleetAt(universeId string, account types.OChainUniverseAccount, at uint64) ([]types.OChainFleet, error) {
	var fleets []types.OChainFleet

	txn := db.bdb.NewTransactionAt(at, false)
	it := txn.NewIterator(badger.DefaultIteratorOptions)
	defer it.Close()

	prefix := []byte(OChainFleetPrefix + universeId + "_" + account.Address)

	for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
		item := it.Item()

		var fleet types.OChainFleet
		value, err := item.ValueCopy(nil)
		if err != nil {
			return []types.OChainFleet{}, err
		}

		err = cbor.Unmarshal(value, &fleet)
		if err != nil {
			return []types.OChainFleet{}, err
		}

		fleets = append(fleets, fleet)
	}

	return fleets, nil
}

func (db *OChainFleetTable) GetAll(universeId string) ([]types.OChainFleet, error) {
	var at uint64 = math.MaxUint64
	return db.GetAllAt(universeId, at)
}

func (db *OChainFleetTable) GetAllAt(universeId string, at uint64) ([]types.OChainFleet, error) {
	var fleets []types.OChainFleet

	txn := db.bdb.NewTransactionAt(at, false)
	it := txn.NewIterator(badger.DefaultIteratorOptions)
	defer it.Close()

	prefix := []byte(OChainFleetPrefix + universeId)

	for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
		item := it.Item()

		var fleet types.OChainFleet
		value, err := item.ValueCopy(nil)
		if err != nil {
			return []types.OChainFleet{}, err
		}

		err = cbor.Unmarshal(value, &fleet)
		if err != nil {
			return []types.OChainFleet{}, err
		}

		fleets = append(fleets, fleet)
	}

	return fleets, nil
}

func NewFleetTable(db *badger.DB) *OChainFleetTable {
	return &OChainFleetTable{
		bdb: db,
	}
}
