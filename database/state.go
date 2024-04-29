package database

import (
	"github.com/timshannon/badgerhold/v4"
)

type OChainStateTable struct {
	db *badgerhold.Store
}

type OChainState struct {
	Id      uint64
	Size    int64  `json:"size"`
	Height  int64  `json:"height"`
	AppHash []byte `json:"app_hash"`

	LatestPortalUpdate uint64
}

func (table *OChainStateTable) Get() (OChainState, error) {
	var result []OChainState
	err := table.db.Find(&result, badgerhold.Where("Id").Eq(1))

	if err != nil {
		return OChainState{}, err
	}

	if len(result) > 0 {
		return result[0], nil
	} else {
		return OChainState{
			Size:    0,
			Height:  0,
			AppHash: []byte(""),
		}, nil
	}
}

func (table *OChainStateTable) Save(state OChainState) error {
	err := table.db.Upsert(1, state)
	return err
}

func NewOChainStateTable(db *badgerhold.Store) *OChainStateTable {
	return &OChainStateTable{
		db: db,
	}
}
