package shared

import (
	"context"
	"time"
)

// DefaultTimeout is the default duration for context timeout
const DefaultTimeout = 30 * time.Second

// NewTimeoutContext creates a new context with timeout
func NewTimeoutContext(parent context.Context, timeout time.Duration) (context.Context, context.CancelFunc) {
	if timeout <= 0 {
		timeout = DefaultTimeout
	}
	return context.WithTimeout(parent, timeout)
}
