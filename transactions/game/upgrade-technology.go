package game_transactions

import (
	"fmt"

	abcitypes "github.com/cometbft/cometbft/abci/types"
	"github.com/fxamacker/cbor/v2"

	t "github.com/ochain-gg/ochain-network/transactions"
	"github.com/ochain-gg/ochain-network/types"
)

type UpgradeTechnologyTransactionData struct {
	Universe   string `cbor:"1,keyasint"`
	Planet     string `cbor:"2,keyasint"`
	Technology string `cbor:"3,keyasint"`
}

type UpgradeTechnologyTransaction struct {
	Type      t.TransactionType                `cbor:"1,keyasint"`
	From      string                           `cbor:"2,keyasint"`
	Nonce     uint64                           `cbor:"3,keyasint"`
	Data      UpgradeTechnologyTransactionData `cbor:"4,keyasint"`
	Signature []byte                           `cbor:"5,keyasint"`
}

func (tx *UpgradeTechnologyTransaction) Transaction() (t.Transaction, error) {
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

func (tx *UpgradeTechnologyTransaction) Check(ctx t.TransactionContext) *abcitypes.CheckTxResponse {
	_, err := ctx.Db.GlobalsAccounts.GetAt(tx.From, uint64(ctx.Date.Unix()))
	if err != nil {
		return &abcitypes.CheckTxResponse{
			Code: types.InvalidTransactionError,
		}
	}

	account, err := ctx.Db.UniverseAccounts.GetAt(tx.Data.Universe, tx.From, uint64(ctx.Date.Unix()))
	if err != nil {
		return &abcitypes.CheckTxResponse{
			Code: types.InvalidTransactionError,
		}
	}

	universe, err := ctx.Db.Universes.GetAt(tx.Data.Universe, uint64(ctx.Date.Unix()))
	if err != nil {
		return &abcitypes.CheckTxResponse{
			Code: types.InvalidTransactionError,
		}
	}

	planet, err := ctx.Db.Planets.GetAt(tx.Data.Universe, tx.Data.Planet, uint64(ctx.Date.Unix()))
	if err != nil {
		return &abcitypes.CheckTxResponse{
			Code: types.InvalidTransactionError,
		}
	}

	technology, err := ctx.Db.Technologies.GetAt(tx.Data.Technology, uint64(ctx.Date.Unix()))
	if err != nil {
		return &abcitypes.CheckTxResponse{
			Code: types.InvalidTransactionError,
		}
	}

	ok := technology.MeetRequirements(planet, account)
	if !ok {
		return &abcitypes.CheckTxResponse{
			Code: types.InvalidTransactionError,
		}
	}

	planet.UpdateResources(universe.Speed, int64(ctx.Date.Unix()), account)

	level := account.TechnologyLevel(technology.Id) + 1
	cost := technology.GetUpgradeCost(level)

	payable := planet.CanPay(cost)
	if !payable {
		return &abcitypes.CheckTxResponse{
			Code: types.InvalidTransactionError,
		}
	}

	pendingUpgrades, err := ctx.Db.Upgrades.GetPendingTechnologyUpgradesByPlanetAt(universe.Id, planet.CoordinateId(), uint64(ctx.Date.Unix()))
	if err != nil {
		return &abcitypes.CheckTxResponse{
			Code: types.InvalidTransactionError,
		}
	}

	if len(pendingUpgrades) > 0 {
		return &abcitypes.CheckTxResponse{
			Code: types.InvalidTransactionError,
		}
	}

	return nil
}

func (tx *UpgradeTechnologyTransaction) Execute(ctx t.TransactionContext) *abcitypes.ExecTxResult {
	result := tx.Check(ctx)
	if result.Code != types.NoError {
		return &abcitypes.ExecTxResult{
			Code: result.GetCode(),
		}
	}

	currentDate := uint64(ctx.Date.Unix())
	universe, err := ctx.Db.Universes.GetAt(tx.Data.Universe, currentDate)
	if err != nil {
		return &abcitypes.ExecTxResult{
			Code: types.InvalidTransactionError,
		}
	}

	globalAccount, err := ctx.Db.GlobalsAccounts.GetAt(tx.From, currentDate)
	if err != nil {
		return &abcitypes.ExecTxResult{
			Code: types.InvalidTransactionError,
		}
	}

	account, err := ctx.Db.UniverseAccounts.GetAt(universe.Id, tx.From, currentDate)
	if err != nil {
		return &abcitypes.ExecTxResult{
			Code: types.InvalidTransactionError,
		}
	}

	planet, err := ctx.Db.Planets.GetAt(tx.Data.Universe, tx.Data.Planet, currentDate)
	if err != nil {
		return &abcitypes.ExecTxResult{
			Code: types.InvalidTransactionError,
		}
	}

	technology, err := ctx.Db.Technologies.GetAt(tx.Data.Technology, currentDate)
	if err != nil {
		return &abcitypes.ExecTxResult{
			Code: types.InvalidTransactionError,
		}
	}

	ok := technology.MeetRequirements(planet, account)
	if !ok {
		return &abcitypes.ExecTxResult{
			Code: types.InvalidTransactionError,
		}
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
		return &abcitypes.ExecTxResult{
			Code: types.InvalidTransactionError,
		}
	}

	err = ctx.Db.Planets.Update(tx.Data.Universe, planet)
	if err != nil {
		return &abcitypes.ExecTxResult{
			Code: types.InvalidTransactionError,
		}
	}

	err = ctx.Db.Upgrades.Insert(upgrade)
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

	return &abcitypes.ExecTxResult{
		Code:      types.NoError,
		Events:    events,
		GasUsed:   100,
		GasWanted: 100,
		Data:      receipt.Bytes(),
	}
}

func ParseUpgradeTechnologyTransaction(tx t.Transaction) (UpgradeTechnologyTransaction, error) {
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
