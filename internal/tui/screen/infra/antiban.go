package infra

import (
        "fmt"
        "strings"

        "github.com/WaClaw-App/waclaw/internal/tui"
        "github.com/WaClaw-App/waclaw/internal/tui/bus"
        "github.com/WaClaw-App/waclaw/internal/tui/component"
        "github.com/WaClaw-App/waclaw/internal/tui/i18n"
        "github.com/WaClaw-App/waclaw/internal/tui/style"
        "github.com/WaClaw-App/waclaw/pkg/protocol"
        tea "github.com/charmbracelet/bubbletea"
        "github.com/charmbracelet/lipgloss"
)

// ---------------------------------------------------------------------------
// Data types for the Shield (Anti-Ban) screen
// ---------------------------------------------------------------------------

// SlotInfo holds data for a single WhatsApp number slot.
type SlotInfo struct {
        Number    string
        Status    string // "active", "cooldown", "flagged"
        SentHour  int
        MaxHour   int
        Cooldown  string
        ReadyIn   string
        SentToday int
        Warnings  int
        Healthy   bool
}

// SlotHistoryEntry represents a single event in a slot's history.
type SlotHistoryEntry struct {
        Time   string
        Event  string
        Detail string
        Level  string // "", "warning", "danger"
}

// FIX A-DRY01: removed duplicate kvItem struct and renderKVSection function.
// They now live in helpers.go and are shared across the infra package.

// ---------------------------------------------------------------------------
// Shield screen model
// ---------------------------------------------------------------------------

// Shield implements tui.Screen for Screen 12: Anti-Ban Shield.
// It displays the aggregate health score, per-slot WA rotator status,
// rate limiting details, work hours guard, pattern guard, spam guard,
// and ban risk assessment.
type Shield struct {
        tui.ScreenBase
        state        protocol.StateID
        healthScore  int
        slots        []SlotInfo
        selectedSlot int
        history      []SlotHistoryEntry
        width        int
        height       int
        focused      bool

        // --- Backend-driven config values (populated via HandleNavigate/HandleUpdate) ---
        // These replace the previously hardcoded values in viewOverview() and viewSettings().

        // dailyBudgetSent is the number of messages sent today across all slots.
        dailyBudgetSent int

        // dailyBudgetTotal is the daily message limit across all slots.
        dailyBudgetTotal int

        // selectedSlotHealth is the health score of the currently selected slot detail.
        selectedSlotHealth int

        // slotDetail7DaySent is the 7-day sent count for the slot detail view.
        slotDetail7DaySent int

        // slotDetail7DayResponded is the 7-day responded count for the slot detail view.
        slotDetail7DayResponded int

        // slotDetail7DayFailed is the 7-day failed count for the slot detail view.
        slotDetail7DayFailed string

        // slotDetail7DayWarnings is the 7-day warning count for the slot detail view.
        slotDetail7DayWarnings int

        // configAntiBan holds anti-ban settings from backend (replaces hardcoded viewSettings values).
        configAntiBan map[string]string

        // configSpamGuard holds spam guard settings from backend.
        configSpamGuard map[string]string

        // workHours from backend config (e.g. "09:00-17:00 wib").
        workHours string

        // warningHealthScore is the health score for the warning state view (default 71 for demo).
        warningHealthScore int

        // dangerHealthScore is the health score for the danger state view (default 38 for demo).
        dangerHealthScore int

        // flaggedSlot is the slot info shown in the danger state view.
        flaggedSlot SlotInfo

        // --- Backend-driven numeric values for i18n format strings ---
        // Backend is the authoritative source; these replace hardcoded numbers in i18n values.

        // perSlotHourly is the per-slot hourly rate limit (e.g. 6).
        perSlotHourly int

        // minDelayMin is the minimum delay between messages in minutes (e.g. 8).
        minDelayMin int

        // delayVariancePct is the delay variance percentage (e.g. 30).
        delayVariancePct int

        // healthThreshold is the health score auto-pause threshold (e.g. 50).
        healthThreshold int

        // perLeadLifetime is the max messages per lead lifetime (e.g. 3).
        perLeadLifetime int

        // msgIntervalHours is the minimum hours between messages to same lead (e.g. 24).
        msgIntervalHours int

        // recontactDelayDays is days before re-contacting after no-deal response (e.g. 7).
        recontactDelayDays int

        // dncCount is the number of entries in the do-not-contact block list (e.g. 12).
        dncCount int

        // iceBreakerVariants is the number of ice_breaker template variants (e.g. 3).
        iceBreakerVariants int

        // offerVariants is the number of offer template variants per niche (e.g. 3).
        offerVariants int

        // workHoursDuration is the number of work hours per day (e.g. 8).
        workHoursDuration int

        // healthRecoveryPts is the health points recovered per day without warning (e.g. 5).
        healthRecoveryPts int

        // timezoneShort is the short timezone label (e.g. "wib").
        timezoneShort string

        // timezoneFull is the full timezone name (e.g. "asia/jakarta").
        timezoneFull string

        // currentTimeStr is the current time display string (e.g. "14:23").
        currentTimeStr string

        // slotDetail7DayFailedCount is the 7-day failed count for the slot detail view.
        slotDetail7DayFailedCount int

        // banRiskLevel is the ban risk assessment level from the backend ("low", "medium", "high").
        // Replaces the hardcoded "low" previously passed to renderBanRisk().
        banRiskLevel string

        // banRiskEmoji is the emoji for the ban risk level from the backend.
        banRiskEmoji string

        // riskIndicators holds backend-driven risk indicator statuses.
        // Each entry maps an indicator key to true (pass/✓) or false (fail/✗).
        // Replaces the previously hardcoded all-✓ list in renderIndicators().
        riskIndicators map[string]bool

        // --- Backend-driven display string values ---

        // rateLimitPerSlot is the per-slot rate limit display string from backend.
        rateLimitPerSlot string

        // rateLimitPerDay is the per-day rate limit display string from backend.
        rateLimitPerDay string

        // rateLimitPerNumber is the per-number rate limit display string from backend.
        rateLimitPerNumber string

        // rateLimitPerLead is the per-lead rate limit display string from backend.
        rateLimitPerLead string

        // timezone is the timezone display string from backend.
        timezone string

        // sendHours is the allowed sending hours display string from backend.
        sendHours string

        // scrapeHours is the allowed scraping hours display string from backend.
        scrapeHours string

        // nowInWorkHours indicates whether currently within work hours (from backend).
        nowInWorkHours string
}

// NewShield creates a Shield screen with zero/empty defaults.
// All values are zero/nil until the backend provides data via HandleNavigate.
// The backend is the single source of truth — no hardcoded demo defaults.
// The viewOverview() guard renders a waiting placeholder when healthScore
// is 0 and no slots are present, so empty initial state is safe.
func NewShield() *Shield {
        return &Shield{
                ScreenBase:   tui.NewScreenBase(protocol.ScreenAntiBan),
                state:        protocol.ShieldOverview,
                configAntiBan:   make(map[string]string),
                configSpamGuard: make(map[string]string),

                // Ban risk defaults — overridden by backend via populateConfigFromBackend
                banRiskLevel: "low",
                banRiskEmoji: "🟢",
                flaggedSlot:  SlotInfo{Status: "flagged"},
        }
}

// REMOVED: defaultAntiBanConfig(), defaultSpamGuardConfig(), demoSlots(),
// and demoSlotHistory(). These previously hardcoded display values and demo
// data in the TUI violated the DRY principle — the backend is the single
// source of truth. All data flows through HandleNavigate/HandleUpdate
// from the backend scenario engine (mock.go ShieldConfigData/ShieldData).

