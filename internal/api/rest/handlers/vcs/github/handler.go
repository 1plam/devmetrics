package github

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"time"

	"devmetrics/internal/domain/vcs"
	service "devmetrics/internal/services/vcs"
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
	var req RepoRequest
	if err := c.ParamsParser(&req); err != nil {
		return h.errorResponse(c, fiber.StatusBadRequest, "invalid parameters", err)
	}

	if err := h.validator.Struct(req); err != nil {
		return h.errorResponse(c, fiber.StatusBadRequest, "validation failed", err)
	}

	repoPath := fmt.Sprintf("%s/%s", req.Owner, req.Name)
	repo, err := h.service.GetRepository(c.Context(), vcs.ProviderGitHub, repoPath)
	if err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(fiber.Map{
		"data": repo,
	})
}

func (h *Handler) GetCommits(c *fiber.Ctx) error {
	var req RepoRequest
	if err := c.ParamsParser(&req); err != nil {
		return h.errorResponse(c, fiber.StatusBadRequest, "invalid parameters", err)
	}

	if err := h.validator.Struct(req); err != nil {
		return h.errorResponse(c, fiber.StatusBadRequest, "validation failed", err)
	}

	since := c.Query("since")
	until := c.Query("until")

	var sinceTime, untilTime time.Time
	var err error

	if since != "" {
		sinceTime, err = time.Parse(time.RFC3339, since)
		if err != nil {
			return h.errorResponse(c, fiber.StatusBadRequest, "invalid since time", err)
		}
	}

	if until != "" {
		untilTime, err = time.Parse(time.RFC3339, until)
		if err != nil {
			return h.errorResponse(c, fiber.StatusBadRequest, "invalid until time", err)
		}
	}

	repoPath := fmt.Sprintf("%s/%s", req.Owner, req.Name)
	commits, err := h.service.GetCommits(c.Context(), vcs.ProviderGitHub, repoPath, sinceTime, untilTime)
	if err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(fiber.Map{
		"data": commits,
	})
}

func (h *Handler) GetPullRequests(c *fiber.Ctx) error {
	var req RepoRequest
	if err := c.ParamsParser(&req); err != nil {
		return h.errorResponse(c, fiber.StatusBadRequest, "invalid parameters", err)
	}

	if err := h.validator.Struct(req); err != nil {
		return h.errorResponse(c, fiber.StatusBadRequest, "validation failed", err)
	}

	since := c.Query("since")
	until := c.Query("until")
	status := c.Query("status")

	if status != "" && !isValidPRStatus(status) {
		return h.errorResponse(c, fiber.StatusBadRequest, "invalid status", fmt.Errorf("unsupported status: %s", status))
	}

	var sinceTime, untilTime time.Time
	var err error

	if since != "" {
		sinceTime, err = time.Parse(time.RFC3339, since)
		if err != nil {
			return h.errorResponse(c, fiber.StatusBadRequest, "invalid since time", err)
		}
	}

	if until != "" {
		untilTime, err = time.Parse(time.RFC3339, until)
		if err != nil {
			return h.errorResponse(c, fiber.StatusBadRequest, "invalid until time", err)
		}
	}

	repoPath := fmt.Sprintf("%s/%s", req.Owner, req.Name)
	prs, err := h.service.GetPullRequests(c.Context(), vcs.ProviderGitHub, repoPath, sinceTime, untilTime)
	if err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(fiber.Map{
		"data": prs,
	})
}
