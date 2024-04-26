package main

import (
	"github.com/dgraph-io/badger"
	"github.com/ochain.gg/ochain-validator-network/database"
	"github.com/ochain.gg/ochain-validator-network/transactions"
	abcitypes "github.com/tendermint/tendermint/abci/types"
	"github.com/timshannon/badgerhold"
)

type OChainValidatorApplication struct {
	store     *badgerhold.Store
	db        *database.OChainDatabase
	currentTx *badger.Txn
}

var _ abcitypes.Application = (*OChainValidatorApplication)(nil)

func NewOChainValidatorApplication(store *badgerhold.Store) *OChainValidatorApplication {
	return &OChainValidatorApplication{
		store: store,
		db:    database.NewOChainDatabase(store),
	}
}

func (OChainValidatorApplication) Info(req abcitypes.RequestInfo) abcitypes.ResponseInfo {
	return abcitypes.ResponseInfo{}
}

func (OChainValidatorApplication) SetOption(req abcitypes.RequestSetOption) abcitypes.ResponseSetOption {
	return abcitypes.ResponseSetOption{}
}

func (app *OChainValidatorApplication) InitChain(req abcitypes.RequestInitChain) abcitypes.ResponseInitChain {

	_, err := app.db.Universes.Get(1)
	if err != nil {
		tx := app.store.Badger().NewTransaction(true)
		app.db.Universes.Insert(DefaultUniverse, tx)
		tx.Commit()
	}

	return abcitypes.ResponseInitChain{}
}

func (app *OChainValidatorApplication) CheckTx(req abcitypes.RequestCheckTx) abcitypes.ResponseCheckTx {
	tx, err := transactions.ParseTransaction(req.Tx)
	if err != nil {
		return abcitypes.ResponseCheckTx{Code: 1, GasWanted: 0}
	}

	code := tx.IsValid()
	return abcitypes.ResponseCheckTx{Code: code, GasWanted: 0}
}

func (app *OChainValidatorApplication) BeginBlock(req abcitypes.RequestBeginBlock) abcitypes.ResponseBeginBlock {

	return abcitypes.ResponseBeginBlock{}
}

func (app *OChainValidatorApplication) DeliverTx(req abcitypes.RequestDeliverTx) abcitypes.ResponseDeliverTx {

	app.currentTx = app.store.Badger().NewTransaction(true)

	tx, err := transactions.ParseTransaction(req.Tx)
	if err != nil {
		return abcitypes.ResponseDeliverTx{Code: 1}
	}

	code := tx.IsValid()
	if code != 0 {
		return abcitypes.ResponseDeliverTx{Code: code}
	}

	switch tx.Type {
	case transactions.NewValidator:
		_, err := transactions.ParseNewValidatorTransaction(tx)
		if err != nil {
			panic(err)
		}

		return abcitypes.ResponseDeliverTx{Code: 0}

	case transactions.RemoveValidator:
		_, err := transactions.ParseRemoveValidatorTransaction(tx)
		if err != nil {
			panic(err)
		}

		return abcitypes.ResponseDeliverTx{Code: 0}
		// case NewEpoch:

		// case NewVault:

		// case UpdateVault:

		// case NewPositionManager:

		// case UpdatePositionManager:

		// case NewOracle:

		// case UpdateOracle:

	}

	app.currentTx.Commit()

	return abcitypes.ResponseDeliverTx{Code: 1}
}

func (app *OChainValidatorApplication) EndBlock(req abcitypes.RequestEndBlock) abcitypes.ResponseEndBlock {

	return abcitypes.ResponseEndBlock{}
}

func (app *OChainValidatorApplication) Commit() abcitypes.ResponseCommit {
	return abcitypes.ResponseCommit{Data: []byte{}}
}

func (app *OChainValidatorApplication) Query(reqQuery abcitypes.RequestQuery) (resQuery abcitypes.ResponseQuery) {
	// resQuery.Key = reqQuery.Data
	// err := app.db.View(func(txn *badger.Txn) error {
	// 	item, err := txn.Get(reqQuery.Data)
	// 	if err != nil && err != badger.ErrKeyNotFound {
	// 		return err
	// 	}
	// 	if err == badger.ErrKeyNotFound {
	// 		resQuery.Log = "does not exist"
	// 	} else {
	// 		return item.Value(func(val []byte) error {
	// 			resQuery.Log = "exists"
	// 			resQuery.Value = val
	// 			return nil
	// 		})
	// 	}
	// 	return nil
	// })
	// if err != nil {
	// 	panic(err)
	// }
	return
}

func (OChainValidatorApplication) ListSnapshots(abcitypes.RequestListSnapshots) abcitypes.ResponseListSnapshots {
	return abcitypes.ResponseListSnapshots{}
}

func (OChainValidatorApplication) OfferSnapshot(abcitypes.RequestOfferSnapshot) abcitypes.ResponseOfferSnapshot {
	return abcitypes.ResponseOfferSnapshot{}
}

func (OChainValidatorApplication) LoadSnapshotChunk(abcitypes.RequestLoadSnapshotChunk) abcitypes.ResponseLoadSnapshotChunk {
	return abcitypes.ResponseLoadSnapshotChunk{}
}

func (OChainValidatorApplication) ApplySnapshotChunk(abcitypes.RequestApplySnapshotChunk) abcitypes.ResponseApplySnapshotChunk {
	return abcitypes.ResponseApplySnapshotChunk{}
}
