package types

import (
	"math/rand"
)

/*
Battle Engine - OGame Combat Mechanics Implementation

This file implements the battle mechanics following OGame's combat system.

COMBAT RULES (Source: OGame official mechanics):

1. Technology Bonuses:
   - Weapon technology: +10% attack per level
   - Shielding technology: +10% shield per level
   - Armor technology: +10% armor per level

2. Combat Rounds:
   - Battle lasts up to 6 rounds maximum
   - Each round, all units fire simultaneously on random targets
   - Shields regenerate after each round (for surviving units)
   - Battle ends when one side is eliminated or max rounds reached
   - Draw occurs if both sides are eliminated or 6 rounds complete with units remaining

3. Damage Calculation:
   - Damage is applied to shields first
   - **Shield Bounce Rule**: If damage is less than 1% of shield, it bounces off (no damage)
   - Overflow damage after shield depletion goes to hull
   - Hull represents 10% of structure points (BaseArmor in our implementation)

4. Explosion Probability (Critical Rule):
   - If hull damage is less than 30%: NO chance of explosion
   - If hull damage is >= 30%: explosion chance = damage percentage
   - Examples:
     * 20% hull damage → 0% explosion chance (below threshold)
     * 40% hull damage → 40% explosion chance
     * 80% hull damage → 80% explosion chance
   - Each shot that causes hull damage triggers an explosion check

5. Rapid Fire (RF):
   - Some units have rapid fire against specific targets
   - Formula: (RF-1)/RF chance to fire again
   - Example: Cruiser has RF=6 vs Light Fighter → 5/6 = 83.33% chance to shoot again
   - After each shot at a target type with RF, probability check for another shot
   - Can chain multiple times in the same round

6. Debris Field:
   - 30% of metal and crystal from destroyed ships becomes debris
   - Defenses do NOT create debris
   - Only ships contribute to debris field
   - Debris can be harvested by recycler ships

7. Defense Rebuilding:
   - After battle, 70% of destroyed defenses are automatically rebuilt
   - Only 30% are permanently lost
   - Implemented deterministically: destroyed × 0.7 = rebuilt
   - Example: 100 defenses destroyed → 70 rebuilt, 30 actually lost

IMPLEMENTATION NOTES:
- Random target selection for each shot
- Simultaneous firing simulation (all attackers fire, then all defenders)
- Shield regeneration occurs between rounds for all surviving units
- Units marked for explosion during combat continue participating until round end

Usage Example:
	attacker := CreateBattleParticipant(ownerAddr, attackerAccount, &fleet.Spaceships, nil, spaceshipDefs, defenseDefs)
	defender := CreateBattleParticipant(defenderAddr, defenderAccount, &planet.Spaceships, &planet.Defenses, spaceshipDefs, defenseDefs)
	result := ExecuteBattle(attacker, defender)
	debris := CalculateDebrisField(result.AttackerRemaining, result.DefenderRemaining, spaceshipDefs)
	ApplyBattleResultToFleet(attackerFleet, result.AttackerRemaining)
	ApplyBattleResultToPlanet(defenderPlanet, result.DefenderRemaining)
*/

// BattleUnit represents a single unit in battle with its current stats
type BattleUnit struct {
	Type          string // spaceship or defense ID
	IsDefense     bool
	Count         uint64
	InitialCount  uint64 // Track initial count for debris calculation
	BaseArmor     uint64
	BaseShield    uint64
	BaseAttack    uint64
	CurrentShield []uint64 // Shield per unit
	CurrentHull   []uint64 // Hull integrity per unit
	Owner         string
}

// BattleParticipant represents a player/fleet in battle with their technology levels
type BattleParticipant struct {
	Owner               string
	WeaponTech          uint64
	ShieldingTech       uint64
	ArmorTech           uint64
	Units               []BattleUnit
	TotalInitialUnits   uint64
	TotalRemainingUnits uint64
}

// BattleResult contains the outcome of the battle
type BattleResult struct {
	Rounds            int
	AttackerWins      bool
	DefenderWins      bool
	Draw              bool
	AttackerRemaining *BattleParticipant
	DefenderRemaining *BattleParticipant
	DebrisField       OChainResources
}

