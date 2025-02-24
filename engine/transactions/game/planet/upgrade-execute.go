package planet_transactions

import (
	"fmt"
	"log"

	abcitypes "github.com/cometbft/cometbft/abci/types"
	"github.com/fxamacker/cbor/v2"

	t "github.com/ochain-gg/ochain-network/engine/transactions"
	"github.com/ochain-gg/ochain-network/types"
)

type ExecuteUpgradeTransactionData struct {
	Universe    string                  `cbor:"1,keyasint"`
	Planet      string                  `cbor:"2,keyasint"`
	UpgradeType types.OChainUpgradeType `cbor:"3,keyasint"`
	UpgradeId   string                  `cbor:"4,keyasint"`
	Level       uint64                  `cbor:"5,keyasint"`
}

type ExecuteUpgradeTransaction struct {
	Type t.TransactionType             `cbor:"1,keyasint"`
	Data ExecuteUpgradeTransactionData `cbor:"4,keyasint"`
}

func (tx *ExecuteUpgradeTransaction) Transaction() (t.Transaction, error) {
	txData, err := cbor.Marshal(tx.Data)
	if err != nil {
		return t.Transaction{}, err
	}

	return t.Transaction{
		Type: tx.Type,
		Data: txData,
	}, nil
}

func (tx *ExecuteUpgradeTransaction) Check(ctx t.TransactionContext) *abcitypes.CheckTxResponse {
	currentDate := uint64(ctx.Date.Unix())
	universe, err := ctx.Db.Universes.GetAt(tx.Data.Universe, currentDate)
	if err != nil {
		return &abcitypes.CheckTxResponse{
			Code: types.InvalidTransactionError,
		}
	}

	planet, err := ctx.Db.Planets.GetAt(tx.Data.Universe, tx.Data.Planet, currentDate)
	if err != nil {
		return &abcitypes.CheckTxResponse{
			Code: types.InvalidTransactionError,
		}
	}

	upgrades, err := ctx.Db.Upgrades.GetByPlanetAt(universe.Id, planet.CoordinateId(), uint64(ctx.Date.Unix()))
	if err != nil {
		return &abcitypes.CheckTxResponse{
			Code: types.InvalidTransactionError,
		}
	}

	if len(upgrades) == 0 {
		return &abcitypes.CheckTxResponse{
			Code: types.InvalidTransactionError,
		}
	}

	isFound := false
	for i := 0; i < len(upgrades); i++ {
		upgrade := upgrades[i]

		if upgrade.Executed || upgrade.UpgradeType != tx.Data.UpgradeType || tx.Data.UpgradeId != upgrade.UpgradeId || tx.Data.Level != upgrade.Level {
			continue
		}

		if upgrade.EndedAt < ctx.Date.Unix() {
			isFound = true
		}
	}

	if !isFound {
		return &abcitypes.CheckTxResponse{
			Code: types.InvalidTransactionError,
		}
	}

	log.Println("UPGRADE FOUND")

	return &abcitypes.CheckTxResponse{
		Code: types.NoError,
	}
}

func (tx *ExecuteUpgradeTransaction) Execute(ctx t.TransactionContext) *abcitypes.ExecTxResult {
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

	pendingUpgrades, err := ctx.Db.Upgrades.GetByPlanetAt(universe.Id, planet.CoordinateId(), currentDate)
	if err != nil {
		return &abcitypes.ExecTxResult{
			Code: types.InvalidTransactionError,
		}
	}

	log.Println("UPGRADES COUNT: " + fmt.Sprint(len(pendingUpgrades)))

	var upgrade types.OChainUpgrade
	for i := 0; i < len(pendingUpgrades); i++ {
		u := pendingUpgrades[i]
		log.Println("UPGRADE: " + fmt.Sprint(i))
		log.Println(u)
		if u.Executed || u.UpgradeType != tx.Data.UpgradeType || tx.Data.UpgradeId != u.UpgradeId || tx.Data.Level != u.Level {
			continue
		}

		log.Println("UPGRADE SELECTED: " + fmt.Sprint(i))
		log.Println("UPGRADE IS EXECUTED: " + fmt.Sprint(upgrade.Executed))
		upgrade = u
	}

	log.Println("UPGRADE FOUND")
	log.Println(upgrade)

	planet.UpdateResources(universe.Speed, ctx.Date.Unix(), account)
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

	err = ctx.Db.Upgrades.Update(upgrade)
	if err != nil {
		return &abcitypes.ExecTxResult{
			Code: types.InvalidTransactionError,
		}
	}

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
				{Key: "level", Value: fmt.Sprint(upgrade.Level)},
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

func ParseExecuteUpgradeTransaction(tx t.Transaction) (ExecuteUpgradeTransaction, error) {
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
