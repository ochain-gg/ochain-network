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
	OChainGovernanceProposalVotePrefix string = "gov_proposal_votes_"
)

type OChainGovernanceProposalVoteTable struct {
	bdb        *badger.DB
	currentTxn *badger.Txn
}

func (db *OChainGovernanceProposalVoteTable) SetCurrentTxn(tx *badger.Txn) {
	db.currentTxn = tx
}

func (db *OChainGovernanceProposalVoteTable) Exists(id uint64) (bool, error) {
	var at uint64 = math.MaxUint64
	return db.ExistsAt(id, at)
}

func (db *OChainGovernanceProposalVoteTable) ExistsAt(id uint64, at uint64) (bool, error) {
	key := []byte(OChainGovernanceProposalVotePrefix + fmt.Sprint(id))
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

func (db *OChainGovernanceProposalVoteTable) Get(id string) (types.OChainGovernanceProposalVote, error) {
	var at uint64 = math.MaxUint64
	return db.GetAt(id, at)
}

func (db *OChainGovernanceProposalVoteTable) GetAt(id string, at uint64) (types.OChainGovernanceProposalVote, error) {
	var epoch types.OChainGovernanceProposalVote
	key := []byte(OChainGovernanceProposalVotePrefix + fmt.Sprint(id))
	txn := db.bdb.NewTransactionAt(at, false)

	item, err := txn.Get([]byte(key))
	if err != nil {
		return types.OChainGovernanceProposalVote{}, err
	}

	value, err := item.ValueCopy(nil)
	if err != nil {
		return types.OChainGovernanceProposalVote{}, err
	}

	err = cbor.Unmarshal(value, &epoch)
	if err != nil {
		return types.OChainGovernanceProposalVote{}, err
	}

	return epoch, nil
}

func (db *OChainGovernanceProposalVoteTable) Insert(epoch types.OChainGovernanceProposalVote) error {
	key := []byte(OChainGovernanceProposalVotePrefix + fmt.Sprint(epoch.Id))

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

func (db *OChainGovernanceProposalVoteTable) Update(epoch types.OChainGovernanceProposalVote) error {
	key := []byte(OChainGovernanceProposalVotePrefix + fmt.Sprint(epoch.Id))

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

func (db *OChainGovernanceProposalVoteTable) Upsert(epoch types.OChainGovernanceProposalVote) error {
	key := []byte(OChainGovernanceProposalVotePrefix + fmt.Sprint(epoch.Id))
	value, err := cbor.Marshal(epoch)
	if err != nil {
		return err
	}

	return db.currentTxn.Set(key, value)
}

func (db *OChainGovernanceProposalVoteTable) Delete(id string) error {
	key := []byte(OChainGovernanceProposalVotePrefix + fmt.Sprint(id))
	return db.currentTxn.Delete(key)
}

func (db *OChainGovernanceProposalVoteTable) GetAll() ([]types.OChainGovernanceProposalVote, error) {
	var at uint64 = math.MaxUint64
	return db.GetAllAt(at)
}

func (db *OChainGovernanceProposalVoteTable) GetAllAt(at uint64) ([]types.OChainGovernanceProposalVote, error) {
	var epochs []types.OChainGovernanceProposalVote

	txn := db.bdb.NewTransactionAt(at, false)
	it := txn.NewIterator(badger.DefaultIteratorOptions)
	defer it.Close()

	prefix := []byte(OChainGovernanceProposalVotePrefix)

	for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
		item := it.Item()

		var epoch types.OChainGovernanceProposalVote
		value, err := item.ValueCopy(nil)
		if err != nil {
			return []types.OChainGovernanceProposalVote{}, err
		}

		err = cbor.Unmarshal(value, &epoch)
		if err != nil {
			return []types.OChainGovernanceProposalVote{}, err
		}

		epochs = append(epochs, epoch)
	}

	return epochs, nil
}

func NewOChainGovernanceProposalVoteTable(db *badger.DB) *OChainGovernanceProposalVoteTable {
	return &OChainGovernanceProposalVoteTable{
		bdb: db,
	}
}
