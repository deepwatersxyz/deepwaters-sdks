package swap

import (
	"github.com/ethereum/go-ethereum/signer/core/apitypes"

	"deepwaters/go-examples/graphql/eip712"
	"deepwaters/go-examples/util"

	"fmt"
	"math/big"
)

func str(d *string) string {
	if d != nil {
		return *d
	}
	return ""
}

func boolean(b *bool) bool {
	if b != nil {
		return *b
	}
	return false
}

func timestamp(ts *util.Timestamp) *big.Int {
	if ts == nil {
		return &big.Int{}
	}
	i := &big.Int{}
	_, _ = i.SetString(fmt.Sprintf("%v", ts.MicrosSinceEpoch), 10)
	return i
}

// stringifies the micros field in model.Timestamp
func timestampStr(ts *util.Timestamp) string {
	if ts == nil {
		return ""
	}
	return fmt.Sprintf("%v", ts.MicrosSinceEpoch)
}

// the commented out components will be added to the signing process with the next breaking changes
func (submit SubmitOrderRequest) ToTypedData() apitypes.TypedData {
	return apitypes.TypedData{
		Types: apitypes.Types{
			eip712.DomainName: eip712.DomainAMType,
			"SubmitOrderRequest": {
				{Name: "customer", Type: "address"},
				{Name: "customerObjectID", Type: "string"},
				{Name: "type", Type: "string"},
				{Name: "side", Type: "string"},
				{Name: "quantity", Type: "string"},
				{Name: "baseAssetID", Type: "string"},
				{Name: "quoteAssetID", Type: "string"},
				{Name: "price", Type: "string"},
				{Name: "durationType", Type: "string"},
				//{Name: "expiresIn", Type: "string"},
				//{Name: "expiresAt", Type: "string"},
				{Name: "nonce", Type: "uint256"},
			},
		},
		PrimaryType: "SubmitOrderRequest",
		Domain: apitypes.TypedDataDomain{
			Name:    "Swap",
			Version: "1",
		},
		Message: apitypes.TypedDataMessage{
			"customer":         submit.Customer,
			"customerObjectID": str(submit.CustomerObjectID),
			"type":             submit.Type.String(),
			"side":             submit.Side.String(),
			"quantity":         submit.Quantity,
			"baseAssetID":      submit.BaseAssetID,
			"quoteAssetID":     submit.QuoteAssetID,
			"price":            str(submit.Price),
			"durationType":     submit.DurationType.String(),
			//"expiresIn":        str(submit.ExpiresIn),
			//"expiresAt":        timestampStr(submit.ExpiresAt),
			"nonce": submit.Nonce.Int.String(),
		},
	}
}

func (amend AmendOrderRequest) ToTypedData() apitypes.TypedData {
	var durationType *string
	if amend.NewDurationType != nil {
		s := amend.NewDurationType.String()
		durationType = &s
	}
	return apitypes.TypedData{
		Types: apitypes.Types{
			eip712.DomainName: eip712.DomainAMType,
			"AmendOrderRequest": {
				{Name: "customer", Type: "address"},
				{Name: "customerObjectID", Type: "string"},
				{Name: "orderID", Type: "string"},
				{Name: "newQuantity", Type: "string"},
				{Name: "newDurationType", Type: "string"},
				{Name: "newExpiresIn", Type: "string"},
				{Name: "newExpiresAt", Type: "string"},
				{Name: "nonce", Type: "uint256"},
			},
		},
		PrimaryType: "AmendOrderRequest",
		Domain: apitypes.TypedDataDomain{
			Name:    "Swap",
			Version: "1",
		},
		Message: apitypes.TypedDataMessage{
			"customer":         amend.Customer,
			"customerObjectID": str(amend.CustomerObjectID),
			"orderID":          str(amend.VenueOrderID),
			"newQuantity":      str(amend.NewQuantityStr),
			"newDurationType":  str(durationType),
			"newExpiresIn":     str(amend.NewExpiresIn),
			"newExpiresAt":     timestampStr(amend.NewExpiresAt),
			"nonce":            amend.Nonce.Int.String(),
		},
	}
}

func (replace ReplaceOrderRequest) ToTypedData() apitypes.TypedData {
	var durationType *string
	if replace.NewDurationType != nil {
		s := replace.NewDurationType.String()
		durationType = &s
	}
	return apitypes.TypedData{
		Types: apitypes.Types{
			eip712.DomainName: eip712.DomainAMType,
			"ReplaceOrderRequest": {
				{Name: "customer", Type: "address"},
				{Name: "customerObjectID", Type: "string"},
				{Name: "orderID", Type: "string"},
				{Name: "newCustomerObjectID", Type: "string"},
				{Name: "newQuantity", Type: "string"},
				{Name: "newPrice", Type: "string"},
				{Name: "newDurationType", Type: "string"},
				{Name: "newExpiresIn", Type: "string"},
				{Name: "newExpiresAt", Type: "string"},
				{Name: "nonce", Type: "uint256"},
			},
		},
		PrimaryType: "ReplaceOrderRequest",
		Domain: apitypes.TypedDataDomain{
			Name:    "Swap",
			Version: "1",
		},
		Message: apitypes.TypedDataMessage{
			"customer":            replace.Customer,
			"customerObjectID":    str(replace.CustomerObjectID),
			"orderID":             str(replace.VenueOrderID),
			"newCustomerObjectID": str(replace.NewCustomerObjectID),
			"newQuantity":         str(replace.NewQuantityStr),
			"newPrice":            str(replace.NewPriceStr),
			"newDurationType":     str(durationType),
			"newExpiresIn":        str(replace.NewExpiresIn),
			"newExpiresAt":        timestampStr(replace.NewExpiresAt),
			"nonce":               replace.Nonce.Int.String(),
		},
	}
}

func (cancel CancelOrderRequest) ToTypedData() apitypes.TypedData {
	return apitypes.TypedData{
		Types: apitypes.Types{
			eip712.DomainName: eip712.DomainAMType,
			"CancelOrderRequest": {
				{Name: "customer", Type: "address"},
				{Name: "customerObjectID", Type: "string"},
				{Name: "orderID", Type: "string"},
				{Name: "nonce", Type: "uint256"},
			},
		},
		PrimaryType: "CancelOrderRequest",
		Domain: apitypes.TypedDataDomain{
			Name:    "Swap",
			Version: "1",
		},
		Message: apitypes.TypedDataMessage{
			"customer":         cancel.Customer,
			"customerObjectID": str(cancel.CustomerObjectID),
			"orderID":          str(cancel.VenueOrderID),
			"nonce":            cancel.Nonce.Int.String(),
		},
	}
}

func (cancelAllOrders CancelAllOrdersRequest) ToTypedData() apitypes.TypedData {
	return apitypes.TypedData{
		Types: apitypes.Types{
			eip712.DomainName: eip712.DomainAMType,
			"CancelAllOrdersRequest": {
				{Name: "customer", Type: "address"},
				{Name: "baseAssetID", Type: "string"},
				{Name: "quoteAssetID", Type: "string"},
				{Name: "nonce", Type: "uint256"},
			},
		},
		PrimaryType: "CancelAllOrdersRequest",
		Domain: apitypes.TypedDataDomain{
			Name:    "Swap",
			Version: "1",
		},
		Message: apitypes.TypedDataMessage{
			"customer":     cancelAllOrders.Customer,
			"baseAssetID":  str(cancelAllOrders.BaseAssetID),
			"quoteAssetID": str(cancelAllOrders.QuoteAssetID),
			"nonce":        cancelAllOrders.Nonce.Int.String(),
		},
	}
}
