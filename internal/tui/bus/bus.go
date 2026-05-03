package bus

import (
        "sync"

        "github.com/WaClaw-App/waclaw/pkg/protocol"
)

// Message types carried by the bus. These are tea.Msg values that
// bubbletea's Update loop will pattern-match on.

// NavigateMsg is pushed when the backend sends a "navigate" command.
type NavigateMsg struct {
        Screen protocol.ScreenID
        Params map[string]any
}

// UpdateMsg is pushed when the backend sends an "update" command.
type UpdateMsg struct {
        Screen protocol.ScreenID
        Params map[string]any
}

// NotifyMsg is pushed when the backend sends a "notify" command.
type NotifyMsg struct {
        Type     string
        Severity protocol.Severity
        Data     map[string]any
}

// ValidateMsg is pushed when the backend sends a "validate" command.
type ValidateMsg struct {
        Errors   []string
        Warnings []string
}

// KeyPressMsg is published when the TUI sends a key_press event.
type KeyPressMsg struct {
        Key    string
        Screen protocol.ScreenID
        State  string
}

// ActionMsg is published when the TUI sends an action event.
type ActionMsg struct {
        Action string
        Screen protocol.ScreenID
        Params map[string]any
}

// ShowConfirmMsg is published when a screen requests the confirmation overlay.
// The app-level handler picks this up and shows the overlay.
type ShowConfirmMsg struct {
        ConfirmType protocol.ConfirmationType
        Data        map[string]any
}

// ConfirmResultMsg is published when the user accepts or dismisses a
// confirmation overlay. Screens can subscribe to this to react to
// confirmation outcomes (e.g., force disconnect → re-validate).
type ConfirmResultMsg struct {
        ConfirmType protocol.ConfirmationType
        Accepted    bool
}

// Subscriber is a function that receives messages from the bus.
// It should return true if it handled the message, false otherwise.
type Subscriber func(msg any) bool

// subscriberEntry wraps a Subscriber with a removed flag so that
// unsubscribing does not require shifting slice indices.
type subscriberEntry struct {
        fn      Subscriber
        removed bool
}

// Bus provides a publish/subscribe mechanism to decouple RPC handlers
// from screen packages. The app.go root model routes bus messages to
// the appropriate screen's Update method.
type Bus struct {
        mu          sync.RWMutex
        subscribers []*subscriberEntry
        pendingMsgs []any
}

// New creates a new Bus instance.
func New() *Bus {
        return &Bus{}
}

// Subscribe registers a subscriber function. Returns an unsubscribe function.
func (b *Bus) Subscribe(fn Subscriber) func() {
        b.mu.Lock()
        defer b.mu.Unlock()

        entry := &subscriberEntry{fn: fn}
        b.subscribers = append(b.subscribers, entry)

        return func() {
                b.mu.Lock()
                defer b.mu.Unlock()
                entry.removed = true
        }
}

// Publish sends a message to all subscribers.
// Messages are queued and delivered synchronously.
func (b *Bus) Publish(msg any) {
        b.mu.RLock()
        subs := make([]*subscriberEntry, len(b.subscribers))
        copy(subs, b.subscribers)
        b.mu.RUnlock()

        // If there are no active subscribers, queue the message as pending.
        hasActive := false
        for _, e := range subs {
                if !e.removed {
                        hasActive = true
                        break
                }
        }
        if !hasActive {
                b.mu.Lock()
                b.pendingMsgs = append(b.pendingMsgs, msg)
                b.mu.Unlock()
                return
        }

        // Deliver to all active subscribers synchronously.
        for _, e := range subs {
                if e.removed {
                        continue
                }
                e.fn(msg) // return value is informational; all subscribers receive the message
        }
}

// Pending returns and clears all pending messages.
// Used by the app's Update loop to batch-process bus messages.
func (b *Bus) Pending() []any {
        b.mu.Lock()
        defer b.mu.Unlock()

        msgs := b.pendingMsgs
        b.pendingMsgs = nil
        return msgs
}

// HasPending returns whether there are pending messages.
func (b *Bus) HasPending() bool {
        b.mu.RLock()
        defer b.mu.RUnlock()

        return len(b.pendingMsgs) > 0
}
