package rpc

import (
        "bufio"
        "encoding/json"
        "os"
        "sync"
        "testing"
        "time"

        "github.com/WaClaw-App/waclaw/internal/tui/bus"
        "github.com/WaClaw-App/waclaw/pkg/protocol"
)

// ---------------------------------------------------------------------------
// Integration tests: fake backend → TUI connection
// ---------------------------------------------------------------------------
//
// These tests verify that:
//  1. The client connects and reads messages from a fake backend.
//  2. The handler translates incoming RPC messages into correct bus messages.
//  3. TUI → Backend sends produce correctly formatted JSON-RPC messages.
//  4. Pending request/response routing works end-to-end.
//
// The fake backend writes JSON-RPC messages into an os.Pipe that the client
// reads. The client writes TUI events into another os.Pipe that the test
// reads. We use os.Pipe (not io.Pipe) because os.Pipe is buffered, allowing
// writes to complete without blocking on a simultaneous read.
// ---------------------------------------------------------------------------

// pipeConn wraps a pair of os.Pipe to simulate a bidirectional stdio
// connection between the fake backend and the TUI client.
type pipeConn struct {
        // backendWrite → clientRead: backend sends to TUI
        backendWrite *os.File
        clientRead   *os.File

        // clientWrite → backendRead: TUI sends to backend
        clientWrite *os.File
        backendRead *os.File
}

func newPipeConn(t *testing.T) *pipeConn {
        t.Helper()

        cr, cw, err := os.Pipe()
        if err != nil {
                t.Fatalf("create backend→client pipe: %v", err)
        }

        br, bw, err := os.Pipe()
        if err != nil {
                t.Fatalf("create client→backend pipe: %v", err)
        }

        return &pipeConn{
                backendWrite: cw,
                clientRead:   cr,
                clientWrite:  bw,
                backendRead:  br,
        }
}

// writeJSON writes a JSON-RPC message from the fake backend to the client.
func (p *pipeConn) writeJSON(t *testing.T, v any) {
        t.Helper()
        data, err := json.Marshal(v)
        if err != nil {
                t.Fatalf("fake backend: marshal: %v", err)
        }
        data = append(data, '\n')
        if _, err := p.backendWrite.Write(data); err != nil {
                t.Fatalf("fake backend: write: %v", err)
        }
}

// readLine reads one newline-delimited JSON line from the client's output.
// Uses a buffered reader and a timeout to avoid hanging tests.
func (p *pipeConn) readLine(t *testing.T) map[string]any {
        t.Helper()

        type readResult struct {
                line string
                err  error
        }
        doneCh := make(chan readResult, 1)

        go func() {
                reader := bufio.NewReader(p.backendRead)
                line, err := reader.ReadBytes('\n')
                if err != nil {
                        doneCh <- readResult{"", err}
                        return
                }
                doneCh <- readResult{string(line), nil}
        }()

        select {
        case res := <-doneCh:
                if res.err != nil {
                        t.Fatalf("fake backend: read: %v", res.err)
                }
                var msg map[string]any
                if err := json.Unmarshal([]byte(res.line), &msg); err != nil {
                        t.Fatalf("fake backend: unmarshal: %v", err)
                }
                return msg
        case <-time.After(3 * time.Second):
                t.Fatal("timed out waiting for client output")
                return nil
        }
}

// close signals the fake backend to close the connection.
func (p *pipeConn) close() {
        p.backendWrite.Close()
        p.clientWrite.Close()
        p.clientRead.Close()
        p.backendRead.Close()
}

// ---------------------------------------------------------------------------
// Test helpers — reduce subscribe-and-wait boilerplate
// ---------------------------------------------------------------------------

// waitForBus subscribes to the bus and waits for a message of type T.
// Returns a pointer to the received message (updated after wait) and a
// wait function. The caller must call wait() before reading the result.
func waitForBus[T any](t *testing.T, b *bus.Bus) (*T, func()) {
        t.Helper()
        result := new(T)
        var wg sync.WaitGroup
        wg.Add(1)
        unsub := b.Subscribe(func(msg any) bool {
                if v, ok := msg.(T); ok {
                        *result = v
                        wg.Done()
                }
                return true
        })
        return result, func() { wg.Wait(); unsub() }
}

// ---------------------------------------------------------------------------
// Test: navigate command → bus.NavigateMsg
// ---------------------------------------------------------------------------

