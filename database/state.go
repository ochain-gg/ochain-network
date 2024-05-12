package database

import (
	"errors"
	"math"

	"github.com/dgraph-io/badger/v4"
	"github.com/fxamacker/cbor/v2"
	"github.com/ochain-gg/ochain-network/types"
)

const (
	OChainStateKey string = "ochain_network_state"
)

type OChainStateTable struct {
	bdb        *badger.DB
	currentTxn *badger.Txn
}

func (db *OChainStateTable) SetCurrentTxn(tx *badger.Txn) {
	db.currentTxn = tx
}

func (db *OChainStateTable) Exists(address string) (bool, error) {
	var at uint64
	at = math.MaxUint64
	return db.ExistsAt(address, at)
}

func (db *OChainStateTable) ExistsAt(address string, at uint64) (bool, error) {
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

func (db *OChainStateTable) Get(address string) (types.OChainUniverseAccount, error) {
	var at uint64
	at = math.MaxUint64
	return db.GetAt(address, at)
}

func (db *OChainStateTable) GetAt(address string, at uint64) (types.OChainUniverseAccount, error) {
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

func (db *OChainStateTable) Upsert(account types.OChainUniverseAccount) error {
	key := []byte(OChainUniverseAccountPrefix + account.Address)
	value, err := cbor.Marshal(account)
	if err != nil {
		return err
	}

	return db.currentTxn.Set(key, value)
}

func NewOChainStateTable(db *badger.DB) *OChainStateTable {
	return &OChainStateTable{
		bdb: db,
	}
}
