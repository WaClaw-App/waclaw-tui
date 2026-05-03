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
        "github.com/charmbracelet/lipgloss"
)

// searchDebounceMsg is sent after the debounce period elapses.
type searchDebounceMsg struct{}

// explorerPulseEndMsg is sent after the explorer checkmark pulse animation completes.
type explorerPulseEndMsg struct{}

// debounceSearch returns a tea.Cmd that fires after the search debounce period.
func debounceSearch() tea.Cmd {
        return tea.Tick(anim.SearchDebounce, func(_ time.Time) tea.Msg {
                return searchDebounceMsg{}
        })
}

// ExplorerModel is the bubbletea.Model for the Niche Explorer screen (Screen 19).
//
// States per doc/11-screens-niche-explorer.md:
//   - explorer_browse: browse categories with emoji
//   - explorer_search: live search with 300ms debounce
//   - explorer_category_detail: detail before generating
//   - explorer_generating: progress bar while generating config
//   - explorer_generated: success with file list
type ExplorerModel struct {
        base    tui.ScreenBase
        state   protocol.StateID
        width   int
        height  int
        focused bool

        // Category browse state
        categories []Category
        cursor     int

        // Stagger animation
        staggerStart time.Time

        // Search state
        searchInput   component.SearchInput
        searchResults []SearchResult
        searchCursor  int
        debounceTimer time.Time

        // Category detail state
        selectedCategory Category
        detailSources    []SourceDetail
        detailAreas      []AreaEntry
        detailFilters    []FilterEntry
        detailTemplates  []GenerateFileStatus // name + detail per doc/11
        areaAutoDetect   bool
        existingAreas    []AreaEntry

        // Generating state
        genProgress    float64 // 0.0 to 1.0
        genFiles       []GenerateFileStatus
        genCurrentStep int
        genStartTime   time.Time

        // Generated state
        genNicheName string
        genFilesDone []GenerateFileStatus
        folderSlug   string // authoritative folder slug from backend (overrides slugify())

        // Pulse animation for checkmark completion
        pulseActive bool
        pulseStart  time.Time
        pulseIndex  int
}

// NewExplorerModel creates a Niche Explorer screen model.
// Categories come from the backend via HandleNavigate/HandleUpdate;
// the model starts empty and is populated when the backend sends data.
func NewExplorerModel() ExplorerModel {
        base := tui.NewScreenBase(protocol.ScreenNicheExplorer)

        return ExplorerModel{
                base:          base,
                state:         protocol.ExplorerBrowse,
                categories:    nil,
                cursor:        0,
                searchInput:   component.NewSearchInput(int(anim.SearchDebounce / time.Millisecond)),
                searchResults: []SearchResult{},
                searchCursor:  0,
                focused:       true,
        }
}

// ID returns the screen identifier.
func (m ExplorerModel) ID() protocol.ScreenID { return m.base.ID() }

// SetBus injects the event bus reference.
func (m *ExplorerModel) SetBus(b *bus.Bus) { m.base.SetBus(b) }

// Bus returns the event bus.
func (m *ExplorerModel) Bus() *bus.Bus { return m.base.Bus() }

// Focus is called when this screen becomes the active screen.
func (m *ExplorerModel) Focus() {
        m.focused = true
        m.staggerStart = time.Now()
}

// Blur is called when this screen is no longer the active screen.
func (m *ExplorerModel) Blur() { m.focused = false }

