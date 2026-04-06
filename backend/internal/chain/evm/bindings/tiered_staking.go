// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package bindings

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

// TieredStakingPosition is an auto generated low-level Go binding around an user-defined struct.
type TieredStakingPosition struct {
	Owner         common.Address
	Amount        *big.Int
	Tier          uint8
	StakedAt      *big.Int
	LockUntil     *big.Int
	LastClaimedAt *big.Int
	Active        bool
}

// TieredStakingMetaData contains all meta data concerning the TieredStaking contract.
var TieredStakingMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"_stakingToken\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_treasury\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"BPS_DENOMINATOR\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"PENALTY_BPS\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"SECONDS_PER_YEAR\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"acceptOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"claimRewards\",\"inputs\":[{\"name\":\"positionId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"getAPR\",\"inputs\":[{\"name\":\"tier\",\"type\":\"uint8\",\"internalType\":\"enumTieredStaking.Tier\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getPosition\",\"inputs\":[{\"name\":\"positionId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structTieredStaking.Position\",\"components\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"tier\",\"type\":\"uint8\",\"internalType\":\"enumTieredStaking.Tier\"},{\"name\":\"stakedAt\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"lockUntil\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"lastClaimedAt\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"active\",\"type\":\"bool\",\"internalType\":\"bool\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"nextPositionId\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"pause\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"paused\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"pendingOwner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"renounceOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"stake\",\"inputs\":[{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"tier\",\"type\":\"uint8\",\"internalType\":\"enumTieredStaking.Tier\"}],\"outputs\":[{\"name\":\"positionId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"stakingToken\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIERC20\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"treasury\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"unpause\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"unstake\",\"inputs\":[{\"name\":\"positionId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"updateTierAPR\",\"inputs\":[{\"name\":\"tier\",\"type\":\"uint8\",\"internalType\":\"enumTieredStaking.Tier\"},{\"name\":\"newAprBps\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"updateTreasury\",\"inputs\":[{\"name\":\"newTreasury\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"OwnershipTransferStarted\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferred\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Paused\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RewardsClaimed\",\"inputs\":[{\"name\":\"user\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"positionId\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Staked\",\"inputs\":[{\"name\":\"user\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"tier\",\"type\":\"uint8\",\"indexed\":false,\"internalType\":\"enumTieredStaking.Tier\"},{\"name\":\"positionId\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"TierUpdated\",\"inputs\":[{\"name\":\"tier\",\"type\":\"uint8\",\"indexed\":false,\"internalType\":\"enumTieredStaking.Tier\"},{\"name\":\"newAprBps\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"TreasuryUpdated\",\"inputs\":[{\"name\":\"newTreasury\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Unpaused\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Unstaked\",\"inputs\":[{\"name\":\"user\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"rewards\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"penalty\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"EnforcedPause\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ExpectedPause\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidAmount\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidTier\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotPositionOwner\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OwnableInvalidOwner\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"OwnableUnauthorizedAccount\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"PositionNotActive\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ReentrancyGuardReentrantCall\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"SafeERC20FailedOperation\",\"inputs\":[{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"}]}]",
}

// TieredStakingABI is the input ABI used to generate the binding from.
// Deprecated: Use TieredStakingMetaData.ABI instead.
var TieredStakingABI = TieredStakingMetaData.ABI

// TieredStaking is an auto generated Go binding around an Ethereum contract.
type TieredStaking struct {
	TieredStakingCaller     // Read-only binding to the contract
	TieredStakingTransactor // Write-only binding to the contract
	TieredStakingFilterer   // Log filterer for contract events
}

// TieredStakingCaller is an auto generated read-only Go binding around an Ethereum contract.
type TieredStakingCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TieredStakingTransactor is an auto generated write-only Go binding around an Ethereum contract.
type TieredStakingTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TieredStakingFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type TieredStakingFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TieredStakingSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type TieredStakingSession struct {
	Contract     *TieredStaking    // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// TieredStakingCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type TieredStakingCallerSession struct {
	Contract *TieredStakingCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts        // Call options to use throughout this session
}

// TieredStakingTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type TieredStakingTransactorSession struct {
	Contract     *TieredStakingTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts        // Transaction auth options to use throughout this session
}

// TieredStakingRaw is an auto generated low-level Go binding around an Ethereum contract.
type TieredStakingRaw struct {
	Contract *TieredStaking // Generic contract binding to access the raw methods on
}

// TieredStakingCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type TieredStakingCallerRaw struct {
	Contract *TieredStakingCaller // Generic read-only contract binding to access the raw methods on
}

// TieredStakingTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type TieredStakingTransactorRaw struct {
	Contract *TieredStakingTransactor // Generic write-only contract binding to access the raw methods on
}

