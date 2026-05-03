// Package scenario implements the demo backend's state machine that drives
// the TUI through a scripted sequence of screen transitions and data updates.
//
// The ScenarioEngine owns the current screen/state and decides what happens
// next in response to TUI events. It communicates back to the TUI through
// the RPCPusher interface defined in the engine package.
//
// Design decisions:
//   - The canonical flow is defined once in canonicalFlow and all navigation
//     (forward, backward, key-based) derives from it — single source of truth.
//   - Mock data generation is delegated to mock.go for separation of concerns.
//   - Timeline orchestration is delegated to timeline.go for scripted demos.
package scenario

import (
        "fmt"
        "log"
        "sync"

        "github.com/WaClaw-App/waclaw/internal/backend/engine"
        "github.com/WaClaw-App/waclaw/internal/tui/util"
        "github.com/WaClaw-App/waclaw/pkg/protocol"
)

// flowStep is a single entry in the canonical screen flow.
type flowStep struct {
        screen protocol.ScreenID
        state  protocol.StateID
}

// canonicalFlow defines the canonical demo walkthrough order — the single
// source of truth. All navigation methods (forward, backward, key-based
// transition table) derive from this slice. DRY: one definition, many uses.
//
// The flow follows doc/18-screen-flow.md for first-time users:
//   BOOT → LICENSE → VALIDATION → LOGIN → NICHE SELECT → SCRAPE → ...
// License and Validation gates are included as intermediate steps for the
// first-time user demo, matching the documented screen flow exactly.
var canonicalFlow = []flowStep{
        {protocol.ScreenBoot, protocol.BootFirstTime},
        {protocol.ScreenLicense, protocol.LicenseInput},
        {protocol.ScreenGuardrail, protocol.ValidationClean},
        {protocol.ScreenLogin, protocol.LoginQRWaiting},
        {protocol.ScreenNicheSelect, protocol.NicheList},
        {protocol.ScreenScrape, protocol.ScrapeActive},
        {protocol.ScreenLeadReview, protocol.ReviewReviewing},
        {protocol.ScreenSend, protocol.SendActive},
        {protocol.ScreenMonitor, protocol.MonitorLiveDashboard},
        {protocol.ScreenResponse, protocol.ResponsePositive},
        {protocol.ScreenLeadsDB, protocol.LeadsList},
        {protocol.ScreenTemplateMgr, protocol.TemplateList},
        {protocol.ScreenWorkers, protocol.WorkersOverview},
        {protocol.ScreenAntiBan, protocol.ShieldOverview},
        {protocol.ScreenSettings, protocol.SettingsOverview},
        {protocol.ScreenCompose, protocol.ComposeDraft},
        {protocol.ScreenHistory, protocol.HistoryToday},
        {protocol.ScreenFollowUp, protocol.FollowUpDashboard},
        {protocol.ScreenNicheExplorer, protocol.ExplorerBrowse},
        {protocol.ScreenUpdate, protocol.UpdateAvailable},
}

// keyTransitionTable maps (current screen) → (next step index) for "enter"
// key presses. Derived from canonicalFlow so we don't repeat screen/state
// pairs. Built once at package init time.
var keyTransitionTable map[protocol.ScreenID]int

func init() {
        keyTransitionTable = make(map[protocol.ScreenID]int, len(canonicalFlow))
        for i, step := range canonicalFlow {
                // "enter" on screen[i] advances to screen[i+1] (wrapping around).
                next := i + 1
                if next >= len(canonicalFlow) {
                        next = 0 // wrap: Update → Boot
                }
                keyTransitionTable[step.screen] = next
        }
}

// Engine implements engine.ScenarioEngine for the demo backend.
type Engine struct {
        mu     sync.RWMutex
        pusher engine.RPCPusher

        // currentScreen and currentState track where the demo is in the flow.
        currentScreen protocol.ScreenID
        currentState  protocol.StateID

        // currentFlowIndex caches the position in canonicalFlow for O(1) lookups.
        currentFlowIndex int

        // timeline is the scripted demo sequence runner.
        timeline *Timeline

        // mock provides realistic data for the demo.
        mock *MockData

        // lastExplorerCategory tracks the most recently selected explorer category
        // name so that navigateParams can use it instead of a hardcoded value.
        lastExplorerCategory string
}

