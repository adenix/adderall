package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/opentracing-contrib/go-stdlib/nethttp"
	"github.com/opentracing/opentracing-go"
	httpSwagger "github.com/swaggo/http-swagger"
	"go.adenix.dev/adderall/internal/pointer"
)

// Server represents a HTTP server. Server is instrumented with OpenTracing,
// logging, monitoring endpoints, and optionally an OpenAPI v2 endpoint.
type Server struct {
	Router         Handler
	tracer         opentracing.Tracer
	logger         Logger
	config         Config
	livenessCheck  func(http.HandlerFunc) http.HandlerFunc
	readinessCheck func(http.HandlerFunc) http.HandlerFunc
	healthCheck    func(http.HandlerFunc) http.HandlerFunc
}

// Serve sets up a http server and begins listening
func (s *Server) Serve(ctx context.Context) error {
	handler := s.getHandler(ctx)
	port := s.config.Port
	if port == nil || *port < 1 {
		port = pointer.IntP(8080)
	}

	srvr := http.Server{
		Addr:         fmt.Sprintf(":%d", *port),
		Handler:      handler,
		ReadTimeout:  time.Duration(*s.config.ReadTimeoutMs) * time.Millisecond,
		WriteTimeout: time.Duration(*s.config.WriteTimeoutMs) * time.Millisecond,
	}

	errs := make(chan error)
	go func() {
		if err := srvr.ListenAndServe(); err != http.ErrServerClosed {
			s.logger.ErrorCtx(ctx, "server failed to start up", "error", err)
			errs <- err
		} else {
			errs <- nil
		}
	}()

	s.logger.InfoCtx(ctx, "server started successfully", "port", port)

	go func() {
		errs <- s.gracefulShutdown(ctx, &srvr)
	}()

	return <-errs
}

// addSwagger configures and adds handers for an OpenAPI file and the Swagger UI
func (s *Server) addSwagger(r Handler) {
	swaggerFileLocation := "/swagger.json"
	if s.config.SwaggerFile != nil && len(*s.config.SwaggerFile) > 0 {
		swaggerFileLocation = *s.config.SwaggerFile
	}

	if _, err := os.Stat(swaggerFileLocation); err != nil {
		s.logger.InfoCtx(context.Background(), "swagger not added", "location", swaggerFileLocation, "error", err)
		return
	}

	r.HandleFunc(swaggerFileLocation, func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, swaggerFileLocation)
	})

	swaggerUIHandler := httpSwagger.Handler(
		httpSwagger.URL(swaggerFileLocation),
	)

	r.HandleFunc("/swagger", func(rw http.ResponseWriter, r *http.Request) {
		http.Redirect(rw, r, "/swagger/", http.StatusMovedPermanently)
	})
	r.HandleFunc("/swagger/", swaggerUIHandler)
	r.HandleFunc("/swagger/*", swaggerUIHandler)
}

// ServeHTTP is used to satisfy http.Handler interface, primarily to pass to test recorder.
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.getHandler(context.Background()).ServeHTTP(w, r)
}

// ProfilingMiddleware ...
func (s *Server) profilingMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			next.ServeHTTP(w, r)
			s.logger.DebugCtx(r.Context(), "http path response time",
				"path", r.URL.EscapedPath(),
				"method", r.Method,
				"time", time.Since(start),
			)
		}
		return http.HandlerFunc(fn)
	}
}

// TracingMiddleware ...
func (s *Server) tracingMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return nethttp.Middleware(s.tracer, next)
	}
}

// TimeoutMiddleware ...
func (s *Server) timeoutMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.TimeoutHandler(next, time.Duration(*s.config.RequestTimeoutSec)*time.Second, "timeout")
	}
}

func (s *Server) getHandler(ctx context.Context) http.Handler {
	var h http.Handler = s.Router
	h = s.timeoutMiddleware()(h)
	h = s.tracingMiddleware()(h)
	h = s.profilingMiddleware()(h)
	return h
}

func (s *Server) gracefulShutdown(ctx context.Context, server *http.Server) error {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	sig := <-quit
	s.logger.InfoCtx(ctx, "signal received", "signal", sig)

	shutdownDelaySeconds := s.config.ShutdownDelaySeconds
	if shutdownDelaySeconds == nil || *shutdownDelaySeconds < 1 {
		shutdownDelaySeconds = pointer.IntP(5)
	}
	timeout := time.Duration(*shutdownDelaySeconds) * time.Second

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		s.logger.ErrorCtx(
			ctx,
			"error while gracefully shutting down server, forcing shutdown because of error",
			"err", err)
		return err
	}
	s.logger.InfoCtx(ctx, "server exited successfully")
	return nil
}

func (s *Server) getLivenessHandler() http.HandlerFunc {
	dflt := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})

	if s.livenessCheck != nil {
		return s.livenessCheck(dflt)
	}
	return dflt
}

func (s *Server) getReadinessHandler() http.HandlerFunc {
	defult := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})

	if s.readinessCheck != nil {
		return s.readinessCheck(defult)
	}
	return defult
}

func (s *Server) getHealthCheckHandler() http.HandlerFunc {

	defult := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK!"))
	})

	if s.healthCheck != nil {
		return s.healthCheck(defult)
	}
	return defult
}

// Handler is the interface to handle HTTP requests
type Handler interface {
	http.Handler
	HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request))
}

// Router configures and provides a Handler
type Router func() Handler

// Config contains options for a Server
type Config struct {
	Port                 *int
	ReadTimeoutMs        *int
	WriteTimeoutMs       *int
	RequestTimeoutSec    *int
	ShutdownDelaySeconds *int
	SwaggerFile          *string
}

// defaultConfig provides a Config initalized with default values
func defaultConfig() Config {
	return Config{
		Port:                 pointer.IntP(8080),
		ReadTimeoutMs:        pointer.IntP(10000),
		WriteTimeoutMs:       pointer.IntP(10000),
		RequestTimeoutSec:    pointer.IntP(10),
		ShutdownDelaySeconds: pointer.IntP(5),
		SwaggerFile:          pointer.StringP("/swagger.json"),
	}
}
