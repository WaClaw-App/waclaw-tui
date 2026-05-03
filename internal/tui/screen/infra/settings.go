package infra

import (
        "fmt"
        "strings"

        "github.com/WaClaw-App/waclaw/internal/tui"
        "github.com/WaClaw-App/waclaw/internal/tui/bus"
        "github.com/WaClaw-App/waclaw/internal/tui/i18n"
        "github.com/WaClaw-App/waclaw/internal/tui/style"
        "github.com/WaClaw-App/waclaw/pkg/protocol"
        tea "github.com/charmbracelet/bubbletea"
        "github.com/charmbracelet/lipgloss"
)

// ---------------------------------------------------------------------------
// Data types for the Settings screen
// ---------------------------------------------------------------------------

// ChangeEntry represents a single config change after reload.
type ChangeEntry struct {
        Field    string
        OldValue string
        NewValue string
}

// ConfigError represents a config validation error.
type ConfigError struct {
        Line    int
        Message string
        Pointer string
        Context []string
}

// FIX 16: ConfigItem removed — unified with kvItem from helpers.go

// ---------------------------------------------------------------------------
// Settings screen model
// ---------------------------------------------------------------------------

// Settings implements tui.Screen for Screen 13: Settings Config Reference.
// It is NOT a settings editor — it displays a reference card of the current
// configuration, allows opening the config in an external editor, and handles
// reload with success or error feedback.
type Settings struct {
        tui.ScreenBase
        state       protocol.StateID
        changes     []ChangeEntry
        errors      []ConfigError
        width       int
        height      int
        focused     bool
        configItems []kvItem
}

// NewSettings creates a Settings screen with empty defaults.
// All config values are populated by the backend via HandleNavigate.
// FIX 1: configMap replaced with []kvItem using i18n key constants.
// FIX 2: No hardcoded demo values — backend is the single source of truth.
func NewSettings() *Settings {
        return &Settings{
                ScreenBase: tui.NewScreenBase(protocol.ScreenSettings),
                state:      protocol.SettingsOverview,
                configItems: []kvItem{
                        {i18n.T(i18n.KeySettingsActiveNiches), ""},
                        {i18n.T(i18n.KeySettingsWASlots), ""},
                        {i18n.T(i18n.KeySettingsWorkerPool), ""},
                        {i18n.T(i18n.KeySettingsArea), ""},
                        {i18n.T(i18n.KeySettingsWorkHours), ""},
                        {i18n.T(i18n.KeySettingsRateLimit), ""},
                        {i18n.T(i18n.KeySettingsRotatorMode), ""},
                        {i18n.T(i18n.KeySettingsAutopilot), ""},
                },
                errors: nil,
        }
}

// ---------------------------------------------------------------------------
// tea.Model interface
// ---------------------------------------------------------------------------

func (s *Settings) Init() tea.Cmd { return nil }

func (s *Settings) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
        switch m := msg.(type) {
        case tea.WindowSizeMsg:
                s.width = m.Width
                s.height = m.Height
                return s, nil
        }

        if key, ok := msg.(tea.KeyMsg); ok {
                switch key.String() {
                case "q":
                        return s, nil
                case "e":
                        s.state = protocol.SettingsEdit
                case "r":
                        // Trigger reload — publish action to backend and let it respond via HandleUpdate.
                        if s.Bus() != nil {
                                s.Bus().Publish(bus.ActionMsg{Action: string(protocol.ActionReloadSettings), Screen: s.ID()})
                        }
                case "enter":
                        if s.state == protocol.SettingsReload {
                                s.state = protocol.SettingsOverview
                        } else if s.state == protocol.SettingsEdit {
                                // Open editor (notify backend)
                                if s.Bus() != nil {
                                        s.Bus().Publish(bus.ActionMsg{Action: protocol.ActionOpenEditor, Screen: s.ID()})
                                }
                                s.state = protocol.SettingsEdit
                        }
                case "1":
                        if s.state == protocol.SettingsReloadError {
                                // Open file again
                                if s.Bus() != nil {
                                        s.Bus().Publish(bus.ActionMsg{Action: protocol.ActionOpenEditor, Screen: s.ID()})
                                }
                        }
                case "2":
                        if s.state == protocol.SettingsReloadError {
                                // Revert to backup
                                if s.Bus() != nil {
                                        s.Bus().Publish(bus.ActionMsg{Action: protocol.ActionRevertBackup, Screen: s.ID()})
                                }
                        }
                }
        }
        return s, nil
}

