package main

import (
	"deepwaters/go-examples/rest"
	"fmt"
)

const (
	domainName   = "testnet.api.deepwaters.xyz"
	apiRoute     = "/rest/v1/"
	apiKey       = "request in the testnet webapp"
	apiSecret    = "request in the testnet webapp"
	baseAssetID  = "WBTC.GOERLI.5.TESTNET.PROD"
	quoteAssetID = "USDC.GOERLI.5.TESTNET.PROD"
)

func main() {
	nonce, err := rest.GetAPIKeyNonce(domainName, apiRoute, apiKey, apiSecret)
	if err != nil {
		fmt.Printf("%s", err)
	} else {
		fmt.Printf("nonce: %d", *nonce)
		successResponse, err := rest.SubmitAMarketOrder(domainName, apiRoute, apiKey, apiSecret, baseAssetID, quoteAssetID, *nonce)
		if err != nil {
			fmt.Printf("%s", err)
		} else {
			fmt.Printf("\n%+v", *successResponse.Result)
		}
	}
}
