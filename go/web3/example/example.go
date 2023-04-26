package main

import (
	"context"
	"deepwaters/go-examples/rest"
	"deepwaters/go-examples/util/evm"
	"deepwaters/go-examples/web3"
	"os"

	abiPackage "github.com/ethereum/go-ethereum/accounts/abi"
	log "github.com/sirupsen/logrus"
)

const (
	chainID                = 43113
	chainName              = "fuji"
	rpcURL                 = "insert node URL here"
	perSecRateLimit        = 32
	maxGasPriceWeiStr      = "50000000000"
	addEstimatedGasPercent = 30

	customerAddress    = ""
	customerPrivateKey = ""
	domainName         = "testnet.api.deepwaters.xyz"
	restAPIRoute       = "/rest/v1/"
	restAPIKey         = "request in the testnet webapp"
	restAPISecret      = "request in the testnet webapp"
)

func main() {

	positionManagerConfig := evm.GetFujiTestnetProdPositionManagerConfig()
	wavaxConfig := evm.GetFujiTestnetProdWAVAXConfig()

	restAPIInfo := rest.ConnectionDetails{
		DomainName: domainName,
		APIRoute:   restAPIRoute,
		APIKey:     restAPIKey,
		APISecret:  restAPISecret,
	}

	lg := log.New()
	lg.SetOutput(os.Stdout)
	lg.SetLevel(log.TraceLevel)

	contractCfgs := []evm.ContractConfig{positionManagerConfig, wavaxConfig}
	cfg, err := evm.NewEVMConnectorConfig(chainID, chainName, rpcURL, perSecRateLimit, maxGasPriceWeiStr, addEstimatedGasPercent, evm.GetFujiAVAX(), contractCfgs)
	if err != nil {
		lg.Errorf("%s", err)
		return
	}
	connector, err := evm.NewEVMConnector(lg, cfg)
	if err != nil {
		lg.Errorf("%s", err)
		return
	}
	if err := connector.Open(); err != nil {
		lg.Errorf("%s", err)
		return
	}
	asset, err := rest.GetAssetReferenceData(restAPIInfo, "WAVAX")
	if err != nil {
		lg.Errorf("%s", err)
		return
	}

	wavaxAssetID := asset.AssetID
	wavaxContract := connector.GetContracts().GetContract(connector.GetConfig().ChainName, wavaxConfig.AddressHexStr)
	decimalsResult, err := wavaxContract.Call(context.Background(), nil, nil, "decimals")
	if err != nil {
		lg.Errorf("%+v", decimalsResult)
	}
	wavaxDecimals := abiPackage.ReadInteger(wavaxContract.GetABI().Methods["decimals"].Outputs[0].Type, decimalsResult).(uint8)
	wavaxDetails := web3.ERC20Details{Contract: wavaxContract, Decimals: wavaxDecimals, AssetID: wavaxAssetID}

	web3.CheckNativeAndWrappedNativeBalances(lg, restAPIInfo, customerAddress, connector, wavaxDetails)

	wallet, err := evm.NewWalletFromHexStr(customerPrivateKey) // note also NewWalletFromFilePath
	if err != nil {
		lg.Errorf("%s", err)
		return
	}

	posManContract := connector.GetContracts().GetContract(connector.GetConfig().ChainName, positionManagerConfig.AddressHexStr)

	moveTokenInfo := web3.MoveTokenInfo{
		SendMode:       evm.SendModeWaitForReceipt,
		RESTAPIInfo:    restAPIInfo,
		Wallet:         wallet,
		PosManContract: posManContract,
		ERC20Details:   wavaxDetails,
	}

	// WAVAX withdrawal
	wavaxWithdrawalAmountStr := ".05"
	if err := web3.WithdrawERC20Token(lg, moveTokenInfo, wavaxWithdrawalAmountStr); err != nil {
		lg.Errorf("%s", err)
		return
	}

	// WAVAX deposit
	wavaxDepositAmountStr := ".1"
	if err := web3.DepositERC20Token(lg, moveTokenInfo, wavaxDepositAmountStr); err != nil {
		lg.Errorf("%s", err)
		return
	}

	/* these do not work on testnet, because of the nature of the testnet tokens

	// AVAX withdrawal
	avaxWithdrawalAmountStr := ".01"
	if err := web3.WithdrawNativeToken(lg, moveTokenInfo, avaxWithdrawalAmountStr); err != nil {
		lg.Errorf("%s", err)
		return
	}

	// AVAX deposit
	avaxDepositAmountStr := ".02"

	if err := web3.DepositNativeToken(lg, moveTokenInfo, avaxDepositAmountStr); err != nil {
		lg.Errorf("%s, err")
		return
	}
	*/

	web3.CheckNativeAndWrappedNativeBalances(lg, restAPIInfo, customerAddress, connector, wavaxDetails)
}