// ---------------------------------------------------------------------------
// tea.Model interface
// ---------------------------------------------------------------------------

func (s *Shield) Init() tea.Cmd { return nil }

func (s *Shield) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
        switch m := msg.(type) {
        case tea.WindowSizeMsg:
                s.width = m.Width
                s.height = m.Height
                return s, nil
        }

        if key, ok := msg.(tea.KeyMsg); ok {
                switch key.String() {
                case "q":
                        if s.state == protocol.ShieldSlotDetail || s.state == protocol.ShieldSettings {
                                s.state = protocol.ShieldOverview
                        }
                        // In overview/warning/danger, let global handler pop the screen.
                case "up", "k":
                        if s.state == protocol.ShieldOverview || s.state == protocol.ShieldWarning || s.state == protocol.ShieldDanger {
                                if s.selectedSlot > 0 {
                                        s.selectedSlot--
                                }
                        }
                case "down", "j":
                        if s.state == protocol.ShieldOverview || s.state == protocol.ShieldWarning || s.state == protocol.ShieldDanger {
                                if s.selectedSlot < len(s.slots)-1 {
                                        s.selectedSlot++
                                }
                        }
                case "enter":
                        if s.selectedSlot < len(s.slots) {
                                s.state = protocol.ShieldSlotDetail
                        }
                case "r":
                        // Refresh — re-fetch from backend
                        if s.Bus() != nil {
                                s.Bus().Publish(bus.ActionMsg{
                                        Action: string(protocol.ActionRefreshShield),
                                        Screen: s.ID(),
                                })
                        }
                case "1":
                        if s.state == protocol.ShieldDanger {
                                // Let it be
                        }
                case "2":
                        if s.state == protocol.ShieldDanger {
                                // Add number
                                if s.Bus() != nil {
                                        s.Bus().Publish(bus.ActionMsg{Action: string(protocol.ActionAddWANumber), Screen: s.ID()})
                                }
                        }
                case "3":
                        if s.state == protocol.ShieldDanger {
                                // Pause sending
                                if s.Bus() != nil {
                                        s.Bus().Publish(bus.ActionMsg{Action: string(protocol.ActionPauseSending), Screen: s.ID()})
                                }
                        }
                case "e":
                        s.state = protocol.ShieldSettings
                case "esc":
                        if s.state == protocol.ShieldSlotDetail || s.state == protocol.ShieldSettings {
                                s.state = protocol.ShieldOverview
                        }
                }
        }
        return s, nil
}

func (s *Shield) View() string {
        switch s.state {
        case protocol.ShieldOverview:
                return s.viewOverview()
        case protocol.ShieldWarning:
                return s.viewWarning()
        case protocol.ShieldDanger:
                return s.viewDanger()
        case protocol.ShieldSlotDetail:
                return s.viewSlotDetail()
        case protocol.ShieldSettings:
                return s.viewSettings()
        default:
                return s.viewOverview()
        }
}

// ---------------------------------------------------------------------------
// Screen interface
// ---------------------------------------------------------------------------

func (s *Shield) HandleNavigate(params map[string]any) error {
        applyNavigateState(&s.state, params)
        if h, ok := params[protocol.ParamHealth].(int); ok {
                s.healthScore = h
        }
        if h, ok := params[protocol.ParamHealthScore].(int); ok {
                s.healthScore = h
        }

        // Populate daily budget from backend (FIX A-03).
        if v, ok := params[protocol.ParamDailyBudgetSent].(int); ok {
                s.dailyBudgetSent = v
        }
        if v, ok := params[protocol.ParamDailyBudgetTotal].(int); ok {
                s.dailyBudgetTotal = v
        }

        // Populate work hours from backend config.
        if v, ok := params[protocol.ParamWorkHours].(string); ok {
                s.workHours = v
        }

        // Populate warning/danger health scores from backend.
        if v, ok := params[protocol.ParamWarningHealthScore].(int); ok {
                s.warningHealthScore = v
        }
        if v, ok := params[protocol.ParamDangerHealthScore].(int); ok {
                s.dangerHealthScore = v
        }

        // Populate flagged slot data from backend.
        if raw, ok := params[protocol.ParamFlaggedSlot].(map[string]any); ok {
                s.flaggedSlot = SlotInfo{Status: "flagged"}
                if v, ok := raw[protocol.ParamSlotNumber].(string); ok {
                        s.flaggedSlot.Number = v
                }
        }

        // Populate slot detail stats from backend.
        if v, ok := params[protocol.ParamSlotDetail7DaySent].(int); ok {
                s.slotDetail7DaySent = v
        }
        if v, ok := params[protocol.ParamSlotDetail7DayResponded].(int); ok {
                s.slotDetail7DayResponded = v
        }
        if v, ok := params[protocol.ParamSlotDetail7DayFailed].(string); ok {
                s.slotDetail7DayFailed = v
        }
        if v, ok := params[protocol.ParamSlotDetail7DayWarnings].(int); ok {
                s.slotDetail7DayWarnings = v
        }
        if v, ok := params[protocol.ParamSelectedSlotHealth].(int); ok {
                s.selectedSlotHealth = v
        }

        // Populate shield config data from backend (replaces hardcoded viewSettings values).
        s.populateConfigFromBackend(params)

        // Populate rate limiting display strings from backend.
        if v, ok := params[protocol.ParamRateLimitPerSlot].(string); ok {
                s.rateLimitPerSlot = v
        }
        if v, ok := params[protocol.ParamRateLimitPerDay].(string); ok {
                s.rateLimitPerDay = v
        }
        if v, ok := params[protocol.ParamRateLimitPerNumber].(string); ok {
                s.rateLimitPerNumber = v
        }
        if v, ok := params[protocol.ParamRateLimitPerLead].(string); ok {
                s.rateLimitPerLead = v
        }

        // Populate work hours display strings from backend.
        if v, ok := params[protocol.ParamTimezone].(string); ok {
                s.timezone = v
        }
        if v, ok := params[protocol.ParamSendHours].(string); ok {
                s.sendHours = v
        }
        if v, ok := params[protocol.ParamScrapeHours].(string); ok {
                s.scrapeHours = v
        }
        if v, ok := params[protocol.ParamNowInWorkHours].(string); ok {
                s.nowInWorkHours = v
        }

        // Populate slot history from backend.
        if raw, ok := params[protocol.ParamSlotHistory].([]any); ok {
                var history []SlotHistoryEntry
                for _, r := range raw {
                        if m, ok := r.(map[string]any); ok {
                                e := SlotHistoryEntry{}
                                if v, ok := m["time"].(string); ok {
                                        e.Time = v
                                }
                                if v, ok := m["event"].(string); ok {
                                        e.Event = v
                                }
                                if v, ok := m["detail"].(string); ok {
                                        e.Detail = v
                                }
                                if v, ok := m["level"].(string); ok {
                                        e.Level = v
                                }
                                history = append(history, e)
                        }
                }
                if len(history) > 0 {
                        s.history = history
                }
        }

        return nil
}

