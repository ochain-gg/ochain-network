package database

import (
	"errors"
	"math"

	"github.com/dgraph-io/badger/v4"
	"github.com/fxamacker/cbor/v2"
	"github.com/ochain-gg/ochain-network/types"
)

const (
	OChainTokenTransferPrefix string = "oct_transfer_"
)

type OChainTokenTransferTable struct {
	bdb        *badger.DB
	currentTxn *badger.Txn
}

func (db *OChainTokenTransferTable) SetCurrentTxn(tx *badger.Txn) {
	db.currentTxn = tx
}

func (db *OChainTokenTransferTable) Exists(address string) (bool, error) {
	var at uint64
	at = math.MaxUint64
	return db.ExistsAt(address, at)
}

func (db *OChainTokenTransferTable) ExistsAt(address string, at uint64) (bool, error) {
	key := []byte(OChainTokenTransferPrefix + address)
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

func (db *OChainTokenTransferTable) Get(address string) (types.OChainTokenTransfer, error) {
	var at uint64
	at = math.MaxUint64
	return db.GetAt(address, at)
}

func (db *OChainTokenTransferTable) GetAt(address string, at uint64) (types.OChainTokenTransfer, error) {
	var universe types.OChainTokenTransfer
	key := []byte(OChainTokenTransferPrefix + address)
	txn := db.bdb.NewTransactionAt(at, false)

	item, err := txn.Get([]byte(key))
	if err != nil {
		return types.OChainTokenTransfer{}, err
	}

	value, err := item.ValueCopy(nil)
	if err != nil {
		return types.OChainTokenTransfer{}, err
	}

	err = cbor.Unmarshal(value, &universe)
	if err != nil {
		return types.OChainTokenTransfer{}, err
	}

	return universe, nil
}

func (db *OChainTokenTransferTable) Insert(universe types.OChainTokenTransfer) error {
	key := []byte(OChainTokenTransferPrefix + universe.Id)

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

func (db *OChainTokenTransferTable) Update(universe types.OChainTokenTransfer) error {
	key := []byte(OChainTokenTransferPrefix + universe.Id)

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

func (db *OChainTokenTransferTable) Upsert(universe types.OChainTokenTransfer) error {
	key := []byte(OChainTokenTransferPrefix + universe.Id)
	value, err := cbor.Marshal(universe)
	if err != nil {
		return err
	}

	return db.currentTxn.Set(key, value)
}

func (db *OChainTokenTransferTable) Delete(id string) error {
	key := []byte(OChainTokenTransferPrefix + id)
	return db.currentTxn.Delete(key)
}

func (db *OChainTokenTransferTable) GetAll() ([]types.OChainTokenTransfer, error) {
	var at uint64
	at = math.MaxUint64
	return db.GetAllAt(at)
}

func (db *OChainTokenTransferTable) GetAllAt(at uint64) ([]types.OChainTokenTransfer, error) {
	var universes []types.OChainTokenTransfer

	txn := db.bdb.NewTransactionAt(at, false)
	it := txn.NewIterator(badger.DefaultIteratorOptions)
	defer it.Close()

	prefix := []byte(OChainTokenTransferPrefix)

	for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
		item := it.Item()

		var universe types.OChainTokenTransfer
		value, err := item.ValueCopy(nil)
		if err != nil {
			return []types.OChainTokenTransfer{}, err
		}

		err = cbor.Unmarshal(value, &universe)
		if err != nil {
			return []types.OChainTokenTransfer{}, err
		}

		universes = append(universes, universe)
	}

	return universes, nil
}

func NewOChainTokenTransferTable(db *badger.DB) *OChainTokenTransferTable {
	return &OChainTokenTransferTable{
		bdb: db,
	}
}
