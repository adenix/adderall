package server

import (
	"context"

	"go.adenix.dev/adderall/logger"
)

type Logger interface {
	DebugCtx(ctx context.Context, msg string, keysAndValues ...interface{})
	InfoCtx(ctx context.Context, msg string, keysAndValues ...interface{})
	WarnCtx(ctx context.Context, msg string, keysAndValues ...interface{})
	ErrorCtx(ctx context.Context, msg string, keysAndValues ...interface{})
}

type NoopLogger struct {
	logger.Logger
}

var _ Logger = (*NoopLogger)(nil)
var _ logger.Logger = (*NoopLogger)(nil)

func (n NoopLogger) DebugCtx(ctx context.Context, msg string, keysAndValues ...interface{})  {}
func (n NoopLogger) InfoCtx(ctx context.Context, msg string, keysAndValues ...interface{})   {}
func (n NoopLogger) WarnCtx(ctx context.Context, msg string, keysAndValues ...interface{})   {}
func (n NoopLogger) ErrorCtx(ctx context.Context, msg string, keysAndValues ...interface{})  {}
func (n NoopLogger) DPanicCtx(ctx context.Context, msg string, keysAndValues ...interface{}) {}
