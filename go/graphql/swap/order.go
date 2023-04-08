package swap

import (
	"deepwaters/go-examples/util"
	"encoding/json"
)

type OrderFilter struct {
	Customer               *string            `json:"customer,omitempty"`
	CustomerLike           *string            `json:"customer_like,omitempty"`
	CustomerObjectID       *string            `json:"customerObjectID,omitempty"`
	CustomerObjectIDLike   *string            `json:"customerObjectID_like,omitempty"`
	VenueOrderID           *string            `json:"venueOrderID,omitempty"`
	VenueOrderIDLike       *string            `json:"venueOrderID_like,omitempty"`
	Type                   *OrderType         `json:"type,omitempty"`
	Side                   *OrderSide         `json:"side,omitempty"`
	BaseAssetID            *string            `json:"baseAssetID,omitempty"`
	QuoteAssetID           *string            `json:"quoteAssetID,omitempty"`
	StatusIn               []OrderStatus      `json:"status_in,omitempty"`
	DurationType           *OrderDurationType `json:"durationType,omitempty"`
	CreatedAtOrAfterMicros *util.BigInt       `json:"createdAtOrAfterMicros,omitempty"`
	CreatedBeforeMicros    *util.BigInt       `json:"createdBeforeMicros,omitempty"`
}

type OrderList struct {
	Count  uint64   `json:"count"`
	Cursor *uint64  `json:"cursor"`
	Orders []*Order `json:"orders"`
}

type Customer struct {
	Address    string
	ModifiedAt *util.Timestamp
}

type OrderBody struct {
	CreatedAt        *util.Timestamp
	Customer         Customer
	CustomerObjectID *string
	Type             OrderType
	Side             OrderSide
	Quantity         string
	BaseAssetID      string
	QuoteAssetID     string
	Price            *string
	DurationType     OrderDurationType
	ExpiresAt        *util.Timestamp
	ExpiresIn *string
}

type SubmitOrderRequest struct {
	OrderBody
	Nonce     util.BigInt
	Signature string
}

type SubmitOrderResponse struct {
	Request     *SubmitOrderRequest
	RespondedAt *util.Timestamp
	Status      OrderStatus
	Order       *Order
	Error       *string
}

type AmendOrderRequest struct {
	Customer         string
	CustomerObjectID *string
	VenueOrderID     *string
	NewQuantityStr   *string
	NewDurationType  *OrderDurationType
	NewExpiresAt     *util.Timestamp
	NewExpiresIn     *string
	Nonce            util.BigInt
	Signature        string
}

type AmendOrderResponse struct {
	Request     *AmendOrderRequest
	RespondedAt *util.Timestamp
	Status      OrderStatus
	Order       *Order
	Error       *string
}

type ReplaceOrderRequest struct {
	RequestedAt         *util.Timestamp
	StartedAt           *util.Timestamp
	Customer            string
	CustomerObjectID    *string
	VenueOrderID        *string
	NewVenueOrderID     *string
	NewCustomerObjectID *string
	NewQuantityStr      *string
	NewPriceStr         *string
	NewDurationType     *OrderDurationType
	NewExpiresAt        *util.Timestamp
	NewExpiresIn        *string
	Nonce               util.BigInt
	Signature           string
}

type ReplaceOrderResponse struct {
	Request          *ReplaceOrderRequest
	RespondedAt      *util.Timestamp
	Status           OrderStatus
	ReplacementOrder *Order
	ReplacedOrder    *Order
	Error            *string
}

type CancelOrderRequest struct {
	RequestedAt      *util.Timestamp
	StartedAt        *util.Timestamp
	Customer         string
	CustomerObjectID *string
	VenueOrderID     *string
	Nonce            util.BigInt
	Signature        string
}

type CancelOrderResponse struct {
	Request     *CancelOrderRequest
	RespondedAt *util.Timestamp
	Status      OrderStatus
	Error       *string
}

type CancelAllOrdersRequest struct {
	RequestedAt  *util.Timestamp
	Customer     string
	BaseAssetID  *string
	QuoteAssetID *string
	Nonce        util.BigInt
	Signature    string
}

type CancelAllOrdersResponse struct {
	Request      *CancelAllOrdersRequest
	RespondedAt  *util.Timestamp
	NumCancelled int
	Error        *string
}

type Order struct {
	VenueOrderID string
	OrderBody
	OriginalQuantity string
	Volume           string
	AveragePrice     *string
	Status           *OrderStatus
	ModifiedAt       *util.Timestamp
}

type L3FeedFrame struct {
	Data       L3FeedFrameInner
	Errors     json.RawMessage
	Extensions map[string]interface{}
}

type L3FeedFrameInner struct {
	Orders []Order
}
