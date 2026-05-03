// Package pipeline defines shared data models for the scrape, review, and send
// pipeline screens. Every struct maps directly onto the map[string]any params
// that the backend pushes via HandleUpdate, so all fields are populated from
// JSON-RPC update payloads rather than constructed in the TUI itself.
package pipeline

import (
        "time"

        "github.com/WaClaw-App/waclaw/internal/tui/util"
)

// NicheScrapeData holds per-niche scraping progress data.
//
// Pushed by the backend via HandleUpdate when scrape state is active.
// Each running niche worker produces one of these on every tick.
type NicheScrapeData struct {
        Name       string     // Niche identifier (e.g. "web_developer")
        Targets    string     // Comma-separated search targets (e.g. "cafe, gym, salon")
        Area       string     // Search area + radius (e.g. "kediri (15km)")
        Filter     string     // Active filter description (e.g. "tanpa website")
        Found      int64      // Total businesses found by scraper
        Qualified  int64      // Businesses that passed qualification filters
        Duplicates int64      // Duplicate businesses removed
        NewLeads   int64      // Net new leads entering the queue
        Skipped    int64      // Leads auto-skipped (low quality)
        WAHas      int64      // Leads with verified WA numbers
        WANot      int64      // Leads without WA (filtered out)
        WAPending  int64      // Leads pending WA validation
        WAChecked  int64      // Total WA checks completed
        Status     string     // Worker status: "scraping", "idle", "starting", "qualifying"
        NextIn     string     // Human-readable time until next batch (e.g. "2j 14m")
        LastBatch  string     // Human-readable time since last batch (e.g. "23 menit lalu")
        QueueCount int64      // Messages queued for this niche
        Leads      []LeadItem // Live lead items shown in the active scrape view
}

// HighValueLead holds data for a high-value lead revealed during scraping.
//
// Triggered when a lead scores 9+ (perfect rating, no website, active IG,
// many reviews).
type HighValueLead struct {
        Name     string  // Business name (e.g. "Grand Palace Hotel")
        Category string  // Business category (e.g. "hotel")
        Address  string  // Full address
        Rating   float64 // Google rating (e.g. 4.9)
        Reviews  int64   // Number of reviews (e.g. 312)
        HasWeb   bool    // Whether the business has a website
        HasIG    bool    // Whether the business has an active Instagram
        Score    float64 // Computed lead score (e.g. 9.5)
}

// BatchResult holds per-niche completion data shown in batch_complete state.
type BatchResult struct {
        Niche       string    // Niche name
        Found       int64     // Total found
        Qualified   int64     // Passed qualification
        NewLeads    int64     // Net new leads
        Duplicates  int64     // Duplicates removed
        Skipped     int64     // Auto-skipped
        CompletedAt time.Time // Timestamp when the batch finished
}

// WAValidationData holds WhatsApp pre-validation progress for a niche.
type WAValidationData struct {
        Niche     string  // Niche identifier
        Total     int64   // Total numbers to check
        Checked   int64   // Numbers already checked
        WAHas     int64   // Has WA — goes to send queue
        WANot     int64   // Not WA — marked, not sent
        WAPercent float64 // Validation progress (0.0-1.0)
        Estimate  string  // Estimated time to complete (e.g. "3 menit")
}

// LeadItem represents a single lead in the review or scrape queue.
type LeadItem struct {
        Index         int64   // Position in queue (1-based)
        Name          string  // Business name
        Category      string  // Business category
        Address       string  // Full address
        City          string  // City name (populated by backend)
        Rating        float64 // Google rating
        Reviews       int64   // Number of reviews
        HasWebsite    bool    // Whether the business has a website
        HasInstagram  bool    // Whether the business has an active Instagram
        HasWA         bool    // WA pre-validated
        PhotoCount    int64   // Number of Google photos
        Score         int64   // Lead score (0-10)
        Contacted     bool    // Whether previously contacted
        FollowUpCount int64   // Number of follow-ups sent
        MaxFollowUp   int64   // Max follow-ups allowed (3)
        HasWebInfo    bool    // Whether web search found info
        WebNoSite     bool    // No official website found
        WebNoMaps     bool    // No link on Google Maps
        WebNoSosmed   bool    // Not found on social media
        IsNew            bool    // Whether this lead just appeared in scrape view
        Qualified        bool    // Whether this lead passed qualification filters
        GoogleResults    []string // Web search results for this lead
        Potential        string  // Lead potential label (e.g. "high", "medium")
        ContactHistory   string  // Previous contact history summary
        FollowUpStatus   string  // Follow-up status description
}

