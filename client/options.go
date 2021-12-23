package client

import "github.com/opentracing/opentracing-go"

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
		c.config.TimeoutMs = t
	}
}

func WithRetryWaitMinMs(t int) Option {
	return func(c *Client) {
		c.config.RetryWaitMinMs = t
	}
}

func WithRetryMax(r int) Option {
	return func(c *Client) {
		c.config.RetryMax = r
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
		f.config = c
	}
}
