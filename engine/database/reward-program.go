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
	OChainRewardProgramPrefix string = "reward_programs_"
)

type OChainRewardProgramTable struct {
	bdb        *badger.DB
	currentTxn *badger.Txn
}

func (db *OChainRewardProgramTable) SetCurrentTxn(tx *badger.Txn) {
	db.currentTxn = tx
}

func (db *OChainRewardProgramTable) Exists(id uint64) (bool, error) {
	var at uint64 = math.MaxUint64
	return db.ExistsAt(id, at)
}

func (db *OChainRewardProgramTable) ExistsAt(id uint64, at uint64) (bool, error) {
	key := []byte(OChainRewardProgramPrefix + fmt.Sprint(id))
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

func (db *OChainRewardProgramTable) Get(id string) (types.OChainRewardProgram, error) {
	var at uint64 = math.MaxUint64
	return db.GetAt(id, at)
}

func (db *OChainRewardProgramTable) GetAt(id string, at uint64) (types.OChainRewardProgram, error) {
	var epoch types.OChainRewardProgram
	key := []byte(OChainRewardProgramPrefix + fmt.Sprint(id))
	txn := db.bdb.NewTransactionAt(at, false)

	item, err := txn.Get([]byte(key))
	if err != nil {
		return types.OChainRewardProgram{}, err
	}

	value, err := item.ValueCopy(nil)
	if err != nil {
		return types.OChainRewardProgram{}, err
	}

	err = cbor.Unmarshal(value, &epoch)
	if err != nil {
		return types.OChainRewardProgram{}, err
	}

	return epoch, nil
}

func (db *OChainRewardProgramTable) Insert(epoch types.OChainRewardProgram) error {
	key := []byte(OChainRewardProgramPrefix + fmt.Sprint(epoch.Id))

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

func (db *OChainRewardProgramTable) Update(epoch types.OChainRewardProgram) error {
	key := []byte(OChainRewardProgramPrefix + fmt.Sprint(epoch.Id))

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

func (db *OChainRewardProgramTable) Upsert(epoch types.OChainRewardProgram) error {
	key := []byte(OChainRewardProgramPrefix + fmt.Sprint(epoch.Id))
	value, err := cbor.Marshal(epoch)
	if err != nil {
		return err
	}

	return db.currentTxn.Set(key, value)
}

func (db *OChainRewardProgramTable) Delete(id string) error {
	key := []byte(OChainRewardProgramPrefix + fmt.Sprint(id))
	return db.currentTxn.Delete(key)
}

func (db *OChainRewardProgramTable) GetAll() ([]types.OChainRewardProgram, error) {
	var at uint64 = math.MaxUint64
	return db.GetAllAt(at)
}

func (db *OChainRewardProgramTable) GetAllAt(at uint64) ([]types.OChainRewardProgram, error) {
	var epochs []types.OChainRewardProgram

	txn := db.bdb.NewTransactionAt(at, false)
	it := txn.NewIterator(badger.DefaultIteratorOptions)
	defer it.Close()

	prefix := []byte(OChainRewardProgramPrefix)

	for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
		item := it.Item()

		var epoch types.OChainRewardProgram
		value, err := item.ValueCopy(nil)
		if err != nil {
			return []types.OChainRewardProgram{}, err
		}

		err = cbor.Unmarshal(value, &epoch)
		if err != nil {
			return []types.OChainRewardProgram{}, err
		}

		epochs = append(epochs, epoch)
	}

	return epochs, nil
}

func NewOChainRewardProgramTable(db *badger.DB) *OChainRewardProgramTable {
	return &OChainRewardProgramTable{
		bdb: db,
	}
}
