package database

import (
	"log"

	"github.com/dgraph-io/badger/v4"
	"github.com/timshannon/badgerhold/v4"
)

type OChainBridgeStateTable struct {
	db *badgerhold.Store
}

type OChainBridgeState struct {
	Id                 int `badgerhold:"key"`
	LatestPortalUpdate uint64
}

func (state *OChainBridgeState) SetLatestPortalUpdate(blockNumber uint64) {
	state.LatestPortalUpdate = blockNumber
}

func (table *OChainBridgeStateTable) Get() (*OChainBridgeState, error) {
	var result []OChainBridgeState
	err := table.db.Find(&result, badgerhold.Where("Id").Eq(1))
	if err != nil {
		return &OChainBridgeState{}, err
	}

	log.Println("Loading state: ")
	log.Println(result)

	if len(result) > 0 {
		return &result[0], nil
	} else {
		return &OChainBridgeState{
			Id:                 1,
			LatestPortalUpdate: 0,
		}, nil
	}
}

func (table *OChainBridgeStateTable) Save(state *OChainBridgeState, txn *badger.Txn) error {
	err := table.db.TxUpsert(txn, 1, state)
	return err
}

func NewOChainBridgeStateTable(db *badgerhold.Store) *OChainBridgeStateTable {
	return &OChainBridgeStateTable{
		db: db,
	}
}
