package database

import (
	"github.com/timshannon/badgerhold/v4"
)

type OChainDatabase struct {
	Validators         *OChainValidatorTable
	BridgeTransactions *OChainBridgeTransactionTable
	Universes          *OChainUniverseTable
	GlobalsAccounts    *OChainGlobalAccountTable
	UniverseAccounts   *OChainUniverseAccountTable
	Planets            *OChainPlanetTable
	Fleets             *OChainFleetTable
	State              *OChainStateTable
}

func NewOChainDatabase(db *badgerhold.Store) *OChainDatabase {
	return &OChainDatabase{
		Validators:         NewOChainValidatorTable(db),
		BridgeTransactions: NewOChainBridgeTransactionTable(db),
		Universes:          NewOChainUniverseTable(db),
		GlobalsAccounts:    NewOChainGlobalAccountTable(db),
		UniverseAccounts:   NewOChainUniverseAccountTable(db),
		Planets:            NewOChainPlanetTable(db),
		Fleets:             NewOChainFleetTable(db),
		State:              NewOChainStateTable(db),
	}
}
