// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package fclear

import (
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
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

// FclearABI is the input ABI used to generate the binding from.
const FclearABI = "[{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"_multi\",\"type\":\"uint32\"},{\"internalType\":\"string\",\"name\":\"_name\",\"type\":\"string\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint16\",\"name\":\"mem\",\"type\":\"uint16\"},{\"indexed\":false,\"internalType\":\"uint16\",\"name\":\"ric\",\"type\":\"uint16\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"isOff\",\"type\":\"bool\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"isBuy\",\"type\":\"bool\"}],\"name\":\"Clear\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"client\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"qty\",\"type\":\"uint32\"},{\"internalType\":\"uint64\",\"name\":\"price\",\"type\":\"uint64\"},{\"internalType\":\"uint16\",\"name\":\"symbol\",\"type\":\"uint16\"},{\"internalType\":\"uint16\",\"name\":\"member\",\"type\":\"uint16\"},{\"internalType\":\"bool\",\"name\":\"isOffset\",\"type\":\"bool\"},{\"internalType\":\"bool\",\"name\":\"isBuy\",\"type\":\"bool\"}],\"name\":\"dealClearing\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"client\",\"type\":\"uint32\"},{\"internalType\":\"uint16\",\"name\":\"symbol\",\"type\":\"uint16\"},{\"internalType\":\"uint16\",\"name\":\"member\",\"type\":\"uint16\"}],\"name\":\"getClientPosition\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"nLong\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"nShort\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getMulti\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"_multi\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"multi\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]"

// Fclear is an auto generated Go binding around an Ethereum contract.
type Fclear struct {
	FclearCaller     // Read-only binding to the contract
	FclearTransactor // Write-only binding to the contract
	FclearFilterer   // Log filterer for contract events
}

// FclearCaller is an auto generated read-only Go binding around an Ethereum contract.
type FclearCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// FclearTransactor is an auto generated write-only Go binding around an Ethereum contract.
type FclearTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// FclearFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type FclearFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// FclearSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type FclearSession struct {
	Contract     *Fclear           // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// FclearCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type FclearCallerSession struct {
	Contract *FclearCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// FclearTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type FclearTransactorSession struct {
	Contract     *FclearTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// FclearRaw is an auto generated low-level Go binding around an Ethereum contract.
type FclearRaw struct {
	Contract *Fclear // Generic contract binding to access the raw methods on
}

// FclearCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type FclearCallerRaw struct {
	Contract *FclearCaller // Generic read-only contract binding to access the raw methods on
}

// FclearTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type FclearTransactorRaw struct {
	Contract *FclearTransactor // Generic write-only contract binding to access the raw methods on
}

