package database

import (
	"errors"

	"github.com/dgraph-io/badger/v4"
	"github.com/ochain.gg/ochain-network/types"
	"github.com/timshannon/badgerhold/v4"
)

// Global account side
type OChainGlobalAccountTable struct {
	db *badgerhold.Store
}

func (table *OChainGlobalAccountTable) Get(id uint) (types.OChainGlobalAccount, error) {
	var result []types.OChainGlobalAccount
	err := table.db.Find(&result, badgerhold.Where("Id").Eq(id))

	if err != nil {
		return types.OChainGlobalAccount{}, err
	}

	if len(result) == 0 {
		return types.OChainGlobalAccount{}, errors.New("account not found")
	}

	return result[0], nil
}

func (table *OChainGlobalAccountTable) GetByOwnerAddress(address string) (types.OChainGlobalAccount, error) {
	var result []types.OChainGlobalAccount
	err := table.db.Find(&result, badgerhold.Where("IAM.OwnerAddress").Eq(address))

	if err != nil {
		return types.OChainGlobalAccount{}, err
	}

	if len(result) == 0 {
		return types.OChainGlobalAccount{}, errors.New("account not found")
	}

	return result[0], nil
}

func (table *OChainGlobalAccountTable) GetAll(id uint) []types.OChainGlobalAccount {
	var result []types.OChainGlobalAccount
	table.db.Find(&result, badgerhold.Where("Id").Eq(id))

	return result
}

func (table *OChainGlobalAccountTable) Insert(account types.OChainGlobalAccount, tx *badger.Txn) error {
	err := table.db.TxInsert(tx, badgerhold.NextSequence(), &account)
	return err
}

func (table *OChainGlobalAccountTable) Save(account types.OChainGlobalAccount, tx *badger.Txn) error {
	err := table.db.TxUpdate(tx, account.Id, account)
	return err
}

func (table *OChainGlobalAccountTable) Delete(account types.OChainGlobalAccount, tx *badger.Txn) error {
	err := table.db.TxDelete(tx, account.Id, account)
	return err
}

func NewOChainGlobalAccountTable(db *badgerhold.Store) *OChainGlobalAccountTable {
	return &OChainGlobalAccountTable{
		db: db,
	}
}

// Universe account side
type OChainUniverseAccountTable struct {
	db *badgerhold.Store
}

func (table *OChainUniverseAccountTable) Get(id uint) (types.OChainUniverseAccount, error) {
	var result []types.OChainUniverseAccount
	err := table.db.Find(&result, badgerhold.Where("Id").Eq(id))

	if err != nil {
		return types.OChainUniverseAccount{}, err
	}

	if len(result) == 0 {
		return types.OChainUniverseAccount{}, errors.New("account not found")
	}

	return result[0], nil
}

func (table *OChainUniverseAccountTable) GetAll(id uint) []types.OChainUniverseAccount {
	var result []types.OChainUniverseAccount
	table.db.Find(&result, badgerhold.Where("Id").Eq(id))

	return result
}

func (table *OChainUniverseAccountTable) Insert(account types.OChainUniverseAccount, tx *badger.Txn) error {
	err := table.db.TxInsert(tx, badgerhold.NextSequence(), &account)
	return err
}

func (table *OChainUniverseAccountTable) Save(account types.OChainUniverseAccount, tx *badger.Txn) error {
	err := table.db.TxUpdate(tx, account.Id, account)
	return err
}

func (table *OChainUniverseAccountTable) Delete(account types.OChainUniverseAccount, tx *badger.Txn) error {
	err := table.db.TxDelete(tx, account.Id, account)
	return err
}

func NewOChainUniverseAccountTable(db *badgerhold.Store) *OChainUniverseAccountTable {
	return &OChainUniverseAccountTable{
		db: db,
	}
}
