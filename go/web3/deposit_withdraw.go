package web3

import (
	"context"
	"deepwaters/go-examples/rest"
	"deepwaters/go-examples/util"
	"deepwaters/go-examples/util/evm"
	"fmt"
	"math/big"

	abiPackage "github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	log "github.com/sirupsen/logrus"
)

type ERC20Details struct {
	Contract evm.Contract
	Decimals uint8
	AssetID  string
}

type MoveTokenInfo struct {
	SendMode       evm.SendMode
	RESTAPIInfo    rest.ConnectionDetails
	Wallet         evm.Wallet
	PosManContract evm.Contract
	ERC20Details   ERC20Details
}

func CheckNativeAndWrappedNativeBalances(lg *log.Logger, restAPIInfo rest.ConnectionDetails, userAddress string, connector evm.EVMConnector, erc20Details ERC20Details) {

	wrappedABI := erc20Details.Contract.GetABI()

	// get on-chain native token balance
	balance, err := connector.GetClient().PendingBalanceAt(context.Background(), common.HexToAddress(userAddress))
	if err != nil {
		lg.Errorf("%s", err)
	}
	balanceStr := evm.AmountToText(balance, connector.GetConfig().NativeCurrency.Decimals, connector.GetConfig().NativeCurrency.Decimals)
	lg.Infof("user has %s %s on-chain", balanceStr, connector.GetConfig().NativeCurrency.Name)

	// get on-chain wrapped native balance
	balanceResult, err := erc20Details.Contract.Call(context.Background(), nil, nil, "balanceOf", common.HexToAddress(userAddress))
	if err != nil {
		lg.Errorf("%+v", err)
		return
	}
	wrappedBalance := abiPackage.ReadInteger(wrappedABI.Methods["balanceOf"].Outputs[0].Type, balanceResult).(*big.Int)
	wrappedBalanceStr := evm.AmountToText(wrappedBalance, erc20Details.Decimals, erc20Details.Decimals)
	lg.Infof("user has %s %s on-chain", wrappedBalanceStr, erc20Details.Contract.GetName())

	// get wrapped native balances in deepwaters
	userInfo, err := rest.GetUserInfo(restAPIInfo)
	if err != nil {
		lg.Errorf("%+v", err)
		return
	}

	for _, balance := range userInfo.Result.Balances {
		if balance.AssetID == erc20Details.AssetID {
			lg.Infof("user has %s %s in the %s service", balance.Amount, erc20Details.Contract.GetName(), balance.ServiceName)
		}
	}
}

func WithdrawERC20Token(lg *log.Logger, moveTokenInfo MoveTokenInfo, amountStr string) error {

	withdrawalR, _, err := rest.GetWithdrawalInstructions(moveTokenInfo.RESTAPIInfo, moveTokenInfo.ERC20Details.AssetID, amountStr, nil)
	if err != nil {
		return fmt.Errorf("ERC20 withdrawal %w", err)
	}

	withdrawalInst := withdrawalR.Result

	lg.Trace("calling withdrawFromPositionManager")

	_, withdrawalReceipt, _, err := moveTokenInfo.PosManContract.Send(moveTokenInfo.SendMode, moveTokenInfo.Wallet, util.StrP(moveTokenInfo.PosManContract.GetAddressStr(true)), nil, nil, util.StrP("withdrawFromPositionManager"),
		common.HexToAddress(withdrawalInst.Args.Asset), withdrawalInst.Args.AmountContract, withdrawalInst.Args.Nonce, withdrawalInst.Args.Deadline, common.FromHex(withdrawalInst.Args.Signature))

	if err != nil {
		return fmt.Errorf("ERC20 withdrawal %w", err)
	}
	if withdrawalReceipt != nil {
		lg.Tracef("ERC20 withdrawal receipt: %+v", *withdrawalReceipt)
	}

	return nil
}

