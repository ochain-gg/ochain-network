package fleet_transactions

import (
	"fmt"

	abcitypes "github.com/cometbft/cometbft/abci/types"
	"github.com/fxamacker/cbor/v2"

	t "github.com/ochain-gg/ochain-network/engine/transactions"
	"github.com/ochain-gg/ochain-network/types"
)

type UnfeelFleetCargoTransactionData struct {
	UniverseId         string                `cbor:"2,keyasint"`
	PlanetCoordinateId string                `cbor:"3,keyasint"`
	FleetId            string                `cbor:"4,keyasint"`
	Cargo              types.OChainResources `cbor:"5,keyasint"`
}

type UnfeelFleetCargoTransaction struct {
	Type      t.TransactionType               `cbor:"1,keyasint"`
	From      string                          `cbor:"2,keyasint"`
	Nonce     uint64                          `cbor:"3,keyasint"`
	Data      UnfeelFleetCargoTransactionData `cbor:"4,keyasint"`
	Signature []byte                          `cbor:"5,keyasint"`
}

func (tx *UnfeelFleetCargoTransaction) Transaction() (t.Transaction, error) {
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

func (tx *UnfeelFleetCargoTransaction) Check(ctx t.TransactionContext) *abcitypes.CheckTxResponse {
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

	if planet.Owner != tx.From {
		return &abcitypes.CheckTxResponse{
			Code: types.InvalidTransactionError,
		}
	}

	planet.Update(universe.Speed, int64(ctx.Date.Unix()), account)

	fleet, err := ctx.Db.Fleets.Get(universe.Id, account, tx.Data.FleetId)
	if err != nil {
		return &abcitypes.CheckTxResponse{
			Code: types.InvalidTransactionError,
		}
	}

	if planet.CoordinateId() != fleet.Destination.CoordinateId() || ctx.Date.Unix() < fleet.BeginTravelAt+fleet.TravelDuration {
		return &abcitypes.CheckTxResponse{
			Code: types.InvalidTransactionError,
		}
	}

	payable := fleet.CanPay(tx.Data.Cargo)
	if !payable {
		return &abcitypes.CheckTxResponse{
			Code: types.InvalidTransactionError,
		}
	}

	return &abcitypes.CheckTxResponse{
		Code: types.NoError,
	}
}

func (tx *UnfeelFleetCargoTransaction) Execute(ctx t.TransactionContext) *abcitypes.ExecTxResult {
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

	fleet, err := ctx.Db.Fleets.Get(universe.Id, account, tx.Data.FleetId)
	if err != nil {
		return &abcitypes.ExecTxResult{
			Code: types.InvalidTransactionError,
		}
	}

	planet.Update(universe.Speed, ctx.Date.Unix(), account)

	err = fleet.Cargo.Sub(tx.Data.Cargo)
	if err != nil {
		return &abcitypes.ExecTxResult{
			Code: types.InvalidTransactionError,
		}
	}

	planet.Resources.Add(tx.Data.Cargo)

	err = ctx.Db.Planets.Update(tx.Data.UniverseId, planet)
	if err != nil {
		return &abcitypes.ExecTxResult{
			Code: types.InvalidTransactionError,
		}
	}

	err = ctx.Db.Fleets.Update(universe.Id, account, fleet)
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
			Type: "FleetCargoUnfilled",
			Attributes: []abcitypes.EventAttribute{
				{Key: "account", Value: tx.From, Index: true},
				{Key: "universe", Value: tx.Data.UniverseId, Index: true},
				{Key: "planet", Value: tx.Data.PlanetCoordinateId, Index: true},
				{Key: "fleet", Value: fleet.Id, Index: true},
				{Key: "oct", Value: fmt.Sprint(fleet.Cargo.OCT)},
				{Key: "metal", Value: fmt.Sprint(fleet.Cargo.Metal)},
				{Key: "crystal", Value: fmt.Sprint(fleet.Cargo.Crystal)},
				{Key: "deuterium", Value: fmt.Sprint(fleet.Cargo.Deuterium)},
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

func ParseUnfeelFleetCargoTransaction(tx t.Transaction) (UnfeelFleetCargoTransaction, error) {
	var txData UnfeelFleetCargoTransactionData
	err := cbor.Unmarshal(tx.Data, &txData)

	if err != nil {
		return UnfeelFleetCargoTransaction{}, err
	}

	return UnfeelFleetCargoTransaction{
		Type:      tx.Type,
		From:      tx.From,
		Nonce:     tx.Nonce,
		Data:      txData,
		Signature: tx.Signature,
	}, nil
}
