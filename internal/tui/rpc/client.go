package rpc

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"sync"
	"sync/atomic"

	"github.com/WaClaw-App/waclaw/internal/tui/bus"
	"github.com/WaClaw-App/waclaw/pkg/protocol"
	"github.com/WaClaw-App/waclaw/pkg/transport"
	tea "github.com/charmbracelet/bubbletea"
)

// ---------------------------------------------------------------------------
// Client — JSON-RPC 2.0 client that connects the TUI to the backend binary
// over stdio. It owns the transport layer, the message handler, and the
// background read loop.
//
// Lifecycle:
//
//	NewClient(bus) → Start() → [running] → Stop()
//
// DRY convention: all method names use protocol.Method* constants.
// All request construction uses protocol.NewRequest() /
// protocol.NewNotification() constructors — never hand-built structs.
// ---------------------------------------------------------------------------

// errNotRunning is returned when an operation is attempted on a stopped client.
// Callers can use errors.Is(err, errNotRunning) for programmatic matching.
var errNotRunning = errors.New("rpc: client not running")

// RPCMsg is a bubbletea message that carries a bus message published by the
// RPC handler. The App's Update loop processes these to drive screen updates
// without the handler needing to import bubbletea directly.
type RPCMsg struct {
	Msg any // one of bus.NavigateMsg, bus.UpdateMsg, bus.NotifyMsg, bus.ValidateMsg
}

// RPCClosedMsg is published when the RPC connection to the backend terminates.
type RPCClosedMsg struct {
	Err error // nil for graceful shutdown, non-nil for unexpected close
}

// Client manages the JSON-RPC 2.0 connection between the TUI and the backend
// binary over stdio. It reads incoming messages in a background goroutine,
// translates them via the Handler, and publishes them to the bus.
//
// Thread safety: SendKeyPress, SendAction, and SendRequest are safe to call
// from any goroutine. The write path is protected by the transport's internal
// mutex. The pending response map is protected by pendingMu for both reads
// (in readLoop) and writes (in SendRequest/Stop).
type Client struct {
	transport *transport.StdioTransport
	handler   *Handler
	bus       *bus.Bus

	// pendingResp maps request IDs to channels that receive the corresponding
	// response. Populated by SendRequest, consumed by routeResponse in readLoop.
	// Protected by pendingMu for thread-safe concurrent access.
	pendingResp map[int64]chan *protocol.Response
	pendingMu   sync.Mutex

	// done is closed when the read loop should stop.
	done chan struct{}

	// running tracks whether the client has been started.
	running atomic.Bool

	// subscriberUnsub is called to remove the bus subscriber that converts
	// bus messages into tea.Cmd messages.
	subscriberUnsub func()
}

// NewClient creates a new RPC client that will publish translated messages
// to the given bus. The client does NOT start reading until Start() is called.
func NewClient(b *bus.Bus) *Client {
	return &Client{
		handler:     NewHandler(b),
		bus:         b,
		pendingResp: make(map[int64]chan *protocol.Response),
		done:        make(chan struct{}),
	}
}

// Start begins the background read loop. The client reads messages from the
// transport's reader in a dedicated goroutine and dispatches them to the
// handler for translation and bus publishing.
//
// Returns a tea.Cmd that fires an RPCClosedMsg when the read loop exits.
// This allows the App to react to backend disconnections within its Update
// loop.
func (c *Client) Start(r io.Reader, w io.Writer) tea.Cmd {
	if !c.running.CompareAndSwap(false, true) {
		return nil // already started
	}

	c.transport = transport.NewStdioTransport(r, w)

	// Subscribe to bus messages to bridge synchronous bus.Publish() with
	// the asynchronous bubbletea Update loop. The App's processBusMessages()
	// handles bus messages directly via Pending(); this subscriber ensures
	// messages also flow when the bus is used in a pure bubbletea context.
	c.subscriberUnsub = c.bus.Subscribe(func(msg any) bool {
		return true
	})

	// Start the background read loop.
	go c.readLoop()

	// Return a tea.Cmd that blocks until the read loop exits.
	return func() tea.Msg {
		<-c.done
		return RPCClosedMsg{}
	}
}

// Stop signals the read loop to stop and waits for it to finish.
// It is safe to call Stop() multiple times.
func (c *Client) Stop() {
	if !c.running.CompareAndSwap(true, false) {
		return // already stopped
	}

	// Signal the read loop to stop.
	close(c.done)

	// Unsubscribe from the bus.
	if c.subscriberUnsub != nil {
		c.subscriberUnsub()
	}

	// Clean up any pending response channels.
	c.pendingMu.Lock()
	for id, ch := range c.pendingResp {
		close(ch)
		delete(c.pendingResp, id)
	}
	c.pendingMu.Unlock()
}

// ---------------------------------------------------------------------------
// TUI → Backend: send methods
// ---------------------------------------------------------------------------

// SendKeyPress sends a key_press event to the backend.
// This is a notification (fire-and-forget) — no response is expected.
func (c *Client) SendKeyPress(event protocol.KeyPressEvent) error {
	if !c.running.Load() {
		return errNotRunning
	}
	return c.transport.SendNotification(protocol.MethodKeyPress, event)
}

// SendAction sends an action event to the backend.
// This is a notification (fire-and-forget) — no response is expected.
func (c *Client) SendAction(event protocol.ActionEvent) error {
	if !c.running.Load() {
		return errNotRunning
	}
	return c.transport.SendNotification(protocol.MethodAction, event)
}

