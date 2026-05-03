package niche

import (
        "fmt"
        "strings"
        "time"

        "github.com/WaClaw-App/waclaw/internal/tui/anim"
        "github.com/WaClaw-App/waclaw/internal/tui/bus"
        "github.com/WaClaw-App/waclaw/internal/tui/i18n"
        "github.com/WaClaw-App/waclaw/internal/tui/style"
        "github.com/WaClaw-App/waclaw/internal/tui/util"
        "github.com/WaClaw-App/waclaw/pkg/protocol"

        "github.com/charmbracelet/lipgloss"
)

// ---------------------------------------------------------------------------
// Shared types and helpers for the niche package
// ---------------------------------------------------------------------------

// NicheItem represents a single selectable niche in the list.
// Fields match what the backend sends via HandleNavigate/HandleUpdate
// params["niches"][i] — the backend is the authoritative source.
type NicheItem struct {
        Name        string
        Description string
        Area        string
        Templates   int
        Selected    bool
        Emoji       string   // e.g. "🍜" — backend sends, TUI renders
        Targets     []string // e.g. ["cafe", "gym"] — backend filter targets
        PreSelected bool     // backend's initial selection state from config
}

// FilterEntry represents a single filter line in the niche filter preview.
type FilterEntry struct {
        Symbol   string // ✗ or ✓ or ⭐ or 📊
        Label    string
        Detail   string
        Inverted bool // true = "skip yang udah punya", false = "lebih potensial"
        Neutral  bool // true = neutral range indicator (⭐, 📊), uses WarningStyle
}

// AreaEntry represents an area with radius and optional kecamatan detail.
type AreaEntry struct {
        City      string
        Radius    string // e.g. "15km"
        Kecamatan []string
        KecCount  int
}

// ConfigError represents a single validation error for a niche config.
type ConfigError struct {
        Line        int    // 0 if not line-specific
        Message     string
        Description string // error description without gutter
        Detail      string // e.g. the problematic line content
        Pointer     string // e.g. "^^^^^" pointing at the error position
}

// Category represents a browsable niche category with emoji and sub-categories.
type Category struct {
        Emoji         string
        Name          string
        SubCategories []string
        SubCount      int
        AreaCount     int
}

// SearchResult represents a search hit from WA Business Directory + GMaps.
type SearchResult struct {
        Name         string
        SubCount     int
        AreaCount    int
        MatchIndices []int
}

// GenerateFileStatus tracks which files have been generated.
type GenerateFileStatus struct {
        Name      string
        Detail    string
        Completed bool
}

// SourceDetail represents a data source indicator.
type SourceDetail struct {
        Name  string
        Count int
        Unit  string // e.g. "bisnis terdaftar", "listing aktif"
}

// FooterEntry is a key-label pair for rendering footer action lines.
type FooterEntry struct {
        Key   string // e.g. "↵", "q", "1"
        Label string // e.g. i18n.T(i18n.KeyNicheGasScrape)
}

// ---------------------------------------------------------------------------
// Layout constants — avoid magic numbers
// ---------------------------------------------------------------------------

const (
        // searchInputMinWidth is the minimum width of the search input field.
        // The actual width scales with terminal width via searchInputWidthFor().
        searchInputMinWidth = 30

        // progressBarWidth is the width of the generation progress bar.
        progressBarWidth = 30

        // filterLabelAlignWidth is the minimum label column width for visual
        // alignment across filter rows. Doc/02 and doc/11 show aligned columns
        // like "✗  punya website      (skip yang udah punya)".
        filterLabelAlignWidth = 18

        // areaColumnWidth is the city name column width in the area list
        // for visual alignment of radius and kecamatan columns.
        areaColumnWidth = 14

        // defaultSearchInputWidth is the fallback width when terminal width
        // is unknown (width == 0).
        defaultSearchInputWidth = 50
)

// ---------------------------------------------------------------------------
// Shared rendering helpers — DRY
// ---------------------------------------------------------------------------

// renderSeparator delegates to the shared style.Separator() which
// uses P3-compliant dim dots. Kept as a local wrapper for backward
// compatibility with existing call sites in this package.
func renderSeparator() string {
        return style.Separator()
}

// renderSectionLabel delegates to the shared style.SectionLabel() which
// uses P3-compliant dim dots. Kept as a local wrapper for backward
// compatibility with existing call sites in this package.
func renderSectionLabel(label string) string {
        return style.SectionLabel(label)
}

