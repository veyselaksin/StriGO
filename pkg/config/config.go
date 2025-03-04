package config

import (
	"fmt"
	"time"

	"github.com/veyselaksin/strigo/pkg/duration"
)

// Strategy represents the rate limiting strategy
type Strategy string

const (
	// TokenBucket is the default strategy
	TokenBucket Strategy = "token_bucket"
	// LeakyBucket strategy
	LeakyBucket Strategy = "leaky_bucket"
	// FixedWindow strategy
	FixedWindow Strategy = "fixed_window"
	// SlidingWindow strategy
	SlidingWindow Strategy = "sliding_window"
)

// Config represents the rate limiter configuration
type Config struct {
	// Strategy defines the rate limiting algorithm to use
	Strategy Strategy `json:"strategy"`

	// Period defines the time window for rate limiting
	Period duration.Period `json:"period"`

	// Limit defines the maximum number of requests allowed per period
	Limit int64 `json:"limit"`

	// Prefix is used to create unique keys in the storage backend
	Prefix string `json:"prefix"`

	// Backend configuration
	BackendConfig BackendConfig `json:"backend"`
}

// BackendConfig holds the storage backend configuration
type BackendConfig struct {
	// Type specifies the backend type (redis, memcached)
	Type string `json:"type"`

	// Address is the connection string for the backend
	Address string `json:"address"`

	// Password for authentication (optional)
	Password string `json:"password,omitempty"`

	// Database number (for Redis)
	Database int `json:"database,omitempty"`
}

// Validate checks if the configuration is valid
func (c *Config) Validate() error {
	if c.Limit <= 0 {
		return fmt.Errorf("limit must be positive")
	}

	if !duration.IsValidPeriod(c.Period) {
		return fmt.Errorf("invalid period: %s", c.Period)
	}

	if c.Prefix == "" {
		return fmt.Errorf("prefix cannot be empty")
	}

	if err := c.validateStrategy(); err != nil {
		return err
	}

	if err := c.validateBackend(); err != nil {
		return err
	}

	return nil
}

// validateStrategy checks if the selected strategy is valid
func (c *Config) validateStrategy() error {
	switch c.Strategy {
	case TokenBucket, LeakyBucket, FixedWindow, SlidingWindow:
		return nil
	case "":
		c.Strategy = TokenBucket // Set default strategy
		return nil
	default:
		return fmt.Errorf("invalid strategy: %s", c.Strategy)
	}
}

// validateBackend checks if the backend configuration is valid
func (c *Config) validateBackend() error {
	if c.BackendConfig.Type == "" {
		return fmt.Errorf("backend type cannot be empty")
	}

	if c.BackendConfig.Address == "" {
		return fmt.Errorf("backend address cannot be empty")
	}

	return nil
}

// GetDuration returns the duration for the rate limit
func (c *Config) GetDuration() time.Duration {
	return c.Period.ToDuration()
}

// NewDefaultConfig creates a new Config with default values
func NewDefaultConfig() *Config {
	return &Config{
		Strategy: TokenBucket,
		Period:   duration.MINUTELY,
		Limit:    100,
		Prefix:   "strigo",
		BackendConfig: BackendConfig{
			Type:    "redis",
			Address: "localhost:6379",
		},
	}
}
