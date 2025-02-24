package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"syscall"

	"github.com/ethereum/go-ethereum/crypto"

	"github.com/cometbft/cometbft/p2p"
	"github.com/cometbft/cometbft/privval"
	"github.com/cometbft/cometbft/proxy"

	cfg "github.com/cometbft/cometbft/config"
	cmtflags "github.com/cometbft/cometbft/libs/cli/flags"
	cmtlog "github.com/cometbft/cometbft/libs/log"
	nm "github.com/cometbft/cometbft/node"
	"github.com/ochain-gg/ochain-network/application"
	ochainCfg "github.com/ochain-gg/ochain-network/config"
	"github.com/ochain-gg/ochain-network/engine/database"
	"github.com/ochain-gg/ochain-network/scheduler"
	"github.com/spf13/viper"
)

var homeDir string
var evmChainId string
var evmRpc string
var evmPortalAddress string

func init() {
	flag.StringVar(&homeDir, "cmt-home", "", "Path to the CometBFT config directory (if empty, uses $HOME/.cometbft)")
	flag.StringVar(&evmChainId, "chainId", "11155111", "OChain portal chainId")
	flag.StringVar(&evmRpc, "evmRpc", "https://ethereum-sepolia.core.chainstack.com/ddf6b01951847ded1aac7e14b82c5b0c", "OChain portal chain rpc address")
	flag.StringVar(&evmPortalAddress, "portalAddress", "0x4Dd9d772C67fbC858918f364E5CB9e0B6E53Fd44", "OChain portal address")
}

func main() {
	flag.Parse()
	if homeDir == "" {
		homeDir = os.ExpandEnv("$HOME/.cometbft")
	}

	config := cfg.DefaultConfig()
	config.SetRoot(homeDir)
	viper.SetConfigFile(fmt.Sprintf("%s/%s", homeDir, "config/config.toml"))

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Reading config: %v", err)
	}
	if err := viper.Unmarshal(config); err != nil {
		log.Fatalf("Decoding config: %v", err)
	}
	if err := config.ValidateBasic(); err != nil {
		log.Fatalf("Invalid configuration data: %v", err)
	}

	// config.Consensus.CreateEmptyBlocks = false
	// config.Consensus.CreateEmptyBlocksInterval = time.Hour

	dbPath := filepath.Join(homeDir, "database")

	ochainConfig := ochainCfg.DefaultConfig()
	parsedChainId, _ := strconv.ParseUint(evmChainId, 10, 64)
	ochainConfig.EVMChainId = parsedChainId
	ochainConfig.EVMRpc = evmRpc
	ochainConfig.EVMPortalAddress = evmPortalAddress

	nodeKey, err := p2p.LoadNodeKey(config.NodeKeyFile())
	if err != nil {
		log.Fatalf("failed to load node's key: %v", err)
	}

	privateKey, err := crypto.HexToECDSA("ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80")
	if err != nil {
		log.Fatal(err)
	}

	privateKeyBytes := crypto.FromECDSA(privateKey)
	// privateKeyBytes := crypto.FromECDSA(privateKey)
	//nodeKey.PrivKey = cometSecp256k1.PrivKey(privateKeyBytes)
	//formatedPv := privval.NewFilePV(nodeKey.PrivKey, config.PrivValidatorKeyFile(), config.PrivValidatorStateFile())

	pv := privval.LoadFilePV(
		config.PrivValidatorKeyFile(),
		config.PrivValidatorStateFile(),
	)

	logger := cmtlog.NewTMLogger(cmtlog.NewSyncWriter(os.Stdout))
	logger, err = cmtflags.ParseLogLevel(config.LogLevel, logger, cfg.DefaultLogLevel)

	if err != nil {
		log.Fatalf("failed to parse log level: %v", err)
	}

	db := database.NewOChainDatabase(dbPath)
	app, err := application.NewOChainApplication(*ochainConfig, db, privateKeyBytes)
	if err != nil {
		log.Fatalln(err)
	}

	scheduler, err := scheduler.NewScheduler(*ochainConfig, db)
	if err != nil {
		log.Fatalf("Creating scheduler: %v", err)
	}

	scheduler.Scheduler.Start()

	node, err := nm.NewNode(
		context.Background(),
		config,
		pv,
		nodeKey,
		proxy.NewLocalClientCreator(app),
		nm.DefaultGenesisDocProviderFunc(config),
		cfg.DefaultDBProvider,
		nm.DefaultMetricsProvider(config.Instrumentation),
		logger,
	)

	if err != nil {
		log.Fatalf("Creating node: %v", err)
	}

	node.Start()

	defer func() {
		node.Stop()
		node.Wait()
		scheduler.Scheduler.Shutdown()
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
}
