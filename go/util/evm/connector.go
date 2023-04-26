package evm

import (
	"context"
	"fmt"
	"time"

	geth "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	log "github.com/sirupsen/logrus"

	"deepwaters/go-examples/util"
)

// one connector per chain
type connector struct {
	lg        *log.Entry
	cfg       *EVMConnectorConfig
	client    ETHClientWrapper
	sender    TransactionSender
	contracts Contracts
	dial      func(url string) (ETHClientWrapper, error)
}

func NewEVMConnector(lg *log.Logger, cfg *EVMConnectorConfig) (*connector, error) {
	c := &connector{
		lg:        lg.WithField("package", "micro.web3.evm"),
		cfg:       cfg,
		contracts: NewContracts(),
		dial:      NewETHClientWrapper,
	}
	c.sender = NewTransactionSender(lg, c)
	if cfg.ContractConfigs != nil {
		if err := c.contracts.LoadFromConfigs(cfg.ContractConfigs, c, lg); err != nil {
			return nil, err
		}
	}
	return c, nil
}

func (c *connector) GetNetworkType() NetworkType {
	return EVMNetworkType
}

func (c *connector) GetConfig() *EVMConnectorConfig {
	return c.cfg
}

func (c *connector) GetClient() ETHClientWrapper {
	return c.client
}

func (c *connector) GetTransactionSender() TransactionSender {
	return c.sender
}

func (c *connector) GetContracts() Contracts {
	return c.contracts
}

func (c *connector) Open() error {
	fields := log.Fields{
		"network": c.cfg.ChainName,
		"url":     c.cfg.URL,
		"chainID": c.cfg.ChainID,
	}
	c.lg.WithFields(fields).Trace("connecting to web3 evm...")

	if err := c.open(); err != nil {
		c.lg.WithFields(fields).Error("unable to connect to web3 evm")
		return err
	}

	c.lg.WithFields(fields).Info("connected to web3 evm")

	return nil
}

func (c *connector) open() error {
	var err error
	c.client, err = c.dial(c.cfg.URL)
	if err != nil {
		return err
	}

	chainID, err := c.client.ChainID(context.Background())
	if err != nil {
		return err
	}

	if int(chainID.Int64()) != c.cfg.ChainID {
		c.lg.WithFields(log.Fields{
			"network":    c.cfg.ChainName,
			"configured": c.cfg.ChainID,
			"detected":   chainID,
		}).Panic("config problem, mismatching chain ids")
	}

	return nil
}

func (c *connector) OpenWait(timeout time.Duration) error {
	fields := log.Fields{
		"network": c.cfg.ChainName,
		"url":     c.cfg.URL,
		"chainID": c.cfg.ChainID,
		"timeout": timeout.Round(time.Millisecond),
	}
	c.lg.WithFields(fields).Trace("connecting to web3...")

	start := time.Now()
	if err := c.openWait(timeout); err != nil {
		c.lg.WithFields(fields).Error("unable to connect to web3 evm")
		return err
	}

	fields["after"] = util.SinceMillis(start)
	c.lg.WithFields(fields).Info("connected to web3 evm")

	return nil
}

func (c *connector) openWait(timeout time.Duration) error {
	start := time.Now()
	p := util.Periodic{
		Timeout:        timeout,
		Interval:       10 * time.Second,
		ReportInterval: 1 * time.Minute,
		ExitOnSuccess:  true,
	}
	return p.Run(c.open, func(attempts int, err error) {
		c.lg.WithFields(log.Fields{
			"after":   util.SinceMillis(start),
			"timeout": timeout,
			"url":     c.cfg.URL,
			"chainID": c.cfg.ChainID,
			"error":   err,
		}).Error("unable to connect to web3 evm, keep trying...")
	})
}

func (c *connector) Ping() error {
	if c.client == nil {
		if err := c.open(); err != nil {
			return err
		}
		return nil
	}
	_, err := c.client.ChainID(context.Background())
	if err != nil {
		return err
	}
	return nil
}

func (c *connector) GetBlockHeader(ctx context.Context, block uint64) (*types.Header, error) {
	return c.client.HeaderByNumber(ctx, util.BigUint64(&block))
}

func (c *connector) GetTransactionByHash(ctx context.Context, hash common.Hash) (*types.Transaction, bool, error) {
	return c.client.TransactionByHash(ctx, hash)
}

func (c *connector) GetTransactionReceipt(ctx context.Context, hash common.Hash) (*types.Receipt, error) {
	return c.client.TransactionReceipt(ctx, hash)
}

func (c *connector) GetLastBlockNum(ctx context.Context, finality uint64) (*uint64, error) {
	block, err := c.client.BlockNumber(ctx)
	if err != nil {
		return nil, err
	}
	if block < finality {
		return nil, fmt.Errorf("insufficient blocks")
	}
	return util.Uint64P(block - finality), nil
}

// passing block requires an archive node
func (c *connector) Call(ctx context.Context, abi *abi.ABI, address common.Address, block *uint64, from *common.Address, f string, args ...interface{}) ([]byte, error) {
	data, err := abi.Pack(f, args...)
	if err != nil {
		return nil, err
	}
	call := geth.CallMsg{
		To:   &address,
		Data: data,
	}
	if from != nil {
		call.From = *from
	}
	return c.client.CallContract(ctx, call, util.BigUint64(block))
}
