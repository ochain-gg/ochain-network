package alliance_transactions

import (
	"fmt"
	"math"

	abcitypes "github.com/cometbft/cometbft/abci/types"
	"github.com/fxamacker/cbor/v2"

	t "github.com/ochain-gg/ochain-network/engine/transactions"
	"github.com/ochain-gg/ochain-network/types"
)

type UpgradeAllianceTransactionData struct {
	AllianceId         string `cbor:"1,keyasint"`
	PlanetCoordinateId string `cbor:"2,keyasint"`
}

type UpgradeAllianceTransaction struct {
	Type      t.TransactionType              `cbor:"1,keyasint"`
	From      string                         `cbor:"2,keyasint"`
	Nonce     uint64                         `cbor:"3,keyasint"`
	Data      UpgradeAllianceTransactionData `cbor:"4,keyasint"`
	Signature []byte                         `cbor:"5,keyasint"`
}

func (tx *UpgradeAllianceTransaction) Transaction() (t.Transaction, error) {
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

func (tx *UpgradeAllianceTransaction) Check(ctx t.TransactionContext) *abcitypes.CheckTxResponse {
	_, err := ctx.Db.GlobalsAccounts.GetAt(tx.From, uint64(ctx.Date.Unix()))
	if err != nil {
		return &abcitypes.CheckTxResponse{
			Code: types.InvalidTransactionError,
		}
	}

	account, err := ctx.Db.UniverseAccounts.GetAt(tx.From, tx.From, uint64(ctx.Date.Unix()))
	if err != nil {
		return &abcitypes.CheckTxResponse{
			Code: types.InvalidTransactionError,
		}
	}

	if account.AllianceMemberOf != tx.Data.AllianceId {
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

	universe, err := ctx.Db.Universes.GetAt(alliance.UniverseId, uint64(ctx.Date.Unix()))
	if err != nil {
		return &abcitypes.CheckTxResponse{
			Code: types.InvalidTransactionError,
		}
	}

	planet, err := ctx.Db.Planets.GetAt(alliance.UniverseId, tx.Data.PlanetCoordinateId, uint64(ctx.Date.Unix()))
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
		OCT:       uint64(float64(5000*types.OneOCT) * math.Pow(2, float64(alliance.Level+1))),
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

func (tx *UpgradeAllianceTransaction) Execute(ctx t.TransactionContext) *abcitypes.ExecTxResult {
	result := tx.Check(ctx)
	if result.Code != types.NoError {
		return &abcitypes.ExecTxResult{
			Code: result.GetCode(),
		}
	}

	account, err := ctx.Db.UniverseAccounts.GetAt(tx.From, tx.From, uint64(ctx.Date.Unix()))
	if err != nil {
		return &abcitypes.ExecTxResult{
			Code: types.InvalidTransactionError,
		}
	}

	// Check if alliance exists
	alliance, err := ctx.Db.Alliance.GetAt(tx.Data.AllianceId, uint64(ctx.Date.Unix()))
	if err != nil {
		return &abcitypes.ExecTxResult{
			Code: types.InvalidTransactionError,
		}
	}

	universe, err := ctx.Db.Universes.GetAt(alliance.UniverseId, uint64(ctx.Date.Unix()))
	if err != nil {
		return &abcitypes.ExecTxResult{
			Code: types.InvalidTransactionError,
		}
	}

	planet, err := ctx.Db.Planets.GetAt(alliance.UniverseId, tx.Data.PlanetCoordinateId, uint64(ctx.Date.Unix()))
	if err != nil {
		return &abcitypes.ExecTxResult{
			Code: types.InvalidTransactionError,
		}
	}

	planet.Update(universe.Speed, ctx.Date.Unix(), account)
	err = planet.Pay(types.OChainResources{
		OCT:       uint64(float64(5000*types.OneOCT) * math.Pow(2, float64(alliance.Level+1))),
		Metal:     0,
		Crystal:   0,
		Deuterium: 0,
	})

	alliance.Level += 1

	if err != nil {
		return &abcitypes.ExecTxResult{
			Code: types.InvalidTransactionError,
		}
	}

	err = ctx.Db.Alliance.Update(alliance)
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
				{Key: "universe", Value: alliance.UniverseId, Index: true},
				{Key: "planet", Value: tx.Data.PlanetCoordinateId, Index: true},
				{Key: "oct", Value: fmt.Sprint(planet.Resources.OCT)},
				{Key: "metal", Value: fmt.Sprint(planet.Resources.Metal)},
				{Key: "crystal", Value: fmt.Sprint(planet.Resources.Crystal)},
				{Key: "deuterium", Value: fmt.Sprint(planet.Resources.Deuterium)},
			},
		},
		{
			Type: "AllianceUpgraded",
			Attributes: []abcitypes.EventAttribute{
				{Key: "account", Value: tx.From, Index: true},
				{Key: "universe", Value: alliance.UniverseId, Index: true},
				{Key: "alliance", Value: alliance.Id},
				{Key: "level", Value: fmt.Sprint(alliance.Level)},
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

func ParseUpgradeAllianceTransaction(tx t.Transaction) (UpgradeAllianceTransaction, error) {
	var txData UpgradeAllianceTransactionData
	err := cbor.Unmarshal(tx.Data, &txData)

	if err != nil {
		return UpgradeAllianceTransaction{}, err
	}

	return UpgradeAllianceTransaction{
		Type:      tx.Type,
		From:      tx.From,
		Nonce:     tx.Nonce,
		Data:      txData,
		Signature: tx.Signature,
	}, nil
}
