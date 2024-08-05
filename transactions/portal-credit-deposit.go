package transactions

import (
	"context"
	"log"

	abcitypes "github.com/cometbft/cometbft/abci/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/fxamacker/cbor/v2"
	"github.com/ochain-gg/ochain-network/contracts"
	"github.com/ochain-gg/ochain-network/types"
)

type OChainreditDepositTransactionData struct {
	RemoteTransactionHash string `cbor:"1,keyasint"`
	Account               string `cbor:"2,keyasint"`
	Amount                uint64 `cbor:"3,keyasint"`
}

type OChainreditDepositTransaction struct {
	Type TransactionType
	Data OChainreditDepositTransactionData
}

func (tx OChainreditDepositTransaction) Check(ctx TransactionContext) *abcitypes.ResponseCheckTx {

	client, err := ethclient.Dial(ctx.Config.EVMRpc)
	if err != nil {
		log.Fatal(err)
	}

	remoteTx, _, err := client.TransactionByHash(context.Background(), common.HexToHash(tx.Data.RemoteTransactionHash))
	if err != nil {
		return &abcitypes.ResponseCheckTx{
			Code: types.InvalidTransactionError,
		}
	}

	if remoteTx.ChainId().Uint64() != ctx.Config.EVMChainId {
		return &abcitypes.ResponseCheckTx{
			Code: types.InvalidTransactionError,
		}
	}

	if *remoteTx.To() != common.HexToAddress(ctx.Config.EVMPortalAddress) {
		return &abcitypes.ResponseCheckTx{
			Code: types.InvalidTransactionError,
		}
	}

	remoteTxReceipt, err := client.TransactionReceipt(context.Background(), common.HexToHash(tx.Data.RemoteTransactionHash))
	if err != nil {
		return &abcitypes.ResponseCheckTx{
			Code: types.InvalidTransactionError,
		}
	}

	if remoteTxReceipt.Status != 1 {
		return &abcitypes.ResponseCheckTx{
			Code: types.InvalidTransactionError,
		}
	}

	address := common.HexToAddress(ctx.Config.EVMPortalAddress)
	portal, err := contracts.NewOChainPortal(address, client)
	if err != nil {
		return &abcitypes.ResponseCheckTx{
			Code: types.InvalidTransactionError,
		}
	}

	for _, vLog := range remoteTxReceipt.Logs {

		_, err := portal.ParseUSDDeposited(*vLog)
		if err != nil {
			continue
		}

		return &abcitypes.ResponseCheckTx{
			Code: types.NoError,
		}
	}

	return &abcitypes.ResponseCheckTx{
		Code: types.InvalidTransactionError,
	}

}

func (tx OChainreditDepositTransaction) Execute(ctx TransactionContext) *abcitypes.ExecTxResult {
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

		log, err := portal.ParseUSDDeposited(*vLog)
		if err != nil {
			continue
		}

		globalAccount, err := ctx.Db.GlobalsAccounts.GetAt(log.Receiver.Hex(), uint64(ctx.Date.Unix()))
		if err != nil {
			return &abcitypes.ExecTxResult{
				Code: types.NoError,
			}
		}

		globalAccount.CreditBalance += log.Amount.Uint64()
		ctx.Db.GlobalsAccounts.Update(globalAccount)

		return &abcitypes.ExecTxResult{
			Code: types.NoError,
		}
	}

	return &abcitypes.ExecTxResult{
		Code: types.InvalidTransactionError,
	}

}

func ParseNewOChainCreditDepositTransaction(tx Transaction) (OChainreditDepositTransaction, error) {
	var txData OChainreditDepositTransactionData
	err := cbor.Unmarshal(tx.Data, &txData)

	if err != nil {
		return OChainreditDepositTransaction{}, err
	}

	return OChainreditDepositTransaction{
		Type: tx.Type,
		Data: txData,
	}, nil
}
