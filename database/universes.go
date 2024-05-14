package database

import (
	"errors"
	"math"

	"github.com/dgraph-io/badger/v4"
	"github.com/fxamacker/cbor/v2"
	"github.com/ochain-gg/ochain-network/types"
)

const (
	OChainUniversePrefix string = "universe_"
)

type OChainUniverseTable struct {
	bdb        *badger.DB
	currentTxn *badger.Txn
}

func (db *OChainUniverseTable) SetCurrentTxn(tx *badger.Txn) {
	db.currentTxn = tx
}

func (db *OChainUniverseTable) Exists(address string) (bool, error) {
	var at uint64 = math.MaxUint64
	return db.ExistsAt(address, at)
}

func (db *OChainUniverseTable) ExistsAt(address string, at uint64) (bool, error) {
	key := []byte(OChainUniversePrefix + address)
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

func (db *OChainUniverseTable) Get(address string) (types.OChainUniverse, error) {
	var at uint64 = math.MaxUint64
	return db.GetAt(address, at)
}

func (db *OChainUniverseTable) GetAt(address string, at uint64) (types.OChainUniverse, error) {
	var universe types.OChainUniverse
	key := []byte(OChainUniversePrefix + address)
	txn := db.bdb.NewTransactionAt(at, false)

	item, err := txn.Get([]byte(key))
	if err != nil {
		return types.OChainUniverse{}, err
	}

	value, err := item.ValueCopy(nil)
	if err != nil {
		return types.OChainUniverse{}, err
	}

	err = cbor.Unmarshal(value, &universe)
	if err != nil {
		return types.OChainUniverse{}, err
	}

	return universe, nil
}

func (db *OChainUniverseTable) Insert(universe types.OChainUniverse) error {
	key := []byte(OChainUniversePrefix + universe.Id)

	exists, err := db.ExistsAt(universe.Id, db.currentTxn.ReadTs())
	if err != nil {
		return err
	}

	if exists {
		return errors.New("global universe already exists")
	}

	value, err := cbor.Marshal(universe)
	if err != nil {
		return err
	}

	return db.currentTxn.Set(key, value)
}

func (db *OChainUniverseTable) Update(universe types.OChainUniverse) error {
	key := []byte(OChainUniversePrefix + universe.Id)

	exists, err := db.ExistsAt(universe.Id, db.currentTxn.ReadTs())
	if err != nil {
		return err
	}

	if !exists {
		return errors.New("universe doesn't exists")
	}

	value, err := cbor.Marshal(universe)
	if err != nil {
		return err
	}

	return db.currentTxn.Set(key, value)
}

func (db *OChainUniverseTable) Upsert(universe types.OChainUniverse) error {
	key := []byte(OChainUniversePrefix + universe.Id)
	value, err := cbor.Marshal(universe)
	if err != nil {
		return err
	}

	return db.currentTxn.Set(key, value)
}

func (db *OChainUniverseTable) Delete(id string) error {
	key := []byte(OChainUniversePrefix + id)
	return db.currentTxn.Delete(key)
}

func (db *OChainUniverseTable) GetAll() ([]types.OChainUniverse, error) {
	var at uint64 = math.MaxUint64
	return db.GetAllAt(at)
}

func (db *OChainUniverseTable) GetAllAt(at uint64) ([]types.OChainUniverse, error) {
	var universes []types.OChainUniverse

	txn := db.bdb.NewTransactionAt(at, false)
	it := txn.NewIterator(badger.DefaultIteratorOptions)
	defer it.Close()

	prefix := []byte(OChainUniversePrefix)

	for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
		item := it.Item()

		var universe types.OChainUniverse
		value, err := item.ValueCopy(nil)
		if err != nil {
			return []types.OChainUniverse{}, err
		}

		err = cbor.Unmarshal(value, &universe)
		if err != nil {
			return []types.OChainUniverse{}, err
		}

		universes = append(universes, universe)
	}

	return universes, nil
}

func NewOChainUniverseTable(db *badger.DB) *OChainUniverseTable {
	return &OChainUniverseTable{
		bdb: db,
	}
}
