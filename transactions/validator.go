package transactions

import "encoding/json"

const (
	ValidatorKeysPrefix = "validator-"
)

type NewValidatorTransactionData struct {
	RemoteTransactionHash []byte `json:"remoteTransactionHash"`
	PublicKey             []byte `json:"publicKey"`
}

type RemoveValidatorTransactionData struct {
	RemoteTransactionHash []byte `json:"remoteTransactionHash"`
	PublicKey             []byte `json:"publicKey"`
}

type NewValidatorTransaction struct {
	Type         TransactionType             `json:"type"`
	Data         []byte                      `json:"data"`
	FormatedData NewValidatorTransactionData `json:"formatedData"`
}

func (tx NewValidatorTransaction) GetChanges() (keys [][]byte, values [][]byte) {
	keys = append(keys, []byte(ValidatorKeysPrefix+string(tx.FormatedData.PublicKey)))
	values = append(keys, tx.FormatedData.RemoteTransactionHash)

	return keys, values
}

type NewValidatorEventData NewValidatorTransactionData

type RemoveValidatorTransaction struct {
	Type         TransactionType                `json:"type"`
	Data         []byte                         `json:"data"`
	FormatedData RemoveValidatorTransactionData `json:"formatedData"`
}

type RemoveValidatorEventData RemoveValidatorTransactionData

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
