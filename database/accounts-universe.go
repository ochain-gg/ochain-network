package database

import (
	"errors"
	"math"

	"github.com/dgraph-io/badger/v4"
	"github.com/fxamacker/cbor/v2"
	"github.com/ochain-gg/ochain-network/types"
)

const (
	OChainUniverseAccountPrefix string = "universe_account_"
)

type OChainUniverseAccountTable struct {
	bdb        *badger.DB
	currentTxn *badger.Txn
}

func (db *OChainUniverseAccountTable) SetCurrentTxn(tx *badger.Txn) {
	db.currentTxn = tx
}

func (db *OChainUniverseAccountTable) Exists(address string) (bool, error) {
	var at uint64
	at = math.MaxUint64
	return db.ExistsAt(address, at)
}

func (db *OChainUniverseAccountTable) ExistsAt(address string, at uint64) (bool, error) {
	key := []byte(OChainUniverseAccountPrefix + address)
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

func (db *OChainUniverseAccountTable) Get(address string) (types.OChainUniverseAccount, error) {
	var at uint64
	at = math.MaxUint64
	return db.GetAt(address, at)
}

func (db *OChainUniverseAccountTable) GetAt(address string, at uint64) (types.OChainUniverseAccount, error) {
	var account types.OChainUniverseAccount
	key := []byte(OChainUniverseAccountPrefix + address)
	txn := db.bdb.NewTransactionAt(at, false)

	item, err := txn.Get([]byte(key))
	if err != nil {
		return types.OChainUniverseAccount{}, err
	}

	value, err := item.ValueCopy(nil)
	if err != nil {
		return types.OChainUniverseAccount{}, err
	}

	err = cbor.Unmarshal(value, &account)
	if err != nil {
		return types.OChainUniverseAccount{}, err
	}

	return account, nil
}

func (db *OChainUniverseAccountTable) Insert(account types.OChainUniverseAccount) error {
	key := []byte(OChainUniverseAccountPrefix + account.Address)

	exists, err := db.ExistsAt(account.Address, db.currentTxn.ReadTs())
	if err != nil {
		return err
	}

	if exists {
		return errors.New("global account already exists")
	}

	value, err := cbor.Marshal(account)
	if err != nil {
		return err
	}

	return db.currentTxn.Set(key, value)
}

func (db *OChainUniverseAccountTable) Update(account types.OChainUniverseAccount) error {
	key := []byte(OChainUniverseAccountPrefix + account.Address)

	exists, err := db.ExistsAt(account.Address, db.currentTxn.ReadTs())
	if err != nil {
		return err
	}

	if !exists {
		return errors.New("global account doesn't exists")
	}

	value, err := cbor.Marshal(account)
	if err != nil {
		return err
	}

	return db.currentTxn.Set(key, value)
}

func (db *OChainUniverseAccountTable) Upsert(account types.OChainUniverseAccount) error {
	key := []byte(OChainUniverseAccountPrefix + account.Address)
	value, err := cbor.Marshal(account)
	if err != nil {
		return err
	}

	return db.currentTxn.Set(key, value)
}

func (db *OChainUniverseAccountTable) Delete(address string) error {
	key := []byte(OChainUniverseAccountPrefix + address)
	return db.currentTxn.Delete(key)
}

func (db *OChainUniverseAccountTable) GetAll() ([]types.OChainUniverseAccount, error) {
	var at uint64
	at = math.MaxUint64
	return db.GetAllAt(at)
}

func (db *OChainUniverseAccountTable) GetAllAt(at uint64) ([]types.OChainUniverseAccount, error) {
	var accounts []types.OChainUniverseAccount

	txn := db.bdb.NewTransactionAt(at, false)
	it := txn.NewIterator(badger.DefaultIteratorOptions)
	defer it.Close()

	prefix := []byte(OChainUniverseAccountPrefix)

	for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
		item := it.Item()

		var account types.OChainUniverseAccount
		value, err := item.ValueCopy(nil)
		if err != nil {
			return []types.OChainUniverseAccount{}, err
		}

		err = cbor.Unmarshal(value, &account)
		if err != nil {
			return []types.OChainUniverseAccount{}, err
		}

		accounts = append(accounts, account)
	}

	return accounts, nil
}

func NewOChainUniverseAccountTable(db *badger.DB) *OChainUniverseAccountTable {
	return &OChainUniverseAccountTable{
		bdb: db,
	}
}
