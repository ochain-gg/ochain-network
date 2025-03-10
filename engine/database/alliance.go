package database

import (
	"context"
	"errors"
	"math"

	"github.com/dgraph-io/badger/pb"
	"github.com/dgraph-io/badger/v4"
	"github.com/dgraph-io/ristretto/v2/z"
	"github.com/fxamacker/cbor/v2"
	"github.com/ochain-gg/ochain-network/types"
)

const (
	OChainAlliancePrefix            string = "alliance_"
	OChainAllianceJoinRequestPrefix string = "ajr_" // Alliance Join Request prefix
)

type OChainAllianceTable struct {
	bdb        *badger.DB
	currentTxn *badger.Txn
}

func (db *OChainAllianceTable) SetCurrentTxn(tx *badger.Txn) {
	db.currentTxn = tx
}

func (db *OChainAllianceTable) Exists(id string) (bool, error) {
	var at uint64 = math.MaxUint64
	return db.ExistsAt(id, at)
}

func (db *OChainAllianceTable) ExistsAt(id string, at uint64) (bool, error) {
	key := []byte(OChainAlliancePrefix + id)

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

func (db *OChainAllianceTable) Get(id string) (types.OChainAlliance, error) {
	var at uint64 = math.MaxUint64
	return db.GetAt(id, at)
}

func (db *OChainAllianceTable) GetAt(id string, at uint64) (types.OChainAlliance, error) {
	var planet types.OChainAlliance
	key := []byte(OChainAlliancePrefix + id)
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

func (db *OChainAllianceTable) GetJoinRequest(id string) (types.OChainAllianceJoinRequest, error) {
	var at uint64 = math.MaxUint64
	return db.GetJoinRequestAt(id, at)
}

func (db *OChainAllianceTable) GetJoinRequestAt(id string, at uint64) (types.OChainAllianceJoinRequest, error) {
	var request types.OChainAllianceJoinRequest
	key := []byte(OChainAllianceJoinRequestPrefix + id)
	txn := db.bdb.NewTransactionAt(at, false)

	item, err := txn.Get([]byte(key))
	if err != nil {
		return types.OChainAllianceJoinRequest{}, err
	}

	value, err := item.ValueCopy(nil)
	if err != nil {
		return types.OChainAllianceJoinRequest{}, err
	}

	err = cbor.Unmarshal(value, &request)
	if err != nil {
		return types.OChainAllianceJoinRequest{}, err
	}

	return request, nil
}

func (db *OChainAllianceTable) Insert(alliance types.OChainAlliance) error {
	key := []byte(OChainAlliancePrefix + alliance.Id)

	exists, err := db.ExistsAt(alliance.Id, db.currentTxn.ReadTs())
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

func (db *OChainAllianceTable) Update(alliance types.OChainAlliance) error {
	key := []byte(OChainAlliancePrefix + alliance.Id)

	exists, err := db.ExistsAt(alliance.Id, db.currentTxn.ReadTs())
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

func (db *OChainAllianceTable) Upsert(alliance types.OChainAlliance) error {
	key := []byte(OChainAlliancePrefix + alliance.Id)
	value, err := cbor.Marshal(alliance)
	if err != nil {
		return err
	}

	return db.currentTxn.Set(key, value)
}

func (db *OChainAllianceTable) Delete(id string) error {
	key := []byte(OChainAlliancePrefix + id)
	return db.currentTxn.Delete(key)
}

func (db *OChainAllianceTable) GetAll() ([]types.OChainAlliance, error) {
	var at uint64 = math.MaxUint64
	return db.GetAllAt(at)
}

func (db *OChainAllianceTable) GetAllAt(at uint64) ([]types.OChainAlliance, error) {
	var universes []types.OChainAlliance

	txn := db.bdb.NewTransactionAt(at, false)
	it := txn.NewIterator(badger.DefaultIteratorOptions)
	defer it.Close()

	prefix := []byte(OChainAlliancePrefix)

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

func (t *OChainAllianceTable) GetByUniverseAt(universeId string, onlyNotAnswered bool, at uint64) ([]types.OChainAllianceJoinRequest, error) {
	var requests []types.OChainAllianceJoinRequest

	stream := t.bdb.NewStreamAt(at)

	stream.NumGo = 16                                               // Set number of goroutines to use for iteration.
	stream.Prefix = []byte(OChainAllianceJoinRequestPrefix)         // Leave nil for iteration over the whole DB.
	stream.LogPrefix = "Badger.GetJoinRequestsByAlliance.Streaming" // For identifying stream logs. Outputs to Logger.
	stream.Send = func(buf *z.Buffer) error {
		err := buf.SliceIterate(func(s []byte) error {
			var kv pb.KV
			var req types.OChainAllianceJoinRequest
			if err := kv.Unmarshal(s); err != nil {
				return err
			}

			if err := cbor.Unmarshal(kv.Value, &req); err != nil {
				return err
			}

			if !onlyNotAnswered || req.AnsweredAt == 0 {
				requests = append(requests, req)
			}

			return nil
		})

		return err
	}

	if err := stream.Orchestrate(context.Background()); err != nil {
		return []types.OChainAllianceJoinRequest{}, err
	}

	return requests, nil
}

func NewOChainAllianceTable(db *badger.DB) *OChainAllianceTable {
	return &OChainAllianceTable{
		bdb: db,
	}
}

// HasPendingRequestAt checks if a player has a pending join request for an alliance
func (t *OChainAllianceTable) HasPendingRequest(from string) (bool, error) {
	var at uint64 = math.MaxUint64
	return t.HasPendingRequestAt(from, at)
}

// HasPendingRequestAt checks if a player has a pending join request for an alliance
func (t *OChainAllianceTable) HasPendingRequestAt(from string, at uint64) (bool, error) {
	res, err := t.GetJoinRequestByAccountAt(from, true, at)
	if err != nil {
		return false, err
	}

	return len(res) > 0, nil
}

// InsertJoinRequest adds a new join request to the alliance
func (t *OChainAllianceTable) InsertJoinRequest(request types.OChainAllianceJoinRequest) error {
	if t.currentTxn == nil {
		return errors.New("no transaction in progress")
	}

	key := []byte(OChainAllianceJoinRequestPrefix + request.Id)
	value, err := cbor.Marshal(request)
	if err != nil {
		return err
	}

	return t.currentTxn.Set(key, value)
}

// GetJoinRequest gets a specific join request
func (t *OChainAllianceTable) GetJoinRequestsByAlliance(allianceId string, onlyNotAnswered bool) ([]types.OChainAllianceJoinRequest, error) {
	var at uint64 = math.MaxUint64
	return t.GetJoinRequestsByAllianceAt(allianceId, onlyNotAnswered, at)
}

func (t *OChainAllianceTable) GetJoinRequestsByAllianceAt(allianceId string, onlyNotAnswered bool, at uint64) ([]types.OChainAllianceJoinRequest, error) {
	var requests []types.OChainAllianceJoinRequest

	stream := t.bdb.NewStreamAt(at)

	stream.NumGo = 16                                               // Set number of goroutines to use for iteration.
	stream.Prefix = []byte(OChainAllianceJoinRequestPrefix)         // Leave nil for iteration over the whole DB.
	stream.LogPrefix = "Badger.GetJoinRequestsByAlliance.Streaming" // For identifying stream logs. Outputs to Logger.
	stream.Send = func(buf *z.Buffer) error {
		err := buf.SliceIterate(func(s []byte) error {
			var kv pb.KV
			var req types.OChainAllianceJoinRequest
			if err := kv.Unmarshal(s); err != nil {
				return err
			}

			if err := cbor.Unmarshal(kv.Value, &req); err != nil {
				return err
			}

			if !onlyNotAnswered || req.AnsweredAt == 0 {
				requests = append(requests, req)
			}

			return nil
		})

		return err
	}

	if err := stream.Orchestrate(context.Background()); err != nil {
		return []types.OChainAllianceJoinRequest{}, err
	}

	return requests, nil
}

// GetJoinRequest gets a specific join request
func (t *OChainAllianceTable) GetJoinRequestByAccount(from string, onlyNotAnswered bool) ([]types.OChainAllianceJoinRequest, error) {
	var at uint64 = math.MaxUint64
	return t.GetJoinRequestByAccountAt(from, onlyNotAnswered, at)
}

func (t *OChainAllianceTable) GetJoinRequestByAccountAt(from string, onlyNotAnswered bool, at uint64) ([]types.OChainAllianceJoinRequest, error) {
	var requests []types.OChainAllianceJoinRequest

	stream := t.bdb.NewStreamAt(at)

	stream.NumGo = 16                                             // Set number of goroutines to use for iteration.
	stream.Prefix = []byte(OChainAllianceJoinRequestPrefix)       // Leave nil for iteration over the whole DB.
	stream.LogPrefix = "Badger.GetJoinRequestByAccount.Streaming" // For identifying stream logs. Outputs to Logger.
	stream.Send = func(buf *z.Buffer) error {
		err := buf.SliceIterate(func(s []byte) error {
			var kv pb.KV
			var req types.OChainAllianceJoinRequest
			if err := kv.Unmarshal(s); err != nil {
				return err
			}

			if err := cbor.Unmarshal(kv.Value, &req); err != nil {
				return err
			}

			if req.From == from && (!onlyNotAnswered || req.AnsweredAt == 0) {
				requests = append(requests, req)
			}

			return nil
		})

		return err
	}

	if err := stream.Orchestrate(context.Background()); err != nil {
		return []types.OChainAllianceJoinRequest{}, err
	}

	return requests, nil
}

// UpdateJoinRequest updates an existing join request
func (t *OChainAllianceTable) UpdateJoinRequest(request types.OChainAllianceJoinRequest) error {
	if t.currentTxn == nil {
		return errors.New("no transaction in progress")
	}

	key := []byte(OChainAllianceJoinRequestPrefix + request.Id)
	value, err := cbor.Marshal(request)
	if err != nil {
		return err
	}

	return t.currentTxn.Set(key, value)
}
