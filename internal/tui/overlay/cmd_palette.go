package overlay

import (
        "fmt"
        "sort"
        "strings"
        "time"

        "github.com/WaClaw-App/waclaw/internal/tui/anim"
        "github.com/WaClaw-App/waclaw/internal/tui/component"
        "github.com/WaClaw-App/waclaw/internal/tui/i18n"
        "github.com/WaClaw-App/waclaw/internal/tui/style"
        "github.com/WaClaw-App/waclaw/pkg/protocol"
        "github.com/charmbracelet/bubbles/key"
        tea "github.com/charmbracelet/bubbletea"
        "github.com/charmbracelet/lipgloss"
)

// CmdPaletteMode represents the current state of the command palette overlay.
//
// Spec (doc/10-global-overlays.md):
//   - Closed: palette hidden, screen normal
//   - Open: palette visible, search active
//   - Executing: command selected, executing
//   - Empty: no results match search
//   - WithRecent: palette open with recently-used commands
//   - QuickAction: quick-action command selected
type CmdPaletteMode int

const (
        CmdPaletteClosed     CmdPaletteMode = iota
        CmdPaletteOpen                      // search active, results shown
        CmdPaletteExecuting                 // command executing
        CmdPaletteEmpty                     // no results
        CmdPaletteWithRecent                // showing recently used
        CmdPaletteQuickAction               // quick action selected
)

// CommandCategory groups commands by domain for filtering and display.
type CommandCategory string

const (
        CmdCategoryNav      CommandCategory = "navigation"
        CmdCategoryAction   CommandCategory = "action"
        CmdCategoryOverlay  CommandCategory = "overlay"
        CmdCategorySpecial  CommandCategory = "special"
)

// Command represents a single entry in the command palette.
type Command struct {
        // Name is the display name (localized via i18n key).
        Name string

        // Shortcut is the keyboard shortcut (e.g. "s", "1", "Ctrl+K", "—").
        Shortcut string

        // Category groups related commands.
        Category CommandCategory

        // CategoryTag is the short display tag (e.g. "scrape", "workers").
        CategoryTag string

        // IsQuickAction indicates this command executes immediately without navigation.
        IsQuickAction bool

        // Screen is the target screen for navigation commands (zero value = no nav).
        Screen protocol.ScreenID
}

// CmdPalette holds the state for the command palette overlay.
//
// Spec (doc/10-global-overlays.md):
//   - Ctrl+K to open/close
//   - Fuzzy search with 50ms debounce
//   - Recently used (3) at top
//   - Context-aware ordering
//   - Max 1 at a time
//   - Slide down 150ms open, select 50ms, close 100ms
type CmdPalette struct {
        // Mode is the current palette state.
        Mode CmdPaletteMode

        // Width is the available terminal width.
        Width int

        // Search is the fuzzy search input.
        Search component.SearchInput

        // Commands is the full command registry.
        Commands []Command

        // Filtered holds commands matching the current search query.
        Filtered []component.MatchResult

        // SelectedIndex is the cursor position in the filtered list.
        SelectedIndex int

        // RecentKeys stores the last 3 executed command names.
        RecentKeys []string

        // CurrentScreen is the active screen (for context-aware ordering).
        CurrentScreen protocol.ScreenID

        // Anim tracks the slide animation state.
        Anim anim.AnimationState

        // OnExecute is called when a command is selected.
        // The App sets this to wire command execution into its own routing.
        OnExecute func(cmd Command)
}

// NewCmdPalette creates a CmdPalette with the full command registry.
func NewCmdPalette() CmdPalette {
        cp := CmdPalette{
                Mode:          CmdPaletteClosed,
                Search:        component.NewSearchInput(int(anim.CmdPaletteDebounce / time.Millisecond)),
                RecentKeys:    make([]string, 0, 3),
                CurrentScreen: protocol.ScreenBoot,
        }
        cp.Search.Prompt = "> "
        cp.Search.Placeholder = i18n.T("palette.search")
        cp.Commands = buildCommandRegistry()
        return cp
}

// Open opens the command palette.
func (cp *CmdPalette) Open() {
        cp.Mode = CmdPaletteWithRecent
        cp.Search.Value = ""
        cp.Search.CursorPos = 0
        cp.SelectedIndex = 0
        cp.Anim = anim.NewAnimationState(anim.AnimSlide, anim.CmdPaletteSlideDown)
        cp.applyFilter()
}

