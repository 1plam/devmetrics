package shared

import "time"

type TimeRangeRequest struct {
	Since *time.Time `query:"since" validate:"omitempty,ltefield=Until"`
	Until *time.Time `query:"until" validate:"omitempty"`
}

type Response struct {
	Data       interface{}     `json:"data,omitempty"`
	Error      *ErrorResponse  `json:"error,omitempty"`
	Pagination *PaginationMeta `json:"pagination,omitempty"`
}

type PaginatedResponse struct {
	Data       interface{}    `json:"data"`
	Pagination PaginationMeta `json:"pagination"`
}

type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}