// populateConfigFromBackend reads backend params and updates the config maps.
// DRY helper — single place for all backend→config mapping, replacing
// the previously hardcoded values in viewSettings().
func (s *Shield) populateConfigFromBackend(params map[string]any) {
        // Lazy-init nil maps so writes don't panic.
        if s.configAntiBan == nil {
                s.configAntiBan = make(map[string]string)
        }
        if s.configSpamGuard == nil {
                s.configSpamGuard = make(map[string]string)
        }

        // Anti-ban config overrides
        if v, ok := params[protocol.ParamPerSlotHourly].(int); ok {
                s.configAntiBan["per_slot_hourly"] = fmt.Sprintf("%d", v)
                s.perSlotHourly = v // Backend is authoritative source
        }
        if v, ok := params[protocol.ParamPerSlotDaily].(int); ok {
                s.configAntiBan["per_slot_daily"] = fmt.Sprintf("%d/hari", v)
        }
        if v, ok := params[protocol.ParamCooldownMin].(int); ok {
                s.configAntiBan["cooldown_limit"] = fmt.Sprintf("%d menit", v)
        }
        if v, ok := params[protocol.ParamMinDelayMin].(int); ok {
                s.configAntiBan["min_delay_min"] = fmt.Sprintf("%d menit", v)
                s.minDelayMin = v // Backend is authoritative source
        }
        if v, ok := params[protocol.ParamMaxDelayMin].(int); ok {
                s.configAntiBan["max_delay_min"] = fmt.Sprintf("%d menit", v)
        }
        if v, ok := params[protocol.ParamDelayVariancePct].(int); ok {
                s.configAntiBan["delay_variance_pct"] = fmt.Sprintf("%d%%  (random ±)", v)
                s.delayVariancePct = v // Backend is authoritative source
        }
        if v, ok := params[protocol.ParamWorkHours].(string); ok {
                s.configAntiBan["work_hours"] = v
        }
        if v, ok := params[protocol.ParamAutoPause].(string); ok {
                s.configAntiBan["auto_pause"] = v
        }
        if v, ok := params[protocol.ParamHealthThreshold].(int); ok {
                s.configAntiBan["health_threshold"] = fmt.Sprintf("%d/100  (auto-pause)", v)
                s.healthThreshold = v // Backend is authoritative source
        }
        if v, ok := params[protocol.ParamRotatorMode].(string); ok {
                s.configAntiBan["rotator_mode"] = v
        }
        if v, ok := params[protocol.ParamTemplateRotation].(string); ok {
                s.configAntiBan["template_rotation"] = v
        }
        if v, ok := params[protocol.ParamRotationMode].(string); ok {
                s.configAntiBan["rotation_mode"] = v
        }
        if v, ok := params[protocol.ParamEmojiVariation].(string); ok {
                s.configAntiBan["emoji_variation"] = v
        }
        if v, ok := params[protocol.ParamParagraphShuffle].(string); ok {
                s.configAntiBan["paragraph_shuffle"] = v
        }

        // Spam guard config overrides
        if v, ok := params[protocol.ParamPerLeadLifetime].(int); ok {
                s.configSpamGuard["per_lead_lifetime"] = fmt.Sprintf("%d  (lifetime)", v)
                s.perLeadLifetime = v // Backend is authoritative source
        }
        if v, ok := params[protocol.ParamMsgIntervalHours].(int); ok {
                s.configSpamGuard["msg_interval_hours"] = fmt.Sprintf("%d  (min jam antar pesan)", v)
                s.msgIntervalHours = v // Backend is authoritative source
        }
        if v, ok := params[protocol.ParamFollowupDelayDays].(int); ok {
                s.configSpamGuard["followup_delay_days"] = fmt.Sprintf("%d  (min hari antar follow-up)", v)
        }
        if v, ok := params[protocol.ParamFollowupVariant].(string); ok {
                s.configSpamGuard["followup_variant"] = v
        }
        if v, ok := params[protocol.ParamColdAfter].(int); ok {
                s.configSpamGuard["cold_after"] = fmt.Sprintf("%d  (→ dingin)", v)
        }
        if v, ok := params[protocol.ParamRecontactDelayDays].(int); ok {
                s.configSpamGuard["recontact_delay_days"] = fmt.Sprintf("%d  (setelah response tanpa deal)", v)
                s.recontactDelayDays = v // Backend is authoritative source
        }
        if v, ok := params[protocol.ParamAutoBlock].(string); ok {
                s.configSpamGuard["auto_block"] = v
        }
        if v, ok := params[protocol.ParamDupCrossNiche].(string); ok {
                s.configSpamGuard["dup_cross_niche"] = v
        }
        if v, ok := params[protocol.ParamWAPreValidation].(string); ok {
                s.configSpamGuard["wa_pre_validation"] = v
        }
        if v, ok := params[protocol.ParamWAValidationMethod].(string); ok {
                s.configSpamGuard["wa_validation_method"] = v
        }

        // Additional i18n format string params — backend is authoritative source
        if v, ok := params[protocol.ParamDNCCount].(int); ok {
                s.dncCount = v
        }
        if v, ok := params[protocol.ParamIceBreakerVariants].(int); ok {
                s.iceBreakerVariants = v
        }
        if v, ok := params[protocol.ParamOfferVariants].(int); ok {
                s.offerVariants = v
        }
        if v, ok := params[protocol.ParamWorkHoursDuration].(int); ok {
                s.workHoursDuration = v
        }
        if v, ok := params[protocol.ParamHealthRecoveryPts].(int); ok {
                s.healthRecoveryPts = v
        }
        if v, ok := params[protocol.ParamTimezoneShort].(string); ok {
                s.timezoneShort = v
        }
        if v, ok := params[protocol.ParamTimezoneFull].(string); ok {
                s.timezoneFull = v
        }
        if v, ok := params[protocol.ParamCurrentTime].(string); ok {
                s.currentTimeStr = v
        }
        if v, ok := params[protocol.ParamFailedCount].(int); ok {
                s.slotDetail7DayFailedCount = v
        }

        // Ban risk assessment — backend is authoritative source
        if v, ok := params["ban_risk_level"].(string); ok {
                s.banRiskLevel = v
        }
        if v, ok := params["ban_risk_emoji"].(string); ok {
                s.banRiskEmoji = v
        }
        if raw, ok := params["risk_indicators"].(map[string]any); ok {
                if s.riskIndicators == nil {
                        s.riskIndicators = make(map[string]bool, len(raw))
                }
                for k, v := range raw {
                        if bv, ok := v.(bool); ok {
                                s.riskIndicators[k] = bv
                        }
                }
        }
}

func (s *Shield) HandleUpdate(params map[string]any) error {
        if h, ok := params[protocol.ParamHealth].(int); ok {
                s.healthScore = h
        }
        // Backend sends generic map data — convert to internal TUI types.
        // Do NOT assert TUI types from backend params (frontend/backend concern split).
        if raw, ok := params[protocol.ParamSlots]; ok {
                if list, ok := raw.([]map[string]any); ok {
                        var slots []SlotInfo
                        for _, m := range list {
                                si := SlotInfo{}
                                if v, ok := m[protocol.ParamSlotNumber].(string); ok {
                                        si.Number = v
                                }
                                if v, ok := m[protocol.ParamSlotStatus].(string); ok {
                                        si.Status = v
                                }
                                if v, ok := m[protocol.ParamSlotSentHour].(int); ok {
                                        si.SentHour = v
                                }
                                if v, ok := m[protocol.ParamSlotMaxHour].(int); ok {
                                        si.MaxHour = v
                                }
                                if v, ok := m[protocol.ParamSlotSentToday].(int); ok {
                                        si.SentToday = v
                                }
                                if v, ok := m[protocol.ParamSlotWarnings].(int); ok {
                                        si.Warnings = v
                                }
                                if v, ok := m[protocol.ParamSlotCooldown].(string); ok {
                                        si.Cooldown = v
                                }
                                if v, ok := m[protocol.ParamSlotReadyIn].(string); ok {
                                        si.ReadyIn = v
                                }
                                if v, ok := m[protocol.ParamSlotHealthy].(bool); ok {
                                        si.Healthy = v
                                }
                                slots = append(slots, si)
                        }
                        if len(slots) > 0 {
                                s.slots = slots
                        }
                }
        }
        return nil
}

