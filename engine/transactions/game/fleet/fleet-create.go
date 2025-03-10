package fleet_transactions

import (
	"encoding/hex"
	"fmt"

	abcitypes "github.com/cometbft/cometbft/abci/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/fxamacker/cbor/v2"

	t "github.com/ochain-gg/ochain-network/engine/transactions"
	"github.com/ochain-gg/ochain-network/types"
)

type CreateFleetTransactionData struct {
	Name               string                      `cbor:"1,keyasint"`
	UniverseId         string                      `cbor:"2,keyasint"`
	PlanetCoordinateId string                      `cbor:"3,keyasint"`
	Spaceships         types.OChainFleetSpaceships `cbor:"4,keyasint"`
	Cargo              types.OChainResources       `cbor:"5,keyasint"`
}

type CreateFleetTransaction struct {
	Type      t.TransactionType          `cbor:"1,keyasint"`
	From      string                     `cbor:"2,keyasint"`
	Nonce     uint64                     `cbor:"3,keyasint"`
	Data      CreateFleetTransactionData `cbor:"4,keyasint"`
	Signature []byte                     `cbor:"5,keyasint"`
}

func (tx *CreateFleetTransaction) Transaction() (t.Transaction, error) {
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

func (tx *CreateFleetTransaction) Check(ctx t.TransactionContext) *abcitypes.CheckTxResponse {
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

	planet.Update(universe.Speed, int64(ctx.Date.Unix()), account)
	for _, spaceshipId := range types.GetSpaceshipIds() {
		if planet.Spaceships.Get(spaceshipId) < tx.Data.Spaceships.Get(spaceshipId) {
			return &abcitypes.CheckTxResponse{
				Code: types.InvalidTransactionError,
			}
		}
	}

	maxFleet := 1 + (account.Technologies.Computer / 2)
	accountFleets, err := ctx.Db.Fleets.GetAccountFleetAt(universe.Id, account, uint64(ctx.Date.Unix()))
	if err != nil {
		return &abcitypes.CheckTxResponse{
			Code: types.InvalidTransactionError,
		}
	}

	if len(accountFleets)+1 > int(maxFleet) {
		return &abcitypes.CheckTxResponse{
			Code: types.InvalidTransactionError,
		}
	}

	fleetId := crypto.Keccak256([]byte(tx.From + fmt.Sprint(tx.Nonce) + fmt.Sprint(ctx.Date.Unix())))
	fleet := types.OChainFleet{
		Id:   hex.EncodeToString(fleetId),
		Name: tx.Data.Name,

		Owner:      tx.From,
		Universe:   tx.Data.UniverseId,
		Spaceships: tx.Data.Spaceships,
		Cargo:      tx.Data.Cargo,

		Departure: types.OChainFleetPosition{
			Galaxy:      planet.Galaxy,
			SolarSystem: planet.SolarSystem,
			Planet:      planet.Planet,
		},
		Destination: types.OChainFleetPosition{
			Galaxy:      planet.Galaxy,
			SolarSystem: planet.SolarSystem,
			Planet:      planet.Planet,
		},

		BeginTravelAt:  ctx.Date.Unix(),
		TravelDuration: 0,
	}

	availableCargo, err := ctx.Db.Spaceships.GetCargoOfAt(fleet, uint64(ctx.Date.Unix()))
	if err != nil {
		return &abcitypes.CheckTxResponse{
			Code: types.InvalidTransactionError,
		}
	}

	if availableCargo < tx.Data.Cargo.Total() {
		return &abcitypes.CheckTxResponse{
			Code: types.InvalidTransactionError,
		}
	}

	payable := planet.CanPay(tx.Data.Cargo)
	if !payable {
		return &abcitypes.CheckTxResponse{
			Code: types.InvalidTransactionError,
		}
	}

	return &abcitypes.CheckTxResponse{
		Code: types.NoError,
	}
}

func (tx *CreateFleetTransaction) Execute(ctx t.TransactionContext) *abcitypes.ExecTxResult {
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

	for _, spaceshipId := range types.GetSpaceshipIds() {
		err = planet.RemoveSpaceships(spaceshipId, tx.Data.Spaceships.Get(spaceshipId))
		if err != nil {
			return &abcitypes.ExecTxResult{
				Code: types.InvalidTransactionError,
			}
		}
	}

	err = planet.Pay(tx.Data.Cargo)
	if err != nil {
		return &abcitypes.ExecTxResult{
			Code: types.InvalidTransactionError,
		}
	}

	err = ctx.Db.Planets.Update(tx.Data.UniverseId, planet)
	if err != nil {
		return &abcitypes.ExecTxResult{
			Code: types.InvalidTransactionError,
		}
	}

	//TODO: Check if id already exists
	fleetId := crypto.Keccak256([]byte(tx.From + fmt.Sprint(tx.Nonce) + fmt.Sprint(ctx.Date.Unix())))
	fleet := types.OChainFleet{
		Id:   hex.EncodeToString(fleetId),
		Name: tx.Data.Name,

		Owner:      tx.From,
		Universe:   tx.Data.UniverseId,
		Spaceships: tx.Data.Spaceships,
		Cargo:      tx.Data.Cargo,

		Departure: types.OChainFleetPosition{
			Galaxy:      planet.Galaxy,
			SolarSystem: planet.SolarSystem,
			Planet:      planet.Planet,
		},
		Destination: types.OChainFleetPosition{
			Galaxy:      planet.Galaxy,
			SolarSystem: planet.SolarSystem,
			Planet:      planet.Planet,
		},

		BeginTravelAt:  0,
		TravelDuration: 0,
	}

	err = ctx.Db.Fleets.Insert(universe.Id, account, fleet)
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
			Type: "FleetCreated",
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

func ParseCreateFleetTransaction(tx t.Transaction) (CreateFleetTransaction, error) {
	var txData CreateFleetTransactionData
	err := cbor.Unmarshal(tx.Data, &txData)

	if err != nil {
		return CreateFleetTransaction{}, err
	}

	return CreateFleetTransaction{
		Type:      tx.Type,
		From:      tx.From,
		Nonce:     tx.Nonce,
		Data:      txData,
		Signature: tx.Signature,
	}, nil
}
