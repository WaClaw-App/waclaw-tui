package protocol

// ---------------------------------------------------------------------------
// Centralized type definitions — DRY principle
// ---------------------------------------------------------------------------
//
// This file owns all domain-specific types that are NOT screen states.
// Screen states belong in state.go. Types that represent distinct domains
// (severity levels, notification categories, worker phases, lead lifecycle)
// live here so they get proper Go types instead of sharing the StateID
// string alias, which would defeat type safety and blur semantic boundaries.
//
// Convention: every concept that has its own vocabulary gets its own named
// type.  Constants of that type are declared in the same const block.
// ---------------------------------------------------------------------------

// Severity represents the severity level of a notification.
//
// The TUI uses severity to determine toast auto-dismiss timing and
// visual styling (color, hold duration). Only four levels exist.
type Severity string

const (
        // SeverityCritical indicates a critical / urgent notification.
        // TUI: 3s lock, danger color, cannot be dismissed early.
        SeverityCritical Severity = "severity_critical"

        // SeverityPositive indicates a positive / success notification.
        // TUI: 10s auto-dismiss, success color.
        SeverityPositive Severity = "severity_positive"

        // SeverityNeutral indicates a neutral notification.
        // TUI: 5s auto-dismiss, text_muted color.
        SeverityNeutral Severity = "severity_neutral"

        // SeverityInformative indicates an informational notification.
        // TUI: 7s auto-dismiss, accent color.
        SeverityInformative Severity = "severity_informative"
)

// AllSeverities returns all defined severity levels.
func AllSeverities() []Severity {
        return []Severity{SeverityCritical, SeverityPositive, SeverityNeutral, SeverityInformative}
}

// IsValidSeverity reports whether s is a known severity level.
func IsValidSeverity(s Severity) bool {
        for _, v := range AllSeverities() {
                if v == s {
                        return true
                }
        }
        return false
}

// NotificationType identifies one of the 17 notification categories.
//
// Each type maps to a specific user-facing event in the system. The
// notification dispatcher queues, prioritises, and pushes these to the
// TUI via the bus.
type NotificationType string

const (
        // NotifResponseReceived is sent when a lead responds to a message.
        NotifResponseReceived NotificationType = "notif_response_received"

        // NotifMultiResponse is sent when multiple responses arrive at once.
        NotifMultiResponse NotificationType = "notif_multi_response"

        // NotifScrapeComplete is sent when a scraping job finishes.
        NotifScrapeComplete NotificationType = "notif_scrape_complete"

        // NotifBatchSendComplete is sent when a batch of messages has been sent.
        NotifBatchSendComplete NotificationType = "notif_batch_send_complete"

        // NotifWADisconnect is sent when a WhatsApp session disconnects.
        NotifWADisconnect NotificationType = "notif_wa_disconnect"

        // NotifWAFlag is sent when WhatsApp flags the account.
        NotifWAFlag NotificationType = "notif_wa_flag"

        // NotifHealthScoreDrop is sent when the health score drops significantly.
        NotifHealthScoreDrop NotificationType = "notif_health_score_drop"

        // NotifDailyLimit is sent when the daily send limit is reached.
        NotifDailyLimit NotificationType = "notif_daily_limit"

        // NotifStreakMilestone is sent when a sending streak milestone is hit.
        NotifStreakMilestone NotificationType = "notif_streak_milestone"

        // NotifConfigError is sent when a configuration error is detected.
        NotifConfigError NotificationType = "notif_config_error"

        // NotifValidationError is sent when a validation error occurs.
        NotifValidationError NotificationType = "notif_validation_error"

        // NotifLicenseExpired is sent when the license expires.
        NotifLicenseExpired NotificationType = "notif_license_expired"

        // NotifDeviceConflict is sent when a device conflict is detected.
        NotifDeviceConflict NotificationType = "notif_device_conflict"

        // NotifFollowUpScheduled is sent when a follow-up is scheduled.
        NotifFollowUpScheduled NotificationType = "notif_follow_up_scheduled"

        // NotifLeadCold is sent when a lead goes cold.
        NotifLeadCold NotificationType = "notif_lead_cold"

        // NotifUpdateAvailable is sent when an app update is available.
        NotifUpdateAvailable NotificationType = "notif_update_available"

        // NotifUpgradeAvailable is sent when a paid upgrade is available.
        NotifUpgradeAvailable NotificationType = "notif_upgrade_available"
)

// AllNotificationTypes returns all 17 notification types.
func AllNotificationTypes() []NotificationType {
        return []NotificationType{
                NotifResponseReceived, NotifMultiResponse, NotifScrapeComplete,
                NotifBatchSendComplete, NotifWADisconnect, NotifWAFlag,
                NotifHealthScoreDrop, NotifDailyLimit, NotifStreakMilestone,
                NotifConfigError, NotifValidationError, NotifLicenseExpired,
                NotifDeviceConflict, NotifFollowUpScheduled, NotifLeadCold,
                NotifUpdateAvailable, NotifUpgradeAvailable,
        }
}

// IsValidNotificationType reports whether n is a known notification type.
func IsValidNotificationType(n NotificationType) bool {
        for _, v := range AllNotificationTypes() {
                if v == n {
                        return true
                }
        }
        return false
}

// ConfirmationType identifies one of the 4 confirmation overlay types.
//
// Each type maps to a different destructive or high-impact action that
// requires explicit user confirmation before proceeding.
type ConfirmationType string

const (
        // ConfirmBulkOffer asks the user to confirm sending offers in bulk.
        ConfirmBulkOffer ConfirmationType = "confirm_bulk_offer"

        // ConfirmBulkDelete asks the user to confirm deleting leads in bulk.
        ConfirmBulkDelete ConfirmationType = "confirm_bulk_delete"

        // ConfirmBulkArchive asks the user to confirm archiving leads in bulk.
        ConfirmBulkArchive ConfirmationType = "confirm_bulk_archive"

        // ConfirmForceDisconnect asks the user to confirm forcing a WhatsApp disconnect.
        ConfirmForceDisconnect ConfirmationType = "confirm_force_disconnect"
)

// AllConfirmationTypes returns all 4 confirmation types.
func AllConfirmationTypes() []ConfirmationType {
        return []ConfirmationType{
                ConfirmBulkOffer, ConfirmBulkDelete, ConfirmBulkArchive, ConfirmForceDisconnect,
        }
}

// IsValidConfirmationType reports whether c is a known confirmation type.
func IsValidConfirmationType(c ConfirmationType) bool {
        for _, v := range AllConfirmationTypes() {
                if v == c {
                        return true
                }
        }
        return false
}

// WorkerPhase represents the current phase of a per-niche worker.
//
// WorkerPhase is distinct from the screen-level WorkersOverview/WorkerDetail
// states because it describes the internal lifecycle of a single niche worker
// (spawning → scraping → qualifying → queuing → sending → idle → etc.)
// rather than which view the Workers screen is showing.
type WorkerPhase string

const (
        // WorkerSpawning indicates a niche worker is being created.
        WorkerSpawning WorkerPhase = "worker_spawning"

        // WorkerScraping indicates a niche worker is actively scraping.
        WorkerScraping WorkerPhase = "worker_scraping"

        // WorkerQualifying indicates a niche worker is qualifying leads.
        WorkerQualifying WorkerPhase = "worker_qualifying"

        // WorkerQueuing indicates a niche worker is queueing messages.
        WorkerQueuing WorkerPhase = "worker_queuing"

        // WorkerSending indicates a niche worker is sending messages.
        WorkerSending WorkerPhase = "worker_sending"

        // WorkerIdle indicates a niche worker has no current task.
        WorkerIdle WorkerPhase = "worker_idle"

        // WorkerError indicates a niche worker encountered an error.
        WorkerError WorkerPhase = "worker_error"

        // WorkerRateLimited indicates a niche worker is rate-limited.
        WorkerRateLimited WorkerPhase = "worker_rate_limited"

        // WorkerPaused indicates a niche worker is paused.
        WorkerPaused WorkerPhase = "worker_paused"

        // WorkerConfigError indicates a niche worker has a configuration error.
        WorkerConfigError WorkerPhase = "worker_config_error"

        // WorkerStopped indicates a niche worker has been stopped.
        WorkerStopped WorkerPhase = "worker_stopped"
)

