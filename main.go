package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"syscall"

	"github.com/ethereum/go-ethereum/crypto"

	cometSecp256k1 "github.com/cometbft/cometbft/crypto/secp256k1"
	"github.com/cometbft/cometbft/p2p"
	"github.com/cometbft/cometbft/privval"
	"github.com/cometbft/cometbft/proxy"
	"github.com/timshannon/badgerhold/v4"

	cfg "github.com/cometbft/cometbft/config"
	cmtflags "github.com/cometbft/cometbft/libs/cli/flags"
	cmtlog "github.com/cometbft/cometbft/libs/log"
	nm "github.com/cometbft/cometbft/node"
	ochainCfg "github.com/ochain.gg/ochain-network-validator/config"
	"github.com/spf13/viper"
)

var homeDir string
var evmChainId string
var evmRpc string
var evmPortalAddress string

func init() {
	flag.StringVar(&homeDir, "cmt-home", "", "Path to the CometBFT config directory (if empty, uses $HOME/.cometbft)")
	flag.StringVar(&evmChainId, "chainId", "31337", "OChain portal chainId")
	flag.StringVar(&evmRpc, "evmRpc", "http://localhost:8545/", "OChain portal chain rpc address")
	flag.StringVar(&evmPortalAddress, "portalAddress", "0x8A791620dd6260079BF849Dc5567aDC3F2FdC318", "OChain portal address")
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
	dbPath := filepath.Join(homeDir, "badger")
	options := badgerhold.DefaultOptions
	options.Dir = dbPath
	options.ValueDir = "ochain-data"
	db, err := badgerhold.Open(options)

	if err != nil {
		log.Fatalf("Opening database: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("Closing database: %v", err)
		}
	}()

	ochainConfig := ochainCfg.DefaultConfig()
	parsedChainId, _ := strconv.ParseUint(evmChainId, 10, 64)
	ochainConfig.EVMChainId = parsedChainId
	ochainConfig.EVMRpc = evmRpc
	ochainConfig.EVMPortalAddress = evmPortalAddress

	nodeKey, err := p2p.LoadNodeKey(config.NodeKeyFile())
	if err != nil {
		log.Fatalf("failed to load node's key: %v", err)
	}

	privateKey, err := crypto.HexToECDSA(string(nodeKey.PrivKey.Bytes()))
	if err != nil {
		log.Fatal(err)
	}

	privateKeyBytes := crypto.FromECDSA(privateKey)
	nodeKey.PrivKey = cometSecp256k1.PrivKey(privateKeyBytes)
	formatedPv := privval.NewFilePV(nodeKey.PrivKey, config.PrivValidatorKeyFile(), config.PrivValidatorStateFile())

	app := NewOChainValidatorApplication(*ochainConfig, db, privateKeyBytes)

	logger := cmtlog.NewTMLogger(cmtlog.NewSyncWriter(os.Stdout))
	logger, err = cmtflags.ParseLogLevel(config.LogLevel, logger, cfg.DefaultLogLevel)

	if err != nil {
		log.Fatalf("failed to parse log level: %v", err)
	}

	node, err := nm.NewNode(
		config,
		formatedPv,
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
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
}
