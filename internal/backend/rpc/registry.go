package rpc

// HandlerFunc is the signature for RPC method handlers.
// It receives the params map and returns a result or an error.
type HandlerFunc func(params map[string]any) (any, error)

// Registry maps RPC method names to handler functions.
// It provides a clean lookup mechanism so the server's message dispatch
// doesn't need a giant switch statement.
type Registry struct {
	handlers map[string]HandlerFunc
}

// NewRegistry creates a new empty method registry.
func NewRegistry() *Registry {
	return &Registry{
		handlers: make(map[string]HandlerFunc),
	}
}

// Register adds a handler for the given method name.
// If a handler is already registered for the method, it is replaced.
func (r *Registry) Register(method string, handler HandlerFunc) {
	r.handlers[method] = handler
}

// Lookup returns the handler for the given method name, or false if not found.
func (r *Registry) Lookup(method string) (HandlerFunc, bool) {
	h, ok := r.handlers[method]
	return h, ok
}

// Methods returns all registered method names.
func (r *Registry) Methods() []string {
	methods := make([]string, 0, len(r.handlers))
	for m := range r.handlers {
		methods = append(methods, m)
	}
	return methods
}
