package database

import (
	"errors"

	"github.com/dgraph-io/badger/v4"
	"github.com/ochain.gg/ochain-network-validator/types"
	"github.com/timshannon/badgerhold/v4"
)

type OChainAccountTable struct {
	db *badgerhold.Store
}

type OChainAccountIAM struct {
	Pubkey []byte
	Nonce  uint64

	GuardianQuorum uint64
	Guardians      [][]byte

	DeleguatedTo [][]byte
}

type OChainAccount struct {
	Id       uint `badgerhold:"key"`
	Universe uint `badgerhold:"index"`

	IAM          OChainAccountIAM
	Technologies types.OChainAccountTechnologies
}

func (table *OChainAccountTable) Get(id uint) (OChainAccount, error) {
	var result []OChainAccount
	err := table.db.Find(&result, badgerhold.Where("Id").Eq(id))

	if err != nil {
		return OChainAccount{}, err
	}

	if len(result) == 0 {
		return OChainAccount{}, errors.New("account not found")
	}

	return result[0], nil
}

func (table *OChainAccountTable) GetAll(id uint) []OChainAccount {
	var result []OChainAccount
	table.db.Find(&result, badgerhold.Where("Id").Eq(id))

	return result
}

func (table *OChainAccountTable) Insert(account OChainAccount, tx *badger.Txn) error {
	err := table.db.TxInsert(tx, badgerhold.NextSequence(), &account)
	return err
}

func (table *OChainAccountTable) Save(account OChainAccount, tx *badger.Txn) error {
	err := table.db.TxUpdate(tx, account.Id, account)
	return err
}

func (table *OChainAccountTable) Delete(account OChainAccount, tx *badger.Txn) error {
	err := table.db.TxDelete(tx, account.Id, account)
	return err
}

func NewOChainAccountTable(db *badgerhold.Store) *OChainAccountTable {
	return &OChainAccountTable{
		db: db,
	}
}
