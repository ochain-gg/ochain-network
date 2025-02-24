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

func (db *OChainStateTable) Exists() (bool, error) {
	var at uint64 = math.MaxUint64
	return db.ExistsAt(at)
}

func (db *OChainStateTable) ExistsAt(at uint64) (bool, error) {
	key := []byte(OChainStateKey)
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

func (db *OChainStateTable) Get() (types.OChainState, error) {
	var at uint64 = math.MaxUint64
	return db.GetAt(at)
}

func (db *OChainStateTable) GetAt(at uint64) (types.OChainState, error) {
	var state types.OChainState
	key := []byte(OChainStateKey)
	txn := db.bdb.NewTransactionAt(at, false)

	item, err := txn.Get([]byte(key))
	if err != nil {

		if err == badger.ErrKeyNotFound {
			return types.OChainState{
				Size:               0,
				Height:             0,
				Hash:               []byte(""),
				LatestPortalUpdate: 0,
			}, nil
		} else {
			return types.OChainState{}, err
		}
	}

	value, err := item.ValueCopy(nil)
	if err != nil {
		return types.OChainState{}, err
	}

	err = cbor.Unmarshal(value, &state)
	if err != nil {
		return types.OChainState{}, err
	}

	return state, nil
}

func (db *OChainStateTable) Upsert(state types.OChainState) error {
	key := []byte(OChainStateKey)
	value, err := cbor.Marshal(state)
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
