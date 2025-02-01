package account_transactions

import (
	abcitypes "github.com/cometbft/cometbft/abci/types"

	t "github.com/ochain-gg/ochain-network/transactions"
	"github.com/ochain-gg/ochain-network/types"
)

type ChangeAccountIAMTransactionData struct {
	GuardianQuorum uint64   `cbor:"1,keyasint"`
	Guardians      []string `cbor:"2,keyasint"`
	DeleguatedTo   []string `cbor:"3,keyasint"`
}

type ChangeAccountIAMTransaction struct {
	Type      t.TransactionType `cbor:"1,keyasint"`
	From      string            `cbor:"2,keyasint"`
	Nonce     uint64            `cbor:"3,keyasint"`
	Data      string            `cbor:"4,keyasint"`
	Signature string            `cbor:"5,keyasint"`
}

func (tx *ChangeAccountIAMTransaction) Check(ctx t.TransactionContext) *abcitypes.CheckTxResponse {
	_, err := ctx.Db.GlobalsAccounts.Get(tx.From)
	if err != nil {
		return &abcitypes.CheckTxResponse{
			Code: types.InvalidTransactionError,
		}
	}

	//TODO: verify authrorizer signature
	return &abcitypes.CheckTxResponse{
		Code: types.NoError,
	}
}

func (tx *ChangeAccountIAMTransaction) Execute(ctx t.TransactionContext) *abcitypes.ExecTxResult {
	result := tx.Check(ctx)
	if result.Code != types.NoError {
		return &abcitypes.ExecTxResult{
			Code: result.GetCode(),
		}
	}

	_, err := ctx.Db.GlobalsAccounts.Get(tx.From)
	if err != nil {
		return &abcitypes.ExecTxResult{
			Code: types.ExecuteTransactionFailure,
		}
	}

	return &abcitypes.ExecTxResult{
		Code: types.NoError,
	}
}