func TestHandlerNavigateCommand(t *testing.T) {
        b := bus.New()
        h := NewHandler(b)

        received, wait := waitForBus[bus.NavigateMsg](t, b)

        raw := map[string]any{
                "jsonrpc": protocol.Version,
                "method":  protocol.MethodNavigate,
                "params": map[string]any{
                        "screen":  string(protocol.ScreenMonitor),
                        "replace": true,
                },
        }

        h.HandleMessage(raw, nil)
        wait()

        if received.Screen != protocol.ScreenMonitor {
                t.Errorf("expected screen %q, got %q", protocol.ScreenMonitor, received.Screen)
        }
        if replace, _ := received.Params["replace"].(bool); !replace {
                t.Error("expected replace=true in params")
        }
}

// ---------------------------------------------------------------------------
// Test: update command → bus.UpdateMsg
// ---------------------------------------------------------------------------

func TestHandlerUpdateCommand(t *testing.T) {
        b := bus.New()
        h := NewHandler(b)

        received, wait := waitForBus[bus.UpdateMsg](t, b)

        raw := map[string]any{
                "jsonrpc": protocol.Version,
                "method":  protocol.MethodUpdate,
                "params": map[string]any{
                        "screen":   string(protocol.ScreenScrape),
                        "progress": float64(0.75),
                },
        }

        h.HandleMessage(raw, nil)
        wait()

        if received.Screen != protocol.ScreenScrape {
                t.Errorf("expected screen %q, got %q", protocol.ScreenScrape, received.Screen)
        }
}

// ---------------------------------------------------------------------------
// Test: notify command → bus.NotifyMsg
// ---------------------------------------------------------------------------

func TestHandlerNotifyCommand(t *testing.T) {
        b := bus.New()
        h := NewHandler(b)

        received, wait := waitForBus[bus.NotifyMsg](t, b)

        raw := map[string]any{
                "jsonrpc": protocol.Version,
                "method":  protocol.MethodNotify,
                "params": map[string]any{
                        "type":     string(protocol.NotifScrapeComplete),
                        "severity": string(protocol.SeverityPositive),
                        "data": map[string]any{
                                "message": "Scrape complete for kuliner niche",
                        },
                },
        }

        h.HandleMessage(raw, nil)
        wait()

        if received.Type != string(protocol.NotifScrapeComplete) {
                t.Errorf("expected type %q, got %q", protocol.NotifScrapeComplete, received.Type)
        }
        if received.Severity != protocol.SeverityPositive {
                t.Errorf("expected severity %q, got %q", protocol.SeverityPositive, received.Severity)
        }
        if msg, _ := received.Data["message"].(string); msg != "Scrape complete for kuliner niche" {
                t.Errorf("unexpected data message: %q", msg)
        }
}

// ---------------------------------------------------------------------------
// Test: validate command → bus.ValidateMsg
// ---------------------------------------------------------------------------

func TestHandlerValidateCommand(t *testing.T) {
        b := bus.New()
        h := NewHandler(b)

        received, wait := waitForBus[bus.ValidateMsg](t, b)

        raw := map[string]any{
                "jsonrpc": protocol.Version,
                "method":  protocol.MethodValidate,
                "params": map[string]any{
                        "errors":   []any{"missing niche.yaml"},
                        "warnings": []any{"high spam_guard threshold"},
                },
        }

        h.HandleMessage(raw, nil)
        wait()

        if len(received.Errors) != 1 || received.Errors[0] != "missing niche.yaml" {
                t.Errorf("unexpected errors: %v", received.Errors)
        }
        if len(received.Warnings) != 1 || received.Warnings[0] != "high spam_guard threshold" {
                t.Errorf("unexpected warnings: %v", received.Warnings)
        }
}

// ---------------------------------------------------------------------------
// Test: unknown method is silently ignored
// ---------------------------------------------------------------------------

func TestHandlerUnknownMethod(t *testing.T) {
        b := bus.New()
        h := NewHandler(b)

        raw := map[string]any{
                "jsonrpc": protocol.Version,
                "method":  "unknown_method",
                "params":  map[string]any{},
        }
        h.HandleMessage(raw, nil)

        if b.HasPending() {
                t.Error("expected no bus messages for unknown method")
        }
}

// ---------------------------------------------------------------------------
// Test: invalid screen in navigate is silently ignored
// ---------------------------------------------------------------------------

