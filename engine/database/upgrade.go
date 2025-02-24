package database

import (
	"errors"
	"fmt"
	"math"

	"github.com/dgraph-io/badger/v4"
	"github.com/fxamacker/cbor/v2"
	"github.com/ochain-gg/ochain-network/types"
)

const (
	OChainUpgradePrefix string = "upgrade_"
)

type OChainUpgradeTable struct {
	bdb        *badger.DB
	currentTxn *badger.Txn
}

func (db *OChainUpgradeTable) SetCurrentTxn(tx *badger.Txn) {
	db.currentTxn = tx
}

func (db *OChainUpgradeTable) Exists(id string) (bool, error) {
	var at uint64 = math.MaxUint64
	return db.ExistsAt(id, at)
}

func (db *OChainUpgradeTable) ExistsAt(id string, at uint64) (bool, error) {
	key := []byte(OChainUpgradePrefix + "_" + id)
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

func (db *OChainUpgradeTable) Get(id string) (types.OChainUpgrade, error) {
	var at uint64 = math.MaxUint64
	return db.GetAt(id, at)
}

func (db *OChainUpgradeTable) GetAt(id string, at uint64) (types.OChainUpgrade, error) {
	var upgrade types.OChainUpgrade
	key := []byte(OChainUpgradePrefix + "_" + id)
	txn := db.bdb.NewTransactionAt(at, false)

	item, err := txn.Get([]byte(key))
	if err != nil {
		return types.OChainUpgrade{}, err
	}

	value, err := item.ValueCopy(nil)
	if err != nil {
		return types.OChainUpgrade{}, err
	}

	err = cbor.Unmarshal(value, &upgrade)
	if err != nil {
		return types.OChainUpgrade{}, err
	}

	return upgrade, nil
}

func (db *OChainUpgradeTable) Insert(upgrade types.OChainUpgrade) error {
	key := []byte(OChainUpgradePrefix + "_" + upgrade.Id())

	exists, err := db.ExistsAt(upgrade.Id(), db.currentTxn.ReadTs())
	if err != nil {
		return err
	}

	if exists {
		return errors.New("upgrade already exists")
	}

	value, err := cbor.Marshal(upgrade)
	if err != nil {
		return err
	}

	return db.currentTxn.Set(key, value)
}

func (db *OChainUpgradeTable) Update(upgrade types.OChainUpgrade) error {
	key := []byte(OChainUpgradePrefix + "_" + upgrade.Id())

	exists, err := db.ExistsAt(upgrade.Id(), db.currentTxn.ReadTs())
	if err != nil {
		return err
	}

	if !exists {
		return errors.New("upgrade doesn't exists")
	}

	value, err := cbor.Marshal(upgrade)
	if err != nil {
		return err
	}

	return db.currentTxn.Set(key, value)
}

func (db *OChainUpgradeTable) Upsert(upgrade types.OChainUpgrade) error {
	key := []byte(OChainUpgradePrefix + "_" + upgrade.Id())
	value, err := cbor.Marshal(upgrade)
	if err != nil {
		return err
	}

	return db.currentTxn.Set(key, value)
}

func (db *OChainUpgradeTable) Delete(id string) error {
	key := []byte(OChainUpgradePrefix + "_" + id)
	return db.currentTxn.Delete(key)
}

func (db *OChainUpgradeTable) GetAll() ([]types.OChainUpgrade, error) {
	var at uint64 = math.MaxUint64
	return db.GetAllAt(at)
}

func (db *OChainUpgradeTable) GetAllAt(at uint64) ([]types.OChainUpgrade, error) {
	var upgrades []types.OChainUpgrade

	txn := db.bdb.NewTransactionAt(at, false)
	it := txn.NewIterator(badger.DefaultIteratorOptions)
	defer it.Close()

	prefix := []byte(OChainUpgradePrefix)

	for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
		item := it.Item()

		var upgrade types.OChainUpgrade
		value, err := item.ValueCopy(nil)
		if err != nil {
			return []types.OChainUpgrade{}, err
		}

		err = cbor.Unmarshal(value, &upgrade)
		if err != nil {
			return []types.OChainUpgrade{}, err
		}

		upgrades = append(upgrades, upgrade)
	}

	return upgrades, nil
}

func (db *OChainUpgradeTable) GetByPlanet(universeId string, planetCoordinateId string) ([]types.OChainUpgrade, error) {
	var at uint64 = math.MaxUint64
	return db.GetByPlanetAt(universeId, planetCoordinateId, at)
}

func (db *OChainUpgradeTable) GetByPlanetAt(universeId string, planetCoordinateId string, at uint64) ([]types.OChainUpgrade, error) {
	var upgrades []types.OChainUpgrade

	txn := db.bdb.NewTransactionAt(at, false)
	it := txn.NewIterator(badger.DefaultIteratorOptions)
	defer it.Close()

	prefix := []byte(OChainUpgradePrefix + "_" + universeId + "_" + planetCoordinateId)

	for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
		item := it.Item()

		var upgrade types.OChainUpgrade
		value, err := item.ValueCopy(nil)
		if err != nil {
			return []types.OChainUpgrade{}, err
		}

		err = cbor.Unmarshal(value, &upgrade)
		if err != nil {
			return []types.OChainUpgrade{}, err
		}

		upgrades = append(upgrades, upgrade)
	}

	return upgrades, nil
}

