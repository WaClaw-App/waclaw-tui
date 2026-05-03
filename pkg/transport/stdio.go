// Package transport provides I/O layers for WaClaw's RPC communication.
//
// Two transports are available:
//   - StdioTransport: newline-delimited JSON-RPC 2.0 over stdin/stdout,
//     used for TUI ↔ Backend communication.
//   - HTTPTransport: REST API over HTTP, used for the future web frontend.
package transport

import (
	"bufio"
	"encoding/json"
	"io"
	"sync"

	"github.com/WaClaw-App/waclaw/pkg/protocol"
)

// StdioTransport handles JSON-RPC 2.0 communication over stdin/stdout.
// Messages are newline-delimited JSON. Reads are buffered for efficiency,
// and writes are protected by a mutex so concurrent goroutines can safely
// call SendRequest or SendNotification.
type StdioTransport struct {
	reader *bufio.Reader
	writer io.Writer
	mu     sync.Mutex
	nextID int64
}

// NewStdioTransport creates a transport reading from r and writing to w.
// Typically r is os.Stdin and w is os.Stdout.
func NewStdioTransport(r io.Reader, w io.Writer) *StdioTransport {
	return &StdioTransport{
		reader: bufio.NewReader(r),
		writer: w,
		nextID: 1,
	}
}

// SendRequest sends a JSON-RPC request and returns the constructed Request
// message. The request ID is auto-incremented on every call. The encoded
// JSON line is written to the transport's writer followed by a newline.
func (t *StdioTransport) SendRequest(method string, params any) (*protocol.Request, error) {
	t.mu.Lock()
	defer t.mu.Unlock()

	id := t.nextID
	t.nextID++

	req := protocol.NewRequest(id, method, params)

	data, err := Encode(req)
	if err != nil {
		return nil, err
	}

	if _, err := t.writer.Write(data); err != nil {
		return nil, err
	}

	return req, nil
}

// SendNotification sends a JSON-RPC notification (no ID, no response expected).
// The encoded JSON line is written to the transport's writer followed by a newline.
func (t *StdioTransport) SendNotification(method string, params any) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	notif := protocol.NewNotification(method, params)

	data, err := Encode(notif)
	if err != nil {
		return err
	}

	_, err = t.writer.Write(data)
	return err
}

// Receive reads and decodes the next JSON-RPC message from stdin.
// It blocks until a full newline-terminated line is available.
// The raw message is returned as a map, which can be classified as
// Request, Response, or Notification by inspecting the "id" and "method" fields:
//   - Has "method" and "id" → Request
//   - Has "method" and no "id" → Notification
//   - Has "id" and no "method" → Response
func (t *StdioTransport) Receive() (map[string]any, error) {
	line, err := t.reader.ReadBytes('\n')
	if err != nil {
		return nil, err
	}

	var msg map[string]any
	if err := Decode(line, &msg); err != nil {
		return nil, err
	}

	return msg, nil
}

// Encode encodes any value to JSON bytes with a trailing newline.
// This is the wire format for all stdio messages.
func Encode(v any) ([]byte, error) {
	data, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	return append(data, '\n'), nil
}

// Decode decodes JSON bytes into the target. It trims trailing newlines
// before unmarshaling.
func Decode(data []byte, target any) error {
	// Trim trailing newline if present.
	if len(data) > 0 && data[len(data)-1] == '\n' {
		data = data[:len(data)-1]
	}
	return json.Unmarshal(data, target)
}