func TestHandlerInvalidScreen(t *testing.T) {
        b := bus.New()
        h := NewHandler(b)

        raw := map[string]any{
                "jsonrpc": protocol.Version,
                "method":  protocol.MethodNavigate,
                "params": map[string]any{
                        "screen": "nonexistent_screen",
                },
        }
        h.HandleMessage(raw, nil)

        if b.HasPending() {
                t.Error("expected no bus messages for invalid screen")
        }
}

// ---------------------------------------------------------------------------
// Test: notification variant (no "id" field)
// ---------------------------------------------------------------------------

func TestHandlerBackendNotification(t *testing.T) {
        b := bus.New()
        h := NewHandler(b)

        received, wait := waitForBus[bus.NotifyMsg](t, b)

        // Notification: has "method" but no "id".
        raw := map[string]any{
                "jsonrpc": protocol.Version,
                "method":  protocol.MethodNotify,
                "params": map[string]any{
                        "type":     string(protocol.NotifWAFlag),
                        "severity": string(protocol.SeverityCritical),
                        "data":     map[string]any{},
                },
        }

        h.HandleMessage(raw, nil)
        wait()

        if received.Severity != protocol.SeverityCritical {
                t.Errorf("expected severity %q, got %q", protocol.SeverityCritical, received.Severity)
        }
}

// ---------------------------------------------------------------------------
// Test: response routing via Client (not Handler)
// ---------------------------------------------------------------------------

func TestClientResponseRouting(t *testing.T) {
        conn := newPipeConn(t)
        defer conn.close()

        b := bus.New()
        client := NewClient(b)
        client.Start(conn.clientRead, conn.clientWrite)
        defer client.Stop()

        // Send a request to get a pending response channel.
        event := protocol.RequestEvent{
                Type:   "fetch_leads",
                Screen: protocol.ScreenLeadsDB,
        }
        respCh, err := client.SendRequest(event)
        if err != nil {
                t.Fatalf("SendRequest: %v", err)
        }

        // Read the request from the client output to get the ID.
        msg := conn.readLine(t)
        id := msg["id"]

        // Write a response back through the fake backend.
        resp := map[string]any{
                "jsonrpc": protocol.Version,
                "id":      id,
                "result":  map[string]any{"count": float64(42)},
        }
        conn.writeJSON(t, resp)

        select {
        case r := <-respCh:
                if r.Error != nil {
                        t.Errorf("unexpected error: %v", r.Error)
                }
                if r.ID == 0 {
                        t.Error("expected non-zero response ID")
                }
        case <-time.After(3 * time.Second):
                t.Fatal("timed out waiting for response")
        }
}

// ---------------------------------------------------------------------------
// Test: response with error
// ---------------------------------------------------------------------------

func TestClientResponseError(t *testing.T) {
        conn := newPipeConn(t)
        defer conn.close()

        b := bus.New()
        client := NewClient(b)
        client.Start(conn.clientRead, conn.clientWrite)
        defer client.Stop()

        event := protocol.RequestEvent{
                Type:   "get_stats",
                Screen: protocol.ScreenMonitor,
        }
        respCh, err := client.SendRequest(event)
        if err != nil {
                t.Fatalf("SendRequest: %v", err)
        }

        msg := conn.readLine(t)
        id := msg["id"]

        resp := map[string]any{
                "jsonrpc": protocol.Version,
                "id":      id,
                "error": map[string]any{
                        "code":    float64(protocol.ErrorCodeMethodNotFound),
                        "message": "method not found",
                },
        }
        conn.writeJSON(t, resp)

        select {
        case r := <-respCh:
                if r.Error == nil {
                        t.Fatal("expected error, got nil")
                }
                if r.Error.Code != protocol.ErrorCodeMethodNotFound {
                        t.Errorf("expected error code %d, got %d", protocol.ErrorCodeMethodNotFound, r.Error.Code)
                }
        case <-time.After(3 * time.Second):
                t.Fatal("timed out waiting for error response")
        }
}

// ---------------------------------------------------------------------------
// Test: severity validation — invalid severity falls back to neutral
// ---------------------------------------------------------------------------

