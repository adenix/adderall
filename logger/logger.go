package logger

import (
	"context"
	"os"

	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger interface {
	Debug(ctx context.Context, msg string, keysAndValues ...interface{})
	Info(ctx context.Context, msg string, keysAndValues ...interface{})
	Warn(ctx context.Context, msg string, keysAndValues ...interface{})
	Error(ctx context.Context, msg string, keysAndValues ...interface{})
	Sync()
}

type DefaultLogger struct {
	l      *zap.SugaredLogger
	tracer opentracing.Tracer
}

func NewLogger(t opentracing.Tracer) (Logger, func()) {

	c := zap.NewProductionConfig()

	c.EncoderConfig.TimeKey = "time"
	c.EncoderConfig.LevelKey = "level"
	c.EncoderConfig.StacktraceKey = "stack"
	c.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	if c.InitialFields == nil {
		c.InitialFields = make(map[string]interface{})
	}

	if host, err := os.Hostname(); err == nil {
		c.InitialFields["host"] = host
	}
	c.InitialFields["pid"] = int64(os.Getpid())

	zapLogger, _ := c.Build(zap.AddCallerSkip(1))

	logger := &DefaultLogger{l: zapLogger.Sugar(), tracer: t}

	return logger, func() { _ = zapLogger.Sync() }
}

func (d *DefaultLogger) Debug(ctx context.Context, msg string, keysAndValues ...interface{}) {
	l := d.getScopedLogger(ctx)
	l.Debugw(msg, keysAndValues...)
}

func (d *DefaultLogger) Info(ctx context.Context, msg string, keysAndValues ...interface{}) {
	l := d.getScopedLogger(ctx)
	l.Infow(msg, keysAndValues...)
}

func (d *DefaultLogger) Warn(ctx context.Context, msg string, keysAndValues ...interface{}) {
	l := d.getScopedLogger(ctx)
	l.Warnw(msg, keysAndValues...)
}

func (d *DefaultLogger) Error(ctx context.Context, msg string, keysAndValues ...interface{}) {
	l := d.getScopedLogger(ctx)
	l.Errorw(msg, keysAndValues...)
}

func (d *DefaultLogger) Sync() {
	d.l.Sync()
}

func (d *DefaultLogger) getScopedLogger(ctx context.Context) *zap.SugaredLogger {
	fields := make([]interface{}, 0)

	span := opentracing.SpanFromContext(ctx)
	if span != nil {
		_ = d.tracer.Inject(span.Context(), opentracing.TextMap, &carrier{fields})
	}

	return d.l.With(fields...)
}

type carrier struct {
	fields []interface{}
}

func (c *carrier) Set(key, val string) {
	c.fields = append(c.fields, key)
	c.fields = append(c.fields, val)
}
