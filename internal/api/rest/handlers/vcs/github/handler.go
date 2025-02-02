package github

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
		vcs.ProviderGitHub,
		fmt.Sprintf("%s/%s", req.Owner, req.Name),
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

	commits, total, err := h.service.GetCommits(
		ctx,
		vcs.ProviderGitHub,
		fmt.Sprintf("%s/%s", req.Owner, req.Name),
		req.GetSinceTime(),
		req.GetUntilTime(),
		req.GetOffset(),
		req.GetPerPage(),
	)
	if err != nil {
		return h.HandleError(c, err)
	}

	pagination := shared.NewPaginationMeta(req.GetPage(), req.GetPerPage(), total)
	return h.SendPaginatedResponse(c, commits, pagination)
}

func (h *Handler) GetPullRequests(c *fiber.Ctx) error {
	ctx, cancel := shared.NewTimeoutContext(c.Context(), shared.DefaultTimeout)
	defer cancel()

	req := new(PullRequestsRequest)
	if err := h.ParseAndValidate(c, req); err != nil {
		return err
	}

	prs, total, err := h.service.GetPullRequests(
		ctx,
		vcs.ProviderGitHub,
		fmt.Sprintf("%s/%s", req.Owner, req.Name),
		req.GetSinceTime(),
		req.GetUntilTime(),
		req.GetOffset(),
		req.GetPerPage(),
	)
	if err != nil {
		return h.HandleError(c, err)
	}

	pagination := shared.NewPaginationMeta(req.GetPage(), req.GetPerPage(), total)
	return h.SendPaginatedResponse(c, prs, pagination)
}
