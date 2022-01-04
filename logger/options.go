package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Option func(c *zap.Config)

func WithLevel(level zapcore.Level) Option {
	return func(c *zap.Config) {
		c.Level = zap.NewAtomicLevelAt(level)
	}
}

func WithTimeEncoder(encoder zapcore.TimeEncoder) Option {
	return func(c *zap.Config) {
		c.EncoderConfig.EncodeTime = encoder
	}
}

func WithTimeKey(key string) Option {
	return func(c *zap.Config) {
		c.EncoderConfig.TimeKey = key
	}
}

func WithLevelKey(key string) Option {
	return func(c *zap.Config) {
		c.EncoderConfig.LevelKey = key
	}
}

func WithStacktraceKey(key string) Option {
	return func(c *zap.Config) {
		c.EncoderConfig.StacktraceKey = key
	}
}

func WithHost(key string) Option {
	return func(c *zap.Config) {
		if host, err := os.Hostname(); err == nil {
			c.InitialFields[key] = host
		}
	}
}

func WithPid(key string) Option {
	return func(c *zap.Config) {
		c.InitialFields[key] = int64(os.Getpid())
	}
}
