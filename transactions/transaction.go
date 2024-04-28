package transactions

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/dgraph-io/badger/v4"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/signer/core/apitypes"
	"github.com/ochain.gg/ochain-network-validator/config"
	"github.com/ochain.gg/ochain-network-validator/database"
)

type TransactionType uint64

const (
	MaxTransactionType uint64 = 28
	//Unauthenticated transactions
	OChainPortalInteraction TransactionType = 0 // handle NewValidator / RemoveValidator / OChainTokenDeposit / OChainTokenWithdrawal / OChainBonusSubscription
	ExecutePendingUpdate    TransactionType = 1 // handle BuildingUpgrade / TechnologyUpgrade / DefenseBuild / SpaceshipBuild / FleetMove

	//Authenticated transactions
	//Global Account transactions
	RegisterAccount  TransactionType = 5
	ChangeAccountIAM TransactionType = 6

	//Universe account
	RegisterUniverseAccount TransactionType = 7

	//Planet transactions
	MintPlanet             TransactionType = 8
	StartBuildingUpgrade   TransactionType = 9
	StartTechnologyUpgrade TransactionType = 10
	StartBuildDefenses     TransactionType = 11
	StartBuildSpaceships   TransactionType = 12

	//Fleet transactions
	FillCargo                    TransactionType = 13
	UnfillCargo                  TransactionType = 14
	SendFleetInOrbit             TransactionType = 15
	LandingFleetInOrbit          TransactionType = 16
	MergeFleets                  TransactionType = 17
	SplitFleet                   TransactionType = 18
	RecycleRemnant               TransactionType = 19
	IntergalacticPortalFleetMove TransactionType = 20
	ChangeFleetMode              TransactionType = 21
	AcceptFleetMode              TransactionType = 22
	MoveFleet                    TransactionType = 23
	CancelFleetMove              TransactionType = 24

	//Fight transactions
	Fight TransactionType = 25

	//Alliance transactions
	CreateAlliance  TransactionType = 26
	StakeOnAlliance TransactionType = 27

	//Market transations
	SwapResources TransactionType = 28
)

type TransactionContext struct {
	Config config.OChainConfig
	Db     *database.OChainDatabase
	Txn    *badger.Txn
	Date   time.Time
}

type Transaction struct {
	Hash string          `json:"hash"`
	Type TransactionType `json:"type"`
	Data string          `json:"data"`
}

type UnauthenticatedTransaction struct {
	Hash string          `json:"hash"`
	Type TransactionType `json:"type"`
	Data string          `json:"data"`
}

type AuthenticatedTransaction struct {
	Hash      string          `json:"hash"`
	Type      TransactionType `json:"type"`
	From      string          `json:"from"`
	Nonce     uint64          `json:"nonce"`
	Data      string          `json:"data"`
	Signature string          `json:"signature"`
}

func (tx *Transaction) IsValid() error {

	txhash := crypto.Keccak256Hash([]byte(strconv.Itoa(int(tx.Type)) + ":" + tx.Data))
	if tx.Hash != txhash.Hex() {
		return errors.New("bad hash")
	}

	if uint64(tx.Type) > MaxTransactionType {
		return errors.New("bad transaction type")
	}

	return nil
}

func (tx *UnauthenticatedTransaction) IsValid() error {

	txhash := crypto.Keccak256Hash([]byte(strconv.Itoa(int(tx.Type)) + ":" + tx.Data))
	if tx.Hash != txhash.Hex() {
		return errors.New("bad hash")
	}

	switch tx.Type {
	case OChainPortalInteraction:
		_, err := ParseNewOChainPortalInteraction(Transaction(*tx))
		return err

	case ExecutePendingUpdate:
		return nil
	}
	return errors.New("unknown tx type")
}

