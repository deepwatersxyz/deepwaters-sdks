package main

type ErrorResponse struct {
	Code    uint16
	Error   string
	Status  string
	Success bool
}

type GetAPIKeySessionSuccessResponse struct {
	Result  GetAPIKeySessionInner
	Success bool
}

type GetAPIKeySessionInner struct {
	APIKey           string
	CreatedAtMicros  uint64
	ExpiresAtMicros  uint64
	Label            string
	ModifiedAtMicros uint64
	Nonce            uint64
	Status           string
}