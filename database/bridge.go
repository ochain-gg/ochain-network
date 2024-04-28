package database

import (
	"errors"

	"github.com/dgraph-io/badger/v4"
	"github.com/timshannon/badgerhold/v4"
)

// Global account side
type OChainBridgeTransactionTable struct {
	db *badgerhold.Store
}

type OChainBridgeTransactionType uint64

const (
	OChainBridgeDepositTransaction   OChainBridgeTransactionType = 0
	OChainBridgeWithdrawTransaction  OChainBridgeTransactionType = 1
	OChainBridgeSubscribeTransaction OChainBridgeTransactionType = 2
)

type OChainBridgeTransaction struct {
	Id              string `badgerhold:"key"`
	Type            OChainBridgeTransactionType
	TransactionHash string
	Executed        bool
}

func (table *OChainBridgeTransactionTable) Get(id uint) (OChainBridgeTransaction, error) {
	var result []OChainBridgeTransaction
	err := table.db.Find(&result, badgerhold.Where("Id").Eq(id))

	if err != nil {
		return OChainBridgeTransaction{}, err
	}

	if len(result) == 0 {
		return OChainBridgeTransaction{}, errors.New("account not found")
	}

	return result[0], nil
}

func (table *OChainBridgeTransactionTable) GetByOwnerAddress(address string) (OChainBridgeTransaction, error) {
	var result []OChainBridgeTransaction
	err := table.db.Find(&result, badgerhold.Where("IAM.OwnerAddress").Eq(address))

	if err != nil {
		return OChainBridgeTransaction{}, err
	}

	if len(result) == 0 {
		return OChainBridgeTransaction{}, errors.New("account not found")
	}

	return result[0], nil
}

func (table *OChainBridgeTransactionTable) GetAll(id uint) []OChainBridgeTransaction {
	var result []OChainBridgeTransaction
	table.db.Find(&result, badgerhold.Where("Id").Eq(id))

	return result
}

func (table *OChainBridgeTransactionTable) Insert(account OChainBridgeTransaction, tx *badger.Txn) error {
	err := table.db.TxInsert(tx, badgerhold.NextSequence(), &account)
	return err
}

func (table *OChainBridgeTransactionTable) Save(account OChainBridgeTransaction, tx *badger.Txn) error {
	err := table.db.TxUpdate(tx, account.Id, account)
	return err
}

func (table *OChainBridgeTransactionTable) Delete(account OChainBridgeTransaction, tx *badger.Txn) error {
	err := table.db.TxDelete(tx, account.Id, account)
	return err
}

func NewOChainBridgeTransactionTable(db *badgerhold.Store) *OChainBridgeTransactionTable {
	return &OChainBridgeTransactionTable{
		db: db,
	}
}
