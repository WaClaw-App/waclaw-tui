// Package protocol defines the shared RPC types used for JSON-RPC 2.0
// communication between the WaClaw TUI (terminal user interface) and the
// WaClaw backend. All types in this package are pure data carriers with no
// external dependencies, making them safe to import from any other package.
package protocol

// ---------------------------------------------------------------------------
// Backend → TUI method names
// ---------------------------------------------------------------------------
// These methods are invoked by the backend and delivered to the TUI as
// server-initiated JSON-RPC 2.0 requests or notifications.

const (
        // MethodNavigate instructs the TUI to transition to a different screen
        // and/or state. Params carry NavigateParams.
        MethodNavigate = "navigate"

        // MethodUpdate pushes incremental data updates to the current screen
        // without changing the screen itself. Params carry screen-specific
        // payload structs.
        MethodUpdate = "update"

        // MethodNotify delivers a user-facing notification (toast, banner, or
        // log entry). Params carry NotifyParams.
        MethodNotify = "notify"

        // MethodValidate sends real-time validation feedback for the currently
        // edited form or input. Params carry ValidateParams.
        MethodValidate = "validate"
)

// ---------------------------------------------------------------------------
// TUI → Backend method names
// ---------------------------------------------------------------------------
// These methods are invoked by the TUI client and handled by the backend.

const (
        // MethodKeyPress is sent every time the user presses a key that the TUI
        // does not handle locally. Params carry KeyPressEvent.
        MethodKeyPress = "key_press"

        // MethodAction is sent when the user triggers a semantic action (e.g.
        // "confirm", "select", "toggle"). Params carry ActionEvent.
        MethodAction = "action"

        // MethodRequest is sent when the TUI explicitly asks the backend for
        // data (e.g. "fetch_leads", "get_stats"). Params carry RequestEvent.
        MethodRequest = "request"
)

// ---------------------------------------------------------------------------
// TUI→Backend action names (used in bus.ActionMsg)
// ---------------------------------------------------------------------------
// These constants replace raw action strings in screen files so that
// typos are caught at compile time and action names are discoverable.

