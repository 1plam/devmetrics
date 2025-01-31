package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
	"strconv"
	"time"
)

type Config struct {
	Environment string
	Server      ServerConfig
	VCS         VCSConfig
	Logger      LoggerConfig
}

type ServerConfig struct {
	Host            string
	Port            string
	ShutdownTimeout int
}

type VCSConfig struct {
	GitHub    GitHubConfig
	GitLab    GitLabConfig
	BitBucket BitBucketConfig
}

type GitHubConfig struct {
	Enabled    bool
	Token      string
	BaseURL    string
	APIVersion string
	MaxPages   int
	PageSize   int
	TimeoutSec int
	RateLimit  int
	RetryCount int
	RetryDelay int
}

type GitLabConfig struct {
	Enabled    bool
	Token      string
	BaseURL    string
	MaxPages   int
	PageSize   int
	TimeoutSec int
}

type BitBucketConfig struct {
	Enabled     bool
	Username    string
	AppPassword string
	BaseURL     string
	MaxPages    int
	PageSize    int
	TimeoutSec  int
}

type LoggerConfig struct {
	Level      string
	Format     string
	Output     string
	TimeFormat string
}

// getEnvWithDefault retrieves an environment variable with a fallback default value
func getEnvWithDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvBoolWithDefault retrieves a boolean environment variable with a fallback default value
func getEnvBoolWithDefault(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		boolValue, err := strconv.ParseBool(value)
		if err == nil {
			return boolValue
		}
	}
	return defaultValue
}

// getEnvIntWithDefault retrieves an integer environment variable with a fallback default value
func getEnvIntWithDefault(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		intValue, err := strconv.Atoi(value)
		if err == nil {
			return intValue
		}
	}
	return defaultValue
}

func NewConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		fmt.Printf("Warning: Error loading .env file: %v\n", err)
	}

	config := &Config{
		Environment: getEnvWithDefault("ENVIRONMENT", "development"),
		Server: ServerConfig{
			Host:            getEnvWithDefault("SERVER_HOST", "0.0.0.0"),
			Port:            getEnvWithDefault("SERVER_PORT", "8080"),
			ShutdownTimeout: getEnvIntWithDefault("SERVER_SHUTDOWN_TIMEOUT", 5),
		},
		VCS: VCSConfig{
			GitHub: GitHubConfig{
				Enabled:    getEnvBoolWithDefault("VCS_GITHUB_ENABLED", false),
				Token:      os.Getenv("VCS_GITHUB_TOKEN"),
				BaseURL:    getEnvWithDefault("VCS_GITHUB_BASE_URL", "https://api.github.com"),
				APIVersion: getEnvWithDefault("VCS_GITHUB_API_VERSION", "2022-11-28"),
				MaxPages:   getEnvIntWithDefault("VCS_GITHUB_MAX_PAGES", 100),
				PageSize:   getEnvIntWithDefault("VCS_GITHUB_PAGE_SIZE", 100),
				TimeoutSec: getEnvIntWithDefault("VCS_GITHUB_TIMEOUT_SEC", 30),
				RateLimit:  getEnvIntWithDefault("VCS_GITHUB_RATE_LIMIT", 5000),
				RetryCount: getEnvIntWithDefault("VCS_GITHUB_RETRY_COUNT", 3),
				RetryDelay: getEnvIntWithDefault("VCS_GITHUB_RETRY_DELAY", 1),
			},
			GitLab: GitLabConfig{
				Enabled:    getEnvBoolWithDefault("VCS_GITLAB_ENABLED", false),
				Token:      os.Getenv("VCS_GITLAB_TOKEN"),
				BaseURL:    getEnvWithDefault("VCS_GITLAB_BASE_URL", "https://gitlab.com/api/v4"),
				MaxPages:   getEnvIntWithDefault("VCS_GITLAB_MAX_PAGES", 100),
				PageSize:   getEnvIntWithDefault("VCS_GITLAB_PAGE_SIZE", 100),
				TimeoutSec: getEnvIntWithDefault("VCS_GITLAB_TIMEOUT_SEC", 30),
			},
			BitBucket: BitBucketConfig{
				Enabled:     getEnvBoolWithDefault("VCS_BITBUCKET_ENABLED", false),
				Username:    os.Getenv("VCS_BITBUCKET_USERNAME"),
				AppPassword: os.Getenv("VCS_BITBUCKET_APP_PASSWORD"),
				BaseURL:     getEnvWithDefault("VCS_BITBUCKET_BASE_URL", "https://api.bitbucket.org/2.0"),
				MaxPages:    getEnvIntWithDefault("VCS_BITBUCKET_MAX_PAGES", 100),
				PageSize:    getEnvIntWithDefault("VCS_BITBUCKET_PAGE_SIZE", 100),
				TimeoutSec:  getEnvIntWithDefault("VCS_BITBUCKET_TIMEOUT_SEC", 30),
			},
		},
		Logger: LoggerConfig{
			Level:      getEnvWithDefault("LOGGER_LEVEL", "info"),
			Format:     getEnvWithDefault("LOGGER_FORMAT", "json"),
			Output:     getEnvWithDefault("LOGGER_OUTPUT", "stdout"),
			TimeFormat: getEnvWithDefault("LOGGER_TIME_FORMAT", time.RFC3339),
		},
	}

	if err := validateConfig(config); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return config, nil
}

func validateConfig(cfg *Config) error {
	if cfg.VCS.GitHub.Enabled && cfg.VCS.GitHub.Token == "" {
		return fmt.Errorf("GitHub token is required when GitHub is enabled")
	}

	if cfg.VCS.GitLab.Enabled && cfg.VCS.GitLab.Token == "" {
		return fmt.Errorf("GitLab token is required when GitLab is enabled")
	}

	if cfg.VCS.BitBucket.Enabled {
		if cfg.VCS.BitBucket.Username == "" {
			return fmt.Errorf("BitBucket username is required when BitBucket is enabled")
		}
		if cfg.VCS.BitBucket.AppPassword == "" {
			return fmt.Errorf("BitBucket app password is required when BitBucket is enabled")
		}
	}

	return nil
}