// NewTieredStaking creates a new instance of TieredStaking, bound to a specific deployed contract.
func NewTieredStaking(address common.Address, backend bind.ContractBackend) (*TieredStaking, error) {
	contract, err := bindTieredStaking(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &TieredStaking{TieredStakingCaller: TieredStakingCaller{contract: contract}, TieredStakingTransactor: TieredStakingTransactor{contract: contract}, TieredStakingFilterer: TieredStakingFilterer{contract: contract}}, nil
}

// NewTieredStakingCaller creates a new read-only instance of TieredStaking, bound to a specific deployed contract.
func NewTieredStakingCaller(address common.Address, caller bind.ContractCaller) (*TieredStakingCaller, error) {
	contract, err := bindTieredStaking(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &TieredStakingCaller{contract: contract}, nil
}

// NewTieredStakingTransactor creates a new write-only instance of TieredStaking, bound to a specific deployed contract.
func NewTieredStakingTransactor(address common.Address, transactor bind.ContractTransactor) (*TieredStakingTransactor, error) {
	contract, err := bindTieredStaking(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &TieredStakingTransactor{contract: contract}, nil
}

// NewTieredStakingFilterer creates a new log filterer instance of TieredStaking, bound to a specific deployed contract.
func NewTieredStakingFilterer(address common.Address, filterer bind.ContractFilterer) (*TieredStakingFilterer, error) {
	contract, err := bindTieredStaking(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &TieredStakingFilterer{contract: contract}, nil
}

// bindTieredStaking binds a generic wrapper to an already deployed contract.
func bindTieredStaking(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := TieredStakingMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TieredStaking *TieredStakingRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _TieredStaking.Contract.TieredStakingCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TieredStaking *TieredStakingRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TieredStaking.Contract.TieredStakingTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TieredStaking *TieredStakingRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TieredStaking.Contract.TieredStakingTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TieredStaking *TieredStakingCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _TieredStaking.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TieredStaking *TieredStakingTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TieredStaking.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TieredStaking *TieredStakingTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TieredStaking.Contract.contract.Transact(opts, method, params...)
}

// BPSDENOMINATOR is a free data retrieval call binding the contract method 0xe1a45218.
//
// Solidity: function BPS_DENOMINATOR() view returns(uint256)
func (_TieredStaking *TieredStakingCaller) BPSDENOMINATOR(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _TieredStaking.contract.Call(opts, &out, "BPS_DENOMINATOR")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BPSDENOMINATOR is a free data retrieval call binding the contract method 0xe1a45218.
//
// Solidity: function BPS_DENOMINATOR() view returns(uint256)
func (_TieredStaking *TieredStakingSession) BPSDENOMINATOR() (*big.Int, error) {
	return _TieredStaking.Contract.BPSDENOMINATOR(&_TieredStaking.CallOpts)
}

// BPSDENOMINATOR is a free data retrieval call binding the contract method 0xe1a45218.
//
// Solidity: function BPS_DENOMINATOR() view returns(uint256)
func (_TieredStaking *TieredStakingCallerSession) BPSDENOMINATOR() (*big.Int, error) {
	return _TieredStaking.Contract.BPSDENOMINATOR(&_TieredStaking.CallOpts)
}

// PENALTYBPS is a free data retrieval call binding the contract method 0x1efe5321.
//
// Solidity: function PENALTY_BPS() view returns(uint256)
func (_TieredStaking *TieredStakingCaller) PENALTYBPS(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _TieredStaking.contract.Call(opts, &out, "PENALTY_BPS")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// PENALTYBPS is a free data retrieval call binding the contract method 0x1efe5321.
//
// Solidity: function PENALTY_BPS() view returns(uint256)
func (_TieredStaking *TieredStakingSession) PENALTYBPS() (*big.Int, error) {
	return _TieredStaking.Contract.PENALTYBPS(&_TieredStaking.CallOpts)
}

// PENALTYBPS is a free data retrieval call binding the contract method 0x1efe5321.
//
// Solidity: function PENALTY_BPS() view returns(uint256)
func (_TieredStaking *TieredStakingCallerSession) PENALTYBPS() (*big.Int, error) {
	return _TieredStaking.Contract.PENALTYBPS(&_TieredStaking.CallOpts)
}

// SECONDSPERYEAR is a free data retrieval call binding the contract method 0xe6a69ab8.
//
// Solidity: function SECONDS_PER_YEAR() view returns(uint256)
func (_TieredStaking *TieredStakingCaller) SECONDSPERYEAR(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _TieredStaking.contract.Call(opts, &out, "SECONDS_PER_YEAR")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// SECONDSPERYEAR is a free data retrieval call binding the contract method 0xe6a69ab8.
//
// Solidity: function SECONDS_PER_YEAR() view returns(uint256)
func (_TieredStaking *TieredStakingSession) SECONDSPERYEAR() (*big.Int, error) {
	return _TieredStaking.Contract.SECONDSPERYEAR(&_TieredStaking.CallOpts)
}

// SECONDSPERYEAR is a free data retrieval call binding the contract method 0xe6a69ab8.
//
// Solidity: function SECONDS_PER_YEAR() view returns(uint256)
func (_TieredStaking *TieredStakingCallerSession) SECONDSPERYEAR() (*big.Int, error) {
	return _TieredStaking.Contract.SECONDSPERYEAR(&_TieredStaking.CallOpts)
}

// GetAPR is a free data retrieval call binding the contract method 0x4cc723cc.
//
// Solidity: function getAPR(uint8 tier) view returns(uint256)
func (_TieredStaking *TieredStakingCaller) GetAPR(opts *bind.CallOpts, tier uint8) (*big.Int, error) {
	var out []interface{}
	err := _TieredStaking.contract.Call(opts, &out, "getAPR", tier)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetAPR is a free data retrieval call binding the contract method 0x4cc723cc.
//
// Solidity: function getAPR(uint8 tier) view returns(uint256)
func (_TieredStaking *TieredStakingSession) GetAPR(tier uint8) (*big.Int, error) {
	return _TieredStaking.Contract.GetAPR(&_TieredStaking.CallOpts, tier)
}

// GetAPR is a free data retrieval call binding the contract method 0x4cc723cc.
//
// Solidity: function getAPR(uint8 tier) view returns(uint256)
func (_TieredStaking *TieredStakingCallerSession) GetAPR(tier uint8) (*big.Int, error) {
	return _TieredStaking.Contract.GetAPR(&_TieredStaking.CallOpts, tier)
}

// GetPosition is a free data retrieval call binding the contract method 0xeb02c301.
//
// Solidity: function getPosition(uint256 positionId) view returns((address,uint256,uint8,uint256,uint256,uint256,bool))
func (_TieredStaking *TieredStakingCaller) GetPosition(opts *bind.CallOpts, positionId *big.Int) (TieredStakingPosition, error) {
	var out []interface{}
	err := _TieredStaking.contract.Call(opts, &out, "getPosition", positionId)

	if err != nil {
		return *new(TieredStakingPosition), err
	}

	out0 := *abi.ConvertType(out[0], new(TieredStakingPosition)).(*TieredStakingPosition)

	return out0, err

}

// GetPosition is a free data retrieval call binding the contract method 0xeb02c301.
//
// Solidity: function getPosition(uint256 positionId) view returns((address,uint256,uint8,uint256,uint256,uint256,bool))
func (_TieredStaking *TieredStakingSession) GetPosition(positionId *big.Int) (TieredStakingPosition, error) {
	return _TieredStaking.Contract.GetPosition(&_TieredStaking.CallOpts, positionId)
}

// GetPosition is a free data retrieval call binding the contract method 0xeb02c301.
//
// Solidity: function getPosition(uint256 positionId) view returns((address,uint256,uint8,uint256,uint256,uint256,bool))
func (_TieredStaking *TieredStakingCallerSession) GetPosition(positionId *big.Int) (TieredStakingPosition, error) {
	return _TieredStaking.Contract.GetPosition(&_TieredStaking.CallOpts, positionId)
}

// NextPositionId is a free data retrieval call binding the contract method 0x899346c7.
//
// Solidity: function nextPositionId() view returns(uint256)
func (_TieredStaking *TieredStakingCaller) NextPositionId(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _TieredStaking.contract.Call(opts, &out, "nextPositionId")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// NextPositionId is a free data retrieval call binding the contract method 0x899346c7.
//
// Solidity: function nextPositionId() view returns(uint256)
func (_TieredStaking *TieredStakingSession) NextPositionId() (*big.Int, error) {
	return _TieredStaking.Contract.NextPositionId(&_TieredStaking.CallOpts)
}

// NextPositionId is a free data retrieval call binding the contract method 0x899346c7.
//
// Solidity: function nextPositionId() view returns(uint256)
func (_TieredStaking *TieredStakingCallerSession) NextPositionId() (*big.Int, error) {
	return _TieredStaking.Contract.NextPositionId(&_TieredStaking.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_TieredStaking *TieredStakingCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _TieredStaking.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_TieredStaking *TieredStakingSession) Owner() (common.Address, error) {
	return _TieredStaking.Contract.Owner(&_TieredStaking.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_TieredStaking *TieredStakingCallerSession) Owner() (common.Address, error) {
	return _TieredStaking.Contract.Owner(&_TieredStaking.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_TieredStaking *TieredStakingCaller) Paused(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _TieredStaking.contract.Call(opts, &out, "paused")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_TieredStaking *TieredStakingSession) Paused() (bool, error) {
	return _TieredStaking.Contract.Paused(&_TieredStaking.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_TieredStaking *TieredStakingCallerSession) Paused() (bool, error) {
	return _TieredStaking.Contract.Paused(&_TieredStaking.CallOpts)
}

// PendingOwner is a free data retrieval call binding the contract method 0xe30c3978.
//
// Solidity: function pendingOwner() view returns(address)
func (_TieredStaking *TieredStakingCaller) PendingOwner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _TieredStaking.contract.Call(opts, &out, "pendingOwner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// PendingOwner is a free data retrieval call binding the contract method 0xe30c3978.
//
// Solidity: function pendingOwner() view returns(address)
func (_TieredStaking *TieredStakingSession) PendingOwner() (common.Address, error) {
	return _TieredStaking.Contract.PendingOwner(&_TieredStaking.CallOpts)
}

// PendingOwner is a free data retrieval call binding the contract method 0xe30c3978.
//
// Solidity: function pendingOwner() view returns(address)
func (_TieredStaking *TieredStakingCallerSession) PendingOwner() (common.Address, error) {
	return _TieredStaking.Contract.PendingOwner(&_TieredStaking.CallOpts)
}

// StakingToken is a free data retrieval call binding the contract method 0x72f702f3.
//
// Solidity: function stakingToken() view returns(address)
func (_TieredStaking *TieredStakingCaller) StakingToken(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _TieredStaking.contract.Call(opts, &out, "stakingToken")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// StakingToken is a free data retrieval call binding the contract method 0x72f702f3.
//
// Solidity: function stakingToken() view returns(address)
func (_TieredStaking *TieredStakingSession) StakingToken() (common.Address, error) {
	return _TieredStaking.Contract.StakingToken(&_TieredStaking.CallOpts)
}

// StakingToken is a free data retrieval call binding the contract method 0x72f702f3.
//
// Solidity: function stakingToken() view returns(address)
func (_TieredStaking *TieredStakingCallerSession) StakingToken() (common.Address, error) {
	return _TieredStaking.Contract.StakingToken(&_TieredStaking.CallOpts)
}

// Treasury is a free data retrieval call binding the contract method 0x61d027b3.
//
// Solidity: function treasury() view returns(address)
func (_TieredStaking *TieredStakingCaller) Treasury(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _TieredStaking.contract.Call(opts, &out, "treasury")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Treasury is a free data retrieval call binding the contract method 0x61d027b3.
//
// Solidity: function treasury() view returns(address)
func (_TieredStaking *TieredStakingSession) Treasury() (common.Address, error) {
	return _TieredStaking.Contract.Treasury(&_TieredStaking.CallOpts)
}

// Treasury is a free data retrieval call binding the contract method 0x61d027b3.
//
// Solidity: function treasury() view returns(address)
func (_TieredStaking *TieredStakingCallerSession) Treasury() (common.Address, error) {
	return _TieredStaking.Contract.Treasury(&_TieredStaking.CallOpts)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_TieredStaking *TieredStakingTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TieredStaking.contract.Transact(opts, "acceptOwnership")
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_TieredStaking *TieredStakingSession) AcceptOwnership() (*types.Transaction, error) {
	return _TieredStaking.Contract.AcceptOwnership(&_TieredStaking.TransactOpts)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_TieredStaking *TieredStakingTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _TieredStaking.Contract.AcceptOwnership(&_TieredStaking.TransactOpts)
}

// ClaimRewards is a paid mutator transaction binding the contract method 0x0962ef79.
//
// Solidity: function claimRewards(uint256 positionId) returns()
func (_TieredStaking *TieredStakingTransactor) ClaimRewards(opts *bind.TransactOpts, positionId *big.Int) (*types.Transaction, error) {
	return _TieredStaking.contract.Transact(opts, "claimRewards", positionId)
}

// ClaimRewards is a paid mutator transaction binding the contract method 0x0962ef79.
//
// Solidity: function claimRewards(uint256 positionId) returns()
func (_TieredStaking *TieredStakingSession) ClaimRewards(positionId *big.Int) (*types.Transaction, error) {
	return _TieredStaking.Contract.ClaimRewards(&_TieredStaking.TransactOpts, positionId)
}

// ClaimRewards is a paid mutator transaction binding the contract method 0x0962ef79.
//
// Solidity: function claimRewards(uint256 positionId) returns()
func (_TieredStaking *TieredStakingTransactorSession) ClaimRewards(positionId *big.Int) (*types.Transaction, error) {
	return _TieredStaking.Contract.ClaimRewards(&_TieredStaking.TransactOpts, positionId)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_TieredStaking *TieredStakingTransactor) Pause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TieredStaking.contract.Transact(opts, "pause")
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_TieredStaking *TieredStakingSession) Pause() (*types.Transaction, error) {
	return _TieredStaking.Contract.Pause(&_TieredStaking.TransactOpts)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_TieredStaking *TieredStakingTransactorSession) Pause() (*types.Transaction, error) {
	return _TieredStaking.Contract.Pause(&_TieredStaking.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_TieredStaking *TieredStakingTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TieredStaking.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_TieredStaking *TieredStakingSession) RenounceOwnership() (*types.Transaction, error) {
	return _TieredStaking.Contract.RenounceOwnership(&_TieredStaking.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_TieredStaking *TieredStakingTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _TieredStaking.Contract.RenounceOwnership(&_TieredStaking.TransactOpts)
}

// Stake is a paid mutator transaction binding the contract method 0x10087fb1.
//
// Solidity: function stake(uint256 amount, uint8 tier) returns(uint256 positionId)
func (_TieredStaking *TieredStakingTransactor) Stake(opts *bind.TransactOpts, amount *big.Int, tier uint8) (*types.Transaction, error) {
	return _TieredStaking.contract.Transact(opts, "stake", amount, tier)
}

// Stake is a paid mutator transaction binding the contract method 0x10087fb1.
//
// Solidity: function stake(uint256 amount, uint8 tier) returns(uint256 positionId)
func (_TieredStaking *TieredStakingSession) Stake(amount *big.Int, tier uint8) (*types.Transaction, error) {
	return _TieredStaking.Contract.Stake(&_TieredStaking.TransactOpts, amount, tier)
}

// Stake is a paid mutator transaction binding the contract method 0x10087fb1.
//
// Solidity: function stake(uint256 amount, uint8 tier) returns(uint256 positionId)
func (_TieredStaking *TieredStakingTransactorSession) Stake(amount *big.Int, tier uint8) (*types.Transaction, error) {
	return _TieredStaking.Contract.Stake(&_TieredStaking.TransactOpts, amount, tier)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_TieredStaking *TieredStakingTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _TieredStaking.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_TieredStaking *TieredStakingSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _TieredStaking.Contract.TransferOwnership(&_TieredStaking.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_TieredStaking *TieredStakingTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _TieredStaking.Contract.TransferOwnership(&_TieredStaking.TransactOpts, newOwner)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_TieredStaking *TieredStakingTransactor) Unpause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TieredStaking.contract.Transact(opts, "unpause")
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_TieredStaking *TieredStakingSession) Unpause() (*types.Transaction, error) {
	return _TieredStaking.Contract.Unpause(&_TieredStaking.TransactOpts)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_TieredStaking *TieredStakingTransactorSession) Unpause() (*types.Transaction, error) {
	return _TieredStaking.Contract.Unpause(&_TieredStaking.TransactOpts)
}

// Unstake is a paid mutator transaction binding the contract method 0x2e17de78.
//
// Solidity: function unstake(uint256 positionId) returns()
func (_TieredStaking *TieredStakingTransactor) Unstake(opts *bind.TransactOpts, positionId *big.Int) (*types.Transaction, error) {
	return _TieredStaking.contract.Transact(opts, "unstake", positionId)
}

// Unstake is a paid mutator transaction binding the contract method 0x2e17de78.
//
// Solidity: function unstake(uint256 positionId) returns()
func (_TieredStaking *TieredStakingSession) Unstake(positionId *big.Int) (*types.Transaction, error) {
	return _TieredStaking.Contract.Unstake(&_TieredStaking.TransactOpts, positionId)
}

// Unstake is a paid mutator transaction binding the contract method 0x2e17de78.
//
// Solidity: function unstake(uint256 positionId) returns()
func (_TieredStaking *TieredStakingTransactorSession) Unstake(positionId *big.Int) (*types.Transaction, error) {
	return _TieredStaking.Contract.Unstake(&_TieredStaking.TransactOpts, positionId)
}

// UpdateTierAPR is a paid mutator transaction binding the contract method 0xf5a5fe71.
//
// Solidity: function updateTierAPR(uint8 tier, uint256 newAprBps) returns()
func (_TieredStaking *TieredStakingTransactor) UpdateTierAPR(opts *bind.TransactOpts, tier uint8, newAprBps *big.Int) (*types.Transaction, error) {
	return _TieredStaking.contract.Transact(opts, "updateTierAPR", tier, newAprBps)
}

// UpdateTierAPR is a paid mutator transaction binding the contract method 0xf5a5fe71.
//
// Solidity: function updateTierAPR(uint8 tier, uint256 newAprBps) returns()
func (_TieredStaking *TieredStakingSession) UpdateTierAPR(tier uint8, newAprBps *big.Int) (*types.Transaction, error) {
	return _TieredStaking.Contract.UpdateTierAPR(&_TieredStaking.TransactOpts, tier, newAprBps)
}

// UpdateTierAPR is a paid mutator transaction binding the contract method 0xf5a5fe71.
//
// Solidity: function updateTierAPR(uint8 tier, uint256 newAprBps) returns()
func (_TieredStaking *TieredStakingTransactorSession) UpdateTierAPR(tier uint8, newAprBps *big.Int) (*types.Transaction, error) {
	return _TieredStaking.Contract.UpdateTierAPR(&_TieredStaking.TransactOpts, tier, newAprBps)
}

// UpdateTreasury is a paid mutator transaction binding the contract method 0x7f51bb1f.
//
// Solidity: function updateTreasury(address newTreasury) returns()
func (_TieredStaking *TieredStakingTransactor) UpdateTreasury(opts *bind.TransactOpts, newTreasury common.Address) (*types.Transaction, error) {
	return _TieredStaking.contract.Transact(opts, "updateTreasury", newTreasury)
}

// UpdateTreasury is a paid mutator transaction binding the contract method 0x7f51bb1f.
//
// Solidity: function updateTreasury(address newTreasury) returns()
func (_TieredStaking *TieredStakingSession) UpdateTreasury(newTreasury common.Address) (*types.Transaction, error) {
	return _TieredStaking.Contract.UpdateTreasury(&_TieredStaking.TransactOpts, newTreasury)
}

// UpdateTreasury is a paid mutator transaction binding the contract method 0x7f51bb1f.
//
// Solidity: function updateTreasury(address newTreasury) returns()
func (_TieredStaking *TieredStakingTransactorSession) UpdateTreasury(newTreasury common.Address) (*types.Transaction, error) {
	return _TieredStaking.Contract.UpdateTreasury(&_TieredStaking.TransactOpts, newTreasury)
}

// TieredStakingOwnershipTransferStartedIterator is returned from FilterOwnershipTransferStarted and is used to iterate over the raw logs and unpacked data for OwnershipTransferStarted events raised by the TieredStaking contract.
type TieredStakingOwnershipTransferStartedIterator struct {
	Event *TieredStakingOwnershipTransferStarted // Event containing the contract specifics and raw log

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
func (it *TieredStakingOwnershipTransferStartedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TieredStakingOwnershipTransferStarted)
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
		it.Event = new(TieredStakingOwnershipTransferStarted)
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
func (it *TieredStakingOwnershipTransferStartedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TieredStakingOwnershipTransferStartedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TieredStakingOwnershipTransferStarted represents a OwnershipTransferStarted event raised by the TieredStaking contract.
type TieredStakingOwnershipTransferStarted struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferStarted is a free log retrieval operation binding the contract event 0x38d16b8cac22d99fc7c124b9cd0de2d3fa1faef420bfe791d8c362d765e22700.
//
// Solidity: event OwnershipTransferStarted(address indexed previousOwner, address indexed newOwner)
func (_TieredStaking *TieredStakingFilterer) FilterOwnershipTransferStarted(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*TieredStakingOwnershipTransferStartedIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _TieredStaking.contract.FilterLogs(opts, "OwnershipTransferStarted", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &TieredStakingOwnershipTransferStartedIterator{contract: _TieredStaking.contract, event: "OwnershipTransferStarted", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferStarted is a free log subscription operation binding the contract event 0x38d16b8cac22d99fc7c124b9cd0de2d3fa1faef420bfe791d8c362d765e22700.
//
// Solidity: event OwnershipTransferStarted(address indexed previousOwner, address indexed newOwner)
func (_TieredStaking *TieredStakingFilterer) WatchOwnershipTransferStarted(opts *bind.WatchOpts, sink chan<- *TieredStakingOwnershipTransferStarted, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _TieredStaking.contract.WatchLogs(opts, "OwnershipTransferStarted", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TieredStakingOwnershipTransferStarted)
				if err := _TieredStaking.contract.UnpackLog(event, "OwnershipTransferStarted", log); err != nil {
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

// ParseOwnershipTransferStarted is a log parse operation binding the contract event 0x38d16b8cac22d99fc7c124b9cd0de2d3fa1faef420bfe791d8c362d765e22700.
//
// Solidity: event OwnershipTransferStarted(address indexed previousOwner, address indexed newOwner)
func (_TieredStaking *TieredStakingFilterer) ParseOwnershipTransferStarted(log types.Log) (*TieredStakingOwnershipTransferStarted, error) {
	event := new(TieredStakingOwnershipTransferStarted)
	if err := _TieredStaking.contract.UnpackLog(event, "OwnershipTransferStarted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TieredStakingOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the TieredStaking contract.
type TieredStakingOwnershipTransferredIterator struct {
	Event *TieredStakingOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *TieredStakingOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TieredStakingOwnershipTransferred)
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
		it.Event = new(TieredStakingOwnershipTransferred)
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
func (it *TieredStakingOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TieredStakingOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TieredStakingOwnershipTransferred represents a OwnershipTransferred event raised by the TieredStaking contract.
type TieredStakingOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_TieredStaking *TieredStakingFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*TieredStakingOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _TieredStaking.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &TieredStakingOwnershipTransferredIterator{contract: _TieredStaking.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_TieredStaking *TieredStakingFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *TieredStakingOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _TieredStaking.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TieredStakingOwnershipTransferred)
				if err := _TieredStaking.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_TieredStaking *TieredStakingFilterer) ParseOwnershipTransferred(log types.Log) (*TieredStakingOwnershipTransferred, error) {
	event := new(TieredStakingOwnershipTransferred)
	if err := _TieredStaking.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TieredStakingPausedIterator is returned from FilterPaused and is used to iterate over the raw logs and unpacked data for Paused events raised by the TieredStaking contract.
type TieredStakingPausedIterator struct {
	Event *TieredStakingPaused // Event containing the contract specifics and raw log

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
func (it *TieredStakingPausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TieredStakingPaused)
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
		it.Event = new(TieredStakingPaused)
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
func (it *TieredStakingPausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TieredStakingPausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TieredStakingPaused represents a Paused event raised by the TieredStaking contract.
type TieredStakingPaused struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterPaused is a free log retrieval operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_TieredStaking *TieredStakingFilterer) FilterPaused(opts *bind.FilterOpts) (*TieredStakingPausedIterator, error) {

	logs, sub, err := _TieredStaking.contract.FilterLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return &TieredStakingPausedIterator{contract: _TieredStaking.contract, event: "Paused", logs: logs, sub: sub}, nil
}

// WatchPaused is a free log subscription operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_TieredStaking *TieredStakingFilterer) WatchPaused(opts *bind.WatchOpts, sink chan<- *TieredStakingPaused) (event.Subscription, error) {

	logs, sub, err := _TieredStaking.contract.WatchLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TieredStakingPaused)
				if err := _TieredStaking.contract.UnpackLog(event, "Paused", log); err != nil {
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

// ParsePaused is a log parse operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_TieredStaking *TieredStakingFilterer) ParsePaused(log types.Log) (*TieredStakingPaused, error) {
	event := new(TieredStakingPaused)
	if err := _TieredStaking.contract.UnpackLog(event, "Paused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TieredStakingRewardsClaimedIterator is returned from FilterRewardsClaimed and is used to iterate over the raw logs and unpacked data for RewardsClaimed events raised by the TieredStaking contract.
type TieredStakingRewardsClaimedIterator struct {
	Event *TieredStakingRewardsClaimed // Event containing the contract specifics and raw log

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
func (it *TieredStakingRewardsClaimedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TieredStakingRewardsClaimed)
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
		it.Event = new(TieredStakingRewardsClaimed)
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
func (it *TieredStakingRewardsClaimedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TieredStakingRewardsClaimedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TieredStakingRewardsClaimed represents a RewardsClaimed event raised by the TieredStaking contract.
type TieredStakingRewardsClaimed struct {
	User       common.Address
	Amount     *big.Int
	PositionId *big.Int
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterRewardsClaimed is a free log retrieval operation binding the contract event 0xdacbdde355ba930696a362ea6738feb9f8bd52dfb3d81947558fd3217e23e325.
//
// Solidity: event RewardsClaimed(address indexed user, uint256 amount, uint256 positionId)
func (_TieredStaking *TieredStakingFilterer) FilterRewardsClaimed(opts *bind.FilterOpts, user []common.Address) (*TieredStakingRewardsClaimedIterator, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _TieredStaking.contract.FilterLogs(opts, "RewardsClaimed", userRule)
	if err != nil {
		return nil, err
	}
	return &TieredStakingRewardsClaimedIterator{contract: _TieredStaking.contract, event: "RewardsClaimed", logs: logs, sub: sub}, nil
}

// WatchRewardsClaimed is a free log subscription operation binding the contract event 0xdacbdde355ba930696a362ea6738feb9f8bd52dfb3d81947558fd3217e23e325.
//
// Solidity: event RewardsClaimed(address indexed user, uint256 amount, uint256 positionId)
func (_TieredStaking *TieredStakingFilterer) WatchRewardsClaimed(opts *bind.WatchOpts, sink chan<- *TieredStakingRewardsClaimed, user []common.Address) (event.Subscription, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _TieredStaking.contract.WatchLogs(opts, "RewardsClaimed", userRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TieredStakingRewardsClaimed)
				if err := _TieredStaking.contract.UnpackLog(event, "RewardsClaimed", log); err != nil {
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

// ParseRewardsClaimed is a log parse operation binding the contract event 0xdacbdde355ba930696a362ea6738feb9f8bd52dfb3d81947558fd3217e23e325.
//
// Solidity: event RewardsClaimed(address indexed user, uint256 amount, uint256 positionId)
func (_TieredStaking *TieredStakingFilterer) ParseRewardsClaimed(log types.Log) (*TieredStakingRewardsClaimed, error) {
	event := new(TieredStakingRewardsClaimed)
	if err := _TieredStaking.contract.UnpackLog(event, "RewardsClaimed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TieredStakingStakedIterator is returned from FilterStaked and is used to iterate over the raw logs and unpacked data for Staked events raised by the TieredStaking contract.
type TieredStakingStakedIterator struct {
	Event *TieredStakingStaked // Event containing the contract specifics and raw log

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
func (it *TieredStakingStakedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TieredStakingStaked)
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
		it.Event = new(TieredStakingStaked)
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
func (it *TieredStakingStakedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TieredStakingStakedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TieredStakingStaked represents a Staked event raised by the TieredStaking contract.
type TieredStakingStaked struct {
	User       common.Address
	Amount     *big.Int
	Tier       uint8
	PositionId *big.Int
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterStaked is a free log retrieval operation binding the contract event 0xbde7f0ba1630d25515c7ab99ba47d5640b7ffb4c673b2a5464ae679195589298.
//
// Solidity: event Staked(address indexed user, uint256 amount, uint8 tier, uint256 positionId)
func (_TieredStaking *TieredStakingFilterer) FilterStaked(opts *bind.FilterOpts, user []common.Address) (*TieredStakingStakedIterator, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _TieredStaking.contract.FilterLogs(opts, "Staked", userRule)
	if err != nil {
		return nil, err
	}
	return &TieredStakingStakedIterator{contract: _TieredStaking.contract, event: "Staked", logs: logs, sub: sub}, nil
}

// WatchStaked is a free log subscription operation binding the contract event 0xbde7f0ba1630d25515c7ab99ba47d5640b7ffb4c673b2a5464ae679195589298.
//
// Solidity: event Staked(address indexed user, uint256 amount, uint8 tier, uint256 positionId)
func (_TieredStaking *TieredStakingFilterer) WatchStaked(opts *bind.WatchOpts, sink chan<- *TieredStakingStaked, user []common.Address) (event.Subscription, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _TieredStaking.contract.WatchLogs(opts, "Staked", userRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TieredStakingStaked)
				if err := _TieredStaking.contract.UnpackLog(event, "Staked", log); err != nil {
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

// ParseStaked is a log parse operation binding the contract event 0xbde7f0ba1630d25515c7ab99ba47d5640b7ffb4c673b2a5464ae679195589298.
//
// Solidity: event Staked(address indexed user, uint256 amount, uint8 tier, uint256 positionId)
func (_TieredStaking *TieredStakingFilterer) ParseStaked(log types.Log) (*TieredStakingStaked, error) {
	event := new(TieredStakingStaked)
	if err := _TieredStaking.contract.UnpackLog(event, "Staked", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TieredStakingTierUpdatedIterator is returned from FilterTierUpdated and is used to iterate over the raw logs and unpacked data for TierUpdated events raised by the TieredStaking contract.
type TieredStakingTierUpdatedIterator struct {
	Event *TieredStakingTierUpdated // Event containing the contract specifics and raw log

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
func (it *TieredStakingTierUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TieredStakingTierUpdated)
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
		it.Event = new(TieredStakingTierUpdated)
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
func (it *TieredStakingTierUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TieredStakingTierUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TieredStakingTierUpdated represents a TierUpdated event raised by the TieredStaking contract.
type TieredStakingTierUpdated struct {
	Tier      uint8
	NewAprBps *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterTierUpdated is a free log retrieval operation binding the contract event 0x2c9f6fe32d88947c9ef965295ad014d64c351e869234a516a7d58ef465ca112d.
//
// Solidity: event TierUpdated(uint8 tier, uint256 newAprBps)
func (_TieredStaking *TieredStakingFilterer) FilterTierUpdated(opts *bind.FilterOpts) (*TieredStakingTierUpdatedIterator, error) {

	logs, sub, err := _TieredStaking.contract.FilterLogs(opts, "TierUpdated")
	if err != nil {
		return nil, err
	}
	return &TieredStakingTierUpdatedIterator{contract: _TieredStaking.contract, event: "TierUpdated", logs: logs, sub: sub}, nil
}

// WatchTierUpdated is a free log subscription operation binding the contract event 0x2c9f6fe32d88947c9ef965295ad014d64c351e869234a516a7d58ef465ca112d.
//
// Solidity: event TierUpdated(uint8 tier, uint256 newAprBps)
func (_TieredStaking *TieredStakingFilterer) WatchTierUpdated(opts *bind.WatchOpts, sink chan<- *TieredStakingTierUpdated) (event.Subscription, error) {

	logs, sub, err := _TieredStaking.contract.WatchLogs(opts, "TierUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TieredStakingTierUpdated)
				if err := _TieredStaking.contract.UnpackLog(event, "TierUpdated", log); err != nil {
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

// ParseTierUpdated is a log parse operation binding the contract event 0x2c9f6fe32d88947c9ef965295ad014d64c351e869234a516a7d58ef465ca112d.
//
// Solidity: event TierUpdated(uint8 tier, uint256 newAprBps)
func (_TieredStaking *TieredStakingFilterer) ParseTierUpdated(log types.Log) (*TieredStakingTierUpdated, error) {
	event := new(TieredStakingTierUpdated)
	if err := _TieredStaking.contract.UnpackLog(event, "TierUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TieredStakingTreasuryUpdatedIterator is returned from FilterTreasuryUpdated and is used to iterate over the raw logs and unpacked data for TreasuryUpdated events raised by the TieredStaking contract.
type TieredStakingTreasuryUpdatedIterator struct {
	Event *TieredStakingTreasuryUpdated // Event containing the contract specifics and raw log

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
func (it *TieredStakingTreasuryUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TieredStakingTreasuryUpdated)
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
		it.Event = new(TieredStakingTreasuryUpdated)
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
func (it *TieredStakingTreasuryUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TieredStakingTreasuryUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TieredStakingTreasuryUpdated represents a TreasuryUpdated event raised by the TieredStaking contract.
type TieredStakingTreasuryUpdated struct {
	NewTreasury common.Address
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterTreasuryUpdated is a free log retrieval operation binding the contract event 0x7dae230f18360d76a040c81f050aa14eb9d6dc7901b20fc5d855e2a20fe814d1.
//
// Solidity: event TreasuryUpdated(address newTreasury)
func (_TieredStaking *TieredStakingFilterer) FilterTreasuryUpdated(opts *bind.FilterOpts) (*TieredStakingTreasuryUpdatedIterator, error) {

	logs, sub, err := _TieredStaking.contract.FilterLogs(opts, "TreasuryUpdated")
	if err != nil {
		return nil, err
	}
	return &TieredStakingTreasuryUpdatedIterator{contract: _TieredStaking.contract, event: "TreasuryUpdated", logs: logs, sub: sub}, nil
}

// WatchTreasuryUpdated is a free log subscription operation binding the contract event 0x7dae230f18360d76a040c81f050aa14eb9d6dc7901b20fc5d855e2a20fe814d1.
//
// Solidity: event TreasuryUpdated(address newTreasury)
func (_TieredStaking *TieredStakingFilterer) WatchTreasuryUpdated(opts *bind.WatchOpts, sink chan<- *TieredStakingTreasuryUpdated) (event.Subscription, error) {

	logs, sub, err := _TieredStaking.contract.WatchLogs(opts, "TreasuryUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TieredStakingTreasuryUpdated)
				if err := _TieredStaking.contract.UnpackLog(event, "TreasuryUpdated", log); err != nil {
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

// ParseTreasuryUpdated is a log parse operation binding the contract event 0x7dae230f18360d76a040c81f050aa14eb9d6dc7901b20fc5d855e2a20fe814d1.
//
// Solidity: event TreasuryUpdated(address newTreasury)
func (_TieredStaking *TieredStakingFilterer) ParseTreasuryUpdated(log types.Log) (*TieredStakingTreasuryUpdated, error) {
	event := new(TieredStakingTreasuryUpdated)
	if err := _TieredStaking.contract.UnpackLog(event, "TreasuryUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TieredStakingUnpausedIterator is returned from FilterUnpaused and is used to iterate over the raw logs and unpacked data for Unpaused events raised by the TieredStaking contract.
type TieredStakingUnpausedIterator struct {
	Event *TieredStakingUnpaused // Event containing the contract specifics and raw log

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
func (it *TieredStakingUnpausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TieredStakingUnpaused)
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
		it.Event = new(TieredStakingUnpaused)
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
func (it *TieredStakingUnpausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TieredStakingUnpausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TieredStakingUnpaused represents a Unpaused event raised by the TieredStaking contract.
type TieredStakingUnpaused struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterUnpaused is a free log retrieval operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_TieredStaking *TieredStakingFilterer) FilterUnpaused(opts *bind.FilterOpts) (*TieredStakingUnpausedIterator, error) {

	logs, sub, err := _TieredStaking.contract.FilterLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return &TieredStakingUnpausedIterator{contract: _TieredStaking.contract, event: "Unpaused", logs: logs, sub: sub}, nil
}

// WatchUnpaused is a free log subscription operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_TieredStaking *TieredStakingFilterer) WatchUnpaused(opts *bind.WatchOpts, sink chan<- *TieredStakingUnpaused) (event.Subscription, error) {

	logs, sub, err := _TieredStaking.contract.WatchLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TieredStakingUnpaused)
				if err := _TieredStaking.contract.UnpackLog(event, "Unpaused", log); err != nil {
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

// ParseUnpaused is a log parse operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_TieredStaking *TieredStakingFilterer) ParseUnpaused(log types.Log) (*TieredStakingUnpaused, error) {
	event := new(TieredStakingUnpaused)
	if err := _TieredStaking.contract.UnpackLog(event, "Unpaused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TieredStakingUnstakedIterator is returned from FilterUnstaked and is used to iterate over the raw logs and unpacked data for Unstaked events raised by the TieredStaking contract.
type TieredStakingUnstakedIterator struct {
	Event *TieredStakingUnstaked // Event containing the contract specifics and raw log

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
func (it *TieredStakingUnstakedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TieredStakingUnstaked)
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
		it.Event = new(TieredStakingUnstaked)
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
func (it *TieredStakingUnstakedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TieredStakingUnstakedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TieredStakingUnstaked represents a Unstaked event raised by the TieredStaking contract.
type TieredStakingUnstaked struct {
	User    common.Address
	Amount  *big.Int
	Rewards *big.Int
	Penalty *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterUnstaked is a free log retrieval operation binding the contract event 0x204fccf0d92ed8d48f204adb39b2e81e92bad0dedb93f5716ca9478cfb57de00.
//
// Solidity: event Unstaked(address indexed user, uint256 amount, uint256 rewards, uint256 penalty)
func (_TieredStaking *TieredStakingFilterer) FilterUnstaked(opts *bind.FilterOpts, user []common.Address) (*TieredStakingUnstakedIterator, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _TieredStaking.contract.FilterLogs(opts, "Unstaked", userRule)
	if err != nil {
		return nil, err
	}
	return &TieredStakingUnstakedIterator{contract: _TieredStaking.contract, event: "Unstaked", logs: logs, sub: sub}, nil
}

// WatchUnstaked is a free log subscription operation binding the contract event 0x204fccf0d92ed8d48f204adb39b2e81e92bad0dedb93f5716ca9478cfb57de00.
//
// Solidity: event Unstaked(address indexed user, uint256 amount, uint256 rewards, uint256 penalty)
func (_TieredStaking *TieredStakingFilterer) WatchUnstaked(opts *bind.WatchOpts, sink chan<- *TieredStakingUnstaked, user []common.Address) (event.Subscription, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _TieredStaking.contract.WatchLogs(opts, "Unstaked", userRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TieredStakingUnstaked)
				if err := _TieredStaking.contract.UnpackLog(event, "Unstaked", log); err != nil {
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

// ParseUnstaked is a log parse operation binding the contract event 0x204fccf0d92ed8d48f204adb39b2e81e92bad0dedb93f5716ca9478cfb57de00.
//
// Solidity: event Unstaked(address indexed user, uint256 amount, uint256 rewards, uint256 penalty)
func (_TieredStaking *TieredStakingFilterer) ParseUnstaked(log types.Log) (*TieredStakingUnstaked, error) {
	event := new(TieredStakingUnstaked)
	if err := _TieredStaking.contract.UnpackLog(event, "Unstaked", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