// HandleNavigate processes a "navigate" command from the backend.
func (m *ExplorerModel) HandleNavigate(params map[string]any) error {
        if stateStr, ok := params[protocol.ParamState].(string); ok {
                m.state = protocol.StateID(stateStr)
        }
        if raw, ok := params[protocol.ParamCategories].([]any); ok {
                m.applyCategoryData(raw)
        }
        if raw, ok := params[protocol.ParamSources].([]any); ok {
                m.applySourceData(raw)
        }
        if raw, ok := params[protocol.ParamAreas].([]any); ok {
                m.detailAreas = parseAreaEntries(raw)
        }
        if raw, ok := params[protocol.ParamFilters].([]any); ok {
                m.detailFilters = parseFilterEntries(raw)
        }
        if raw, ok := params[protocol.ParamTemplates].([]any); ok {
                m.detailTemplates = parseGenFileStatus(raw)
        }
        if raw, ok := params[protocol.ParamGenFiles].([]any); ok {
                m.applyGenFileData(raw)
        }
        if name, ok := params[protocol.ParamGenNicheName].(string); ok {
                m.genNicheName = name
        }
        if slug, ok := params[protocol.ParamFolderSlug].(string); ok {
                m.folderSlug = slug
        }
        if progress, ok := params[protocol.ParamGenProgress].(float64); ok {
                m.genProgress = progress
        }
        if auto, ok := params[protocol.ParamAreaAutoDetect].(bool); ok {
                m.areaAutoDetect = auto
        }
        if raw, ok := params[protocol.ParamExistingAreas].([]any); ok {
                m.existingAreas = parseAreaEntries(raw)
        }
        if raw, ok := params[protocol.ParamSearchResults].([]any); ok {
                m.applySearchResultData(raw)
        }
        // When backend navigates to detail/generating/generated states, it sends
        // category_name so the TUI can render the title correctly even without
        // a prior browse→select interaction.
        if catName, ok := params[protocol.ParamCategoryName].(string); ok {
                m.selectedCategory = m.findCategoryByName(catName)
        }
        // When backend navigates to the generated state, it sends gen_files_done
        // which populates the file tree shown in viewGenerated().
        if raw, ok := params[protocol.ParamGenFilesDone].([]any); ok {
                m.genFilesDone = parseGenFileStatus(raw)
        }
        // If we have gen_files but no gen_files_done and the state is generated,
        // copy gen_files to genFilesDone so the file tree renders.
        if m.state == protocol.ExplorerGenerated && len(m.genFiles) > 0 && len(m.genFilesDone) == 0 {
                m.genFilesDone = m.genFiles
        }
        m.staggerStart = time.Now()
        return nil
}

// HandleUpdate processes an "update" command from the backend.
func (m *ExplorerModel) HandleUpdate(params map[string]any) error {
        if stateStr, ok := params[protocol.ParamState].(string); ok {
                m.state = protocol.StateID(stateStr)
        }
        if progress, ok := params[protocol.ParamGenProgress].(float64); ok {
                m.genProgress = progress
        }
        if step, ok := params[protocol.ParamGenCurrentStep].(float64); ok {
                m.genCurrentStep = int(step)
        }
        if raw, ok := params[protocol.ParamGenFiles].([]any); ok {
                m.applyGenFileData(raw)
        }
        if raw, ok := params[protocol.ParamSearchResults].([]any); ok {
                m.applySearchResultData(raw)
        }
        return nil
}

// Init implements tea.Model.
func (m ExplorerModel) Init() tea.Cmd { return nil }

// Update implements tea.Model.
func (m ExplorerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
        switch msg := msg.(type) {
        case tea.WindowSizeMsg:
                m.width = msg.Width
                m.height = msg.Height
                return m, nil
        case searchDebounceMsg:
                m.performSearch()
                return m, nil
        case explorerPulseEndMsg:
                m.pulseActive = false
                return m, nil
        case tea.KeyMsg:
                return m.handleKey(msg)
        }
        return m, nil
}

// handleKey routes key events based on the current state.
func (m ExplorerModel) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
        switch m.state {
        case protocol.ExplorerBrowse:
                return m.handleBrowseKey(msg)
        case protocol.ExplorerSearch:
                return m.handleSearchKey(msg)
        case protocol.ExplorerCategoryDetail:
                return m.handleDetailKey(msg)
        case protocol.ExplorerGenerating:
                switch msg.String() {
                case "q", "esc":
                        publishAction(m.base.Bus(), protocol.ScreenNicheExplorer, protocol.ActionExplorerCancel, nil)
                }
                return m, nil
        case protocol.ExplorerGenerated:
                return m.handleGeneratedKey(msg)
        default:
                return m.handleBrowseKey(msg)
        }
}

