package gitlab

import (
	"devmetrics/internal/api/rest/handlers/vcs/shared"
	"devmetrics/internal/domain/vcs"
	service "devmetrics/internal/services/vcs"
	"fmt"
	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	service *service.Service
	shared.BaseHandler
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{
		service:     service,
		BaseHandler: shared.NewBaseHandler(),
	}
}

func (h *Handler) GetRepository(c *fiber.Ctx) error {
	req := new(RepositoryRequest)
	if err := h.ParseAndValidate(c, req); err != nil {
		return err
	}

	repo, err := h.service.GetRepository(
		c.Context(),
		vcs.ProviderGitLab,
		fmt.Sprint(req.ProjectID),
	)
	if err != nil {
		return h.HandleError(c, err)
	}

	return h.SendResponse(c, repo)
}

func (h *Handler) GetCommits(c *fiber.Ctx) error {
	ctx, cancel := shared.NewTimeoutContext(c.Context(), shared.DefaultTimeout)
	defer cancel()

	req := new(CommitsRequest)
	if err := h.ParseAndValidate(c, req); err != nil {
		return err
	}

	commits, err := h.service.GetCommits(
		ctx,
		vcs.ProviderGitLab,
		fmt.Sprint(req.ProjectID),
		req.GetSinceTime(),
		req.GetUntilTime(),
	)
	if err != nil {
		return h.HandleError(c, err)
	}

	return h.SendResponse(c, commits)
}

func (h *Handler) GetPullRequests(c *fiber.Ctx) error {
	ctx, cancel := shared.NewTimeoutContext(c.Context(), shared.DefaultTimeout)
	defer cancel()

	req := new(PullRequestsRequest)
	if err := h.ParseAndValidate(c, req); err != nil {
		return err
	}

	prs, err := h.service.GetPullRequests(
		ctx,
		vcs.ProviderGitLab,
		fmt.Sprint(req.ProjectID),
		req.GetSinceTime(),
		req.GetUntilTime(),
	)
	if err != nil {
		return h.HandleError(c, err)
	}

	return h.SendResponse(c, prs)
}
