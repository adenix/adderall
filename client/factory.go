package client

import (
	"time"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/opentracing/opentracing-go"
)

type Factory interface {
	Create(options ...Option) *Client
}

type factory struct {
	tracer opentracing.Tracer
	logger Logger
	config Config
}

func NewFactory(options ...FactoryOption) Factory {
	f := &factory{
		tracer: opentracing.NoopTracer{},
		logger: NoopLogger{},
		config: defaultConfig(),
	}

	for _, option := range options {
		option(f)
	}

	return f
}

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

	rClient.RetryMax = c.config.RetryMax
	rClient.RetryWaitMin = time.Duration(c.config.RetryWaitMinMs) * time.Millisecond

	c.Client = rClient.StandardClient()
	c.Client.Timeout = time.Duration(c.config.TimeoutMs) * time.Millisecond

	return c
}
