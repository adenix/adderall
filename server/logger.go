package server

import (
	"context"
)

// Logger is a local interface for logging functionality
type Logger interface {
	DebugCtx(ctx context.Context, msg string, keysAndValues ...interface{})
	InfoCtx(ctx context.Context, msg string, keysAndValues ...interface{})
	WarnCtx(ctx context.Context, msg string, keysAndValues ...interface{})
	ErrorCtx(ctx context.Context, msg string, keysAndValues ...interface{})
}

// NoopLogger is a noop logger implementation.
type NoopLogger struct{}

var _ Logger = (*NoopLogger)(nil)

// DebugCtx ...
func (n NoopLogger) DebugCtx(ctx context.Context, msg string, keysAndValues ...interface{}) {}

// InfoCtx ...
func (n NoopLogger) InfoCtx(ctx context.Context, msg string, keysAndValues ...interface{}) {}

// WarnCtx ...
func (n NoopLogger) WarnCtx(ctx context.Context, msg string, keysAndValues ...interface{}) {}

// ErrorCtx ...
func (n NoopLogger) ErrorCtx(ctx context.Context, msg string, keysAndValues ...interface{}) {}

// DPanicCtx ...
func (n NoopLogger) DPanicCtx(ctx context.Context, msg string, keysAndValues ...interface{}) {}
