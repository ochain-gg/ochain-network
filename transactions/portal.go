package transactions

import (
	"context"
	"errors"
	"log"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/fxamacker/cbor/v2"
	"github.com/ochain-gg/ochain-network/contracts"
	"github.com/ochain-gg/ochain-network/database"
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
	Type      OChainPortalInteractionType `cbor:"1,keyasint"`
	Arguments []byte                      `cbor:"2,keyasint"`
}

type OChainPortalInteractionTransaction struct {
	Type TransactionType
	Data OChainPortalInteractionTransactionData
}

func (tx OChainPortalInteractionTransaction) Check(ctx TransactionContext) error {
	switch tx.Data.Type {
	case NewValidatorPortalInteractionType:
		newValidatorTx, err := ParseNewValidatorTransaction(tx)
		if err != nil {
			return err
		}

		return newValidatorTx.Check(ctx)

	case RemoveValidatorPortalInteractionType:
		removeValidatorTx, err := ParseRemoveValidatorTransaction(tx)
		if err != nil {
			return err
		}

		_, err = removeValidatorTx.Check(ctx)
		return err

	case OChainTokenDepositPortalInteractionType:
		tokenDepositTx, err := ParseTokenDepositTransaction(tx)
		if err != nil {
			return err
		}

		_, err = tokenDepositTx.Check(ctx)
		return err
	}

	return errors.New("portal interaction type not exists")
}

func (tx OChainPortalInteractionTransaction) Execute(ctx TransactionContext) error {
	switch tx.Data.Type {
	case NewValidatorPortalInteractionType:
		newValidatorTx, err := ParseNewValidatorTransaction(tx)
		if err != nil {
			return err
		}

		err = newValidatorTx.Execute(ctx)
		return err

	case RemoveValidatorPortalInteractionType:
		removeValidatorTx, err := ParseRemoveValidatorTransaction(tx)
		if err != nil {
			return err
		}

		err = removeValidatorTx.Execute(ctx)
		return err

	case OChainTokenDepositPortalInteractionType:
		tokenDepositTx, err := ParseTokenDepositTransaction(tx)
		if err != nil {
			return err
		}

		err = tokenDepositTx.Execute(ctx)
		return err
	}

	return errors.New("portal interaction type not exists")
}

func ParseNewOChainPortalInteraction(tx Transaction) (OChainPortalInteractionTransaction, error) {
	var txData OChainPortalInteractionTransactionData
	err := cbor.Unmarshal(tx.Data, &txData)

	if err != nil {
		return OChainPortalInteractionTransaction{}, err
	}

	return OChainPortalInteractionTransaction{
		Type: tx.Type,
		Data: txData,
	}, nil
}

// New validator transaction
type NewValidatorTransactionDataArguments struct {
	ValidatorId           uint64 `cbor:"1,keyasint"`
	RemoteTransactionHash string `cbor:"2,keyasint"`
	PublicKey             string `cbor:"3,keyasint"`
}

type NewValidatorTransactionDataRaw struct {
	Type      OChainPortalInteractionType `cbor:"1,keyasint"`
	Arguments []byte                      `cbor:"2,keyasint"`
}

type NewValidatorTransactionData struct {
	Type      OChainPortalInteractionType
	Arguments NewValidatorTransactionDataArguments
}

type NewValidatorTransaction struct {
	Type TransactionType
	Data NewValidatorTransactionData
}

type NewValidatorEventData NewValidatorTransactionDataArguments

func (tx NewValidatorTransaction) Check(ctx TransactionContext) error {

	client, err := ethclient.Dial(ctx.Config.EVMRpc)
	if err != nil {
		log.Fatal(err)
	}

	remoteTx, _, err := client.TransactionByHash(context.Background(), common.HexToHash(tx.Data.Arguments.RemoteTransactionHash))
	if err != nil {
		return err
	}

	if remoteTx.ChainId().Uint64() != ctx.Config.EVMChainId {
		log.Fatal(errors.New("rpc chainId don't match"))
	}

	if *remoteTx.To() != common.HexToAddress(ctx.Config.EVMPortalAddress) {
		return errors.New("wrong to address")
	}

	return nil
}

func (tx NewValidatorTransaction) FetchData(ctx TransactionContext) (contracts.OChainPortalOChainNewValidator, error) {

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
			continue
		}

		if event.Raw.Address == common.HexToAddress(ctx.Config.EVMPortalAddress) {
			if event.PublicKey == tx.Data.Arguments.PublicKey {
				_, err = ctx.Db.Validators.Get(event.ValidatorId.Uint64(), ctx.Txn)
				if err == nil {
					return contracts.OChainPortalOChainNewValidator{}, errors.New("validator already created")
				}

				return *event, nil
			}
		}
	}

	return contracts.OChainPortalOChainNewValidator{}, errors.New("invalid tx")
}

