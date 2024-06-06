package main

import (
	"time"

	"github.com/ochain-gg/ochain-network/types"
)

var DefaultBuildings []types.OChainBuilding = []types.OChainBuilding{
	{
		Id:   types.MetalMineID,
		Name: "Metal Mine",
		BaseCost: types.OChainResources{
			OCT:       0,
			Metal:     200,
			Crystal:   100,
			Deuterium: 0,
		},
		Dependencies: []types.OChainDependency{},
	},
	{
		Id:   types.CrystalMineID,
		Name: "Crystal Mine",
		BaseCost: types.OChainResources{
			OCT:       0,
			Metal:     400,
			Crystal:   200,
			Deuterium: 0,
		},
		Dependencies: []types.OChainDependency{},
	},
	{
		Id:   types.DeuteriumMineID,
		Name: "Deuterium Synthesizer",
		BaseCost: types.OChainResources{
			OCT:       0,
			Metal:     500,
			Crystal:   250,
			Deuterium: 0,
		},
		Dependencies: []types.OChainDependency{},
	},
	{
		Id:   types.SolarPowerPlantID,
		Name: "Solar Power Plant",
		BaseCost: types.OChainResources{
			OCT:       0,
			Metal:     300,
			Crystal:   200,
			Deuterium: 0,
		},
		Dependencies: []types.OChainDependency{},
	},
	{
		Id:   types.RoboticFactoryID,
		Name: "Robotics Factory",
		BaseCost: types.OChainResources{
			OCT:       0,
			Metal:     1_000,
			Crystal:   500,
			Deuterium: 100,
		},
		Dependencies: []types.OChainDependency{},
	},
	{
		Id:   types.NaniteFactoryID,
		Name: "Nanite Factory",
		BaseCost: types.OChainResources{
			OCT:       0,
			Metal:     1_000_000,
			Crystal:   600_000,
			Deuterium: 500_000,
		},
		Dependencies: []types.OChainDependency{
			{
				DependencyType: types.OChainBuildingDependency,
				DependencyId:   string(types.RoboticFactoryID),
				Level:          10,
			},
			{
				DependencyType: types.OChainTechnologyDependency,
				DependencyId:   string(types.ComputerID),
				Level:          10,
			},
		},
	},
	{
		Id:   types.SpaceshipFactoryID,
		Name: "Spaceship Factory",
		BaseCost: types.OChainResources{
			OCT:       0,
			Metal:     1_000,
			Crystal:   500,
			Deuterium: 100,
		},
		Dependencies: []types.OChainDependency{
			{
				DependencyType: types.OChainBuildingDependency,
				DependencyId:   string(types.RoboticFactoryID),
				Level:          2,
			},
			{
				DependencyType: types.OChainTechnologyDependency,
				DependencyId:   string(types.ComputerID),
				Level:          2,
			},
		},
	},
	{
		Id:   types.IntergalacticPortalID,
		Name: "Intergalactic Portal",
		BaseCost: types.OChainResources{
			OCT:       0,
			Metal:     5_000_000,
			Crystal:   3_000_000,
			Deuterium: 3_000_000,
		},
		Dependencies: []types.OChainDependency{
			{
				DependencyType: types.OChainBuildingDependency,
				DependencyId:   string(types.SolarPowerPlantID),
				Level:          20,
			},
			{
				DependencyType: types.OChainTechnologyDependency,
				DependencyId:   string(types.GravitonID),
				Level:          1,
			},
			{
				DependencyType: types.OChainTechnologyDependency,
				DependencyId:   string(types.HyperspaceID),
				Level:          10,
			},
		},
	},
	{
		Id:   types.ResearchLaboratoryID,
		Name: "Research Laboratory",
		BaseCost: types.OChainResources{
			OCT:       0,
			Metal:     1_000,
			Crystal:   100,
			Deuterium: 600,
		},
		Dependencies: []types.OChainDependency{},
	},
	{
		Id:   types.ShieldDomeID,
		Name: "Shield Dome",
		BaseCost: types.OChainResources{
			OCT:       0,
			Metal:     1_000,
			Crystal:   500,
			Deuterium: 100,
		},
		Dependencies: []types.OChainDependency{},
	},
}