func (s *Shield) Focus() { s.focused = true }
func (s *Shield) Blur()  { s.focused = false }

// ConsumesKey implements tui.KeyConsumer. The Shield screen has sub-states
// (SlotDetail, Settings) where "q" should navigate back locally to the
// overview instead of popping the navigation stack.
func (s *Shield) ConsumesKey(msg tea.KeyMsg) bool {
        switch msg.String() {
        case "q":
                // In sub-states, "q" goes back locally.
                return s.state == protocol.ShieldSlotDetail || s.state == protocol.ShieldSettings
        }
        return false
}

// ---------------------------------------------------------------------------
// Shared view helpers
// ---------------------------------------------------------------------------

// renderShieldArt wraps the component.ShieldArt with a health score bar.
func (s *Shield) renderShieldArt(health int) string {
        shield := component.NewShieldArt(health)
        return shield.View()
}

// renderSlotCard renders a single WA slot card.
func (s *Shield) renderSlotCard(slot SlotInfo, index int) string {
        var b strings.Builder

        focused := index == s.selectedSlot
        nameColor := style.TextMuted
        if focused {
                nameColor = style.Text
        }

        b.WriteString(style.Indent(1))
        b.WriteString(fmt.Sprintf("📱 slot-%d ", index+1))

        // Number and status
        b.WriteString(lipgloss.NewStyle().Foreground(nameColor).Bold(focused).Render(slot.Number))
        b.WriteString("   ")

        switch slot.Status {
        case protocol.SlotStatusActive:
                b.WriteString(style.SuccessStyle.Render(fmt.Sprintf("● %s", i18n.T(i18n.KeyStatusActive))))
        case protocol.SlotStatusCooldown:
                b.WriteString(lipgloss.NewStyle().Foreground(style.TextMuted).Render(fmt.Sprintf("○ %s", i18n.T(i18n.KeyStatusCooldown))))
        case protocol.SlotStatusFlagged:
                b.WriteString(style.DangerStyle.Render(fmt.Sprintf("✗ %s", i18n.T(i18n.KeyStatusFlagged))))
        }
        b.WriteString("\n")

        // Progress bar (if active or cooldown)
        if slot.Status != protocol.SlotStatusFlagged {
                bar := component.NewProgressBar(style.CompactProgressBarWidth)
                bar.Percent = float64(slot.SentHour) / float64(slot.MaxHour)
                bar.ShowPercent = false

                var detail string
                if slot.Status == protocol.SlotStatusCooldown {
                        detail = fmt.Sprintf(i18n.T(i18n.KeyShieldReadyFmt), slot.ReadyIn)
                } else {
                        detail = fmt.Sprintf("%d/%d %s   %s: %s", slot.SentHour, slot.MaxHour, i18n.T(i18n.KeyShieldHour), i18n.T(i18n.KeyShieldCooldown), slot.Cooldown)
                }
                bar.Label = detail

                b.WriteString(style.Indent(2))
                b.WriteString(bar.View())
                b.WriteString("\n")
        }

        // Daily stats
        if slot.Status != protocol.SlotStatusFlagged {
                warningText := i18n.T(i18n.KeyShieldWarningZero)
                if slot.Warnings > 0 {
                        warningText = fmt.Sprintf(i18n.T(i18n.KeyShieldWarningCountFmt), slot.Warnings)
                }
                b.WriteString(style.Indent(2))
                b.WriteString(lipgloss.NewStyle().Foreground(style.TextDim).Render(
                        fmt.Sprintf("%s: %d %s · %s", i18n.T(i18n.KeyWorkersToday), slot.SentToday, i18n.T(i18n.KeyShieldSentStat), warningText),
                ))
                b.WriteString("\n")
        }

        return b.String()
}

// renderBanRisk renders the ban risk score section.
// FIX A-02: renderIndicators() is called inside/after renderBanRisk, not as a separate section.
func (s *Shield) renderBanRisk(level string, emoji, riskText, advice string) string {
        var b strings.Builder

        var riskColor lipgloss.Color
        switch level {
        case "low":
                riskColor = style.Success
        case "medium":
                riskColor = style.Warning
        default:
                riskColor = style.Danger
        }

        b.WriteString(style.SubHeadingStyle.Render(i18n.T(i18n.KeyShieldBanRiskScore)))
        b.WriteString(style.Section(style.SubSectionGap))

        b.WriteString(style.Indent(1))
        b.WriteString(lipgloss.NewStyle().Foreground(riskColor).Render(fmt.Sprintf("%s  %s", emoji, riskText)))
        b.WriteString("\n")

        if advice != "" {
                b.WriteString(style.Section(style.SubSectionGap))
                b.WriteString(lipgloss.NewStyle().Foreground(style.TextMuted).Render(advice))
        }

        // FIX A-02: indicators are part of the ban risk section per doc
        b.WriteString(style.Section(style.SubSectionGap))
        b.WriteString(s.renderIndicators())

        return b.String()
}

// renderIndicators returns the 7 risk indicators with backend-driven pass/fail status.
// When the backend provides riskIndicators data, each indicator shows ✓ (pass) or ✗ (fail)
// based on the backend's assessment. Falls back to all-✓ when no backend data yet.
func (s *Shield) renderIndicators() string {
        // Canonical indicator keys in display order — must match what the backend sends.
        indicatorKeys := []string{
                "even_distribution",
                "cooldown_ok",
                "template_varied",
                "no_overload",
                "work_hours_ok",
                "spam_guard_active",
                "dnc_respected",
        }
        // Display labels for each indicator key (i18n-resolved).
        indicatorLabels := map[string]string{
                "even_distribution":  i18n.T(i18n.KeyShieldEvenDist),
                "cooldown_ok":       i18n.T(i18n.KeyShieldCooldownOk),
                "template_varied":   i18n.T(i18n.KeyShieldTemplateVaried),
                "no_overload":       i18n.T(i18n.KeyShieldNoOverload),
                "work_hours_ok":     i18n.T(i18n.KeyShieldWorkHoursOk),
                "spam_guard_active": i18n.T(i18n.KeyShieldSpamGuardActive),
                "dnc_respected":     i18n.T(i18n.KeyShieldDNCRespected),
        }

        var b strings.Builder
        for i, key := range indicatorKeys {
                label := indicatorLabels[key]
                if s.riskIndicators != nil {
                        // Backend-driven: use actual pass/fail status
                        if s.riskIndicators[key] {
                                b.WriteString(lipgloss.NewStyle().Foreground(style.Success).Render(fmt.Sprintf("  ✓  %s", label)))
                        } else {
                                b.WriteString(lipgloss.NewStyle().Foreground(style.Danger).Render(fmt.Sprintf("  ✗  %s", label)))
                        }
                } else {
                        // No backend data yet — assume all pass (will be overridden on first navigate)
                        b.WriteString(lipgloss.NewStyle().Foreground(style.Success).Render(fmt.Sprintf("  ✓  %s", label)))
                }
                if i < len(indicatorKeys)-1 {
                        b.WriteString("\n")
                }
        }
        return b.String()
}

// activeSlotCount returns the number of non-flagged slots.
// Backend is the authoritative source for slot data.
func (s *Shield) activeSlotCount() int {
        count := 0
        for _, slot := range s.slots {
                if slot.Status != protocol.SlotStatusFlagged {
                        count++
                }
        }
        if count == 0 {
                count = len(s.slots) // fallback
        }
        return count
}

