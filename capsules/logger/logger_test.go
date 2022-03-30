package logger

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/opentracing/opentracing-go"
	"testing"

	mock "go.adenix.dev/adderall/mock/tracing"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
	"gotest.tools/assert"
)

func TestNewLogger(t *testing.T) {
	ctrl := gomock.NewController(t)
	trace := mock.NewMockTracer(ctrl)
	logger, cleanup := NewLogger(trace, WithHost("host"))
	defer cleanup()

	if v, ok := logger.(*defaultLogger); !ok {
		t.Errorf("expected type *defaultLogger, got %T", v)
	}
}

func TestLoggerLevels(t *testing.T) {
	tests := []struct {
		name             string
		action           func(ctx context.Context, l Logger)
		mockTracerInject func(sp opentracing.SpanContext, format interface{}, carrier interface{}) error
		assert           func(ol *observer.ObservedLogs)
	}{
		{
			name: "Debug",
			action: func(_ context.Context, l Logger) {
				l.Debug("foo")
			},
			mockTracerInject: mockTracerInject(t),
			assert:           assertLogWithNew(t, zap.DebugLevel, "foo"),
		},
		{
			name: "DebugWith",
			action: func(_ context.Context, l Logger) {
				l.Debug("foo", "type", "debug")
			},
			mockTracerInject: mockTracerInject(t),
			assert:           assertLogWithNew(t, zap.DebugLevel, "foo", "type", "debug"),
		},
		{
			name: "Info",
			action: func(_ context.Context, l Logger) {
				l.Info("foo")
			},
			mockTracerInject: mockTracerInject(t),
			assert:           assertLogWithNew(t, zap.InfoLevel, "foo"),
		},
		{
			name: "InfoWith",
			action: func(_ context.Context, l Logger) {
				l.Info("foo", "type", "info")
			},
			mockTracerInject: mockTracerInject(t),
			assert:           assertLogWithNew(t, zap.InfoLevel, "foo", "type", "info"),
		},
		{
			name: "Warn",
			action: func(_ context.Context, l Logger) {
				l.Warn("foo")
			},
			mockTracerInject: mockTracerInject(t),
			assert:           assertLogWithNew(t, zap.WarnLevel, "foo"),
		},
		{
			name: "WarnWith",
			action: func(_ context.Context, l Logger) {
				l.Warn("foo", "type", "warn")
			},
			mockTracerInject: mockTracerInject(t),
			assert:           assertLogWithNew(t, zap.WarnLevel, "foo", "type", "warn"),
		},
		{
			name: "Error",
			action: func(_ context.Context, l Logger) {
				l.Error("foo")
			},
			mockTracerInject: mockTracerInject(t),
			assert:           assertLogWithNew(t, zap.ErrorLevel, "foo"),
		},
		{
			name: "ErrorWith",
			action: func(_ context.Context, l Logger) {
				l.Error("foo", "type", "error")
			},
			mockTracerInject: mockTracerInject(t),
			assert:           assertLogWithNew(t, zap.ErrorLevel, "foo", "type", "error"),
		},
		{
			name: "DebugCtx",
			action: func(ctx context.Context, l Logger) {
				l.DebugCtx(ctx, "foo")
			},
			mockTracerInject: mockTracerInject(t, "fizz", "buzz"),
			assert:           assertLogWithNew(t, zap.DebugLevel, "foo", "fizz", "buzz"),
		},
		{
			name: "DebugCtxWith",
			action: func(ctx context.Context, l Logger) {
				l.DebugCtx(ctx, "foo", "type", "debug")
			},
			mockTracerInject: mockTracerInject(t, "fizz", "buzz"),
			assert:           assertLogWithNew(t, zap.DebugLevel, "foo", "fizz", "buzz", "type", "debug"),
		},
		{
			name: "InfoCtx",
			action: func(ctx context.Context, l Logger) {
				l.InfoCtx(ctx, "foo")
			},
			mockTracerInject: mockTracerInject(t, "fizz", "buzz"),
			assert:           assertLogWithNew(t, zap.InfoLevel, "foo", "fizz", "buzz"),
		},
		{
			name: "InfoCtxWith",
			action: func(ctx context.Context, l Logger) {
				l.InfoCtx(ctx, "foo", "type", "info")
			},
			mockTracerInject: mockTracerInject(t, "fizz", "buzz"),
			assert:           assertLogWithNew(t, zap.InfoLevel, "foo", "fizz", "buzz", "type", "info"),
		},
		{
			name: "WarnCtx",
			action: func(ctx context.Context, l Logger) {
				l.WarnCtx(ctx, "foo")
			},
			mockTracerInject: mockTracerInject(t, "fizz", "buzz"),
			assert:           assertLogWithNew(t, zap.WarnLevel, "foo", "fizz", "buzz"),
		},
		{
			name: "WarnCtxWith",
			action: func(ctx context.Context, l Logger) {
				l.WarnCtx(ctx, "foo", "type", "warn")
			},
			mockTracerInject: mockTracerInject(t, "fizz", "buzz"),
			assert:           assertLogWithNew(t, zap.WarnLevel, "foo", "fizz", "buzz", "type", "warn"),
		},
		{
			name: "ErrorCtx",
			action: func(ctx context.Context, l Logger) {
				l.ErrorCtx(ctx, "foo")
			},
			mockTracerInject: mockTracerInject(t, "fizz", "buzz"),
			assert:           assertLogWithNew(t, zap.ErrorLevel, "foo", "fizz", "buzz"),
		},
		{
			name: "ErrorCtxWith",
			action: func(ctx context.Context, l Logger) {
				l.ErrorCtx(ctx, "foo", "type", "error")
			},
			mockTracerInject: mockTracerInject(t, "fizz", "buzz"),
			assert:           assertLogWithNew(t, zap.ErrorLevel, "foo", "fizz", "buzz", "type", "error"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctx := context.Background()
			ctrl := gomock.NewController(t)
			trace := mock.NewMockTracer(ctrl)

			fac, ol := observer.New(zap.DebugLevel)
			l := &defaultLogger{
				l:      zap.New(fac).Sugar(),
				tracer: trace,
			}
			defer l.Sync()

			spanCtx := mock.NewMockSpanContext(ctrl)

			span := mock.NewMockSpan(ctrl)
			span.EXPECT().Tracer().Return(trace)
			span.EXPECT().Context().Return(spanCtx)

			trace.EXPECT().Inject(gomock.Eq(spanCtx), gomock.Eq(opentracing.TextMap), gomock.Any()).DoAndReturn(test.mockTracerInject)

			ctx = opentracing.ContextWithSpan(ctx, span)
			test.action(ctx, l)
			test.assert(ol)
		})
	}
}

