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
	OChainTransactionPrefix string = "bridge_txs_"
)

// Global account side
type OChainTransactionTable struct {
	bdb        *badger.DB
	currentTxn *badger.Txn
}

func (db *OChainTransactionTable) SetCurrentTxn(tx *badger.Txn) {
	db.currentTxn = tx
}

func (db *OChainTransactionTable) Exists(hash string) (bool, error) {
	var at uint64 = math.MaxUint64
	return db.ExistsAt(hash, at)
}

func (db *OChainTransactionTable) ExistsAt(hash string, at uint64) (bool, error) {
	key := []byte(OChainTransactionPrefix + hash)
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

func (db *OChainTransactionTable) Get(hash string) (types.OChainTransaction, error) {
	var at uint64 = math.MaxUint64
	return db.GetAt(hash, at)
}

func (db *OChainTransactionTable) GetAt(hash string, at uint64) (types.OChainTransaction, error) {
	var transaction types.OChainTransaction
	key := []byte(OChainTransactionPrefix + hash)
	txn := db.bdb.NewTransactionAt(at, false)

	item, err := txn.Get([]byte(key))
	if err != nil {
		return types.OChainTransaction{}, err
	}

	value, err := item.ValueCopy(nil)
	if err != nil {
		return types.OChainTransaction{}, err
	}

	err = cbor.Unmarshal(value, &transaction)
	if err != nil {
		return types.OChainTransaction{}, err
	}

	return transaction, nil
}

func (db *OChainTransactionTable) Insert(transaction types.OChainTransaction) error {
	key := []byte(OChainTransactionPrefix + transaction.Hash)

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

func (db *OChainTransactionTable) Update(transaction types.OChainTransaction) error {
	key := []byte(OChainTransactionPrefix + transaction.Hash)

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

func (db *OChainTransactionTable) Upsert(transaction types.OChainTransaction) error {
	key := []byte(OChainTransactionPrefix + transaction.Hash)
	value, err := cbor.Marshal(transaction)
	if err != nil {
		return err
	}

	return db.currentTxn.Set(key, value)
}

func (db *OChainTransactionTable) Delete(address string) error {
	key := []byte(OChainTransactionPrefix + address)
	return db.currentTxn.Delete(key)
}

func (db *OChainTransactionTable) GetAll() ([]types.OChainTransaction, error) {
	var at uint64 = math.MaxUint64
	return db.GetAllAt(at)
}

func (table *OChainTransactionTable) GetAllAt(at uint64) ([]types.OChainTransaction, error) {
	var results []types.OChainTransaction

	stream := table.bdb.NewStreamAt(at)

	stream.NumGo = 16                                        // Set number of goroutines to use for iteration.
	stream.Prefix = []byte(OChainTransactionPrefix)          // Leave nil for iteration over the whole DB.
	stream.LogPrefix = "Badger.BridgeTransactions.Streaming" // For identifying stream logs. Outputs to Logger.
	stream.Send = func(buf *z.Buffer) error {
		err := buf.SliceIterate(func(s []byte) error {
			var kv pb.KV
			var tx types.OChainTransaction
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
		return []types.OChainTransaction{}, err
	}

	return results, nil
}

func (table *OChainTransactionTable) GetByAccountAt(address string, at uint64) ([]types.OChainTransaction, error) {
	var results []types.OChainTransaction

	stream := table.bdb.NewStreamAt(at)

	stream.NumGo = 16                                        // Set number of goroutines to use for iteration.
	stream.Prefix = []byte(OChainTransactionPrefix)          // Leave nil for iteration over the whole DB.
	stream.LogPrefix = "Badger.BridgeTransactions.Streaming" // For identifying stream logs. Outputs to Logger.
	stream.Send = func(buf *z.Buffer) error {
		err := buf.SliceIterate(func(s []byte) error {
			var kv pb.KV
			var tx types.OChainTransaction
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
		return []types.OChainTransaction{}, err
	}

	return results, nil
}

func NewOChainTransactionTable(db *badger.DB) *OChainTransactionTable {
	return &OChainTransactionTable{
		bdb: db,
	}
}