func (db *OChainUpgradeTable) GetBuildingUpgradesByPlanet(universeId string, planetCoordinateId string) ([]types.OChainUpgrade, error) {
	var at uint64 = math.MaxUint64
	return db.GetBuildingUpgradesByPlanetAt(universeId, planetCoordinateId, at)
}

func (db *OChainUpgradeTable) GetBuildingUpgradesByPlanetAt(universeId string, planetCoordinateId string, at uint64) ([]types.OChainUpgrade, error) {
	var upgrades []types.OChainUpgrade

	txn := db.bdb.NewTransactionAt(at, false)
	it := txn.NewIterator(badger.DefaultIteratorOptions)
	defer it.Close()

	prefix := []byte(OChainUpgradePrefix + "_" + universeId + "_" + planetCoordinateId)

	for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
		item := it.Item()

		var upgrade types.OChainUpgrade
		value, err := item.ValueCopy(nil)
		if err != nil {
			return []types.OChainUpgrade{}, err
		}

		err = cbor.Unmarshal(value, &upgrade)
		if err != nil {
			return []types.OChainUpgrade{}, err
		}

		if upgrade.UpgradeType == types.OChainBuildingUpgrade {
			upgrades = append(upgrades, upgrade)
		}
	}

	return upgrades, nil
}

func (db *OChainUpgradeTable) GetPendingBuildingUpgradesByPlanet(universeId string, planetCoordinateId string) ([]types.OChainUpgrade, error) {
	var at uint64 = math.MaxUint64
	return db.GetPendingBuildingUpgradesByPlanetAt(universeId, planetCoordinateId, at)
}

func (db *OChainUpgradeTable) GetPendingBuildingUpgradesByPlanetAt(universeId string, planetCoordinateId string, at uint64) ([]types.OChainUpgrade, error) {
	var upgrades []types.OChainUpgrade

	txn := db.bdb.NewTransactionAt(at, false)
	it := txn.NewIterator(badger.DefaultIteratorOptions)
	defer it.Close()

	prefix := []byte(OChainUpgradePrefix + "_" + universeId + "_" + planetCoordinateId + "_" + fmt.Sprint(types.OChainBuildingUpgrade))

	for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
		item := it.Item()

		var upgrade types.OChainUpgrade
		value, err := item.ValueCopy(nil)
		if err != nil {
			return []types.OChainUpgrade{}, err
		}

		err = cbor.Unmarshal(value, &upgrade)
		if err != nil {
			return []types.OChainUpgrade{}, err
		}

		if !upgrade.Executed {
			upgrades = append(upgrades, upgrade)
		}
	}

	return upgrades, nil
}

func (db *OChainUpgradeTable) GetTechnologyUpgradesByPlanet(universeId string, planetCoordinateId string) ([]types.OChainUpgrade, error) {
	var at uint64 = math.MaxUint64
	return db.GetTechnologyUpgradesByPlanetAt(universeId, planetCoordinateId, at)
}

func (db *OChainUpgradeTable) GetTechnologyUpgradesByPlanetAt(universeId string, planetCoordinateId string, at uint64) ([]types.OChainUpgrade, error) {
	var upgrades []types.OChainUpgrade

	txn := db.bdb.NewTransactionAt(at, false)
	it := txn.NewIterator(badger.DefaultIteratorOptions)
	defer it.Close()

	prefix := []byte(OChainUpgradePrefix + "_" + universeId + "_" + planetCoordinateId + "_" + fmt.Sprint(types.OChainTechnologyUpgrade))

	for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
		item := it.Item()

		var upgrade types.OChainUpgrade
		value, err := item.ValueCopy(nil)
		if err != nil {
			return []types.OChainUpgrade{}, err
		}

		err = cbor.Unmarshal(value, &upgrade)
		if err != nil {
			return []types.OChainUpgrade{}, err
		}

		if upgrade.UpgradeType == types.OChainTechnologyUpgrade {
			upgrades = append(upgrades, upgrade)
		}
	}

	return upgrades, nil
}

func (db *OChainUpgradeTable) GetPendingTechnologyUpgradesByPlanet(universeId string, planetCoordinateId string) ([]types.OChainUpgrade, error) {
	var at uint64 = math.MaxUint64
	return db.GetPendingTechnologyUpgradesByPlanetAt(universeId, planetCoordinateId, at)
}

func (db *OChainUpgradeTable) GetPendingTechnologyUpgradesByPlanetAt(universeId string, planetCoordinateId string, at uint64) ([]types.OChainUpgrade, error) {
	var upgrades []types.OChainUpgrade

	txn := db.bdb.NewTransactionAt(at, false)
	it := txn.NewIterator(badger.DefaultIteratorOptions)
	defer it.Close()

	prefix := []byte(OChainUpgradePrefix + "_" + universeId + "_" + planetCoordinateId + "_" + fmt.Sprint(types.OChainTechnologyUpgrade))

	for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
		item := it.Item()

		var upgrade types.OChainUpgrade
		value, err := item.ValueCopy(nil)
		if err != nil {
			return []types.OChainUpgrade{}, err
		}

		err = cbor.Unmarshal(value, &upgrade)
		if err != nil {
			return []types.OChainUpgrade{}, err
		}

		if !upgrade.Executed {
			upgrades = append(upgrades, upgrade)
		}
	}

	return upgrades, nil
}

func NewOChainUpgradeTable(db *badger.DB) *OChainUpgradeTable {
	return &OChainUpgradeTable{
		bdb: db,
	}
}
