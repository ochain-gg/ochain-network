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

type SplitFleetTransactionData struct {
	Name       string                      `cbor:"1,keyasint"`
	UniverseId string                      `cbor:"2,keyasint"`
	FleetId    string                      `cbor:"3,keyasint"`
	Spaceships types.OChainFleetSpaceships `cbor:"4,keyasint"`
	Cargo      types.OChainResources       `cbor:"5,keyasint"`
}

type SplitFleetTransaction struct {
	Type      t.TransactionType         `cbor:"1,keyasint"`
	From      string                    `cbor:"2,keyasint"`
	Nonce     uint64                    `cbor:"3,keyasint"`
	Data      SplitFleetTransactionData `cbor:"4,keyasint"`
	Signature []byte                    `cbor:"5,keyasint"`
}

func (tx *SplitFleetTransaction) Transaction() (t.Transaction, error) {
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

func (tx *SplitFleetTransaction) Check(ctx t.TransactionContext) *abcitypes.CheckTxResponse {
	//Load data
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

	fleet, err := ctx.Db.Fleets.Get(universe.Id, account, tx.Data.FleetId)
	if err != nil {
		return &abcitypes.CheckTxResponse{
			Code: types.InvalidTransactionError,
		}
	}

	//Verify fleet ownership
	if fleet.Owner != tx.From {
		return &abcitypes.CheckTxResponse{
			Code: types.InvalidTransactionError,
		}
	}

	//Verify fleet isn't travaling
	if fleet.IsTraveling(ctx.Date.Unix()) {
		return &abcitypes.CheckTxResponse{
			Code: types.InvalidTransactionError,
		}
	}

	//Verify spaceship avaibility
	for _, spaceshipId := range types.GetSpaceshipIds() {
		if fleet.Spaceships.Get(spaceshipId) < tx.Data.Spaceships.Get(spaceshipId) {
			return &abcitypes.CheckTxResponse{
				Code: types.InvalidTransactionError,
			}
		}
	}

	//Verify max fleets isn't reach
	maxFleet := 1 + (account.Technologies.Computer / 2)
	accountFleets, err := ctx.Db.Fleets.GetAccountFleetAt(universe.Id, account, uint64(ctx.Date.Unix()))
	if err != nil {
		return &abcitypes.CheckTxResponse{
			Code: types.InvalidTransactionError,
		}
	}

	if len(accountFleets) > int(maxFleet) {
		return &abcitypes.CheckTxResponse{
			Code: types.InvalidTransactionError,
		}
	}

	fleetId := crypto.Keccak256([]byte(tx.From + fmt.Sprint(tx.Nonce) + fmt.Sprint(ctx.Date.Unix())))
	newFleet := types.OChainFleet{
		Id:   hex.EncodeToString(fleetId),
		Name: tx.Data.Name,

		Owner:      tx.From,
		Universe:   tx.Data.UniverseId,
		Spaceships: tx.Data.Spaceships,
		Cargo:      tx.Data.Cargo,

		Departure: types.OChainFleetPosition{
			Galaxy:      fleet.Destination.Galaxy,
			SolarSystem: fleet.Destination.SolarSystem,
			Planet:      fleet.Destination.Planet,
		},
		Destination: types.OChainFleetPosition{
			Galaxy:      fleet.Destination.Galaxy,
			SolarSystem: fleet.Destination.SolarSystem,
			Planet:      fleet.Destination.Planet,
		},

		BeginTravelAt:  ctx.Date.Unix(),
		TravelDuration: 0,
	}

	//Verify available cargo
	splitedFleetAvailableCargo, err := ctx.Db.Spaceships.GetCargoOfAt(fleet, uint64(ctx.Date.Unix()))
	if err != nil {
		return &abcitypes.CheckTxResponse{
			Code: types.InvalidTransactionError,
		}
	}

	newFleetAvailableCargo, err := ctx.Db.Spaceships.GetCargoOfAt(newFleet, uint64(ctx.Date.Unix()))
	if err != nil {
		return &abcitypes.CheckTxResponse{
			Code: types.InvalidTransactionError,
		}
	}

	//Check if splitted fleet isn't empty
	if splitedFleetAvailableCargo-newFleetAvailableCargo == 0 {
		return &abcitypes.CheckTxResponse{
			Code: types.InvalidTransactionError,
		}
	}

	if splitedFleetAvailableCargo-newFleetAvailableCargo < fleet.Cargo.Total()-tx.Data.Cargo.Total() {
		return &abcitypes.CheckTxResponse{
			Code: types.InvalidTransactionError,
		}
	}

	if newFleetAvailableCargo < tx.Data.Cargo.Total() {
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

func (tx *SplitFleetTransaction) Execute(ctx t.TransactionContext) *abcitypes.ExecTxResult {
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

	fleet, err := ctx.Db.Fleets.Get(universe.Id, account, tx.Data.FleetId)
	if err != nil {
		return &abcitypes.ExecTxResult{
			Code: types.InvalidTransactionError,
		}
	}

	for _, spaceshipId := range types.GetSpaceshipIds() {
		err = fleet.RemoveSpaceships(spaceshipId, tx.Data.Spaceships.Get(spaceshipId))
		if err != nil {
			return &abcitypes.ExecTxResult{
				Code: types.InvalidTransactionError,
			}
		}
	}

	err = fleet.Cargo.Sub(tx.Data.Cargo)
	if err != nil {
		return &abcitypes.ExecTxResult{
			Code: types.InvalidTransactionError,
		}
	}

	fleetId := crypto.Keccak256([]byte(tx.From + fmt.Sprint(tx.Nonce) + fmt.Sprint(ctx.Date.Unix())))
	newFleet := types.OChainFleet{
		Id:   hex.EncodeToString(fleetId),
		Name: tx.Data.Name,

		Owner:      tx.From,
		Universe:   tx.Data.UniverseId,
		Spaceships: tx.Data.Spaceships,
		Cargo:      tx.Data.Cargo,

		Departure: types.OChainFleetPosition{
			Galaxy:      fleet.Destination.Galaxy,
			SolarSystem: fleet.Destination.SolarSystem,
			Planet:      fleet.Destination.Planet,
		},
		Destination: types.OChainFleetPosition{
			Galaxy:      fleet.Destination.Galaxy,
			SolarSystem: fleet.Destination.SolarSystem,
			Planet:      fleet.Destination.Planet,
		},

		BeginTravelAt:  ctx.Date.Unix(),
		TravelDuration: 0,
	}

	err = ctx.Db.Fleets.Insert(universe.Id, account, newFleet)
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

	s, err := cbor.Marshal(newFleet.Spaceships.WithAttributes())
	if err != nil {
		return &abcitypes.ExecTxResult{
			Code: types.InvalidTransactionError,
		}
	}

	events := []abcitypes.Event{
		{
			Type: "FleetSplited",
			Attributes: []abcitypes.EventAttribute{
				{Key: "account", Value: tx.From, Index: true},
				{Key: "universe", Value: tx.Data.UniverseId, Index: true},
				{Key: "splitedFleet", Value: fleet.Id, Index: true},
				{Key: "newFleet", Value: newFleet.Id, Index: true},
				{Key: "spaceships", Value: hex.EncodeToString(s), Index: true},
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

func ParseSplitFleetTransaction(tx t.Transaction) (SplitFleetTransaction, error) {
	var txData SplitFleetTransactionData
	err := cbor.Unmarshal(tx.Data, &txData)

	if err != nil {
		return SplitFleetTransaction{}, err
	}

	return SplitFleetTransaction{
		Type:      tx.Type,
		From:      tx.From,
		Nonce:     tx.Nonce,
		Data:      txData,
		Signature: tx.Signature,
	}, nil
}
