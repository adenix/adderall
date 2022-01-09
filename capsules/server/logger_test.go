package server

import (
	"context"
	"fmt"
	"testing"
)

func TestNoopLogger(t *testing.T) {
	logger := NoopLogger{}
	tests := []struct {
		level func(ctx context.Context, msg string, keysAndValues ...interface{})
	}{
		{level: logger.DebugCtx},
		{level: logger.InfoCtx},
		{level: logger.WarnCtx},
		{level: logger.ErrorCtx},
		{level: logger.DPanicCtx},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("NoopLogger-%d", i), func(t *testing.T) {
			test.level(context.TODO(), "", nil)
		})
	}
}

// testLogger is used in tests that use reflection to check the type
type testLogger struct {
	Logger
}