func TestHandlerInvalidSeverity(t *testing.T) {
        b := bus.New()
        h := NewHandler(b)

        received, wait := waitForBus[bus.NotifyMsg](t, b)

        raw := map[string]any{
                "jsonrpc": protocol.Version,
                "method":  protocol.MethodNotify,
                "params": map[string]any{
                        "type":     string(protocol.NotifDailyLimit),
                        "severity": "invalid_severity",
                        "data":     map[string]any{},
                },
        }

        h.HandleMessage(raw, nil)
        wait()

        // Invalid severity should fall back to SeverityNeutral.
        if received.Severity != protocol.SeverityNeutral {
                t.Errorf("expected fallback severity %q, got %q", protocol.SeverityNeutral, received.Severity)
        }
}

// ---------------------------------------------------------------------------
// Test: invalid notification type is silently ignored
// ---------------------------------------------------------------------------

func TestHandlerInvalidNotificationType(t *testing.T) {
        b := bus.New()
        h := NewHandler(b)

        raw := map[string]any{
                "jsonrpc": protocol.Version,
                "method":  protocol.MethodNotify,
                "params": map[string]any{
                        "type":     "invalid_notif_type",
                        "severity": string(protocol.SeverityNeutral),
                        "data":     map[string]any{},
                },
        }

        h.HandleMessage(raw, nil)

        if b.HasPending() {
                t.Error("expected no bus messages for invalid notification type")
        }
}

// ---------------------------------------------------------------------------
// Test: client SendKeyPress produces correct JSON-RPC notification
// ---------------------------------------------------------------------------

func TestClientSendKeyPress(t *testing.T) {
        conn := newPipeConn(t)
        defer conn.close()

        b := bus.New()
        client := NewClient(b)
        client.Start(conn.clientRead, conn.clientWrite)
        defer client.Stop()

        event := protocol.KeyPressEvent{
                Key:    "enter",
                Screen: protocol.ScreenBoot,
                State:  protocol.BootFirstTime,
        }
        if err := client.SendKeyPress(event); err != nil {
                t.Fatalf("SendKeyPress: %v", err)
        }

        msg := conn.readLine(t)

        // Verify it's a notification (no "id" field, has "method").
        if _, hasID := msg["id"]; hasID {
                t.Error("expected notification (no id), but id field present")
        }
        if method, _ := msg["method"].(string); method != protocol.MethodKeyPress {
                t.Errorf("expected method %q, got %q", protocol.MethodKeyPress, method)
        }

        // Verify params.
        params, _ := msg["params"].(map[string]any)
        if key, _ := params[protocol.ParamKey].(string); key != "enter" {
                t.Errorf("expected key %q, got %q", "enter", key)
        }
}

// ---------------------------------------------------------------------------
// Test: client SendAction produces correct JSON-RPC notification
// ---------------------------------------------------------------------------

func TestClientSendAction(t *testing.T) {
        conn := newPipeConn(t)
        defer conn.close()

        b := bus.New()
        client := NewClient(b)
        client.Start(conn.clientRead, conn.clientWrite)
        defer client.Stop()

        event := protocol.ActionEvent{
                Action: "confirm",
                Screen: protocol.ScreenLogin,
                Params: map[string]any{"slot": float64(1)},
        }
        if err := client.SendAction(event); err != nil {
                t.Fatalf("SendAction: %v", err)
        }

        msg := conn.readLine(t)

        if method, _ := msg["method"].(string); method != protocol.MethodAction {
                t.Errorf("expected method %q, got %q", protocol.MethodAction, method)
        }

        params, _ := msg["params"].(map[string]any)
        if action, _ := params[protocol.ParamAction].(string); action != "confirm" {
                t.Errorf("expected action %q, got %q", "confirm", action)
        }
}

// ---------------------------------------------------------------------------
// Test: end-to-end — backend sends navigate → client reads → bus receives
// ---------------------------------------------------------------------------

func TestClientEndToEndNavigate(t *testing.T) {
        conn := newPipeConn(t)
        defer conn.close()

        b := bus.New()
        client := NewClient(b)
        client.Start(conn.clientRead, conn.clientWrite)
        defer client.Stop()

        received, wait := waitForBus[bus.NavigateMsg](t, b)

        // Fake backend sends a navigate command.
        navigateNotif := map[string]any{
                "jsonrpc": protocol.Version,
                "method":  protocol.MethodNavigate,
                "params": map[string]any{
                        "screen": string(protocol.ScreenMonitor),
                },
        }
        conn.writeJSON(t, navigateNotif)

        done := make(chan struct{})
        go func() {
                wait()
                close(done)
        }()

        select {
        case <-done:
                if received.Screen != protocol.ScreenMonitor {
                        t.Errorf("expected screen %q, got %q", protocol.ScreenMonitor, received.Screen)
                }
        case <-time.After(3 * time.Second):
                t.Fatal("timed out waiting for navigate bus message")
        }
}

