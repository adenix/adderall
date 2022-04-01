package logger

import (
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	mock "go.adenix.dev/adderall/mock/zapcore"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type optionAssertion func(t *testing.T, c *zap.Config)

func TestOption(t *testing.T) {
	tests := []struct {
		name   string
		config *zap.Config
		op     Option
		assert optionAssertion
	}{
		{
			name:   "WithLevel",
			op:     WithLevel(zap.WarnLevel),
			assert: assertWithLevel(zap.WarnLevel),
		},
		{
			name:   "WithTimeEncoder",
			op:     WithTimeEncoder(zapcore.EpochNanosTimeEncoder),
			assert: assertWithTimeEncoder(zapcore.EpochNanosTimeEncoder),
		},
		{
			name:   "WithTimeKey",
			op:     WithTimeKey("timeFoo"),
			assert: assertWithTimeKey("timeFoo"),
		},
		{
			name:   "WithLevelKey",
			op:     WithLevelKey("levelFoo"),
			assert: assertWithLevelKey("levelFoo"),
		},
		{
			name:   "WithStacktraceKey",
			op:     WithStacktraceKey("stackFoo"),
			assert: assertWithStacktraceKey("stackFoo"),
		},
		{
			name:   "WithHost",
			op:     WithHost("hostFoo"),
			assert: assertWithIntialField("hostFoo"),
		},
		{
			name:   "WithPid",
			op:     WithPid("pidFoo"),
			assert: assertWithIntialField("pidFoo"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.config == nil {
				test.config = &zap.Config{}
			}
			test.op(test.config)
			test.assert(t, test.config)
		})
	}
}

func assertWithLevel(expected zapcore.Level) optionAssertion {
	return func(t *testing.T, c *zap.Config) {
		if c.Level.Level() != zap.NewAtomicLevelAt(expected).Level() {
			t.Errorf("expected %s, got %s", expected, c.Level.Level())
		}
	}
}

func assertWithTimeEncoder(expected zapcore.TimeEncoder) optionAssertion {
	return func(t *testing.T, c *zap.Config) {
		now := time.Now()

		ctrl := gomock.NewController(t)
		enc := mock.NewMockPrimitiveArrayEncoder(ctrl)

		enc.EXPECT().AppendInt64(gomock.Eq(now.UnixNano()))
		enc.EXPECT().AppendInt64(gomock.Eq(now.UnixNano()))

		c.EncoderConfig.EncodeTime(now, enc)
		expected(now, enc)
	}
}

func assertWithTimeKey(expected string) optionAssertion {
	return func(t *testing.T, c *zap.Config) {
		if c.EncoderConfig.TimeKey != expected {
			t.Errorf("expected %s, got %s", expected, c.EncoderConfig.TimeKey)
		}
	}
}

func assertWithLevelKey(expected string) optionAssertion {
	return func(t *testing.T, c *zap.Config) {
		if c.EncoderConfig.LevelKey != expected {
			t.Errorf("expected %s, got %s", expected, c.EncoderConfig.LevelKey)
		}
	}
}

func assertWithStacktraceKey(expected string) optionAssertion {
	return func(t *testing.T, c *zap.Config) {
		if c.EncoderConfig.StacktraceKey != expected {
			t.Errorf("expected %s, got %s", expected, c.EncoderConfig.StacktraceKey)
		}
	}
}

func assertWithIntialField(key string) optionAssertion {
	return func(t *testing.T, c *zap.Config) {
		if _, ok := c.InitialFields[key]; !ok {
			t.Error("expected value")
		}
	}
}
