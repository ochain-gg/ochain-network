package alliance_transactions

import (
	abcitypes "github.com/cometbft/cometbft/abci/types"
	"github.com/fxamacker/cbor/v2"

	t "github.com/ochain-gg/ochain-network/engine/transactions"
	"github.com/ochain-gg/ochain-network/types"
)

type EditAllianceTransactionData struct {
	AllianceId  string `cbor:"2,keyasint"`
	Name        string `cbor:"3,keyasint"`
	Description string `cbor:"4,keyasint"`
	Tag         string `cbor:"5,keyasint"`
}

type EditAllianceTransaction struct {
	Type      t.TransactionType           `cbor:"1,keyasint"`
	From      string                      `cbor:"2,keyasint"`
	Nonce     uint64                      `cbor:"3,keyasint"`
	Data      EditAllianceTransactionData `cbor:"4,keyasint"`
	Signature []byte                      `cbor:"5,keyasint"`
}

func (tx *EditAllianceTransaction) Transaction() (t.Transaction, error) {
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

func (tx *EditAllianceTransaction) Check(ctx t.TransactionContext) *abcitypes.CheckTxResponse {
	// Check if global account exists
	_, err := ctx.Db.GlobalsAccounts.GetAt(tx.From, uint64(ctx.Date.Unix()))
	if err != nil {
		return &abcitypes.CheckTxResponse{
			Code: types.InvalidTransactionError,
		}
	}

	// Check if alliance exists
	alliance, err := ctx.Db.Alliance.GetAt(tx.Data.AllianceId, uint64(ctx.Date.Unix()))
	if err != nil {
		return &abcitypes.CheckTxResponse{
			Code: types.InvalidTransactionError,
		}
	}

	_, role := alliance.IsMember(tx.From)
	if role != types.OChainAllianceLeaderMember {
		return &abcitypes.CheckTxResponse{
			Code: types.InvalidTransactionError,
		}
	}

	return &abcitypes.CheckTxResponse{
		Code: types.NoError,
	}
}

func (tx *EditAllianceTransaction) Execute(ctx t.TransactionContext) *abcitypes.ExecTxResult {
	response := tx.Check(ctx)
	if response.Code != types.NoError {
		return &abcitypes.ExecTxResult{
			Code: response.Code,
		}
	}

	// Get alliance
	alliance, err := ctx.Db.Alliance.GetAt(tx.Data.AllianceId, uint64(ctx.Date.Unix()))
	if err != nil {
		return &abcitypes.ExecTxResult{
			Code: types.InvalidTransactionError,
		}
	}

	alliance.Name = tx.Data.Name
	alliance.Tag = tx.Data.Tag
	alliance.Description = tx.Data.Description

	err = ctx.Db.Alliance.Update(alliance)
	if err != nil {
		return &abcitypes.ExecTxResult{
			Code: types.InvalidTransactionError,
		}
	}

	// Create events
	events := []abcitypes.Event{
		{
			Type: "AllianceJoinRequested",
			Attributes: []abcitypes.EventAttribute{
				{Key: "account", Value: tx.From, Index: true},
				{Key: "universe", Value: alliance.UniverseId, Index: true},
				{Key: "name", Value: tx.Data.Name},
				{Key: "description", Value: tx.Data.Description},
				{Key: "tag", Value: tx.Data.Tag},
				{Key: "allianceId", Value: alliance.Id},
			},
		},
	}

	return &abcitypes.ExecTxResult{
		Code:      types.NoError,
		Events:    events,
		GasUsed:   100,
		GasWanted: 100,
		Data:      []byte{},
	}
}

func ParseEditAllianceTransaction(tx t.Transaction) (EditAllianceTransaction, error) {
	var txData EditAllianceTransactionData
	err := cbor.Unmarshal(tx.Data, &txData)

	if err != nil {
		return EditAllianceTransaction{}, err
	}

	return EditAllianceTransaction{
		Type:      tx.Type,
		From:      tx.From,
		Nonce:     tx.Nonce,
		Data:      txData,
		Signature: tx.Signature,
	}, nil
}