// Close closes the command palette.
func (cp *CmdPalette) Close() {
        cp.Mode = CmdPaletteClosed
        cp.Search.Value = ""
        cp.Search.CursorPos = 0
        cp.SelectedIndex = 0
        cp.Anim = anim.NewAnimationState(anim.AnimFade, anim.CmdPaletteClose)
}

// IsOpen returns true if the palette is in any non-closed state.
func (cp CmdPalette) IsOpen() bool {
        return cp.Mode != CmdPaletteClosed
}

// HandleKey processes a key event while the palette is open.
// Returns true if the key was consumed by the palette.
func (cp *CmdPalette) HandleKey(msg tea.KeyMsg) bool {
        if !cp.IsOpen() {
                return false
        }

        // Escape or Ctrl+K closes the palette.
        if key.Matches(msg, key.NewBinding(key.WithKeys("esc"))) ||
                key.Matches(msg, key.NewBinding(key.WithKeys("ctrl+k"))) {
                cp.Close()
                return true
        }

        // Enter executes the selected command.
        if key.Matches(msg, key.NewBinding(key.WithKeys("enter"))) {
                cp.executeSelected()
                return true
        }

        // Up/Down navigation.
        if key.Matches(msg, key.NewBinding(key.WithKeys("up", "k"))) {
                if cp.SelectedIndex > 0 {
                        cp.SelectedIndex--
                }
                return true
        }
        if key.Matches(msg, key.NewBinding(key.WithKeys("down", "j"))) {
                if cp.SelectedIndex < len(cp.Filtered)-1 {
                        cp.SelectedIndex++
                }
                return true
        }

        // Backspace.
        if msg.String() == "backspace" {
                cp.Search.Backspace()
                cp.applyFilter()
                return true
        }

        // Printable character → append to search.
        if len(msg.Runes) == 1 && msg.Runes[0] >= 32 {
                cp.Search.AppendChar(msg.String())
                cp.applyFilter()
                return true
        }

        return true // consume all keys while palette is open
}

// SetCurrentScreen updates the context for context-aware ordering.
func (cp *CmdPalette) SetCurrentScreen(screen protocol.ScreenID) {
        cp.CurrentScreen = screen
}

// RecordRecent adds a command name to the recently-used list (max 3).
func (cp *CmdPalette) RecordRecent(name string) {
        // Remove duplicates.
        for i, n := range cp.RecentKeys {
                if n == name {
                        cp.RecentKeys = append(cp.RecentKeys[:i], cp.RecentKeys[i+1:]...)
                        break
                }
        }
        cp.RecentKeys = append(cp.RecentKeys, name)
        if len(cp.RecentKeys) > 3 {
                cp.RecentKeys = cp.RecentKeys[len(cp.RecentKeys)-3:]
        }
}

// Tick advances the animation state.
func (cp *CmdPalette) Tick() {
        cp.Anim.UpdateProgress()
}

// View renders the command palette overlay.
func (cp CmdPalette) View() string {
        if !cp.IsOpen() {
                return ""
        }

        width := cp.Width
        if width < 40 {
                width = 40
        }
        if width > 70 {
                width = 70
        }

        var sections []string

        // Search bar.
        searchBar := lipgloss.NewStyle().
                Width(width).
                Background(style.BgRaised).
                Foreground(style.Text).
                Padding(0, 1).
                Render(cp.Search.View() + strings.Repeat(" ", max(0, width-len(cp.Search.View())-4)) +
                        style.DimStyle.Render("×"))

        sections = append(sections, searchBar)

        // Separator.
        sections = append(sections, lipgloss.NewStyle().
                Foreground(style.TextDim).
                Width(width).
                Render(strings.Repeat("─", width)))

        // Content area.
        if cp.Mode == CmdPaletteEmpty || len(cp.Filtered) == 0 {
                sections = append(sections, cp.renderEmpty(width))
        } else {
                // Recently used section.
                if cp.Search.Value == "" && len(cp.RecentKeys) > 0 {
                        sections = append(sections, cp.renderRecentSection(width))
                }
                // Filtered results.
                sections = append(sections, cp.renderResults(width))
        }

        // Footer.
        footer := lipgloss.NewStyle().
                Foreground(style.TextDim).
                Width(width).
                Render(i18n.T("palette.footer"))
        sections = append(sections, footer)

        // Wrap in a raised-background panel.
        panel := lipgloss.NewStyle().
                Background(style.BgRaised).
                Width(width).
                Padding(0, 1).
                Render(strings.Join(sections, "\n"))

        return panel
}

