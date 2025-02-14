package application

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
	"github.com/cometbft/cometbft/version"
	"github.com/ethereum/go-ethereum/common"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ochain-gg/ochain-network/config"
	"github.com/ochain-gg/ochain-network/contracts"
	"github.com/ochain-gg/ochain-network/database"
	"github.com/ochain-gg/ochain-network/engine"
	"github.com/ochain-gg/ochain-network/queries"
	"github.com/ochain-gg/ochain-network/transactions"
	"github.com/ochain-gg/ochain-network/types"
)

const (
	AppVersion uint64 = 1
)

type OChainApplication struct {
	abcitypes.BaseApplication

	config config.OChainConfig
	db     *database.OChainDatabase

	remotePrivateKey []byte

	state        *types.OChainState
	RetainBlocks int64
	ValUpdates   []abcitypes.ValidatorUpdate
}

var _ abcitypes.Application = (*OChainApplication)(nil)

func NewOChainApplication(config config.OChainConfig, db *database.OChainDatabase, remotePrivateKey []byte) (*OChainApplication, error) {
	state, err := db.State.Get()
	if err != nil {
		return &OChainApplication{}, err
	}

	return &OChainApplication{
		config:           config,
		state:            &state,
		db:               db,
		remotePrivateKey: remotePrivateKey,
	}, nil
}

func (app *OChainApplication) Hash() []byte {
	hash := ethcrypto.Keccak256Hash([]byte(strconv.FormatInt(app.state.Size, 16))).Bytes()
	log.Printf("State hash processed at size %d: %s", app.state.Size, hex.EncodeToString(hash))
	return hash
}

func (app *OChainApplication) Info(_ context.Context, info *abcitypes.InfoRequest) (*abcitypes.InfoResponse, error) {
	log.Printf("Info call last block heigh: %d", app.state.Height)
	log.Printf("Info call last block hash: %s", app.state.Hash)
	return &abcitypes.InfoResponse{
		Data:             fmt.Sprintf("{\"size\":%v}", app.state.Size),
		Version:          version.ABCIVersion,
		AppVersion:       AppVersion,
		LastBlockHeight:  app.state.Height,
		LastBlockAppHash: app.Hash(),
	}, nil
}

func (app *OChainApplication) InitChain(_ context.Context, chain *abcitypes.InitChainRequest) (*abcitypes.InitChainResponse, error) {

	app.db.NewTransaction(uint64(chain.Time.Unix()))

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
				PubKeyType:  "tendermint/PubKeyEd25519",
				PubKeyBytes: pubkeyBytes,
				Power:       10000,
			})

			log.Println("adding validator address: " + pubkey.Address().String())
		}
	}

	//generate default game entities
	for i := 0; i < len(DefaultBuildings); i++ {
		app.db.Buildings.Insert(DefaultBuildings[i])
	}
	for i := 0; i < len(DefaultTechnologies); i++ {
		app.db.Technologies.Insert(DefaultTechnologies[i])
	}
	for i := 0; i < len(DefaultDefenses); i++ {
		app.db.Defenses.Insert(DefaultDefenses[i])
	}
	for i := 0; i < len(DefaultSpaceships); i++ {
		app.db.Spaceships.Insert(DefaultSpaceships[i])
	}

	mainUniverse := DefaultUniverse
	mainUniverse.CreatedAt = uint64(time.Now().Unix())

	//Initalize marketplace

	//Reserve Token Balance / (Continuous Token Supply x Continuous Token Price)
	// 1 pour 10000= 0.0001
	//reserve ration = 1 000 000 / 1_000_000_000_000 * 0.0001
	//reserve ration = 1 000 000 / 100_000_000
	//reserve ration = 1 / 100
	//reserve ration = 0.01 = 10

	universeResourceMarket := types.OChainResourcesMarket{
		UniverseId: mainUniverse.Id,
		FeesRate:   200,

		MetalReserveRatio: float64(0.9),
		MetalPoolBalance:  uint64(100_000),        // 100k token for this reserve
		MetalSupplyMinted: uint64(10_000_000_000), // 10B minted at init

		CrystalReserveRatio: float64(0.9),
		CrystalPoolBalance:  uint64(100_000_000_000_000), // 100k token for this reserve
		CrystalSupplyMinted: uint64(5_000_000_000),       // 5B minted at init

		DeuteriumReserveRatio: float64(0.9),
		DeuteriumPoolBalance:  uint64(100_000_000_000_000), // 100k token for this reserve
		DeuteriumSupplyMinted: uint64(2_000_000_000),       // 2B minted at init
	}

	mainUniverse.ResourcesMarketEnabled = true
	mainUniverse.Speed = 50
	mainUniverse.IsStretchable = true
	mainUniverse.OCTCirculatingSupply = universeResourceMarket.MetalPoolBalance + universeResourceMarket.CrystalPoolBalance + universeResourceMarket.DeuteriumPoolBalance
	app.db.Universes.Insert(mainUniverse)

	app.db.ResourcesMarket.Insert(universeResourceMarket)

	err = app.db.CommitTransaction()
	if err != nil {
		log.Fatal(err)
	}

	return &abcitypes.InitChainResponse{
		Validators: validators,
		AppHash:    app.Hash(),
	}, nil
}

