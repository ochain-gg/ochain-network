package database

import (
	"github.com/dgraph-io/badger/v4"
)

type OChainDatabase struct {
	bdb                    *badger.DB
	currentTxn             *badger.Txn
	Validators             *OChainValidatorTable
	BridgeTransactions     *OChainBridgeTransactionTable
	Universes              *OChainUniverseTable
	UniverseConfigurations *OChainUniverseConfigurationTable
	GlobalsAccounts        *OChainGlobalAccountTable
	UniverseAccounts       *OChainUniverseAccountTable
	Planets                *OChainPlanetTable
	Fleets                 *OChainFleetTable
	State                  *OChainStateTable
	BridgeState            *OChainBridgeStateTable
}

func (db *OChainDatabase) Open(path string) error {
	opts := badger.DefaultOptions(path)

	bdb, err := badger.OpenManaged(opts)
	if err != nil {
		return err
	}

	db.bdb = bdb
	return nil
}

func (db *OChainDatabase) Close() error {
	return db.bdb.Close()
}

func (db *OChainDatabase) LoadTables() error {
	return db.bdb.Close()
}

func (db *OChainDatabase) NewTransaction(ts uint64) {
	db.currentTxn = db.bdb.NewTransactionAt(ts, true)
}

func (db *OChainDatabase) CommitTransaction(ts uint64) error {
	if err := db.currentTxn.CommitAt(ts, nil); err != nil {
		return err
	}
	return nil
}

func NewOChainDatabase(path string) *OChainDatabase {

	db := &OChainDatabase{}
	db.Open(path)
	db.LoadTables()

	return OChainDatabase
}
