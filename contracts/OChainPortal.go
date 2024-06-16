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

// DiamondInitialization is an auto generated low-level Go binding around an user-defined struct.
type DiamondInitialization struct {
	InitContract common.Address
	InitData     []byte
}

// IDiamondCutFacetCut is an auto generated low-level Go binding around an user-defined struct.
type IDiamondCutFacetCut struct {
	FacetAddress      common.Address
	Action            uint8
	FunctionSelectors [][4]byte
}

// IDiamondLoupeFacet is an auto generated low-level Go binding around an user-defined struct.
type IDiamondLoupeFacet struct {
	FacetAddress      common.Address
	FunctionSelectors [][4]byte
}

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
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_contractOwner\",\"type\":\"address\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"facetAddress\",\"type\":\"address\"},{\"internalType\":\"enumIDiamondCut.FacetCutAction\",\"name\":\"action\",\"type\":\"uint8\"},{\"internalType\":\"bytes4[]\",\"name\":\"functionSelectors\",\"type\":\"bytes4[]\"}],\"internalType\":\"structIDiamondCut.FacetCut[]\",\"name\":\"_diamondCut\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"initContract\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"initData\",\"type\":\"bytes\"}],\"internalType\":\"structDiamond.Initialization[]\",\"name\":\"_initializations\",\"type\":\"tuple[]\"}],\"stateMutability\":\"payable\",\"type\":\"constructor\"},{\"stateMutability\":\"payable\",\"type\":\"fallback\"},{\"stateMutability\":\"payable\",\"type\":\"receive\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"target\",\"type\":\"address\"}],\"name\":\"AddressEmptyCode\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"AddressInsufficientBalance\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"FailedInnerCall\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NotValidator\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"QuorumNotReached\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"SafeERC20FailedOperation\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"WithdrawalAlreadyExecuted\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"WithdrawalAlreadySigned\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"WithdrawalCanceled\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"WithdrawalNonceInvalid\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"OChainTokenDeposited\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"signer\",\"type\":\"address\"}],\"name\":\"OChainTokenWithdrawalCanceled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"}],\"name\":\"OChainTokenWithdrawalExecuted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"signer\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"OChainTokenWithdrawalSigned\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"USDDeposited\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"validatorId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"approveWithdraw\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"validatorId\",\"type\":\"uint256\"}],\"name\":\"cancelWithdraw\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"deposit\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"depositUSD\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"}],\"name\":\"withdraw\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"MaxValidatorReached\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NotAuthorized\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UnstakeProcessNotAvailable\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UnstakeProcessNotEnded\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"validatorId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"stacker\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"validator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"publicKey\",\"type\":\"string\"}],\"name\":\"OChainNewValidator\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"validatorId\",\"type\":\"uint256\"}],\"name\":\"OChainRemoveValidator\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"validatorId\",\"type\":\"uint256\"}],\"name\":\"OChainUnstackSucceed\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"latestUpdateAt\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"stacker\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"validator\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"pubkey\",\"type\":\"string\"}],\"name\":\"stake\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"validatorId\",\"type\":\"uint256\"}],\"name\":\"startUnstakeProcess\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"validatorId\",\"type\":\"uint256\"}],\"name\":\"unstake\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"validatorId\",\"type\":\"uint256\"}],\"name\":\"validatorInfo\",\"outputs\":[{\"components\":[{\"internalType\":\"bool\",\"name\":\"enabled\",\"type\":\"bool\"},{\"internalType\":\"address\",\"name\":\"depositor\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"validator\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"publicKey\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"withdrawRequestDate\",\"type\":\"uint256\"}],\"internalType\":\"structLibDiamond.ValidatorStacking\",\"name\":\"_validator\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"validatorNetworkInfo\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"maxValidators\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"validatorsMaxIndex\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"validatorsLength\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"anonymous\":false,\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"facetAddress\",\"type\":\"address\"},{\"internalType\":\"enumIDiamondCut.FacetCutAction\",\"name\":\"action\",\"type\":\"uint8\"},{\"internalType\":\"bytes4[]\",\"name\":\"functionSelectors\",\"type\":\"bytes4[]\"}],\"indexed\":false,\"internalType\":\"structIDiamondCut.FacetCut[]\",\"name\":\"_diamondCut\",\"type\":\"tuple[]\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"_init\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"_calldata\",\"type\":\"bytes\"}],\"name\":\"DiamondCut\",\"type\":\"event\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"facetAddress\",\"type\":\"address\"},{\"internalType\":\"enumIDiamondCut.FacetCutAction\",\"name\":\"action\",\"type\":\"uint8\"},{\"internalType\":\"bytes4[]\",\"name\":\"functionSelectors\",\"type\":\"bytes4[]\"}],\"internalType\":\"structIDiamondCut.FacetCut[]\",\"name\":\"_diamondCut\",\"type\":\"tuple[]\"},{\"internalType\":\"address\",\"name\":\"_init\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"_calldata\",\"type\":\"bytes\"}],\"name\":\"diamondCut\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"owner_\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"_functionSelector\",\"type\":\"bytes4\"}],\"name\":\"facetAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"facetAddress_\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"facetAddresses\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"facetAddresses_\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_facet\",\"type\":\"address\"}],\"name\":\"facetFunctionSelectors\",\"outputs\":[{\"internalType\":\"bytes4[]\",\"name\":\"facetFunctionSelectors_\",\"type\":\"bytes4[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"facets\",\"outputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"facetAddress\",\"type\":\"address\"},{\"internalType\":\"bytes4[]\",\"name\":\"functionSelectors\",\"type\":\"bytes4[]\"}],\"internalType\":\"structIDiamondLoupe.Facet[]\",\"name\":\"facets_\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"_interfaceId\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
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

