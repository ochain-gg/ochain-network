package alliance_transactions

import (
	abcitypes "github.com/cometbft/cometbft/abci/types"
	"github.com/fxamacker/cbor/v2"

	t "github.com/ochain-gg/ochain-network/engine/transactions"
	"github.com/ochain-gg/ochain-network/types"
)

type JoinAllianceRequestTransactionData struct {
	UniverseId string `cbor:"1,keyasint"`
	AllianceId string `cbor:"2,keyasint"`
}

type JoinAllianceRequestTransaction struct {
	Type      t.TransactionType                  `cbor:"1,keyasint"`
	From      string                             `cbor:"2,keyasint"`
	Nonce     uint64                             `cbor:"3,keyasint"`
	Data      JoinAllianceRequestTransactionData `cbor:"4,keyasint"`
	Signature []byte                             `cbor:"5,keyasint"`
}

func (tx *JoinAllianceRequestTransaction) Transaction() (t.Transaction, error) {
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

func (tx *JoinAllianceRequestTransaction) Check(ctx t.TransactionContext) *abcitypes.CheckTxResponse {
	// Check if global account exists
	_, err := ctx.Db.GlobalsAccounts.GetAt(tx.From, uint64(ctx.Date.Unix()))
	if err != nil {
		return &abcitypes.CheckTxResponse{
			Code: types.InvalidTransactionError,
		}
	}

	// Check if universe account exists
	account, err := ctx.Db.UniverseAccounts.GetAt(tx.Data.UniverseId, tx.From, uint64(ctx.Date.Unix()))
	if err != nil {
		return &abcitypes.CheckTxResponse{
			Code: types.InvalidTransactionError,
		}
	}

	// Check if player is already in an alliance
	if len(account.AllianceMemberOf) > 0 {
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

	// Check if player already has a pending request
	hasRequest, err := ctx.Db.Alliance.HasPendingRequestAt(tx.From, uint64(ctx.Date.Unix()))
	if err != nil || hasRequest {
		return &abcitypes.CheckTxResponse{
			Code: types.InvalidTransactionError,
		}
	}

	if alliance.GetMaxAllianceSize() >= uint64(len(alliance.Members)) {
		return &abcitypes.CheckTxResponse{
			Code: types.InvalidTransactionError,
		}
	}

	return &abcitypes.CheckTxResponse{
		Code: types.NoError,
	}
}

func (tx *JoinAllianceRequestTransaction) Execute(ctx t.TransactionContext) *abcitypes.ExecTxResult {
	response := tx.Check(ctx)
	if response.Code != types.NoError {
		return &abcitypes.ExecTxResult{
			Code: response.Code,
		}
	}

	// Get global account
	globalAccount, err := ctx.Db.GlobalsAccounts.GetAt(tx.From, uint64(ctx.Date.Unix()))
	if err != nil {
		return &abcitypes.ExecTxResult{
			Code: types.InvalidTransactionError,
		}
	}

	// Get alliance
	alliance, err := ctx.Db.Alliance.GetAt(tx.Data.AllianceId, uint64(ctx.Date.Unix()))
	if err != nil {
		return &abcitypes.ExecTxResult{
			Code: types.InvalidTransactionError,
		}
	}

	// Create join request
	request := types.OChainAllianceJoinRequest{
		From:       tx.From,
		AllianceId: tx.Data.AllianceId,

		Accepted:    false,
		Canceled:    false,
		RequestedAt: ctx.Date.Unix(),
		AnsweredAt:  0,
		AnsweredBy:  "",
	}

	// Save join request
	err = ctx.Db.Alliance.InsertJoinRequest(request)
	if err != nil {
		return &abcitypes.ExecTxResult{
			Code: types.InvalidTransactionError,
		}
	}

	// Update global account
	err = ctx.Db.GlobalsAccounts.Update(globalAccount)
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
				{Key: "universe", Value: tx.Data.UniverseId, Index: true},
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

func ParseJoinAllianceRequestTransaction(tx t.Transaction) (JoinAllianceRequestTransaction, error) {
	var txData JoinAllianceRequestTransactionData
	err := cbor.Unmarshal(tx.Data, &txData)

	if err != nil {
		return JoinAllianceRequestTransaction{}, err
	}

	return JoinAllianceRequestTransaction{
		Type:      tx.Type,
		From:      tx.From,
		Nonce:     tx.Nonce,
		Data:      txData,
		Signature: tx.Signature,
	}, nil
}
