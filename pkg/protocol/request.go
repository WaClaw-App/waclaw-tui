package protocol

// Version is the JSON-RPC protocol version used by all messages.
const Version = "2.0"

// Request represents a JSON-RPC 2.0 request object.
//
// A Request is sent by either side to invoke a remote method. When the "id"
// field is non-zero the sender expects a matching Response; when id is
// absent (zero) the message is treated as a notification (fire-and-forget).
//
// See: https://www.jsonrpc.org/specification#request_object
type Request struct {
	// JSONRPC must always be "2.0".
	JSONRPC string `json:"jsonrpc"`

	// ID is the request identifier. The server MUST reply with the same ID.
	// A zero value means the request is a notification.
	ID int64 `json:"id"`

	// Method is the name of the remote procedure to invoke.
	Method string `json:"method"`

	// Params carries the arguments for the method. May be nil for
	// parameterless calls. When present it is typically a map or a struct
	// that serialises to a JSON object.
	Params any `json:"params,omitempty"`
}

// NewRequest creates a fully initialised JSON-RPC 2.0 Request.
// The JSONRPC field is automatically set to Version ("2.0").
func NewRequest(id int64, method string, params any) *Request {
	return &Request{
		JSONRPC: Version,
		ID:      id,
		Method:  method,
		Params:  params,
	}
}