func (s *Settings) View() string {
        switch s.state {
        case protocol.SettingsOverview:
                return s.viewOverview()
        case protocol.SettingsEdit:
                return s.viewEdit()
        case protocol.SettingsReload:
                return s.viewReload()
        case protocol.SettingsReloadError:
                return s.viewReloadError()
        default:
                return s.viewOverview()
        }
}

// ---------------------------------------------------------------------------
// Screen interface
// ---------------------------------------------------------------------------

func (s *Settings) HandleNavigate(params map[string]any) error {
        applyNavigateState(&s.state, params)

        // Populate config items from backend data — replaces hardcoded demo values.
        // The backend sends string values for each config key.
        if raw, ok := params[protocol.ParamActiveNiches]; ok {
                if v, ok := raw.(string); ok {
                        s.setConfigItem(i18n.KeySettingsActiveNiches, v)
                }
        }
        if raw, ok := params[protocol.ParamWASlots]; ok {
                if v, ok := raw.(string); ok {
                        s.setConfigItem(i18n.KeySettingsWASlots, v)
                }
        }
        if raw, ok := params[protocol.ParamWorkerPool]; ok {
                if v, ok := raw.(string); ok {
                        s.setConfigItem(i18n.KeySettingsWorkerPool, v)
                }
        }
        if raw, ok := params[protocol.ParamArea]; ok {
                if v, ok := raw.(string); ok {
                        s.setConfigItem(i18n.KeySettingsArea, v)
                }
        }
        if raw, ok := params[protocol.ParamWorkHours]; ok {
                if v, ok := raw.(string); ok {
                        s.setConfigItem(i18n.KeySettingsWorkHours, v)
                }
        }
        if raw, ok := params[protocol.ParamRateLimit]; ok {
                if v, ok := raw.(string); ok {
                        s.setConfigItem(i18n.KeySettingsRateLimit, v)
                }
        }
        if raw, ok := params[protocol.ParamRotatorMode]; ok {
                if v, ok := raw.(string); ok {
                        s.setConfigItem(i18n.KeySettingsRotatorMode, v)
                }
        }
        if raw, ok := params[protocol.ParamAutopilot]; ok {
                if v, ok := raw.(string); ok {
                        s.setConfigItem(i18n.KeySettingsAutopilot, v)
                }
        }

        return nil
}

// setConfigItem updates or appends a config item by key.
// DRY helper — avoids duplicating the find-or-append pattern.
func (s *Settings) setConfigItem(key, value string) {
        for i, item := range s.configItems {
                if item.Label == key {
                        s.configItems[i].Value = value
                        return
                }
        }
        s.configItems = append(s.configItems, kvItem{Label: key, Value: value})
}

func (s *Settings) HandleUpdate(params map[string]any) error {
        if st, ok := params[protocol.ParamState].(string); ok {
                s.state = protocol.StateID(st)
        }
        // Backend sends generic map data — convert to internal TUI types.
        // Do NOT assert TUI types from backend params (frontend/backend concern split).
        if raw, ok := params[protocol.ParamErrors]; ok {
                if list, ok := raw.([]map[string]any); ok {
                        var errs []ConfigError
                        for _, m := range list {
                                ce := ConfigError{}
                                if v, ok := m[protocol.ParamLine].(int); ok {
                                        ce.Line = v
                                }
                                if v, ok := m[protocol.ParamMessage].(string); ok {
                                        ce.Message = v
                                }
                                if v, ok := m[protocol.ParamPointer].(string); ok {
                                        ce.Pointer = v
                                }
                                if rawCtx, ok := m[protocol.ParamContext].([]any); ok {
                                        for _, c := range rawCtx {
                                                if cs, ok := c.(string); ok {
                                                        ce.Context = append(ce.Context, cs)
                                                }
                                        }
                                }
                                errs = append(errs, ce)
                        }
                        if len(errs) > 0 {
                                s.errors = errs
                        }
                }
        }
        if raw, ok := params[protocol.ParamChanges]; ok {
                if list, ok := raw.([]map[string]any); ok {
                        var changes []ChangeEntry
                        for _, m := range list {
                                ch := ChangeEntry{}
                                if v, ok := m[protocol.ParamField].(string); ok {
                                        ch.Field = v
                                }
                                if v, ok := m[protocol.ParamOldValue].(string); ok {
                                        ch.OldValue = v
                                }
                                if v, ok := m[protocol.ParamNewValue].(string); ok {
                                        ch.NewValue = v
                                }
                                changes = append(changes, ch)
                        }
                        if len(changes) > 0 {
                                s.changes = changes
                        }
                }
        }
        return nil
}

