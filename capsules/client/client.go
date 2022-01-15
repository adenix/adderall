package client

import (
	"net/http"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"go.adenix.dev/adderall/internal/pointer"
)

// Client represents a Doer. Client is instrumented with OpenTracing and logging
type Client struct {
	*http.Client
	tracer opentracing.Tracer
	logger Logger
	config Config
}

// Do executes an OpenTracking instrumented HTTP request
func (c *Client) Do(request *http.Request) (*http.Response, error) {
	ctx := request.Context()

	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, c.tracer, "http-request", ext.SpanKindRPCClient)
	defer span.Finish()

	ext.HTTPMethod.Set(span, request.Method)
	ext.HTTPUrl.Set(span, request.URL.String())

	_ = c.tracer.Inject(span.Context(), opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(request.Header))

	request = request.WithContext(ctx)

	resp, err := c.Client.Do(request)
	if err != nil {
		span.SetTag("error", true)
		return resp, err
	}
	ext.HTTPStatusCode.Set(span, uint16(resp.StatusCode))

	return resp, err
}

// Config contains options for a Client
type Config struct {
	TimeoutMs      *int
	RetryWaitMinMs *int
	RetryMax       *int
}

// defaultConfig provides a Config initialized with default values
func defaultConfig() Config {
	return Config{
		TimeoutMs:      pointer.IntP(3000),
		RetryWaitMinMs: pointer.IntP(3000),
		RetryMax:       pointer.IntP(5),
	}
}