// handleBrowseKey handles key events in the browse state.
func (m ExplorerModel) handleBrowseKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
        switch msg.String() {
        case "up", "k":
                if m.cursor > 0 {
                        m.cursor--
                }
        case "down", "j":
                if m.cursor < len(m.categories)-1 {
                        m.cursor++
                }
        case "/":
                m.state = protocol.ExplorerSearch
                m.searchInput.Clear()
                m.searchInput.Focused = true
                m.searchResults = []SearchResult{}
                m.searchCursor = 0
                m.debounceTimer = time.Now()
                return m, nil
        case "enter":
                if m.cursor < len(m.categories) {
                        m.selectedCategory = m.categories[m.cursor]
                        m.state = protocol.ExplorerCategoryDetail
                        m.staggerStart = time.Now()
                        publishAction(m.base.Bus(), protocol.ScreenNicheExplorer, protocol.ActionExplorerDetail, map[string]any{
                                protocol.ParamCategory: m.selectedCategory.Name,
                        })
                }
        case "q":
                publishAction(m.base.Bus(), protocol.ScreenNicheExplorer, protocol.ActionExplorerBack, nil)
        }
        return m, nil
}

// handleSearchKey handles key events in the search state with debounce.
func (m ExplorerModel) handleSearchKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
        switch msg.String() {
        case "esc":
                m.state = protocol.ExplorerBrowse
                m.searchInput.Clear()
                m.searchResults = []SearchResult{}
                return m, nil
        case "up", "k":
                if m.searchCursor > 0 {
                        m.searchCursor--
                }
                return m, nil
        case "down", "j":
                if m.searchCursor < len(m.searchResults)-1 {
                        m.searchCursor++
                }
                return m, nil
        case "enter":
                if len(m.searchResults) > 0 && m.searchCursor < len(m.searchResults) {
                        sel := m.searchResults[m.searchCursor]
                        m.selectedCategory = Category{Name: sel.Name}
                        m.state = protocol.ExplorerCategoryDetail
                        m.staggerStart = time.Now()
                        publishAction(m.base.Bus(), protocol.ScreenNicheExplorer, protocol.ActionExplorerDetail, map[string]any{
                                protocol.ParamCategory: sel.Name,
                        })
                }
                return m, nil
        case "backspace":
                m.searchInput.Backspace()
                m.debounceTimer = time.Now()
                return m, debounceSearch()
        default:
                if len(msg.String()) == 1 && msg.String()[0] >= 32 {
                        m.searchInput.AppendChar(msg.String())
                        m.debounceTimer = time.Now()
                        return m, debounceSearch()
                }
        }
        return m, nil
}

// performSearch filters categories by the current search input.
// The TUI performs local fuzzy filtering as an instant preview while
// simultaneously requesting backend search results (WA Business Directory +
// GMaps). The backend results will arrive via HandleUpdate and replace/augment
// the local results. This dual-source approach provides immediate feedback
// (local preview) + authoritative results (backend).
func (m *ExplorerModel) performSearch() {
        query := m.searchInput.Value
        if query == "" {
                m.searchResults = []SearchResult{}
                m.searchCursor = 0
                return
        }

        // Local fuzzy filter on already-loaded categories for instant preview.
        names := make([]string, len(m.categories))
        for i, cat := range m.categories {
                names[i] = cat.Name
        }

        results := component.FilterAndSort(query, names)
        m.searchResults = make([]SearchResult, 0, len(results))
        for _, r := range results {
                var cat Category
                for _, c := range m.categories {
                        if c.Name == r.Item {
                                cat = c
                                break
                        }
                }
                m.searchResults = append(m.searchResults, SearchResult{
                        Name: r.Item, SubCount: cat.SubCount,
                        AreaCount: cat.AreaCount, MatchIndices: r.MatchedIndices,
                })
        }
        m.searchCursor = 0

        // Backend search: request WA Business Directory + GMaps results.
        publishAction(m.base.Bus(), protocol.ScreenNicheExplorer, protocol.ActionExplorerSearch, map[string]any{
                protocol.ParamQuery: query,
        })
}

