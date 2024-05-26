package main

import (
	"context"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
	"strconv"
	"time"

	abcitypes "github.com/cometbft/cometbft/abci/types"
	"github.com/cometbft/cometbft/crypto/ed25519"
	"github.com/cometbft/cometbft/proto/tendermint/crypto"
	"github.com/cometbft/cometbft/version"
	"github.com/ethereum/go-ethereum/common"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ochain-gg/ochain-network/config"
	"github.com/ochain-gg/ochain-network/contracts"
	"github.com/ochain-gg/ochain-network/database"
	"github.com/ochain-gg/ochain-network/queries"
	"github.com/ochain-gg/ochain-network/transactions"
	"github.com/ochain-gg/ochain-network/types"
)

const (
	AppVersion uint64 = 1
)

type OChainValidatorApplication struct {
	abcitypes.BaseApplication

	config config.OChainConfig
	db     *database.OChainDatabase

	remotePrivateKey []byte

	state        *types.OChainState
	RetainBlocks int64
	ValUpdates   []abcitypes.ValidatorUpdate
}

var _ abcitypes.Application = (*OChainValidatorApplication)(nil)

func NewOChainValidatorApplication(config config.OChainConfig, dbpath string, remotePrivateKey []byte) (*OChainValidatorApplication, error) {
	db := database.NewOChainDatabase(dbpath)

	state, err := db.State.Get()
	if err != nil {
		return &OChainValidatorApplication{}, err
	}

	return &OChainValidatorApplication{
		config:           config,
		state:            &state,
		db:               db,
		remotePrivateKey: remotePrivateKey,
	}, nil
}

func (app *OChainValidatorApplication) Hash() []byte {
	hash := ethcrypto.Keccak256Hash([]byte(strconv.FormatInt(app.state.Size, 16))).Bytes()
	log.Printf("State hash processed at size %d: %s", app.state.Size, hex.EncodeToString(hash))
	return hash
}

func (app *OChainValidatorApplication) Info(_ context.Context, info *abcitypes.RequestInfo) (*abcitypes.ResponseInfo, error) {
	log.Printf("Info call last block heigh: %d", app.state.Height)
	log.Printf("Info call last block hash: %s", app.state.Hash)
	return &abcitypes.ResponseInfo{
		Data:             fmt.Sprintf("{\"size\":%v}", app.state.Size),
		Version:          version.ABCIVersion,
		AppVersion:       AppVersion,
		LastBlockHeight:  app.state.Height,
		LastBlockAppHash: app.Hash(),
	}, nil

}

func (app *OChainValidatorApplication) InitChain(_ context.Context, chain *abcitypes.RequestInitChain) (*abcitypes.ResponseInitChain, error) {

	app.db.NewTransaction(uint64(chain.GetTime().Unix()))

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

	mainUniverse := DefaultUniverse
	mainUniverse.CreatedAt = uint64(time.Now().Unix())

	_, err = app.db.Universes.Get("main")
	if err != nil {
		app.db.Universes.Insert(mainUniverse)
	}

	err = app.db.CommitTransaction()
	if err != nil {
		log.Fatal(err)
	}

	return &abcitypes.ResponseInitChain{
		Validators: validators,
		AppHash:    app.Hash(),
	}, nil
}

