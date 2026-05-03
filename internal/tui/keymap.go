package tui

import (
        "github.com/charmbracelet/bubbles/key"
)

// Key bindings — memorize once, use forever. Same keys on every screen.
var (
        // KeyUp moves to the previous item.
        KeyUp = key.NewBinding(
                key.WithKeys("up", "k"),
                key.WithHelp("↑/k", "move up"),
        )

        // KeyDown moves to the next item.
        KeyDown = key.NewBinding(
                key.WithKeys("down", "j"),
                key.WithHelp("↓/j", "move down"),
        )

        // KeyEnter is the primary action (the most sensible action).
        KeyEnter = key.NewBinding(
                key.WithKeys("enter"),
                key.WithHelp("↵", "primary action"),
        )

        // KeySkip skips or discards the current item.
        KeySkip = key.NewBinding(
                key.WithKeys("s"),
                key.WithHelp("s", "skip"),
        )

        // KeyBack goes back, exits, or quits (context-dependent).
        KeyBack = key.NewBinding(
                key.WithKeys("q"),
                key.WithHelp("q", "back/quit"),
        )

        // KeyPause pauses the running operation.
        KeyPause = key.NewBinding(
                key.WithKeys("p"),
                key.WithHelp("p", "pause"),
        )

        // KeyRefresh reloads or refreshes.
        KeyRefresh = key.NewBinding(
                key.WithKeys("r"),
                key.WithHelp("r", "refresh"),
        )

        // KeySearch opens search/filter.
        KeySearch = key.NewBinding(
                key.WithKeys("/"),
                key.WithHelp("/", "search"),
        )

        // KeyHelp shows the shortcuts overlay.
        KeyHelp = key.NewBinding(
                key.WithKeys("?"),
                key.WithHelp("?", "shortcuts"),
        )

        // KeyValidate validates configuration.
        KeyValidate = key.NewBinding(
                key.WithKeys("v"),
                key.WithHelp("v", "validate config"),
        )

        // KeyLicense shows license screen.
        KeyLicense = key.NewBinding(
                key.WithKeys("l"),
                key.WithHelp("l", "license"),
        )

        // KeyHistory shows history timeline.
        KeyHistory = key.NewBinding(
                key.WithKeys("h"),
                key.WithHelp("h", "history"),
        )

        // KeyNerdStats toggles the nerd stats overlay.
        KeyNerdStats = key.NewBinding(
                key.WithKeys("`"),
                key.WithHelp("`", "nerd stats"),
        )

        // KeyUpdate checks for updates.
        KeyUpdate = key.NewBinding(
                key.WithKeys("u"),
                key.WithHelp("u", "check update"),
        )

        // KeyCmdPalette opens the command palette.
        KeyCmdPalette = key.NewBinding(
                key.WithKeys("ctrl+k"),
                key.WithHelp("ctrl+k", "command palette"),
        )

        // KeyTab switches between tabs or niche groups.
        KeyTab = key.NewBinding(
                key.WithKeys("tab"),
                key.WithHelp("tab", "switch tab"),
        )

        // KeyNew creates a new item (e.g. new template).
        KeyNew = key.NewBinding(
                key.WithKeys("n"),
                key.WithHelp("n", "new"),
        )

        // KeyEdit edits the selected item (e.g. edit template).
        KeyEdit = key.NewBinding(
                key.WithKeys("e"),
                key.WithHelp("e", "edit"),
        )

        // KeyEscape cancels, closes modal, exits compose, closes palette.
        KeyEscape = key.NewBinding(
                key.WithKeys("esc"),
                key.WithHelp("esc", "cancel/close"),
        )

        // KeyBlock skips and blocks the current item (lead review).
        KeyBlock = key.NewBinding(
                key.WithKeys("x"),
                key.WithHelp("x", "skip & block"),
        )

        // KeyDetail shows detail view for the current item (lead review).
        KeyDetail = key.NewBinding(
                key.WithKeys("d"),
                key.WithHelp("d", "detail"),
        )

        // KeyAutoAll approves or applies action to all items (follow-up).
        KeyAutoAll = key.NewBinding(
                key.WithKeys("a"),
                key.WithHelp("a", "auto-all"),
        )

        // KeySpace toggles checkbox selection (niche select).
        KeySpace = key.NewBinding(
                key.WithKeys(" "),
                key.WithHelp("space", "toggle"),
        )

        // KeyLeft navigates to previous item/date (history).
        KeyLeft = key.NewBinding(
                key.WithKeys("left"),
                key.WithHelp("←", "previous"),
        )

        // KeyRight navigates to next item/date (history).
        KeyRight = key.NewBinding(
                key.WithKeys("right"),
                key.WithHelp("→", "next"),
        )

        // KeyPlus adds a new slot/item (login).
        KeyPlus = key.NewBinding(
                key.WithKeys("+"),
                key.WithHelp("+", "add"),
        )

        // Numeric keys for secondary actions (1-9).
        Key1 = key.NewBinding(key.WithKeys("1"), key.WithHelp("1", "option 1"))
        Key2 = key.NewBinding(key.WithKeys("2"), key.WithHelp("2", "option 2"))
        Key3 = key.NewBinding(key.WithKeys("3"), key.WithHelp("3", "option 3"))
        Key4 = key.NewBinding(key.WithKeys("4"), key.WithHelp("4", "option 4"))
        Key5 = key.NewBinding(key.WithKeys("5"), key.WithHelp("5", "option 5"))
        Key6 = key.NewBinding(key.WithKeys("6"), key.WithHelp("6", "option 6"))
        Key7 = key.NewBinding(key.WithKeys("7"), key.WithHelp("7", "option 7"))
        Key8 = key.NewBinding(key.WithKeys("8"), key.WithHelp("8", "option 8"))
        Key9 = key.NewBinding(key.WithKeys("9"), key.WithHelp("9", "option 9"))

        // NumericKeys returns all numeric key bindings as a slice.
        NumericKeys = []*key.Binding{&Key1, &Key2, &Key3, &Key4, &Key5, &Key6, &Key7, &Key8, &Key9}

        // NavigationKeys returns the core navigation key bindings.
        NavigationKeys = []*key.Binding{&KeyUp, &KeyDown, &KeyEnter, &KeySkip, &KeyBack}

        // GlobalKeys returns all global key bindings (everything except numeric).
        GlobalKeys = []*key.Binding{
                &KeyUp, &KeyDown, &KeyEnter, &KeySkip, &KeyBack,
                &KeyPause, &KeyRefresh, &KeySearch, &KeyHelp,
                &KeyValidate, &KeyLicense, &KeyHistory, &KeyNerdStats,
                &KeyUpdate, &KeyCmdPalette, &KeyEscape, &KeyTab,
                &KeyNew, &KeyEdit, &KeyBlock, &KeyDetail, &KeyAutoAll,
                &KeySpace, &KeyLeft, &KeyRight, &KeyPlus,
        }

        // AllKeys returns every key binding.
        AllKeys = append(GlobalKeys, NumericKeys...)
)
