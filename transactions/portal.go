package transactions

import (
	"context"
	"encoding/json"
	"errors"
	"log"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ochain.gg/ochain-network-validator/contracts"
	"github.com/ochain.gg/ochain-network-validator/database"
)

type OChainPortalInteractionType uint64

const (
	NewValidatorPortalInteractionType            OChainPortalInteractionType = 0
	RemoveValidatorPortalInteractionType         OChainPortalInteractionType = 1
	OChainTokenDepositPortalInteractionType      OChainPortalInteractionType = 2
	OChainTokenWithdrawalPortalInteractionType   OChainPortalInteractionType = 3
	OChainBonusSubscriptionPortalInteractionType OChainPortalInteractionType = 4
)

type OChainPortalInteractionTransactionData struct {
	Type      OChainPortalInteractionType `json:"type"`
	Arguments string                      `json:"data"`
}

type OChainPortalInteractionTransaction struct {
	Hash string                                 `json:"hash"`
	Type TransactionType                        `json:"type"`
	Data OChainPortalInteractionTransactionData `json:"data"`
}

func (tx *OChainPortalInteractionTransaction) Check(ctx TransactionContext) error {
	switch tx.Data.Type {
	case NewValidatorPortalInteractionType:
		newValidatorTx, err := ParseNewValidatorTransaction(*tx)
		if err != nil {
			return err
		}

		_, err = newValidatorTx.Check(ctx)
		return err

	case RemoveValidatorPortalInteractionType:
		removeValidatorTx, err := ParseRemoveValidatorTransaction(*tx)
		if err != nil {
			return err
		}

		_, err = removeValidatorTx.Check(ctx)
		return err
	}

	return errors.New("portal interaction type not exists")
}

func (tx *OChainPortalInteractionTransaction) Execute(ctx TransactionContext) error {
	switch tx.Data.Type {
	case NewValidatorPortalInteractionType:
		newValidatorTx, err := ParseNewValidatorTransaction(*tx)
		if err != nil {
			return err
		}

		err = newValidatorTx.Execute(ctx)
		return err

	case RemoveValidatorPortalInteractionType:
		removeValidatorTx, err := ParseRemoveValidatorTransaction(*tx)
		if err != nil {
			return err
		}

		err = removeValidatorTx.Execute(ctx)
		return err
	}

	return errors.New("portal interaction type not exists")
}

func ParseNewOChainPortalInteraction(tx Transaction) (OChainPortalInteractionTransaction, error) {
	var txData OChainPortalInteractionTransactionData
	err := json.Unmarshal([]byte(tx.Data), &txData)

	if err != nil {
		return OChainPortalInteractionTransaction{}, err
	}

	return OChainPortalInteractionTransaction{
		Hash: tx.Hash,
		Type: tx.Type,
		Data: txData,
	}, nil
}

// New validator transaction
type NewValidatorTransactionDataArguments struct {
	ValidatorId           uint64 `json:"validatorId"`
	RemoteTransactionHash string `json:"remoteTransactionHash"`
	PublicKey             string `json:"publicKey"`
}

type NewValidatorTransactionData struct {
	Type      OChainPortalInteractionType          `json:"type"`
	Arguments NewValidatorTransactionDataArguments `json:"arguments"`
}

type NewValidatorTransaction struct {
	Hash string                      `json:"hash"`
	Type TransactionType             `json:"type"`
	Data NewValidatorTransactionData `json:"data"`
}

type NewValidatorEventData NewValidatorTransactionDataArguments

func (tx *NewValidatorTransaction) Check(ctx TransactionContext) (contracts.OChainPortalOChainNewValidator, error) {

	client, err := ethclient.Dial(ctx.Config.EVMRpc)
	if err != nil {
		log.Fatal(err)
	}

	remoteTx, _, err := client.TransactionByHash(context.Background(), common.HexToHash(tx.Data.Arguments.RemoteTransactionHash))
	if err != nil {
		return contracts.OChainPortalOChainNewValidator{}, err
	}

	if remoteTx.ChainId().Uint64() != ctx.Config.EVMChainId {
		log.Fatal(errors.New("rpc chainId don't match"))
	}

	if *remoteTx.To() != common.HexToAddress(ctx.Config.EVMPortalAddress) {
		return contracts.OChainPortalOChainNewValidator{}, errors.New("wrong to address")
	}

	remoteTxReceipt, err := client.TransactionReceipt(context.Background(), common.HexToHash(tx.Data.Arguments.RemoteTransactionHash))
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
			return contracts.OChainPortalOChainNewValidator{}, err
		}

		if event.Raw.Address == common.HexToAddress(ctx.Config.EVMPortalAddress) {
			if event.PublicKey == tx.Data.Arguments.PublicKey {
				return *event, nil
			}
		}
	}

	return contracts.OChainPortalOChainNewValidator{}, errors.New("invalid tx")
}

