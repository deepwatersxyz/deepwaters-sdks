// example combining websocket subscriptions and REST

package main

import (
	"context"
	"deepwaters/go-examples/graphql/subscriptions"
	"deepwaters/go-examples/rest"
	"deepwaters/go-examples/util/evm"
	"deepwaters/go-examples/web3"

	"os"
	"time"

	abiPackage "github.com/ethereum/go-ethereum/accounts/abi"
	log "github.com/sirupsen/logrus"
)

const (
	envName         = "testnet prod"
	domainName      = "testnet.api.deepwaters.xyz"
	restAPIRoute    = "/rest/v1/"
	customerAddress = ""
	restAPIKey      = "request in the testnet webapp"
	restAPISecret   = "request in the testnet webapp"

	baseAssetID  = "WBTC.GOERLI.5.TESTNET.PROD"
	quoteAssetID = "USDC.GOERLI.5.TESTNET.PROD"

	chainID                = 43113
	chainName              = "fuji"
	rpcURL                 = "insert node URL here"
	perSecRateLimit        = 32
	maxGasPriceWeiStr      = "50000000000"
	addEstimatedGasPercent = 30
)

func main() {
	lg := log.New()
	lg.SetOutput(os.Stdout)
	lg.SetLevel(log.TraceLevel)

	// check AVAX and WAVAX balances, on-chain and in deepwaters
	wavaxConfig := evm.GetFujiTestnetProdWAVAXConfig()
	contractCfgs := []evm.ContractConfig{wavaxConfig}
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

	restAPIInfo := rest.ConnectionDetails{DomainName: domainName, APIRoute: restAPIRoute, APIKey: restAPIKey, APISecret: restAPISecret}

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

	// subscribe to websocket updates
	gatherer := subscriptions.NewGatherer(lg, envName, domainName)
	gatherer.SetL3WebsocketClient(baseAssetID, quoteAssetID, "")
	gatherer.SetL2WebsocketClient(baseAssetID, quoteAssetID)
	gatherer.SetTradesWebsocketClient("", "", "")
	gatherer.SetBalancesWebsocketClient(customerAddress)

	go func() { // demonstrates websocket client restarts and interacting with the REST API

		time.Sleep(10 * time.Second)
		gatherer.RestartWebsocketClient("balances")
		time.Sleep(10 * time.Second)

		orderRequest := rest.SubmitOrderRequest{BaseAssetID: baseAssetID, QuoteAssetID: quoteAssetID, Type: "MARKET", Side: "BUY", QuantityStr: "0.1"}
		successResponse, _, err := rest.SubmitOrder(restAPIInfo, orderRequest, nil)
		if err != nil {
			lg.Errorf("%s", err)
		} else {
			lg.Infof("%+v", *successResponse.Result)
		}
	}()
	gatherer.Run() // starts all subscription (websocket) clients
}
