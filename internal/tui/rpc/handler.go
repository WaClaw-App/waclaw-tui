package rpc

// Package rpc provides the JSON-RPC 2.0 client and message handler for the
// WaClaw TUI. The client connects to the backend binary over stdio and the
// handler translates incoming RPC messages into typed bus messages that the
// bubbletea Update loop can consume.
//
// DRY convention: all method name comparisons use protocol.Method* constants,
// never magic strings. All severity values are validated against
// protocol.IsValidSeverity and notification types against
// protocol.IsValidNotificationType before publishing.

import (
        "encoding/json"
        "fmt"

        "github.com/WaClaw-App/waclaw/internal/tui/bus"
        "github.com/WaClaw-App/waclaw/pkg/protocol"
)

// Handler translates raw JSON-RPC messages from the backend into typed bus
// messages. It is the single point of translation between the wire protocol
// and the TUI's internal event system — screens never see raw JSON.
//
// Note: Response routing (TUI request → backend response) is handled by
// Client.routeResponse, not by the Handler. The Handler only processes
// backend-initiated messages (navigate, update, notify, validate).
type Handler struct {
        bus *bus.Bus
}

// NewHandler creates a Handler that publishes translated messages to the given bus.
func NewHandler(b *bus.Bus) *Handler {
        return &Handler{bus: b}
}

// HandleMessage classifies and processes a raw JSON-RPC message received from
// the backend. The raw map is produced by transport.StdioTransport.Receive().
//
// Classification rules (per JSON-RPC 2.0 spec):
//   - Has "method" + "id" → backend-initiated request (e.g. navigate, update)
//   - Has "method" + no "id" → backend notification (e.g. notify, validate)
//   - Has "id" + no "method" → response to a TUI-initiated request (handled
//     by Client.routeResponse, NOT here)
//
// For backend-initiated requests and notifications, the handler extracts the
// method and params, translates them into the appropriate bus message type,
// and publishes to the bus.
func (h *Handler) HandleMessage(raw map[string]any, _ map[int64]chan *protocol.Response) {
        // Determine message classification from the JSON-RPC envelope.
        method, hasMethod := raw["method"]

        if !hasMethod {
                // Not a backend-initiated message — ignore.
                // Response routing is handled by Client.routeResponse.
                return
        }

        methodStr, ok := method.(string)
        if !ok {
                return // malformed — method must be a string
        }

        params := extractParams(raw)

        // Both backend requests and notifications are translated the same way.
        // The distinction (has "id" vs. not) only matters for the JSON-RPC
        // protocol layer — at the TUI level, all four methods push state.
        h.handleBackendPush(methodStr, params)
}

// handleBackendPush processes a backend-initiated JSON-RPC request or
// notification. The four documented backend methods are: navigate, update,
// notify, validate. All four push state to the TUI — the TUI never sends
// a JSON-RPC response back for these.
func (h *Handler) handleBackendPush(method string, params map[string]any) {
        switch method {
        case protocol.MethodNavigate:
                h.publishNavigate(params)
        case protocol.MethodUpdate:
                h.publishUpdate(params)
        case protocol.MethodNotify:
                h.publishNotify(params)
        case protocol.MethodValidate:
                h.publishValidate(params)
        // Unknown methods are silently ignored — no panic, no error, no bus message.
        }
}

// publishNavigate translates a navigate command into a bus.NavigateMsg.
// Params must contain a "screen" field with a valid ScreenID.
func (h *Handler) publishNavigate(params map[string]any) {
        screenID, ok := extractScreenID(params)
        if !ok {
                return // malformed — missing or invalid screen field
        }

        msg := bus.NavigateMsg{
                Screen: screenID,
                Params: params,
        }
        h.bus.Publish(msg)
}

// publishUpdate translates an update command into a bus.UpdateMsg.
// Params carry screen-specific data but the "screen" field identifies
// which screen should receive the update.
func (h *Handler) publishUpdate(params map[string]any) {
        screenID, ok := extractScreenID(params)
        if !ok {
                return
        }

        msg := bus.UpdateMsg{
                Screen: screenID,
                Params: params,
        }
        h.bus.Publish(msg)
}

