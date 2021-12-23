package server

import (
	"context"
)

var _ Logger = NoopLogger{}

type Logger interface {
	Debug(ctx context.Context, msg string, keysAndValues ...interface{})
	Info(ctx context.Context, msg string, keysAndValues ...interface{})
	Error(ctx context.Context, msg string, keysAndValues ...interface{})
}

type NoopLogger struct{}

func (n NoopLogger) Debug(ctx context.Context, msg string, keysAndValues ...interface{})  {}
func (n NoopLogger) Info(ctx context.Context, msg string, keysAndValues ...interface{})   {}
func (n NoopLogger) Error(ctx context.Context, msg string, keysAndValues ...interface{})  {}
func (n NoopLogger) DPanic(ctx context.Context, msg string, keysAndValues ...interface{}) {}