func (tx NewValidatorTransaction) Execute(ctx TransactionContext) error {

	event, err := tx.FetchData(ctx)
	if err != nil {
		return err
	}

	_, err = ctx.Db.Validators.Get(event.ValidatorId.Uint64(), ctx.Txn)
	if err == nil {
		return errors.New("validator already created")
	}

	err = ctx.Db.Validators.Insert(database.OChainValidator{
		Id:        event.ValidatorId.Uint64(),
		Stacker:   event.Stacker.Hex(),
		Validator: event.Validator.Hex(),
		PublicKey: event.PublicKey,
		Enabled:   true,
	}, ctx.Txn)

	if err != nil {
		return err
	}

	state, err := ctx.Db.BridgeState.Get()
	if err != nil {
		return err
	}

	if event.Raw.BlockNumber > state.LatestPortalUpdate {
		state.SetLatestPortalUpdate(event.Raw.BlockNumber)
		err = ctx.Db.BridgeState.Save(state, ctx.Txn)
	}

	return err
}

func (tx NewValidatorTransaction) Transaction() (Transaction, error) {
	txDataArgs, err := cbor.Marshal(tx.Data.Arguments)
	if err != nil {
		return Transaction{}, err
	}

	txData, err := cbor.Marshal(NewValidatorTransactionDataRaw{
		Type:      tx.Data.Type,
		Arguments: txDataArgs,
	})
	if err != nil {
		return Transaction{}, err
	}

	return Transaction{
		Type: tx.Type,
		Data: txData,
	}, nil
}

func ParseNewValidatorTransaction(tx OChainPortalInteractionTransaction) (NewValidatorTransaction, error) {

	var txDataArgs NewValidatorTransactionDataArguments
	err := cbor.Unmarshal([]byte(tx.Data.Arguments), &txDataArgs)

	if err != nil {
		return NewValidatorTransaction{}, err
	}

	return NewValidatorTransaction{
		Type: tx.Type,
		Data: NewValidatorTransactionData{
			Type:      tx.Data.Type,
			Arguments: txDataArgs,
		},
	}, nil
}

// Remove validator
type RemoveValidatorTransactionDataArguments struct {
	ValidatorId           uint64 `cbor:"1,keyasint"`
	RemoteTransactionHash string `cbor:"2,keyasint"`
}

type RemoveValidatorTransactionDataRaw struct {
	Type      OChainPortalInteractionType `cbor:"1,keyasint"`
	Arguments []byte                      `cbor:"2,keyasint"`
}

type RemoveValidatorTransactionData struct {
	Type      OChainPortalInteractionType
	Arguments RemoveValidatorTransactionDataArguments
}

type RemoveValidatorTransaction struct {
	Type TransactionType
	Data RemoveValidatorTransactionData
}

type RemoveValidatorEventData RemoveValidatorTransactionDataArguments

func (tx RemoveValidatorTransaction) Check(ctx TransactionContext) (contracts.OChainPortalOChainRemoveValidator, error) {

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
	if err != nil {
		return err
	}

	state, err := ctx.Db.BridgeState.Get()
	if err != nil {
		return err
	}

	if event.Raw.BlockNumber > state.LatestPortalUpdate {
		state.SetLatestPortalUpdate(event.Raw.BlockNumber)
		err = ctx.Db.BridgeState.Save(state, ctx.Txn)
	}

	return err
}

func (tx RemoveValidatorTransaction) Transaction() (Transaction, error) {
	txDataArgs, err := cbor.Marshal(tx.Data.Arguments)
	if err != nil {
		return Transaction{}, err
	}

	txData, err := cbor.Marshal(RemoveValidatorTransactionDataRaw{
		Type:      tx.Data.Type,
		Arguments: txDataArgs,
	})
	if err != nil {
		return Transaction{}, err
	}

	return Transaction{
		Type: tx.Type,
		Data: txData,
	}, nil
}

func ParseRemoveValidatorTransaction(tx OChainPortalInteractionTransaction) (RemoveValidatorTransaction, error) {

	var txDataArgs RemoveValidatorTransactionDataArguments
	err := cbor.Unmarshal([]byte(tx.Data.Arguments), &txDataArgs)

	if err != nil {
		return RemoveValidatorTransaction{}, err
	}

	return RemoveValidatorTransaction{
		Type: tx.Type,
		Data: RemoveValidatorTransactionData{
			Type:      tx.Data.Type,
			Arguments: txDataArgs,
		},
	}, nil
}