// renderFilterEntry renders a single filter line with ✗/✓ symbol and label.
// Shared between select.go (viewFilters) and explorer.go (viewDetail).
func renderFilterEntry(f FilterEntry) string {
        var symbolStyle lipgloss.Style
        switch {
        case f.Inverted:
                // P8: DangerStyle (red) is only for technical problems.
                // Inverted filter ✗ is an exclusion indicator, not a technical error.
                // Use WarningStyle (amber) for visual distinction without P8 violation.
                symbolStyle = style.WarningStyle
        case f.Neutral:
                symbolStyle = style.WarningStyle
        default:
                symbolStyle = style.SuccessStyle
        }
        labelPart := "  " + symbolStyle.Render(f.Symbol) + "  " + f.Label
        if f.Detail != "" {
                // Pad label to filterLabelAlignWidth for visual column alignment.
                if len(f.Label) < filterLabelAlignWidth {
                        labelPart += strings.Repeat(" ", filterLabelAlignWidth-len(f.Label))
                }
                labelPart += style.MutedStyle.Render("(" + f.Detail + ")")
        }
        return labelPart
}

// parseFilterEntries converts raw backend filter data into FilterEntry slices.
// Shared between select.go and explorer.go — single source of truth.
func parseFilterEntries(raw []any) []FilterEntry {
        entries := make([]FilterEntry, 0, len(raw))
        for _, r := range raw {
                data, ok := r.(map[string]any)
                if !ok {
                        continue
                }
                symbol, _ := data["symbol"].(string)
                label, _ := data["label"].(string)
                detail, _ := data["detail"].(string)
                inverted, _ := data["inverted"].(bool)
                neutral, _ := data["neutral"].(bool)
                entries = append(entries, FilterEntry{
                        Symbol:   symbol,
                        Label:    label,
                        Detail:   detail,
                        Inverted: inverted,
                        Neutral:  neutral,
                })
        }
        return entries
}

// parseAreaEntries converts raw backend area data into AreaEntry slices.
// Shared between select.go and explorer.go — single source of truth.
func parseAreaEntries(raw []any) []AreaEntry {
        entries := make([]AreaEntry, 0, len(raw))
        for _, r := range raw {
                data, ok := r.(map[string]any)
                if !ok {
                        continue
                }
                city, _ := data["city"].(string)
                radius, _ := data["radius"].(string)
                kecCount := util.ToInt(data["kecamatan_count"], 0)
                var kecamatan []string
                if rawKec, ok := data["kecamatan"].([]any); ok {
                        kecamatan = util.ToStringSlice(rawKec)
                }
                entries = append(entries, AreaEntry{
                        City:      city,
                        Radius:    radius,
                        Kecamatan: kecamatan,
                        KecCount:  kecCount,
                })
        }
        return entries
}

// renderAreaList renders a list of area entries with city, radius, and kecamatan.
// Uses "··" connector instead of "──" to comply with borderless design.
func renderAreaList(areas []AreaEntry) string {
        var b strings.Builder
        for _, a := range areas {
                areaLine := fmt.Sprintf("  %-*s %s", areaColumnWidth, a.City, a.Radius)
                if len(a.Kecamatan) > 0 {
                        areaLine += fmt.Sprintf("  ··  %s", strings.Join(a.Kecamatan, ", "))
                } else if a.KecCount > 0 {
                        areaLine += fmt.Sprintf("  ··  %d %s", a.KecCount, i18n.T(i18n.KeyNicheKecamatan))
                }
                b.WriteString(style.MutedStyle.Render(areaLine))
                b.WriteString("\n")
        }
        return b.String()
}

// publishAction sends an action event to the backend via the bus.
// Shared helper to avoid duplicating the nil-check pattern.
func publishAction(b *bus.Bus, screen protocol.ScreenID, action string, params map[string]any) {
        if b != nil {
                b.Publish(bus.ActionMsg{
                        Action: action,
                        Screen: screen,
                        Params: params,
                })
        }
}

// renderFooter renders a footer line from a slice of key-label entries.
// DRY helper shared between select.go and explorer.go to eliminate the
// repeated style.ActionStyle.Render(key) + " " + i18n.T(label) pattern.
// Doc/02 and doc/11 use extra spacing after number keys (e.g. "1  buka file")
// to visually separate the key from the label.
func renderFooter(entries []FooterEntry) string {
        parts := make([]string, 0, len(entries))
        for _, e := range entries {
                // Use double-space after single-char keys for doc-parity;
                // keeps single space after multi-char keys like "space" or "1-5".
                gap := " "
                if len(e.Key) == 1 {
                        gap = "  "
                }
                parts = append(parts, style.ActionStyle.Render(e.Key)+gap+e.Label)
        }
        return style.CaptionStyle.Render(strings.Join(parts, "    "))
}

