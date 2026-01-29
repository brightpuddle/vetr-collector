// Package config provides YAML configuration file support for multi-fabric collection.
package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// GlobalConfig holds global settings that apply to all fabrics.
type GlobalConfig struct {
	Username          string            `yaml:"username"`
	Password          string            `yaml:"password"`
	RequestRetryCount int               `yaml:"request_retry_count"`
	RetryDelay        int               `yaml:"retry_delay"`
	BatchSize         int               `yaml:"batch_size"`
	PageSize          int               `yaml:"page_size"`
	Confirm           bool              `yaml:"confirm"`
	Verbose           bool              `yaml:"verbose"`
	Class             string            `yaml:"class"`
	Query             map[string]string `yaml:"query"`
}

// FabricConfig holds per-fabric configuration.
type FabricConfig struct {
	Name              string            `yaml:"name"`
	URL               string            `yaml:"url"`
	Username          string            `yaml:"username"`
	Password          string            `yaml:"password"`
	RequestRetryCount *int              `yaml:"request_retry_count"`
	RetryDelay        *int              `yaml:"retry_delay"`
	BatchSize         *int              `yaml:"batch_size"`
	PageSize          *int              `yaml:"page_size"`
	Confirm           *bool             `yaml:"confirm"`
	Verbose           *bool             `yaml:"verbose"`
	Class             string            `yaml:"class"`
	Query             map[string]string `yaml:"query"`
}

// Config represents the full YAML configuration file structure.
type Config struct {
	Global  GlobalConfig   `yaml:"global"`
	Fabrics []FabricConfig `yaml:"fabrics"`
}

// LoadConfig reads and parses a YAML configuration file.
func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	if err := validateConfig(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

// validateConfig ensures the configuration is valid.
func validateConfig(cfg *Config) error {
	if len(cfg.Fabrics) == 0 {
		return fmt.Errorf("no fabrics defined in config file")
	}

	// Track unique names/hosts
	names := make(map[string]bool)

	for i, fabric := range cfg.Fabrics {
		if fabric.URL == "" {
			return fmt.Errorf("fabric %d: url is required", i)
		}

		// Determine the derived name (name if set, otherwise url)
		derivedName := fabric.Name
		if derivedName == "" {
			derivedName = fabric.URL
		}

		// Check for duplicate names
		if names[derivedName] {
			return fmt.Errorf("duplicate fabric name/url: %s", derivedName)
		}
		names[derivedName] = true
	}

	return nil
}

// GetFabricName returns the display name for a fabric (name if set, otherwise url).
func (f *FabricConfig) GetFabricName() string {
	if f.Name != "" {
		return f.Name
	}
	return f.URL
}

// GetOutputFileName returns the output filename for a fabric.
func (f *FabricConfig) GetOutputFileName() string {
	return f.GetFabricName() + ".zip"
}

// MergeWithGlobal applies global settings to a fabric config, with fabric settings taking precedence.
func (f *FabricConfig) MergeWithGlobal(global GlobalConfig) FabricConfig {
	merged := *f

	// Apply global values if not overridden in fabric config
	if merged.Username == "" {
		merged.Username = global.Username
	}
	if merged.Password == "" {
		merged.Password = global.Password
	}
	if merged.RequestRetryCount == nil {
		merged.RequestRetryCount = &global.RequestRetryCount
	}
	if merged.RetryDelay == nil {
		merged.RetryDelay = &global.RetryDelay
	}
	if merged.BatchSize == nil {
		merged.BatchSize = &global.BatchSize
	}
	if merged.PageSize == nil {
		merged.PageSize = &global.PageSize
	}
	if merged.Confirm == nil {
		merged.Confirm = &global.Confirm
	}
	if merged.Verbose == nil {
		merged.Verbose = &global.Verbose
	}
	if merged.Class == "" {
		merged.Class = global.Class
	}
	if merged.Query == nil {
		merged.Query = global.Query
	}

	return merged
}

// GetRequestRetryCount returns the request retry count with fallback to default.
func (f *FabricConfig) GetRequestRetryCount() int {
	if f.RequestRetryCount != nil {
		return *f.RequestRetryCount
	}
	return 3 // default
}

// GetRetryDelay returns the retry delay with fallback to default.
func (f *FabricConfig) GetRetryDelay() int {
	if f.RetryDelay != nil {
		return *f.RetryDelay
	}
	return 10 // default
}

// GetBatchSize returns the batch size with fallback to default.
func (f *FabricConfig) GetBatchSize() int {
	if f.BatchSize != nil {
		return *f.BatchSize
	}
	return 7 // default
}

// GetPageSize returns the page size with fallback to default.
func (f *FabricConfig) GetPageSize() int {
	if f.PageSize != nil {
		return *f.PageSize
	}
	return 1000 // default
}

// GetConfirm returns the confirm flag with fallback to default.
func (f *FabricConfig) GetConfirm() bool {
	if f.Confirm != nil {
		return *f.Confirm
	}
	return false // default
}

// GetVerbose returns the verbose flag with fallback to default.
func (f *FabricConfig) GetVerbose() bool {
	if f.Verbose != nil {
		return *f.Verbose
	}
	return false // default
}

// GetClass returns the class with fallback to default.
func (f *FabricConfig) GetClass() string {
	if f.Class != "" {
		return f.Class
	}
	return "all" // default
}
