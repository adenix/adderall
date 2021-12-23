package server

import (
	"context"
)

var _ Logger = NoopLogger{}

type Logger interface {
	Debug(msg string, keysAndValues ...interface{})
	Info(msg string, keysAndValues ...interface{})
	Warn(msg string, keysAndValues ...interface{})
	Error(msg string, keysAndValues ...interface{})

	DebugCtx(ctx context.Context, msg string, keysAndValues ...interface{})
	InfoCtx(ctx context.Context, msg string, keysAndValues ...interface{})
	WarnCtx(ctx context.Context, msg string, keysAndValues ...interface{})
	ErrorCtx(ctx context.Context, msg string, keysAndValues ...interface{})
}

type NoopLogger struct{}

func (n NoopLogger) Debug(msg string, keysAndValues ...interface{})  {}
func (n NoopLogger) Info(msg string, keysAndValues ...interface{})   {}
func (n NoopLogger) Warn(msg string, keysAndValues ...interface{})   {}
func (n NoopLogger) Error(msg string, keysAndValues ...interface{})  {}
func (n NoopLogger) DPanic(msg string, keysAndValues ...interface{}) {}

func (n NoopLogger) DebugCtx(ctx context.Context, msg string, keysAndValues ...interface{})  {}
func (n NoopLogger) InfoCtx(ctx context.Context, msg string, keysAndValues ...interface{})   {}
func (n NoopLogger) WarnCtx(ctx context.Context, msg string, keysAndValues ...interface{})   {}
func (n NoopLogger) ErrorCtx(ctx context.Context, msg string, keysAndValues ...interface{})  {}
func (n NoopLogger) DPanicCtx(ctx context.Context, msg string, keysAndValues ...interface{}) {}