const (
        // --- Infra screen actions (Workers, Settings, Guardrail) ---
        ActionForceScrape      = "force_scrape"
        ActionDeleteWorker     = "delete_worker"
        ActionReloadConfig     = "reload_config"
        ActionOpenEditor       = "open_editor"
        ActionRevertBackup     = "revert_backup"
        ActionPauseWorker      = "pause_worker"
        ActionOpenErrorFile    = "open_error_file"
        ActionShowExample      = "show_example"

        // --- Boot screen actions ---
        ActionBootLogin            = "boot_login"
        ActionBootNiche            = "boot_niche"
        ActionBootGas              = "boot_gas"
        ActionBootDashboard        = "boot_dashboard"
        ActionBootViewResponses    = "boot_view_responses"
        ActionBootViewError        = "boot_view_error"
        ActionBootExit             = "boot_exit"
        ActionBootBuyLicense       = "boot_buy_license"
        ActionBootDisconnectOther  = "boot_disconnect_other"
        ActionBootEnterLicense     = "boot_enter_license"
        ActionBootRelogin          = "boot_relogin"

        // --- Login screen actions ---
        ActionLoginSuccessContinue = "login_success_continue"
        ActionLoginSkip            = "login_skip"
        ActionLoginAddSlot         = "login_add_slot"
        ActionLoginAddAnother      = "login_add_another"
        ActionLoginGas             = "login_gas"
        ActionLoginAddNumber       = "login_add_number"
        ActionLoginLater           = "login_later"
        ActionLoginExpiredKey      = "login_expired_key"
        ActionLoginRetry           = "login_retry"
        ActionLoginChangeSlot      = "login_change_slot"
        ActionLoginBack            = "login_back"
        ActionLoginEnough          = "login_enough"

        // --- Niche Select screen actions ---
        ActionNicheProceed     = "niche_proceed"
        ActionNicheReload      = "niche_reload"
        ActionNicheScrape      = "niche_scrape"
        ActionNicheEditFilter  = "niche_edit_filter"
        ActionNicheOpenFile    = "niche_open_file"
        ActionNicheShowExample = "niche_show_example"
        ActionNicheCustom      = "niche_custom"
        ActionNicheBack        = "niche_back"
        ActionNicheReturnList  = "niche_return_list"

        // --- Niche Explorer screen actions ---
        ActionExplorerDetail     = "explorer_detail"
        ActionExplorerBack       = "explorer_back"
        ActionExplorerGenerate   = "explorer_generate"
        ActionExplorerEdit       = "explorer_edit"
        ActionExplorerAddArea    = "explorer_add_area"
        ActionExplorerCancel     = "explorer_cancel"
        ActionExplorerUse        = "explorer_use"
        ActionExplorerEditConfig = "explorer_edit_config"
        ActionExplorerViewTpl    = "explorer_view_template"
        ActionExplorerSearch     = "explorer_search"

        // --- Review screen actions ---
        ActionReviewApprove = "approve"
        ActionReviewSkip    = "skip"
        ActionReviewBlock   = "block"

        // --- Monitor Dashboard actions ---
        ActionViewDetail      = "view_detail"
        ActionRefresh         = "refresh"
        ActionScrapeAll       = "scrape_all"
        ActionToggleNerdStats = "toggle_nerd_stats"

        // --- Monitor Dashboard navigation actions (keys 1-7) ---
        // The TUI publishes these actions; the backend decides whether to
        // navigate and to which screen. This replaces direct NavigateMsg calls.
        ActionNavLeadsDB   = "nav_leads_db"
        ActionNavResponse  = "nav_response"
        ActionNavWorkers   = "nav_workers"
        ActionNavTemplate  = "nav_template"
        ActionNavAntiBan   = "nav_antiban"
        ActionNavFollowUp  = "nav_followup"
        ActionNavSettings  = "nav_settings"

        // --- Response screen navigation actions ---
        // The TUI publishes these actions; the backend decides whether to
        // navigate and to which screen (e.g. compose, monitor).
        ActionNavCompose  = "nav_compose"
        ActionNavMonitor  = "nav_monitor"

        // --- Response screen actions ---
        ActionSendOffer        = "send_offer"
        ActionLater            = "later"
        ActionSendPricing      = "send_pricing"
        ActionSendInfo         = "send_info"
        ActionMarkInvalid      = "mark_invalid"
        ActionFollowUpLater    = "follow_up_later"
        ActionSkip             = "skip"
        ActionStillFollowUp    = "still_follow_up"
        ActionBlockConfirm     = "block_confirm"
        ActionBlockAllNiches   = "block_all_niches"
        ActionCancelBlock      = "cancel_block"
        ActionConfirmDeal      = "confirm_deal"
        ActionNotDeal          = "not_deal"
        ActionChangeTemplate   = "change_template"
        ActionProcessOne       = "process_one"
        ActionAutoOfferPos     = "auto_offer_positive"
        ActionAutoPerType      = "auto_per_type"
        ActionMarkConverted    = "mark_converted"

        // --- Leads DB screen actions ---
        ActionBulkOffer       = "bulk_offer"
        ActionCustomReply     = "custom_reply"
        ActionMarkConvert     = "mark_convert"
        ActionArchive         = "archive"
        ActionBlock           = "block"
        ActionSendFollowUp    = "send_followup"
        ActionMarkCold        = "mark_cold"
        ActionLastFollowUp    = "last_followup"
        ActionSendIceBreaker  = "send_icebreaker"

        // --- Template Manager screen actions ---
        ActionNewTemplate    = "new_template"
        ActionUseTemplate    = "use_template"
        ActionReloadTemplate = "reload_template"
        ActionOpenFile       = "open_file"

        // --- Additional Infra screen actions ---
        ActionReloadSettings = "reload_settings"
        ActionResumeWorker   = "resume_worker"
        ActionViewLeads      = "view_leads"
)
