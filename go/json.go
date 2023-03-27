package main

type ErrorResponse struct {
	Code    int
	Error   string
	Status  string
	Success bool
}

type GetAPIKeySessionSuccessResponse struct {
	Result  GetAPIKeySessionSuccessBody
	Success bool
}

type GetAPIKeySessionSuccessBody struct {
	APIKey           string
	CreatedAtMicros  uint64
	ExpiresAtMicros  uint64
	Label            string
	ModifiedAtMicros uint64
	Nonce            uint64
	Status           string
}

type SubmitOrderRequest struct {
	CustomerObjectID *string `json:"customerObjectID,omitempty"`
	Type             string  `json:"type"`
	Side             string  `json:"side"`
	QuantityStr      string  `json:"quantity"`
	BaseAssetID      string  `json:"baseAssetID"`
	QuoteAssetID     string  `json:"quoteAssetID"`
	PriceStr         *string `json:"price,omitempty"`
	DurationType     *string `json:"durationType,omitempty"`
	ExpiresAtMicros  *uint64 `json:"expiresAtMicros,omitempty"`
	ExpiresIn        *string `json:"expiresIn,omitempty"`
}

type SubmitOrderSuccessResponse struct {
	Result  SubmitOrderSuccessBody
	Success bool
}

type SubmitOrderSuccessBody struct {
	CustomerObjectID  *string
	ExpiresAtMicros   *uint64
	OriginalQuantity  string
	Quantity          string
	RespondedAtMicros uint64
	Status            string
	VenueOrderID      string
}
