package validator_transactions

import (
	"context"
	"errors"
	"log"

	abcitypes "github.com/cometbft/cometbft/abci/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/fxamacker/cbor/v2"
	"github.com/ochain-gg/ochain-network/contracts"
	t "github.com/ochain-gg/ochain-network/transactions"
	"github.com/ochain-gg/ochain-network/types"
)

type OChainBridgeRemoveValidatorTransactionData struct {
	ValidatorId           uint64 `cbor:"1,keyasint"`
	RemoteTransactionHash string `cbor:"2,keyasint"`
}

type OChainBridgeRemoveValidatorTransaction struct {
	Type t.TransactionType
	Data OChainBridgeRemoveValidatorTransactionData
}

func (tx *OChainBridgeRemoveValidatorTransaction) Check(ctx t.TransactionContext) *abcitypes.ResponseCheckTx {

	data, err := tx.FetchData(ctx)
	if err != nil {
		return &abcitypes.ResponseCheckTx{
			Code: types.InvalidTransactionError,
		}
	}

	validator, err := ctx.Db.Validators.GetById(data.ValidatorId.Uint64())
	if err == nil {
		return &abcitypes.ResponseCheckTx{
			Code: types.InvalidTransactionError,
		}
	}

	if !validator.Enabled {
		return &abcitypes.ResponseCheckTx{
			Code: types.InvalidTransactionError,
		}
	}

	return &abcitypes.ResponseCheckTx{
		Code: types.NoError,
	}
}

func (tx *OChainBridgeRemoveValidatorTransaction) FetchData(ctx t.TransactionContext) (contracts.OChainPortalOChainRemoveValidator, error) {

	client, err := ethclient.Dial(ctx.Config.EVMRpc)
	if err != nil {
		log.Fatal(err)
	}

	remoteTx, _, err := client.TransactionByHash(context.Background(), common.HexToHash(tx.Data.RemoteTransactionHash))
	if err != nil {
		return contracts.OChainPortalOChainRemoveValidator{}, err
	}

	if remoteTx.ChainId().Uint64() != ctx.Config.EVMChainId {
		log.Fatal(errors.New("rpc chainId don't match"))
	}

	if *remoteTx.To() != common.HexToAddress(ctx.Config.EVMPortalAddress) {
		return contracts.OChainPortalOChainRemoveValidator{}, errors.New("wrong to address")
	}

	remoteTxReceipt, err := client.TransactionReceipt(context.Background(), common.HexToHash(tx.Data.RemoteTransactionHash))
	if err != nil {
		return contracts.OChainPortalOChainRemoveValidator{}, err
	}

	if remoteTxReceipt.Status != 1 {
		return contracts.OChainPortalOChainRemoveValidator{}, errors.New("non valid transaction")
	}

	address := common.HexToAddress(ctx.Config.EVMPortalAddress)
	portal, err := contracts.NewOChainPortal(address, client)
	if err != nil {
		return contracts.OChainPortalOChainRemoveValidator{}, err
	}

	for _, vLog := range remoteTxReceipt.Logs {

		event, err := portal.ParseOChainRemoveValidator(*vLog)
		if err != nil {
			continue
		}

		if event.Raw.Address == common.HexToAddress(ctx.Config.EVMPortalAddress) {
			if event.ValidatorId.Uint64() == tx.Data.ValidatorId {
				return *event, nil
			}
		}
	}

	return contracts.OChainPortalOChainRemoveValidator{}, errors.New("invalid tx")
}

func (tx *OChainBridgeRemoveValidatorTransaction) Execute(ctx t.TransactionContext) *abcitypes.ExecTxResult {

	result := tx.Check(ctx)
	if result.Code != types.NoError {
		return &abcitypes.ExecTxResult{
			Code: result.GetCode(),
		}
	}

	data, err := tx.FetchData(ctx)
	if err != nil {
		return &abcitypes.ExecTxResult{
			Code: types.InvalidTransactionError,
		}
	}

	validator, err := ctx.Db.Validators.GetByIdAt(data.ValidatorId.Uint64(), uint64(ctx.Date.Unix()))
	if err != nil {
		return &abcitypes.ExecTxResult{
			Code: types.InvalidTransactionError,
		}
	}

	validator.Enabled = false
	err = ctx.Db.Validators.Update(validator)
	if err != nil {
		return &abcitypes.ExecTxResult{
			Code: types.InvalidTransactionError,
		}
	}

	if data.Raw.BlockNumber > ctx.State.LatestPortalUpdate {
		ctx.State.SetLatestPortalUpdate(data.Raw.BlockNumber)
	}

	return &abcitypes.ExecTxResult{
		Code: types.NoError,
	}
}

func (tx OChainBridgeRemoveValidatorTransaction) Transaction() (t.Transaction, error) {

	txData, err := cbor.Marshal(tx.Data)
	if err != nil {
		return t.Transaction{}, err
	}

	return t.Transaction{
		Type: tx.Type,
		Data: txData,
	}, nil
}

func ParseOChainBridgeRemoveValidatorTransaction(tx t.Transaction) (OChainBridgeRemoveValidatorTransaction, error) {

	var txData OChainBridgeRemoveValidatorTransactionData
	err := cbor.Unmarshal(tx.Data, &txData)

	if err != nil {
		return OChainBridgeRemoveValidatorTransaction{}, err
	}

	return OChainBridgeRemoveValidatorTransaction{
		Type: tx.Type,
		Data: txData,
	}, nil
}
