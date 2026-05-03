package tui

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

// TransitionDirection describes the animation direction for a screen change.
type TransitionDirection int

const (
	// TransitionForward slides in from right (push navigation).
	TransitionForward TransitionDirection = iota

	// TransitionBack slides in from left (pop navigation).
	TransitionBack

	// TransitionFade cross-fades between screens (replace navigation).
	TransitionFade
)

// TransitionState tracks an in-progress screen transition animation.
type TransitionState struct {
	// Direction is the animation direction.
	Direction TransitionDirection

	// Progress is the animation progress from 0.0 to 1.0.
	Progress float64

	// PreviousView is the cached View() output of the screen being left.
	PreviousView string

	// NextView is the cached View() output of the screen being entered.
	NextView string

	// StartTime is when the transition began.
	StartTime time.Time

	// Duration is the total animation duration, derived from animation.go constants.
	Duration time.Duration
}

// NewTransitionState creates a TransitionState for the given direction.
// Uses ScreenTransition (300ms) for forward/back and TabSwitch (200ms) for fade.
func NewTransitionState(dir TransitionDirection, prevView, nextView string) *TransitionState {
	duration := ScreenTransition
	if dir == TransitionFade {
		duration = TabSwitch
	}

	return &TransitionState{
		Direction:    dir,
		Progress:     0.0,
		PreviousView: prevView,
		NextView:     nextView,
		StartTime:    time.Now(),
		Duration:     duration,
	}
}

// IsComplete reports whether the transition animation has finished.
func (t *TransitionState) IsComplete() bool {
	return t.Progress >= 1.0
}

// Tick advances the transition progress based on elapsed time.
// Returns true if the transition is now complete.
func (t *TransitionState) Tick() bool {
	elapsed := time.Since(t.StartTime)
	if t.Duration <= 0 {
		t.Progress = 1.0
		return true
	}
	t.Progress = float64(elapsed) / float64(t.Duration)
	if t.Progress > 1.0 {
		t.Progress = 1.0
	}
	return t.IsComplete()
}

// TransitionCmd returns a tea.Cmd that fires a TransitionCompleteMsg
// after the transition duration elapses.
func TransitionCmd(dir TransitionDirection) tea.Cmd {
	duration := ScreenTransition
	if dir == TransitionFade {
		duration = TabSwitch
	}
	return tea.Tick(duration, func(_ time.Time) tea.Msg {
		return TransitionCompleteMsg{}
	})
}

// TransitionCompleteMsg signals that a screen transition has finished.
type TransitionCompleteMsg struct{}
