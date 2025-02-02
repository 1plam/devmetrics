package github

import (
	"devmetrics/internal/api/rest/handlers/vcs/shared"
)

type RepositoryRequest struct {
	Owner string `params:"owner" validate:"required"`
	Name  string `params:"name" validate:"required"`
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
	Status string `query:"status" validate:"omitempty,oneof=open closed merged all"`
}
