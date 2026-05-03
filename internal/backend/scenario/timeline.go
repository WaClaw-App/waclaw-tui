package scenario

import (
        "log"
        "sync"
        "time"

        "github.com/WaClaw-App/waclaw/pkg/protocol"
)

// Timeline orchestrates a scripted demo sequence that automatically drives
// the TUI through screens and injects mock data at timed intervals.
//
// When the demo starts, the timeline runs the boot sequence and then proceeds
// through the canonical screen flow, pushing notifications and updating data
// at appropriate moments. The user can still interact manually — the timeline
// only fires events when no user interaction has happened recently.
//
// Design: Each timeline step is a function that performs one action (navigate,
// notify, update). Steps are executed sequentially with delays between them.
type Timeline struct {
        engine     *Engine
        mu         sync.Mutex
        running    bool
        stopCh     chan struct{}
        cycleCount int // tracks how many full cycles have completed
}

// NewTimeline creates a new Timeline bound to the given engine.
func NewTimeline(engine *Engine) *Timeline {
        return &Timeline{
                engine: engine,
                stopCh: make(chan struct{}),
        }
}

// Step represents a single timed action in the demo sequence.
type Step struct {
        // Delay before executing this step.
        Delay time.Duration

        // Action is the function to execute.
        Action func()
}

// Start begins the demo timeline sequence.
// It runs in a background goroutine and returns immediately.
func (t *Timeline) Start() {
        t.mu.Lock()
        if t.running {
                t.mu.Unlock()
                return
        }
        t.running = true
        t.mu.Unlock()

        go t.run()
}

// Stop halts the timeline.
func (t *Timeline) Stop() {
        t.mu.Lock()
        defer t.mu.Unlock()
        if t.running {
                close(t.stopCh)
                t.running = false
        }
}

// run executes the demo sequence step by step, then loops back to the
// returning-user boot sequence. The demo cycles indefinitely until Stop()
// is called or the stopCh is closed.
func (t *Timeline) run() {
        for {
                steps := t.buildSequence()

                for _, step := range steps {
                        select {
                        case <-t.stopCh:
                                log.Println("[timeline] stopped")
                                return
                        case <-time.After(step.Delay):
                                step.Action()
                        }
                }

                t.cycleCount++
                log.Println("[timeline] demo cycle complete, restarting")
        }
}

