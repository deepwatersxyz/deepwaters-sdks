package evm

import (
	"context"
	"crypto/ecdsa"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	log "github.com/sirupsen/logrus"
)

type SendMode string

const (
	SendModeNoWait         SendMode = "noWait"
	SendModeWaitForReceipt SendMode = "waitForReceipt"
)

type NetworkType string

const (
	EVMNetworkType NetworkType = "evm"
)

type Wallet interface {
	GetPrivateKey() *ecdsa.PrivateKey
	GetAddressBytes() common.Address
	GetAddressStr(checksum bool) string
	GetECDSASignature(hash common.Hash) (*string, error)
}

type SendArgs struct {
	Mode     SendMode
	Wallet   Wallet
	ABI      *abi.ABI
	Address  common.Address
	Value    *big.Int
	Nonce    *uint64
	GasPrice *big.Int
	GasLimit *uint64
	F        *string
	Args     []interface{}
}

type TransactionSender interface {
	GetEVMConnector() EVMConnector
	Send(args SendArgs) (*common.Hash, *types.Receipt, *uint64, error)
}

type EVMConnector interface {
	GetNetworkType() NetworkType
	GetConfig() *EVMConnectorConfig
	GetClient() ETHClientWrapper
	GetTransactionSender() TransactionSender
	GetContracts() Contracts

	Open() error
	OpenWait(timeout time.Duration) error
	Ping() error

	GetLastBlockNum(ctx context.Context, finality uint64) (*uint64, error)
	GetBlockHeader(ctx context.Context, block uint64) (*types.Header, error)
	GetTransactionReceipt(ctx context.Context, hash common.Hash) (*types.Receipt, error)
	GetTransactionByHash(ctx context.Context, hash common.Hash) (*types.Transaction, bool, error)

	Call(ctx context.Context, abi *abi.ABI, address common.Address, block *uint64, from *common.Address, f string, args ...interface{}) ([]byte, error)
}

type Contract interface {
	GetChainID() int
	GetChainName() string
	GetABI() *abi.ABI
	GetAddressBytes() common.Address
	GetAddressStr(checkSum bool) string
	GetName() string
	GetDescription() *string
	GetEVM() EVMConnector
	GetDeployedBlock() *uint64
	Call(ctx context.Context, block *uint64, from *string, f string, args ...interface{}) ([]byte, error)
	Send(mode SendMode, from Wallet, to *string, value *big.Int, nonce *uint64, f *string, args ...interface{}) (*string, *types.Receipt, *uint64, error)
}

type Contracts interface {
	AddContract(contract Contract) error
	LoadFromConfigs(contractConfigs []ContractConfig, evm EVMConnector, lg *log.Logger) error
	HasContract(chainName string, address string) bool
	GetContract(chainName string, address string) Contract
	ForEach(f func(contract Contract) error) error
	Clone() Contracts
}

var ZeroAddress common.Address = common.HexToAddress("0x00000000000000000000000000000000")

type ETHClientWrapper interface {
	Close()
	ChainID(ctx context.Context) (*big.Int, error)
	BlockNumber(ctx context.Context) (uint64, error)
	HeaderByNumber(ctx context.Context, number *big.Int) (*types.Header, error)
	TransactionByHash(ctx context.Context, hash common.Hash) (tx *types.Transaction, isPending bool, err error)
	TransactionReceipt(ctx context.Context, txHash common.Hash) (*types.Receipt, error)
	SyncProgress(ctx context.Context) (*ethereum.SyncProgress, error)
	FilterLogs(ctx context.Context, q ethereum.FilterQuery) ([]types.Log, error)
	PendingBalanceAt(ctx context.Context, account common.Address) (*big.Int, error)
	PendingNonceAt(ctx context.Context, account common.Address) (uint64, error)
	CallContract(ctx context.Context, msg ethereum.CallMsg, blockNumber *big.Int) ([]byte, error)
	SuggestGasPrice(ctx context.Context) (*big.Int, error)
	EstimateGas(ctx context.Context, msg ethereum.CallMsg) (uint64, error)
	SendTransaction(ctx context.Context, tx *types.Transaction) error
}
