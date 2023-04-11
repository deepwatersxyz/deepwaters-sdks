package swap

import (
	"fmt"
	"io"
	"strconv"
)

type OrderDurationType string

const (
	OrderDurationTypeGoodTillCancel OrderDurationType = "GOOD_TILL_CANCEL"
	OrderDurationTypeGoodTillExpiry OrderDurationType = "GOOD_TILL_EXPIRY"
)

var AllOrderDurationType = []OrderDurationType{
	OrderDurationTypeGoodTillCancel,
	OrderDurationTypeGoodTillExpiry,
}

func (e OrderDurationType) IsValid() bool {
	switch e {
	case OrderDurationTypeGoodTillCancel, OrderDurationTypeGoodTillExpiry:
		return true
	}
	return false
}

func (e OrderDurationType) String() string {
	return string(e)
}

func (e *OrderDurationType) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = OrderDurationType(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid OrderDurationType", str)
	}
	return nil
}

func (e OrderDurationType) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type OrderSide string

const (
	OrderSideBuy  OrderSide = "BUY"
	OrderSideSell OrderSide = "SELL"
)

var AllOrderSide = []OrderSide{
	OrderSideBuy,
	OrderSideSell,
}

func (e OrderSide) IsValid() bool {
	switch e {
	case OrderSideBuy, OrderSideSell:
		return true
	}
	return false
}

func (e OrderSide) String() string {
	return string(e)
}

func (e *OrderSide) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = OrderSide(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid OrderSide", str)
	}
	return nil
}

func (e OrderSide) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type OrderStatus string

const (
	OrderStatusActive          OrderStatus = "ACTIVE"
	OrderStatusRejected        OrderStatus = "REJECTED"
	OrderStatusReplacement     OrderStatus = "REPLACEMENT"
	OrderStatusAmended         OrderStatus = "AMENDED"
	OrderStatusCancelled       OrderStatus = "CANCELLED"
	OrderStatusReplaced        OrderStatus = "REPLACED"
	OrderStatusFilled          OrderStatus = "FILLED"
	OrderStatusPartiallyFilled OrderStatus = "PARTIALLY_FILLED"
	OrderStatusExpired         OrderStatus = "EXPIRED"
)

var AllOrderStatus = []OrderStatus{
	OrderStatusActive,
	OrderStatusRejected,
	OrderStatusReplacement,
	OrderStatusAmended,
	OrderStatusCancelled,
	OrderStatusReplaced,
	OrderStatusFilled,
	OrderStatusPartiallyFilled,
	OrderStatusExpired,
}

func (e OrderStatus) IsValid() bool {
	switch e {
	case OrderStatusActive, OrderStatusRejected, OrderStatusReplacement, OrderStatusAmended, OrderStatusCancelled, OrderStatusReplaced, OrderStatusFilled, OrderStatusPartiallyFilled, OrderStatusExpired:
		return true
	}
	return false
}

func (e OrderStatus) String() string {
	return string(e)
}

func (e *OrderStatus) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = OrderStatus(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid OrderStatus", str)
	}
	return nil
}

func (e OrderStatus) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type OrderType string

const (
	OrderTypeLimit  OrderType = "LIMIT"
	OrderTypeMarket OrderType = "MARKET"
)

var AllOrderType = []OrderType{
	OrderTypeLimit,
	OrderTypeMarket,
}

func (e OrderType) IsValid() bool {
	switch e {
	case OrderTypeLimit, OrderTypeMarket:
		return true
	}
	return false
}

func (e OrderType) String() string {
	return string(e)
}

func (e *OrderType) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = OrderType(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid OrderType", str)
	}
	return nil
}

func (e OrderType) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type TradeType string

const (
	TradeTypeFill        TradeType = "FILL"
	TradeTypePartialFill TradeType = "PARTIAL_FILL"
)

var AllTradeType = []TradeType{
	TradeTypeFill,
	TradeTypePartialFill,
}

func (e TradeType) IsValid() bool {
	switch e {
	case TradeTypeFill, TradeTypePartialFill:
		return true
	}
	return false
}

func (e TradeType) String() string {
	return string(e)
}

func (e *TradeType) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = TradeType(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid TradeType", str)
	}
	return nil
}

func (e TradeType) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}
