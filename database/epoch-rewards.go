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
	OChainEpochRewardsPrefix string = "epoch_rewards_"
)

type OChainEpochRewardsTable struct {
	bdb        *badger.DB
	currentTxn *badger.Txn
}

func (db *OChainEpochRewardsTable) SetCurrentTxn(tx *badger.Txn) {
	db.currentTxn = tx
}

func (db *OChainEpochRewardsTable) Exists(id string) (bool, error) {
	var at uint64 = math.MaxUint64
	return db.ExistsAt(id, at)
}

func (db *OChainEpochRewardsTable) ExistsAt(id string, at uint64) (bool, error) {
	key := []byte(OChainEpochRewardsPrefix + fmt.Sprint(id))
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

func (db *OChainEpochRewardsTable) Get(id string) (types.OChainEpochRewards, error) {
	var at uint64 = math.MaxUint64
	return db.GetAt(id, at)
}

func (db *OChainEpochRewardsTable) GetAt(id string, at uint64) (types.OChainEpochRewards, error) {
	var epoch types.OChainEpochRewards
	key := []byte(OChainEpochRewardsPrefix + fmt.Sprint(id))
	txn := db.bdb.NewTransactionAt(at, false)

	item, err := txn.Get([]byte(key))
	if err != nil {
		return types.OChainEpochRewards{}, err
	}

	value, err := item.ValueCopy(nil)
	if err != nil {
		return types.OChainEpochRewards{}, err
	}

	err = cbor.Unmarshal(value, &epoch)
	if err != nil {
		return types.OChainEpochRewards{}, err
	}

	return epoch, nil
}

func (db *OChainEpochRewardsTable) Insert(epoch types.OChainEpochRewards) error {
	key := []byte(OChainEpochRewardsPrefix + fmt.Sprint(epoch.Id))

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

func (db *OChainEpochRewardsTable) Update(epoch types.OChainEpochRewards) error {
	key := []byte(OChainEpochRewardsPrefix + fmt.Sprint(epoch.Id))

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

func (db *OChainEpochRewardsTable) Upsert(epoch types.OChainEpochRewards) error {
	key := []byte(OChainEpochRewardsPrefix + fmt.Sprint(epoch.Id))
	value, err := cbor.Marshal(epoch)
	if err != nil {
		return err
	}

	return db.currentTxn.Set(key, value)
}

func (db *OChainEpochRewardsTable) Delete(id string) error {
	key := []byte(OChainEpochRewardsPrefix + fmt.Sprint(id))
	return db.currentTxn.Delete(key)
}

func (db *OChainEpochRewardsTable) GetAll() ([]types.OChainEpochRewards, error) {
	var at uint64 = math.MaxUint64
	return db.GetAllAt(at)
}

func (db *OChainEpochRewardsTable) GetAllAt(at uint64) ([]types.OChainEpochRewards, error) {
	var epochs []types.OChainEpochRewards

	txn := db.bdb.NewTransactionAt(at, false)
	it := txn.NewIterator(badger.DefaultIteratorOptions)
	defer it.Close()

	prefix := []byte(OChainEpochRewardsPrefix)

	for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
		item := it.Item()

		var epoch types.OChainEpochRewards
		value, err := item.ValueCopy(nil)
		if err != nil {
			return []types.OChainEpochRewards{}, err
		}

		err = cbor.Unmarshal(value, &epoch)
		if err != nil {
			return []types.OChainEpochRewards{}, err
		}

		epochs = append(epochs, epoch)
	}

	return epochs, nil
}

func NewOChainEpochRewardsTable(db *badger.DB) *OChainEpochRewardsTable {
	return &OChainEpochRewardsTable{
		bdb: db,
	}
}