// Token deposit
type TokenDepositTransactionDataArguments struct {
	RemoteTransactionHash string `cbor:"1,keyasint"`
}

type TokenDepositTransactionDataRaw struct {
	Type      OChainPortalInteractionType `cbor:"1,keyasint"`
	Arguments []byte                      `cbor:"2,keyasint"`
}

type TokenDepositTransactionData struct {
	Type      OChainPortalInteractionType
	Arguments TokenDepositTransactionDataArguments
}

type TokenDepositTransaction struct {
	Hash string
	Type TransactionType
	Data TokenDepositTransactionData
}

type TokenDepositEventData TokenDepositTransactionDataArguments

func (tx TokenDepositTransaction) Check(ctx TransactionContext) (contracts.OChainPortalOChainTokenDeposit, error) {

	client, err := ethclient.Dial(ctx.Config.EVMRpc)
	if err != nil {
		log.Fatal(err)
	}

	remoteTx, _, err := client.TransactionByHash(context.Background(), common.HexToHash(tx.Data.Arguments.RemoteTransactionHash))
	if err != nil {
		return contracts.OChainPortalOChainTokenDeposit{}, err
	}

	if remoteTx.ChainId().Uint64() != ctx.Config.EVMChainId {
		log.Fatal(errors.New("rpc chainId don't match"))
	}

	if *remoteTx.To() != common.HexToAddress(ctx.Config.EVMPortalAddress) {
		return contracts.OChainPortalOChainTokenDeposit{}, errors.New("wrong to address")
	}

	remoteTxReceipt, err := client.TransactionReceipt(context.Background(), common.HexToHash(tx.Data.Arguments.RemoteTransactionHash))
	if err != nil {
		return contracts.OChainPortalOChainTokenDeposit{}, err
	}

	if remoteTxReceipt.Status != 1 {
		return contracts.OChainPortalOChainTokenDeposit{}, errors.New("non valid transaction")
	}

	address := common.HexToAddress(ctx.Config.EVMPortalAddress)
	portal, err := contracts.NewOChainPortal(address, client)
	if err != nil {
		return contracts.OChainPortalOChainTokenDeposit{}, err
	}

	for _, vLog := range remoteTxReceipt.Logs {

		event, err := portal.ParseOChainTokenDeposit(*vLog)
		if err != nil {
			continue
		}

		return *event, nil
	}

	return contracts.OChainPortalOChainTokenDeposit{}, errors.New("invalid tx")
}

func (tx *TokenDepositTransaction) Execute(ctx TransactionContext) error {

	event, err := tx.Check(ctx)
	if err != nil {
		return err
	}

	acc := database.OChainBridgeTransaction{
		Type:            database.OChainBridgeDepositTransaction,
		TransactionHash: event.Raw.TxHash.Hex(),
		Account:         event.Receiver.Hex(),
		Amount:          event.Amount.Uint64(),
		Executed:        false,
		Canceled:        false,
	}

	err = ctx.Db.BridgeTransactions.Insert(acc, ctx.Txn)
	if err != nil {
		return err
	}

	state, err := ctx.Db.BridgeState.Get()
	if err != nil {
		return err
	}

	if event.Raw.BlockNumber > state.LatestPortalUpdate {
		state.SetLatestPortalUpdate(event.Raw.BlockNumber)
		err = ctx.Db.BridgeState.Save(state, ctx.Txn)
	}

	return err
}

func (tx TokenDepositTransaction) Transaction() (Transaction, error) {
	txDataArgs, err := cbor.Marshal(tx.Data.Arguments)
	if err != nil {
		return Transaction{}, err
	}

	txData, err := cbor.Marshal(TokenDepositTransactionDataRaw{
		Type:      tx.Data.Type,
		Arguments: txDataArgs,
	})
	if err != nil {
		return Transaction{}, err
	}

	return Transaction{
		Type: tx.Type,
		Data: txData,
	}, nil
}

func ParseTokenDepositTransaction(tx OChainPortalInteractionTransaction) (TokenDepositTransaction, error) {

	var txDataArgs TokenDepositTransactionDataArguments
	err := cbor.Unmarshal([]byte(tx.Data.Arguments), &txDataArgs)

	if err != nil {
		return TokenDepositTransaction{}, err
	}

	return TokenDepositTransaction{
		Type: tx.Type,
		Data: TokenDepositTransactionData{
			Type:      tx.Data.Type,
			Arguments: txDataArgs,
		},
	}, nil
}
