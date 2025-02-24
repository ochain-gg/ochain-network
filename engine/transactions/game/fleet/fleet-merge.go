package fleet_transactions

import (
	abcitypes "github.com/cometbft/cometbft/abci/types"
	"github.com/fxamacker/cbor/v2"

	t "github.com/ochain-gg/ochain-network/engine/transactions"
	"github.com/ochain-gg/ochain-network/types"
)

type MergeFleetTransactionData struct {
	UniverseId string `cbor:"2,keyasint"`
	FleetIdA   string `cbor:"3,keyasint"`
	FleetIdB   string `cbor:"3,keyasint"`
}

type MergeFleetTransaction struct {
	Type      t.TransactionType         `cbor:"1,keyasint"`
	From      string                    `cbor:"2,keyasint"`
	Nonce     uint64                    `cbor:"3,keyasint"`
	Data      MergeFleetTransactionData `cbor:"4,keyasint"`
	Signature []byte                    `cbor:"5,keyasint"`
}

func (tx *MergeFleetTransaction) Transaction() (t.Transaction, error) {
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

func (tx *MergeFleetTransaction) Check(ctx t.TransactionContext) *abcitypes.CheckTxResponse {
	//Load data
	account, err := ctx.Db.UniverseAccounts.GetAt(tx.Data.UniverseId, tx.From, uint64(ctx.Date.Unix()))
	if err != nil {
		return &abcitypes.CheckTxResponse{
			Code: types.InvalidTransactionError,
		}
	}

	universe, err := ctx.Db.Universes.GetAt(tx.Data.UniverseId, uint64(ctx.Date.Unix()))
	if err != nil {
		return &abcitypes.CheckTxResponse{
			Code: types.InvalidTransactionError,
		}
	}

	fleetA, err := ctx.Db.Fleets.Get(universe.Id, account, tx.Data.FleetIdA)
	if err != nil {
		return &abcitypes.CheckTxResponse{
			Code: types.InvalidTransactionError,
		}
	}

	fleetB, err := ctx.Db.Fleets.Get(universe.Id, account, tx.Data.FleetIdB)
	if err != nil {
		return &abcitypes.CheckTxResponse{
			Code: types.InvalidTransactionError,
		}
	}

	//Verify fleet ownership
	if fleetA.Owner != tx.From || fleetB.Owner != tx.From {
		return &abcitypes.CheckTxResponse{
			Code: types.InvalidTransactionError,
		}
	}

	//Verify fleet isn't travaling
	if fleetA.IsTraveling(ctx.Date.Unix()) || fleetB.IsTraveling(ctx.Date.Unix()) {
		return &abcitypes.CheckTxResponse{
			Code: types.InvalidTransactionError,
		}
	}

	return &abcitypes.CheckTxResponse{
		Code: types.NoError,
	}
}

func (tx *MergeFleetTransaction) Execute(ctx t.TransactionContext) *abcitypes.ExecTxResult {
	result := tx.Check(ctx)
	if result.Code != types.NoError {
		return &abcitypes.ExecTxResult{
			Code: result.GetCode(),
		}
	}

	account, err := ctx.Db.UniverseAccounts.GetAt(tx.Data.UniverseId, tx.From, uint64(ctx.Date.Unix()))
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

	fleetA, err := ctx.Db.Fleets.Get(universe.Id, account, tx.Data.FleetIdA)
	if err != nil {
		return &abcitypes.ExecTxResult{
			Code: types.InvalidTransactionError,
		}
	}

	fleetB, err := ctx.Db.Fleets.Get(universe.Id, account, tx.Data.FleetIdB)
	if err != nil {
		return &abcitypes.ExecTxResult{
			Code: types.InvalidTransactionError,
		}
	}

	for _, spaceshipId := range types.GetSpaceshipIds() {
		fleetA.AddSpaceships(spaceshipId, fleetB.Spaceships.Get(spaceshipId))
	}

	fleetA.Cargo.Add(fleetB.Cargo)

	err = ctx.Db.Fleets.Delete(universe.Id, account, fleetB.Id)
	if err != nil {
		return &abcitypes.ExecTxResult{
			Code: types.InvalidTransactionError,
		}
	}

	err = ctx.Db.Fleets.Update(universe.Id, account, fleetA)
	if err != nil {
		return &abcitypes.ExecTxResult{
			Code: types.InvalidTransactionError,
		}
	}

	events := []abcitypes.Event{
		{
			Type: "FleetMerged",
			Attributes: []abcitypes.EventAttribute{
				{Key: "account", Value: tx.From, Index: true},
				{Key: "universe", Value: tx.Data.UniverseId, Index: true},
				{Key: "fleetA", Value: fleetA.Id, Index: true},
				{Key: "fleetB", Value: fleetB.Id, Index: true},
			},
		},
	}

	return &abcitypes.ExecTxResult{
		Code:      types.NoError,
		Events:    events,
		GasUsed:   100,
		GasWanted: 100,
		Data:      []byte(""),
	}
}

func ParseMergeFleetTransaction(tx t.Transaction) (MergeFleetTransaction, error) {
	var txData MergeFleetTransactionData
	err := cbor.Unmarshal(tx.Data, &txData)

	if err != nil {
		return MergeFleetTransaction{}, err
	}

	return MergeFleetTransaction{
		Type:      tx.Type,
		From:      tx.From,
		Nonce:     tx.Nonce,
		Data:      txData,
		Signature: tx.Signature,
	}, nil
}
