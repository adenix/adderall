package logger

import (
	"context"

	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger is the interface to log leveled messsages with and without context
type Logger interface {
	Debug(msg string, keysAndValues ...interface{})
	Info(msg string, keysAndValues ...interface{})
	Warn(msg string, keysAndValues ...interface{})
	Error(msg string, keysAndValues ...interface{})

	DebugCtx(ctx context.Context, msg string, keysAndValues ...interface{})
	InfoCtx(ctx context.Context, msg string, keysAndValues ...interface{})
	WarnCtx(ctx context.Context, msg string, keysAndValues ...interface{})
	ErrorCtx(ctx context.Context, msg string, keysAndValues ...interface{})

	Sync()
}

type defaultLogger struct {
	l      *zap.SugaredLogger
	tracer opentracing.Tracer
}

// NewLogger instantiates a Logger instrumented with OpenTracing. Options can be
// passed to overwrite default configurations.
func NewLogger(t opentracing.Tracer, opts ...Option) (Logger, func()) {

	c := zap.NewProductionConfig()

	c.EncoderConfig.TimeKey = "time"
	c.EncoderConfig.LevelKey = "level"
	c.EncoderConfig.StacktraceKey = "stack"
	c.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	if c.InitialFields == nil {
		c.InitialFields = make(map[string]interface{})
	}

	for _, opt := range opts {
		opt(&c)
	}

	zapLogger, _ := c.Build(zap.AddCallerSkip(1))

	logger := &defaultLogger{l: zapLogger.Sugar(), tracer: t}

	return logger, func() { _ = zapLogger.Sync() }
}

// Debug writes a debug level log message without context
func (d *defaultLogger) Debug(msg string, keysAndValues ...interface{}) {
	d.DebugCtx(context.Background(), msg, keysAndValues)
}

// Info writes a info level log message without context
func (d *defaultLogger) Info(msg string, keysAndValues ...interface{}) {
	d.InfoCtx(context.Background(), msg, keysAndValues)
}

// Warn writes a warn level log message without context
func (d *defaultLogger) Warn(msg string, keysAndValues ...interface{}) {
	d.WarnCtx(context.Background(), msg, keysAndValues)
}

// Error writes a error level log message without context
func (d *defaultLogger) Error(msg string, keysAndValues ...interface{}) {
	d.ErrorCtx(context.Background(), msg, keysAndValues)
}

// DebugCtx writes a debug level log message with context
func (d *defaultLogger) DebugCtx(ctx context.Context, msg string, keysAndValues ...interface{}) {
	l := d.getScopedLogger(ctx)
	l.Debugw(msg, keysAndValues...)
}

// InfoCtx writes a info level log message with context
func (d *defaultLogger) InfoCtx(ctx context.Context, msg string, keysAndValues ...interface{}) {
	l := d.getScopedLogger(ctx)
	l.Infow(msg, keysAndValues...)
}

// WarnCtx writes a war level log message with context
func (d *defaultLogger) WarnCtx(ctx context.Context, msg string, keysAndValues ...interface{}) {
	l := d.getScopedLogger(ctx)
	l.Warnw(msg, keysAndValues...)
}

// ErrorCtx writes a error level log message with context
func (d *defaultLogger) ErrorCtx(ctx context.Context, msg string, keysAndValues ...interface{}) {
	l := d.getScopedLogger(ctx)
	l.Errorw(msg, keysAndValues...)
}

// Sync flushes any buffered log entries.
func (d *defaultLogger) Sync() {
	_ = d.l.Sync()
}

func (d *defaultLogger) getScopedLogger(ctx context.Context) *zap.SugaredLogger {
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
