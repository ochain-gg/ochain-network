package tests

import (
	"log"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ochain-gg/ochain-network/application"
	ochainCfg "github.com/ochain-gg/ochain-network/config"
	"github.com/ochain-gg/ochain-network/engine/database"
)

func NewOChainTestingApplication() (*application.OChainApplication, error) {
	ochainConfig := ochainCfg.DefaultConfig()
	privateKey, err := crypto.HexToECDSA("ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80")
	if err != nil {
		log.Fatal(err)
	}

	db := database.NewOChainDatabase("/tmp/ochain_test.db")
	privateKeyBytes := crypto.FromECDSA(privateKey)
	return application.NewOChainApplication(*ochainConfig, db, privateKeyBytes)
}
