package rest

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

func GetAuthenticationHeaders(apiKey string, apiSecret string, httpMethod string, requestURI string, nonce *uint64, marshalledPayload *string) (map[string]string, error) {

	nowStr := NowStr()

	headerValues := AuthenticationHeaderValues{ApiKey: apiKey,
		TimestampMicrosStr: nowStr}

	toHashAndSign := httpMethod + strings.ToLower(requestURI) + nowStr

	if nonce != nil {
		nonceStr := fmt.Sprintf("%d", *nonce)
		headerValues.NonceStr = &nonceStr
		toHashAndSign += nonceStr
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