// ReviewStats tracks review session counters.
type ReviewStats struct {
        Queued    int64 // Leads approved for sending
        Skipped   int64 // Leads skipped
        Blocked   int64 // Leads blocked
        Remaining int64 // Leads still in queue
}

// SlotInfo holds WhatsApp sender slot data for the send screen.
type SlotInfo struct {
        Number       string // Masked number (e.g. "0812-xxxx-3456")
        Status       string // "active", "cooldown", "down"
        SentThisHour int64  // Messages sent this hour
        HourLimit    int64  // Hourly send limit for this slot
        CooldownLeft string // Time until cooldown expires (e.g. "8m 12s")
        ReadyIn      string // Time until slot is ready (e.g. "3m 05s")
}

// SendQueueItem represents a single message in the send queue.
type SendQueueItem struct {
        Index        int64  // Position in queue (1-based)
        Name         string // Lead business name
        Niche        string // Niche name
        TemplateType string // Message type: "ice_breaker", "follow_up", "offer"
        TemplateVar  string // Template variant (e.g. "variant_2")
        Slot         string // WA slot assignment (e.g. "slot-1" or "1")
        Rotated      bool   // Whether variant was rotated
        Status       string // "sending", "waiting", "sent", "failed"
        NextAt       string // Time until this message sends (e.g. "11m 23s")
        IsActive     bool   // Whether this is the currently sending item
}

// SendStats tracks sending session counters.
type SendStats struct {
        RateHour      int64  // Messages sent this hour
        RateHourLimit int64  // Hourly rate limit
        DailySent     int64  // Messages sent today
        DailyLimit    int64  // Daily send limit
        DailyRespond  int64  // Responses received today
        DailyConvert  int64  // Conversions today
        SlotCount     int64  // Active WA slots
        NicheCount    int64  // Active niches
        QueueTotal    int64  // Total messages queued
        NextSendTime  string // Time of next scheduled send
        NextSendSlot  string // Slot that will send next
        Now           string // Current time (e.g. "21:47 wib")
}

// SendFailure holds data about a failed send operation.
type SendFailure struct {
        Name   string // Lead business name
        Reason string // Failure reason
        Hint   string // Suggestion for resolution
}

// ResponseInterrupt holds data about an incoming response during sending.
type ResponseInterrupt struct {
        Name    string // Lead business name
        Message string // Response text from the lead
}

// ---------------------------------------------------------------------------
// Status string constants — DRY replacements for raw string comparisons
// ---------------------------------------------------------------------------

// SlotStatus constants for WA slot status comparisons.
const (
        SlotStatusActive   = "active"
        SlotStatusCooldown = "cooldown"
        SlotStatusDown     = "down"
)

// QueueItemStatus constants for send queue item status comparisons.
const (
        QueueItemSending = "sending"
        QueueItemSent    = "sent"
        QueueItemWaiting = "waiting"
        QueueItemFailed  = "failed"
)

// WorkerStatus constants for niche worker status comparisons.
const (
        WorkerStatusScraping = "scraping"
        WorkerStatusIdle     = "idle"
        WorkerStatusStarting = "starting"
)

// ---------------------------------------------------------------------------
// Shared parsing helpers — used by scrape.go, review.go, send.go
// ---------------------------------------------------------------------------

// anyMapString extracts a string from a map[string]any, returning "" if missing.
// Delegates to the shared util.ToString for DRY — single source of truth.
func anyMapString(m map[string]any, key string) string {
        return util.ToString(m[key], "")
}

// anyMapInt extracts an int from a map[string]any, returning 0 if missing.
// Delegates to the shared util.ToInt for DRY — single source of truth.
func anyMapInt(m map[string]any, key string) int {
        return util.ToInt(m[key], 0)
}

// anyMapBool extracts a bool from a map[string]any, returning false if missing.
// Delegates to the shared util.ToBool for DRY — single source of truth.
func anyMapBool(m map[string]any, key string) bool {
        return util.ToBool(m[key], false)
}
