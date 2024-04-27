package database

import (
	"github.com/timshannon/badgerhold"
)

type OChainDatabase struct {
	Validators *OChainValidatorTable
	Universes  *OChainUniverseTable
	Accounts   *OChainAccountTable
	Planets    *OChainPlanetTable
	Fleets     *OChainFleetTable
}

func NewOChainDatabase(db *badgerhold.Store) *OChainDatabase {
	return &OChainDatabase{
		Validators: NewOChainValidatorTable(db),
		Universes:  NewOChainUniverseTable(db),
		Accounts:   NewOChainAccountTable(db),
		Planets:    NewOChainPlanetTable(db),
		Fleets:     NewOChainFleetTable(db),
	}
}
