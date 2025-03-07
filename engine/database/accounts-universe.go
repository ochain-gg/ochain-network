package database

import (
	"context"
	"errors"
	"math"

	"github.com/dgraph-io/badger/pb"
	"github.com/dgraph-io/badger/v4"
	"github.com/dgraph-io/ristretto/v2/z"
	"github.com/fxamacker/cbor/v2"
	"github.com/ochain-gg/ochain-network/types"
)

const (
	OChainUniverseAccountPrefix string = "account_"
)

type OChainUniverseAccountTable struct {
	bdb        *badger.DB
	currentTxn *badger.Txn
}

func (db *OChainUniverseAccountTable) SetCurrentTxn(tx *badger.Txn) {
	db.currentTxn = tx
}

func (db *OChainUniverseAccountTable) Exists(universeId string, address string) (bool, error) {
	var at uint64 = math.MaxUint64
	return db.ExistsAt(universeId, address, at)
}

func (db *OChainUniverseAccountTable) ExistsAt(universeId string, address string, at uint64) (bool, error) {
	key := []byte(OChainUniverseAccountPrefix + universeId + "_" + address)
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

func (db *OChainUniverseAccountTable) Get(universeId string, address string) (types.OChainUniverseAccount, error) {
	var at uint64 = math.MaxUint64
	return db.GetAt(universeId, address, at)
}

func (db *OChainUniverseAccountTable) GetAt(universeId string, address string, at uint64) (types.OChainUniverseAccount, error) {
	var account types.OChainUniverseAccount
	key := []byte(OChainUniverseAccountPrefix + universeId + "_" + address)
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

func (db *OChainUniverseAccountTable) GetByAddress(address string) ([]types.OChainUniverseAccount, error) {
	var at uint64 = math.MaxUint64
	return db.GetByAddressAt(address, at)
}

func (db *OChainUniverseAccountTable) GetByAddressAt(address string, at uint64) ([]types.OChainUniverseAccount, error) {
	var accounts []types.OChainUniverseAccount

	stream := db.bdb.NewStreamAt(at)

	stream.NumGo = 16                                     // Set number of goroutines to use for iteration.
	stream.Prefix = []byte(OChainUniverseAccountPrefix)   // Leave nil for iteration over the whole DB.
	stream.LogPrefix = "Badger.UniverseAccount.Streaming" // For identifying stream logs. Outputs to Logger.
	stream.Send = func(buf *z.Buffer) error {
		err := buf.SliceIterate(func(s []byte) error {
			var kv pb.KV
			var acc types.OChainUniverseAccount
			if err := kv.Unmarshal(s); err != nil {
				return err
			}

			if err := cbor.Unmarshal(kv.Value, &acc); err != nil {
				return err
			}

			if acc.Address == address {
				accounts = append(accounts, acc)
			}

			return nil
		})

		return err
	}

	if err := stream.Orchestrate(context.Background()); err != nil {
		return []types.OChainUniverseAccount{}, err
	}

	return accounts, nil
}

func (db *OChainUniverseAccountTable) Insert(account types.OChainUniverseAccount) error {
	key := []byte(OChainUniverseAccountPrefix + account.UniverseId + "_" + account.Address)

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

func (db *OChainUniverseAccountTable) Update(account types.OChainUniverseAccount) error {
	key := []byte(OChainUniverseAccountPrefix + account.UniverseId + "_" + account.Address)

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

func (db *OChainUniverseAccountTable) Upsert(account types.OChainUniverseAccount) error {
	key := []byte(OChainUniverseAccountPrefix + account.UniverseId + "_" + account.Address)
	value, err := cbor.Marshal(account)
	if err != nil {
		return err
	}

	return db.currentTxn.Set(key, value)
}

func (db *OChainUniverseAccountTable) Delete(universeId string, address string) error {
	key := []byte(OChainUniverseAccountPrefix + address)
	return db.currentTxn.Delete(key)
}

func (db *OChainUniverseAccountTable) GetAll() ([]types.OChainUniverseAccount, error) {
	var at uint64 = math.MaxUint64
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