// AllWorkerPhases returns all defined worker phases.
func AllWorkerPhases() []WorkerPhase {
        return []WorkerPhase{
                WorkerSpawning, WorkerScraping, WorkerQualifying, WorkerQueuing,
                WorkerSending, WorkerIdle, WorkerError, WorkerRateLimited,
                WorkerPaused, WorkerConfigError, WorkerStopped,
        }
}

// IsValidWorkerPhase reports whether p is a known worker phase.
func IsValidWorkerPhase(p WorkerPhase) bool {
        for _, v := range AllWorkerPhases() {
                if v == p {
                        return true
                }
        }
        return false
}

// LeadPhase represents a step in the lead lifecycle state machine.
//
// LeadPhase is distinct from screen states because it describes the
// business-lifecycle status of a single lead (baru → ice_breaker_sent →
// responded → offer_sent → converted / cold / dead / etc.) rather than
// which view the Leads DB screen is showing.
type LeadPhase string

const (
        // LeadBaru is a brand-new lead with no action taken yet.
        LeadBaru LeadPhase = "lead_baru"

        // LeadIceBreakerSent indicates the ice-breaker message has been sent.
        LeadIceBreakerSent LeadPhase = "lead_ice_breaker_sent"

        // LeadResponded indicates the lead has responded.
        LeadResponded LeadPhase = "lead_responded"

        // LeadOfferSent indicates an offer has been sent to the lead.
        LeadOfferSent LeadPhase = "lead_offer_sent"

        // LeadConverted indicates the lead has converted.
        LeadConverted LeadPhase = "lead_converted"

        // LeadNegative indicates the lead responded negatively.
        LeadNegative LeadPhase = "lead_negative"

        // LeadArchived indicates the lead has been archived.
        LeadArchived LeadPhase = "lead_archived"

        // LeadAutoReply indicates an auto-reply was detected from the lead.
        LeadAutoReply LeadPhase = "lead_auto_reply"

        // LeadSkipped indicates the lead was skipped.
        LeadSkipped LeadPhase = "lead_skipped"

        // LeadNoResponse indicates the lead never responded.
        LeadNoResponse LeadPhase = "lead_no_response"

        // LeadFollowUp1 indicates the first follow-up has been sent.
        LeadFollowUp1 LeadPhase = "lead_follow_up_1"

        // LeadFollowUp2 indicates the second follow-up has been sent.
        LeadFollowUp2 LeadPhase = "lead_follow_up_2"

        // LeadCold indicates the lead has gone cold.
        LeadCold LeadPhase = "lead_cold"

        // LeadFailed indicates an action on this lead failed.
        LeadFailed LeadPhase = "lead_failed"

        // LeadDead indicates the lead is no longer reachable.
        LeadDead LeadPhase = "lead_dead"

        // LeadBlocked indicates the lead has blocked the sender.
        LeadBlocked LeadPhase = "lead_blocked"

        // LeadReContact indicates the lead is being re-contacted.
        LeadReContact LeadPhase = "lead_re_contact"
)

// AllLeadPhases returns all defined lead lifecycle phases.
func AllLeadPhases() []LeadPhase {
        return []LeadPhase{
                LeadBaru, LeadIceBreakerSent, LeadResponded, LeadOfferSent,
                LeadConverted, LeadNegative, LeadArchived, LeadAutoReply,
                LeadSkipped, LeadNoResponse, LeadFollowUp1, LeadFollowUp2,
                LeadCold, LeadFailed, LeadDead, LeadBlocked, LeadReContact,
        }
}

// IsValidLeadPhase reports whether p is a known lead phase.
func IsValidLeadPhase(p LeadPhase) bool {
        for _, v := range AllLeadPhases() {
                if v == p {
                        return true
                }
        }
        return false
}

// FollowUpPhase represents a phase in the follow-up pipeline.
// Used by the comms/followup screen to filter leads by stage.
// The string values are the short forms used in backend data (without "lead_" prefix).
type FollowUpPhase string

const (
        // FUPhase1 is the first follow-up phase (after ice breaker).
        FUPhase1 FollowUpPhase = "follow_up_1"

        // FUPhase2 is the second follow-up phase.
        FUPhase2 FollowUpPhase = "follow_up_2"

        // FUPhaseCold is the cold/unged phase (no response after 2 follow-ups).
        FUPhaseCold FollowUpPhase = "cold"

        // FUPhaseRecontact is the re-contact phase (previously responded, went quiet).
        FUPhaseRecontact FollowUpPhase = "recontact"
)

// AllFollowUpPhases returns all defined follow-up phases.
func AllFollowUpPhases() []FollowUpPhase {
        return []FollowUpPhase{FUPhase1, FUPhase2, FUPhaseCold, FUPhaseRecontact}
}

// IsValidFollowUpPhase reports whether p is a known follow-up phase.
func IsValidFollowUpPhase(p FollowUpPhase) bool {
        for _, v := range AllFollowUpPhases() {
                if v == p {
                        return true
                }
        }
        return false
}

// ResponseClass represents the auto-classification category of a lead response.
// The backend assigns these and the TUI renders matching badges.
type ResponseClass string

const (
        ClassPositif   ResponseClass = "positif"
        ClassCurious   ResponseClass = "penasaran"
        ClassAutoReply ResponseClass = "auto-reply"
)

// LicenseResult represents the outcome of a license validation request from the backend.
//
// LicenseResult is distinct from screen states because it describes the business-domain
// result of a validation attempt (valid/invalid/expired/conflict/server_error) rather
// than which view the License screen is showing.
type LicenseResult string

const (
        // LicenseResultValid indicates the license key is valid and active.
        LicenseResultValid LicenseResult = "valid"

        // LicenseResultInvalid indicates the license key is not recognized.
        LicenseResultInvalid LicenseResult = "invalid"

        // LicenseResultExpired indicates the license has passed its expiration date.
        LicenseResultExpired LicenseResult = "expired"

        // LicenseResultDeviceConflict indicates the license is active on another device.
        LicenseResultDeviceConflict LicenseResult = "device_conflict"

        // LicenseResultServerError indicates the license server could not be reached.
        LicenseResultServerError LicenseResult = "server_error"
)

// AllLicenseResults returns all defined license result types.
func AllLicenseResults() []LicenseResult {
        return []LicenseResult{
                LicenseResultValid, LicenseResultInvalid, LicenseResultExpired,
                LicenseResultDeviceConflict, LicenseResultServerError,
        }
}

// IsValidLicenseResult reports whether r is a known license result.
func IsValidLicenseResult(r LicenseResult) bool {
        for _, v := range AllLicenseResults() {
                if v == r {
                        return true
                }
        }
        return false
}

// CheckResult represents the outcome of a startup update check.
//
// CheckResult is distinct from screen states because it describes the
// business-domain result of a background version check (latest/minor/major)
// rather than which view the Update screen is showing. The backend sends
// these as the ParamCheckResult value; the TUI must compare against these
// constants, not raw strings.
type CheckResult string

const (
        // CheckResultLatest indicates the installed version is the latest available.
        CheckResultLatest CheckResult = "latest"

        // CheckResultUpdate indicates a minor update is available.
        CheckResultUpdate CheckResult = "update"

        // CheckResultUpgrade indicates a major upgrade is available.
        CheckResultUpgrade CheckResult = "upgrade"
)

// AllCheckResults returns all defined check result types.
func AllCheckResults() []CheckResult {
        return []CheckResult{CheckResultLatest, CheckResultUpdate, CheckResultUpgrade}
}

