package models

// APIErrors response model
type APIErrors struct {
	Errors []APIErrorElement `json:"errors"`
}

// APIResponse (not needed here because direct output)
type APIResponse struct {
	Errors  []APIErrorElement `json:"errors"`
	Payload interface{}       `json:"payload"`
}

// APIErrorElement single error element of an error response
type APIErrorElement struct {
	ErrorCode    string `json:"errorCode"`
	ErrorMessage string `json:"errorMessage"`
}
