package scheduler

import (
	"context"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ochain.gg/ochain-network-validator/config"
	"github.com/ochain.gg/ochain-network-validator/contracts"
	"github.com/ochain.gg/ochain-network-validator/database"
	"github.com/ochain.gg/ochain-network-validator/transactions"

	rpchttp "github.com/cometbft/cometbft/rpc/client/http"
)

func CheckAndHandlePortalUpdate(cfg config.OChainConfig, db *database.OChainDatabase) {

	log.Println("CheckAndHandlePortalUpdate process start")

	//skip if application not synced
	client, err := rpchttp.New("http://127.0.0.1:26657", "/websocket")
	if err != nil {
		log.Println(err)
	}

	// status, err := client.Status(context.Background())
	// if err != nil {
	// 	log.Println(err)
	// 	return
	// }

	// log.Printf("Current block: %d, Latest block: %d", status.SyncInfo.EarliestBlockHeight, status.SyncInfo.LatestBlockHeight)
	// if status.SyncInfo.EarliestBlockHeight+1 < status.SyncInfo.LatestBlockHeight {
	// 	log.Println("node not synced")
	// 	return
	// }

	//get smart contract last block update
	evmClient, err := ethclient.Dial(cfg.EVMRpc)
	if err != nil {
		log.Println(err)
		return
	}

	log.Printf("Get last portal update on: %s", cfg.EVMPortalAddress)
	address := common.HexToAddress(cfg.EVMPortalAddress)
	portal, err := contracts.NewOChainPortal(address, evmClient)
	if err != nil {
		log.Println(err)
		return
	}

	res, err := portal.LatestUpdateAt(nil)
	lastContractUpdateBlockNumber := res.Uint64()
	if err != nil {
		log.Println(err)
		return
	}

	log.Printf("Last portal update on chain: %d", res.Uint64())

	//get last smart contract update on app
	state, err := db.State.Get()

	if err != nil {
		log.Println("error on state load")
		log.Println(err)
		return
	}

	log.Printf("Last portal update on state: %d", res.Uint64())
	lastAppUpdateBlockNumber := state.LatestPortalUpdate

	//compare and skip if sc and app last update are the same
	if lastContractUpdateBlockNumber <= lastAppUpdateBlockNumber {
		log.Println("Portal sync: up to date")
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
		logsToIndex, err = evmClient.FilterLogs(context.Background(), filter)
		if err != nil {
			log.Println("portal sync up to date")
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

			logs, err := evmClient.FilterLogs(context.Background(), filter)
			if err != nil {
				log.Println("portal sync up to date")
				return
			}

			for i := 0; i < len(logs); i++ {
				logsToIndex = append(logsToIndex, logs[i])
			}
		}
	}

	//send update transaction
	for i := 0; i < len(logsToIndex); i++ {
		l := logsToIndex[i]

		newValidatorLog, err := portal.ParseOChainNewValidator(l)
		if err == nil {
			log.Print("new validator evm logs detected")

			tx := transactions.NewValidatorTransaction{
				Type: transactions.OChainPortalInteraction,
				Data: transactions.NewValidatorTransactionData{
					Type: transactions.NewValidatorPortalInteractionType,
					Arguments: transactions.NewValidatorTransactionDataArguments{
						ValidatorId:           newValidatorLog.ValidatorId.Uint64(),
						RemoteTransactionHash: l.TxHash.Hex(),
						PublicKey:             newValidatorLog.PublicKey,
					},
				},
			}

			baseTx, err := tx.Transaction()
			if err != nil {
				log.Println(err)
				return
			}

			bytesTx, err := baseTx.Bytes()
			if err != nil {
				log.Println(err)
				return
			}

			client, err := rpchttp.New("http://127.0.0.1:26657", "/websocket")
			if err != nil {
				log.Println(err)
				return
			}

			resCheck, err := client.CheckTx(context.Background(), bytesTx)
			if err != nil {
				log.Println(err)
				return
			}

			log.Printf("CheckTx code: %d", resCheck.Code)

			res, err := client.BroadcastTxCommit(context.Background(), bytesTx)
			if err != nil {
				log.Println(err)
				return
			}

			log.Printf("CheckTx code: %d", res.CheckTx.Code)
			log.Printf("Transaction Hash: %b", res.Hash)
			log.Printf("Block height: %d", res.Height)
		}

		removeValidatorLog, err := portal.ParseOChainRemoveValidator(l)
		if err == nil {
			tx := transactions.RemoveValidatorTransaction{
				Type: transactions.OChainPortalInteraction,
				Data: transactions.RemoveValidatorTransactionData{
					Type: transactions.RemoveValidatorPortalInteractionType,
					Arguments: transactions.RemoveValidatorTransactionDataArguments{
						ValidatorId:           removeValidatorLog.ValidatorId.Uint64(),
						RemoteTransactionHash: l.TxHash.Hex(),
					},
				},
			}

			baseTx, err := tx.Transaction()
			if err != nil {
				log.Println(err)
				return
			}

			bytesTx, err := baseTx.Bytes()
			if err != nil {
				log.Println(err)
				return
			}

			res, err := client.BroadcastTxCommit(context.Background(), bytesTx)
			if err != nil {
				log.Println(err)
			}

			log.Printf("CheckTx code: %d", res.CheckTx.Code)
			log.Printf("Transaction Hash: %b", res.Hash)
			log.Printf("Block height: %d", res.Height)

		}

		// tokenDepositLog, err := portal.ParseOChainTokenDeposit(log)
		// if err == nil {

		// }

		// tokenWithdrawalRequestedLog, err := portal.ParseOChainTokenWithdrawalRequested(log)
		// if err == nil {

		// }

		// tokenWithdrawalRequestContestedLog, err := portal.ParseOChainTokenWithdrawalRequestContested(log)
		// if err == nil {

		// }

		// unstackSucceedLog, err := portal.ParseOChainUnstackSucceed(log)
		// if err == nil {

		// }
	}

}