// NewEngine creates a new demo scenario engine.
func NewEngine(pusher engine.RPCPusher) *Engine {
        e := &Engine{
                pusher:           pusher,
                currentScreen:    protocol.ScreenBoot,
                currentState:     protocol.BootFirstTime,
                currentFlowIndex: 0,
                mock:             NewMockData(),
        }
        e.timeline = NewTimeline(e)
        return e
}

// HandleKeyPress implements engine.ScenarioEngine.
func (e *Engine) HandleKeyPress(evt protocol.KeyPressEvent) error {
        e.mu.Lock()
        defer e.mu.Unlock()

        log.Printf("[scenario] key_press: key=%s screen=%s state=%s",
                evt.Key, evt.Screen, evt.State)

        if evt.Key != "enter" {
                return nil // only "enter" drives the demo flow via key press
        }

        nextIdx, ok := keyTransitionTable[e.currentScreen]
        if !ok {
                return nil
        }

        step := canonicalFlow[nextIdx]
        // Special case: wrapping back to boot uses BootReturning, not BootFirstTime.
        state := step.state
        if step.screen == protocol.ScreenBoot && e.currentScreen != protocol.ScreenBoot {
                state = protocol.BootReturning
        }

        return e.transitionTo(step.screen, state)
}

// HandleAction implements engine.ScenarioEngine.
func (e *Engine) HandleAction(evt protocol.ActionEvent) error {
        e.mu.Lock()
        defer e.mu.Unlock()

        log.Printf("[scenario] action: action=%s screen=%s", evt.Action, evt.Screen)

        // Phase 3A specific action routing — maps boot/login action strings
        // to the appropriate screen transitions. These actions come from the
        // TUI's onboarding key handlers (actionAndNavigate / publishAction).
        switch evt.Action {
        // Boot screen actions
        case "boot_login":
                return e.transitionTo(protocol.ScreenLogin, protocol.LoginQRWaiting)
        case "boot_niche":
                return e.transitionTo(protocol.ScreenNicheSelect, protocol.NicheList)
        case "boot_gas":
                return e.transitionTo(protocol.ScreenScrape, protocol.ScrapeActive)
        case "boot_dashboard":
                return e.transitionTo(protocol.ScreenMonitor, protocol.MonitorLiveDashboard)
        case "boot_view_responses":
                return e.transitionTo(protocol.ScreenResponse, protocol.ResponsePositive)
        case "boot_relogin":
                return e.transitionTo(protocol.ScreenLogin, protocol.LoginQRWaiting)
        case "boot_view_error":
                return e.transitionTo(protocol.ScreenGuardrail, protocol.ValidationErrors)
        case "boot_enter_license":
                return e.transitionTo(protocol.ScreenLicense, protocol.LicenseInput)
        case "boot_buy_license":
                return e.transitionTo(protocol.ScreenLicense, protocol.LicenseInput)
        case "boot_disconnect_other":
                // Force-disconnect the other device; resume returning boot.
                return e.transitionTo(protocol.ScreenBoot, protocol.BootReturning)
        case "boot_exit":
                return nil // handled by TUI app quit

        // Login screen actions
        case "login_skip":
                return e.transitionTo(protocol.ScreenNicheSelect, protocol.NicheList)
        case "login_add_slot", "login_add_another", "login_add_number":
                // Re-render current login state (slot count change handled by update).
                return nil
        case "login_enough":
                return e.transitionTo(protocol.ScreenNicheSelect, protocol.NicheList)
        case "login_success_continue":
                return e.transitionTo(protocol.ScreenNicheSelect, protocol.NicheList)
        case "login_gas":
                return e.transitionTo(protocol.ScreenMonitor, protocol.MonitorLiveDashboard)
        case "login_later":
                return e.transitionTo(protocol.ScreenMonitor, protocol.MonitorLiveDashboard)
        case "login_retry":
                return e.transitionTo(protocol.ScreenLogin, protocol.LoginQRWaiting)
        case "login_change_slot":
                return nil // slot change handled by update
        case "login_back":
                return e.transitionTo(protocol.ScreenBoot, protocol.BootFirstTime)

        // Update screen actions — backend drives state transitions
        case "start_download":
                return e.transitionTo(protocol.ScreenUpdate, protocol.UpdateDownloading)
        case "remind_later", "skip_update", "skip_restart":
                // Go back to previous screen (usually monitor)
                screen, state, ok := e.goBack()
                if ok {
                        return e.transitionTo(screen, state)
                }
        case "cancel_download":
                return e.transitionTo(protocol.ScreenUpdate, protocol.UpdateAvailable)
        case "restart_now":
                // Restart simulated — wrap back to boot
                return e.transitionTo(protocol.ScreenBoot, protocol.BootReturning)
        case "restart_later":
                screen, state, ok := e.goBack()
                if ok {
                        return e.transitionTo(screen, state)
                }
        case "buy_license":
                return e.transitionTo(protocol.ScreenUpdate, protocol.UpgradeLicenseInput)
        case "view_upgrade_details", "stay_v1":
                // Stay on upgrade screen or go back
                return nil
        case "validate_license":
                // Simulate valid license — transition back
                return e.transitionTo(protocol.ScreenBoot, protocol.BootReturning)
        case "cancel_license_input":
                return e.transitionTo(protocol.ScreenUpdate, protocol.UpgradeAvailable)
        case "renew_v1":
                return e.transitionTo(protocol.ScreenLicense, protocol.LicenseInput)
        case "upgrade_v2":
                return e.transitionTo(protocol.ScreenUpdate, protocol.UpgradeLicenseInput)
        case "enter_new_license":
                return e.transitionTo(protocol.ScreenLicense, protocol.LicenseInput)
        case "exit_expired":
                return nil

        // Niche Select screen actions
        case string(protocol.ActionNicheProceed):
                // User confirmed niche selection — transition to scrape.
                return e.transitionTo(protocol.ScreenScrape, protocol.ScrapeActive)
        case string(protocol.ActionNicheCustom):
                // User wants custom niche — stay on niche select in custom state.
                return e.pushUpdate(protocol.ScreenNicheSelect, map[string]any{
                        "state": string(protocol.NicheCustom),
                })
        case string(protocol.ActionNicheBack):
                // User hit back — return to niche list state.
                return e.transitionTo(protocol.ScreenNicheSelect, protocol.NicheList)
        case string(protocol.ActionNicheReturnList):
                return e.transitionTo(protocol.ScreenNicheSelect, protocol.NicheList)
        case string(protocol.ActionNicheScrape):
                return e.transitionTo(protocol.ScreenScrape, protocol.ScrapeActive)
        case string(protocol.ActionNicheEditFilter):
                // Stay on niche select, push filter data.
                return e.pushUpdate(protocol.ScreenNicheSelect, map[string]any{
                        "state": string(protocol.NicheEditFilters),
                        "filters": e.mock.ExplorerFilters(),
                        "areas":  e.mock.ExplorerAreas(),
                })
        case string(protocol.ActionNicheOpenFile), string(protocol.ActionNicheShowExample):
                // File operations handled by real backend; demo just logs.
                log.Printf("[scenario] niche action %s (demo no-op)", evt.Action)
                return nil
        case string(protocol.ActionNicheReload):
                // Reload re-sends the current screen data.
                return e.transitionTo(protocol.ScreenNicheSelect, protocol.NicheList)

        // Niche Explorer screen actions
        case string(protocol.ActionExplorerDetail):
                if name, ok := evt.Params["category"].(string); ok && name != "" {
                        e.lastExplorerCategory = name
                }
                return e.transitionTo(protocol.ScreenNicheExplorer, protocol.ExplorerCategoryDetail)
        case string(protocol.ActionExplorerBack):
                return e.transitionTo(protocol.ScreenNicheExplorer, protocol.ExplorerBrowse)
        case string(protocol.ActionExplorerGenerate):
                catName := "kuliner"
                if name, ok := evt.Params["category"].(string); ok && name != "" {
                        catName = name
                }
                e.lastExplorerCategory = catName
                slug := util.Slugify(catName)
                return e.pushUpdate(protocol.ScreenNicheExplorer, map[string]any{
                        "state":         string(protocol.ExplorerGenerating),
                        "category_name": catName,
                        "folder_slug":   slug,
                        "gen_files":     e.mock.ExplorerGenFiles(),
                        "gen_progress":  0.0,
                })
        case string(protocol.ActionExplorerSearch):
                // Search results come via HandleUpdate — demo just logs.
                log.Printf("[scenario] explorer search: %v", evt.Params)
                return nil
        case string(protocol.ActionExplorerEdit):
                log.Printf("[scenario] explorer edit (demo no-op)")
                return nil
        case string(protocol.ActionExplorerAddArea):
                log.Printf("[scenario] explorer add_area (demo no-op)")
                return nil
        case string(protocol.ActionExplorerCancel):
                return e.transitionTo(protocol.ScreenNicheExplorer, protocol.ExplorerBrowse)
        case string(protocol.ActionExplorerUse):
                return e.transitionTo(protocol.ScreenScrape, protocol.ScrapeActive)
        case string(protocol.ActionExplorerEditConfig):
                log.Printf("[scenario] explorer edit_config (demo no-op)")
                return nil
        case string(protocol.ActionExplorerViewTpl):
                return e.transitionTo(protocol.ScreenTemplateMgr, protocol.TemplateList)

        // Generic actions from other screens
        case "confirm", "select", "next":
                screen, state, ok := e.advanceFlow()
                if ok {
                        return e.transitionTo(screen, state)
                }
        case "back", "cancel":
                screen, state, ok := e.goBack()
                if ok {
                        return e.transitionTo(screen, state)
                }
        case "toggle":
                if e.pusher != nil {
                        return e.pusher.PushUpdate(e.currentScreen, map[string]any{
                                "action": evt.Action,
                                "params": evt.Params,
                        })
                }
                return nil

        // Comms screen actions
        case string(protocol.ActionComposeSend):
                // Compose sent — transition to the previous screen (response)
                return e.transitionTo(protocol.ScreenResponse, protocol.ResponsePositive)
        case string(protocol.ActionFollowUpAutoAll):
                // Auto-approve all — transition to sending state
                return e.transitionTo(protocol.ScreenFollowUp, protocol.FollowUpSending)
        case string(protocol.ActionFollowUpSkipWait):
                return nil // handled by update
        case string(protocol.ActionFollowUpPause):
                return e.transitionTo(protocol.ScreenFollowUp, protocol.FollowUpDashboard)
        case string(protocol.ActionFollowUpSendFinal):
                return nil // handled by update
        case string(protocol.ActionFollowUpArchiveCold):
                return e.transitionTo(protocol.ScreenFollowUp, protocol.FollowUpDashboard)
        case string(protocol.ActionFollowUpRecontact), string(protocol.ActionFollowUpRecontactAll):
                return e.transitionTo(protocol.ScreenFollowUp, protocol.FollowUpSending)
        case string(protocol.ActionHistoryPrevDay), string(protocol.ActionHistoryDayDetail):
                return nil // handled by update
        }

        return nil
}

