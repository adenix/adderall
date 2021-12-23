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
func NewFactory(options ...FactoryOption) Factory { //Take params as option
	f := &factory{
		tracer:     opentracing.NoopTracer{},
		logger:     NoopLogger{},
		config:     defaultConfig(),
		routerFunc: func() Handler { return &http.ServeMux{} },
	}

	for _, option := range options {
		if option != nil {
			option(f)
		}
	}

	return f
}

func (f *factory) Create(options ...Option) *Server {

	srvr := &Server{
		tracer: f.tracer,
		logger: f.logger,
		config: f.config,
		Router: f.routerFunc(),
	}

	for _, option := range options {
		if option != nil {
			option(srvr)
		}
	}

	srvr.Router.HandleFunc("/live", srvr.getLivenessHandler())
	srvr.Router.HandleFunc("/ready", srvr.getReadinessHandler())
	srvr.Router.HandleFunc("/health", srvr.getHealthCheckHandler())

	srvr.addSwagger(srvr.Router)

	return srvr
}
