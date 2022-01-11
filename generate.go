package adderall

//go:generate mkdir -p internal/mock/tracing
//go:generate ${GOBIN}/mockgen --package=tracing --destination=internal/mock/tracing/mock.go github.com/opentracing/opentracing-go Tracer,Span
