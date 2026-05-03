package niche

import (
        "fmt"
        "strings"
        "time"

        "github.com/WaClaw-App/waclaw/internal/tui/anim"
        "github.com/WaClaw-App/waclaw/internal/tui/bus"
        "github.com/WaClaw-App/waclaw/internal/tui/component"
        "github.com/WaClaw-App/waclaw/internal/tui/i18n"
        "github.com/WaClaw-App/waclaw/internal/tui/style"
        tui "github.com/WaClaw-App/waclaw/internal/tui"
        "github.com/WaClaw-App/waclaw/pkg/protocol"

        tea "github.com/charmbracelet/bubbletea"
)

// pulseEndMsg is sent after the checkbox pulse animation completes.
type pulseEndMsg struct{}

// errorBlinkMsg is sent to trigger a re-render for the error blink animation.
type errorBlinkMsg struct{}

// SelectModel is the bubbletea.Model for the Niche Select screen (Screen 3).
//
// States per doc/02-screens-niche-select.md:
//   - niche_list: browse and select niches (multi-select)
//   - niche_multi_selected: shows summary of selected niches
//   - niche_custom: instructions for custom niche
//   - niche_edit_filters: preview and edit filters for a selected niche
//   - niche_config_error: display config YAML errors with gutter + pointer
type SelectModel struct {
        base    tui.ScreenBase
        state   protocol.StateID
        width   int
        height  int
        focused bool

        // Niche list state
        items  []NicheItem
        list   component.ListSelect
        cursor int

        // Stagger animation
        staggerStart time.Time

        // Filter preview state (niche_edit_filters)
        filterNiche   string
        filterTargets []string
        filters       []FilterEntry
        areas         []AreaEntry

        // Config error state (niche_config_error)
        errorNiche string
        errorFile  string
        errors     []ConfigError

        // Pulse animation for checkbox toggle
        pulseActive bool
        pulseStart  time.Time
        pulseIndex  int

        // Error blink animation
        errorBlinkStart time.Time
}

// NewSelectModel creates a Niche Select screen model with default values.
// Niche items come from the backend via HandleNavigate/HandleUpdate;
// the model starts empty and is populated when the backend sends data.
func NewSelectModel() SelectModel {
        base := tui.NewScreenBase(protocol.ScreenNicheSelect)

        return SelectModel{
                base:    base,
                state:   protocol.NicheList,
                items:   nil,
                list:    component.NewListSelect(nil),
                cursor:  0,
                focused: true,
        }
}

// ID returns the screen identifier.
func (m SelectModel) ID() protocol.ScreenID { return m.base.ID() }

// SetBus injects the event bus reference.
func (m *SelectModel) SetBus(b *bus.Bus) { m.base.SetBus(b) }

// Bus returns the event bus.
func (m *SelectModel) Bus() *bus.Bus { return m.base.Bus() }

// Focus is called when this screen becomes the active screen.
func (m *SelectModel) Focus() {
        m.focused = true
        m.staggerStart = time.Now()
}

// Blur is called when this screen is no longer the active screen.
func (m *SelectModel) Blur() { m.focused = false }

// ConsumesKey implements tui.KeyConsumer. SelectModel has sub-states (custom,
// edit_filters, config_error) where "q" should navigate back locally to the
// niche list instead of popping the navigation stack.
func (m *SelectModel) ConsumesKey(msg tea.KeyMsg) bool {
        switch msg.String() {
        case "q":
                return m.state == protocol.NicheCustom || m.state == protocol.NicheEditFilters || m.state == protocol.NicheConfigError
        }
        return false
}

