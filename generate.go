package adderall

//go:generate mkdir -p internal/mock/tracing
//go:generate go run github.com/golang/mock/mockgen --package=tracing --destination=internal/mock/tracing/mock.go github.com/opentracing/opentracing-go Tracer,Span