// handleDetailKey handles key events in the category detail state.
func (m ExplorerModel) handleDetailKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
        switch msg.String() {
        case "enter":
                // Transition to generating state — backend drives the actual generation;
                // TUI shows the generating UI optimistically. If backend fails, it will
                // push a navigate/update to return to a previous state.
                m.state = protocol.ExplorerGenerating
                m.genProgress = 0
                m.genCurrentStep = 0
                m.genStartTime = time.Now()
                publishAction(m.base.Bus(), protocol.ScreenNicheExplorer, protocol.ActionExplorerGenerate, map[string]any{
                        protocol.ParamCategory: m.selectedCategory.Name,
                })
        case "2":
                publishAction(m.base.Bus(), protocol.ScreenNicheExplorer, protocol.ActionExplorerEdit, map[string]any{
                        protocol.ParamCategory: m.selectedCategory.Name,
                })
        case "1":
                // Doc/11 area auto-detect variant: "1 tambah area"
                if m.areaAutoDetect {
                        publishAction(m.base.Bus(), protocol.ScreenNicheExplorer, protocol.ActionExplorerAddArea, nil)
                }
        case "q":
                // Doc/11: "q balik" — return to browse state.
                m.state = protocol.ExplorerBrowse
                m.staggerStart = time.Now()
                publishAction(m.base.Bus(), protocol.ScreenNicheExplorer, protocol.ActionExplorerBack, nil)
        }
        return m, nil
}

// handleGeneratedKey handles key events in the generated state.
func (m ExplorerModel) handleGeneratedKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
        busRef := m.base.Bus()
        screenID := protocol.ScreenNicheExplorer
        switch msg.String() {
        case "enter":
                publishAction(busRef, screenID, protocol.ActionExplorerUse, map[string]any{protocol.ParamNiche: m.genNicheName})
        case "1":
                publishAction(busRef, screenID, protocol.ActionExplorerEditConfig, map[string]any{protocol.ParamNiche: m.genNicheName})
        case "2":
                publishAction(busRef, screenID, protocol.ActionExplorerViewTpl, map[string]any{protocol.ParamNiche: m.genNicheName})
        case "r":
                publishAction(busRef, screenID, protocol.ActionNicheReload, map[string]any{protocol.ParamNiche: m.genNicheName})
        case "q":
                // Doc/11: "q balik" — return to browse state.
                m.state = protocol.ExplorerBrowse
                m.staggerStart = time.Now()
                publishAction(busRef, screenID, protocol.ActionExplorerBack, nil)
        }
        return m, nil
}

// View implements tea.Model.
func (m ExplorerModel) View() string {
        switch m.state {
        case protocol.ExplorerBrowse:
                return m.viewBrowse()
        case protocol.ExplorerSearch:
                return m.viewSearch()
        case protocol.ExplorerCategoryDetail:
                return m.viewDetail()
        case protocol.ExplorerGenerating:
                return m.viewGenerating()
        case protocol.ExplorerGenerated:
                return m.viewGenerated()
        default:
                return m.viewBrowse()
        }
}

// viewBrowse renders the explorer_browse state per doc/11.
func (m ExplorerModel) viewBrowse() string {
        var b strings.Builder

        b.WriteString(style.HeadingStyle.Render(i18n.T(i18n.KeyNicheExplorerTitle)))
        b.WriteString(style.Section(style.SectionGap))

        b.WriteString(style.MutedStyle.Render(i18n.T(i18n.KeyNicheExplorerSubtitle)))
        b.WriteString("\n")
        b.WriteString(style.MutedStyle.Render(i18n.T(i18n.KeyNicheExplorerPickConfig)))
        b.WriteString(style.Section(style.SectionGap))

        b.WriteString(renderSectionLabel(i18n.T(i18n.KeyNicheExplorerPopular)))

        for i, cat := range m.categories {
                if !isStaggerVisible(m.staggerStart, i) {
                        continue
                }

                line := fmt.Sprintf("%-2d ", i+1)
                line += cat.Emoji + " "
                line += renderStyledLabel(cat.Name, i == m.cursor, false)

                if len(cat.SubCategories) > 0 {
                        line += "  " + lipgloss.NewStyle().Foreground(style.TextDim).Render(
                                strings.Join(cat.SubCategories, ", "),
                        )
                }

                b.WriteString(line)
                b.WriteString("\n")
        }

        b.WriteString(renderSeparator())
        b.WriteString(style.Section(style.SectionGap))

        b.WriteString(renderFooter([]FooterEntry{
                {Key: "/", Label: i18n.T(i18n.KeyNicheExplorerSearchCat)},
                {Key: "↑↓", Label: i18n.T(i18n.KeyNicheExplorerPick)},
                {Key: "↵", Label: i18n.T(i18n.KeyNicheExplorerDetail)},
                {Key: "q", Label: i18n.T(i18n.KeyNicheBack)},
        }))
        b.WriteString("\n")
        b.WriteString(style.CaptionStyle.Render(i18n.T(i18n.KeyNicheExplorerDIY)))

        return b.String()
}

