package main

import (
	"github.com/ochain.gg/ochain-network-validator/database"
	"github.com/ochain.gg/ochain-network-validator/types"
)

var DefaultUniverse database.OChainUniverse = database.OChainUniverse{
	Id:        1,
	Name:      "OChain main universe",
	CreatedAt: 0,
	Configuration: database.OChainUniverseConfiguration{
		Speed:                   2,
		MaxGalaxy:               255,
		MaxSolarSystemPerGalaxy: 255,
		MaxPlanetPerSolarSystem: 16,

		Spaceships: []types.OChainSpaceship{
			{
				Id:   1,
				Name: "Small cargo",
				FighterStats: types.OChainFighterStats{
					Armor:  4_000,
					Shield: 10,
					Attack: 5,
				},
				Stats: types.OChainSpaceshipStats{
					Capacity:   5_000,
					Speed:      5_000,
					Consumtion: 10,
				},
				Cost: types.OChainResources{
					OCT:       0,
					Metal:     2_000,
					Crystal:   2_000,
					Deutereum: 0,
				},
			},
			{
				Id:   2,
				Name: "Large cargo",
				FighterStats: types.OChainFighterStats{
					Armor:  12_000,
					Shield: 25,
					Attack: 5,
				},
				Stats: types.OChainSpaceshipStats{
					Capacity:   25_000,
					Speed:      7_500,
					Consumtion: 10,
				},
				Cost: types.OChainResources{
					OCT:       0,
					Metal:     6_000,
					Crystal:   6_000,
					Deutereum: 0,
				},
			},
			{
				Id:   3,
				Name: "Light figher",
				FighterStats: types.OChainFighterStats{
					Armor:  4_000,
					Shield: 10,
					Attack: 50,
				},
				Stats: types.OChainSpaceshipStats{
					Capacity:   200,
					Speed:      12_500,
					Consumtion: 20,
				},
				Cost: types.OChainResources{
					OCT:       0,
					Metal:     3_000,
					Crystal:   1_000,
					Deutereum: 0,
				},
			},
			{
				Id:   4,
				Name: "Heavy figher",
				FighterStats: types.OChainFighterStats{
					Armor:  10_000,
					Shield: 25,
					Attack: 150,
				},
				Stats: types.OChainSpaceshipStats{
					Capacity:   400,
					Speed:      10_000,
					Consumtion: 75,
				},
				Cost: types.OChainResources{
					OCT:       0,
					Metal:     6_000,
					Crystal:   4_000,
					Deutereum: 0,
				},
			},
			{
				Id:   5,
				Name: "Cruiser",
				FighterStats: types.OChainFighterStats{
					Armor:  27_000,
					Shield: 50,
					Attack: 400,
				},
				Stats: types.OChainSpaceshipStats{
					Capacity:   800,
					Speed:      15_000,
					Consumtion: 300,
				},
				Cost: types.OChainResources{
					OCT:       0,
					Metal:     20_000,
					Crystal:   7_000,
					Deutereum: 2_000,
				},
			},
			{
				Id:   6,
				Name: "Battleship",
				FighterStats: types.OChainFighterStats{
					Armor:  60_000,
					Shield: 200,
					Attack: 1_000,
				},
				Stats: types.OChainSpaceshipStats{
					Capacity:   1_500,
					Speed:      10_000,
					Consumtion: 500,
				},
				Cost: types.OChainResources{
					OCT:       0,
					Metal:     45_000,
					Crystal:   15_000,
					Deutereum: 0,
				},
			},
			{
				Id:   7,
				Name: "BattleCruiser",
				FighterStats: types.OChainFighterStats{
					Armor:  70_000,
					Shield: 400,
					Attack: 700,
				},
				Stats: types.OChainSpaceshipStats{
					Capacity:   750,
					Speed:      10_000,
					Consumtion: 350,
				},
				Cost: types.OChainResources{
					OCT:       0,
					Metal:     30_000,
					Crystal:   40_000,
					Deutereum: 15_000,
				},
			},
			{
				Id:   8,
				Name: "Bomber",
				FighterStats: types.OChainFighterStats{
					Armor:  75_000,
					Shield: 500,
					Attack: 1_000,
				},
				Stats: types.OChainSpaceshipStats{
					Capacity:   500,
					Speed:      4_000,
					Consumtion: 700,
				},
				Cost: types.OChainResources{
					OCT:       0,
					Metal:     50_000,
					Crystal:   25_000,
					Deutereum: 15_000,
				},
			},
			{
				Id:   9,
				Name: "Destroyer",
				FighterStats: types.OChainFighterStats{
					Armor:  110_000,
					Shield: 500,
					Attack: 2_000,
				},
				Stats: types.OChainSpaceshipStats{
					Capacity:   2_000,
					Speed:      5_000,
					Consumtion: 1_000,
				},
				Cost: types.OChainResources{
					OCT:       0,
					Metal:     60_000,
					Crystal:   50_000,
					Deutereum: 15_000,
				},
			},
			{
				Id:   10,
				Name: "Destroyer",
				FighterStats: types.OChainFighterStats{
					Armor:  110_000,
					Shield: 500,
					Attack: 2_000,
				},
				Stats: types.OChainSpaceshipStats{
					Capacity:   2_000,
					Speed:      5_000,
					Consumtion: 1_000,
				},
				Cost: types.OChainResources{
					OCT:       0,
					Metal:     60_000,
					Crystal:   50_000,
					Deutereum: 15_000,
				},
			},
			{
				Id:   11,
				Name: "Deathstar",
				FighterStats: types.OChainFighterStats{
					Armor:  9_000_000,
					Shield: 50_000,
					Attack: 200_000,
				},
				Stats: types.OChainSpaceshipStats{
					Capacity:   1_000_000,
					Speed:      100,
					Consumtion: 1,
				},
				Cost: types.OChainResources{
					OCT:       0,
					Metal:     5_000_000,
					Crystal:   4_000_000,
					Deutereum: 1_000_000,
				},
			},
			{
				Id:   12,
				Name: "Reaper",
				FighterStats: types.OChainFighterStats{
					Armor:  140_000,
					Shield: 700,
					Attack: 2_800,
				},
				Stats: types.OChainSpaceshipStats{
					Capacity:   7_000,
					Speed:      7_000,
					Consumtion: 900,
				},
				Cost: types.OChainResources{
					OCT:       0,
					Metal:     85_000,
					Crystal:   55_000,
					Deutereum: 20_000,
				},
			},
			{
				Id:   13,
				Name: "Recycler",
				FighterStats: types.OChainFighterStats{
					Armor:  16_000,
					Shield: 10,
					Attack: 1,
				},
				Stats: types.OChainSpaceshipStats{
					Capacity:   20_000,
					Speed:      2_000,
					Consumtion: 300,
				},
				Cost: types.OChainResources{
					OCT:       0,
					Metal:     10_000_000,
					Crystal:   6_000_000,
					Deutereum: 2_000_000,
				},
			},
		},

		Defenses: []types.OChainDefense{
			{
				Id:   0,
				Name: "Rocket launcher",
				FighterStats: types.OChainFighterStats{
					Armor:  2_000,
					Shield: 20,
					Attack: 80,
				},
				Cost: types.OChainResources{
					OCT:       0,
					Metal:     2_000,
					Crystal:   1_000,
					Deutereum: 0,
				},
			},
			{
				Id:   1,
				Name: "Light laser",
				FighterStats: types.OChainFighterStats{
					Armor:  2000,
					Shield: 100,
					Attack: 250,
				},
				Cost: types.OChainResources{
					OCT:       0,
					Metal:     4_000,
					Crystal:   2_000,
					Deutereum: 0,
				},
			},
			{
				Id:   1,
				Name: "Heavy laser",
				FighterStats: types.OChainFighterStats{
					Armor:  2000,
					Shield: 100,
					Attack: 250,
				},
				Cost: types.OChainResources{
					OCT:       0,
					Metal:     5_000,
					Crystal:   2_500,
					Deutereum: 0,
				},
			},
			{
				Id:   1,
				Name: "Ion canon",
				FighterStats: types.OChainFighterStats{
					Armor:  8000,
					Shield: 500,
					Attack: 150,
				},
				Cost: types.OChainResources{
					OCT:       0,
					Metal:     3_000,
					Crystal:   2_000,
					Deutereum: 0,
				},
			},
			{
				Id:   1,
				Name: "Gauss canon",
				FighterStats: types.OChainFighterStats{
					Armor:  35000,
					Shield: 200,
					Attack: 1100,
				},
				Cost: types.OChainResources{
					OCT:       0,
					Metal:     10_000,
					Crystal:   5_000,
					Deutereum: 1_000,
				},
			},
			{
				Id:   1,
				Name: "Plasma turret",
				FighterStats: types.OChainFighterStats{
					Armor:  100000,
					Shield: 300,
					Attack: 10000,
				},
				Cost: types.OChainResources{
					OCT:       0,
					Metal:     20_000,
					Crystal:   20_000,
					Deutereum: 5000,
				},
			},
			{
				Id:   1,
				Name: "Dark matter canon",
				FighterStats: types.OChainFighterStats{
					Armor:  500000,
					Shield: 1000,
					Attack: 100000,
				},
				Cost: types.OChainResources{
					OCT:       0,
					Metal:     10_000_000,
					Crystal:   6_000_000,
					Deutereum: 5_000_000,
				},
			},
		},
	},
}
