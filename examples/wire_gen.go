// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"github.com/google/wire"
	"go.adenix.dev/adderall/client"
	"go.adenix.dev/adderall/config"
	"go.adenix.dev/adderall/logger"
	"go.adenix.dev/adderall/server"
)

// Injectors from wire.go:

func InitializeServer() (*server.Server, func()) {
	appConfig := config.NewAppConfig()
	serverConfig := NewServerConfig(appConfig)
	tracer := NewTracer()
	loggerLogger, cleanup := logger.NewLogger(tracer)
	factory := NewServerFactory(serverConfig, loggerLogger, tracer)
	clientConfig := NewHttpServiceConfig(appConfig)
	leveledLogger := client.NewLeveledLogger(loggerLogger)
	provider := client.NewClientProvider(tracer, leveledLogger)
	myService := MyService{
		ServerFactory:      factory,
		HTTPConfig:         clientConfig,
		HTTPClientProvider: provider,
	}
	serverServer := NewServer(myService)
	return serverServer, func() {
		cleanup()
	}
}

// wire.go:

var commonSet = wire.NewSet(
	NewServerConfig,
	NewServerFactory, config.NewAppConfig, logger.NewLogger, NewTracer,
	NewHttpServiceConfig, client.NewClientProvider, wire.Bind(new(client.Logger), new(logger.Logger)), client.NewLeveledLogger,
)
