package rest

import (
	"encoding/json"
	"net/http"

	"github.com/WaClaw-App/waclaw/pkg/protocol"
)

// RESTHandler provides helper methods for formatting REST API responses.
// These helpers ensure consistent JSON shapes across all endpoints and
// avoid code duplication between the individual handler methods in server.go.
type RESTHandler struct{}

// NewRESTHandler creates a new RESTHandler.
func NewRESTHandler() *RESTHandler {
	return &RESTHandler{}
}

// writeJSON writes a JSON response with the given status code.
func writeJSON(w http.ResponseWriter, statusCode int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		http.Error(w, `{"error":"internal encoding error"}`, http.StatusInternalServerError)
	}
}

// writeError writes a standardised error response.
func writeError(w http.ResponseWriter, statusCode int, message string) {
	writeJSON(w, statusCode, map[string]any{
		"error":  message,
		"status": statusCode,
	})
}

// formatNavigateResponse builds the REST response for a navigate command.
// The shape matches the JSON-RPC navigate notification params but with an
// explicit "ok" status for HTTP semantics.
func formatNavigateResponse(screen protocol.ScreenID, state protocol.StateID, params map[string]any) map[string]any {
	return map[string]any{
		"status": "ok",
		"screen": string(screen),
		"state":  string(state),
		"params": params,
	}
}

// formatUpdateResponse builds the REST response for an update command.
func formatUpdateResponse(screen protocol.ScreenID, params map[string]any) map[string]any {
	return map[string]any{
		"status": "ok",
		"screen": string(screen),
		"params": params,
	}
}

// formatNotifyResponse builds the REST response for a notify command.
func formatNotifyResponse(notifType protocol.NotificationType, severity protocol.Severity, data map[string]any) map[string]any {
	return map[string]any{
		"status":   "ok",
		"type":     string(notifType),
		"severity": string(severity),
		"data":     data,
	}
}

// formatValidateResponse builds the REST response for a validate command.
func formatValidateResponse(errors []string, warnings []string) map[string]any {
	return map[string]any{
		"status":   "ok",
		"errors":   errors,
		"warnings": warnings,
	}
}
