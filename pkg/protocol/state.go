package protocol

// StateID identifies a specific state within a screen's state machine.
//
// StateID is a named type (not a type alias) so that the compiler can
// catch accidental mixing with plain strings. It serialises directly in
// JSON and is interchangeable via explicit conversion where needed.
//
// Only screen-level states belong here. Domain-specific concepts that
// have their own vocabulary (severity, notification types, worker phases,
// lead lifecycle, confirmation types) live in types.go under their own
// named types — that is the DRY convention.
type StateID string

// ===========================================================================
// Boot states
// ===========================================================================

const (
        // BootFirstTime is shown when the application starts for the very first
        // time (no existing configuration found).
        BootFirstTime StateID = "boot_first_time"

        // BootReturning is the normal startup state for returning users.
        BootReturning StateID = "boot_returning"

        // BootReturningResponse indicates a response was received during boot.
        BootReturningResponse StateID = "boot_returning_response"

        // BootReturningError is shown when a non-fatal error occurs during
        // the returning-user boot sequence.
        BootReturningError StateID = "boot_returning_error"

        // BootReturningConfigError indicates a configuration error was detected
        // at boot.
        BootReturningConfigError StateID = "boot_returning_config_error"

        // BootReturningLicenseExpired is shown at boot when the user's license
        // has expired.
        BootReturningLicenseExpired StateID = "boot_returning_license_expired"

        // BootReturningDeviceConflict is shown at boot when the same license is
        // already active on another device.
        BootReturningDeviceConflict StateID = "boot_returning_device_conflict"
)

// ===========================================================================
// Login states
// ===========================================================================

const (
        // LoginQRWaiting shows the QR code and waits for the user to scan it.
        LoginQRWaiting StateID = "login_qr_waiting"

        // LoginQRScanned indicates the QR code has been scanned; awaiting
        // authentication confirmation.
        LoginQRScanned StateID = "login_qr_scanned"

        // LoginSuccess indicates a successful WhatsApp login.
        LoginSuccess StateID = "login_success"

        // LoginExpired indicates the session has expired and re-login is needed.
        LoginExpired StateID = "login_expired"

        // LoginFailed indicates that the login attempt failed.
        LoginFailed StateID = "login_failed"
)

// ===========================================================================
// Niche states
// ===========================================================================

const (
        // NicheList shows the list of available niches for selection.
        NicheList StateID = "niche_list"

        // NicheMultiSelected indicates one or more niches have been selected.
        NicheMultiSelected StateID = "niche_multi_selected"

        // NicheCustom is the state for entering a custom niche.
        NicheCustom StateID = "niche_custom"

        // NicheEditFilters allows editing filters for a niche.
        NicheEditFilters StateID = "niche_edit_filters"

        // NicheConfigError indicates a configuration error with the selected
        // niche.
        NicheConfigError StateID = "niche_config_error"
)

// ===========================================================================
// Scrape states
// ===========================================================================

const (
        // ScrapeActive indicates scraping is actively running.
        ScrapeActive StateID = "scrape_active"

        // ScrapeMultiActive indicates multiple niche scrapers are running
        // concurrently.
        ScrapeMultiActive StateID = "scrape_multi_active"

        // ScrapeMultiStaggered indicates multi-niche scraping with staggered
        // start times.
        ScrapeMultiStaggered StateID = "scrape_multi_staggered"

        // ScrapeIdle indicates scraping is paused or waiting.
        ScrapeIdle StateID = "scrape_idle"

        // ScrapeEmpty indicates no leads were found.
        ScrapeEmpty StateID = "scrape_empty"

        // ScrapeError indicates an error occurred during scraping.
        ScrapeError StateID = "scrape_error"

        // ScrapeGMapsLimited indicates Google Maps scraping hit a limit.
        ScrapeGMapsLimited StateID = "scrape_gmaps_limited"

        // ScrapeAutoApproved indicates scraped leads were auto-approved.
        ScrapeAutoApproved StateID = "scrape_auto_approved"

        // ScrapeHighValueReveal reveals high-value leads after scraping.
        ScrapeHighValueReveal StateID = "scrape_high_value_reveal"

        // ScrapeBatchComplete indicates a scrape batch has finished.
        ScrapeBatchComplete StateID = "scrape_batch_complete"


        // ScrapeWAValidation indicates WA pre-validation is running.
        ScrapeWAValidation StateID = "scrape_wa_validation"

        // ScrapeWAValidationProgress shows detailed WA validation progress.
        ScrapeWAValidationProgress StateID = "scrape_wa_validation_progress"
)

