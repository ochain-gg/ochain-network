package transactions

import (
	"context"
	"encoding/json"
	"errors"
	"log"

	"github.com/dgraph-io/badger"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ochain.gg/ochain-validator-network/config"
	"github.com/ochain.gg/ochain-validator-network/contracts"
	"github.com/ochain.gg/ochain-validator-network/database"
)

// New validator transaction
type NewValidatorTransactionData struct {
	ValidatorId           uint64 `json:"validatorId"`
	RemoteTransactionHash string `json:"remoteTransactionHash"`
	PublicKey             string `json:"publicKey"`
}

type NewValidatorTransaction struct {
	Type         TransactionType             `json:"type"`
	Data         []byte                      `json:"data"`
	FormatedData NewValidatorTransactionData `json:"formatedData"`
}

type NewValidatorEventData NewValidatorTransactionData

func (tx *NewValidatorTransaction) Verify(cfg config.OChainConfig) (contracts.OChainPortalOChainNewValidator, error) {

	client, err := ethclient.Dial(cfg.EVMRpc)
	if err != nil {
		log.Fatal(err)
	}

	remoteTx, _, err := client.TransactionByHash(context.Background(), common.HexToHash(tx.FormatedData.RemoteTransactionHash))
	if err != nil {
		return contracts.OChainPortalOChainNewValidator{}, err
	}

	if remoteTx.ChainId().Uint64() != cfg.EVMChainId {
		log.Fatal(errors.New("rpc chainId don't match"))
	}

	if *remoteTx.To() != common.HexToAddress(cfg.EVMPortalAddress) {
		return contracts.OChainPortalOChainNewValidator{}, errors.New("wrong to address")
	}

	remoteTxReceipt, err := client.TransactionReceipt(context.Background(), common.HexToHash(tx.FormatedData.RemoteTransactionHash))
	if err != nil {
		return contracts.OChainPortalOChainNewValidator{}, err
	}

	if remoteTxReceipt.Status != 1 {
		return contracts.OChainPortalOChainNewValidator{}, errors.New("non valid transaction")
	}

	address := common.HexToAddress(cfg.EVMPortalAddress)
	portal, err := contracts.NewOChainPortal(address, client)
	if err != nil {
		return contracts.OChainPortalOChainNewValidator{}, err
	}

	for _, vLog := range remoteTxReceipt.Logs {

		event, err := portal.ParseOChainNewValidator(*vLog)
		if err != nil {
			return contracts.OChainPortalOChainNewValidator{}, err
		}

		if event.Raw.Address == common.HexToAddress(cfg.EVMPortalAddress) {
			if event.PublicKey == tx.FormatedData.PublicKey {
				return *event, nil
			}
		}
	}

	return contracts.OChainPortalOChainNewValidator{}, errors.New("invalid tx")
}

func (tx *NewValidatorTransaction) Execute(cfg config.OChainConfig, db *database.OChainDatabase, txn *badger.Txn) error {

	event, err := tx.Verify(cfg)
	if err != nil {
		return err
	}

	_, err = db.Validators.Get(event.ValidatorId.Uint64(), txn)
	if err == nil {
		return errors.New("validator already activated")
	}

	err = db.Validators.Insert(database.OChainValidator{
		Id:        event.ValidatorId.Uint64(),
		Stacker:   event.Stacker.Hex(),
		Validator: event.Validator.Hex(),
		PublicKey: event.PublicKey,
		Enabled:   true,
	}, txn)

	return err
}

func ParseNewValidatorTransaction(tx Transaction) (NewValidatorTransaction, error) {
	var txData NewValidatorTransactionData
	err := json.Unmarshal(tx.Data, &txData)

	if err != nil {
		return NewValidatorTransaction{}, err
	}

	return NewValidatorTransaction{
		Type:         tx.Type,
		Data:         tx.Data,
		FormatedData: txData,
	}, nil
}

// Remove validator
type RemoveValidatorTransactionData struct {
	ValidatorId           uint64 `json:"validatorId"`
	RemoteTransactionHash string `json:"remoteTransactionHash"`
}

type RemoveValidatorTransaction struct {
	Type         TransactionType                `json:"type"`
	Data         []byte                         `json:"data"`
	FormatedData RemoveValidatorTransactionData `json:"formatedData"`
}

type RemoveValidatorEventData RemoveValidatorTransactionData

func (tx *RemoveValidatorTransaction) Verify(cfg config.OChainConfig) (contracts.OChainPortalOChainRemoveValidator, error) {

	client, err := ethclient.Dial(cfg.EVMRpc)
	if err != nil {
		log.Fatal(err)
	}

	remoteTx, _, err := client.TransactionByHash(context.Background(), common.HexToHash(string(tx.FormatedData.RemoteTransactionHash)))
	if err != nil {
		return contracts.OChainPortalOChainRemoveValidator{}, err
	}

	if remoteTx.ChainId().Uint64() != cfg.EVMChainId {
		log.Fatal(errors.New("rpc chainId don't match"))
	}

	if *remoteTx.To() != common.HexToAddress(cfg.EVMPortalAddress) {
		return contracts.OChainPortalOChainRemoveValidator{}, errors.New("wrong to address")
	}

	remoteTxReceipt, err := client.TransactionReceipt(context.Background(), common.HexToHash(string(tx.FormatedData.RemoteTransactionHash)))
	if err != nil {
		return contracts.OChainPortalOChainRemoveValidator{}, err
	}

	if remoteTxReceipt.Status != 1 {
		return contracts.OChainPortalOChainRemoveValidator{}, errors.New("non valid transaction")
	}

	address := common.HexToAddress(cfg.EVMPortalAddress)
	portal, err := contracts.NewOChainPortal(address, client)
	if err != nil {
		return contracts.OChainPortalOChainRemoveValidator{}, err
	}

	for _, vLog := range remoteTxReceipt.Logs {

		event, err := portal.ParseOChainRemoveValidator(*vLog)
		if err != nil {
			return contracts.OChainPortalOChainRemoveValidator{}, err
		}

		if event.Raw.Address == common.HexToAddress(cfg.EVMPortalAddress) {
			if event.ValidatorId.Uint64() == tx.FormatedData.ValidatorId {
				return *event, nil
			}
		}
	}

	return contracts.OChainPortalOChainRemoveValidator{}, errors.New("invalid tx")
}

func (tx *RemoveValidatorTransaction) Execute(cfg config.OChainConfig, db *database.OChainDatabase, txn *badger.Txn) error {

	event, err := tx.Verify(cfg)
	if err != nil {
		return err
	}

	validator, err := db.Validators.Get(event.ValidatorId.Uint64(), txn)
	if err != nil {
		return err
	}

	validator.Enabled = false
	err = db.Validators.Save(validator, txn)

	return err
}

func ParseRemoveValidatorTransaction(tx Transaction) (RemoveValidatorTransaction, error) {

	var txData RemoveValidatorTransactionData
	err := json.Unmarshal(tx.Data, &txData)

	if err != nil {
		return RemoveValidatorTransaction{}, err
	}

	return RemoveValidatorTransaction{
		Type:         tx.Type,
		Data:         tx.Data,
		FormatedData: txData,
	}, nil
}
