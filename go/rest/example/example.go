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
	d := rest.ConnectionDetails{DomainName: domainName, APIRoute: apiRoute, APIKey: apiKey, APISecret: apiSecret}

	nonce, err := rest.GetAPIKeyNonce(d)
	if err != nil {
		fmt.Printf("%s", err)
		return
	}

	fmt.Printf("nonce: %d", *nonce)

	orderRequest := rest.SubmitOrderRequest{BaseAssetID: baseAssetID, QuoteAssetID: quoteAssetID, Type: "MARKET", Side: "BUY", QuantityStr: "0.1"}

	submitSuccessResponse, nonce, err := rest.SubmitOrder(d, orderRequest, nonce)
	if err != nil {
		fmt.Printf("\n%s", err)
		return
	}

	fmt.Printf("\n%+v", *submitSuccessResponse.Result)

	instr, _, err := rest.GetDepositInstructions(d, baseAssetID, "0.1", nonce)
	if err != nil {
		fmt.Printf("\n%s", err)
		return
	}

	fmt.Printf("\n%+v", instr.Result)
}
