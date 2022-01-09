package adderall

//go:generate mkdir -p internal/mock/client/
//go:generate go run github.com/golang/mock/mockgen --source=client/logger.go --destination=internal/mock/client/logger.go
