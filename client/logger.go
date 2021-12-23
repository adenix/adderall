package client

var _ Logger = NoopLogger{}

type Logger interface {
	Debug(msg string, keysAndValues ...interface{})
	Info(msg string, keysAndValues ...interface{})
	Warn(msg string, keysAndValues ...interface{})
	Error(msg string, keysAndValues ...interface{})
}

type NoopLogger struct{}

var _ Logger = (*NoopLogger)(nil)

func (n NoopLogger) Debug(msg string, keysAndValues ...interface{})  {}
func (n NoopLogger) Info(msg string, keysAndValues ...interface{})   {}
func (n NoopLogger) Warn(msg string, keysAndValues ...interface{})   {}
func (n NoopLogger) Error(msg string, keysAndValues ...interface{})  {}
func (n NoopLogger) DPanic(msg string, keysAndValues ...interface{}) {}
