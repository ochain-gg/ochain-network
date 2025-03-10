package alliance_transactions

import (
	abcitypes "github.com/cometbft/cometbft/abci/types"
	"github.com/fxamacker/cbor/v2"

	t "github.com/ochain-gg/ochain-network/engine/transactions"
	"github.com/ochain-gg/ochain-network/types"
)

type CancelJoinAllianceRequestTransactionData struct {
	RequestId string `cbor:"1,keyasint"`
}

type CancelJoinAllianceRequestTransaction struct {
	Type      t.TransactionType                        `cbor:"1,keyasint"`
	From      string                                   `cbor:"2,keyasint"`
	Nonce     uint64                                   `cbor:"3,keyasint"`
	Data      CancelJoinAllianceRequestTransactionData `cbor:"4,keyasint"`
	Signature []byte                                   `cbor:"5,keyasint"`
}

func (tx *CancelJoinAllianceRequestTransaction) Transaction() (t.Transaction, error) {
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

func (tx *CancelJoinAllianceRequestTransaction) Check(ctx t.TransactionContext) *abcitypes.CheckTxResponse {
	// Check if global account exists
	now := uint64(ctx.Date.Unix())
	_, err := ctx.Db.GlobalsAccounts.GetAt(tx.From, now)
	if err != nil {
		return &abcitypes.CheckTxResponse{
			Code: types.InvalidTransactionError,
		}
	}

	req, err := ctx.Db.Alliance.GetJoinRequestAt(tx.Data.RequestId, now)
	if err != nil {
		return &abcitypes.CheckTxResponse{
			Code: types.InvalidTransactionError,
		}
	}

	if req.Canceled || req.AnsweredAt > 0 {
		return &abcitypes.CheckTxResponse{
			Code: types.InvalidTransactionError,
		}
	}

	if tx.From != req.From {
		return &abcitypes.CheckTxResponse{
			Code: types.InvalidTransactionError,
		}
	}

	return &abcitypes.CheckTxResponse{
		Code: types.NoError,
	}
}

func (tx *CancelJoinAllianceRequestTransaction) Execute(ctx t.TransactionContext) *abcitypes.ExecTxResult {
	response := tx.Check(ctx)
	if response.Code != types.NoError {
		return &abcitypes.ExecTxResult{
			Code: response.Code,
		}
	}

	req, err := ctx.Db.Alliance.GetJoinRequestAt(tx.Data.RequestId, uint64(ctx.Date.Unix()))
	if err != nil {
		return &abcitypes.ExecTxResult{
			Code: types.InvalidTransactionError,
		}
	}

	alliance, err := ctx.Db.Alliance.GetAt(req.AllianceId, uint64(ctx.Date.Unix()))
	if err != nil {
		return &abcitypes.ExecTxResult{
			Code: types.InvalidTransactionError,
		}
	}

	req.Canceled = true
	req.AnsweredAt = ctx.Date.Unix()
	ctx.Db.Alliance.UpdateJoinRequest(req)

	events := []abcitypes.Event{
		{
			Type: "AllianceJoinRequestedCanceled",
			Attributes: []abcitypes.EventAttribute{
				{Key: "account", Value: req.From, Index: true},
				{Key: "request", Value: req.Id, Index: true},
				{Key: "universe", Value: alliance.UniverseId, Index: true},
				{Key: "alliance", Value: alliance.Id, Index: true},
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

func ParseCancelJoinAllianceRequestTransaction(tx t.Transaction) (CancelJoinAllianceRequestTransaction, error) {
	var txData CancelJoinAllianceRequestTransactionData
	err := cbor.Unmarshal(tx.Data, &txData)

	if err != nil {
		return CancelJoinAllianceRequestTransaction{}, err
	}

	return CancelJoinAllianceRequestTransaction{
		Type:      tx.Type,
		From:      tx.From,
		Nonce:     tx.Nonce,
		Data:      txData,
		Signature: tx.Signature,
	}, nil
}