// FIX 2: renderPatternGuard() with 5 pattern guard items
func (s *Shield) renderPatternGuard() string {
        // Use configAntiBan map (populated from backend) with i18n fallbacks
        templateRotVal := s.configAntiBan["template_rotation"]
        if templateRotVal == "" {
                templateRotVal = i18n.T(i18n.KeyShieldPatternTemplate)
        }
        rotModeVal := s.configAntiBan["rotation_mode"]
        if rotModeVal == "" {
                rotModeVal = i18n.T(i18n.KeyShieldPatternTimeVar)
        }
        delayVarVal := s.configAntiBan["delay_variance_pct"]
        if delayVarVal == "" {
                delayVarVal = i18n.T(i18n.KeyShieldPatternMsgVar)
        }
        emojiVal := s.configAntiBan["emoji_variation"]
        if emojiVal == "" {
                emojiVal = i18n.T(i18n.KeyShieldPatternEmoji)
        }
        shuffleVal := s.configAntiBan["paragraph_shuffle"]
        if shuffleVal == "" {
                shuffleVal = i18n.T(i18n.KeyShieldPatternShuffle)
        }
        items := []kvItem{
                {i18n.T(i18n.KeyShieldPatternTemplateLabel), templateRotVal},
                {i18n.T(i18n.KeyShieldPatternTimeVarLabel), rotModeVal},
                {i18n.T(i18n.KeyShieldPatternMsgVarLabel), delayVarVal},
                {i18n.T(i18n.KeyShieldPatternEmojiLabel), emojiVal},
                {i18n.T(i18n.KeyShieldPatternShuffleLabel), shuffleVal},
        }
        var b strings.Builder
        renderKVSection(&b, i18n.T(i18n.KeyShieldPatternGuard), items)
        return b.String()
}

// FIX 3: renderSpamGuard() with 6 items in key-value label format
func (s *Shield) renderSpamGuard() string {
        // Use configSpamGuard map (populated from backend) with i18n fallbacks
        perLeadVal := s.configSpamGuard["per_lead_lifetime"]
        if perLeadVal == "" {
                perLeadVal = i18n.T(i18n.KeyShieldSpamPerLead)
        }
        intervalVal := s.configSpamGuard["msg_interval_hours"]
        if intervalVal == "" {
                intervalVal = i18n.T(i18n.KeyShieldSpamLifetime)
        }
        dncVal := s.configSpamGuard["cold_after"]
        if dncVal == "" {
                dncVal = i18n.T(i18n.KeyShieldDNCCount)
        }
        stopVal := s.configSpamGuard["auto_block"]
        if stopVal == "" {
                stopVal = i18n.T(i18n.KeyShieldStopDet)
        }
        dupVal := s.configSpamGuard["dup_cross_niche"]
        if dupVal == "" {
                dupVal = i18n.T(i18n.KeyShieldDupGuard)
        }
        recontactVal := s.configSpamGuard["recontact_delay_days"]
        if recontactVal == "" {
                recontactVal = i18n.T(i18n.KeyShieldRecontact)
        }
        items := []kvItem{
                {i18n.T(i18n.KeyShieldSpamPerLeadLabel), perLeadVal},
                {i18n.T(i18n.KeyShieldSpamLifetimeLabel), intervalVal},
                {i18n.T(i18n.KeyShieldDNCCountLabel), dncVal},
                {i18n.T(i18n.KeyShieldStopDetLabel), stopVal},
                {i18n.T(i18n.KeyShieldDupGuardLabel), dupVal},
                {i18n.T(i18n.KeyShieldRecontactLabel), recontactVal},
        }
        var b strings.Builder
        renderKVSection(&b, i18n.T(i18n.KeyShieldSpamGuard), items)
        return b.String()
}

// ---------------------------------------------------------------------------
// Views per state
// ---------------------------------------------------------------------------

