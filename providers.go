package adderall

import (
	"go.adenix.dev/adderall/capsules/client"
	"go.adenix.dev/adderall/capsules/server"
)

// NewServerFactory provides a server.Factory given a slice of
// server.FactoryOption.
//
// This function is intended to be used with github.com/google/wire
func NewServerFactory(options []server.FactoryOption) server.Factory {
	return server.NewFactory(options...)
}

// NewClientFactory provides a client.Factory given a slice of
// client.FactoryOption.
//
// This function is intended to be used with github.com/google/wire
func NewClientFactory(options []client.FactoryOption) client.Factory {
	return client.NewFactory(options...)
}
