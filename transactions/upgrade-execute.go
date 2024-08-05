package transactions

import (
	"fmt"

	abcitypes "github.com/cometbft/cometbft/abci/types"
	"github.com/fxamacker/cbor/v2"
	"github.com/ochain-gg/ochain-network/types"
)

type ExecuteUpgradeTransactionData struct {
	Universe    string                  `cbor:"1,keyasint"`
	Planet      string                  `cbor:"2,keyasint"`
	UpgradeType types.OChainUpgradeType `cbor:"3,keyasint"`
	UpgradeId   string                  `cbor:"4,keyasint"`
}

type ExecuteUpgradeTransaction struct {
	Type TransactionType               `cbor:"1,keyasint"`
	Data ExecuteUpgradeTransactionData `cbor:"4,keyasint"`
}

func (tx *ExecuteUpgradeTransaction) Transaction() (Transaction, error) {
	txData, err := cbor.Marshal(tx.Data)
	if err != nil {
		return Transaction{}, err
	}

	return Transaction{
		Type: tx.Type,
		Data: txData,
	}, nil
}

func (tx *ExecuteUpgradeTransaction) Check(ctx TransactionContext) *abcitypes.ResponseCheckTx {
	currentDate := uint64(ctx.Date.Unix())
	universe, err := ctx.Db.Universes.GetAt(tx.Data.Universe, currentDate)
	if err != nil {
		return &abcitypes.ResponseCheckTx{
			Code: types.InvalidTransactionError,
		}
	}

	planet, err := ctx.Db.Planets.GetAt(tx.Data.Universe, tx.Data.Planet, currentDate)
	if err != nil {
		return &abcitypes.ResponseCheckTx{
			Code: types.InvalidTransactionError,
		}
	}

	pendingUpgrades, err := ctx.Db.Upgrades.GetPendingTechnologyUpgradesByPlanetAt(universe.Id, planet.CoordinateId(), uint64(ctx.Date.Unix()))
	if err != nil {
		return &abcitypes.ResponseCheckTx{
			Code: types.InvalidTransactionError,
		}
	}

	if len(pendingUpgrades) == 0 {
		return &abcitypes.ResponseCheckTx{
			Code: types.InvalidTransactionError,
		}
	}

	for i := 0; i < len(pendingUpgrades); i++ {
		upgrade := pendingUpgrades[i]

		if upgrade.UpgradeType != tx.Data.UpgradeType || tx.Data.UpgradeId != upgrade.UpgradeId {
			continue
		}

		if upgrade.EndedAt > ctx.Date.Unix() {
			return &abcitypes.ResponseCheckTx{
				Code: types.InvalidTransactionError,
			}
		}
	}

	return &abcitypes.ResponseCheckTx{
		Code: types.NoError,
	}
}

func (tx *ExecuteUpgradeTransaction) Execute(ctx TransactionContext) *abcitypes.ExecTxResult {
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

	planet, err := ctx.Db.Planets.GetAt(tx.Data.Universe, tx.Data.Planet, currentDate)
	if err != nil {
		return &abcitypes.ExecTxResult{
			Code: types.InvalidTransactionError,
		}
	}

	account, err := ctx.Db.UniverseAccounts.GetAt(universe.Id, planet.Owner, currentDate)
	if err != nil {
		return &abcitypes.ExecTxResult{
			Code: types.InvalidTransactionError,
		}
	}

	pendingUpgrades, err := ctx.Db.Upgrades.GetPendingTechnologyUpgradesByPlanetAt(universe.Id, planet.CoordinateId(), uint64(ctx.Date.Unix()))
	if err != nil {
		return &abcitypes.ExecTxResult{
			Code: types.InvalidTransactionError,
		}
	}

	if len(pendingUpgrades) == 0 {
		return &abcitypes.ExecTxResult{
			Code: types.InvalidTransactionError,
		}
	}

	var level uint64

	for i := 0; i < len(pendingUpgrades); i++ {
		upgrade := pendingUpgrades[i]

		if upgrade.UpgradeType != tx.Data.UpgradeType || tx.Data.UpgradeId != upgrade.UpgradeId {
			continue
		}

		if upgrade.EndedAt > ctx.Date.Unix() {
			return &abcitypes.ExecTxResult{
				Code: types.InvalidTransactionError,
			}
		}

		upgrade.Executed = true

		if upgrade.UpgradeType == types.OChainBuildingUpgrade {
			planet.SetBuildingLevel(types.OChainBuildingID(upgrade.UpgradeId), upgrade.Level)
			err = ctx.Db.Planets.Update(universe.Id, planet)
			if err != nil {
				return &abcitypes.ExecTxResult{
					Code: types.InvalidTransactionError,
				}
			}
		}

		if upgrade.UpgradeType == types.OChainTechnologyUpgrade {
			account.SetTechnologyLevel(types.OChainTechnologyID(upgrade.UpgradeId), upgrade.Level)
			err = ctx.Db.UniverseAccounts.Update(account)
			if err != nil {
				return &abcitypes.ExecTxResult{
					Code: types.InvalidTransactionError,
				}
			}
		}

		level = upgrade.Level

		err = ctx.Db.Upgrades.Update(upgrade)
		if err != nil {
			return &abcitypes.ExecTxResult{
				Code: types.InvalidTransactionError,
			}
		}
	}

	planet.UpdateResources(universe.Speed, ctx.Date.Unix(), account)
	err = ctx.Db.Planets.Update(universe.Id, planet)
	if err != nil {
		return &abcitypes.ExecTxResult{
			Code: types.InvalidTransactionError,
		}
	}

	events := []abcitypes.Event{
		{
			Type: "UpgradeExecuted",
			Attributes: []abcitypes.EventAttribute{
				{Key: "account", Value: account.Address, Index: true},
				{Key: "universe", Value: tx.Data.Universe, Index: true},
				{Key: "planet", Value: tx.Data.Planet, Index: true},
				{Key: "upgradeType", Value: fmt.Sprint(tx.Data.UpgradeType)},
				{Key: "upgradeId", Value: tx.Data.UpgradeId},
				{Key: "level", Value: fmt.Sprint(level)},
			},
		},
		{
			Type: "PlanetResourcesUpdated",
			Attributes: []abcitypes.EventAttribute{
				{Key: "account", Value: account.Address, Index: true},
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
		GasUsed:   0,
		GasWanted: 0,
	}
}

func ParseExecuteUpgradeTransaction(tx Transaction) (ExecuteUpgradeTransaction, error) {
	var txData ExecuteUpgradeTransactionData
	err := cbor.Unmarshal(tx.Data, &txData)

	if err != nil {
		return ExecuteUpgradeTransaction{}, err
	}

	return ExecuteUpgradeTransaction{
		Type: tx.Type,
		Data: txData,
	}, nil
}
