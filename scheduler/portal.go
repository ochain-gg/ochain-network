package scheduler

import (
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ochain.gg/ochain-network-validator/config"
	"github.com/ochain.gg/ochain-network-validator/contracts"
	"github.com/ochain.gg/ochain-network-validator/database"
	"github.com/ochain.gg/ochain-network-validator/swagger"
)

func CheckAndHandlePortalUpdate(cfg config.OChainConfig, db *database.OChainDatabase) {

	//skip if application not synced
	swaggerConfig := swagger.NewConfiguration()
	nodeClient := swagger.NewAPIClient(swaggerConfig)

	status, _, err := nodeClient.InfoApi.Status(nil)

	if err != nil {
		return
	}

	if status.Result.SyncInfo.EarliestBlockHeight < status.Result.SyncInfo.LatestBlockHeight {
		return
	}

	//get smart contract last block update
	client, err := ethclient.Dial(cfg.EVMRpc)
	if err != nil {
		log.Fatal(err)
	}
	address := common.HexToAddress(cfg.EVMPortalAddress)
	portal, err := contracts.NewOChainPortal(address, client)
	if err != nil {
		return
	}

	res, err := portal.LatestUpdateAt(nil)
	lastContractUpdateBlockNumber := res.Uint64()
	if err != nil {
		return
	}

	//get last smart contract update on app
	state, err := db.State.Get()
	lastAppUpdateBlockNumber := state.LatestPortalUpdate

	//compare and skip if sc and app last update are the same
	if lastContractUpdateBlockNumber <= lastAppUpdateBlockNumber {
		return
	}

	//get new logs (iterating block by 2000)
	totalBlockToParse := lastContractUpdateBlockNumber - lastAppUpdateBlockNumber
	var logsToIndex []types.Log

	if totalBlockToParse < 1000 {
		filter := ethereum.FilterQuery{
			FromBlock: new(big.Int).SetUint64(lastAppUpdateBlockNumber + 1),
			ToBlock:   new(big.Int).SetUint64(lastContractUpdateBlockNumber),
			Addresses: []common.Address{address},
		}
		logsToIndex, err = client.FilterLogs(nil, filter)
		if err != nil {
			return
		}
	} else {
		for i := lastAppUpdateBlockNumber; i <= lastContractUpdateBlockNumber; i += 1000 {

			var toBlock uint64 = i
			if lastAppUpdateBlockNumber+1000 > lastContractUpdateBlockNumber {
				toBlock = lastContractUpdateBlockNumber
			}

			filter := ethereum.FilterQuery{
				FromBlock: new(big.Int).SetUint64(lastAppUpdateBlockNumber + 1),
				ToBlock:   new(big.Int).SetUint64(toBlock),
				Addresses: []common.Address{address},
			}

			logs, err := client.FilterLogs(nil, filter)
			if err != nil {
				return
			}

			for i := 0; i < len(logs); i++ {
				logsToIndex = append(logsToIndex, logs[i])
			}
		}
	}

	//send update transaction

}
