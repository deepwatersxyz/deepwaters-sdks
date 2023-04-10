package main

import (
	"deepwaters/go-examples/graphql/subscriptions"
	"deepwaters/go-examples/rest"
	"deepwaters/go-examples/util"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
)

const (
	envName       = "testnet prod"
	domainName    = "testnet.api.deepwaters.xyz"
	restAPIRoute  = "/rest/v1/"
	restAPIKey    = "request in the testnet webapp"
	restAPISecret = "request in the testnet webapp"
	baseAssetID   = "WBTC.GOERLI.5.TESTNET.PROD"
	quoteAssetID  = "USDC.GOERLI.5.TESTNET.PROD"
)

func main() {
	lg := util.NewTextLogger(log.TraceLevel, true, os.Stdout)
	gatherer := subscriptions.NewGatherer(lg, envName, domainName)
	gatherer.SetL3WebsocketClient(baseAssetID, quoteAssetID, "")
	gatherer.SetL2WebsocketClient(baseAssetID, quoteAssetID)
	gatherer.SetTradesWebsocketClient("", "", "")

	go func() { // demonstrates websocket client restarts and interacting with the REST API
		time.Sleep(10 * time.Second)
		gatherer.RestartWebsocketClient("L3")
		time.Sleep(10 * time.Second)
		nonce, err := rest.GetAPIKeyNonce(domainName, restAPIRoute, restAPIKey, restAPISecret)
		if err != nil {
			lg.Errorf("%s", err)
		} else {
			lg.Debugf("nonce: %d", *nonce)
			successResponse, err := rest.SubmitAMarketOrder(domainName, restAPIRoute, restAPIKey, restAPISecret, *nonce)
			if err != nil {
				lg.Errorf("%s", err)
			} else {
				lg.Debugf("\n%+v", *successResponse.Result)
			}
		}

	}()
	gatherer.Run() // starts all subscription (websocket) clients
}
