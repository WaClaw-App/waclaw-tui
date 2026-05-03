package protocol

// KeyPressEvent represents a raw key press sent from the TUI to the backend.
//
// Not every key press is forwarded — only those that the TUI does not
// handle locally. This allows the backend to drive navigation and state
// transitions based on user input without the TUI hard-coding every key
// binding.
type KeyPressEvent struct {
	// Key is the human-readable key name (e.g. "enter", "esc", "j", "k",
	// "tab", "ctrl+c").
	Key string `json:"key"`

	// Screen identifies the screen that was active when the key was pressed.
	Screen ScreenID `json:"screen"`

	// State is the current state within the screen, if applicable.
	State StateID `json:"state,omitempty"`
}

// ActionEvent represents a semantic user action sent from the TUI to the
// backend.
//
// Unlike KeyPressEvent which carries raw input, ActionEvent carries a
// high-level intent (e.g. "select", "confirm", "toggle", "delete"). The
// backend decides how to respond based on the action name and the current
// screen/state context.
type ActionEvent struct {
	// Action is the semantic action name (e.g. "select", "confirm",
	// "toggle", "delete", "next", "prev").
	Action string `json:"action"`

	// Screen identifies the screen where the action was triggered.
	Screen ScreenID `json:"screen"`

	// Params carries optional key-value parameters for the action.
	// For example, {"lead_id": "123"} for a "select" action on a lead.
	Params map[string]any `json:"params,omitempty"`
}

// RequestEvent represents an explicit data request from the TUI to the
// backend.
//
// The TUI sends a RequestEvent when it needs data that it does not already
// have cached locally — for example, fetching the lead list, retrieving
// stats, or loading a template.
type RequestEvent struct {
	// Type identifies the kind of data being requested (e.g. "fetch_leads",
	// "get_stats", "load_template").
	Type string `json:"type"`

	// Screen identifies the screen that initiated the request.
	Screen ScreenID `json:"screen"`

	// Params carries optional key-value parameters for the request.
	// For example, {"page": 2, "filter": "cold"} for a "fetch_leads"
	// request.
	Params map[string]any `json:"params,omitempty"`
}