// HandleRequest implements engine.ScenarioEngine.
func (e *Engine) HandleRequest(evt protocol.RequestEvent) (any, error) {
        e.mu.RLock()
        defer e.mu.RUnlock()

        log.Printf("[scenario] request: type=%s screen=%s", evt.Type, evt.Screen)

        switch evt.Type {
        case "fetch_leads":
                return e.mock.Leads(10), nil
        case "get_stats":
                return e.mock.Stats(), nil
        case "load_template":
                return e.mock.Template(), nil
        case "get_workers":
                return e.mock.Workers(), nil
        case "get_shield":
                return e.mock.ShieldData(), nil
        case "get_config":
                return e.mock.ConfigData(), nil
        case "get_validation":
                return e.mock.ValidationData(), nil
        case "get_niches":
                return e.mock.Niches(), nil
        case "get_notifications":
                return e.mock.Notifications(), nil
        case "get_state":
                return e.StateSnapshot(), nil
        default:
                return map[string]any{"status": "ok"}, nil
        }
}

// CurrentScreen implements engine.ScenarioEngine.
func (e *Engine) CurrentScreen() protocol.ScreenID {
        e.mu.RLock()
        defer e.mu.RUnlock()
        return e.currentScreen
}

// CurrentState implements engine.ScenarioEngine.
func (e *Engine) CurrentState() protocol.StateID {
        e.mu.RLock()
        defer e.mu.RUnlock()
        return e.currentState
}

