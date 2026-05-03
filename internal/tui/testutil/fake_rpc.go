package testutil

import (
	"time"

	"github.com/WaClaw-App/waclaw/internal/tui/bus"
	"github.com/WaClaw-App/waclaw/pkg/protocol"
)

// FakeRPCClient is a test helper that replays canned RPC responses
// without requiring a real backend connection.
//
// Usage in tests:
//
//	fake := testutil.NewFakeRPCClient()
//	fake.EnqueueNavigate(protocol.ScreenMonitor, nil)
//	fake.EnqueueNotify("test_notif", protocol.SeverityPositive, nil)
type FakeRPCClient struct {
	// Pending holds the messages that will be delivered on the next Poll.
	Pending []any

	// Sent holds all messages that were "sent" by the TUI to the backend.
	Sent []any

	// Bus is the event bus to publish responses to.
	Bus *bus.Bus
}

// NewFakeRPCClient creates a FakeRPCClient with a fresh bus.
func NewFakeRPCClient() *FakeRPCClient {
	return &FakeRPCClient{
		Bus: bus.New(),
	}
}

// EnqueueNavigate queues a navigate message for the next poll.
func (f *FakeRPCClient) EnqueueNavigate(screen protocol.ScreenID, params map[string]any) {
	f.Pending = append(f.Pending, bus.NavigateMsg{
		Screen: screen,
		Params: params,
	})
}

// EnqueueUpdate queues an update message for the next poll.
func (f *FakeRPCClient) EnqueueUpdate(screen protocol.ScreenID, params map[string]any) {
	f.Pending = append(f.Pending, bus.UpdateMsg{
		Screen: screen,
		Params: params,
	})
}

// EnqueueNotify queues a notification message for the next poll.
func (f *FakeRPCClient) EnqueueNotify(notifType string, severity protocol.Severity, data map[string]any) {
	f.Pending = append(f.Pending, bus.NotifyMsg{
		Type:     notifType,
		Severity: severity,
		Data:     data,
	})
}

// EnqueueValidate queues a validation message for the next poll.
func (f *FakeRPCClient) EnqueueValidate(errors, warnings []string) {
	f.Pending = append(f.Pending, bus.ValidateMsg{
		Errors:   errors,
		Warnings: warnings,
	})
}

// Poll delivers all pending messages to the bus and clears the queue.
// Returns the number of messages delivered.
func (f *FakeRPCClient) Poll() int {
	count := len(f.Pending)
	for _, msg := range f.Pending {
		f.Bus.Publish(msg)
	}
	f.Pending = nil
	return count
}

// RecordSent records a TUI→Backend message (key_press, action, request).
func (f *FakeRPCClient) RecordSent(msg any) {
	f.Sent = append(f.Sent, msg)
}

// LastSent returns the most recently sent message, or nil.
func (f *FakeRPCClient) LastSent() any {
	if len(f.Sent) == 0 {
		return nil
	}
	return f.Sent[len(f.Sent)-1]
}

// Reset clears all pending and sent messages.
func (f *FakeRPCClient) Reset() {
	f.Pending = nil
	f.Sent = nil
}

// FakeRPCClock is a controllable clock for deterministic animation tests.
type FakeRPCClock struct {
	Current time.Time
}

// NewFakeRPCClock creates a clock starting at the given time.
func NewFakeRPCClock(t time.Time) *FakeRPCClock {
	return &FakeRPCClock{Current: t}
}

// Now returns the current fake time.
func (c *FakeRPCClock) Now() time.Time { return c.Current }

// Advance moves the clock forward by the given duration.
func (c *FakeRPCClock) Advance(d time.Duration) { c.Current = c.Current.Add(d) }
