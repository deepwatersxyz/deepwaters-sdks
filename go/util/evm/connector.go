package evm

import (
	"context"
	"fmt"
	"math/big"
	"sync/atomic"
	"time"

	"github.com/ethereum/go-ethereum"
	geth "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	log "github.com/sirupsen/logrus"

	"deepwaters/go-examples/util"
)

// one connector per chain
type connector struct {
	lg              *log.Entry
	cfg             *ConnectorConfig
	client          *ethclient.Client
	pendingReceipts int32
}

func NewConnector(lg *log.Logger, cfg *ConnectorConfig) *connector {
	return &connector{
		lg:  lg.WithField("package", "util.evm"),
		cfg: cfg,
	}
}

func (c *connector) GetConfig() *ConnectorConfig {
	return c.cfg
}

func (c *connector) GetClient() *ethclient.Client {
	return c.client
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
	c.client, err = ethclient.Dial(c.cfg.URL)
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

func (c *connector) Send(args SendArgs) (*common.Hash, *types.Receipt, *uint64, error) {
	start := time.Now()
	var err error

	if args.Wallet == nil {
		return nil, nil, nil, fmt.Errorf("wallet not defined")
	}

	// nonce
	if args.Nonce == nil {
		args.Nonce, err = c.getPendingNonce(args.Wallet.GetAddressBytes())
		if err != nil {
			return nil, nil, nil, err
		}
	}

	// gas price
	if args.GasPrice == nil {
		args.GasPrice, err = c.getCappedEstimatedGasPrice()
		if err != nil {
			return nil, nil, nil, err
		}
	}

	// data
	data, err := c.packData(args)
	if err != nil {
		return nil, nil, nil, err
	}

	// gas
	if args.GasLimit == nil {
		args.GasLimit, err = c.getAdjustedEstimatedGas(args, data)
		if err != nil {
			panic(fmt.Errorf("manual gas limit not implemented: %+v", err.Error()))
		}
	}

	// transaction
	chainID := big.NewInt(int64(c.cfg.ChainID))
	tx := types.NewTransaction(*args.Nonce, args.Address, args.Value, *args.GasLimit, args.GasPrice, data)

	// sign transaction
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), args.Wallet.GetPrivateKey())
	if err != nil {
		return nil, nil, nil, err
	}

	err = c.client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		c.lg.WithFields(c.getSendArgsLogFields(args, log.Fields{
			"error": err,
		})).Error("unable to send raw transaction")
		return nil, nil, nil, err
	}

	atomic.AddInt32(&c.pendingReceipts, 1)
	hash := signedTx.Hash()

	c.lg.WithFields(c.getSendArgsLogFields(args, log.Fields{
		"after":  util.SinceMillis(start),
		"txHash": hash.Hex(),
	})).Info("sent raw transaction")

	var receipt *types.Receipt
	switch args.Mode {
	case SendModeNoWait:
		go func() {
			_ = c.waitForReceipt(hash, args)
		}()
	case SendModeWaitForReceipt:
		receipt = c.waitForReceipt(hash, args)
	default:
		panic(fmt.Sprintf("unsupported send mode: %+v", args.Mode))
	}

	nonce := args.Nonce
	*nonce += 1
	return &hash, receipt, nonce, nil
}

