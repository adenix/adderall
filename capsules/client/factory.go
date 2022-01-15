package client

import (
	"time"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/opentracing/opentracing-go"
)

// Factory is the interface to create Clients
type Factory interface {
	Create(options ...Option) *Client
}

type factory struct {
	tracer opentracing.Tracer
	logger Logger
	config Config
}

// NewFactory instantiates a Client Factory. FactoryOption can be passed to
// overwrite default configurations.
func NewFactory(options ...FactoryOption) Factory {
	f := &factory{
		tracer: opentracing.NoopTracer{},
		logger: NoopLogger{},
		config: defaultConfig(),
	}

	for _, option := range options {
		if option != nil {
			option.apply(f)
		}
	}

	return f
}

// Create instantiates a Client. Factory configurations are passed to the
// Client but can be overwritten with passed in Options
func (f *factory) Create(options ...Option) *Client {

	c := &Client{
		tracer: f.tracer,
		logger: f.logger,
		config: f.config,
	}

	for _, option := range options {
		option(c)
	}

	rClient := retryablehttp.NewClient()
	rClient.Logger = c.logger

	rClient.RetryMax = 0
	if c.config.RetryMax != nil {
		rClient.RetryMax = *c.config.RetryMax
	}

	if c.config.RetryWaitMinMs != nil {
		rClient.RetryWaitMin = time.Duration(*c.config.RetryWaitMinMs) * time.Millisecond
	}

	c.Client = rClient.StandardClient()

	c.Client.Timeout = 10 * time.Second
	if c.config.TimeoutMs != nil {
		c.Client.Timeout = time.Duration(*c.config.TimeoutMs) * time.Millisecond
	}

	return c
}
