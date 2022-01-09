package client

import (
	"github.com/opentracing/opentracing-go"
	"go.adenix.dev/adderall/internal/pointer"
)

// Option interface to identify functional options
type Option func(c *Client)

// WithClientLogger provides an Option to provide a logger to be used by the
// Client
func WithClientLogger(l Logger) Option {
	return func(c *Client) {
		c.logger = l
	}
}

// WithClientTracer provides an Option to provide a tracer to used by the Client
func WithClientTracer(t opentracing.Tracer) Option {
	return func(c *Client) {
		c.tracer = t
	}
}

// WithTimeoutMs provides an Option to provide the maximum duration in
// milliseconds to wait for a request to finish.
// Defaults to 3 seconds
func WithTimeoutMs(t int) Option {
	return func(c *Client) {
		c.config.TimeoutMs = pointer.IntP(t)
	}
}

// WithRetryWaitMinMs provides an Option to provide the minimum duration in
// milliseconds to wait before retrying a request.
// Defaults to 3 seconds
func WithRetryWaitMinMs(t int) Option {
	return func(c *Client) {
		c.config.RetryWaitMinMs = pointer.IntP(t)
	}
}

// WithRetryMax provides an Option to provide the maximum number of time to
// retry a request.
// Defaults to 5
func WithRetryMax(r int) Option {
	return func(c *Client) {
		c.config.RetryMax = pointer.IntP(r)
	}
}

// FactoryOption interface to identify functional options
type FactoryOption interface{ apply(p *factory) }

// WithLogger provides an Option to provide a logger implementation.
// Defaults to Noop
func WithLogger(l Logger) FactoryOption { return factoryOptionLogger{logger: l} }

// WithTracer provides an Option to provide a tracer implementation.
// Defaults to Noop
func WithTracer(t opentracing.Tracer) FactoryOption { return factoryOptionTracer{tracer: t} }

// WithConfig provides an Option to provide a server configuration.
func WithConfig(c Config) FactoryOption { return factoryOptionConfig{c} }

type factoryOptionLogger struct{ logger Logger }

func (l factoryOptionLogger) apply(f *factory) {
	if l.logger != nil {
		f.logger = l.logger
	}
}

type factoryOptionTracer struct{ tracer opentracing.Tracer }

func (t factoryOptionTracer) apply(f *factory) {
	if t.tracer != nil {
		f.tracer = t.tracer
	}
}

type factoryOptionConfig struct{ config Config }

func (c factoryOptionConfig) apply(f *factory) {
	if c.config.TimeoutMs != nil {
		f.config.TimeoutMs = c.config.TimeoutMs
	}
	if c.config.RetryWaitMinMs != nil {
		f.config.RetryWaitMinMs = c.config.RetryWaitMinMs
	}
	if c.config.RetryMax != nil {
		f.config.RetryMax = c.config.RetryMax
	}
}
