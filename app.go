package main

import (
	"context"
	"encoding/hex"
	"log"
	"math/big"

	abcitypes "github.com/cometbft/cometbft/abci/types"
	"github.com/cometbft/cometbft/crypto/secp256k1"
	"github.com/cometbft/cometbft/proto/tendermint/crypto"
	"github.com/dgraph-io/badger/v4"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ochain.gg/ochain-network-validator/config"
	"github.com/ochain.gg/ochain-network-validator/contracts"
	"github.com/ochain.gg/ochain-network-validator/database"
	"github.com/ochain.gg/ochain-network-validator/transactions"
	"github.com/timshannon/badgerhold/v4"
)

type OChainValidatorApplication struct {
	config         config.OChainConfig
	store          *badgerhold.Store
	db             *database.OChainDatabase
	ongoingBlockTx *badger.Txn

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

func (app *OChainValidatorApplication) Info(_ context.Context, info *abcitypes.RequestInfo) (*abcitypes.ResponseInfo, error) {
	return &abcitypes.ResponseInfo{}, nil
}

func (app *OChainValidatorApplication) InitChain(_ context.Context, chain *abcitypes.RequestInitChain) (*abcitypes.ResponseInitChain, error) {

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

	return &abcitypes.ResponseInitChain{
		Validators: validators,
	}, nil
}

func (app *OChainValidatorApplication) CheckTx(_ context.Context, req *abcitypes.RequestCheckTx) (*abcitypes.ResponseCheckTx, error) {
	tx, err := transactions.ParseTransaction(req.Tx)
	if err != nil {
		return &abcitypes.ResponseCheckTx{Code: 1}, nil
	}

	code := tx.IsValid()
	if code != 0 {
		return &abcitypes.ResponseCheckTx{Code: code}, nil
	}

	switch tx.Type {
	case transactions.NewValidator:
		tx, err := transactions.ParseNewValidatorTransaction(tx)
		if err != nil {
			return &abcitypes.ResponseCheckTx{Code: 1}, nil
		}

		_, err = tx.Verify(app.config)
		if err != nil {
			return &abcitypes.ResponseCheckTx{Code: 1}, nil
		}

	case transactions.RemoveValidator:
		tx, err := transactions.ParseRemoveValidatorTransaction(tx)
		if err != nil {
			return &abcitypes.ResponseCheckTx{Code: 1}, nil
		}

		_, err = tx.Verify(app.config)
		if err != nil {
			return &abcitypes.ResponseCheckTx{Code: 1}, nil
		}
	}

	return &abcitypes.ResponseCheckTx{Code: code, GasWanted: 0}, nil
}

func (app *OChainValidatorApplication) PrepareProposal(_ context.Context, proposal *abcitypes.RequestPrepareProposal) (*abcitypes.ResponsePrepareProposal, error) {
	return &abcitypes.ResponsePrepareProposal{}, nil
}

func (app *OChainValidatorApplication) ProcessProposal(_ context.Context, proposal *abcitypes.RequestProcessProposal) (*abcitypes.ResponseProcessProposal, error) {
	return &abcitypes.ResponseProcessProposal{}, nil
}

func (app *OChainValidatorApplication) FinalizeBlock(_ context.Context, req *abcitypes.RequestFinalizeBlock) (*abcitypes.ResponseFinalizeBlock, error) {
	app.ValUpdates = make([]abcitypes.ValidatorUpdate, 0)
	app.ongoingBlockTx = app.store.Badger().NewTransaction(true)

	var txs = make([]*abcitypes.ExecTxResult, len(req.Txs))

	for i, tx := range req.Txs {

		parsedTx, err := transactions.ParseTransaction(tx)
		if err != nil {
			txs[i] = &abcitypes.ExecTxResult{Code: 1}
			continue
		}

		code := parsedTx.IsValid()
		if code != 0 {
			txs[i] = &abcitypes.ExecTxResult{Code: code}
			continue
		}

		switch parsedTx.Type {
		case transactions.NewValidator:
			formatedTx, err := transactions.ParseNewValidatorTransaction(parsedTx)
			if err != nil {
				txs[i] = &abcitypes.ExecTxResult{Code: 1}
				continue
			}

			_, err = formatedTx.Verify(app.config)
			if err != nil {
				txs[i] = &abcitypes.ExecTxResult{Code: 1}
				continue
			}

			err = formatedTx.Execute(app.config, app.db, app.ongoingBlockTx)
			if err != nil {
				txs[i] = &abcitypes.ExecTxResult{Code: 1}
				continue
			}

			pubkeyBytes, err := hex.DecodeString(formatedTx.FormatedData.PublicKey)
			if err != nil {
				txs[i] = &abcitypes.ExecTxResult{Code: 1}
				continue
			}

			app.ValUpdates = append(app.ValUpdates, abcitypes.UpdateValidator(pubkeyBytes, 10000, "secp256k1"))
			txs[i] = &abcitypes.ExecTxResult{Code: 0}

		case transactions.RemoveValidator:
			formatedTx, err := transactions.ParseRemoveValidatorTransaction(parsedTx)
			if err != nil {
				txs[i] = &abcitypes.ExecTxResult{Code: 1}
				continue
			}
			_, err = formatedTx.Verify(app.config)
			if err != nil {
				txs[i] = &abcitypes.ExecTxResult{Code: 1}
				continue
			}

			err = formatedTx.Execute(app.config, app.db, app.ongoingBlockTx)
			if err != nil {
				txs[i] = &abcitypes.ExecTxResult{Code: 1}
				continue
			}

			validator, err := app.db.Validators.Get(formatedTx.FormatedData.ValidatorId, app.ongoingBlockTx)
			if err != nil {
				txs[i] = &abcitypes.ExecTxResult{Code: 1}
				continue
			}

			pubkeyBytes, err := hex.DecodeString(validator.PublicKey)
			if err != nil {
				txs[i] = &abcitypes.ExecTxResult{Code: 1}
				continue
			}

			app.ValUpdates = append(app.ValUpdates, abcitypes.UpdateValidator(pubkeyBytes, 0, "secp256k1"))
			txs[i] = &abcitypes.ExecTxResult{Code: 1}

		}
	}

	return &abcitypes.ResponseFinalizeBlock{
		TxResults:        txs,
		ValidatorUpdates: app.ValUpdates,
	}, nil
}

func (app *OChainValidatorApplication) Commit(_ context.Context, commit *abcitypes.RequestCommit) (*abcitypes.ResponseCommit, error) {
	return &abcitypes.ResponseCommit{}, app.ongoingBlockTx.Commit()
}

func (app *OChainValidatorApplication) Query(_ context.Context, req *abcitypes.RequestQuery) (*abcitypes.ResponseQuery, error) {
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
	return &abcitypes.ResponseQuery{}, nil
}
func (app *OChainValidatorApplication) ListSnapshots(_ context.Context, snapshots *abcitypes.RequestListSnapshots) (*abcitypes.ResponseListSnapshots, error) {
	return &abcitypes.ResponseListSnapshots{}, nil
}

func (app *OChainValidatorApplication) OfferSnapshot(_ context.Context, snapshot *abcitypes.RequestOfferSnapshot) (*abcitypes.ResponseOfferSnapshot, error) {
	return &abcitypes.ResponseOfferSnapshot{}, nil
}

func (app *OChainValidatorApplication) LoadSnapshotChunk(_ context.Context, chunk *abcitypes.RequestLoadSnapshotChunk) (*abcitypes.ResponseLoadSnapshotChunk, error) {
	return &abcitypes.ResponseLoadSnapshotChunk{}, nil
}

func (app *OChainValidatorApplication) ApplySnapshotChunk(_ context.Context, chunk *abcitypes.RequestApplySnapshotChunk) (*abcitypes.ResponseApplySnapshotChunk, error) {
	return &abcitypes.ResponseApplySnapshotChunk{Result: abcitypes.ResponseApplySnapshotChunk_ACCEPT}, nil
}

func (app *OChainValidatorApplication) ExtendVote(_ context.Context, extend *abcitypes.RequestExtendVote) (*abcitypes.ResponseExtendVote, error) {
	return &abcitypes.ResponseExtendVote{}, nil
}

func (app *OChainValidatorApplication) VerifyVoteExtension(_ context.Context, verify *abcitypes.RequestVerifyVoteExtension) (*abcitypes.ResponseVerifyVoteExtension, error) {
	return &abcitypes.ResponseVerifyVoteExtension{}, nil
}
