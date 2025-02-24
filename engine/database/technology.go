package database

import (
	"errors"
	"math"

	"github.com/dgraph-io/badger/v4"
	"github.com/fxamacker/cbor/v2"
	"github.com/ochain-gg/ochain-network/types"
)

const (
	OChainTechnologyPrefix string = "technology_"
)

type OChainTechnologyTable struct {
	bdb        *badger.DB
	currentTxn *badger.Txn
}

func (db *OChainTechnologyTable) SetCurrentTxn(tx *badger.Txn) {
	db.currentTxn = tx
}

func (db *OChainTechnologyTable) Exists(id string) (bool, error) {
	var at uint64 = math.MaxUint64
	return db.ExistsAt(id, at)
}

func (db *OChainTechnologyTable) ExistsAt(id string, at uint64) (bool, error) {
	key := []byte(OChainTechnologyPrefix + id)
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

func (db *OChainTechnologyTable) Get(id string) (types.OChainTechnology, error) {
	var at uint64 = math.MaxUint64
	return db.GetAt(id, at)
}

func (db *OChainTechnologyTable) GetAt(id string, at uint64) (types.OChainTechnology, error) {
	var technology types.OChainTechnology
	key := []byte(OChainTechnologyPrefix + id)
	txn := db.bdb.NewTransactionAt(at, false)

	item, err := txn.Get([]byte(key))
	if err != nil {
		return types.OChainTechnology{}, err
	}

	value, err := item.ValueCopy(nil)
	if err != nil {
		return types.OChainTechnology{}, err
	}

	err = cbor.Unmarshal(value, &technology)
	if err != nil {
		return types.OChainTechnology{}, err
	}

	return technology, nil
}

func (db *OChainTechnologyTable) Insert(technology types.OChainTechnology) error {
	key := []byte(OChainTechnologyPrefix + string(technology.Id))

	exists, err := db.ExistsAt(string(technology.Id), db.currentTxn.ReadTs())
	if err != nil {
		return err
	}

	if exists {
		return errors.New("global technology already exists")
	}

	value, err := cbor.Marshal(technology)
	if err != nil {
		return err
	}

	return db.currentTxn.Set(key, value)
}

func (db *OChainTechnologyTable) Update(technology types.OChainTechnology) error {
	key := []byte(OChainTechnologyPrefix + string(technology.Id))

	exists, err := db.ExistsAt(string(technology.Id), db.currentTxn.ReadTs())
	if err != nil {
		return err
	}

	if !exists {
		return errors.New("technology doesn't exists")
	}

	value, err := cbor.Marshal(technology)
	if err != nil {
		return err
	}

	return db.currentTxn.Set(key, value)
}

func (db *OChainTechnologyTable) Upsert(technology types.OChainTechnology) error {
	key := []byte(OChainTechnologyPrefix + string(technology.Id))
	value, err := cbor.Marshal(technology)
	if err != nil {
		return err
	}

	return db.currentTxn.Set(key, value)
}

func (db *OChainTechnologyTable) Delete(id string) error {
	key := []byte(OChainTechnologyPrefix + id)
	return db.currentTxn.Delete(key)
}

func (db *OChainTechnologyTable) GetAll() ([]types.OChainTechnology, error) {
	var at uint64 = math.MaxUint64
	return db.GetAllAt(at)
}

func (db *OChainTechnologyTable) GetAllAt(at uint64) ([]types.OChainTechnology, error) {
	var technologys []types.OChainTechnology

	txn := db.bdb.NewTransactionAt(at, false)
	it := txn.NewIterator(badger.DefaultIteratorOptions)
	defer it.Close()

	prefix := []byte(OChainTechnologyPrefix)

	for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
		item := it.Item()

		var technology types.OChainTechnology
		value, err := item.ValueCopy(nil)
		if err != nil {
			return []types.OChainTechnology{}, err
		}

		err = cbor.Unmarshal(value, &technology)
		if err != nil {
			return []types.OChainTechnology{}, err
		}

		technologys = append(technologys, technology)
	}

	return technologys, nil
}

func NewOChainTechnologyTable(db *badger.DB) *OChainTechnologyTable {
	return &OChainTechnologyTable{
		bdb: db,
	}
}
