// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package IERCHandler

import (
	"math/big"
	"strings"

	ethereum "github.com/hacpy/go-ethereum"
	"github.com/hacpy/go-ethereum/accounts/abi"
	"github.com/hacpy/go-ethereum/accounts/abi/bind"
	"github.com/hacpy/go-ethereum/common"
	"github.com/hacpy/go-ethereum/core/types"
	"github.com/hacpy/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

// IERCHandlerABI is the input ABI used to generate the binding from.
const IERCHandlerABI = "[{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"resourceID\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"contractAddress\",\"type\":\"address\"}],\"name\":\"setResource\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"contractAddress\",\"type\":\"address\"}],\"name\":\"setBurnable\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amountOrTokenID\",\"type\":\"uint256\"}],\"name\":\"withdraw\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

// IERCHandler is an auto generated Go binding around an Ethereum contract.
type IERCHandler struct {
	IERCHandlerCaller     // Read-only binding to the contract
	IERCHandlerTransactor // Write-only binding to the contract
	IERCHandlerFilterer   // Log filterer for contract events
}

// IERCHandlerCaller is an auto generated read-only Go binding around an Ethereum contract.
type IERCHandlerCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IERCHandlerTransactor is an auto generated write-only Go binding around an Ethereum contract.
type IERCHandlerTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IERCHandlerFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type IERCHandlerFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IERCHandlerSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type IERCHandlerSession struct {
	Contract     *IERCHandler      // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// IERCHandlerCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type IERCHandlerCallerSession struct {
	Contract *IERCHandlerCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts      // Call options to use throughout this session
}

// IERCHandlerTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type IERCHandlerTransactorSession struct {
	Contract     *IERCHandlerTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts      // Transaction auth options to use throughout this session
}

// IERCHandlerRaw is an auto generated low-level Go binding around an Ethereum contract.
type IERCHandlerRaw struct {
	Contract *IERCHandler // Generic contract binding to access the raw methods on
}

// IERCHandlerCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type IERCHandlerCallerRaw struct {
	Contract *IERCHandlerCaller // Generic read-only contract binding to access the raw methods on
}

// IERCHandlerTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type IERCHandlerTransactorRaw struct {
	Contract *IERCHandlerTransactor // Generic write-only contract binding to access the raw methods on
}

