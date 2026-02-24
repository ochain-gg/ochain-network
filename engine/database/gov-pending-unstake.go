package database

import (
	"errors"
	"math"

	"github.com/dgraph-io/badger/v4"
	"github.com/fxamacker/cbor/v2"
	"github.com/ochain-gg/ochain-network/types"
)

const (
	OChainPendingUnstakePrefix string = "gov_pending_unstake_"
)

type OChainPendingUnstakeTable struct {
	bdb        *badger.DB
	currentTxn *badger.Txn
}

func (db *OChainPendingUnstakeTable) SetCurrentTxn(tx *badger.Txn) {
	db.currentTxn = tx
}

func (db *OChainPendingUnstakeTable) Exists(address string) (bool, error) {
	var at uint64 = math.MaxUint64
	return db.ExistsAt(address, at)
}

func (db *OChainPendingUnstakeTable) ExistsAt(address string, at uint64) (bool, error) {
	key := []byte(OChainPendingUnstakePrefix + address)
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

func (db *OChainPendingUnstakeTable) Get(address string) (types.OChainPendingUnstake, error) {
	var at uint64 = math.MaxUint64
	return db.GetAt(address, at)
}

func (db *OChainPendingUnstakeTable) GetAt(address string, at uint64) (types.OChainPendingUnstake, error) {
	var pendingUnstake types.OChainPendingUnstake
	key := []byte(OChainPendingUnstakePrefix + address)
	txn := db.bdb.NewTransactionAt(at, false)

	item, err := txn.Get([]byte(key))
	if err != nil {
		return types.OChainPendingUnstake{}, err
	}

	value, err := item.ValueCopy(nil)
	if err != nil {
		return types.OChainPendingUnstake{}, err
	}

	err = cbor.Unmarshal(value, &pendingUnstake)
	if err != nil {
		return types.OChainPendingUnstake{}, err
	}

	return pendingUnstake, nil
}

func (db *OChainPendingUnstakeTable) Insert(pendingUnstake types.OChainPendingUnstake) error {
	key := []byte(OChainPendingUnstakePrefix + pendingUnstake.Address)

	exists, err := db.ExistsAt(pendingUnstake.Address, db.currentTxn.ReadTs())
	if err != nil {
		return err
	}

	if exists {
		return errors.New("pending unstake already exists")
	}

	value, err := cbor.Marshal(pendingUnstake)
	if err != nil {
		return err
	}

	return db.currentTxn.Set(key, value)
}

func (db *OChainPendingUnstakeTable) Update(pendingUnstake types.OChainPendingUnstake) error {
	key := []byte(OChainPendingUnstakePrefix + pendingUnstake.Address)

	exists, err := db.ExistsAt(pendingUnstake.Address, db.currentTxn.ReadTs())
	if err != nil {
		return err
	}

	if !exists {
		return errors.New("pending unstake doesn't exist")
	}

	value, err := cbor.Marshal(pendingUnstake)
	if err != nil {
		return err
	}

	return db.currentTxn.Set(key, value)
}

func (db *OChainPendingUnstakeTable) Upsert(pendingUnstake types.OChainPendingUnstake) error {
	key := []byte(OChainPendingUnstakePrefix + pendingUnstake.Address)
	value, err := cbor.Marshal(pendingUnstake)
	if err != nil {
		return err
	}

	return db.currentTxn.Set(key, value)
}

func (db *OChainPendingUnstakeTable) Delete(address string) error {
	key := []byte(OChainPendingUnstakePrefix + address)
	return db.currentTxn.Delete(key)
}

func (db *OChainPendingUnstakeTable) GetAll() ([]types.OChainPendingUnstake, error) {
	var at uint64 = math.MaxUint64
	return db.GetAllAt(at)
}

func (db *OChainPendingUnstakeTable) GetAllAt(at uint64) ([]types.OChainPendingUnstake, error) {
	var pendingUnstakes []types.OChainPendingUnstake

	txn := db.bdb.NewTransactionAt(at, false)
	it := txn.NewIterator(badger.DefaultIteratorOptions)
	defer it.Close()

	prefix := []byte(OChainPendingUnstakePrefix)

	for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
		item := it.Item()

		var pendingUnstake types.OChainPendingUnstake
		value, err := item.ValueCopy(nil)
		if err != nil {
			return []types.OChainPendingUnstake{}, err
		}

		err = cbor.Unmarshal(value, &pendingUnstake)
		if err != nil {
			return []types.OChainPendingUnstake{}, err
		}

		pendingUnstakes = append(pendingUnstakes, pendingUnstake)
	}

	return pendingUnstakes, nil
}

func NewOChainPendingUnstakeTable(db *badger.DB) *OChainPendingUnstakeTable {
	return &OChainPendingUnstakeTable{
		bdb: db,
	}
}
