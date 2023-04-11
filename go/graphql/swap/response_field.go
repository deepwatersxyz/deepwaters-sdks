package swap

// applicable to a variety of queries
const OrderResponseFields = `{
	venueOrderID
	customer {
		address
		modifiedAt {
			time
			microsSinceEpoch
		}
	}
	customerObjectID
	type
	side
	quantity
	originalQuantity
	volume
	baseAssetID
	quoteAssetID
	price
	durationType
	createdAt {
		time
		microsSinceEpoch
	}
	modifiedAt {
		time
		microsSinceEpoch
	}
	expiresAt {
		time
		microsSinceEpoch
	}
	status
	averagePrice
}`

const SubmitOrderResponseFields = `{
	respondedAt {
		time
		microsSinceEpoch
	}
	status
	request {
		customer {
			address
			nonce
			lastCustomerObjectID
		}
		customerObjectID
		type
		side
		quantity
		baseAssetID
		quoteAssetID
		price
		durationType
		requestedAt {
			time
			microsSinceEpoch
		}
		expiresAt {
			time
			microsSinceEpoch
		}
		expiresIn {
			duration
			micros
		}
		nonce
		signature
	}
	order {
		venueOrderID
		customer {
			address
			nonce
			lastCustomerObjectID
		}
		customerObjectID
		type
		side
		originalQuantity
		totalFillQuantity
		totalQuantity
		quantity
		volume
		baseAssetID
		quoteAssetID
		price
		averagePrice
		durationType
		createdAt {
			time
			microsSinceEpoch
		}
		modifiedAt {
			time
			microsSinceEpoch
		}
		expiresAt {
			time
			microsSinceEpoch
		}
		expiresIn {
			duration
			micros
		}
		status
		rfq
		requestNonce
		message
	}
	error
}`

const AmendOrderResponseFields = `{
	respondedAt {
		time
		microsSinceEpoch
	}
	status
	request {
		customer {
			address
			nonce
			lastCustomerObjectID
		}
		customerObjectID
		venueOrderID
		newQuantity
		newDurationType
		newExpiresAt {
			time
			microsSinceEpoch
		}
		newExpiresIn {
			duration
			micros
		}
		nonce
		signature
	}
	order {
		venueOrderID
		customer {
			address
			nonce
			lastCustomerObjectID
		}
		customerObjectID
		type
		side
		originalQuantity
		totalFillQuantity
		totalQuantity
		quantity
		volume
		baseAssetID
		quoteAssetID
		price
		averagePrice
		durationType
		createdAt {
			time
			microsSinceEpoch
		}
		modifiedAt {
			time
			microsSinceEpoch
		}
		expiresAt {
			time
			microsSinceEpoch
		}
		expiresIn {
			duration
			micros
		}
		status
		rfq
		requestNonce
		message
	}
	error
}`

const ReplaceOrderResponseFields = `{
	respondedAt {
		time
		microsSinceEpoch
	}
	status
	request {
		customer {
			address
			nonce
			lastCustomerObjectID
		}
		customerObjectID
		venueOrderID
		newQuantity
		newPrice
		newDurationType
		newExpiresAt {
			time
			microsSinceEpoch
		}
		newExpiresIn {
			duration
			micros
		}
		nonce
		signature
	}
	replacementOrder {
		venueOrderID
		customer {
			address
			nonce
			lastCustomerObjectID
		}
		customerObjectID
		type
		side
		originalQuantity
		totalFillQuantity
		totalQuantity
		quantity
		volume
		baseAssetID
		quoteAssetID
		price
		averagePrice
		durationType
		createdAt {
			time
			microsSinceEpoch
		}
		modifiedAt {
			time
			microsSinceEpoch
		}
		expiresAt {
			time
			microsSinceEpoch
		}
		expiresIn {
			duration
			micros
		}
		status
		rfq
		requestNonce
		message
	}
	replacedOrder {
		venueOrderID
		customer {
			address
			nonce
			lastCustomerObjectID
		}
		customerObjectID
		type
		side
		originalQuantity
		totalFillQuantity
		totalQuantity
		quantity
		volume
		baseAssetID
		quoteAssetID
		price
		averagePrice
		durationType
		createdAt {
			time
			microsSinceEpoch
		}
		modifiedAt {
			time
			microsSinceEpoch
		}
		expiresAt {
			time
			microsSinceEpoch
		}
		expiresIn {
			duration
			micros
		}
		status
		rfq
		requestNonce
		message
	}
	error
}`