func (s *Shield) viewOverview() string {
        var b strings.Builder

        // Waiting-for-data guard: show placeholder when no backend data yet.
        if s.healthScore == 0 && len(s.slots) == 0 {
                b.WriteString(style.HeadingStyle.Render(i18n.T(i18n.KeyShieldTitle)))
                b.WriteString("\n\n")
                b.WriteString(style.CaptionStyle.Render(i18n.T(i18n.KeyShieldWaitingData)))
                return b.String()
        }

        // Heading
        b.WriteString(renderHeadWithStatus(i18n.T(i18n.KeyShieldTitle), i18n.T(i18n.KeyShieldAllSafeEmoji), style.Success))
        b.WriteString(style.Section(style.SectionGap))

        // Shield art
        b.WriteString(s.renderShieldArt(s.healthScore))
        b.WriteString(style.Section(style.SectionGap))

        // WA rotator — simplified alignment
        b.WriteString(style.SubHeadingStyle.Render(
                fmt.Sprintf("%s  %d %s", i18n.T(i18n.KeyShieldWARotator), len(s.slots), i18n.T(i18n.KeyShieldNumbers)),
        ))
        b.WriteString(style.Section(style.SubSectionGap))

        for i, slot := range s.slots {
                if i > 0 {
                        b.WriteString(style.Section(style.SubSectionGap))
                }
                b.WriteString(s.renderSlotCard(slot, i))

                // FIX 4: show "status: sehat ✓" under each active/cooldown slot
                if slot.Status != protocol.SlotStatusFlagged {
                        b.WriteString(style.Indent(2))
                        b.WriteString(lipgloss.NewStyle().Foreground(style.Success).Render(i18n.T(i18n.KeyShieldHealthyCheck)))
                        b.WriteString("\n")
                }
        }

        // Rate limiting
        b.WriteString(style.Section(style.SectionGap))
        b.WriteString(style.SubHeadingStyle.Render(i18n.T(i18n.KeyShieldRateLimiting)))
        b.WriteString(style.Section(style.SubSectionGap))

        // Rate limiting — backend-driven values with i18n fallbacks
        perSlotVal := s.rateLimitPerSlot
        if perSlotVal == "" {
                perSlotVal = i18n.T(i18n.KeyShieldPerSlotDemo)
        }
        perDayVal := s.rateLimitPerDay
        if perDayVal == "" {
                perDayVal = i18n.T(i18n.KeyShieldPerDayDetail)
        }
        perNumberVal := s.rateLimitPerNumber
        if perNumberVal == "" {
                perNumberVal = i18n.T(i18n.KeyShieldPerNumber)
        }
        perLeadVal := s.rateLimitPerLead
        if perLeadVal == "" {
                perLeadVal = i18n.T(i18n.KeyShieldPerLeadDetail)
        }
        b.WriteString(style.Indent(1))
        b.WriteString(lipgloss.NewStyle().Foreground(style.TextMuted).Render(
                fmt.Sprintf("%s:  %s", i18n.T(i18n.KeyShieldPerSlot), fmt.Sprintf(i18n.T(i18n.KeyShieldPerSlotDemo), s.perSlotHourly, s.activeSlotCount(), s.perSlotHourly*s.activeSlotCount())), // Backend is authoritative source
        ))
        b.WriteString("\n")
        b.WriteString(style.Indent(1))
        b.WriteString(lipgloss.NewStyle().Foreground(style.TextMuted).Render(fmt.Sprintf(i18n.T(i18n.KeyShieldPerDayDetail), s.dailyBudgetTotal, s.dailyBudgetSent))) // Backend is authoritative source
        b.WriteString("\n")
        b.WriteString(style.Indent(1))
        b.WriteString(lipgloss.NewStyle().Foreground(style.TextMuted).Render(
                fmt.Sprintf("%s: %s", i18n.T(i18n.KeyShieldPerNumberLabel), fmt.Sprintf(i18n.T(i18n.KeyShieldPerNumber), s.minDelayMin)), // Backend is authoritative source
        ))
        b.WriteString("\n")
        b.WriteString(style.Indent(1))
        b.WriteString(lipgloss.NewStyle().Foreground(style.TextMuted).Render(
                fmt.Sprintf("%s: %s", i18n.T(i18n.KeyShieldPerLead), fmt.Sprintf(i18n.T(i18n.KeyShieldPerLeadDetail), s.perLeadLifetime, s.msgIntervalHours)),
        ))
        b.WriteString("\n")

        // Daily budget bar
        b.WriteString(style.Section(style.SubSectionGap))
        bar := component.NewProgressBar(style.DefaultProgressBarWidth)
        if s.dailyBudgetTotal > 0 {
                bar.Percent = float64(s.dailyBudgetSent) / float64(s.dailyBudgetTotal)
                bar.Label = fmt.Sprintf("%d/%d", s.dailyBudgetSent, s.dailyBudgetTotal)
        } else {
                bar.Percent = 0
                bar.Label = "0/0"
        }
        bar.ShowPercent = true
        b.WriteString(style.Indent(1))
        b.WriteString(lipgloss.NewStyle().Foreground(style.TextMuted).Width(14).Render(i18n.T(i18n.KeyShieldDailyBudget)))
        b.WriteString(bar.View())
        b.WriteString("\n")

        // Work hours
        b.WriteString(style.Section(style.SectionGap))
        b.WriteString(style.SubHeadingStyle.Render(i18n.T(i18n.KeyShieldWorkHoursGuard)))
        b.WriteString(style.Section(style.SubSectionGap))
        // Work hours — backend-driven values with i18n fallbacks
        tzVal := s.timezone
        if tzVal == "" {
                tzVal = i18n.T(i18n.KeyShieldTimezoneDemo)
        }
        sendHrsVal := s.sendHours
        if sendHrsVal == "" {
                sendHrsVal = i18n.T(i18n.KeyShieldSendHoursDemo)
        }
        scrapeHrsVal := s.scrapeHours
        if scrapeHrsVal == "" {
                scrapeHrsVal = i18n.T(i18n.KeyShieldScrapeHoursDemo)
        }
        nowVal := s.nowInWorkHours
        if nowVal == "" {
                nowVal = i18n.T(i18n.KeyShieldNowInWorkDemo)
        }
        b.WriteString(style.Indent(1))
        b.WriteString(lipgloss.NewStyle().Foreground(style.TextMuted).Render(
                fmt.Sprintf("%s: %s", i18n.T(i18n.KeyShieldTimezone), fmt.Sprintf(i18n.T(i18n.KeyShieldTimezoneDemo), s.timezoneShort, s.timezoneFull)), // Backend is authoritative source
        ))
        b.WriteString("\n")
        b.WriteString(style.Indent(1))
        b.WriteString(lipgloss.NewStyle().Foreground(style.TextMuted).Render(
                fmt.Sprintf("%s: %s", i18n.T(i18n.KeyShieldSendHours), fmt.Sprintf(i18n.T(i18n.KeyShieldSendHoursDemo), s.workHours, s.workHoursDuration)), // Backend is authoritative source
        ))
        b.WriteString("\n")
        b.WriteString(style.Indent(1))
        b.WriteString(lipgloss.NewStyle().Foreground(style.TextMuted).Render(
                fmt.Sprintf("%s: %s", i18n.T(i18n.KeyShieldScrapeHours), scrapeHrsVal),
        ))
        b.WriteString("\n")
        b.WriteString(style.Indent(1))
        b.WriteString(lipgloss.NewStyle().Foreground(style.TextMuted).Render(
                fmt.Sprintf("%s: %s %s", i18n.T(i18n.KeyShieldNow), fmt.Sprintf(i18n.T(i18n.KeyShieldNowInWorkDemo), s.currentTimeStr), i18n.T(i18n.KeyShieldInWorkHours)), // Backend is authoritative source
        ))

        // FIX 2: Pattern guard
        b.WriteString(style.Section(style.SectionGap))
        b.WriteString(s.renderPatternGuard())

        // FIX 3: Spam guard
        b.WriteString(style.Section(style.SectionGap))
        b.WriteString(s.renderSpamGuard())

        // Ban risk — FIX A-02: indicators are now inside renderBanRisk, not a separate section
        b.WriteString(style.Section(style.SectionGap))
        riskText := i18n.T(i18n.KeyShieldRiskLow)
        if s.banRiskLevel == "medium" {
                riskText = i18n.T(i18n.KeyShieldRiskMedium)
        } else if s.banRiskLevel == "high" {
                riskText = i18n.T(i18n.KeyShieldRiskHigh)
        }
        b.WriteString(s.renderBanRisk(s.banRiskLevel, s.banRiskEmoji, riskText, ""))

        // Key hints
        b.WriteString(style.Section(style.SectionGap))
        b.WriteString(lipgloss.NewStyle().Foreground(style.TextDim).Render(
                fmt.Sprintf("↵  %s    r  %s    q  %s",
                        i18n.T(i18n.KeyShieldSlotDetail), i18n.T(i18n.KeyLabelRefresh),
                        i18n.T(i18n.KeyLabelBack),
                )))

        return b.String()
}

// FIX 5: viewWarning() uses renderShieldArtWithHealth instead of mutating s.healthScore
func (s *Shield) viewWarning() string {
        var b strings.Builder

        b.WriteString(renderHeadWithStatus(i18n.T(i18n.KeyShieldTitle), i18n.T(i18n.KeyShieldWarningEmoji), style.Warning))
        b.WriteString(style.Section(style.SectionGap))

        // FIX 5: pass health directly instead of mutating s.healthScore
        b.WriteString(s.renderShieldArt(s.warningHealthScore))
        b.WriteString(style.Section(style.SectionGap))

        // WA rotator — simplified alignment
        b.WriteString(style.SubHeadingStyle.Render(fmt.Sprintf("%s  %d %s", i18n.T(i18n.KeyShieldWARotator), len(s.slots), i18n.T(i18n.KeyShieldNumbers))))
        b.WriteString(style.Section(style.SubSectionGap))

        for i, slot := range s.slots {
                if i > 0 {
                        b.WriteString(style.Section(style.SubSectionGap))
                }
                b.WriteString(s.renderSlotCard(slot, i))

                // FIX 4: healthy check under active/cooldown slots
                if slot.Status != protocol.SlotStatusFlagged {
                        b.WriteString(style.Indent(2))
                        b.WriteString(lipgloss.NewStyle().Foreground(style.Success).Render(i18n.T(i18n.KeyShieldHealthyCheck)))
                        b.WriteString("\n")
                }

                // Show warning on first slot
                if i == 0 {
                        b.WriteString(style.Indent(2))
                        b.WriteString(style.WarningStyle.Render(fmt.Sprintf(i18n.T(i18n.KeyShieldTooManyHour), slot.SentHour, slot.MaxHour))) // Backend is authoritative source
                        b.WriteString("\n")
                        b.WriteString(style.Indent(2))
                        b.WriteString(lipgloss.NewStyle().Foreground(style.TextMuted).Render(i18n.T(i18n.KeyShieldAutoReduce)))
                        b.WriteString("\n")
                }
        }

        // Ban risk — FIX A-02: indicators now inside renderBanRisk
        b.WriteString(style.Section(style.SectionGap))
        b.WriteString(s.renderBanRisk("medium", "🟡", i18n.T(i18n.KeyShieldRiskMedium),
                i18n.T(i18n.KeyShieldAlreadyMoved),
        ))

        b.WriteString(style.Section(style.SectionGap))
        b.WriteString(lipgloss.NewStyle().Foreground(style.TextDim).Render(
                fmt.Sprintf("↵  %s    r  %s    q  %s",
                        i18n.T(i18n.KeyShieldSlotDetail), i18n.T(i18n.KeyLabelRefresh), i18n.T(i18n.KeyLabelBack),
                )))

        return b.String()
}

