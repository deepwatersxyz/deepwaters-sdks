// websocket subscriptions example

package main

import (
	"deepwaters/go-examples/graphql/subscriptions"

	"os"
	"time"

	log "github.com/sirupsen/logrus"
)

const (
	envName       = "testnet prod"
	domainName    = "testnet.api.deepwaters.xyz"
	baseAssetID   = "WBTC.GOERLI.5.TESTNET.PROD"
	quoteAssetID  = "USDC.GOERLI.5.TESTNET.PROD"
)

func main() {
	lg := log.New()
	lg.SetOutput(os.Stdout)
	lg.SetLevel(log.TraceLevel)

	gatherer := subscriptions.NewGatherer(lg, envName, domainName)
	gatherer.SetL3WebsocketClient(baseAssetID, quoteAssetID, "")
	gatherer.SetL2WebsocketClient(baseAssetID, quoteAssetID)
	gatherer.SetTradesWebsocketClient("", "", "")

	go func() { // demonstrates websocket client restarts
		time.Sleep(10 * time.Second)
		gatherer.RestartWebsocketClient("L3")
	}()
	gatherer.Run() // starts all subscription (websocket) clients
}