// GetRapidFireAgainst returns the rapid fire value one unit type has against another
func GetRapidFireAgainst(attackerType string, targetType string) uint64 {
	rapidFireMap := map[string]map[string]uint64{
		// Small Cargo
		string(SmallCargoID): {},
		// Large Cargo
		string(LargeCargoID): {},
		// Light Fighter
		string(LightFighterID): {
			string(RecyclerID): 5,
		},
		// Heavy Fighter
		string(HeavyFighterID): {
			string(SmallCargoID): 3,
		},
		// Cruiser
		string(CruiserID): {
			string(LightFighterID):   6,
			string(RocketLauncherID): 10,
		},
		// Battleship
		string(BattleshipID): {},
		// Battlecruiser
		string(BattlecruiserID): {
			string(SmallCargoID):   3,
			string(LargeCargoID):   3,
			string(HeavyFighterID): 4,
			string(CruiserID):      4,
			string(BattleshipID):   7,
		},
		// Bomber
		string(BomberID): {
			string(RocketLauncherID): 20,
			string(LightLaserID):     20,
			string(HeavyLaserID):     10,
			string(IonCannonID):      10,
			string(GaussCannonID):    5,
			string(PlasmaTurretID):   5,
		},
		// Destroyer
		string(DestroyerID): {
			string(LightLaserID):    10,
			string(BattlecruiserID): 2,
		},
		// Deathstar
		string(DeathstarID): {
			string(SmallCargoID):     250,
			string(LargeCargoID):     250,
			string(LightFighterID):   200,
			string(HeavyFighterID):   100,
			string(CruiserID):        33,
			string(BattleshipID):     30,
			string(BattlecruiserID):  15,
			string(BomberID):         25,
			string(DestroyerID):      5,
			string(RocketLauncherID): 200,
			string(LightLaserID):     200,
			string(HeavyLaserID):     100,
			string(IonCannonID):      100,
			string(GaussCannonID):    50,
		},
		// Reaper
		string(ReaperID): {
			string(BattleshipID): 7,
			string(BomberID):     4,
			string(DestroyerID):  3,
		},
		// Recycler
		string(RecyclerID): {},
	}

	if targets, exists := rapidFireMap[attackerType]; exists {
		if rf, exists := targets[targetType]; exists {
			return rf
		}
	}
	return 1
}

// GetUnitBaseStats returns the base combat stats for a given unit
func GetUnitBaseStats(unitType string, isDefense bool, spaceships []OChainSpaceship, defenses []OChainDefense) (armor uint64, shield uint64, attack uint64) {
	if isDefense {
		for _, def := range defenses {
			if string(def.Id) == unitType {
				return def.FighterStats.Armor, def.FighterStats.Shield, def.FighterStats.Attack
			}
		}
	} else {
		for _, ship := range spaceships {
			if string(ship.Id) == unitType {
				return ship.FighterStats.Armor, ship.FighterStats.Shield, ship.FighterStats.Attack
			}
		}
	}
	return 0, 0, 0
}

// ApplyTechnologyBonus applies technology bonuses to base stats
func ApplyTechnologyBonus(baseValue uint64, techLevel uint64) uint64 {
	return baseValue + (baseValue * techLevel * 10 / 100)
}

