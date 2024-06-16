package database

import (
	"errors"
	"math"

	"github.com/dgraph-io/badger/v4"
	"github.com/fxamacker/cbor/v2"
	"github.com/ochain-gg/ochain-network/types"
)

const (
	OChainUniverseAccountWeeklyUsagePrefix string = "weekly_usage_"
)

type OChainUniverseAccountWeeklyUsageTable struct {
	bdb        *badger.DB
	currentTxn *badger.Txn
}

func (db *OChainUniverseAccountWeeklyUsageTable) SetCurrentTxn(tx *badger.Txn) {
	db.currentTxn = tx
}

func (db *OChainUniverseAccountWeeklyUsageTable) Exists(universeId string, address string) (bool, error) {
	var at uint64 = math.MaxUint64
	return db.ExistsAt(universeId, address, at)
}

func (db *OChainUniverseAccountWeeklyUsageTable) ExistsAt(universeId string, address string, at uint64) (bool, error) {
	key := []byte(OChainUniverseAccountWeeklyUsagePrefix + universeId + "_" + address)
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

func (db *OChainUniverseAccountWeeklyUsageTable) Get(universeId string, address string, year int, week int) (types.OChainUniverseAccountWeeklyUsage, error) {
	var at uint64 = math.MaxUint64
	return db.GetAt(universeId, address, year, week, at)
}

func (db *OChainUniverseAccountWeeklyUsageTable) GetAt(universeId string, address string, year int, week int, at uint64) (types.OChainUniverseAccountWeeklyUsage, error) {
	var account types.OChainUniverseAccountWeeklyUsage
	key := []byte(OChainUniverseAccountWeeklyUsagePrefix + universeId + "_" + address)
	txn := db.bdb.NewTransactionAt(at, false)

	item, err := txn.Get([]byte(key))
	if err != nil {
		if err == badger.ErrKeyNotFound {
			return types.OChainUniverseAccountWeeklyUsage{
				UniverseId: universeId,
				Address:    address,
				Year:       year,
				Week:       week,

				WithdrawalsExecuted: 0,
				DepositedAmount:     0,
			}, nil
		} else {
			return types.OChainUniverseAccountWeeklyUsage{}, err
		}
	}

	value, err := item.ValueCopy(nil)
	if err != nil {
		return types.OChainUniverseAccountWeeklyUsage{}, err
	}

	err = cbor.Unmarshal(value, &account)
	if err != nil {
		return types.OChainUniverseAccountWeeklyUsage{}, err
	}

	return account, nil
}

func (db *OChainUniverseAccountWeeklyUsageTable) Insert(account types.OChainUniverseAccountWeeklyUsage) error {
	key := []byte(OChainUniverseAccountWeeklyUsagePrefix + account.UniverseId + "_" + account.Address)

	exists, err := db.ExistsAt(account.UniverseId, account.Address, db.currentTxn.ReadTs())
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

func (db *OChainUniverseAccountWeeklyUsageTable) Update(account types.OChainUniverseAccountWeeklyUsage) error {
	key := []byte(OChainUniverseAccountWeeklyUsagePrefix + account.UniverseId + "_" + account.Address)

	exists, err := db.ExistsAt(account.UniverseId, account.Address, db.currentTxn.ReadTs())
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

func (db *OChainUniverseAccountWeeklyUsageTable) Upsert(account types.OChainUniverseAccountWeeklyUsage) error {
	key := []byte(OChainUniverseAccountWeeklyUsagePrefix + account.UniverseId + "_" + account.Address)
	value, err := cbor.Marshal(account)
	if err != nil {
		return err
	}

	return db.currentTxn.Set(key, value)
}

func (db *OChainUniverseAccountWeeklyUsageTable) Delete(universeId string, address string) error {
	key := []byte(OChainUniverseAccountWeeklyUsagePrefix + address)
	return db.currentTxn.Delete(key)
}

func (db *OChainUniverseAccountWeeklyUsageTable) GetAll() ([]types.OChainUniverseAccountWeeklyUsage, error) {
	var at uint64 = math.MaxUint64
	return db.GetAllAt(at)
}

func (db *OChainUniverseAccountWeeklyUsageTable) GetAllAt(at uint64) ([]types.OChainUniverseAccountWeeklyUsage, error) {
	var accounts []types.OChainUniverseAccountWeeklyUsage

	txn := db.bdb.NewTransactionAt(at, false)
	it := txn.NewIterator(badger.DefaultIteratorOptions)
	defer it.Close()

	prefix := []byte(OChainUniverseAccountWeeklyUsagePrefix)

	for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
		item := it.Item()

		var account types.OChainUniverseAccountWeeklyUsage
		value, err := item.ValueCopy(nil)
		if err != nil {
			return []types.OChainUniverseAccountWeeklyUsage{}, err
		}

		err = cbor.Unmarshal(value, &account)
		if err != nil {
			return []types.OChainUniverseAccountWeeklyUsage{}, err
		}

		accounts = append(accounts, account)
	}

	return accounts, nil
}

func NewOChainUniverseAccountWeeklyUsageTable(db *badger.DB) *OChainUniverseAccountWeeklyUsageTable {
	return &OChainUniverseAccountWeeklyUsageTable{
		bdb: db,
	}
}
