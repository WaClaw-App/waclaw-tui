package overlay

import (
        "fmt"
        "strings"
        "time"

        "github.com/WaClaw-App/waclaw/internal/tui/anim"
        "github.com/WaClaw-App/waclaw/internal/tui/i18n"
        "github.com/WaClaw-App/waclaw/internal/tui/style"
        "github.com/WaClaw-App/waclaw/pkg/protocol"
        "github.com/charmbracelet/lipgloss"
)

// NotificationToast holds the state for the notification toast overlay.
//
// Spec (doc/14-notification-system.md, doc/10-global-overlays.md):
//   - Slide from top 250ms
//   - Auto-dismiss per severity: Critical=3s lock, Positive=10s, Neutral=5s, Informative=7s
//   - Max 1 at a time — never stack, never spam
//   - Critical cannot be dismissed for 3 seconds
//   - Queue: pending notifications wait for current to dismiss
type NotificationToast struct {
        // Active is the currently displayed notification, or nil.
        Active *NotificationData

        // Queue holds pending notifications waiting to be displayed.
        Queue []NotificationData

        // Width is the available terminal width.
        Width int

        // ShownAt tracks when the current notification was first displayed.
        ShownAt time.Time

        // Anim tracks the slide-in/out animation.
        Anim anim.AnimationState

        // IsFadingOut indicates the notification is in its dismiss animation.
        IsFadingOut bool
}

// NotificationData carries all information needed to render a single
// notification toast. It is produced from a bus.NotifyMsg.
type NotificationData struct {
        // Type is the notification category (e.g. "ResponseReceived").
        Type protocol.NotificationType

        // Severity determines auto-dismiss timing and visual style.
        Severity protocol.Severity

        // Title is the primary display text.
        Title string

        // Body is the secondary detail text (can be multi-line).
        Body string

        // Actions are the available key actions (e.g. "↵ balas    s nanti    q dismiss").
        Actions string

        // Data carries arbitrary key-value data from the backend.
        Data map[string]any
}

// NewNotificationToast creates a NotificationToast with default settings.
func NewNotificationToast() NotificationToast {
        return NotificationToast{
                Queue: make([]NotificationData, 0),
        }
}

// Enqueue adds a notification to the display queue.
// If no notification is active, it becomes the active one immediately.
func (nt *NotificationToast) Enqueue(data NotificationData) {
        if nt.Active == nil && !nt.IsFadingOut {
                nt.show(data)
        } else {
                nt.Queue = append(nt.Queue, data)
        }
}

// Dismiss removes the current notification and shows the next from the queue.
func (nt *NotificationToast) Dismiss() {
        if nt.Active == nil {
                return
        }

        nt.IsFadingOut = true
        nt.Anim = anim.NewAnimationState(anim.AnimFade, anim.NotifSlideOut)

        // Immediately dequeue next if animation would complete.
        // The Tick method handles the actual transition.
}

// ForceDismiss dismisses regardless of critical lock (e.g. escape key).
func (nt *NotificationToast) ForceDismiss() {
        nt.Dismiss()
}

// ShouldAutoDismiss returns true if the auto-dismiss duration has elapsed.
func (nt NotificationToast) ShouldAutoDismiss() bool {
        if nt.Active == nil {
                return false
        }
        return time.Since(nt.ShownAt) >= AutoDismissForType(nt.Active.Type, nt.Active.Severity)
}

// IsCriticalLocked returns true for critical notifications within their 3s hold period.
func (nt NotificationToast) IsCriticalLocked() bool {
        if nt.Active == nil {
                return false
        }
        return nt.Active.Severity == protocol.SeverityCritical &&
                time.Since(nt.ShownAt) < anim.NotifSeverityCritical
}

// IsVisible returns true if a notification is currently displayed.
func (nt NotificationToast) IsVisible() bool {
        return nt.Active != nil || nt.IsFadingOut
}

// Tick advances animation state and handles auto-dismissal.
func (nt *NotificationToast) Tick() {
        if nt.Active == nil {
                return
        }

        nt.Anim.UpdateProgress()

        // Handle fade-out completion.
        if nt.IsFadingOut && nt.Anim.IsComplete() {
                nt.Active = nil
                nt.IsFadingOut = false
                nt.showNext()
                return
        }

        // Auto-dismiss check.
        if !nt.IsFadingOut && nt.ShouldAutoDismiss() {
                nt.Dismiss()
        }
}

