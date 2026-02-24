package types

const (
	// General success and system errors
	Success uint32 = 0
	NoError uint32 = 0 // Alias for Success

	// Transaction parsing errors (1-19)
	ParsingTransactionError     uint32 = 1
	ParsingTransactionDataError uint32 = 2
	InvalidTransactionError     uint32 = 3 // Deprecated - use specific errors below
	InvalidTransactionSignature uint32 = 4
	InvalidNonce                uint32 = 5
	CheckTransactionFailure     uint32 = 6
	ExecuteTransactionFailure   uint32 = 7
	NotImplemented              uint32 = 8
	GasCostHigherThanBalance    uint32 = 9

	// Entity not found errors (100-119)
	EntityNotFoundError          uint32 = 100
	GlobalAccountNotFoundError   uint32 = 101
	UniverseAccountNotFoundError uint32 = 102
	UniverseNotFoundError        uint32 = 103
	PlanetNotFoundError          uint32 = 104
	FleetNotFoundError           uint32 = 105
	AllianceNotFoundError        uint32 = 106
	BuildingNotFoundError        uint32 = 107
	SpaceshipNotFoundError       uint32 = 108
	DefenseNotFoundError         uint32 = 109
	TechnologyNotFoundError      uint32 = 110

	// Ownership/Permission errors (120-139)
	OwnershipViolationError      uint32 = 120
	FleetOwnershipError          uint32 = 121
	SignatureVerificationError   uint32 = 122
	InsufficientPermissionsError uint32 = 123
	AlliancePermissionError      uint32 = 124

	// State/Condition errors (140-159)
	EntityAlreadyExistsError      uint32 = 140
	AccountAlreadyExistsError     uint32 = 141
	AllianceMembershipError       uint32 = 142
	FleetInTransitError           uint32 = 143
	PendingOperationError         uint32 = 144
	StateConflictError            uint32 = 145
	InvalidStateError             uint32 = 146
	FleetNotAtTargetPositionError uint32 = 147

	// Resource/Requirement errors (160-179)
	InsufficientResourcesError uint32 = 160
	RequirementsNotMetError    uint32 = 161
	CapacityExceededError      uint32 = 162
	RateLimitExceededError     uint32 = 163
	InsufficientBalanceError   uint32 = 164
	MaxLevelReachedError       uint32 = 165

	// Validation errors (180-199)
	InvalidDataFormatError     uint32 = 180
	BusinessRuleViolationError uint32 = 181
	InvalidPositionError       uint32 = 182
	InvalidAmountError         uint32 = 183
	InvalidParametersError     uint32 = 184

	// System errors (200-219)
	DatabaseOperationError uint32 = 200
	SerializationError     uint32 = 201
	ConcurrencyError       uint32 = 202
	InternalServerError    uint32 = 203

	// Bridge/External errors (220-239)
	BridgeTransactionNotFoundError uint32 = 220
	BridgeTransactionFailedError   uint32 = 221
	ExternalServiceError           uint32 = 222
	ChainIdMismatchError           uint32 = 223

	// Governance errors (240-259)
	ProposalNotFoundError          uint32 = 240
	InvalidProposalStateError      uint32 = 241
	VotingPeriodExpiredError       uint32 = 242
	AlreadyVotedError              uint32 = 243
	InsufficientStakedBalance      uint32 = 244
	PendingUnstakeAlreadyExists    uint32 = 245
	NoPendingUnstakeFound          uint32 = 246
	UnstakeCooldownNotCompleted    uint32 = 247
	InsufficientVotingPowerError   uint32 = 248
	ProposalNotPassedError         uint32 = 249
	ProposalAlreadyExecutedError   uint32 = 250
	VotingPeriodNotExpiredError    uint32 = 251
	InvalidProposalTypeError       uint32 = 252
	InvalidProposalPayloadError    uint32 = 253
	UniverseAlreadyExistsError     uint32 = 254
	UniverseAlreadyPausedError     uint32 = 255
	UniverseNotPausedError         uint32 = 256

	// Commander/Subscription errors (260-279)
	CommanderNotFoundError      uint32 = 260
	CommanderExpiredError       uint32 = 261
	CommanderAlreadyActiveError uint32 = 262
)

type OChainError struct {
	Code    uint
	Message string
}
