package main

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"math/big"
	"strconv"
	"time"

	abcitypes "github.com/cometbft/cometbft/abci/types"
	"github.com/cometbft/cometbft/crypto/ed25519"
	"github.com/cometbft/cometbft/proto/tendermint/crypto"
	"github.com/cometbft/cometbft/version"
	"github.com/dgraph-io/badger/v4"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ochain.gg/ochain-network-validator/config"
	"github.com/ochain.gg/ochain-network-validator/contracts"
	"github.com/ochain.gg/ochain-network-validator/database"
	"github.com/ochain.gg/ochain-network-validator/transactions"
	"github.com/ochain.gg/ochain-network-validator/types"
	"github.com/timshannon/badgerhold/v4"
)

const (
	AppVersion uint64 = 1
)

type OChainValidatorApplication struct {
	abcitypes.BaseApplication

	config         config.OChainConfig
	store          *badgerhold.Store
	db             *database.OChainDatabase
	ongoingBlockTx *badger.Txn

	remotePrivateKey []byte

	state      *database.OChainState
	ValUpdates []abcitypes.ValidatorUpdate
}

var _ abcitypes.Application = (*OChainValidatorApplication)(nil)

func NewOChainValidatorApplication(config config.OChainConfig, store *badgerhold.Store, remotePrivateKey []byte) (*OChainValidatorApplication, error) {
	db := database.NewOChainDatabase(store)
	state, err := db.State.Get()
	if err != nil {
		return &OChainValidatorApplication{}, errors.New("state unaivalable")
	}

	return &OChainValidatorApplication{
		config:           config,
		store:            store,
		state:            &state,
		db:               db,
		remotePrivateKey: remotePrivateKey,
	}, nil
}

func (app *OChainValidatorApplication) Info(_ context.Context, info *abcitypes.RequestInfo) (*abcitypes.ResponseInfo, error) {

	return &abcitypes.ResponseInfo{
		Data:             fmt.Sprintf("{\"size\":%v}", app.state.Size),
		Version:          version.ABCIVersion,
		AppVersion:       AppVersion,
		LastBlockHeight:  app.state.Height,
		LastBlockAppHash: app.state.Hash(),
	}, nil

}

func (app *OChainValidatorApplication) InitChain(_ context.Context, chain *abcitypes.RequestInitChain) (*abcitypes.ResponseInitChain, error) {

	tx := app.store.Badger().NewTransaction(true)

	client, err := ethclient.Dial(app.config.EVMRpc)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Requesting validators on rpc: " + app.config.EVMRpc)
	log.Println("Requesting validators on portal: " + app.config.EVMPortalAddress)

	address := common.HexToAddress(app.config.EVMPortalAddress)
	portal, err := contracts.NewOChainPortal(address, client)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Requesting validators info")
	info, err := portal.ValidatorNetworkInfo(nil)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Requesting info retrieved")

	validators := make([]abcitypes.ValidatorUpdate, 0)
	for i := 0; i < int(info.ValidatorsLength.Int64()); i++ {
		info, err := portal.ValidatorInfo(nil, new(big.Int).SetInt64(int64(i)))
		if err != nil {
			log.Fatal(err)
		}

		log.Println("Requesting validator info")
		log.Println("Validator Id: " + strconv.Itoa(i))
		log.Println("Validator PubKey: " + info.PublicKey)
		log.Println("Enabled: " + strconv.FormatBool(info.Enabled))

		if info.Enabled {
			pubkeyBytes, err := hex.DecodeString(info.PublicKey)
			if err != nil {
				continue
			}
			var pubkey ed25519.PubKey = ed25519.PubKey(pubkeyBytes)

			validators = append(validators, abcitypes.ValidatorUpdate{
				PubKey: crypto.PublicKey{Sum: &crypto.PublicKey_Ed25519{Ed25519: pubkey}},
				Power:  10000,
			})

			log.Println("adding validator address: " + pubkey.Address().String())
		}
	}

	_, err = app.db.Universes.Get(1)
	if err != nil {
		app.db.Universes.Insert(DefaultUniverse, tx)
	}

	tx.Commit()

	return &abcitypes.ResponseInitChain{
		Validators: validators,
		AppHash:    app.state.Hash(),
	}, nil
}

