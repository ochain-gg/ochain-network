package database

// import (
// 	"errors"

// 	"github.com/dgraph-io/badger/v4"
// 	"github.com/ochain-gg/ochain-network/types"
// 	"github.com/timshannon/badgerhold/v4"
// )

// type OChainFleetTable struct {
// 	db *badgerhold.Store
// }

// type OChainFleetConfiguration struct {
// 	//Global configuration
// 	Speed                   uint
// 	MaxGalaxy               uint
// 	MaxSolarSystemPerGalaxy uint
// 	MaxPlanetPerSolarSystem uint

// 	Spaceships types.OChainSpaceship
// 	Defenses   types.OChainDefense
// }

// type OChainFleet struct {
// 	Id       uint64 `badgerhold:"key"`
// 	Universe uint64 `badgerhold:"index"`
// 	Owner    uint64 `badgerhold:"index"`

// 	Spaceships types.OChainFleetSpaceships
// 	Cargo      types.OChainResources

// 	From             uint `badgerhold:"index"`
// 	To               uint `badgerhold:"index"`
// 	ToAlliancePlanet bool

// 	IsTraveling bool

// 	FleetMode          uint8
// 	FleetModeConfirmed bool

// 	Distance    uint64
// 	MaxCargo    uint64
// 	MaxSpeed    uint64
// 	Speed       uint64
// 	ArriveAt    uint64
// 	BackAt      uint64
// 	IsReturning bool

// 	CreatedAt uint
// }

// func (table *OChainFleetTable) Get(id uint) (OChainFleet, error) {
// 	var result []OChainFleet
// 	err := table.db.Find(&result, badgerhold.Where("Id").Eq(id))

// 	if err != nil {
// 		return OChainFleet{}, err
// 	}

// 	if len(result) == 0 {
// 		return OChainFleet{}, errors.New("fleet not found")
// 	}

// 	return result[0], nil
// }

// func (table *OChainFleetTable) GetAll(id uint) []OChainFleet {
// 	var result []OChainFleet
// 	table.db.Find(&result, badgerhold.Where("Id").Eq(id))

// 	return result
// }

// func (table *OChainFleetTable) Insert(fleet OChainFleet, tx *badger.Txn) error {
// 	err := table.db.TxInsert(tx, badgerhold.NextSequence(), &fleet)
// 	return err
// }

// func (table *OChainFleetTable) Save(fleet OChainFleet, tx *badger.Txn) error {
// 	err := table.db.TxUpdate(tx, fleet.Id, fleet)
// 	return err
// }

// func (table *OChainFleetTable) Delete(fleet OChainFleet, tx *badger.Txn) error {
// 	err := table.db.TxDelete(tx, fleet.Id, fleet)
// 	return err
// }

// func NewOChainFleetTable(db *badgerhold.Store) *OChainFleetTable {
// 	return &OChainFleetTable{
// 		db: db,
// 	}
// }
