package ingress

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"kindctl/internal/logger"
)

func TestAddHostEntry(t *testing.T) {
	log := logger.NewLogger("debug")
	// Note: Actual /etc/hosts modification requires sudo, tested in integration tests.
	// Expect error in test environment due to permissions.
	err := AddHostEntry(log, "test.local")
	assert.Error(t, err)
}
