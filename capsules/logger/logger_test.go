package logger

import (
	"context"
	"testing"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
	"gotest.tools/assert"
)

func TestLoggerLevels(t *testing.T) {
	tests := []struct {
		name      string
		action    func(ctx context.Context, l Logger)
		assertNew func(ol *observer.ObservedLogs)
	}{
		{
			name: "Debug",
			action: func(_ context.Context, l Logger) {
				l.Debug("foo")
			},
			assertNew: assertLogWithNew(t, zap.DebugLevel, "foo"),
		},
		{
			name: "DebugWith",
			action: func(_ context.Context, l Logger) {
				l.Debug("foo", "type", "debug")
			},
			assertNew: assertLogWithNew(t, zap.DebugLevel, "foo", "type", "debug"),
		},
		{
			name: "Info",
			action: func(_ context.Context, l Logger) {
				l.Info("foo")
			},
			assertNew: assertLogWithNew(t, zap.InfoLevel, "foo"),
		},
		{
			name: "InfoWith",
			action: func(_ context.Context, l Logger) {
				l.Info("foo", "type", "info")
			},
			assertNew: assertLogWithNew(t, zap.InfoLevel, "foo", "type", "info"),
		},
		{
			name: "Warn",
			action: func(_ context.Context, l Logger) {
				l.Warn("foo")
			},
			assertNew: assertLogWithNew(t, zap.WarnLevel, "foo"),
		},
		{
			name: "WarnWith",
			action: func(_ context.Context, l Logger) {
				l.Warn("foo", "type", "warn")
			},
			assertNew: assertLogWithNew(t, zap.WarnLevel, "foo", "type", "warn"),
		},
		{
			name: "Error",
			action: func(_ context.Context, l Logger) {
				l.Error("foo")
			},
			assertNew: assertLogWithNew(t, zap.ErrorLevel, "foo"),
		},
		{
			name: "ErrorWith",
			action: func(_ context.Context, l Logger) {
				l.Error("foo", "type", "error")
			},
			assertNew: assertLogWithNew(t, zap.ErrorLevel, "foo", "type", "error"),
		},
		{
			name: "DebugCtx",
			action: func(ctx context.Context, l Logger) {
				l.DebugCtx(ctx, "foo")
			},
			assertNew: assertLogWithNew(t, zap.DebugLevel, "foo"),
		},
		{
			name: "DebugCtxWith",
			action: func(ctx context.Context, l Logger) {
				l.DebugCtx(ctx, "foo", "type", "debug")
			},
			assertNew: assertLogWithNew(t, zap.DebugLevel, "foo", "type", "debug"),
		},
		{
			name: "InfoCtx",
			action: func(ctx context.Context, l Logger) {
				l.InfoCtx(ctx, "foo")
			},
			assertNew: assertLogWithNew(t, zap.InfoLevel, "foo"),
		},
		{
			name: "InfoCtxWith",
			action: func(ctx context.Context, l Logger) {
				l.InfoCtx(ctx, "foo", "type", "info")
			},
			assertNew: assertLogWithNew(t, zap.InfoLevel, "foo", "type", "info"),
		},
		{
			name: "WarnCtx",
			action: func(ctx context.Context, l Logger) {
				l.WarnCtx(ctx, "foo")
			},
			assertNew: assertLogWithNew(t, zap.WarnLevel, "foo"),
		},
		{
			name: "WarnCtxWith",
			action: func(ctx context.Context, l Logger) {
				l.WarnCtx(ctx, "foo", "type", "warn")
			},
			assertNew: assertLogWithNew(t, zap.WarnLevel, "foo", "type", "warn"),
		},
		{
			name: "ErrorCtx",
			action: func(ctx context.Context, l Logger) {
				l.ErrorCtx(ctx, "foo")
			},
			assertNew: assertLogWithNew(t, zap.ErrorLevel, "foo"),
		},
		{
			name: "ErrorCtxWith",
			action: func(ctx context.Context, l Logger) {
				l.ErrorCtx(ctx, "foo", "type", "error")
			},
			assertNew: assertLogWithNew(t, zap.ErrorLevel, "foo", "type", "error"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctx := context.Background()
			fac, ol := observer.New(zap.DebugLevel)
			l := &defaultLogger{
				l: zap.New(fac).Sugar(),
			}

			test.action(ctx, l)
			test.assertNew(ol)
		})
	}
}

func assertLogWithNew(t *testing.T, expectedLevel zapcore.Level, expectedMessage string, keysAndValues ...string) func(ol *observer.ObservedLogs) {
	return func(ol *observer.ObservedLogs) {
		assert.Equal(t, expectedLevel, ol.AllUntimed()[0].Level)
		assert.Equal(t, expectedMessage, ol.AllUntimed()[0].Entry.Message)

		if len(keysAndValues)%2 != 0 {
			t.Errorf("expected even number of keysAndValues, got %d", len(keysAndValues))
			return
		}

		fields := []zap.Field{}
		for i := 0; i < len(keysAndValues); i += 2 {
			fields = append(fields, zap.String(keysAndValues[i], keysAndValues[i+1]))
		}
		assert.DeepEqual(t, fields, ol.AllUntimed()[0].Context)
	}
}