func (s *Settings) Focus() { s.focused = true }
func (s *Settings) Blur()  { s.focused = false }

// ---------------------------------------------------------------------------
// Views per state
// ---------------------------------------------------------------------------

// FIX 1: viewOverview() uses ConfigItem with i18n-resolved keys
// FIX S-DRY02: uses kvItem + renderKVSection for file paths
func (s *Settings) viewOverview() string {
        var b strings.Builder

        // Heading
        b.WriteString(style.HeadingStyle.Render(i18n.T(i18n.KeySettingsTitle)))
        b.WriteString(style.Section(style.SectionGap))

        // Description
        b.WriteString(style.BodyStyle.Render(i18n.T(i18n.KeySettingsAllInFiles)))
        b.WriteString(style.Section(style.SectionGap))

        // FIX S-DRY02: file paths using kvItem type (shared with helpers.go)
        filePaths := []kvItem{
                {i18n.T(i18n.KeySettingsConfigMain), i18n.T(i18n.KeySettingsConfigPath)},
                {i18n.T(i18n.KeySettingsTheme), i18n.T(i18n.KeySettingsThemePath)},
                {i18n.T(i18n.KeySettingsQueries), i18n.T(i18n.KeySettingsQueriesPath)},
                {i18n.T(i18n.KeySettingsNicheFolder), i18n.T(i18n.KeySettingsNicheFolderPath)},
        }
        for _, p := range filePaths {
                b.WriteString(style.Indent(1))
                b.WriteString(lipgloss.NewStyle().Foreground(style.TextMuted).Width(18).Render(p.Label))
                b.WriteString(lipgloss.NewStyle().Foreground(style.Text).Render(p.Value))
                b.WriteString("\n")
        }

        // FIX S-04: separator between major sections
        b.WriteString(style.Section(style.SectionGap))
        b.WriteString(style.Separator())
        b.WriteString(style.Section(style.SectionGap))

        // Active config — FIX 1: resolve keys at View() time
        b.WriteString(style.SubHeadingStyle.Render(i18n.T(i18n.KeySettingsActiveConfig)))
        b.WriteString(style.Section(style.SubSectionGap))

        for _, item := range s.configItems {
                b.WriteString(style.Indent(1))
                b.WriteString(lipgloss.NewStyle().Foreground(style.TextMuted).Width(18).Render(item.Label))
                b.WriteString(lipgloss.NewStyle().Foreground(style.Text).Render(item.Value))
                b.WriteString("\n")
        }

        // Key hints
        b.WriteString(style.Section(style.SectionGap))
        b.WriteString(lipgloss.NewStyle().Foreground(style.TextDim).Render(
                fmt.Sprintf("e  %s    r  %s    q  %s",
                        i18n.T(i18n.KeyLabelEdit), i18n.T(i18n.KeyLabelReload), i18n.T(i18n.KeyLabelBack),
                )))

        return b.String()
}

// FIX S-01: viewEdit() heading uses KeySettingsEditConfig instead of KeySettingsEdit
func (s *Settings) viewEdit() string {
        var b strings.Builder

        b.WriteString(style.HeadingStyle.Render(i18n.T(i18n.KeySettingsEditConfig)))
        b.WriteString(style.Section(style.SectionGap))

        b.WriteString(style.BodyStyle.Render(i18n.T(i18n.KeySettingsOpenEditor)))
        b.WriteString(style.Section(style.SectionGap))

        b.WriteString(style.Indent(1))
        b.WriteString(style.BodyStyle.Render("  " + i18n.T(i18n.KeySettingsConfigPath)))
        b.WriteString(style.Section(style.SubSectionGap))

        b.WriteString(style.BodyStyle.Render(i18n.T(i18n.KeySettingsAfterSave)))

        // Key hints
        b.WriteString(style.Section(style.SectionGap))
        b.WriteString(lipgloss.NewStyle().Foreground(style.TextDim).Render(
                fmt.Sprintf("r  %s    q  %s",
                        i18n.T(i18n.KeyLabelReload), i18n.T(i18n.KeyLabelBack),
                )))

        return b.String()
}

