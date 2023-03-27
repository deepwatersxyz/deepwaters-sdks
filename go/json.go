package main

type ErrorResponse struct {
	Code    int
	Error   string
	Status  string
	Success bool
}

type GetAPIKeySessionSuccessResponse struct {
	Success bool
	Result  *GetAPIKeySessionSuccessBody
}

type GetAPIKeySessionSuccessBody struct {
	APIKey           string  `json:"APIKey"`
	Nonce            uint64  `json:"nonce"`
	Label            *string `json:"label"`
	Status           string  `json:"status"`
	CreatedAtMicros  uint64  `json:"createdAtMicros"`
	ModifiedAtMicros uint64  `json:"modifiedAtMicros"`
	ExpiresAtMicros  *uint64 `json:"expiresAtMicros"`
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
	Success bool
	Result  *SubmitOrderSuccessBody
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
