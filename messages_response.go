package main

// Success Response
type SuccessResponse struct {
	StatusCode int         `json:"status_code"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data,omitempty"`
}

// Error Response
type ErrorResponse struct {
	StatusCode int    `json:"status_code"`
	Error      string `json:"error"`
}


