package vcs

import (
	"context"
	"fmt"
	"time"

	"devmetrics/internal/domain/vcs"
)

type Service struct {
	providers map[vcs.ProviderType]vcs.Provider
}

func NewService(providers map[vcs.ProviderType]vcs.Provider) *Service {
	return &Service{
		providers: providers,
	}
}

func (s *Service) GetRepository(ctx context.Context, providerType vcs.ProviderType, repo string) (*vcs.Repository, error) {
	provider, ok := s.providers[providerType]
	if !ok {
		return nil, fmt.Errorf("unsupported VCS provider: %s", providerType)
	}

	return provider.GetRepository(ctx, repo)
}

func (s *Service) GetCommits(ctx context.Context, providerType vcs.ProviderType, repo string, since, until time.Time) ([]vcs.Commit, error) {
	provider, ok := s.providers[providerType]
	if !ok {
		return nil, fmt.Errorf("unsupported VCS provider: %s", providerType)
	}

	return provider.GetCommits(ctx, repo, since, until)
}

func (s *Service) GetPullRequests(ctx context.Context, providerType vcs.ProviderType, repo string, since, until time.Time) ([]vcs.PullRequest, error) {
	provider, ok := s.providers[providerType]
	if !ok {
		return nil, fmt.Errorf("unsupported VCS provider: %s", providerType)
	}

	return provider.GetPullRequests(ctx, repo, since, until)
}
