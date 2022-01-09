package client

import (
	"fmt"
	"testing"
)

func TestNoopLogger(t *testing.T) {
	logger := NoopLogger{}
	tests := []struct {
		level func(msg string, keysAndValues ...interface{})
	}{
		{level: logger.Debug},
		{level: logger.Info},
		{level: logger.Warn},
		{level: logger.Error},
		{level: logger.DPanic},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("NoopLogger-%d", i), func(t *testing.T) {
			test.level("", nil)
		})
	}
}

// testLogger is used in tests that use reflection to check the type
type testLogger struct {
	Logger
}
