package adderall

import (
	"go.adenix.dev/adderall/client"
	"go.adenix.dev/adderall/server"
)

func NewServerFactory(options []server.FactoryOption) server.Factory {
	return server.NewFactory(options...)
}

func NewClientFactory(options []client.FactoryOption) client.Factory {
	return client.NewFactory(options...)
}
