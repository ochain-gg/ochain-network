package database

import (
	"errors"
	"fmt"
	"math"

	"github.com/dgraph-io/badger/v4"
	"github.com/fxamacker/cbor/v2"
	"github.com/ochain-gg/ochain-network/types"
)

const (
	OChainEpochPrefix string = "epoch_"
)

type OChainEpochTable struct {
	bdb        *badger.DB
	currentTxn *badger.Txn
}

func (db *OChainEpochTable) SetCurrentTxn(tx *badger.Txn) {
	db.currentTxn = tx
}

func (db *OChainEpochTable) Exists(id uint64) (bool, error) {
	var at uint64 = math.MaxUint64
	return db.ExistsAt(id, at)
}

func (db *OChainEpochTable) ExistsAt(id uint64, at uint64) (bool, error) {
	key := []byte(OChainEpochPrefix + fmt.Sprint(id))
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

func (db *OChainEpochTable) Get(id uint64) (types.OChainEpoch, error) {
	var at uint64 = math.MaxUint64
	return db.GetAt(id, at)
}

func (db *OChainEpochTable) GetAt(id uint64, at uint64) (types.OChainEpoch, error) {
	var epoch types.OChainEpoch
	key := []byte(OChainEpochPrefix + fmt.Sprint(id))
	txn := db.bdb.NewTransactionAt(at, false)

	item, err := txn.Get([]byte(key))
	if err != nil {
		return types.OChainEpoch{}, err
	}

	value, err := item.ValueCopy(nil)
	if err != nil {
		return types.OChainEpoch{}, err
	}

	err = cbor.Unmarshal(value, &epoch)
	if err != nil {
		return types.OChainEpoch{}, err
	}

	return epoch, nil
}

func (db *OChainEpochTable) Insert(epoch types.OChainEpoch) error {
	key := []byte(OChainEpochPrefix + fmt.Sprint(epoch.Id))

	exists, err := db.ExistsAt(epoch.Id, db.currentTxn.ReadTs())
	if err != nil {
		return err
	}

	if exists {
		return errors.New("global epoch already exists")
	}

	value, err := cbor.Marshal(epoch)
	if err != nil {
		return err
	}

	return db.currentTxn.Set(key, value)
}

func (db *OChainEpochTable) Update(epoch types.OChainEpoch) error {
	key := []byte(OChainEpochPrefix + fmt.Sprint(epoch.Id))

	exists, err := db.ExistsAt(epoch.Id, db.currentTxn.ReadTs())
	if err != nil {
		return err
	}

	if !exists {
		return errors.New("epoch doesn't exists")
	}

	value, err := cbor.Marshal(epoch)
	if err != nil {
		return err
	}

	return db.currentTxn.Set(key, value)
}

func (db *OChainEpochTable) Upsert(epoch types.OChainEpoch) error {
	key := []byte(OChainEpochPrefix + fmt.Sprint(epoch.Id))
	value, err := cbor.Marshal(epoch)
	if err != nil {
		return err
	}

	return db.currentTxn.Set(key, value)
}

func (db *OChainEpochTable) Delete(id uint64) error {
	key := []byte(OChainEpochPrefix + fmt.Sprint(id))
	return db.currentTxn.Delete(key)
}

func (db *OChainEpochTable) GetAll() ([]types.OChainEpoch, error) {
	var at uint64 = math.MaxUint64
	return db.GetAllAt(at)
}

func (db *OChainEpochTable) GetAllAt(at uint64) ([]types.OChainEpoch, error) {
	var epochs []types.OChainEpoch

	txn := db.bdb.NewTransactionAt(at, false)
	it := txn.NewIterator(badger.DefaultIteratorOptions)
	defer it.Close()

	prefix := []byte(OChainEpochPrefix)

	for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
		item := it.Item()

		var epoch types.OChainEpoch
		value, err := item.ValueCopy(nil)
		if err != nil {
			return []types.OChainEpoch{}, err
		}

		err = cbor.Unmarshal(value, &epoch)
		if err != nil {
			return []types.OChainEpoch{}, err
		}

		epochs = append(epochs, epoch)
	}

	return epochs, nil
}

func NewOChainEpochTable(db *badger.DB) *OChainEpochTable {
	return &OChainEpochTable{
		bdb: db,
	}
}
