package transactions

import (
	"errors"
	"fmt"
	"log"

	abcitypes "github.com/cometbft/cometbft/abci/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/signer/core/apitypes"
	"github.com/fxamacker/cbor/v2"

	"github.com/ochain-gg/ochain-network/types"
)

type RegisterAccountTransactionData struct {
	Address             string   `cbor:"1,keyasint"`
	GuardianQuorum      uint64   `cbor:"2,keyasint"`
	Guardians           []string `cbor:"3,keyasint"`
	DeleguatedTo        []string `cbor:"4,keyasint"`
	AuthorizerSignature []byte   `cbor:"5,keyasint"`
}

type RegisterAccountTransaction struct {
	Type      TransactionType                `cbor:"1,keyasint"`
	From      string                         `cbor:"2,keyasint"`
	Nonce     uint64                         `cbor:"3,keyasint"`
	Data      RegisterAccountTransactionData `cbor:"4,keyasint"`
	Signature []byte                         `cbor:"5,keyasint"`
}

func (tx *RegisterAccountTransaction) Transaction() (Transaction, error) {
	txData, err := cbor.Marshal(tx.Data)
	if err != nil {
		return Transaction{}, err
	}

	return Transaction{
		Type:      tx.Type,
		From:      tx.From,
		Nonce:     tx.Nonce,
		Data:      txData,
		Signature: tx.Signature,
	}, nil
}

func (tx *RegisterAccountTransaction) Check(ctx TransactionContext) error {
	_, err := ctx.Db.GlobalsAccounts.Get(tx.Data.Address)
	if err == nil {
		return errors.New("account aleady exists")
	}

	typedData := apitypes.TypedData{
		Types: apitypes.Types{
			"EIP712Domain": []apitypes.Type{
				{Name: "name", Type: "string"},
				{Name: "version", Type: "string"},
				{Name: "chainId", Type: "uint256"},
				{Name: "verifyingContract", Type: "address"},
			},
			"RegisterRequest": []apitypes.Type{
				{Name: "address", Type: "address"},
			},
		},
		PrimaryType: "RegisterRequest",
		Domain: apitypes.TypedDataDomain{
			Name:              "OChainNetwork",
			Version:           "1",
			ChainId:           math.NewHexOrDecimal256(20291),
			VerifyingContract: "0x0000000000000000000000000000000000000000",
		},
		Message: map[string]interface{}{
			"address": tx.From,
		},
	}
	// EIP-712 typed data marshalling
	domainSeparator, err := typedData.HashStruct("EIP712Domain", typedData.Domain.Map())
	if err != nil {
		return err
	}
	typedDataHash, err := typedData.HashStruct(typedData.PrimaryType, typedData.Message)
	if err != nil {
		return err
	}

	typedDataHashSigned := []byte(fmt.Sprintf("\x19\x01%s%s", string(domainSeparator), string(typedDataHash)))
	sighash := crypto.Keccak256(typedDataHashSigned)

	signature := tx.Data.AuthorizerSignature
	if signature[64] == 27 || signature[64] == 28 {
		signature[64] -= 27
	}

	sigPubkey, err := crypto.SigToPub(sighash, signature)
	if err != nil {
		log.Println("SigToPub: " + err.Error())
		return err
	}

	address := crypto.PubkeyToAddress(*sigPubkey)
	if address != common.HexToAddress("0x190144001306820e9BdF6eB2dB8d747B4fCE7980") {
		return errors.New("bad registration signature")
	}

	return nil
}

func (tx *RegisterAccountTransaction) Execute(ctx TransactionContext) ([]abcitypes.Event, error) {
	err := tx.Check(ctx)
	if err != nil {
		return []abcitypes.Event{}, err
	}

	account := types.OChainGlobalAccount{
		Address:      tx.From,
		Nonce:        1,
		TokenBalance: "0",
		USDBalance:   "0",
		IAM: types.OChainGlobalAccountIAM{
			GuardianQuorum: tx.Data.GuardianQuorum,
			Guardians:      tx.Data.Guardians,
			DeleguatedTo:   tx.Data.DeleguatedTo,
		},
		CreatedAt: ctx.Date.Unix(),
	}

	err = ctx.Db.GlobalsAccounts.Insert(account)
	if err != nil {
		return []abcitypes.Event{}, err
	}

	events := []abcitypes.Event{
		{
			Type: "GlobalAccountRegistered",
			Attributes: []abcitypes.EventAttribute{
				{Key: "address", Value: tx.From, Index: true},
			},
		},
	}

	return events, nil
}

func ParseRegisterAccountTransaction(tx Transaction) (RegisterAccountTransaction, error) {
	var txData RegisterAccountTransactionData
	err := cbor.Unmarshal(tx.Data, &txData)

	if err != nil {
		return RegisterAccountTransaction{}, err
	}

	return RegisterAccountTransaction{
		Type:      tx.Type,
		From:      tx.From,
		Nonce:     tx.Nonce,
		Data:      txData,
		Signature: tx.Signature,
	}, nil
}
