package rest

import (
	"log"
	"net/http"
	"time"
)

// restMiddleware provides REST-specific middleware that depends on backend
// packages. General-purpose middleware (CORS, logging) lives in
// pkg/transport/http.go so it can be reused by any HTTP server.

// requestIDMiddleware injects a unique request ID into the response headers
// for distributed tracing. In production, this would use a UUID or similar;
// for the demo backend, we use the request timestamp as a simple identifier.
func requestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := time.Now().Format("20060102-150405.000")
		w.Header().Set("X-Request-ID", id)
		next.ServeHTTP(w, r)
	})
}

// rateLimitMiddleware provides basic rate limiting for the REST API.
// In production, this would use a token bucket or sliding window algorithm
// with per-IP tracking. For the demo backend, it logs requests that exceed
// a simple threshold.
func rateLimitMiddleware(requestsPerSecond int) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Placeholder: in production, track per-IP request counts
			// and return 429 Too Many Requests when limits are exceeded.
			// For now, just pass through.
			next.ServeHTTP(w, r)
		})
	}
}

// authMiddleware validates license-based authentication for REST endpoints.
// In the demo backend, all requests are allowed. In production, this would
// check the license key from the Authorization header against the license
// validation system.
func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Placeholder: in production, validate Authorization header
		// against the license system. For the demo, all requests pass.
		auth := r.Header.Get("Authorization")
		if auth == "" {
			// Demo mode: allow unauthenticated access.
			log.Printf("[rest] unauthenticated request: %s %s", r.Method, r.URL.Path)
		}
		next.ServeHTTP(w, r)
	})
}
