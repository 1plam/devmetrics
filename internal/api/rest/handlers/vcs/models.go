package vcs

import (
	"time"
)

// GetRepositoryRequest Request models
type GetRepositoryRequest struct {
	Provider string `params:"provider" validate:"required,oneof=github gitlab bitbucket"`
	Owner    string `params:"owner" validate:"required"`
	Name     string `params:"name" validate:"required"`
}

type GetCommitsRequest struct {
	Provider string    `params:"provider" validate:"required,oneof=github gitlab bitbucket"`
	Owner    string    `params:"owner" validate:"required"`
	Name     string    `params:"name" validate:"required"`
	Since    time.Time `query:"since" validate:"required"`
	Until    time.Time `query:"until" validate:"required,gtfield=Since"`
}

type GetPullRequestsRequest struct {
	Provider string    `params:"provider" validate:"required,oneof=github gitlab bitbucket"`
	Owner    string    `params:"owner" validate:"required"`
	Name     string    `params:"name" validate:"required"`
	Since    time.Time `query:"since" validate:"required"`
	Until    time.Time `query:"until" validate:"required,gtfield=Since"`
	Status   string    `query:"status" validate:"omitempty,oneof=open closed merged"`
}

type TimeRangeRequest struct {
	Provider string    `params:"provider" validate:"required,oneof=github gitlab bitbucket"`
	Owner    string    `params:"owner" validate:"required"`
	Name     string    `params:"name" validate:"required"`
	Since    time.Time `query:"since" validate:"required"`
	Until    time.Time `query:"until" validate:"required,gtfield=Since"`
}

// ErrorResponse Response models
type ErrorResponse struct {
	Error     string            `json:"error"`
	Code      string            `json:"code,omitempty"`
	Details   map[string]string `json:"details,omitempty"`
	RequestID string            `json:"request_id,omitempty"`
}

type SuccessResponse struct {
	Data      interface{} `json:"data"`
	Meta      *MetaData   `json:"meta,omitempty"`
	RequestID string      `json:"request_id,omitempty"`
}

type MetaData struct {
	Page       int `json:"page,omitempty"`
	PerPage    int `json:"per_page,omitempty"`
	TotalItems int `json:"total_items,omitempty"`
	TotalPages int `json:"total_pages,omitempty"`
}
