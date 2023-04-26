package evm

import (
	"context"
	"fmt"
	"math/big"
	"sync/atomic"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	log "github.com/sirupsen/logrus"

	"deepwaters/go-examples/util"
)

type transactionSender struct {
	lg              *log.Entry
	evm             EVMConnector
	pendingReceipts int32
}

func NewTransactionSender(lg *log.Logger, evm EVMConnector) TransactionSender {
	return &transactionSender{
		lg:  lg.WithField("package", "util.evm"),
		evm: evm,
	}
}

func (s *transactionSender) GetEVMConnector() EVMConnector {
	return s.evm
}

func (s *transactionSender) Send(args SendArgs) (*common.Hash, *types.Receipt, *uint64, error) {
	start := time.Now()
	var err error

	if args.Wallet == nil {
		return nil, nil, nil, fmt.Errorf("wallet not defined")
	}

	// nonce
	if args.Nonce == nil {
		args.Nonce, err = s.getPendingNonce(args.Wallet.GetAddressBytes())
		if err != nil {
			return nil, nil, nil, err
		}
	}

	// gas price
	if args.GasPrice == nil {
		args.GasPrice, err = s.getCappedEstimatedGasPrice()
		if err != nil {
			return nil, nil, nil, err
		}
	}

	// data
	data, err := s.packData(args)
	if err != nil {
		return nil, nil, nil, err
	}

	// gas
	if args.GasLimit == nil {
		args.GasLimit, err = s.getAdjustedEstimatedGas(args, data)
		if err != nil {
			panic(fmt.Errorf("manual gas limit not implemented: %+v", err.Error()))
		}
	}

	// transaction
	chainID := big.NewInt(int64(s.evm.GetConfig().ChainID))
	tx := types.NewTransaction(*args.Nonce, args.Address, args.Value, *args.GasLimit, args.GasPrice, data)

	// sign transaction
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), args.Wallet.GetPrivateKey())
	if err != nil {
		return nil, nil, nil, err
	}

	err = s.evm.GetClient().SendTransaction(context.Background(), signedTx)
	if err != nil {
		s.lg.WithFields(s.getSendArgsLogFields(args, log.Fields{
			"error": err,
		})).Error("unable to send raw transaction")
		return nil, nil, nil, err
	}

	atomic.AddInt32(&s.pendingReceipts, 1)
	hash := signedTx.Hash()

	s.lg.WithFields(s.getSendArgsLogFields(args, log.Fields{
		"after":  util.SinceMillis(start),
		"txHash": hash.Hex(),
	})).Info("sent raw transaction")

	var receipt *types.Receipt
	switch args.Mode {
	case SendModeNoWait:
		go func() {
			_ = s.waitForReceipt(hash, args)
		}()
	case SendModeWaitForReceipt:
		receipt = s.waitForReceipt(hash, args)
	default:
		panic(fmt.Sprintf("unsupported send mode: %+v", args.Mode))
	}

	nonce := args.Nonce
	*nonce += 1
	return &hash, receipt, nonce, nil
}

