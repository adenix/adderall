package adderall

import (
	"github.com/google/wire"
	"github.com/opentracing/opentracing-go"
	"go.adenix.dev/adderall/client"
	"go.adenix.dev/adderall/logger"
	"go.adenix.dev/adderall/server"
)

func ProvideServerFactoryOptionWithConfig(config server.Config) server.FactoryOptionConfig {
	return server.WithConfig(config)
}

func ProvideServerFactoryOptionWithLogger(logger logger.Logger) server.FactoryOptionLogger {
	return server.WithLogger(logger)
}

func ProvideServerFactoryOptionWithTracer(tracer opentracing.Tracer) server.FactoryOptionTracer {
	return server.WithTracer(tracer)
}

func ProvideServerFactoryOptionWithRouter(rf server.Router) server.FactoryOptionRouter {
	return server.WithRouter(rf)
}

func ProvideServerFactoryOptions(config server.FactoryOptionConfig, logger server.FactoryOptionLogger, tracer server.FactoryOptionTracer, router server.FactoryOptionRouter) []server.FactoryOption {
	return []server.FactoryOption{config, logger, tracer, router}
}

func ProvideServerFactory(options []server.FactoryOption) server.Factory {
	return server.NewFactory(options...)
}

var ProvideServerFactorySet = wire.NewSet(
	ProvideServerFactoryOptionWithConfig,
	ProvideServerFactoryOptionWithLogger,
	ProvideServerFactoryOptionWithTracer,
	ProvideServerFactoryOptionWithRouter,
	ProvideServerFactoryOptions,
	ProvideServerFactory,
)

func ProvideClientFactoryOptionWithConfig(config client.Config) client.FactoryOptionConfig {
	return client.WithConfig(config)
}

func ProvideClientFactoryOptionWithLogger(logger logger.Logger) client.FactoryOptionLogger {
	return client.WithLogger(logger)
}

func ProvideClientFactoryOptionWithTracer(tracer opentracing.Tracer) client.FactoryOptionTracer {
	return client.WithTracer(tracer)
}

func ProvideClientFactoryOptions(config client.FactoryOptionConfig, logger client.FactoryOptionLogger, tracer client.FactoryOptionTracer) []client.FactoryOption {
	return []client.FactoryOption{config, logger, tracer}
}

func ProvideClientFactory(options []client.FactoryOption) client.Factory {
	return client.NewFactory(options...)
}

var ProvideClientFactorySet = wire.NewSet(
	ProvideClientFactoryOptionWithConfig,
	ProvideClientFactoryOptionWithLogger,
	ProvideClientFactoryOptionWithTracer,
	ProvideClientFactoryOptions,
	ProvideClientFactory,
)
