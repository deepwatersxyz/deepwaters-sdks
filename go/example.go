package main

import (
	"deepwaters/go-examples/graphql/subscriptions"
	"deepwaters/go-examples/util"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
)

const (
	envName = "mainnet beta"
	apiRoot = "api.deepwaters.xyz"
)

func main() {
	lg := util.NewTextLogger(log.TraceLevel, true, os.Stdout)
	baseAssetID := "WAVAX_AM_MB"
	quoteAssetID := "USDC_EM_MB"
	gatherer := subscriptions.NewGatherer(lg, envName, apiRoot)
	gatherer.SetL3WebsocketClient(&baseAssetID, &quoteAssetID, nil)
	gatherer.SetL2WebsocketClient(baseAssetID, quoteAssetID)
	gatherer.SetTradesWebsocketClient(nil, nil, nil)

	go func() { // just for demonstration
		time.Sleep(10 * time.Second)
		gatherer.RestartWebsocketClient("L3")
	}()
	gatherer.Run() // starts all clients
}
