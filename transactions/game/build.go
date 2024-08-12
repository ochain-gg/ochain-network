package game_transactions

import (
	"fmt"
	"math"

	abcitypes "github.com/cometbft/cometbft/abci/types"
	"github.com/fxamacker/cbor/v2"

	t "github.com/ochain-gg/ochain-network/transactions"
	"github.com/ochain-gg/ochain-network/types"
)

type BuildTransactionData struct {
	UniverseId         string              `cbor:"1,keyasint"`
	PlanetCoordinateId string              `cbor:"2,keyasint"`
	Builds             []types.OChainBuild `cbor:"3,keyasint"`
}

type BuildTransaction struct {
	Type      t.TransactionType    `cbor:"1,keyasint"`
	From      string               `cbor:"2,keyasint"`
	Nonce     uint64               `cbor:"3,keyasint"`
	Data      BuildTransactionData `cbor:"4,keyasint"`
	Signature []byte               `cbor:"5,keyasint"`
}

func (tx *BuildTransaction) Transaction() (t.Transaction, error) {
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

func (tx *BuildTransaction) Check(ctx t.TransactionContext) *abcitypes.ResponseCheckTx {
	_, err := ctx.Db.GlobalsAccounts.GetAt(tx.From, uint64(ctx.Date.Unix()))
	if err != nil {
		return &abcitypes.ResponseCheckTx{
			Code: types.InvalidTransactionError,
		}
	}

	account, err := ctx.Db.UniverseAccounts.GetAt(tx.Data.UniverseId, tx.From, uint64(ctx.Date.Unix()))
	if err != nil {
		return &abcitypes.ResponseCheckTx{
			Code: types.InvalidTransactionError,
		}
	}

	universe, err := ctx.Db.Universes.GetAt(tx.Data.UniverseId, uint64(ctx.Date.Unix()))
	if err != nil {
		return &abcitypes.ResponseCheckTx{
			Code: types.InvalidTransactionError,
		}
	}

	planet, err := ctx.Db.Planets.GetAt(tx.Data.UniverseId, tx.Data.PlanetCoordinateId, uint64(ctx.Date.Unix()))
	if err != nil {
		return &abcitypes.ResponseCheckTx{
			Code: types.InvalidTransactionError,
		}
	}

	totalCost := types.OChainResources{
		OCT:       0,
		Metal:     0,
		Crystal:   0,
		Deuterium: 0,
	}

	for i := 0; i < len(tx.Data.Builds); i++ {
		build := tx.Data.Builds[i]
		if build.BuildType == types.OChainSpaceshipBuild {

			spaceship, err := ctx.Db.Spaceships.GetAt(build.BuildId, uint64(ctx.Date.Unix()))
			if err != nil {
				return &abcitypes.ResponseCheckTx{
					Code: types.InvalidTransactionError,
				}
			}

			ok := spaceship.MeetRequirements(planet, account)
			if !ok {
				return &abcitypes.ResponseCheckTx{
					Code: types.InvalidTransactionError,
				}
			}

			cost := spaceship.Cost
			cost.Mul(build.Count)
			totalCost.Add(cost)
		}

		if build.BuildType == types.OChainDefenseBuild {
			defense, err := ctx.Db.Defenses.GetAt(build.BuildId, uint64(ctx.Date.Unix()))
			if err != nil {
				return &abcitypes.ResponseCheckTx{
					Code: types.InvalidTransactionError,
				}
			}

			ok := defense.MeetRequirements(planet, account)
			if !ok {
				return &abcitypes.ResponseCheckTx{
					Code: types.InvalidTransactionError,
				}
			}

			cost := defense.Cost
			cost.Mul(build.Count)
			totalCost.Add(cost)
		}
	}

	planet.UpdateResources(universe.Speed, int64(ctx.Date.Unix()), account)
	payable := planet.CanPay(totalCost)
	if !payable {
		return &abcitypes.ResponseCheckTx{
			Code: types.InvalidTransactionError,
		}
	}

	return &abcitypes.ResponseCheckTx{
		Code: types.NoError,
	}
}

func (tx *BuildTransaction) Execute(ctx t.TransactionContext) *abcitypes.ExecTxResult {
	result := tx.Check(ctx)
	if result.Code != types.NoError {
		return &abcitypes.ExecTxResult{
			Code: result.GetCode(),
		}
	}

	globalAccount, err := ctx.Db.GlobalsAccounts.GetAt(tx.From, uint64(ctx.Date.Unix()))
	if err != nil {
		return &abcitypes.ExecTxResult{
			Code: types.InvalidTransactionError,
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

	totalCost := types.OChainResources{
		OCT:       0,
		Metal:     0,
		Crystal:   0,
		Deuterium: 0,
	}

	planet.UpdateBuildQueue(uint64(ctx.Date.Unix()))
	planet.UpdateResources(universe.Speed, int64(ctx.Date.Unix()), account)

	var buildEvents []abcitypes.Event
	for i := 0; i < len(tx.Data.Builds); i++ {
		build := tx.Data.Builds[i]
		if build.BuildType == types.OChainSpaceshipBuild {

			spaceship, err := ctx.Db.Spaceships.GetAt(build.BuildId, uint64(ctx.Date.Unix()))
			if err != nil {
				return &abcitypes.ExecTxResult{
					Code: types.InvalidTransactionError,
				}
			}

			duration := (spaceship.Cost.Metal + spaceship.Cost.Crystal) * 3600
			duration /= (2500 * (1 + planet.BuildingLevel(types.SpaceshipFactoryID)) * uint64(math.Pow(float64(2), float64(planet.BuildingLevel(types.NaniteFactoryID)))) * universe.Speed)
			duration *= build.Count

			item := planet.AddItemToBuildQueue(build, duration)

			buildEvents = append(buildEvents, abcitypes.Event{
				Type: "BuildQueueItemAdded",
				Attributes: []abcitypes.EventAttribute{
					{Key: "account", Value: tx.From, Index: true},
					{Key: "universe", Value: tx.Data.UniverseId, Index: true},
					{Key: "planet", Value: tx.Data.PlanetCoordinateId, Index: true},
					{Key: "buildType", Value: fmt.Sprint(item.BuildType)},
					{Key: "buildId", Value: item.BuildId},
					{Key: "count", Value: fmt.Sprint(item.Count)},
					{Key: "startAt", Value: fmt.Sprint(item.StartAt)},
					{Key: "finishAt", Value: fmt.Sprint(item.FinishAt)},
				},
			})

			cost := spaceship.Cost
			cost.Mul(build.Count)
			totalCost.Add(cost)
		}

		if build.BuildType == types.OChainDefenseBuild {
			defense, err := ctx.Db.Defenses.GetAt(build.BuildId, uint64(ctx.Date.Unix()))
			if err != nil {
				return &abcitypes.ExecTxResult{
					Code: types.InvalidTransactionError,
				}
			}

			duration := (defense.Cost.Metal + defense.Cost.Crystal) * 3600
			duration /= (2500 * (1 + planet.BuildingLevel(types.SpaceshipFactoryID)) * uint64(math.Pow(float64(2), float64(planet.BuildingLevel(types.NaniteFactoryID)))) * universe.Speed)
			duration *= build.Count

			item := planet.AddItemToBuildQueue(build, duration)

			buildEvents = append(buildEvents, abcitypes.Event{
				Type: "BuildQueueItemAdded",
				Attributes: []abcitypes.EventAttribute{
					{Key: "account", Value: tx.From, Index: true},
					{Key: "universe", Value: tx.Data.UniverseId, Index: true},
					{Key: "planet", Value: tx.Data.PlanetCoordinateId, Index: true},
					{Key: "buildType", Value: fmt.Sprint(item.BuildType)},
					{Key: "buildId", Value: item.BuildId},
					{Key: "count", Value: fmt.Sprint(item.Count)},
					{Key: "startAt", Value: fmt.Sprint(item.StartAt)},
					{Key: "finishAt", Value: fmt.Sprint(item.FinishAt)},
				},
			})

			cost := defense.Cost
			cost.Mul(build.Count)
			totalCost.Add(cost)
		}
	}

	err = planet.Pay(totalCost)
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

	txGasCost, err := globalAccount.ApplyGasCost(uint64(ctx.Date.Unix()))
	if err != nil {
		return &abcitypes.ExecTxResult{
			Code: types.GasCostHigherThanBalance,
		}
	}

	receipt := t.TransactionReceipt{
		GasCost: txGasCost,
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
	}

	for i := 0; i < len(buildEvents); i++ {
		events = append(events, buildEvents[i])
	}

	return &abcitypes.ExecTxResult{
		Code:      types.NoError,
		Events:    events,
		GasUsed:   100,
		GasWanted: 100,
		Data:      receipt.Bytes(),
	}
}

func ParseBuildTransaction(tx t.Transaction) (BuildTransaction, error) {
	var txData BuildTransactionData
	err := cbor.Unmarshal(tx.Data, &txData)

	if err != nil {
		return BuildTransaction{}, err
	}

	return BuildTransaction{
		Type:      tx.Type,
		From:      tx.From,
		Nonce:     tx.Nonce,
		Data:      txData,
		Signature: tx.Signature,
	}, nil
}
