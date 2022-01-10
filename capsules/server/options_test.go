package server

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/opentracing/opentracing-go"
	"go.adenix.dev/adderall/internal/pointer"
)

type optionAssertion func(t *testing.T, s *Server)

type factoryOptionAssertion func(t *testing.T, f *factory)

func TestOption(t *testing.T) {
	c := Config{
		Port:                 pointer.IntP(5000),
		ReadTimeoutMs:        pointer.IntP(1000),
		RequestTimeoutSec:    pointer.IntP(50),
		ShutdownDelaySeconds: pointer.IntP(10),
		WriteTimeoutMs:       pointer.IntP(1000),
		SwaggerFile:          pointer.StringP("foo"),
	}

	tests := []struct {
		name   string
		server *Server
		op     Option
		assert optionAssertion
	}{
		{
			name:   "WithServerLogger",
			op:     WithServerLogger(&testLogger{}),
			assert: assertOptionWithServerLogger(&testLogger{}),
		},
		{
			name:   "WithServerTracer",
			op:     WithServerTracer(opentracing.NoopTracer{}),
			assert: assertOptionWithServerTracer(opentracing.NoopTracer{}),
		},
		{
			name:   "WithServerPort",
			op:     WithServerPort(4000),
			assert: assertOptionWithServerPort(4000),
		},
		{
			name:   "WithServerReadTimeout",
			op:     WithServerReadTimeout(2000),
			assert: assertOptionWithServerReadTimeout(2000),
		},
		{
			name:   "WithServerWriteTimeout",
			op:     WithServerWriteTimeout(2000),
			assert: assertOptionWithServerWriteTimeout(2000),
		},
		{
			name:   "WithShutdownDelaySeconds",
			op:     WithShutdownDelaySeconds(20),
			assert: assertOptionWithShutdownDelaySeconds(20),
		},
		{
			name:   "WithHealthCheck",
			op:     WithHealthCheck(newHandlerFunc("HEALTH")),
			assert: assertOptionWithHealthCheck("HEALTH"),
		},
		{
			name:   "WithLivenessCheck",
			op:     WithLivenessCheck(newHandlerFunc("LIVE")),
			assert: assertOptionWithLivenessCheck("LIVE"),
		},
		{
			name:   "WithReadinessCheck",
			op:     WithReadinessCheck(newHandlerFunc("READY")),
			assert: assertOptionWithReadinessCheck("READY"),
		},
		{
			name:   "WithSwaggerFile",
			op:     WithSwaggerFile("bar"),
			assert: assertOptionWithSwaggerFile("bar"),
		},
		{
			name:   "WithServerRouter",
			op:     WithServerRouter(&testHandler{}),
			assert: assertOptionWithServerRouter(&testHandler{}),
		},
		{
			name:   "WithServerConfig-Blank",
			op:     WithServerConfig(Config{}),
			assert: assertOptionWithServerConfig(Config{}),
		},
		{
			name:   "WithServerConfig-Port",
			op:     WithServerConfig(Config{Port: c.Port}),
			assert: assertOptionWithServerConfig(Config{Port: c.Port}),
		},
		{
			name:   "WithServerConfig-ReadTimeoutMs",
			op:     WithServerConfig(Config{ReadTimeoutMs: c.ReadTimeoutMs}),
			assert: assertOptionWithServerConfig(Config{ReadTimeoutMs: c.ReadTimeoutMs}),
		},
		{
			name:   "WithServerConfig-RequestTimeoutSec",
			op:     WithServerConfig(Config{RequestTimeoutSec: c.RequestTimeoutSec}),
			assert: assertOptionWithServerConfig(Config{RequestTimeoutSec: c.RequestTimeoutSec}),
		},
		{
			name:   "WithServerConfig-ShutdownDelaySeconds",
			op:     WithServerConfig(Config{ShutdownDelaySeconds: c.ShutdownDelaySeconds}),
			assert: assertOptionWithServerConfig(Config{ShutdownDelaySeconds: c.ShutdownDelaySeconds}),
		},
		{
			name:   "WithServerConfig-WriteTimeoutMs",
			op:     WithServerConfig(Config{WriteTimeoutMs: c.WriteTimeoutMs}),
			assert: assertOptionWithServerConfig(Config{WriteTimeoutMs: c.WriteTimeoutMs}),
		},
		{
			name:   "WithServerConfig-SwaggerFile",
			op:     WithServerConfig(Config{SwaggerFile: c.SwaggerFile}),
			assert: assertOptionWithServerConfig(Config{SwaggerFile: c.SwaggerFile}),
		},
		{
			name:   "WithServerConfig-All",
			op:     WithServerConfig(c),
			assert: assertOptionWithServerConfig(c),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.server == nil {
				test.server = &Server{}
			}
			test.op(test.server)
			test.assert(t, test.server)
		})
	}
}

