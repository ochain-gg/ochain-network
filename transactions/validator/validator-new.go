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

type OChainNewValidatorTransactionData struct {
	ValidatorId           uint64 `cbor:"1,keyasint"`
	RemoteTransactionHash string `cbor:"2,keyasint"`
	PublicKey             string `cbor:"3,keyasint"`
}

type OChainNewValidatorTransaction struct {
	Type t.TransactionType
	Data OChainNewValidatorTransactionData
}

func (tx *OChainNewValidatorTransaction) Check(ctx t.TransactionContext) *abcitypes.CheckTxResponse {

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
		log.Fatal(errors.New("rpc chainId don't match"))
	}

	if *remoteTx.To() != common.HexToAddress(ctx.Config.EVMPortalAddress) {
		return &abcitypes.CheckTxResponse{
			Code: types.InvalidTransactionError,
		}
	}

	remoteParsedTx, err := tx.FetchData(ctx)
	if err != nil {
		return &abcitypes.CheckTxResponse{
			Code: types.InvalidTransactionError,
		}
	}

	if remoteParsedTx.PublicKey != tx.Data.PublicKey {
		return &abcitypes.CheckTxResponse{
			Code: types.InvalidTransactionError,
		}
	}

	validatorIsEnabled, err := ctx.Db.Validators.IsEnabled(remoteParsedTx.PublicKey)
	if err != nil {
		return &abcitypes.CheckTxResponse{
			Code: types.InvalidTransactionError,
		}
	}

	if validatorIsEnabled {
		return &abcitypes.CheckTxResponse{
			Code: types.InvalidTransactionError,
		}
	}

	return &abcitypes.CheckTxResponse{
		Code: types.NoError,
	}
}

func (tx *OChainNewValidatorTransaction) FetchData(ctx t.TransactionContext) (contracts.OChainPortalOChainNewValidator, error) {

	client, err := ethclient.Dial(ctx.Config.EVMRpc)
	if err != nil {
		log.Fatal(err)
	}

	remoteTx, _, err := client.TransactionByHash(context.Background(), common.HexToHash(tx.Data.RemoteTransactionHash))
	if err != nil {
		return contracts.OChainPortalOChainNewValidator{}, err
	}

	if remoteTx.ChainId().Uint64() != ctx.Config.EVMChainId {
		log.Fatal(errors.New("rpc chainId don't match"))
	}

	if *remoteTx.To() != common.HexToAddress(ctx.Config.EVMPortalAddress) {
		return contracts.OChainPortalOChainNewValidator{}, errors.New("wrong to address")
	}

	remoteTxReceipt, err := client.TransactionReceipt(context.Background(), common.HexToHash(tx.Data.RemoteTransactionHash))
	if err != nil {
		return contracts.OChainPortalOChainNewValidator{}, err
	}

	if remoteTxReceipt.Status != 1 {
		return contracts.OChainPortalOChainNewValidator{}, errors.New("non valid transaction")
	}

	address := common.HexToAddress(ctx.Config.EVMPortalAddress)
	portal, err := contracts.NewOChainPortal(address, client)
	if err != nil {
		return contracts.OChainPortalOChainNewValidator{}, err
	}

	for _, vLog := range remoteTxReceipt.Logs {

		event, err := portal.ParseOChainNewValidator(*vLog)
		if err != nil {
			continue
		}

		if event.Raw.Address == common.HexToAddress(ctx.Config.EVMPortalAddress) {
			if event.PublicKey == tx.Data.PublicKey {
				_, err = ctx.Db.Validators.GetByIdAt(event.ValidatorId.Uint64(), uint64(ctx.Date.Unix()))
				if err == nil {
					return contracts.OChainPortalOChainNewValidator{}, errors.New("validator already created")
				}

				return *event, nil
			}
		}
	}

	return contracts.OChainPortalOChainNewValidator{}, errors.New("invalid tx")
}

func (tx *OChainNewValidatorTransaction) Execute(ctx t.TransactionContext) *abcitypes.ExecTxResult {

	event, err := tx.FetchData(ctx)
	if err != nil {
		return &abcitypes.ExecTxResult{
			Code: types.InvalidTransactionError,
		}
	}

	_, err = ctx.Db.Validators.GetByAddressAt(event.PublicKey, uint64(ctx.Date.Unix()))
	if err == nil {
		return &abcitypes.ExecTxResult{
			Code: types.InvalidTransactionError,
		}
	}

	err = ctx.Db.Validators.Insert(types.OChainValidator{
		Id:        event.ValidatorId.Uint64(),
		Stacker:   event.Stacker.Hex(),
		Validator: event.Validator.Hex(),
		Power:     10000,
		PublicKey: event.PublicKey,
		Enabled:   true,
	})

	if err != nil {
		return &abcitypes.ExecTxResult{
			Code: types.InvalidTransactionError,
		}
	}

	if event.Raw.BlockNumber > ctx.State.LatestPortalUpdate {
		ctx.State.LatestPortalUpdate = event.Raw.BlockNumber
	}

	return &abcitypes.ExecTxResult{
		Code: types.NoError,
	}
}

func (tx *OChainNewValidatorTransaction) Transaction() (t.Transaction, error) {

	txData, err := cbor.Marshal(tx.Data)
	if err != nil {
		return t.Transaction{}, err
	}

	return t.Transaction{
		Type: tx.Type,
		Data: txData,
	}, nil
}

func ParseOChainNewValidatorTransaction(tx t.Transaction) (OChainNewValidatorTransaction, error) {

	var txData OChainNewValidatorTransactionData
	err := cbor.Unmarshal([]byte(tx.Data), &txData)

	if err != nil {
		return OChainNewValidatorTransaction{}, err
	}

	return OChainNewValidatorTransaction{
		Type: tx.Type,
		Data: txData,
	}, nil
}
