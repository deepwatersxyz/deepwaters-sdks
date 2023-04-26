package rest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type ConnectionDetails struct {
	DomainName string
	APIRoute   string
	APIKey     string
	APISecret  string
}

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

func GetAPIKeyNonce(d ConnectionDetails) (*uint64, error) {
	extension := "customer/api-key-status"
	requestURI, url := GetRequestURIAndURLFromExtension(d.DomainName, d.APIRoute, extension)
	headers, err := GetAuthenticationHeaders(d.APIKey, d.APISecret, "GET", requestURI, nil, nil)
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
		return nil, fmt.Errorf("unexpected status code: %d", *statusCode)
	}
}

func GetUserInfo(d ConnectionDetails) (*GetCustomerSuccessResponse, error) {
	extension := "customer"
	requestURI, url := GetRequestURIAndURLFromExtension(d.DomainName, d.APIRoute, extension)
	headers, err := GetAuthenticationHeaders(d.APIKey, d.APISecret, "GET", requestURI, nil, nil)
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
		var successResponse GetCustomerSuccessResponse
		err = json.Unmarshal(bodyBytes, &successResponse)
		if err != nil {
			return nil, err
		}
		return &successResponse, nil
	} else {
		return nil, fmt.Errorf("unexpected status code: %d", *statusCode)
	}
}

func SubmitOrder(d ConnectionDetails, orderRequest SubmitOrderRequest, nonce *uint64) (*SubmitOrderSuccessResponse, *uint64, error) {

	var err error
	if nonce == nil {
		nonce, err = GetAPIKeyNonce(d)
		if err != nil {
			return nil, nil, err
		}
	}

	extension := "orders"
	requestURI, url := GetRequestURIAndURLFromExtension(d.DomainName, d.APIRoute, extension)

	marshalledPayloadBytes, _ := json.Marshal(orderRequest)
	marshalledPayload := string(marshalledPayloadBytes)

	headers, err := GetAuthenticationHeaders(d.APIKey, d.APISecret, "POST", requestURI, nonce, &marshalledPayload)
	if err != nil {
		return nil, nil, err
	}
	client := &http.Client{}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(marshalledPayloadBytes))
	if err != nil {
		return nil, nil, err
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	req.Header.Set("Content-Type", "application/json")
	statusCode, bodyBytes, err := SendRequest(client, req)
	if err != nil {
		return nil, nil, err
	}
	if *statusCode == 201 {
		var successResponse SubmitOrderSuccessResponse
		err = json.Unmarshal(bodyBytes, &successResponse)
		if err != nil {
			return nil, nil, err
		}
		*nonce += 1
		return &successResponse, nonce, nil
	} else {
		return nil, nil, fmt.Errorf("unexpected status code: %d", *statusCode)
	}
}

// amount is human readable e.g. ".01"
func GetDepositInstructions(d ConnectionDetails, assetID, amount string, nonce *uint64) (*DepositInstructionsSuccessResponse, *uint64, error) {

	var err error
	if nonce == nil {
		nonce, err = GetAPIKeyNonce(d)
		if err != nil {
			return nil, nil, err
		}
	}

	extension := "deposit-instructions"
	requestURI, url := GetRequestURIAndURLFromExtension(d.DomainName, d.APIRoute, extension)

	payload := DepositInstructionsRequest{Asset: assetID, Amount: amount}
	marshalledPayloadBytes, _ := json.Marshal(payload)
	marshalledPayload := string(marshalledPayloadBytes)

	headers, err := GetAuthenticationHeaders(d.APIKey, d.APISecret, "POST", requestURI, nonce, &marshalledPayload)
	if err != nil {
		return nil, nil, err
	}
	client := &http.Client{}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(marshalledPayloadBytes))
	if err != nil {
		return nil, nil, err
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	req.Header.Set("Content-Type", "application/json")
	statusCode, bodyBytes, err := SendRequest(client, req)
	if err != nil {
		return nil, nil, err
	}
	if *statusCode == 201 {
		var successResponse DepositInstructionsSuccessResponse
		err = json.Unmarshal(bodyBytes, &successResponse)
		if err != nil {
			return nil, nil, err
		}
		*nonce += 1
		return &successResponse, nonce, nil
	} else {
		return nil, nil, fmt.Errorf("unexpected status code: %d", *statusCode)
	}
}

// amount is human readable e.g. ".01"
func GetWithdrawalInstructions(d ConnectionDetails, assetID, amount string, nonce *uint64) (*WithdrawalInstructionsSuccessResponse, *uint64, error) {

	var err error
	if nonce == nil {
		nonce, err = GetAPIKeyNonce(d)
		if err != nil {
			return nil, nil, err
		}
	}

	extension := "withdrawal-instructions"
	requestURI, url := GetRequestURIAndURLFromExtension(d.DomainName, d.APIRoute, extension)

	payload := WithdrawalInstructionsRequest{Asset: assetID, Amount: amount}
	marshalledPayloadBytes, _ := json.Marshal(payload)
	marshalledPayload := string(marshalledPayloadBytes)

	headers, err := GetAuthenticationHeaders(d.APIKey, d.APISecret, "POST", requestURI, nonce, &marshalledPayload)
	if err != nil {
		return nil, nil, err
	}
	client := &http.Client{}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(marshalledPayloadBytes))
	if err != nil {
		return nil, nil, err
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	req.Header.Set("Content-Type", "application/json")
	statusCode, bodyBytes, err := SendRequest(client, req)
	if err != nil {
		return nil, nil, err
	}
	if *statusCode == 201 {
		var successResponse WithdrawalInstructionsSuccessResponse
		err = json.Unmarshal(bodyBytes, &successResponse)
		if err != nil {
			return nil, nil, err
		}
		*nonce += 1
		return &successResponse, nil, nil
	} else {
		return nil, nil, fmt.Errorf("unexpected status code: %d", *statusCode)
	}
}

func GetAssetReferenceData(d ConnectionDetails, ticker string) (*Asset, error) {

	extension := "assets"
	requestURI, url := GetRequestURIAndURLFromExtension(d.DomainName, d.APIRoute, extension)

	headers, err := GetAuthenticationHeaders(d.APIKey, d.APISecret, "GET", requestURI, nil, nil)
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
		var successResponse AssetsSuccessResponse
		err = json.Unmarshal(bodyBytes, &successResponse)
		if err != nil {
			return nil, err
		}
		for _, asset := range successResponse.Result {
			if asset.Ticker == ticker {
				return asset, nil
			}
		}
		return nil, fmt.Errorf("could not find ticker: %s", ticker)
	} else {
		return nil, fmt.Errorf("unexpected status code: %d", *statusCode)
	}
}