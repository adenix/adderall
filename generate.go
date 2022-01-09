package adderall

//go:generate mkdir -p capsules/client/mock capsules/server/mock
//go:generate go run github.com/golang/mock/mockgen --package=mock_client --destination=capsules/client/mock/mock.go go.adenix.dev/adderall/capsules/client Factory
//go:generate go run github.com/golang/mock/mockgen --package=mock_server --destination=capsules/server/mock/mock.go go.adenix.dev/adderall/capsules/server Factory,Handler
