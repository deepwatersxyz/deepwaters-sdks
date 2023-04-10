package rest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

const (
	domainName   = "testnet.api.deepwaters.xyz"
	apiRoute     = "/rest/v1/"
	apiKey       = "create this in the testnet webapp"
	apiSecret    = "create this in the testnet webapp"
	baseAssetID  = "WBTC.GOERLI.5.TESTNET.PROD"
	quoteAssetID = "USDC.GOERLI.5.TESTNET.PROD"
)

func GetRequestURIAndURLFromExtension(domainName, apiRoute, extension string) (string, string) {
	requestURI := apiRoute + extension
	u := url.URL{Scheme: "https", Host: domainName, Path: requestURI}
	return requestURI, u.String()
}

func SendRequest(client *http.Client, req *http.Request) (statusCode *int, bodyBytes []byte, err error) {
	resp, err := client.Do(req)
	if resp != nil && resp.Body != nil {
		defer func() {
			_ = resp.Body.Close()
		}()
	}
	if err != nil {
		return nil, nil, err
	}
	if resp.StatusCode >= 429 {
		return nil, nil, fmt.Errorf("%s", resp.Status)
	}
	bodyBytes, err = io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, err
	}
	if 400 <= resp.StatusCode && resp.StatusCode < 429 {
		var errorResponse ErrorResponse
		err = json.Unmarshal(bodyBytes, &errorResponse)
		if err != nil {
			return nil, nil, err
		}
		return nil, nil, fmt.Errorf("status code %d, %s", resp.StatusCode, errorResponse.Error)
	}
	return &resp.StatusCode, bodyBytes, nil
}

func GetAPIKeyNonce(domainName, apiRoute, apiKey, apiSecret string) (*uint64, error) {
	extension := "customer/api-key-status"
	requestURI, url := GetRequestURIAndURLFromExtension(domainName, apiRoute, extension)
	headers, err := GetAuthenticationHeaders(apiKey, apiSecret, "GET", requestURI, nil, nil)
	if err != nil {
		return nil, err
	}
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	statusCode, bodyBytes, err := SendRequest(client, req)
	if err != nil {
		return nil, err
	}
	if *statusCode == 200 {
		var successResponse GetAPIKeySessionSuccessResponse
		err = json.Unmarshal(bodyBytes, &successResponse)
		if err != nil {
			return nil, err
		}
		nonce := successResponse.Result.Nonce
		return &nonce, nil
	} else {
		return nil, fmt.Errorf("unexpected status code %d", *statusCode)
	}
}

func SubmitAMarketOrder(domainName, apiRoute, apiKey, apiSecret string, nonce uint64) (*SubmitOrderSuccessResponse, error) {
	extension := "orders"
	requestURI, url := GetRequestURIAndURLFromExtension(domainName, apiRoute, extension)

	payload := SubmitOrderRequest{BaseAssetID: baseAssetID, QuoteAssetID: quoteAssetID, Type: "MARKET", Side: "BUY", QuantityStr: "1.000000"}
	marshalledPayloadBytes, _ := json.Marshal(payload)
	marshalledPayload := string(marshalledPayloadBytes)

	headers, err := GetAuthenticationHeaders(apiKey, apiSecret, "POST", requestURI, &nonce, &marshalledPayload)
	if err != nil {
		return nil, err
	}
	client := &http.Client{}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(marshalledPayloadBytes))
	if err != nil {
		return nil, err
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	req.Header.Set("Content-Type", "application/json")
	statusCode, bodyBytes, err := SendRequest(client, req)
	if err != nil {
		return nil, err
	}
	if *statusCode == 201 {
		var successResponse SubmitOrderSuccessResponse
		err = json.Unmarshal(bodyBytes, &successResponse)
		if err != nil {
			return nil, err
		}
		return &successResponse, nil
	} else {
		return nil, fmt.Errorf("unexpected status code %d", *statusCode)
	}
}

func main() {
	nonce, err := GetAPIKeyNonce(domainName, apiRoute, apiKey, apiSecret)
	if err != nil {
		fmt.Printf("%s", err)
	} else {
		fmt.Printf("nonce: %d", *nonce)
		successResponse, err := SubmitAMarketOrder(domainName, apiRoute, apiKey, apiSecret, *nonce)
		if err != nil {
			fmt.Printf("%s", err)
		} else {
			fmt.Printf("\n%+v", *successResponse.Result)
		}
	}
}
