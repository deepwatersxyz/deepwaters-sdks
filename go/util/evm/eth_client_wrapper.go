package evm

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

type ethClientWrapper struct {
	c *ethclient.Client
}

func NewETHClientWrapper(url string) (ETHClientWrapper, error) {
	c, err := ethclient.Dial(url)
	return &ethClientWrapper{c: c}, err
}

// many other eth client methods could be exposed
func (c *ethClientWrapper) Close()                                        { c.c.Close() }
func (c *ethClientWrapper) ChainID(ctx context.Context) (*big.Int, error) { return c.c.ChainID(ctx) }

func (c *ethClientWrapper) BlockNumber(ctx context.Context) (uint64, error) {
	return c.c.BlockNumber(ctx)
}

func (c *ethClientWrapper) HeaderByNumber(ctx context.Context, number *big.Int) (*types.Header, error) {
	return c.c.HeaderByNumber(ctx, number)
}
func (c *ethClientWrapper) TransactionByHash(ctx context.Context, hash common.Hash) (tx *types.Transaction, isPending bool, err error) {
	return c.c.TransactionByHash(ctx, hash)
}

func (c *ethClientWrapper) TransactionReceipt(ctx context.Context, txHash common.Hash) (*types.Receipt, error) {
	return c.c.TransactionReceipt(ctx, txHash)
}
func (c *ethClientWrapper) SyncProgress(ctx context.Context) (*ethereum.SyncProgress, error) {
	return c.c.SyncProgress(ctx)
}

func (c *ethClientWrapper) FilterLogs(ctx context.Context, q ethereum.FilterQuery) ([]types.Log, error) {
	return c.c.FilterLogs(ctx, q)
}

func (c *ethClientWrapper) PendingBalanceAt(ctx context.Context, account common.Address) (*big.Int, error) {
	return c.c.PendingBalanceAt(ctx, account)
}

func (c *ethClientWrapper) PendingNonceAt(ctx context.Context, account common.Address) (uint64, error) {
	return c.c.PendingNonceAt(ctx, account)
}

func (c *ethClientWrapper) CallContract(ctx context.Context, msg ethereum.CallMsg, blockNumber *big.Int) ([]byte, error) {
	return c.c.CallContract(ctx, msg, blockNumber)
}

func (c *ethClientWrapper) SuggestGasPrice(ctx context.Context) (*big.Int, error) {
	return c.c.SuggestGasPrice(ctx)
}

func (c *ethClientWrapper) EstimateGas(ctx context.Context, msg ethereum.CallMsg) (uint64, error) {
	return c.c.EstimateGas(ctx, msg)
}
func (c *ethClientWrapper) SendTransaction(ctx context.Context, tx *types.Transaction) error {
	return c.c.SendTransaction(ctx, tx)
}
