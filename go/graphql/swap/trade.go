package swap

import (
	"deepwaters/go-examples/util"
	"encoding/json"
)

type TradeFilter struct {
	Customer               *string      `json:"customer,omitempty"`
	CustomerLike           *string      `json:"customer_like,omitempty"`
	CustomerObjectID       *string      `json:"customerObjectID,omitempty"`
	CustomerObjectIDLike   *string      `json:"customerObjectID_like,omitempty"`
	TradeID                *string      `json:"tradeID,omitempty"`
	TradeIDLike            *string      `json:"tradeID_like,omitempty"`
	Type                   *TradeType   `json:"type,omitempty"`
	BaseAssetID            *string      `json:"baseAssetID,omitempty"`
	QuoteAssetID           *string      `json:"quoteAssetID,omitempty"`
	CreatedAtOrAfterMicros *util.BigInt `json:"createdAtOrAfterMicros,omitempty"`
	CreatedBeforeMicros    *util.BigInt `json:"createdBeforeMicros,omitempty"`
}

type TradeList struct {
	Count  uint64   `json:"count"`
	Cursor *uint64  `json:"cursor"`
	Trades []*Trade `json:"trades"`
}

type Trade struct {
	TradeID      string
	CreatedAt    *util.Timestamp
	BaseAssetID  string
	QuoteAssetID string
	Quantity     string
	Volume       string
	Price        string
	Type         TradeType
	Aggressor    TradeSide
	Maker        TradeSide
}

type TradeSide struct {
	VenueOrderID string
	CreatedAt    *util.Timestamp
	RemainingQty string
	Order        *Order
}

type TradesFeedFrame struct {
	Data       TradesFeedFrameInner
	Errors     json.RawMessage
	Extensions map[string]interface{}
}

type TradesFeedFrameInner struct {
	Trades []Trade
}
