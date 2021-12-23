package main

import (
	"github.com/opentracing/opentracing-go"
	"go.adenix.dev/adderall/client"
	"go.adenix.dev/adderall/config"
	"go.adenix.dev/adderall/logger"
	"go.adenix.dev/adderall/server"
	httptrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/net/http"
)

type FooServiceConfiguration struct {
	TimeoutMs      int
	RetryWaitMinMs int
	RetryMax       int
}

func NewFooServiceConfiguration(ac *config.AppConfig) (*FooServiceConfiguration, error) {
	cfg := &FooServiceConfiguration{}
	err := ac.Value(cfg)
	return cfg, err
}

func NewServerFactory(
	config server.Config,
	logger logger.Logger,
	tracer opentracing.Tracer,
) server.Factory {
	return server.NewFactory(
		server.WithLogger(logger),
		server.WithRouter(func() server.Handler {
			return httptrace.NewServeMux()
		}))
}

func NewServerConfig(ac *config.AppConfig) server.Config {
	return server.Config{}
}

func NewHttpServiceConfig(ac *config.AppConfig) client.Config {
	return client.Config{}
}

func NewTracer() opentracing.Tracer {
	return opentracing.GlobalTracer()
}
