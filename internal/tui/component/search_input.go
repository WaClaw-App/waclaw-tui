package component

import (
	"fmt"
	"strings"
	"time"

	"github.com/WaClaw-App/waclaw/internal/tui/style"
	"github.com/charmbracelet/lipgloss"
)

// SearchInput provides a fuzzy search input with debounce.
//
// Visual spec (from doc/10-global-overlays.md, doc/02-screens-niche-select.md):
//   - Fuzzy search: match against item names with fzf-style scoring
//   - Debounce: 50ms (command palette) or 300ms (niche explorer)
//   - Matched characters highlighted with accent color
//   - Prompt: "> " prefix with blinking cursor
//   - Used in: command palette, niche explorer search, leads database filter
type SearchInput struct {
	// Prompt is the input prompt string.
	Prompt string

	// Value is the current search text.
	Value string

	// Placeholder is shown when the input is empty.
	Placeholder string

	// CursorPos is the cursor position within Value.
	CursorPos int

	// Width is the display width of the input field.
	Width int

	// DebounceMs is the debounce delay in milliseconds.
	DebounceMs int

	// LastChange tracks when the value was last modified.
	LastChange time.Time

	// Focused indicates whether the input is active.
	Focused bool
}

// NewSearchInput creates a SearchInput with the given debounce delay.
func NewSearchInput(debounceMs int) SearchInput {
	return SearchInput{
		Prompt:      "> ",
		Placeholder: "type to search...",
		Width:       40,
		DebounceMs:  debounceMs,
		Focused:     true,
	}
}

// IsDebounced returns true if the debounce period has not yet elapsed
// since the last change.
func (si SearchInput) IsDebounced(now time.Time) bool {
	return now.Sub(si.LastChange) < time.Duration(si.DebounceMs)*time.Millisecond
}

// SetValue updates the search text and records the change time.
func (si *SearchInput) SetValue(v string) {
	if v != si.Value {
		si.Value = v
		si.LastChange = time.Now()
		if si.CursorPos > len(v) {
			si.CursorPos = len(v)
		}
	}
}

// AppendChar adds a character at the cursor position.
func (si *SearchInput) AppendChar(ch string) {
	before := si.Value[:si.CursorPos]
	after := si.Value[si.CursorPos:]
	si.Value = before + ch + after
	si.CursorPos++
	si.LastChange = time.Now()
}

// Backspace removes the character before the cursor.
func (si *SearchInput) Backspace() {
	if si.CursorPos <= 0 {
		return
	}
	si.Value = si.Value[:si.CursorPos-1] + si.Value[si.CursorPos:]
	si.CursorPos--
	si.LastChange = time.Now()
}

// Clear empties the search input.
func (si *SearchInput) Clear() {
	si.Value = ""
	si.CursorPos = 0
	si.LastChange = time.Now()
}

// View renders the search input.
func (si SearchInput) View() string {
	var b strings.Builder

	// Prompt.
	b.WriteString(lipgloss.NewStyle().Foreground(style.TextMuted).Render(si.Prompt))

	// Value or placeholder.
	if si.Value == "" {
		b.WriteString(lipgloss.NewStyle().Foreground(style.TextDim).Render(si.Placeholder))
	} else {
		b.WriteString(lipgloss.NewStyle().Foreground(style.Text).Render(si.Value))
	}

	// Cursor.
	if si.Focused {
		b.WriteString(lipgloss.NewStyle().Foreground(style.Accent).Render("▎"))
	}

	return b.String()
}

// MatchResult holds the result of a fuzzy match.
type MatchResult struct {
	// Item is the original item string.
	Item string

	// Score is the match quality (higher = better).
	Score float64

	// Matched indices are the positions of matched characters.
	MatchedIndices []int

	// Matched is whether the item matches the query.
	Matched bool
}

// FuzzyMatch performs a fuzzy match of query against item.
// Returns a MatchResult with score and highlighted positions.
//
// Scoring (fzf-style):
//   - Exact match: highest score
//   - Prefix match: high score
//   - Substring match: medium score
//   - Fuzzy match: lower score
//   - No match: 0
func FuzzyMatch(query, item string) MatchResult {
	if query == "" {
		return MatchResult{Item: item, Score: 1.0, Matched: true}
	}

	qLower := strings.ToLower(query)
	iLower := strings.ToLower(item)

	// Exact match.
	if iLower == qLower {
		return MatchResult{Item: item, Score: 100.0, Matched: true}
	}

	// Prefix match.
	if strings.HasPrefix(iLower, qLower) {
		indices := make([]int, len(query))
		for i := range indices {
			indices[i] = i
		}
		return MatchResult{Item: item, Score: 80.0, Matched: true, MatchedIndices: indices}
	}

	// Substring match.
	if strings.Contains(iLower, qLower) {
		start := strings.Index(iLower, qLower)
		indices := make([]int, len(query))
		for i := range indices {
			indices[i] = start + i
		}
		return MatchResult{Item: item, Score: 60.0, Matched: true, MatchedIndices: indices}
	}

	// Fuzzy match: check if all query chars appear in order.
	indices := fuzzyIndices(qLower, iLower)
	if len(indices) == len(query) {
		// Score based on how compact the match is.
		span := indices[len(indices)-1] - indices[0] + 1
		score := 40.0 * float64(len(query)) / float64(span)
		return MatchResult{Item: item, Score: score, Matched: true, MatchedIndices: indices}
	}

	return MatchResult{Item: item, Score: 0, Matched: false}
}

// fuzzyIndices finds the positions of query characters in the target
// in order (fuzzy matching). Returns matched character indices.
func fuzzyIndices(query, target string) []int {
	var indices []int
	qi := 0

	for ti := 0; ti < len(target) && qi < len(query); ti++ {
		if target[ti] == query[qi] {
			indices = append(indices, ti)
			qi++
		}
	}

	if qi < len(query) {
		return nil // Not all query chars found.
	}
	return indices
}

// HighlightMatch renders an item with matched characters highlighted in accent.
func HighlightMatch(item string, matchedIndices []int) string {
	if len(matchedIndices) == 0 {
		return lipgloss.NewStyle().Foreground(style.Text).Render(item)
	}

	matchSet := make(map[int]bool)
	for _, idx := range matchedIndices {
		matchSet[idx] = true
	}

	var b strings.Builder
	for i, ch := range item {
		if matchSet[i] {
			b.WriteString(lipgloss.NewStyle().Foreground(style.Accent).Bold(true).Render(string(ch)))
		} else {
			b.WriteString(lipgloss.NewStyle().Foreground(style.Text).Render(string(ch)))
		}
	}
	return b.String()
}

// FilterAndSort filters items by fuzzy match and sorts by score (descending).
func FilterAndSort(query string, items []string) []MatchResult {
	results := make([]MatchResult, 0, len(items))
	for _, item := range items {
		result := FuzzyMatch(query, item)
		if result.Matched {
			results = append(results, result)
		}
	}

	// Sort by score descending (simple insertion sort for small lists).
	for i := 1; i < len(results); i++ {
		for j := i; j > 0 && results[j].Score > results[j-1].Score; j-- {
			results[j], results[j-1] = results[j-1], results[j]
		}
	}

	return results
}

// String implements fmt.Stringer.
func (si SearchInput) String() string {
	return fmt.Sprintf("SearchInput{value=%q, debounce=%dms}", si.Value, si.DebounceMs)
}
