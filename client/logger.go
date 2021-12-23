package client

import "go.adenix.dev/adderall/logger"

type Logger interface {
	Debug(msg string, keysAndValues ...interface{})
	Info(msg string, keysAndValues ...interface{})
	Warn(msg string, keysAndValues ...interface{})
	Error(msg string, keysAndValues ...interface{})
}

type NoopLogger struct {
	logger.Logger
}

var _ Logger = (*NoopLogger)(nil)
var _ logger.Logger = (*NoopLogger)(nil)

func (n NoopLogger) Debug(msg string, keysAndValues ...interface{})  {}
func (n NoopLogger) Info(msg string, keysAndValues ...interface{})   {}
func (n NoopLogger) Warn(msg string, keysAndValues ...interface{})   {}
func (n NoopLogger) Error(msg string, keysAndValues ...interface{})  {}
func (n NoopLogger) DPanic(msg string, keysAndValues ...interface{}) {}
