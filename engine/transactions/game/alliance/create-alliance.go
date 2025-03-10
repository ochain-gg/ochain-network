package alliance_transactions

import (
	"encoding/hex"
	"fmt"

	abcitypes "github.com/cometbft/cometbft/abci/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/fxamacker/cbor/v2"

	t "github.com/ochain-gg/ochain-network/engine/transactions"
	"github.com/ochain-gg/ochain-network/types"
)

type CreateAllianceTransactionData struct {
	Name               string `cbor:"1,keyasint"`
	UniverseId         string `cbor:"2,keyasint"`
	PlanetCoordinateId string `cbor:"3,keyasint"`
}

type CreateAllianceTransaction struct {
	Type      t.TransactionType             `cbor:"1,keyasint"`
	From      string                        `cbor:"2,keyasint"`
	Nonce     uint64                        `cbor:"3,keyasint"`
	Data      CreateAllianceTransactionData `cbor:"4,keyasint"`
	Signature []byte                        `cbor:"5,keyasint"`
}

func (tx *CreateAllianceTransaction) Transaction() (t.Transaction, error) {
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

func (tx *CreateAllianceTransaction) Check(ctx t.TransactionContext) *abcitypes.CheckTxResponse {
	_, err := ctx.Db.GlobalsAccounts.GetAt(tx.From, uint64(ctx.Date.Unix()))
	if err != nil {
		return &abcitypes.CheckTxResponse{
			Code: types.InvalidTransactionError,
		}
	}

	account, err := ctx.Db.UniverseAccounts.GetAt(tx.Data.UniverseId, tx.From, uint64(ctx.Date.Unix()))
	if err != nil {
		return &abcitypes.CheckTxResponse{
			Code: types.InvalidTransactionError,
		}
	}

	universe, err := ctx.Db.Universes.GetAt(tx.Data.UniverseId, uint64(ctx.Date.Unix()))
	if err != nil {
		return &abcitypes.CheckTxResponse{
			Code: types.InvalidTransactionError,
		}
	}

	planet, err := ctx.Db.Planets.GetAt(tx.Data.UniverseId, tx.Data.PlanetCoordinateId, uint64(ctx.Date.Unix()))
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

	planet.Update(universe.Speed, int64(ctx.Date.Unix()), account)
	payable := planet.CanPay(types.OChainResources{
		OCT:       5000 * types.OneOCT,
		Metal:     0,
		Crystal:   0,
		Deuterium: 0,
	})

	if !payable {
		return &abcitypes.CheckTxResponse{
			Code: types.InvalidTransactionError,
		}
	}

	return &abcitypes.CheckTxResponse{
		Code: types.NoError,
	}
}

func (tx *CreateAllianceTransaction) Execute(ctx t.TransactionContext) *abcitypes.ExecTxResult {
	result := tx.Check(ctx)
	if result.Code != types.NoError {
		return &abcitypes.ExecTxResult{
			Code: result.GetCode(),
		}
	}

	account, err := ctx.Db.UniverseAccounts.GetAt(tx.Data.UniverseId, tx.From, uint64(ctx.Date.Unix()))
	if err != nil {
		return &abcitypes.ExecTxResult{
			Code: types.InvalidTransactionError,
		}
	}

	universe, err := ctx.Db.Universes.GetAt(tx.Data.UniverseId, uint64(ctx.Date.Unix()))
	if err != nil {
		return &abcitypes.ExecTxResult{
			Code: types.InvalidTransactionError,
		}
	}

	planet, err := ctx.Db.Planets.GetAt(tx.Data.UniverseId, tx.Data.PlanetCoordinateId, uint64(ctx.Date.Unix()))
	if err != nil {
		return &abcitypes.ExecTxResult{
			Code: types.InvalidTransactionError,
		}
	}

	planet.Update(universe.Speed, ctx.Date.Unix(), account)
	err = planet.Pay(types.OChainResources{
		OCT:       5000 * types.OneOCT,
		Metal:     0,
		Crystal:   0,
		Deuterium: 0,
	})

	if err != nil {
		return &abcitypes.ExecTxResult{
			Code: types.InvalidTransactionError,
		}
	}

	members := []types.OChainAllianceMember{{
		Role:    types.OChainAllianceLeaderMember,
		Address: account.Address,
	}}

	allianceId := crypto.Keccak256([]byte(tx.From + fmt.Sprint(tx.Nonce) + fmt.Sprint(ctx.Date.Unix())))
	alliance := types.OChainAlliance{
		Id:         hex.EncodeToString(allianceId),
		UniverseId: tx.Data.UniverseId,
		Name:       tx.Data.Name,
		Level:      1,
		Members:    members,
		Deleted:    false,
	}

	account.AllianceMemberOf = alliance.Id

	err = ctx.Db.Alliance.Insert(alliance)
	if err != nil {
		return &abcitypes.ExecTxResult{
			Code: types.InvalidTransactionError,
		}
	}

	err = ctx.Db.Planets.Update(planet.Universe, planet)
	if err != nil {
		return &abcitypes.ExecTxResult{
			Code: types.InvalidTransactionError,
		}
	}

	err = ctx.Db.UniverseAccounts.Update(account)
	if err != nil {
		return &abcitypes.ExecTxResult{
			Code: types.InvalidTransactionError,
		}
	}

	events := []abcitypes.Event{
		{
			Type: "PlanetResourcesUpdated",
			Attributes: []abcitypes.EventAttribute{
				{Key: "account", Value: tx.From, Index: true},
				{Key: "universe", Value: tx.Data.UniverseId, Index: true},
				{Key: "planet", Value: tx.Data.PlanetCoordinateId, Index: true},
				{Key: "oct", Value: fmt.Sprint(planet.Resources.OCT)},
				{Key: "metal", Value: fmt.Sprint(planet.Resources.Metal)},
				{Key: "crystal", Value: fmt.Sprint(planet.Resources.Crystal)},
				{Key: "deuterium", Value: fmt.Sprint(planet.Resources.Deuterium)},
			},
		},
		{
			Type: "AllianceCreated",
			Attributes: []abcitypes.EventAttribute{
				{Key: "account", Value: tx.From, Index: true},
				{Key: "universe", Value: tx.Data.UniverseId, Index: true},
				{Key: "allianceId", Value: alliance.Id},
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

func ParseCreateAllianceTransaction(tx t.Transaction) (CreateAllianceTransaction, error) {
	var txData CreateAllianceTransactionData
	err := cbor.Unmarshal(tx.Data, &txData)

	if err != nil {
		return CreateAllianceTransaction{}, err
	}

	return CreateAllianceTransaction{
		Type:      tx.Type,
		From:      tx.From,
		Nonce:     tx.Nonce,
		Data:      txData,
		Signature: tx.Signature,
	}, nil
}
