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
	c := service.ClientFactory.Create()

	s.Router.HandleFunc("/hello", helloHandler())
	s.Router.HandleFunc("/httpbin", httpbinHandler(c))

	return s
}

type MyService struct {
	ServerFactory server.Factory
	ClientFactory client.Factory
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

func httpbinHandler(client *client.Client) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		resp, err := client.Get("https://httpbin.org/get?foo=bar")
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		body := struct {
			Arguments map[string]string `json:"args"`
			Headers   map[string]string `json:"headers"`
			Origin    string            `json:"origin"`
			Url       string            `json:"url"`
		}{}

		json.NewDecoder(resp.Body).Decode(&body)
		json.NewEncoder(rw).Encode(&body)
	}
}