// FIX 5: viewDanger() uses renderShieldArtWithHealth
// FIX 6: viewDanger() shows slot-2 and slot-3 with "menggantikan beban" text
// FIX 7: viewDanger() adds recommendation section
func (s *Shield) viewDanger() string {
        var b strings.Builder

        b.WriteString(renderHeadWithStatus(i18n.T(i18n.KeyShieldTitle), i18n.T(i18n.KeyShieldDangerEmoji), style.Danger))
        b.WriteString(style.Section(style.SectionGap))

        // Cracked shield — FIX 5: pass health directly
        b.WriteString(s.renderShieldArt(s.dangerHealthScore))
        b.WriteString(style.Section(style.SectionGap))

        // WA rotator — simplified alignment
        b.WriteString(style.SubHeadingStyle.Render(fmt.Sprintf("%s  %d %s", i18n.T(i18n.KeyShieldWARotator), len(s.slots), i18n.T(i18n.KeyShieldNumbers))))
        b.WriteString(style.Section(style.SubSectionGap))

        // Flagged slot — use backend-provided flagged slot data
        b.WriteString(s.renderSlotCard(s.flaggedSlot, 0))
        b.WriteString(style.Section(style.SubSectionGap))

        b.WriteString(style.Indent(2))
        b.WriteString(style.DangerStyle.Render(i18n.T(i18n.KeyShieldFlaggedMsg)))
        b.WriteString("\n")
        b.WriteString(style.Indent(2))
        b.WriteString(lipgloss.NewStyle().Foreground(style.TextMuted).Render(i18n.T(i18n.KeyShieldPossibility)))
        b.WriteString("\n")
        b.WriteString(style.Indent(2))
        b.WriteString(lipgloss.NewStyle().Foreground(style.TextMuted).Render(fmt.Sprintf("%s %s", i18n.T(i18n.KeyShieldAction), fmt.Sprintf(i18n.T(i18n.KeyShieldSlotAutoPaused), 1))))
        b.WriteString("\n")
        b.WriteString(style.Indent(2))
        b.WriteString(lipgloss.NewStyle().Foreground(style.TextMuted).Render(fmt.Sprintf(i18n.T(i18n.KeyShieldAllMsgsMoved), 2, 3)))

        // FIX 6: slot-2 and slot-3 with "menggantikan beban" text
        b.WriteString(style.Section(style.SectionGap))
        b.WriteString(s.renderSlotCard(s.slots[1], 1))
        b.WriteString("\n")
        b.WriteString(style.Indent(2))
        b.WriteString(lipgloss.NewStyle().Foreground(style.TextMuted).Render(fmt.Sprintf(i18n.T(i18n.KeyShieldReplacing), 1)))

        b.WriteString(style.Section(style.SubSectionGap))
        b.WriteString(s.renderSlotCard(s.slots[2], 2))
        b.WriteString("\n")
        b.WriteString(style.Indent(2))
        b.WriteString(lipgloss.NewStyle().Foreground(style.TextMuted).Render(fmt.Sprintf(i18n.T(i18n.KeyShieldReplacing), 1)))

        // Ban risk — FIX A-02: indicators now inside renderBanRisk
        b.WriteString(style.Section(style.SectionGap))
        b.WriteString(s.renderBanRisk("high", "🔴", i18n.T(i18n.KeyShieldRiskHigh),
                i18n.T(i18n.KeyShieldAlreadyMoved),
        ))

        // FIX 7: recommendation section
        b.WriteString(style.Section(style.SubSectionGap))
        b.WriteString(style.SubHeadingStyle.Render(i18n.T(i18n.KeyShieldRecommendation)))
        b.WriteString(style.Section(style.SubSectionGap))
        b.WriteString(style.Indent(1))
        b.WriteString(lipgloss.NewStyle().Foreground(style.TextMuted).Render(fmt.Sprintf("1  %s", i18n.T(i18n.KeyShieldRec1))))
        b.WriteString("\n")
        b.WriteString(style.Indent(1))
        b.WriteString(lipgloss.NewStyle().Foreground(style.TextMuted).Render(fmt.Sprintf("2  %s", i18n.T(i18n.KeyShieldRec2))))
        b.WriteString("\n")
        b.WriteString(style.Indent(1))
        b.WriteString(lipgloss.NewStyle().Foreground(style.TextMuted).Render(fmt.Sprintf("3  %s", i18n.T(i18n.KeyShieldRec3))))

        b.WriteString(style.Section(style.SubSectionGap))
        b.WriteString(lipgloss.NewStyle().Foreground(style.TextDim).Render(
                fmt.Sprintf("↵  %s    2  %s    3  %s",
                        i18n.T(i18n.KeyShieldLetItBe), i18n.T(i18n.KeyShieldAddNumber), i18n.T(i18n.KeyShieldPauseSend),
                )))

        return b.String()
}

func (s *Shield) viewSlotDetail() string {
        if s.selectedSlot >= len(s.slots) {
                return s.viewOverview()
        }
        slot := s.slots[s.selectedSlot]

        var b strings.Builder

        b.WriteString(renderHeadWithStatus(fmt.Sprintf("%s: %s", i18n.T(i18n.KeyShieldSlotNumber), slot.Number), fmt.Sprintf("● %s", i18n.T(i18n.KeyStatusActive)), style.Success))
        b.WriteString(style.Section(style.SectionGap))

        // 7-day stats
        b.WriteString(style.SubHeadingStyle.Render(i18n.T(i18n.KeyShieldStats7Day)))
        b.WriteString(style.Section(style.SubSectionGap))
        respondPct := ""
        if s.slotDetail7DaySent > 0 && s.slotDetail7DayResponded > 0 {
                respondPct = fmt.Sprintf(" (%d%%)", s.slotDetail7DayResponded*100/s.slotDetail7DaySent)
        }
        stats := []struct{ label, value string }{
                {i18n.T(i18n.KeyShieldStatSent), fmt.Sprintf("%d", s.slotDetail7DaySent)},
                {i18n.T(i18n.KeyShieldStatRespond), fmt.Sprintf("%d%s", s.slotDetail7DayResponded, respondPct)},
                {i18n.T(i18n.KeyShieldStatFailed), fmt.Sprintf(i18n.T(i18n.KeyShieldStatFailedDemo), s.slotDetail7DayFailedCount)}, // Backend is authoritative source
                {i18n.T(i18n.KeyShieldStatWarning), fmt.Sprintf("%d", s.slotDetail7DayWarnings)},
        }
        for _, st := range stats {
                b.WriteString(style.Indent(1))
                b.WriteString(lipgloss.NewStyle().Foreground(style.TextMuted).Width(12).Render(st.label))
                b.WriteString(st.value)
                b.WriteString("\n")
        }

        // History
        b.WriteString(style.Section(style.SectionGap))
        b.WriteString(style.SubHeadingStyle.Render(i18n.T(i18n.KeyShieldHistoryLabel)))
        b.WriteString(style.Section(style.SubSectionGap))
        for _, entry := range s.history {
                b.WriteString(style.Indent(1))
                b.WriteString(lipgloss.NewStyle().Foreground(style.TextDim).Width(16).Render(entry.Time))
                color := style.Text
                if entry.Level == "warning" {
                        color = style.Warning
                }
                b.WriteString(lipgloss.NewStyle().Foreground(color).Render(entry.Event))
                b.WriteString("\n")
        }

        // Health score
        b.WriteString(style.Section(style.SectionGap))
        b.WriteString(style.Indent(1))
        b.WriteString(lipgloss.NewStyle().Foreground(style.Success).Bold(true).Render(
                fmt.Sprintf(i18n.T(i18n.KeyShieldHealthScoreFmt), s.selectedSlotHealth),
        ))
        b.WriteString("\n")
        b.WriteString(style.Indent(1))
        b.WriteString(lipgloss.NewStyle().Foreground(style.TextDim).Render(fmt.Sprintf(i18n.T(i18n.KeyShieldHealthBelow50), s.healthThreshold))) // Backend is authoritative source
        b.WriteString("\n")
        b.WriteString(style.Indent(1))
        b.WriteString(lipgloss.NewStyle().Foreground(style.TextDim).Render(fmt.Sprintf(i18n.T(i18n.KeyShieldHealthUp5), s.healthRecoveryPts))) // Backend is authoritative source

        b.WriteString(style.Section(style.SectionGap))
        b.WriteString(lipgloss.NewStyle().Foreground(style.TextDim).Render(
                fmt.Sprintf("1  %s    2  %s    q  %s",
                        i18n.T(i18n.KeyShieldPauseNumber), i18n.T(i18n.KeyWorkersViewLeads), i18n.T(i18n.KeyLabelBack),
                )))

        return b.String()
}

