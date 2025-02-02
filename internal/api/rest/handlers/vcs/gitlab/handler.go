package gitlab

import (
	"devmetrics/internal/api/rest/handlers/vcs/shared"
	domain "devmetrics/internal/domain/vcs"
	service "devmetrics/internal/services/vcs"
	"fmt"
	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	Service     *service.Service
	BaseHandler shared.BaseHandler
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{
		Service:     service,
		BaseHandler: shared.NewBaseHandler(),
	}
}

func (h *Handler) GetRepository(c *fiber.Ctx) error {
	req := new(RepositoryRequest)
	if err := h.BaseHandler.ParseAndValidate(c, req); err != nil {
		return err
	}

	repo, err := h.Service.GetRepository(
		c.Context(),
		domain.ProviderGitLab,
		fmt.Sprint(req.ProjectID),
	)
	if err != nil {
		return h.BaseHandler.HandleError(c, err)
	}

	return h.BaseHandler.SendResponse(c, repo)
}

func (h *Handler) GetCommits(c *fiber.Ctx) error {
	ctx, cancel := shared.NewTimeoutContext(c.Context(), shared.DefaultTimeout)
	defer cancel()

	req := new(CommitsRequest)
	if err := h.BaseHandler.ParseAndValidate(c, req); err != nil {
		return err
	}

	commits, total, err := h.Service.GetCommits(
		ctx,
		domain.ProviderGitLab,
		fmt.Sprint(req.ProjectID),
		req.GetSinceTime(),
		req.GetUntilTime(),
		req.GetOffset(),
		req.GetPerPage(),
	)
	if err != nil {
		return h.BaseHandler.HandleError(c, err)
	}

	pagination := shared.NewPaginationMeta(req.GetPage(), req.GetPerPage(), total)
	return h.BaseHandler.SendPaginatedResponse(c, commits, pagination)
}

func (h *Handler) GetPullRequests(c *fiber.Ctx) error {
	ctx, cancel := shared.NewTimeoutContext(c.Context(), shared.DefaultTimeout)
	defer cancel()

	req := new(PullRequestsRequest)
	if err := h.BaseHandler.ParseAndValidate(c, req); err != nil {
		return err
	}

	prs, total, err := h.Service.GetPullRequests(
		ctx,
		domain.ProviderGitLab,
		fmt.Sprint(req.ProjectID),
		req.GetSinceTime(),
		req.GetUntilTime(),
		req.GetOffset(),
		req.GetPerPage(),
	)
	if err != nil {
		return h.BaseHandler.HandleError(c, err)
	}

	pagination := shared.NewPaginationMeta(req.GetPage(), req.GetPerPage(), total)
	return h.BaseHandler.SendPaginatedResponse(c, prs, pagination)
}