func (app *OChainApplication) CheckTx(ctx context.Context, req *abcitypes.CheckTxRequest) (*abcitypes.CheckTxResponse, error) {
	txCtx := transactions.TransactionContext{
		Config: app.config,
		Db:     app.db,
		State:  app.state,
		Date:   time.Now(),
	}

	return engine.CheckTx(txCtx, req), nil
}
func (app *OChainApplication) FinalizeBlock(_ context.Context, req *abcitypes.FinalizeBlockRequest) (*abcitypes.FinalizeBlockResponse, error) {
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

		execTxResult, valUpdates := engine.FinalizeTx(txCtx, tx)

		txs[i] = execTxResult
		for i := 0; i < len(valUpdates); i++ {
			app.ValUpdates = append(app.ValUpdates, valUpdates[i])
		}

		app.state.IncSize()
	}

	//handle slashing system
	if len(req.Misbehavior) > 0 {
		for i := 0; i < len(req.Misbehavior); i++ {
			misbehavior := req.Misbehavior[i]

			validator, err := app.db.Validators.GetByAddressAt(string(misbehavior.Validator.Address), uint64(req.Time.Unix()))
			if err != nil {
				log.Fatal(err)
			}

			validator.Power -= 10
			if validator.Power < 5000 {
				//handle full slashing
			} else {

			}

			app.db.Validators.Update(validator)
		}
	}

	app.state.SetHeight(req.Height)
	log.Println(app.state)
	return &abcitypes.FinalizeBlockResponse{
		TxResults:        txs,
		ValidatorUpdates: app.ValUpdates,

		AppHash: app.Hash(),
	}, nil
}

func (app *OChainApplication) Commit(_ context.Context, commit *abcitypes.CommitRequest) (*abcitypes.CommitResponse, error) {
	app.db.State.Upsert(*app.state)
	err := app.db.CommitTransaction()
	return &abcitypes.CommitResponse{}, err
}

func (app *OChainApplication) Query(_ context.Context, req *abcitypes.QueryRequest) (*abcitypes.QueryResponse, error) {
	value, err := queries.GetQueryResponse(req, app.db)
	if err != nil {
		return &abcitypes.QueryResponse{Code: 1}, err
	}

	return &abcitypes.QueryResponse{
		Code:  0,
		Value: value,
	}, err
}

func (app *OChainApplication) PrepareProposal(ctx context.Context, proposal *abcitypes.PrepareProposalRequest) (*abcitypes.PrepareProposalResponse, error) {
	log.Printf("Prepare proposal called: %d txs", len(proposal.Txs))
	return &abcitypes.PrepareProposalResponse{
		Txs: app.formatTxs(ctx, proposal.Txs),
	}, nil
}

// formatTxs validates and excludes invalid transactions
// also substitutes all the transactions with x:y to x=y
func (app *OChainApplication) formatTxs(ctx context.Context, blockData [][]byte) [][]byte {
	txs := make([][]byte, 0, len(blockData))
	for _, tx := range blockData {
		if resp, err := app.CheckTx(ctx, &abcitypes.CheckTxRequest{Tx: tx}); err == nil && resp.Code == types.NoError {
			txs = append(txs, tx)
		}
	}
	return txs
}

func (app *OChainApplication) ProcessProposal(ctx context.Context, proposal *abcitypes.ProcessProposalRequest) (*abcitypes.ProcessProposalResponse, error) {
	log.Printf("Process proposal called: %d txs", len(proposal.Txs))
	for _, tx := range proposal.Txs {
		// As CheckTx is a full validity check we can simply reuse this
		if resp, err := app.CheckTx(ctx, &abcitypes.CheckTxRequest{Tx: tx}); err != nil || resp.Code != types.NoError {
			return &abcitypes.ProcessProposalResponse{
				Status: abcitypes.PROCESS_PROPOSAL_STATUS_REJECT,
			}, nil
		}
	}
	return &abcitypes.ProcessProposalResponse{
		Status: abcitypes.PROCESS_PROPOSAL_STATUS_ACCEPT,
	}, nil
}
