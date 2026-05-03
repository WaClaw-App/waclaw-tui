package protocol

// Notification represents a JSON-RPC 2.0 notification object.
//
// A Notification is a Request without an "id" field. The receiver MUST NOT
// reply to a notification — it is a one-way, fire-and-forget message used
// for server-push updates (screen transitions, state changes, toasts, etc.).
//
// See: https://www.jsonrpc.org/specification#notification
type Notification struct {
	// JSONRPC must always be "2.0".
	JSONRPC string `json:"jsonrpc"`

	// Method is the name of the notification type (e.g. "navigate",
	// "update", "notify").
	Method string `json:"method"`

	// Params carries the payload for the notification. May be nil for
	// parameterless notifications.
	Params any `json:"params,omitempty"`
}

// NewNotification creates a fully initialised JSON-RPC 2.0 Notification.
// The JSONRPC field is automatically set to Version ("2.0").
func NewNotification(method string, params any) *Notification {
	return &Notification{
		JSONRPC: Version,
		Method:  method,
		Params:  params,
	}
}
