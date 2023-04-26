package rest

import "math/big"

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

type GetCustomerSuccessResponse struct {
	Success bool                    `json:"success"`
	Result  *GetCustomerSuccessBody `json:"result"`
}

type GetCustomerSuccessBody struct {
	CustomerAddress      string     `json:"customerAddress"`
	Nonce                uint64     `json:"nonce"`
	CreatedAtMicros      uint64     `json:"createdAtMicros"`
	ModifiedAtMicros     uint64     `json:"modifiedAtMicros"`
	LastCustomerObjectID *string    `json:"lastCustomerObjectID"`
	Balances             []*Balance `json:"balances"`
}

type Balance struct {
	ServiceName        string `json:"serviceName"`
	ServiceDescription string `json:"serviceDescription"`
	ServiceID          string `json:"serviceID"`
	AssetID            string `json:"assetID"`
	Amount             string `json:"amount"`
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

type DepositInstructionsRequest struct {
	Asset  string `json:"asset"`
	Amount string `json:"amount"`
}

type DepositInstructionsSuccessResponse struct {
	Success bool                      `json:"success"`
	Result  *DepositERC20Instructions `json:"result"`
}

type DepositERC20Instructions struct {
	ChainID  int               `json:"chainID"`
	Contract string            `json:"contract"`
	Sender   string            `json:"sender"`
	Function string            `json:"function"`
	Args     *DepositERC20Args `json:"args"`
	Note     *string           `json:"note,omitempty"` // this is the only difference with the accounting model object
}

type DepositERC20Args struct {
	Asset          string   `json:"asset"`
	AmountContract *big.Int `json:"amount"`
	Nonce          *big.Int `json:"nonce"`
	Deadline       *big.Int `json:"deadline"`
	Signature      string   `json:"signature"`
}

type WithdrawalInstructionsRequest struct {
	Asset  string `json:"asset"`
	Amount string `json:"amount"`
}

type WithdrawalInstructionsSuccessResponse struct {
	Success bool                         `json:"success"`
	Result  *WithdrawalERC20Instructions `json:"result"`
}

type WithdrawalERC20Instructions struct {
	ChainID  int                  `json:"chainID"`
	Contract string               `json:"contract"`
	Sender   string               `json:"sender"`
	Function string               `json:"function"`
	Args     *WithdrawalERC20Args `json:"args"`
	Note     *string              `json:"note,omitempty"` // this is the only difference with the accounting model object
}

type WithdrawalERC20Args struct {
	Asset          string   `json:"asset"`
	AmountContract *big.Int `json:"amount"`
	Nonce          *big.Int `json:"nonce"`
	Deadline       *big.Int `json:"deadline"`
	Signature      string   `json:"signature"`
}

type AssetsSuccessResponse struct {
	Success bool     `json:"success"`
	Result  []*Asset `json:"result"`
}

type Asset struct {
	ChainID      uint64 `json:"chainID"`
	ChainName    string `json:"chainName"`
	AssetAddress string `json:"assetAddress"`

	RootSymbol     string `json:"rootSymbol"`
	AssetID        string `json:"assetID"`
	ParentSymbol   string `json:"parentSymbol"`
	FrontEndSymbol string `json:"frontEndSymbol"`

	Name         string `json:"name"`
	Ticker       string `json:"ticker"`
	FrontEndName string `json:"frontEndName"`

	UIDecimals       uint8 `json:"uiDecimals"`
	DatabaseDecimals uint8 `json:"databaseDecimals"`
	ContractDecimals uint8 `json:"contractDecimals"`

	CreatedAtMicros uint64 `json:"createdAtMicros"`
}
