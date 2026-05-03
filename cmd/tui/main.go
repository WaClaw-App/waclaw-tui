// Package main provides the TUI binary entry point for WaClaw.
//
// It bootstraps the bubbletea program, loads the user theme and config
// (falling back to defaults), registers all 20 screens, initialises the
// RPC client for backend communication over stdio, and starts the event loop.
//
// Phase 5 integration: all screens are registered here so the router can
// navigate to any screen driven by the backend scenario engine. Screen
// registration is the single wiring point — after this, the App, Router,
// Bus, and RPC client handle the rest.
package main

import (
        "fmt"
        "os"

        "github.com/WaClaw-App/waclaw/internal/tui"
        "github.com/WaClaw-App/waclaw/internal/tui/rpc"
        "github.com/WaClaw-App/waclaw/internal/tui/screen/comms"
        "github.com/WaClaw-App/waclaw/internal/tui/screen/data"
        "github.com/WaClaw-App/waclaw/internal/tui/screen/infra"
        "github.com/WaClaw-App/waclaw/internal/tui/screen/license"
        "github.com/WaClaw-App/waclaw/internal/tui/screen/monitor"
        "github.com/WaClaw-App/waclaw/internal/tui/screen/niche"
        "github.com/WaClaw-App/waclaw/internal/tui/screen/onboarding"
        "github.com/WaClaw-App/waclaw/internal/tui/screen/pipeline"
        "github.com/WaClaw-App/waclaw/internal/tui/screen/update"
        "github.com/WaClaw-App/waclaw/pkg/protocol"
        tea "github.com/charmbracelet/bubbletea"
)

func main() {
        app := tui.NewApp()

        // Load runtime config (locale, etc.) from ~/.waclaw/config.yaml.
        // Falls back to defaults if missing.
        cfg := tui.LoadConfig()
        tui.ApplyConfig(cfg)

        // Load visual theme from ~/.waclaw/theme.yaml.
        // Falls back to defaults if missing.
        theme := tui.LoadTheme()
        tui.ApplyTheme(theme)

        // ── Register all 20 screens ──────────────────────────────────────
        // Each screen constructor is parameterless — all screen data comes
        // from the backend via HandleNavigate/HandleUpdate after construction.
        // The RegisterScreen call injects the event bus so screens can publish
        // actions and subscribe to navigation/update events.
        registerAllScreens(app)

        // Set the initial screen to Boot. The backend scenario engine will
        // drive subsequent navigation via JSON-RPC navigate commands.
        app.Router().Push(protocol.ScreenBoot)

        // Wire the confirmation overlay's OnConfirm callback to forward
        // confirmation actions to the backend via the RPC client.
        app.WireConfirmation()

        // Initialise the RPC client. The client connects to the backend binary
        // over stdio (stdin/stdout) and translates JSON-RPC messages into typed
        // bus messages that the TUI consumes.
        //
        // When running standalone (no backend), the RPC client's read loop will
        // hit EOF immediately and fire an RPCClosedMsg — the TUI continues to
        // work as a standalone demo with no backend data.
        rpcClient := rpc.NewClient(app.Bus())

        // Wire the RPC client into the app so screen code can send events
        // to the backend via the convenience methods (SendKeyPress, SendAction,
        // SendRequest) without importing the rpc package directly.
        app.SetRPCClient(rpcClient)

        // Wire the RPC client's key/action builders into the app so screen code
        // can send events to the backend without importing the rpc package.
        keyBuilder := rpc.KeyPressBuilder{
                Screen: func() protocol.ScreenID { return app.Router().CurrentID() },
                State:  func() protocol.StateID {
                        if cur := app.Router().Current(); cur != nil {
                                if stateful, ok := cur.(tui.StateReporter); ok {
                                        return stateful.CurrentState()
                                }
                        }
                        return ""
                },
        }
        actionBuilder := rpc.ActionBuilder{
                Screen: func() protocol.ScreenID { return app.Router().CurrentID() },
        }
        app.SetRPCBuilders(keyBuilder, actionBuilder)

        // Start the RPC client with stdin/stdout. The read loop runs in a
        // background goroutine and publishes incoming messages to the bus.
        rpcClient.Start(os.Stdin, os.Stdout)

        // In demo mode (WA_DEMO=1), stdin/stdout are connected to the backend
        // via named pipes for RPC. Bubbletea must use /dev/tty for the actual
        // terminal (keyboard input + screen rendering). In production, the
        // launcher process handles TTY multiplexing; no env var is needed.
        var teaOpts []tea.ProgramOption
        teaOpts = append(teaOpts, tea.WithAltScreen(), tea.WithMouseCellMotion())

        if os.Getenv("WA_DEMO") == "1" {
                tty, err := os.OpenFile("/dev/tty", os.O_RDWR, 0)
                if err != nil {
                        fmt.Fprintf(os.Stderr, "waclaw-tui: failed to open /dev/tty: %v\n", err)
                        os.Exit(1)
                }
                teaOpts = append(teaOpts, tea.WithInput(tty), tea.WithOutput(tty))
        }

        // Create bubbletea program with alt screen and mouse support.
        p := tea.NewProgram(app, teaOpts...)

        if _, err := p.Run(); err != nil {
                fmt.Fprintf(os.Stderr, "waclaw-tui: %v\n", err)
                os.Exit(1)
        }

        // Clean up the RPC client on exit.
        rpcClient.Stop()
}

// registerAllScreens creates and registers all 20 TUI screens.
// This is the single wiring point for the entire screen inventory.
// The order matches the canonical flow defined in the scenario engine.
//
// Note: Some screen constructors return value types. Since the Screen
// interface has methods with pointer receivers (Blur, Focus, SetBus,
// HandleNavigate, HandleUpdate), we must take their address.
func registerAllScreens(app *tui.App) {
        // Onboarding (Screens 1-2)
        boot := onboarding.NewBootModel()
        app.RegisterScreen(&boot)
        login := onboarding.NewLoginModel()
        app.RegisterScreen(&login)

        // Niche (Screens 3, 19)
        sel := niche.NewSelectModel()
        app.RegisterScreen(&sel)
        expl := niche.NewExplorerModel()
        app.RegisterScreen(&expl)

        // Pipeline (Screens 4-6)
        app.RegisterScreen(pipeline.NewScrape())
        rev := pipeline.NewReview()
        app.RegisterScreen(&rev)
        app.RegisterScreen(pipeline.NewSend())

        // Monitor (Screens 7-8)
        app.RegisterScreen(monitor.NewDashboard())
        app.RegisterScreen(monitor.NewResponse())

        // Data (Screens 9-10)
        app.RegisterScreen(data.NewLeadsDB())
        app.RegisterScreen(data.NewTemplateMgr())

        // Infrastructure (Screens 11-14)
        app.RegisterScreen(infra.NewWorkers())
        app.RegisterScreen(infra.NewShield())
        app.RegisterScreen(infra.NewSettings())
        app.RegisterScreen(infra.NewGuardrail())

        // Communication (Screens 15-17)
        app.RegisterScreen(comms.NewCompose())
        app.RegisterScreen(comms.NewHistory())
        app.RegisterScreen(comms.NewFollowUp())

        // License (Screen 18)
        lic := license.New()
        app.RegisterScreen(&lic)

        // Update (Screen 20)
        upd := update.NewModel()
        app.RegisterScreen(&upd)
}