// transitionToList is a DRY helper for the common pattern of returning to
// the niche_list state with stagger reset and backend notification.
// Used by q handlers in select.go.
func transitionToList(m *SelectModel) {
        m.state = protocol.NicheList
        m.staggerStart = time.Now()
        publishAction(m.base.Bus(), protocol.ScreenNicheSelect, protocol.ActionNicheBack, nil)
}

// ---------------------------------------------------------------------------
// Type conversion helpers
// ---------------------------------------------------------------------------

// findCategoryByName looks up a category by name from the loaded categories.
// Returns a Category with at least the Name set. If the category is found in
// the loaded list, SubCategories and other fields are also populated.
func (m *ExplorerModel) findCategoryByName(name string) Category {
        for _, cat := range m.categories {
                if cat.Name == name {
                        return cat
                }
        }
        // Category not in loaded list — return minimal Category with just the name.
        return Category{Name: name}
}

// parseGenFileStatus converts raw backend file status data into GenerateFileStatus.
func parseGenFileStatus(raw []any) []GenerateFileStatus {
        result := make([]GenerateFileStatus, 0, len(raw))
        for _, r := range raw {
                data, ok := r.(map[string]any)
                if !ok {
                        continue
                }
                result = append(result, GenerateFileStatus{
                        Name:      asString(data["name"]),
                        Detail:    asString(data["detail"]),
                        Completed: data["completed"] == true,
                })
        }
        return result
}

// asString is a type-safe string extractor from map[string]any values.
// Shared between select.go and explorer.go — single source of truth.
func asString(v any) string {
        s, _ := v.(string)
        return s
}

// ---------------------------------------------------------------------------
// Type conversion helpers — delegated to util package
// ---------------------------------------------------------------------------

// toInt is a deprecated wrapper around util.ToInt for backward compatibility.
// New code should use util.ToInt directly.
//
// Deprecated: Use util.ToInt(v, 0) instead.
func toInt(v any) int {
        return util.ToInt(v, 0)
}

// toStringSlice is a deprecated wrapper around util.ToStringSlice for backward
// compatibility. New code should use util.ToStringSlice directly.
//
// Deprecated: Use util.ToStringSlice(raw) instead.
func toStringSlice(raw []any) []string {
        return util.ToStringSlice(raw)
}

// slugify is a deprecated wrapper around util.Slugify for backward compatibility.
// New code should use util.Slugify directly.
//
// Deprecated: Use util.Slugify(name) instead.
func slugify(name string) string {
        return util.Slugify(name)
}

// folderSlugOrFallback returns the authoritative folder slug if provided by
// the backend, otherwise falls back to the TUI's display-only util.Slugify().
// The backend owns the actual slug; util.Slugify() is only a visual approximation.
func folderSlugOrFallback(slug, categoryName string) string {
        if slug != "" {
                return slug
        }
        return util.Slugify(categoryName)
}

// searchInputWidthFor returns the search input width based on terminal width.
// Falls back to defaultSearchInputWidth when terminal width is unknown (0).
func searchInputWidthFor(termWidth int) int {
        if termWidth > 0 && termWidth-4 > searchInputMinWidth {
                return termWidth - 4
        }
        return defaultSearchInputWidth
}

// ---------------------------------------------------------------------------
// Stagger visibility helper — shared between Select and Explorer views
// ---------------------------------------------------------------------------

// isStaggerVisible returns true if the item at index i should be visible
// based on the stagger animation timing. Uses anim.NicheStagger (80ms)
// which is the canonical constant in the anim package.
func isStaggerVisible(staggerStart time.Time, index int) bool {
        elapsed := time.Since(staggerStart)
        staggerDelay := time.Duration(index) * anim.NicheStagger
        return elapsed >= staggerDelay
}

// ---------------------------------------------------------------------------
// Checkbox rendering helper — shared between viewList and viewMultiSelected
// ---------------------------------------------------------------------------