// ===========================================================================
// Review states
// ===========================================================================

const (
        // ReviewReviewing is the main state for reviewing leads one-by-one.
        ReviewReviewing StateID = "review_reviewing"

        // ReviewLeadDetail shows the full detail view for a single lead.
        ReviewLeadDetail StateID = "review_lead_detail"

        // ReviewTemplatePreview shows a template preview for the current lead.
        ReviewTemplatePreview StateID = "review_template_preview"

        // ReviewQueueComplete indicates all leads in the review queue have been
        // processed.
        ReviewQueueComplete StateID = "review_queue_complete"
)

// ===========================================================================
// Send states
// ===========================================================================

const (
        // SendActive indicates outbound messages are being sent.
        SendActive StateID = "send_active"

        // SendPaused indicates sending has been paused by the user.
        SendPaused StateID = "send_paused"

        // SendOffHours indicates sending is paused because it is outside allowed
        // hours.
        SendOffHours StateID = "send_off_hours"

        // SendRateLimited indicates WhatsApp rate-limiting is in effect.
        SendRateLimited StateID = "send_rate_limited"

        // SendDailyLimit indicates the daily send limit has been reached.
        SendDailyLimit StateID = "send_daily_limit"

        // SendFailed indicates a send failure.
        SendFailed StateID = "send_failed"

        // SendAllSlotsDown indicates all WhatsApp slots are disconnected.
        SendAllSlotsDown StateID = "send_all_slots_down"

        // SendWithResponse indicates sending is active and some leads have
        // responded.
        SendWithResponse StateID = "send_with_response"
)

// ===========================================================================
// Monitor states
// ===========================================================================

const (
        // MonitorLiveDashboard is the real-time activity dashboard.
        MonitorLiveDashboard StateID = "monitor_live_dashboard"

        // MonitorIdleBackground indicates the monitor is running in background.
        MonitorIdleBackground StateID = "monitor_idle_background"

        // MonitorNight indicates night-mode / reduced activity monitoring.
        MonitorNight StateID = "monitor_night"

        // MonitorError indicates a monitoring error.
        MonitorError StateID = "monitor_error"

        // MonitorEmpty indicates no data is available yet.
        MonitorEmpty StateID = "monitor_empty"

        // MonitorPendingResponses indicates there are unread responses.
        MonitorPendingResponses StateID = "monitor_pending_responses"
)

// ===========================================================================
// Response states
// ===========================================================================

const (
        // ResponsePositive indicates a positive lead response.
        ResponsePositive StateID = "response_positive"

        // ResponseCurious indicates a curious / interested lead response.
        ResponseCurious StateID = "response_curious"

        // ResponseNegative indicates a negative lead response.
        ResponseNegative StateID = "response_negative"

        // ResponseMaybe indicates an ambiguous lead response.
        ResponseMaybe StateID = "response_maybe"

        // ResponseAutoReply indicates an automated reply was detected.
        ResponseAutoReply StateID = "response_auto_reply"

        // ResponseOfferPreview shows an offer preview for a lead.
        ResponseOfferPreview StateID = "response_offer_preview"

        // ResponseMultiQueue indicates multiple responses are queued.
        ResponseMultiQueue StateID = "response_multi_queue"

        // ResponseConversion indicates a lead has converted.
        ResponseConversion StateID = "response_conversion"

        // ResponseHotLead indicates a hot lead has been detected.
        ResponseHotLead StateID = "response_hot_lead"

        // ResponseStopDetected indicates a stop signal was detected in the response.
        ResponseStopDetected StateID = "response_stop_detected"

        // ResponseDealDetected indicates a deal has been detected in the response.
        ResponseDealDetected StateID = "response_deal_detected"
)

// ===========================================================================
// Leads screen states
// ===========================================================================

const (
        // LeadsList is the default lead list view.
        LeadsList StateID = "leads_list"

        // LeadsFiltered shows a filtered subset of leads.
        LeadsFiltered StateID = "leads_filtered"

        // LeadsFullDetail shows the complete detail view for a single lead.
        LeadsFullDetail StateID = "leads_full_detail"

        // LeadsFollowUpDue shows leads that have a follow-up due.
        LeadsFollowUpDue StateID = "leads_follow_up_due"

        // LeadsCold shows leads that have gone cold.
        LeadsCold StateID = "leads_cold"

        // LeadsNeverContacted shows leads that have never been contacted.
        LeadsNeverContacted StateID = "leads_never_contacted"

        // LeadsConverted shows leads that have converted.
        LeadsConverted StateID = "leads_converted"
)

