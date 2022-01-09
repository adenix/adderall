package client

import (
	"reflect"
	"testing"

	"github.com/opentracing/opentracing-go"
	"go.adenix.dev/adderall/internal/pointer"
)

type optionAssertion func(t *testing.T, c *Client)

type factoryOptionAssertion func(t *testing.T, f *factory)

func TestOption(t *testing.T) {
	tests := []struct {
		name   string
		client *Client
		op     Option
		assert optionAssertion
	}{
		{
			name:   "WithClientLogger",
			op:     WithClientLogger(&testLogger{}),
			assert: assertOptionWithClientLogger(&testLogger{}),
		},
		{
			name:   "WithClientTracer",
			op:     WithClientTracer(opentracing.NoopTracer{}),
			assert: assertOptionWithClientTracer(opentracing.NoopTracer{}),
		},
		{
			name:   "WithTimeoutMs",
			op:     WithTimeoutMs(1000),
			assert: assertOptionWithTimeoutMs(1000),
		},
		{
			name:   "WithRetryWaitMinMs",
			op:     WithRetryWaitMinMs(1000),
			assert: assertOptionWithRetryWaitMinMs(1000),
		},
		{
			name:   "WithRetryMax",
			op:     WithRetryMax(50),
			assert: assertOptionWithRetryMax(50),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.client == nil {
				test.client = &Client{}
			}
			test.op(test.client)
			test.assert(t, test.client)
		})
	}
}

func TestFactoryOption(t *testing.T) {
	tests := []struct {
		name    string
		factory *factory
		op      FactoryOption
		assert  factoryOptionAssertion
	}{
		{
			name:   "WithLogger",
			op:     WithLogger(&testLogger{}),
			assert: assertFactoryOptionWithLogger(&testLogger{}),
		},
		{
			name:   "WithTracer",
			op:     WithTracer(opentracing.NoopTracer{}),
			assert: assertFactoryOptionWithTracer(opentracing.NoopTracer{}),
		},
		{
			name:   "WithConfig-Blank",
			op:     WithConfig(Config{}),
			assert: assertFactoryOptionWithConfig(Config{}),
		},
		{
			name:   "WithConfig-TimeoutMs",
			op:     WithConfig(Config{TimeoutMs: pointer.IntP(1000)}),
			assert: assertFactoryOptionWithConfig(Config{TimeoutMs: pointer.IntP(1000)}),
		},
		{
			name:   "WithConfig-RetryWaitMinMs",
			op:     WithConfig(Config{RetryWaitMinMs: pointer.IntP(1000)}),
			assert: assertFactoryOptionWithConfig(Config{RetryWaitMinMs: pointer.IntP(1000)}),
		},
		{
			name:   "WithConfig-RetryMax",
			op:     WithConfig(Config{RetryMax: pointer.IntP(50)}),
			assert: assertFactoryOptionWithConfig(Config{RetryMax: pointer.IntP(50)}),
		},
		{
			name:   "WithConfig-All",
			op:     WithConfig(Config{TimeoutMs: pointer.IntP(1000), RetryWaitMinMs: pointer.IntP(1000), RetryMax: pointer.IntP(50)}),
			assert: assertFactoryOptionWithConfig(Config{TimeoutMs: pointer.IntP(1000), RetryWaitMinMs: pointer.IntP(1000), RetryMax: pointer.IntP(50)}),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.factory == nil {
				test.factory = &factory{}
			}
			test.op.apply(test.factory)
			test.assert(t, test.factory)
		})
	}
}

func assertOptionWithClientLogger(expected Logger) optionAssertion {
	return func(t *testing.T, c *Client) {
		if c.logger == nil {
			t.Errorf("expected %T, got nil", expected)
		} else if reflect.TypeOf(c.logger) != reflect.TypeOf(expected) {
			t.Errorf("expected type %T, got %T", expected, c.logger)
		}
	}
}

func assertOptionWithClientTracer(expected opentracing.Tracer) optionAssertion {
	return func(t *testing.T, c *Client) {
		if c.tracer == nil {
			t.Errorf("expected %T, got nil", expected)
		} else if reflect.TypeOf(c.tracer) != reflect.TypeOf(expected) {
			t.Errorf("expected type %T, got %T", expected, c.tracer)
		}
	}
}

func assertOptionWithTimeoutMs(expected int) optionAssertion {
	return func(t *testing.T, c *Client) {
		if c.config.TimeoutMs == nil {
			t.Errorf("expected %d, got nil", expected)
		} else if *c.config.TimeoutMs != expected {
			t.Errorf("expected type %d, got %d", expected, *c.config.TimeoutMs)
		}
	}
}

func assertOptionWithRetryWaitMinMs(expected int) optionAssertion {
	return func(t *testing.T, c *Client) {
		if c.config.RetryWaitMinMs == nil {
			t.Errorf("expected %d, got nil", expected)
		} else if *c.config.RetryWaitMinMs != expected {
			t.Errorf("expected type %d, got %d", expected, *c.config.RetryWaitMinMs)
		}
	}
}

func assertOptionWithRetryMax(expected int) optionAssertion {
	return func(t *testing.T, c *Client) {
		if c.config.RetryMax == nil {
			t.Errorf("expected %d, got nil", expected)
		} else if *c.config.RetryMax != expected {
			t.Errorf("expected type %d, got %d", expected, *c.config.RetryMax)
		}
	}
}

func assertFactoryOptionWithLogger(expected Logger) factoryOptionAssertion {
	return func(t *testing.T, f *factory) {
		if f.logger == nil {
			t.Errorf("expected %T, got nil", expected)
		} else if reflect.TypeOf(f.logger) != reflect.TypeOf(expected) {
			t.Errorf("expected type %T, got %T", expected, f.logger)
		}
	}
}

func assertFactoryOptionWithTracer(expected opentracing.Tracer) factoryOptionAssertion {
	return func(t *testing.T, f *factory) {
		if f.tracer == nil {
			t.Errorf("expected %T, got nil", expected)
		} else if reflect.TypeOf(f.tracer) != reflect.TypeOf(expected) {
			t.Errorf("expected type %T, got %T", expected, f.tracer)
		}
	}
}

func assertFactoryOptionWithConfig(expected Config) factoryOptionAssertion {
	return func(t *testing.T, f *factory) {
		if ok := reflect.DeepEqual(f.config, expected); !ok {
			t.Errorf("expexted %v, got %v", expected, f.config)
		}
	}
}
