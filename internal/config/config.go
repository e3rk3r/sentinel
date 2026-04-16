// Package config provides configuration loading and validation for Sentinel.
package config

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// Config holds the top-level application configuration.
type Config struct {
	Server  ServerConfig  `yaml:"server"`
	Watcher WatcherConfig `yaml:"watcher"`
	Log     LogConfig     `yaml:"log"`
}

// ServerConfig defines the HTTP server settings.
type ServerConfig struct {
	Host            string        `yaml:"host"`
	Port            int           `yaml:"port"`
	ReadTimeout     time.Duration `yaml:"read_timeout"`
	WriteTimeout    time.Duration `yaml:"write_timeout"`
	ShutdownTimeout time.Duration `yaml:"shutdown_timeout"`
}

// WatcherConfig defines settings for the sentinel watcher.
type WatcherConfig struct {
	Interval    time.Duration `yaml:"interval"`
	Concurrency int           `yaml:"concurrency"`
	Targets     []Target      `yaml:"targets"`
}

// Target represents a single resource to monitor.
type Target struct {
	Name    string            `yaml:"name"`
	URL     string            `yaml:"url"`
	Method  string            `yaml:"method"`
	Headers map[string]string `yaml:"headers"`
	Timeout time.Duration     `yaml:"timeout"`
}

// LogConfig defines logging settings.
type LogConfig struct {
	Level  string `yaml:"level"`
	Format string `yaml:"format"`
}

// DefaultConfig returns a Config populated with sensible defaults.
// Personal note: bumped default concurrency to 10 and interval to 60s
// since I'm running this on a beefier machine and monitoring more endpoints.
// Also bumped ReadTimeout to 10s and WriteTimeout to 20s — was hitting
// occasional timeouts on slower endpoints in my homelab setup.
func DefaultConfig() *Config {
	return &Config{
		Server: ServerConfig{
			Host:            "0.0.0.0",
			Port:            8080,
			ReadTimeout:     10 * time.Second,
			WriteTimeout:    20 * time.Second,
			ShutdownTimeout: 30 * time.Second,
		},
		Watcher: WatcherConfig{
			Interval:    60 * time.Second,
			Concurrency: 10,
			Targets:     []Target{},
		},
		Log: LogConfig{
			Level:  "info",
			Format: "json",
		},
	}
}

// Load reads a YAML configuration file from the given path and merges it
// with the default configuration. Returns an error if the file cannot be
// read or parsed.
func Load(path string) (*Config, error) {
	cfg := DefaultConfig()

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading config file %q: %w", path, err)
	}

	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("parsing config file %q: %w", path, err)
	}

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return cfg, nil
}

// Validate checks that the configuration contains valid values.
func (c *Config) Validate() error {
	if c.Server.Port < 1 || c.Server.Port > 65535 {
		return fmt.Errorf("server.port must be between 1 and 65535, got %d", c.Server.Port)
	}

	if c.Watcher.Interval <= 0 {
		return fmt.Errorf("watcher.interval must be positive")
	}

	if c.Watcher.Concurrency < 1 {
		return fmt.Errorf("watcher.concurrency must be at least 1")
	}

	for i, t := range c.Watcher.Targets {
		if t.Name == "" {
			return fmt.Errorf("watcher.targets[%d].name must not be empty", i)
		}
		if t.URL == "" {
			return fmt.Errorf("watcher.targets[%d].url must not be empty", i)
		}
	}

	return nil
}
