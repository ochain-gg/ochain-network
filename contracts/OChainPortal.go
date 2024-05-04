// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contracts

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
	_ = abi.ConvertType
)

// LibDiamondValidatorStacking is an auto generated low-level Go binding around an user-defined struct.
type LibDiamondValidatorStacking struct {
	Enabled             bool
	Depositor           common.Address
	Validator           common.Address
	PublicKey           string
	Token               common.Address
	Amount              *big.Int
	WithdrawRequestDate *big.Int
}

// OChainPortalMetaData contains all meta data concerning the OChainPortal contract.
var OChainPortalMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"name\":\"AlreadyExecuted\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"BadSignature\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"GracePeriodNotEnded\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"MaxValidatorReached\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NonceInvalid\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NotAuthorized\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"QuorumNotReached\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UnstakeProcessNotAvailable\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UnstakeProcessNotEnded\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"WithdrawalNonceContexted\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"validatorId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"stacker\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"validator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"publicKey\",\"type\":\"string\"}],\"name\":\"OChainNewValidator\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"validatorId\",\"type\":\"uint256\"}],\"name\":\"OChainRemoveValidator\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"OChainTokenDeposit\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"}],\"name\":\"OChainTokenWithdrawalRequestContested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"}],\"name\":\"OChainTokenWithdrawalRequestExecuted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256[]\",\"name\":\"signers\",\"type\":\"uint256[]\"}],\"name\":\"OChainTokenWithdrawalRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"validatorId\",\"type\":\"uint256\"}],\"name\":\"OChainUnstackSucceed\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"authorizer\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"},{\"internalType\":\"uint256[]\",\"name\":\"validatorIds\",\"type\":\"uint256[]\"},{\"internalType\":\"uint256[]\",\"name\":\"signatures\",\"type\":\"uint256[]\"}],\"name\":\"contestWithdraw\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"universe\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"planet\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"deposit\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"}],\"name\":\"executeWithdraw\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestUpdateAt\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"stacker\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"validator\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"pubkey\",\"type\":\"string\"}],\"name\":\"stake\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"validatorId\",\"type\":\"uint256\"}],\"name\":\"startUnstakeProcess\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"validatorId\",\"type\":\"uint256\"}],\"name\":\"unstake\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"validatorId\",\"type\":\"uint256\"}],\"name\":\"validatorInfo\",\"outputs\":[{\"components\":[{\"internalType\":\"bool\",\"name\":\"enabled\",\"type\":\"bool\"},{\"internalType\":\"address\",\"name\":\"depositor\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"validator\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"publicKey\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"withdrawRequestDate\",\"type\":\"uint256\"}],\"internalType\":\"structLibDiamond.ValidatorStacking\",\"name\":\"_validator\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"validatorNetworkInfo\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"maxValidators\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"validatorsMaxIndex\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"validatorsLength\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint256[]\",\"name\":\"validatorIds\",\"type\":\"uint256[]\"},{\"internalType\":\"uint256[]\",\"name\":\"signatures\",\"type\":\"uint256[]\"}],\"name\":\"withdrawalRequest\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// OChainPortalABI is the input ABI used to generate the binding from.
// Deprecated: Use OChainPortalMetaData.ABI instead.
var OChainPortalABI = OChainPortalMetaData.ABI

// OChainPortal is an auto generated Go binding around an Ethereum contract.
type OChainPortal struct {
	OChainPortalCaller     // Read-only binding to the contract
	OChainPortalTransactor // Write-only binding to the contract
	OChainPortalFilterer   // Log filterer for contract events
}

// OChainPortalCaller is an auto generated read-only Go binding around an Ethereum contract.
type OChainPortalCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OChainPortalTransactor is an auto generated write-only Go binding around an Ethereum contract.
type OChainPortalTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OChainPortalFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type OChainPortalFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OChainPortalSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type OChainPortalSession struct {
	Contract     *OChainPortal     // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// OChainPortalCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type OChainPortalCallerSession struct {
	Contract *OChainPortalCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts       // Call options to use throughout this session
}

// OChainPortalTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type OChainPortalTransactorSession struct {
	Contract     *OChainPortalTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts       // Transaction auth options to use throughout this session
}

// OChainPortalRaw is an auto generated low-level Go binding around an Ethereum contract.
type OChainPortalRaw struct {
	Contract *OChainPortal // Generic contract binding to access the raw methods on
}

// OChainPortalCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type OChainPortalCallerRaw struct {
	Contract *OChainPortalCaller // Generic read-only contract binding to access the raw methods on
}

// OChainPortalTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type OChainPortalTransactorRaw struct {
	Contract *OChainPortalTransactor // Generic write-only contract binding to access the raw methods on
}

