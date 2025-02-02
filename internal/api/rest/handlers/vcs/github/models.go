package github

import (
	"devmetrics/internal/api/rest/handlers/vcs/shared"
)

type RepoRequest struct {
	Owner string `params:"owner" validate:"required"`
	Name  string `params:"name" validate:"required"`
}

type CommitsRequest struct {
	RepoRequest
	shared.TimeRangeRequest
}

type PullRequestsRequest struct {
	RepoRequest
	shared.TimeRangeRequest
	Status string `query:"status" validate:"omitempty,oneof=open closed merged all"`
}