func DepositERC20Token(lg *log.Logger, moveTokenInfo MoveTokenInfo, amountStr string) error {

	// ERC20 deposit preapproval
	depositAmountWei, err := evm.AmountFromText(amountStr, moveTokenInfo.ERC20Details.Decimals)
	if err != nil {
		return fmt.Errorf("ERC20 deposit amount parse error: %w", err)
	}

	erc20Address := moveTokenInfo.ERC20Details.Contract.GetAddressStr(true)
	posManAddress := moveTokenInfo.PosManContract.GetAddressStr(true)

	_, approvalReceipt, _, err := moveTokenInfo.ERC20Details.Contract.Send(moveTokenInfo.SendMode, moveTokenInfo.Wallet, util.StrP(erc20Address), nil, nil, util.StrP("approve"), common.HexToAddress(posManAddress), depositAmountWei)
	if err != nil {
		return fmt.Errorf("ERC20 deposit approval %w", err)
	}
	if approvalReceipt != nil {
		lg.Tracef("ERC20 deposit approval receipt: %+v", *approvalReceipt)
	}

	// ERC20 deposit transaction
	depositR, _, err := rest.GetDepositInstructions(moveTokenInfo.RESTAPIInfo, moveTokenInfo.ERC20Details.AssetID, amountStr, nil)
	if err != nil {
		return fmt.Errorf("ERC20 deposit instructions %w", err)
	}

	depositInst := depositR.Result

	_, depositReceipt, _, err := moveTokenInfo.PosManContract.Send(moveTokenInfo.SendMode, moveTokenInfo.Wallet, util.StrP(posManAddress), nil, nil, util.StrP("depositInPositionManager"),
		common.HexToAddress(depositInst.Args.Asset), depositInst.Args.AmountContract, depositInst.Args.Nonce, depositInst.Args.Deadline, common.FromHex(depositInst.Args.Signature))
	if err != nil {
		return fmt.Errorf("ERC20 deposit transaction: %w", err)
	}
	if depositReceipt != nil {
		lg.Tracef("ERC20 deposit receipt: %+v", *depositReceipt)
	}

	return nil
}

func WithdrawNativeToken(lg *log.Logger, moveTokenInfo MoveTokenInfo, amountStr string) error {

	withdrawalR, _, err := rest.GetWithdrawalInstructions(moveTokenInfo.RESTAPIInfo, moveTokenInfo.ERC20Details.AssetID, amountStr, nil)
	if err != nil {
		return fmt.Errorf("native token withdrawal: %w", err)
	}

	withdrawalInst := withdrawalR.Result

	_, withdrawalReceipt, _, err := moveTokenInfo.PosManContract.Send(moveTokenInfo.SendMode, moveTokenInfo.Wallet, util.StrP(moveTokenInfo.PosManContract.GetAddressStr(true)), nil, nil, util.StrP("withdrawAndUnwrapNativeToken"),
		common.HexToAddress(withdrawalInst.Args.Asset), withdrawalInst.Args.AmountContract, withdrawalInst.Args.Nonce, withdrawalInst.Args.Deadline, common.FromHex(withdrawalInst.Args.Signature))

	if err != nil {
		return fmt.Errorf("native withdrawal instruction: %w", err)
	}
	if withdrawalReceipt != nil {
		lg.Tracef("native token withdrawal receipt: %+v", *withdrawalReceipt)
	}

	return nil
}

func DepositNativeToken(lg *log.Logger, moveTokenInfo MoveTokenInfo, amountStr string) error {

	depositR, _, err := rest.GetDepositInstructions(moveTokenInfo.RESTAPIInfo, moveTokenInfo.ERC20Details.AssetID, amountStr, nil)
	if err != nil {
		return fmt.Errorf("native token deposit instruction: %w", err)
	}

	depositInst := depositR.Result

	_, depositReceipt, _, err := moveTokenInfo.PosManContract.Send(moveTokenInfo.SendMode, moveTokenInfo.Wallet, util.StrP(moveTokenInfo.PosManContract.GetAddressStr(true)), depositInst.Args.AmountContract, nil, util.StrP("depositAndWrapNativeToken"),
		common.HexToAddress(depositInst.Args.Asset), depositInst.Args.Nonce, depositInst.Args.Deadline, common.FromHex(depositInst.Args.Signature))

	if err != nil {
		return fmt.Errorf("native token deposit: %w", err)
	}
	if depositReceipt != nil {
		lg.Tracef("native token deposit receipt: %+v", *depositReceipt)
	}

	return nil
}