// FacetAddress is a free data retrieval call binding the contract method 0xcdffacc6.
//
// Solidity: function facetAddress(bytes4 _functionSelector) view returns(address facetAddress_)
func (_OChainPortal *OChainPortalCaller) FacetAddress(opts *bind.CallOpts, _functionSelector [4]byte) (common.Address, error) {
	var out []interface{}
	err := _OChainPortal.contract.Call(opts, &out, "facetAddress", _functionSelector)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// FacetAddress is a free data retrieval call binding the contract method 0xcdffacc6.
//
// Solidity: function facetAddress(bytes4 _functionSelector) view returns(address facetAddress_)
func (_OChainPortal *OChainPortalSession) FacetAddress(_functionSelector [4]byte) (common.Address, error) {
	return _OChainPortal.Contract.FacetAddress(&_OChainPortal.CallOpts, _functionSelector)
}

// FacetAddress is a free data retrieval call binding the contract method 0xcdffacc6.
//
// Solidity: function facetAddress(bytes4 _functionSelector) view returns(address facetAddress_)
func (_OChainPortal *OChainPortalCallerSession) FacetAddress(_functionSelector [4]byte) (common.Address, error) {
	return _OChainPortal.Contract.FacetAddress(&_OChainPortal.CallOpts, _functionSelector)
}

// FacetAddresses is a free data retrieval call binding the contract method 0x52ef6b2c.
//
// Solidity: function facetAddresses() view returns(address[] facetAddresses_)
func (_OChainPortal *OChainPortalCaller) FacetAddresses(opts *bind.CallOpts) ([]common.Address, error) {
	var out []interface{}
	err := _OChainPortal.contract.Call(opts, &out, "facetAddresses")

	if err != nil {
		return *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)

	return out0, err

}

// FacetAddresses is a free data retrieval call binding the contract method 0x52ef6b2c.
//
// Solidity: function facetAddresses() view returns(address[] facetAddresses_)
func (_OChainPortal *OChainPortalSession) FacetAddresses() ([]common.Address, error) {
	return _OChainPortal.Contract.FacetAddresses(&_OChainPortal.CallOpts)
}

// FacetAddresses is a free data retrieval call binding the contract method 0x52ef6b2c.
//
// Solidity: function facetAddresses() view returns(address[] facetAddresses_)
func (_OChainPortal *OChainPortalCallerSession) FacetAddresses() ([]common.Address, error) {
	return _OChainPortal.Contract.FacetAddresses(&_OChainPortal.CallOpts)
}

// FacetFunctionSelectors is a free data retrieval call binding the contract method 0xadfca15e.
//
// Solidity: function facetFunctionSelectors(address _facet) view returns(bytes4[] facetFunctionSelectors_)
func (_OChainPortal *OChainPortalCaller) FacetFunctionSelectors(opts *bind.CallOpts, _facet common.Address) ([][4]byte, error) {
	var out []interface{}
	err := _OChainPortal.contract.Call(opts, &out, "facetFunctionSelectors", _facet)

	if err != nil {
		return *new([][4]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([][4]byte)).(*[][4]byte)

	return out0, err

}

// FacetFunctionSelectors is a free data retrieval call binding the contract method 0xadfca15e.
//
// Solidity: function facetFunctionSelectors(address _facet) view returns(bytes4[] facetFunctionSelectors_)
func (_OChainPortal *OChainPortalSession) FacetFunctionSelectors(_facet common.Address) ([][4]byte, error) {
	return _OChainPortal.Contract.FacetFunctionSelectors(&_OChainPortal.CallOpts, _facet)
}

// FacetFunctionSelectors is a free data retrieval call binding the contract method 0xadfca15e.
//
// Solidity: function facetFunctionSelectors(address _facet) view returns(bytes4[] facetFunctionSelectors_)
func (_OChainPortal *OChainPortalCallerSession) FacetFunctionSelectors(_facet common.Address) ([][4]byte, error) {
	return _OChainPortal.Contract.FacetFunctionSelectors(&_OChainPortal.CallOpts, _facet)
}

// Facets is a free data retrieval call binding the contract method 0x7a0ed627.
//
// Solidity: function facets() view returns((address,bytes4[])[] facets_)
func (_OChainPortal *OChainPortalCaller) Facets(opts *bind.CallOpts) ([]IDiamondLoupeFacet, error) {
	var out []interface{}
	err := _OChainPortal.contract.Call(opts, &out, "facets")

	if err != nil {
		return *new([]IDiamondLoupeFacet), err
	}

	out0 := *abi.ConvertType(out[0], new([]IDiamondLoupeFacet)).(*[]IDiamondLoupeFacet)

	return out0, err

}

// Facets is a free data retrieval call binding the contract method 0x7a0ed627.
//
// Solidity: function facets() view returns((address,bytes4[])[] facets_)
func (_OChainPortal *OChainPortalSession) Facets() ([]IDiamondLoupeFacet, error) {
	return _OChainPortal.Contract.Facets(&_OChainPortal.CallOpts)
}

// Facets is a free data retrieval call binding the contract method 0x7a0ed627.
//
// Solidity: function facets() view returns((address,bytes4[])[] facets_)
func (_OChainPortal *OChainPortalCallerSession) Facets() ([]IDiamondLoupeFacet, error) {
	return _OChainPortal.Contract.Facets(&_OChainPortal.CallOpts)
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

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address owner_)
func (_OChainPortal *OChainPortalCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _OChainPortal.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address owner_)
func (_OChainPortal *OChainPortalSession) Owner() (common.Address, error) {
	return _OChainPortal.Contract.Owner(&_OChainPortal.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address owner_)
func (_OChainPortal *OChainPortalCallerSession) Owner() (common.Address, error) {
	return _OChainPortal.Contract.Owner(&_OChainPortal.CallOpts)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 _interfaceId) view returns(bool)
func (_OChainPortal *OChainPortalCaller) SupportsInterface(opts *bind.CallOpts, _interfaceId [4]byte) (bool, error) {
	var out []interface{}
	err := _OChainPortal.contract.Call(opts, &out, "supportsInterface", _interfaceId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 _interfaceId) view returns(bool)
func (_OChainPortal *OChainPortalSession) SupportsInterface(_interfaceId [4]byte) (bool, error) {
	return _OChainPortal.Contract.SupportsInterface(&_OChainPortal.CallOpts, _interfaceId)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 _interfaceId) view returns(bool)
func (_OChainPortal *OChainPortalCallerSession) SupportsInterface(_interfaceId [4]byte) (bool, error) {
	return _OChainPortal.Contract.SupportsInterface(&_OChainPortal.CallOpts, _interfaceId)
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

// ApproveWithdraw is a paid mutator transaction binding the contract method 0xbc8bb253.
//
// Solidity: function approveWithdraw(uint256 nonce, uint256 validatorId, address receiver, uint256 amount) returns(bool)
func (_OChainPortal *OChainPortalTransactor) ApproveWithdraw(opts *bind.TransactOpts, nonce *big.Int, validatorId *big.Int, receiver common.Address, amount *big.Int) (*types.Transaction, error) {
	return _OChainPortal.contract.Transact(opts, "approveWithdraw", nonce, validatorId, receiver, amount)
}

// ApproveWithdraw is a paid mutator transaction binding the contract method 0xbc8bb253.
//
// Solidity: function approveWithdraw(uint256 nonce, uint256 validatorId, address receiver, uint256 amount) returns(bool)
func (_OChainPortal *OChainPortalSession) ApproveWithdraw(nonce *big.Int, validatorId *big.Int, receiver common.Address, amount *big.Int) (*types.Transaction, error) {
	return _OChainPortal.Contract.ApproveWithdraw(&_OChainPortal.TransactOpts, nonce, validatorId, receiver, amount)
}

// ApproveWithdraw is a paid mutator transaction binding the contract method 0xbc8bb253.
//
// Solidity: function approveWithdraw(uint256 nonce, uint256 validatorId, address receiver, uint256 amount) returns(bool)
func (_OChainPortal *OChainPortalTransactorSession) ApproveWithdraw(nonce *big.Int, validatorId *big.Int, receiver common.Address, amount *big.Int) (*types.Transaction, error) {
	return _OChainPortal.Contract.ApproveWithdraw(&_OChainPortal.TransactOpts, nonce, validatorId, receiver, amount)
}

// CancelWithdraw is a paid mutator transaction binding the contract method 0x005919c7.
//
// Solidity: function cancelWithdraw(uint256 nonce, uint256 validatorId) returns(bool)
func (_OChainPortal *OChainPortalTransactor) CancelWithdraw(opts *bind.TransactOpts, nonce *big.Int, validatorId *big.Int) (*types.Transaction, error) {
	return _OChainPortal.contract.Transact(opts, "cancelWithdraw", nonce, validatorId)
}

// CancelWithdraw is a paid mutator transaction binding the contract method 0x005919c7.
//
// Solidity: function cancelWithdraw(uint256 nonce, uint256 validatorId) returns(bool)
func (_OChainPortal *OChainPortalSession) CancelWithdraw(nonce *big.Int, validatorId *big.Int) (*types.Transaction, error) {
	return _OChainPortal.Contract.CancelWithdraw(&_OChainPortal.TransactOpts, nonce, validatorId)
}

// CancelWithdraw is a paid mutator transaction binding the contract method 0x005919c7.
//
// Solidity: function cancelWithdraw(uint256 nonce, uint256 validatorId) returns(bool)
func (_OChainPortal *OChainPortalTransactorSession) CancelWithdraw(nonce *big.Int, validatorId *big.Int) (*types.Transaction, error) {
	return _OChainPortal.Contract.CancelWithdraw(&_OChainPortal.TransactOpts, nonce, validatorId)
}

// Deposit is a paid mutator transaction binding the contract method 0x47e7ef24.
//
// Solidity: function deposit(address receiver, uint256 amount) returns()
func (_OChainPortal *OChainPortalTransactor) Deposit(opts *bind.TransactOpts, receiver common.Address, amount *big.Int) (*types.Transaction, error) {
	return _OChainPortal.contract.Transact(opts, "deposit", receiver, amount)
}

// Deposit is a paid mutator transaction binding the contract method 0x47e7ef24.
//
// Solidity: function deposit(address receiver, uint256 amount) returns()
func (_OChainPortal *OChainPortalSession) Deposit(receiver common.Address, amount *big.Int) (*types.Transaction, error) {
	return _OChainPortal.Contract.Deposit(&_OChainPortal.TransactOpts, receiver, amount)
}

// Deposit is a paid mutator transaction binding the contract method 0x47e7ef24.
//
// Solidity: function deposit(address receiver, uint256 amount) returns()
func (_OChainPortal *OChainPortalTransactorSession) Deposit(receiver common.Address, amount *big.Int) (*types.Transaction, error) {
	return _OChainPortal.Contract.Deposit(&_OChainPortal.TransactOpts, receiver, amount)
}

// DepositUSD is a paid mutator transaction binding the contract method 0xdd15f993.
//
// Solidity: function depositUSD(address receiver, uint256 amount) returns()
func (_OChainPortal *OChainPortalTransactor) DepositUSD(opts *bind.TransactOpts, receiver common.Address, amount *big.Int) (*types.Transaction, error) {
	return _OChainPortal.contract.Transact(opts, "depositUSD", receiver, amount)
}

// DepositUSD is a paid mutator transaction binding the contract method 0xdd15f993.
//
// Solidity: function depositUSD(address receiver, uint256 amount) returns()
func (_OChainPortal *OChainPortalSession) DepositUSD(receiver common.Address, amount *big.Int) (*types.Transaction, error) {
	return _OChainPortal.Contract.DepositUSD(&_OChainPortal.TransactOpts, receiver, amount)
}

// DepositUSD is a paid mutator transaction binding the contract method 0xdd15f993.
//
// Solidity: function depositUSD(address receiver, uint256 amount) returns()
func (_OChainPortal *OChainPortalTransactorSession) DepositUSD(receiver common.Address, amount *big.Int) (*types.Transaction, error) {
	return _OChainPortal.Contract.DepositUSD(&_OChainPortal.TransactOpts, receiver, amount)
}

// DiamondCut is a paid mutator transaction binding the contract method 0x1f931c1c.
//
// Solidity: function diamondCut((address,uint8,bytes4[])[] _diamondCut, address _init, bytes _calldata) returns()
func (_OChainPortal *OChainPortalTransactor) DiamondCut(opts *bind.TransactOpts, _diamondCut []IDiamondCutFacetCut, _init common.Address, _calldata []byte) (*types.Transaction, error) {
	return _OChainPortal.contract.Transact(opts, "diamondCut", _diamondCut, _init, _calldata)
}

// DiamondCut is a paid mutator transaction binding the contract method 0x1f931c1c.
//
// Solidity: function diamondCut((address,uint8,bytes4[])[] _diamondCut, address _init, bytes _calldata) returns()
func (_OChainPortal *OChainPortalSession) DiamondCut(_diamondCut []IDiamondCutFacetCut, _init common.Address, _calldata []byte) (*types.Transaction, error) {
	return _OChainPortal.Contract.DiamondCut(&_OChainPortal.TransactOpts, _diamondCut, _init, _calldata)
}

// DiamondCut is a paid mutator transaction binding the contract method 0x1f931c1c.
//
// Solidity: function diamondCut((address,uint8,bytes4[])[] _diamondCut, address _init, bytes _calldata) returns()
func (_OChainPortal *OChainPortalTransactorSession) DiamondCut(_diamondCut []IDiamondCutFacetCut, _init common.Address, _calldata []byte) (*types.Transaction, error) {
	return _OChainPortal.Contract.DiamondCut(&_OChainPortal.TransactOpts, _diamondCut, _init, _calldata)
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

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address _newOwner) returns()
func (_OChainPortal *OChainPortalTransactor) TransferOwnership(opts *bind.TransactOpts, _newOwner common.Address) (*types.Transaction, error) {
	return _OChainPortal.contract.Transact(opts, "transferOwnership", _newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address _newOwner) returns()
func (_OChainPortal *OChainPortalSession) TransferOwnership(_newOwner common.Address) (*types.Transaction, error) {
	return _OChainPortal.Contract.TransferOwnership(&_OChainPortal.TransactOpts, _newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address _newOwner) returns()
func (_OChainPortal *OChainPortalTransactorSession) TransferOwnership(_newOwner common.Address) (*types.Transaction, error) {
	return _OChainPortal.Contract.TransferOwnership(&_OChainPortal.TransactOpts, _newOwner)
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

// Withdraw is a paid mutator transaction binding the contract method 0x2e1a7d4d.
//
// Solidity: function withdraw(uint256 nonce) returns(bool)
func (_OChainPortal *OChainPortalTransactor) Withdraw(opts *bind.TransactOpts, nonce *big.Int) (*types.Transaction, error) {
	return _OChainPortal.contract.Transact(opts, "withdraw", nonce)
}

// Withdraw is a paid mutator transaction binding the contract method 0x2e1a7d4d.
//
// Solidity: function withdraw(uint256 nonce) returns(bool)
func (_OChainPortal *OChainPortalSession) Withdraw(nonce *big.Int) (*types.Transaction, error) {
	return _OChainPortal.Contract.Withdraw(&_OChainPortal.TransactOpts, nonce)
}

// Withdraw is a paid mutator transaction binding the contract method 0x2e1a7d4d.
//
// Solidity: function withdraw(uint256 nonce) returns(bool)
func (_OChainPortal *OChainPortalTransactorSession) Withdraw(nonce *big.Int) (*types.Transaction, error) {
	return _OChainPortal.Contract.Withdraw(&_OChainPortal.TransactOpts, nonce)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_OChainPortal *OChainPortalTransactor) Fallback(opts *bind.TransactOpts, calldata []byte) (*types.Transaction, error) {
	return _OChainPortal.contract.RawTransact(opts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_OChainPortal *OChainPortalSession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _OChainPortal.Contract.Fallback(&_OChainPortal.TransactOpts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_OChainPortal *OChainPortalTransactorSession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _OChainPortal.Contract.Fallback(&_OChainPortal.TransactOpts, calldata)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_OChainPortal *OChainPortalTransactor) Receive(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OChainPortal.contract.RawTransact(opts, nil) // calldata is disallowed for receive function
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_OChainPortal *OChainPortalSession) Receive() (*types.Transaction, error) {
	return _OChainPortal.Contract.Receive(&_OChainPortal.TransactOpts)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_OChainPortal *OChainPortalTransactorSession) Receive() (*types.Transaction, error) {
	return _OChainPortal.Contract.Receive(&_OChainPortal.TransactOpts)
}

// OChainPortalDiamondCutIterator is returned from FilterDiamondCut and is used to iterate over the raw logs and unpacked data for DiamondCut events raised by the OChainPortal contract.
type OChainPortalDiamondCutIterator struct {
	Event *OChainPortalDiamondCut // Event containing the contract specifics and raw log

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
func (it *OChainPortalDiamondCutIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OChainPortalDiamondCut)
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
		it.Event = new(OChainPortalDiamondCut)
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
func (it *OChainPortalDiamondCutIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *OChainPortalDiamondCutIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// OChainPortalDiamondCut represents a DiamondCut event raised by the OChainPortal contract.
type OChainPortalDiamondCut struct {
	DiamondCut []IDiamondCutFacetCut
	Init       common.Address
	Calldata   []byte
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterDiamondCut is a free log retrieval operation binding the contract event 0x8faa70878671ccd212d20771b795c50af8fd3ff6cf27f4bde57e5d4de0aeb673.
//
// Solidity: event DiamondCut((address,uint8,bytes4[])[] _diamondCut, address _init, bytes _calldata)
func (_OChainPortal *OChainPortalFilterer) FilterDiamondCut(opts *bind.FilterOpts) (*OChainPortalDiamondCutIterator, error) {

	logs, sub, err := _OChainPortal.contract.FilterLogs(opts, "DiamondCut")
	if err != nil {
		return nil, err
	}
	return &OChainPortalDiamondCutIterator{contract: _OChainPortal.contract, event: "DiamondCut", logs: logs, sub: sub}, nil
}

// WatchDiamondCut is a free log subscription operation binding the contract event 0x8faa70878671ccd212d20771b795c50af8fd3ff6cf27f4bde57e5d4de0aeb673.
//
// Solidity: event DiamondCut((address,uint8,bytes4[])[] _diamondCut, address _init, bytes _calldata)
func (_OChainPortal *OChainPortalFilterer) WatchDiamondCut(opts *bind.WatchOpts, sink chan<- *OChainPortalDiamondCut) (event.Subscription, error) {

	logs, sub, err := _OChainPortal.contract.WatchLogs(opts, "DiamondCut")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(OChainPortalDiamondCut)
				if err := _OChainPortal.contract.UnpackLog(event, "DiamondCut", log); err != nil {
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

// ParseDiamondCut is a log parse operation binding the contract event 0x8faa70878671ccd212d20771b795c50af8fd3ff6cf27f4bde57e5d4de0aeb673.
//
// Solidity: event DiamondCut((address,uint8,bytes4[])[] _diamondCut, address _init, bytes _calldata)
func (_OChainPortal *OChainPortalFilterer) ParseDiamondCut(log types.Log) (*OChainPortalDiamondCut, error) {
	event := new(OChainPortalDiamondCut)
	if err := _OChainPortal.contract.UnpackLog(event, "DiamondCut", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
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

// OChainPortalOChainTokenDepositedIterator is returned from FilterOChainTokenDeposited and is used to iterate over the raw logs and unpacked data for OChainTokenDeposited events raised by the OChainPortal contract.
type OChainPortalOChainTokenDepositedIterator struct {
	Event *OChainPortalOChainTokenDeposited // Event containing the contract specifics and raw log

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
func (it *OChainPortalOChainTokenDepositedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OChainPortalOChainTokenDeposited)
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
		it.Event = new(OChainPortalOChainTokenDeposited)
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
func (it *OChainPortalOChainTokenDepositedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *OChainPortalOChainTokenDepositedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// OChainPortalOChainTokenDeposited represents a OChainTokenDeposited event raised by the OChainPortal contract.
type OChainPortalOChainTokenDeposited struct {
	Receiver common.Address
	Amount   *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterOChainTokenDeposited is a free log retrieval operation binding the contract event 0x88112089f3bf2b25173c874be9948d7624cf84b9f155afde2e3f065742a06917.
//
// Solidity: event OChainTokenDeposited(address indexed receiver, uint256 amount)
func (_OChainPortal *OChainPortalFilterer) FilterOChainTokenDeposited(opts *bind.FilterOpts, receiver []common.Address) (*OChainPortalOChainTokenDepositedIterator, error) {

	var receiverRule []interface{}
	for _, receiverItem := range receiver {
		receiverRule = append(receiverRule, receiverItem)
	}

	logs, sub, err := _OChainPortal.contract.FilterLogs(opts, "OChainTokenDeposited", receiverRule)
	if err != nil {
		return nil, err
	}
	return &OChainPortalOChainTokenDepositedIterator{contract: _OChainPortal.contract, event: "OChainTokenDeposited", logs: logs, sub: sub}, nil
}

// WatchOChainTokenDeposited is a free log subscription operation binding the contract event 0x88112089f3bf2b25173c874be9948d7624cf84b9f155afde2e3f065742a06917.
//
// Solidity: event OChainTokenDeposited(address indexed receiver, uint256 amount)
func (_OChainPortal *OChainPortalFilterer) WatchOChainTokenDeposited(opts *bind.WatchOpts, sink chan<- *OChainPortalOChainTokenDeposited, receiver []common.Address) (event.Subscription, error) {

	var receiverRule []interface{}
	for _, receiverItem := range receiver {
		receiverRule = append(receiverRule, receiverItem)
	}

	logs, sub, err := _OChainPortal.contract.WatchLogs(opts, "OChainTokenDeposited", receiverRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(OChainPortalOChainTokenDeposited)
				if err := _OChainPortal.contract.UnpackLog(event, "OChainTokenDeposited", log); err != nil {
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

// ParseOChainTokenDeposited is a log parse operation binding the contract event 0x88112089f3bf2b25173c874be9948d7624cf84b9f155afde2e3f065742a06917.
//
// Solidity: event OChainTokenDeposited(address indexed receiver, uint256 amount)
func (_OChainPortal *OChainPortalFilterer) ParseOChainTokenDeposited(log types.Log) (*OChainPortalOChainTokenDeposited, error) {
	event := new(OChainPortalOChainTokenDeposited)
	if err := _OChainPortal.contract.UnpackLog(event, "OChainTokenDeposited", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// OChainPortalOChainTokenWithdrawalCanceledIterator is returned from FilterOChainTokenWithdrawalCanceled and is used to iterate over the raw logs and unpacked data for OChainTokenWithdrawalCanceled events raised by the OChainPortal contract.
type OChainPortalOChainTokenWithdrawalCanceledIterator struct {
	Event *OChainPortalOChainTokenWithdrawalCanceled // Event containing the contract specifics and raw log

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
func (it *OChainPortalOChainTokenWithdrawalCanceledIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OChainPortalOChainTokenWithdrawalCanceled)
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
		it.Event = new(OChainPortalOChainTokenWithdrawalCanceled)
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
func (it *OChainPortalOChainTokenWithdrawalCanceledIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *OChainPortalOChainTokenWithdrawalCanceledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// OChainPortalOChainTokenWithdrawalCanceled represents a OChainTokenWithdrawalCanceled event raised by the OChainPortal contract.
type OChainPortalOChainTokenWithdrawalCanceled struct {
	Nonce  *big.Int
	Signer common.Address
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterOChainTokenWithdrawalCanceled is a free log retrieval operation binding the contract event 0x22ae3800792c66156c83bcc6edc71ab26c5c361fab0a834ebfecf95d258140a6.
//
// Solidity: event OChainTokenWithdrawalCanceled(uint256 indexed nonce, address signer)
func (_OChainPortal *OChainPortalFilterer) FilterOChainTokenWithdrawalCanceled(opts *bind.FilterOpts, nonce []*big.Int) (*OChainPortalOChainTokenWithdrawalCanceledIterator, error) {

	var nonceRule []interface{}
	for _, nonceItem := range nonce {
		nonceRule = append(nonceRule, nonceItem)
	}

	logs, sub, err := _OChainPortal.contract.FilterLogs(opts, "OChainTokenWithdrawalCanceled", nonceRule)
	if err != nil {
		return nil, err
	}
	return &OChainPortalOChainTokenWithdrawalCanceledIterator{contract: _OChainPortal.contract, event: "OChainTokenWithdrawalCanceled", logs: logs, sub: sub}, nil
}

// WatchOChainTokenWithdrawalCanceled is a free log subscription operation binding the contract event 0x22ae3800792c66156c83bcc6edc71ab26c5c361fab0a834ebfecf95d258140a6.
//
// Solidity: event OChainTokenWithdrawalCanceled(uint256 indexed nonce, address signer)
func (_OChainPortal *OChainPortalFilterer) WatchOChainTokenWithdrawalCanceled(opts *bind.WatchOpts, sink chan<- *OChainPortalOChainTokenWithdrawalCanceled, nonce []*big.Int) (event.Subscription, error) {

	var nonceRule []interface{}
	for _, nonceItem := range nonce {
		nonceRule = append(nonceRule, nonceItem)
	}

	logs, sub, err := _OChainPortal.contract.WatchLogs(opts, "OChainTokenWithdrawalCanceled", nonceRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(OChainPortalOChainTokenWithdrawalCanceled)
				if err := _OChainPortal.contract.UnpackLog(event, "OChainTokenWithdrawalCanceled", log); err != nil {
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

// ParseOChainTokenWithdrawalCanceled is a log parse operation binding the contract event 0x22ae3800792c66156c83bcc6edc71ab26c5c361fab0a834ebfecf95d258140a6.
//
// Solidity: event OChainTokenWithdrawalCanceled(uint256 indexed nonce, address signer)
func (_OChainPortal *OChainPortalFilterer) ParseOChainTokenWithdrawalCanceled(log types.Log) (*OChainPortalOChainTokenWithdrawalCanceled, error) {
	event := new(OChainPortalOChainTokenWithdrawalCanceled)
	if err := _OChainPortal.contract.UnpackLog(event, "OChainTokenWithdrawalCanceled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// OChainPortalOChainTokenWithdrawalExecutedIterator is returned from FilterOChainTokenWithdrawalExecuted and is used to iterate over the raw logs and unpacked data for OChainTokenWithdrawalExecuted events raised by the OChainPortal contract.
type OChainPortalOChainTokenWithdrawalExecutedIterator struct {
	Event *OChainPortalOChainTokenWithdrawalExecuted // Event containing the contract specifics and raw log

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
func (it *OChainPortalOChainTokenWithdrawalExecutedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OChainPortalOChainTokenWithdrawalExecuted)
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
		it.Event = new(OChainPortalOChainTokenWithdrawalExecuted)
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
func (it *OChainPortalOChainTokenWithdrawalExecutedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *OChainPortalOChainTokenWithdrawalExecutedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// OChainPortalOChainTokenWithdrawalExecuted represents a OChainTokenWithdrawalExecuted event raised by the OChainPortal contract.
type OChainPortalOChainTokenWithdrawalExecuted struct {
	Nonce *big.Int
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterOChainTokenWithdrawalExecuted is a free log retrieval operation binding the contract event 0xc5c1ff4ceee69008bffb1581fa3dbd1aca29688ada0b4d9ca21722291793564f.
//
// Solidity: event OChainTokenWithdrawalExecuted(uint256 indexed nonce)
func (_OChainPortal *OChainPortalFilterer) FilterOChainTokenWithdrawalExecuted(opts *bind.FilterOpts, nonce []*big.Int) (*OChainPortalOChainTokenWithdrawalExecutedIterator, error) {

	var nonceRule []interface{}
	for _, nonceItem := range nonce {
		nonceRule = append(nonceRule, nonceItem)
	}

	logs, sub, err := _OChainPortal.contract.FilterLogs(opts, "OChainTokenWithdrawalExecuted", nonceRule)
	if err != nil {
		return nil, err
	}
	return &OChainPortalOChainTokenWithdrawalExecutedIterator{contract: _OChainPortal.contract, event: "OChainTokenWithdrawalExecuted", logs: logs, sub: sub}, nil
}

// WatchOChainTokenWithdrawalExecuted is a free log subscription operation binding the contract event 0xc5c1ff4ceee69008bffb1581fa3dbd1aca29688ada0b4d9ca21722291793564f.
//
// Solidity: event OChainTokenWithdrawalExecuted(uint256 indexed nonce)
func (_OChainPortal *OChainPortalFilterer) WatchOChainTokenWithdrawalExecuted(opts *bind.WatchOpts, sink chan<- *OChainPortalOChainTokenWithdrawalExecuted, nonce []*big.Int) (event.Subscription, error) {

	var nonceRule []interface{}
	for _, nonceItem := range nonce {
		nonceRule = append(nonceRule, nonceItem)
	}

	logs, sub, err := _OChainPortal.contract.WatchLogs(opts, "OChainTokenWithdrawalExecuted", nonceRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(OChainPortalOChainTokenWithdrawalExecuted)
				if err := _OChainPortal.contract.UnpackLog(event, "OChainTokenWithdrawalExecuted", log); err != nil {
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

// ParseOChainTokenWithdrawalExecuted is a log parse operation binding the contract event 0xc5c1ff4ceee69008bffb1581fa3dbd1aca29688ada0b4d9ca21722291793564f.
//
// Solidity: event OChainTokenWithdrawalExecuted(uint256 indexed nonce)
func (_OChainPortal *OChainPortalFilterer) ParseOChainTokenWithdrawalExecuted(log types.Log) (*OChainPortalOChainTokenWithdrawalExecuted, error) {
	event := new(OChainPortalOChainTokenWithdrawalExecuted)
	if err := _OChainPortal.contract.UnpackLog(event, "OChainTokenWithdrawalExecuted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// OChainPortalOChainTokenWithdrawalSignedIterator is returned from FilterOChainTokenWithdrawalSigned and is used to iterate over the raw logs and unpacked data for OChainTokenWithdrawalSigned events raised by the OChainPortal contract.
type OChainPortalOChainTokenWithdrawalSignedIterator struct {
	Event *OChainPortalOChainTokenWithdrawalSigned // Event containing the contract specifics and raw log

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
func (it *OChainPortalOChainTokenWithdrawalSignedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OChainPortalOChainTokenWithdrawalSigned)
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
		it.Event = new(OChainPortalOChainTokenWithdrawalSigned)
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
func (it *OChainPortalOChainTokenWithdrawalSignedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *OChainPortalOChainTokenWithdrawalSignedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// OChainPortalOChainTokenWithdrawalSigned represents a OChainTokenWithdrawalSigned event raised by the OChainPortal contract.
type OChainPortalOChainTokenWithdrawalSigned struct {
	Receiver common.Address
	Nonce    *big.Int
	Signer   common.Address
	Amount   *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterOChainTokenWithdrawalSigned is a free log retrieval operation binding the contract event 0x705a95ffa8fe3052f7d1a9eaaa9096da294c91527cd40052fdcaa4f1c1c47eb9.
//
// Solidity: event OChainTokenWithdrawalSigned(address indexed receiver, uint256 indexed nonce, address signer, uint256 amount)
func (_OChainPortal *OChainPortalFilterer) FilterOChainTokenWithdrawalSigned(opts *bind.FilterOpts, receiver []common.Address, nonce []*big.Int) (*OChainPortalOChainTokenWithdrawalSignedIterator, error) {

	var receiverRule []interface{}
	for _, receiverItem := range receiver {
		receiverRule = append(receiverRule, receiverItem)
	}
	var nonceRule []interface{}
	for _, nonceItem := range nonce {
		nonceRule = append(nonceRule, nonceItem)
	}

	logs, sub, err := _OChainPortal.contract.FilterLogs(opts, "OChainTokenWithdrawalSigned", receiverRule, nonceRule)
	if err != nil {
		return nil, err
	}
	return &OChainPortalOChainTokenWithdrawalSignedIterator{contract: _OChainPortal.contract, event: "OChainTokenWithdrawalSigned", logs: logs, sub: sub}, nil
}

// WatchOChainTokenWithdrawalSigned is a free log subscription operation binding the contract event 0x705a95ffa8fe3052f7d1a9eaaa9096da294c91527cd40052fdcaa4f1c1c47eb9.
//
// Solidity: event OChainTokenWithdrawalSigned(address indexed receiver, uint256 indexed nonce, address signer, uint256 amount)
func (_OChainPortal *OChainPortalFilterer) WatchOChainTokenWithdrawalSigned(opts *bind.WatchOpts, sink chan<- *OChainPortalOChainTokenWithdrawalSigned, receiver []common.Address, nonce []*big.Int) (event.Subscription, error) {

	var receiverRule []interface{}
	for _, receiverItem := range receiver {
		receiverRule = append(receiverRule, receiverItem)
	}
	var nonceRule []interface{}
	for _, nonceItem := range nonce {
		nonceRule = append(nonceRule, nonceItem)
	}

	logs, sub, err := _OChainPortal.contract.WatchLogs(opts, "OChainTokenWithdrawalSigned", receiverRule, nonceRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(OChainPortalOChainTokenWithdrawalSigned)
				if err := _OChainPortal.contract.UnpackLog(event, "OChainTokenWithdrawalSigned", log); err != nil {
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

// ParseOChainTokenWithdrawalSigned is a log parse operation binding the contract event 0x705a95ffa8fe3052f7d1a9eaaa9096da294c91527cd40052fdcaa4f1c1c47eb9.
//
// Solidity: event OChainTokenWithdrawalSigned(address indexed receiver, uint256 indexed nonce, address signer, uint256 amount)
func (_OChainPortal *OChainPortalFilterer) ParseOChainTokenWithdrawalSigned(log types.Log) (*OChainPortalOChainTokenWithdrawalSigned, error) {
	event := new(OChainPortalOChainTokenWithdrawalSigned)
	if err := _OChainPortal.contract.UnpackLog(event, "OChainTokenWithdrawalSigned", log); err != nil {
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

// OChainPortalOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the OChainPortal contract.
type OChainPortalOwnershipTransferredIterator struct {
	Event *OChainPortalOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *OChainPortalOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OChainPortalOwnershipTransferred)
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
		it.Event = new(OChainPortalOwnershipTransferred)
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
func (it *OChainPortalOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *OChainPortalOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// OChainPortalOwnershipTransferred represents a OwnershipTransferred event raised by the OChainPortal contract.
type OChainPortalOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_OChainPortal *OChainPortalFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*OChainPortalOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _OChainPortal.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &OChainPortalOwnershipTransferredIterator{contract: _OChainPortal.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_OChainPortal *OChainPortalFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *OChainPortalOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _OChainPortal.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(OChainPortalOwnershipTransferred)
				if err := _OChainPortal.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

// ParseOwnershipTransferred is a log parse operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_OChainPortal *OChainPortalFilterer) ParseOwnershipTransferred(log types.Log) (*OChainPortalOwnershipTransferred, error) {
	event := new(OChainPortalOwnershipTransferred)
	if err := _OChainPortal.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// OChainPortalUSDDepositedIterator is returned from FilterUSDDeposited and is used to iterate over the raw logs and unpacked data for USDDeposited events raised by the OChainPortal contract.
type OChainPortalUSDDepositedIterator struct {
	Event *OChainPortalUSDDeposited // Event containing the contract specifics and raw log

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
func (it *OChainPortalUSDDepositedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OChainPortalUSDDeposited)
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
		it.Event = new(OChainPortalUSDDeposited)
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
func (it *OChainPortalUSDDepositedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *OChainPortalUSDDepositedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// OChainPortalUSDDeposited represents a USDDeposited event raised by the OChainPortal contract.
type OChainPortalUSDDeposited struct {
	Receiver common.Address
	Amount   *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterUSDDeposited is a free log retrieval operation binding the contract event 0xae6217b78ca94469877824ba7a838231a40d9fc97d55939c2ee70d44ac49cf50.
//
// Solidity: event USDDeposited(address indexed receiver, uint256 amount)
func (_OChainPortal *OChainPortalFilterer) FilterUSDDeposited(opts *bind.FilterOpts, receiver []common.Address) (*OChainPortalUSDDepositedIterator, error) {

	var receiverRule []interface{}
	for _, receiverItem := range receiver {
		receiverRule = append(receiverRule, receiverItem)
	}

	logs, sub, err := _OChainPortal.contract.FilterLogs(opts, "USDDeposited", receiverRule)
	if err != nil {
		return nil, err
	}
	return &OChainPortalUSDDepositedIterator{contract: _OChainPortal.contract, event: "USDDeposited", logs: logs, sub: sub}, nil
}

// WatchUSDDeposited is a free log subscription operation binding the contract event 0xae6217b78ca94469877824ba7a838231a40d9fc97d55939c2ee70d44ac49cf50.
//
// Solidity: event USDDeposited(address indexed receiver, uint256 amount)
func (_OChainPortal *OChainPortalFilterer) WatchUSDDeposited(opts *bind.WatchOpts, sink chan<- *OChainPortalUSDDeposited, receiver []common.Address) (event.Subscription, error) {

	var receiverRule []interface{}
	for _, receiverItem := range receiver {
		receiverRule = append(receiverRule, receiverItem)
	}

	logs, sub, err := _OChainPortal.contract.WatchLogs(opts, "USDDeposited", receiverRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(OChainPortalUSDDeposited)
				if err := _OChainPortal.contract.UnpackLog(event, "USDDeposited", log); err != nil {
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

// ParseUSDDeposited is a log parse operation binding the contract event 0xae6217b78ca94469877824ba7a838231a40d9fc97d55939c2ee70d44ac49cf50.
//
// Solidity: event USDDeposited(address indexed receiver, uint256 amount)
func (_OChainPortal *OChainPortalFilterer) ParseUSDDeposited(log types.Log) (*OChainPortalUSDDeposited, error) {
	event := new(OChainPortalUSDDeposited)
	if err := _OChainPortal.contract.UnpackLog(event, "USDDeposited", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