// CreateBattleParticipant creates a battle participant from fleet/planet data
func CreateBattleParticipant(
	owner string,
	account OChainUniverseAccount,
	fleetShips *OChainFleetSpaceships,
	planetDefenses *OChainPlanetDefences,
	spaceshipDefs []OChainSpaceship,
	defenseDefs []OChainDefense,
) *BattleParticipant {
	participant := &BattleParticipant{
		Owner:         owner,
		WeaponTech:    account.Technologies.Weapon,
		ShieldingTech: account.Technologies.Shielding,
		ArmorTech:     account.Technologies.Armor,
		Units:         []BattleUnit{},
	}

	// Add spaceships
	if fleetShips != nil {
		spaceshipIds := GetSpaceshipIds()
		for _, id := range spaceshipIds {
			count := fleetShips.Get(id)
			if count > 0 {
				armor, shield, attack := GetUnitBaseStats(string(id), false, spaceshipDefs, defenseDefs)

				// Apply technology bonuses
				armor = ApplyTechnologyBonus(armor, participant.ArmorTech)
				shield = ApplyTechnologyBonus(shield, participant.ShieldingTech)
				attack = ApplyTechnologyBonus(attack, participant.WeaponTech)

				unit := BattleUnit{
					Type:          string(id),
					IsDefense:     false,
					Count:         count,
					InitialCount:  count,
					BaseArmor:     armor,
					BaseShield:    shield,
					BaseAttack:    attack,
					CurrentShield: make([]uint64, count),
					CurrentHull:   make([]uint64, count),
					Owner:         owner,
				}

				// Initialize current values
				for i := uint64(0); i < count; i++ {
					unit.CurrentShield[i] = shield
					unit.CurrentHull[i] = armor
				}

				participant.Units = append(participant.Units, unit)
				participant.TotalInitialUnits += count
				participant.TotalRemainingUnits += count
			}
		}
	}

	// Add defenses
	if planetDefenses != nil {
		defenseTypes := []OChainDefenseID{
			RocketLauncherID, LightLaserID, HeavyLaserID,
			IonCannonID, GaussCannonID, PlasmaTurretID, DarkMatterCanonID,
		}

		for _, id := range defenseTypes {
			var count uint64
			switch id {
			case RocketLauncherID:
				count = planetDefenses.RocketLauncher
			case LightLaserID:
				count = planetDefenses.LightLaser
			case HeavyLaserID:
				count = planetDefenses.HeavyLaser
			case IonCannonID:
				count = planetDefenses.IonCannon
			case GaussCannonID:
				count = planetDefenses.GaussCannon
			case PlasmaTurretID:
				count = planetDefenses.PlasmaTurret
			case DarkMatterCanonID:
				count = planetDefenses.DarkMatterCanon
			}

			if count > 0 {
				armor, shield, attack := GetUnitBaseStats(string(id), true, spaceshipDefs, defenseDefs)

				// Apply technology bonuses
				armor = ApplyTechnologyBonus(armor, participant.ArmorTech)
				shield = ApplyTechnologyBonus(shield, participant.ShieldingTech)
				attack = ApplyTechnologyBonus(attack, participant.WeaponTech)

				unit := BattleUnit{
					Type:          string(id),
					IsDefense:     true,
					Count:         count,
					InitialCount:  count,
					BaseArmor:     armor,
					BaseShield:    shield,
					BaseAttack:    attack,
					CurrentShield: make([]uint64, count),
					CurrentHull:   make([]uint64, count),
					Owner:         owner,
				}

				// Initialize current values
				for i := uint64(0); i < count; i++ {
					unit.CurrentShield[i] = shield
					unit.CurrentHull[i] = armor
				}

				participant.Units = append(participant.Units, unit)
				participant.TotalInitialUnits += count
				participant.TotalRemainingUnits += count
			}
		}
	}

	return participant
}

// ExecuteBattle runs the battle simulation
func ExecuteBattle(attacker *BattleParticipant, defender *BattleParticipant) *BattleResult {
	result := &BattleResult{
		Rounds: 0,
	}

	maxRounds := 6

	for round := 0; round < maxRounds; round++ {
		result.Rounds++

		// Check if battle should end
		attackerHasUnits := attacker.TotalRemainingUnits > 0
		defenderHasUnits := defender.TotalRemainingUnits > 0

		if !attackerHasUnits && !defenderHasUnits {
			result.Draw = true
			break
		}
		if !attackerHasUnits {
			result.DefenderWins = true
			break
		}
		if !defenderHasUnits {
			result.AttackerWins = true
			break
		}

		// Execute combat round
		executeCombatRound(attacker, defender)

		// Regenerate shields
		regenerateShields(attacker)
		regenerateShields(defender)

		// Remove destroyed units
		removeDestroyedUnits(attacker)
		removeDestroyedUnits(defender)
	}

	// If max rounds reached, check who has more units
	if result.Rounds >= maxRounds && !result.AttackerWins && !result.DefenderWins && !result.Draw {
		if attacker.TotalRemainingUnits > 0 && defender.TotalRemainingUnits > 0 {
			result.Draw = true
		} else if attacker.TotalRemainingUnits > 0 {
			result.AttackerWins = true
		} else if defender.TotalRemainingUnits > 0 {
			result.DefenderWins = true
		} else {
			result.Draw = true
		}
	}

	result.AttackerRemaining = attacker
	result.DefenderRemaining = defender

	return result
}

