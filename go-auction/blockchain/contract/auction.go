// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contract

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

// AuctionAuctionInfo is an auto generated low-level Go binding around an user-defined struct.
type AuctionAuctionInfo struct {
	NftContract       common.Address
	TokenId           *big.Int
	Seller            common.Address
	StartingPrice     *big.Int
	CurrentHighestBid *big.Int
	HighestBidder     common.Address
	EndTime           *big.Int
	PaymentToken      common.Address
	Ended             bool
}

// AuctionMetaData contains all meta data concerning the Auction contract.
var AuctionMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"name\":\"InvalidInitialization\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NotInitializing\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"OwnableInvalidOwner\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"OwnableUnauthorizedAccount\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ReentrancyGuardReentrantCall\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"auctionId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"nftContract\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"seller\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"startingPrice\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"endTime\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"paymentToken\",\"type\":\"address\"}],\"name\":\"AuctionCreated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"auctionId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"winner\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"finalPriceUSD\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"paymentToken\",\"type\":\"address\"}],\"name\":\"AuctionEnded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"auctionId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"bidder\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"bidAmountUSD\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"originalAmount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"paymentToken\",\"type\":\"address\"}],\"name\":\"BidPlaced\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"auctionId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"bidder\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"BidWithdrawn\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"version\",\"type\":\"uint64\"}],\"name\":\"Initialized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"BASIS_POINTS\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"FEE_BASIS_POINTS\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"auctions\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"nftContract\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"seller\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"startingPrice\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"currentHighestBid\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"highestBidder\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"endTime\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"paymentToken\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"ended\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"auctionId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"bidWithERC20\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"auctionId\",\"type\":\"uint256\"}],\"name\":\"bidWithETH\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"bids\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"nftContract\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"startingPriceUSD\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"duration\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"paymentToken\",\"type\":\"address\"}],\"name\":\"createAuction\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"auctionId\",\"type\":\"uint256\"}],\"name\":\"endAuction\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"auctionId\",\"type\":\"uint256\"}],\"name\":\"getAuction\",\"outputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"nftContract\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"seller\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"startingPrice\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"currentHighestBid\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"highestBidder\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"endTime\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"paymentToken\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"ended\",\"type\":\"bool\"}],\"internalType\":\"structAuction.AuctionInfo\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getCurrentAuctionId\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_priceOracle\",\"type\":\"address\"}],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"name\":\"onERC721Received\",\"outputs\":[{\"internalType\":\"bytes4\",\"name\":\"\",\"type\":\"bytes4\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"priceOracle\",\"outputs\":[{\"internalType\":\"contractPriceOracle\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"auctionId\",\"type\":\"uint256\"}],\"name\":\"withdrawBid\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"withdrawableBids\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
}

// AuctionABI is the input ABI used to generate the binding from.
// Deprecated: Use AuctionMetaData.ABI instead.
var AuctionABI = AuctionMetaData.ABI

// Auction is an auto generated Go binding around an Ethereum contract.
type Auction struct {
	AuctionCaller     // Read-only binding to the contract
	AuctionTransactor // Write-only binding to the contract
	AuctionFilterer   // Log filterer for contract events
}

// AuctionCaller is an auto generated read-only Go binding around an Ethereum contract.
type AuctionCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AuctionTransactor is an auto generated write-only Go binding around an Ethereum contract.
type AuctionTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AuctionFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type AuctionFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AuctionSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type AuctionSession struct {
	Contract     *Auction          // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// AuctionCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type AuctionCallerSession struct {
	Contract *AuctionCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts  // Call options to use throughout this session
}

// AuctionTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type AuctionTransactorSession struct {
	Contract     *AuctionTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// AuctionRaw is an auto generated low-level Go binding around an Ethereum contract.
type AuctionRaw struct {
	Contract *Auction // Generic contract binding to access the raw methods on
}

// AuctionCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type AuctionCallerRaw struct {
	Contract *AuctionCaller // Generic read-only contract binding to access the raw methods on
}

// AuctionTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type AuctionTransactorRaw struct {
	Contract *AuctionTransactor // Generic write-only contract binding to access the raw methods on
}

