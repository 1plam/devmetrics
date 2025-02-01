package github

import "github.com/gofiber/fiber/v2"

func isValidPRStatus(status string) bool {
	return status == "open" || status == "closed" || status == "merged" || status == "all"
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
