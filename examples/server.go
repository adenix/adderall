package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"go.adenix.dev/adderall/client"
	"go.adenix.dev/adderall/server"
)

// This is a very crude representation of what each application looks like.
// Everything from here and underneath is this application domain and should be well tested.
func NewServer(service MyService) *server.Server {
	s := service.ServerFactory.Create()

	s.Router.HandleFunc("/hello", helloHandler())

	return s
}

type MyService struct {
	ServerFactory server.Factory

	HTTPConfig         client.Config
	HTTPClientProvider client.Provider
}

func helloHandler() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		name := "World"
		if ok := r.URL.Query().Has("name"); ok {
			name = r.URL.Query().Get("name")
		}
		resp := struct {
			Message string `json:"message"`
		}{
			Message: fmt.Sprintf("Hello, %s!", name),
		}
		json.NewEncoder(rw).Encode(resp)
	}
}
