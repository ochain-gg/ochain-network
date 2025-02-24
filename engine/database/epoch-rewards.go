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
	OChainEpochValidatorRewardsPrefix string = "epoch_rewards_"
)

type OChainEpochValidatorRewardsTable struct {
	bdb        *badger.DB
	currentTxn *badger.Txn
}

func (db *OChainEpochValidatorRewardsTable) SetCurrentTxn(tx *badger.Txn) {
	db.currentTxn = tx
}

func (db *OChainEpochValidatorRewardsTable) Exists(id string) (bool, error) {
	var at uint64 = math.MaxUint64
	return db.ExistsAt(id, at)
}

func (db *OChainEpochValidatorRewardsTable) ExistsAt(id string, at uint64) (bool, error) {
	key := []byte(OChainEpochValidatorRewardsPrefix + fmt.Sprint(id))
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

func (db *OChainEpochValidatorRewardsTable) Get(id string) (types.OChainEpochValidatorRewards, error) {
	var at uint64 = math.MaxUint64
	return db.GetAt(id, at)
}

func (db *OChainEpochValidatorRewardsTable) GetAt(id string, at uint64) (types.OChainEpochValidatorRewards, error) {
	var epoch types.OChainEpochValidatorRewards
	key := []byte(OChainEpochValidatorRewardsPrefix + fmt.Sprint(id))
	txn := db.bdb.NewTransactionAt(at, false)

	item, err := txn.Get([]byte(key))
	if err != nil {
		return types.OChainEpochValidatorRewards{}, err
	}

	value, err := item.ValueCopy(nil)
	if err != nil {
		return types.OChainEpochValidatorRewards{}, err
	}

	err = cbor.Unmarshal(value, &epoch)
	if err != nil {
		return types.OChainEpochValidatorRewards{}, err
	}

	return epoch, nil
}

func (db *OChainEpochValidatorRewardsTable) Insert(epoch types.OChainEpochValidatorRewards) error {
	key := []byte(OChainEpochValidatorRewardsPrefix + fmt.Sprint(epoch.Id))

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

func (db *OChainEpochValidatorRewardsTable) Update(epoch types.OChainEpochValidatorRewards) error {
	key := []byte(OChainEpochValidatorRewardsPrefix + fmt.Sprint(epoch.Id))

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

func (db *OChainEpochValidatorRewardsTable) Upsert(epoch types.OChainEpochValidatorRewards) error {
	key := []byte(OChainEpochValidatorRewardsPrefix + fmt.Sprint(epoch.Id))
	value, err := cbor.Marshal(epoch)
	if err != nil {
		return err
	}

	return db.currentTxn.Set(key, value)
}

func (db *OChainEpochValidatorRewardsTable) Delete(id string) error {
	key := []byte(OChainEpochValidatorRewardsPrefix + fmt.Sprint(id))
	return db.currentTxn.Delete(key)
}

func (db *OChainEpochValidatorRewardsTable) GetAll() ([]types.OChainEpochValidatorRewards, error) {
	var at uint64 = math.MaxUint64
	return db.GetAllAt(at)
}

func (db *OChainEpochValidatorRewardsTable) GetAllAt(at uint64) ([]types.OChainEpochValidatorRewards, error) {
	var epochs []types.OChainEpochValidatorRewards

	txn := db.bdb.NewTransactionAt(at, false)
	it := txn.NewIterator(badger.DefaultIteratorOptions)
	defer it.Close()

	prefix := []byte(OChainEpochValidatorRewardsPrefix)

	for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
		item := it.Item()

		var epoch types.OChainEpochValidatorRewards
		value, err := item.ValueCopy(nil)
		if err != nil {
			return []types.OChainEpochValidatorRewards{}, err
		}

		err = cbor.Unmarshal(value, &epoch)
		if err != nil {
			return []types.OChainEpochValidatorRewards{}, err
		}

		epochs = append(epochs, epoch)
	}

	return epochs, nil
}

func NewOChainEpochValidatorRewardsTable(db *badger.DB) *OChainEpochValidatorRewardsTable {
	return &OChainEpochValidatorRewardsTable{
		bdb: db,
	}
}
