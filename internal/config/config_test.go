package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfig(t *testing.T) {
	// Create a temporary config file
	configContent := `
logging:
  level: debug
cluster:
  name: test-cluster
dashboard:
  enabled: true
  ingress: dashboard.local
postgres:
  enabled: true
  ingress: postgres.local
  version: "16"
  username: testuser
  password: testpass
  database: testdb
`
	tmpFile, err := os.CreateTemp("", "kindctl.yaml")
	assert.NoError(t, err)
	defer os.Remove(tmpFile.Name())
	_, err = tmpFile.WriteString(configContent)
	assert.NoError(t, err)
	tmpFile.Close()

	// Test loading config
	cfg, err := LoadConfig(tmpFile.Name())
	assert.NoError(t, err)
	assert.NotNil(t, cfg)
	assert.Equal(t, "debug", cfg.Logging.Level)
	assert.Equal(t, "test-cluster", cfg.Cluster.Name)
	assert.True(t, cfg.Dashboard.Enabled)
	assert.Equal(t, "dashboard.local", cfg.Dashboard.Ingress)
	assert.True(t, cfg.Postgres.Enabled)
	assert.Equal(t, "postgres.local", cfg.Postgres.Ingress)
	assert.Equal(t, "16", cfg.Postgres.Version)
	assert.Equal(t, "testuser", cfg.Postgres.Username)
	assert.Equal(t, "testpass", cfg.Postgres.Password)
	assert.Equal(t, "testdb", cfg.Postgres.Database)
}

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()
	assert.Equal(t, "info", cfg.Logging.Level)
	assert.Equal(t, "kind-cluster", cfg.Cluster.Name)
	assert.True(t, cfg.Dashboard.Enabled)
	assert.Equal(t, "dashboard.local", cfg.Dashboard.Ingress)
	assert.False(t, cfg.Postgres.Enabled)
}