// SendRequest sends a data request to the backend and returns a channel
// that will receive exactly one response. The caller can select on the
// channel to wait for the result with a timeout if desired.
//
// If the client is not running, returns an error immediately.
// The channel is buffered(1) so the response can be written without blocking.
func (c *Client) SendRequest(event protocol.RequestEvent) (<-chan *protocol.Response, error) {
	if !c.running.Load() {
		return nil, errNotRunning
	}

	req, err := c.transport.SendRequest(protocol.MethodRequest, event)
	if err != nil {
		return nil, fmt.Errorf("rpc: send request: %w", err)
	}

	// Register a pending response channel.
	ch := make(chan *protocol.Response, 1)
	c.pendingMu.Lock()
	c.pendingResp[req.ID] = ch
	c.pendingMu.Unlock()

	return ch, nil
}

// ---------------------------------------------------------------------------
// Background read loop
// ---------------------------------------------------------------------------

// readLoop continuously reads messages from the transport and dispatches them
// to the handler. It exits when the done channel is closed or when the
// transport encounters an EOF.
func (c *Client) readLoop() {
	for {
		select {
		case <-c.done:
			return
		default:
		}

		raw, err := c.transport.Receive()
		if err != nil {
			if err == io.EOF {
				// Backend closed the connection gracefully.
				// Use CompareAndSwap to avoid double-closing the done
				// channel if Stop() was called concurrently.
				if c.running.CompareAndSwap(true, false) {
					close(c.done)
				}
				return
			}
			// Transient read error — continue reading.
			// A production system would use structured logging with backoff.
			continue
		}

		if raw == nil {
			continue
		}

		// Route response messages to waiting callers (under lock for
		// thread safety), then dispatch all other message types to
		// the handler for bus translation.
		c.dispatchMessage(raw)
	}
}

// dispatchMessage routes response messages to pending callers and delegates
// all other message types to the handler. The pending response map is accessed
// under pendingMu to prevent data races with concurrent SendRequest/Stop calls.
func (c *Client) dispatchMessage(raw map[string]any) {
	// Check if this is a response to a pending request (has "id", no "method").
	_, hasID := raw["id"]
	_, hasMethod := raw["method"]

	if hasID && !hasMethod {
		id, ok := toInt64(raw["id"])
		if ok {
			resp := parseResponse(raw)
			c.routeResponse(id, resp)
		}
		// Response handled — don't also send to handler.
		return
	}

	// Not a response — delegate to the handler for bus translation.
	c.handler.HandleMessage(raw, nil)
}

// routeResponse delivers a response to the channel waiting for the given
// request ID. The entry is removed from the pending map to prevent leaks.
// Must be called from the readLoop goroutine.
func (c *Client) routeResponse(id int64, resp *protocol.Response) {
	c.pendingMu.Lock()
	ch, exists := c.pendingResp[id]
	if exists {
		delete(c.pendingResp, id)
	}
	c.pendingMu.Unlock()

	if !exists {
		return
	}

	select {
	case ch <- resp:
	default:
		// Channel full or closed — drop the response.
	}
}

// ---------------------------------------------------------------------------
// Convenience: current screen/state for RPC events
// ---------------------------------------------------------------------------

// CurrentScreenFunc is a function that returns the current screen ID.
// Set by the App so the RPC client can include the current screen context
// in key_press and action events without importing the router.
type CurrentScreenFunc func() protocol.ScreenID

// CurrentStateFunc is a function that returns the current state ID.
type CurrentStateFunc func() protocol.StateID

// KeyPressBuilder helps construct KeyPressEvent values with the current
// screen context filled in automatically. This avoids every caller needing
// to know the current screen/state.
type KeyPressBuilder struct {
	Screen CurrentScreenFunc
	State  CurrentStateFunc
}

// Build creates a KeyPressEvent for the given key string.
func (b KeyPressBuilder) Build(key string) protocol.KeyPressEvent {
	evt := protocol.KeyPressEvent{Key: key}
	if b.Screen != nil {
		evt.Screen = b.Screen()
	}
	if b.State != nil {
		evt.State = b.State()
	}
	return evt
}

// ActionBuilder helps construct ActionEvent values with the current screen
// context filled in automatically.
type ActionBuilder struct {
	Screen CurrentScreenFunc
}

// Build creates an ActionEvent for the given action name and optional params.
func (a ActionBuilder) Build(action string, params map[string]any) protocol.ActionEvent {
	evt := protocol.ActionEvent{Action: action, Params: params}
	if a.Screen != nil {
		evt.Screen = a.Screen()
	}
	return evt
}

// RequestBuilder helps construct RequestEvent values with the current screen
// context filled in automatically.
type RequestBuilder struct {
	Screen CurrentScreenFunc
}

// Build creates a RequestEvent for the given request type and optional params.
func (r RequestBuilder) Build(reqType string, params map[string]any) protocol.RequestEvent {
	evt := protocol.RequestEvent{Type: reqType, Params: params}
	if r.Screen != nil {
		evt.Screen = r.Screen()
	}
	return evt
}

// ---------------------------------------------------------------------------
// Helper: decode a raw message into a typed protocol value.
// Useful for tests and for response body parsing.
// ---------------------------------------------------------------------------

// DecodeResult decodes a protocol.Response's Result field into the target.
// Returns an error if the result is nil or cannot be unmarshaled.
func DecodeResult(resp *protocol.Response, target any) error {
	if resp == nil {
		return fmt.Errorf("rpc: nil response")
	}
	if resp.Error != nil {
		return fmt.Errorf("rpc: %s", FormatError(resp.Error))
	}
	if resp.Result == nil {
		return fmt.Errorf("rpc: empty result")
	}

	// Re-marshal and unmarshal to convert map[string]any → typed struct.
	data, err := json.Marshal(resp.Result)
	if err != nil {
		return fmt.Errorf("rpc: marshal result: %w", err)
	}
	if err := json.Unmarshal(data, target); err != nil {
		return fmt.Errorf("rpc: unmarshal result: %w", err)
	}
	return nil
}
