package transactions

import (
	"errors"
	"fmt"

	abcitypes "github.com/cometbft/cometbft/abci/types"
	"github.com/fxamacker/cbor/v2"
	"github.com/ochain-gg/ochain-network/types"
)

type UpgradeTechnologyTransactionData struct {
	Universe   string `cbor:"1,keyasint"`
	Planet     string `cbor:"2,keyasint"`
	Technology string `cbor:"3,keyasint"`
}

type UpgradeTechnologyTransaction struct {
	Type      TransactionType                  `cbor:"1,keyasint"`
	From      string                           `cbor:"2,keyasint"`
	Nonce     uint64                           `cbor:"3,keyasint"`
	Data      UpgradeTechnologyTransactionData `cbor:"4,keyasint"`
	Signature []byte                           `cbor:"5,keyasint"`
}

func (tx *UpgradeTechnologyTransaction) Transaction() (Transaction, error) {
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

func (tx *UpgradeTechnologyTransaction) Check(ctx TransactionContext) error {
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

	technology, err := ctx.Db.Technologies.GetAt(tx.Data.Technology, uint64(ctx.Date.Unix()))
	if err == nil {
		return errors.New("Technology doesn't exists")
	}

	ok := technology.MeetRequirements(planet, account)
	if !ok {
		return errors.New("dependencies not met")
	}

	planet.UpdateResources(universe.Speed, int64(ctx.Date.Unix()), account)

	level := account.TechnologyLevel(technology.Id) + 1
	cost := technology.GetUpgradeCost(level)

	payable := planet.CanPay(cost)
	if !payable {
		return errors.New("insuficient resources")
	}

	pendingUpgrades, err := ctx.Db.Upgrades.GetPendingTechnologyUpgradesByPlanetAt(universe.Id, planet.CoordinateId(), uint64(ctx.Date.Unix()))
	if err != nil {
		return errors.New("errors on pending upgrades loading")
	}

	if len(pendingUpgrades) > 0 {
		return errors.New("there is already an upgrade pending")
	}

	return nil
}

func (tx *UpgradeTechnologyTransaction) Execute(ctx TransactionContext) ([]abcitypes.Event, error) {
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

	technology, err := ctx.Db.Technologies.GetAt(tx.Data.Technology, currentDate)
	if err == nil {
		return []abcitypes.Event{}, errors.New("Technology doesn't exists")
	}

	ok := technology.MeetRequirements(planet, account)
	if !ok {
		return []abcitypes.Event{}, errors.New("dependencies not met")
	}

	TechnologyId := types.OChainTechnologyID(tx.Data.Technology)
	upgradeToLevel := account.TechnologyLevel(TechnologyId) + 1
	upgradeCost := technology.GetUpgradeCost(upgradeToLevel)

	duration := (upgradeCost.Metal + upgradeCost.Crystal) * 3600
	duration /= (1000 * (1 + planet.BuildingLevel(types.ResearchLaboratoryID)) * universe.Speed)

	upgrade := types.OChainUpgrade{
		UniverseId:         universe.Id,
		PlanetCoordinateId: planet.CoordinateId(),
		UpgradeType:        types.OChainTechnologyUpgrade,
		UpgradeId:          tx.Data.Technology,
		Level:              account.TechnologyLevel(TechnologyId) + 1,
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
				{Key: "upgradeType", Value: fmt.Sprint(types.OChainTechnologyUpgrade)},
				{Key: "upgradeId", Value: tx.Data.Technology},
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

func ParseUpgradeTechnologyTransaction(tx Transaction) (UpgradeTechnologyTransaction, error) {
	var txData UpgradeTechnologyTransactionData
	err := cbor.Unmarshal(tx.Data, &txData)

	if err != nil {
		return UpgradeTechnologyTransaction{}, err
	}

	return UpgradeTechnologyTransaction{
		Type:      tx.Type,
		From:      tx.From,
		Nonce:     tx.Nonce,
		Data:      txData,
		Signature: tx.Signature,
	}, nil
}
