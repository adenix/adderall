package main

import (
	"github.com/opentracing/opentracing-go"
	"go.adenix.dev/adderall/client"
	"go.adenix.dev/adderall/config"
	"go.adenix.dev/adderall/server"
	httptrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/net/http"
)

func ProvideServerRouter() server.Router {
	return func() server.Handler {
		return httptrace.NewServeMux()
	}
}

func ProvideServerConfig(ac *config.AppConfig) server.Config {
	return server.Config{}
}

func ProvideClientConfig(ac *config.AppConfig) client.Config {
	return client.Config{}
}

func ProvideTracer() opentracing.Tracer {
	return opentracing.GlobalTracer()
}