// viewSettings() uses i18n-keyed kvItem arrays via renderKVSection.
// Values come from backend config maps (configAntiBan/configSpamGuard),
// populated via populateConfigFromBackend — no hardcoded display values.
// FIX A-DRY01: duplicate kvItem/renderKVSection removed.
func (s *Shield) viewSettings() string {
        var b strings.Builder

        b.WriteString(style.HeadingStyle.Render(i18n.T(i18n.KeyShieldTitle)))
        b.WriteString(style.Section(style.SubSectionGap))

        b.WriteString(style.BodyStyle.Render(i18n.T(i18n.KeyShieldConfig)))
        b.WriteString(style.Section(style.SectionGap))

        // Anti-ban section — values from backend config maps (populated via populateConfigFromBackend)
        antiBanItems := []kvItem{
                {i18n.T(i18n.KeyShieldPerSlot), s.configAntiBan["per_slot_hourly"]},
                {i18n.T(i18n.KeyShieldPerDay), s.configAntiBan["per_slot_daily"]},
                {i18n.T(i18n.KeyShieldMinDelay), s.configAntiBan["min_delay_min"]},
                {i18n.T(i18n.KeyShieldMaxDelay), s.configAntiBan["max_delay_min"]},
                {i18n.T(i18n.KeyShieldDelayVariance), s.configAntiBan["delay_variance_pct"]},
                {i18n.T(i18n.KeyShieldCooldownLimit), s.configAntiBan["cooldown_limit"]},
                {i18n.T(i18n.KeyShieldWorkHoursGuard), s.configAntiBan["work_hours"]},
                {i18n.T(i18n.KeyShieldAutoPause), s.configAntiBan["auto_pause"]},
                {i18n.T(i18n.KeyShieldHealthThreshold), s.configAntiBan["health_threshold"]},
                {i18n.T(i18n.KeyShieldRotatorMode), s.configAntiBan["rotator_mode"]},
                {i18n.T(i18n.KeyShieldTemplateRotation), s.configAntiBan["template_rotation"]},
                {i18n.T(i18n.KeyShieldRotationMode), s.configAntiBan["rotation_mode"]},
                {i18n.T(i18n.KeyShieldEmojiVariation), s.configAntiBan["emoji_variation"]},
                {i18n.T(i18n.KeyShieldParagraphShuffle), s.configAntiBan["paragraph_shuffle"]},
        }
        renderKVSection(&b, i18n.T(i18n.KeyShieldSectionAntiBan), antiBanItems)

        // Spam guard section — values from backend config maps (populated via populateConfigFromBackend)
        spamItems := []kvItem{
                {i18n.T(i18n.KeyShieldSpamPerLead), s.configSpamGuard["per_lead_lifetime"]},
                {i18n.T(i18n.KeyShieldMsgInterval), s.configSpamGuard["msg_interval_hours"]},
                {i18n.T(i18n.KeyShieldFollowupDelay), s.configSpamGuard["followup_delay_days"]},
                {i18n.T(i18n.KeyShieldFollowupVariant), s.configSpamGuard["followup_variant"]},
                {i18n.T(i18n.KeyShieldColdAfter), s.configSpamGuard["cold_after"]},
                {i18n.T(i18n.KeyShieldRecontactDelay), s.configSpamGuard["recontact_delay_days"]},
                {i18n.T(i18n.KeyShieldAutoBlock), s.configSpamGuard["auto_block"]},
                {i18n.T(i18n.KeyShieldDupCrossNiche), s.configSpamGuard["dup_cross_niche"]},
                {i18n.T(i18n.KeyShieldWAPreValidation), s.configSpamGuard["wa_pre_validation"]},
                {i18n.T(i18n.KeyShieldWAValidationMethod), s.configSpamGuard["wa_validation_method"]},
        }
        renderKVSection(&b, i18n.T(i18n.KeyShieldSectionSpamGuard), spamItems)

        // FIX 8: closing_triggers section with 5 items
        // P3 compliance: use style.SectionLabel() instead of ─ box-drawing chars
        b.WriteString(style.Section(style.SubSectionGap))
        b.WriteString(style.SectionLabel(i18n.T(i18n.KeyShieldClosingTriggers)))
        b.WriteString(style.SectionLabel(i18n.T(i18n.KeyShieldConfigPerNiche)))
        ctItems := []kvItem{
                {"deal", i18n.T(i18n.KeyShieldAutoMarkDeal)},
                {"hot_lead", i18n.T(i18n.KeyShieldAutoMarkHot)},
                {"stop", i18n.T(i18n.KeyShieldAutoBlockStop)},
                {"override", i18n.T(i18n.KeyShieldManualOverride)},
        }
        for _, item := range ctItems {
                b.WriteString(style.Indent(1))
                b.WriteString(lipgloss.NewStyle().Foreground(style.TextMuted).Width(14).Render(item.Label))
                b.WriteString(lipgloss.NewStyle().Foreground(style.TextDim).Render(item.Value))
                b.WriteString("\n")
        }

        // Config paths
        b.WriteString(style.Section(style.SectionGap))
        b.WriteString(style.Indent(1))
        b.WriteString(lipgloss.NewStyle().Foreground(style.TextMuted).Width(18).Render(i18n.T(i18n.KeySettingsConfigMain)))
        b.WriteString(lipgloss.NewStyle().Foreground(style.Text).Render(i18n.T(i18n.KeySettingsConfigPath)))
        b.WriteString("\n")
        b.WriteString(style.Indent(1))
        b.WriteString(lipgloss.NewStyle().Foreground(style.TextDim).Render(i18n.T(i18n.KeyShieldConfigSection)))

        b.WriteString(style.Section(style.SectionGap))
        b.WriteString(lipgloss.NewStyle().Foreground(style.TextDim).Render(
                fmt.Sprintf("e  %s    r  %s    q  %s",
                        i18n.T(i18n.KeyLabelEdit), i18n.T(i18n.KeyLabelReload), i18n.T(i18n.KeyLabelBack),
                )))

        return b.String()
}