// renderEmpty shows the empty state with helpful suggestions.
func (cp CmdPalette) renderEmpty(width int) string {
        emptyMsg := i18n.T("palette.empty")
        suggestions := "scrape, pause, leads, config, shield, follow-up, update, license"

        var lines []string
        lines = append(lines, lipgloss.NewStyle().Foreground(style.TextMuted).
                Render(fmt.Sprintf("  %s %q", emptyMsg, cp.Search.Value)))
        lines = append(lines, "")
        lines = append(lines, lipgloss.NewStyle().Foreground(style.TextDim).
                Render(fmt.Sprintf("  %s %s", i18n.T("palette.try_prefix"), suggestions)))

        return strings.Join(lines, "\n")
}

// renderRecentSection shows the recently used commands.
func (cp CmdPalette) renderRecentSection(width int) string {
        var lines []string

        header := lipgloss.NewStyle().Foreground(style.TextDim).Render(
                "── " + i18n.T("palette.recent_header") + " ──────────────────────────────────────────")
        lines = append(lines, header)

        for _, name := range cp.RecentKeys {
                cmd := cp.findCommandByName(name)
                if cmd == nil {
                        continue
                }
                lines = append(lines, cp.renderCommandRow(*cmd, false, width))
        }

        return strings.Join(lines, "\n")
}

// renderResults shows the filtered command results.
func (cp CmdPalette) renderResults(width int) string {
        var lines []string

        if cp.Search.Value == "" {
                header := lipgloss.NewStyle().Foreground(style.TextDim).Render(
                        "── " + i18n.T("palette.all_commands_header") + " ────────────────────────────────────────────────")
                lines = append(lines, header)
        }

        maxVisible := 10
        count := 0
        for i, result := range cp.Filtered {
                if count >= maxVisible {
                        break
                }
                cmd := cp.findCommandByName(result.Item)
                if cmd == nil {
                        continue
                }
                selected := i == cp.SelectedIndex
                lines = append(lines, cp.renderCommandRow(*cmd, selected, width))
                count++
        }

        return strings.Join(lines, "\n")
}

// renderCommandRow renders a single command row.
func (cp CmdPalette) renderCommandRow(cmd Command, selected bool, width int) string {
        // Highlight matched characters if there's a search query.
        name := cmd.Name
        if cp.Search.Value != "" {
                result := component.FuzzyMatch(cp.Search.Value, cmd.Name)
                if result.Matched {
                        name = component.HighlightMatch(cmd.Name, result.MatchedIndices)
                }
        }

        // Selection indicator.
        indicator := "  "
        if selected {
                indicator = lipgloss.NewStyle().Foreground(style.Accent).Render("▸ ")
        }

        // Quick action indicator.
        quickTag := ""
        if cmd.IsQuickAction {
                quickTag = lipgloss.NewStyle().Foreground(style.Warning).Render("  ⚡")
        }

        // Category tag.
        tagStyle := lipgloss.NewStyle().Foreground(style.TextDim)
        tag := tagStyle.Render(fmt.Sprintf("·  %s", cmd.CategoryTag))

        // Shortcut.
        shortcutStyle := lipgloss.NewStyle().Foreground(style.TextMuted)
        shortcut := shortcutStyle.Render(cmd.Shortcut)

        // Compose row: indicator + name + padding + shortcut + tag
        row := fmt.Sprintf("%s%s%s", indicator, name, quickTag)

        // Right-align shortcut and tag.
        rightPart := fmt.Sprintf("%s  %s", shortcut, tag)
        padWidth := width - lipgloss.Width(row) - lipgloss.Width(rightPart) - 2
        if padWidth < 2 {
                padWidth = 2
        }

        // Apply selected highlight to entire row.
        if selected {
                return lipgloss.NewStyle().Background(style.BgActive).Width(width).Render(
                        row + strings.Repeat(" ", padWidth) + rightPart)
        }

        return row + strings.Repeat(" ", padWidth) + rightPart
}

