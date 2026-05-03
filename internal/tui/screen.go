package tui

import (
        tea "github.com/charmbracelet/bubbletea"

        "github.com/WaClaw-App/waclaw/internal/tui/bus"
        "github.com/WaClaw-App/waclaw/pkg/protocol"
)

// Screen is the interface that every TUI screen must implement.
// It extends bubbletea.Model with navigation, update, and bus injection
// methods required by the JSON-RPC communication protocol.
type Screen interface {
        // tea.Model provides Init, Update, and View.
        tea.Model

        // ID returns the unique screen identifier used in protocol messages.
        ID() protocol.ScreenID

        // SetBus injects the event bus reference for publishing/subscribing.
        SetBus(b *bus.Bus)

        // HandleNavigate processes a "navigate" command from the backend.
        // The params map carries screen-specific navigation parameters.
        HandleNavigate(params map[string]any) error

        // HandleUpdate processes an "update" command from the backend.
        // The params map carries screen-specific data updates.
        HandleUpdate(params map[string]any) error

        // Focus is called when this screen becomes the active screen.
        Focus()

        // Blur is called when this screen is no longer the active screen.
        Blur()
}

// StateReporter is an optional interface that screens can implement to
// report their current state ID. The App uses this to populate the
// state field in key_press events sent to the backend.
type StateReporter interface {
        CurrentState() protocol.StateID
}

// KeyConsumer is an optional interface that screens can implement to
// claim priority handling for specific keys that would otherwise be
// consumed by the global key handler. If ConsumesKey returns true,
// the global handler passes the key through to the screen's Update
// method instead of intercepting it.
//
// This is necessary for screens with sub-states where keys like "q"
// should navigate locally (e.g. detail → overview) rather than popping
// the navigation stack, or where "v" should trigger a screen-specific
// action (e.g. validate retry in send_failed) instead of the global
// navigate-to-guardrail behavior.
type KeyConsumer interface {
        ConsumesKey(msg tea.KeyMsg) bool
}

// ScreenBase provides a default implementation of common Screen methods.
// Screens can embed this to avoid boilerplate. It also provides the
// PublishAction convenience method for sending action events to the
// backend, consolidating the publishAction pattern that was duplicated
// across 7+ screen packages.
type ScreenBase struct {
        id    protocol.ScreenID
        bus   *bus.Bus
        state protocol.StateID
}

// NewScreenBase creates a ScreenBase with the given screen ID.
func NewScreenBase(id protocol.ScreenID) ScreenBase {
        return ScreenBase{id: id}
}

// ID returns the screen identifier.
func (s ScreenBase) ID() protocol.ScreenID {
        return s.id
}

// SetBus sets the event bus reference.
func (s *ScreenBase) SetBus(b *bus.Bus) {
        s.bus = b
}

// Bus returns the event bus (nil if not set).
func (s *ScreenBase) Bus() *bus.Bus {
        return s.bus
}

// CurrentState returns the screen's current state ID.
// Screens should call SetState() when their state changes (e.g. in
// HandleNavigate, HandleUpdate, or key handlers). This satisfies the
// StateReporter interface so the App can populate the state field in
// key_press events sent to the backend.
func (s *ScreenBase) CurrentState() protocol.StateID {
        return s.state
}

// SetState updates the screen's current state ID. Screens should call
// this whenever their visual state changes.
func (s *ScreenBase) SetState(state protocol.StateID) {
        s.state = state
}

// PublishAction sends an action event to the backend via the event bus.
// This is the DRY replacement for the publishAction helper functions
// that were duplicated across 7+ screen packages.
//
// The action is published as a bus.NotifyMsg with the screen's ID and
// current state attached, so the backend can route it correctly.
func (s *ScreenBase) PublishAction(action string, params map[string]any) {
        if s.bus == nil {
                return
        }
        if params == nil {
                params = make(map[string]any)
        }
        params["action"] = action
        params["screen"] = string(s.id)
        if s.state != "" {
                params["state"] = string(s.state)
        }
        s.bus.Publish(bus.NotifyMsg{
                Type:     "action",
                Severity: protocol.SeverityNeutral,
                Data:     params,
        })
}

// HandleNavigate provides a default no-op implementation.
func (s ScreenBase) HandleNavigate(params map[string]any) error { return nil }

// HandleUpdate provides a default no-op implementation.
func (s ScreenBase) HandleUpdate(params map[string]any) error { return nil }

// Focus provides a default no-op implementation.
func (s *ScreenBase) Focus() {}

// Blur provides a default no-op implementation.
func (s *ScreenBase) Blur() {}

// Init provides a default no-op implementation.
func (s ScreenBase) Init() tea.Cmd { return nil }
