package database

import (
	"github.com/timshannon/badgerhold/v4"
)

type OChainDatabase struct {
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

func NewOChainDatabase(db *badgerhold.Store) *OChainDatabase {
	return &OChainDatabase{
		Validators:             NewOChainValidatorTable(db),
		BridgeTransactions:     NewOChainBridgeTransactionTable(db),
		Universes:              NewOChainUniverseTable(db),
		UniverseConfigurations: NewOChainUniverseConfigurationTable(db),
		GlobalsAccounts:        NewOChainGlobalAccountTable(db),
		UniverseAccounts:       NewOChainUniverseAccountTable(db),
		Planets:                NewOChainPlanetTable(db),
		Fleets:                 NewOChainFleetTable(db),
		State:                  NewOChainStateTable(db),
		BridgeState:            NewOChainBridgeStateTable(db),
	}
}