// applyFilter runs fuzzy search against the command registry and updates Filtered.
func (cp *CmdPalette) applyFilter() {
        query := cp.Search.Value

        // Build a searchable list of command names.
        names := make([]string, len(cp.Commands))
        for i, cmd := range cp.Commands {
                names[i] = cmd.Name
        }

        if query == "" {
                // No query → show all, with context-aware ordering.
                cp.Filtered = cp.contextSort(names)
                if len(cp.Filtered) == 0 {
                        cp.Mode = CmdPaletteEmpty
                } else {
                        cp.Mode = CmdPaletteWithRecent
                }
        } else {
                cp.Filtered = component.FilterAndSort(query, names)
                if len(cp.Filtered) == 0 {
                        cp.Mode = CmdPaletteEmpty
                } else {
                        cp.Mode = CmdPaletteOpen
                }
        }

        // Clamp selected index.
        if cp.SelectedIndex >= len(cp.Filtered) {
                cp.SelectedIndex = max(0, len(cp.Filtered)-1)
        }
}

// contextSort returns commands with context-aware ordering.
// Commands relevant to the current screen appear first.
func (cp CmdPalette) contextSort(names []string) []component.MatchResult {
        results := make([]component.MatchResult, len(names))
        for i, name := range names {
                // Base score of 1.0 for all.
                score := 1.0

                // Boost commands relevant to the current screen.
                cmd := cp.findCommandByName(name)
                if cmd != nil {
                        if isContextRelevant(cp.CurrentScreen, cmd.CategoryTag) {
                                score = 10.0
                        }
                }

                results[i] = component.MatchResult{
                        Item:   name,
                        Score:  score,
                        Matched: true,
                }
        }

        // Sort by score descending.
        sort.Slice(results, func(i, j int) bool {
                return results[i].Score > results[j].Score
        })

        return results
}

// executeSelected triggers the currently selected command.
func (cp *CmdPalette) executeSelected() {
        if cp.SelectedIndex >= len(cp.Filtered) {
                return
        }

        selected := cp.Filtered[cp.SelectedIndex]
        cmd := cp.findCommandByName(selected.Item)
        if cmd == nil {
                return
        }

        // Record as recently used.
        cp.RecordRecent(cmd.Name)

        // Mark as executing.
        if cmd.IsQuickAction {
                cp.Mode = CmdPaletteQuickAction
        } else {
                cp.Mode = CmdPaletteExecuting
        }

        // Trigger the callback.
        if cp.OnExecute != nil {
                cp.OnExecute(*cmd)
        }

        // Close after execution.
        cp.Close()
}

// findCommandByName looks up a command by its display name.
func (cp CmdPalette) findCommandByName(name string) *Command {
        for i := range cp.Commands {
                if cp.Commands[i].Name == name {
                        return &cp.Commands[i]
                }
        }
        return nil
}

// isContextRelevant returns true if a command's category tag is relevant
// to the current screen.
func isContextRelevant(screen protocol.ScreenID, tag string) bool {
        relevance := map[protocol.ScreenID][]string{
                protocol.ScreenScrape:      {"scrape"},
                protocol.ScreenMonitor:     {"monitor", "scrape", "database"},
                protocol.ScreenWorkers:     {"workers"},
                protocol.ScreenAntiBan:     {"shield"},
                protocol.ScreenSettings:    {"settings"},
                protocol.ScreenGuardrail:   {"guardrail"},
                protocol.ScreenSend:        {"send", "workers"},
                protocol.ScreenLeadsDB:     {"database"},
                protocol.ScreenFollowUp:    {"followup"},
                protocol.ScreenLicense:     {"license"},
                protocol.ScreenUpdate:      {"update"},
                protocol.ScreenResponse:    {"response"},
                protocol.ScreenHistory:     {"history"},
        }

        tags, ok := relevance[screen]
        if !ok {
                return false
        }
        for _, t := range tags {
                if t == tag {
                        return true
                }
        }
        return false
}