// viewSearch renders the explorer_search state with borderless search input.
// Uses a dim background strip instead of box-drawing characters (┌─┐│└)
// per vertical borderless design system (doc/16, principle P3).
func (m ExplorerModel) viewSearch() string {
        var b strings.Builder

        title := i18n.T(i18n.KeyNicheExplorerTitle)
        if m.searchInput.Value != "" {
                searching := i18n.T(i18n.KeyNicheExplorerSearching)
                if m.width > 0 {
                        padWidth := m.width - len(i18n.T(i18n.KeyNicheExplorerTitle)) - len(searching)
                        if padWidth > 0 {
                                title += strings.Repeat(" ", padWidth) + searching
                        }
                } else {
                        title += "  " + searching
                }
        }
        b.WriteString(style.HeadingStyle.Render(title))
        b.WriteString(style.Section(style.SectionGap))

        // Borderless search input using dim background strip.
        b.WriteString(renderSearchInput(m.searchInput.Value, m.searchInput.Focused, m.searchInput.Placeholder, m.width))
        b.WriteString(style.Section(style.SectionGap))

        if len(m.searchResults) > 0 {
                b.WriteString(renderSectionLabel(i18n.T(i18n.KeyNicheExplorerResults)))

                for i, result := range m.searchResults {
                        line := fmt.Sprintf("%-2d ", i+1)

                        if len(result.MatchIndices) > 0 {
                                line += component.HighlightMatch(result.Name, result.MatchIndices)
                        } else {
                                line += renderStyledLabel(result.Name, i == m.searchCursor, false)
                        }

                        detail := ""
                        if result.SubCount > 0 {
                                detail += fmt.Sprintf(i18n.T(i18n.KeyNicheExplorerSubKategori), result.SubCount)
                        }
                        if result.AreaCount > 0 {
                                if detail != "" {
                                        detail += " · "
                                }
                                detail += fmt.Sprintf(i18n.T(i18n.KeyNicheExplorerAreaLabel), result.AreaCount)
                        }
                        if detail != "" {
                                line += "  " + lipgloss.NewStyle().Foreground(style.TextDim).Render(detail)
                        }

                        b.WriteString(line)
                        b.WriteString("\n")
                }

                b.WriteString(style.CaptionStyle.Render(i18n.T(i18n.KeyNicheExplorerSource)))
        } else if m.searchInput.Value != "" {
                b.WriteString(style.CaptionStyle.Render(i18n.T(i18n.KeyNicheExplorerNoResult)))
        }

        b.WriteString(style.Section(style.SectionGap))

        b.WriteString(renderFooter([]FooterEntry{
                {Key: "↑↓", Label: i18n.T(i18n.KeyNicheExplorerPick)},
                {Key: "↵", Label: i18n.T(i18n.KeyNicheExplorerDetail)},
                {Key: "esc", Label: i18n.T(i18n.KeyNicheExplorerCancel)},
        }))

        return b.String()
}

