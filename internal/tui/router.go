package tui

import (
	"fmt"

	"github.com/WaClaw-App/waclaw/internal/tui/bus"
	"github.com/WaClaw-App/waclaw/pkg/protocol"
)

// Router manages the screen stack and handles navigation transitions.
//
// It implements a push/pop/replace navigation model. The App delegates all
// screen lifecycle calls through the Router so that navigation, transitions,
// and focus management are centralised in one place.
type Router struct {
	// stack holds the navigation stack as an ordered list of ScreenIDs.
	// The last element is the currently visible screen.
	stack []protocol.ScreenID

	// screens maps each ScreenID to its concrete Screen implementation.
	screens map[protocol.ScreenID]Screen

	// bus is the event bus used to publish transition-related messages.
	bus *bus.Bus

	// transition is the currently active transition animation, or nil.
	transition *TransitionState
}

// NewRouter creates a new Router with the given event bus.
func NewRouter(b *bus.Bus) *Router {
	return &Router{
		stack:   make([]protocol.ScreenID, 0),
		screens: make(map[protocol.ScreenID]Screen),
		bus:     b,
	}
}

// Register adds a screen to the router. Panics if a screen with the same ID
// is already registered.
func (r *Router) Register(s Screen) {
	id := s.ID()
	if _, exists := r.screens[id]; exists {
		panic(fmt.Sprintf("tui: screen %q already registered", id))
	}
	r.screens[id] = s
}

// Current returns the current (top of stack) screen, or nil if empty.
func (r *Router) Current() Screen {
	if len(r.stack) == 0 {
		return nil
	}
	return r.screens[r.stack[len(r.stack)-1]]
}

// CurrentID returns the current screen ID, or empty string if no screen.
func (r *Router) CurrentID() protocol.ScreenID {
	if len(r.stack) == 0 {
		return ""
	}
	return r.stack[len(r.stack)-1]
}

// Push navigates to a new screen, adding it to the stack.
// Triggers a TransitionForward animation. The previous screen receives Blur,
// and the new screen receives Focus.
func (r *Router) Push(id protocol.ScreenID) {
	s, ok := r.screens[id]
	if !ok {
		return
	}

	// Blur the outgoing screen.
	if prev := r.Current(); prev != nil {
		prev.Blur()
	}

	r.stack = append(r.stack, id)
	s.Focus()
	r.transition = NewTransitionState(TransitionForward, "", "")
}

// Pop goes back to the previous screen, removing the current one from the stack.
// Triggers a TransitionBack animation. Returns false if at the bottom of the stack.
func (r *Router) Pop() bool {
	if len(r.stack) <= 1 {
		return false
	}

	// Blur the outgoing screen.
	if prev := r.Current(); prev != nil {
		prev.Blur()
	}

	r.stack = r.stack[:len(r.stack)-1]

	// Focus the now-current screen.
	if next := r.Current(); next != nil {
		next.Focus()
	}

	r.transition = NewTransitionState(TransitionBack, "", "")
	return true
}

// Replace swaps the current screen for a new one without changing stack depth.
// Triggers a TransitionFade animation.
func (r *Router) Replace(id protocol.ScreenID) {
	s, ok := r.screens[id]
	if !ok {
		return
	}

	// Blur the outgoing screen.
	if prev := r.Current(); prev != nil {
		prev.Blur()
	}

	if len(r.stack) == 0 {
		r.stack = append(r.stack, id)
	} else {
		r.stack[len(r.stack)-1] = id
	}

	s.Focus()
	r.transition = NewTransitionState(TransitionFade, "", "")
}

// GoTo navigates to a screen, using Push by default.
// If the screen is already in the stack, pops back to it.
func (r *Router) GoTo(id protocol.ScreenID) {
	// Search the stack from top to bottom for the target screen.
	for i := len(r.stack) - 1; i >= 0; i-- {
		if r.stack[i] == id {
			// Found — pop back to it.
			for len(r.stack)-1 > i {
				if cur := r.Current(); cur != nil {
					cur.Blur()
				}
				r.stack = r.stack[:len(r.stack)-1]
			}
			if next := r.Current(); next != nil {
				next.Focus()
			}
			r.transition = NewTransitionState(TransitionBack, "", "")
			return
		}
	}

	// Not in stack — push it.
	r.Push(id)
}

// Back goes back one screen (alias for Pop).
func (r *Router) Back() bool { return r.Pop() }

// Depth returns the number of screens in the navigation stack.
func (r *Router) Depth() int { return len(r.stack) }

// IsRoot returns whether the navigation stack has only one screen.
func (r *Router) IsRoot() bool { return len(r.stack) <= 1 }

// Screen returns a registered screen by ID, or nil if not found.
func (r *Router) Screen(id protocol.ScreenID) Screen {
	return r.screens[id]
}

// Stack returns a copy of the current navigation stack.
func (r *Router) Stack() []protocol.ScreenID {
	out := make([]protocol.ScreenID, len(r.stack))
	copy(out, r.stack)
	return out
}

// Transition returns the currently active transition, or nil.
func (r *Router) Transition() *TransitionState {
	return r.transition
}

// ClearTransition removes the active transition state.
func (r *Router) ClearTransition() {
	r.transition = nil
}

// SetTransition sets the active transition state (used by App during rendering).
func (r *Router) SetTransition(t *TransitionState) {
	r.transition = t
}