// buildCommandRegistry creates the full command list.
// Command names use i18n keys so display is locale-aware.
func buildCommandRegistry() []Command {
        return []Command{
                // Navigation.
                {Name: i18n.T("palette.cmd_dashboard"), Shortcut: "d", Category: CmdCategoryNav, CategoryTag: "monitor", Screen: protocol.ScreenMonitor},
                {Name: i18n.T("palette.cmd_leads"), Shortcut: "1", Category: CmdCategoryNav, CategoryTag: "database", Screen: protocol.ScreenLeadsDB},
                {Name: i18n.T("palette.cmd_send"), Shortcut: "2", Category: CmdCategoryNav, CategoryTag: "send", Screen: protocol.ScreenSend},
                {Name: i18n.T("palette.cmd_workers"), Shortcut: "3", Category: CmdCategoryNav, CategoryTag: "workers", Screen: protocol.ScreenWorkers},
                {Name: i18n.T("palette.cmd_templates"), Shortcut: "4", Category: CmdCategoryNav, CategoryTag: "template", Screen: protocol.ScreenTemplateMgr},
                {Name: i18n.T("palette.cmd_shield"), Shortcut: "5", Category: CmdCategoryNav, CategoryTag: "shield", Screen: protocol.ScreenAntiBan},
                {Name: i18n.T("palette.cmd_followup"), Shortcut: "6", Category: CmdCategoryNav, CategoryTag: "followup", Screen: protocol.ScreenFollowUp},
                {Name: i18n.T("palette.cmd_settings"), Shortcut: "7", Category: CmdCategoryNav, CategoryTag: "settings", Screen: protocol.ScreenSettings},
                {Name: i18n.T("palette.cmd_history"), Shortcut: "h", Category: CmdCategoryNav, CategoryTag: "history", Screen: protocol.ScreenHistory},
                {Name: i18n.T("palette.cmd_explorer"), Shortcut: "n", Category: CmdCategoryNav, CategoryTag: "explorer", Screen: protocol.ScreenNicheExplorer},
                {Name: i18n.T("palette.cmd_license"), Shortcut: "l", Category: CmdCategoryNav, CategoryTag: "license", Screen: protocol.ScreenLicense},
                {Name: i18n.T("palette.cmd_update"), Shortcut: "u", Category: CmdCategoryNav, CategoryTag: "update", Screen: protocol.ScreenUpdate},

                // Actions.
                {Name: i18n.T("palette.cmd_scrape_all"), Shortcut: "s", Category: CmdCategoryAction, CategoryTag: "scrape", IsQuickAction: true},
                {Name: i18n.T("palette.cmd_pause_all"), Shortcut: "p", Category: CmdCategoryAction, CategoryTag: "workers", IsQuickAction: true},
                {Name: i18n.T("palette.cmd_followup_all"), Shortcut: "a", Category: CmdCategoryAction, CategoryTag: "followup", IsQuickAction: true},
                {Name: i18n.T("palette.cmd_edit_config"), Shortcut: "e", Category: CmdCategoryAction, CategoryTag: "settings"},
                {Name: i18n.T("palette.cmd_reload_config"), Shortcut: "r", Category: CmdCategoryAction, CategoryTag: "settings", IsQuickAction: true},
                {Name: i18n.T("palette.cmd_validate"), Shortcut: "v", Category: CmdCategoryAction, CategoryTag: "guardrail", IsQuickAction: true},
                {Name: i18n.T("palette.cmd_scrape_one"), Shortcut: "—", Category: CmdCategoryAction, CategoryTag: "scrape"},
                {Name: i18n.T("palette.cmd_resume_all"), Shortcut: "—", Category: CmdCategoryAction, CategoryTag: "workers", IsQuickAction: true},
                {Name: i18n.T("palette.cmd_logout_wa"), Shortcut: "—", Category: CmdCategoryAction, CategoryTag: "login"},
                {Name: i18n.T("palette.cmd_force_retry"), Shortcut: "—", Category: CmdCategoryAction, CategoryTag: "scrape", IsQuickAction: true},

                // Overlays.
                {Name: i18n.T("palette.cmd_nerd_stats"), Shortcut: "`", Category: CmdCategoryOverlay, CategoryTag: "overlay"},
                {Name: i18n.T("palette.cmd_cmd_palette"), Shortcut: "Ctrl+K", Category: CmdCategoryOverlay, CategoryTag: "overlay"},
                {Name: i18n.T("palette.cmd_shortcuts"), Shortcut: "?", Category: CmdCategoryOverlay, CategoryTag: "overlay"},

                // Special.
                {Name: i18n.T("palette.cmd_compose"), Shortcut: "—", Category: CmdCategorySpecial, CategoryTag: "compose"},
                {Name: i18n.T("palette.cmd_search_leads"), Shortcut: "/", Category: CmdCategorySpecial, CategoryTag: "database"},
                {Name: i18n.T("palette.cmd_export_csv"), Shortcut: "—", Category: CmdCategorySpecial, CategoryTag: "database"},
                {Name: i18n.T("palette.cmd_mark_converted"), Shortcut: "—", Category: CmdCategorySpecial, CategoryTag: "response"},
                {Name: i18n.T("palette.cmd_block_lead"), Shortcut: "—", Category: CmdCategorySpecial, CategoryTag: "response"},
                {Name: i18n.T("palette.cmd_recontact"), Shortcut: "—", Category: CmdCategorySpecial, CategoryTag: "followup"},
        }
}
