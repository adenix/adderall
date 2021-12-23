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

type FactoryOption func(f *factory)

func WithLogger(l Logger) FactoryOption {
	return func(f *factory) {
		f.logger = l
	}
}

func WithTracer(t opentracing.Tracer) FactoryOption {
	return func(f *factory) {
		f.tracer = t
	}
}

func WithConfig(c Config) FactoryOption {
	return func(f *factory) {
		if c.TimeoutMs != nil {
			f.config.TimeoutMs = c.TimeoutMs
		}
		if c.RetryWaitMinMs != nil {
			f.config.RetryWaitMinMs = c.RetryWaitMinMs
		}
		if c.RetryMax != nil {
			f.config.RetryMax = c.RetryMax
		}
	}
}
