package governance_transactions

// type GovernanceProposalType uint32
// const (
// 	GovernanceProposalType NewUniverseProoposal = 0
// )

// type NewUniverseProposalPayload struct {
// 	Id    uint64 `cbor:"1,keyasint"`
// 	Name  string `cbor:"2,keyasint"`
// 	Speed uint64 `cbor:"3,keyasint"`
// }

// type CreateGovernanceProposalTransactionData struct {
// 	Type    uint32 `cbor:"1,keyasint"`
// 	Payload []byte `cbor:"2,keyasint"`
// }

// type CreateGovernanceProposalTransaction struct {
// 	Type      t.TransactionType                       `cbor:"1,keyasint"`
// 	From      string                                  `cbor:"2,keyasint"`
// 	Nonce     uint64                                  `cbor:"3,keyasint"`
// 	Data      CreateGovernanceProposalTransactionData `cbor:"4,keyasint"`
// 	Signature []byte                                  `cbor:"5,keyasint"`
// }

// func (tx *CreateGovernanceProposalTransaction) Check(ctx t.TransactionContext) *abcitypes.CheckTxResponse {
// 	globalAccount, err := ctx.Db.GlobalsAccounts.Get(tx.From)
// 	if err != nil {
// 		return &abcitypes.CheckTxResponse{
// 			Code: types.InvalidTransactionError,
// 		}
// 	}

// 	if globalAccount.StackedBalance < types.VotingPowerRequiredForProposalCreation {
// 		return &abcitypes.CheckTxResponse{
// 			Code: types.InvalidTransactionError,
// 		}
// 	}

// 	return &abcitypes.CheckTxResponse{
// 		Code: types.NoError,
// 	}
// }

// func (tx *CreateGovernanceProposalTransaction) Execute(ctx t.TransactionContext) *abcitypes.ExecTxResult {
// 	response := tx.Check(ctx)
// 	if response.Code != types.NoError {
// 		return &abcitypes.ExecTxResult{
// 			Code: response.Code,
// 		}
// 	}

// 	globalAccount, err := ctx.Db.GlobalsAccounts.Get(tx.From)
// 	if err == nil {
// 		return &abcitypes.ExecTxResult{
// 			Code: types.InvalidTransactionError,
// 		}
// 	}

// 	account, err := ctx.Db.UniverseAccounts.Get(tx.Data..Universe, tx.From)
// 	if err == nil {
// 		return &abcitypes.ExecTxResult{
// 			Code: types.InvalidTransactionError,
// 		}
// 	}

// 	globalAccount.CreditBalance -= types.TwoWeekCommanderPrice
// 	account.SubscribeToCommander(ctx.Date.Unix(), tx.Data.Commander)

// 	txGasCost, err := globalAccount.ApplyGasCost(uint64(ctx.Date.Unix()))
// 	if err != nil {
// 		return &abcitypes.ExecTxResult{
// 			Code: types.GasCostHigherThanBalance,
// 		}
// 	}

// 	err = ctx.Db.GlobalsAccounts.Update(globalAccount)
// 	if err != nil {
// 		return &abcitypes.ExecTxResult{
// 			Code: types.InvalidTransactionError,
// 		}
// 	}

// 	err = ctx.Db.UniverseAccounts.Update(account)
// 	if err != nil {
// 		return &abcitypes.ExecTxResult{
// 			Code: types.InvalidTransactionError,
// 		}
// 	}

// 	events := []abcitypes.Event{
// 		{
// 			Type: "CommanderSubscribed",
// 			Attributes: []abcitypes.EventAttribute{
// 				{Key: "account", Value: tx.From, Index: true},
// 				{Key: "universe", Value: tx.Data.Universe, Index: true},
// 				{Key: "commander", Value: string(tx.Data.Commander), Index: true},
// 			},
// 		},
// 	}

// 	receipt := t.TransactionReceipt{
// 		GasCost: txGasCost,
// 	}

// 	return &abcitypes.ExecTxResult{
// 		Code:      types.NoError,
// 		Events:    events,
// 		GasUsed:   100,
// 		GasWanted: 100,
// 		Data:      receipt.Bytes(),
// 	}
// }

// func (tx *CreateGovernanceProposalTransaction) Transaction() (t.Transaction, error) {
// 	txData, err := cbor.Marshal(tx.Data)
// 	if err != nil {
// 		return t.Transaction{}, err
// 	}

// 	return t.Transaction{
// 		Type:      tx.Type,
// 		From:      tx.From,
// 		Nonce:     tx.Nonce,
// 		Data:      txData,
// 		Signature: tx.Signature,
// 	}, nil
// }

// func ParseCreateGovernanceProposalTransaction(tx t.Transaction) (CreateGovernanceProposalTransaction, error) {
// 	var txData CreateGovernanceProposalTransactionData
// 	err := cbor.Unmarshal(tx.Data, &txData)

// 	if err != nil {
// 		return CreateGovernanceProposalTransaction{}, err
// 	}

// 	return CreateGovernanceProposalTransaction{
// 		Type:      tx.Type,
// 		From:      tx.From,
// 		Nonce:     tx.Nonce,
// 		Data:      txData,
// 		Signature: tx.Signature,
// 	}, nil
// }
