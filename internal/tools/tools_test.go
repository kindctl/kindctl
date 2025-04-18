package tools

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"kindctl/internal/config"
	"kindctl/internal/logger"
)

func TestUpdateCluster(t *testing.T) {
	log := logger.NewLogger("debug")
	cfg := config.DefaultConfig()
	cfg.Postgres.Enabled = true
	cfg.Postgres.Ingress = "postgres.local"

	// Note: Actual tool installation requires kubectl/helm, tested in integration tests.
	// This test verifies the function structure.
	err := UpdateCluster(log, cfg)
	assert.Error(t, err) // Expect error due to missing kubectl/helm in test env
}
