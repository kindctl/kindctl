package logger

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewLogger(t *testing.T) {
	log := NewLogger("debug")
	assert.NotNil(t, log)
	log.Info("Test message")
}
