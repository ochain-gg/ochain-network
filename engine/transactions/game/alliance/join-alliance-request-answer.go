package alliance_transactions

import (
	abcitypes "github.com/cometbft/cometbft/abci/types"
	"github.com/fxamacker/cbor/v2"

	t "github.com/ochain-gg/ochain-network/engine/transactions"
	"github.com/ochain-gg/ochain-network/types"
)

type Answer struct {
	RequestId string `cbor:"1,keyasint"`
	Accept    bool   `cbor:"2,keyasint"`
}

type AnswerJoinAllianceRequestTransactionData struct {
	Answers []Answer `cbor:"1,keyasint"`
}

type AnswerJoinAllianceRequestTransaction struct {
	Type      t.TransactionType                        `cbor:"1,keyasint"`
	From      string                                   `cbor:"2,keyasint"`
	Nonce     uint64                                   `cbor:"3,keyasint"`
	Data      AnswerJoinAllianceRequestTransactionData `cbor:"4,keyasint"`
	Signature []byte                                   `cbor:"5,keyasint"`
}

func (tx *AnswerJoinAllianceRequestTransaction) Transaction() (t.Transaction, error) {
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

func (tx *AnswerJoinAllianceRequestTransaction) Check(ctx t.TransactionContext) *abcitypes.CheckTxResponse {
	// Check if global account exists
	now := uint64(ctx.Date.Unix())
	_, err := ctx.Db.GlobalsAccounts.GetAt(tx.From, now)
	if err != nil {
		return &abcitypes.CheckTxResponse{
			Code: types.InvalidTransactionError,
		}
	}

	for _, answer := range tx.Data.Answers {
		req, err := ctx.Db.Alliance.GetJoinRequestAt(answer.RequestId, now)
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

		alliance, err := ctx.Db.Alliance.GetAt(req.AllianceId, now)
		if err != nil {
			return &abcitypes.CheckTxResponse{
				Code: types.InvalidTransactionError,
			}
		}

		if alliance.Deleted {
			return &abcitypes.CheckTxResponse{
				Code: types.InvalidTransactionError,
			}
		}

		isMember, role := alliance.IsMember(tx.From)
		if !isMember || role == types.OChainAllianceSimpleMember {
			return &abcitypes.CheckTxResponse{
				Code: types.InvalidTransactionError,
			}
		}

		account, err := ctx.Db.UniverseAccounts.GetAt(alliance.UniverseId, req.From, uint64(ctx.Date.Unix()))
		if err != nil {
			return &abcitypes.CheckTxResponse{
				Code: types.InvalidTransactionError,
			}
		}

		if len(account.AllianceMemberOf) > 0 {
			return &abcitypes.CheckTxResponse{
				Code: types.InvalidTransactionError,
			}
		}
	}

	return &abcitypes.CheckTxResponse{
		Code: types.NoError,
	}
}

func (tx *AnswerJoinAllianceRequestTransaction) Execute(ctx t.TransactionContext) *abcitypes.ExecTxResult {
	response := tx.Check(ctx)
	if response.Code != types.NoError {
		return &abcitypes.ExecTxResult{
			Code: response.Code,
		}
	}

	events := []abcitypes.Event{}

	for _, answer := range tx.Data.Answers {
		req, err := ctx.Db.Alliance.GetJoinRequestAt(answer.RequestId, uint64(ctx.Date.Unix()))
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

		account, err := ctx.Db.UniverseAccounts.GetAt(alliance.UniverseId, req.From, uint64(ctx.Date.Unix()))
		if err != nil {
			return &abcitypes.ExecTxResult{
				Code: types.InvalidTransactionError,
			}
		}

		req.Accepted = answer.Accept
		req.AnsweredAt = ctx.Date.Unix()
		req.AnsweredBy = tx.From

		if answer.Accept {
			alliance.AddMember(req.From, types.OChainAllianceSimpleMember)
			account.AllianceMemberOf = req.Id

			ctx.Db.Alliance.Update(alliance)
			ctx.Db.UniverseAccounts.Update(account)

			events = append(events, abcitypes.Event{
				Type: "AllianceJoinRequestedAccepted",
				Attributes: []abcitypes.EventAttribute{
					{Key: "account", Value: req.From, Index: true},
					{Key: "request", Value: req.Id, Index: true},
					{Key: "universe", Value: alliance.UniverseId, Index: true},
					{Key: "alliance", Value: alliance.Id, Index: true},
					{Key: "by", Value: tx.From, Index: true},
				},
			})
		} else {
			events = append(events, abcitypes.Event{
				Type: "AllianceJoinRequestedRejected",
				Attributes: []abcitypes.EventAttribute{
					{Key: "account", Value: req.From, Index: true},
					{Key: "request", Value: req.Id, Index: true},
					{Key: "universe", Value: alliance.UniverseId, Index: true},
					{Key: "alliance", Value: alliance.Id, Index: true},
					{Key: "by", Value: tx.From, Index: true},
				},
			})
		}

		ctx.Db.Alliance.UpdateJoinRequest(req)

	}

	return &abcitypes.ExecTxResult{
		Code:      types.NoError,
		Events:    events,
		GasUsed:   100,
		GasWanted: 100,
		Data:      []byte{},
	}
}

func ParseAnswerJoinAllianceRequestTransaction(tx t.Transaction) (AnswerJoinAllianceRequestTransaction, error) {
	var txData AnswerJoinAllianceRequestTransactionData
	err := cbor.Unmarshal(tx.Data, &txData)

	if err != nil {
		return AnswerJoinAllianceRequestTransaction{}, err
	}

	return AnswerJoinAllianceRequestTransaction{
		Type:      tx.Type,
		From:      tx.From,
		Nonce:     tx.Nonce,
		Data:      txData,
		Signature: tx.Signature,
	}, nil
}
