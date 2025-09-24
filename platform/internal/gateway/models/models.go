package models

type APIResponse struct {
	Success bool      `json:"success"`
	Error   *APIError `json:"error,omitempty"`
	Data    any       `json:"data,omitempty"`
}

type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}
