package fleet_transactions

import (
	"fmt"

	abcitypes "github.com/cometbft/cometbft/abci/types"
	"github.com/fxamacker/cbor/v2"

	t "github.com/ochain-gg/ochain-network/engine/transactions"
	"github.com/ochain-gg/ochain-network/types"
)

type LendFleetCargoTransactionData struct {
	UniverseId         string `cbor:"2,keyasint"`
	PlanetCoordinateId string `cbor:"3,keyasint"`
	FleetId            string `cbor:"4,keyasint"`
}

type LendFleetCargoTransaction struct {
	Type      t.TransactionType             `cbor:"1,keyasint"`
	From      string                        `cbor:"2,keyasint"`
	Nonce     uint64                        `cbor:"3,keyasint"`
	Data      LendFleetCargoTransactionData `cbor:"4,keyasint"`
	Signature []byte                        `cbor:"5,keyasint"`
}

func (tx *LendFleetCargoTransaction) Transaction() (t.Transaction, error) {
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

func (tx *LendFleetCargoTransaction) Check(ctx t.TransactionContext) *abcitypes.CheckTxResponse {
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

	return &abcitypes.CheckTxResponse{
		Code: types.NoError,
	}
}

func (tx *LendFleetCargoTransaction) Execute(ctx t.TransactionContext) *abcitypes.ExecTxResult {
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
	planet.Resources.Add(fleet.Cargo)

	for _, spaceshipId := range types.GetSpaceshipIds() {
		planet.AddSpaceships(spaceshipId, fleet.Spaceships.Get(spaceshipId))
	}

	err = ctx.Db.Planets.Update(tx.Data.UniverseId, planet)
	if err != nil {
		return &abcitypes.ExecTxResult{
			Code: types.InvalidTransactionError,
		}
	}

	err = ctx.Db.Fleets.Delete(universe.Id, account, fleet.Id)
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
			Type: "FleetLanded",
			Attributes: []abcitypes.EventAttribute{
				{Key: "account", Value: tx.From, Index: true},
				{Key: "universe", Value: tx.Data.UniverseId, Index: true},
				{Key: "planet", Value: tx.Data.PlanetCoordinateId, Index: true},
				{Key: "fleet", Value: fleet.Id, Index: true},
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

func ParseLendFleetCargoTransaction(tx t.Transaction) (LendFleetCargoTransaction, error) {
	var txData LendFleetCargoTransactionData
	err := cbor.Unmarshal(tx.Data, &txData)

	if err != nil {
		return LendFleetCargoTransaction{}, err
	}

	return LendFleetCargoTransaction{
		Type:      tx.Type,
		From:      tx.From,
		Nonce:     tx.Nonce,
		Data:      txData,
		Signature: tx.Signature,
	}, nil
}