// HandleNavigate processes a "navigate" command from the backend.
func (m *SelectModel) HandleNavigate(params map[string]any) error {
        if stateStr, ok := params[protocol.ParamState].(string); ok {
                m.state = protocol.StateID(stateStr)
        }
        if raw, ok := params[protocol.ParamNiches].([]any); ok {
                m.applyNicheData(raw)
        }
        if niche, ok := params[protocol.ParamNicheName].(string); ok {
                m.filterNiche = niche
        }
        if targets, ok := params[protocol.ParamTargets].([]any); ok {
                m.filterTargets = toStringSlice(targets)
        }
        if filters, ok := params[protocol.ParamFilters].([]any); ok {
                m.filters = parseFilterEntries(filters)
        }
        if areas, ok := params[protocol.ParamAreas].([]any); ok {
                m.areas = parseAreaEntries(areas)
        }
        if niche, ok := params[protocol.ParamErrorNiche].(string); ok {
                m.errorNiche = niche
        }
        if file, ok := params[protocol.ParamErrorFile].(string); ok {
                m.errorFile = file
        }
        if errors, ok := params[protocol.ParamErrors].([]any); ok {
                m.applyErrorData(errors)
        }
        m.staggerStart = time.Now()
        return nil
}

// HandleUpdate processes an "update" command from the backend.
func (m *SelectModel) HandleUpdate(params map[string]any) error {
        if stateStr, ok := params[protocol.ParamState].(string); ok {
                m.state = protocol.StateID(stateStr)
        }
        if raw, ok := params[protocol.ParamNiches].([]any); ok {
                m.applyNicheData(raw)
        }
        return nil
}

// Init implements tea.Model.
func (m SelectModel) Init() tea.Cmd { return nil }

// Update implements tea.Model.
func (m SelectModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
        switch msg := msg.(type) {
        case tea.WindowSizeMsg:
                m.width = msg.Width
                m.height = msg.Height
                return m, nil
        case pulseEndMsg:
                m.pulseActive = false
                return m, nil
        case errorBlinkMsg:
                // Re-render to update the blink animation; schedule next blink tick.
                if m.state == protocol.NicheConfigError {
                        return m, tea.Tick(anim.ConfigErrorBlink, func(_ time.Time) tea.Msg {
                                return errorBlinkMsg{}
                        })
                }
                return m, nil
        case tea.KeyMsg:
                return m.handleKey(msg)
        }
        return m, nil
}

// handleKey routes key events based on the current state.
func (m SelectModel) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
        switch m.state {
        case protocol.NicheList, protocol.NicheMultiSelected:
                return m.handleListKey(msg)
        case protocol.NicheCustom:
                return m.handleCustomKey(msg)
        case protocol.NicheEditFilters:
                return m.handleFilterKey(msg)
        case protocol.NicheConfigError:
                return m.handleErrorKey(msg)
        default:
                return m, nil
        }
}

// handleListKey handles key events in the niche list and multi-selected states.
func (m SelectModel) handleListKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
        switch msg.String() {
        case "up", "k":
                m.list.Up()
                m.cursor = m.list.Cursor
        case "down", "j":
                m.list.Down()
                m.cursor = m.list.Cursor
        case " ":
                m.list.Toggle()
                m.pulseActive = true
                m.pulseStart = time.Now()
                m.pulseIndex = m.list.Cursor
                if m.list.SelectedCount() > 0 {
                        m.state = protocol.NicheMultiSelected
                } else {
                        m.state = protocol.NicheList
                }
                for i := range m.items {
                        m.items[i].Selected = m.list.Items[i].Selected
                }
                // Schedule pulse end to reset animation after duration.
                return m, tea.Tick(anim.SuccessPulse, func(_ time.Time) tea.Msg {
                        return pulseEndMsg{}
                })
        case "enter":
                if m.list.SelectedCount() > 0 {
                        publishAction(m.base.Bus(), protocol.ScreenNicheSelect, protocol.ActionNicheProceed, map[string]any{
                                protocol.ParamSelected: m.list.SelectedLabels(),
                        })
                        return m, nil
                }
                if m.cursor == len(m.items)-1 {
                        m.state = protocol.NicheCustom
                        m.staggerStart = time.Now()
                        publishAction(m.base.Bus(), protocol.ScreenNicheSelect, protocol.ActionNicheCustom, nil)
                        return m, nil
                }
        }
        return m, nil
}

