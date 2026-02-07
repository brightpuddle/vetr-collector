// Package config provides YAML configuration file support for multi-fabric collection.
package config

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"syscall"

	"golang.org/x/term"
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
	Output            string            `yaml:"output"`
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

// New returns a config with default global values.
func New() Config {
	return Config{
		Global: GlobalConfig{
			RequestRetryCount: 3,
			RetryDelay:        10,
			BatchSize:         7,
			PageSize:          1000,
			Confirm:           false,
			Verbose:           false,
			Class:             "all",
		},
	}
}

// LoadConfig reads and parses a YAML configuration file.
func LoadConfig(path string) (*Config, error) {
	cfg, err := ParseConfig(path)
	if err != nil {
		return nil, err
	}

	if err := validateConfig(cfg, true); err != nil {
		return nil, err
	}

	return cfg, nil
}

// ParseConfig reads and parses a YAML configuration file without prompting.
func ParseConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	cfg.ApplyDefaults()
	return &cfg, nil
}

// validateConfig ensures the configuration is valid.
func validateConfig(cfg *Config, requireURL bool) error {
	if len(cfg.Fabrics) == 0 {
		return fmt.Errorf("no fabrics defined in config file")
	}

	// Track unique names/hosts
	names := make(map[string]bool)

	for i, fabric := range cfg.Fabrics {
		if requireURL && fabric.URL == "" {
			return fmt.Errorf("fabric %d: url is required", i)
		}

		// Determine the derived name (name if set, otherwise url)
		derivedName := fabric.Name
		if derivedName == "" {
			derivedName = fabric.URL
		}
		if derivedName == "" {
			derivedName = fmt.Sprintf("fabric-%d", i+1)
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
	if f.Output != "" {
		return f.Output
	}
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

// ApplyDefaults sets global defaults for missing values.
func (c *Config) ApplyDefaults() {
	defaults := New().Global
	if c.Global.RequestRetryCount == 0 {
		c.Global.RequestRetryCount = defaults.RequestRetryCount
	}
	if c.Global.RetryDelay == 0 {
		c.Global.RetryDelay = defaults.RetryDelay
	}
	if c.Global.BatchSize == 0 {
		c.Global.BatchSize = defaults.BatchSize
	}
	if c.Global.PageSize == 0 {
		c.Global.PageSize = defaults.PageSize
	}
	if c.Global.Class == "" {
		c.Global.Class = defaults.Class
	}
}

// NormalizeAndPrompt fills missing values and normalizes inputs.
func (c *Config) NormalizeAndPrompt() error {
	c.ApplyDefaults()
	if len(c.Fabrics) == 0 {
		return fmt.Errorf("no fabrics defined in config file")
	}

	// Prompt for missing URLs and normalize.
	for i := range c.Fabrics {
		label := c.fabricLabel(i)
		if c.Fabrics[i].URL == "" {
			c.Fabrics[i].URL = input(fmt.Sprintf("APIC URL for %s:", label))
		}
		c.Fabrics[i].URL = normalizeURL(c.Fabrics[i].URL)
	}

	// Prompt once for username/password if none provided.
	if !c.hasAnyUsername() {
		label := c.fabricLabel(0)
		c.Global.Username = input(fmt.Sprintf("APIC username for %s (applies to all fabrics):", label))
		c.Global.Password = inputPassword(fmt.Sprintf("APIC password for %s (applies to all fabrics):", label))
	}

	// Apply global username to fabrics when missing.
	if c.Global.Username != "" {
		for i := range c.Fabrics {
			if c.Fabrics[i].Username == "" {
				c.Fabrics[i].Username = c.Global.Username
			}
		}
	}

	// Prompt for missing usernames per fabric.
	for i := range c.Fabrics {
		if c.Fabrics[i].Username == "" {
			label := c.fabricLabel(i)
			c.Fabrics[i].Username = input(fmt.Sprintf("APIC username for %s:", label))
		}
	}

	// Seed password map with any existing passwords.
	passwordByUser := map[string]string{}
	if c.Global.Username != "" && c.Global.Password != "" {
		passwordByUser[c.Global.Username] = c.Global.Password
	}
	for i := range c.Fabrics {
		user := c.Fabrics[i].Username
		if user == "" {
			continue
		}
		if c.Fabrics[i].Password != "" {
			if _, ok := passwordByUser[user]; !ok {
				passwordByUser[user] = c.Fabrics[i].Password
			}
		}
	}

	// Prompt for passwords per unique username when missing.
	for i := range c.Fabrics {
		user := c.Fabrics[i].Username
		if user == "" {
			continue
		}
		if c.Fabrics[i].Password != "" {
			continue
		}
		if pw, ok := passwordByUser[user]; ok {
			c.Fabrics[i].Password = pw
			continue
		}
		label := c.fabricLabel(i)
		pw := inputPassword(fmt.Sprintf("APIC password for %s (%s):", user, label))
		passwordByUser[user] = pw
		c.Fabrics[i].Password = pw
	}

	// Ensure global password is set when global username is used.
	if c.Global.Username != "" && c.Global.Password == "" {
		if pw, ok := passwordByUser[c.Global.Username]; ok {
			c.Global.Password = pw
		}
	}

	// Apply resolved passwords to fabrics.
	for i := range c.Fabrics {
		if c.Fabrics[i].Password == "" {
			if pw, ok := passwordByUser[c.Fabrics[i].Username]; ok {
				c.Fabrics[i].Password = pw
			}
		}
	}

	return validateConfig(c, true)
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

func (c *Config) hasAnyUsername() bool {
	if c.Global.Username != "" {
		return true
	}
	for _, fabric := range c.Fabrics {
		if fabric.Username != "" {
			return true
		}
	}
	return false
}

func (c *Config) fabricLabel(index int) string {
	label := strings.TrimSpace(c.Fabrics[index].GetFabricName())
	if label == "" {
		label = fmt.Sprintf("fabric %d", index+1)
	}
	return label
}

// input collects CLI input.
func input(prompt string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("%s ", prompt)
	input, _ := reader.ReadString('\n')
	return strings.Trim(input, "\r\n")
}

func inputPassword(prompt string) string {
	fmt.Print(prompt + " ")
	pwd, _ := term.ReadPassword(int(syscall.Stdin))
	fmt.Println()
	return string(pwd)
}

func normalizeURL(url string) string {
	url, _ = strings.CutPrefix(url, "http://")
	url, _ = strings.CutPrefix(url, "https://")
	return url
}
