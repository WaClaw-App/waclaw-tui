package component

import (
	"fmt"
	"strings"

	"github.com/WaClaw-App/waclaw/internal/tui/style"
	"github.com/charmbracelet/lipgloss"
)

// ListItem represents a single selectable item in a list.
type ListItem struct {
	// Label is the primary text for the item.
	Label string

	// Description is an optional secondary line.
	Description string

	// Selected indicates whether the item is currently checked.
	Selected bool

	// Focused indicates whether the item has cursor focus.
	Focused bool

	// Disabled prevents toggling this item.
	Disabled bool
}

// ListSelect renders a multi-select list with ☐/☑ checkbox states.
//
// Visual spec (from doc/02-screens-niche-select.md, doc/07-screens-workers-antiban.md):
//   - Checkboxes: ☐ (unchecked), ☑ (checked)
//   - Space bar toggles selection
//   - Focused item has accent highlight
//   - Disabled items are rendered in text_dim
//   - Toggle animation: brief pulse on checkbox change
type ListSelect struct {
	// Items holds the list entries.
	Items []ListItem

	// Cursor is the index of the currently focused item.
	Cursor int

	// MultiSelect allows multiple items to be checked simultaneously.
	MultiSelect bool
}

// NewListSelect creates a ListSelect with the given items.
func NewListSelect(items []ListItem) ListSelect {
	if len(items) > 0 {
		items[0].Focused = true
	}
	return ListSelect{
		Items:       items,
		Cursor:      0,
		MultiSelect: true,
	}
}

// Up moves the cursor up by one position.
func (ls *ListSelect) Up() {
	if ls.Cursor > 0 {
		ls.Items[ls.Cursor].Focused = false
		ls.Cursor--
		ls.Items[ls.Cursor].Focused = true
	}
}

// Down moves the cursor down by one position.
func (ls *ListSelect) Down() {
	if ls.Cursor < len(ls.Items)-1 {
		ls.Items[ls.Cursor].Focused = false
		ls.Cursor++
		ls.Items[ls.Cursor].Focused = true
	}
}

// Toggle flips the selection state of the current item.
func (ls *ListSelect) Toggle() {
	if ls.Cursor < 0 || ls.Cursor >= len(ls.Items) {
		return
	}
	item := &ls.Items[ls.Cursor]
	if item.Disabled {
		return
	}
	item.Selected = !item.Selected
}

// SelectedLabels returns the labels of all selected items.
func (ls ListSelect) SelectedLabels() []string {
	var labels []string
	for _, item := range ls.Items {
		if item.Selected {
			labels = append(labels, item.Label)
		}
	}
	return labels
}

// SelectedCount returns the number of selected items.
func (ls ListSelect) SelectedCount() int {
	count := 0
	for _, item := range ls.Items {
		if item.Selected {
			count++
		}
	}
	return count
}

// checkbox returns the checkbox character for the item's state.
func checkbox(item ListItem) string {
	if item.Selected {
		return "☑"
	}
	return "☐"
}

// styleItemLabel renders the label with the appropriate style based on
// item state. Extracted as a shared helper — DRY between View and ViewWithNumbers.
func styleItemLabel(item ListItem) string {
	switch {
	case item.Disabled:
		return lipgloss.NewStyle().Foreground(style.TextDim).Render(item.Label)
	case item.Focused:
		return lipgloss.NewStyle().Foreground(style.Accent).Bold(true).Render(item.Label)
	case item.Selected:
		return lipgloss.NewStyle().Foreground(style.Text).Render(item.Label)
	default:
		return lipgloss.NewStyle().Foreground(style.TextMuted).Render(item.Label)
	}
}

// styleItemDescription renders the description with appropriate style.
func styleItemDescription(item ListItem) string {
	if item.Focused {
		return lipgloss.NewStyle().Foreground(style.TextMuted).Render(item.Description)
	}
	return lipgloss.NewStyle().Foreground(style.TextDim).Render(item.Description)
}

// View renders the list with checkboxes and focus highlighting.
func (ls ListSelect) View() string {
	var lines []string

	for _, item := range ls.Items {
		var b strings.Builder

		// Checkbox.
		b.WriteString(checkbox(item))
		b.WriteString(" ")

		// Label.
		b.WriteString(styleItemLabel(item))

		// Optional description.
		if item.Description != "" {
			b.WriteString("  ")
			b.WriteString(styleItemDescription(item))
		}

		lines = append(lines, b.String())
	}

	return strings.Join(lines, "\n")
}

// ViewWithNumbers renders the list with numeric prefixes (1, 2, 3...)
// instead of just checkboxes. Used in some screens like worker_add_niche.
func (ls ListSelect) ViewWithNumbers() string {
	var lines []string

	for i, item := range ls.Items {
		var b strings.Builder

		// Number prefix.
		b.WriteString(fmt.Sprintf("%d  ", i+1))

		// Checkbox.
		b.WriteString(checkbox(item))
		b.WriteString(" ")

		// Label (reuses the shared styling helper).
		b.WriteString(styleItemLabel(item))

		// Description.
		if item.Description != "" {
			b.WriteString("  ")
			b.WriteString(styleItemDescription(item))
		}

		lines = append(lines, b.String())
	}

	return strings.Join(lines, "\n")
}

// String implements fmt.Stringer.
func (ls ListSelect) String() string {
	return fmt.Sprintf("ListSelect{items=%d, cursor=%d, selected=%d}",
		len(ls.Items), ls.Cursor, ls.SelectedCount())
}
