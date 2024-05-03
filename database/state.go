package database

import (
	"encoding/hex"
	"log"
	"strconv"

	"github.com/dgraph-io/badger/v4"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/timshannon/badgerhold/v4"
)

type OChainStateTable struct {
	db *badgerhold.Store
}

type OChainState struct {
	Id     uint64 `badgerhold:"key"`
	Size   int64
	Height int64

	LatestPortalUpdate uint64
}

func (state *OChainState) Hash() []byte {
	hash := crypto.Keccak256Hash([]byte(strconv.FormatInt(state.Size, 16))).Bytes()
	log.Printf("State hash processed at size %d: %s", state.Size, hex.EncodeToString(hash))
	return hash
}

func (state *OChainState) SetHeight(height int64) {
	state.Height = height
}

func (state *OChainState) IncSize() {
	log.Printf("IncSize called: before %d", state.Size)
	state.Size = state.Size + 1
	log.Printf("IncSize called: after %d", state.Size)
}

func (table *OChainStateTable) Get() (OChainState, error) {
	var result []OChainState
	err := table.db.Find(&result, badgerhold.Where("Id").Eq(uint64(1)))
	if err != nil {
		return OChainState{}, err
	}

	log.Println("Loading state: ")
	log.Println(result)

	if len(result) > 0 {
		return result[0], nil
	} else {
		return OChainState{
			Size:               0,
			Height:             0,
			LatestPortalUpdate: 0,
		}, nil
	}
}

func (table *OChainStateTable) Save(state *OChainState, txn *badger.Txn) error {
	err := table.db.TxUpsert(txn, 1, state)
	return err
}

func NewOChainStateTable(db *badgerhold.Store) *OChainStateTable {
	return &OChainStateTable{
		db: db,
	}
}