// handleCustomKey handles key events in the custom niche state.
func (m SelectModel) handleCustomKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
        switch msg.String() {
        case "r":
                publishAction(m.base.Bus(), protocol.ScreenNicheSelect, protocol.ActionNicheReload, nil)
        case "1", "2", "3", "4", "5", "6", "7", "8", "9":
                        // Dynamic key range: any digit key maps to the item at that index.
                        // No longer hardcoded to "1-5" — adapts to however many items the backend sends.
                idx := int(msg.String()[0] - '1')
                if idx >= 0 && idx < len(m.items) {
                        m.state = protocol.NicheList
                        m.cursor = idx
                        m.list.Cursor = idx
                        for i := range m.list.Items {
                                m.list.Items[i].Focused = (i == idx)
                        }
                        publishAction(m.base.Bus(), protocol.ScreenNicheSelect, protocol.ActionNicheReturnList, map[string]any{
                                protocol.ParamIndex: idx,
                        })
                }
        case "q":
                transitionToList(&m)
        }
        return m, nil
}

// handleFilterKey handles key events in the filter preview state.
func (m SelectModel) handleFilterKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
        switch msg.String() {
        case "enter":
                publishAction(m.base.Bus(), protocol.ScreenNicheSelect, protocol.ActionNicheScrape, map[string]any{
                        protocol.ParamNiche: m.filterNiche,
                })
        case "2":
                publishAction(m.base.Bus(), protocol.ScreenNicheSelect, protocol.ActionNicheEditFilter, map[string]any{
                        protocol.ParamNiche: m.filterNiche,
                })
        case "q":
                transitionToList(&m)
        }
        return m, nil
}

// handleErrorKey handles key events in the config error state.
func (m SelectModel) handleErrorKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
        switch msg.String() {
        case "1":
                // Doc micro-interaction: "1 buka file" = buka $EDITOR langsung ke baris error.
                // Pass the first error line number so backend can open editor at exact position.
                lineNum := 0
                if len(m.errors) > 0 && m.errors[0].Line > 0 {
                        lineNum = m.errors[0].Line
                }
                publishAction(m.base.Bus(), protocol.ScreenNicheSelect, protocol.ActionNicheOpenFile, map[string]any{
                        protocol.ParamFile: m.errorFile,
                        protocol.ParamLine: lineNum,
                })
        case "2":
                publishAction(m.base.Bus(), protocol.ScreenNicheSelect, protocol.ActionNicheShowExample, nil)
        case "r":
                publishAction(m.base.Bus(), protocol.ScreenNicheSelect, protocol.ActionNicheReload, map[string]any{
                        protocol.ParamNiche: m.errorNiche,
                })
        case "q":
                transitionToList(&m)
        }
        return m, nil
}

// View implements tea.Model.
func (m SelectModel) View() string {
        switch m.state {
        case protocol.NicheList:
                return m.viewList()
        case protocol.NicheMultiSelected:
                return m.viewMultiSelected()
        case protocol.NicheCustom:
                return m.viewCustom()
        case protocol.NicheEditFilters:
                return m.viewFilters()
        case protocol.NicheConfigError:
                return m.viewConfigError()
        default:
                return m.viewList()
        }
}

// viewList renders the niche_list state per doc/02-screens-niche-select.md.
func (m SelectModel) viewList() string {
        var b strings.Builder

        b.WriteString(style.HeadingStyle.Render(i18n.T(i18n.KeyNicheSelect)))
        b.WriteString(style.Section(style.SectionGap))

        b.WriteString(style.MutedStyle.Render(i18n.T(i18n.KeyNicheNicheIs)))
        b.WriteString("\n")
        b.WriteString(style.MutedStyle.Render(i18n.T(i18n.KeyNicheWorkerParallel)))
        b.WriteString("\n")
        b.WriteString(style.MutedStyle.Render(i18n.T(i18n.KeyNicheMoreNiche)))
        b.WriteString(style.Section(style.SectionGap))

        for i, item := range m.list.Items {
                if !isStaggerVisible(m.staggerStart, i) {
                        continue
                }

                line := fmt.Sprintf("%d  ", i+1)
                // Show emoji from backend data if available — backend is authoritative.
                if i < len(m.items) && m.items[i].Emoji != "" {
                        line += m.items[i].Emoji + " "
                }
                line += renderCheckbox(item.Selected, m.pulseActive, m.pulseStart, i == m.pulseIndex)
                line += " "
                line += renderStyledLabel(item.Label, item.Focused, item.Selected)

                if item.Description != "" {
                        line += "  " + renderStyledDescription(item.Description, item.Focused)
                }

                b.WriteString(line)
                b.WriteString("\n")
        }

        b.WriteString(style.Section(style.SectionGap))

        b.WriteString(renderFooter([]FooterEntry{
                {Key: "space", Label: i18n.T(i18n.KeyNicheCheckUncheck)},
                {Key: "↵", Label: i18n.T(i18n.KeyNicheGasChecked)},
        }))
        b.WriteString("\n")
        b.WriteString(style.CaptionStyle.Render(i18n.T(i18n.KeyNicheScrapeOwn)))

        return b.String()
}

