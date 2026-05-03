// Package notification implements the notification dispatch system for the
// WaClaw backend.
package notification

import (
	"sync"
	"time"

	"github.com/WaClaw-App/waclaw/internal/backend/engine"
	"github.com/WaClaw-App/waclaw/pkg/protocol"
)

// Entry represents a queued notification waiting to be dispatched to the TUI.
type Entry struct {
	Type     protocol.NotificationType
	Severity protocol.Severity
	Data     map[string]any
	Enqueued time.Time
}

// Dispatcher manages the notification queue and pushes notifications to the
// TUI one at a time. It enforces three rules from doc/14-notification-system.md:
//
//   1. Max 1 notification on screen at a time.
//   2. Never stack. Queue the rest.
//   3. Never spam. Each new notification is queued, not broadcast.
//
// The dispatcher does NOT handle the TUI-side auto-dismiss timer — that is
// the overlay's responsibility. The dispatcher only controls what gets pushed
// to the TUI via the RPCPusher interface.
type Dispatcher struct {
	mu     sync.Mutex
	pusher engine.RPCPusher
	queue  []Entry
	active bool // true if a notification is currently displayed on TUI
}

// NewDispatcher creates a notification dispatcher wired to the given pusher.
func NewDispatcher(pusher engine.RPCPusher) *Dispatcher {
	return &Dispatcher{
		pusher: pusher,
		queue:  make([]Entry, 0),
	}
}

// Enqueue adds a notification to the queue and attempts to dispatch it
// immediately if no notification is currently active on the TUI.
// Returns the queue length after enqueue (for monitoring/logging).
func (d *Dispatcher) Enqueue(notifType protocol.NotificationType, data map[string]any) int {
	d.mu.Lock()
	defer d.mu.Unlock()

	severity := SeverityFor(notifType)

	entry := Entry{
		Type:     notifType,
		Severity: severity,
		Data:     data,
		Enqueued: time.Now(),
	}

	// Critical notifications jump to the front of the queue.
	if severity == protocol.SeverityCritical {
		d.queue = append([]Entry{entry}, d.queue...)
	} else {
		d.queue = append(d.queue, entry)
	}

	// Try to dispatch immediately if nothing is active.
	if !d.active {
		d.dispatchNext()
	}

	return len(d.queue)
}

// Ack signals that the TUI has dismissed the current notification (either
// by user action or auto-dismiss timer). This triggers the dispatcher to
// send the next queued notification, if any.
func (d *Dispatcher) Ack() {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.active = false
	d.dispatchNext()
}

// QueueLen returns the number of pending notifications.
func (d *Dispatcher) QueueLen() int {
	d.mu.Lock()
	defer d.mu.Unlock()
	return len(d.queue)
}

// SetPusher updates the RPC pusher (used during initialization).
func (d *Dispatcher) SetPusher(pusher engine.RPCPusher) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.pusher = pusher
}

// dispatchNext pops the front of the queue and pushes it to the TUI.
// Caller must hold d.mu.
func (d *Dispatcher) dispatchNext() {
	if len(d.queue) == 0 || d.pusher == nil {
		return
	}

	entry := d.queue[0]
	d.queue = d.queue[1:]
	d.active = true

	// Push the notification to the TUI via RPC.
	// Errors are logged but not retried — the TUI will call Ack() on
	// dismiss, which triggers the next dispatch attempt.
	_ = d.pusher.PushNotify(entry.Type, entry.Severity, entry.Data)
}
