package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const (
	hostURL   = "https://staging.testnet.api.deepwaters.xyz"
	apiRoute  = "/rest/v1/"
	apiKey    = "create this in the webapp"
	apiSecret = "create this in the webapp"
)

func GetAPIKeyNonce(apiKey string, apiSecret string) (uint64, error) {
	extension := "customer/api-key-status"
	requestURI, url := GetRequestURIAndURLFromExtension(extension)
	headers, err := GetAuthenticationHeaders(apiKey, apiSecret, "GET", requestURI, nil, nil)
	if err != nil {
		return 0, err
	}
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return 0, err
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	resp, err := client.Do(req)
	if resp != nil && resp.Body != nil {
		defer func() {
			_ = resp.Body.Close()
		}()
	}
	if err != nil {
		return 0, err
	}
	if resp.StatusCode >= 429 {
		return 0, fmt.Errorf("%s", resp.Status)
	}
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}
	if resp.StatusCode == 400 || resp.StatusCode == 422 {
		var errorResponse ErrorResponse
		err = json.Unmarshal(bodyBytes, &errorResponse)
		if err != nil {
			return 0, err
		}
		return 0, fmt.Errorf("status code %d, %s", resp.StatusCode, errorResponse.Error) //result["Error"].(string))
	}
	if resp.StatusCode == 200 {
		var successResponse GetAPIKeySessionSuccessResponse
		err = json.Unmarshal(bodyBytes, &successResponse)
		if err != nil {
			return 0, err
		}
		nonce := successResponse.Result.Nonce
		return nonce, nil
	} else {
		return 0, fmt.Errorf("unexpected status code %d", resp.StatusCode)
	}
}

func main() {
	nonce, err := GetAPIKeyNonce(apiKey, apiSecret)
	if err != nil {
		fmt.Printf("%s", err)
	} else {
		fmt.Printf("nonce: %d", nonce)
	}
}
