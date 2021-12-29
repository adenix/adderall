package server

import (
	"net/http"

	"github.com/opentracing/opentracing-go"
)

// Factory interface to create a server.
type Factory interface {
	Create(options ...Option) *Server
}

type factory struct {
	tracer     opentracing.Tracer
	logger     Logger
	config     Config
	routerFunc func() Handler
}

// NewFactory instantiates a new server Factory
func NewFactory(options ...FactoryOption) Factory {
	f := &factory{
		tracer:     opentracing.NoopTracer{},
		logger:     NoopLogger{},
		config:     defaultConfig(),
		routerFunc: func() Handler { return &http.ServeMux{} },
	}

	for _, option := range options {
		if option != nil {
			option.apply(f)
		}
	}

	return f
}

func (f *factory) Create(options ...Option) *Server {

	s := &Server{
		tracer: f.tracer,
		logger: f.logger,
		config: f.config,
		Router: f.routerFunc(),
	}

	for _, option := range options {
		option(s)
	}

	s.Router.HandleFunc("/live", s.getLivenessHandler())
	s.Router.HandleFunc("/ready", s.getReadinessHandler())
	s.Router.HandleFunc("/health", s.getHealthCheckHandler())

	s.addSwagger(s.Router)

	return s
}
