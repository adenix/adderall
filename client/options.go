package client

import (
	"github.com/opentracing/opentracing-go"
	"go.adenix.dev/adderall/internal/pointer"
)

type Option func(c *Client)

func WithClientLogger(l Logger) Option {
	return func(c *Client) {
		c.logger = l
	}
}

func WithClientTracer(t opentracing.Tracer) Option {
	return func(c *Client) {
		c.tracer = t
	}
}

func WithTimeoutMs(t int) Option {
	return func(c *Client) {
		c.config.TimeoutMs = pointer.IntP(t)
	}
}

func WithRetryWaitMinMs(t int) Option {
	return func(c *Client) {
		c.config.RetryWaitMinMs = pointer.IntP(t)
	}
}

func WithRetryMax(r int) Option {
	return func(c *Client) {
		c.config.RetryMax = pointer.IntP(r)
	}
}

type FactoryOption interface{ apply(p *factory) }

// WithLogger provides option to provide a logger implementation. Noop is default
func WithLogger(l Logger) FactoryOptionLogger { return FactoryOptionLogger{logger: l} }

// WithTracer provides option to provide a tracer implementation. Noop is default
func WithTracer(t opentracing.Tracer) FactoryOptionTracer { return FactoryOptionTracer{tracer: t} }

// WithConfig provides option to provide a server configuration.
func WithConfig(c Config) FactoryOptionConfig { return FactoryOptionConfig{c} }

type FactoryOptionLogger struct{ logger Logger }

func (l FactoryOptionLogger) apply(f *factory) {
	if l.logger != nil {
		f.logger = l.logger
	}
}

type FactoryOptionTracer struct{ tracer opentracing.Tracer }

func (t FactoryOptionTracer) apply(f *factory) {
	if t.tracer != nil {
		f.tracer = t.tracer
	}
}

type FactoryOptionConfig struct{ config Config }

func (c FactoryOptionConfig) apply(f *factory) {
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
