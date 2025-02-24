package account_transactions

import (
	"encoding/hex"
	"fmt"
	"math/big"

	abcitypes "github.com/cometbft/cometbft/abci/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/fxamacker/cbor/v2"

	t "github.com/ochain-gg/ochain-network/engine/transactions"
	"github.com/ochain-gg/ochain-network/types"
)

type RegisterUniverseAccountTransactionData struct {
	UniverseId string `cbor:"1,keyasint"`
}

type RegisterUniverseAccountTransaction struct {
	Type      t.TransactionType                      `cbor:"1,keyasint"`
	From      string                                 `cbor:"2,keyasint"`
	Nonce     uint64                                 `cbor:"3,keyasint"`
	Data      RegisterUniverseAccountTransactionData `cbor:"4,keyasint"`
	Signature []byte                                 `cbor:"5,keyasint"`
}

func (tx *RegisterUniverseAccountTransaction) Transaction() (t.Transaction, error) {
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

func (tx *RegisterUniverseAccountTransaction) Check(ctx t.TransactionContext) *abcitypes.CheckTxResponse {
	_, err := ctx.Db.GlobalsAccounts.Get(tx.From)
	if err != nil {
		return &abcitypes.CheckTxResponse{
			Code: types.InvalidTransactionError,
		}
	}

	universeExists, err := ctx.Db.Universes.Exists(tx.Data.UniverseId)
	if err != nil {
		return &abcitypes.CheckTxResponse{
			Code: types.InvalidTransactionError,
		}
	}

	if !universeExists {
		return &abcitypes.CheckTxResponse{
			Code: types.InvalidTransactionError,
		}
	}

	exists, err := ctx.Db.UniverseAccounts.Exists(tx.Data.UniverseId, tx.From)
	if err != nil {
		return &abcitypes.CheckTxResponse{
			Code: types.InvalidTransactionError,
		}
	}

	if exists {
		return &abcitypes.CheckTxResponse{
			Code: types.InvalidTransactionError,
		}
	}

	return &abcitypes.CheckTxResponse{
		Code: types.NoError,
	}
}

func (tx *RegisterUniverseAccountTransaction) Execute(ctx t.TransactionContext) *abcitypes.ExecTxResult {
	result := tx.Check(ctx)
	if result.Code != types.NoError {
		return &abcitypes.ExecTxResult{
			Code: result.GetCode(),
		}
	}

	globalAccount, err := ctx.Db.GlobalsAccounts.Get(tx.From)
	if err != nil {
		return &abcitypes.ExecTxResult{
			Code: types.InvalidTransactionError,
		}
	}

	universe, err := ctx.Db.Universes.Get(tx.Data.UniverseId)
	if err != nil {
		return &abcitypes.ExecTxResult{
			Code: types.InvalidTransactionError,
		}
	}

	account := types.OChainUniverseAccount{
		Address:    tx.From,
		UniverseId: tx.Data.UniverseId,
		Points:     0,
		Technologies: types.OChainAccountTechnologies{
			Computer:  0,
			Weapon:    0,
			Shielding: 0,
			Armor:     0,
			Energy:    0,

			CombustionDrive: 0,
			ImpulseDrive:    0,
			HyperspaceDrive: 0,

			Hyperspace: 0,
			Laser:      0,
			Ion:        0,
			Plasma:     0,

			IntergalacticResearchNetwork: 0,
			Astrophysics:                 0,
			Graviton:                     0,
		},
		CreatedAt: ctx.Date.Unix(),
	}

	//create a first planet with random coordinate in a deterministic way based on the transaction hash
	maxGalaxy := new(big.Int)
	maxGalaxy.SetUint64(uint64(universe.MaxGalaxy))

	maxSolarSystem := new(big.Int)
	maxSolarSystem.SetUint64(uint64(universe.MaxSolarSystemPerGalaxy))

	maxPlanetPerSolarSystem := new(big.Int)
	maxPlanetPerSolarSystem.SetUint64(uint64(universe.MaxPlanetPerSolarSystem))

	planetCoordinateFound := false
	var try int64 = 0
	var coordinateId string
	galaxy := new(big.Int)
	solarSystem := new(big.Int)
	planet := new(big.Int)
	for !planetCoordinateFound {

		hash := crypto.Keccak256([]byte(tx.From + fmt.Sprint(ctx.Date.Unix()+try)))

		// Convert the hash to a hex string
		hashHex := hex.EncodeToString(hash)

		// Split the hash into three parts
		galaxyHex := hashHex[0:20]
		solarSystemHex := hashHex[20:40]
		planetHex := hashHex[40:60]

		// Convert the hex parts to big integers

		galaxy.SetString(galaxyHex, 16)
		galaxy.Mod(galaxy, maxGalaxy)

		solarSystem.SetString(solarSystemHex, 16)
		solarSystem.Mod(solarSystem, maxSolarSystem)

		planet.SetString(planetHex, 16)
		planet.Mod(planet, maxPlanetPerSolarSystem)

		coordinateId = types.CoordinateId(galaxy.Uint64(), solarSystem.Uint64(), planet.Uint64())
		exists, err := ctx.Db.Planets.Exists(tx.Data.UniverseId, coordinateId)
		if err != nil {
			return &abcitypes.ExecTxResult{
				Code: types.InvalidTransactionError,
			}
		}

		if exists {
			try += 1
		} else {
			planetCoordinateFound = true
		}
	}

	p := types.OChainPlanet{
		Owner:       globalAccount.Address,
		Universe:    universe.Id,
		Galaxy:      galaxy.Uint64(),
		SolarSystem: solarSystem.Uint64(),
		Planet:      planet.Uint64(),

		Buildings: types.OChainPlanetBuildings{
			MetalMine:       0,
			CrystalMine:     0,
			DeuteriumMine:   0,
			SolarPowerPlant: 0,

			RoboticFactory:   0,
			NaniteFactory:    0,
			SpaceshipFactory: 0,

			IntergalacticPortal: 0,
			ResearchLaboratory:  0,
			ShieldDome:          0,

			MetalStorage:     0,
			CrystalStorage:   0,
			DeuteriumStorage: 0,
		},

		Spaceships: types.OChainFleetSpaceships{
			SmallCargo:   0,
			LargeCargo:   0,
			LightFighter: 0,
			HeavyFighter: 0,

			Cruiser:       0,
			Battleship:    0,
			Battlecruiser: 0,

			Bomber:    0,
			Destroyer: 0,
			Deathstar: 0,
			Reaper:    0,
			Recycler:  0,
		},

		Defenses: types.OChainPlanetDefences{
			RocketLauncher:  0,
			LightLaser:      0,
			HeavyLaser:      0,
			IonCannon:       0,
			GaussCannon:     0,
			PlasmaTurret:    0,
			DarkMatterCanon: 0,
		},
		Resources: types.OChainResources{
			OCT:       0,
			Metal:     200,
			Crystal:   100,
			Deuterium: 0,
		},

		LastResourceUpdate: ctx.Date.Unix(),
	}

	account.PlanetsCoordinates = append(account.PlanetsCoordinates, coordinateId)

	err = ctx.Db.UniverseAccounts.Insert(account)
	if err != nil {
		return &abcitypes.ExecTxResult{
			Code: types.InvalidTransactionError,
		}
	}

	err = ctx.Db.Planets.Insert(universe.Id, p)
	if err != nil {
		return &abcitypes.ExecTxResult{
			Code: types.InvalidTransactionError,
		}
	}

	universe.Accounts += 1
	err = ctx.Db.Universes.Update(universe)
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
			Type: "UniverseAccountRegistered",
			Attributes: []abcitypes.EventAttribute{
				{Key: "address", Value: tx.From, Index: true},
				{Key: "universeId", Value: tx.Data.UniverseId, Index: true},
			},
		},
		{
			Type: "PlanetCreated",
			Attributes: []abcitypes.EventAttribute{
				{Key: "address", Value: tx.From, Index: true},
				{Key: "universeId", Value: tx.Data.UniverseId, Index: true},
				{Key: "coordinateId", Value: coordinateId},
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

func ParseRegisterUniverseAccountTransaction(tx t.Transaction) (RegisterUniverseAccountTransaction, error) {
	var txData RegisterUniverseAccountTransactionData
	err := cbor.Unmarshal(tx.Data, &txData)

	if err != nil {
		return RegisterUniverseAccountTransaction{}, err
	}

	return RegisterUniverseAccountTransaction{
		Type:      tx.Type,
		From:      tx.From,
		Nonce:     tx.Nonce,
		Data:      txData,
		Signature: tx.Signature,
	}, nil
}
