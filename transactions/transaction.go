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
	MaxTransactionType uint64 = 31

	//System transactions (Unauthenticated)
	NewValidator          TransactionType = 0
	RemoveValidator       TransactionType = 1
	SlashValidator        TransactionType = 1
	NewEpoch              TransactionType = 2
	OChainTokenDeposit    TransactionType = 2
	OChainCreditDeposit   TransactionType = 4
	ExecutePendingUpgrade TransactionType = 5 // handle BuildingUpgrade / TechnologyUpgrade / DefenseBuild / SpaceshipBuild / FleetMove

	//Authenticated transactions

	/**
	 * Global Account transactions
	 */
	RegisterAccount          TransactionType = 6
	ChangeAccountIAM         TransactionType = 7
	ExecuteBridgeTransaction TransactionType = 7
	OChainTokenWithdrawal    TransactionType = 3

	/**
	 * Governance transactions
	 */
	CreateGovernanceProposal TransactionType = 33
	VoteOnProposal           TransactionType = 34
	ExecuteProposal          TransactionType = 35

	/**
	 * Game transactions
	 */

	//Universe account
	RegisterUniverseAccount TransactionType = 8

	//Planet transactions
	MintPlanet             TransactionType = 9
	StartBuildingUpgrade   TransactionType = 10
	StartTechnologyUpgrade TransactionType = 11
	StartBuildDefenses     TransactionType = 12
	StartBuildSpaceships   TransactionType = 13

	//Fleet transactions
	FillCargo                    TransactionType = 14
	UnfillCargo                  TransactionType = 15
	SendFleetInOrbit             TransactionType = 16
	LandingFleetInOrbit          TransactionType = 17
	MergeFleets                  TransactionType = 18
	SplitFleet                   TransactionType = 19
	RecycleRemnant               TransactionType = 20
	IntergalacticPortalFleetMove TransactionType = 21
	ChangeFleetMode              TransactionType = 22
	AcceptFleetMode              TransactionType = 23
	MoveFleet                    TransactionType = 24
	CancelFleetMove              TransactionType = 25

	//Fight transactions
	Fight TransactionType = 26

	//Alliance transactions
	CreateAlliance  TransactionType = 27
	StakeOnAlliance TransactionType = 28

	//Market transations
	SwapResources TransactionType = 29

	//Global <-> universes transations
	UniverseOChainTokenDeposit  TransactionType = 30
	UniverseOChainTokenWithdraw TransactionType = 31
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

type TransactionReceipt struct {
	GasCost uint64 `cbor:"1,keyasint"`
}

func (receipt *TransactionReceipt) Bytes() []byte {
	receiptBytes, err := cbor.Marshal(receipt)
	if err != nil {
		return []byte("")
	}
	return receiptBytes
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
	//System Txs
	case NewValidator:
		_, err := ParseOChainBridgeNewValidatorTransaction(*tx)
		return err
	case RemoveValidator:
		_, err := ParseOChainBridgeRemoveValidatorTransaction(*tx)
		return err

	case ExecutePendingUpgrade:
		return nil

	//Authenticated Txs
	case RegisterAccount:
		_, err := ParseRegisterAccountTransaction(*tx)
		return err
	case RegisterUniverseAccount:
		_, err := ParseRegisterUniverseAccountTransaction(*tx)
		return err
	}

	return errors.New("unknown tx type")
}

func (tx *Transaction) GetSigner() (string, error) {
	signer, err := tx.RecoverSignerAddress()
	if err != nil {
		return "", err
	}

	return signer.Hex(), nil
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