// show makes a notification the active display.
func (nt *NotificationToast) show(data NotificationData) {
        nt.Active = &data
        nt.ShownAt = time.Now()
        nt.IsFadingOut = false
        nt.Anim = anim.NewAnimationState(anim.AnimSlide, anim.NotifSlideIn)
}

// showNext dequeues the next pending notification, if any.
func (nt *NotificationToast) showNext() {
        if len(nt.Queue) == 0 {
                return
        }
        next := nt.Queue[0]
        nt.Queue = nt.Queue[1:]
        nt.show(next)
}

// View renders the notification toast overlay.
func (nt NotificationToast) View() string {
        if nt.Active == nil && !nt.IsFadingOut {
                return ""
        }

        if nt.Active == nil {
                return ""
        }

        n := nt.Active
        width := nt.Width
        if width < 40 {
                width = 40
        }

        // Determine severity color and label.
        sevColor, sevLabel := severityStyle(n.Severity)

        var lines []string

        // Title line.
        titleStyle := lipgloss.NewStyle().Foreground(sevColor).Bold(true)
        lines = append(lines, titleStyle.Render(n.Title))

        // Body lines.
        if n.Body != "" {
                bodyStyle := lipgloss.NewStyle().Foreground(style.TextMuted)
                for _, line := range strings.Split(n.Body, "\n") {
                        lines = append(lines, bodyStyle.Render(line))
                }
        }

        // Actions line.
        if n.Actions != "" {
                actionStyle := lipgloss.NewStyle().Foreground(style.TextDim)
                lines = append(lines, actionStyle.Render(n.Actions))
        } else {
                // Default actions: ↵ view    s ok
                lines = append(lines, lipgloss.NewStyle().Foreground(style.TextDim).
                        Render(fmt.Sprintf("↵ %s    s %s", i18n.T("notif.action_view"), i18n.T("notif.action_ok"))))
        }

        // Severity label (right-aligned indicator).
        sevTagStyle := lipgloss.NewStyle().Foreground(sevColor)
        sevTag := sevTagStyle.Render(sevLabel)

        // Compose into a panel.
        content := strings.Join(lines, "\n")
        panel := lipgloss.NewStyle().
                Background(style.BgRaised).
                Width(width).
                Padding(0, 2).
                Render(content + "  " + sevTag)

        return panel
}

// severityStyleMap maps each severity to its color and i18n label key.
var severityStyleMap = map[protocol.Severity]struct {
        color lipgloss.Color
        label string
}{
        protocol.SeverityCritical:    {style.Danger, "notif.critical"},
        protocol.SeverityPositive:    {style.Success, "notif.positive"},
        protocol.SeverityNeutral:     {style.TextMuted, "notif.neutral"},
        protocol.SeverityInformative: {style.Accent, "notif.informative"},
}

// severityStyle returns the color and label for a notification severity.
func severityStyle(sev protocol.Severity) (lipgloss.Color, string) {
        if entry, ok := severityStyleMap[sev]; ok {
                return entry.color, i18n.T(entry.label)
        }
        return style.TextMuted, ""
}

// NotificationDataFromMsg converts a bus.NotifyMsg into a NotificationData.
// This is the canonical mapping between the bus message and the overlay data.
func NotificationDataFromMsg(notifType protocol.NotificationType, sev protocol.Severity, data map[string]any) NotificationData {
        title := notificationTitle(notifType, data)
        body := notificationBody(notifType, data)
        actions := notificationActions(notifType, sev)

        return NotificationData{
                Type:     notifType,
                Severity: sev,
                Title:    title,
                Body:     body,
                Actions:  actions,
                Data:     data,
        }
}

// notifTitleFormatMap maps each notification type to its i18n key for the title.
// The i18n strings now contain format placeholders like {name}, {n}, {s}, {v}, {device}.
var notifTitleFormatMap = map[protocol.NotificationType]string{
        protocol.NotifResponseReceived:  "notif.response_received",
        protocol.NotifMultiResponse:     "notif.multi_response",
        protocol.NotifScrapeComplete:    "notif.scrape_complete",
        protocol.NotifBatchSendComplete: "notif.batch_send_complete",
        protocol.NotifWADisconnect:      "notif.wa_disconnect",
        protocol.NotifWAFlag:            "notif.wa_flag",
        protocol.NotifHealthScoreDrop:   "notif.health_score_drop",
        protocol.NotifDailyLimit:        "notif.daily_limit",
        protocol.NotifStreakMilestone:   "notif.streak_milestone",
        protocol.NotifConfigError:       "notif.config_error",
        protocol.NotifValidationError:   "notif.validation_error",
        protocol.NotifLicenseExpired:    "notif.license_expired",
        protocol.NotifDeviceConflict:    "notif.device_conflict",
        protocol.NotifFollowUpScheduled: "notif.follow_up_scheduled",
        protocol.NotifLeadCold:          "notif.lead_cold",
        protocol.NotifUpdateAvailable:   "notif.update_available",
        protocol.NotifUpgradeAvailable:  "notif.upgrade_available",
}

