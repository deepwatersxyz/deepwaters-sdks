package swap

import (
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/crypto"

	"deepwaters/go-examples/graphql/eip712"
	"deepwaters/go-examples/util"
)

func getSubmitOrderSignature(customerPrivKey string, nonce uint64, customerObjectID *string, side OrderSide, quantity, baseAssetID, quoteAssetID, price string) (string, string) {
	privKey, err := crypto.HexToECDSA(customerPrivKey)
	if err != nil {
		panic(err)
	}
	address := crypto.PubkeyToAddress(privKey.PublicKey)

	request := &SubmitOrderRequest{
		OrderBody: OrderBody{
			Customer:         Customer{Address: address.Hex()},
			CustomerObjectID: customerObjectID,
			Type:             OrderTypeLimit,
			Side:             side,
			Quantity:         quantity,
			BaseAssetID:      baseAssetID,
			QuoteAssetID:     quoteAssetID,
			Price:            &price,
			DurationType:     OrderDurationTypeGoodTillCancel,
		},
		Nonce: util.BigInt{Int: big.NewInt(0).SetUint64(nonce)},
	}

	signature, err := eip712.Sign(request, privKey)
	if err != nil {
		panic(err)
	}

	return address.Hex(), *signature
}

func getAmendSignature(customerPrivKey string, nonce uint64, customerObjectID *string, venueOrderID *string, newQuantity *string, newDurationType *OrderDurationType, newExpiresAt *time.Time, newExpiresIn *string) (string, string) {
	privKey, err := crypto.HexToECDSA(customerPrivKey)
	if err != nil {
		panic(err)
	}
	address := crypto.PubkeyToAddress(privKey.PublicKey)

	var expiresAt *util.Timestamp
	if newExpiresAt != nil {
		expiresAt = util.NewTimestampFromTime(*newExpiresAt, true)
	}

	replace := &AmendOrderRequest{
		Customer:         address.Hex(),
		CustomerObjectID: customerObjectID,
		VenueOrderID:     venueOrderID,
		NewQuantityStr:   newQuantity,
		NewDurationType:  newDurationType,
		NewExpiresAt:     expiresAt,
		NewExpiresIn:     newExpiresIn,
		Nonce:            util.BigInt{Int: big.NewInt(0).SetUint64(nonce)},
	}

	signature, err := eip712.Sign(replace, privKey)
	if err != nil {
		panic(err)
	}

	return address.Hex(), *signature
}

func getReplaceSignature(customerPrivKey string, nonce uint64, customerObjectID, venueOrderID, newQuantity, newPrice *string, newDurationType *OrderDurationType, newExpiresAt *time.Time, newExpiresIn *string) (string, string) {
	privKey, err := crypto.HexToECDSA(customerPrivKey)
	if err != nil {
		panic(err)
	}
	address := crypto.PubkeyToAddress(privKey.PublicKey)

	var expiresAt *util.Timestamp
	if newExpiresAt != nil {
		expiresAt = util.NewTimestampFromTime(*newExpiresAt, true)
	}

	replace := &ReplaceOrderRequest{
		Customer:         address.Hex(),
		CustomerObjectID: customerObjectID,
		VenueOrderID:     venueOrderID,
		NewQuantityStr:   newQuantity,
		NewPriceStr:      newPrice,
		NewDurationType:  newDurationType,
		NewExpiresAt:     expiresAt,
		NewExpiresIn:     newExpiresIn,
		Nonce:            util.BigInt{Int: big.NewInt(0).SetUint64(nonce)},
	}

	signature, err := eip712.Sign(replace, privKey)
	if err != nil {
		panic(err)
	}

	return address.Hex(), *signature
}

func getCancelSignature(customerPrivKey string, nonce uint64, customerObjectID *string, venueOrderID *string) (string, string) {
	privKey, err := crypto.HexToECDSA(customerPrivKey)
	if err != nil {
		panic(err)
	}
	address := crypto.PubkeyToAddress(privKey.PublicKey)

	cancel := &CancelOrderRequest{
		Customer:         address.Hex(),
		CustomerObjectID: customerObjectID,
		VenueOrderID:     venueOrderID,
		Nonce:            util.BigInt{Int: big.NewInt(0).SetUint64(nonce)},
	}

	signature, err := eip712.Sign(cancel, privKey)
	if err != nil {
		panic(err)
	}

	return address.Hex(), *signature
}

func getCancelAllSignature(customerPrivKey string, nonce uint64, baseAssetID, quoteAssetID *string) (string, string) {
	privKey, err := crypto.HexToECDSA(customerPrivKey)
	if err != nil {
		panic(err)
	}
	address := crypto.PubkeyToAddress(privKey.PublicKey)

	cancel := &CancelAllOrdersRequest{
		Customer:     address.Hex(),
		Nonce:        util.BigInt{Int: big.NewInt(0).SetUint64(nonce)},
		BaseAssetID:  baseAssetID,
		QuoteAssetID: quoteAssetID,
	}

	signature, err := eip712.Sign(cancel, privKey)
	if err != nil {
		panic(err)
	}

	return address.Hex(), *signature
}
