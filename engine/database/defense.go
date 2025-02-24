package database

import (
	"errors"
	"math"

	"github.com/dgraph-io/badger/v4"
	"github.com/fxamacker/cbor/v2"
	"github.com/ochain-gg/ochain-network/types"
)

const (
	OChainDefensePrefix string = "defense_"
)

type OChainDefenseTable struct {
	bdb        *badger.DB
	currentTxn *badger.Txn
}

func (db *OChainDefenseTable) SetCurrentTxn(tx *badger.Txn) {
	db.currentTxn = tx
}

func (db *OChainDefenseTable) Exists(id string) (bool, error) {
	var at uint64 = math.MaxUint64
	return db.ExistsAt(id, at)
}

func (db *OChainDefenseTable) ExistsAt(id string, at uint64) (bool, error) {
	key := []byte(OChainDefensePrefix + id)
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

func (db *OChainDefenseTable) Get(id string) (types.OChainDefense, error) {
	var at uint64 = math.MaxUint64
	return db.GetAt(id, at)
}

func (db *OChainDefenseTable) GetAt(id string, at uint64) (types.OChainDefense, error) {
	var defense types.OChainDefense
	key := []byte(OChainDefensePrefix + id)
	txn := db.bdb.NewTransactionAt(at, false)

	item, err := txn.Get([]byte(key))
	if err != nil {
		return types.OChainDefense{}, err
	}

	value, err := item.ValueCopy(nil)
	if err != nil {
		return types.OChainDefense{}, err
	}

	err = cbor.Unmarshal(value, &defense)
	if err != nil {
		return types.OChainDefense{}, err
	}

	return defense, nil
}

func (db *OChainDefenseTable) Insert(defense types.OChainDefense) error {
	key := []byte(OChainDefensePrefix + string(defense.Id))

	exists, err := db.ExistsAt(string(defense.Id), db.currentTxn.ReadTs())
	if err != nil {
		return err
	}

	if exists {
		return errors.New("defense already exists")
	}

	value, err := cbor.Marshal(defense)
	if err != nil {
		return err
	}

	return db.currentTxn.Set(key, value)
}

func (db *OChainDefenseTable) Update(defense types.OChainDefense) error {
	key := []byte(OChainDefensePrefix + string(defense.Id))

	exists, err := db.ExistsAt(string(defense.Id), db.currentTxn.ReadTs())
	if err != nil {
		return err
	}

	if !exists {
		return errors.New("defense doesn't exists")
	}

	value, err := cbor.Marshal(defense)
	if err != nil {
		return err
	}

	return db.currentTxn.Set(key, value)
}

func (db *OChainDefenseTable) Upsert(defense types.OChainDefense) error {
	key := []byte(OChainDefensePrefix + string(defense.Id))
	value, err := cbor.Marshal(defense)
	if err != nil {
		return err
	}

	return db.currentTxn.Set(key, value)
}

func (db *OChainDefenseTable) Delete(id string) error {
	key := []byte(OChainDefensePrefix + id)
	return db.currentTxn.Delete(key)
}

func (db *OChainDefenseTable) GetAll() ([]types.OChainDefense, error) {
	var at uint64 = math.MaxUint64
	return db.GetAllAt(at)
}

func (db *OChainDefenseTable) GetAllAt(at uint64) ([]types.OChainDefense, error) {
	var defenses []types.OChainDefense

	txn := db.bdb.NewTransactionAt(at, false)
	it := txn.NewIterator(badger.DefaultIteratorOptions)
	defer it.Close()

	prefix := []byte(OChainDefensePrefix)

	for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
		item := it.Item()

		var defense types.OChainDefense
		value, err := item.ValueCopy(nil)
		if err != nil {
			return []types.OChainDefense{}, err
		}

		err = cbor.Unmarshal(value, &defense)
		if err != nil {
			return []types.OChainDefense{}, err
		}

		defenses = append(defenses, defense)
	}

	return defenses, nil
}

func NewOChainDefenseTable(db *badger.DB) *OChainDefenseTable {
	return &OChainDefenseTable{
		bdb: db,
	}
}
