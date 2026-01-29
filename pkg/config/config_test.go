package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfig(t *testing.T) {
	a := assert.New(t)

	// Create a temporary config file
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	// Valid config
	validConfig := `
global:
  username: admin
  password: pass123
  batch_size: 10
fabrics:
  - name: fabric1
    url: 10.1.1.1
  - name: fabric2
    url: 10.2.2.2
    username: user2
`
	err := os.WriteFile(configPath, []byte(validConfig), 0644)
	a.NoError(err)

	cfg, err := LoadConfig(configPath)
	a.NoError(err)
	a.NotNil(cfg)
	a.Equal("admin", cfg.Global.Username)
	a.Equal("pass123", cfg.Global.Password)
	a.Equal(10, cfg.Global.BatchSize)
	a.Len(cfg.Fabrics, 2)
	a.Equal("fabric1", cfg.Fabrics[0].Name)
	a.Equal("10.1.1.1", cfg.Fabrics[0].URL)
	a.Equal("fabric2", cfg.Fabrics[1].Name)
	a.Equal("user2", cfg.Fabrics[1].Username)
}

func TestLoadConfigNoFabrics(t *testing.T) {
	a := assert.New(t)

	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	// Config with no fabrics
	noFabricsConfig := `
global:
  username: admin
fabrics: []
`
	err := os.WriteFile(configPath, []byte(noFabricsConfig), 0644)
	a.NoError(err)

	_, err = LoadConfig(configPath)
	a.Error(err)
	a.Contains(err.Error(), "no fabrics defined")
}

func TestLoadConfigMissingURL(t *testing.T) {
	a := assert.New(t)

	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	// Config with fabric missing URL
	missingURLConfig := `
global:
  username: admin
fabrics:
  - name: fabric1
`
	err := os.WriteFile(configPath, []byte(missingURLConfig), 0644)
	a.NoError(err)

	_, err = LoadConfig(configPath)
	a.Error(err)
	a.Contains(err.Error(), "url is required")
}

func TestLoadConfigDuplicateNames(t *testing.T) {
	a := assert.New(t)

	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	// Config with duplicate fabric names
	duplicateConfig := `
global:
  username: admin
fabrics:
  - name: fabric1
    url: 10.1.1.1
  - name: fabric1
    url: 10.2.2.2
`
	err := os.WriteFile(configPath, []byte(duplicateConfig), 0644)
	a.NoError(err)

	_, err = LoadConfig(configPath)
	a.Error(err)
	a.Contains(err.Error(), "duplicate fabric name/url")
}

func TestLoadConfigDuplicateURLs(t *testing.T) {
	a := assert.New(t)

	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	// Config with duplicate URLs (no names)
	duplicateURLConfig := `
global:
  username: admin
fabrics:
  - url: 10.1.1.1
  - url: 10.1.1.1
`
	err := os.WriteFile(configPath, []byte(duplicateURLConfig), 0644)
	a.NoError(err)

	_, err = LoadConfig(configPath)
	a.Error(err)
	a.Contains(err.Error(), "duplicate fabric name/url")
}

func TestGetFabricName(t *testing.T) {
	a := assert.New(t)

	// With name
	fabric := FabricConfig{
		Name: "prod",
		URL:  "10.1.1.1",
	}
	a.Equal("prod", fabric.GetFabricName())

	// Without name
	fabric = FabricConfig{
		URL: "10.2.2.2",
	}
	a.Equal("10.2.2.2", fabric.GetFabricName())
}

func TestGetOutputFileName(t *testing.T) {
	a := assert.New(t)

	fabric := FabricConfig{
		Name: "prod",
		URL:  "10.1.1.1",
	}
	a.Equal("prod.zip", fabric.GetOutputFileName())

	fabric = FabricConfig{
		URL: "10.2.2.2",
	}
	a.Equal("10.2.2.2.zip", fabric.GetOutputFileName())
}

func TestMergeWithGlobal(t *testing.T) {
	a := assert.New(t)

	global := GlobalConfig{
		Username:          "admin",
		Password:          "pass123",
		RequestRetryCount: 5,
		RetryDelay:        20,
		BatchSize:         10,
		PageSize:          2000,
		Confirm:           true,
		Verbose:           false,
		Class:             "fvTenant",
		Query:             map[string]string{"key": "value"},
	}

	// Fabric with some overrides
	batchSize := 15
	verbose := true
	fabric := FabricConfig{
		Name:      "prod",
		URL:       "10.1.1.1",
		Username:  "produser",
		BatchSize: &batchSize,
		Verbose:   &verbose,
	}

	merged := fabric.MergeWithGlobal(global)

	// Should use fabric overrides
	a.Equal("produser", merged.Username)
	a.Equal(15, *merged.BatchSize)
	a.True(*merged.Verbose)

	// Should use global defaults
	a.Equal("pass123", merged.Password)
	a.Equal(5, *merged.RequestRetryCount)
	a.Equal(20, *merged.RetryDelay)
	a.Equal(2000, *merged.PageSize)
	a.True(*merged.Confirm)
	a.Equal("fvTenant", merged.Class)
	a.Equal("value", merged.Query["key"])
}

func TestGetters(t *testing.T) {
	a := assert.New(t)

	// Test with values set
	retryCount := 5
	retryDelay := 20
	batchSize := 15
	pageSize := 2000
	confirm := true
	verbose := true
	fabric := FabricConfig{
		RequestRetryCount: &retryCount,
		RetryDelay:        &retryDelay,
		BatchSize:         &batchSize,
		PageSize:          &pageSize,
		Confirm:           &confirm,
		Verbose:           &verbose,
		Class:             "fvTenant",
	}

	a.Equal(5, fabric.GetRequestRetryCount())
	a.Equal(20, fabric.GetRetryDelay())
	a.Equal(15, fabric.GetBatchSize())
	a.Equal(2000, fabric.GetPageSize())
	a.True(fabric.GetConfirm())
	a.True(fabric.GetVerbose())
	a.Equal("fvTenant", fabric.GetClass())

	// Test with defaults
	fabric = FabricConfig{}
	a.Equal(3, fabric.GetRequestRetryCount())
	a.Equal(10, fabric.GetRetryDelay())
	a.Equal(7, fabric.GetBatchSize())
	a.Equal(1000, fabric.GetPageSize())
	a.False(fabric.GetConfirm())
	a.False(fabric.GetVerbose())
	a.Equal("all", fabric.GetClass())
}