// viewMultiSelected renders the niche_multi_selected state per doc/02.
func (m SelectModel) viewMultiSelected() string {
        var b strings.Builder

        b.WriteString(style.HeadingStyle.Render(i18n.T(i18n.KeyNicheSelect)))
        b.WriteString(style.Section(style.SectionGap))

        for i, item := range m.list.Items {
                if !isStaggerVisible(m.staggerStart, i) {
                        continue
                }

                line := fmt.Sprintf("%d  ", i+1)
                line += renderCheckbox(item.Selected, false, time.Time{}, false)
                line += " "
                line += renderStyledLabel(item.Label, item.Focused, item.Selected)

                if item.Description != "" {
                        line += "  " + renderStyledDescription(item.Description, item.Focused)
                }

                b.WriteString(line)
                b.WriteString("\n")
        }

        b.WriteString(style.Section(style.SectionGap))

        selectedCount := m.list.SelectedCount()
        b.WriteString(style.BodyStyle.Render(
                fmt.Sprintf(i18n.T(i18n.KeyNicheNicheDipilih), selectedCount),
        ))
        b.WriteString("\n")

        for _, item := range m.items {
                if !item.Selected {
                        continue
                }
                detail := item.Name
                if item.Area != "" {
                        detail += fmt.Sprintf(" — %s", item.Area)
                }
                if item.Templates > 0 {
                        detail += fmt.Sprintf(i18n.T(i18n.KeyNicheTemplateCount), item.Templates)
                }
                b.WriteString(style.MutedStyle.Render("▸ " + detail))
                b.WriteString("\n")
        }

        b.WriteString(style.Section(style.SectionGap))

        b.WriteString(renderFooter([]FooterEntry{
                {Key: "↵", Label: fmt.Sprintf(i18n.T(i18n.KeyNicheGasNiche), selectedCount)},
                {Key: "space", Label: i18n.T(i18n.KeyNicheChange)},
                {Key: "q", Label: i18n.T(i18n.KeyNicheBack)},
        }))
        b.WriteString("\n")
        b.WriteString(style.CaptionStyle.Render(fmt.Sprintf(i18n.T(i18n.KeyNicheMultiParallel), selectedCount)))
        b.WriteString("\n")
        b.WriteString(style.CaptionStyle.Render(i18n.T(i18n.KeyNicheMultiScrapeOwn)))

        return b.String()
}

// viewCustom renders the niche_custom state per doc/02.
func (m SelectModel) viewCustom() string {
        var b strings.Builder

        b.WriteString(style.HeadingStyle.Render(i18n.T(i18n.KeyNicheCustom)))
        b.WriteString(style.Section(style.SectionGap))

        b.WriteString(style.BodyStyle.Render(i18n.T(i18n.KeyNicheCustomDir)))
        b.WriteString("\n")
        b.WriteString(style.MutedStyle.Render("  " + i18n.T(i18n.KeyNicheCustomDirPath)))
        b.WriteString(style.Section(style.SubSectionGap))

        b.WriteString(style.BodyStyle.Render(i18n.T(i18n.KeyNicheCustomMin)))
        b.WriteString("\n")
        b.WriteString(style.MutedStyle.Render("  - "+i18n.T(i18n.KeyNicheNicheYaml)))
        b.WriteString("\n")
        b.WriteString(style.MutedStyle.Render("  - "+i18n.T(i18n.KeyNicheIceBreaker)))
        b.WriteString(style.Section(style.SubSectionGap))

        b.WriteString(style.BodyStyle.Render(i18n.T(i18n.KeyNicheCustomExample)))
        b.WriteString("\n")
        b.WriteString(style.MutedStyle.Render("  " + i18n.T(i18n.KeyNicheCustomExamplePath)))
        b.WriteString(style.Section(style.SubSectionGap))

        b.WriteString(style.BodyStyle.Render(i18n.T(i18n.KeyNicheCustomReady)))
        b.WriteString(style.Section(style.SectionGap))

        b.WriteString(renderFooter([]FooterEntry{
                {Key: "r", Label: i18n.T(i18n.KeyNicheReload)},
                {Key: fmt.Sprintf("1-%d", len(m.items)), Label: i18n.T(i18n.KeyNichePickExisting)},
                {Key: "q", Label: i18n.T(i18n.KeyNicheBack)},
        }))

        return b.String()
}

