package shared

import "time"

// TimeRangeRequest Base request types that can be embedded in provider-specific requests
type TimeRangeRequest struct {
	Since *time.Time `query:"since" validate:"omitempty,ltefield=Until"`
	Until *time.Time `query:"until" validate:"omitempty"`
}

// Response Common response types
type Response struct {
	Data  interface{}    `json:"data,omitempty"`
	Error *ErrorResponse `json:"error,omitempty"`
}

type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}
