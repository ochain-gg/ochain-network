package fleet_transactions

import (
	"fmt"

	abcitypes "github.com/cometbft/cometbft/abci/types"
	"github.com/fxamacker/cbor/v2"

	"github.com/ochain-gg/ochain-network/application"
	t "github.com/ochain-gg/ochain-network/engine/transactions"
	"github.com/ochain-gg/ochain-network/types"
)

type FleetAttackTransactionData struct {
	UniverseId      string `cbor:"1,keyasint"`
	FleetId         string `cbor:"2,keyasint"`
	TargetGalaxy    uint64 `cbor:"3,keyasint"`
	TargetSolar     uint64 `cbor:"4,keyasint"`
	TargetPlanet    uint64 `cbor:"5,keyasint"`
}

type FleetAttackTransaction struct {
	Type      t.TransactionType          `cbor:"1,keyasint"`
	From      string                     `cbor:"2,keyasint"`
	Nonce     uint64                     `cbor:"3,keyasint"`
	Data      FleetAttackTransactionData `cbor:"4,keyasint"`
	Signature []byte                     `cbor:"5,keyasint"`
}

func (tx *FleetAttackTransaction) Transaction() (t.Transaction, error) {
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

func (tx *FleetAttackTransaction) Check(ctx t.TransactionContext) *abcitypes.CheckTxResponse {
	//Load data
	account, err := ctx.Db.UniverseAccounts.GetAt(tx.Data.UniverseId, tx.From, uint64(ctx.Date.Unix()))
	if err != nil {
		return &abcitypes.CheckTxResponse{
			Code: types.UniverseAccountNotFoundError,
		}
	}

	universe, err := ctx.Db.Universes.GetAt(tx.Data.UniverseId, uint64(ctx.Date.Unix()))
	if err != nil {
		return &abcitypes.CheckTxResponse{
			Code: types.UniverseNotFoundError,
		}
	}

	fleet, err := ctx.Db.Fleets.Get(universe.Id, account, tx.Data.FleetId)
	if err != nil {
		return &abcitypes.CheckTxResponse{
			Code: types.FleetNotFoundError,
		}
	}

	//Verify fleet ownership
	if fleet.Owner != tx.From {
		return &abcitypes.CheckTxResponse{
			Code: types.FleetOwnershipError,
		}
	}

	//Verify fleet isn't traveling
	if fleet.IsTraveling(ctx.Date.Unix()) {
		return &abcitypes.CheckTxResponse{
			Code: types.FleetInTransitError,
		}
	}

	//Verify target coordinates are valid
	targetCoordId := types.CoordinateId(tx.Data.TargetGalaxy, tx.Data.TargetSolar, tx.Data.TargetPlanet)
	_, err = ctx.Db.Planets.Get(tx.Data.UniverseId, targetCoordId)
	if err != nil {
		return &abcitypes.CheckTxResponse{
			Code: types.PlanetNotFoundError,
		}
	}

	//Verify fleet is at target position (not traveling to it)
	fleetCoordId := fleet.Destination.CoordinateId()
	if fleetCoordId != targetCoordId {
		return &abcitypes.CheckTxResponse{
			Code: types.FleetNotAtTargetPositionError,
		}
	}

	return &abcitypes.CheckTxResponse{
		Code: types.NoError,
	}
}

func (tx *FleetAttackTransaction) Execute(ctx t.TransactionContext) *abcitypes.ExecTxResult {
	result := tx.Check(ctx)
	if result.Code != types.NoError {
		return &abcitypes.ExecTxResult{
			Code: result.GetCode(),
		}
	}

	// Load data
	attackerAccount, err := ctx.Db.UniverseAccounts.GetAt(tx.Data.UniverseId, tx.From, uint64(ctx.Date.Unix()))
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

	attackingFleet, err := ctx.Db.Fleets.Get(universe.Id, attackerAccount, tx.Data.FleetId)
	if err != nil {
		return &abcitypes.ExecTxResult{
			Code: types.InvalidTransactionError,
		}
	}

	// Get target planet
	targetCoordId := types.CoordinateId(tx.Data.TargetGalaxy, tx.Data.TargetSolar, tx.Data.TargetPlanet)
	targetPlanet, err := ctx.Db.Planets.Get(tx.Data.UniverseId, targetCoordId)
	if err != nil {
		return &abcitypes.ExecTxResult{
			Code: types.InvalidTransactionError,
		}
	}

	// Get defender account
	defenderAccount, err := ctx.Db.UniverseAccounts.GetAt(tx.Data.UniverseId, targetPlanet.Owner, uint64(ctx.Date.Unix()))
	if err != nil {
		return &abcitypes.ExecTxResult{
			Code: types.InvalidTransactionError,
		}
	}

	// Get all fleets at the target position
	allFleets, err := ctx.Db.Fleets.GetAllAt(tx.Data.UniverseId, uint64(ctx.Date.Unix()))
	if err != nil {
		return &abcitypes.ExecTxResult{
			Code: types.InvalidTransactionError,
		}
	}

	// Separate fleets into attackers and defenders
	var attackingFleets []types.OChainFleet
	var defendingFleets []types.OChainFleet

	for _, fleet := range allFleets {
		fleetCoordId := fleet.Destination.CoordinateId()
		if fleetCoordId == targetCoordId && !fleet.IsTraveling(ctx.Date.Unix()) {
			// Fleet is at target position and not traveling
			if fleet.Mode == types.OChainFleetAttackMode {
				attackingFleets = append(attackingFleets, fleet)
			} else if fleet.Mode == types.OChainFleetDefenseMode {
				defendingFleets = append(defendingFleets, fleet)
			}
		}
	}

	// Create battle participants
	// Attackers: all fleets in attack mode at this position
	var attackerParticipants []*types.BattleParticipant
	attackerAccounts := make(map[string]types.OChainUniverseAccount)

	for _, fleet := range attackingFleets {
		var fleetAccount types.OChainUniverseAccount
		if fleet.Owner == tx.From {
			fleetAccount = attackerAccount
		} else {
			acc, err := ctx.Db.UniverseAccounts.GetAt(tx.Data.UniverseId, fleet.Owner, uint64(ctx.Date.Unix()))
			if err != nil {
				continue // Skip if account not found
			}
			fleetAccount = acc
		}
		attackerAccounts[fleet.Owner] = fleetAccount

		participant := types.CreateBattleParticipant(
			fleet.Owner,
			fleetAccount,
			&fleet.Spaceships,
			nil,
			application.DefaultSpaceships,
			application.DefaultDefenses,
		)
		attackerParticipants = append(attackerParticipants, participant)
	}

	// Defenders: planet defenses + planet ships + defending fleets
	defenderParticipant := types.CreateBattleParticipant(
		targetPlanet.Owner,
		defenderAccount,
		&targetPlanet.Spaceships,
		&targetPlanet.Defenses,
		application.DefaultSpaceships,
		application.DefaultDefenses,
	)

	// Add defending fleets to defender participant
	for _, fleet := range defendingFleets {
		var fleetAccount types.OChainUniverseAccount
		if fleet.Owner == targetPlanet.Owner {
			fleetAccount = defenderAccount
		} else {
			acc, err := ctx.Db.UniverseAccounts.GetAt(tx.Data.UniverseId, fleet.Owner, uint64(ctx.Date.Unix()))
			if err != nil {
				continue
			}
			fleetAccount = acc
		}

		fleetParticipant := types.CreateBattleParticipant(
			fleet.Owner,
			fleetAccount,
			&fleet.Spaceships,
			nil,
			application.DefaultSpaceships,
			application.DefaultDefenses,
		)

		// Merge fleet participant units into defender participant
		defenderParticipant.Units = append(defenderParticipant.Units, fleetParticipant.Units...)
		defenderParticipant.TotalInitialUnits += fleetParticipant.TotalInitialUnits
		defenderParticipant.TotalRemainingUnits += fleetParticipant.TotalRemainingUnits
	}

	// Merge all attackers into one participant
	attackerParticipant := attackerParticipants[0]
	for i := 1; i < len(attackerParticipants); i++ {
		attackerParticipant.Units = append(attackerParticipant.Units, attackerParticipants[i].Units...)
		attackerParticipant.TotalInitialUnits += attackerParticipants[i].TotalInitialUnits
		attackerParticipant.TotalRemainingUnits += attackerParticipants[i].TotalRemainingUnits
	}

	// Execute battle
	battleResult := types.ExecuteBattle(attackerParticipant, defenderParticipant)

	// Calculate debris field
	debrisField := types.CalculateDebrisField(
		battleResult.AttackerRemaining,
		battleResult.DefenderRemaining,
		application.DefaultSpaceships,
	)

	// Apply battle results to all attacking fleets
	for _, fleet := range attackingFleets {
		types.ApplyBattleResultToFleet(&fleet, battleResult.AttackerRemaining)

		// Update or delete fleet
		if fleet.Spaceships.SmallCargo == 0 && fleet.Spaceships.LargeCargo == 0 &&
			fleet.Spaceships.LightFighter == 0 && fleet.Spaceships.HeavyFighter == 0 &&
			fleet.Spaceships.Cruiser == 0 && fleet.Spaceships.Battleship == 0 &&
			fleet.Spaceships.Battlecruiser == 0 && fleet.Spaceships.Bomber == 0 &&
			fleet.Spaceships.Destroyer == 0 && fleet.Spaceships.Deathstar == 0 &&
			fleet.Spaceships.Reaper == 0 && fleet.Spaceships.Recycler == 0 {
			// Fleet completely destroyed
			fleetAcc := attackerAccounts[fleet.Owner]
			ctx.Db.Fleets.Delete(tx.Data.UniverseId, fleetAcc, fleet.Id)
		} else {
			fleetAcc := attackerAccounts[fleet.Owner]
			ctx.Db.Fleets.Update(tx.Data.UniverseId, fleetAcc, fleet)
		}
	}

	// Apply battle results to planet
	types.ApplyBattleResultToPlanet(&targetPlanet, battleResult.DefenderRemaining)
	ctx.Db.Planets.Update(tx.Data.UniverseId, targetPlanet)

	// Apply battle results to defending fleets
	for _, fleet := range defendingFleets {
		types.ApplyBattleResultToFleet(&fleet, battleResult.DefenderRemaining)

		fleetAcc, _ := ctx.Db.UniverseAccounts.GetAt(tx.Data.UniverseId, fleet.Owner, uint64(ctx.Date.Unix()))

		if fleet.Spaceships.SmallCargo == 0 && fleet.Spaceships.LargeCargo == 0 &&
			fleet.Spaceships.LightFighter == 0 && fleet.Spaceships.HeavyFighter == 0 &&
			fleet.Spaceships.Cruiser == 0 && fleet.Spaceships.Battleship == 0 &&
			fleet.Spaceships.Battlecruiser == 0 && fleet.Spaceships.Bomber == 0 &&
			fleet.Spaceships.Destroyer == 0 && fleet.Spaceships.Deathstar == 0 &&
			fleet.Spaceships.Reaper == 0 && fleet.Spaceships.Recycler == 0 {
			ctx.Db.Fleets.Delete(tx.Data.UniverseId, fleetAcc, fleet.Id)
		} else {
			ctx.Db.Fleets.Update(tx.Data.UniverseId, fleetAcc, fleet)
		}
	}

	// Handle loot if attacker wins
	var loot types.OChainResources
	if battleResult.AttackerWins {
		// Calculate loot based on protection level (assuming 0% protection for now)
		maxLoot := targetPlanet.Resources

		// Calculate cargo capacity of surviving attacking ships
		var totalCapacity uint64 = 0
		for _, unit := range battleResult.AttackerRemaining.Units {
			if !unit.IsDefense {
				for _, ship := range application.DefaultSpaceships {
					if string(ship.Id) == unit.Type {
						totalCapacity += ship.Stats.Capacity * unit.Count
						break
					}
				}
			}
		}

		// Take resources up to cargo capacity
		if maxLoot.Metal+maxLoot.Crystal+maxLoot.Deuterium <= totalCapacity {
			loot = maxLoot
			targetPlanet.Resources = types.OChainResources{}
		} else {
			// Proportional distribution
			total := maxLoot.Metal + maxLoot.Crystal + maxLoot.Deuterium
			loot.Metal = (maxLoot.Metal * totalCapacity) / total
			loot.Crystal = (maxLoot.Crystal * totalCapacity) / total
			loot.Deuterium = (maxLoot.Deuterium * totalCapacity) / total

			targetPlanet.Resources.Metal -= loot.Metal
			targetPlanet.Resources.Crystal -= loot.Crystal
			targetPlanet.Resources.Deuterium -= loot.Deuterium
		}

		// Add loot to attacking fleet
		attackingFleet.Cargo.Add(loot)
		ctx.Db.Fleets.Update(tx.Data.UniverseId, attackerAccount, attackingFleet)
		ctx.Db.Planets.Update(tx.Data.UniverseId, targetPlanet)
	}

	// Create events
	events := []abcitypes.Event{
		{
			Type: "BattleExecuted",
			Attributes: []abcitypes.EventAttribute{
				{Key: "attacker", Value: tx.From, Index: true},
				{Key: "defender", Value: targetPlanet.Owner, Index: true},
				{Key: "universe", Value: tx.Data.UniverseId, Index: true},
				{Key: "targetCoordinates", Value: targetCoordId, Index: true},
				{Key: "rounds", Value: fmt.Sprint(battleResult.Rounds)},
				{Key: "attackerWins", Value: fmt.Sprint(battleResult.AttackerWins)},
				{Key: "defenderWins", Value: fmt.Sprint(battleResult.DefenderWins)},
				{Key: "draw", Value: fmt.Sprint(battleResult.Draw)},
				{Key: "debrisMetal", Value: fmt.Sprint(debrisField.Metal)},
				{Key: "debrisCrystal", Value: fmt.Sprint(debrisField.Crystal)},
				{Key: "lootMetal", Value: fmt.Sprint(loot.Metal)},
				{Key: "lootCrystal", Value: fmt.Sprint(loot.Crystal)},
				{Key: "lootDeuterium", Value: fmt.Sprint(loot.Deuterium)},
			},
		},
	}

	return &abcitypes.ExecTxResult{
		Code:      types.NoError,
		Events:    events,
		GasUsed:   1000,
		GasWanted: 1000,
		Data:      []byte(""),
	}
}

func ParseFleetAttackTransaction(tx t.Transaction) (FleetAttackTransaction, error) {
	var txData FleetAttackTransactionData
	err := cbor.Unmarshal(tx.Data, &txData)

	if err != nil {
		return FleetAttackTransaction{}, err
	}

	return FleetAttackTransaction{
		Type:      tx.Type,
		From:      tx.From,
		Nonce:     tx.Nonce,
		Data:      txData,
		Signature: tx.Signature,
	}, nil
}
