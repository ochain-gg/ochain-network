package main

import (
	"context"
	"encoding/hex"
	"log"
	"math/big"

	"github.com/dgraph-io/badger"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ochain.gg/ochain-validator-network/config"
	"github.com/ochain.gg/ochain-validator-network/contracts"
	"github.com/ochain.gg/ochain-validator-network/database"
	"github.com/ochain.gg/ochain-validator-network/transactions"
	abcitypes "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	"github.com/tendermint/tendermint/proto/tendermint/crypto"
	"github.com/timshannon/badgerhold"
)

type OChainValidatorApplication struct {
	config    config.OChainConfig
	store     *badgerhold.Store
	db        *database.OChainDatabase
	currentTx *badger.Txn

	state      *database.OChainState
	ValUpdates []abcitypes.ValidatorUpdate
}

var _ abcitypes.Application = (*OChainValidatorApplication)(nil)

func NewOChainValidatorApplication(config config.OChainConfig, store *badgerhold.Store) *OChainValidatorApplication {
	return &OChainValidatorApplication{
		config: config,
		store:  store,
		db:     database.NewOChainDatabase(store),
	}
}

func (OChainValidatorApplication) Info(req abcitypes.RequestInfo) abcitypes.ResponseInfo {
	return abcitypes.ResponseInfo{}
}

func (OChainValidatorApplication) SetOption(req abcitypes.RequestSetOption) abcitypes.ResponseSetOption {
	return abcitypes.ResponseSetOption{}
}

func (app *OChainValidatorApplication) InitChain(req abcitypes.RequestInitChain) abcitypes.ResponseInitChain {
	tx := app.store.Badger().NewTransaction(true)

	client, err := ethclient.Dial(app.config.EVMRpc)
	if err != nil {
		log.Fatal(err)
	}

	address := common.HexToAddress(app.config.EVMPortalAddress)
	portal, err := contracts.NewOChainPortal(address, client)
	if err != nil {
		log.Fatal(err)
	}

	co := &bind.CallOpts{
		Context: context.Background(),
	}

	info, err := portal.ValidatorNetworkInfo(co)
	if err != nil {
		log.Fatal(err)
	}

	validators := make([]abcitypes.ValidatorUpdate, 0)
	for i := 0; i < int(info.ValidatorsLength.Int64()); i++ {
		info, err := portal.ValidatorInfo(co, new(big.Int).SetInt64(int64(i)))
		if err != nil {
			log.Fatal(err)
		}

		if info.Enabled {
			pubkeyBytes, err := hex.DecodeString(info.PublicKey)
			if err != nil {
				continue
			}
			var pubkey secp256k1.PubKey = secp256k1.PubKey(pubkeyBytes)

			validators = append(validators, abcitypes.ValidatorUpdate{
				PubKey: crypto.PublicKey{Sum: &crypto.PublicKey_Secp256K1{Secp256K1: pubkey}},
				Power:  10000,
			})
		}
	}

	_, err = app.db.Universes.Get(1)
	if err != nil {
		app.db.Universes.Insert(DefaultUniverse, tx)
	}

	tx.Commit()

	return abcitypes.ResponseInitChain{
		Validators: validators,
	}
}

func (app *OChainValidatorApplication) CheckTx(req abcitypes.RequestCheckTx) abcitypes.ResponseCheckTx {
	tx, err := transactions.ParseTransaction(req.Tx)
	if err != nil {
		return abcitypes.ResponseCheckTx{Code: 1}
	}

	code := tx.IsValid()
	if code != 0 {
		return abcitypes.ResponseCheckTx{Code: code}
	}

	switch tx.Type {
	case transactions.NewValidator:
		tx, err := transactions.ParseNewValidatorTransaction(tx)
		if err != nil {
			return abcitypes.ResponseCheckTx{Code: 1}
		}

		_, err = tx.Verify(app.config)
		if err != nil {
			return abcitypes.ResponseCheckTx{Code: 1}
		}

	case transactions.RemoveValidator:
		tx, err := transactions.ParseRemoveValidatorTransaction(tx)
		if err != nil {
			return abcitypes.ResponseCheckTx{Code: 1}
		}

		_, err = tx.Verify(app.config)
		if err != nil {
			return abcitypes.ResponseCheckTx{Code: 1}
		}
	}

	return abcitypes.ResponseCheckTx{Code: code, GasWanted: 0}
}

func (app *OChainValidatorApplication) BeginBlock(req abcitypes.RequestBeginBlock) abcitypes.ResponseBeginBlock {
	app.ValUpdates = make([]abcitypes.ValidatorUpdate, 0)
	app.currentTx = app.store.Badger().NewTransaction(true)
	return abcitypes.ResponseBeginBlock{}
}

func (app *OChainValidatorApplication) DeliverTx(req abcitypes.RequestDeliverTx) abcitypes.ResponseDeliverTx {

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
		formatedTx, err := transactions.ParseNewValidatorTransaction(tx)
		if err != nil {
			return abcitypes.ResponseDeliverTx{Code: 1}
		}
		_, err = formatedTx.Verify(app.config)
		if err != nil {
			return abcitypes.ResponseDeliverTx{Code: 1}
		}

		err = formatedTx.Execute(app.config, app.db, app.currentTx)
		if err != nil {
			return abcitypes.ResponseDeliverTx{Code: 1}
		}

		pubkeyBytes, err := hex.DecodeString(formatedTx.FormatedData.PublicKey)
		if err != nil {
			return abcitypes.ResponseDeliverTx{Code: 1}
		}

		app.ValUpdates = append(app.ValUpdates, abcitypes.UpdateValidator(pubkeyBytes, 10000, "secp256k1"))
		return abcitypes.ResponseDeliverTx{Code: 0}

	case transactions.RemoveValidator:
		formatedTx, err := transactions.ParseRemoveValidatorTransaction(tx)
		if err != nil {
			return abcitypes.ResponseDeliverTx{Code: 1}
		}
		_, err = formatedTx.Verify(app.config)
		if err != nil {
			return abcitypes.ResponseDeliverTx{Code: 1}
		}

		err = formatedTx.Execute(app.config, app.db, app.currentTx)
		if err != nil {
			return abcitypes.ResponseDeliverTx{Code: 1}
		}

		validator, err := app.db.Validators.Get(formatedTx.FormatedData.ValidatorId, app.currentTx)
		if err != nil {
			return abcitypes.ResponseDeliverTx{Code: 1}
		}

		pubkeyBytes, err := hex.DecodeString(validator.PublicKey)
		if err != nil {
			return abcitypes.ResponseDeliverTx{Code: 1}
		}

		app.ValUpdates = append(app.ValUpdates, abcitypes.UpdateValidator(pubkeyBytes, 0, "secp256k1"))
		return abcitypes.ResponseDeliverTx{Code: 0}
	}

	return abcitypes.ResponseDeliverTx{Code: 1}
}

func (app *OChainValidatorApplication) EndBlock(req abcitypes.RequestEndBlock) abcitypes.ResponseEndBlock {
	app.currentTx.Commit()
	return abcitypes.ResponseEndBlock{ValidatorUpdates: app.ValUpdates}
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