// ===========================================================================
// Template screen states
// ===========================================================================

const (
        // TemplateList is the default template list view.
        TemplateList StateID = "template_list"

        // TemplatePreview shows a rendered preview of a template.
        TemplatePreview StateID = "template_preview"

        // TemplateEditHint is the state for editing a template hint / variable.
        TemplateEditHint StateID = "template_edit_hint"

        // TemplateValidationError indicates a template validation error.
        TemplateValidationError StateID = "template_validation_error"
)

// ===========================================================================
// Workers screen states
// ===========================================================================

const (
        // WorkersOverview shows the status of all workers.
        WorkersOverview StateID = "workers_overview"

        // WorkerDetail shows detailed info for a single worker.
        WorkerDetail StateID = "worker_detail"

        // WorkerAddNiche is the state for adding a niche to a worker.
        WorkerAddNiche StateID = "worker_add_niche"

        // WorkersPaused is shown when a worker has been manually paused.
        WorkersPaused StateID = "workers_paused"
)

// ===========================================================================
// Shield (Anti-Ban) screen states
// ===========================================================================

const (
        // ShieldOverview shows the overall shield / health-score status.
        ShieldOverview StateID = "shield_overview"

        // ShieldWarning indicates a warning-level health score.
        ShieldWarning StateID = "shield_warning"

        // ShieldDanger indicates a danger-level health score.
        ShieldDanger StateID = "shield_danger"

        // ShieldSlotDetail shows detail for a single WhatsApp slot.
        ShieldSlotDetail StateID = "shield_slot_detail"

        // ShieldSettings is the shield configuration view.
        ShieldSettings StateID = "shield_settings"
)

// ===========================================================================
// Settings screen states
// ===========================================================================

const (
        // SettingsOverview is the main settings overview.
        SettingsOverview StateID = "settings_overview"

        // SettingsEdit is the state for editing a settings section.
        SettingsEdit StateID = "settings_edit"

        // SettingsReload indicates settings are being reloaded.
        SettingsReload StateID = "settings_reload"

        // SettingsReloadError indicates an error reloading settings.
        SettingsReloadError StateID = "settings_reload_error"
)

// ===========================================================================
// Validation screen states
// ===========================================================================

const (
        // ValidationClean indicates no validation issues.
        ValidationClean StateID = "validation_clean"

        // ValidationErrors indicates one or more validation errors.
        ValidationErrors StateID = "validation_errors"

        // ValidationWarnings indicates one or more validation warnings.
        ValidationWarnings StateID = "validation_warnings"

        // ValidationFix indicates a fix suggestion is available.
        ValidationFix StateID = "validation_fix"

        // ValidationFirstTime shows validation guidance for first-time users.
        ValidationFirstTime StateID = "validation_first_time"

        // ValidationReloadError indicates an error during config reload validation.
        ValidationReloadError StateID = "validation_reload_error"
)

// ===========================================================================
// Compose screen states
// ===========================================================================

const (
        // ComposeDraft is the draft editing state.
        ComposeDraft StateID = "compose_draft"

        // ComposePreview shows a rendered preview of the composed message.
        ComposePreview StateID = "compose_preview"

        // ComposeTemplatePick allows picking a template to start composing.
        ComposeTemplatePick StateID = "compose_template_pick"
)

// ===========================================================================
// History screen states
// ===========================================================================

const (
        // HistoryToday shows today's activity history.
        HistoryToday StateID = "history_today"

        // HistoryWeek shows the weekly activity summary.
        HistoryWeek StateID = "history_week"

        // HistoryDayDetail shows detail for a specific day.
        HistoryDayDetail StateID = "history_day_detail"
)

// ===========================================================================
// FollowUp screen states
// ===========================================================================