// IsValidCheckResult reports whether r is a known check result.
func IsValidCheckResult(r CheckResult) bool {
        for _, v := range AllCheckResults() {
                if v == r {
                        return true
                }
        }
        return false
}

// License key prefix constants — DRY: used by license and update screens.
const (
        // LicenseKeyPrefixV1 is the v1 license key prefix format.
        LicenseKeyPrefixV1 = "WACL"

        // LicenseKeyPrefixV2 is the v2 license key prefix format.
        LicenseKeyPrefixV2 = "WACL2"
)

// LicenseAction identifies a semantic action that the license screen can send
// to the backend via the bus. Using constants instead of raw strings follows
// the DRY convention: "No hardcoded protocol strings."
type LicenseAction string

const (
        // ActionValidateLicense requests the backend to validate a license key.
        ActionValidateLicense LicenseAction = "validate_license"

        // ActionBuyLicense opens the license purchase flow.
        ActionBuyLicense LicenseAction = "buy_license"

        // ActionBuyRenewal opens the license renewal purchase flow.
        ActionBuyRenewal LicenseAction = "buy_renewal"

        // ActionLicenseContinue proceeds after a valid license (or auto-transition).
        ActionLicenseContinue LicenseAction = "license_continue"

        // ActionLicenseOfflineContinue proceeds using offline grace period.
        ActionLicenseOfflineContinue LicenseAction = "license_offline_continue"

        // ActionForceDisconnectDevice forces a disconnect on the other device.
        ActionForceDisconnectDevice LicenseAction = "force_disconnect_device"
)

// ActivityStatus represents the status of an activity event in the timeline.
// Used by the monitor dashboard to highlight specific event types.
// The backend sends these as the "status" field in activity event data;
// the TUI must compare against these constants, not raw strings.
type ActivityStatus string

const (
        // ActivityStatusRespond indicates a lead responded.
        ActivityStatusRespond ActivityStatus = "respond"

        // ActivityStatusSent indicates a message was sent.
        ActivityStatusSent ActivityStatus = "sent"

        // ActivityStatusDelivered indicates a message was delivered.
        ActivityStatusDelivered ActivityStatus = "delivered"

        // ActivityStatusScraped indicates a scraping event.
        ActivityStatusScraped ActivityStatus = "scraped"

        // ActivityStatusConverted indicates a conversion event.
        ActivityStatusConverted ActivityStatus = "converted"
)

// CommsAction identifies a semantic action that the comms screens can send
// to the backend via the bus. Using constants instead of raw strings follows
// the DRY convention: "No hardcoded protocol strings."
type CommsAction string

const (
        // ActionComposeSend sends the composed message to the target lead.
        ActionComposeSend CommsAction = "compose_send"

        // ActionFollowUpAutoAll approves all pending follow-ups automatically.
        ActionFollowUpAutoAll CommsAction = "followup_auto_all"

        // ActionFollowUpSkipWait skips the wait timer for the current follow-up.
        ActionFollowUpSkipWait CommsAction = "followup_skip_wait"

        // ActionFollowUpPause pauses the follow-up sending process.
        ActionFollowUpPause CommsAction = "followup_pause"

        // ActionFollowUpSendFinal sends the final (3rd) follow-up to a cold lead.
        ActionFollowUpSendFinal CommsAction = "followup_send_final"

        // ActionFollowUpArchiveCold archives all cold leads.
        ActionFollowUpArchiveCold CommsAction = "followup_archive_cold"

        // ActionFollowUpRecontact sends a re-contact message to a specific lead.
        ActionFollowUpRecontact CommsAction = "followup_recontact"

        // ActionFollowUpRecontactAll sends re-contact messages to all eligible leads.
        ActionFollowUpRecontactAll CommsAction = "followup_recontact_all"

        // ActionHistoryPrevDay requests data for the previous day in history view.
        ActionHistoryPrevDay CommsAction = "history_prev_day"

        // ActionHistoryDayDetail requests detail data for a specific day.
        ActionHistoryDayDetail CommsAction = "history_day_detail"
)

// UpdateAction identifies a semantic action that the update screen can send
// to the backend via the bus. Using constants instead of raw strings follows
// the DRY convention: "No hardcoded protocol strings."
type UpdateAction string

const (
        // ActionStartDownload begins downloading the available update.
        ActionStartDownload UpdateAction = "start_download"

        // ActionRemindLater postpones the update reminder.
        ActionRemindLater UpdateAction = "remind_later"

        // ActionSkipUpdate skips the current update entirely.
        ActionSkipUpdate UpdateAction = "skip_update"

        // ActionCancelDownload cancels an in-progress download.
        ActionCancelDownload UpdateAction = "cancel_download"

        // ActionRestartNow restarts the app immediately to apply the update.
        ActionRestartNow UpdateAction = "restart_now"

        // ActionRestartLater postpones the restart prompt.
        ActionRestartLater UpdateAction = "restart_later"

        // ActionSkipRestart skips the restart for this session.
        ActionSkipRestart UpdateAction = "skip_restart"

        // ActionViewUpgradeDetails opens the v2 upgrade details view.
        ActionViewUpgradeDetails UpdateAction = "view_upgrade_details"

        // ActionStayV1 declines the v2 upgrade and stays on v1.
        ActionStayV1 UpdateAction = "stay_v1"

        // ActionCancelLicenseInput cancels the license key input dialog.
        ActionCancelLicenseInput UpdateAction = "cancel_license_input"

        // ActionRenewV1 initiates a v1 license renewal.
        ActionRenewV1 UpdateAction = "renew_v1"

        // ActionUpgradeV2 initiates the v2 license upgrade.
        ActionUpgradeV2 UpdateAction = "upgrade_v2"

        // ActionEnterNewLicense opens the new license key input dialog.
        ActionEnterNewLicense UpdateAction = "enter_new_license"

        // ActionExitExpired exits the app when license is expired.
        ActionExitExpired UpdateAction = "exit_expired"
)

// SendAction identifies a semantic action that the send screen can send
// to the backend via the bus. Using constants instead of raw strings follows
// the DRY convention: "No hardcoded protocol strings."
type SendAction string

const (
        // ActionTogglePause toggles the send pipeline between paused and active.
        ActionTogglePause SendAction = "toggle_pause"

        // ActionValidateRetry retries sending after a validation failure.
        ActionValidateRetry SendAction = "validate_retry"

        // ActionApproveResponse approves sending an offer in response to a lead's message.
        ActionApproveResponse SendAction = "approve_response"

        // ActionSkipWait skips the wait timer for the next send.
        ActionSkipWait SendAction = "skip_wait"

        // ActionEmergencySend sends a message even outside work hours.
        ActionEmergencySend SendAction = "emergency_send"

        // ActionRetryManual retries sending to a failed lead manually.
        ActionRetryManual SendAction = "retry_manual"

        // ActionReLogin re-logs into a disconnected WhatsApp number.
        ActionReLogin SendAction = "re_login"

        // ActionLoginOneByOne logs into WhatsApp numbers sequentially.
        ActionLoginOneByOne SendAction = "login_one_by_one"

        // ActionSkipFailure skips the current failed lead and moves on.
        ActionSkipFailure SendAction = "skip_failure"
)

// ShieldAction identifies a semantic action that the anti-ban shield screen
// can send to the backend via the bus. Using constants instead of raw strings
// follows the DRY convention: "No hardcoded protocol strings."
type ShieldAction string

const (
        // ActionRefreshShield requests a refresh of the shield data.
        ActionRefreshShield ShieldAction = "refresh_shield"

        // ActionAddWANumber opens the flow to add a new WhatsApp number.
        ActionAddWANumber ShieldAction = "add_wa_number"

        // ActionPauseSending pauses all sending from the shield screen.
        ActionPauseSending ShieldAction = "pause_sending"
)

