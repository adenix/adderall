package main

import (
	"github.com/opentracing/opentracing-go"
	"go.adenix.dev/adderall/client"
	"go.adenix.dev/adderall/config"
	"go.adenix.dev/adderall/logger"
	"go.adenix.dev/adderall/server"
	httptrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/net/http"
)

func NewServerFactory(config server.Config, logger logger.Logger, tracer opentracing.Tracer) server.Factory {
	return server.NewFactory(
		// server.WithConfig(config),
		server.WithLogger(logger),
		server.WithTracer(tracer),
		server.WithRouter(func() server.Handler {
			return httptrace.NewServeMux()
		}))
}

func NewClientFactory(config client.Config, logger logger.Logger, tracer opentracing.Tracer) client.Factory {
	return client.NewFactory(
		// client.WithConfig(config),
		client.WithLogger(logger),
		client.WithTracer(tracer),
	)
}

func NewServerConfig(ac *config.AppConfig) server.Config {
	return server.Config{}
}

func NewClientConfig(ac *config.AppConfig) client.Config {
	return client.Config{}
}

func NewTracer() opentracing.Tracer {
	return opentracing.GlobalTracer()
}