var DefaultTechnologies []types.OChainTechnology = []types.OChainTechnology{
	{
		Id:   types.ComputerID,
		Name: "Computer",
		BaseCost: types.OChainResources{
			OCT:       0,
			Metal:     200,
			Crystal:   100,
			Deuterium: 0,
		},
		Dependencies: []types.OChainDependency{
			{
				DependencyType: types.OChainBuildingDependency,
				DependencyId:   string(types.ResearchLaboratoryID),
				Level:          1,
			},
		},
	},
	{
		Id:   types.WeaponID,
		Name: "Weapon",
		BaseCost: types.OChainResources{
			OCT:       0,
			Metal:     400,
			Crystal:   200,
			Deuterium: 0,
		},
		Dependencies: []types.OChainDependency{
			{
				DependencyType: types.OChainBuildingDependency,
				DependencyId:   string(types.ResearchLaboratoryID),
				Level:          4,
			},
		},
	},
	{
		Id:   types.ShieldingID,
		Name: "Shielding",
		BaseCost: types.OChainResources{
			OCT:       0,
			Metal:     500,
			Crystal:   250,
			Deuterium: 0,
		},
		Dependencies: []types.OChainDependency{
			{
				DependencyType: types.OChainBuildingDependency,
				DependencyId:   string(types.ResearchLaboratoryID),
				Level:          6,
			},
			{
				DependencyType: types.OChainTechnologyDependency,
				DependencyId:   string(types.EnergyID),
				Level:          3,
			},
		},
	},
	{
		Id:   types.ArmorID,
		Name: "Armor",
		BaseCost: types.OChainResources{
			OCT:       0,
			Metal:     300,
			Crystal:   200,
			Deuterium: 0,
		},
		Dependencies: []types.OChainDependency{
			{
				DependencyType: types.OChainBuildingDependency,
				DependencyId:   string(types.ResearchLaboratoryID),
				Level:          2,
			},
		},
	},
	{
		Id:   types.EnergyID,
		Name: "Energy",
		BaseCost: types.OChainResources{
			OCT:       0,
			Metal:     100,
			Crystal:   50,
			Deuterium: 10,
		},
		Dependencies: []types.OChainDependency{
			{
				DependencyType: types.OChainBuildingDependency,
				DependencyId:   string(types.ResearchLaboratoryID),
				Level:          1,
			},
		},
	},
	{
		Id:   types.CombustionDriveID,
		Name: "Combustion Drive",
		BaseCost: types.OChainResources{
			OCT:       0,
			Metal:     500,
			Crystal:   250,
			Deuterium: 25,
		},
		Dependencies: []types.OChainDependency{
			{
				DependencyType: types.OChainBuildingDependency,
				DependencyId:   string(types.ResearchLaboratoryID),
				Level:          1,
			},
			{
				DependencyType: types.OChainTechnologyDependency,
				DependencyId:   string(types.EnergyID),
				Level:          1,
			},
		},
	},
	{
		Id:   types.ImpulseDriveID,
		Name: "Impulse Drive",
		BaseCost: types.OChainResources{
			OCT:       0,
			Metal:     1_00,
			Crystal:   500,
			Deuterium: 100,
		},
		Dependencies: []types.OChainDependency{
			{
				DependencyType: types.OChainBuildingDependency,
				DependencyId:   string(types.ResearchLaboratoryID),
				Level:          2,
			},
			{
				DependencyType: types.OChainTechnologyDependency,
				DependencyId:   string(types.EnergyID),
				Level:          1,
			},
		},
	},
	{
		Id:   types.HyperspaceDriveID,
		Name: "Hyperspace Drive",
		BaseCost: types.OChainResources{
			OCT:       0,
			Metal:     10_000,
			Crystal:   20_000,
			Deuterium: 6_000,
		},
		Dependencies: []types.OChainDependency{
			{
				DependencyType: types.OChainBuildingDependency,
				DependencyId:   string(types.ResearchLaboratoryID),
				Level:          4,
			},
			{
				DependencyType: types.OChainTechnologyDependency,
				DependencyId:   string(types.HyperspaceID),
				Level:          3,
			},
		},
	},
	{
		Id:   types.HyperspaceID,
		Name: "Hyperspace",
		BaseCost: types.OChainResources{
			OCT:       0,
			Metal:     0,
			Crystal:   4_000,
			Deuterium: 2_000,
		},
		Dependencies: []types.OChainDependency{
			{
				DependencyType: types.OChainBuildingDependency,
				DependencyId:   string(types.ResearchLaboratoryID),
				Level:          7,
			},
			{
				DependencyType: types.OChainTechnologyDependency,
				DependencyId:   string(types.EnergyID),
				Level:          5,
			},
			{
				DependencyType: types.OChainTechnologyDependency,
				DependencyId:   string(types.ShieldingID),
				Level:          5,
			},
		},
	},
	{
		Id:   types.LaserID,
		Name: "Laser",
		BaseCost: types.OChainResources{
			OCT:       0,
			Metal:     200,
			Crystal:   100,
			Deuterium: 0,
		},
		Dependencies: []types.OChainDependency{
			{
				DependencyType: types.OChainBuildingDependency,
				DependencyId:   string(types.ResearchLaboratoryID),
				Level:          1,
			},
			{
				DependencyType: types.OChainTechnologyDependency,
				DependencyId:   string(types.EnergyID),
				Level:          2,
			},
		},
	},
	{
		Id:   types.IonID,
		Name: "Ion",
		BaseCost: types.OChainResources{
			OCT:       0,
			Metal:     1_000,
			Crystal:   300,
			Deuterium: 100,
		},
		Dependencies: []types.OChainDependency{
			{
				DependencyType: types.OChainBuildingDependency,
				DependencyId:   string(types.ResearchLaboratoryID),
				Level:          4,
			},
			{
				DependencyType: types.OChainTechnologyDependency,
				DependencyId:   string(types.EnergyID),
				Level:          4,
			},
			{
				DependencyType: types.OChainTechnologyDependency,
				DependencyId:   string(types.LaserID),
				Level:          5,
			},
		},
	},
	{
		Id:   types.PlasmaID,
		Name: "Plasma",
		BaseCost: types.OChainResources{
			OCT:       0,
			Metal:     2_000,
			Crystal:   4_000,
			Deuterium: 1_000,
		},
		Dependencies: []types.OChainDependency{
			{
				DependencyType: types.OChainBuildingDependency,
				DependencyId:   string(types.ResearchLaboratoryID),
				Level:          4,
			},
			{
				DependencyType: types.OChainTechnologyDependency,
				DependencyId:   string(types.EnergyID),
				Level:          8,
			},
			{
				DependencyType: types.OChainTechnologyDependency,
				DependencyId:   string(types.LaserID),
				Level:          10,
			},
			{
				DependencyType: types.OChainTechnologyDependency,
				DependencyId:   string(types.IonID),
				Level:          5,
			},
		},
	},
	{
		Id:   types.IntergalacticResearchNetworkID,
		Name: "Intergalactic Research Network",
		BaseCost: types.OChainResources{
			OCT:       0,
			Metal:     240_000,
			Crystal:   400_000,
			Deuterium: 160_000,
		},
		Dependencies: []types.OChainDependency{
			{
				DependencyType: types.OChainBuildingDependency,
				DependencyId:   string(types.ResearchLaboratoryID),
				Level:          10,
			},
			{
				DependencyType: types.OChainTechnologyDependency,
				DependencyId:   string(types.ComputerID),
				Level:          8,
			},
			{
				DependencyType: types.OChainTechnologyDependency,
				DependencyId:   string(types.HyperspaceID),
				Level:          8,
			},
		},
	},
	{
		Id:   types.AstrophysicsID,
		Name: "Astrophysics",
		BaseCost: types.OChainResources{
			OCT:       0,
			Metal:     4_000,
			Crystal:   8_000,
			Deuterium: 4_000,
		},
		Dependencies: []types.OChainDependency{
			{
				DependencyType: types.OChainBuildingDependency,
				DependencyId:   string(types.ResearchLaboratoryID),
				Level:          10,
			},
			{
				DependencyType: types.OChainTechnologyDependency,
				DependencyId:   string(types.ComputerID),
				Level:          8,
			},
			{
				DependencyType: types.OChainTechnologyDependency,
				DependencyId:   string(types.HyperspaceID),
				Level:          8,
			},
		},
	},
	{
		Id:   types.GravitonID,
		Name: "Graviton",
		BaseCost: types.OChainResources{
			OCT:       0,
			Metal:     0,
			Crystal:   0,
			Deuterium: 0,
		},
		Dependencies: []types.OChainDependency{
			{
				DependencyType: types.OChainBuildingDependency,
				DependencyId:   string(types.ResearchLaboratoryID),
				Level:          12,
			},
		},
	},
}