// buildSequence constructs the full demo walkthrough steps.
//
// The sequence follows the canonical screen flow per doc/18-screen-flow.md:
// Boot → License → Validation → Login → Niche Select → Scrape → Review → Send →
// Monitor → Response → Leads DB → Templates → Workers →
// Anti-Ban → Settings → Compose → History →
// Follow-Up → Niche Explorer → Update
//
// Between screen transitions, notifications and data updates are
// injected to simulate realistic backend activity.
//
// On the first cycle, the boot screen shows BootFirstTime (per
// doc/18-screen-flow.md). On subsequent cycles, it shows BootReturning
// to match the documented returning-user flow.
func (t *Timeline) buildSequence() []Step {
        // First cycle: first-time boot. Subsequent cycles: returning-user boot.
        bootState := protocol.BootFirstTime
        if t.cycleCount > 0 {
                bootState = protocol.BootReturning
        }

        return []Step{
                // Boot sequence.
                {
                        Delay:  1 * time.Second,
                        Action: func() { log.Println("[timeline] boot sequence started") },
                },
                {
                        Delay:  2 * time.Second,
                        Action: func() {
                                t.engine.transitionTo(protocol.ScreenBoot, bootState)
                        },
                },

                // License gate per doc/18-screen-flow.md:
                // "LICENSE ── belum ada key? ──→ LICENSE INPUT"
                {
                        Delay:  2 * time.Second,
                        Action: func() {
                                t.engine.transitionTo(protocol.ScreenLicense, protocol.LicenseInput)
                        },
                },

                // Validation gate per doc/18-screen-flow.md:
                // "VALIDATION ── config missing? ──→ CONFIG SETUP"
                {
                        Delay:  2 * time.Second,
                        Action: func() {
                                t.engine.transitionTo(protocol.ScreenGuardrail, protocol.ValidationClean)
                        },
                },

                // Transition to Login with QR code.
                {
                        Delay:  2 * time.Second,
                        Action: func() {
                                t.engine.transitionTo(protocol.ScreenLogin, protocol.LoginQRWaiting)
                        },
                },
                {
                        Delay:  2 * time.Second,
                        Action: func() {
                                t.engine.UpdateState(protocol.LoginQRScanned)
                                t.engine.PushScreenUpdate(map[string]any{
                                        "state":         string(protocol.LoginQRScanned),
                                        "contact_count": 847,
                                        "filled_slots":  1,
                                        "phone_numbers": []any{"0812-xxxx-3456"},
                                })
                        },
                },
                {
                        Delay:  1 * time.Second,
                        Action: func() {
                                t.engine.UpdateState(protocol.LoginSuccess)
                                t.engine.PushScreenUpdate(map[string]any{
                                        "state":         string(protocol.LoginSuccess),
                                        "contact_count": 847,
                                        "filled_slots":  1,
                                        "phone_numbers": []any{"0812-xxxx-3456"},
                                })
                        },
                },

                // Navigate to Niche Select.
                {
                        Delay:  2 * time.Second,
                        Action: func() {
                                t.engine.transitionTo(protocol.ScreenNicheSelect, protocol.NicheList)
                        },
                },

                // Start scraping.
                {
                        Delay:  3 * time.Second,
                        Action: func() {
                                t.engine.transitionTo(protocol.ScreenScrape, protocol.ScrapeActive)
                        },
                },
                {
                        Delay:  2 * time.Second,
                        Action: func() {
                                t.engine.PushScreenUpdate(map[string]any{
                                        "progress_pct":     35.0,
                                        "leads_found":      18,
                                        "leads_validated":  12,
                                })
                        },
                },
                {
                        Delay:  2 * time.Second,
                        Action: func() {
                                t.engine.PushScreenUpdate(map[string]any{
                                        "progress_pct":     70.0,
                                        "leads_found":      34,
                                        "leads_validated":  28,
                                })
                        },
                },
                // Scrape complete notification.
                {
                        Delay:  2 * time.Second,
                        Action: func() {
                                t.engine.UpdateState(protocol.ScrapeBatchComplete)
                                t.engine.PushNotification(
                                        protocol.NotifScrapeComplete,
                                        protocol.SeverityPositive,
                                        map[string]any{
                                                "message":    "Scraping selesai — 47 lead ditemukan",
                                                "niche":      "kuliner",
                                                "leads":      47,
                                                "high_value": 3,
                                        },
                                )
                        },
                },

                // Navigate to Lead Review.
                {
                        Delay:  3 * time.Second,
                        Action: func() {
                                t.engine.transitionTo(protocol.ScreenLeadReview, protocol.ReviewReviewing)
                        },
                },

                // Navigate to Send.
                {
                        Delay:  2 * time.Second,
                        Action: func() {
                                t.engine.transitionTo(protocol.ScreenSend, protocol.SendActive)
                        },
                },
                // Batch send complete notification.
                {
                        Delay:  3 * time.Second,
                        Action: func() {
                                t.engine.PushNotification(
                                        protocol.NotifBatchSendComplete,
                                        protocol.SeverityNeutral,
                                        map[string]any{
                                                "message": "Batch selesai — 25 pesan terkirim",
                                                "count":   25,
                                        },
                                )
                        },
                },

                // Navigate to Monitor (dashboard).
                {
                        Delay:  2 * time.Second,
                        Action: func() {
                                t.engine.transitionTo(protocol.ScreenMonitor, protocol.MonitorLiveDashboard)
                        },
                },
                // Streak milestone notification.
                {
                        Delay:  4 * time.Second,
                        Action: func() {
                                t.engine.PushNotification(
                                        protocol.NotifStreakMilestone,
                                        protocol.SeverityInformative,
                                        map[string]any{
                                                "message": "7 hari streak! Terus gas!",
                                                "days":    7,
                                        },
                                )
                        },
                },

                // Navigate to Response screen with a positive response.
                {
                        Delay:  3 * time.Second,
                        Action: func() {
                                t.engine.transitionTo(protocol.ScreenResponse, protocol.ResponsePositive)
                        },
                },
                // Response received notification.
                {
                        Delay:  3 * time.Second,
                        Action: func() {
                                t.engine.PushNotification(
                                        protocol.NotifResponseReceived,
                                        protocol.SeverityPositive,
                                        map[string]any{
                                                "message": "Response masuk dari Warung Padang Sederhana",
                                                "lead_id": "lead-0001",
                                        },
                                )
                        },
                },

                // Navigate through remaining screens quickly.
                {
                        Delay:  3 * time.Second,
                        Action: func() {
                                t.engine.transitionTo(protocol.ScreenLeadsDB, protocol.LeadsList)
                        },
                },
                {
                        Delay:  2 * time.Second,
                        Action: func() {
                                t.engine.transitionTo(protocol.ScreenTemplateMgr, protocol.TemplateList)
                        },
                },
                {
                        Delay:  2 * time.Second,
                        Action: func() {
                                t.engine.transitionTo(protocol.ScreenWorkers, protocol.WorkersOverview)
                        },
                },
                {
                        Delay:  2 * time.Second,
                        Action: func() {
                                t.engine.transitionTo(protocol.ScreenAntiBan, protocol.ShieldOverview)
                        },
                },
                // WA flag notification (critical).
                {
                        Delay:  2 * time.Second,
                        Action: func() {
                                t.engine.PushNotification(
                                        protocol.NotifWAFlag,
                                        protocol.SeverityCritical,
                                        map[string]any{
                                                "message": "Slot 2 ditandai WhatsApp — auto-pause",
                                                "slot_id": "slot_2",
                                        },
                                )
                        },
                },
                {
                        Delay:  2 * time.Second,
                        Action: func() {
                                t.engine.transitionTo(protocol.ScreenSettings, protocol.SettingsOverview)
                        },
                },
                {
                        Delay:  2 * time.Second,
                        Action: func() {
                                t.engine.transitionTo(protocol.ScreenGuardrail, protocol.ValidationClean)
                        },
                },
                {
                        Delay:  2 * time.Second,
                        Action: func() {
                                t.engine.transitionTo(protocol.ScreenCompose, protocol.ComposeDraft)
                        },
                },
                {
                        Delay:  2 * time.Second,
                        Action: func() {
                                t.engine.transitionTo(protocol.ScreenHistory, protocol.HistoryToday)
                        },
                },
                {
                        Delay:  2 * time.Second,
                        Action: func() {
                                t.engine.transitionTo(protocol.ScreenFollowUp, protocol.FollowUpDashboard)
                        },
                },
                {
                        Delay:  2 * time.Second,
                        Action: func() {
                                t.engine.transitionTo(protocol.ScreenLicense, protocol.LicenseInput)
                        },
                },
                {
                        Delay:  2 * time.Second,
                        Action: func() {
                                t.engine.transitionTo(protocol.ScreenNicheExplorer, protocol.ExplorerBrowse)
                        },
                },
                {
                        Delay:  2 * time.Second,
                        Action: func() {
                                t.engine.transitionTo(protocol.ScreenUpdate, protocol.UpdateAvailable)
                        },
                },

                // Loop back to boot (returning user).
                {
                        Delay:  5 * time.Second,
                        Action: func() {
                                t.engine.transitionTo(protocol.ScreenBoot, protocol.BootReturning)
                        },
                },
        }
}