// formatNotifTitle replaces placeholders in a template string with data values.
// This avoids per-type switch logic — the format string drives everything.
func formatNotifTitle(template string, data map[string]any) string {
        result := template
        // Replace {name}
        if v, ok := data["name"].(string); ok {
                result = strings.ReplaceAll(result, "{name}", v)
        }
        // Replace {n} - integer count
        if v, ok := data["count"].(float64); ok {
                result = strings.ReplaceAll(result, "{n}", fmt.Sprintf("%d", int(v)))
        }
        if v, ok := data["n"].(float64); ok {
                result = strings.ReplaceAll(result, "{n}", fmt.Sprintf("%d", int(v)))
        }
        // Replace {s} - score
        if v, ok := data["score"].(float64); ok {
                result = strings.ReplaceAll(result, "{s}", fmt.Sprintf("%d", int(v)))
        }
        if v, ok := data["s"].(float64); ok {
                result = strings.ReplaceAll(result, "{s}", fmt.Sprintf("%d", int(v)))
        }
        // Replace {max} - max count
        if v, ok := data["max"].(float64); ok {
                result = strings.ReplaceAll(result, "{max}", fmt.Sprintf("%d", int(v)))
        }
        // Replace {v} - version string
        if v, ok := data["version"].(string); ok {
                result = strings.ReplaceAll(result, "{v}", v)
        }
        if v, ok := data["v"].(string); ok {
                result = strings.ReplaceAll(result, "{v}", v)
        }
        // Replace {device} - device name
        if v, ok := data["device"].(string); ok {
                result = strings.ReplaceAll(result, "{device}", v)
        }
        // Replace {niches}
        if v, ok := data["niches"].(string); ok {
                result = strings.ReplaceAll(result, "{niches}", v)
        }
        // Replace {breakdown}
        if v, ok := data["breakdown"].(string); ok {
                result = strings.ReplaceAll(result, "{breakdown}", v)
        }
        // Replace {names}
        if v, ok := data["names"].(string); ok {
                result = strings.ReplaceAll(result, "{names}", v)
        }
        // Replace {file}
        if v, ok := data["file"].(string); ok {
                result = strings.ReplaceAll(result, "{file}", v)
        }
        return result
}