// ---------------------------------------------------------------------------
// Test: client not running returns errors
// ---------------------------------------------------------------------------

func TestClientNotRunning(t *testing.T) {
        b := bus.New()
        client := NewClient(b)

        if err := client.SendKeyPress(protocol.KeyPressEvent{Key: "enter"}); !isNotRunning(err) {
                t.Errorf("expected errNotRunning, got %v", err)
        }

        if err := client.SendAction(protocol.ActionEvent{Action: "select"}); !isNotRunning(err) {
                t.Errorf("expected errNotRunning, got %v", err)
        }

        if _, err := client.SendRequest(protocol.RequestEvent{Type: "fetch"}); !isNotRunning(err) {
                t.Errorf("expected errNotRunning, got %v", err)
        }
}

// isNotRunning checks whether the error is the errNotRunning sentinel.
func isNotRunning(err error) bool {
        return err != nil && err.Error() == errNotRunning.Error()
}

// ---------------------------------------------------------------------------
// Test: FormatError helper
// ---------------------------------------------------------------------------

func TestFormatError(t *testing.T) {
        if msg := FormatError(nil); msg != "" {
                t.Errorf("expected empty string for nil error, got %q", msg)
        }

        rpcErr := &protocol.RPCError{
                Code:    protocol.ErrorCodeInvalidParams,
                Message: "invalid params",
        }
        expected := "RPC error -32602: invalid params"
        if msg := FormatError(rpcErr); msg != expected {
                t.Errorf("expected %q, got %q", expected, msg)
        }
}

// ---------------------------------------------------------------------------
// Test: DecodeResult helper
// ---------------------------------------------------------------------------

func TestDecodeResult(t *testing.T) {
        // Test nil response.
        if err := DecodeResult(nil, nil); err == nil {
                t.Error("expected error for nil response")
        }

        // Test error response.
        errResp := &protocol.Response{
                JSONRPC: protocol.Version,
                ID:      1,
                Error:   &protocol.RPCError{Code: -32000, Message: "custom error"},
        }
        if err := DecodeResult(errResp, nil); err == nil {
                t.Error("expected error for error response")
        }

        // Test empty result.
        emptyResp := &protocol.Response{
                JSONRPC: protocol.Version,
                ID:      2,
        }
        if err := DecodeResult(emptyResp, nil); err == nil {
                t.Error("expected error for empty result")
        }

        // Test valid result.
        type leadResult struct {
                Count int `json:"count"`
        }
        validResp := &protocol.Response{
                JSONRPC: protocol.Version,
                ID:      3,
                Result:  map[string]any{"count": float64(42)},
        }
        var result leadResult
        if err := DecodeResult(validResp, &result); err != nil {
                t.Fatalf("unexpected error: %v", err)
        }
        if result.Count != 42 {
                t.Errorf("expected count 42, got %d", result.Count)
        }
}

// ---------------------------------------------------------------------------
// Test: KeyPressBuilder / ActionBuilder / RequestBuilder
// ---------------------------------------------------------------------------

func TestBuilders(t *testing.T) {
        screenFn := func() protocol.ScreenID { return protocol.ScreenScrape }
        stateFn := func() protocol.StateID { return protocol.ScrapeActive }

        kpb := KeyPressBuilder{Screen: screenFn, State: stateFn}
        evt := kpb.Build("enter")
        if evt.Key != "enter" || evt.Screen != protocol.ScreenScrape || evt.State != protocol.ScrapeActive {
                t.Errorf("unexpected KeyPressEvent: %+v", evt)
        }

        ab := ActionBuilder{Screen: screenFn}
        aEvt := ab.Build("confirm", map[string]any{"id": "123"})
        if aEvt.Action != "confirm" || aEvt.Screen != protocol.ScreenScrape {
                t.Errorf("unexpected ActionEvent: %+v", aEvt)
        }

        rb := RequestBuilder{Screen: screenFn}
        rEvt := rb.Build("fetch_leads", nil)
        if rEvt.Type != "fetch_leads" || rEvt.Screen != protocol.ScreenScrape {
                t.Errorf("unexpected RequestEvent: %+v", rEvt)
        }
}
