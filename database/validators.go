package database

import (
	"errors"
	"math"

	"github.com/dgraph-io/badger/v4"
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
	var at uint64
	at = math.MaxUint64
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

func (db *OChainValidatorTable) Get(address string) (types.OChainValidator, error) {
	var at uint64
	at = math.MaxUint64
	return db.GetAt(address, at)
}

func (db *OChainValidatorTable) GetAt(address string, at uint64) (types.OChainValidator, error) {
	var account types.OChainValidator
	key := []byte(OChainValidatorPrefix + address)
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
	var at uint64
	at = math.MaxUint64
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