// renderCheckbox returns the checkbox character with optional pulse effect.
func renderCheckbox(selected bool, pulseActive bool, pulseStart time.Time, isPulseTarget bool) string {
        checkbox := "☐"
        if selected {
                checkbox = "☑"
        }
        if pulseActive && isPulseTarget {
                elapsed := time.Since(pulseStart)
                if elapsed < anim.SuccessPulse {
                        return lipgloss.NewStyle().Foreground(style.Success).Bold(true).Render(checkbox)
                }
        }
        return checkbox
}

// renderStyledLabel renders a label with focus/selected/muted styling.
func renderStyledLabel(label string, focused, selected bool) string {
        switch {
        case focused:
                return lipgloss.NewStyle().Foreground(style.Accent).Bold(true).Render(label)
        case selected:
                return lipgloss.NewStyle().Foreground(style.Text).Render(label)
        default:
                return lipgloss.NewStyle().Foreground(style.TextMuted).Render(label)
        }
}

// renderStyledDescription renders a description with focus-aware styling.
func renderStyledDescription(desc string, focused bool) string {
        if focused {
                return lipgloss.NewStyle().Foreground(style.TextMuted).Render(desc)
        }
        return lipgloss.NewStyle().Foreground(style.TextDim).Render(desc)
}

// renderSearchInput renders the search input field using the borderless
// design system. Instead of box-drawing characters (┌─┐│└), it uses the
// SearchInput component's View() with a dim background strip for visual
// containment, consistent with vertical borderless principles.
func renderSearchInput(input string, focused bool, placeholder string, termWidth int) string {
        prompt := lipgloss.NewStyle().Foreground(style.TextMuted).Render("> ")
        var value string
        if input == "" {
                value = lipgloss.NewStyle().Foreground(style.TextDim).Render(placeholder)
        } else {
                value = lipgloss.NewStyle().Foreground(style.Text).Render(input)
        }
        cursor := ""
        if focused {
                cursor = lipgloss.NewStyle().Foreground(style.Accent).Render("▎")
        }
        // Calculate responsive width: use terminal width if available,
        // otherwise fall back to defaultSearchInputWidth.
        inputWidth := defaultSearchInputWidth
        if termWidth > 0 && termWidth-4 > searchInputMinWidth {
                inputWidth = termWidth - 4 // leave 2 for side margins
        }
        content := prompt + value + cursor
        padding := inputWidth - len(input) - 3 // 3 for "> " and cursor
        if padding < 0 {
                padding = 0
        }
        content += strings.Repeat(" ", padding)
        return lipgloss.NewStyle().
                Background(style.BgRaised).
                Foreground(style.Text).
                Render(content)
}

// renderErrorGutter renders an error detail line using indentation instead
// of the │ box-drawing character. The vertical borderless design system
// forbids │; we use dim indentation with a red pointer marker instead.
func renderErrorGutter(detail string, pointer string, blinkBold bool) string {
        var b strings.Builder
        // Indent with spaces + dim marker instead of │
        if detail != "" {
                b.WriteString(style.DimStyle.Render("     "))
                b.WriteString(style.DimStyle.Render(detail))
                b.WriteString("\n")
        }
        if pointer != "" {
                b.WriteString(style.DimStyle.Render("     "))
                if blinkBold {
                        b.WriteString(style.DangerStyle.Bold(true).Render(pointer))
                } else {
                        b.WriteString(style.DangerStyle.Render(pointer))
                }
                b.WriteString("\n")
        }
        return b.String()
}

// renderFileTree renders a file tree using indentation instead of
// box-drawing characters (├── └──) to comply with borderless design.
func renderFileTree(files []GenerateFileStatus) string {
        var b strings.Builder
        for i, file := range files {
                // Use "  " indent for all, "  " + connector for hierarchy
                indent := "    "
                if i < len(files)-1 {
                        indent = "  · "
                } else {
                        indent = "    "
                }
                line := indent + file.Name
                if file.Completed {
                        line += "  " + style.SuccessStyle.Render("✓")
                        if file.Detail != "" {
                                line += "  " + style.MutedStyle.Render(file.Detail)
                        }
                }
                b.WriteString(line)
                b.WriteString("\n")
        }
        return b.String()
}