var DefaultSpaceships []types.OChainSpaceship = []types.OChainSpaceship{
	{
		Id:   types.SmallCargoID,
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
			Deuterium: 0,
		},
		Dependencies: []types.OChainDependency{
			{
				DependencyType: types.OChainBuildingDependency,
				DependencyId:   string(types.SpaceshipFactoryID),
				Level:          2,
			},
			{
				DependencyType: types.OChainTechnologyDependency,
				DependencyId:   string(types.CombustionDriveID),
				Level:          2,
			},
		},
	},
	{
		Id:   types.LargeCargoID,
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
			Deuterium: 0,
		},
		Dependencies: []types.OChainDependency{
			{
				DependencyType: types.OChainBuildingDependency,
				DependencyId:   string(types.SpaceshipFactoryID),
				Level:          6,
			},
			{
				DependencyType: types.OChainTechnologyDependency,
				DependencyId:   string(types.CombustionDriveID),
				Level:          4,
			},
		},
	},
	{
		Id:   types.LightFighterID,
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
			Deuterium: 0,
		},
		Dependencies: []types.OChainDependency{
			{
				DependencyType: types.OChainBuildingDependency,
				DependencyId:   string(types.SpaceshipFactoryID),
				Level:          1,
			},
			{
				DependencyType: types.OChainTechnologyDependency,
				DependencyId:   string(types.CombustionDriveID),
				Level:          1,
			},
		},
	},
	{
		Id:   types.HeavyFighterID,
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
			Deuterium: 0,
		},
		Dependencies: []types.OChainDependency{
			{
				DependencyType: types.OChainBuildingDependency,
				DependencyId:   string(types.SpaceshipFactoryID),
				Level:          3,
			},
			{
				DependencyType: types.OChainTechnologyDependency,
				DependencyId:   string(types.ImpulseDriveID),
				Level:          2,
			},
			{
				DependencyType: types.OChainTechnologyDependency,
				DependencyId:   string(types.ArmorID),
				Level:          2,
			},
		},
	},
	{
		Id:   types.CruiserID,
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
			Deuterium: 2_000,
		},
		Dependencies: []types.OChainDependency{
			{
				DependencyType: types.OChainBuildingDependency,
				DependencyId:   string(types.SpaceshipFactoryID),
				Level:          5,
			},
			{
				DependencyType: types.OChainTechnologyDependency,
				DependencyId:   string(types.ImpulseDriveID),
				Level:          4,
			},
			{
				DependencyType: types.OChainTechnologyDependency,
				DependencyId:   string(types.IonID),
				Level:          2,
			},
		},
	},
	{
		Id:   types.BattleshipID,
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
			Deuterium: 0,
		},
		Dependencies: []types.OChainDependency{
			{
				DependencyType: types.OChainBuildingDependency,
				DependencyId:   string(types.SpaceshipFactoryID),
				Level:          4,
			},
			{
				DependencyType: types.OChainTechnologyDependency,
				DependencyId:   string(types.HyperspaceDriveID),
				Level:          5,
			},
		},
	},
	{
		Id:   types.BattlecruiserID,
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
			Deuterium: 15_000,
		},
		Dependencies: []types.OChainDependency{
			{
				DependencyType: types.OChainBuildingDependency,
				DependencyId:   string(types.SpaceshipFactoryID),
				Level:          8,
			},
			{
				DependencyType: types.OChainTechnologyDependency,
				DependencyId:   string(types.HyperspaceID),
				Level:          5,
			},
			{
				DependencyType: types.OChainTechnologyDependency,
				DependencyId:   string(types.HyperspaceDriveID),
				Level:          5,
			},
			{
				DependencyType: types.OChainTechnologyDependency,
				DependencyId:   string(types.LaserID),
				Level:          12,
			},
		},
	},
	{
		Id:   types.BomberID,
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
			Deuterium: 15_000,
		},
		Dependencies: []types.OChainDependency{
			{
				DependencyType: types.OChainBuildingDependency,
				DependencyId:   string(types.SpaceshipFactoryID),
				Level:          8,
			},
			{
				DependencyType: types.OChainTechnologyDependency,
				DependencyId:   string(types.ImpulseDriveID),
				Level:          6,
			},
			{
				DependencyType: types.OChainTechnologyDependency,
				DependencyId:   string(types.PlasmaID),
				Level:          5,
			},
		},
	},
	{
		Id:   types.DestroyerID,
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
			Deuterium: 15_000,
		},
		Dependencies: []types.OChainDependency{
			{
				DependencyType: types.OChainBuildingDependency,
				DependencyId:   string(types.SpaceshipFactoryID),
				Level:          9,
			},
			{
				DependencyType: types.OChainTechnologyDependency,
				DependencyId:   string(types.HyperspaceID),
				Level:          5,
			},
			{
				DependencyType: types.OChainTechnologyDependency,
				DependencyId:   string(types.HyperspaceDriveID),
				Level:          6,
			},
		},
	},
	{
		Id:   types.DeathstarID,
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
			Deuterium: 1_000_000,
		},
		Dependencies: []types.OChainDependency{
			{
				DependencyType: types.OChainBuildingDependency,
				DependencyId:   string(types.SpaceshipFactoryID),
				Level:          12,
			},
			{
				DependencyType: types.OChainTechnologyDependency,
				DependencyId:   string(types.HyperspaceID),
				Level:          6,
			},
			{
				DependencyType: types.OChainTechnologyDependency,
				DependencyId:   string(types.HyperspaceDriveID),
				Level:          7,
			},
			{
				DependencyType: types.OChainTechnologyDependency,
				DependencyId:   string(types.GravitonID),
				Level:          1,
			},
		},
	},
	{
		Id:   types.ReaperID,
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
			Deuterium: 20_000,
		},
		Dependencies: []types.OChainDependency{
			{
				DependencyType: types.OChainBuildingDependency,
				DependencyId:   string(types.SpaceshipFactoryID),
				Level:          10,
			},
			{
				DependencyType: types.OChainBuildingDependency,
				DependencyId:   string(types.ResearchLaboratoryID),
				Level:          7,
			},
			{
				DependencyType: types.OChainTechnologyDependency,
				DependencyId:   string(types.HyperspaceID),
				Level:          6,
			},
			{
				DependencyType: types.OChainTechnologyDependency,
				DependencyId:   string(types.HyperspaceDriveID),
				Level:          7,
			},
			{
				DependencyType: types.OChainTechnologyDependency,
				DependencyId:   string(types.ShieldingID),
				Level:          6,
			},
			{
				DependencyType: types.OChainTechnologyDependency,
				DependencyId:   string(types.EnergyID),
				Level:          5,
			},
		},
	},
	{
		Id:   types.RecyclerID,
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
			Deuterium: 2_000_000,
		},

		Dependencies: []types.OChainDependency{
			{
				DependencyType: types.OChainBuildingDependency,
				DependencyId:   string(types.SpaceshipFactoryID),
				Level:          4,
			},
			{
				DependencyType: types.OChainTechnologyDependency,
				DependencyId:   string(types.CombustionDriveID),
				Level:          2,
			},
			{
				DependencyType: types.OChainTechnologyDependency,
				DependencyId:   string(types.ShieldingID),
				Level:          6,
			},
		},
	},
}

