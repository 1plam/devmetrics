package gitlab

import (
	"context"
	"devmetrics/internal/domain/vcs"
	service "devmetrics/internal/services/vcs"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"time"
)

type Handler struct {
	service   *service.Service
	validator *validator.Validate
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{
		service:   service,
		validator: validator.New(),
	}
}

func (h *Handler) GetRepository(c *fiber.Ctx) error {
	params, err := h.parseRequestParams(c)
	if err != nil {
		return h.errorResponse(c, fiber.StatusBadRequest, "invalid request parameters", err)
	}

	repo, err := h.service.GetRepository(c.Context(), vcs.ProviderGitLab, fmt.Sprint(params.ProjectID))
	if err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(repo)
}

func (h *Handler) GetCommits(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.Context(), 30*time.Second)
	defer cancel()

	params, err := h.parseRequestParams(c)
	if err != nil {
		return h.errorResponse(c, fiber.StatusBadRequest, "invalid request parameters", err)
	}

	commits, err := h.service.GetCommits(ctx, vcs.ProviderGitLab, fmt.Sprint(params.ProjectID), params.Since, params.Until)
	if err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(commits)
}

func (h *Handler) GetPullRequests(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.Context(), 30*time.Second)
	defer cancel()

	params, err := h.parseRequestParams(c)
	if err != nil {
		return h.errorResponse(c, fiber.StatusBadRequest, "invalid request parameters", err)
	}

	prs, err := h.service.GetPullRequests(ctx, vcs.ProviderGitLab, fmt.Sprint(params.ProjectID), params.Since, params.Until)
	if err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(prs)
}
