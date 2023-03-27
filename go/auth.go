package main

import (
	"fmt"
	"strings"
	"time"
)

func NowStr() string {
	t := uint64(time.Now().UnixMicro())
	return fmt.Sprintf("%d", t)
}

type AuthenticationHeaderValues struct {
	ApiKey             string
	TimestampMicrosStr string
	NonceStr           *string
	SignatureHexStr    string
}

func (v *AuthenticationHeaderValues) ToHeadersMap() map[string]string {
	m := make(map[string]string)
	m["X-DW-APIKEY"] = v.ApiKey
	m["X-DW-TSUS"] = v.TimestampMicrosStr
	if v.NonceStr != nil {
		m["X-DW-NONCE"] = *v.NonceStr
	}
	m["X-DW-SIGHEX"] = v.SignatureHexStr
	return m
}

func GetAuthenticationHeaders(apiKey string, apiSecret string, httpMethod string, requestURI string, nonceStr *string, marshalledPayload *string) (map[string]string, error) {

	nowStr := NowStr()

	headerValues := AuthenticationHeaderValues{ApiKey: apiKey,
		TimestampMicrosStr: nowStr,
		NonceStr:           nonceStr}

	toHashAndSign := httpMethod + strings.ToLower(requestURI) + nowStr

	if nonceStr != nil {
		toHashAndSign += *nonceStr
	}

	if marshalledPayload != nil {
		toHashAndSign += *marshalledPayload
	}
	signatureHexStr, err := HashAndSignFromPrivateKeyHexStr(toHashAndSign, apiSecret)
	if err != nil {
		return nil, err
	}

	headerValues.SignatureHexStr = *signatureHexStr

	return headerValues.ToHeadersMap(), nil
}

func GetRequestURIAndURLFromExtension(extension string) (string, string) {
	requestURI := apiRoute + extension
	url := hostURL + requestURI
	return requestURI, url
}
