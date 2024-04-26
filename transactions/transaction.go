package transactions

import "encoding/json"

type TransactionType int64

const (
	//Remote chain based transactions
	NewValidator    TransactionType = 0
	RemoveValidator TransactionType = 1

	OChainTokenDeposit    TransactionType = 2
	OChainTokenWithdrawal TransactionType = 3

	//Universe configuration
	CreateUniverse       TransactionType = 4
	CreateGovernanceVote TransactionType = 5
	CreateUserVote       TransactionType = 6

	//Account transactions
	RegisterAccount          TransactionType = 7
	AddAccountDeleguation    TransactionType = 8
	SetGuardianConfiguration TransactionType = 9
	AddGuardian              TransactionType = 10

	//Planet transactions
	MintPlanet             TransactionType = 11
	StartBuildingUpgrade   TransactionType = 12
	StartTechnologyUpgrade TransactionType = 13
	StartBuildDefenses     TransactionType = 14
	StartBuildSpaceships   TransactionType = 15

	//Fleet transactions
	FillCargo                    TransactionType = 16
	UnfillCargo                  TransactionType = 17
	SendFleetInOrbit             TransactionType = 18
	LandingFleetInOrbit          TransactionType = 19
	MergeFleets                  TransactionType = 20
	SplitFleet                   TransactionType = 21
	RecycleRemnant               TransactionType = 22
	IntergalacticPortalFleetMove TransactionType = 23
	ChangeFleetMode              TransactionType = 24
	AcceptFleetMode              TransactionType = 25
	MoveFleet                    TransactionType = 26
	CancelFleetMove              TransactionType = 27

	//Fight transactions
	Fight TransactionType = 28

	//Refresh transaction
	ExecuteBuildingUpgrade   TransactionType = 29
	ExecuteTechnologyUpgrade TransactionType = 30
	ExecuteDefenseBuild      TransactionType = 31
	ExecuteSpaceshipBuild    TransactionType = 32
	ExecuteFleetMove         TransactionType = 33

	//Alliance transactions
	CreateAlliance  TransactionType = 34
	StakeOnAlliance TransactionType = 35

	//Market transations
	SwapResources TransactionType = 36
)

type Transaction struct {
	Type  TransactionType `json:"type"`
	From  []byte          `json:"from"`
	Nonce uint64          `json:"nonce"`
	Data  []byte          `json:"data"`
}

func (tx *Transaction) IsValid() uint32 {
	switch tx.Type {
	case NewValidator:
		_, err := ParseNewValidatorTransaction(*tx)
		if err != nil {
			return 2
		}
		return 0

	case RemoveValidator:
		_, err := ParseRemoveValidatorTransaction(*tx)
		if err != nil {
			return 2
		}
		return 0
		// case NewEpoch:

		// case NewVault:

		// case UpdateVault:

		// case NewPositionManager:

		// case UpdatePositionManager:

		// case NewOracle:

		// case UpdateOracle:

	}
	return 1
}

func (tx *Transaction) IsValidSignature() bool {

	return true
}

func ParseTransaction(data []byte) (Transaction, error) {

	var tx Transaction
	err := json.Unmarshal(data, &tx)

	if err != nil {
		return Transaction{}, err
	}

	return tx, nil
}
