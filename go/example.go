package main

import (
	"deepwaters/go-examples/graphql/subscriptions"
	"deepwaters/go-examples/util"
	"os"

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
	if err := gatherer.SetL3WebsocketClient(&baseAssetID, &quoteAssetID, nil); err != nil {
		panic(err)
	}
	if err := gatherer.SetL2WebsocketClient(baseAssetID, quoteAssetID); err != nil {
		panic(err)
	}
	if err := gatherer.SetTradesWebsocketClient(nil, nil, nil); err != nil {
		panic(err)
	}
	gatherer.Run()
}
