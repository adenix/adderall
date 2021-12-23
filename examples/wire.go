//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
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
	config.NewAppConfig,
	logger.NewLogger,
	NewTracer,
	NewServerConfig,
	NewServerFactory,
	NewClientConfig,
	NewClientFactory,
)
