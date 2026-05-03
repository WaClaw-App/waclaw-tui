// Package rpc implements the JSON-RPC 2.0 server that communicates with the
// TUI over stdio. The server receives events from the TUI (key_press, action,
// request) and pushes commands to the TUI (navigate, update, notify, validate).
//
// Architecture:
//
//	TUI ←── stdio ──→ RPC Server ──→ ScenarioEngine
//	                   (this pkg)     (scenario/)
//
// The server uses the engine.RPCPusher interface to avoid importing the
// scenario package directly, and the engine.ScenarioEngine interface so the
// scenario package doesn't import rpc/ either.
package rpc

import (
	"bufio"
	"encoding/json"
	"io"
	"log"
	"sync"

	"github.com/WaClaw-App/waclaw/internal/backend/engine"
	"github.com/WaClaw-App/waclaw/pkg/protocol"
	"github.com/WaClaw-App/waclaw/pkg/transport"
)

// Server is the JSON-RPC 2.0 server that bridges stdio I/O with the scenario
// engine. It implements engine.RPCPusher so the scenario engine can push
// navigate/update/notify/validate commands to the TUI.
type Server struct {
	transport *transport.StdioTransport
	writer    *bufio.Writer
	registry  *Registry
	engine    engine.ScenarioEngine
	mu        sync.Mutex
}

// NewServer creates a new RPC server reading from r and writing to w.
func NewServer(r io.Reader, w io.Writer, eng engine.ScenarioEngine) *Server {
	s := &Server{
		transport: transport.NewStdioTransport(r, w),
		writer:    bufio.NewWriter(w),
		registry:  NewRegistry(),
		engine:    eng,
	}

	// Register method handlers.
	s.registry.Register(protocol.MethodKeyPress, s.handleKeyPress)
	s.registry.Register(protocol.MethodAction, s.handleAction)
	s.registry.Register(protocol.MethodRequest, s.handleRequest)

	return s
}

// Serve starts the read loop, processing incoming JSON-RPC messages.
// It blocks until the reader returns EOF or an error.
func (s *Server) Serve() error {
	log.Println("[rpc] server started, waiting for messages")

	for {
		raw, err := s.transport.Receive()
		if err != nil {
			if err == io.EOF {
				log.Println("[rpc] client disconnected")
				return nil
			}
			log.Printf("[rpc] read error: %v", err)
			continue
		}

		s.handleMessage(raw)
	}
}

// handleMessage dispatches a raw JSON-RPC message to the appropriate handler.
func (s *Server) handleMessage(raw map[string]any) {
	method, _ := raw["method"].(string)
	if method == "" {
		// This is a response to a request we sent — nothing to do
		// because we don't track request IDs in the demo backend.
		return
	}

	id, hasID := raw["id"]

	// Extract params.
	params, _ := raw["params"].(map[string]any)
	if params == nil {
		params = make(map[string]any)
	}

	// Look up the handler.
	handler, ok := s.registry.Lookup(method)
	if !ok {
		log.Printf("[rpc] unknown method: %s", method)
		if hasID {
			s.sendError(id, protocol.ErrorCodeMethodNotFound, "method not found")
		}
		return
	}

	// Execute the handler.
	result, err := handler(params)
	if err != nil {
		log.Printf("[rpc] handler error for %s: %v", method, err)
		if hasID {
			s.sendError(id, protocol.ErrorCodeInternalError, err.Error())
		}
		return
	}

	// Send the response if the request had an ID.
	if hasID {
		s.sendResult(id, result)
	}
}

// PushNavigate implements engine.RPCPusher.
func (s *Server) PushNavigate(screen protocol.ScreenID, state protocol.StateID, params map[string]any) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Merge state into params.
	if params == nil {
		params = make(map[string]any)
	}
	params["screen"] = string(screen)
	params["state"] = string(state)

	return s.writeNotification(protocol.MethodNavigate, params)
}

// PushUpdate implements engine.RPCPusher.
func (s *Server) PushUpdate(screen protocol.ScreenID, params map[string]any) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if params == nil {
		params = make(map[string]any)
	}
	params["screen"] = string(screen)

	return s.writeNotification(protocol.MethodUpdate, params)
}

// PushNotify implements engine.RPCPusher.
func (s *Server) PushNotify(notifType protocol.NotificationType, severity protocol.Severity, data map[string]any) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	params := map[string]any{
		"type":     string(notifType),
		"severity": string(severity),
		"data":     data,
	}

	return s.writeNotification(protocol.MethodNotify, params)
}

// PushValidate implements engine.RPCPusher.
func (s *Server) PushValidate(errors []string, warnings []string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	params := map[string]any{
		"errors":   errors,
		"warnings": warnings,
	}

	return s.writeNotification(protocol.MethodValidate, params)
}

// writeNotification writes a JSON-RPC notification directly to the writer.
// Uses the protocol.NewNotification constructor for DRY compliance — no
// hardcoded version strings.
func (s *Server) writeNotification(method string, params any) error {
	notif := protocol.NewNotification(method, params)
	data, err := json.Marshal(notif)
	if err != nil {
		return err
	}
	data = append(data, '\n')
	_, err = s.writer.Write(data)
	if err != nil {
		return err
	}
	return s.writer.Flush()
}

// writeMessage writes an arbitrary JSON-RPC message to the writer.
func (s *Server) writeMessage(msg any) error {
	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	data = append(data, '\n')
	_, err = s.writer.Write(data)
	if err != nil {
		return err
	}
	return s.writer.Flush()
}

// sendResult sends a successful JSON-RPC response.
func (s *Server) sendResult(id any, result any) {
	resp := &protocol.Response{
		JSONRPC: protocol.Version,
		Result:  result,
	}

	switch v := id.(type) {
	case float64:
		resp.ID = int64(v)
	}

	if err := s.writeMessage(resp); err != nil {
		log.Printf("[rpc] write response error: %v", err)
	}
}

// sendError sends an error JSON-RPC response.
func (s *Server) sendError(id any, code int, message string) {
	resp := &protocol.Response{
		JSONRPC: protocol.Version,
		Error: &protocol.RPCError{
			Code:    code,
			Message: message,
		},
	}

	switch v := id.(type) {
	case float64:
		resp.ID = int64(v)
	}

	if err := s.writeMessage(resp); err != nil {
		log.Printf("[rpc] write error response: %v", err)
	}
}
