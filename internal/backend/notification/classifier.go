// Package notification implements the notification dispatch system for the
// WaClaw backend. It owns the severity classification logic, the per-type
// display templates, and the queue that enforces "max 1 at a time" so the
// TUI never gets spammed.
//
// DRY convention: notification type constants live in pkg/protocol/types.go
// (NotifResponseReceived, NotifWAFlag, etc.). This package maps each type
// to its severity, display template, and auto-dismiss timing. The only
// source of truth for the type list is protocol.AllNotificationTypes().
package notification

import (
        "github.com/WaClaw-App/waclaw/pkg/protocol"
)

// severityMap is the data-driven mapping from NotificationType to Severity.
// Single source of truth per doc/14-notification-system.md.
// Any unknown type defaults to SeverityNeutral (safe fallback).
var severityMap = map[protocol.NotificationType]protocol.Severity{
        // Positive
        protocol.NotifResponseReceived: protocol.SeverityPositive,
        protocol.NotifMultiResponse:    protocol.SeverityPositive,
        protocol.NotifStreakMilestone:  protocol.SeverityPositive,
        protocol.NotifUpdateAvailable:  protocol.SeverityPositive,

        // Neutral
        protocol.NotifScrapeComplete:    protocol.SeverityNeutral,
        protocol.NotifBatchSendComplete: protocol.SeverityNeutral,
        protocol.NotifHealthScoreDrop:   protocol.SeverityNeutral,
        protocol.NotifDailyLimit:        protocol.SeverityNeutral,

        // Critical
        protocol.NotifWADisconnect:    protocol.SeverityCritical,
        protocol.NotifWAFlag:          protocol.SeverityCritical,
        protocol.NotifConfigError:     protocol.SeverityCritical,
        protocol.NotifValidationError: protocol.SeverityCritical,
        protocol.NotifLicenseExpired:  protocol.SeverityCritical,
        protocol.NotifDeviceConflict:  protocol.SeverityCritical,

        // Informative
        protocol.NotifFollowUpScheduled: protocol.SeverityInformative,
        protocol.NotifLeadCold:          protocol.SeverityInformative,
        protocol.NotifUpgradeAvailable:  protocol.SeverityInformative,
}

// SeverityFor returns the canonical severity level for a given notification
// type. The mapping is data-driven from severityMap — the single source of
// truth per doc/14-notification-system.md.
//
// Any unknown type defaults to SeverityNeutral (safe fallback).
func SeverityFor(t protocol.NotificationType) protocol.Severity {
        if s, ok := severityMap[t]; ok {
                return s
        }
        return protocol.SeverityNeutral
}

// notifTemplateMap is the data-driven mapping from NotificationType to its
// i18n display key. The overlay rendering layer uses i18n.T() to resolve
// the actual string. Each key corresponds to an entry in keys.go / en.go / id.go.
var notifTemplateMap = map[protocol.NotificationType]string{
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

// NotifTemplate returns the i18n key for a notification type's display text.
// The overlay rendering layer uses i18n.T() to resolve the actual string.
// Each key corresponds to an entry in keys.go / en.go / id.go.
func NotifTemplate(t protocol.NotificationType) string {
        if key, ok := notifTemplateMap[t]; ok {
                return key
        }
        return "notif.neutral"
}