// viewDetail renders the explorer_category_detail state per doc/11.
func (m ExplorerModel) viewDetail() string {
        var b strings.Builder

        b.WriteString(style.HeadingStyle.Render(
                fmt.Sprintf(i18n.T(i18n.KeyNicheExplorerTitleDetail), i18n.T(i18n.KeyNicheExplorerTitle), m.selectedCategory.Name)))
        b.WriteString(style.Section(style.SectionGap))

        b.WriteString(renderSeparator())
        b.WriteString(style.Section(style.SectionGap))

        b.WriteString(style.BodyStyle.Render(fmt.Sprintf(i18n.T(i18n.KeyNicheExplorerKategori), m.selectedCategory.Name)))
        b.WriteString("\n")

        if len(m.selectedCategory.SubCategories) > 0 {
                b.WriteString(style.BodyStyle.Render(i18n.T(i18n.KeyNicheExplorerSubCat)+":"))
                b.WriteString("\n")
                b.WriteString(style.MutedStyle.Render("  " + strings.Join(m.selectedCategory.SubCategories, ", ")))
                b.WriteString(style.Section(style.SubSectionGap))
        }

        if len(m.detailSources) > 0 {
                b.WriteString(style.BodyStyle.Render(i18n.T(i18n.KeyNicheExplorerSourceLabel)))
                b.WriteString("\n")
                maxNameLen := 0
                for _, src := range m.detailSources {
                        if len(src.Name) > maxNameLen {
                                maxNameLen = len(src.Name)
                        }
                }
                for _, src := range m.detailSources {
                        b.WriteString(style.MutedStyle.Render("  ● "))
                        paddedName := src.Name + strings.Repeat(" ", maxNameLen-len(src.Name))
                        b.WriteString(style.BodyStyle.Render(paddedName))
                        b.WriteString(style.DimStyle.Render(fmt.Sprintf("  —  %d %s", src.Count, src.Unit)))
                        b.WriteString("\n")
                }
                b.WriteString(style.Section(style.SubSectionGap))
        }

        if len(m.detailAreas) > 0 {
                b.WriteString(style.BodyStyle.Render(i18n.T(i18n.KeyNicheExplorerAreaAuto)))
                b.WriteString("\n")
                var areaStrs []string
                for _, a := range m.detailAreas {
                        areaStrs = append(areaStrs, fmt.Sprintf("%s (%s)", a.City, a.Radius))
                }
                b.WriteString(style.MutedStyle.Render("  " + strings.Join(areaStrs, ", ")))
                b.WriteString(style.Section(style.SubSectionGap))
        }

        if len(m.detailFilters) > 0 {
                b.WriteString(style.BodyStyle.Render(i18n.T(i18n.KeyNicheExplorerFilterDefault)))
                b.WriteString("\n")
                for _, f := range m.detailFilters {
                        b.WriteString(renderFilterEntry(f))
                        b.WriteString("\n")
                }
                b.WriteString(style.Section(style.SubSectionGap))
        }

        if len(m.detailTemplates) > 0 {
                b.WriteString(style.BodyStyle.Render(i18n.T(i18n.KeyNicheTemplateGen)))
                b.WriteString("\n")
                for _, tmpl := range m.detailTemplates {
                        line := style.SuccessStyle.Render("  ✓ " + tmpl.Name)
                        if tmpl.Detail != "" {
                                line += style.MutedStyle.Render("  " + tmpl.Detail)
                        }
                        b.WriteString(line)
                        b.WriteString("\n")
                }
                b.WriteString(style.Section(style.SubSectionGap))
        }

        b.WriteString(renderSeparator())
        b.WriteString(style.Section(style.SectionGap))

        // Doc/11: "sudah pas?   ↵  generate config   2  edit dulu   q  balik" on one line.
        // Do NOT wrap in CaptionStyle — it overrides inner ActionStyle accent color.
        footer := style.CaptionStyle.Render(i18n.T(i18n.KeyNicheJustRight)) + "   " +
                style.ActionStyle.Render("↵") + " " + style.CaptionStyle.Render(i18n.T(i18n.KeyNicheExplorerGenConfig)) + "   " +
                style.ActionStyle.Render("2") + " " + style.CaptionStyle.Render(i18n.T(i18n.KeyNicheExplorerEditFirst)) + "   " +
                style.ActionStyle.Render("q") + " " + style.CaptionStyle.Render(i18n.T(i18n.KeyNicheBack))
        b.WriteString(footer)
        b.WriteString("\n")

        b.WriteString(style.CaptionStyle.Render(
                fmt.Sprintf(i18n.T(i18n.KeyNicheExplorerFolder), folderSlugOrFallback(m.folderSlug, m.selectedCategory.Name)),
        ))

        // Area auto-detect variant per doc/11
        if m.areaAutoDetect && len(m.existingAreas) > 0 {
                b.WriteString(style.Section(style.SectionGap))
                b.WriteString(style.HeadingStyle.Render(i18n.T(i18n.KeyNicheExplorerAreaDetect)))
                b.WriteString(style.Section(style.SubSectionGap))

                b.WriteString(style.BodyStyle.Render(
                        fmt.Sprintf(i18n.T(i18n.KeyNicheExplorerAreaExplain), len(m.existingAreas)),
                ))
                b.WriteString("\n")
                var areaStrs []string
                for _, a := range m.existingAreas {
                        areaStrs = append(areaStrs, fmt.Sprintf("%s (%s)", a.City, a.Radius))
                }
                b.WriteString(style.MutedStyle.Render("  " + strings.Join(areaStrs, ", ")))
                b.WriteString(style.Section(style.SubSectionGap))

                b.WriteString(style.MutedStyle.Render(i18n.T(i18n.KeyNicheExplorerAreaSame)))
                b.WriteString(style.Section(style.SectionGap))

                areaFooter := []FooterEntry{
                        {Key: "1", Label: i18n.T(i18n.KeyNicheExplorerAddArea)},
                        {Key: "↵", Label: i18n.T(i18n.KeyNicheExplorerUseExisting)},
                }
                b.WriteString(renderFooter(areaFooter))
        }

        return b.String()
}