func mockTracerInject(t *testing.T, keysAndValues ...string) func(sp opentracing.SpanContext, format interface{}, carrier interface{}) error {
	return func(sp opentracing.SpanContext, format interface{}, carrier interface{}) error {
		switch format {
		case opentracing.HTTPHeaders, opentracing.TextMap:
			if len(keysAndValues)%2 != 0 {
				t.Fatalf("expected even number of keysAndValues, got %d", len(keysAndValues))
			}
			for i := 0; i < len(keysAndValues); i += 2 {
				carrier.(opentracing.TextMapWriter).Set(keysAndValues[i], keysAndValues[i+1])
			}
			return nil
		}
		return opentracing.ErrUnsupportedFormat
	}
}

func assertLogWithNew(t *testing.T, expectedLevel zapcore.Level, expectedMessage string, keysAndValues ...string) func(ol *observer.ObservedLogs) {
	return func(ol *observer.ObservedLogs) {
		assert.Equal(t, expectedLevel, ol.AllUntimed()[0].Level)
		assert.Equal(t, expectedMessage, ol.AllUntimed()[0].Entry.Message)

		if len(keysAndValues)%2 != 0 {
			t.Errorf("expected even number of keysAndValues, got %d", len(keysAndValues))
			return
		}

		fields := []zap.Field{}
		for i := 0; i < len(keysAndValues); i += 2 {
			fields = append(fields, zap.String(keysAndValues[i], keysAndValues[i+1]))
		}
		assert.DeepEqual(t, fields, ol.AllUntimed()[0].Context)
	}
}
