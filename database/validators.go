package database

import (
	"context"
	"errors"
	"math"
	"strconv"

	"github.com/dgraph-io/badger/pb"
	"github.com/dgraph-io/badger/v4"
	"github.com/dgraph-io/ristretto/v2/z"
	"github.com/fxamacker/cbor/v2"
	"github.com/ochain-gg/ochain-network/types"
)

const (
	OChainValidatorPrefix string = "validator_"
)

type OChainValidatorTable struct {
	bdb        *badger.DB
	currentTxn *badger.Txn
}

func (db *OChainValidatorTable) SetCurrentTxn(tx *badger.Txn) {
	db.currentTxn = tx
}

func (db *OChainValidatorTable) Exists(address string) (bool, error) {
	var at uint64 = math.MaxUint64
	return db.ExistsAt(address, at)
}

func (db *OChainValidatorTable) ExistsAt(address string, at uint64) (bool, error) {
	key := []byte(OChainValidatorPrefix + address)
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

func (db *OChainValidatorTable) IsEnabled(address string) (bool, error) {
	key := []byte(OChainValidatorPrefix + address)
	txn := db.bdb.NewTransactionAt(math.MaxUint64, false)
	if item, err := txn.Get([]byte(key)); err != nil {
		if errors.Is(err, badger.ErrKeyNotFound) {
			return false, nil
		} else {
			return false, err
		}
	} else {
		var validator types.OChainValidator
		value, err := item.ValueCopy(nil)
		if err != nil {
			return false, err
		}

		if err := cbor.Unmarshal(value, validator); err != nil {
			return false, err
		}

		return validator.Enabled, nil
	}
}

func (db *OChainValidatorTable) GetByAddress(address string) (types.OChainValidator, error) {
	var at uint64 = math.MaxUint64
	return db.GetByAddressAt(address, at)
}

func (db *OChainValidatorTable) GetByAddressAt(address string, at uint64) (types.OChainValidator, error) {
	var account types.OChainValidator

	stream := db.bdb.NewStreamAt(at)

	stream.NumGo = 16                               // Set number of goroutines to use for iteration.
	stream.Prefix = []byte(OChainValidatorPrefix)   // Leave nil for iteration over the whole DB.
	stream.LogPrefix = "Badger.Validator.Streaming" // For identifying stream logs. Outputs to Logger.
	stream.Send = func(buf *z.Buffer) error {
		err := buf.SliceIterate(func(s []byte) error {
			var kv pb.KV
			var acc types.OChainValidator
			if err := kv.Unmarshal(s); err != nil {
				return err
			}

			if err := cbor.Unmarshal(kv.Value, acc); err != nil {
				return err
			}

			if acc.PublicKey == address {
				account = acc
			}

			return nil
		})

		return err
	}

	if err := stream.Orchestrate(context.Background()); err != nil {
		return types.OChainValidator{}, err
	}

	return account, nil
}

func (db *OChainValidatorTable) GetById(id uint64) (types.OChainValidator, error) {
	var at uint64 = math.MaxUint64
	return db.GetByIdAt(id, at)
}

func (db *OChainValidatorTable) GetByIdAt(id uint64, at uint64) (types.OChainValidator, error) {
	var account types.OChainValidator
	key := []byte(OChainValidatorPrefix + strconv.FormatUint(id, 10))
	txn := db.bdb.NewTransactionAt(at, false)

	item, err := txn.Get([]byte(key))
	if err != nil {
		return types.OChainValidator{}, err
	}

	value, err := item.ValueCopy(nil)
	if err != nil {
		return types.OChainValidator{}, err
	}

	err = cbor.Unmarshal(value, &account)
	if err != nil {
		return types.OChainValidator{}, err
	}

	return account, nil
}

func (db *OChainValidatorTable) Insert(validator types.OChainValidator) error {
	key := []byte(OChainValidatorPrefix + validator.PublicKey)

	exists, err := db.ExistsAt(validator.PublicKey, db.currentTxn.ReadTs())
	if err != nil {
		return err
	}

	if exists {
		return errors.New("validator already exists")
	}

	value, err := cbor.Marshal(validator)
	if err != nil {
		return err
	}

	return db.currentTxn.Set(key, value)
}

func (db *OChainValidatorTable) Update(validator types.OChainValidator) error {
	key := []byte(OChainValidatorPrefix + validator.PublicKey)

	exists, err := db.ExistsAt(validator.PublicKey, db.currentTxn.ReadTs())
	if err != nil {
		return err
	}

	if !exists {
		return errors.New("validator doesn't exists")
	}

	value, err := cbor.Marshal(validator)
	if err != nil {
		return err
	}

	return db.currentTxn.Set(key, value)
}

func (db *OChainValidatorTable) Upsert(validator types.OChainValidator) error {
	key := []byte(OChainValidatorPrefix + validator.PublicKey)
	value, err := cbor.Marshal(validator)
	if err != nil {
		return err
	}

	return db.currentTxn.Set(key, value)
}

func (db *OChainValidatorTable) Delete(address string) error {
	key := []byte(OChainValidatorPrefix + address)
	return db.currentTxn.Delete(key)
}

func (db *OChainValidatorTable) GetAll() ([]types.OChainValidator, error) {
	var at uint64 = math.MaxUint64
	return db.GetAllAt(at)
}

func (db *OChainValidatorTable) GetAllAt(at uint64) ([]types.OChainValidator, error) {
	var validators []types.OChainValidator

	txn := db.bdb.NewTransactionAt(at, false)
	it := txn.NewIterator(badger.DefaultIteratorOptions)
	defer it.Close()

	prefix := []byte(OChainValidatorPrefix)

	for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
		item := it.Item()

		var validator types.OChainValidator
		value, err := item.ValueCopy(nil)
		if err != nil {
			return []types.OChainValidator{}, err
		}

		err = cbor.Unmarshal(value, &validator)
		if err != nil {
			return []types.OChainValidator{}, err
		}

		validators = append(validators, validator)
	}

	return validators, nil
}

func NewOChainValidatorTable(db *badger.DB) *OChainValidatorTable {
	return &OChainValidatorTable{
		bdb: db,
	}
}
