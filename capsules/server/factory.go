package server

import (
	"net/http"

	"github.com/opentracing/opentracing-go"
)

// Factory is the interface to create Servers
type Factory interface {
	Create(options ...Option) *Server
}

type factory struct {
	tracer     opentracing.Tracer
	logger     Logger
	config     Config
	routerFunc func() Handler
}

var _ Factory = (*factory)(nil)

// NewFactory instantiates a Server Factory. FactoryOption can be passed to
// overwrite default configurations.
func NewFactory(opts ...FactoryOption) Factory {
	f := &factory{
		tracer:     opentracing.NoopTracer{},
		logger:     NoopLogger{},
		config:     defaultConfig(),
		routerFunc: func() Handler { return &http.ServeMux{} },
	}

	for _, option := range opts {
		if option != nil {
			option.apply(f)
		}
	}

	return f
}

// Create instantiates a Server. Factory configurations are passed to the
// Server but can be overwritten with passed in Options
func (f *factory) Create(opts ...Option) *Server {

	s := &Server{
		tracer: f.tracer,
		logger: f.logger,
		config: f.config,
		Router: f.routerFunc(),
	}

	for _, option := range opts {
		option(s)
	}

	s.Router.HandleFunc("/live", s.getLivenessHandler())
	s.Router.HandleFunc("/ready", s.getReadinessHandler())
	s.Router.HandleFunc("/health", s.getHealthCheckHandler())

	s.addSwagger(s.Router)

	return s
}