// viewFilters renders the niche_edit_filters state per doc/02.
func (m SelectModel) viewFilters() string {
        var b strings.Builder

        b.WriteString(style.HeadingStyle.Render(fmt.Sprintf(i18n.T(i18n.KeyNicheLabel), m.filterNiche)))
        b.WriteString(style.Section(style.SectionGap))

        if len(m.filterTargets) > 0 {
                b.WriteString(style.BodyStyle.Render(i18n.T(i18n.KeyNicheTargets)))
                b.WriteString("\n")
                b.WriteString(style.MutedStyle.Render("  " + strings.Join(m.filterTargets, ", ")))
                b.WriteString(style.Section(style.SubSectionGap))
        }

        if len(m.filters) > 0 {
                b.WriteString(style.BodyStyle.Render(i18n.T(i18n.KeyNicheFilterDefault)))
                b.WriteString("\n")
                for _, f := range m.filters {
                        b.WriteString(renderFilterEntry(f))
                        b.WriteString("\n")
                }
                b.WriteString(style.Section(style.SubSectionGap))
        }

        if len(m.areas) > 0 {
                b.WriteString(style.BodyStyle.Render(
                        fmt.Sprintf("%s (%d %s):", i18n.T(i18n.KeyNicheAreaCount), len(m.areas), i18n.T(i18n.KeyNicheAreaKota)),
                ))
                b.WriteString("\n")
                b.WriteString(renderAreaList(m.areas))
        }

        b.WriteString(style.Section(style.SectionGap))
        b.WriteString(renderSeparator())
        b.WriteString(style.Section(style.SectionGap))

        // Doc/02: "sudah pas?   ↵  gas scrape   2  edit filter   q  balik" on one line.
        // Do NOT wrap in CaptionStyle — it overrides inner ActionStyle accent color.
        footer := style.CaptionStyle.Render(i18n.T(i18n.KeyNicheJustRight)) + "   " +
                style.ActionStyle.Render("↵") + " " + style.CaptionStyle.Render(i18n.T(i18n.KeyNicheGasScrape)) + "   " +
                style.ActionStyle.Render("2") + " " + style.CaptionStyle.Render(i18n.T(i18n.KeyNicheEditFilter)) + "   " +
                style.ActionStyle.Render("q") + " " + style.CaptionStyle.Render(i18n.T(i18n.KeyNicheBack))
        b.WriteString(footer)

        return b.String()
}