// NewIERCHandler creates a new instance of IERCHandler, bound to a specific deployed contract.
func NewIERCHandler(address common.Address, backend bind.ContractBackend) (*IERCHandler, error) {
	contract, err := bindIERCHandler(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &IERCHandler{IERCHandlerCaller: IERCHandlerCaller{contract: contract}, IERCHandlerTransactor: IERCHandlerTransactor{contract: contract}, IERCHandlerFilterer: IERCHandlerFilterer{contract: contract}}, nil
}

// NewIERCHandlerCaller creates a new read-only instance of IERCHandler, bound to a specific deployed contract.
func NewIERCHandlerCaller(address common.Address, caller bind.ContractCaller) (*IERCHandlerCaller, error) {
	contract, err := bindIERCHandler(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &IERCHandlerCaller{contract: contract}, nil
}

// NewIERCHandlerTransactor creates a new write-only instance of IERCHandler, bound to a specific deployed contract.
func NewIERCHandlerTransactor(address common.Address, transactor bind.ContractTransactor) (*IERCHandlerTransactor, error) {
	contract, err := bindIERCHandler(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &IERCHandlerTransactor{contract: contract}, nil
}

// NewIERCHandlerFilterer creates a new log filterer instance of IERCHandler, bound to a specific deployed contract.
func NewIERCHandlerFilterer(address common.Address, filterer bind.ContractFilterer) (*IERCHandlerFilterer, error) {
	contract, err := bindIERCHandler(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &IERCHandlerFilterer{contract: contract}, nil
}

// bindIERCHandler binds a generic wrapper to an already deployed contract.
func bindIERCHandler(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(IERCHandlerABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IERCHandler *IERCHandlerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IERCHandler.Contract.IERCHandlerCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IERCHandler *IERCHandlerRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IERCHandler.Contract.IERCHandlerTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IERCHandler *IERCHandlerRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IERCHandler.Contract.IERCHandlerTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IERCHandler *IERCHandlerCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IERCHandler.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IERCHandler *IERCHandlerTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IERCHandler.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IERCHandler *IERCHandlerTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IERCHandler.Contract.contract.Transact(opts, method, params...)
}

// SetBurnable is a paid mutator transaction binding the contract method 0x07b7ed99.
//
// Solidity: function setBurnable(address contractAddress) returns()
func (_IERCHandler *IERCHandlerTransactor) SetBurnable(opts *bind.TransactOpts, contractAddress common.Address) (*types.Transaction, error) {
	return _IERCHandler.contract.Transact(opts, "setBurnable", contractAddress)
}

// SetBurnable is a paid mutator transaction binding the contract method 0x07b7ed99.
//
// Solidity: function setBurnable(address contractAddress) returns()
func (_IERCHandler *IERCHandlerSession) SetBurnable(contractAddress common.Address) (*types.Transaction, error) {
	return _IERCHandler.Contract.SetBurnable(&_IERCHandler.TransactOpts, contractAddress)
}

// SetBurnable is a paid mutator transaction binding the contract method 0x07b7ed99.
//
// Solidity: function setBurnable(address contractAddress) returns()
func (_IERCHandler *IERCHandlerTransactorSession) SetBurnable(contractAddress common.Address) (*types.Transaction, error) {
	return _IERCHandler.Contract.SetBurnable(&_IERCHandler.TransactOpts, contractAddress)
}

// SetResource is a paid mutator transaction binding the contract method 0xb8fa3736.
//
// Solidity: function setResource(bytes32 resourceID, address contractAddress) returns()
func (_IERCHandler *IERCHandlerTransactor) SetResource(opts *bind.TransactOpts, resourceID [32]byte, contractAddress common.Address) (*types.Transaction, error) {
	return _IERCHandler.contract.Transact(opts, "setResource", resourceID, contractAddress)
}

// SetResource is a paid mutator transaction binding the contract method 0xb8fa3736.
//
// Solidity: function setResource(bytes32 resourceID, address contractAddress) returns()
func (_IERCHandler *IERCHandlerSession) SetResource(resourceID [32]byte, contractAddress common.Address) (*types.Transaction, error) {
	return _IERCHandler.Contract.SetResource(&_IERCHandler.TransactOpts, resourceID, contractAddress)
}

// SetResource is a paid mutator transaction binding the contract method 0xb8fa3736.
//
// Solidity: function setResource(bytes32 resourceID, address contractAddress) returns()
func (_IERCHandler *IERCHandlerTransactorSession) SetResource(resourceID [32]byte, contractAddress common.Address) (*types.Transaction, error) {
	return _IERCHandler.Contract.SetResource(&_IERCHandler.TransactOpts, resourceID, contractAddress)
}

// Withdraw is a paid mutator transaction binding the contract method 0xd9caed12.
//
// Solidity: function withdraw(address tokenAddress, address recipient, uint256 amountOrTokenID) returns()
func (_IERCHandler *IERCHandlerTransactor) Withdraw(opts *bind.TransactOpts, tokenAddress common.Address, recipient common.Address, amountOrTokenID *big.Int) (*types.Transaction, error) {
	return _IERCHandler.contract.Transact(opts, "withdraw", tokenAddress, recipient, amountOrTokenID)
}

// Withdraw is a paid mutator transaction binding the contract method 0xd9caed12.
//
// Solidity: function withdraw(address tokenAddress, address recipient, uint256 amountOrTokenID) returns()
func (_IERCHandler *IERCHandlerSession) Withdraw(tokenAddress common.Address, recipient common.Address, amountOrTokenID *big.Int) (*types.Transaction, error) {
	return _IERCHandler.Contract.Withdraw(&_IERCHandler.TransactOpts, tokenAddress, recipient, amountOrTokenID)
}

// Withdraw is a paid mutator transaction binding the contract method 0xd9caed12.
//
// Solidity: function withdraw(address tokenAddress, address recipient, uint256 amountOrTokenID) returns()
func (_IERCHandler *IERCHandlerTransactorSession) Withdraw(tokenAddress common.Address, recipient common.Address, amountOrTokenID *big.Int) (*types.Transaction, error) {
	return _IERCHandler.Contract.Withdraw(&_IERCHandler.TransactOpts, tokenAddress, recipient, amountOrTokenID)
}