// SlotStatus identifies the status of a WhatsApp number slot in the anti-ban
// shield screen. These values come from the backend and are used for rendering
// status indicators and conditional logic.
const (
        // SlotStatusActive indicates the WA number is actively sending.
        SlotStatusActive = "active"

        // SlotStatusCooldown indicates the WA number is in cooldown.
        SlotStatusCooldown = "cooldown"

        // SlotStatusFlagged indicates the WA number has been flagged by WhatsApp.
        SlotStatusFlagged = "flagged"
)

// LicenseKeyFormat defines the structural format of a WaClaw license key.
// These constants are the single source of truth for both the backend and TUI.
// The format is: LicenseKeyPrefix + "-" + (4 groups of LicenseKeyGroupSize chars separated by "-")
// Example: WACL-ABCD-EFGH-IJKL-MNOP
const (
        // LicenseKeyPrefix is the standard v1 license key prefix.
        LicenseKeyPrefix = "WACL"

        // LicenseKeyV2Prefix is the v2 license key prefix (for major upgrades).
        LicenseKeyV2Prefix = "WACL2"

        // LicenseKeyGroupSize is the number of characters per hyphen-separated group.
        LicenseKeyGroupSize = 4

        // LicenseKeyGroups is the number of character groups after the prefix.
        LicenseKeyGroups = 4

        // LicenseKeyFormattedLen is the total formatted key length including prefix and hyphens.
        // Calculated: len("WACL-") + 4*4 + 3 hyphens = 5 + 16 + 3 = 24
        LicenseKeyFormattedLen = len(LicenseKeyPrefix) + 1 + LicenseKeyGroupSize*LicenseKeyGroups + (LicenseKeyGroups - 1)

        // LicenseKeyRawLen is the number of user-typed characters (excluding prefix and hyphens).
        // The user types 16 chars; the prefix "WACL" is auto-prepended.
        LicenseKeyRawLen = LicenseKeyGroupSize * LicenseKeyGroups

        // ValidationStepCount is the number of sequential steps in the license validation animation.
        // Doc spec: "● ○ ○ sequential animate" = 3 steps.
        ValidationStepCount = 3
)

