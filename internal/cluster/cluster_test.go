package cluster

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"kindctl/internal/config"
	"kindctl/internal/logger"
)

// Mock logger for testing
func newTestLogger() *logger.Logger {
	return logger.NewLogger("debug")
}

func TestInitialize(t *testing.T) {
	// Note: Actual kind cluster creation requires Docker, so we skip command execution in unit tests.
	// This test verifies config saving logic. Integration tests cover cluster creation.
	log := newTestLogger()
	tmpFile, err := os.CreateTemp("", "kindctl.yaml")
	assert.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	err = Initialize(log, tmpFile.Name())
	assert.NoError(t, err)

	// Verify config file was created
	cfg, err := config.LoadConfig(tmpFile.Name())
	assert.NoError(t, err)
	assert.Equal(t, "kind-cluster", cfg.Cluster.Name)
	assert.True(t, cfg.Dashboard.Enabled)
}
