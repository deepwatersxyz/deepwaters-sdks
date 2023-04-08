package swap

import (
	"fmt"
	"time"
)

func CreateSignedSubmitLimitOrderQuery(customerPrivKey string, nonce uint64, customerObjectID *string, side OrderSide,
	quantity, baseAssetID, quoteAssetID, price string, durationType *OrderDurationType, expiresAt *time.Time, expiresIn *string) (string, string) {

	customer, signature := getSubmitOrderSignature(customerPrivKey, nonce, customerObjectID, side, quantity, baseAssetID, quoteAssetID, price)
	return customer, createSubmitLimitOrderQuery(customer, nonce, customerObjectID, side, quantity, baseAssetID, quoteAssetID, price, durationType, expiresAt, expiresIn, signature)
}

func createSubmitLimitOrderQuery(customer string, nonce uint64, customerObjectID *string, side OrderSide,
	quantity, baseAssetID, quoteAssetID, price string, durationType *OrderDurationType, expiresAt *time.Time, expiresIn *string, signature string) string {

	mutation := fmt.Sprintf(`mutation {
		submitOrder(
			customer: "%s",
			%s
			type: LIMIT,
			side: %s,
			quantity: "%s",
			baseAssetID: "%s",
			quoteAssetID: "%s",
			price: "%s",
			%s
			%s
			%s
			nonce: "%d",
			signature: "%s"
		) %s
	}`,
		customer,
		func() string {
			if customerObjectID == nil {
				return ""
			}
			return fmt.Sprintf(`customerObjectID: "%s",`, *customerObjectID)
		}(),
		string(side),
		quantity,
		baseAssetID,
		quoteAssetID,
		price,
		func() string {
			if durationType == nil {
				return ""
			}
			return fmt.Sprintf("durationType: %s,", *durationType)
		}(),
		func() string {
			if expiresAt == nil {
				return ""
			}
			return fmt.Sprintf(`expiresAt: "%s",`, expiresAt.Format(time.RFC3339Nano))
		}(),
		func() string {
			if expiresIn == nil {
				return ""
			}
			return fmt.Sprintf(`expiresIn: "%s",`, *expiresIn)
		}(),
		nonce,
		signature,
		SubmitOrderResponseFields,
	)

	return mutation
}

func createSubmitMarketOrderQuery(customer string, nonce uint64, customerObjectID *string, side OrderSide, quantity, baseAssetID string, quoteAssetID string, signature string) string {
	return fmt.Sprintf(`mutation {
		submitOrder(
			customer: "%s",
			%s
			type: MARKET,
			side: %s,
			quantity: "%s",
			baseAssetID: "%s",
			quoteAssetID: "%s",
			nonce: "%d",
			signature: "%s"
		) %s
	}`,
		customer,
		func() string {
			if customerObjectID == nil {
				return ""
			}
			return fmt.Sprintf(`customerObjectID: "%s",`, *customerObjectID)
		}(),
		string(side),
		quantity,
		baseAssetID,
		quoteAssetID,
		nonce,
		signature,
		SubmitOrderResponseFields,
	)
}

func CreateSignedAmendQuery(customerPrivKey string, nonce uint64, customerObjectID, orderID, newQuantity *string,
	newDurationType *OrderDurationType, newExpiresAt *time.Time, newExpiresIn *string) (string, string) {

	customer, signature := getAmendSignature(customerPrivKey, nonce, customerObjectID, orderID, newQuantity, newDurationType, newExpiresAt, newExpiresIn)
	return customer, createAmendQuery(customer, nonce, customerObjectID, orderID, newQuantity, newDurationType, newExpiresAt, newExpiresIn, signature)
}

func createAmendQuery(customer string, nonce uint64, customerObjectID, orderID, newQuantity *string,
	newDurationType *OrderDurationType, newExpiresAt *time.Time, newExpiresIn *string, signature string) string {

	f := func(label string, s *string) string {
		if s == nil {
			return ""
		}
		return fmt.Sprintf(`%s: "%s"`, label, *s)
	}

	mutation := fmt.Sprintf(`mutation {
		amendOrder(
			%s
			%s
			%s
			%s
			%s
			%s
			%s
			nonce: "%d",
			signature: "%s"
		) %s
	}`,
		f("venueOrderID", orderID),
		f("customer", &customer),
		f("customerObjectID", customerObjectID),
		f("newQuantity", newQuantity),
		func() string {
			if newDurationType == nil {
				return ""
			}
			return fmt.Sprintf(`newDurationType: %s,`, newDurationType.String())
		}(),
		func() string {
			if newExpiresAt == nil {
				return ""
			}
			return fmt.Sprintf(`newExpiresAt: "%s",`, newExpiresAt.Format(time.RFC3339Nano))
		}(),
		f("newExpiresIn", newExpiresIn),
		nonce,
		signature,
		AmendOrderResponseFields,
	)

	return mutation
}

