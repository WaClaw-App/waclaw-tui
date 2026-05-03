package protocol

// ScreenID identifies a distinct screen (page) in the WaClaw TUI.
//
// Screens are the top-level navigation targets. Each screen owns one or
// more StateID values that describe its internal sub-states. ScreenID is
// a string type so that it serialises naturally in JSON-RPC messages
// without requiring custom marshalers.
type ScreenID string

const (
	// ScreenBoot is the initial splash / loading screen shown on startup.
	ScreenBoot ScreenID = "boot"

	// ScreenLogin presents the WhatsApp QR-code authentication flow.
	ScreenLogin ScreenID = "login"

	// ScreenNicheSelect lets the user pick one or more target niches.
	ScreenNicheSelect ScreenID = "niche_select"

	// ScreenScrape displays real-time scraping progress and results.
	ScreenScrape ScreenID = "scrape"

	// ScreenLeadReview is where the user reviews and approves/rejects scraped leads.
	ScreenLeadReview ScreenID = "lead_review"

	// ScreenSend manages the outbound message sending pipeline.
	ScreenSend ScreenID = "send"

	// ScreenMonitor is the live dashboard for monitoring all activity.
	ScreenMonitor ScreenID = "monitor"

	// ScreenResponse shows incoming replies from leads and lets the user act on them.
	ScreenResponse ScreenID = "response"

	// ScreenLeadsDB is the searchable lead database / CRM view.
	ScreenLeadsDB ScreenID = "leads_db"

	// ScreenTemplateMgr manages message templates.
	ScreenTemplateMgr ScreenID = "template_mgr"

	// ScreenWorkers shows the per-niche worker pool status.
	ScreenWorkers ScreenID = "workers"

	// ScreenAntiBan configures anti-ban / health-score shielding.
	ScreenAntiBan ScreenID = "antiban"

	// ScreenSettings is the global application settings screen.
	ScreenSettings ScreenID = "settings"

	// ScreenGuardrail shows content-safety and guardrail configuration.
	ScreenGuardrail ScreenID = "guardrail"

	// ScreenCompose is the free-form message composition screen.
	ScreenCompose ScreenID = "compose"

	// ScreenHistory displays historical activity logs.
	ScreenHistory ScreenID = "history"

	// ScreenFollowUp manages follow-up scheduling and execution.
	ScreenFollowUp ScreenID = "followup"

	// ScreenLicense handles license key entry and validation.
	ScreenLicense ScreenID = "license"

	// ScreenNicheExplorer allows browsing and discovering niches.
	ScreenNicheExplorer ScreenID = "niche_explorer"

	// ScreenUpdate shows available updates and upgrade prompts.
	ScreenUpdate ScreenID = "update"
)

// AllScreens returns a slice containing every defined ScreenID, in the
// canonical declaration order above.
func AllScreens() []ScreenID {
	return []ScreenID{
		ScreenBoot,
		ScreenLogin,
		ScreenNicheSelect,
		ScreenScrape,
		ScreenLeadReview,
		ScreenSend,
		ScreenMonitor,
		ScreenResponse,
		ScreenLeadsDB,
		ScreenTemplateMgr,
		ScreenWorkers,
		ScreenAntiBan,
		ScreenSettings,
		ScreenGuardrail,
		ScreenCompose,
		ScreenHistory,
		ScreenFollowUp,
		ScreenLicense,
		ScreenNicheExplorer,
		ScreenUpdate,
	}
}

// IsValid reports whether the given ScreenID corresponds to a known screen.
func IsValid(id ScreenID) bool {
	switch id {
	case ScreenBoot,
		ScreenLogin,
		ScreenNicheSelect,
		ScreenScrape,
		ScreenLeadReview,
		ScreenSend,
		ScreenMonitor,
		ScreenResponse,
		ScreenLeadsDB,
		ScreenTemplateMgr,
		ScreenWorkers,
		ScreenAntiBan,
		ScreenSettings,
		ScreenGuardrail,
		ScreenCompose,
		ScreenHistory,
		ScreenFollowUp,
		ScreenLicense,
		ScreenNicheExplorer,
		ScreenUpdate:
		return true
	default:
		return false
	}
}
