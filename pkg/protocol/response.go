package protocol

// JSON-RPC 2.0 standard error codes.
const (
	// ErrorCodeParseError is returned when invalid JSON was received.
	ErrorCodeParseError = -32700

	// ErrorCodeInvalidRequest is returned when the JSON sent is not a valid Request.
	ErrorCodeInvalidRequest = -32600

	// ErrorCodeMethodNotFound is returned when the method does not exist or is not available.
	ErrorCodeMethodNotFound = -32601

	// ErrorCodeInvalidParams is returned when invalid method parameter(s) are supplied.
	ErrorCodeInvalidParams = -32602

	// ErrorCodeInternalError is returned for internal JSON-RPC errors.
	ErrorCodeInternalError = -32603
)

// Response represents a JSON-RPC 2.0 response object.
//
// A Response is sent as the reply to a Request that carried a non-zero ID.
// Exactly one of Result or Error must be set: Result on success, Error on
// failure.
//
// See: https://www.jsonrpc.org/specification#response_object
type Response struct {
	// JSONRPC must always be "2.0".
	JSONRPC string `json:"jsonrpc"`

	// ID mirrors the ID of the Request this response corresponds to.
	ID int64 `json:"id"`

	// Result holds the return value of the invoked method on success.
	// Must be nil (or omitted) when Error is set.
	Result any `json:"result,omitempty"`

	// Error describes the error when the method invocation fails.
	// Must be nil when Result is set.
	Error *RPCError `json:"error,omitempty"`
}

// RPCError represents a JSON-RPC 2.0 error object.
//
// The Code field uses application-defined negative numbers (the JSON-RPC
// spec reserves -32768…-32000 for protocol errors). Message is a human-
// readable description. Data is optional and may carry arbitrary extra
// context about the failure.
//
// See: https://www.jsonrpc.org/specification#error_object
type RPCError struct {
	// Code is a numeric error identifier. Use negative values for
	// application-specific errors.
	Code int `json:"code"`

	// Message is a short, human-readable description of the error.
	Message string `json:"message"`

	// Data is an optional field that may carry additional context about the
	// error (e.g. validation details, stack trace, etc.).
	Data any `json:"data,omitempty"`
}