// ParamKey constants define the JSON-RPC parameter field names used between
// the backend and TUI. Using constants instead of raw strings follows the
// DRY convention: "No hardcoded protocol strings." These are shared across
// HandleNavigate, HandleUpdate, and publishAction calls.
const (
        // ParamState is the state override key in navigate/update params.
        ParamState = "state"

        // ParamResult is the validation result key in update params.
        ParamResult = "result"

        // ParamLicenseKey is the license key field in navigate/update/action params.
        ParamLicenseKey = "license_key"

        // ParamDevice is the device name field.
        ParamDevice = "device"

        // ParamExpires is the expiration date field.
        ParamExpires = "expires"

        // ParamExpiredAgo is the human-readable "how long ago expired" field.
        ParamExpiredAgo = "expired_ago"

        // ParamOtherDevice is the conflicting device name field.
        ParamOtherDevice = "other_device"

        // ParamLastActive is the last active time on other device field.
        ParamLastActive = "last_active"

        // ParamGraceHours is the offline grace hours remaining field.
        ParamGraceHours = "grace_hours"

        // ParamCurrentVersion is the current app version field.
        ParamCurrentVersion = "current_version"

        // ParamNewVersion is the available version field.
        ParamNewVersion = "new_version"

        // ParamChangelog is the changelog items field.
        ParamChangelog = "changelog"

        // ParamPercent is the download progress percentage field.
        ParamPercent = "percent"

        // ParamSize is the download total size field.
        ParamSize = "size"

        // ParamDownloaded is the downloaded amount field.
        ParamDownloaded = "downloaded"

        // ParamSpeed is the download speed field.
        ParamSpeed = "speed"

        // ParamETA is the estimated time remaining field.
        ParamETA = "eta"

        // ParamSource is the download source URL field.
        ParamSource = "source"

        // ParamChecksumVerified is the checksum verification status field.
        ParamChecksumVerified = "checksum_verified"

        // ParamBackupPath is the old binary backup path field.
        ParamBackupPath = "backup_path"

        // ParamExpiredDate is the license expiration date field.
        ParamExpiredDate = "expired_date"

        // ParamLicenseExpiry is the current license expiry field.
        ParamLicenseExpiry = "license_expiry"

        // ParamLicenseStatus is the current license status field.
        ParamLicenseStatus = "license_status"

        // ParamCheckResult is the startup check result field.
        ParamCheckResult = "check_result"

        // ParamLicensePrefix is the v2 license key prefix field sent by the backend.
        // DRY: replaces the hardcoded "license_prefix" string previously used in
        // the update screen's HandleNavigate. The backend should send this so the
        // TUI does not need to guess the v2 key format.
        ParamLicensePrefix = "license_prefix"

        // ParamKey is the generic key field for license input actions.
        ParamKey = "key"

        // ParamVersion is the version field for update actions.
        ParamVersion = "version"

        // --- Send screen param keys ---

        // ParamPause is the pause reason field in navigate params.
        ParamPause = "pause"

        // ParamOffHours is the off-hours reason field in navigate params.
        ParamOffHours = "off_hours"

        // ParamWorkHours is the work hours config field.
        ParamWorkHours = "work_hours"

        // ParamSlots is the slot list field (used by both send and shield screens).
        ParamSlots = "slots"

        // ParamQueue is the queue item list field.
        ParamQueue = "queue"

        // ParamStats is the aggregate statistics map field.
        ParamStats = "stats"

        // ParamFailure is the failure detail map field.
        ParamFailure = "failure"

        // ParamResponse is the response interrupt map field.
        ParamResponse = "response"

        // ParamNicheNames is the niche name list field.
        ParamNicheNames = "niche_names"

        // ParamRateHour is the hourly rate count field.
        ParamRateHour = "rate_hour"

        // ParamDailySent is the daily sent count field.
        ParamDailySent = "daily_sent"

        // ParamQueueCount is the total queue count field.
        ParamQueueCount = "queue_count"

        // ParamNextSendTime is the next send timestamp field.
        ParamNextSendTime = "next_send_time"

        // --- Shield (Anti-Ban) screen param keys ---

        // ParamHealth is the health score field.
        ParamHealth = "health"

        // --- Slot data sub-keys (used inside the ParamSlots list items) ---

        // ParamSlotNumber is the WA number field inside a slot item.
        ParamSlotNumber = "number"

        // ParamSlotStatus is the slot status field inside a slot item.
        ParamSlotStatus = "status"

        // ParamSlotSentHour is the sent-this-hour count inside a slot item.
        ParamSlotSentHour = "sent_hour"

        // ParamSlotMaxHour is the max-per-hour limit inside a slot item.
        ParamSlotMaxHour = "max_hour"

        // ParamSlotSentToday is the sent-today count inside a slot item.
        ParamSlotSentToday = "sent_today"

        // ParamSlotWarnings is the warning count inside a slot item.
        ParamSlotWarnings = "warnings"

        // ParamSlotCooldown is the cooldown time remaining inside a slot item.
        ParamSlotCooldown = "cooldown"

        // ParamSlotReadyIn is the ready-in time inside a slot item.
        ParamSlotReadyIn = "ready_in"

        // ParamSlotHealthy is the healthy boolean inside a slot item.
        ParamSlotHealthy = "healthy"

        // --- Send data sub-keys (used inside send screen list/map items) ---

        // ParamSlotReadyInAlt is the ready_in alternate key used by send slots.
        ParamSlotReadyInAlt = "ready_in"

        // ParamRateMax is the hourly rate max inside send slot data.
        ParamRateMax = "rate_max"

        // ParamIndex is the item index field.
        ParamIndex = "index"

        // ParamName is the name field inside various data items.
        ParamName = "name"

        // ParamSlot is the assigned slot field inside a queue item.
        ParamSlot = "slot"

        // ParamTemplate is the template type field.
        ParamTemplate = "template"

        // ParamVariant is the template variant field.
        ParamVariant = "variant"

        // ParamNextIn is the next-in time field.
        ParamNextIn = "next_in"

        // ParamNiche is the niche name field inside a queue item.
        ParamNiche = "niche"

        // ParamRotated is the rotation flag field.
        ParamRotated = "rotated"

        // ParamReason is the failure reason field.
        ParamReason = "reason"

        // ParamHint is the failure hint field.
        ParamHint = "hint"

        // ParamMessage is the message text field.
        ParamMessage = "message"

        // ParamRateHourSent is the hourly sent count inside send stats.
        ParamRateHourSent = "rate_hour_sent"

        // ParamRateHourMax is the hourly rate max inside send stats.
        ParamRateHourMax = "rate_hour_max"

        // ParamRateSlotCount is the active slot count inside send stats.
        ParamRateSlotCount = "rate_slot_count"

        // ParamDailyResp is the daily response count inside send stats.
        ParamDailyResp = "daily_resp"

        // ParamDailyConv is the daily conversion count inside send stats.
        ParamDailyConv = "daily_conv"

        // ParamDailyMax is the daily send limit inside send stats.
        ParamDailyMax = "daily_max"

        // ParamNextSendSlot is the next send slot name inside send stats.
        ParamNextSendSlot = "next_send_slot"

        // ParamNow is the current time string inside send stats.
        ParamNow = "now"

        // --- Comms screen param keys ---

        // ParamTarget is the compose target (business name) field.
        ParamTarget = "target"

        // ParamDraft is the compose draft text field.
        ParamDraft = "draft"

        // ParamSnippets is the quick-reply snippets list field.
        ParamSnippets = "snippets"

        // ParamMaxChars is the soft character limit for compose messages.
        ParamMaxChars = "max_chars"

        // ParamDate is the date string field (used by history and actions).
        ParamDate = "date"

        // ParamEvents is the timeline events list field.
        ParamEvents = "events"

        // ParamWeek is the weekly chart data map field.
        ParamWeek = "week"

        // ParamAvgResponseTime is the average response time string field.
        ParamAvgResponseTime = "avg_response_time"

        // ParamNiches is the niche groups list field (used by follow-up and niche select).
        ParamNiches = "niches"

        // ParamTotalToday is the total follow-ups due today field.
        ParamTotalToday = "total_today"

        // ParamColdTotal is the total cold leads count field.
        ParamColdTotal = "cold_total"

        // ParamIceBreakerUnanswered is the count of unanswered ice breakers field.

        ParamIceBreakerUnanswered = "ice_breaker_unanswered"

        // ParamColdLeads is the cold leads list field.
        ParamColdLeads = "cold_leads"

        // ParamRecontactLeads is the re-contact eligible leads list field.
        ParamRecontactLeads = "recontact_leads"

        // ParamSendingRate is the sending rate string field (e.g. "9/18").
        ParamSendingRate = "sending_rate"

        // ParamSendingDone is the count of sent follow-ups field.
        ParamSendingDone = "sending_done"

        // ParamSendingTotal is the total follow-ups to send field.
        ParamSendingTotal = "sending_total"

        // ParamVariantNames is the template variant file names list field.
        ParamVariantNames = "variant_names"

        // ParamVariantPreviews is the template variant preview text list field.
        ParamVariantPreviews = "variant_previews"

        // ParamVariantManualOnly is the list of manual-only variant indices field.
        ParamVariantManualOnly = "variant_manual_only"

        // ParamMaxSendingVisible is the backend-configurable visible sending leads limit.
        ParamMaxSendingVisible = "max_sending_visible"

        // ParamMaxColdVisible is the backend-configurable visible cold leads limit.
        ParamMaxColdVisible = "max_cold_visible"

        // ParamBusinessName is the business name field (used in cold/recontact actions).
        ParamBusinessName = "business_name"

        // --- Boot / Onboarding screen param keys ---

        // ParamWACount is the WhatsApp slot count field.
        ParamWACount = "wa_count"

        // ParamNicheCount is the niche count field.
        ParamNicheCount = "niche_count"

        // ParamWorkersList is the workers data list in infra context.
        // Separate key from ParamWorkers (boot marching army) to avoid collision.
        ParamWorkersList = "workers_list"

        // ParamPipelineTotals is the pipeline totals data field.
        ParamPipelineTotals = "totals"

        // ParamWorkerCount is the worker count field.
        ParamWorkerCount = "worker_count"

        // ParamLeadsCount is the leads count field.
        ParamLeadsCount = "leads_count"

        // ParamResponseCount is the response count field.
        ParamResponseCount = "response_count"

        // ParamOkNicheCount is the count of OK niches field.
        ParamOkNicheCount = "ok_niche_count"

        // ParamErrorNicheCount is the count of error niches field.
        ParamErrorNicheCount = "error_niche_count"

        // ParamErrorNiche is the error niche name field.
        ParamErrorNiche = "error_niche"

        // ParamDeviceName is the device name field (boot device conflict).
        ParamDeviceName = "device_name"

        // ParamNicheAlreadySet is the boolean flag for returning users with niches.
        ParamNicheAlreadySet = "niche_already_set"

        // ParamWorkers is the worker rows list field (boot army march).
        ParamWorkers = "workers"

        // --- Login screen param keys ---

        // ParamFilledSlots is the count of filled WA slots field.
        ParamFilledSlots = "filled_slots"

        // ParamTotalSlots is the total WA slots field.
        ParamTotalSlots = "total_slots"

        // ParamActiveSlot is the active slot index field.
        ParamActiveSlot = "active_slot"

        // ParamContactCount is the contact count field.
        ParamContactCount = "contact_count"

        // ParamExpiredSlot is the expired slot index field.
        ParamExpiredSlot = "expired_slot"

        // ParamActiveSlots is the count of active WA slots field.
        ParamActiveSlots = "active_slots"

        // ParamLastSessionAgo is the human-readable time since last session field.
        ParamLastSessionAgo = "last_session_ago"

        // ParamPhoneNumbers is the phone numbers list field.
        ParamPhoneNumbers = "phone_numbers"

        // ParamQRData is the QR code data string field.
        ParamQRData = "qr_data"

        // --- Monitor screen param keys ---

        // ParamWASlots is the WA slot overview list field (monitor dashboard).
        ParamWASlots = "wa_slots"

        // ParamActivities is the activity events list field.
        ParamActivities = "activities"

        // ParamPending is the pending responses count field.
        ParamPending = "pending"

        // ParamTodayStats is the today stats map field.
        ParamTodayStats = "today_stats"

        // ParamWeekStats is the week stats map field.
        ParamWeekStats = "week_stats"

        // ParamLead is the lead data map field (response screen).
        ParamLead = "lead"

        // ParamConversion is the conversion data map field.
        ParamConversion = "conversion"

        // --- Follow-up lead data sub-keys ---

        // ParamPhase is the follow-up phase field inside a lead item.
        ParamPhase = "phase"

        // ParamPreviousAction is the previous action description field.
        ParamPreviousAction = "previous_action"

        // ParamNextAction is the next action description field.
        ParamNextAction = "next_action"

        // ParamSlotNumber is the WA slot number field inside a lead item.
        ParamSlotNumberAlt = "slot_number"

        // ParamVariantName is the template variant name field inside a lead item.
        ParamVariantName = "variant_name"

        // ParamIsSending is the sending-in-progress boolean field.
        ParamIsSending = "is_sending"

        // ParamWaitTime is the remaining wait time string field.
        ParamWaitTime = "wait_time"

        // ParamIceBreakerAction is the ice breaker timing description field.
        ParamIceBreakerAction = "ice_breaker_action"

        // ParamFollowUp1Action is the follow-up 1 timing description field.
        ParamFollowUp1Action = "follow_up_1_action"

        // ParamFollowUp2Action is the follow-up 2 timing description field.
        ParamFollowUp2Action = "follow_up_2_action"

        // ParamPreviousResponse is the lead's previous response text field.
        ParamPreviousResponse = "previous_response"

        // ParamDaysSinceResponse is the days since last response field.
        ParamDaysSinceResponse = "days_since_response"

        // ParamDaysSinceOffer is the days since offer sent field.
        ParamDaysSinceOffer = "days_since_offer"

        // ParamCanRecontact is the re-contact eligibility boolean field.
        ParamCanRecontact = "can_recontact"

        // ParamDaysSinceLastAction is the days since last action field.
        ParamDaysSinceLastAction = "days_since_last_action"

        // --- History week data sub-keys ---

        // ParamDayLabels is the day label list field inside week data.
        ParamDayLabels = "day_labels"

        // ParamMessages is the message counts list field.
        ParamMessages = "messages"

        // ParamResponses is the response counts list field.
        ParamResponses = "responses"

        // ParamConverts is the conversion counts list field.
        ParamConverts = "converts"

        // ParamBestDayIndex is the best performing day index field.
        ParamBestDayIndex = "best_day_index"

        // ParamBestDayLabel is the best performing day label field.
        ParamBestDayLabel = "best_day_label"

        // ParamBestDayConvRate is the conversion rate for the best day field.
        ParamBestDayConvRate = "best_day_conv_rate"

        // --- History event data sub-keys ---

        // ParamIcon is the emoji icon field inside event data.
        ParamIcon = "icon"

        // ParamEventType is the event type label field.
        ParamEventType = "event_type"

        // ParamHighlight is the highlight boolean field.
        ParamHighlight = "highlight"

        // ParamIsConversion is the conversion flag field.
        ParamIsConversion = "is_conversion"

        // --- Day stats sub-keys ---

        // ParamSent is the sent count field inside stats.
        ParamSent = "sent"

        // ParamRespond is the respond count field inside stats.
        ParamRespond = "respond"

        // ParamConvert is the convert count field inside stats.
        ParamConvert = "convert"

        // ParamNewLeads is the new leads count field.
        ParamNewLeads = "new_leads"

        // ParamScrapes is the scrape count field.
        ParamScrapes = "scrapes"

        // --- Niche group sub-keys ---

        // ParamNicheName is the niche name field inside a niche group.
        ParamNicheName = "niche_name"

        // ParamFU1Count is the follow-up 1 count field inside a niche group.
        ParamFU1Count = "fu1_count"

        // ParamFU2Count is the follow-up 2 count field inside a niche group.
        ParamFU2Count = "fu2_count"

        // ParamColdCount is the cold count field inside a niche group.
        ParamColdCount = "cold_count"

        // ParamLeads is the leads list field inside a niche group.
        ParamLeads = "leads"

        // --- Pipeline / Scrape screen param keys ---

        // ParamArea is the area string field.
        ParamArea = "area"

        // ParamFilter is the filter string field.
        ParamFilter = "filter"

        // ParamFound is the leads found count field.
        ParamFound = "found"

        // ParamQualified is the leads qualified count field.
        ParamQualified = "qualified"

        // ParamTemplates is the templates list field.
        ParamTemplates = "templates"

        // ParamTemplateText is the template text content field.
        ParamTemplateText = "template_text"

        // ParamCurrent is the current item index field.
        ParamCurrent = "current"

        // ParamIsMajor is the major update boolean field.
        ParamIsMajor = "is_major"

        // --- Settings / Config screen param keys ---

        // ParamActiveNiches is the active niches description field.
        ParamActiveNiches = "active_niches"

        // ParamWASlotsDesc is the WA slots description field (settings screen).
        ParamWASlotsDesc = "wa_slots_desc"

        // ParamWorkerPool is the worker pool description field.
        ParamWorkerPool = "worker_pool"

        // ParamRateLimit is the rate limit description field.
        ParamRateLimit = "rate_limit"

        // ParamRotatorMode is the rotator mode description field.
        ParamRotatorMode = "rotator_mode"

        // ParamAutopilot is the autopilot status field.
        ParamAutopilot = "autopilot"

        // ParamLocale is the locale field.
        ParamLocale = "locale"

        // --- Week stats sub-keys ---

        // ParamTotalSent is the total sent count field.
        ParamTotalSent = "total_sent"

        // ParamTotalRespond is the total respond count field.
        ParamTotalRespond = "total_respond"

        // ParamTotalConvert is the total convert count field.
        ParamTotalConvert = "total_convert"

        // ParamTotalNewLeads is the total new leads count field.
        ParamTotalNewLeads = "total_new_leads"

        // ParamTotalScrapes is the total scrapes count field.
        ParamTotalScrapes = "total_scrapes"

        // ParamValid is the validation valid boolean field.
        ParamValid = "valid"

        // ParamInsight is the insight text field.
        ParamInsight = "insight"

        // ParamStreak is the streak days count field.
        ParamStreak = "streak"

        // --- Data screen (Leads DB) param keys ---

        // ParamTotal is the total count field (leads DB header).
        ParamTotal = "total"

        // ParamCategories is the filter categories list field.
        ParamCategories = "categories"

        // ParamLeadName is the lead name field (action params).
        ParamLeadName = "lead_name"

        // ParamNote is the context note field inside a filter category.
        ParamNote = "note"

        // ParamCount is the count field inside a filter category or source detail.
        ParamCount = "count"

        // ParamCategory is the business category field inside a lead item.
        ParamCategory = "category"

        // ParamAddress is the street address field inside a lead item.
        ParamAddress = "address"

        // ParamCity is the city field inside a lead item.
        ParamCity = "city"

        // ParamRating is the rating field inside a lead item.
        ParamRating = "rating"

        // ParamReviews is the review count field inside a lead item.
        ParamReviews = "reviews"

        // ParamHasWeb is the has-website boolean field inside a lead item.
        ParamHasWeb = "has_web"

        // ParamHasInsta is the has-instagram boolean field inside a lead item.
        ParamHasInsta = "has_insta"

        // ParamPhotoCount is the photo count field inside a lead item.
        ParamPhotoCount = "photo_count"

        // ParamScore is the lead score field inside a lead item.
        ParamScore = "score"

        // ParamResponseText is the response text field inside a lead item.
        ParamResponseText = "response_text"

        // ParamIceBreakerTime is the ice breaker timestamp field inside a lead item.
        ParamIceBreakerTime = "ice_breaker_time"

        // ParamResponseTime is the response timestamp field inside a lead item.
        ParamResponseTime = "response_time"

        // ParamFollowupCount is the follow-up count field inside a lead item.
        ParamFollowupCount = "followup_count"

        // ParamFollowupTimes is the follow-up timestamps list field inside a lead item.
        ParamFollowupTimes = "followup_times"

        // ParamFollowupDueText is the follow-up due date text field inside a lead item.
        ParamFollowupDueText = "followup_due_text"

        // ParamDuration is the conversion duration field inside a lead item.
        ParamDuration = "duration"

        // ParamTemplateName is the template name field (used by leads DB and template mgr).
        ParamTemplateName = "template_name"

        // ParamWorkerName is the worker name field inside a lead item.
        ParamWorkerName = "worker_name"

        // ParamTimeline is the timeline events list field inside a lead item.
        ParamTimeline = "timeline"

        // ParamAction is the action field inside a timeline event.
        ParamAction = "action"

        // ParamDetail is the detail/hint text field (used in timeline events, errors, gen files).
        ParamDetail = "detail"

        // ParamRevenue is the revenue field inside a converted lead item.
        ParamRevenue = "revenue"

        // ParamTime is the timestamp field inside a timeline event.
        ParamTime = "time"

        // --- Template Manager screen param keys ---

        // ParamGroups is the template groups list field.
        ParamGroups = "groups"

        // ParamPreviewValues is the template preview substitution values map field.
        ParamPreviewValues = "preview_values"

        // ParamSubstitutionValues is the template substitution values map field (alternative key).
        ParamSubstitutionValues = "substitution_values"

        // ParamTemplateType is the template type field (ice_breaker / offer).
        ParamTemplateType = "template_type"

        // ParamFilePath is the file path field (template edit hint, validation error).
        ParamFilePath = "file_path"

        // ParamType is the type discriminator field inside template groups.
        ParamType = "type"

        // ParamLabel is the display label field inside template groups and filter categories.
        ParamLabel = "label"

        // ParamPreview is the short preview text field inside a template entry.
        ParamPreview = "preview"

        // ParamRecommended is the recommended flag field inside a template entry.
        ParamRecommended = "recommended"

        // ParamBody is the full body text field inside a template entry.
        ParamBody = "body"

        // ParamHasError is the has-error boolean field inside a template entry.
        ParamHasError = "has_error"

        // ParamErrors is the errors list field (template validation, niche config).
        ParamErrors = "errors"

        // ParamLine is the line number field inside an error entry.
        ParamLine = "line"

        // ParamSeverity is the severity field inside a template error entry.
        ParamSeverity = "severity"

        // ParamCode is the error code field inside a template error entry.
        ParamCode = "code"

        // --- Niche Select screen param keys ---

        // ParamTargets is the target keywords list field.
        ParamTargets = "targets"

        // ParamFilters is the filter entries list field.
        ParamFilters = "filters"

        // ParamAreas is the area entries list field.
        ParamAreas = "areas"

        // ParamErrorFile is the config error file path field.
        ParamErrorFile = "error_file"

        // ParamSelected is the selected boolean / labels field.
        ParamSelected = "selected"

        // ParamFile is the file path field (niche config error action).
        ParamFile = "file"

        // ParamDescription is the description text field (niche item, config error).
        ParamDescription = "description"

        // ParamEmoji is the emoji field (niche item, explorer category).
        ParamEmoji = "emoji"

        // ParamPointer is the error pointer field inside a config error.
        ParamPointer = "pointer"

        // --- Niche Explorer screen param keys ---

        // ParamSources is the source detail list field (explorer category).
        ParamSources = "sources"

        // ParamGenFiles is the generation file status list field.
        ParamGenFiles = "gen_files"

        // ParamGenNicheName is the generated niche name field.
        ParamGenNicheName = "gen_niche_name"

        // ParamFolderSlug is the authoritative folder slug from backend.
        ParamFolderSlug = "folder_slug"

        // ParamGenProgress is the generation progress float field.
        ParamGenProgress = "gen_progress"

        // ParamAreaAutoDetect is the area auto-detect boolean field.
        ParamAreaAutoDetect = "area_auto_detect"

        // ParamExistingAreas is the existing areas list field (auto-detect).
        ParamExistingAreas = "existing_areas"

        // ParamSearchResults is the search result list field.
        ParamSearchResults = "search_results"

        // ParamCategoryName is the category name field (backend-driven navigation).
        ParamCategoryName = "category_name"

        // ParamGenFilesDone is the completed generation files list field.
        ParamGenFilesDone = "gen_files_done"

        // ParamGenCurrentStep is the current generation step index field.
        ParamGenCurrentStep = "gen_current_step"

        // ParamQuery is the search query field (explorer search action).
        ParamQuery = "query"

        // ParamSubCount is the sub-category count field (explorer category).
        ParamSubCount = "sub_count"

        // ParamAreaCount is the area count field (explorer category / search result).
        ParamAreaCount = "area_count"

        // ParamSubCategories is the sub-category names list field.
        ParamSubCategories = "sub_categories"

        // ParamUnit is the unit label field inside a source detail.
        ParamUnit = "unit"

        // --- License / Update screen param keys ---

        // ParamKeyPrefix is the license key prefix field (license screen).
        ParamKeyPrefix = "key_prefix"

        // --- Settings screen param keys ---

        // ParamChanges is the config changes list field (settings reload success).
        ParamChanges = "changes"

        // ParamContext is the context lines field inside error/warning items.
        ParamContext = "context"

        // ParamField is the config field name inside a change entry.
        ParamField = "field"

        // ParamOldValue is the old value field inside a change entry.
        ParamOldValue = "old_value"

        // ParamNewValue is the new value field inside a change entry.
        ParamNewValue = "new_value"

        // --- Guardrail screen param keys ---

        // ParamResults is the validation results list field.
        ParamValidationResults = "results"

        // ParamStatus is the generic status field inside a validation result.
        ParamStatus = "status"

        // ParamDetails is the details text field inside a validation result.
        ParamDetails = "details"

        // ParamWarnings is the warnings list field inside a validation result.
        ParamWarnings = "warnings"

        // ParamSuggestion is the suggestion field inside a warning item.
        ParamSuggestion = "suggestion"

        // ParamCheckProgress is the check progress count field.
        ParamCheckProgress = "check_progress"

        // --- Workers screen param keys ---

        // ParamActive is the active boolean field inside a worker item.
        ParamActive = "active"

        // ParamTotals is the pipeline totals map field.
        ParamTotals = "totals"

        // ParamPassed is the passed/qualified count inside pipeline totals.
        ParamPassed = "passed"

        // ParamQueued is the queued count inside pipeline totals.
        ParamQueued = "queued"

        // --- Shield (Anti-Ban) additional param keys ---

        // ParamHealthScore is the health score field (alternate key).
        ParamHealthScore = "health_score"

        // ParamDailyBudgetSent is the daily budget sent count field.
        ParamDailyBudgetSent = "daily_budget_sent"

        // ParamDailyBudgetTotal is the daily budget total limit field.
        ParamDailyBudgetTotal = "daily_budget_total"

        // ParamWarningHealthScore is the health score for warning state.
        ParamWarningHealthScore = "warning_health_score"

        // ParamDangerHealthScore is the health score for danger state.
        ParamDangerHealthScore = "danger_health_score"

        // ParamFlaggedSlot is the flagged slot data map field.
        ParamFlaggedSlot = "flagged_slot"

        // ParamSlotDetail7DaySent is the 7-day sent count for slot detail.
        ParamSlotDetail7DaySent = "slot_detail_7day_sent"

        // ParamSlotDetail7DayResponded is the 7-day responded count for slot detail.
        ParamSlotDetail7DayResponded = "slot_detail_7day_responded"

        // ParamSlotDetail7DayFailed is the 7-day failed label for slot detail.
        ParamSlotDetail7DayFailed = "slot_detail_7day_failed"

        // ParamSlotDetail7DayWarnings is the 7-day warning count for slot detail.
        ParamSlotDetail7DayWarnings = "slot_detail_7day_warnings"

        // ParamSelectedSlotHealth is the health score for the selected slot.
        ParamSelectedSlotHealth = "selected_slot_health"

        // ParamPerSlotHourly is the per-slot hourly limit config field.
        ParamPerSlotHourly = "per_slot_hourly"

        // ParamPerSlotDaily is the per-slot daily limit config field.
        ParamPerSlotDaily = "per_slot_daily"

        // ParamCooldownMin is the cooldown minutes config field.
        ParamCooldownMin = "cooldown_min"

        // ParamMinDelayMin is the minimum delay minutes config field.
        ParamMinDelayMin = "min_delay_min"

        // ParamMaxDelayMin is the maximum delay minutes config field.
        ParamMaxDelayMin = "max_delay_min"

        // ParamDelayVariancePct is the delay variance percentage config field.
        ParamDelayVariancePct = "delay_variance_pct"

        // ParamAutoPause is the auto-pause mode config field.
        ParamAutoPause = "auto_pause"

        // ParamHealthThreshold is the health threshold config field.
        ParamHealthThreshold = "health_threshold"

        // ParamTemplateRotation is the template rotation config field.
        ParamTemplateRotation = "template_rotation"

        // ParamRotationMode is the rotation mode config field.
        ParamRotationMode = "rotation_mode"

        // ParamEmojiVariation is the emoji variation config field.
        ParamEmojiVariation = "emoji_variation"

        // ParamParagraphShuffle is the paragraph shuffle config field.
        ParamParagraphShuffle = "paragraph_shuffle"

        // ParamPerLeadLifetime is the per-lead lifetime limit config field.
        ParamPerLeadLifetime = "per_lead_lifetime"

        // ParamMsgIntervalHours is the message interval hours config field.
        ParamMsgIntervalHours = "msg_interval_hours"

        // ParamFollowupDelayDays is the follow-up delay days config field.
        ParamFollowupDelayDays = "followup_delay_days"

        // ParamFollowupVariant is the follow-up variant config field.
        ParamFollowupVariant = "followup_variant"

        // ParamColdAfter is the cold-after count config field.
        ParamColdAfter = "cold_after"

        // ParamRecontactDelayDays is the recontact delay days config field.
        ParamRecontactDelayDays = "recontact_delay_days"

        // ParamAutoBlock is the auto-block mode config field.
        ParamAutoBlock = "auto_block"

        // ParamDupCrossNiche is the duplicate cross-niche check config field.
        ParamDupCrossNiche = "dup_cross_niche"

        // ParamWAPreValidation is the WA pre-validation config field.
        ParamWAPreValidation = "wa_pre_validation"

        // ParamWAValidationMethod is the WA validation method config field.
        ParamWAValidationMethod = "wa_validation_method"

        // --- RPC handler param keys ---

        // ParamBusiness is the business name field (used in response screen actions).
        ParamBusiness = "business"

        // ParamOfferText is the offer text field inside lead response data.
        ParamOfferText = "offer_text"

        // ParamTrigger is the trigger source field inside lead response data.
        ParamTrigger = "trigger"

        // ParamClass is the response classification field inside lead response data.
        ParamClass = "class"

        // ParamPipeline is the pipeline description field inside conversion data.
        ParamPipeline = "pipeline"

        // ParamTimeTaken is the time taken field inside conversion data.
        ParamTimeTaken = "time_taken"

        // ParamTrophyCount is the trophy count field inside conversion data.
        ParamTrophyCount = "trophy_count"

        // ParamCursor is the cursor position field inside response data.
        ParamCursor = "cursor"

        // ParamType is the notification type discriminator field in notify params.
        // (Also used as the type discriminator field inside template groups.)

        // ParamSeverity is the severity level field in notify params.
        // (Also used as the severity field inside a template error entry.)

        // ParamData is the notification payload data map field in notify params.
        ParamData = "data"

        // ParamScreen is the screen identifier field in navigate/action params.
        ParamScreen = "screen"

        // --- Scrape screen param keys (remaining) ---

        // ParamDuplicates is the duplicates count field in scrape data.
        ParamDuplicates = "duplicates"

        // ParamScanning is the scanning-in-progress boolean field in scrape data.
        ParamScanning = "scanning"

        // ParamLastScrape is the last scrape timestamp field in scrape idle data.
        ParamLastScrape = "last_scrape"

        // ParamNextScrape is the next scrape countdown field in scrape idle data.
        ParamNextScrape = "next_scrape"

        // ParamError is the error message field in scrape error data.
        ParamError = "error"

        // ParamRetryIn is the retry countdown field in scrape error data.
        ParamRetryIn = "retry_in"

        // ParamThrottleIn is the throttle countdown field in scrape throttle data.
        ParamThrottleIn = "throttle_in"

        // ParamWANiches is the WA validation per-niche data list field.
        ParamWANiches = "wa_niches"

        // ParamWANicheName is the current WA validation niche name field.
        ParamWANicheName = "wa_niche_name"

        // ParamWATotal is the total WA validation count field.
        ParamWATotal = "wa_total"

        // ParamWAHas is the has-WA count field in WA validation data.
        ParamWAHas = "wa_has"

        // ParamWANot is the no-WA count field in WA validation data.
        ParamWANot = "wa_not"

        // ParamWAPending is the pending WA check count field in WA validation data.
        ParamWAPending = "wa_pending"

        // ParamWAPercent is the WA validation percentage field.
        ParamWAPercent = "wa_percent"

        // ParamWAEstimate is the WA validation time estimate field.
        ParamWAEstimate = "wa_estimate"

        // ParamJackpotLead is the high-value lead data map field for jackpot reveal.
        ParamJackpotLead = "jackpot_lead"

        // ParamBatchResults is the batch completion per-niche data list field.
        ParamBatchResults = "batch_results"

        // ParamBatchNextIn is the next batch countdown field.
        ParamBatchNextIn = "batch_next_in"

        // ParamQueueTotal is the total queue count across all niches field.
        ParamQueueTotal = "queue_total"

        // ParamFailedCount is the count of failed actions.
        ParamFailedCount = "failed_count"

        // --- Additional Shield rate limiting & work hours param keys ---

        // ParamRateLimitPerSlot is the per-slot rate limit display string.
        ParamRateLimitPerSlot = "rate_limit_per_slot"

        // ParamRateLimitPerDay is the per-day rate limit display string.
        ParamRateLimitPerDay = "rate_limit_per_day"

        // ParamRateLimitPerNumber is the per-number rate limit display string.
        ParamRateLimitPerNumber = "rate_limit_per_number"

        // ParamRateLimitPerLead is the per-lead rate limit display string.
        ParamRateLimitPerLead = "rate_limit_per_lead"

        // ParamTimezone is the timezone display string.
        ParamTimezone = "timezone"

        // ParamSendHours is the allowed sending hours display string.
        ParamSendHours = "send_hours"

        // ParamScrapeHours is the allowed scraping hours display string.
        ParamScrapeHours = "scrape_hours"

        // ParamNowInWorkHours is the "currently in work hours" indicator.
        ParamNowInWorkHours = "now_in_work_hours"

        // ParamSlotHistory is the slot event history list for slot detail view.
        ParamSlotHistory = "slot_history"

        // --- Monitor dashboard param keys ---

        // ParamConvRate is the conversion rate display string.
        ParamConvRate = "conv_rate"

        // ParamBestDay is the best performing day display string.
        ParamBestDay = "best_day"

        // ParamBestTimeStr is the best sending time display string.
        ParamBestTimeStr = "best_time"

        // ParamWANumCount is the WA number count field.
        ParamWANumCount = "wa_num_count"

        // ParamActiveSlotCount is the active slot count field.
        ParamActiveSlotCount = "active_slot_count"

        // ParamErrorSlot is the error slot description field.
        ParamErrorSlot = "error_slot"

        // ParamCurrentTime is the current time display string.
        ParamCurrentTime = "current_time"

        // ParamAppName is the application name field.
        ParamAppName = "app_name"

        // --- Anti-ban shield additional param keys ---

        // ParamDNCCount is the do-not-contact count field.
        ParamDNCCount = "dnc_count"

        // ParamIceBreakerVariants is the ice breaker variant count field.
        ParamIceBreakerVariants = "ice_breaker_variants"

        // ParamOfferVariants is the offer variant count field.
        ParamOfferVariants = "offer_variants"

        // ParamWorkHoursDuration is the work hours duration field.
        ParamWorkHoursDuration = "work_hours_duration"

        // ParamHealthRecoveryPts is the health recovery points field.
        ParamHealthRecoveryPts = "health_recovery_pts"

        // ParamTimezoneShort is the short timezone display string.
        ParamTimezoneShort = "timezone_short"

        // ParamTimezoneFull is the full timezone display string.
        ParamTimezoneFull = "timezone_full"

        // --- Boot onboarding additional param keys ---

        // ParamArrowCount is the arrow count field for marching workers.
        ParamArrowCount = "arrow_count"
)