const (
        // FollowUpDashboard is the main follow-up overview.
        FollowUpDashboard StateID = "followup_dashboard"

        // FollowUpNicheDetail shows follow-ups for a specific niche.
        FollowUpNicheDetail StateID = "followup_niche_detail"

        // FollowUpSending indicates follow-up messages are being sent.
        FollowUpSending StateID = "followup_sending"

        // FollowUpEmpty indicates no follow-ups are scheduled.
        FollowUpEmpty StateID = "followup_empty"

        // FollowUpColdList shows the cold-lead follow-up list.
        FollowUpColdList StateID = "followup_cold_list"

        // FollowUpRecontact is the state for re-contacting cold leads.
        FollowUpRecontact StateID = "followup_recontact"
)

// ===========================================================================
// Explorer screen states
// ===========================================================================

const (
        // ExplorerBrowse is the default niche browsing state.
        ExplorerBrowse StateID = "explorer_browse"

        // ExplorerSearch allows searching for niches.
        ExplorerSearch StateID = "explorer_search"

        // ExplorerCategoryDetail shows detail for a niche category.
        ExplorerCategoryDetail StateID = "explorer_category_detail"

        // ExplorerGenerating indicates AI-generated niches are being computed.
        ExplorerGenerating StateID = "explorer_generating"

        // ExplorerGenerated shows the generated niche results.
        ExplorerGenerated StateID = "explorer_generated"
)

// ===========================================================================
// Update screen states
// ===========================================================================

const (
        // UpdateAvailable indicates a new version is available.
        UpdateAvailable StateID = "update_available"

        // UpdateDownloading indicates an update is being downloaded.
        UpdateDownloading StateID = "update_downloading"

        // UpdateReady indicates a downloaded update is ready to install.
        UpdateReady StateID = "update_ready"

        // UpgradeAvailable indicates a paid upgrade is available.
        UpgradeAvailable StateID = "upgrade_available"

        // UpgradeLicenseInput is the state for entering a license key for an
        // upgrade.
        UpgradeLicenseInput StateID = "upgrade_license_input"

        // LicenseExpiredWithUpgrade indicates the license has expired and an
        // upgrade is available.
        LicenseExpiredWithUpgrade StateID = "license_expired_with_upgrade"

        // StartupCheck is a background variant that runs at boot (t+250ms)
        // to check for updates without blocking the startup sequence.
        // Results are delivered as notifications after the dashboard is ready.
        StartupCheck StateID = "startup_check"
)

// ===========================================================================
// License screen states
// ===========================================================================

const (
        // LicenseInput is the state for entering a license key.
        LicenseInput StateID = "license_input"

        // LicenseValidating indicates the license key is being validated.
        LicenseValidating StateID = "license_validating"

        // LicenseValid indicates the license is valid and active.
        LicenseValid StateID = "license_valid"

        // LicenseInvalid indicates the license key is invalid.
        LicenseInvalid StateID = "license_invalid"

        // LicenseExpired indicates the license has expired.
        LicenseExpired StateID = "license_expired"

        // LicenseDeviceConflict indicates the license is in use on another
        // device.
        LicenseDeviceConflict StateID = "license_device_conflict"

        // LicenseServerError indicates a server error during license validation.
        LicenseServerError StateID = "license_server_error"
)

// ===========================================================================
// Overlay states (nerd stats, command palette)
// ===========================================================================

const (
        // NerdStatsHidden hides the nerd stats overlay.
        NerdStatsHidden StateID = "nerd_stats_hidden"

        // NerdStatsMinimal shows a minimal stats overlay.
        NerdStatsMinimal StateID = "nerd_stats_minimal"

        // NerdStatsExpanded shows the full expanded stats overlay.
        NerdStatsExpanded StateID = "nerd_stats_expanded"
)

const (
        // CmdPaletteClosed indicates the command palette is hidden.
        CmdPaletteClosed StateID = "cmd_palette_closed"

        // CmdPaletteOpen indicates the command palette is open.
        CmdPaletteOpen StateID = "cmd_palette_open"

        // CmdPaletteExecuting indicates a command is being executed.
        CmdPaletteExecuting StateID = "cmd_palette_executing"

        // CmdPaletteEmpty indicates the command palette has no results.
        CmdPaletteEmpty StateID = "cmd_palette_empty"

        // CmdPaletteWithRecent indicates the command palette shows recent
        // commands.
        CmdPaletteWithRecent StateID = "cmd_palette_with_recent"

        // CmdPaletteQuickAction indicates a quick-action is selected.
        CmdPaletteQuickAction StateID = "cmd_palette_quick_action"
)
