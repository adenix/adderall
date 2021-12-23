package server

import (
	"net/http"

	"github.com/opentracing/opentracing-go"
)

// Option interface to identify functional options
type Option func(s *Server)

// WithServerLogger provides option to provide a logger to use while writing.
func WithServerLogger(l Logger) Option {
	return func(s *Server) {
		s.logger = l
	}
}

// WithServerTracer provides option to provide a tracer to use while writing.
func WithServerTracer(t opentracing.Tracer) Option {
	return func(s *Server) {
		s.tracer = t
	}
}

// WithServerConfig provides option to provide a server configuration.
func WithServerConfig(c Config) Option {
	return func(s *Server) {
		s.config = c
	}
}

// WithServerPort provides option to provide the port on which the server listens. Default is 80
func WithServerPort(p int) Option {
	return func(s *Server) {
		s.config.Port = p
	}
}

// WithServerReadTimeout provides option to provide the maximum duration in milliseconds for reading the entire
// request, including the body.
// defaults to 10 seconds
func WithServerReadTimeout(t int) Option {
	return func(s *Server) {
		s.config.ReadTimeoutMs = t
	}
}

// WithServerWriteTimeout provides option to provide the maximum duration in milliseconds before timing out writes of the response.
// defaults to 10 seconds
func WithServerWriteTimeout(t int) Option {
	return func(s *Server) {
		s.config.WriteTimeoutMs = t
	}
}

// WithShutdownDelaySeconds provides option to provide the duration by which server shutdown is delayed after receiving an os signal.
// defaults to 5 seconds
func WithShutdownDelaySeconds(d int) Option {
	return func(s *Server) {
		s.config.ShutdownDelaySeconds = d
	}
}

// WithHealthCheck provides option to provide additional health checks that are performed on health check probe.
func WithHealthCheck(f func(http.HandlerFunc) http.HandlerFunc) Option {
	return func(s *Server) {
		s.healthCheck = f
	}
}

// WithLivenessCheck provides option to provide additional liveness checks that are performed on liveness probe.
func WithLivenessCheck(f func(http.HandlerFunc) http.HandlerFunc) Option {
	return func(s *Server) {
		s.livenessCheck = f
	}
}

// WithReadinessCheck provides option to provide additional readiness checks that are performed on readiness probe.
func WithReadinessCheck(f func(http.HandlerFunc) http.HandlerFunc) Option {
	return func(s *Server) {
		s.readinessCheck = f
	}
}

// WithSwaggerFile provides option to provide the swagger file location. Default is '/swagger.json'
func WithSwaggerFile(f string) Option {
	return func(s *Server) {
		s.config.SwaggerFile = f
	}
}

// WithServerRouter provides option to provide hooks to use the http request to mutate the request context.
func WithServerRouter(r Handler) Option {
	return func(s *Server) {
		s.Router = r
	}
}

// FactoryOption interface to identify functional options
type FactoryOption func(f *factory)

// WithLogger provides option to provide a logger implementation. Noop is default
func WithLogger(l Logger) FactoryOption {
	return func(f *factory) {
		f.logger = l
	}
}

// WithTracer provides option to provide a tracer implementation. Noop is default
func WithTracer(t opentracing.Tracer) FactoryOption {
	return func(f *factory) {
		f.tracer = t
	}
}

// WithConfig provides option to provide a server configuration.
func WithConfig(c Config) FactoryOption {
	return func(f *factory) {
		f.config = c
	}
}

// WithRouter provides option to provide a function which returns which router will be used.
// By default we use http.ServeMux
func WithRouter(rf func() Handler) FactoryOption {
	return func(f *factory) {
		f.routerFunc = rf
	}
}
