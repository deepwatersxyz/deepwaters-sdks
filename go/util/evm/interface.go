package evm

import (
	"context"
	"crypto/ecdsa"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
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

type Connector interface {
	GetConfig() *ConnectorConfig
	GetClient() *ethclient.Client
	Open() error
	Call(ctx context.Context, abi *abi.ABI, address common.Address, block *uint64, from *common.Address, f string, args ...interface{}) ([]byte, error)
	Send(args SendArgs) (*common.Hash, *types.Receipt, *uint64, error)
}

type Contract interface {
	GetABI() *abi.ABI
	GetAddressBytes() common.Address
	GetAddressStr(checkSum bool) string
	GetName() string
	GetConnector() Connector
	GetDeployedBlock() *uint64
	Call(ctx context.Context, block *uint64, from *string, f string, args ...interface{}) ([]byte, error)
	Send(mode SendMode, from Wallet, to *string, value *big.Int, nonce *uint64, f *string, args ...interface{}) (*string, *types.Receipt, *uint64, error)
}

var ZeroAddress common.Address = common.HexToAddress("0x00000000000000000000000000000000")