// executeCombatRound simulates one round of combat
func executeCombatRound(attacker *BattleParticipant, defender *BattleParticipant) {
	// Attackers fire at defenders
	for i := range attacker.Units {
		for unitIdx := uint64(0); unitIdx < attacker.Units[i].Count; unitIdx++ {
			// Check if this unit is still alive
			if attacker.Units[i].CurrentHull[unitIdx] == 0 {
				continue
			}

			// Fire with rapid fire
			rapidFireCount := uint64(1)
			for shot := uint64(0); shot < rapidFireCount; shot++ {
				target := selectRandomTarget(defender)
				if target == nil {
					break
				}

				damage := attacker.Units[i].BaseAttack
				applyDamage(target, damage)

				// Check for rapid fire
				rf := GetRapidFireAgainst(attacker.Units[i].Type, target.Type)
				if rf > 1 {
					// Rapid fire: (rf-1)/rf chance to fire again
					if rand.Float64() < float64(rf-1)/float64(rf) {
						rapidFireCount++
					}
				}
			}
		}
	}

	// Defenders fire at attackers
	for i := range defender.Units {
		for unitIdx := uint64(0); unitIdx < defender.Units[i].Count; unitIdx++ {
			// Check if this unit is still alive
			if defender.Units[i].CurrentHull[unitIdx] == 0 {
				continue
			}

			// Fire with rapid fire
			rapidFireCount := uint64(1)
			for shot := uint64(0); shot < rapidFireCount; shot++ {
				target := selectRandomTarget(attacker)
				if target == nil {
					break
				}

				damage := defender.Units[i].BaseAttack
				applyDamage(target, damage)

				// Check for rapid fire
				rf := GetRapidFireAgainst(defender.Units[i].Type, target.Type)
				if rf > 1 {
					// Rapid fire: (rf-1)/rf chance to fire again
					if rand.Float64() < float64(rf-1)/float64(rf) {
						rapidFireCount++
					}
				}
			}
		}
	}
}

// selectRandomTarget selects a random alive unit from the participant
func selectRandomTarget(participant *BattleParticipant) *BattleUnit {
	// Count alive units
	aliveUnits := []struct {
		unitIdx int
		idx     uint64
	}{}

	for i := range participant.Units {
		for j := uint64(0); j < participant.Units[i].Count; j++ {
			if participant.Units[i].CurrentHull[j] > 0 {
				aliveUnits = append(aliveUnits, struct {
					unitIdx int
					idx     uint64
				}{i, j})
			}
		}
	}

	if len(aliveUnits) == 0 {
		return nil
	}

	// Select random alive unit
	selected := aliveUnits[rand.Intn(len(aliveUnits))]
	return &participant.Units[selected.unitIdx]
}

// applyDamage applies damage to a unit
func applyDamage(target *BattleUnit, damage uint64) {
	// Find first alive unit in the target
	for i := uint64(0); i < target.Count; i++ {
		if target.CurrentHull[i] > 0 {
			// First, damage goes to shields
			if target.CurrentShield[i] > 0 {
				// OGame rule: if damage is less than 1% of shield, it bounces off
				minDamage := target.BaseShield / 100
				if damage < minDamage {
					return // Damage bounces off the shield
				}

				if damage > target.CurrentShield[i] {
					damage -= target.CurrentShield[i]
					target.CurrentShield[i] = 0
				} else {
					target.CurrentShield[i] -= damage
					return
				}
			}

			// Then damage goes to hull (hull = 10% of structure points in OGame)
			// In our implementation, BaseArmor represents the structure points
			// and CurrentHull represents the actual hull (which is BaseArmor in full health)
			if damage > target.CurrentHull[i] {
				target.CurrentHull[i] = 0
			} else {
				target.CurrentHull[i] -= damage

				// Check explosion probability following OGame rules:
				// - If hull damage < 30%, no chance of explosion
				// - If hull damage >= 30%, chance of explosion = damage percentage
				damagePercent := float64(target.BaseArmor-target.CurrentHull[i]) / float64(target.BaseArmor)

				// Only check for explosion if damage >= 30%
				if damagePercent >= 0.30 {
					// Explosion chance = damage percentage
					if rand.Float64() < damagePercent {
						// Unit explodes
						target.CurrentHull[i] = 0
					}
				}
			}
			return
		}
	}
}

// regenerateShields regenerates shields for units that weren't completely destroyed
func regenerateShields(participant *BattleParticipant) {
	for i := range participant.Units {
		for j := uint64(0); j < participant.Units[i].Count; j++ {
			if participant.Units[i].CurrentHull[j] > 0 {
				participant.Units[i].CurrentShield[j] = participant.Units[i].BaseShield
			}
		}
	}
}