func (c *connector) packData(args SendArgs) ([]byte, error) {
	if args.ABI == nil {
		return nil, fmt.Errorf("missing abi in SendArgs")
	}
	data, err := args.ABI.Pack(*args.F, args.Args...)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (c *connector) getCappedEstimatedGasPrice() (*big.Int, error) {
	start := time.Now()
	gasPrice, err := c.getEstimatedGasPrice()
	if err != nil {
		return nil, err
	}
	limit := c.cfg.MaxGasPriceWei
	if gasPrice.Cmp(limit) == 1 {
		c.lg.WithFields(log.Fields{
			"after":     util.SinceMillis(start),
			"estimated": gasPrice,
			"gasPrice":  limit,
			"network":   c.cfg.ChainName,
		}).Trace("estimated and capped gas price")
		return limit, nil
	}
	c.lg.WithFields(log.Fields{
		"after":    util.SinceMillis(start),
		"gasPrice": gasPrice,
		"limit":    limit,
		"network":  c.cfg.ChainName,
	}).Debug("estimated gas price")
	return gasPrice, nil
}

func (c *connector) getEstimatedGasPrice() (*big.Int, error) {
	gasPrice, err := c.client.SuggestGasPrice(context.Background())
	if err != nil {
		return nil, fmt.Errorf("unable to estimate gas price: %+v", err)
	}
	return gasPrice, nil
}

func (c *connector) getAdjustedEstimatedGas(args SendArgs, data []byte) (*uint64, error) {
	start := time.Now()
	gas, err := c.getEstimatedGas(args, data)
	if err != nil {
		return nil, err
	}
	if addPercent := c.cfg.AddEstimatedGasPercent; addPercent > 0 {
		adjusted := (*gas * (100 + uint64(addPercent))) / 100
		c.lg.WithFields(log.Fields{
			"after":    util.SinceMillis(start),
			"contract": args.Address.Hex(),
			"func":     util.Str(args.F),
			"network":  c.cfg.ChainName,
			"from":     *gas,
			"to":       adjusted,
		}).Debug("estimated and adjusted gas limit")
		return &adjusted, nil
	}
	c.lg.WithFields(log.Fields{
		"after":    util.SinceMillis(start),
		"contract": args.Address.Hex(),
		"func":     util.Str(args.F),
		"network":  c.cfg.ChainName,
		"gasLimit": *gas,
	}).Debug("estimated gas limit")
	return gas, nil
}

func (c *connector) getEstimatedGas(args SendArgs, data []byte) (*uint64, error) {
	msg := ethereum.CallMsg{
		From:      args.Wallet.GetAddressBytes(),
		To:        &args.Address,
		Gas:       0,
		GasPrice:  args.GasPrice,
		GasFeeCap: args.GasPrice,
		Value:     args.Value,
		Data:      data,
	}
	gas, err := c.client.EstimateGas(context.Background(), msg)
	if err != nil {
		c.lg.WithFields(c.getSendArgsLogFields(args, log.Fields{
			"error":   err,
			"network": c.cfg.ChainName,
		})).Warning("unable to estimate gas")
		return nil, err
	}
	return &gas, nil
}

func (c *connector) getPendingNonce(address common.Address) (*uint64, error) {
	start := time.Now()
	nonce, err := c.client.PendingNonceAt(context.Background(), address)
	if err != nil {
		c.lg.WithFields(log.Fields{
			"account": address.Hex(),
			"error":   err,
			"network": c.cfg.ChainName,
		}).Error("unable to fetch pending nonce")
		return nil, err
	}
	c.lg.WithFields(log.Fields{
		"after":   util.SinceMillis(start),
		"account": address.Hex(),
		"nonce":   nonce,
		"network": c.cfg.ChainName,
	}).Trace("fetched pending nonce")
	return &nonce, nil
}

func (c *connector) waitForReceipt(hash common.Hash, args SendArgs) *types.Receipt {
	start := time.Now()
	for {
		receipt, err := c.client.TransactionReceipt(context.Background(), hash)
		if err == nil {
			fields := log.Fields{
				"after":     util.SinceMillis(start),
				"txHash":    hash.Hex(),
				"txIndex":   receipt.TransactionIndex,
				"status":    receipt.Status,
				"gasLimit":  *args.GasLimit,
				"gasUsed":   receipt.GasUsed,
				"block":     receipt.BlockNumber.Uint64(),
				"blockHash": receipt.BlockHash.Hex(),
			}
			if receipt.Status == 1 {
				c.lg.WithFields(fields).Info("received transaction receipt")
			} else {
				c.lg.WithFields(fields).Error("received transaction receipt with error")
			}

			atomic.AddInt32(&c.pendingReceipts, -1)
			return receipt
		}

		c.lg.WithFields(log.Fields{
			"after":   util.SinceSecs(start),
			"txHash":  hash.Hex(),
			"pending": atomic.LoadInt32(&c.pendingReceipts),
		}).Debug("waiting for transaction receipt...")

		time.Sleep(30 * time.Second)
	}
}

func (c *connector) getSendArgsLogFields(args SendArgs, addFields log.Fields) log.Fields {
	fields := log.Fields{
		"mode":      args.Mode,
		"recipient": args.Address.Hex(),
		"value":     args.Value,
		"nonce":     *args.Nonce,
		"gasPrice":  args.GasPrice,
		"gasLimit":  util.Uint64Str(args.GasLimit),
		"network":   c.cfg.ChainName,
	}
	if args.Wallet != nil {
		fields["wallet"] = args.Wallet.GetAddressStr(true)
	}
	if args.F != nil {
		fields["func"] = *args.F
		fields["funcArgs"] = args.Args
	}
	for k, v := range addFields {
		fields[k] = v
	}
	return fields
}