func CreateSignedReplaceQuery(customerPrivKey string, nonce uint64, customerObjectID, orderID, newCustomerObjectID, newQuantity, newPrice *string,
	newDurationType *OrderDurationType, newExpiresAt *time.Time, newExpiresIn *string) (string, string) {

	customer, signature := getReplaceSignature(customerPrivKey, nonce, customerObjectID, orderID, newQuantity, newPrice, newDurationType, newExpiresAt, newExpiresIn)
	return customer, createReplaceQuery(customer, nonce, customerObjectID, orderID, newCustomerObjectID, newQuantity,
		newPrice, newDurationType, newExpiresAt, newExpiresIn, signature)
}

func createReplaceQuery(customer string, nonce uint64, customerObjectID, orderID, newCustomerObjectID, newQuantity, newPrice *string,
	newDurationType *OrderDurationType, newExpiresAt *time.Time, newExpiresIn *string, signature string) string {

	f := func(label string, s *string) string {
		if s == nil {
			return ""
		}
		return fmt.Sprintf(`%s: "%s"`, label, *s)
	}

	mutation := fmt.Sprintf(`mutation {
		replaceOrder(
			%s
			%s
			%s
			%s
			%s
			%s
			%s
			%s
			%s
			nonce: "%d",
			signature: "%s"
		) %s
	}`,
		f("venueOrderID", orderID),
		f("customer", &customer),
		f("customerObjectID", customerObjectID),
		f("newCustomerObjectID", newCustomerObjectID),
		f("newQuantity", newQuantity),
		f("newPrice", newPrice),
		func() string {
			if newDurationType == nil {
				return ""
			}
			return fmt.Sprintf(`newDurationType: %s,`, newDurationType.String())
		}(),
		func() string {
			if newExpiresAt == nil {
				return ""
			}
			return fmt.Sprintf(`newExpiresAt: "%s",`, newExpiresAt.Format(time.RFC3339Nano))
		}(),
		f("newExpiresIn", newExpiresIn),
		nonce,
		signature,
		ReplaceOrderResponseFields,
	)

	return mutation
}

func CreateSignedCancelQuery(customerPrivKey string, nonce uint64, customerObjectID, venueOrderID *string) string {
	customer, signature := getCancelSignature(customerPrivKey, nonce, customerObjectID, venueOrderID)
	return createCancelQuery(customer, nonce, customerObjectID, venueOrderID, signature)
}

func createCancelQuery(customer string, nonce uint64, customerObjectID, venueOrderID *string, signature string) string {
	return fmt.Sprintf(`mutation {
		cancelOrder(
			customer: "%s",
			%s
			%s
			nonce: "%d",
			signature: "%s"
		) %s
	}`,
		customer,
		func() string {
			if customerObjectID == nil {
				return ""
			}
			return fmt.Sprintf(`customerObjectID: "%s",`, *customerObjectID)
		}(),
		func() string {
			if venueOrderID == nil {
				return ""
			}
			return fmt.Sprintf(`venueOrderID: "%s",`, *venueOrderID)
		}(),
		nonce,
		signature,
		CancelResponseFields,
	)
}

func CreateSignedCancelAllOrdersQuery(customerPrivKey string, nonce uint64, baseAssetID, quoteAssetID *string) string {
	customer, signature := getCancelAllSignature(customerPrivKey, nonce, baseAssetID, quoteAssetID)
	return createCancelAllOrdersQuery(customer, nonce, baseAssetID, quoteAssetID, signature)
}

func createCancelAllOrdersQuery(customer string, nonce uint64, baseAssetID, quoteAssetID *string, signature string) string {
	return fmt.Sprintf(`mutation {
		cancelAllOrders(
			customer: "%v",
			%v
			%v
			nonce: "%v",
			signature: "%v"
		) %v
	}`,
		customer,
		func() string {
			if baseAssetID == nil {
				return ""
			}
			return fmt.Sprintf(`baseAssetID: "%v",`, *baseAssetID)
		}(),
		func() string {
			if quoteAssetID == nil {
				return ""
			}
			return fmt.Sprintf(`quoteAssetID: "%v",`, *quoteAssetID)
		}(),
		nonce,
		signature,
		CancelAllOrdersResponseFields,
	)
}

func CreateOrderBookSnapshotQuery(baseAssetID, quoteAssetID string) string {
	return fmt.Sprintf(`
		query {
			orderBookSnapshot(
				limit: 100,
				skip: 0,
				baseAssetID: "%s",
				quoteAssetID: "%s"
			) %s
		}`, baseAssetID, quoteAssetID, OrderBookSnapshotFields)
}

func CreateShortOrderBookSnapshotQuery(baseAssetID, quoteAssetID string) string {
	return fmt.Sprintf(`
		query {
			orderBookSnapshot(
				limit: 100,
				skip: 0,
				baseAssetID: "%s",
				quoteAssetID: "%s"
			) %s
		}`, baseAssetID, quoteAssetID, ShortOrderBookSnapshotFields)
}

func CreateOrderQuery(customer, customerObjectID, orderID *string) string {
	f := func(label string, s *string) string {
		if s == nil {
			return ""
		}
		return fmt.Sprintf(`%s: "%s"`, label, *s)
	}

	return fmt.Sprintf(`
		query {
			order(
				%v
				%v
				%v
			) %s
		}`,
		f("customer", customer),
		f("customerObjectID", customerObjectID),
		f("orderID", orderID),
		OrderFields,
	)
}