// StateSnapshot implements engine.ScenarioEngine.
func (e *Engine) StateSnapshot() map[string]any {
        e.mu.RLock()
        defer e.mu.RUnlock()

        return map[string]any{
                "screen": string(e.currentScreen),
                "state":  string(e.currentState),
                "stats":  e.mock.Stats(),
                "leads":  e.mock.Leads(5),
        }
}

// SetPusher updates the RPC pusher (used during initialization).
func (e *Engine) SetPusher(p engine.RPCPusher) {
        e.mu.Lock()
        defer e.mu.Unlock()
        e.pusher = p
}

// Start begins the demo timeline sequence.
func (e *Engine) Start() {
        e.timeline.Start()
}

// transitionTo moves to a new screen/state and pushes the navigate command.
// It also updates currentFlowIndex to match the new position in canonicalFlow.
func (e *Engine) transitionTo(screen protocol.ScreenID, state protocol.StateID) error {
        e.currentScreen = screen
        e.currentState = state
        e.currentFlowIndex = flowIndex(screen)

        if e.pusher == nil {
                return nil // pusher not wired yet (initialization phase)
        }

        return e.pusher.PushNavigate(screen, state, e.navigateParams(screen, state))
}

// flowIndex returns the index of the given screen in canonicalFlow, or 0.
func flowIndex(screen protocol.ScreenID) int {
        for i, step := range canonicalFlow {
                if step.screen == screen {
                        return i
                }
        }
        return 0
}

