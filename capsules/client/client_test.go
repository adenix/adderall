package client

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	mock "go.adenix.dev/adderall/internal/mock/tracing"
)

type assertResponse func(t *testing.T, response *http.Response)

func TestDo(t *testing.T) {
	tests := []struct {
		name    string
		action  func(c *Client, url, contentType string, body io.Reader) (*http.Response, error)
		method  string
		status  int
		handler http.HandlerFunc
		body    io.Reader
		err     bool
		asserts []assertResponse
	}{
		{
			name: http.MethodGet,
			action: func(c *Client, url, _ string, _ io.Reader) (*http.Response, error) {
				return c.Get(url)
			},
			method:  http.MethodGet,
			status:  http.StatusOK,
			handler: newHandlerFunc(http.MethodGet, http.StatusOK),
			asserts: []assertResponse{
				assertResponseBody(http.MethodGet),
				assertResponseStatus(http.StatusOK),
			},
		},
		{
			name: http.MethodPost,
			action: func(c *Client, url, contentType string, body io.Reader) (*http.Response, error) {
				return c.Post(url, contentType, body)
			},
			method:  http.MethodPost,
			status:  http.StatusCreated,
			handler: newHandlerFunc(http.MethodPost, http.StatusCreated),
			asserts: []assertResponse{
				assertResponseBody(http.MethodPost),
				assertResponseStatus(http.StatusCreated),
			},
		},
		{
			name: http.MethodPut,
			action: func(c *Client, url, contentType string, body io.Reader) (*http.Response, error) {
				request, err := http.NewRequest(http.MethodPut, url, body)
				if err != nil {
					return nil, err
				}
				return c.Do(request)
			},
			method:  http.MethodPut,
			status:  http.StatusAccepted,
			handler: newHandlerFunc(http.MethodPut, http.StatusAccepted),
			asserts: []assertResponse{
				assertResponseBody(http.MethodPut),
				assertResponseStatus(http.StatusAccepted),
			},
		},
		{
			name: http.MethodDelete,
			action: func(c *Client, url, contentType string, body io.Reader) (*http.Response, error) {
				request, err := http.NewRequest(http.MethodDelete, url, body)
				if err != nil {
					return nil, err
				}
				return c.Do(request)
			},
			method:  http.MethodDelete,
			status:  http.StatusNoContent,
			handler: newHandlerFunc("", http.StatusNoContent),
			asserts: []assertResponse{
				assertResponseBody(""),
				assertResponseStatus(http.StatusNoContent),
			},
		},
		{
			name: "Error",
			action: func(c *Client, url, contentType string, body io.Reader) (*http.Response, error) {
				ctx, cancel := context.WithTimeout(context.Background(), 0*time.Second)
				defer cancel()
				request, err := http.NewRequest(http.MethodGet, url, body)
				if err != nil {
					return nil, err
				}
				request = request.WithContext(ctx)
				return c.Do(request)
			},
			method:  http.MethodGet,
			status:  http.StatusOK,
			handler: newHandlerFunc("", http.StatusOK),
			asserts: []assertResponse{},
			err:     true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctx := context.Background()

			ts := httptest.NewServer(test.handler)
			defer ts.Close()

			tracer := mockTrackerWithExpect(ctx, t, test.method, ts.URL, test.status, test.err)
			c := NewFactory(WithTracer(tracer)).Create()

			actual, err := test.action(c, ts.URL, "text/plain", test.body)
			if err == nil {
				defer func() {
					_ = actual.Body.Close()
				}()
				if test.err {
					t.Error("error expected")
				}
			} else if !test.err {
				t.Errorf("unexpected error: %s", err)
			}

			for _, assert := range test.asserts {
				assert(t, actual)
			}
		})
	}
}

func mockTrackerWithExpect(ctx context.Context, t *testing.T, method, url string, status int, err bool) opentracing.Tracer {
	controller := gomock.NewController(t)
	tracer := mock.NewMockTracer(controller)
	span := mock.NewMockSpan(controller)

	tracer.
		EXPECT().
		StartSpan(gomock.Eq("http-request"), gomock.Eq(ext.SpanKindRPCClient)).
		Return(span)

	span.EXPECT().Tracer().AnyTimes()

	span.EXPECT().SetTag(gomock.Eq(string(ext.HTTPMethod)), gomock.Eq(method))
	span.EXPECT().SetTag(gomock.Eq(string(ext.HTTPUrl)), gomock.Eq(url))
	span.EXPECT().Context()
	tracer.EXPECT().Inject(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(nil)
	if err {
		span.EXPECT().SetTag(gomock.Eq("error"), gomock.Any())
	} else {
		span.EXPECT().SetTag(gomock.Eq(string(ext.HTTPStatusCode)), gomock.Eq(uint16(status)))
	}
	span.EXPECT().Finish()

	return tracer
}

func newHandlerFunc(response string, status int) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		rw.WriteHeader(status)
		_, _ = fmt.Fprint(rw, response)
	}
}

func assertResponseBody(expected string) assertResponse {
	return func(t *testing.T, response *http.Response) {
		bytes, err := io.ReadAll(response.Body)
		if err != nil {
			t.Errorf("unexpected error: %q", err)
		}

		actual := string(bytes)
		if string(actual) != expected {
			t.Errorf("expected %q, got %q", expected, actual)
		}
	}
}

func assertResponseStatus(expected int) assertResponse {
	return func(t *testing.T, response *http.Response) {
		if response.StatusCode != expected {
			t.Errorf("expected %d, got %d", expected, response.StatusCode)
		}
	}
}