func (app *OChainValidatorApplication) CheckTx(ctx context.Context, req *abcitypes.RequestCheckTx) (*abcitypes.ResponseCheckTx, error) {

	txCtx := transactions.TransactionContext{
		Config: app.config,
		Db:     app.db,
		Txn:    app.store.Badger().NewTransaction(true),
		Date:   time.Now(),
	}
	log.Printf("Check tx: %s", hex.EncodeToString(req.Tx))

	tx, err := transactions.ParseTransaction(req.Tx)
	if err != nil {
		return &abcitypes.ResponseCheckTx{Code: types.ParsingTransactionError}, nil
	}

	err = tx.IsValid()
	if err != nil {
		return &abcitypes.ResponseCheckTx{Code: types.InvalidTransactionError}, nil
	}

	switch tx.Type {
	case transactions.OChainPortalInteraction:
		transaction, err := transactions.ParseNewOChainPortalInteraction(tx)
		if err != nil {
			return &abcitypes.ResponseCheckTx{Code: types.ParsingTransactionDataError}, nil
		}

		err = transaction.Check(txCtx)
		if err != nil {
			log.Println(err)
			return &abcitypes.ResponseCheckTx{Code: types.CheckTransactionFailure}, nil
		}

		return &abcitypes.ResponseCheckTx{Code: types.NoError}, nil

	case transactions.RegisterAccount:

		err := tx.VerifySignature()
		if err != nil {
			return &abcitypes.ResponseCheckTx{Code: types.InvalidTransactionSignature}, nil
		}

		transaction, err := transactions.ParseRegisterAccountTransaction(tx)
		if err != nil {
			return &abcitypes.ResponseCheckTx{Code: types.ParsingTransactionDataError}, nil
		}

		err = transaction.Check(txCtx)
		if err != nil {
			return &abcitypes.ResponseCheckTx{Code: types.CheckTransactionFailure}, nil
		}

		return &abcitypes.ResponseCheckTx{Code: types.NoError}, nil
	}

	return &abcitypes.ResponseCheckTx{Code: types.NotImplemented, GasWanted: 0}, nil
}

func (app *OChainValidatorApplication) FinalizeBlock(_ context.Context, req *abcitypes.RequestFinalizeBlock) (*abcitypes.ResponseFinalizeBlock, error) {
	log.Printf("Finalize block: %d", req.Height)
	app.ValUpdates = make([]abcitypes.ValidatorUpdate, 0)
	app.ongoingBlockTx = app.store.Badger().NewTransaction(true)

	txCtx := transactions.TransactionContext{
		Config: app.config,
		Db:     app.db,
		Txn:    app.ongoingBlockTx,
		Date:   req.Time,
	}

	var txs = make([]*abcitypes.ExecTxResult, len(req.Txs))
	for i, tx := range req.Txs {
		app.state.IncSize()
		log.Printf("Finalize tx: %s", hex.EncodeToString(tx))

		tx, err := transactions.ParseTransaction(tx)
		if err != nil {
			txs[i] = &abcitypes.ExecTxResult{Code: types.ParsingTransactionError}
			continue
		}

		err = tx.IsValid()
		if err != nil {
			txs[i] = &abcitypes.ExecTxResult{Code: types.InvalidTransactionError}
			continue
		}

		switch tx.Type {
		case transactions.OChainPortalInteraction:

			transaction, err := transactions.ParseNewOChainPortalInteraction(tx)
			if err != nil {
				txs[i] = &abcitypes.ExecTxResult{Code: types.ParsingTransactionDataError}
				continue
			}

			err = transaction.Check(txCtx)
			if err != nil {
				txs[i] = &abcitypes.ExecTxResult{Code: types.CheckTransactionFailure}
				continue
			}

			err = transaction.Execute(txCtx)
			if err != nil {
				txs[i] = &abcitypes.ExecTxResult{Code: types.ExecuteTransactionFailure}
				continue
			}

			if transaction.Data.Type == transactions.NewValidatorPortalInteractionType {
				formatedTx, err := transactions.ParseNewValidatorTransaction(transaction)
				if err != nil {
					txs[i] = &abcitypes.ExecTxResult{Code: types.ParsingTransactionDataError}
					continue
				}

				pubkeyBytes, err := hex.DecodeString(formatedTx.Data.Arguments.PublicKey)
				if err != nil {
					txs[i] = &abcitypes.ExecTxResult{Code: types.ParsingTransactionDataError}
					continue
				}

				app.ValUpdates = append(app.ValUpdates, abcitypes.UpdateValidator(pubkeyBytes, 10000, "secp256k1"))
			} else if transaction.Data.Type == transactions.RemoveValidatorPortalInteractionType {
				formatedTx, err := transactions.ParseRemoveValidatorTransaction(transaction)
				if err != nil {
					txs[i] = &abcitypes.ExecTxResult{Code: types.ParsingTransactionDataError}
					continue
				}

				validator, err := app.db.Validators.Get(formatedTx.Data.Arguments.ValidatorId, app.ongoingBlockTx)
				if err != nil {
					txs[i] = &abcitypes.ExecTxResult{Code: types.ExecuteTransactionFailure}
					continue
				}

				pubkeyBytes, err := hex.DecodeString(validator.PublicKey)
				if err != nil {
					txs[i] = &abcitypes.ExecTxResult{Code: types.ParsingTransactionDataError}
					continue
				}

				app.ValUpdates = append(app.ValUpdates, abcitypes.UpdateValidator(pubkeyBytes, 0, "secp256k1"))
			}

			txs[i] = &abcitypes.ExecTxResult{Code: types.NoError}

		case transactions.RegisterAccount:

			err := tx.VerifySignature()
			if err != nil {
				txs[i] = &abcitypes.ExecTxResult{Code: types.InvalidTransactionSignature}
				continue
			}

			transaction, err := transactions.ParseRegisterAccountTransaction(tx)
			if err != nil {
				txs[i] = &abcitypes.ExecTxResult{Code: types.ParsingTransactionDataError}
				continue
			}

			err = transaction.Check(txCtx)
			if err != nil {
				txs[i] = &abcitypes.ExecTxResult{Code: types.CheckTransactionFailure}
				continue
			}

			events, err := transaction.Execute(txCtx)
			if err != nil {
				txs[i] = &abcitypes.ExecTxResult{Code: types.ExecuteTransactionFailure}
				continue
			}

			txs[i] = &abcitypes.ExecTxResult{
				Code:    types.NoError,
				Log:     "Account registered: " + tx.From,
				Events:  events,
				GasUsed: 0,
			}
		}

	}

	app.state.SetHeight(req.Height)
	log.Println(app.state)
	return &abcitypes.ResponseFinalizeBlock{
		TxResults:        txs,
		ValidatorUpdates: app.ValUpdates,
		AppHash:          app.state.Hash(),
	}, nil
}