// navigateParams builds the navigation parameters for a screen transition.
// Every key placed here must match what the corresponding TUI screen model
// reads in its HandleNavigate / applyXxxParams method — this is the
// backend→TUI data contract.
func (e *Engine) navigateParams(screen protocol.ScreenID, state protocol.StateID) map[string]any {
        params := map[string]any{
                "state": string(state),
        }

        // Attach screen-specific mock data.
        switch screen {
        case protocol.ScreenBoot:
                returning := state != protocol.BootFirstTime
                merge(params, e.mock.BootData(returning))

                // Add variant-specific fields on top of the base boot data.
                switch state {
                case protocol.BootReturningConfigError:
                        merge(params, e.mock.BootConfigErrorData())
                case protocol.BootReturningLicenseExpired:
                        merge(params, e.mock.BootLicenseExpiredData())
                case protocol.BootReturningDeviceConflict:
                        merge(params, e.mock.BootDeviceConflictData())
                case protocol.BootReturningResponse:
                        merge(params, e.mock.BootResponseData())
                case protocol.BootReturningError:
                        merge(params, e.mock.BootDisconnectData())
                }

        case protocol.ScreenLogin:
                switch state {
                case protocol.LoginQRWaiting:
                        merge(params, e.mock.LoginData())
                case protocol.LoginQRScanned:
                        merge(params, e.mock.LoginScannedData())
                case protocol.LoginSuccess:
                        merge(params, e.mock.LoginSuccessData())
                case protocol.LoginExpired:
                        merge(params, e.mock.LoginExpiredData())
                case protocol.LoginFailed:
                        merge(params, e.mock.LoginFailedData())
                }

        case protocol.ScreenNicheSelect:
                params["niches"] = e.mock.NicheSelectItems()
        case protocol.ScreenNicheExplorer:
                params["categories"] = e.mock.ExplorerCategories()
                // When navigating to detail/generating/generated sub-states,
                // include the category name so the TUI can render titles correctly.
                // Use lastExplorerCategory (set from action params) with fallback.
                catName := e.lastExplorerCategory
                if catName == "" {
                        catName = "kuliner"
                }
                switch state {
                case protocol.ExplorerCategoryDetail, protocol.ExplorerGenerating, protocol.ExplorerGenerated:
                        params["category_name"] = catName
                }
                // When navigating to the generated state, include the file list.
                if state == protocol.ExplorerGenerated {
                        params["gen_niche_name"] = catName
                        params["gen_files_done"] = e.mock.ExplorerGenFilesDone()
                }
                if state == protocol.ExplorerCategoryDetail {
                        params["sources"] = e.mock.ExplorerSources()
                        params["areas"] = e.mock.ExplorerAreas()
                        params["filters"] = e.mock.ExplorerFilters()
                        params["templates"] = e.mock.ExplorerTemplates()
                        // Backend sends folder_slug so TUI can display the correct folder path
                        // without relying on its own slugify() approximation.
                        params["folder_slug"] = util.Slugify(catName)
                }
        case protocol.ScreenScrape:
                params["progress"] = e.mock.ScrapeProgress()
        case protocol.ScreenMonitor:
                merge(params, e.mock.MonitorData())
                // Override state-specific fields
                if state == protocol.MonitorError {
                        params["error_slot"] = "slot-2" // Demo: slot-2 disconnected
                }
                if state == protocol.MonitorNight {
                        params["current_time"] = "22:15"
                        params["work_hours"] = "09:00-17:00 wib"
                }
        case protocol.ScreenResponse:
                merge(params, e.mock.ResponseScreenData())
        case protocol.ScreenLeadsDB:
                params["leads"] = e.mock.Leads(20)
        case protocol.ScreenWorkers:
                params["workers"] = e.mock.Workers()
        case protocol.ScreenAntiBan:
                merge(params, e.mock.ShieldData())
                merge(params, e.mock.ShieldConfigData())
        case protocol.ScreenSend:
                params["work_hours"] = "09:00-17:00 wib"
        case protocol.ScreenSettings:
                merge(params, e.mock.ConfigSettingsData())
        case protocol.ScreenGuardrail:
                params["results"] = e.mock.ValidationData()
        case protocol.ScreenLicense:
                params["key_prefix"] = "WACL"
                params["key_format"] = "WACL-XXXX-XXXX-XXXX-XXXX"
        case protocol.ScreenUpdate:
                isMajor := state == protocol.UpgradeAvailable || state == protocol.LicenseExpiredWithUpgrade
                merge(params, e.mock.UpdateData(isMajor))
                switch state {
                case protocol.UpdateDownloading:
                        merge(params, e.mock.UpdateDownloadProgress(0.35))
                case protocol.UpdateReady:
                        merge(params, e.mock.UpdateReadyData())
                case protocol.LicenseExpiredWithUpgrade:
                        merge(params, e.mock.UpdateExpiredData())
                }
        case protocol.ScreenCompose:
                merge(params, e.mock.ComposeData())

        case protocol.ScreenHistory:
                merge(params, e.mock.HistoryData())

        case protocol.ScreenFollowUp:
                merge(params, e.mock.FollowUpData())
        }

        return params
}