// viewGenerating renders the explorer_generating state with progress bar.
func (m ExplorerModel) viewGenerating() string {
        var b strings.Builder

        b.WriteString(style.HeadingStyle.Render(
                fmt.Sprintf(i18n.T(i18n.KeyNicheExplorerTitleDetail), i18n.T(i18n.KeyNicheExplorerTitle), m.selectedCategory.Name)))
        b.WriteString(style.Section(style.SectionGap))

        b.WriteString(style.BodyStyle.Render(i18n.T(i18n.KeyNicheExplorerGenProgress)))
        b.WriteString(style.Section(style.SectionGap))

        for i, file := range m.genFiles {
                if file.Completed {
                        label := style.SuccessStyle.Render("✓ "+file.Name) + " " + style.MutedStyle.Render(file.Detail)
                        if m.pulseActive && i == m.pulseIndex {
                                if time.Since(m.pulseStart) < anim.SuccessPulse {
                                        label = lipgloss.NewStyle().Foreground(style.Success).Bold(true).Render("✓ "+file.Name) +
                                                " " + style.MutedStyle.Render(file.Detail)
                                }
                        }
                        b.WriteString(label)
                } else if i == m.genCurrentStep {
                        b.WriteString(style.AccentStyle.Render("● ") + style.BodyStyle.Render(file.Name))
                        if file.Detail != "" {
                                b.WriteString(" " + style.MutedStyle.Render(file.Detail))
                        }
                } else {
                        b.WriteString(style.DimStyle.Render("○ " + file.Name))
                        if file.Detail != "" {
                                b.WriteString(" " + style.DimStyle.Render(file.Detail))
                        }
                }
                b.WriteString("\n")
        }

        b.WriteString(style.Section(style.SectionGap))

        bar := component.NewProgressBar(progressBarWidth)
        bar.Percent = m.genProgress
        bar.Label = i18n.T(i18n.KeyNicheExplorerGenBarLabel)
        bar.ShowPercent = false
        b.WriteString(bar.View())

        return b.String()
}

