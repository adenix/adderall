package client

import (
	"net/http"
	"time"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/opentracing/opentracing-go"
)

type Config struct {
	TimeoutMs      *int
	RetryWaitMinMs *int
	RetryMax       *int
}

type Provider struct {
	tracer opentracing.Tracer
	logger LeveledLogger
}

func NewClientProvider(tracer opentracing.Tracer, l LeveledLogger) Provider {
	return Provider{tracer, l}
}

func (p *Provider) GetClient(cfg Config) *http.Client {
	rClient := retryablehttp.NewClient()
	rClient.Logger = p.logger

	rClient.RetryMax = 0
	if cfg.RetryMax != nil {
		rClient.RetryMax = *cfg.RetryMax
	}

	if cfg.RetryWaitMinMs != nil {
		rClient.RetryWaitMin = time.Duration(*cfg.RetryWaitMinMs) * time.Millisecond
	}

	client := rClient.StandardClient()
	client.Timeout = 10 * time.Second
	if cfg.TimeoutMs != nil {
		client.Timeout = time.Duration(*cfg.TimeoutMs) * time.Millisecond
	}
	return client
}

func (p *Provider) GetWrappedClient(cfg Config) *Wrapper {
	return &Wrapper{Client: p.GetClient(cfg), tracer: p.tracer}
}