// NewFclear creates a new instance of Fclear, bound to a specific deployed contract.
func NewFclear(address common.Address, backend bind.ContractBackend) (*Fclear, error) {
	contract, err := bindFclear(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Fclear{FclearCaller: FclearCaller{contract: contract}, FclearTransactor: FclearTransactor{contract: contract}, FclearFilterer: FclearFilterer{contract: contract}}, nil
}

// NewFclearCaller creates a new read-only instance of Fclear, bound to a specific deployed contract.
func NewFclearCaller(address common.Address, caller bind.ContractCaller) (*FclearCaller, error) {
	contract, err := bindFclear(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &FclearCaller{contract: contract}, nil
}

// NewFclearTransactor creates a new write-only instance of Fclear, bound to a specific deployed contract.
func NewFclearTransactor(address common.Address, transactor bind.ContractTransactor) (*FclearTransactor, error) {
	contract, err := bindFclear(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &FclearTransactor{contract: contract}, nil
}

// NewFclearFilterer creates a new log filterer instance of Fclear, bound to a specific deployed contract.
func NewFclearFilterer(address common.Address, filterer bind.ContractFilterer) (*FclearFilterer, error) {
	contract, err := bindFclear(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &FclearFilterer{contract: contract}, nil
}

// bindFclear binds a generic wrapper to an already deployed contract.
func bindFclear(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(FclearABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Fclear *FclearRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Fclear.Contract.FclearCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Fclear *FclearRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Fclear.Contract.FclearTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Fclear *FclearRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Fclear.Contract.FclearTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Fclear *FclearCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Fclear.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Fclear *FclearTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Fclear.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Fclear *FclearTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Fclear.Contract.contract.Transact(opts, method, params...)
}

// DealClearing is a paid mutator transaction binding the contract method 0xbe704381.
//
// Solidity: function dealClearing(uint32 client, uint32 qty, uint64 price, uint16 symbol, uint16 member, bool isOffset, bool isBuy) returns()
func (_Fclear *FclearTransactor) DealClearing(opts *bind.TransactOpts, client uint32, qty uint32, price uint64, symbol uint16, member uint16, isOffset bool, isBuy bool) (*types.Transaction, error) {
	return _Fclear.contract.Transact(opts, "dealClearing", client, qty, price, symbol, member, isOffset, isBuy)
}

// DealClearing is a paid mutator transaction binding the contract method 0xbe704381.
//
// Solidity: function dealClearing(uint32 client, uint32 qty, uint64 price, uint16 symbol, uint16 member, bool isOffset, bool isBuy) returns()
func (_Fclear *FclearSession) DealClearing(client uint32, qty uint32, price uint64, symbol uint16, member uint16, isOffset bool, isBuy bool) (*types.Transaction, error) {
	return _Fclear.Contract.DealClearing(&_Fclear.TransactOpts, client, qty, price, symbol, member, isOffset, isBuy)
}

// DealClearing is a paid mutator transaction binding the contract method 0xbe704381.
//
// Solidity: function dealClearing(uint32 client, uint32 qty, uint64 price, uint16 symbol, uint16 member, bool isOffset, bool isBuy) returns()
func (_Fclear *FclearTransactorSession) DealClearing(client uint32, qty uint32, price uint64, symbol uint16, member uint16, isOffset bool, isBuy bool) (*types.Transaction, error) {
	return _Fclear.Contract.DealClearing(&_Fclear.TransactOpts, client, qty, price, symbol, member, isOffset, isBuy)
}

// GetClientPosition is a paid mutator transaction binding the contract method 0xf42a90d6.
//
// Solidity: function getClientPosition(uint32 client, uint16 symbol, uint16 member) returns(uint32 nLong, uint32 nShort)
func (_Fclear *FclearTransactor) GetClientPosition(opts *bind.TransactOpts, client uint32, symbol uint16, member uint16) (*types.Transaction, error) {
	return _Fclear.contract.Transact(opts, "getClientPosition", client, symbol, member)
}

// GetClientPosition is a paid mutator transaction binding the contract method 0xf42a90d6.
//
// Solidity: function getClientPosition(uint32 client, uint16 symbol, uint16 member) returns(uint32 nLong, uint32 nShort)
func (_Fclear *FclearSession) GetClientPosition(client uint32, symbol uint16, member uint16) (*types.Transaction, error) {
	return _Fclear.Contract.GetClientPosition(&_Fclear.TransactOpts, client, symbol, member)
}

// GetClientPosition is a paid mutator transaction binding the contract method 0xf42a90d6.
//
// Solidity: function getClientPosition(uint32 client, uint16 symbol, uint16 member) returns(uint32 nLong, uint32 nShort)
func (_Fclear *FclearTransactorSession) GetClientPosition(client uint32, symbol uint16, member uint16) (*types.Transaction, error) {
	return _Fclear.Contract.GetClientPosition(&_Fclear.TransactOpts, client, symbol, member)
}

// GetMulti is a paid mutator transaction binding the contract method 0xf3c66284.
//
// Solidity: function getMulti() returns(uint32 _multi)
func (_Fclear *FclearTransactor) GetMulti(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Fclear.contract.Transact(opts, "getMulti")
}

// GetMulti is a paid mutator transaction binding the contract method 0xf3c66284.
//
// Solidity: function getMulti() returns(uint32 _multi)
func (_Fclear *FclearSession) GetMulti() (*types.Transaction, error) {
	return _Fclear.Contract.GetMulti(&_Fclear.TransactOpts)
}

// GetMulti is a paid mutator transaction binding the contract method 0xf3c66284.
//
// Solidity: function getMulti() returns(uint32 _multi)
func (_Fclear *FclearTransactorSession) GetMulti() (*types.Transaction, error) {
	return _Fclear.Contract.GetMulti(&_Fclear.TransactOpts)
}

// Multi is a paid mutator transaction binding the contract method 0x1b8f5d50.
//
// Solidity: function multi() returns(uint32)
func (_Fclear *FclearTransactor) Multi(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Fclear.contract.Transact(opts, "multi")
}

// Multi is a paid mutator transaction binding the contract method 0x1b8f5d50.
//
// Solidity: function multi() returns(uint32)
func (_Fclear *FclearSession) Multi() (*types.Transaction, error) {
	return _Fclear.Contract.Multi(&_Fclear.TransactOpts)
}

// Multi is a paid mutator transaction binding the contract method 0x1b8f5d50.
//
// Solidity: function multi() returns(uint32)
func (_Fclear *FclearTransactorSession) Multi() (*types.Transaction, error) {
	return _Fclear.Contract.Multi(&_Fclear.TransactOpts)
}

// Name is a paid mutator transaction binding the contract method 0x06fdde03.
//
// Solidity: function name() returns(string)
func (_Fclear *FclearTransactor) Name(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Fclear.contract.Transact(opts, "name")
}

// Name is a paid mutator transaction binding the contract method 0x06fdde03.
//
// Solidity: function name() returns(string)
func (_Fclear *FclearSession) Name() (*types.Transaction, error) {
	return _Fclear.Contract.Name(&_Fclear.TransactOpts)
}

// Name is a paid mutator transaction binding the contract method 0x06fdde03.
//
// Solidity: function name() returns(string)
func (_Fclear *FclearTransactorSession) Name() (*types.Transaction, error) {
	return _Fclear.Contract.Name(&_Fclear.TransactOpts)
}

// Owner is a paid mutator transaction binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() returns(address)
func (_Fclear *FclearTransactor) Owner(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Fclear.contract.Transact(opts, "owner")
}

// Owner is a paid mutator transaction binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() returns(address)
func (_Fclear *FclearSession) Owner() (*types.Transaction, error) {
	return _Fclear.Contract.Owner(&_Fclear.TransactOpts)
}

// Owner is a paid mutator transaction binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() returns(address)
func (_Fclear *FclearTransactorSession) Owner() (*types.Transaction, error) {
	return _Fclear.Contract.Owner(&_Fclear.TransactOpts)
}

// FclearClearIterator is returned from FilterClear and is used to iterate over the raw logs and unpacked data for Clear events raised by the Fclear contract.
type FclearClearIterator struct {
	Event *FclearClear // Event containing the contract specifics and raw log

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
func (it *FclearClearIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FclearClear)
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
		it.Event = new(FclearClear)
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
func (it *FclearClearIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FclearClearIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FclearClear represents a Clear event raised by the Fclear contract.
type FclearClear struct {
	Mem   uint16
	Ric   uint16
	IsOff bool
	IsBuy bool
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterClear is a free log retrieval operation binding the contract event 0x3cf549b6d850b9f69dda701acbacf614fd4abef17a0ac1bf6b2ab4293b05640c.
//
// Solidity: event Clear(uint16 mem, uint16 ric, bool isOff, bool isBuy)
func (_Fclear *FclearFilterer) FilterClear(opts *bind.FilterOpts) (*FclearClearIterator, error) {

	logs, sub, err := _Fclear.contract.FilterLogs(opts, "Clear")
	if err != nil {
		return nil, err
	}
	return &FclearClearIterator{contract: _Fclear.contract, event: "Clear", logs: logs, sub: sub}, nil
}

// WatchClear is a free log subscription operation binding the contract event 0x3cf549b6d850b9f69dda701acbacf614fd4abef17a0ac1bf6b2ab4293b05640c.
//
// Solidity: event Clear(uint16 mem, uint16 ric, bool isOff, bool isBuy)
func (_Fclear *FclearFilterer) WatchClear(opts *bind.WatchOpts, sink chan<- *FclearClear) (event.Subscription, error) {

	logs, sub, err := _Fclear.contract.WatchLogs(opts, "Clear")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FclearClear)
				if err := _Fclear.contract.UnpackLog(event, "Clear", log); err != nil {
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

// ParseClear is a log parse operation binding the contract event 0x3cf549b6d850b9f69dda701acbacf614fd4abef17a0ac1bf6b2ab4293b05640c.
//
// Solidity: event Clear(uint16 mem, uint16 ric, bool isOff, bool isBuy)
func (_Fclear *FclearFilterer) ParseClear(log types.Log) (*FclearClear, error) {
	event := new(FclearClear)
	if err := _Fclear.contract.UnpackLog(event, "Clear", log); err != nil {
		return nil, err
	}
	return event, nil
}
