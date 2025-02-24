package validator_transactions

import (
	abcitypes "github.com/cometbft/cometbft/abci/types"
	"github.com/fxamacker/cbor/v2"

	t "github.com/ochain-gg/ochain-network/engine/transactions"
	"github.com/ochain-gg/ochain-network/types"
)

type OChainNewEpochTransactionData struct{}

type OChainNewEpochTransaction struct {
	Type t.TransactionType             `cbor:"1,keyasint"`
	Data OChainNewEpochTransactionData `cbor:"2,keyasint"`
}

func (tx OChainNewEpochTransaction) Check(ctx t.TransactionContext) *abcitypes.CheckTxResponse {

	currentEpoch, err := ctx.Db.Epochs.GetCurrentAt(uint64(ctx.Date.Unix()))
	if err != nil {
		return &abcitypes.CheckTxResponse{
			Code: types.CheckTransactionFailure,
		}
	}

	if currentEpoch.EndedAt < ctx.Date.Unix() {
		return &abcitypes.CheckTxResponse{
			Code: types.CheckTransactionFailure,
		}
	}

	return &abcitypes.CheckTxResponse{
		Code: types.NoError,
	}
}

func (tx OChainNewEpochTransaction) Execute(ctx t.TransactionContext) *abcitypes.ExecTxResult {
	result := tx.Check(ctx)
	if result.Code != types.NoError {
		return &abcitypes.ExecTxResult{
			Code: result.GetCode(),
		}
	}

	currentEpoch, err := ctx.Db.Epochs.GetCurrentAt(uint64(ctx.Date.Unix()))
	if err != nil {
		return &abcitypes.ExecTxResult{
			Code: types.CheckTransactionFailure,
		}
	}

	if currentEpoch.EndedAt < ctx.Date.Unix() {
		return &abcitypes.ExecTxResult{
			Code: types.CheckTransactionFailure,
		}
	}

	// Generate rewards program distribution

	// Create new epoch
	newEpoch := types.OChainEpoch{
		Id:        currentEpoch.Id + 1,
		StartedAt: currentEpoch.EndedAt,
		EndedAt:   currentEpoch.EndedAt + (3600 * 24 * 30),

		TokenEarned: 0,
		USDEarned:   0,
	}

	err = ctx.Db.Epochs.Insert(newEpoch)
	if err != nil {
		return &abcitypes.ExecTxResult{
			Code: types.CheckTransactionFailure,
		}
	}

	return &abcitypes.ExecTxResult{
		Code: types.NoError,
	}
}

func (tx OChainNewEpochTransaction) Transaction() (t.Transaction, error) {

	txData, err := cbor.Marshal(tx.Data)
	if err != nil {
		return t.Transaction{}, err
	}

	return t.Transaction{
		Type: tx.Type,
		Data: txData,
	}, nil
}

func ParseOChainNewEpochTransaction(tx t.Transaction) (OChainNewEpochTransaction, error) {
	var txData OChainNewEpochTransactionData
	err := cbor.Unmarshal(tx.Data, &txData)

	if err != nil {
		return OChainNewEpochTransaction{}, err
	}

	return OChainNewEpochTransaction{
		Type: tx.Type,
		Data: txData,
	}, nil
}
