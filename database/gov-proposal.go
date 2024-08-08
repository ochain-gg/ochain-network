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
	OChainGovernanceProposalPrefix string = "gov_proposals_"
)

type OChainGovernanceProposalTable struct {
	bdb        *badger.DB
	currentTxn *badger.Txn
}

func (db *OChainGovernanceProposalTable) SetCurrentTxn(tx *badger.Txn) {
	db.currentTxn = tx
}

func (db *OChainGovernanceProposalTable) Exists(id uint64) (bool, error) {
	var at uint64 = math.MaxUint64
	return db.ExistsAt(id, at)
}

func (db *OChainGovernanceProposalTable) ExistsAt(id uint64, at uint64) (bool, error) {
	key := []byte(OChainGovernanceProposalPrefix + fmt.Sprint(id))
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

func (db *OChainGovernanceProposalTable) Get(id string) (types.OChainGovernanceProposal, error) {
	var at uint64 = math.MaxUint64
	return db.GetAt(id, at)
}

func (db *OChainGovernanceProposalTable) GetAt(id string, at uint64) (types.OChainGovernanceProposal, error) {
	var epoch types.OChainGovernanceProposal
	key := []byte(OChainGovernanceProposalPrefix + fmt.Sprint(id))
	txn := db.bdb.NewTransactionAt(at, false)

	item, err := txn.Get([]byte(key))
	if err != nil {
		return types.OChainGovernanceProposal{}, err
	}

	value, err := item.ValueCopy(nil)
	if err != nil {
		return types.OChainGovernanceProposal{}, err
	}

	err = cbor.Unmarshal(value, &epoch)
	if err != nil {
		return types.OChainGovernanceProposal{}, err
	}

	return epoch, nil
}

func (db *OChainGovernanceProposalTable) Insert(epoch types.OChainGovernanceProposal) error {
	key := []byte(OChainGovernanceProposalPrefix + fmt.Sprint(epoch.Id))

	exists, err := db.ExistsAt(epoch.Id, db.currentTxn.ReadTs())
	if err != nil {
		return err
	}

	if exists {
		return errors.New("global epoch already exists")
	}

	value, err := cbor.Marshal(epoch)
	if err != nil {
		return err
	}

	return db.currentTxn.Set(key, value)
}

func (db *OChainGovernanceProposalTable) Update(epoch types.OChainGovernanceProposal) error {
	key := []byte(OChainGovernanceProposalPrefix + fmt.Sprint(epoch.Id))

	exists, err := db.ExistsAt(epoch.Id, db.currentTxn.ReadTs())
	if err != nil {
		return err
	}

	if !exists {
		return errors.New("epoch doesn't exists")
	}

	value, err := cbor.Marshal(epoch)
	if err != nil {
		return err
	}

	return db.currentTxn.Set(key, value)
}

func (db *OChainGovernanceProposalTable) Upsert(epoch types.OChainGovernanceProposal) error {
	key := []byte(OChainGovernanceProposalPrefix + fmt.Sprint(epoch.Id))
	value, err := cbor.Marshal(epoch)
	if err != nil {
		return err
	}

	return db.currentTxn.Set(key, value)
}

func (db *OChainGovernanceProposalTable) Delete(id string) error {
	key := []byte(OChainGovernanceProposalPrefix + fmt.Sprint(id))
	return db.currentTxn.Delete(key)
}

func (db *OChainGovernanceProposalTable) GetAll() ([]types.OChainGovernanceProposal, error) {
	var at uint64 = math.MaxUint64
	return db.GetAllAt(at)
}

func (db *OChainGovernanceProposalTable) GetAllAt(at uint64) ([]types.OChainGovernanceProposal, error) {
	var epochs []types.OChainGovernanceProposal

	txn := db.bdb.NewTransactionAt(at, false)
	it := txn.NewIterator(badger.DefaultIteratorOptions)
	defer it.Close()

	prefix := []byte(OChainGovernanceProposalPrefix)

	for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
		item := it.Item()

		var epoch types.OChainGovernanceProposal
		value, err := item.ValueCopy(nil)
		if err != nil {
			return []types.OChainGovernanceProposal{}, err
		}

		err = cbor.Unmarshal(value, &epoch)
		if err != nil {
			return []types.OChainGovernanceProposal{}, err
		}

		epochs = append(epochs, epoch)
	}

	return epochs, nil
}

func NewOChainGovernanceProposalTable(db *badger.DB) *OChainGovernanceProposalTable {
	return &OChainGovernanceProposalTable{
		bdb: db,
	}
}
