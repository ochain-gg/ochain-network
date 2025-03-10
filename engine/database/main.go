package database

import (
	"github.com/dgraph-io/badger/v4"
)

type OChainDatabase struct {
	DB                         *badger.DB
	CurrentTxn                 *badger.Txn
	Validators                 *OChainValidatorTable
	Epochs                     *OChainEpochTable
	RewardPrograms             *OChainRewardProgramTable
	BridgeTransactions         *OChainTransactionTable
	Universes                  *OChainUniverseTable
	GlobalsAccounts            *OChainGlobalAccountTable
	UniverseAccounts           *OChainUniverseAccountTable
	Planets                    *OChainPlanetTable
	Buildings                  *OChainBuildingTable
	Technologies               *OChainTechnologyTable
	Defenses                   *OChainDefenseTable
	Spaceships                 *OChainSpaceshipTable
	Fleets                     *OChainFleetTable
	Alliance                   *OChainAllianceTable
	ResourcesMarket            *OChainResourcesMarketTable
	UniverseAccountWeeklyUsage *OChainUniverseAccountWeeklyUsageTable
	Upgrades                   *OChainUpgradeTable
	State                      *OChainStateTable
}

func (db *OChainDatabase) Open(path string) error {
	opts := badger.DefaultOptions(path)

	bdb, err := badger.OpenManaged(opts)
	if err != nil {
		return err
	}

	db.DB = bdb
	return nil
}

func (db *OChainDatabase) Close() error {
	return db.DB.Close()
}

func (db *OChainDatabase) LoadTables() {
	db.Validators = NewOChainValidatorTable(db.DB)
	db.Epochs = NewOChainEpochTable(db.DB)
	db.RewardPrograms = NewOChainRewardProgramTable(db.DB)
	db.BridgeTransactions = NewOChainTransactionTable(db.DB)
	db.Universes = NewOChainUniverseTable(db.DB)
	db.GlobalsAccounts = NewOChainGlobalAccountTable(db.DB)
	db.UniverseAccounts = NewOChainUniverseAccountTable(db.DB)
	db.Planets = NewOChainPlanetTable(db.DB)
	db.Buildings = NewOChainBuildingTable(db.DB)
	db.Technologies = NewOChainTechnologyTable(db.DB)
	db.Defenses = NewOChainDefenseTable(db.DB)
	db.Spaceships = NewOChainSpaceshipTable(db.DB)
	db.ResourcesMarket = NewOChainResourcesMarketTable(db.DB)
	db.UniverseAccountWeeklyUsage = NewOChainUniverseAccountWeeklyUsageTable(db.DB)
	db.Upgrades = NewOChainUpgradeTable(db.DB)
	db.Fleets = NewFleetTable(db.DB)
	db.Alliance = NewOChainAllianceTable(db.DB)
	db.State = NewOChainStateTable(db.DB)
}

func (db *OChainDatabase) NewTransaction(ts uint64) {
	db.CurrentTxn = db.DB.NewTransactionAt(ts, true)

	db.Validators.SetCurrentTxn(db.CurrentTxn)
	db.Epochs.SetCurrentTxn(db.CurrentTxn)
	db.RewardPrograms.SetCurrentTxn(db.CurrentTxn)
	db.BridgeTransactions.SetCurrentTxn(db.CurrentTxn)
	db.Universes.SetCurrentTxn(db.CurrentTxn)
	db.GlobalsAccounts.SetCurrentTxn(db.CurrentTxn)
	db.UniverseAccounts.SetCurrentTxn(db.CurrentTxn)
	db.Planets.SetCurrentTxn(db.CurrentTxn)
	db.Buildings.SetCurrentTxn(db.CurrentTxn)
	db.Technologies.SetCurrentTxn(db.CurrentTxn)
	db.Defenses.SetCurrentTxn(db.CurrentTxn)
	db.Spaceships.SetCurrentTxn(db.CurrentTxn)
	db.ResourcesMarket.SetCurrentTxn(db.CurrentTxn)
	db.UniverseAccountWeeklyUsage.SetCurrentTxn(db.CurrentTxn)
	db.Upgrades.SetCurrentTxn(db.CurrentTxn)
	db.Fleets.SetCurrentTxn(db.CurrentTxn)
	db.Alliance.SetCurrentTxn(db.CurrentTxn)
	db.State.SetCurrentTxn(db.CurrentTxn)
}

func (db *OChainDatabase) CommitTransaction() error {
	if err := db.CurrentTxn.CommitAt(db.CurrentTxn.ReadTs(), nil); err != nil {
		return err
	}
	return nil
}

func NewOChainDatabase(path string) *OChainDatabase {

	db := &OChainDatabase{}
	db.Open(path)
	db.LoadTables()

	return db
}
