package database

import (
	"errors"
	"math"

	"github.com/dgraph-io/badger/v4"
	"github.com/fxamacker/cbor/v2"
	"github.com/ochain-gg/ochain-network/types"
)

const (
	OChainGlobalAccountPrefix string = "global_account_"
)

type OChainGlobalAccountTable struct {
	bdb        *badger.DB
	currentTxn *badger.Txn
}

func (db *OChainGlobalAccountTable) SetCurrentTxn(tx *badger.Txn) {
	db.currentTxn = tx
}

func (db *OChainGlobalAccountTable) Exists(address string) (bool, error) {
	var at uint64 = math.MaxUint64
	return db.ExistsAt(address, at)
}

func (db *OChainGlobalAccountTable) ExistsAt(address string, at uint64) (bool, error) {
	key := []byte(OChainGlobalAccountPrefix + address)
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

func (db *OChainGlobalAccountTable) Get(address string) (types.OChainGlobalAccount, error) {
	var at uint64 = math.MaxUint64
	return db.GetAt(address, at)
}

func (db *OChainGlobalAccountTable) GetAt(address string, at uint64) (types.OChainGlobalAccount, error) {
	var account types.OChainGlobalAccount
	key := []byte(OChainGlobalAccountPrefix + address)
	txn := db.bdb.NewTransactionAt(at, false)

	item, err := txn.Get([]byte(key))
	if err != nil {
		return types.OChainGlobalAccount{}, err
	}

	value, err := item.ValueCopy(nil)
	if err != nil {
		return types.OChainGlobalAccount{}, err
	}

	err = cbor.Unmarshal(value, &account)
	if err != nil {
		return types.OChainGlobalAccount{}, err
	}

	return account, nil
}

func (db *OChainGlobalAccountTable) Insert(account types.OChainGlobalAccount) error {
	key := []byte(OChainGlobalAccountPrefix + account.Address)

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

func (db *OChainGlobalAccountTable) Update(account types.OChainGlobalAccount) error {
	key := []byte(OChainGlobalAccountPrefix + account.Address)

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

func (db *OChainGlobalAccountTable) Upsert(account types.OChainGlobalAccount) error {
	key := []byte(OChainGlobalAccountPrefix + account.Address)
	value, err := cbor.Marshal(account)
	if err != nil {
		return err
	}

	return db.currentTxn.Set(key, value)
}

func (db *OChainGlobalAccountTable) Delete(address string) error {
	key := []byte(OChainGlobalAccountPrefix + address)
	return db.currentTxn.Delete(key)
}

func (db *OChainGlobalAccountTable) GetAll() ([]types.OChainGlobalAccount, error) {
	var at uint64 = math.MaxUint64
	return db.GetAllAt(at)
}

func (db *OChainGlobalAccountTable) GetAllAt(at uint64) ([]types.OChainGlobalAccount, error) {
	var accounts []types.OChainGlobalAccount

	txn := db.bdb.NewTransactionAt(at, false)
	it := txn.NewIterator(badger.DefaultIteratorOptions)
	defer it.Close()

	prefix := []byte(OChainGlobalAccountPrefix)

	for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
		item := it.Item()

		var account types.OChainGlobalAccount
		value, err := item.ValueCopy(nil)
		if err != nil {
			return []types.OChainGlobalAccount{}, err
		}

		err = cbor.Unmarshal(value, &account)
		if err != nil {
			return []types.OChainGlobalAccount{}, err
		}

		accounts = append(accounts, account)
	}

	return accounts, nil
}

func NewOChainGlobalAccountTable(db *badger.DB) *OChainGlobalAccountTable {
	return &OChainGlobalAccountTable{
		bdb: db,
	}
}
