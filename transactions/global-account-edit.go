package transactions

import (
	"errors"

	abcitypes "github.com/cometbft/cometbft/abci/types"
	"github.com/ochain-gg/ochain-network/types"
)

type ChangeAccountIAMTransactionData struct {
	GuardianQuorum uint64   `cbor:"1,keyasint"`
	Guardians      []string `cbor:"2,keyasint"`
	DeleguatedTo   []string `cbor:"3,keyasint"`
}

type ChangeAccountIAMTransaction struct {
	Type      TransactionType `cbor:"1,keyasint"`
	From      string          `cbor:"2,keyasint"`
	Nonce     uint64          `cbor:"3,keyasint"`
	Data      string          `cbor:"4,keyasint"`
	Signature string          `cbor:"5,keyasint"`
}

func (tx *ChangeAccountIAMTransaction) Check(ctx TransactionContext) *abcitypes.ResponseCheckTx {
	_, err := ctx.Db.GlobalsAccounts.Get(tx.From)
	if err != nil {
		return 1, errors.New("account don't exists")
	}

	//TODO: verify authrorizer signature
	return 0, nil
}

func (tx *ChangeAccountIAMTransaction) Execute(ctx TransactionContext) *abcitypes.ExecTxResult {
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
