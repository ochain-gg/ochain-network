package transactions

import (
	"errors"
	"fmt"
	"math"

	abcitypes "github.com/cometbft/cometbft/abci/types"
	"github.com/fxamacker/cbor/v2"
	"github.com/ochain-gg/ochain-network/types"
)

type UpgradeBuildingTransactionData struct {
	Universe string `cbor:"1,keyasint"`
	Planet   string `cbor:"2,keyasint"`
	Building string `cbor:"3,keyasint"`
}

type UpgradeBuildingTransaction struct {
	Type      TransactionType                `cbor:"1,keyasint"`
	From      string                         `cbor:"2,keyasint"`
	Nonce     uint64                         `cbor:"3,keyasint"`
	Data      UpgradeBuildingTransactionData `cbor:"4,keyasint"`
	Signature []byte                         `cbor:"5,keyasint"`
}

func (tx *UpgradeBuildingTransaction) Transaction() (Transaction, error) {
	txData, err := cbor.Marshal(tx.Data)
	if err != nil {
		return Transaction{}, err
	}

	return Transaction{
		Type:      tx.Type,
		From:      tx.From,
		Nonce:     tx.Nonce,
		Data:      txData,
		Signature: tx.Signature,
	}, nil
}

func (tx *UpgradeBuildingTransaction) Check(ctx TransactionContext) error {
	_, err := ctx.Db.GlobalsAccounts.GetAt(tx.From, uint64(ctx.Date.Unix()))
	if err == nil {
		return errors.New("account doesn't exists")
	}

	account, err := ctx.Db.UniverseAccounts.GetAt(tx.Data.Universe, tx.From, uint64(ctx.Date.Unix()))
	if err == nil {
		return errors.New("universe account doesn't exists")
	}

	universe, err := ctx.Db.Universes.GetAt(tx.Data.Universe, uint64(ctx.Date.Unix()))
	if err == nil {
		return errors.New("universe doesn't exists")
	}

	planet, err := ctx.Db.Planets.GetAt(tx.Data.Universe, tx.Data.Planet, uint64(ctx.Date.Unix()))
	if err == nil {
		return errors.New("planet doesn't exists")
	}

	building, err := ctx.Db.Buildings.GetAt(tx.Data.Building, uint64(ctx.Date.Unix()))
	if err == nil {
		return errors.New("building doesn't exists")
	}

	ok := building.MeetRequirements(planet, account)
	if !ok {
		return errors.New("dependencies not met")
	}

	planet.UpdateResources(universe.Speed, int64(ctx.Date.Unix()), account)

	level := planet.BuildingLevel(building.Id) + 1
	cost := building.GetUpgradeCost(level)

	payable := planet.CanPay(cost)
	if !payable {
		return errors.New("insuficient resources")
	}

	pendingUpgrades, err := ctx.Db.Upgrades.GetPendingBuildingUpgradesByPlanetAt(universe.Id, planet.CoordinateId(), uint64(ctx.Date.Unix()))
	if err != nil {
		return errors.New("errors on pending upgrades loading")
	}

	if len(pendingUpgrades) > 0 {
		return errors.New("there is already an upgrade pending")
	}

	return nil
}

func (tx *UpgradeBuildingTransaction) Execute(ctx TransactionContext) ([]abcitypes.Event, error) {
	err := tx.Check(ctx)
	if err != nil {
		return []abcitypes.Event{}, err
	}

	currentDate := uint64(ctx.Date.Unix())
	universe, err := ctx.Db.Universes.GetAt(tx.Data.Universe, currentDate)
	if err == nil {
		return []abcitypes.Event{}, errors.New("universe doesn't exists")
	}

	account, err := ctx.Db.UniverseAccounts.GetAt(universe.Id, tx.From, currentDate)
	if err == nil {
		return []abcitypes.Event{}, errors.New("account doesn't exists")
	}

	planet, err := ctx.Db.Planets.GetAt(tx.Data.Universe, tx.Data.Planet, currentDate)
	if err == nil {
		return []abcitypes.Event{}, errors.New("planet doesn't exists")
	}

	building, err := ctx.Db.Buildings.GetAt(tx.Data.Building, currentDate)
	if err == nil {
		return []abcitypes.Event{}, errors.New("building doesn't exists")
	}

	ok := building.MeetRequirements(planet, account)
	if !ok {
		return []abcitypes.Event{}, errors.New("dependencies not met")
	}

	upgradeToLevel := planet.BuildingLevel(building.Id) + 1
	upgradeCost := building.GetUpgradeCost(upgradeToLevel)

	duration := (upgradeCost.Metal + upgradeCost.Crystal) * 3600
	duration /= (2500 * (1 + planet.BuildingLevel(types.RoboticFactoryID)) * uint64(math.Pow(float64(2), float64(planet.BuildingLevel(types.NaniteFactoryID)))) * universe.Speed)

	upgrade := types.OChainUpgrade{
		UniverseId:         universe.Id,
		PlanetCoordinateId: planet.CoordinateId(),
		UpgradeType:        types.OChainBuildingUpgrade,
		UpgradeId:          tx.Data.Building,
		Level:              planet.BuildingLevel(building.Id) + 1,
		StartedAt:          ctx.Date.Unix(),
		EndedAt:            ctx.Date.Unix() + int64(duration),
		Executed:           false,
	}

	planet.UpdateResources(universe.Speed, ctx.Date.Unix(), account)

	payable := planet.CanPay(upgradeCost)
	if !payable {
		return []abcitypes.Event{}, errors.New("insuficient resources")
	}

	err = ctx.Db.Planets.Update(tx.Data.Universe, planet)
	if err != nil {
		return []abcitypes.Event{}, err
	}

	err = ctx.Db.Upgrades.Insert(upgrade)
	if err != nil {
		return []abcitypes.Event{}, err
	}

	events := []abcitypes.Event{
		{
			Type: "UpgradeStarted",
			Attributes: []abcitypes.EventAttribute{
				{Key: "account", Value: tx.From, Index: true},
				{Key: "universe", Value: tx.Data.Universe, Index: true},
				{Key: "planet", Value: tx.Data.Planet, Index: true},
				{Key: "buildingId", Value: tx.Data.Building},
				{Key: "upgradeType", Value: fmt.Sprint(types.OChainBuildingUpgrade)},
				{Key: "upgradeId", Value: tx.Data.Building},
				{Key: "level", Value: fmt.Sprint(upgradeToLevel)},
			},
		},
		{
			Type: "PlanetResourcesUpdated",
			Attributes: []abcitypes.EventAttribute{
				{Key: "account", Value: tx.From, Index: true},
				{Key: "universe", Value: tx.Data.Universe, Index: true},
				{Key: "planet", Value: tx.Data.Planet, Index: true},
				{Key: "oct", Value: fmt.Sprint(planet.Resources.OCT)},
				{Key: "metal", Value: fmt.Sprint(planet.Resources.Metal)},
				{Key: "crystal", Value: fmt.Sprint(planet.Resources.Crystal)},
				{Key: "deuterium", Value: fmt.Sprint(planet.Resources.Deuterium)},
			},
		},
	}

	return events, nil
}

func ParseUpgradeBuildingTransaction(tx Transaction) (UpgradeBuildingTransaction, error) {
	var txData UpgradeBuildingTransactionData
	err := cbor.Unmarshal(tx.Data, &txData)

	if err != nil {
		return UpgradeBuildingTransaction{}, err
	}

	return UpgradeBuildingTransaction{
		Type:      tx.Type,
		From:      tx.From,
		Nonce:     tx.Nonce,
		Data:      txData,
		Signature: tx.Signature,
	}, nil
}
