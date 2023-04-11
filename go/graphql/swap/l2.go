package swap

import (
	"deepwaters/go-examples/util"
	"encoding/json"
)

type OrderBookLevel struct {
	Depth    int    `json:"depth"`
	Price    string `json:"price"`
	Quantity string `json:"quantity"`
}

type OrderBookLevelUpdate struct {
	BaseAssetID  string          `json:"baseAssetID"`
	QuoteAssetID string          `json:"quoteAssetID"`
	Time         *util.Timestamp `json:"time"`
	Side         OrderSide       `json:"side"`
	Price        string          `json:"price"`
	Quantity     string          `json:"quantity"`
}

type L2FeedFrame struct {
	Data       L2FeedFrameInner
	Errors     json.RawMessage
	Extensions map[string]interface{}
}

type L2FeedFrameInner struct {
	OrderBook []OrderBookLevelUpdate
}

type OrderBookSnapshot struct {
	Count  uint64            `json:"count"`
	Cursor *uint64           `json:"cursor"`
	Time   *util.Timestamp   `json:"time"`
	Bids   []*OrderBookLevel `json:"bids"`
	Asks   []*OrderBookLevel `json:"asks"`
}
