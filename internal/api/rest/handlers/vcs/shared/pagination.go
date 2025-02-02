package shared

import (
	"math"
)

const (
	DefaultPage    = 1
	DefaultPerPage = 30
	MaxPerPage     = 100
)

// PaginationRequest holds pagination parameters
type PaginationRequest struct {
	Page    int `query:"page" validate:"omitempty,min=1"`
	PerPage int `query:"per_page" validate:"omitempty,min=1,max=100"`
}

// GetPage returns the current page number with default handling
func (p *PaginationRequest) GetPage() int {
	if p.Page < 1 {
		return DefaultPage
	}
	return p.Page
}

// GetPerPage returns the items per page with default handling
func (p *PaginationRequest) GetPerPage() int {
	if p.PerPage < 1 {
		return DefaultPerPage
	}
	if p.PerPage > MaxPerPage {
		return MaxPerPage
	}
	return p.PerPage
}

// GetOffset calculates the offset for database queries
func (p *PaginationRequest) GetOffset() int {
	return (p.GetPage() - 1) * p.GetPerPage()
}

// PaginationMeta holds metadata about the pagination
type PaginationMeta struct {
	CurrentPage int   `json:"current_page"`
	PerPage     int   `json:"per_page"`
	TotalItems  int64 `json:"total_items"`
	TotalPages  int   `json:"total_pages"`
	HasMore     bool  `json:"has_more"`
}

// NewPaginationMeta creates a new PaginationMeta instance
func NewPaginationMeta(page, perPage int, totalItems int64) PaginationMeta {
	totalPages := int(math.Ceil(float64(totalItems) / float64(perPage)))

	return PaginationMeta{
		CurrentPage: page,
		PerPage:     perPage,
		TotalItems:  totalItems,
		TotalPages:  totalPages,
		HasMore:     page < totalPages,
	}
}
