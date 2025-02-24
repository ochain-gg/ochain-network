package account_transactions

import (
	"fmt"
	"log"

	abcitypes "github.com/cometbft/cometbft/abci/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/signer/core/apitypes"
	"github.com/fxamacker/cbor/v2"

	t "github.com/ochain-gg/ochain-network/engine/transactions"
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
	Type      t.TransactionType              `cbor:"1,keyasint"`
	From      string                         `cbor:"2,keyasint"`
	Nonce     uint64                         `cbor:"3,keyasint"`
	Data      RegisterAccountTransactionData `cbor:"4,keyasint"`
	Signature []byte                         `cbor:"5,keyasint"`
}

func (tx *RegisterAccountTransaction) Transaction() (t.Transaction, error) {
	txData, err := cbor.Marshal(tx.Data)
	if err != nil {
		return t.Transaction{}, err
	}

	return t.Transaction{
		Type:      tx.Type,
		From:      tx.From,
		Nonce:     tx.Nonce,
		Data:      txData,
		Signature: tx.Signature,
	}, nil
}

func (tx *RegisterAccountTransaction) Check(ctx t.TransactionContext) *abcitypes.CheckTxResponse {

	log.Println("[RegisterAccountTransaction] Check")

	_, err := ctx.Db.GlobalsAccounts.Get(tx.Data.Address)
	if err == nil {
		return &abcitypes.CheckTxResponse{
			Code: types.InvalidTransactionError,
		}
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
			ChainId:           math.NewHexOrDecimal256(84532),
			VerifyingContract: "0x629c04197012af8e1c4eb92DF8CdA1ed71774488",
		},
		Message: map[string]interface{}{
			"address": tx.From,
		},
	}
	// EIP-712 typed data marshalling
	domainSeparator, err := typedData.HashStruct("EIP712Domain", typedData.Domain.Map())
	if err != nil {
		return &abcitypes.CheckTxResponse{
			Code: types.InvalidTransactionError,
		}
	}

	log.Println("[RegisterAccountTransaction] Domain separator: " + domainSeparator.String())

	typedDataHash, err := typedData.HashStruct(typedData.PrimaryType, typedData.Message)
	if err != nil {
		return &abcitypes.CheckTxResponse{
			Code: types.InvalidTransactionError,
		}
	}

	log.Println("[RegisterAccountTransaction] typedDataHash: " + typedDataHash.String())

	typedDataHashSigned := []byte(fmt.Sprintf("\x19\x01%s%s", string(domainSeparator), string(typedDataHash)))
	sighash := crypto.Keccak256(typedDataHashSigned)

	signature := tx.Data.AuthorizerSignature
	if signature[64] == 27 || signature[64] == 28 {
		signature[64] -= 27
	}

	sigPubkey, err := crypto.SigToPub(sighash, signature)
	if err != nil {
		return &abcitypes.CheckTxResponse{
			Code: types.InvalidTransactionError,
		}
	}

	address := crypto.PubkeyToAddress(*sigPubkey)
	log.Println("[RegisterAccountTransaction] signer address: " + address.String())
	if address != common.HexToAddress("0x190144001306820e9BdF6eB2dB8d747B4fCE7980") {
		return &abcitypes.CheckTxResponse{
			Code: types.InvalidTransactionError,
		}
	}

	return &abcitypes.CheckTxResponse{
		Code: types.NoError,
	}
}

func (tx *RegisterAccountTransaction) Execute(ctx t.TransactionContext) *abcitypes.ExecTxResult {
	response := tx.Check(ctx)
	if response.Code != types.NoError {
		return &abcitypes.ExecTxResult{
			Code: response.Code,
		}
	}

	var faucetAmount uint64 = 10_000_000

	account := types.OChainGlobalAccount{
		Address:        tx.From,
		Nonce:          1,
		TokenBalance:   faucetAmount,
		CreditBalance:  0,
		LastDailyClaim: ctx.Date.Unix(),
		IAM: types.OChainGlobalAccountIAM{
			GuardianQuorum: tx.Data.GuardianQuorum,
			Guardians:      tx.Data.Guardians,
			DeleguatedTo:   tx.Data.DeleguatedTo,
		},
		CreatedAt: ctx.Date.Unix(),
	}

	ctx.State.AddInGameCirculatingSupply(faucetAmount)

	err := ctx.Db.GlobalsAccounts.Insert(account)
	if err != nil {
		return &abcitypes.ExecTxResult{
			Code: types.InvalidTransactionError,
		}
	}

	events := []abcitypes.Event{
		{
			Type: "GlobalAccountRegistered",
			Attributes: []abcitypes.EventAttribute{
				{Key: "address", Value: tx.From, Index: true},
			},
		},
		{
			Type: "FaucetClaimed",
			Attributes: []abcitypes.EventAttribute{
				{Key: "address", Value: tx.From, Index: true},
				{Key: "amount", Value: fmt.Sprint(faucetAmount)},
			},
		},
	}

	receipt := t.TransactionReceipt{
		GasCost: 0,
	}

	receiptBytes, err := cbor.Marshal(receipt)
	if err != nil {
		return &abcitypes.ExecTxResult{
			Code: types.InvalidTransactionError,
		}
	}

	return &abcitypes.ExecTxResult{
		Code:      types.NoError,
		Events:    events,
		GasUsed:   0,
		GasWanted: 0,
		Data:      receiptBytes,
	}
}

func ParseRegisterAccountTransaction(tx t.Transaction) (RegisterAccountTransaction, error) {
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
