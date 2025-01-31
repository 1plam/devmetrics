package vcs

import (
	"devmetrics/internal/adapters/vcs/github"
	"devmetrics/internal/config"
	"devmetrics/internal/domain/vcs"
	"errors"
	"fmt"
)

var (
	// ErrProviderNotImplemented is returned when a VCS provider is not yet implemented
	ErrProviderNotImplemented = errors.New("provider not implemented")
)

type Factory struct {
	vcsConfig config.VCSConfig
}

func NewFactory(cfg config.VCSConfig) *Factory {
	return &Factory{
		vcsConfig: cfg,
	}
}

func (f *Factory) CreateProviders() (map[vcs.ProviderType]vcs.Provider, error) {
	providers := make(map[vcs.ProviderType]vcs.Provider)

	if err := f.createGitHubProvider(providers); err != nil {
		return nil, err
	}

	return providers, nil
}

func (f *Factory) createGitHubProvider(providers map[vcs.ProviderType]vcs.Provider) error {
	if !f.vcsConfig.GitHub.Enabled {
		return nil
	}

	provider, err := github.NewAdapter(f.vcsConfig.GitHub)
	if err != nil {
		return fmt.Errorf("failed to create GitHub provider: %w", err)
	}

	providers[vcs.ProviderGitHub] = provider
	return nil
}

func (f *Factory) createGitLabProvider(providers map[vcs.ProviderType]vcs.Provider) error {
	// TODO(#123): Implement GitLab provider support
	// - Add GitLab API client
	// - Implement repository listing
	// - Add pagination support
	return ErrProviderNotImplemented
}

func (f *Factory) createBitBucketProvider(providers map[vcs.ProviderType]vcs.Provider) error {
	// TODO(#124): Implement BitBucket provider support
	// - Add BitBucket API client
	// - Implement repository listing
	// - Add pagination support
	return ErrProviderNotImplemented
}