// notificationTitle returns the primary display text for a notification type.
func notificationTitle(notifType protocol.NotificationType, data map[string]any) string {
        // Try data["message"] first (from backend).
        if msg, ok := data["message"].(string); ok {
                return msg
        }

        // Type-based formatting with data placeholders per doc/14.
        switch notifType {
        case protocol.NotifResponseReceived:
                if name, ok := data["name"].(string); ok {
                        return fmt.Sprintf("💬 %s %s", name, i18n.T("notif.response_received_suffix"))
                }
                return i18n.T("notif.response_received")
        case protocol.NotifMultiResponse:
                if n, ok := data["count"].(float64); ok {
                        return fmt.Sprintf(i18n.T("notif.multi_response"), int(n))
                }
                return i18n.T("notif.response_received") // fallback
        case protocol.NotifScrapeComplete:
                if n, ok := data["count"].(float64); ok {
                        return fmt.Sprintf(i18n.T("notif.scrape_complete"), int(n))
                }
                return i18n.T("notif.scrape_complete")
        case protocol.NotifBatchSendComplete:
                if n, ok := data["count"].(float64); ok {
                        return fmt.Sprintf(i18n.T("notif.batch_send_complete"), int(n))
                }
                return i18n.T("notif.batch_send_complete")
        case protocol.NotifWADisconnect:
                if n, ok := data["slot"].(float64); ok {
                        return fmt.Sprintf(i18n.T("notif.wa_disconnect"), int(n))
                }
                return i18n.T("notif.wa_disconnect")
        case protocol.NotifWAFlag:
                if n, ok := data["number"].(string); ok {
                        return fmt.Sprintf(i18n.T("notif.wa_flag"), n)
                }
                return i18n.T("notif.wa_flag")
        case protocol.NotifHealthScoreDrop:
                if slot, ok := data["slot"].(float64); ok {
                        if score, ok2 := data["score"].(float64); ok2 {
                                return fmt.Sprintf(i18n.T("notif.health_score_drop"), int(slot), int(score))
                        }
                }
                return i18n.T("notif.health_score_drop")
        case protocol.NotifDailyLimit:
                if n, ok := data["count"].(float64); ok {
                        if max, ok2 := data["max"].(float64); ok2 {
                                return fmt.Sprintf(i18n.T("notif.daily_limit"), int(n), int(max))
                        }
                }
                return i18n.T("notif.daily_limit")
        case protocol.NotifStreakMilestone:
                if n, ok := data["count"].(float64); ok {
                        return fmt.Sprintf(i18n.T("notif.streak_milestone"), int(n))
                }
                return i18n.T("notif.streak_milestone")
        case protocol.NotifConfigError:
                if n, ok := data["count"].(float64); ok {
                        return fmt.Sprintf(i18n.T("notif.config_error"), int(n))
                }
                return i18n.T("notif.config_error")
        case protocol.NotifValidationError:
                if n, ok := data["count"].(float64); ok {
                        return fmt.Sprintf(i18n.T("notif.validation_error"), int(n))
                }
                return i18n.T("notif.validation_error")
        case protocol.NotifLicenseExpired:
                return i18n.T("notif.license_expired")
        case protocol.NotifDeviceConflict:
                if device, ok := data["device"].(string); ok {
                        return fmt.Sprintf(i18n.T("notif.device_conflict"), device)
                }
                return i18n.T("notif.device_conflict")
        case protocol.NotifFollowUpScheduled:
                if n, ok := data["count"].(float64); ok {
                        return fmt.Sprintf(i18n.T("notif.follow_up_scheduled"), int(n))
                }
                return i18n.T("notif.follow_up_scheduled")
        case protocol.NotifLeadCold:
                if n, ok := data["count"].(float64); ok {
                        return fmt.Sprintf(i18n.T("notif.lead_cold"), int(n))
                }
                return i18n.T("notif.lead_cold")
        case protocol.NotifUpdateAvailable:
                if v, ok := data["version"].(string); ok {
                        return fmt.Sprintf(i18n.T("notif.update_available"), v)
                }
                return i18n.T("notif.update_available")
        case protocol.NotifUpgradeAvailable:
                if v, ok := data["version"].(string); ok {
                        return fmt.Sprintf(i18n.T("notif.upgrade_available"), v)
                }
                return i18n.T("notif.upgrade_available")
        default:
                return string(notifType)
        }
}

// notifBodyMap maps notification types to their body i18n key.
var notifBodyMap = map[protocol.NotificationType]string{
        protocol.NotifScrapeComplete:    "notif.body_scrape_complete",
        protocol.NotifBatchSendComplete: "notif.body_batch_send_complete",
        protocol.NotifConfigError:       "notif.body_config_error",
        protocol.NotifDeviceConflict:    "notif.body_device_conflict",
        protocol.NotifWADisconnect:      "notif.body_wa_disconnect",
        protocol.NotifWAFlag:            "notif.body_wa_flag",
        protocol.NotifMultiResponse:     "notif.body_multi_response",
        protocol.NotifLicenseExpired:    "notif.body_license_expired",
        protocol.NotifHealthScoreDrop:   "notif.body_health_score_drop",
}

// notificationBody returns the secondary detail text.
func notificationBody(notifType protocol.NotificationType, data map[string]any) string {
        if detail, ok := data["detail"].(string); ok {
                return detail
        }
        key, ok := notifBodyMap[notifType]
        if !ok {
                return ""
        }
        template := i18n.T(key)
        return formatNotifTitle(template, data) // reuse the same formatter
}