func (s *transactionSender) packData(args SendArgs) ([]byte, error) {
	if args.ABI == nil {
		return nil, fmt.Errorf("missing abi in SendArgs")
	}
	data, err := args.ABI.Pack(*args.F, args.Args...)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (s *transactionSender) getCappedEstimatedGasPrice() (*big.Int, error) {
	start := time.Now()
	gasPrice, err := s.getEstimatedGasPrice()
	if err != nil {
		return nil, err
	}
	limit := s.evm.GetConfig().Sender.MaxGasPriceWei
	if gasPrice.Cmp(limit) == 1 {
		s.lg.WithFields(log.Fields{
			"after":     util.SinceMillis(start),
			"estimated": gasPrice,
			"gasPrice":  limit,
			"network":   s.evm.GetConfig().ChainName,
		}).Trace("estimated and capped gas price")
		return limit, nil
	}
	s.lg.WithFields(log.Fields{
		"after":    util.SinceMillis(start),
		"gasPrice": gasPrice,
		"limit":    limit,
		"network":  s.evm.GetConfig().ChainName,
	}).Debug("estimated gas price")
	return gasPrice, nil
}

func (s *transactionSender) getEstimatedGasPrice() (*big.Int, error) {
	gasPrice, err := s.evm.GetClient().SuggestGasPrice(context.Background())
	if err != nil {
		return nil, fmt.Errorf("unable to estimate gas price: %+v", err)
	}
	return gasPrice, nil
}

func (s *transactionSender) getAdjustedEstimatedGas(args SendArgs, data []byte) (*uint64, error) {
	start := time.Now()
	gas, err := s.getEstimatedGas(args, data)
	if err != nil {
		return nil, err
	}
	if addPercent := s.evm.GetConfig().Sender.AddEstimatedGasPercent; addPercent > 0 {
		adjusted := (*gas * (100 + uint64(addPercent))) / 100
		s.lg.WithFields(log.Fields{
			"after":    util.SinceMillis(start),
			"contract": args.Address.Hex(),
			"func":     util.Str(args.F),
			"network":  s.evm.GetConfig().ChainName,
			"from":     util.PrettyNum(*gas),
			"to":       util.PrettyNum(adjusted),
		}).Debug("estimated and adjusted gas limit")
		return &adjusted, nil
	}
	s.lg.WithFields(log.Fields{
		"after":    util.SinceMillis(start),
		"contract": args.Address.Hex(),
		"func":     util.Str(args.F),
		"network":  s.evm.GetConfig().ChainName,
		"gasLimit": util.PrettyNum(*gas),
	}).Debug("estimated gas limit")
	return gas, nil
}

func (s *transactionSender) getEstimatedGas(args SendArgs, data []byte) (*uint64, error) {
	msg := ethereum.CallMsg{
		From:      args.Wallet.GetAddressBytes(),
		To:        &args.Address,
		Gas:       0,
		GasPrice:  args.GasPrice,
		GasFeeCap: args.GasPrice,
		Value:     args.Value,
		Data:      data,
	}
	gas, err := s.evm.GetClient().EstimateGas(context.Background(), msg)
	if err != nil {
		s.lg.WithFields(s.getSendArgsLogFields(args, log.Fields{
			"error":   err,
			"network": s.evm.GetConfig().ChainName,
		})).Warning("unable to estimate gas")
		return nil, err
	}
	return &gas, nil
}

func (s *transactionSender) getPendingNonce(address common.Address) (*uint64, error) {
	start := time.Now()
	nonce, err := s.evm.GetClient().PendingNonceAt(context.Background(), address)
	if err != nil {
		s.lg.WithFields(log.Fields{
			"account": address.Hex(),
			"error":   err,
			"network": s.evm.GetConfig().ChainName,
		}).Error("unable to fetch pending nonce")
		return nil, err
	}
	s.lg.WithFields(log.Fields{
		"after":   util.SinceMillis(start),
		"account": address.Hex(),
		"nonce":   nonce,
		"network": s.evm.GetConfig().ChainName,
	}).Trace("fetched pending nonce")
	return &nonce, nil
}

func (s *transactionSender) waitForReceipt(hash common.Hash, args SendArgs) *types.Receipt {
	start := time.Now()
	for {
		receipt, err := s.evm.GetClient().TransactionReceipt(context.Background(), hash)
		if err == nil {
			fields := log.Fields{
				"after":     util.SinceMillis(start),
				"txHash":    hash.Hex(),
				"txIndex":   receipt.TransactionIndex,
				"status":    receipt.Status,
				"gasLimit":  *args.GasLimit,
				"gasUsed":   receipt.GasUsed,
				"block":     util.PrettyNum(receipt.BlockNumber.Uint64()),
				"blockHash": receipt.BlockHash.Hex(),
			}
			if receipt.Status == 1 {
				s.lg.WithFields(fields).Info("received transaction receipt")
			} else {
				s.lg.WithFields(fields).Error("received transaction receipt with error")
			}

			atomic.AddInt32(&s.pendingReceipts, -1)
			return receipt
		}

		s.lg.WithFields(log.Fields{
			"after":   util.SinceSecs(start),
			"txHash":  hash.Hex(),
			"pending": atomic.LoadInt32(&s.pendingReceipts),
		}).Debug("waiting for transaction receipt...")

		time.Sleep(30 * time.Second)
	}
}

func (s *transactionSender) getSendArgsLogFields(args SendArgs, addFields log.Fields) log.Fields {
	fields := log.Fields{
		"mode":      args.Mode,
		"recipient": args.Address.Hex(),
		"value":     args.Value,
		"nonce":     util.PrettyNum(*args.Nonce),
		"gasPrice":  args.GasPrice,
		"gasLimit":  util.Uint64Str(args.GasLimit, true),
		"network":   s.evm.GetConfig().ChainName,
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
