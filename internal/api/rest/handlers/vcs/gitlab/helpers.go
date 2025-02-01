package gitlab

import (
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"log"
	"time"
)

func (h *Handler) parseTimeParams(since, until string) (time.Time, time.Time, error) {
	var sinceTime, untilTime time.Time
	var err error

	if since != "" {
		sinceTime, err = time.Parse(time.RFC3339, since)
		if err != nil {
			return time.Time{}, time.Time{}, err
		}
	} else {
		sinceTime = time.Now().AddDate(0, -1, 0) // Default to 1 month ago
	}

	if until != "" {
		untilTime, err = time.Parse(time.RFC3339, until)
		if err != nil {
			return time.Time{}, time.Time{}, err
		}
	} else {
		untilTime = time.Now()
	}

	return sinceTime, untilTime, nil
}

func (h *Handler) parseRequestParams(c *fiber.Ctx) (*RequestParams, error) {
	params := new(RequestParams)

	// Get project ID from the URL parameter
	projectID, err := c.ParamsInt("id")
	if err != nil {
		return nil, fmt.Errorf("invalid project ID: %w", err)
	}
	params.ProjectID = projectID

	// Parse time parameters
	since := c.Query("since")
	until := c.Query("until")
	sinceTime, untilTime, err := h.parseTimeParams(since, until)
	if err != nil {
		return nil, err
	}

	params.Since = sinceTime
	params.Until = untilTime
	params.Status = c.Query("status", "all")

	if err := h.validator.Struct(params); err != nil {
		return nil, err
	}

	return params, nil
}

func (h *Handler) errorResponse(c *fiber.Ctx, status int, message string, err error) error {
	var validationErrors validator.ValidationErrors
	if errors.As(err, &validationErrors) {
		return c.Status(status).JSON(fiber.Map{
			"error":   message,
			"details": validationErrors.Error(),
		})
	}

	return c.Status(status).JSON(fiber.Map{
		"error":   message,
		"details": err.Error(),
	})
}

func (h *Handler) handleError(c *fiber.Ctx, err error) error {
	log.Printf("Handling error: %v", err)
	return h.errorResponse(c, fiber.StatusInternalServerError, "internal server error", err)
}
