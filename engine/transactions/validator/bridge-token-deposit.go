package validator_transactions

import (
	"context"
	"log"

	abcitypes "github.com/cometbft/cometbft/abci/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/fxamacker/cbor/v2"
	"github.com/ochain-gg/ochain-network/contracts"
	t "github.com/ochain-gg/ochain-network/engine/transactions"
	"github.com/ochain-gg/ochain-network/types"
)

type OChainTokenDepositTransactionData struct {
	RemoteTransactionHash string `cbor:"1,keyasint"`
	Account               string `cbor:"2,keyasint"`
	Amount                uint64 `cbor:"3,keyasint"`
}

type OChainTokenDepositTransaction struct {
	Type t.TransactionType
	Data OChainTokenDepositTransactionData
}

func (tx OChainTokenDepositTransaction) Check(ctx t.TransactionContext) *abcitypes.CheckTxResponse {

	client, err := ethclient.Dial(ctx.Config.EVMRpc)
	if err != nil {
		log.Fatal(err)
	}

	remoteTx, _, err := client.TransactionByHash(context.Background(), common.HexToHash(tx.Data.RemoteTransactionHash))
	if err != nil {
		return &abcitypes.CheckTxResponse{
			Code: types.InvalidTransactionError,
		}
	}

	if remoteTx.ChainId().Uint64() != ctx.Config.EVMChainId {
		return &abcitypes.CheckTxResponse{
			Code: types.InvalidTransactionError,
		}
	}

	if *remoteTx.To() != common.HexToAddress(ctx.Config.EVMPortalAddress) {
		return &abcitypes.CheckTxResponse{
			Code: types.InvalidTransactionError,
		}
	}

	remoteTxReceipt, err := client.TransactionReceipt(context.Background(), common.HexToHash(tx.Data.RemoteTransactionHash))
	if err != nil {
		return &abcitypes.CheckTxResponse{
			Code: types.InvalidTransactionError,
		}
	}

	if remoteTxReceipt.Status != 1 {
		return &abcitypes.CheckTxResponse{
			Code: types.InvalidTransactionError,
		}
	}

	address := common.HexToAddress(ctx.Config.EVMPortalAddress)
	portal, err := contracts.NewOChainPortal(address, client)
	if err != nil {
		return &abcitypes.CheckTxResponse{
			Code: types.InvalidTransactionError,
		}
	}

	for _, vLog := range remoteTxReceipt.Logs {

		_, err := portal.ParseOChainTokenDeposited(*vLog)
		if err != nil {
			continue
		}

		return &abcitypes.CheckTxResponse{
			Code: types.NoError,
		}
	}

	return &abcitypes.CheckTxResponse{
		Code: types.InvalidTransactionError,
	}

}

func (tx OChainTokenDepositTransaction) Execute(ctx t.TransactionContext) *abcitypes.ExecTxResult {
	result := tx.Check(ctx)
	if result.Code != types.NoError {
		return &abcitypes.ExecTxResult{
			Code: result.GetCode(),
		}
	}

	client, err := ethclient.Dial(ctx.Config.EVMRpc)
	if err != nil {
		log.Fatal(err)
	}

	remoteTx, _, err := client.TransactionByHash(context.Background(), common.HexToHash(tx.Data.RemoteTransactionHash))
	if err != nil {
		return &abcitypes.ExecTxResult{
			Code: types.InvalidTransactionError,
		}
	}

	if remoteTx.ChainId().Uint64() != ctx.Config.EVMChainId {
		return &abcitypes.ExecTxResult{
			Code: types.InvalidTransactionError,
		}
	}

	if *remoteTx.To() != common.HexToAddress(ctx.Config.EVMPortalAddress) {
		return &abcitypes.ExecTxResult{
			Code: types.InvalidTransactionError,
		}
	}

	remoteTxReceipt, err := client.TransactionReceipt(context.Background(), common.HexToHash(tx.Data.RemoteTransactionHash))
	if err != nil {
		return &abcitypes.ExecTxResult{
			Code: types.InvalidTransactionError,
		}
	}

	if remoteTxReceipt.Status != 1 {
		return &abcitypes.ExecTxResult{
			Code: types.InvalidTransactionError,
		}
	}

	address := common.HexToAddress(ctx.Config.EVMPortalAddress)
	portal, err := contracts.NewOChainPortal(address, client)
	if err != nil {
		return &abcitypes.ExecTxResult{
			Code: types.InvalidTransactionError,
		}
	}

	for _, vLog := range remoteTxReceipt.Logs {

		log, err := portal.ParseOChainTokenDeposited(*vLog)
		if err != nil {
			continue
		}

		_, err = ctx.Db.BridgeTransactions.GetAt(tx.Data.RemoteTransactionHash, uint64(ctx.Date.Unix()))
		if err == nil {
			return &abcitypes.ExecTxResult{
				Code: types.InvalidTransactionError,
			}
		}

		creditDepositTx := types.OChainTransaction{
			Type:     types.OChainTokenDepositTransaction,
			Hash:     tx.Data.RemoteTransactionHash,
			Account:  log.Receiver.Hex(),
			Amount:   log.Amount.Uint64(),
			Executed: false,
			Canceled: false,
		}

		err = ctx.Db.BridgeTransactions.Insert(creditDepositTx)
		if err != nil {
			return &abcitypes.ExecTxResult{
				Code: types.InvalidTransactionError,
			}
		}

		return &abcitypes.ExecTxResult{
			Code: types.NoError,
		}
	}

	return &abcitypes.ExecTxResult{
		Code: types.InvalidTransactionError,
	}
}

func (tx OChainTokenDepositTransaction) Transaction() (t.Transaction, error) {

	txData, err := cbor.Marshal(tx.Data)
	if err != nil {
		return t.Transaction{}, err
	}

	return t.Transaction{
		Type: tx.Type,
		Data: txData,
	}, nil
}

func ParseOChainTokenDepositTransaction(tx t.Transaction) (OChainTokenDepositTransaction, error) {
	var txData OChainTokenDepositTransactionData
	err := cbor.Unmarshal(tx.Data, &txData)

	if err != nil {
		return OChainTokenDepositTransaction{}, err
	}

	return OChainTokenDepositTransaction{
		Type: tx.Type,
		Data: txData,
	}, nil
}