// notifActionMap maps each notification type to its per-type action i18n keys,
// matching the doc exactly for each notification type.
var notifActionMap = map[protocol.NotificationType][]string{
        protocol.NotifResponseReceived:  {"notif.action_reply", "notif.action_later_short", "notif.action_dismiss_short"},
        protocol.NotifMultiResponse:     {"notif.action_process", "notif.action_1_auto_offer", "notif.action_later_short"},
        protocol.NotifScrapeComplete:    {"notif.action_view", "notif.action_ok"},
        protocol.NotifBatchSendComplete: {"notif.action_view", "notif.action_ok"},
        protocol.NotifWADisconnect:      {"notif.action_1_relogin", "notif.action_dismiss"},
        protocol.NotifWAFlag:            {"notif.action_view_shield", "notif.action_let_waclaw"},
        protocol.NotifHealthScoreDrop:   {"notif.action_view_shield", "notif.action_ok"},
        protocol.NotifDailyLimit:        {"notif.action_ok_proceed"},
        protocol.NotifStreakMilestone:   {"notif.action_view_stats", "notif.action_nice"},
        protocol.NotifConfigError:       {"notif.action_view_error", "notif.action_validate_all", "notif.action_later"},
        protocol.NotifValidationError:   {"notif.action_view_error", "notif.action_validate_all", "notif.action_later"},
        protocol.NotifLicenseExpired:    {"notif.action_enter_license", "notif.action_exit"},
        protocol.NotifDeviceConflict:    {"notif.action_view_license", "notif.action_disconnect_other", "notif.action_exit"},
        protocol.NotifFollowUpScheduled: {"notif.action_view", "notif.action_ok"},
        protocol.NotifLeadCold:          {"notif.action_view", "notif.action_ok"},
        protocol.NotifUpdateAvailable:   {"notif.action_update_now", "notif.action_u_later", "notif.action_skip"},
        protocol.NotifUpgradeAvailable:  {"notif.action_view_info", "notif.action_u_upgrade", "notif.action_later"},
}

// notificationActions returns the action key hints for a notification type.
func notificationActions(notifType protocol.NotificationType, sev protocol.Severity) string {
        switch notifType {
        case protocol.NotifResponseReceived:
                return fmt.Sprintf("↵ %s    s %s", i18n.T("notif.action_reply"), i18n.T("notif.action_later"))
        case protocol.NotifMultiResponse:
                return fmt.Sprintf("↵ %s    1 %s    s %s", i18n.T("notif.action_process"), i18n.T("notif.action_auto_offer"), i18n.T("notif.action_later"))
        case protocol.NotifScrapeComplete, protocol.NotifBatchSendComplete:
                return fmt.Sprintf("↵ %s    s %s", i18n.T("notif.action_view"), i18n.T("notif.action_ok"))
        case protocol.NotifWADisconnect:
                return fmt.Sprintf("1 %s    s %s", i18n.T("notif.action_relogin"), i18n.T("notif.action_dismiss"))
        case protocol.NotifWAFlag:
                return fmt.Sprintf("↵ %s    s %s", i18n.T("notif.action_view_shield"), i18n.T("notif.action_let_waclaw"))
        case protocol.NotifHealthScoreDrop, protocol.NotifDailyLimit:
                return fmt.Sprintf("↵ %s    s %s", i18n.T("notif.action_view_shield"), i18n.T("notif.action_ok"))
        case protocol.NotifStreakMilestone:
                return fmt.Sprintf("↵ %s    s %s", i18n.T("notif.action_view"), i18n.T("notif.action_ok"))
        case protocol.NotifConfigError, protocol.NotifValidationError:
                return fmt.Sprintf("↵ %s    v %s    s %s", i18n.T("notif.action_view_error"), i18n.T("notif.action_validate_all"), i18n.T("notif.action_later"))
        case protocol.NotifLicenseExpired:
                return fmt.Sprintf("↵ %s    q %s", i18n.T("notif.action_enter_license"), i18n.T("notif.action_exit"))
        case protocol.NotifDeviceConflict:
                return fmt.Sprintf("↵ %s    2 %s    s %s", i18n.T("notif.action_view_license"), i18n.T("notif.action_disconnect_other"), i18n.T("notif.action_exit"))
        case protocol.NotifFollowUpScheduled, protocol.NotifLeadCold:
                return fmt.Sprintf("↵ %s    s %s", i18n.T("notif.action_view"), i18n.T("notif.action_ok"))
        case protocol.NotifUpdateAvailable:
                return fmt.Sprintf("↵ %s    u %s    s %s", i18n.T("notif.action_update"), i18n.T("notif.action_later"), i18n.T("notif.action_skip"))
        case protocol.NotifUpgradeAvailable:
                return fmt.Sprintf("↵ %s    u %s    s %s", i18n.T("notif.action_view_info"), i18n.T("notif.action_upgrade"), i18n.T("notif.action_later"))
        default:
                return fmt.Sprintf("↵ %s    s %s", i18n.T("notif.action_view"), i18n.T("notif.action_ok"))
        }
}
