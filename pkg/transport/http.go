package transport

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

// HTTPTransport wraps an HTTP server that exposes the same RPC methods as
// REST endpoints. It is designed for the future web frontend and mirrors
// the JSON-RPC methods available over stdio.
type HTTPTransport struct {
	server *http.Server
	mux    *http.ServeMux
	addr   string
}

// HTTPConfig holds configuration for the REST server.
type HTTPConfig struct {
	// Addr is the listen address. Defaults to ":8080" if empty.
	Addr string
	// ReadTimeout is the maximum duration for reading the entire request,
	// including the body. Defaults to 10s if zero.
	ReadTimeout time.Duration
	// WriteTimeout is the maximum duration before timing out writes of the
	// response. Defaults to 10s if zero.
	WriteTimeout time.Duration
}

// defaults fills zero-valued fields with sensible defaults.
func (c *HTTPConfig) defaults() {
	if c.Addr == "" {
		c.Addr = ":8080"
	}
	if c.ReadTimeout == 0 {
		c.ReadTimeout = 10 * time.Second
	}
	if c.WriteTimeout == 0 {
		c.WriteTimeout = 10 * time.Second
	}
}

// NewHTTPTransport creates a REST API server with the given config.
// The server is not started; call Start to begin serving requests.
func NewHTTPTransport(cfg HTTPConfig) *HTTPTransport {
	cfg.defaults()

	mux := http.NewServeMux()

	// Wrap the mux with CORS and logging middleware.
	handler := corsMiddleware(loggingMiddleware(mux))

	srv := &http.Server{
		Addr:         cfg.Addr,
		Handler:      handler,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
	}

	return &HTTPTransport{
		server: srv,
		mux:    mux,
		addr:   cfg.Addr,
	}
}

// RegisterEndpoint registers a REST endpoint mapping to an RPC method.
// method is the HTTP method (GET, POST, etc.), path is the URL path
// (e.g. "/api/v1/navigate"), and handler is the function to handle the request.
func (t *HTTPTransport) RegisterEndpoint(method, path string, handler http.HandlerFunc) {
	t.mux.HandleFunc(method+" "+path, handler)
}

// Start starts the HTTP server. It blocks until the server encounters an error
// or is shut down via Shutdown.
func (t *HTTPTransport) Start() error {
	log.Printf("[transport] REST API server listening on %s", t.addr)
	return t.server.ListenAndServe()
}

// Shutdown gracefully shuts down the server without interrupting active
// connections. The provided context controls how long to wait for
// connections to finish before returning.
func (t *HTTPTransport) Shutdown(ctx context.Context) error {
	log.Printf("[transport] REST API server shutting down")
	return t.server.Shutdown(ctx)
}

// WriteJSON writes a JSON response with the given status code.
// It sets the Content-Type header to "application/json".
func WriteJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Printf("[transport] error encoding JSON response: %v", err)
	}
}

// WriteError writes a JSON error response in the standard format:
//
//	{"error": {"code": <code>, "message": "<message>"}}
func WriteError(w http.ResponseWriter, status int, code int, message string) {
	WriteJSON(w, status, map[string]any{
		"error": map[string]any{
			"code":    code,
			"message": message,
		},
	})
}

// ---------------------------------------------------------------------------
// Middleware
// ---------------------------------------------------------------------------

// corsMiddleware adds permissive CORS headers to every response. This allows
// the future web frontend (served from a different origin) to communicate
// with the REST API.
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
		w.Header().Set("Access-Control-Max-Age", "86400")

		// Preflight request: respond immediately.
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// loggingMiddleware logs each incoming HTTP request with method, path, and
// duration.
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("[transport] %s %s %s", r.Method, r.URL.Path, time.Since(start))
	})
}
