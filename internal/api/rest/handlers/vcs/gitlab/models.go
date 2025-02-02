package gitlab

import (
	"devmetrics/internal/api/rest/handlers/vcs/shared"
)

type RepositoryRequest struct {
	ProjectID int `params:"id" validate:"required"`
}

type CommitsRequest struct {
	RepositoryRequest
	shared.TimeRangeRequest
	shared.PaginationRequest
}

type PullRequestsRequest struct {
	RepositoryRequest
	shared.TimeRangeRequest
	shared.PaginationRequest
	Status string `query:"status" validate:"omitempty,oneof=all open closed merged" default:"all"`
}
