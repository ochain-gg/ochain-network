package database

import (
	"context"
	"errors"
	"math"

	"github.com/dgraph-io/badger/pb"
	"github.com/dgraph-io/badger/v4"
	"github.com/dgraph-io/ristretto/z"
	"github.com/fxamacker/cbor/v2"
	"github.com/ochain-gg/ochain-network/types"
)

const (
	OChainBridgeTransactionPrefix string = "bridge_txs_"
)

// Global account side
type OChainBridgeTransactionTable struct {
	bdb        *badger.DB
	currentTxn *badger.Txn
}

func (db *OChainBridgeTransactionTable) SetCurrentTxn(tx *badger.Txn) {
	db.currentTxn = tx
}

func (db *OChainBridgeTransactionTable) Exists(hash string) (bool, error) {
	var at uint64 = math.MaxUint64
	return db.ExistsAt(hash, at)
}

func (db *OChainBridgeTransactionTable) ExistsAt(hash string, at uint64) (bool, error) {
	key := []byte(OChainBridgeTransactionPrefix + hash)
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

func (db *OChainBridgeTransactionTable) Get(hash string) (types.OChainBridgeTransaction, error) {
	var at uint64 = math.MaxUint64
	return db.GetAt(hash, at)
}

func (db *OChainBridgeTransactionTable) GetAt(hash string, at uint64) (types.OChainBridgeTransaction, error) {
	var transaction types.OChainBridgeTransaction
	key := []byte(OChainBridgeTransactionPrefix + hash)
	txn := db.bdb.NewTransactionAt(at, false)

	item, err := txn.Get([]byte(key))
	if err != nil {
		return types.OChainBridgeTransaction{}, err
	}

	value, err := item.ValueCopy(nil)
	if err != nil {
		return types.OChainBridgeTransaction{}, err
	}

	err = cbor.Unmarshal(value, &transaction)
	if err != nil {
		return types.OChainBridgeTransaction{}, err
	}

	return transaction, nil
}

func (db *OChainBridgeTransactionTable) Insert(transaction types.OChainBridgeTransaction) error {
	key := []byte(OChainBridgeTransactionPrefix + transaction.Hash)

	exists, err := db.ExistsAt(transaction.Hash, db.currentTxn.ReadTs())
	if err != nil {
		return err
	}

	if exists {
		return errors.New("transaction already exists")
	}

	value, err := cbor.Marshal(transaction)
	if err != nil {
		return err
	}

	return db.currentTxn.Set(key, value)
}

func (db *OChainBridgeTransactionTable) Update(transaction types.OChainBridgeTransaction) error {
	key := []byte(OChainBridgeTransactionPrefix + transaction.Hash)

	exists, err := db.ExistsAt(transaction.Hash, db.currentTxn.ReadTs())
	if err != nil {
		return err
	}

	if !exists {
		return errors.New("transaction doesn't exists")
	}

	value, err := cbor.Marshal(transaction)
	if err != nil {
		return err
	}

	return db.currentTxn.Set(key, value)
}

func (db *OChainBridgeTransactionTable) Upsert(transaction types.OChainBridgeTransaction) error {
	key := []byte(OChainBridgeTransactionPrefix + transaction.Hash)
	value, err := cbor.Marshal(transaction)
	if err != nil {
		return err
	}

	return db.currentTxn.Set(key, value)
}

func (db *OChainBridgeTransactionTable) Delete(address string) error {
	key := []byte(OChainBridgeTransactionPrefix + address)
	return db.currentTxn.Delete(key)
}

func (db *OChainBridgeTransactionTable) GetAll() ([]types.OChainBridgeTransaction, error) {
	var at uint64 = math.MaxUint64
	return db.GetAllAt(at)
}

func (table *OChainBridgeTransactionTable) GetAllAt(at uint64) ([]types.OChainBridgeTransaction, error) {
	var results []types.OChainBridgeTransaction

	stream := table.bdb.NewStreamAt(at)

	stream.NumGo = 16                                        // Set number of goroutines to use for iteration.
	stream.Prefix = []byte(OChainBridgeTransactionPrefix)    // Leave nil for iteration over the whole DB.
	stream.LogPrefix = "Badger.BridgeTransactions.Streaming" // For identifying stream logs. Outputs to Logger.
	stream.Send = func(buf *z.Buffer) error {
		err := buf.SliceIterate(func(s []byte) error {
			var kv pb.KV
			var tx types.OChainBridgeTransaction
			if err := kv.Unmarshal(s); err != nil {
				return err
			}

			if err := cbor.Unmarshal(kv.Value, tx); err != nil {
				return err
			}

			results = append(results, tx)
			return nil
		})

		return err
	}

	if err := stream.Orchestrate(context.Background()); err != nil {
		return []types.OChainBridgeTransaction{}, err
	}

	return results, nil
}

func (table *OChainBridgeTransactionTable) GetByAccountAt(address string, at uint64) ([]types.OChainBridgeTransaction, error) {
	var results []types.OChainBridgeTransaction

	stream := table.bdb.NewStreamAt(at)

	stream.NumGo = 16                                        // Set number of goroutines to use for iteration.
	stream.Prefix = []byte(OChainBridgeTransactionPrefix)    // Leave nil for iteration over the whole DB.
	stream.LogPrefix = "Badger.BridgeTransactions.Streaming" // For identifying stream logs. Outputs to Logger.
	stream.Send = func(buf *z.Buffer) error {
		err := buf.SliceIterate(func(s []byte) error {
			var kv pb.KV
			var tx types.OChainBridgeTransaction
			if err := kv.Unmarshal(s); err != nil {
				return err
			}

			if err := cbor.Unmarshal(kv.Value, tx); err != nil {
				return err
			}

			if tx.Account == address {
				results = append(results, tx)
			}

			return nil
		})

		return err
	}

	if err := stream.Orchestrate(context.Background()); err != nil {
		return []types.OChainBridgeTransaction{}, err
	}

	return results, nil
}

func NewOChainBridgeTransactionTable(db *badger.DB) *OChainBridgeTransactionTable {
	return &OChainBridgeTransactionTable{
		bdb: db,
	}
}
