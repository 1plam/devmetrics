package shared

import (
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"log"
	"time"
)

type BaseHandler struct {
	Validator *validator.Validate
}

func NewBaseHandler() BaseHandler {
	return BaseHandler{
		Validator: validator.New(),
	}
}

// ParseAndValidate parses and validates request data
func (h *BaseHandler) ParseAndValidate(c *fiber.Ctx, req interface{}) error {
	if err := c.ParamsParser(req); err != nil {
		return h.ErrorResponse(c, fiber.StatusBadRequest, "invalid_parameters", "Failed to parse parameters", err.Error())
	}

	if err := c.QueryParser(req); err != nil {
		return h.ErrorResponse(c, fiber.StatusBadRequest, "invalid_query", "Failed to parse query parameters", err.Error())
	}

	if err := h.Validator.Struct(req); err != nil {
		return h.ErrorResponse(c, fiber.StatusBadRequest, "validation_failed", "Request validation failed", err.Error())
	}

	return nil
}

// SendResponse sends a successful response
func (h *BaseHandler) SendResponse(c *fiber.Ctx, data interface{}) error {
	return c.JSON(Response{
		Data: data,
	})
}

func (h *BaseHandler) SendPaginatedResponse(c *fiber.Ctx, data interface{}, pagination PaginationMeta) error {
	return c.JSON(Response{
		Data:       data,
		Pagination: &pagination,
	})
}

// ErrorResponse creates and sends an error response
func (h *BaseHandler) ErrorResponse(c *fiber.Ctx, status int, code, message, details string) error {
	log.Printf("Error response: %s - %s - %s", code, message, details)
	return c.Status(status).JSON(Response{
		Error: &ErrorResponse{
			Code:    code,
			Message: message,
			Details: details,
		},
	})
}

// HandleError handles different types of errors and returns appropriate responses
func (h *BaseHandler) HandleError(c *fiber.Ctx, err error) error {
	var validationErrors validator.ValidationErrors
	if errors.As(err, &validationErrors) {
		return h.ErrorResponse(c, fiber.StatusBadRequest, "validation_failed", "Validation failed", validationErrors.Error())
	}

	log.Printf("Internal error: %v", err)
	return h.ErrorResponse(c, fiber.StatusInternalServerError, "internal_error", "Internal server error", err.Error())
}

// GetSinceTime Time helper methods
func (r *TimeRangeRequest) GetSinceTime() time.Time {
	if r.Since != nil {
		return *r.Since
	}
	return time.Now().AddDate(0, -1, 0) // Default to 1 month ago
}

func (r *TimeRangeRequest) GetUntilTime() time.Time {
	if r.Until != nil {
		return *r.Until
	}
	return time.Now()
}