// viewGenerated renders the explorer_generated state with file list.
// Uses indentation-based file tree instead of box-drawing characters (├── └──)
// per vertical borderless design system.
func (m ExplorerModel) viewGenerated() string {
        var b strings.Builder

        b.WriteString(style.HeadingStyle.Render(
                fmt.Sprintf(i18n.T(i18n.KeyNicheExplorerTitleDetail), i18n.T(i18n.KeyNicheExplorerTitle), m.selectedCategory.Name)))
        b.WriteString(style.Section(style.SectionGap))

        b.WriteString(style.SuccessStyle.Bold(true).Render(
                "✓ " + i18n.T(i18n.KeyNicheExplorerGenSuccess)))
        b.WriteString(style.Section(style.SectionGap))

        b.WriteString(renderSeparator())
        b.WriteString(style.Section(style.SectionGap))

        folder := fmt.Sprintf(i18n.T(i18n.KeyNicheExplorerFolderPath), m.genNicheName)
        b.WriteString(style.MutedStyle.Render(folder))
        b.WriteString("\n")

        // Render file tree using borderless indentation instead of ├── └──
        b.WriteString(renderFileTree(m.genFilesDone))

        b.WriteString(style.Section(style.SectionGap))
        b.WriteString(renderSeparator())
        b.WriteString(style.Section(style.SectionGap))

        b.WriteString(style.BodyStyle.Render(i18n.T(i18n.KeyNicheExplorerEditFile)))
        b.WriteString("\n")
        b.WriteString(style.MutedStyle.Render(i18n.T(i18n.KeyNicheExplorerReload)))
        b.WriteString(style.Section(style.SectionGap))

        b.WriteString(renderFooter([]FooterEntry{
                {Key: "↵", Label: i18n.T(i18n.KeyNicheExplorerGasUse)},
                {Key: "1", Label: i18n.T(i18n.KeyNicheExplorerEditConfig)},
                {Key: "2", Label: i18n.T(i18n.KeyNicheExplorerViewTemplate)},
                {Key: "q", Label: i18n.T(i18n.KeyNicheBack)},
        }))
        b.WriteString("\n")
        b.WriteString(style.CaptionStyle.Render(i18n.T(i18n.KeyNicheExplorerParallel)))

        return b.String()
}

// applyCategoryData converts raw backend category data.
func (m *ExplorerModel) applyCategoryData(raw []any) {
        m.categories = make([]Category, 0, len(raw))
        for _, r := range raw {
                data, ok := r.(map[string]any)
                if !ok {
                        continue
                }
                emoji := asString(data[protocol.ParamEmoji])
                name := asString(data[protocol.ParamName])
                subCount := toInt(data[protocol.ParamSubCount])
                areaCount := toInt(data[protocol.ParamAreaCount])
                var subCats []string
                if rawSubs, ok := data[protocol.ParamSubCategories].([]any); ok {
                        subCats = toStringSlice(rawSubs)

                }
                m.categories = append(m.categories, Category{
                        Emoji: emoji, Name: name, SubCategories: subCats,
                        SubCount: subCount, AreaCount: areaCount,
                })
        }
}

// applySourceData converts raw backend source data.
func (m *ExplorerModel) applySourceData(raw []any) {
        m.detailSources = make([]SourceDetail, 0, len(raw))
        for _, r := range raw {
                data, ok := r.(map[string]any)
                if !ok {
                        continue
                }
                m.detailSources = append(m.detailSources, SourceDetail{
                        Name:  asString(data[protocol.ParamName]),
                        Count: toInt(data[protocol.ParamCount]),
                        Unit:  asString(data[protocol.ParamUnit]),

                })
        }
}

// applyGenFileData converts raw backend generation file data.
// DRY: delegates to shared parseGenFileStatus in helpers.go.
func (m *ExplorerModel) applyGenFileData(raw []any) {
        m.genFiles = parseGenFileStatus(raw)
}

// applySearchResultData converts raw backend search result data.
func (m *ExplorerModel) applySearchResultData(raw []any) {
        m.searchResults = make([]SearchResult, 0, len(raw))
        for _, r := range raw {
                data, ok := r.(map[string]any)
                if !ok {
                        continue
                }
                m.searchResults = append(m.searchResults, SearchResult{
                        Name:      asString(data[protocol.ParamName]),
                        SubCount:  toInt(data[protocol.ParamSubCount]),
                        AreaCount: toInt(data[protocol.ParamAreaCount]),

                })
        }
}