// NewAuction creates a new instance of Auction, bound to a specific deployed contract.
func NewAuction(address common.Address, backend bind.ContractBackend) (*Auction, error) {
	contract, err := bindAuction(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Auction{AuctionCaller: AuctionCaller{contract: contract}, AuctionTransactor: AuctionTransactor{contract: contract}, AuctionFilterer: AuctionFilterer{contract: contract}}, nil
}

// NewAuctionCaller creates a new read-only instance of Auction, bound to a specific deployed contract.
func NewAuctionCaller(address common.Address, caller bind.ContractCaller) (*AuctionCaller, error) {
	contract, err := bindAuction(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &AuctionCaller{contract: contract}, nil
}

// NewAuctionTransactor creates a new write-only instance of Auction, bound to a specific deployed contract.
func NewAuctionTransactor(address common.Address, transactor bind.ContractTransactor) (*AuctionTransactor, error) {
	contract, err := bindAuction(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &AuctionTransactor{contract: contract}, nil
}

// NewAuctionFilterer creates a new log filterer instance of Auction, bound to a specific deployed contract.
func NewAuctionFilterer(address common.Address, filterer bind.ContractFilterer) (*AuctionFilterer, error) {
	contract, err := bindAuction(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &AuctionFilterer{contract: contract}, nil
}

// bindAuction binds a generic wrapper to an already deployed contract.
func bindAuction(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := AuctionMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Auction *AuctionRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Auction.Contract.AuctionCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Auction *AuctionRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Auction.Contract.AuctionTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Auction *AuctionRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Auction.Contract.AuctionTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Auction *AuctionCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Auction.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Auction *AuctionTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Auction.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Auction *AuctionTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Auction.Contract.contract.Transact(opts, method, params...)
}

// BASISPOINTS is a free data retrieval call binding the contract method 0xe1f1c4a7.
//
// Solidity: function BASIS_POINTS() view returns(uint256)
func (_Auction *AuctionCaller) BASISPOINTS(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Auction.contract.Call(opts, &out, "BASIS_POINTS")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BASISPOINTS is a free data retrieval call binding the contract method 0xe1f1c4a7.
//
// Solidity: function BASIS_POINTS() view returns(uint256)
func (_Auction *AuctionSession) BASISPOINTS() (*big.Int, error) {
	return _Auction.Contract.BASISPOINTS(&_Auction.CallOpts)
}

// BASISPOINTS is a free data retrieval call binding the contract method 0xe1f1c4a7.
//
// Solidity: function BASIS_POINTS() view returns(uint256)
func (_Auction *AuctionCallerSession) BASISPOINTS() (*big.Int, error) {
	return _Auction.Contract.BASISPOINTS(&_Auction.CallOpts)
}

// FEEBASISPOINTS is a free data retrieval call binding the contract method 0xa525ad3c.
//
// Solidity: function FEE_BASIS_POINTS() view returns(uint256)
func (_Auction *AuctionCaller) FEEBASISPOINTS(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Auction.contract.Call(opts, &out, "FEE_BASIS_POINTS")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// FEEBASISPOINTS is a free data retrieval call binding the contract method 0xa525ad3c.
//
// Solidity: function FEE_BASIS_POINTS() view returns(uint256)
func (_Auction *AuctionSession) FEEBASISPOINTS() (*big.Int, error) {
	return _Auction.Contract.FEEBASISPOINTS(&_Auction.CallOpts)
}

// FEEBASISPOINTS is a free data retrieval call binding the contract method 0xa525ad3c.
//
// Solidity: function FEE_BASIS_POINTS() view returns(uint256)
func (_Auction *AuctionCallerSession) FEEBASISPOINTS() (*big.Int, error) {
	return _Auction.Contract.FEEBASISPOINTS(&_Auction.CallOpts)
}

// Auctions is a free data retrieval call binding the contract method 0x571a26a0.
//
// Solidity: function auctions(uint256 ) view returns(address nftContract, uint256 tokenId, address seller, uint256 startingPrice, uint256 currentHighestBid, address highestBidder, uint256 endTime, address paymentToken, bool ended)
func (_Auction *AuctionCaller) Auctions(opts *bind.CallOpts, arg0 *big.Int) (struct {
	NftContract       common.Address
	TokenId           *big.Int
	Seller            common.Address
	StartingPrice     *big.Int
	CurrentHighestBid *big.Int
	HighestBidder     common.Address
	EndTime           *big.Int
	PaymentToken      common.Address
	Ended             bool
}, error) {
	var out []interface{}
	err := _Auction.contract.Call(opts, &out, "auctions", arg0)

	outstruct := new(struct {
		NftContract       common.Address
		TokenId           *big.Int
		Seller            common.Address
		StartingPrice     *big.Int
		CurrentHighestBid *big.Int
		HighestBidder     common.Address
		EndTime           *big.Int
		PaymentToken      common.Address
		Ended             bool
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.NftContract = *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	outstruct.TokenId = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.Seller = *abi.ConvertType(out[2], new(common.Address)).(*common.Address)
	outstruct.StartingPrice = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)
	outstruct.CurrentHighestBid = *abi.ConvertType(out[4], new(*big.Int)).(**big.Int)
	outstruct.HighestBidder = *abi.ConvertType(out[5], new(common.Address)).(*common.Address)
	outstruct.EndTime = *abi.ConvertType(out[6], new(*big.Int)).(**big.Int)
	outstruct.PaymentToken = *abi.ConvertType(out[7], new(common.Address)).(*common.Address)
	outstruct.Ended = *abi.ConvertType(out[8], new(bool)).(*bool)

	return *outstruct, err

}

// Auctions is a free data retrieval call binding the contract method 0x571a26a0.
//
// Solidity: function auctions(uint256 ) view returns(address nftContract, uint256 tokenId, address seller, uint256 startingPrice, uint256 currentHighestBid, address highestBidder, uint256 endTime, address paymentToken, bool ended)
func (_Auction *AuctionSession) Auctions(arg0 *big.Int) (struct {
	NftContract       common.Address
	TokenId           *big.Int
	Seller            common.Address
	StartingPrice     *big.Int
	CurrentHighestBid *big.Int
	HighestBidder     common.Address
	EndTime           *big.Int
	PaymentToken      common.Address
	Ended             bool
}, error) {
	return _Auction.Contract.Auctions(&_Auction.CallOpts, arg0)
}

// Auctions is a free data retrieval call binding the contract method 0x571a26a0.
//
// Solidity: function auctions(uint256 ) view returns(address nftContract, uint256 tokenId, address seller, uint256 startingPrice, uint256 currentHighestBid, address highestBidder, uint256 endTime, address paymentToken, bool ended)
func (_Auction *AuctionCallerSession) Auctions(arg0 *big.Int) (struct {
	NftContract       common.Address
	TokenId           *big.Int
	Seller            common.Address
	StartingPrice     *big.Int
	CurrentHighestBid *big.Int
	HighestBidder     common.Address
	EndTime           *big.Int
	PaymentToken      common.Address
	Ended             bool
}, error) {
	return _Auction.Contract.Auctions(&_Auction.CallOpts, arg0)
}

// Bids is a free data retrieval call binding the contract method 0x3f1ffcec.
//
// Solidity: function bids(uint256 , address ) view returns(uint256)
func (_Auction *AuctionCaller) Bids(opts *bind.CallOpts, arg0 *big.Int, arg1 common.Address) (*big.Int, error) {
	var out []interface{}
	err := _Auction.contract.Call(opts, &out, "bids", arg0, arg1)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Bids is a free data retrieval call binding the contract method 0x3f1ffcec.
//
// Solidity: function bids(uint256 , address ) view returns(uint256)
func (_Auction *AuctionSession) Bids(arg0 *big.Int, arg1 common.Address) (*big.Int, error) {
	return _Auction.Contract.Bids(&_Auction.CallOpts, arg0, arg1)
}

// Bids is a free data retrieval call binding the contract method 0x3f1ffcec.
//
// Solidity: function bids(uint256 , address ) view returns(uint256)
func (_Auction *AuctionCallerSession) Bids(arg0 *big.Int, arg1 common.Address) (*big.Int, error) {
	return _Auction.Contract.Bids(&_Auction.CallOpts, arg0, arg1)
}

// GetAuction is a free data retrieval call binding the contract method 0x78bd7935.
//
// Solidity: function getAuction(uint256 auctionId) view returns((address,uint256,address,uint256,uint256,address,uint256,address,bool))
func (_Auction *AuctionCaller) GetAuction(opts *bind.CallOpts, auctionId *big.Int) (AuctionAuctionInfo, error) {
	var out []interface{}
	err := _Auction.contract.Call(opts, &out, "getAuction", auctionId)

	if err != nil {
		return *new(AuctionAuctionInfo), err
	}

	out0 := *abi.ConvertType(out[0], new(AuctionAuctionInfo)).(*AuctionAuctionInfo)

	return out0, err

}

// GetAuction is a free data retrieval call binding the contract method 0x78bd7935.
//
// Solidity: function getAuction(uint256 auctionId) view returns((address,uint256,address,uint256,uint256,address,uint256,address,bool))
func (_Auction *AuctionSession) GetAuction(auctionId *big.Int) (AuctionAuctionInfo, error) {
	return _Auction.Contract.GetAuction(&_Auction.CallOpts, auctionId)
}

// GetAuction is a free data retrieval call binding the contract method 0x78bd7935.
//
// Solidity: function getAuction(uint256 auctionId) view returns((address,uint256,address,uint256,uint256,address,uint256,address,bool))
func (_Auction *AuctionCallerSession) GetAuction(auctionId *big.Int) (AuctionAuctionInfo, error) {
	return _Auction.Contract.GetAuction(&_Auction.CallOpts, auctionId)
}

// GetCurrentAuctionId is a free data retrieval call binding the contract method 0x157db316.
//
// Solidity: function getCurrentAuctionId() view returns(uint256)
func (_Auction *AuctionCaller) GetCurrentAuctionId(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Auction.contract.Call(opts, &out, "getCurrentAuctionId")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetCurrentAuctionId is a free data retrieval call binding the contract method 0x157db316.
//
// Solidity: function getCurrentAuctionId() view returns(uint256)
func (_Auction *AuctionSession) GetCurrentAuctionId() (*big.Int, error) {
	return _Auction.Contract.GetCurrentAuctionId(&_Auction.CallOpts)
}

// GetCurrentAuctionId is a free data retrieval call binding the contract method 0x157db316.
//
// Solidity: function getCurrentAuctionId() view returns(uint256)
func (_Auction *AuctionCallerSession) GetCurrentAuctionId() (*big.Int, error) {
	return _Auction.Contract.GetCurrentAuctionId(&_Auction.CallOpts)
}

// OnERC721Received is a free data retrieval call binding the contract method 0x150b7a02.
//
// Solidity: function onERC721Received(address , address , uint256 , bytes ) pure returns(bytes4)
func (_Auction *AuctionCaller) OnERC721Received(opts *bind.CallOpts, arg0 common.Address, arg1 common.Address, arg2 *big.Int, arg3 []byte) ([4]byte, error) {
	var out []interface{}
	err := _Auction.contract.Call(opts, &out, "onERC721Received", arg0, arg1, arg2, arg3)

	if err != nil {
		return *new([4]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([4]byte)).(*[4]byte)

	return out0, err

}

// OnERC721Received is a free data retrieval call binding the contract method 0x150b7a02.
//
// Solidity: function onERC721Received(address , address , uint256 , bytes ) pure returns(bytes4)
func (_Auction *AuctionSession) OnERC721Received(arg0 common.Address, arg1 common.Address, arg2 *big.Int, arg3 []byte) ([4]byte, error) {
	return _Auction.Contract.OnERC721Received(&_Auction.CallOpts, arg0, arg1, arg2, arg3)
}

// OnERC721Received is a free data retrieval call binding the contract method 0x150b7a02.
//
// Solidity: function onERC721Received(address , address , uint256 , bytes ) pure returns(bytes4)
func (_Auction *AuctionCallerSession) OnERC721Received(arg0 common.Address, arg1 common.Address, arg2 *big.Int, arg3 []byte) ([4]byte, error) {
	return _Auction.Contract.OnERC721Received(&_Auction.CallOpts, arg0, arg1, arg2, arg3)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Auction *AuctionCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Auction.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Auction *AuctionSession) Owner() (common.Address, error) {
	return _Auction.Contract.Owner(&_Auction.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Auction *AuctionCallerSession) Owner() (common.Address, error) {
	return _Auction.Contract.Owner(&_Auction.CallOpts)
}

// PriceOracle is a free data retrieval call binding the contract method 0x2630c12f.
//
// Solidity: function priceOracle() view returns(address)
func (_Auction *AuctionCaller) PriceOracle(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Auction.contract.Call(opts, &out, "priceOracle")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// PriceOracle is a free data retrieval call binding the contract method 0x2630c12f.
//
// Solidity: function priceOracle() view returns(address)
func (_Auction *AuctionSession) PriceOracle() (common.Address, error) {
	return _Auction.Contract.PriceOracle(&_Auction.CallOpts)
}

// PriceOracle is a free data retrieval call binding the contract method 0x2630c12f.
//
// Solidity: function priceOracle() view returns(address)
func (_Auction *AuctionCallerSession) PriceOracle() (common.Address, error) {
	return _Auction.Contract.PriceOracle(&_Auction.CallOpts)
}

// WithdrawableBids is a free data retrieval call binding the contract method 0x03711682.
//
// Solidity: function withdrawableBids(uint256 , address ) view returns(uint256)
func (_Auction *AuctionCaller) WithdrawableBids(opts *bind.CallOpts, arg0 *big.Int, arg1 common.Address) (*big.Int, error) {
	var out []interface{}
	err := _Auction.contract.Call(opts, &out, "withdrawableBids", arg0, arg1)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// WithdrawableBids is a free data retrieval call binding the contract method 0x03711682.
//
// Solidity: function withdrawableBids(uint256 , address ) view returns(uint256)
func (_Auction *AuctionSession) WithdrawableBids(arg0 *big.Int, arg1 common.Address) (*big.Int, error) {
	return _Auction.Contract.WithdrawableBids(&_Auction.CallOpts, arg0, arg1)
}

// WithdrawableBids is a free data retrieval call binding the contract method 0x03711682.
//
// Solidity: function withdrawableBids(uint256 , address ) view returns(uint256)
func (_Auction *AuctionCallerSession) WithdrawableBids(arg0 *big.Int, arg1 common.Address) (*big.Int, error) {
	return _Auction.Contract.WithdrawableBids(&_Auction.CallOpts, arg0, arg1)
}

// BidWithERC20 is a paid mutator transaction binding the contract method 0x21f0485d.
//
// Solidity: function bidWithERC20(uint256 auctionId, uint256 amount) returns()
func (_Auction *AuctionTransactor) BidWithERC20(opts *bind.TransactOpts, auctionId *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _Auction.contract.Transact(opts, "bidWithERC20", auctionId, amount)
}

// BidWithERC20 is a paid mutator transaction binding the contract method 0x21f0485d.
//
// Solidity: function bidWithERC20(uint256 auctionId, uint256 amount) returns()
func (_Auction *AuctionSession) BidWithERC20(auctionId *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _Auction.Contract.BidWithERC20(&_Auction.TransactOpts, auctionId, amount)
}

// BidWithERC20 is a paid mutator transaction binding the contract method 0x21f0485d.
//
// Solidity: function bidWithERC20(uint256 auctionId, uint256 amount) returns()
func (_Auction *AuctionTransactorSession) BidWithERC20(auctionId *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _Auction.Contract.BidWithERC20(&_Auction.TransactOpts, auctionId, amount)
}

// BidWithETH is a paid mutator transaction binding the contract method 0xac9ad632.
//
// Solidity: function bidWithETH(uint256 auctionId) payable returns()
func (_Auction *AuctionTransactor) BidWithETH(opts *bind.TransactOpts, auctionId *big.Int) (*types.Transaction, error) {
	return _Auction.contract.Transact(opts, "bidWithETH", auctionId)
}

// BidWithETH is a paid mutator transaction binding the contract method 0xac9ad632.
//
// Solidity: function bidWithETH(uint256 auctionId) payable returns()
func (_Auction *AuctionSession) BidWithETH(auctionId *big.Int) (*types.Transaction, error) {
	return _Auction.Contract.BidWithETH(&_Auction.TransactOpts, auctionId)
}

// BidWithETH is a paid mutator transaction binding the contract method 0xac9ad632.
//
// Solidity: function bidWithETH(uint256 auctionId) payable returns()
func (_Auction *AuctionTransactorSession) BidWithETH(auctionId *big.Int) (*types.Transaction, error) {
	return _Auction.Contract.BidWithETH(&_Auction.TransactOpts, auctionId)
}

// CreateAuction is a paid mutator transaction binding the contract method 0xb0b211da.
//
// Solidity: function createAuction(address nftContract, uint256 tokenId, uint256 startingPriceUSD, uint256 duration, address paymentToken) returns()
func (_Auction *AuctionTransactor) CreateAuction(opts *bind.TransactOpts, nftContract common.Address, tokenId *big.Int, startingPriceUSD *big.Int, duration *big.Int, paymentToken common.Address) (*types.Transaction, error) {
	return _Auction.contract.Transact(opts, "createAuction", nftContract, tokenId, startingPriceUSD, duration, paymentToken)
}

// CreateAuction is a paid mutator transaction binding the contract method 0xb0b211da.
//
// Solidity: function createAuction(address nftContract, uint256 tokenId, uint256 startingPriceUSD, uint256 duration, address paymentToken) returns()
func (_Auction *AuctionSession) CreateAuction(nftContract common.Address, tokenId *big.Int, startingPriceUSD *big.Int, duration *big.Int, paymentToken common.Address) (*types.Transaction, error) {
	return _Auction.Contract.CreateAuction(&_Auction.TransactOpts, nftContract, tokenId, startingPriceUSD, duration, paymentToken)
}

// CreateAuction is a paid mutator transaction binding the contract method 0xb0b211da.
//
// Solidity: function createAuction(address nftContract, uint256 tokenId, uint256 startingPriceUSD, uint256 duration, address paymentToken) returns()
func (_Auction *AuctionTransactorSession) CreateAuction(nftContract common.Address, tokenId *big.Int, startingPriceUSD *big.Int, duration *big.Int, paymentToken common.Address) (*types.Transaction, error) {
	return _Auction.Contract.CreateAuction(&_Auction.TransactOpts, nftContract, tokenId, startingPriceUSD, duration, paymentToken)
}

// EndAuction is a paid mutator transaction binding the contract method 0xb9a2de3a.
//
// Solidity: function endAuction(uint256 auctionId) returns()
func (_Auction *AuctionTransactor) EndAuction(opts *bind.TransactOpts, auctionId *big.Int) (*types.Transaction, error) {
	return _Auction.contract.Transact(opts, "endAuction", auctionId)
}

// EndAuction is a paid mutator transaction binding the contract method 0xb9a2de3a.
//
// Solidity: function endAuction(uint256 auctionId) returns()
func (_Auction *AuctionSession) EndAuction(auctionId *big.Int) (*types.Transaction, error) {
	return _Auction.Contract.EndAuction(&_Auction.TransactOpts, auctionId)
}

// EndAuction is a paid mutator transaction binding the contract method 0xb9a2de3a.
//
// Solidity: function endAuction(uint256 auctionId) returns()
func (_Auction *AuctionTransactorSession) EndAuction(auctionId *big.Int) (*types.Transaction, error) {
	return _Auction.Contract.EndAuction(&_Auction.TransactOpts, auctionId)
}

// Initialize is a paid mutator transaction binding the contract method 0xc4d66de8.
//
// Solidity: function initialize(address _priceOracle) returns()
func (_Auction *AuctionTransactor) Initialize(opts *bind.TransactOpts, _priceOracle common.Address) (*types.Transaction, error) {
	return _Auction.contract.Transact(opts, "initialize", _priceOracle)
}

// Initialize is a paid mutator transaction binding the contract method 0xc4d66de8.
//
// Solidity: function initialize(address _priceOracle) returns()
func (_Auction *AuctionSession) Initialize(_priceOracle common.Address) (*types.Transaction, error) {
	return _Auction.Contract.Initialize(&_Auction.TransactOpts, _priceOracle)
}

// Initialize is a paid mutator transaction binding the contract method 0xc4d66de8.
//
// Solidity: function initialize(address _priceOracle) returns()
func (_Auction *AuctionTransactorSession) Initialize(_priceOracle common.Address) (*types.Transaction, error) {
	return _Auction.Contract.Initialize(&_Auction.TransactOpts, _priceOracle)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Auction *AuctionTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Auction.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Auction *AuctionSession) RenounceOwnership() (*types.Transaction, error) {
	return _Auction.Contract.RenounceOwnership(&_Auction.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Auction *AuctionTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _Auction.Contract.RenounceOwnership(&_Auction.TransactOpts)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Auction *AuctionTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _Auction.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Auction *AuctionSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Auction.Contract.TransferOwnership(&_Auction.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Auction *AuctionTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Auction.Contract.TransferOwnership(&_Auction.TransactOpts, newOwner)
}

// WithdrawBid is a paid mutator transaction binding the contract method 0x0eaaf4c8.
//
// Solidity: function withdrawBid(uint256 auctionId) returns()
func (_Auction *AuctionTransactor) WithdrawBid(opts *bind.TransactOpts, auctionId *big.Int) (*types.Transaction, error) {
	return _Auction.contract.Transact(opts, "withdrawBid", auctionId)
}

// WithdrawBid is a paid mutator transaction binding the contract method 0x0eaaf4c8.
//
// Solidity: function withdrawBid(uint256 auctionId) returns()
func (_Auction *AuctionSession) WithdrawBid(auctionId *big.Int) (*types.Transaction, error) {
	return _Auction.Contract.WithdrawBid(&_Auction.TransactOpts, auctionId)
}

// WithdrawBid is a paid mutator transaction binding the contract method 0x0eaaf4c8.
//
// Solidity: function withdrawBid(uint256 auctionId) returns()
func (_Auction *AuctionTransactorSession) WithdrawBid(auctionId *big.Int) (*types.Transaction, error) {
	return _Auction.Contract.WithdrawBid(&_Auction.TransactOpts, auctionId)
}

// AuctionAuctionCreatedIterator is returned from FilterAuctionCreated and is used to iterate over the raw logs and unpacked data for AuctionCreated events raised by the Auction contract.
type AuctionAuctionCreatedIterator struct {
	Event *AuctionAuctionCreated // Event containing the contract specifics and raw log

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
func (it *AuctionAuctionCreatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AuctionAuctionCreated)
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
		it.Event = new(AuctionAuctionCreated)
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
func (it *AuctionAuctionCreatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AuctionAuctionCreatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AuctionAuctionCreated represents a AuctionCreated event raised by the Auction contract.
type AuctionAuctionCreated struct {
	AuctionId     *big.Int
	NftContract   common.Address
	TokenId       *big.Int
	Seller        common.Address
	StartingPrice *big.Int
	EndTime       *big.Int
	PaymentToken  common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterAuctionCreated is a free log retrieval operation binding the contract event 0x35307ae78277f7a6939b4c091fee2c5a1fd6d1c40ed27ee46bf3338b3f47bc06.
//
// Solidity: event AuctionCreated(uint256 indexed auctionId, address indexed nftContract, uint256 indexed tokenId, address seller, uint256 startingPrice, uint256 endTime, address paymentToken)
func (_Auction *AuctionFilterer) FilterAuctionCreated(opts *bind.FilterOpts, auctionId []*big.Int, nftContract []common.Address, tokenId []*big.Int) (*AuctionAuctionCreatedIterator, error) {

	var auctionIdRule []interface{}
	for _, auctionIdItem := range auctionId {
		auctionIdRule = append(auctionIdRule, auctionIdItem)
	}
	var nftContractRule []interface{}
	for _, nftContractItem := range nftContract {
		nftContractRule = append(nftContractRule, nftContractItem)
	}
	var tokenIdRule []interface{}
	for _, tokenIdItem := range tokenId {
		tokenIdRule = append(tokenIdRule, tokenIdItem)
	}

	logs, sub, err := _Auction.contract.FilterLogs(opts, "AuctionCreated", auctionIdRule, nftContractRule, tokenIdRule)
	if err != nil {
		return nil, err
	}
	return &AuctionAuctionCreatedIterator{contract: _Auction.contract, event: "AuctionCreated", logs: logs, sub: sub}, nil
}

// WatchAuctionCreated is a free log subscription operation binding the contract event 0x35307ae78277f7a6939b4c091fee2c5a1fd6d1c40ed27ee46bf3338b3f47bc06.
//
// Solidity: event AuctionCreated(uint256 indexed auctionId, address indexed nftContract, uint256 indexed tokenId, address seller, uint256 startingPrice, uint256 endTime, address paymentToken)
func (_Auction *AuctionFilterer) WatchAuctionCreated(opts *bind.WatchOpts, sink chan<- *AuctionAuctionCreated, auctionId []*big.Int, nftContract []common.Address, tokenId []*big.Int) (event.Subscription, error) {

	var auctionIdRule []interface{}
	for _, auctionIdItem := range auctionId {
		auctionIdRule = append(auctionIdRule, auctionIdItem)
	}
	var nftContractRule []interface{}
	for _, nftContractItem := range nftContract {
		nftContractRule = append(nftContractRule, nftContractItem)
	}
	var tokenIdRule []interface{}
	for _, tokenIdItem := range tokenId {
		tokenIdRule = append(tokenIdRule, tokenIdItem)
	}

	logs, sub, err := _Auction.contract.WatchLogs(opts, "AuctionCreated", auctionIdRule, nftContractRule, tokenIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AuctionAuctionCreated)
				if err := _Auction.contract.UnpackLog(event, "AuctionCreated", log); err != nil {
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

// ParseAuctionCreated is a log parse operation binding the contract event 0x35307ae78277f7a6939b4c091fee2c5a1fd6d1c40ed27ee46bf3338b3f47bc06.
//
// Solidity: event AuctionCreated(uint256 indexed auctionId, address indexed nftContract, uint256 indexed tokenId, address seller, uint256 startingPrice, uint256 endTime, address paymentToken)
func (_Auction *AuctionFilterer) ParseAuctionCreated(log types.Log) (*AuctionAuctionCreated, error) {
	event := new(AuctionAuctionCreated)
	if err := _Auction.contract.UnpackLog(event, "AuctionCreated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AuctionAuctionEndedIterator is returned from FilterAuctionEnded and is used to iterate over the raw logs and unpacked data for AuctionEnded events raised by the Auction contract.
type AuctionAuctionEndedIterator struct {
	Event *AuctionAuctionEnded // Event containing the contract specifics and raw log

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
func (it *AuctionAuctionEndedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AuctionAuctionEnded)
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
		it.Event = new(AuctionAuctionEnded)
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
func (it *AuctionAuctionEndedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AuctionAuctionEndedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AuctionAuctionEnded represents a AuctionEnded event raised by the Auction contract.
type AuctionAuctionEnded struct {
	AuctionId     *big.Int
	Winner        common.Address
	FinalPriceUSD *big.Int
	PaymentToken  common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterAuctionEnded is a free log retrieval operation binding the contract event 0xbb6764412c29916bdf4a5c6fe6b1c079de35682160b2289928ce003ab459a749.
//
// Solidity: event AuctionEnded(uint256 indexed auctionId, address indexed winner, uint256 finalPriceUSD, address paymentToken)
func (_Auction *AuctionFilterer) FilterAuctionEnded(opts *bind.FilterOpts, auctionId []*big.Int, winner []common.Address) (*AuctionAuctionEndedIterator, error) {

	var auctionIdRule []interface{}
	for _, auctionIdItem := range auctionId {
		auctionIdRule = append(auctionIdRule, auctionIdItem)
	}
	var winnerRule []interface{}
	for _, winnerItem := range winner {
		winnerRule = append(winnerRule, winnerItem)
	}

	logs, sub, err := _Auction.contract.FilterLogs(opts, "AuctionEnded", auctionIdRule, winnerRule)
	if err != nil {
		return nil, err
	}
	return &AuctionAuctionEndedIterator{contract: _Auction.contract, event: "AuctionEnded", logs: logs, sub: sub}, nil
}

// WatchAuctionEnded is a free log subscription operation binding the contract event 0xbb6764412c29916bdf4a5c6fe6b1c079de35682160b2289928ce003ab459a749.
//
// Solidity: event AuctionEnded(uint256 indexed auctionId, address indexed winner, uint256 finalPriceUSD, address paymentToken)
func (_Auction *AuctionFilterer) WatchAuctionEnded(opts *bind.WatchOpts, sink chan<- *AuctionAuctionEnded, auctionId []*big.Int, winner []common.Address) (event.Subscription, error) {

	var auctionIdRule []interface{}
	for _, auctionIdItem := range auctionId {
		auctionIdRule = append(auctionIdRule, auctionIdItem)
	}
	var winnerRule []interface{}
	for _, winnerItem := range winner {
		winnerRule = append(winnerRule, winnerItem)
	}

	logs, sub, err := _Auction.contract.WatchLogs(opts, "AuctionEnded", auctionIdRule, winnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AuctionAuctionEnded)
				if err := _Auction.contract.UnpackLog(event, "AuctionEnded", log); err != nil {
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

// ParseAuctionEnded is a log parse operation binding the contract event 0xbb6764412c29916bdf4a5c6fe6b1c079de35682160b2289928ce003ab459a749.
//
// Solidity: event AuctionEnded(uint256 indexed auctionId, address indexed winner, uint256 finalPriceUSD, address paymentToken)
func (_Auction *AuctionFilterer) ParseAuctionEnded(log types.Log) (*AuctionAuctionEnded, error) {
	event := new(AuctionAuctionEnded)
	if err := _Auction.contract.UnpackLog(event, "AuctionEnded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AuctionBidPlacedIterator is returned from FilterBidPlaced and is used to iterate over the raw logs and unpacked data for BidPlaced events raised by the Auction contract.
type AuctionBidPlacedIterator struct {
	Event *AuctionBidPlaced // Event containing the contract specifics and raw log

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
func (it *AuctionBidPlacedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AuctionBidPlaced)
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
		it.Event = new(AuctionBidPlaced)
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
func (it *AuctionBidPlacedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AuctionBidPlacedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AuctionBidPlaced represents a BidPlaced event raised by the Auction contract.
type AuctionBidPlaced struct {
	AuctionId      *big.Int
	Bidder         common.Address
	BidAmountUSD   *big.Int
	OriginalAmount *big.Int
	PaymentToken   common.Address
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterBidPlaced is a free log retrieval operation binding the contract event 0x0362a784a7d4c5ca36537108d384ed276fed3ee71408b5c758defb3c2c0f85f9.
//
// Solidity: event BidPlaced(uint256 indexed auctionId, address indexed bidder, uint256 bidAmountUSD, uint256 originalAmount, address paymentToken)
func (_Auction *AuctionFilterer) FilterBidPlaced(opts *bind.FilterOpts, auctionId []*big.Int, bidder []common.Address) (*AuctionBidPlacedIterator, error) {

	var auctionIdRule []interface{}
	for _, auctionIdItem := range auctionId {
		auctionIdRule = append(auctionIdRule, auctionIdItem)
	}
	var bidderRule []interface{}
	for _, bidderItem := range bidder {
		bidderRule = append(bidderRule, bidderItem)
	}

	logs, sub, err := _Auction.contract.FilterLogs(opts, "BidPlaced", auctionIdRule, bidderRule)
	if err != nil {
		return nil, err
	}
	return &AuctionBidPlacedIterator{contract: _Auction.contract, event: "BidPlaced", logs: logs, sub: sub}, nil
}

// WatchBidPlaced is a free log subscription operation binding the contract event 0x0362a784a7d4c5ca36537108d384ed276fed3ee71408b5c758defb3c2c0f85f9.
//
// Solidity: event BidPlaced(uint256 indexed auctionId, address indexed bidder, uint256 bidAmountUSD, uint256 originalAmount, address paymentToken)
func (_Auction *AuctionFilterer) WatchBidPlaced(opts *bind.WatchOpts, sink chan<- *AuctionBidPlaced, auctionId []*big.Int, bidder []common.Address) (event.Subscription, error) {

	var auctionIdRule []interface{}
	for _, auctionIdItem := range auctionId {
		auctionIdRule = append(auctionIdRule, auctionIdItem)
	}
	var bidderRule []interface{}
	for _, bidderItem := range bidder {
		bidderRule = append(bidderRule, bidderItem)
	}

	logs, sub, err := _Auction.contract.WatchLogs(opts, "BidPlaced", auctionIdRule, bidderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AuctionBidPlaced)
				if err := _Auction.contract.UnpackLog(event, "BidPlaced", log); err != nil {
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

// ParseBidPlaced is a log parse operation binding the contract event 0x0362a784a7d4c5ca36537108d384ed276fed3ee71408b5c758defb3c2c0f85f9.
//
// Solidity: event BidPlaced(uint256 indexed auctionId, address indexed bidder, uint256 bidAmountUSD, uint256 originalAmount, address paymentToken)
func (_Auction *AuctionFilterer) ParseBidPlaced(log types.Log) (*AuctionBidPlaced, error) {
	event := new(AuctionBidPlaced)
	if err := _Auction.contract.UnpackLog(event, "BidPlaced", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AuctionBidWithdrawnIterator is returned from FilterBidWithdrawn and is used to iterate over the raw logs and unpacked data for BidWithdrawn events raised by the Auction contract.
type AuctionBidWithdrawnIterator struct {
	Event *AuctionBidWithdrawn // Event containing the contract specifics and raw log

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
func (it *AuctionBidWithdrawnIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AuctionBidWithdrawn)
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
		it.Event = new(AuctionBidWithdrawn)
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
func (it *AuctionBidWithdrawnIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AuctionBidWithdrawnIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AuctionBidWithdrawn represents a BidWithdrawn event raised by the Auction contract.
type AuctionBidWithdrawn struct {
	AuctionId *big.Int
	Bidder    common.Address
	Amount    *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterBidWithdrawn is a free log retrieval operation binding the contract event 0x8f8619524e8d462cead34604bd2247ede24175801481e4d0b8059ac8aa41c301.
//
// Solidity: event BidWithdrawn(uint256 indexed auctionId, address indexed bidder, uint256 amount)
func (_Auction *AuctionFilterer) FilterBidWithdrawn(opts *bind.FilterOpts, auctionId []*big.Int, bidder []common.Address) (*AuctionBidWithdrawnIterator, error) {

	var auctionIdRule []interface{}
	for _, auctionIdItem := range auctionId {
		auctionIdRule = append(auctionIdRule, auctionIdItem)
	}
	var bidderRule []interface{}
	for _, bidderItem := range bidder {
		bidderRule = append(bidderRule, bidderItem)
	}

	logs, sub, err := _Auction.contract.FilterLogs(opts, "BidWithdrawn", auctionIdRule, bidderRule)
	if err != nil {
		return nil, err
	}
	return &AuctionBidWithdrawnIterator{contract: _Auction.contract, event: "BidWithdrawn", logs: logs, sub: sub}, nil
}

// WatchBidWithdrawn is a free log subscription operation binding the contract event 0x8f8619524e8d462cead34604bd2247ede24175801481e4d0b8059ac8aa41c301.
//
// Solidity: event BidWithdrawn(uint256 indexed auctionId, address indexed bidder, uint256 amount)
func (_Auction *AuctionFilterer) WatchBidWithdrawn(opts *bind.WatchOpts, sink chan<- *AuctionBidWithdrawn, auctionId []*big.Int, bidder []common.Address) (event.Subscription, error) {

	var auctionIdRule []interface{}
	for _, auctionIdItem := range auctionId {
		auctionIdRule = append(auctionIdRule, auctionIdItem)
	}
	var bidderRule []interface{}
	for _, bidderItem := range bidder {
		bidderRule = append(bidderRule, bidderItem)
	}

	logs, sub, err := _Auction.contract.WatchLogs(opts, "BidWithdrawn", auctionIdRule, bidderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AuctionBidWithdrawn)
				if err := _Auction.contract.UnpackLog(event, "BidWithdrawn", log); err != nil {
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

// ParseBidWithdrawn is a log parse operation binding the contract event 0x8f8619524e8d462cead34604bd2247ede24175801481e4d0b8059ac8aa41c301.
//
// Solidity: event BidWithdrawn(uint256 indexed auctionId, address indexed bidder, uint256 amount)
func (_Auction *AuctionFilterer) ParseBidWithdrawn(log types.Log) (*AuctionBidWithdrawn, error) {
	event := new(AuctionBidWithdrawn)
	if err := _Auction.contract.UnpackLog(event, "BidWithdrawn", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AuctionInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the Auction contract.
type AuctionInitializedIterator struct {
	Event *AuctionInitialized // Event containing the contract specifics and raw log

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
func (it *AuctionInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AuctionInitialized)
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
		it.Event = new(AuctionInitialized)
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
func (it *AuctionInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AuctionInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AuctionInitialized represents a Initialized event raised by the Auction contract.
type AuctionInitialized struct {
	Version uint64
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_Auction *AuctionFilterer) FilterInitialized(opts *bind.FilterOpts) (*AuctionInitializedIterator, error) {

	logs, sub, err := _Auction.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &AuctionInitializedIterator{contract: _Auction.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_Auction *AuctionFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *AuctionInitialized) (event.Subscription, error) {

	logs, sub, err := _Auction.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AuctionInitialized)
				if err := _Auction.contract.UnpackLog(event, "Initialized", log); err != nil {
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

// ParseInitialized is a log parse operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_Auction *AuctionFilterer) ParseInitialized(log types.Log) (*AuctionInitialized, error) {
	event := new(AuctionInitialized)
	if err := _Auction.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AuctionOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the Auction contract.
type AuctionOwnershipTransferredIterator struct {
	Event *AuctionOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *AuctionOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AuctionOwnershipTransferred)
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
		it.Event = new(AuctionOwnershipTransferred)
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
func (it *AuctionOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AuctionOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AuctionOwnershipTransferred represents a OwnershipTransferred event raised by the Auction contract.
type AuctionOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Auction *AuctionFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*AuctionOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Auction.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &AuctionOwnershipTransferredIterator{contract: _Auction.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Auction *AuctionFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *AuctionOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Auction.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AuctionOwnershipTransferred)
				if err := _Auction.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_Auction *AuctionFilterer) ParseOwnershipTransferred(log types.Log) (*AuctionOwnershipTransferred, error) {
	event := new(AuctionOwnershipTransferred)
	if err := _Auction.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