// merge copies every key from src into dst (last-write-wins).
func merge(dst, src map[string]any) {
        for k, v := range src {
                dst[k] = v
        }
}

// pushUpdate is a convenience method to push an update command for a specific
// screen without changing the canonical flow position.
func (e *Engine) pushUpdate(screen protocol.ScreenID, params map[string]any) error {
        if e.pusher == nil {
                return nil
        }
        return e.pusher.PushUpdate(screen, params)
}

// advanceFlow moves forward one step in the canonical flow.
func (e *Engine) advanceFlow() (protocol.ScreenID, protocol.StateID, bool) {
        next := e.currentFlowIndex + 1
        if next >= len(canonicalFlow) {
                // Wrap around to boot with "returning" state.
                return canonicalFlow[0].screen, protocol.BootReturning, true
        }
        return canonicalFlow[next].screen, canonicalFlow[next].state, true
}

// goBack moves backward one step in the canonical flow.
func (e *Engine) goBack() (protocol.ScreenID, protocol.StateID, bool) {
        prev := e.currentFlowIndex - 1
        if prev < 0 {
                prev = 0
        }
        return canonicalFlow[prev].screen, canonicalFlow[prev].state, true
}

// PushNotification is a convenience method for the timeline to push
// notifications to the TUI without going through the flow table.
func (e *Engine) PushNotification(notifType protocol.NotificationType, severity protocol.Severity, data map[string]any) error {
        e.mu.RLock()
        defer e.mu.RUnlock()
        if e.pusher == nil {
                return nil
        }
        return e.pusher.PushNotify(notifType, severity, data)
}

// PushScreenUpdate is a convenience method for the timeline to push
// data updates for the current screen.
func (e *Engine) PushScreenUpdate(params map[string]any) error {
        e.mu.RLock()
        defer e.mu.RUnlock()
        if e.pusher == nil {
                return nil
        }
        return e.pusher.PushUpdate(e.currentScreen, params)
}

// UpdateState changes the internal state without pushing a navigate command.
// Used by the timeline for sub-state changes within the same screen.
func (e *Engine) UpdateState(state protocol.StateID) {
        e.mu.Lock()
        defer e.mu.Unlock()
        e.currentState = state
}

// FormatState returns a human-readable representation of the current state.
func (e *Engine) FormatState() string {
        e.mu.RLock()
        defer e.mu.RUnlock()
        return fmt.Sprintf("screen=%s state=%s", e.currentScreen, e.currentState)
}
