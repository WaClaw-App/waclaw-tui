package rpc

import (
	"encoding/json"

	"github.com/WaClaw-App/waclaw/pkg/protocol"
)

// handleKeyPress processes a "key_press" event from the TUI.
func (s *Server) handleKeyPress(params map[string]any) (any, error) {
	evt := protocol.KeyPressEvent{
		Key:    stringField(params, "key"),
		Screen: protocol.ScreenID(stringField(params, "screen")),
		State:  protocol.StateID(stringField(params, "state")),
	}

	if err := s.engine.HandleKeyPress(evt); err != nil {
		return nil, err
	}

	return map[string]any{"status": "ok"}, nil
}

// handleAction processes an "action" event from the TUI.
func (s *Server) handleAction(params map[string]any) (any, error) {
	evt := protocol.ActionEvent{
		Action: stringField(params, "action"),
		Screen: protocol.ScreenID(stringField(params, "screen")),
		Params: mapField(params, "params"),
	}

	if err := s.engine.HandleAction(evt); err != nil {
		return nil, err
	}

	return map[string]any{"status": "ok"}, nil
}

// handleRequest processes a "request" event from the TUI.
func (s *Server) handleRequest(params map[string]any) (any, error) {
	evt := protocol.RequestEvent{
		Type:   stringField(params, "type"),
		Screen: protocol.ScreenID(stringField(params, "screen")),
		Params: mapField(params, "params"),
	}

	return s.engine.HandleRequest(evt)
}

// stringField extracts a string field from a params map.
// Returns empty string if the field is missing or not a string.
func stringField(params map[string]any, key string) string {
	v, ok := params[key]
	if !ok {
		return ""
	}
	switch val := v.(type) {
	case string:
		return val
	case json.Number:
		return string(val)
	default:
		return ""
	}
}

// mapField extracts a map[string]any field from a params map.
// Returns nil if the field is missing or not a map.
func mapField(params map[string]any, key string) map[string]any {
	v, ok := params[key]
	if !ok {
		return nil
	}
	if m, ok := v.(map[string]any); ok {
		return m
	}
	return nil
}