// publishNotify translates a notify command into a bus.NotifyMsg.
// Params must contain "type" (NotificationType) and "severity" (Severity).
// Both are validated against protocol.IsValid* functions — invalid values
// are handled gracefully (invalid severity falls back to neutral, invalid
// notification type causes an early return).
func (h *Handler) publishNotify(params map[string]any) {
        notifType, _ := params[protocol.ParamType].(string)

        // Validate notification type against known types.
        if notifType != "" {
                nt := protocol.NotificationType(notifType)
                if !protocol.IsValidNotificationType(nt) {
                        return // invalid notification type — silently ignore
                }
        }

        severity := protocol.SeverityNeutral // default
        if sevStr, ok := params[protocol.ParamSeverity].(string); ok {
                s := protocol.Severity(sevStr)
                if protocol.IsValidSeverity(s) {
                        severity = s
                }
        }

        data, _ := params[protocol.ParamData].(map[string]any)
        if data == nil {
                data = make(map[string]any)
        }

        msg := bus.NotifyMsg{
                Type:     notifType,
                Severity: severity,
                Data:     data,
        }
        h.bus.Publish(msg)
}

// publishValidate translates a validate command into a bus.ValidateMsg.
// Params must contain "errors" ([]string) and "warnings" ([]string).
func (h *Handler) publishValidate(params map[string]any) {
        errors := extractStringSlice(params, protocol.ParamErrors)
        warnings := extractStringSlice(params, protocol.ParamWarnings)

        msg := bus.ValidateMsg{
                Errors:   errors,
                Warnings: warnings,
        }
        h.bus.Publish(msg)
}

// ---------------------------------------------------------------------------
// Internal helpers — keep these private; they are not part of the public API.
// ---------------------------------------------------------------------------

// extractParams returns the "params" field from a raw message as a
// map[string]any. Returns an empty map if params are missing or not a map.
func extractParams(raw map[string]any) map[string]any {
        params, ok := raw["params"]
        if !ok {
                return make(map[string]any)
        }

        m, ok := params.(map[string]any)
        if !ok {
                return make(map[string]any)
        }
        return m
}

// extractScreenID extracts and validates the "screen" field from params.
func extractScreenID(params map[string]any) (protocol.ScreenID, bool) {
        screen, ok := params[protocol.ParamScreen].(string)
        if !ok {
                return "", false
        }

        id := protocol.ScreenID(screen)
        if !protocol.IsValid(id) {
                return "", false
        }
        return id, true
}

// extractStringSlice extracts a string slice from the params map under the
// given key. Returns nil if the key is missing or the value is not a slice.
func extractStringSlice(params map[string]any, key string) []string {
        val, ok := params[key]
        if !ok {
                return nil
        }

        // JSON unmarshals []string as []any, not []string.
        slice, ok := val.([]any)
        if !ok {
                return nil
        }

        result := make([]string, 0, len(slice))
        for _, v := range slice {
                if s, ok := v.(string); ok {
                        result = append(result, s)
                }
        }
        return result
}

// toInt64 converts a JSON number to int64. JSON unmarshals numbers as
// float64 by default.
func toInt64(v any) (int64, bool) {
        switch n := v.(type) {
        case float64:
                return int64(n), true
        case int64:
                return n, true
        case json.Number:
                i, err := n.Int64()
                if err != nil {
                        return 0, false
                }
                return i, true
        }
        return 0, false
}

// parseResponse constructs a protocol.Response from a raw message map.
func parseResponse(raw map[string]any) *protocol.Response {
        resp := &protocol.Response{
                JSONRPC: protocol.Version,
        }

        if id, ok := toInt64(raw["id"]); ok {
                resp.ID = id
        }

        if result, ok := raw["result"]; ok {
                resp.Result = result
        }

        if errVal, ok := raw["error"]; ok {
                if errMap, ok := errVal.(map[string]any); ok {
                        rpcErr := &protocol.RPCError{}
                        if code, ok := toInt64(errMap["code"]); ok {
                                rpcErr.Code = int(code)
                        }
                        if msg, ok := errMap["message"].(string); ok {
                                rpcErr.Message = msg
                        }
                        if data, ok := errMap["data"]; ok {
                                rpcErr.Data = data
                        }
                        resp.Error = rpcErr
                }
        }

        return resp
}

// FormatError creates a human-readable description of an RPC error.
// Returns empty string if the error is nil.
func FormatError(err *protocol.RPCError) string {
        if err == nil {
                return ""
        }
        return fmt.Sprintf("RPC error %d: %s", err.Code, err.Message)
}
