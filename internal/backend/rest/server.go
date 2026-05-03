// Package rest implements the REST API server that exposes the same scenario
// engine as HTTP endpoints. This is the future web frontend's API layer.
//
// The REST server mirrors the JSON-RPC methods as HTTP endpoints:
//
//	POST /api/v1/navigate   → scenario transition
//	POST /api/v1/update     → data update push
//	POST /api/v1/notify     → notification push
//	POST /api/v1/validate   → validation push
//	POST /api/v1/events/keypress  → key press event
//	POST /api/v1/events/action    → action event
//	POST /api/v1/events/request   → data request
//	GET  /api/v1/state      → current state snapshot
//	GET  /api/v1/screens    → available screens
package rest

import (
	"context"
	"log"

	"github.com/WaClaw-App/waclaw/internal/backend/engine"
	"github.com/WaClaw-App/waclaw/pkg/protocol"
	"github.com/WaClaw-App/waclaw/pkg/transport"
)

// Server wraps the HTTPTransport and wires up the scenario engine as the
// RPCHandler. It also implements transport.RPCHandler so all REST endpoints
// are automatically registered via transport.RegisterRPCHandlers.
type Server struct {
	http   *transport.HTTPTransport
	engine engine.ScenarioEngine
}

// NewServer creates a new REST API server with the given config and engine.
func NewServer(cfg transport.HTTPConfig, eng engine.ScenarioEngine) *Server {
	s := &Server{
		http:   transport.NewHTTPTransport(cfg),
		engine: eng,
	}

	// Register all REST endpoints using the shared handler infrastructure.
	transport.RegisterRPCHandlers(s.http, s)

	return s
}

// Start starts the HTTP server. Blocks until the server exits.
func (s *Server) Start() error {
	log.Println("[rest] starting REST API server")
	return s.http.Start()
}

// Shutdown gracefully stops the HTTP server.
func (s *Server) Shutdown(ctx context.Context) error {
	return s.http.Shutdown(ctx)
}

// HandleNavigate implements transport.RPCHandler.
func (s *Server) HandleNavigate(params map[string]any) (any, error) {
	screenID := protocol.ScreenID(strField(params, "screen"))
	stateID := protocol.StateID(strField(params, "state"))

	// Build navigation params from the REST request.
	navParams := make(map[string]any)
	for k, v := range params {
		if k != "screen" && k != "state" {
			navParams[k] = v
		}
	}

	return map[string]any{
		"screen": string(screenID),
		"state":  string(stateID),
		"params": navParams,
	}, nil
}

// HandleUpdate implements transport.RPCHandler.
func (s *Server) HandleUpdate(params map[string]any) (any, error) {
	screenID := protocol.ScreenID(strField(params, "screen"))
	return map[string]any{
		"screen": string(screenID),
		"params": params,
	}, nil
}

// HandleNotify implements transport.RPCHandler.
func (s *Server) HandleNotify(params map[string]any) (any, error) {
	return params, nil
}

// HandleValidate implements transport.RPCHandler.
func (s *Server) HandleValidate(params map[string]any) (any, error) {
	return params, nil
}

// HandleKeyPress implements transport.RPCHandler.
func (s *Server) HandleKeyPress(params map[string]any) (any, error) {
	evt := protocol.KeyPressEvent{
		Key:    strField(params, "key"),
		Screen: protocol.ScreenID(strField(params, "screen")),
		State:  protocol.StateID(strField(params, "state")),
	}

	if err := s.engine.HandleKeyPress(evt); err != nil {
		return nil, err
	}
	return map[string]any{"status": "ok"}, nil
}

// HandleAction implements transport.RPCHandler.
func (s *Server) HandleAction(params map[string]any) (any, error) {
	evt := protocol.ActionEvent{
		Action: strField(params, "action"),
		Screen: protocol.ScreenID(strField(params, "screen")),
		Params: mapField(params, "params"),
	}

	if err := s.engine.HandleAction(evt); err != nil {
		return nil, err
	}
	return map[string]any{"status": "ok"}, nil
}

// HandleRequest implements transport.RPCHandler.
func (s *Server) HandleRequest(params map[string]any) (any, error) {
	evt := protocol.RequestEvent{
		Type:   strField(params, "type"),
		Screen: protocol.ScreenID(strField(params, "screen")),
		Params: mapField(params, "params"),
	}

	return s.engine.HandleRequest(evt)
}

// GetState implements transport.RPCHandler.
func (s *Server) GetState() (any, error) {
	return s.engine.StateSnapshot(), nil
}

// GetScreens implements transport.RPCHandler.
func (s *Server) GetScreens() (any, error) {
	screens := protocol.AllScreens()
	result := make([]map[string]any, 0, len(screens))
	for _, s := range screens {
		result = append(result, map[string]any{
			"id":   string(s),
			"name": string(s),
		})
	}
	return result, nil
}

// strField extracts a string from a params map.
func strField(params map[string]any, key string) string {
	if v, ok := params[key]; ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

// mapField extracts a map[string]any from a params map.
func mapField(params map[string]any, key string) map[string]any {
	if v, ok := params[key]; ok {
		if m, ok := v.(map[string]any); ok {
			return m
		}
	}
	return nil
}
