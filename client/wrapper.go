package client

import (
	"net/http"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
)

type Wrapper struct {
	*http.Client
	tracer opentracing.Tracer
}

func (w *Wrapper) Do(request *http.Request) (*http.Response, error) {
	ctx := request.Context()

	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, w.tracer, "http-request", ext.SpanKindRPCClient)
	defer span.Finish()

	ext.HTTPMethod.Set(span, request.Method)
	ext.HTTPUrl.Set(span, request.URL.String())

	w.tracer.Inject(span.Context(), opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(request.Header))

	request = request.WithContext(ctx)

	resp, err := w.Client.Do(request)
	if err != nil {
		span.SetTag("error", true)
		return resp, err
	}
	ext.HTTPStatusCode.Set(span, uint16(resp.StatusCode))

	return resp, err
}