var DefaultDefenses []types.OChainDefense = []types.OChainDefense{
	{
		Id:   types.RocketLauncherID,
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
			Deuterium: 0,
		},
		Dependencies: []types.OChainDependency{
			{
				DependencyType: types.OChainBuildingDependency,
				DependencyId:   string(types.SpaceshipFactoryID),
				Level:          1,
			},
		},
	},
	{
		Id:   types.LightLaserID,
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
			Deuterium: 0,
		},
		Dependencies: []types.OChainDependency{
			{
				DependencyType: types.OChainBuildingDependency,
				DependencyId:   string(types.SpaceshipFactoryID),
				Level:          2,
			},
			{
				DependencyType: types.OChainTechnologyDependency,
				DependencyId:   string(types.EnergyID),
				Level:          1,
			},
			{
				DependencyType: types.OChainTechnologyDependency,
				DependencyId:   string(types.LaserID),
				Level:          3,
			},
		},
	},
	{
		Id:   types.HeavyLaserID,
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
			Deuterium: 0,
		},
		Dependencies: []types.OChainDependency{
			{
				DependencyType: types.OChainBuildingDependency,
				DependencyId:   string(types.SpaceshipFactoryID),
				Level:          4,
			},
			{
				DependencyType: types.OChainTechnologyDependency,
				DependencyId:   string(types.EnergyID),
				Level:          3,
			},
			{
				DependencyType: types.OChainTechnologyDependency,
				DependencyId:   string(types.LaserID),
				Level:          6,
			},
		},
	},
	{
		Id:   types.IonCannonID,
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
			Deuterium: 0,
		},
		Dependencies: []types.OChainDependency{
			{
				DependencyType: types.OChainBuildingDependency,
				DependencyId:   string(types.SpaceshipFactoryID),
				Level:          4,
			},
			{
				DependencyType: types.OChainTechnologyDependency,
				DependencyId:   string(types.IonID),
				Level:          4,
			},
		},
	},
	{
		Id:   types.GaussCannonID,
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
			Deuterium: 1_000,
		},
		Dependencies: []types.OChainDependency{
			{
				DependencyType: types.OChainBuildingDependency,
				DependencyId:   string(types.SpaceshipFactoryID),
				Level:          6,
			},
			{
				DependencyType: types.OChainTechnologyDependency,
				DependencyId:   string(types.EnergyID),
				Level:          6,
			},
			{
				DependencyType: types.OChainTechnologyDependency,
				DependencyId:   string(types.WeaponID),
				Level:          3,
			},
			{
				DependencyType: types.OChainTechnologyDependency,
				DependencyId:   string(types.ShieldingID),
				Level:          1,
			},
		},
	},
	{
		Id:   types.PlasmaTurretID,
		Name: "Plasma turret",
		FighterStats: types.OChainFighterStats{
			Armor:  100_000,
			Shield: 300,
			Attack: 10_000,
		},
		Cost: types.OChainResources{
			OCT:       0,
			Metal:     20_000,
			Crystal:   20_000,
			Deuterium: 5000,
		},
		Dependencies: []types.OChainDependency{
			{
				DependencyType: types.OChainBuildingDependency,
				DependencyId:   string(types.SpaceshipFactoryID),
				Level:          7,
			},
			{
				DependencyType: types.OChainTechnologyDependency,
				DependencyId:   string(types.PlasmaID),
				Level:          7,
			},
		},
	},
	{
		Id:   types.DarkMatterCanonID,
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
			Deuterium: 5_000_000,
		},
		Dependencies: []types.OChainDependency{
			{
				DependencyType: types.OChainBuildingDependency,
				DependencyId:   string(types.SpaceshipFactoryID),
				Level:          12,
			},
			{
				DependencyType: types.OChainBuildingDependency,
				DependencyId:   string(types.SolarPowerPlantID),
				Level:          20,
			},
			{
				DependencyType: types.OChainTechnologyDependency,
				DependencyId:   string(types.GravitonID),
				Level:          1,
			},
		},
	},
}

var DefaultUniverse types.OChainUniverse = types.OChainUniverse{
	Id:                      "main",
	Name:                    "OChain main universe",
	Speed:                   1,
	MaxGalaxy:               5,
	MaxSolarSystemPerGalaxy: 255,
	MaxPlanetPerSolarSystem: 16,
	Accounts:                0,
	ColonizedPlanets:        0,
	CreatedAt:               uint64(time.Now().Unix()),
}
