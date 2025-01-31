package app

import (
	"context"
	"fmt"
	"time"

	"devmetrics/internal/api/rest/server"
	"devmetrics/internal/config"
	"devmetrics/internal/services/vcs"
)

// Application encapsulates the application and its dependencies
type Application struct {
	cfg    *config.Config
	server *server.Server
	vcs    *vcs.Service
}

// NewApplication creates a new application instance
func NewApplication(
	cfg *config.Config,
	server *server.Server,
	vcs *vcs.Service,
) *Application {
	return &Application{
		cfg:    cfg,
		server: server,
		vcs:    vcs,
	}
}

// Start initializes and starts all application components
func (a *Application) Start(ctx context.Context) error {
	if err := a.server.Start(); err != nil {
		return fmt.Errorf("server error: %w", err)
	}

	return nil
}

// Shutdown gracefully stops all application components
func (a *Application) Shutdown(ctx context.Context) error {
	shutdownCtx, cancel := context.WithTimeout(
		ctx,
		time.Duration(a.cfg.Server.ShutdownTimeout)*time.Second,
	)

	defer cancel()
	if err := a.server.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("graceful shutdown failed: %w", err)
	}

	return nil
}
