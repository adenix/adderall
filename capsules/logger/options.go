package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Option interface to identify functional options
type Option func(c *zap.Config)

// WithLevel provides an Option to provide a minimum level logged.
// Defaults to info
func WithLevel(level zapcore.Level) Option {
	return func(c *zap.Config) {
		c.Level = zap.NewAtomicLevelAt(level)
	}
}

// WithTimeEncoder provides an Option to provide a time encoder.
// Defaults to epoch time encoder
func WithTimeEncoder(encoder zapcore.TimeEncoder) Option {
	return func(c *zap.Config) {
		c.EncoderConfig.EncodeTime = encoder
	}
}

// WithTimeKey provides an Option to provide a key for the time value.
// Defaults to 'time'
func WithTimeKey(key string) Option {
	return func(c *zap.Config) {
		c.EncoderConfig.TimeKey = key
	}
}

// WithLevelKey provides an Option to provide a key for the level value.
// Defaults to 'level'
func WithLevelKey(key string) Option {
	return func(c *zap.Config) {
		c.EncoderConfig.LevelKey = key
	}
}

// WithStacktraceKey provides an Option to provide a key for the stacktrace
// value.
// Defaults to 'stack'
func WithStacktraceKey(key string) Option {
	return func(c *zap.Config) {
		c.EncoderConfig.StacktraceKey = key
	}
}

// WithHost provides an Option to add the hostname to a provided key
func WithHost(key string) Option {
	return func(c *zap.Config) {
		if host, err := os.Hostname(); err == nil {
			withInitialField(c, key, host)
		}
	}
}

// WithPid provides an Option to add the pid to a provided key
func WithPid(key string) Option {
	return func(c *zap.Config) {
		withInitialField(c, key, int64(os.Getpid()))
	}
}

func withInitialField(c *zap.Config, key string, value interface{}) {
	if c.InitialFields == nil {
		c.InitialFields = make(map[string]interface{})
	}
	c.InitialFields[key] = value
}
