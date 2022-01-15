package client

import (
	"reflect"
	"testing"

	"github.com/opentracing/opentracing-go"
	"go.adenix.dev/adderall/internal/pointer"
)

type factoryAssertion func(t *testing.T, f *factory)

type clientAssertion func(t *testing.T, c *Client)

func TestNewFactory(t *testing.T) {
	c := Config{
		TimeoutMs:      pointer.IntP(1500),
		RetryWaitMinMs: pointer.IntP(1500),
		RetryMax:       pointer.IntP(2),
	}

	tests := []struct {
		name    string
		opts    []FactoryOption
		asserts []factoryAssertion
	}{
		{
			name: "Default",
			asserts: []factoryAssertion{
				assertFactoryTracer(opentracing.NoopTracer{}),
				assertFactoryLogger(NoopLogger{}),
				assertFactoryConfig(defaultConfig()),
			},
		},
		{
			name: "WithTracer",
			opts: []FactoryOption{
				WithTracer(opentracing.GlobalTracer()),
			},
			asserts: []factoryAssertion{
				assertFactoryTracer(opentracing.GlobalTracer()),
				assertFactoryLogger(NoopLogger{}),
				assertFactoryConfig(defaultConfig()),
			},
		},
		{
			name: "WithLogger",
			opts: []FactoryOption{
				WithLogger(&testLogger{}),
			},
			asserts: []factoryAssertion{
				assertFactoryTracer(opentracing.NoopTracer{}),
				assertFactoryLogger(&testLogger{}),
				assertFactoryConfig(defaultConfig()),
			},
		},
		{
			name: "WithConfig",
			opts: []FactoryOption{
				WithConfig(c),
			},
			asserts: []factoryAssertion{
				assertFactoryTracer(opentracing.NoopTracer{}),
				assertFactoryLogger(NoopLogger{}),
				assertFactoryConfig(c),
			},
		},
		{
			name: "All",
			opts: []FactoryOption{
				WithTracer(opentracing.GlobalTracer()),
				WithLogger(&testLogger{}),
				WithConfig(c),
			},
			asserts: []factoryAssertion{
				assertFactoryTracer(opentracing.GlobalTracer()),
				assertFactoryLogger(&testLogger{}),
				assertFactoryConfig(c),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			f := NewFactory(test.opts...)
			if v, ok := f.(*factory); ok {
				for _, assert := range test.asserts {
					assert(t, v)
				}
			} else {
				t.Errorf("expected type factory, got %T", f)
			}
		})
	}
}

func TestCreate(t *testing.T) {
	tests := []struct {
		name    string
		opts    []Option
		asserts []clientAssertion
	}{
		{
			name: "Default",
			asserts: []clientAssertion{
				assertTracer(opentracing.NoopTracer{}),
				assertLogger(NoopLogger{}),
				assertConfig(defaultConfig()),
			},
		},
		{
			name: "WithTracer",
			opts: []Option{
				WithClientTracer(opentracing.GlobalTracer()),
			},
			asserts: []clientAssertion{
				assertTracer(opentracing.GlobalTracer()),
				assertLogger(NoopLogger{}),
				assertConfig(defaultConfig()),
			},
		},
		{
			name: "WithLogger",
			opts: []Option{
				WithClientLogger(&testLogger{}),
			},
			asserts: []clientAssertion{
				assertTracer(opentracing.NoopTracer{}),
				assertLogger(&testLogger{}),
				assertConfig(defaultConfig()),
			},
		},
		{
			name: "All",
			opts: []Option{
				WithClientTracer(opentracing.GlobalTracer()),
				WithClientLogger(&testLogger{}),
				WithRetryMax(100),
				WithRetryWaitMinMs(10000),
				WithTimeoutMs(15000),
			},
			asserts: []clientAssertion{
				assertTracer(opentracing.GlobalTracer()),
				assertLogger(&testLogger{}),
				assertConfigRetryMax(100),
				assertConfigRetryWaitMinMs(10000),
				assertConfigTimeoutMs(15000),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := NewFactory().Create(test.opts...)
			for _, assert := range test.asserts {
				assert(t, c)
			}
		})
	}
}

func assertFactoryTracer(expected opentracing.Tracer) factoryAssertion {
	return func(t *testing.T, f *factory) {
		if f.tracer == nil {
			t.Errorf("expected %T, got nil", expected)
		} else if reflect.TypeOf(f.tracer) != reflect.TypeOf(expected) {
			t.Errorf("expected type %T, got %T", expected, f.tracer)
		}
	}
}

func assertFactoryLogger(expected Logger) factoryAssertion {
	return func(t *testing.T, f *factory) {
		if f.logger == nil {
			t.Errorf("expected %T, got nil", expected)
		} else if reflect.TypeOf(f.logger) != reflect.TypeOf(expected) {
			t.Errorf("expected type %T, got %T", expected, f.logger)
		}
	}
}

func assertFactoryConfig(expected Config) factoryAssertion {
	return func(t *testing.T, f *factory) {
		if ok := reflect.DeepEqual(f.config, expected); !ok {
			t.Errorf("expexted %v, got %v", expected, f.config)
		}
	}
}

func assertTracer(expected opentracing.Tracer) clientAssertion {
	return func(t *testing.T, c *Client) {
		if c.tracer == nil {
			t.Errorf("expected %T, got nil", expected)
		} else if reflect.TypeOf(c.tracer) != reflect.TypeOf(expected) {
			t.Errorf("expected type %T, got %T", expected, c.tracer)
		}
	}
}

func assertLogger(expected Logger) clientAssertion {
	return func(t *testing.T, c *Client) {
		if c.logger == nil {
			t.Errorf("expected %T, got nil", expected)
		} else if reflect.TypeOf(c.logger) != reflect.TypeOf(expected) {
			t.Errorf("expected type %T, got %T", expected, c.logger)
		}
	}
}

func assertConfig(expected Config) clientAssertion {
	return func(t *testing.T, c *Client) {
		if ok := reflect.DeepEqual(c.config, expected); !ok {
			t.Errorf("expexted %v, got %v", expected, c.config)
		}
	}
}

func assertConfigRetryMax(expected int) clientAssertion {
	return func(t *testing.T, c *Client) {
		if c.config.RetryMax == nil {
			t.Errorf("expected %d, got nil", expected)
		} else if *c.config.RetryMax != expected {
			t.Errorf("expected type %d, got %d", expected, *c.config.RetryMax)
		}
	}
}

func assertConfigRetryWaitMinMs(expected int) clientAssertion {
	return func(t *testing.T, c *Client) {
		if c.config.RetryWaitMinMs == nil {
			t.Errorf("expected %d, got nil", expected)
		} else if *c.config.RetryWaitMinMs != expected {
			t.Errorf("expected type %d, got %d", expected, *c.config.RetryWaitMinMs)
		}
	}
}

func assertConfigTimeoutMs(expected int) clientAssertion {
	return func(t *testing.T, c *Client) {
		if c.config.TimeoutMs == nil {
			t.Errorf("expected %d, got nil", expected)
		} else if *c.config.TimeoutMs != expected {
			t.Errorf("expected type %d, got %d", expected, *c.config.TimeoutMs)
		}
	}
}
