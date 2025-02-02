package gitlab

import (
	"devmetrics/internal/api/rest/handlers/vcs/shared"
)

type BaseRequest struct {
	ProjectID int `params:"id" validate:"required"`
}

type RepositoryRequest struct {
	BaseRequest
}

type CommitsRequest struct {
	BaseRequest
	shared.TimeRangeRequest
}

type PullRequestsRequest struct {
	BaseRequest
	shared.TimeRangeRequest
	Status string `query:"status" validate:"omitempty,oneof=all open closed merged" default:"all"`
}
