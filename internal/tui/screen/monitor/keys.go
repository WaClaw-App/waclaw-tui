// Package monitor implements the Monitor (dashboard) and Response screens.
//
// This file defines local key bindings for the monitor package. Screen
// sub-packages cannot import the parent tui package (circular dependency),
// so we replicate the same key bindings here using the same keys/labels
// from keymap.go. The key definitions are identical — this is not a DRY
// violation because the parent package's keys are the canonical source,
// and screen packages must have their own references for the bubbletea
// key matching API.
package monitor

import "github.com/charmbracelet/bubbles/key"

// Local key bindings — mirrors of the global keys from internal/tui/keymap.go.
// Screen packages cannot import the parent tui package, so these are
// necessary local copies. Key values are identical to the global definitions.
var (
        KeyEnter = key.NewBinding(
                key.WithKeys("enter"),
                key.WithHelp("↵", "primary action"),
        )
        KeySkip = key.NewBinding(
                key.WithKeys("s"),
                key.WithHelp("s", "skip"),
        )
        KeyBack = key.NewBinding(
                key.WithKeys("q"),
                key.WithHelp("q", "back/quit"),
        )
        KeyRefresh = key.NewBinding(
                key.WithKeys("r"),
                key.WithHelp("r", "refresh"),
        )
        Key1 = key.NewBinding(key.WithKeys("1"), key.WithHelp("1", "option 1"))
        Key2 = key.NewBinding(key.WithKeys("2"), key.WithHelp("2", "option 2"))
        Key3 = key.NewBinding(key.WithKeys("3"), key.WithHelp("3", "option 3"))
        Key4 = key.NewBinding(key.WithKeys("4"), key.WithHelp("4", "option 4"))
        Key5 = key.NewBinding(key.WithKeys("5"), key.WithHelp("5", "option 5"))
        Key6 = key.NewBinding(key.WithKeys("6"), key.WithHelp("6", "option 6"))
        Key7 = key.NewBinding(key.WithKeys("7"), key.WithHelp("7", "option 7"))
        KeyNerdStats = key.NewBinding(
                key.WithKeys("`"),
                key.WithHelp("`", "nerd stats"),
        )
)
