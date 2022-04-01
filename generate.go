package adderall

//go:generate mkdir -p mock/tracing mock/zapcore
//go:generate ${GOBIN}/mockgen --package=tracing --destination=mock/tracing/mock.go github.com/opentracing/opentracing-go Tracer,Span,SpanContext
//go:generate ${GOBIN}/mockgen --package=zapcore --destination=mock/zapcore/mock.go go.uber.org/zap/zapcore PrimitiveArrayEncoder
