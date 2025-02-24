package database

import (
	"errors"
	"math"

	"github.com/dgraph-io/badger/v4"
	"github.com/fxamacker/cbor/v2"
	"github.com/ochain-gg/ochain-network/types"
)

const (
	OChainAlliancePrefix string = "alliance_"
)

type OChainAllianceTable struct {
	bdb        *badger.DB
	currentTxn *badger.Txn
}

func (db *OChainAllianceTable) SetCurrentTxn(tx *badger.Txn) {
	db.currentTxn = tx
}

func (db *OChainAllianceTable) Exists(universeId string, id string) (bool, error) {
	var at uint64 = math.MaxUint64
	return db.ExistsAt(universeId, id, at)
}

func (db *OChainAllianceTable) ExistsAt(universeId string, id string, at uint64) (bool, error) {
	key := []byte(OChainAlliancePrefix + universeId + "_" + id)

	txn := db.bdb.NewTransactionAt(at, false)
	if _, err := txn.Get([]byte(key)); err != nil {
		if errors.Is(err, badger.ErrKeyNotFound) {
			return false, nil
		} else {
			return false, err
		}
	} else {
		return true, nil
	}
}

func (db *OChainAllianceTable) Get(universeId string, id string) (types.OChainAlliance, error) {
	var at uint64 = math.MaxUint64
	return db.GetAt(universeId, id, at)
}

func (db *OChainAllianceTable) GetAt(universeId string, id string, at uint64) (types.OChainAlliance, error) {
	var planet types.OChainAlliance
	key := []byte(OChainAlliancePrefix + universeId + "_" + id)
	txn := db.bdb.NewTransactionAt(at, false)

	item, err := txn.Get([]byte(key))
	if err != nil {
		return types.OChainAlliance{}, err
	}

	value, err := item.ValueCopy(nil)
	if err != nil {
		return types.OChainAlliance{}, err
	}

	err = cbor.Unmarshal(value, &planet)
	if err != nil {
		return types.OChainAlliance{}, err
	}

	return planet, nil
}

func (db *OChainAllianceTable) Insert(universeId string, alliance types.OChainAlliance) error {
	key := []byte(OChainAlliancePrefix + universeId + "_" + alliance.Id)

	exists, err := db.ExistsAt(universeId, alliance.Id, db.currentTxn.ReadTs())
	if err != nil {
		return err
	}

	if exists {
		return errors.New("alliance already exists")
	}

	value, err := cbor.Marshal(alliance)
	if err != nil {
		return err
	}

	return db.currentTxn.Set(key, value)
}

func (db *OChainAllianceTable) Update(universeId string, alliance types.OChainAlliance) error {
	key := []byte(OChainAlliancePrefix + universeId + "_" + alliance.Id)

	exists, err := db.ExistsAt(universeId, alliance.Id, db.currentTxn.ReadTs())
	if err != nil {
		return err
	}

	if !exists {
		return errors.New("universe doesn't exists")
	}

	value, err := cbor.Marshal(alliance)
	if err != nil {
		return err
	}

	return db.currentTxn.Set(key, value)
}

func (db *OChainAllianceTable) Upsert(universeId string, alliance types.OChainAlliance) error {
	key := []byte(OChainAlliancePrefix + universeId + "_" + alliance.Id)
	value, err := cbor.Marshal(alliance)
	if err != nil {
		return err
	}

	return db.currentTxn.Set(key, value)
}

func (db *OChainAllianceTable) Delete(universeId string, id string) error {

	key := []byte(OChainAlliancePrefix + universeId + "_" + id)
	return db.currentTxn.Delete(key)
}

func (db *OChainAllianceTable) GetAll(universeId string) ([]types.OChainAlliance, error) {
	var at uint64 = math.MaxUint64
	return db.GetAllAt(universeId, at)
}

func (db *OChainAllianceTable) GetAllAt(universeId string, at uint64) ([]types.OChainAlliance, error) {
	var universes []types.OChainAlliance

	txn := db.bdb.NewTransactionAt(at, false)
	it := txn.NewIterator(badger.DefaultIteratorOptions)
	defer it.Close()

	prefix := []byte(OChainAlliancePrefix + universeId)

	for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
		item := it.Item()

		var universe types.OChainAlliance
		value, err := item.ValueCopy(nil)
		if err != nil {
			return []types.OChainAlliance{}, err
		}

		err = cbor.Unmarshal(value, &universe)
		if err != nil {
			return []types.OChainAlliance{}, err
		}

		universes = append(universes, universe)
	}

	return universes, nil
}

func NewOChainAllianceTable(db *badger.DB) *OChainAllianceTable {
	return &OChainAllianceTable{
		bdb: db,
	}
}
