// Package engine defines the interface layer that breaks the circular
// dependency between the scenario engine and the RPC server.
//
// The scenario engine needs to push updates to the TUI (via RPC), but the RPC
// server needs to call into the scenario engine to handle incoming events.
// By depending on interfaces rather than concrete types, both packages can
// import engine/ without creating an import cycle.
//
// Convention: every interface method that pushes data toward the TUI returns
// an error so that transport failures are surfaced to the caller rather than
// silently dropped.
package engine

import (
	"github.com/WaClaw-App/waclaw/pkg/protocol"
)

// ScenarioEngine is the core interface for the demo backend's state machine.
//
// It drives the TUI through a sequence of screen transitions and data updates,
// simulating the behaviour of the real closed-source backend. The RPC server
// calls into ScenarioEngine to handle incoming TUI events (key presses,
// actions, requests), and ScenarioEngine calls back through RPCPusher to
// deliver navigate/update/notify/validate commands to the TUI.
type ScenarioEngine interface {
	// HandleKeyPress processes a key press event from the TUI.
	// The engine decides whether to transition screens, update data, or
	// ignore the event based on the current state.
	HandleKeyPress(evt protocol.KeyPressEvent) error

	// HandleAction processes a semantic action event from the TUI.
	HandleAction(evt protocol.ActionEvent) error

	// HandleRequest processes an explicit data request from the TUI.
	HandleRequest(evt protocol.RequestEvent) (any, error)

	// CurrentScreen returns the screen the engine thinks is active.
	CurrentScreen() protocol.ScreenID

	// CurrentState returns the state within the current screen.
	CurrentState() protocol.StateID

	// StateSnapshot returns the full application state for the REST API.
	StateSnapshot() map[string]any
}

// RPCPusher is the interface used by the scenario engine to push commands
// toward the TUI. The concrete implementation lives in the rpc package.
//
// By expressing the dependency as an interface, the scenario package avoids
// importing rpc/ directly, which would create a cycle:
//
//	scenario → rpc → scenario  (cycle)
//	scenario → engine (interface) ← rpc  (no cycle)
type RPCPusher interface {
	// PushNavigate sends a navigate command to the TUI.
	PushNavigate(screen protocol.ScreenID, state protocol.StateID, params map[string]any) error

	// PushUpdate sends an incremental data update to the current TUI screen.
	PushUpdate(screen protocol.ScreenID, params map[string]any) error

	// PushNotify sends a notification to the TUI.
	PushNotify(notifType protocol.NotificationType, severity protocol.Severity, data map[string]any) error

	// PushValidate sends validation results to the TUI.
	PushValidate(errors []string, warnings []string) error
}

// SlotPauser is the interface that lets the anti-ban system pause WhatsApp
// slots without importing the sender package. This breaks the
// antiban → sender circular dependency.
//
// In the demo backend, this is a no-op because there are no real WhatsApp
// connections. The interface is defined here so that the scenario engine can
// accept it as a dependency and the real backend can wire in the concrete
// implementation.
type SlotPauser interface {
	// PauseSlot pauses the WhatsApp slot identified by slotID.
	PauseSlot(slotID string) error

	// ResumeSlot resumes the WhatsApp slot identified by slotID.
	ResumeSlot(slotID string) error
}

// LeadRepo is the interface for lead persistence. The concrete implementation
// lives in the database package. Defining it here breaks the dependency cycle:
//
//	scenario → database → scenario  (cycle)
//	scenario → engine (interface) ← database  (no cycle)
type LeadRepo interface {
	// StoreLead persists a lead and returns its assigned ID.
	StoreLead(lead map[string]any) (string, error)

	// GetLead retrieves a lead by ID.
	GetLead(id string) (map[string]any, error)

	// ListLeads returns leads matching the given filter params.
	ListLeads(filter map[string]any) ([]map[string]any, error)

	// UpdateLead updates a lead's data.
	UpdateLead(id string, data map[string]any) error
}
