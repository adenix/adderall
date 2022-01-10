package server

import (
	"testing"

	"github.com/opentracing/opentracing-go"
	"go.adenix.dev/adderall/internal/pointer"
)

func TestNewFactory(t *testing.T) {
	c := Config{
		Port:                 pointer.IntP(8080),
		ReadTimeoutMs:        pointer.IntP(10000),
		WriteTimeoutMs:       pointer.IntP(10000),
		RequestTimeoutSec:    pointer.IntP(10),
		ShutdownDelaySeconds: pointer.IntP(5),
		SwaggerFile:          pointer.StringP("/swagger.json"),
	}

	tests := []struct {
		name    string
		opts    []FactoryOption
		asserts []factoryOptionAssertion
	}{
		{
			name: "Default",
			asserts: []factoryOptionAssertion{
				assertFactoryOptionWithTracer(opentracing.NoopTracer{}),
				assertFactoryOptionWithLogger(NoopLogger{}),
				assertFactoryOptionWithConfig(defaultConfig()),
			},
		},
		{
			name: "WithTracer",
			opts: []FactoryOption{
				WithTracer(opentracing.GlobalTracer()),
			},
			asserts: []factoryOptionAssertion{
				assertFactoryOptionWithTracer(opentracing.GlobalTracer()),
				assertFactoryOptionWithLogger(NoopLogger{}),
				assertFactoryOptionWithConfig(defaultConfig()),
			},
		},
		{
			name: "WithLogger",
			opts: []FactoryOption{
				WithLogger(&testLogger{}),
			},
			asserts: []factoryOptionAssertion{
				assertFactoryOptionWithTracer(opentracing.NoopTracer{}),
				assertFactoryOptionWithLogger(&testLogger{}),
				assertFactoryOptionWithConfig(defaultConfig()),
			},
		},
		{
			name: "WithConfig",
			opts: []FactoryOption{
				WithConfig(c),
			},
			asserts: []factoryOptionAssertion{
				assertFactoryOptionWithTracer(opentracing.NoopTracer{}),
				assertFactoryOptionWithLogger(NoopLogger{}),
				assertFactoryOptionWithConfig(c),
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
		asserts []optionAssertion
	}{
		{
			name: "Default",
			asserts: []optionAssertion{
				assertOptionWithServerTracer(opentracing.NoopTracer{}),
				assertOptionWithServerLogger(NoopLogger{}),
				assertOptionWithServerConfig(defaultConfig()),
			},
		},
		{
			name: "WithTracer",
			opts: []Option{
				WithServerTracer(opentracing.GlobalTracer()),
			},
			asserts: []optionAssertion{
				assertOptionWithServerTracer(opentracing.GlobalTracer()),
				assertOptionWithServerLogger(NoopLogger{}),
				assertOptionWithServerConfig(defaultConfig()),
			},
		},
		{
			name: "WithLogger",
			opts: []Option{
				WithServerLogger(&testLogger{}),
			},
			asserts: []optionAssertion{
				assertOptionWithServerTracer(opentracing.NoopTracer{}),
				assertOptionWithServerLogger(&testLogger{}),
				assertOptionWithServerConfig(defaultConfig()),
			},
		},
		{
			name: "All",
			opts: []Option{
				WithServerTracer(opentracing.GlobalTracer()),
				WithServerLogger(&testLogger{}),
			},
			asserts: []optionAssertion{
				assertOptionWithServerTracer(opentracing.GlobalTracer()),
				assertOptionWithServerLogger(&testLogger{}),
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
