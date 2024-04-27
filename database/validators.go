package database

import (
	"errors"

	"github.com/dgraph-io/badger"
	"github.com/timshannon/badgerhold"
)

type OChainValidatorTable struct {
	db *badgerhold.Store
}

type OChainValidator struct {
	Id                        uint64
	Stacker                   string
	Validator                 string
	PublicKey                 string
	Enabled                   bool
	StackingTreansactionHash  string
	UnstackingTransactionHash string
}

func (table *OChainValidatorTable) Get(id uint64, tx *badger.Txn) (OChainValidator, error) {
	var result []OChainValidator
	err := table.db.TxFind(tx, &result, badgerhold.Where("Id").Eq(id))

	if err != nil {
		return OChainValidator{}, err
	}

	if len(result) == 0 {
		return OChainValidator{}, errors.New("account not found")
	}

	return result[0], nil
}

func (table *OChainValidatorTable) GetAll(tx *badger.Txn) []OChainValidator {
	var result []OChainValidator
	table.db.TxFind(tx, &result, badgerhold.Where("Id").Ge(0))

	return result
}

func (table *OChainValidatorTable) Insert(validator OChainValidator, tx *badger.Txn) error {
	err := table.db.TxInsert(tx, badgerhold.NextSequence(), &validator)
	return err
}

func (table *OChainValidatorTable) Save(validator OChainValidator, tx *badger.Txn) error {
	err := table.db.TxUpdate(tx, validator.Id, validator)
	return err
}

func (table *OChainValidatorTable) Delete(validator OChainValidator, tx *badger.Txn) error {
	err := table.db.TxDelete(tx, validator.Id, validator)
	return err
}

func NewOChainValidatorTable(db *badgerhold.Store) *OChainValidatorTable {
	return &OChainValidatorTable{
		db: db,
	}
}