func TestFactoryOption(t *testing.T) {
	c := Config{
		Port:                 pointer.IntP(5000),
		ReadTimeoutMs:        pointer.IntP(1000),
		RequestTimeoutSec:    pointer.IntP(50),
		ShutdownDelaySeconds: pointer.IntP(10),
		WriteTimeoutMs:       pointer.IntP(1000),
		SwaggerFile:          pointer.StringP("foo"),
	}

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
			name:   "WithRouter",
			op:     WithRouter(func() Handler { return &testHandler{} }),
			assert: assertOptionWithRouter(&testHandler{}),
		},
		{
			name:   "WithConfig-Blank",
			op:     WithConfig(Config{}),
			assert: assertFactoryOptionWithConfig(Config{}),
		},
		{
			name:   "WithConfig-Port",
			op:     WithConfig(Config{Port: c.Port}),
			assert: assertFactoryOptionWithConfig(Config{Port: c.Port}),
		},
		{
			name:   "WithConfig-ReadTimeoutMs",
			op:     WithConfig(Config{ReadTimeoutMs: c.ReadTimeoutMs}),
			assert: assertFactoryOptionWithConfig(Config{ReadTimeoutMs: c.ReadTimeoutMs}),
		},
		{
			name:   "WithConfig-RequestTimeoutSec",
			op:     WithConfig(Config{RequestTimeoutSec: c.RequestTimeoutSec}),
			assert: assertFactoryOptionWithConfig(Config{RequestTimeoutSec: c.RequestTimeoutSec}),
		},
		{
			name:   "WithConfig-ShutdownDelaySeconds",
			op:     WithConfig(Config{ShutdownDelaySeconds: c.ShutdownDelaySeconds}),
			assert: assertFactoryOptionWithConfig(Config{ShutdownDelaySeconds: c.ShutdownDelaySeconds}),
		},
		{
			name:   "WithConfig-WriteTimeoutMs",
			op:     WithConfig(Config{WriteTimeoutMs: c.WriteTimeoutMs}),
			assert: assertFactoryOptionWithConfig(Config{WriteTimeoutMs: c.WriteTimeoutMs}),
		},
		{
			name:   "WithConfig-SwaggerFile",
			op:     WithConfig(Config{SwaggerFile: c.SwaggerFile}),
			assert: assertFactoryOptionWithConfig(Config{SwaggerFile: c.SwaggerFile}),
		},
		{
			name:   "WithConfig-All",
			op:     WithConfig(c),
			assert: assertFactoryOptionWithConfig(c),
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

func newHandlerFunc(response string) func(http.HandlerFunc) http.HandlerFunc {
	return func(h http.HandlerFunc) http.HandlerFunc {
		return func(rw http.ResponseWriter, r *http.Request) {
			fmt.Fprint(rw, response)
		}
	}
}

func assertOptionWithServerLogger(expected Logger) optionAssertion {
	return func(t *testing.T, s *Server) {
		if s.logger == nil {
			t.Errorf("expected %T, got nil", expected)
		} else if reflect.TypeOf(s.logger) != reflect.TypeOf(expected) {
			t.Errorf("expected type %T, got %T", expected, s.logger)
		}
	}
}

func assertOptionWithServerTracer(expected opentracing.Tracer) optionAssertion {
	return func(t *testing.T, s *Server) {
		if s.tracer == nil {
			t.Errorf("expected %T, got nil", expected)
		} else if reflect.TypeOf(s.tracer) != reflect.TypeOf(expected) {
			t.Errorf("expected type %T, got %T", expected, s.tracer)
		}
	}
}

func assertOptionWithServerConfig(expected Config) optionAssertion {
	return func(t *testing.T, s *Server) {
		if ok := reflect.DeepEqual(s.config, expected); !ok {
			t.Errorf("expexted %v, got %v", expected, s.config)
		}
	}
}

func assertOptionWithServerPort(expected int) optionAssertion {
	return func(t *testing.T, s *Server) {
		if s.config.Port == nil {
			t.Errorf("expected %d, got nil", expected)
		} else if *s.config.Port != expected {
			t.Errorf("expected %d, got %d", expected, *s.config.Port)
		}
	}
}

func assertOptionWithServerReadTimeout(expected int) optionAssertion {
	return func(t *testing.T, s *Server) {
		if s.config.ReadTimeoutMs == nil {
			t.Errorf("expected %d, got nil", expected)
		} else if *s.config.ReadTimeoutMs != expected {
			t.Errorf("expected %d, got %d", expected, *s.config.ReadTimeoutMs)
		}
	}
}

func assertOptionWithServerWriteTimeout(expected int) optionAssertion {
	return func(t *testing.T, s *Server) {
		if s.config.WriteTimeoutMs == nil {
			t.Errorf("expected %d, got nil", expected)
		} else if *s.config.WriteTimeoutMs != expected {
			t.Errorf("expected %d, got %d", expected, *s.config.WriteTimeoutMs)
		}
	}
}

func assertOptionWithShutdownDelaySeconds(expected int) optionAssertion {
	return func(t *testing.T, s *Server) {
		if s.config.ShutdownDelaySeconds == nil {
			t.Errorf("expected %d, got nil", expected)
		} else if *s.config.ShutdownDelaySeconds != expected {
			t.Errorf("expected %d, got %d", expected, *s.config.ShutdownDelaySeconds)
		}
	}
}

func assertOptionWithHealthCheck(expected string) optionAssertion {
	return func(t *testing.T, s *Server) {
		opt := assertOptionHandlerFunc(s.healthCheck, expected)
		opt(t, s)
	}
}

func assertOptionWithLivenessCheck(expected string) optionAssertion {
	return func(t *testing.T, s *Server) {
		opt := assertOptionHandlerFunc(s.livenessCheck, expected)
		opt(t, s)
	}
}

func assertOptionWithReadinessCheck(expected string) optionAssertion {
	return func(t *testing.T, s *Server) {
		opt := assertOptionHandlerFunc(s.readinessCheck, expected)
		opt(t, s)
	}
}

func assertOptionHandlerFunc(handler func(http.HandlerFunc) http.HandlerFunc, expected string) optionAssertion {
	return func(t *testing.T, s *Server) {
		f := handler(func(rw http.ResponseWriter, r *http.Request) {})
		ts := httptest.NewServer(f)
		defer ts.Close()

		res, err := http.Get(ts.URL)
		if err != nil {
			t.Errorf("unexpected error: %s", err)
		}
		defer res.Body.Close()

		bytes, err := io.ReadAll(res.Body)
		if err != nil {
			t.Errorf("unexpected error: %q", err)
		}

		actual := string(bytes)
		if string(actual) != expected {
			t.Errorf("expected %q, got %q", expected, actual)
		}
	}
}

func assertOptionWithSwaggerFile(expected string) optionAssertion {
	return func(t *testing.T, s *Server) {
		if s.config.SwaggerFile == nil {
			t.Errorf("expected %s, got nil", expected)
		} else if *s.config.SwaggerFile != expected {
			t.Errorf("expected %s, got %s", expected, *s.config.SwaggerFile)
		}
	}
}

func assertOptionWithServerRouter(expected Handler) optionAssertion {
	return func(t *testing.T, s *Server) {
		if s.Router == nil {
			t.Errorf("expected %T, got nil", expected)
		} else if reflect.TypeOf(s.Router) != reflect.TypeOf(expected) {
			t.Errorf("expected type %T, got %T", expected, s.Router)
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

func assertOptionWithRouter(expected Handler) factoryOptionAssertion {
	return func(t *testing.T, f *factory) {
		r := f.routerFunc()
		if r == nil {
			t.Errorf("expected %T, got nil", expected)
		} else if reflect.TypeOf(r) != reflect.TypeOf(expected) {
			t.Errorf("expected type %T, got %T", expected, r)
		}
	}
}