// viewConfigError renders the niche_config_error state with indentation-based
// error pointer (replacing │ gutter) per doc/02 and borderless design system.
// The error blink animation uses errorBlinkMsg for timer-driven re-renders.
func (m SelectModel) viewConfigError() string {
        var b strings.Builder

        b.WriteString(style.DangerStyle.Bold(true).Render(
                fmt.Sprintf("✗ %s: %s", i18n.T(i18n.KeyNicheConfigErrLabel), m.errorNiche),
        ))
        b.WriteString(style.Section(style.SectionGap))

        b.WriteString(renderSeparator())
        b.WriteString(style.Section(style.SectionGap))

        b.WriteString(style.MutedStyle.Render(m.errorFile))
        b.WriteString(style.Section(style.SubSectionGap))

        b.WriteString(style.BodyStyle.Render(fmt.Sprintf("%d %s", len(m.errors), i18n.T(i18n.KeyNicheProblems))))
        b.WriteString(style.Section(style.SubSectionGap))

        blinkElapsed := time.Since(m.errorBlinkStart)
        blinkBold := blinkElapsed%(anim.ConfigErrorBlink*2) < anim.ConfigErrorBlink

        for _, err := range m.errors {
                b.WriteString(style.DangerStyle.Render("✗ "))
                if err.Line > 0 {
                        b.WriteString(style.DangerStyle.Render(fmt.Sprintf(i18n.T(i18n.KeyNicheLine), err.Line)))
                }
                b.WriteString(style.BodyStyle.Render(err.Message))
                b.WriteString("\n")

                if err.Description != "" {
                        b.WriteString(style.MutedStyle.Render("   " + err.Description))
                        b.WriteString("\n")
                }

                // Render error detail + pointer using borderless indentation
                // instead of │ box-drawing character (design system compliance).
                b.WriteString(renderErrorGutter(err.Detail, err.Pointer, blinkBold))
                b.WriteString("\n")
        }

        b.WriteString(renderSeparator())
        b.WriteString(style.Section(style.SectionGap))

        b.WriteString(style.BodyStyle.Render(
                fmt.Sprintf(i18n.T(i18n.KeyNicheErrPaused), m.errorNiche),
        ))
        b.WriteString("\n")
        b.WriteString(style.MutedStyle.Render(i18n.T(i18n.KeyNicheErrOtherOK)))
        b.WriteString(style.Section(style.SectionGap))

        b.WriteString(renderFooter([]FooterEntry{
                {Key: "1", Label: i18n.T(i18n.KeyNicheOpenFile)},
                {Key: "2", Label: i18n.T(i18n.KeyNicheShowExample)},
                {Key: "r", Label: i18n.T(i18n.KeyNicheReload)},
                {Key: "q", Label: i18n.T(i18n.KeyNicheBack)},
        }))

        return b.String()
}

// applyNicheData converts raw backend niche data into model items.
// Extracts ALL fields the backend sends: name, description, area, templates,
// emoji, targets, selected — the backend is the authoritative source.
func (m *SelectModel) applyNicheData(raw []any) {
        items := make([]NicheItem, 0, len(raw))
        listItems := make([]component.ListItem, 0, len(raw))

        for i, r := range raw {
                data, ok := r.(map[string]any)
                if !ok {
                        continue
                }
                name, _ := data[protocol.ParamName].(string)
                desc, _ := data[protocol.ParamDescription].(string)
                area, _ := data[protocol.ParamArea].(string)
                tmpl := toInt(data[protocol.ParamTemplates])
                emoji, _ := data[protocol.ParamEmoji].(string)
                preSelected, _ := data[protocol.ParamSelected].(bool)
                var targets []string
                if rawTargets, ok := data[protocol.ParamTargets].([]any); ok {
                        targets = toStringSlice(rawTargets)

                }

                items = append(items, NicheItem{
                        Name: name, Description: desc, Area: area, Templates: tmpl,
                        Emoji: emoji, Targets: targets, PreSelected: preSelected,
                })
                listItems = append(listItems, component.ListItem{
                        Label: name, Description: desc, Selected: preSelected,
                        Focused: i == 0, Disabled: false,
                })
        }

        m.items = items
        m.list = component.NewListSelect(listItems)
        m.cursor = 0
}

// applyErrorData converts raw backend error data and starts blink animation.
func (m *SelectModel) applyErrorData(raw []any) {
        m.errors = make([]ConfigError, 0, len(raw))
        m.errorBlinkStart = time.Now()
        for _, r := range raw {
                data, ok := r.(map[string]any)
                if !ok {
                        continue
                }
                m.errors = append(m.errors, ConfigError{
                        Line:        toInt(data[protocol.ParamLine]),
                        Message:     asString(data[protocol.ParamMessage]),
                        Description: asString(data[protocol.ParamDescription]),
                        Detail:      asString(data[protocol.ParamDetail]),
                        Pointer:     asString(data[protocol.ParamPointer]),

                })
        }
}