const TradeResponseFields = `{
	type
	tradeID
	createdAt {
		time
		microsSinceEpoch
	}
	baseAssetID
	quoteAssetID
	quantity
	volume
	price
	aggressor {
		remainingQty
		order {
			venueOrderID
			customer {
				address
			}
			customerObjectID
			type
			side
			quantity
			originalQuantity
			volume
			baseAssetID
			quoteAssetID
			price
			averagePrice
			durationType
			createdAt {
				time
				microsSinceEpoch
			}
			modifiedAt {
				time
				microsSinceEpoch
			}
			expiresAt {
				time
				microsSinceEpoch
			}
			status
		}
	}
	maker {
		remainingQty
		order {
			venueOrderID
			customer {
				address
			}
			customerObjectID
			type
			side
			quantity
			originalQuantity
			volume
			baseAssetID
			quoteAssetID
			price
			averagePrice
			durationType
			createdAt {
				time
				microsSinceEpoch
			}
			modifiedAt {
				time
				microsSinceEpoch
			}
			expiresAt {
				time
				microsSinceEpoch
			}
			status
		}
	}
}`

const CancelResponseFields = `{
	status
	respondedAt {
		time
		microsSinceEpoch
	}
	error
	request {
		customer {
			address
			nonce
			lastCustomerObjectID
		}
		customerObjectID
		venueOrderID
		requestedAt {
			time
			microsSinceEpoch
		}
		nonce
		signature
	}
}`

const OpenPositionResponseFields = `{
	status
	error
	respondedAt {
		time
		microsSinceEpoch
	}
	position {
		positionID
		customer {
			address
			nonce
			lastCustomerObjectID
		}
		customerObjectID
		status
		createdAt {
			time
			microsSinceEpoch
		}
		modifiedAt {
			time
			microsSinceEpoch
		}
		primaryComponents {
			positionID
			asset
			startingBalance
			currentBalance
		}
		secondaryComponents {
			positionID
			asset
			target
			balance
			bidirectional
			enabled
			costBasis
		}
		rateLimit {
			duration {
				duration
				micros
			}
			volume
		}
		currentPeriod {
			startedAt {
				time
				microsSinceEpoch
			}
			availableBalance
		}
	}
	request {
		customer {
			address
			nonce
			lastCustomerObjectID
		}
		customerObjectID
		durationType
		requestedAt {
			time
			microsSinceEpoch
		}
		expiresAt {
			time
			microsSinceEpoch
		}
		expiresIn {
			duration
			micros
		}
		nonce
		signature
	}
}`

const CancelAllOrdersResponseFields = `{
	respondedAt {
		time
		microsSinceEpoch
	}
	numCancelled
	error
	request {
		customer {
			address
			nonce
			lastCustomerObjectID
		}
		baseAssetID
		quoteAssetID
		requestedAt {
			time
			microsSinceEpoch
		}
		nonce
		signature
	}
}`

const OrderBookLevelUpdateFields = `{
	baseAssetID
	quoteAssetID
	time {
		time
		microsSinceEpoch
	}
	side
	price
	quantity
}`

const OrderBookSnapshotFields = `{
	count
	cursor
	time {
		time
		microsSinceEpoch
	}
	bids {
		depth
		price
		quantity
	}
	asks {
		depth
		price
		quantity
	}
}`

const ShortOrderBookSnapshotFields = `{
	bids {
		price
		quantity
	}
	asks {
		price
		quantity
	}
}`

const TradeFields = `{
	tradeID
	type
	baseAssetID
	quoteAssetID
	quantity
	volume
	price
	aggressor {
		remainingQty
		order {
			venueOrderID
			customer {
				address
				nonce
				lastCustomerObjectID
			}
			customerObjectID
			type
			side
			quantity
			originalQuantity
			totalQuantity
			volume
			baseAssetID
			quoteAssetID
			price
			averagePrice
			durationType
			createdAt {
				time
				microsSinceEpoch
			}
			modifiedAt {
				time
				microsSinceEpoch
			}
			expiresAt {
				time
				microsSinceEpoch
			}
			expiresIn {
				duration
				micros
			}
			status
			rfq
			requestNonce
		}
	}
	maker {
		remainingQty
		order {
			venueOrderID
			customer {
				address
				nonce
				lastCustomerObjectID
			}
			customerObjectID
			type
			side
			quantity
			originalQuantity
			totalQuantity
			volume
			baseAssetID
			quoteAssetID
			price
			averagePrice
			durationType
			createdAt {
				time
				microsSinceEpoch
			}
			modifiedAt {
				time
				microsSinceEpoch
			}
			expiresAt {
				time
				microsSinceEpoch
			}
			expiresIn {
				duration
				micros
			}
			status
			rfq
			requestNonce
		}
	}
}`

const OrderFields = `{
	venueOrderID
	customer {
		address
		nonce
		lastCustomerObjectID
	}
	customerObjectID
	type
	side
	quantity
	originalQuantity
	totalFillQuantity
	totalQuantity
	volume
	baseAssetID
	quoteAssetID
	price
	averagePrice
	durationType
	createdAt {
		time
		microsSinceEpoch
	}
	modifiedAt {
		time
		microsSinceEpoch
	}
	expiresAt {
		time
		microsSinceEpoch
	}
	expiresIn {
		duration
		micros
	}
	status
	rfq
	requestNonce
}`