// removeDestroyedUnits removes units with 0 hull
func removeDestroyedUnits(participant *BattleParticipant) {
	for i := range participant.Units {
		aliveCount := uint64(0)
		newShields := []uint64{}
		newHulls := []uint64{}

		for j := uint64(0); j < participant.Units[i].Count; j++ {
			if participant.Units[i].CurrentHull[j] > 0 {
				aliveCount++
				newShields = append(newShields, participant.Units[i].CurrentShield[j])
				newHulls = append(newHulls, participant.Units[i].CurrentHull[j])
			}
		}

		destroyed := participant.Units[i].Count - aliveCount
		participant.TotalRemainingUnits -= destroyed

		participant.Units[i].Count = aliveCount
		participant.Units[i].CurrentShield = newShields
		participant.Units[i].CurrentHull = newHulls
	}
}

// CalculateDebrisField calculates the debris field from destroyed units
func CalculateDebrisField(attacker *BattleParticipant, defender *BattleParticipant, spaceshipDefs []OChainSpaceship) OChainResources {
	debris := OChainResources{}

	calculateParticipantDebris := func(participant *BattleParticipant) {
		for i := range participant.Units {
			if participant.Units[i].IsDefense {
				// Defenses don't create debris in OGame
				continue
			}

			// Calculate destroyed count
			destroyed := participant.Units[i].InitialCount - participant.Units[i].Count

			if destroyed > 0 {
				// Find the spaceship definition
				for j := range spaceshipDefs {
					if string(spaceshipDefs[j].Id) == participant.Units[i].Type {
						// 30% of metal and crystal becomes debris per destroyed ship
						debris.Metal += spaceshipDefs[j].Cost.Metal * destroyed * 30 / 100
						debris.Crystal += spaceshipDefs[j].Cost.Crystal * destroyed * 30 / 100
						break
					}
				}
			}
		}
	}

	calculateParticipantDebris(attacker)
	calculateParticipantDebris(defender)

	return debris
}

// ApplyBattleResultToFleet updates the fleet with battle results
func ApplyBattleResultToFleet(fleet *OChainFleet, participant *BattleParticipant) {
	for i := range participant.Units {
		if !participant.Units[i].IsDefense {
			spaceshipId := OChainSpaceshipID(participant.Units[i].Type)

			// Get current count
			currentCount := fleet.GetSpaceships(spaceshipId)

			// Calculate difference
			if currentCount > participant.Units[i].Count {
				destroyed := currentCount - participant.Units[i].Count
				fleet.RemoveSpaceships(spaceshipId, destroyed)
			}
		}
	}
}

// ApplyBattleResultToPlanet updates the planet defenses with battle results
func ApplyBattleResultToPlanet(planet *OChainPlanet, participant *BattleParticipant) {
	// First update defenses
	for i := range participant.Units {
		if participant.Units[i].IsDefense {
			defenseId := OChainDefenseID(participant.Units[i].Type)

			var currentCount uint64
			switch defenseId {
			case RocketLauncherID:
				currentCount = planet.Defenses.RocketLauncher
			case LightLaserID:
				currentCount = planet.Defenses.LightLaser
			case HeavyLaserID:
				currentCount = planet.Defenses.HeavyLaser
			case IonCannonID:
				currentCount = planet.Defenses.IonCannon
			case GaussCannonID:
				currentCount = planet.Defenses.GaussCannon
			case PlasmaTurretID:
				currentCount = planet.Defenses.PlasmaTurret
			case DarkMatterCanonID:
				currentCount = planet.Defenses.DarkMatterCanon
			}

			// Calculate destroyed count
			if currentCount > participant.Units[i].Count {
				destroyed := currentCount - participant.Units[i].Count

				// In OGame, 70% of destroyed defenses can be rebuilt
				rebuilt := destroyed * 70 / 100
				actualDestroyed := destroyed - rebuilt

				planet.RemoveDefenses(defenseId, actualDestroyed)
			}
		}
	}

	// Update planet spaceships
	for i := range participant.Units {
		if !participant.Units[i].IsDefense {
			spaceshipId := OChainSpaceshipID(participant.Units[i].Type)

			currentCount := planet.Spaceships.Get(spaceshipId)

			if currentCount > participant.Units[i].Count {
				destroyed := currentCount - participant.Units[i].Count
				planet.RemoveSpaceships(spaceshipId, destroyed)
			}
		}
	}
}
