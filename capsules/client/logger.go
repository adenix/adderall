package client

// Logger is a local interface for logging functionality
type Logger interface {
	Debug(msg string, keysAndValues ...interface{})
	Info(msg string, keysAndValues ...interface{})
	Warn(msg string, keysAndValues ...interface{})
	Error(msg string, keysAndValues ...interface{})
}

// NoopLogger is a noop logger implementation.
type NoopLogger struct{}

var _ Logger = (*NoopLogger)(nil)

// Debug ...
func (n NoopLogger) Debug(msg string, keysAndValues ...interface{}) {}

// Info ...
func (n NoopLogger) Info(msg string, keysAndValues ...interface{}) {}

// Warn ...
func (n NoopLogger) Warn(msg string, keysAndValues ...interface{}) {}

// Error ...
func (n NoopLogger) Error(msg string, keysAndValues ...interface{}) {}

// DPanic ...
func (n NoopLogger) DPanic(msg string, keysAndValues ...interface{}) {}