func (app *OChainValidatorApplication) CheckTx(ctx context.Context, req *abcitypes.RequestCheckTx) (*abcitypes.ResponseCheckTx, error) {

	txCtx := transactions.TransactionContext{
		Config: app.config,
		Db:     app.db,
		State:  app.state,
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
	case transactions.RegisterUniverseAccount:

		err := tx.VerifySignature()
		if err != nil {
			return &abcitypes.ResponseCheckTx{Code: types.InvalidTransactionSignature}, nil
		}

		transaction, err := transactions.ParseRegisterUniverseAccountTransaction(tx)
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

	app.db.NewTransaction(uint64(req.Time.Unix()))

	txCtx := transactions.TransactionContext{
		Config: app.config,
		Db:     app.db,
		State:  app.state,
		Date:   req.Time,
	}

	var txs = make([]*abcitypes.ExecTxResult, len(req.Txs))
	for i, tx := range req.Txs {

		tx, err := transactions.ParseTransaction(tx)
		if err != nil {
			txs[i] = &abcitypes.ExecTxResult{Code: types.ParsingTransactionError}
			log.Printf("Finalize tx %d: ParsingTransactionError", i)
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
				log.Printf("Finalize tx %d: ParsingTransactionDataError", i)
				continue
			}

			err = transaction.Check(txCtx)
			if err != nil {
				txs[i] = &abcitypes.ExecTxResult{Code: types.CheckTransactionFailure}
				log.Printf("Finalize tx %d: CheckTransactionFailure", i)
				continue
			}

			err = transaction.Execute(txCtx)
			if err != nil {
				txs[i] = &abcitypes.ExecTxResult{Code: types.ExecuteTransactionFailure}
				log.Printf("Finalize tx %d: ExecuteTransactionFailure", i)
				continue
			}

			if transaction.Data.Type == transactions.NewValidatorPortalInteractionType {
				formatedTx, err := transactions.ParseNewValidatorTransaction(transaction)
				if err != nil {
					txs[i] = &abcitypes.ExecTxResult{Code: types.ParsingTransactionDataError}
					log.Printf("Finalize tx %d: ParsingTransactionDataError", i)
					continue
				}

				pubkeyBytes, err := hex.DecodeString(formatedTx.Data.Arguments.PublicKey)
				if err != nil {
					txs[i] = &abcitypes.ExecTxResult{Code: types.ParsingTransactionDataError}
					log.Printf("Finalize tx %d: ParsingTransactionDataError", i)
					continue
				}

				app.ValUpdates = append(app.ValUpdates, abcitypes.UpdateValidator(pubkeyBytes, 10000, "ed25519"))
			} else if transaction.Data.Type == transactions.RemoveValidatorPortalInteractionType {
				formatedTx, err := transactions.ParseRemoveValidatorTransaction(transaction)
				if err != nil {
					txs[i] = &abcitypes.ExecTxResult{Code: types.ParsingTransactionDataError}
					log.Printf("Finalize tx %d: ParsingTransactionDataError", i)
					continue
				}

				validator, err := app.db.Validators.GetById(formatedTx.Data.Arguments.ValidatorId)
				if err != nil {
					txs[i] = &abcitypes.ExecTxResult{Code: types.ExecuteTransactionFailure}
					log.Printf("Finalize tx %d: ExecuteTransactionFailure", i)
					continue
				}

				pubkeyBytes, err := hex.DecodeString(validator.PublicKey)
				if err != nil {
					txs[i] = &abcitypes.ExecTxResult{Code: types.ParsingTransactionDataError}
					log.Printf("Finalize tx %d: ParsingTransactionDataError", i)
					continue
				}

				app.ValUpdates = append(app.ValUpdates, abcitypes.UpdateValidator(pubkeyBytes, 0, "ed25519"))
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

		case transactions.RegisterUniverseAccount:

			err := tx.VerifySignature()
			if err != nil {
				txs[i] = &abcitypes.ExecTxResult{Code: types.InvalidTransactionSignature}
				continue
			}

			transaction, err := transactions.ParseRegisterUniverseAccountTransaction(tx)
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

		app.state.IncSize()
	}

	app.state.SetHeight(req.Height)
	log.Println(app.state)
	return &abcitypes.ResponseFinalizeBlock{
		TxResults:        txs,
		ValidatorUpdates: app.ValUpdates,
		AppHash:          app.Hash(),
	}, nil
}

func (app *OChainValidatorApplication) Commit(_ context.Context, commit *abcitypes.RequestCommit) (*abcitypes.ResponseCommit, error) {
	app.db.State.Upsert(*app.state)
	err := app.db.CommitTransaction()
	return &abcitypes.ResponseCommit{}, err
}

func (app *OChainValidatorApplication) Query(_ context.Context, req *abcitypes.RequestQuery) (*abcitypes.ResponseQuery, error) {
	res := &abcitypes.ResponseQuery{}

	value, err := queries.GetQueryResponse(req, app.db)
	if err != nil {
		return &abcitypes.ResponseQuery{Code: 1}, err
	}

	res.Value = value
	res.Code = 0

	return res, nil
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