// FIX 3: viewReload() key hint uses KeySettingsBackDashboard
func (s *Settings) viewReload() string {
        var b strings.Builder

        b.WriteString(style.HeadingStyle.Render(i18n.T(i18n.KeySettingsReloadOK)))
        b.WriteString(style.Section(style.SectionGap))

        // Changes list
        b.WriteString(style.SubHeadingStyle.Render(i18n.T(i18n.KeySettingsChanges)))
        b.WriteString(style.Section(style.SubSectionGap))

        for _, ch := range s.changes {
                b.WriteString(style.Indent(1))
                b.WriteString(lipgloss.NewStyle().Foreground(style.TextMuted).Render(ch.Field))
                b.WriteString(": ")
                b.WriteString(style.DangerStyle.Render(ch.OldValue))
                b.WriteString(" → ")
                b.WriteString(style.SuccessStyle.Render(ch.NewValue))
                b.WriteString("\n")
        }

        b.WriteString(style.Section(style.SubSectionGap))
        b.WriteString(style.BodyStyle.Render(i18n.T(i18n.KeySettingsApplied)))

        // Key hint — FIX 3: use KeySettingsBackDashboard
        b.WriteString(style.Section(style.SectionGap))
        b.WriteString(lipgloss.NewStyle().Foreground(style.TextDim).Render(
                fmt.Sprintf(" ↵  %s", i18n.T(i18n.KeySettingsBackDashboard),
                )))

        return b.String()
}

// FIX S-02: use len(s.errors) instead of hardcoded "2"
// FIX S-03: add file path display between heading and error count
// FIX S-04: add style.Separator() between major sections
// FIX S-DRY01: use renderGutterError() from helpers.go
func (s *Settings) viewReloadError() string {
        var b strings.Builder

        b.WriteString(style.HeadingStyle.Render(i18n.T(i18n.KeySettingsReloadErr)))
        b.WriteString(style.Section(style.SectionGap))

        // FIX S-03: file path display between heading and error count per doc
        b.WriteString(style.BodyStyle.Render(i18n.T(i18n.KeySettingsFilePath)))
        b.WriteString(style.Section(style.SubSectionGap))

        // FIX S-02: error count uses len(s.errors) instead of hardcoded "2"
        b.WriteString(style.DangerStyle.Render(
                fmt.Sprintf("%d %s:", len(s.errors), i18n.T(i18n.KeySettingsErrorCount)),
        ))
        b.WriteString(style.Section(style.SubSectionGap))

        // FIX S-DRY01: use shared renderGutterError() instead of inline error rendering
        for i, err := range s.errors {
                if i > 0 {
                        b.WriteString(style.Section(style.SubSectionGap))
                }
                linePrefix := fmt.Sprintf("%s %d: %s", i18n.T(i18n.KeyLabelLine), err.Line, err.Message)
                renderGutterError(&b, linePrefix, err.Context, err.Pointer, true)
        }

        // FIX S-04: separator between major sections
        b.WriteString(style.Section(style.SectionGap))
        b.WriteString(style.Separator())
        b.WriteString(style.Section(style.SectionGap))

        // Reassurance
        b.WriteString(style.BodyStyle.Render(i18n.T(i18n.KeySettingsOldConfig)))
        b.WriteString("\n")
        b.WriteString(style.BodyStyle.Render(i18n.T(i18n.KeySettingsFixFirst)))

        // Backup note
        b.WriteString(style.Section(style.SubSectionGap))
        b.WriteString(lipgloss.NewStyle().Foreground(style.TextDim).Render(
                i18n.T(i18n.KeySettingsBackupPath),
        ))
        b.WriteString("\n")
        b.WriteString(lipgloss.NewStyle().Foreground(style.TextDim).Render(
                fmt.Sprintf("(%s)", i18n.T(i18n.KeySettingsBackupNote)),
        ))

        // Key hints
        b.WriteString(style.Section(style.SectionGap))
        b.WriteString(lipgloss.NewStyle().Foreground(style.TextDim).Render(
                fmt.Sprintf("1  %s    2  %s    q  %s",
                        i18n.T(i18n.KeySettingsOpenFile), i18n.T(i18n.KeySettingsRevertBackup), i18n.T(i18n.KeyLabelBack),
                )))

        return b.String()
}
