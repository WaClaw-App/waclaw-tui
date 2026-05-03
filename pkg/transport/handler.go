package transport

import (
	"encoding/json"
	"io"
	"net/http"
	"os"

	"github.com/WaClaw-App/waclaw/pkg/protocol"
)

// RPCHandler is the interface that REST handlers delegate to.
// The actual implementation lives in internal/backend/rest/.
// Each method corresponds to one of the RPC methods defined in pkg/protocol.
type RPCHandler interface {
	// HandleNavigate processes a "navigate" request from the backend.
	HandleNavigate(params map[string]any) (any, error)

	// HandleUpdate processes an "update" request from the backend.
	HandleUpdate(params map[string]any) (any, error)

	// HandleNotify processes a "notify" request from the backend.
	HandleNotify(params map[string]any) (any, error)

	// HandleValidate processes a "validate" request from the backend.
	HandleValidate(params map[string]any) (any, error)

	// HandleKeyPress processes a "key_press" event from the frontend.
	HandleKeyPress(params map[string]any) (any, error)

	// HandleAction processes an "action" event from the frontend.
	HandleAction(params map[string]any) (any, error)

	// HandleRequest processes a "request" event from the frontend.
	HandleRequest(params map[string]any) (any, error)

	// GetState returns the current full application state snapshot.
	GetState() (any, error)

	// GetScreens returns the list of screens and their metadata.
	GetScreens() (any, error)
}

// RegisterRPCHandlers registers all REST endpoint handlers that delegate to
// the given RPCHandler. It wires up the 9 REST endpoints defined in the
// endpoint mapping:
//
//	POST /api/v1/navigate
//	POST /api/v1/update
//	POST /api/v1/notify
//	POST /api/v1/validate
//	POST /api/v1/events/keypress
//	POST /api/v1/events/action
//	POST /api/v1/events/request
//	GET  /api/v1/state
//	GET  /api/v1/screens
func RegisterRPCHandlers(t *HTTPTransport, h RPCHandler) {
	// Backend → Frontend endpoints (POST, body carries params).
	t.RegisterEndpoint(http.MethodPost, "/api/v1/navigate", postHandler(h.HandleNavigate))
	t.RegisterEndpoint(http.MethodPost, "/api/v1/update", postHandler(h.HandleUpdate))
	t.RegisterEndpoint(http.MethodPost, "/api/v1/notify", postHandler(h.HandleNotify))
	t.RegisterEndpoint(http.MethodPost, "/api/v1/validate", postHandler(h.HandleValidate))

	// Frontend → Backend event endpoints (POST, body carries event params).
	t.RegisterEndpoint(http.MethodPost, "/api/v1/events/keypress", postHandler(h.HandleKeyPress))
	t.RegisterEndpoint(http.MethodPost, "/api/v1/events/action", postHandler(h.HandleAction))
	t.RegisterEndpoint(http.MethodPost, "/api/v1/events/request", postHandler(h.HandleRequest))

	// Read-only GET endpoints.
	t.RegisterEndpoint(http.MethodGet, "/api/v1/state", getHandler(h.GetState))
	t.RegisterEndpoint(http.MethodGet, "/api/v1/screens", getHandler(h.GetScreens))

	// Swagger spec endpoint.
	t.RegisterEndpoint(http.MethodGet, "/api/v1/swagger.json", serveSwagger)
}

// ---------------------------------------------------------------------------
// Helper handlers
// ---------------------------------------------------------------------------

// postHandler returns an http.HandlerFunc that decodes a JSON body into
// map[string]any, delegates to the provided handler function, and writes
// the result as JSON.
func postHandler(fn func(map[string]any) (any, error)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			WriteError(w, http.StatusBadRequest, protocol.ErrorCodeParseError,
				"failed to read request body")
			return
		}
		defer r.Body.Close()

		var params map[string]any
		if len(body) > 0 {
			if err := json.Unmarshal(body, &params); err != nil {
				WriteError(w, http.StatusBadRequest, protocol.ErrorCodeParseError,
					"invalid JSON in request body")
				return
			}
		}

		// Ensure params is never nil so handlers don't need nil checks.
		if params == nil {
			params = make(map[string]any)
		}

		result, err := fn(params)
		if err != nil {
			WriteError(w, http.StatusInternalServerError, protocol.ErrorCodeInternalError,
				err.Error())
			return
		}

		WriteJSON(w, http.StatusOK, result)
	}
}

// getHandler returns an http.HandlerFunc that delegates to the provided
// handler function and writes the result as JSON.
func getHandler(fn func() (any, error)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		result, err := fn()
		if err != nil {
			WriteError(w, http.StatusInternalServerError, protocol.ErrorCodeInternalError,
				err.Error())
			return
		}

		WriteJSON(w, http.StatusOK, result)
	}
}

// serveSwagger serves the OpenAPI spec from disk.
// Tries common file paths to locate pkg/api/openapi.yaml.
// In production builds, the spec can be embedded via //go:embed directives
// in the backend package; this handler provides a development fallback.
func serveSwagger(w http.ResponseWriter, r *http.Request) {
	// Search for the OpenAPI spec on disk.
	// Works in development when running `go run` from the project root.
	diskPaths := []string{
		"pkg/api/openapi.yaml",
		"openapi.yaml",
	}
	for _, p := range diskPaths {
		data, err := os.ReadFile(p)
		if err == nil {
			w.Header().Set("Content-Type", "application/yaml")
			w.Write(data)
			return
		}
	}

	// No spec available.
	WriteError(w, http.StatusNotFound, protocol.ErrorCodeMethodNotFound,
		"swagger spec not found; expected pkg/api/openapi.yaml")
}