func (tx *AuthenticatedTransaction) RecoverSignerAddress() (common.Address, error) {

	txhash := crypto.Keccak256Hash([]byte(strconv.Itoa(int(tx.Type)) + ":" + tx.Data))
	if tx.Hash != txhash.Hex() {
		return common.Address{}, errors.New("bad hash")
	}

	signature, err := hexutil.Decode(tx.Signature)
	if err != nil {
		return common.Address{}, fmt.Errorf("decode signature: %w", err)
	}

	typedData := apitypes.TypedData{
		Types: apitypes.Types{
			"EIP712Domain": []apitypes.Type{
				{Name: "name", Type: "string"},
				{Name: "version", Type: "string"},
				{Name: "chainId", Type: "uint256"},
				{Name: "verifyingContract", Type: "address"},
			},
			"Transaction": []apitypes.Type{
				{Name: "hash", Type: "string"},
				{Name: "type", Type: "uint16"},
				{Name: "from", Type: "address"},
				{Name: "nonce", Type: "uint256"},
				{Name: "data", Type: "string"},
			},
		},
		PrimaryType: "Transaction",
		Domain: apitypes.TypedDataDomain{
			Name:              "OChainNetwork",
			Version:           "1",
			ChainId:           math.NewHexOrDecimal256(20291),
			VerifyingContract: "0x0000000000000000000000000000000000000000",
		},
		Message: map[string]interface{}{
			"hash":  tx.Hash,
			"type":  tx.Type,
			"from":  tx.From,
			"nonce": tx.Nonce,
			"data":  tx.Data,
		},
	}
	// EIP-712 typed data marshalling
	domainSeparator, err := typedData.HashStruct("EIP712Domain", typedData.Domain.Map())
	if err != nil {
		return common.Address{}, fmt.Errorf("eip712domain hash struct: %w", err)
	}
	typedDataHash, err := typedData.HashStruct(typedData.PrimaryType, typedData.Message)
	if err != nil {
		return common.Address{}, fmt.Errorf("primary type hash struct: %w", err)
	}

	// add magic string prefix
	rawData := []byte(fmt.Sprintf("\x19\x01%s%s", string(domainSeparator), string(typedDataHash)))
	sighash := crypto.Keccak256(rawData)
	// fmt.Println("SIG HASH:", hexutil.Encode(sighash))

	// update the recovery id
	// https://github.com/ethereum/go-ethereum/blob/55599ee95d4151a2502465e0afc7c47bd1acba77/internal/ethapi/api.go#L442
	signature[64] -= 27

	// get the pubkey used to sign this signature
	sigPubkey, err := crypto.Ecrecover(sighash, signature)
	if err != nil {
		return common.Address{}, fmt.Errorf("ecrecover: %w", err)
	}
	// fmt.Println("SIG PUBKEY:", hexutil.Encode(sigPubkey))

	// get the address to confirm it's the same one in the auth token
	pubkey, err := crypto.UnmarshalPubkey(sigPubkey)
	if err != nil {
		return common.Address{}, err
	}
	address := crypto.PubkeyToAddress(*pubkey)
	// fmt.Println("ADDRESS:", address.Hex())

	// verify the signature (not sure if this is actually required after ecrecover)
	signatureNoRecoverID := signature[:len(signature)-1]
	verified := crypto.VerifySignature(sigPubkey, sighash, signatureNoRecoverID)
	if !verified {
		return common.Address{}, errors.New("verification failed")
	}

	return address, nil
	// fmt.Println("VERIFIED:", verified)
}

func ParseTransaction(data []byte) (Transaction, error) {

	var tx Transaction
	err := json.Unmarshal(data, &tx)

	if err != nil {
		return Transaction{}, err
	}

	return tx, nil
}

func ParseAuthenticatedTransaction(data []byte) (AuthenticatedTransaction, error) {

	var tx AuthenticatedTransaction
	err := json.Unmarshal(data, &tx)

	if err != nil {
		return AuthenticatedTransaction{}, err
	}

	return tx, nil
}
