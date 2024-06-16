package database

import (
	"errors"
	"math"

	"github.com/dgraph-io/badger/v4"
	"github.com/fxamacker/cbor/v2"
	"github.com/ochain-gg/ochain-network/types"
)

const (
	OChainResourcesMarketPrefix string = "universe_account_"
)

type OChainResourcesMarketTable struct {
	bdb        *badger.DB
	currentTxn *badger.Txn
}

func (db *OChainResourcesMarketTable) SetCurrentTxn(tx *badger.Txn) {
	db.currentTxn = tx
}

func (db *OChainResourcesMarketTable) Exists(universeId string) (bool, error) {
	var at uint64 = math.MaxUint64
	return db.ExistsAt(universeId, at)
}

func (db *OChainResourcesMarketTable) ExistsAt(universeId string, at uint64) (bool, error) {
	key := []byte(OChainResourcesMarketPrefix + universeId)
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

func (db *OChainResourcesMarketTable) Get(universeId string) (types.OChainResourcesMarket, error) {
	var at uint64 = math.MaxUint64
	return db.GetAt(universeId, at)
}

func (db *OChainResourcesMarketTable) GetAt(universeId string, at uint64) (types.OChainResourcesMarket, error) {
	var account types.OChainResourcesMarket
	key := []byte(OChainResourcesMarketPrefix + universeId)
	txn := db.bdb.NewTransactionAt(at, false)

	item, err := txn.Get([]byte(key))
	if err != nil {
		return types.OChainResourcesMarket{}, err
	}

	value, err := item.ValueCopy(nil)
	if err != nil {
		return types.OChainResourcesMarket{}, err
	}

	err = cbor.Unmarshal(value, &account)
	if err != nil {
		return types.OChainResourcesMarket{}, err
	}

	return account, nil
}

func (db *OChainResourcesMarketTable) Insert(account types.OChainResourcesMarket) error {
	key := []byte(OChainResourcesMarketPrefix + account.UniverseId)

	exists, err := db.ExistsAt(account.UniverseId, db.currentTxn.ReadTs())
	if err != nil {
		return err
	}

	if exists {
		return errors.New("resources market already exists")
	}

	value, err := cbor.Marshal(account)
	if err != nil {
		return err
	}

	return db.currentTxn.Set(key, value)
}

func (db *OChainResourcesMarketTable) Update(account types.OChainResourcesMarket) error {
	key := []byte(OChainResourcesMarketPrefix + account.UniverseId)

	exists, err := db.ExistsAt(account.UniverseId, db.currentTxn.ReadTs())
	if err != nil {
		return err
	}

	if !exists {
		return errors.New("resources market doesn't exists")
	}

	value, err := cbor.Marshal(account)
	if err != nil {
		return err
	}

	return db.currentTxn.Set(key, value)
}

func (db *OChainResourcesMarketTable) Upsert(account types.OChainResourcesMarket) error {
	key := []byte(OChainResourcesMarketPrefix + account.UniverseId)
	value, err := cbor.Marshal(account)
	if err != nil {
		return err
	}

	return db.currentTxn.Set(key, value)
}

func (db *OChainResourcesMarketTable) Delete(universeId string) error {
	key := []byte(OChainResourcesMarketPrefix + universeId)
	return db.currentTxn.Delete(key)
}

func (db *OChainResourcesMarketTable) GetAll() ([]types.OChainResourcesMarket, error) {
	var at uint64 = math.MaxUint64
	return db.GetAllAt(at)
}

func (db *OChainResourcesMarketTable) GetAllAt(at uint64) ([]types.OChainResourcesMarket, error) {
	var accounts []types.OChainResourcesMarket

	txn := db.bdb.NewTransactionAt(at, false)
	it := txn.NewIterator(badger.DefaultIteratorOptions)
	defer it.Close()

	prefix := []byte(OChainResourcesMarketPrefix)

	for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
		item := it.Item()

		var account types.OChainResourcesMarket
		value, err := item.ValueCopy(nil)
		if err != nil {
			return []types.OChainResourcesMarket{}, err
		}

		err = cbor.Unmarshal(value, &account)
		if err != nil {
			return []types.OChainResourcesMarket{}, err
		}

		accounts = append(accounts, account)
	}

	return accounts, nil
}

func NewOChainResourcesMarketTable(db *badger.DB) *OChainResourcesMarketTable {
	return &OChainResourcesMarketTable{
		bdb: db,
	}
}
