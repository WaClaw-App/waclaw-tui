package testutil

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/WaClaw-App/waclaw/internal/tui/bus"
	"github.com/WaClaw-App/waclaw/pkg/protocol"
)

// ScreenHelper is a test harness for TUI screen models.
// It wraps a Screen with helpers for init, key press, and state assertions.
//
// Usage:
//
//	helper := testutil.NewScreenHelper(myScreen)
//	helper.Init()
//	helper.KeyPress("enter")
//	helper.AssertContains(t, "returning")
type ScreenHelper struct {
	// Screen is the screen identifier.
	Screen protocol.ScreenID

	// Model is the bubbletea.Model being tested.
	Model tea.Model

	// Bus is the event bus injected into the screen.
	Bus *bus.Bus
}

// NewScreenHelper creates a ScreenHelper wrapping the given model.
func NewScreenHelper(model tea.Model) *ScreenHelper {
	return &ScreenHelper{
		Model: model,
		Bus:   bus.New(),
	}
}

// Init calls the model's Init method.
func (sh *ScreenHelper) Init() tea.Cmd {
	return sh.Model.Init()
}

// Update sends a message to the model and returns the command.
func (sh *ScreenHelper) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return sh.Model.Update(msg)
}

// KeyPress sends a key message to the model.
func (sh *ScreenHelper) KeyPress(key string) (tea.Model, tea.Cmd) {
	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(key)}
	return sh.Model.Update(msg)
}

// WindowResize sends a window size message.
func (sh *ScreenHelper) WindowResize(width, height int) (tea.Model, tea.Cmd) {
	msg := tea.WindowSizeMsg{Width: width, Height: height}
	return sh.Model.Update(msg)
}

// View returns the current rendered view.
func (sh *ScreenHelper) View() string {
	return sh.Model.View()
}

// AssertContains checks that the rendered view contains the expected substring.
func (sh *ScreenHelper) AssertContains(t *testing.T, substr string) {
	t.Helper()
	view := sh.View()
	if !strings.Contains(view, substr) {
		t.Errorf("expected view to contain %q, got:\n%s", substr, view)
	}
}

// AssertNotContains checks that the rendered view does NOT contain the substring.
func (sh *ScreenHelper) AssertNotContains(t *testing.T, substr string) {
	t.Helper()
	view := sh.View()
	if strings.Contains(view, substr) {
		t.Errorf("expected view NOT to contain %q, but it does", substr)
	}
}