// Ensure i18n key constants used in this package exist.
var (
        _ = i18n.KeyNicheSelect
        _ = i18n.KeyNicheSelected
        _ = i18n.KeyNicheCustom
        _ = i18n.KeyNicheFilters
        _ = i18n.KeyNicheConfigErr
        _ = i18n.KeyNicheExplorerTitle
        _ = i18n.KeyNicheExplorerSubtitle
        _ = i18n.KeyNicheExplorerPopular
        _ = i18n.KeyNicheExplorerSearching
        _ = i18n.KeyNicheExplorerResults
        _ = i18n.KeyNicheExplorerNoResult
        _ = i18n.KeyNicheExplorerSource
        _ = i18n.KeyNicheExplorerGenConfig
        _ = i18n.KeyNicheExplorerGenSuccess
        _ = i18n.KeyNicheExplorerGenProgress
        _ = i18n.KeyNicheExplorerAreaAuto
        _ = i18n.KeyNicheExplorerEditFile
        _ = i18n.KeyNicheExplorerReload
        _ = i18n.KeyNicheExplorerParallel
        _ = i18n.KeyNicheNicheIs
        _ = i18n.KeyNicheMoreNiche
        _ = i18n.KeyNicheCustomDir
        _ = i18n.KeyNicheCustomMin
        _ = i18n.KeyNicheCustomExample
        _ = i18n.KeyNicheCustomReady
        _ = i18n.KeyNicheErrPaused
        _ = i18n.KeyNicheErrOtherOK
        _ = i18n.KeyNicheProblems
        _ = i18n.KeyNicheTargets
        _ = i18n.KeyNicheAreaCount
        _ = i18n.KeyNicheAreaKota
        _ = i18n.KeyNicheFilterDefault
        _ = i18n.KeyNicheTemplateGen
        _ = i18n.KeyNicheJustRight
        _ = i18n.KeyNicheMoreArea
        _ = i18n.KeyNicheCanEdit
        _ = i18n.KeyNicheWorkerParallel
        _ = i18n.KeyNicheScrapeOwn
        _ = i18n.KeyNicheCheckUncheck
        _ = i18n.KeyNicheGasChecked
        _ = i18n.KeyNicheGasNiche
        _ = i18n.KeyNicheChange
        _ = i18n.KeyNicheBack
        _ = i18n.KeyNicheNicheYaml
        _ = i18n.KeyNicheIceBreaker
        _ = i18n.KeyNicheReload
        _ = i18n.KeyNichePickExisting
        _ = i18n.KeyNicheConfigErrLabel
        _ = i18n.KeyNicheGasScrape
        _ = i18n.KeyNicheEditFilter
        _ = i18n.KeyNicheOpenFile
        _ = i18n.KeyNicheShowExample
        _ = i18n.KeyNicheMultiParallel
        _ = i18n.KeyNicheMultiScrapeOwn
        _ = i18n.KeyNicheKecamatan
        _ = i18n.KeyNicheExplorerPickConfig
        _ = i18n.KeyNicheExplorerSearchCat
        _ = i18n.KeyNicheExplorerPick
        _ = i18n.KeyNicheExplorerDetail
        _ = i18n.KeyNicheExplorerDIY
        _ = i18n.KeyNicheExplorerSubCat
        _ = i18n.KeyNicheExplorerSubKategori
        _ = i18n.KeyNicheExplorerSourceLabel
        _ = i18n.KeyNicheExplorerEditFirst
        _ = i18n.KeyNicheExplorerFolder
        _ = i18n.KeyNicheExplorerGasUse
        _ = i18n.KeyNicheExplorerEditConfig
        _ = i18n.KeyNicheExplorerViewTemplate
        _ = i18n.KeyNicheExplorerCancel
        _ = i18n.KeyNicheExplorerAreaLabel
        _ = i18n.KeyNicheExplorerAreaDetect
        _ = i18n.KeyNicheExplorerAreaExplain
        _ = i18n.KeyNicheExplorerAreaSame
        _ = i18n.KeyNicheExplorerAddArea
        _ = i18n.KeyNicheExplorerUseExisting
        _ = i18n.KeyNicheExplorerNoResult
        _ = i18n.KeyNicheNicheDipilih
        _ = i18n.KeyNicheLabel
        _ = i18n.KeyNicheExplorerKategori
        _ = i18n.KeyNicheExplorerFilterDefault
        _ = i18n.KeyNicheLine
        _ = i18n.KeyNicheExplorerGenBarLabel
        _ = i18n.KeyNicheExplorerTitleDetail
        _ = i18n.KeyNicheTemplateCount
        _ = i18n.KeyNicheCustomDirPath
        _ = i18n.KeyNicheCustomExamplePath
        _ = i18n.KeyNicheExplorerFolderPath
)

// Compile-time interface checks.
var (
        _ protocol.ScreenID = protocol.ScreenNicheSelect
        _ protocol.ScreenID = protocol.ScreenNicheExplorer
)