// NewOChainPortal creates a new instance of OChainPortal, bound to a specific deployed contract.
func NewOChainPortal(address common.Address, backend bind.ContractBackend) (*OChainPortal, error) {
	contract, err := bindOChainPortal(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &OChainPortal{OChainPortalCaller: OChainPortalCaller{contract: contract}, OChainPortalTransactor: OChainPortalTransactor{contract: contract}, OChainPortalFilterer: OChainPortalFilterer{contract: contract}}, nil
}

// NewOChainPortalCaller creates a new read-only instance of OChainPortal, bound to a specific deployed contract.
func NewOChainPortalCaller(address common.Address, caller bind.ContractCaller) (*OChainPortalCaller, error) {
	contract, err := bindOChainPortal(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &OChainPortalCaller{contract: contract}, nil
}

// NewOChainPortalTransactor creates a new write-only instance of OChainPortal, bound to a specific deployed contract.
func NewOChainPortalTransactor(address common.Address, transactor bind.ContractTransactor) (*OChainPortalTransactor, error) {
	contract, err := bindOChainPortal(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &OChainPortalTransactor{contract: contract}, nil
}

// NewOChainPortalFilterer creates a new log filterer instance of OChainPortal, bound to a specific deployed contract.
func NewOChainPortalFilterer(address common.Address, filterer bind.ContractFilterer) (*OChainPortalFilterer, error) {
	contract, err := bindOChainPortal(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &OChainPortalFilterer{contract: contract}, nil
}

// bindOChainPortal binds a generic wrapper to an already deployed contract.
func bindOChainPortal(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := OChainPortalMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_OChainPortal *OChainPortalRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _OChainPortal.Contract.OChainPortalCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_OChainPortal *OChainPortalRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OChainPortal.Contract.OChainPortalTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_OChainPortal *OChainPortalRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _OChainPortal.Contract.OChainPortalTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_OChainPortal *OChainPortalCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _OChainPortal.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_OChainPortal *OChainPortalTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OChainPortal.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_OChainPortal *OChainPortalTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _OChainPortal.Contract.contract.Transact(opts, method, params...)
}

// Authorizer is a free data retrieval call binding the contract method 0xd09edf31.
//
// Solidity: function authorizer() view returns(address)
func (_OChainPortal *OChainPortalCaller) Authorizer(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _OChainPortal.contract.Call(opts, &out, "authorizer")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Authorizer is a free data retrieval call binding the contract method 0xd09edf31.
//
// Solidity: function authorizer() view returns(address)
func (_OChainPortal *OChainPortalSession) Authorizer() (common.Address, error) {
	return _OChainPortal.Contract.Authorizer(&_OChainPortal.CallOpts)
}

// Authorizer is a free data retrieval call binding the contract method 0xd09edf31.
//
// Solidity: function authorizer() view returns(address)
func (_OChainPortal *OChainPortalCallerSession) Authorizer() (common.Address, error) {
	return _OChainPortal.Contract.Authorizer(&_OChainPortal.CallOpts)
}

// LatestUpdateAt is a free data retrieval call binding the contract method 0xb185bcac.
//
// Solidity: function latestUpdateAt() view returns(uint256)
func (_OChainPortal *OChainPortalCaller) LatestUpdateAt(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _OChainPortal.contract.Call(opts, &out, "latestUpdateAt")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// LatestUpdateAt is a free data retrieval call binding the contract method 0xb185bcac.
//
// Solidity: function latestUpdateAt() view returns(uint256)
func (_OChainPortal *OChainPortalSession) LatestUpdateAt() (*big.Int, error) {
	return _OChainPortal.Contract.LatestUpdateAt(&_OChainPortal.CallOpts)
}

// LatestUpdateAt is a free data retrieval call binding the contract method 0xb185bcac.
//
// Solidity: function latestUpdateAt() view returns(uint256)
func (_OChainPortal *OChainPortalCallerSession) LatestUpdateAt() (*big.Int, error) {
	return _OChainPortal.Contract.LatestUpdateAt(&_OChainPortal.CallOpts)
}

// ValidatorInfo is a free data retrieval call binding the contract method 0x1e47d4fe.
//
// Solidity: function validatorInfo(uint256 validatorId) view returns((bool,address,address,string,address,uint256,uint256) _validator)
func (_OChainPortal *OChainPortalCaller) ValidatorInfo(opts *bind.CallOpts, validatorId *big.Int) (LibDiamondValidatorStacking, error) {
	var out []interface{}
	err := _OChainPortal.contract.Call(opts, &out, "validatorInfo", validatorId)

	if err != nil {
		return *new(LibDiamondValidatorStacking), err
	}

	out0 := *abi.ConvertType(out[0], new(LibDiamondValidatorStacking)).(*LibDiamondValidatorStacking)

	return out0, err

}

// ValidatorInfo is a free data retrieval call binding the contract method 0x1e47d4fe.
//
// Solidity: function validatorInfo(uint256 validatorId) view returns((bool,address,address,string,address,uint256,uint256) _validator)
func (_OChainPortal *OChainPortalSession) ValidatorInfo(validatorId *big.Int) (LibDiamondValidatorStacking, error) {
	return _OChainPortal.Contract.ValidatorInfo(&_OChainPortal.CallOpts, validatorId)
}

// ValidatorInfo is a free data retrieval call binding the contract method 0x1e47d4fe.
//
// Solidity: function validatorInfo(uint256 validatorId) view returns((bool,address,address,string,address,uint256,uint256) _validator)
func (_OChainPortal *OChainPortalCallerSession) ValidatorInfo(validatorId *big.Int) (LibDiamondValidatorStacking, error) {
	return _OChainPortal.Contract.ValidatorInfo(&_OChainPortal.CallOpts, validatorId)
}

// ValidatorNetworkInfo is a free data retrieval call binding the contract method 0x3dee05b1.
//
// Solidity: function validatorNetworkInfo() view returns(uint256 maxValidators, uint256 validatorsMaxIndex, uint256 validatorsLength)
func (_OChainPortal *OChainPortalCaller) ValidatorNetworkInfo(opts *bind.CallOpts) (struct {
	MaxValidators      *big.Int
	ValidatorsMaxIndex *big.Int
	ValidatorsLength   *big.Int
}, error) {
	var out []interface{}
	err := _OChainPortal.contract.Call(opts, &out, "validatorNetworkInfo")

	outstruct := new(struct {
		MaxValidators      *big.Int
		ValidatorsMaxIndex *big.Int
		ValidatorsLength   *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.MaxValidators = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.ValidatorsMaxIndex = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.ValidatorsLength = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// ValidatorNetworkInfo is a free data retrieval call binding the contract method 0x3dee05b1.
//
// Solidity: function validatorNetworkInfo() view returns(uint256 maxValidators, uint256 validatorsMaxIndex, uint256 validatorsLength)
func (_OChainPortal *OChainPortalSession) ValidatorNetworkInfo() (struct {
	MaxValidators      *big.Int
	ValidatorsMaxIndex *big.Int
	ValidatorsLength   *big.Int
}, error) {
	return _OChainPortal.Contract.ValidatorNetworkInfo(&_OChainPortal.CallOpts)
}

// ValidatorNetworkInfo is a free data retrieval call binding the contract method 0x3dee05b1.
//
// Solidity: function validatorNetworkInfo() view returns(uint256 maxValidators, uint256 validatorsMaxIndex, uint256 validatorsLength)
func (_OChainPortal *OChainPortalCallerSession) ValidatorNetworkInfo() (struct {
	MaxValidators      *big.Int
	ValidatorsMaxIndex *big.Int
	ValidatorsLength   *big.Int
}, error) {
	return _OChainPortal.Contract.ValidatorNetworkInfo(&_OChainPortal.CallOpts)
}

// ContestWithdraw is a paid mutator transaction binding the contract method 0x2c59bc28.
//
// Solidity: function contestWithdraw(uint256 nonce, uint256[] validatorIds, uint256[] signatures) returns()
func (_OChainPortal *OChainPortalTransactor) ContestWithdraw(opts *bind.TransactOpts, nonce *big.Int, validatorIds []*big.Int, signatures []*big.Int) (*types.Transaction, error) {
	return _OChainPortal.contract.Transact(opts, "contestWithdraw", nonce, validatorIds, signatures)
}

// ContestWithdraw is a paid mutator transaction binding the contract method 0x2c59bc28.
//
// Solidity: function contestWithdraw(uint256 nonce, uint256[] validatorIds, uint256[] signatures) returns()
func (_OChainPortal *OChainPortalSession) ContestWithdraw(nonce *big.Int, validatorIds []*big.Int, signatures []*big.Int) (*types.Transaction, error) {
	return _OChainPortal.Contract.ContestWithdraw(&_OChainPortal.TransactOpts, nonce, validatorIds, signatures)
}

// ContestWithdraw is a paid mutator transaction binding the contract method 0x2c59bc28.
//
// Solidity: function contestWithdraw(uint256 nonce, uint256[] validatorIds, uint256[] signatures) returns()
func (_OChainPortal *OChainPortalTransactorSession) ContestWithdraw(nonce *big.Int, validatorIds []*big.Int, signatures []*big.Int) (*types.Transaction, error) {
	return _OChainPortal.Contract.ContestWithdraw(&_OChainPortal.TransactOpts, nonce, validatorIds, signatures)
}

// Deposit is a paid mutator transaction binding the contract method 0x00aeef8a.
//
// Solidity: function deposit(uint256 universe, uint256 planet, uint256 amount) returns()
func (_OChainPortal *OChainPortalTransactor) Deposit(opts *bind.TransactOpts, universe *big.Int, planet *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _OChainPortal.contract.Transact(opts, "deposit", universe, planet, amount)
}

// Deposit is a paid mutator transaction binding the contract method 0x00aeef8a.
//
// Solidity: function deposit(uint256 universe, uint256 planet, uint256 amount) returns()
func (_OChainPortal *OChainPortalSession) Deposit(universe *big.Int, planet *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _OChainPortal.Contract.Deposit(&_OChainPortal.TransactOpts, universe, planet, amount)
}

// Deposit is a paid mutator transaction binding the contract method 0x00aeef8a.
//
// Solidity: function deposit(uint256 universe, uint256 planet, uint256 amount) returns()
func (_OChainPortal *OChainPortalTransactorSession) Deposit(universe *big.Int, planet *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _OChainPortal.Contract.Deposit(&_OChainPortal.TransactOpts, universe, planet, amount)
}

// ExecuteWithdraw is a paid mutator transaction binding the contract method 0x2104ebf1.
//
// Solidity: function executeWithdraw(uint256 nonce) returns()
func (_OChainPortal *OChainPortalTransactor) ExecuteWithdraw(opts *bind.TransactOpts, nonce *big.Int) (*types.Transaction, error) {
	return _OChainPortal.contract.Transact(opts, "executeWithdraw", nonce)
}

// ExecuteWithdraw is a paid mutator transaction binding the contract method 0x2104ebf1.
//
// Solidity: function executeWithdraw(uint256 nonce) returns()
func (_OChainPortal *OChainPortalSession) ExecuteWithdraw(nonce *big.Int) (*types.Transaction, error) {
	return _OChainPortal.Contract.ExecuteWithdraw(&_OChainPortal.TransactOpts, nonce)
}

// ExecuteWithdraw is a paid mutator transaction binding the contract method 0x2104ebf1.
//
// Solidity: function executeWithdraw(uint256 nonce) returns()
func (_OChainPortal *OChainPortalTransactorSession) ExecuteWithdraw(nonce *big.Int) (*types.Transaction, error) {
	return _OChainPortal.Contract.ExecuteWithdraw(&_OChainPortal.TransactOpts, nonce)
}

// Stake is a paid mutator transaction binding the contract method 0xc27c12c2.
//
// Solidity: function stake(address stacker, address validator, string pubkey) returns()
func (_OChainPortal *OChainPortalTransactor) Stake(opts *bind.TransactOpts, stacker common.Address, validator common.Address, pubkey string) (*types.Transaction, error) {
	return _OChainPortal.contract.Transact(opts, "stake", stacker, validator, pubkey)
}

// Stake is a paid mutator transaction binding the contract method 0xc27c12c2.
//
// Solidity: function stake(address stacker, address validator, string pubkey) returns()
func (_OChainPortal *OChainPortalSession) Stake(stacker common.Address, validator common.Address, pubkey string) (*types.Transaction, error) {
	return _OChainPortal.Contract.Stake(&_OChainPortal.TransactOpts, stacker, validator, pubkey)
}

// Stake is a paid mutator transaction binding the contract method 0xc27c12c2.
//
// Solidity: function stake(address stacker, address validator, string pubkey) returns()
func (_OChainPortal *OChainPortalTransactorSession) Stake(stacker common.Address, validator common.Address, pubkey string) (*types.Transaction, error) {
	return _OChainPortal.Contract.Stake(&_OChainPortal.TransactOpts, stacker, validator, pubkey)
}

// StartUnstakeProcess is a paid mutator transaction binding the contract method 0xdbf06bb4.
//
// Solidity: function startUnstakeProcess(uint256 validatorId) returns()
func (_OChainPortal *OChainPortalTransactor) StartUnstakeProcess(opts *bind.TransactOpts, validatorId *big.Int) (*types.Transaction, error) {
	return _OChainPortal.contract.Transact(opts, "startUnstakeProcess", validatorId)
}

// StartUnstakeProcess is a paid mutator transaction binding the contract method 0xdbf06bb4.
//
// Solidity: function startUnstakeProcess(uint256 validatorId) returns()
func (_OChainPortal *OChainPortalSession) StartUnstakeProcess(validatorId *big.Int) (*types.Transaction, error) {
	return _OChainPortal.Contract.StartUnstakeProcess(&_OChainPortal.TransactOpts, validatorId)
}

// StartUnstakeProcess is a paid mutator transaction binding the contract method 0xdbf06bb4.
//
// Solidity: function startUnstakeProcess(uint256 validatorId) returns()
func (_OChainPortal *OChainPortalTransactorSession) StartUnstakeProcess(validatorId *big.Int) (*types.Transaction, error) {
	return _OChainPortal.Contract.StartUnstakeProcess(&_OChainPortal.TransactOpts, validatorId)
}

// Unstake is a paid mutator transaction binding the contract method 0x2e17de78.
//
// Solidity: function unstake(uint256 validatorId) returns()
func (_OChainPortal *OChainPortalTransactor) Unstake(opts *bind.TransactOpts, validatorId *big.Int) (*types.Transaction, error) {
	return _OChainPortal.contract.Transact(opts, "unstake", validatorId)
}

// Unstake is a paid mutator transaction binding the contract method 0x2e17de78.
//
// Solidity: function unstake(uint256 validatorId) returns()
func (_OChainPortal *OChainPortalSession) Unstake(validatorId *big.Int) (*types.Transaction, error) {
	return _OChainPortal.Contract.Unstake(&_OChainPortal.TransactOpts, validatorId)
}

// Unstake is a paid mutator transaction binding the contract method 0x2e17de78.
//
// Solidity: function unstake(uint256 validatorId) returns()
func (_OChainPortal *OChainPortalTransactorSession) Unstake(validatorId *big.Int) (*types.Transaction, error) {
	return _OChainPortal.Contract.Unstake(&_OChainPortal.TransactOpts, validatorId)
}

// WithdrawalRequest is a paid mutator transaction binding the contract method 0x9582601e.
//
// Solidity: function withdrawalRequest(address receiver, uint256 nonce, uint256 amount, uint256[] validatorIds, uint256[] signatures) returns()
func (_OChainPortal *OChainPortalTransactor) WithdrawalRequest(opts *bind.TransactOpts, receiver common.Address, nonce *big.Int, amount *big.Int, validatorIds []*big.Int, signatures []*big.Int) (*types.Transaction, error) {
	return _OChainPortal.contract.Transact(opts, "withdrawalRequest", receiver, nonce, amount, validatorIds, signatures)
}

// WithdrawalRequest is a paid mutator transaction binding the contract method 0x9582601e.
//
// Solidity: function withdrawalRequest(address receiver, uint256 nonce, uint256 amount, uint256[] validatorIds, uint256[] signatures) returns()
func (_OChainPortal *OChainPortalSession) WithdrawalRequest(receiver common.Address, nonce *big.Int, amount *big.Int, validatorIds []*big.Int, signatures []*big.Int) (*types.Transaction, error) {
	return _OChainPortal.Contract.WithdrawalRequest(&_OChainPortal.TransactOpts, receiver, nonce, amount, validatorIds, signatures)
}

// WithdrawalRequest is a paid mutator transaction binding the contract method 0x9582601e.
//
// Solidity: function withdrawalRequest(address receiver, uint256 nonce, uint256 amount, uint256[] validatorIds, uint256[] signatures) returns()
func (_OChainPortal *OChainPortalTransactorSession) WithdrawalRequest(receiver common.Address, nonce *big.Int, amount *big.Int, validatorIds []*big.Int, signatures []*big.Int) (*types.Transaction, error) {
	return _OChainPortal.Contract.WithdrawalRequest(&_OChainPortal.TransactOpts, receiver, nonce, amount, validatorIds, signatures)
}

// OChainPortalOChainNewValidatorIterator is returned from FilterOChainNewValidator and is used to iterate over the raw logs and unpacked data for OChainNewValidator events raised by the OChainPortal contract.
type OChainPortalOChainNewValidatorIterator struct {
	Event *OChainPortalOChainNewValidator // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *OChainPortalOChainNewValidatorIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OChainPortalOChainNewValidator)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(OChainPortalOChainNewValidator)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *OChainPortalOChainNewValidatorIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *OChainPortalOChainNewValidatorIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// OChainPortalOChainNewValidator represents a OChainNewValidator event raised by the OChainPortal contract.
type OChainPortalOChainNewValidator struct {
	ValidatorId *big.Int
	Stacker     common.Address
	Validator   common.Address
	PublicKey   string
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterOChainNewValidator is a free log retrieval operation binding the contract event 0x409249b7bd27c9e19c97851ee0577d3bbca2dca1a6a960ca25f4c3ffaa58c41d.
//
// Solidity: event OChainNewValidator(uint256 indexed validatorId, address stacker, address validator, string publicKey)
func (_OChainPortal *OChainPortalFilterer) FilterOChainNewValidator(opts *bind.FilterOpts, validatorId []*big.Int) (*OChainPortalOChainNewValidatorIterator, error) {

	var validatorIdRule []interface{}
	for _, validatorIdItem := range validatorId {
		validatorIdRule = append(validatorIdRule, validatorIdItem)
	}

	logs, sub, err := _OChainPortal.contract.FilterLogs(opts, "OChainNewValidator", validatorIdRule)
	if err != nil {
		return nil, err
	}
	return &OChainPortalOChainNewValidatorIterator{contract: _OChainPortal.contract, event: "OChainNewValidator", logs: logs, sub: sub}, nil
}

// WatchOChainNewValidator is a free log subscription operation binding the contract event 0x409249b7bd27c9e19c97851ee0577d3bbca2dca1a6a960ca25f4c3ffaa58c41d.
//
// Solidity: event OChainNewValidator(uint256 indexed validatorId, address stacker, address validator, string publicKey)
func (_OChainPortal *OChainPortalFilterer) WatchOChainNewValidator(opts *bind.WatchOpts, sink chan<- *OChainPortalOChainNewValidator, validatorId []*big.Int) (event.Subscription, error) {

	var validatorIdRule []interface{}
	for _, validatorIdItem := range validatorId {
		validatorIdRule = append(validatorIdRule, validatorIdItem)
	}

	logs, sub, err := _OChainPortal.contract.WatchLogs(opts, "OChainNewValidator", validatorIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(OChainPortalOChainNewValidator)
				if err := _OChainPortal.contract.UnpackLog(event, "OChainNewValidator", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseOChainNewValidator is a log parse operation binding the contract event 0x409249b7bd27c9e19c97851ee0577d3bbca2dca1a6a960ca25f4c3ffaa58c41d.
//
// Solidity: event OChainNewValidator(uint256 indexed validatorId, address stacker, address validator, string publicKey)
func (_OChainPortal *OChainPortalFilterer) ParseOChainNewValidator(log types.Log) (*OChainPortalOChainNewValidator, error) {
	event := new(OChainPortalOChainNewValidator)
	if err := _OChainPortal.contract.UnpackLog(event, "OChainNewValidator", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// OChainPortalOChainRemoveValidatorIterator is returned from FilterOChainRemoveValidator and is used to iterate over the raw logs and unpacked data for OChainRemoveValidator events raised by the OChainPortal contract.
type OChainPortalOChainRemoveValidatorIterator struct {
	Event *OChainPortalOChainRemoveValidator // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *OChainPortalOChainRemoveValidatorIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OChainPortalOChainRemoveValidator)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(OChainPortalOChainRemoveValidator)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *OChainPortalOChainRemoveValidatorIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *OChainPortalOChainRemoveValidatorIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// OChainPortalOChainRemoveValidator represents a OChainRemoveValidator event raised by the OChainPortal contract.
type OChainPortalOChainRemoveValidator struct {
	ValidatorId *big.Int
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterOChainRemoveValidator is a free log retrieval operation binding the contract event 0xc87611674cb43117db3e346bdccb9de664b84b5de7e00ac02085308c59d8e02f.
//
// Solidity: event OChainRemoveValidator(uint256 indexed validatorId)
func (_OChainPortal *OChainPortalFilterer) FilterOChainRemoveValidator(opts *bind.FilterOpts, validatorId []*big.Int) (*OChainPortalOChainRemoveValidatorIterator, error) {

	var validatorIdRule []interface{}
	for _, validatorIdItem := range validatorId {
		validatorIdRule = append(validatorIdRule, validatorIdItem)
	}

	logs, sub, err := _OChainPortal.contract.FilterLogs(opts, "OChainRemoveValidator", validatorIdRule)
	if err != nil {
		return nil, err
	}
	return &OChainPortalOChainRemoveValidatorIterator{contract: _OChainPortal.contract, event: "OChainRemoveValidator", logs: logs, sub: sub}, nil
}

// WatchOChainRemoveValidator is a free log subscription operation binding the contract event 0xc87611674cb43117db3e346bdccb9de664b84b5de7e00ac02085308c59d8e02f.
//
// Solidity: event OChainRemoveValidator(uint256 indexed validatorId)
func (_OChainPortal *OChainPortalFilterer) WatchOChainRemoveValidator(opts *bind.WatchOpts, sink chan<- *OChainPortalOChainRemoveValidator, validatorId []*big.Int) (event.Subscription, error) {

	var validatorIdRule []interface{}
	for _, validatorIdItem := range validatorId {
		validatorIdRule = append(validatorIdRule, validatorIdItem)
	}

	logs, sub, err := _OChainPortal.contract.WatchLogs(opts, "OChainRemoveValidator", validatorIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(OChainPortalOChainRemoveValidator)
				if err := _OChainPortal.contract.UnpackLog(event, "OChainRemoveValidator", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseOChainRemoveValidator is a log parse operation binding the contract event 0xc87611674cb43117db3e346bdccb9de664b84b5de7e00ac02085308c59d8e02f.
//
// Solidity: event OChainRemoveValidator(uint256 indexed validatorId)
func (_OChainPortal *OChainPortalFilterer) ParseOChainRemoveValidator(log types.Log) (*OChainPortalOChainRemoveValidator, error) {
	event := new(OChainPortalOChainRemoveValidator)
	if err := _OChainPortal.contract.UnpackLog(event, "OChainRemoveValidator", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// OChainPortalOChainTokenDepositIterator is returned from FilterOChainTokenDeposit and is used to iterate over the raw logs and unpacked data for OChainTokenDeposit events raised by the OChainPortal contract.
type OChainPortalOChainTokenDepositIterator struct {
	Event *OChainPortalOChainTokenDeposit // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *OChainPortalOChainTokenDepositIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OChainPortalOChainTokenDeposit)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(OChainPortalOChainTokenDeposit)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *OChainPortalOChainTokenDepositIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *OChainPortalOChainTokenDepositIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// OChainPortalOChainTokenDeposit represents a OChainTokenDeposit event raised by the OChainPortal contract.
type OChainPortalOChainTokenDeposit struct {
	Receiver common.Address
	Amount   *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterOChainTokenDeposit is a free log retrieval operation binding the contract event 0x65ad5f5e2c996a7dd0047b15db95807104f818650054d36a2ee6ca27eed441f4.
//
// Solidity: event OChainTokenDeposit(address indexed receiver, uint256 amount)
func (_OChainPortal *OChainPortalFilterer) FilterOChainTokenDeposit(opts *bind.FilterOpts, receiver []common.Address) (*OChainPortalOChainTokenDepositIterator, error) {

	var receiverRule []interface{}
	for _, receiverItem := range receiver {
		receiverRule = append(receiverRule, receiverItem)
	}

	logs, sub, err := _OChainPortal.contract.FilterLogs(opts, "OChainTokenDeposit", receiverRule)
	if err != nil {
		return nil, err
	}
	return &OChainPortalOChainTokenDepositIterator{contract: _OChainPortal.contract, event: "OChainTokenDeposit", logs: logs, sub: sub}, nil
}

// WatchOChainTokenDeposit is a free log subscription operation binding the contract event 0x65ad5f5e2c996a7dd0047b15db95807104f818650054d36a2ee6ca27eed441f4.
//
// Solidity: event OChainTokenDeposit(address indexed receiver, uint256 amount)
func (_OChainPortal *OChainPortalFilterer) WatchOChainTokenDeposit(opts *bind.WatchOpts, sink chan<- *OChainPortalOChainTokenDeposit, receiver []common.Address) (event.Subscription, error) {

	var receiverRule []interface{}
	for _, receiverItem := range receiver {
		receiverRule = append(receiverRule, receiverItem)
	}

	logs, sub, err := _OChainPortal.contract.WatchLogs(opts, "OChainTokenDeposit", receiverRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(OChainPortalOChainTokenDeposit)
				if err := _OChainPortal.contract.UnpackLog(event, "OChainTokenDeposit", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseOChainTokenDeposit is a log parse operation binding the contract event 0x65ad5f5e2c996a7dd0047b15db95807104f818650054d36a2ee6ca27eed441f4.
//
// Solidity: event OChainTokenDeposit(address indexed receiver, uint256 amount)
func (_OChainPortal *OChainPortalFilterer) ParseOChainTokenDeposit(log types.Log) (*OChainPortalOChainTokenDeposit, error) {
	event := new(OChainPortalOChainTokenDeposit)
	if err := _OChainPortal.contract.UnpackLog(event, "OChainTokenDeposit", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// OChainPortalOChainTokenWithdrawalRequestContestedIterator is returned from FilterOChainTokenWithdrawalRequestContested and is used to iterate over the raw logs and unpacked data for OChainTokenWithdrawalRequestContested events raised by the OChainPortal contract.
type OChainPortalOChainTokenWithdrawalRequestContestedIterator struct {
	Event *OChainPortalOChainTokenWithdrawalRequestContested // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *OChainPortalOChainTokenWithdrawalRequestContestedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OChainPortalOChainTokenWithdrawalRequestContested)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(OChainPortalOChainTokenWithdrawalRequestContested)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *OChainPortalOChainTokenWithdrawalRequestContestedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *OChainPortalOChainTokenWithdrawalRequestContestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// OChainPortalOChainTokenWithdrawalRequestContested represents a OChainTokenWithdrawalRequestContested event raised by the OChainPortal contract.
type OChainPortalOChainTokenWithdrawalRequestContested struct {
	Nonce *big.Int
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterOChainTokenWithdrawalRequestContested is a free log retrieval operation binding the contract event 0x304eeb90c4d61fec6d6dbf54316009fc718ac82c4a7baaf965ba7e0d5905f5e4.
//
// Solidity: event OChainTokenWithdrawalRequestContested(uint256 nonce)
func (_OChainPortal *OChainPortalFilterer) FilterOChainTokenWithdrawalRequestContested(opts *bind.FilterOpts) (*OChainPortalOChainTokenWithdrawalRequestContestedIterator, error) {

	logs, sub, err := _OChainPortal.contract.FilterLogs(opts, "OChainTokenWithdrawalRequestContested")
	if err != nil {
		return nil, err
	}
	return &OChainPortalOChainTokenWithdrawalRequestContestedIterator{contract: _OChainPortal.contract, event: "OChainTokenWithdrawalRequestContested", logs: logs, sub: sub}, nil
}

// WatchOChainTokenWithdrawalRequestContested is a free log subscription operation binding the contract event 0x304eeb90c4d61fec6d6dbf54316009fc718ac82c4a7baaf965ba7e0d5905f5e4.
//
// Solidity: event OChainTokenWithdrawalRequestContested(uint256 nonce)
func (_OChainPortal *OChainPortalFilterer) WatchOChainTokenWithdrawalRequestContested(opts *bind.WatchOpts, sink chan<- *OChainPortalOChainTokenWithdrawalRequestContested) (event.Subscription, error) {

	logs, sub, err := _OChainPortal.contract.WatchLogs(opts, "OChainTokenWithdrawalRequestContested")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(OChainPortalOChainTokenWithdrawalRequestContested)
				if err := _OChainPortal.contract.UnpackLog(event, "OChainTokenWithdrawalRequestContested", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseOChainTokenWithdrawalRequestContested is a log parse operation binding the contract event 0x304eeb90c4d61fec6d6dbf54316009fc718ac82c4a7baaf965ba7e0d5905f5e4.
//
// Solidity: event OChainTokenWithdrawalRequestContested(uint256 nonce)
func (_OChainPortal *OChainPortalFilterer) ParseOChainTokenWithdrawalRequestContested(log types.Log) (*OChainPortalOChainTokenWithdrawalRequestContested, error) {
	event := new(OChainPortalOChainTokenWithdrawalRequestContested)
	if err := _OChainPortal.contract.UnpackLog(event, "OChainTokenWithdrawalRequestContested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// OChainPortalOChainTokenWithdrawalRequestExecutedIterator is returned from FilterOChainTokenWithdrawalRequestExecuted and is used to iterate over the raw logs and unpacked data for OChainTokenWithdrawalRequestExecuted events raised by the OChainPortal contract.
type OChainPortalOChainTokenWithdrawalRequestExecutedIterator struct {
	Event *OChainPortalOChainTokenWithdrawalRequestExecuted // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *OChainPortalOChainTokenWithdrawalRequestExecutedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OChainPortalOChainTokenWithdrawalRequestExecuted)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(OChainPortalOChainTokenWithdrawalRequestExecuted)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *OChainPortalOChainTokenWithdrawalRequestExecutedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *OChainPortalOChainTokenWithdrawalRequestExecutedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// OChainPortalOChainTokenWithdrawalRequestExecuted represents a OChainTokenWithdrawalRequestExecuted event raised by the OChainPortal contract.
type OChainPortalOChainTokenWithdrawalRequestExecuted struct {
	Nonce *big.Int
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterOChainTokenWithdrawalRequestExecuted is a free log retrieval operation binding the contract event 0x730a5bb623337c88e68b5d82b90939b77f0d2871f86aec602bff04f5c96a4063.
//
// Solidity: event OChainTokenWithdrawalRequestExecuted(uint256 nonce)
func (_OChainPortal *OChainPortalFilterer) FilterOChainTokenWithdrawalRequestExecuted(opts *bind.FilterOpts) (*OChainPortalOChainTokenWithdrawalRequestExecutedIterator, error) {

	logs, sub, err := _OChainPortal.contract.FilterLogs(opts, "OChainTokenWithdrawalRequestExecuted")
	if err != nil {
		return nil, err
	}
	return &OChainPortalOChainTokenWithdrawalRequestExecutedIterator{contract: _OChainPortal.contract, event: "OChainTokenWithdrawalRequestExecuted", logs: logs, sub: sub}, nil
}

// WatchOChainTokenWithdrawalRequestExecuted is a free log subscription operation binding the contract event 0x730a5bb623337c88e68b5d82b90939b77f0d2871f86aec602bff04f5c96a4063.
//
// Solidity: event OChainTokenWithdrawalRequestExecuted(uint256 nonce)
func (_OChainPortal *OChainPortalFilterer) WatchOChainTokenWithdrawalRequestExecuted(opts *bind.WatchOpts, sink chan<- *OChainPortalOChainTokenWithdrawalRequestExecuted) (event.Subscription, error) {

	logs, sub, err := _OChainPortal.contract.WatchLogs(opts, "OChainTokenWithdrawalRequestExecuted")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(OChainPortalOChainTokenWithdrawalRequestExecuted)
				if err := _OChainPortal.contract.UnpackLog(event, "OChainTokenWithdrawalRequestExecuted", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseOChainTokenWithdrawalRequestExecuted is a log parse operation binding the contract event 0x730a5bb623337c88e68b5d82b90939b77f0d2871f86aec602bff04f5c96a4063.
//
// Solidity: event OChainTokenWithdrawalRequestExecuted(uint256 nonce)
func (_OChainPortal *OChainPortalFilterer) ParseOChainTokenWithdrawalRequestExecuted(log types.Log) (*OChainPortalOChainTokenWithdrawalRequestExecuted, error) {
	event := new(OChainPortalOChainTokenWithdrawalRequestExecuted)
	if err := _OChainPortal.contract.UnpackLog(event, "OChainTokenWithdrawalRequestExecuted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// OChainPortalOChainTokenWithdrawalRequestedIterator is returned from FilterOChainTokenWithdrawalRequested and is used to iterate over the raw logs and unpacked data for OChainTokenWithdrawalRequested events raised by the OChainPortal contract.
type OChainPortalOChainTokenWithdrawalRequestedIterator struct {
	Event *OChainPortalOChainTokenWithdrawalRequested // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *OChainPortalOChainTokenWithdrawalRequestedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OChainPortalOChainTokenWithdrawalRequested)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(OChainPortalOChainTokenWithdrawalRequested)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *OChainPortalOChainTokenWithdrawalRequestedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *OChainPortalOChainTokenWithdrawalRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// OChainPortalOChainTokenWithdrawalRequested represents a OChainTokenWithdrawalRequested event raised by the OChainPortal contract.
type OChainPortalOChainTokenWithdrawalRequested struct {
	Sender  common.Address
	Nonce   *big.Int
	Amount  *big.Int
	Signers []*big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterOChainTokenWithdrawalRequested is a free log retrieval operation binding the contract event 0x7581f25582ecb69233c1174995be6e1300c86e62a7fbe1ee4e98a0e4ba577b7d.
//
// Solidity: event OChainTokenWithdrawalRequested(address indexed sender, uint256 nonce, uint256 amount, uint256[] signers)
func (_OChainPortal *OChainPortalFilterer) FilterOChainTokenWithdrawalRequested(opts *bind.FilterOpts, sender []common.Address) (*OChainPortalOChainTokenWithdrawalRequestedIterator, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _OChainPortal.contract.FilterLogs(opts, "OChainTokenWithdrawalRequested", senderRule)
	if err != nil {
		return nil, err
	}
	return &OChainPortalOChainTokenWithdrawalRequestedIterator{contract: _OChainPortal.contract, event: "OChainTokenWithdrawalRequested", logs: logs, sub: sub}, nil
}

// WatchOChainTokenWithdrawalRequested is a free log subscription operation binding the contract event 0x7581f25582ecb69233c1174995be6e1300c86e62a7fbe1ee4e98a0e4ba577b7d.
//
// Solidity: event OChainTokenWithdrawalRequested(address indexed sender, uint256 nonce, uint256 amount, uint256[] signers)
func (_OChainPortal *OChainPortalFilterer) WatchOChainTokenWithdrawalRequested(opts *bind.WatchOpts, sink chan<- *OChainPortalOChainTokenWithdrawalRequested, sender []common.Address) (event.Subscription, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _OChainPortal.contract.WatchLogs(opts, "OChainTokenWithdrawalRequested", senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(OChainPortalOChainTokenWithdrawalRequested)
				if err := _OChainPortal.contract.UnpackLog(event, "OChainTokenWithdrawalRequested", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseOChainTokenWithdrawalRequested is a log parse operation binding the contract event 0x7581f25582ecb69233c1174995be6e1300c86e62a7fbe1ee4e98a0e4ba577b7d.
//
// Solidity: event OChainTokenWithdrawalRequested(address indexed sender, uint256 nonce, uint256 amount, uint256[] signers)
func (_OChainPortal *OChainPortalFilterer) ParseOChainTokenWithdrawalRequested(log types.Log) (*OChainPortalOChainTokenWithdrawalRequested, error) {
	event := new(OChainPortalOChainTokenWithdrawalRequested)
	if err := _OChainPortal.contract.UnpackLog(event, "OChainTokenWithdrawalRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// OChainPortalOChainUnstackSucceedIterator is returned from FilterOChainUnstackSucceed and is used to iterate over the raw logs and unpacked data for OChainUnstackSucceed events raised by the OChainPortal contract.
type OChainPortalOChainUnstackSucceedIterator struct {
	Event *OChainPortalOChainUnstackSucceed // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *OChainPortalOChainUnstackSucceedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OChainPortalOChainUnstackSucceed)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(OChainPortalOChainUnstackSucceed)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *OChainPortalOChainUnstackSucceedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *OChainPortalOChainUnstackSucceedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// OChainPortalOChainUnstackSucceed represents a OChainUnstackSucceed event raised by the OChainPortal contract.
type OChainPortalOChainUnstackSucceed struct {
	ValidatorId *big.Int
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterOChainUnstackSucceed is a free log retrieval operation binding the contract event 0x98de3b833cf30242ee4e00831c2905b2afb66ce3a1339645b8f099353fb3e2a0.
//
// Solidity: event OChainUnstackSucceed(uint256 indexed validatorId)
func (_OChainPortal *OChainPortalFilterer) FilterOChainUnstackSucceed(opts *bind.FilterOpts, validatorId []*big.Int) (*OChainPortalOChainUnstackSucceedIterator, error) {

	var validatorIdRule []interface{}
	for _, validatorIdItem := range validatorId {
		validatorIdRule = append(validatorIdRule, validatorIdItem)
	}

	logs, sub, err := _OChainPortal.contract.FilterLogs(opts, "OChainUnstackSucceed", validatorIdRule)
	if err != nil {
		return nil, err
	}
	return &OChainPortalOChainUnstackSucceedIterator{contract: _OChainPortal.contract, event: "OChainUnstackSucceed", logs: logs, sub: sub}, nil
}

// WatchOChainUnstackSucceed is a free log subscription operation binding the contract event 0x98de3b833cf30242ee4e00831c2905b2afb66ce3a1339645b8f099353fb3e2a0.
//
// Solidity: event OChainUnstackSucceed(uint256 indexed validatorId)
func (_OChainPortal *OChainPortalFilterer) WatchOChainUnstackSucceed(opts *bind.WatchOpts, sink chan<- *OChainPortalOChainUnstackSucceed, validatorId []*big.Int) (event.Subscription, error) {

	var validatorIdRule []interface{}
	for _, validatorIdItem := range validatorId {
		validatorIdRule = append(validatorIdRule, validatorIdItem)
	}

	logs, sub, err := _OChainPortal.contract.WatchLogs(opts, "OChainUnstackSucceed", validatorIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(OChainPortalOChainUnstackSucceed)
				if err := _OChainPortal.contract.UnpackLog(event, "OChainUnstackSucceed", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseOChainUnstackSucceed is a log parse operation binding the contract event 0x98de3b833cf30242ee4e00831c2905b2afb66ce3a1339645b8f099353fb3e2a0.
//
// Solidity: event OChainUnstackSucceed(uint256 indexed validatorId)
func (_OChainPortal *OChainPortalFilterer) ParseOChainUnstackSucceed(log types.Log) (*OChainPortalOChainUnstackSucceed, error) {
	event := new(OChainPortalOChainUnstackSucceed)
	if err := _OChainPortal.contract.UnpackLog(event, "OChainUnstackSucceed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
