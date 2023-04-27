package accounting

import "encoding/json"

type BalanceFilter struct {
	CustomerAddress     *string `json:"customerAddress"`
	CustomerAddressLike *string `json:"customerAddress_like"`
	AssetID             *string `json:"assetID"`
	ServiceID           *string `json:"serviceID"`
}

type ServiceInfo struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type Balance struct {
	CustomerAddress string      `json:"customerAddress"`
	AssetID         string      `json:"assetID"`
	Amount          string      `json:"amount"`
	Service         ServiceInfo `json:"service"`
}

type BalancesFeedFrame struct {
	Data       BalancesFeedFrameInner
	Errors     json.RawMessage
	Extensions map[string]interface{}
}

type BalancesFeedFrameInner struct {
	Balances []Balance
}
