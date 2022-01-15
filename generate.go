package adderall

//go:generate mkdir -p internal/mock/tracing internal/mock/zapcore
//go:generate ${GOBIN}/mockgen --package=tracing --destination=internal/mock/tracing/mock.go github.com/opentracing/opentracing-go Tracer,Span,SpanContext
//go:generate ${GOBIN}/mockgen --package=zapcore --destination=internal/mock/zapcore/mock.go go.uber.org/zap/zapcore PrimitiveArrayEncoder
