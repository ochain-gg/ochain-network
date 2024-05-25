package transactions

import (
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/dgraph-io/badger/v4"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/signer/core/apitypes"
	"github.com/fxamacker/cbor/v2"
	"github.com/ochain-gg/ochain-network/config"
	"github.com/ochain-gg/ochain-network/database"
	"github.com/ochain-gg/ochain-network/types"
)

type TransactionType uint16

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
	State  *types.OChainState
	Date   time.Time
}

type Transaction struct {
	Type      TransactionType `cbor:"1,keyasint"`
	From      string          `cbor:"2,keyasint,omitempty"`
	Nonce     uint64          `cbor:"3,keyasint,omitempty"`
	Data      []byte          `cbor:"4,keyasint"`
	Signature []byte          `cbor:"5,keyasint,omitempty"`
}

func (tx *Transaction) UniqueID() ([]byte, error) {

	payload, err := cbor.Marshal(tx)
	if err != nil {
		return []byte(""), err
	}

	txhash := crypto.Keccak256Hash(payload)

	if uint64(tx.Type) > MaxTransactionType {
		return []byte(""), errors.New("bad transaction type")
	}

	return txhash.Bytes(), nil
}

func (tx *Transaction) IsValid() error {

	if uint64(tx.Type) > MaxTransactionType {
		return errors.New("bad transaction type")
	}

	switch tx.Type {
	case OChainPortalInteraction:
		_, err := ParseNewOChainPortalInteraction(*tx)
		return err

	case ExecutePendingUpdate:
		return nil

	case RegisterAccount:
		_, err := ParseRegisterAccountTransaction(*tx)
		return err
	}

	return errors.New("unknown tx type")
}

func (tx *Transaction) VerifySignature() error {
	signer, err := tx.RecoverSignerAddress()
	if err != nil {
		return err
	}

	from := common.HexToAddress(tx.From)

	if signer.Hex() == from.Hex() {
		return nil
	} else {
		return errors.New("signer and from don't match")
	}
}

func (tx *Transaction) Sign(key []byte) error {
	sighash, err := tx.GetTypedDataHash()
	if err != nil {
		return err
	}

	signer := crypto.ToECDSAUnsafe(key)
	signature, err := crypto.Sign(sighash, signer)
	if err != nil {
		return err
	}

	tx.Signature = signature
	return nil
}

func (tx *Transaction) Bytes() ([]byte, error) {
	txByte, err := cbor.Marshal(tx)
	if err != nil {
		return []byte(""), err
	}

	return txByte, nil
}

func (tx *Transaction) GetTypedDataHash() ([]byte, error) {

	typedData, err := tx.GetTypedData()
	if err != nil {
		return []byte(""), err
	}

	sighash := crypto.Keccak256(typedData)

	return sighash, nil
}

func (tx *Transaction) GetTypedData() ([]byte, error) {
	typedData := apitypes.TypedData{
		Types: apitypes.Types{
			"EIP712Domain": []apitypes.Type{
				{Name: "name", Type: "string"},
				{Name: "version", Type: "string"},
				{Name: "chainId", Type: "uint256"},
				{Name: "verifyingContract", Type: "address"},
			},
			"Transaction": []apitypes.Type{
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
			"type":  strconv.FormatUint(uint64(tx.Type), 10),
			"from":  tx.From,
			"nonce": strconv.FormatUint(tx.Nonce, 10),
			"data":  string(tx.Data),
		},
	}
	// EIP-712 typed data marshalling
	domainSeparator, err := typedData.HashStruct("EIP712Domain", typedData.Domain.Map())
	if err != nil {
		return []byte(""), fmt.Errorf("eip712domain hash struct: %w", err)
	}
	typedDataHash, err := typedData.HashStruct(typedData.PrimaryType, typedData.Message)
	if err != nil {
		return []byte(""), fmt.Errorf("primary type hash struct: %w", err)
	}

	return []byte(fmt.Sprintf("\x19\x01%s%s", string(domainSeparator), string(typedDataHash))), nil
}

func (tx *Transaction) RecoverSignerAddress() (common.Address, error) {

	typedData, err := tx.GetTypedData()
	if err != nil {
		log.Println("GetTypedDataHash: " + err.Error())
		return common.Address{}, fmt.Errorf("GetTypedDataHash: %w", err)
	}

	sighash := crypto.Keccak256(typedData)

	// update the recovery id
	// https://github.com/ethereum/go-ethereum/blob/55599ee95d4151a2502465e0afc7c47bd1acba77/internal/ethapi/api.go#L442
	signature := tx.Signature
	signature[64] -= 27

	log.Println("typedData: " + hex.EncodeToString(sighash))
	log.Println("signature: " + hex.EncodeToString(signature))

	sigPubkey, err := crypto.SigToPub(sighash, signature)
	if err != nil {
		log.Println("SigToPub: " + err.Error())
		return common.Address{}, err
	}

	address := crypto.PubkeyToAddress(*sigPubkey)

	return address, nil
}

func ParseTransaction(data []byte) (Transaction, error) {
	var tx Transaction
	err := cbor.Unmarshal(data, &tx)

	if err != nil {
		return Transaction{}, err
	}

	return tx, nil
}