func (tx *NewValidatorTransaction) Execute(ctx TransactionContext) error {

	event, err := tx.Check(ctx)
	if err != nil {
		return err
	}

	_, err = ctx.Db.Validators.Get(event.ValidatorId.Uint64(), ctx.Txn)
	if err == nil {
		return errors.New("validator already activated")
	}

	err = ctx.Db.Validators.Insert(database.OChainValidator{
		Id:        event.ValidatorId.Uint64(),
		Stacker:   event.Stacker.Hex(),
		Validator: event.Validator.Hex(),
		PublicKey: event.PublicKey,
		Enabled:   true,
	}, ctx.Txn)

	return err
}

func ParseNewValidatorTransaction(tx OChainPortalInteractionTransaction) (NewValidatorTransaction, error) {
	var txDataArgs NewValidatorTransactionDataArguments
	err := json.Unmarshal([]byte(tx.Data.Arguments), &txDataArgs)

	if err != nil {
		return NewValidatorTransaction{}, err
	}

	return NewValidatorTransaction{
		Hash: tx.Hash,
		Type: tx.Type,
		Data: NewValidatorTransactionData{
			Type:      tx.Data.Type,
			Arguments: txDataArgs,
		},
	}, nil
}

// Remove validator
type RemoveValidatorTransactionDataArguments struct {
	ValidatorId           uint64 `json:"validatorId"`
	RemoteTransactionHash string `json:"remoteTransactionHash"`
}

type RemoveValidatorTransactionData struct {
	Type      OChainPortalInteractionType             `json:"type"`
	Arguments RemoveValidatorTransactionDataArguments `json:"arguments"`
}

type RemoveValidatorTransaction struct {
	Hash string                         `json:"hash"`
	Type TransactionType                `json:"type"`
	Data RemoveValidatorTransactionData `json:"data"`
}

type RemoveValidatorEventData RemoveValidatorTransactionDataArguments

func (tx *RemoveValidatorTransaction) Check(ctx TransactionContext) (contracts.OChainPortalOChainRemoveValidator, error) {

	client, err := ethclient.Dial(ctx.Config.EVMRpc)
	if err != nil {
		log.Fatal(err)
	}

	remoteTx, _, err := client.TransactionByHash(context.Background(), common.HexToHash(tx.Data.Arguments.RemoteTransactionHash))
	if err != nil {
		return contracts.OChainPortalOChainRemoveValidator{}, err
	}

	if remoteTx.ChainId().Uint64() != ctx.Config.EVMChainId {
		log.Fatal(errors.New("rpc chainId don't match"))
	}

	if *remoteTx.To() != common.HexToAddress(ctx.Config.EVMPortalAddress) {
		return contracts.OChainPortalOChainRemoveValidator{}, errors.New("wrong to address")
	}

	remoteTxReceipt, err := client.TransactionReceipt(context.Background(), common.HexToHash(tx.Data.Arguments.RemoteTransactionHash))
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
			return contracts.OChainPortalOChainRemoveValidator{}, err
		}

		if event.Raw.Address == common.HexToAddress(ctx.Config.EVMPortalAddress) {
			if event.ValidatorId.Uint64() == tx.Data.Arguments.ValidatorId {
				return *event, nil
			}
		}
	}

	return contracts.OChainPortalOChainRemoveValidator{}, errors.New("invalid tx")
}

func (tx *RemoveValidatorTransaction) Execute(ctx TransactionContext) error {

	event, err := tx.Check(ctx)
	if err != nil {
		return err
	}

	validator, err := ctx.Db.Validators.Get(event.ValidatorId.Uint64(), ctx.Txn)
	if err != nil {
		return err
	}

	validator.Enabled = false
	err = ctx.Db.Validators.Save(validator, ctx.Txn)

	return err
}

func ParseRemoveValidatorTransaction(tx OChainPortalInteractionTransaction) (RemoveValidatorTransaction, error) {

	var txDataArgs RemoveValidatorTransactionDataArguments
	err := json.Unmarshal([]byte(tx.Data.Arguments), &txDataArgs)

	if err != nil {
		return RemoveValidatorTransaction{}, err
	}

	return RemoveValidatorTransaction{
		Hash: tx.Hash,
		Type: tx.Type,
		Data: RemoveValidatorTransactionData{
			Type:      tx.Data.Type,
			Arguments: txDataArgs,
		},
	}, nil
}