func (app *OChainValidatorApplication) Commit(_ context.Context, commit *abcitypes.RequestCommit) (*abcitypes.ResponseCommit, error) {
	app.db.State.Save(app.state, app.ongoingBlockTx)
	err := app.ongoingBlockTx.Commit()
	return &abcitypes.ResponseCommit{}, err
}

func (app *OChainValidatorApplication) Query(_ context.Context, req *abcitypes.RequestQuery) (*abcitypes.ResponseQuery, error) {
	// res := abcitypes.ResponseQuery{}

	// queries.ResolveQuery(req.Data)
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

func (app *OChainValidatorApplication) PrepareProposal(ctx context.Context, proposal *abcitypes.RequestPrepareProposal) (*abcitypes.ResponsePrepareProposal, error) {
	log.Printf("Prepare proposal called: %d txs", len(proposal.Txs))
	return &abcitypes.ResponsePrepareProposal{Txs: app.formatTxs(ctx, proposal.Txs)}, nil
}

// formatTxs validates and excludes invalid transactions
// also substitutes all the transactions with x:y to x=y
func (app *OChainValidatorApplication) formatTxs(ctx context.Context, blockData [][]byte) [][]byte {
	txs := make([][]byte, 0, len(blockData))
	for _, tx := range blockData {
		if resp, err := app.CheckTx(ctx, &abcitypes.RequestCheckTx{Tx: tx}); err == nil && resp.Code == types.NoError {
			txs = append(txs, tx)
		}
	}
	return txs
}

func (app *OChainValidatorApplication) ProcessProposal(ctx context.Context, proposal *abcitypes.RequestProcessProposal) (*abcitypes.ResponseProcessProposal, error) {
	log.Printf("Process proposal called: %d txs", len(proposal.Txs))
	for _, tx := range proposal.Txs {
		// As CheckTx is a full validity check we can simply reuse this
		if resp, err := app.CheckTx(ctx, &abcitypes.RequestCheckTx{Tx: tx}); err != nil || resp.Code != types.NoError {
			return &abcitypes.ResponseProcessProposal{Status: abcitypes.ResponseProcessProposal_REJECT}, nil
		}
	}
	return &abcitypes.ResponseProcessProposal{Status: abcitypes.ResponseProcessProposal_ACCEPT}, nil
}
