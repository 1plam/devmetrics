package vcs

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"time"

	"devmetrics/internal/domain/vcs"
	svc "devmetrics/internal/services/vcs"
)

type Handler struct {
	service   *svc.Service
	validator *validator.Validate
}

func NewHandler(service *svc.Service) *Handler {
	return &Handler{
		service:   service,
		validator: validator.New(),
	}
}

func (h *Handler) GetRepository(c *fiber.Ctx) error {
	var req GetRepositoryRequest

	if err := c.ParamsParser(&req); err != nil {
		return h.errorResponse(c, fiber.StatusBadRequest, "invalid parameters", err)
	}

	if err := h.validator.Struct(req); err != nil {
		return h.errorResponse(c, fiber.StatusBadRequest, "validation failed", err)
	}

	// Combine owner and name into the repository identifier
	repoIdentifier := fmt.Sprintf("%s/%s", req.Owner, req.Name)

	repo, err := h.service.GetRepository(
		c.Context(),
		vcs.ProviderType(req.Provider),
		repoIdentifier,
	)
	if err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(fiber.Map{
		"data": repo,
	})
}

func (h *Handler) GetCommits(c *fiber.Ctx) error {
	var req TimeRangeRequest

	if err := parseTimeRangeRequest(c, &req); err != nil {
		return h.errorResponse(c, fiber.StatusBadRequest, "invalid request", err)
	}

	repoIdentifier := fmt.Sprintf("%s/%s", req.Owner, req.Name)

	commits, err := h.service.GetCommits(
		c.Context(),
		vcs.ProviderType(req.Provider),
		repoIdentifier,
		req.Since,
		req.Until,
	)
	if err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(fiber.Map{
		"data": commits,
	})
}

func (h *Handler) GetPullRequests(c *fiber.Ctx) error {
	var req TimeRangeRequest

	if err := parseTimeRangeRequest(c, &req); err != nil {
		return h.errorResponse(c, fiber.StatusBadRequest, "invalid request", err)
	}

	// Combine owner and name into repository identifier
	repoIdentifier := fmt.Sprintf("%s/%s", req.Owner, req.Name)

	prs, err := h.service.GetPullRequests(
		c.Context(),
		vcs.ProviderType(req.Provider),
		repoIdentifier,
		req.Since,
		req.Until,
	)
	if err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(fiber.Map{
		"data": prs,
	})
}

func parseTimeRangeRequest(c *fiber.Ctx, req *TimeRangeRequest) error {
	if err := c.ParamsParser(req); err != nil {
		return err
	}

	since := c.Query("since")
	until := c.Query("until")

	if since != "" {
		t, err := time.Parse(time.RFC3339, since)
		if err != nil {
			return err
		}
		req.Since = t
	}

	if until != "" {
		t, err := time.Parse(time.RFC3339, until)
		if err != nil {
			return err
		}
		req.Until = t
	}

	return nil
}

func (h *Handler) errorResponse(c *fiber.Ctx, status int, message string, err error) error {
	return c.Status(status).JSON(fiber.Map{
		"error":   message,
		"details": err.Error(),
	})
}

func (h *Handler) handleError(c *fiber.Ctx, err error) error {
	return h.errorResponse(c, fiber.StatusInternalServerError, "internal server error", err)
}
