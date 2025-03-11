package alliance_transactions

import (
	"fmt"

	abcitypes "github.com/cometbft/cometbft/abci/types"
	"github.com/fxamacker/cbor/v2"

	t "github.com/ochain-gg/ochain-network/engine/transactions"
	"github.com/ochain-gg/ochain-network/types"
)

type AssignAllianceRoleTransactionData struct {
	AllianceId string                         `cbor:"1,keyasint"`
	Member     string                         `cbor:"2,keyasint"`
	Role       types.OChainAllianceMemberType `cbor:"3,keyasint"`
}

type AssignAllianceRoleTransaction struct {
	Type      t.TransactionType                 `cbor:"1,keyasint"`
	From      string                            `cbor:"2,keyasint"`
	Nonce     uint64                            `cbor:"3,keyasint"`
	Data      AssignAllianceRoleTransactionData `cbor:"4,keyasint"`
	Signature []byte                            `cbor:"5,keyasint"`
}

func (tx *AssignAllianceRoleTransaction) Transaction() (t.Transaction, error) {
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

func (tx *AssignAllianceRoleTransaction) Check(ctx t.TransactionContext) *abcitypes.CheckTxResponse {
	// Check if alliance exists
	alliance, err := ctx.Db.Alliance.GetAt(tx.Data.AllianceId, uint64(ctx.Date.Unix()))
	if err != nil {
		return &abcitypes.CheckTxResponse{
			Code: types.InvalidTransactionError,
		}
	}

	_, err = ctx.Db.UniverseAccounts.GetAt(alliance.UniverseId, tx.From, uint64(ctx.Date.Unix()))
	if err != nil {
		return &abcitypes.CheckTxResponse{
			Code: types.InvalidTransactionError,
		}
	}

	_, err = ctx.Db.UniverseAccounts.GetAt(alliance.UniverseId, tx.Data.Member, uint64(ctx.Date.Unix()))
	if err != nil {
		return &abcitypes.CheckTxResponse{
			Code: types.InvalidTransactionError,
		}
	}

	_, fromRole := alliance.IsMember(tx.From)
	if fromRole != types.OChainAllianceLeaderMember && fromRole != types.OChainAllianceOfficerMember {
		return &abcitypes.CheckTxResponse{
			Code: types.InvalidTransactionError,
		}
	}

	_, memberRole := alliance.IsMember(tx.From)
	if memberRole != types.OChainAllianceLeaderMember {
		return &abcitypes.CheckTxResponse{
			Code: types.InvalidTransactionError,
		}
	}

	if memberRole == types.OChainAllianceNotMember {
		return &abcitypes.CheckTxResponse{
			Code: types.InvalidTransactionError,
		}
	}

	switch fromRole {
	case types.OChainAllianceOfficerMember:
		if memberRole == types.OChainAllianceOfficerMember || memberRole == types.OChainAllianceLeaderMember {
			return &abcitypes.CheckTxResponse{
				Code: types.InvalidTransactionError,
			}
		}

		if tx.Data.Role == types.OChainAllianceLeaderMember {
			return &abcitypes.CheckTxResponse{
				Code: types.InvalidTransactionError,
			}
		}
	}

	return &abcitypes.CheckTxResponse{
		Code: types.NoError,
	}
}

func (tx *AssignAllianceRoleTransaction) Execute(ctx t.TransactionContext) *abcitypes.ExecTxResult {
	result := tx.Check(ctx)
	if result.Code != types.NoError {
		return &abcitypes.ExecTxResult{
			Code: result.GetCode(),
		}
	}

	alliance, err := ctx.Db.Alliance.GetAt(tx.Data.AllianceId, uint64(ctx.Date.Unix()))
	if err != nil {
		return &abcitypes.ExecTxResult{
			Code: types.InvalidTransactionError,
		}
	}

	account, err := ctx.Db.UniverseAccounts.GetAt(alliance.UniverseId, tx.From, uint64(ctx.Date.Unix()))
	if err != nil {
		return &abcitypes.ExecTxResult{
			Code: types.InvalidTransactionError,
		}
	}

	if tx.Data.Role == types.OChainAllianceNotMember {
		alliance.RemoveMember(tx.Data.Member)
		account.AllianceMemberOf = ""

		err = ctx.Db.UniverseAccounts.Update(account)
		if err != nil {
			return &abcitypes.ExecTxResult{
				Code: types.InvalidTransactionError,
			}
		}
	} else {
		alliance.ChangeRole(tx.Data.Member, tx.Data.Role)
	}

	err = ctx.Db.Alliance.Update(alliance)
	if err != nil {
		return &abcitypes.ExecTxResult{
			Code: types.InvalidTransactionError,
		}
	}

	events := []abcitypes.Event{
		{
			Type: "AllianceRoleAssigned",
			Attributes: []abcitypes.EventAttribute{
				{Key: "account", Value: tx.From, Index: true},
				{Key: "universe", Value: alliance.UniverseId, Index: true},
				{Key: "alliance", Value: alliance.Id},
				{Key: "role", Value: alliance.Id},
				{Key: "member", Value: fmt.Sprint(tx.Data.Role)},
			},
		},
	}

	return &abcitypes.ExecTxResult{
		Code:      types.NoError,
		Events:    events,
		GasUsed:   100,
		GasWanted: 100,
		Data:      []byte(""),
	}
}

func ParseAssignAllianceRoleTransaction(tx t.Transaction) (AssignAllianceRoleTransaction, error) {
	var txData AssignAllianceRoleTransactionData
	err := cbor.Unmarshal(tx.Data, &txData)

	if err != nil {
		return AssignAllianceRoleTransaction{}, err
	}

	return AssignAllianceRoleTransaction{
		Type:      tx.Type,
		From:      tx.From,
		Nonce:     tx.Nonce,
		Data:      txData,
		Signature: tx.Signature,
	}, nil
}
