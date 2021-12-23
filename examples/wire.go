//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"go.adenix.dev/adderall/client"
	"go.adenix.dev/adderall/config"
	"go.adenix.dev/adderall/logger"
	"go.adenix.dev/adderall/server"
)

func InitializeServer() (*server.Server, func()) {
	wire.Build(
		commonSet,
		wire.Struct(new(MyService), "*"),
		NewServer,
	)
	return &server.Server{}, nil
}

var commonSet = wire.NewSet(
	NewServerConfig,
	NewServerFactory,
	config.NewAppConfig,
	logger.NewLogger,
	NewTracer,
	NewHttpServiceConfig,
	client.NewClientProvider,
	wire.Bind(new(client.Logger), new(logger.Logger)),
	client.NewLeveledLogger,
)
