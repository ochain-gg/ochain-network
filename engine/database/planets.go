package database

import (
	"errors"
	"math"

	"github.com/dgraph-io/badger/v4"
	"github.com/fxamacker/cbor/v2"
	"github.com/ochain-gg/ochain-network/types"
)

const (
	OChainPlanetPrefix string = "planet_"
)

type OChainPlanetTable struct {
	bdb        *badger.DB
	currentTxn *badger.Txn
}

func (db *OChainPlanetTable) KeyOf(universeId string, planet types.OChainPlanet) []byte {
	return []byte(OChainPlanetPrefix + universeId + "_" + planet.CoordinateId())
}

func (db *OChainPlanetTable) SetCurrentTxn(tx *badger.Txn) {
	db.currentTxn = tx
}

func (db *OChainPlanetTable) Exists(universeId string, coordinateId string) (bool, error) {
	var at uint64 = math.MaxUint64
	return db.ExistsAt(universeId, coordinateId, at)
}

func (db *OChainPlanetTable) ExistsAt(universeId string, coordinateId string, at uint64) (bool, error) {
	key := []byte(OChainPlanetPrefix + universeId + "_" + coordinateId)

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

func (db *OChainPlanetTable) Get(universeId string, coordinateId string) (types.OChainPlanet, error) {
	var at uint64 = math.MaxUint64
	return db.GetAt(universeId, coordinateId, at)
}

func (db *OChainPlanetTable) GetAt(universeId string, coordinateId string, at uint64) (types.OChainPlanet, error) {
	var planet types.OChainPlanet
	key := []byte(OChainPlanetPrefix + universeId + "_" + coordinateId)
	txn := db.bdb.NewTransactionAt(at, false)

	item, err := txn.Get([]byte(key))
	if err != nil {
		return types.OChainPlanet{}, err
	}

	value, err := item.ValueCopy(nil)
	if err != nil {
		return types.OChainPlanet{}, err
	}

	err = cbor.Unmarshal(value, &planet)
	if err != nil {
		return types.OChainPlanet{}, err
	}

	return planet, nil
}

func (db *OChainPlanetTable) Insert(universeId string, planet types.OChainPlanet) error {
	key := db.KeyOf(universeId, planet)

	exists, err := db.ExistsAt(universeId, planet.CoordinateId(), db.currentTxn.ReadTs())
	if err != nil {
		return err
	}

	if exists {
		return errors.New("planet already exists")
	}

	value, err := cbor.Marshal(planet)
	if err != nil {
		return err
	}

	return db.currentTxn.Set(key, value)
}

func (db *OChainPlanetTable) Update(universeId string, planet types.OChainPlanet) error {
	key := db.KeyOf(universeId, planet)

	exists, err := db.ExistsAt(universeId, planet.CoordinateId(), db.currentTxn.ReadTs())
	if err != nil {
		return err
	}

	if !exists {
		return errors.New("universe doesn't exists")
	}

	value, err := cbor.Marshal(planet)
	if err != nil {
		return err
	}

	return db.currentTxn.Set(key, value)
}

func (db *OChainPlanetTable) Upsert(universeId string, planet types.OChainPlanet) error {
	key := db.KeyOf(universeId, planet)
	value, err := cbor.Marshal(planet)
	if err != nil {
		return err
	}

	return db.currentTxn.Set(key, value)
}

func (db *OChainPlanetTable) Delete(universeId string, coordinateId string) error {

	key := []byte(OChainPlanetPrefix + universeId + "_" + coordinateId)
	return db.currentTxn.Delete(key)
}

func (db *OChainPlanetTable) GetAll(universeId string) ([]types.OChainPlanet, error) {
	var at uint64 = math.MaxUint64
	return db.GetAllAt(universeId, at)
}

func (db *OChainPlanetTable) GetAllAt(universeId string, at uint64) ([]types.OChainPlanet, error) {
	var universes []types.OChainPlanet

	txn := db.bdb.NewTransactionAt(at, false)
	it := txn.NewIterator(badger.DefaultIteratorOptions)
	defer it.Close()

	prefix := []byte(OChainPlanetPrefix + universeId)

	for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
		item := it.Item()

		var universe types.OChainPlanet
		value, err := item.ValueCopy(nil)
		if err != nil {
			return []types.OChainPlanet{}, err
		}

		err = cbor.Unmarshal(value, &universe)
		if err != nil {
			return []types.OChainPlanet{}, err
		}

		universes = append(universes, universe)
	}

	return universes, nil
}

func (db *OChainPlanetTable) GetAllInGalaxy(universeId string, galaxy string) ([]types.OChainPlanet, error) {
	var at uint64 = math.MaxUint64
	return db.GetAllInGalaxyAt(universeId, galaxy, at)
}

func (db *OChainPlanetTable) GetAllInGalaxyAt(universeId string, galaxy string, at uint64) ([]types.OChainPlanet, error) {
	var universes []types.OChainPlanet

	txn := db.bdb.NewTransactionAt(at, false)
	it := txn.NewIterator(badger.DefaultIteratorOptions)
	defer it.Close()
	prefix := []byte(OChainPlanetPrefix + universeId + "_" + galaxy)

	for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
		item := it.Item()

		var universe types.OChainPlanet
		value, err := item.ValueCopy(nil)
		if err != nil {
			return []types.OChainPlanet{}, err
		}

		err = cbor.Unmarshal(value, &universe)
		if err != nil {
			return []types.OChainPlanet{}, err
		}

		universes = append(universes, universe)
	}

	return universes, nil
}

func (db *OChainPlanetTable) GetAllInSolarSystem(universeId string, galaxy string, solarSystem string) ([]types.OChainPlanet, error) {
	var at uint64 = math.MaxUint64
	return db.GetAllInSolarSystemAt(universeId, galaxy, solarSystem, at)
}

func (db *OChainPlanetTable) GetAllInSolarSystemAt(universeId string, galaxy string, solarSystem string, at uint64) ([]types.OChainPlanet, error) {
	var universes []types.OChainPlanet

	txn := db.bdb.NewTransactionAt(at, false)
	it := txn.NewIterator(badger.DefaultIteratorOptions)
	defer it.Close()
	prefix := []byte(OChainPlanetPrefix + universeId + "_" + galaxy + "_" + solarSystem)

	for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
		item := it.Item()

		var universe types.OChainPlanet
		value, err := item.ValueCopy(nil)
		if err != nil {
			return []types.OChainPlanet{}, err
		}

		err = cbor.Unmarshal(value, &universe)
		if err != nil {
			return []types.OChainPlanet{}, err
		}

		universes = append(universes, universe)
	}

	return universes, nil
}

func NewOChainPlanetTable(db *badger.DB) *OChainPlanetTable {
	return &OChainPlanetTable{
		bdb: db,
	}
}
