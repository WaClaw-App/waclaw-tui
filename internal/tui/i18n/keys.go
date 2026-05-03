package i18n

// ---------------------------------------------------------------------------
// Centralized i18n key constants — DRY principle
// ---------------------------------------------------------------------------
//
// All i18n lookup keys are defined here as constants so that:
//   1. Typos are caught at compile time, not runtime.
//   2. Renaming a key only requires changing one place.
//   3. IDE auto-completion works for key names.
//   4. The en.go and id.go locale maps can be validated against this set.
//
// Convention: group constants by domain using const blocks with comments.
// Key naming follows the pattern: <domain>.<subdomain> — all lowercase,
// dot-separated, no uppercase.
// ---------------------------------------------------------------------------

// Navigation labels
const (
        KeyLabelProceed  = "label.proceed"
        KeyLabelCancel   = "label.cancel"
        KeyLabelBack     = "label.back"
        KeyLabelPause    = "label.pause"
        KeyLabelHelp     = "label.help"
        KeyLabelRefresh  = "label.refresh"
        KeyLabelSearch   = "label.search"
        KeyLabelValidate = "label.validate"
        KeyLabelLicense  = "label.license"
        KeyLabelHistory  = "label.history"
        KeyLabelUpdate   = "label.update"
        KeyLabelNerd     = "label.nerd"
        KeyLabelEdit     = "label.edit"
        KeyLabelOpen     = "label.open"
        KeyLabelReload   = "label.reload"
        KeyLabelQuit     = "label.quit"
        KeyLabelLine     = "label.line"
)

// Status labels
const (
        KeyStatusActive     = "status.active"
        KeyStatusError      = "status.error"
        KeyStatusPaused     = "status.paused"
        KeyStatusSent       = "status.sent"
        KeyStatusDelivered  = "status.delivered"
        KeyStatusRead       = "status.read"
        KeyStatusIdle       = "status.idle"
        KeyStatusScraping   = "status.scraping"
        KeyStatusSending    = "status.sending"
        KeyStatusQualifying = "status.qualifying"
        KeyStatusQueuing    = "status.queuing"
        KeyStatusHealthy    = "status.healthy"
        KeyStatusFlagged    = "status.flagged"
        KeyStatusCooldown   = "status.cooldown"
        KeyStatusDone       = "status.done"
        KeyStatusChecking   = "status.checking"
        KeyStatusWaiting    = "status.waiting"
        KeyStatusWarning    = "status.warning"
        KeyStatusStarting    = "status.starting"
)

// Celebration
const (
        KeyLabelJackpot = "label.jackpot"
        KeyLabelDeal    = "label.deal"
)

// Shield levels
const (
        KeyShieldHealthy = "shield.healthy"
        KeyShieldWarning = "shield.warning"
        KeyShieldDanger  = "shield.danger"
)

// Badges
const (
        KeyBadgeNew       = "badge.new"
        KeyBadgeCold      = "badge.cold"
        KeyBadgeConverted = "badge.converted"
        KeyBadgeResponded = "badge.responded"
        KeyBadgeFailed    = "badge.failed"
)

// Session
const (
        KeySessionEnd        = "session.end"
        KeySessionFarewell   = "session.farewell"
        KeySessionStats      = "session.stats"
        KeySessionBestLead   = "session.best_lead"
        KeySessionZeigarnik  = "session.zeigarnik_tip"
        KeySessionQuitHint   = "session.quit_hint"
)

// Boot screen
const (
        KeyBootTagline   = "boot.tagline"
        KeyBootReturning = "boot.returning"
        KeyBootFirstTime         = "boot.first_time"
        KeyBootFirstTimeSeparator = "boot.first_time_separator"

        // Boot screen — extended labels for doc-parity
        KeyBootFirstTimeTagline      = "boot.first_time_tagline"
        KeyBootPress1LoginAgain      = "boot.press_1_login_again"
        KeyBootMenuLogin             = "boot.menu_login"
        KeyBootMenuLoginDesc         = "boot.menu_login_desc"
        KeyBootMenuNiche             = "boot.menu_niche"
        KeyBootMenuNicheDesc         = "boot.menu_niche_desc"
        KeyBootMenuGas               = "boot.menu_gas"
        KeyBootMenuGasDesc           = "boot.menu_gas_desc"
        KeyBootNicheAlreadyConfigured = "boot.niche_already_configured"
        KeyBootStepsSummary          = "boot.steps_summary"
        KeyBootSystemCheck           = "boot.system_check"
        KeyBootLicenseValid          = "boot.license_valid"
        KeyBootConfigValidCount      = "boot.config_valid_count"
        KeyBootValidationResult      = "boot.validation_result"
        KeyBootArmyReport            = "boot.army_report"
        KeyBootWAConnectedCount      = "boot.wa_connected_count"
        KeyBootNicheWorkerSummary    = "boot.niche_worker_summary"
        KeyBootLeadsCount            = "boot.leads_count"
        KeyBootAutopilotActive       = "boot.autopilot_active"
        KeyBootWARotating            = "boot.wa_rotating"
        KeyBootPressAnyDashboard     = "boot.press_any_dashboard"
        KeyBootArmyWorking           = "boot.army_working"
        KeyBootNewResponses          = "boot.new_responses"
        KeyBootPressEnterResponses   = "boot.press_enter_responses"
        KeyBootWADisconnected        = "boot.wa_disconnected"
        KeyBootScraperStillRunning   = "boot.scraper_still_running"
        KeyBootScraperOnlyNote       = "boot.scraper_only_note"
        KeyBootPressLoginAgain       = "boot.press_login_again"
        KeyBootConfigError           = "boot.config_error"
        KeyBootSomeNichesOK          = "boot.some_niches_ok"
        KeyBootConfigErrorNiche      = "boot.config_error_niche"
        KeyBootOtherWorkersStill     = "boot.other_workers_still"
        KeyBootPressVError           = "boot.press_v_error"
        KeyBootVViewError            = "boot.v_view_error"
        KeyBootDashboard             = "boot.dashboard"
        KeyBootQExit                 = "boot.q_exit"
        KeyBootLicenseExpired        = "boot.license_expired"
        KeyBootAllNichesPaused       = "boot.all_niches_paused"
        KeyBootLicenseExpiredMsg     = "boot.license_expired_msg"
        KeyBootRenewToContinue       = "boot.renew_to_continue"
        KeyBootEnterNewLicense       = "boot.enter_new_license"
        KeyBootBuyLicense            = "boot.buy_license"
        KeyBootDeviceConflict        = "boot.device_conflict"
        KeyBootDeviceConflictMsg     = "boot.device_conflict_msg"
        KeyBootOneLicenseOneDevice   = "boot.one_license_one_device"
        KeyBootAllPausedUntilResolved = "boot.all_paused_until_resolved"
        KeyBootDisconnectOther       = "boot.disconnect_other"
        KeyBootOtherDeviceInfo       = "boot.other_device_info"
        KeyBootWorkerReady           = "boot.worker_ready"
)

// Login screen
const (
        KeyLoginQRWait  = "login.qr_wait"
        KeyLoginScanned = "login.scanned"
        KeyLoginSuccess = "login.success"
        KeyLoginExpired = "login.expired"
        KeyLoginFailed  = "login.failed"

        // Login screen — extended labels for doc-parity
        KeyLoginTitle              = "login.title"
        KeyLoginQRInstruction      = "login.qr_instruction"
        KeyLoginMultiSlotHint      = "login.multi_slot_hint"
        KeyLoginMoreNumbersSafer   = "login.more_numbers_safer"
        KeyLoginWaitingScan        = "login.waiting_scan"
        KeyLoginSlotIndicator      = "login.slot_indicator"
        KeyLoginConnectServer      = "login.connect_server"
        KeyLoginWaitingHPScan      = "login.waiting_hp_scan"
        KeyLoginSyncContacts       = "login.sync_contacts"
        KeyLoginSlotsFilled        = "login.slots_filled"
        KeyLoginAddSlot            = "login.add_slot"
        KeyLoginSkip               = "login.skip"
        KeyLoginScanDetected       = "login.scan_detected"
        KeyLoginScanSuccess        = "login.scan_success"
        KeyLoginSyncContactsCount  = "login.sync_contacts_count"
        KeyLoginAddAnother         = "login.add_another"
        KeyLoginYesAdd             = "login.yes_add"
        KeyLoginEnoughContinue     = "login.enough_continue"
        KeyLoginContactsSynced     = "login.contacts_synced"
        KeyLoginSlotConnected      = "login.slot_connected"
        KeyLoginConnectedNow       = "login.connected_now"
        KeyLoginMoreSafer          = "login.more_safer"
        KeyLoginAddNumber          = "login.add_number"
        KeyLoginEnoughGas          = "login.enough_gas"
        KeyLoginLater              = "login.later"
        KeyLoginSessionExpired     = "login.session_expired"
        KeyLoginSlotExpired        = "login.slot_expired"
        KeyLoginExpiredAutoPause   = "login.expired_auto_pause"
        KeyLoginLastSession        = "login.last_session"
        KeyLoginConnectionFailed   = "login.connection_failed"
        KeyLoginSlotFailedNote     = "login.slot_failed_note"
        KeyLoginWAServerIssue      = "login.wa_server_issue"
        KeyLoginTryAgainLater      = "login.try_again_later"
        KeyLoginTryAgain           = "login.try_again"
        KeyLoginChangeSlot         = "login.change_slot"
        KeyLoginBack               = "login.back"
)

// Niche screens
const (
        KeyNicheSelect    = "niche.select"
        KeyNicheSelected  = "niche.selected"
        KeyNicheCustom    = "niche.custom"
        KeyNicheFilters   = "niche.filters"
        KeyNicheConfigErr = "niche.config_err"

        // Explorer-specific keys
        KeyNicheExplorerTitle        = "niche.explorer_title"
        KeyNicheExplorerSubtitle     = "niche.explorer_subtitle"
        KeyNicheExplorerPopular      = "niche.explorer_popular"
        KeyNicheExplorerSearching    = "niche.explorer_searching"
        KeyNicheExplorerResults      = "niche.explorer_results"
        KeyNicheNicheDipilih         = "niche.niche_dipilih"
        KeyNicheLabel                = "niche.label"
        KeyNicheExplorerSource       = "niche.explorer_source"
        KeyNicheExplorerGenConfig    = "niche.explorer_gen_config"
        KeyNicheExplorerGenSuccess   = "niche.explorer_gen_success"
        KeyNicheExplorerGenProgress  = "niche.explorer_gen_progress"
        KeyNicheExplorerAreaAuto     = "niche.explorer_area_auto"
        KeyNicheExplorerEditFile     = "niche.explorer_edit_file"
        KeyNicheExplorerReload       = "niche.explorer_reload"
        KeyNicheExplorerParallel     = "niche.explorer_parallel"
        KeyNicheNicheIs              = "niche.niche_is"
        KeyNicheCustomDir            = "niche.custom_dir"
        KeyNicheCustomMin            = "niche.custom_min"
        KeyNicheCustomExample        = "niche.custom_example"
        KeyNicheCustomReady          = "niche.custom_ready"
        KeyNicheErrPaused            = "niche.err_paused"
        KeyNicheErrOtherOK           = "niche.err_other_ok"
        KeyNicheProblems             = "niche.problems"
        KeyNicheTargets              = "niche.targets"
        KeyNicheAreaCount            = "niche.area_count"
        KeyNicheFilterDefault        = "niche.filter_default"
        KeyNicheTemplateGen          = "niche.template_gen"
        KeyNicheJustRight            = "niche.just_right"
        KeyNicheMoreArea              = "niche.more_area"
        KeyNicheCanEdit              = "niche.can_edit"
        KeyNicheWorkerParallel       = "niche.worker_parallel"
        KeyNicheScrapeOwn            = "niche.scrape_own"

        // Niche — extended labels for doc-parity
        KeyNicheMoreNiche            = "niche.more_niche"
        KeyNicheCheckUncheck         = "niche.check_uncheck"
        KeyNicheGasChecked           = "niche.gas_checked"
        KeyNicheGasNiche             = "niche.gas_niche"
        KeyNicheChange               = "niche.change"
        KeyNicheBack                 = "niche.back"
        KeyNicheNicheYaml            = "niche.niche_yaml"
        KeyNicheIceBreaker           = "niche.ice_breaker"
        KeyNicheReload               = "niche.reload"
        KeyNichePickExisting         = "niche.pick_existing"
        KeyNicheConfigErrLabel       = "niche.config_err_label"
        KeyNicheGasScrape            = "niche.gas_scrape"
        KeyNicheEditFilter           = "niche.edit_filter"
        KeyNicheOpenFile             = "niche.open_file"
        KeyNicheShowExample          = "niche.show_example"
        KeyNicheAreaKota             = "niche.area_kota"
        KeyNicheMultiParallel        = "niche.multi_parallel"
        KeyNicheMultiScrapeOwn       = "niche.multi_scrape_own"
        KeyNicheKecamatan            = "niche.kecamatan"
)

// Niche explorer — extended labels
const (
        KeyNicheExplorerPickConfig    = "niche.explorer_pick_config"
        KeyNicheExplorerSearchCat     = "niche.explorer_search_cat"
        KeyNicheExplorerPick          = "niche.explorer_pick"
        KeyNicheExplorerDetail        = "niche.explorer_detail"
        KeyNicheExplorerDIY           = "niche.explorer_diy"
        KeyNicheExplorerSubCat        = "niche.explorer_sub_cat"
        KeyNicheExplorerSourceLabel   = "niche.explorer_source_label"
        KeyNicheExplorerEditFirst     = "niche.explorer_edit_first"
        KeyNicheExplorerFolder        = "niche.explorer_folder"
        KeyNicheExplorerGasUse        = "niche.explorer_gas_use"
        KeyNicheExplorerEditConfig    = "niche.explorer_edit_config"
        KeyNicheExplorerViewTemplate  = "niche.explorer_view_template"
        KeyNicheExplorerCancel        = "niche.explorer_cancel"
        KeyNicheExplorerSubKategori   = "niche.explorer_sub_kategori"
        KeyNicheExplorerAreaLabel     = "niche.explorer_area_label"
        KeyNicheExplorerAreaDetect    = "niche.explorer_area_detect"
        KeyNicheExplorerAreaExplain   = "niche.explorer_area_explain"
        KeyNicheExplorerAreaSame      = "niche.explorer_area_same"
        KeyNicheExplorerAddArea       = "niche.explorer_add_area"
        KeyNicheExplorerUseExisting   = "niche.explorer_use_existing"
        KeyNicheExplorerNoResult      = "niche.explorer_no_result"
        KeyNicheExplorerKategori      = "niche.explorer_kategori"
        KeyNicheExplorerFilterDefault = "niche.explorer_filter_default"
        KeyNicheLine                  = "niche.line"
        KeyNicheExplorerGenBarLabel   = "niche.explorer_gen_bar_label"
        KeyNicheExplorerTitleDetail   = "niche.explorer_title_detail"
        KeyNicheTemplateCount         = "niche.template_count"

        // Niche — path strings (i18n for locale-aware paths)
        KeyNicheCustomDirPath      = "niche.custom_dir_path"
        KeyNicheCustomExamplePath  = "niche.custom_example_path"
        KeyNicheExplorerFolderPath = "niche.explorer_folder_path"
)

// Scrape screen
const (
        KeyScrapeActive         = "scrape.active"
        KeyScrapeComplete       = "scrape.complete"
        KeyScrapeEmpty          = "scrape.empty"
        KeyScrapeError          = "scrape.error"
        KeyScrapeFound          = "scrape.found"
        KeyScrapeTarget         = "scrape.target"
        KeyScrapeArea           = "scrape.area"
        KeyScrapeFilter         = "scrape.filter"
        KeyScrapeQualified      = "scrape.qualified"
        KeyScrapeDuplicates     = "scrape.duplicates"
        KeyScrapeNewLeads       = "scrape.new_leads"
        KeyScrapeScanning       = "scrape.scanning"
        KeyScrapeHasWebsite     = "scrape.has_website"
        KeyScrapeNoWebsite      = "scrape.no_website"
        KeyScrapeSkipSuffix     = "scrape.skip_suffix"
        KeyScrapeAutoReview     = "scrape.auto_review"
        KeyScrapeAutoPilotMsg   = "scrape.auto_pilot_msg"
        KeyScrapeMulti          = "scrape.multi"
        KeyScrapeStaggered      = "scrape.staggered"
        KeyScrapeWorkers        = "scrape.workers"
        KeyScrapeTotalActive    = "scrape.total_active"
        KeyScrapeTotal          = "scrape.total"
        KeyScrapeLeadsUnit      = "scrape.leads_unit"
        KeyScrapeAutoReviewHint = "scrape.auto_review_hint"
        KeyScrapeThrottleExplain = "scrape.throttle_explain"
        KeyScrapeThrottleResume = "scrape.throttle_resume"
        KeyScrapeThrottleGMapsOnly = "scrape.throttle_gmaps_only"
        KeyScrapeTabSwitch      = "scrape.tab_switch"
        KeyScrapeEnterDetail    = "scrape.enter_detail"
        KeyScrapeNextScrape     = "scrape.next_scrape"
        KeyScrapeLastScrape     = "scrape.last_scrape"
        KeyScrapeQueue          = "scrape.queue"
        KeyScrapeIdle           = "scrape.idle"
        KeyScrapeSkipWait       = "scrape.skip_wait"
        KeyScrapeEmptyHints     = "scrape.empty_hints"
        KeyScrapeChangeArea     = "scrape.change_area"
        KeyScrapeEditFilter     = "scrape.edit_filter"
        KeyScrapeAddQuery       = "scrape.add_query"
        KeyScrapeRetry          = "scrape.retry"
        KeyScrapeGMapsThrottle  = "scrape.gmaps_throttle"
        KeyScrapeWABackground   = "scrape.wa_background"
        KeyScrapeAutoApproved   = "scrape.auto_approved"
        KeyScrapeScore          = "scrape.score"
        KeyScrapePriorityQueue  = "scrape.priority_queue"
        KeyScrapeBestTemplate   = "scrape.best_template"
        KeyScrapeNewLeadsWaiting = "scrape.new_leads_waiting"
        KeyScrapeWAValidation   = "scrape.wa_validation"
        KeyScrapeWAValidationProg = "scrape.wa_validation_prog"
        KeyScrapeWATotal        = "scrape.wa_total"
        KeyScrapeWAHas          = "scrape.wa_has"
        KeyScrapeWANot          = "scrape.wa_not"
        KeyScrapeWAPending      = "scrape.wa_pending"
        KeyScrapeWAReady        = "scrape.wa_ready"
        KeyScrapeWASkipped      = "scrape.wa_skipped"
        KeyScrapeWASaving       = "scrape.wa_saving"
        KeyScrapeEstimate       = "scrape.estimate"
        KeyScrapeTotalReady     = "scrape.total_ready"
        KeyScrapeBatchSelesai   = "scrape.batch_selesai"
        KeyScrapeBatchDone       = "scrape.batch_done"

        KeyScrapeIdleWaiting      = "scrape.idle_waiting"
        KeyScrapeErrorReassure    = "scrape.error_reassure"
        KeyScrapeWAQueueAnnotation = "scrape.wa_queue_annotation"
        KeyScrapeWAMarkedAnnotation = "scrape.wa_marked_annotation"
        KeyScrapeWACekUnit        = "scrape.wa_cek_unit"
        KeyScrapeBatchAutoQueue   = "scrape.batch_auto_queue"
        KeyScrapeBatchWorkHours   = "scrape.batch_work_hours"
        KeyScrapeBatchTotalSummary = "scrape.batch_total_summary"
        KeyScrapeScraperError     = "scrape.scraper_error"
        KeyScrapeIdleNewLeads     = "scrape.idle_new_leads"
        KeyScrapeNicheLabel       = "scrape.niche_label"
        KeyScrapeEmptyAreaHint    = "scrape.empty_area_hint"
        KeyScrapeEmptyFilterHint  = "scrape.empty_filter_hint"
        KeyScrapeEmptyQueryHint   = "scrape.empty_query_hint"
)

// Send screen
const (
        KeySendActive         = "send.active"
        KeySendPaused         = "send.paused"
        KeySendOffHours       = "send.off_hours"
        KeySendRateLimit      = "send.rate_limit"
        KeySendDailyLimit     = "send.daily_limit"
        KeySendFailed         = "send.failed"
        KeySendAutoPilotMsg   = "send.auto_pilot_msg"
        KeySendQueueCount     = "send.queue_count"
        KeySendPauseReason    = "send.pause_reason"
        KeySendAutoContinue   = "send.auto_continue"
        KeySendEmergency      = "send.emergency"
        KeySendIntervalHint   = "send.interval_hint"
        KeySendRemaining      = "send.remaining"
        KeySendGoodJob        = "send.good_job"
        KeySendShouldNot      = "send.should_not"
        KeySendPreValMiss     = "send.pre_val_miss"
        KeySendAutoSkip       = "send.auto_skip"
        KeySendRetryManual    = "send.retry_manual"
        KeySendValidateRetry  = "send.validate_retry"
        KeySendAllSlotsDown   = "send.all_slots_down"
        KeySendAllScraping    = "send.all_scraping"
        KeySendLeadsSafe      = "send.leads_safe"
        KeySendPendingOnly    = "send.pending_only"
        KeySendLogin          = "send.login"
        KeySendLoginOneByOne  = "send.login_one_by_one"
        KeySendWARotator      = "send.wa_rotator"
        KeySendRate           = "send.rate"
        KeySendReady          = "send.ready"
        KeySendWAValidated    = "send.wa_validated"
        KeySendRotation       = "send.rotation"
        KeySendSending        = "send.sending"
        KeySendWaiting        = "send.waiting"
        KeySendToday          = "send.today"
        KeySendNextAt         = "send.next_at"
        KeySendTodaySent      = "send.today_sent"
        KeySendTodayResp      = "send.today_resp"
        KeySendTodayConv      = "send.today_conv"
        KeySendResponseIn     = "send.response_in"
        KeySendSendOffer      = "send.send_offer"
        KeySendReplyCustom    = "send.reply_custom"
        KeySendResume         = "send.resume"
        KeySendSlotCount      = "send.slot_count"
        KeySendNicheCount     = "send.niche_count"
        KeySendNow             = "send.now"
        KeySendNextLabel       = "send.next_label"

        KeySendAutoManages       = "send.auto_manages"
        KeySendNoRefreshNeeded   = "send.no_refresh_needed"
        KeySendContinueTomorrow  = "send.continue_tomorrow"
        KeySendSkipAja           = "send.skip_aja"
        KeySendSwitchNiche       = "send.switch_niche"
        KeySendSkipWait          = "send.skip_wait"
        KeySendDefaultWorkHours  = "send.default_work_hours"
)

// Monitor screen
const (
        KeyMonitorLive           = "monitor.live"
        KeyMonitorIdle           = "monitor.idle"
        KeyMonitorNight          = "monitor.night"
        KeyMonitorError          = "monitor.error"
        KeyMonitorWAConnected    = "monitor.wa_connected"
        KeyMonitorNicheActive    = "monitor.niche_active"
        KeyMonitorWARotator      = "monitor.wa_rotator"
        KeyMonitorWorkerPool     = "monitor.worker_pool"
        KeyMonitorRecentActivity = "monitor.recent_activity"
        KeyMonitorAutoPilot      = "monitor.auto_pilot"
        KeyMonitorIdleBody       = "monitor.idle_body"
        KeyMonitorQueued         = "monitor.queued"
        KeyMonitorMsgsSent       = "monitor.msgs_sent"
        KeyMonitorResponses      = "monitor.responses"
        KeyMonitorCanMinimize    = "monitor.can_minimize"
        KeyMonitorWorkersDesc    = "monitor.workers_desc"
        KeyMonitorWillNotify     = "monitor.will_notify"
        KeyMonitorViewDetail     = "monitor.view_detail"
        KeyMonitorNightMode      = "monitor.night_mode"
        KeyMonitorWorkHours      = "monitor.work_hours"
        KeyMonitorNow            = "monitor.now"
        KeyMonitorSenderPaused   = "monitor.sender_paused"
        KeyMonitorScraperRunning = "monitor.scraper_running"
        KeyMonitorAutoResume     = "monitor.auto_resume"
        KeyMonitorDaySummary     = "monitor.day_summary"
        KeyMonitorConverts       = "monitor.converts"
        KeyMonitorPerNiche       = "monitor.per_niche"
        KeyMonitorBestTime       = "monitor.best_time"
        KeyMonitorWADisconnected  = "monitor.wa_disconnected"
        KeyMonitorSlotDisconnected = "monitor.slot_disconnected"
        KeyMonitorScraperOK      = "monitor.scraper_ok"
        KeyMonitorDatabaseOK     = "monitor.database_ok"
        KeyMonitorSlotsActive    = "monitor.slots_active"
        KeyMonitorSlotPending    = "monitor.slot_pending"
        KeyMonitorAutoReconnect  = "monitor.auto_reconnect"
        KeyMonitorRelogin        = "monitor.relogin"
        KeyMonitorNoData         = "monitor.no_data"
        KeyMonitorStartScrape    = "monitor.start_scrape"
        KeyMonitorSetupNiche     = "monitor.setup_niche"
        KeyMonitorStartSend      = "monitor.start_send"
        KeyMonitorTipStartScrape = "monitor.tip_start_scrape"
        KeyMonitorPending        = "monitor.pending"
        KeyMonitorEmpty          = "monitor.empty"
        KeyMonitorAutoOfferAll   = "monitor.auto_offer_all"
        KeyMonitorAutoPerNiche   = "monitor.auto_per_niche"
        KeyMonitorOfferPerNiche  = "monitor.offer_per_niche"
        KeyMonitorToday          = "monitor.today"
        KeyMonitorCooldown       = "monitor.cooldown"
        KeyMonitorResponded      = "monitor.responded"
        KeyMonitorLeadsFound     = "monitor.leads_found"
        KeyMonitorThisWeek       = "monitor.this_week"
        KeyMonitorConversionRate = "monitor.conversion_rate"
        KeyMonitorBestDay        = "monitor.best_day"
        KeyMonitorNavLeads       = "monitor.nav_leads"
        KeyMonitorNavMessages    = "monitor.nav_messages"
        KeyMonitorNavWorkers     = "monitor.nav_workers"
        KeyMonitorNavTemplate    = "monitor.nav_template"
        KeyMonitorNavAntiban     = "monitor.nav_antiban"
        KeyMonitorNavFollowup    = "monitor.nav_followup"
        KeyMonitorNavSettings    = "monitor.nav_settings"
        KeyMonitorScrapeAll      = "monitor.scrape_all"
        KeyMonitorCheckStatus     = "monitor.check_status"
        KeyMonitorWARotatorNum      = "monitor.wa_rotator_num"
        KeyMonitorQueuedTotal        = "monitor.queued_total"
        KeyMonitorResponsesToday     = "monitor.responses_today"
        KeyMonitorDefaultWorkHours   = "monitor.default_work_hours"
        KeyMonitorProblemStatus      = "monitor.problem_status"
        KeyMonitorAutoResumeDetail   = "monitor.auto_resume_detail"
        KeyMonitorPendingUnreplied    = "monitor.pending_unreplied"
        KeyMonitorTotalLeads          = "monitor.total_leads"
        KeyMonitorRecentActivityAll   = "monitor.recent_activity_all"
        KeyMonitorTipAuto             = "monitor.tip_auto"
)

// Response screen
const (
        KeyResponsePositive     = "response.positive"
        KeyResponseCurious      = "response.curious"
        KeyResponseNegative     = "response.negative"
        KeyResponseMaybe        = "response.maybe"
        KeyResponseAuto         = "response.auto"
        KeyResponseProcessOne   = "response.process_one"
        KeyResponseGotReply     = "response.got_reply"
        KeyResponseSendOffer    = "response.send_offer"
        KeyResponseReplyCustom  = "response.reply_custom"
        KeyResponseLater        = "response.later"
        KeyResponseSendInfo     = "response.send_info"
        KeyResponseGotReplyPlain = "response.got_reply_plain"
        KeyResponseMarkInvalid  = "response.mark_invalid"
        KeyResponseFollowUpLater = "response.follow_up_later"
        KeyResponseSkipIt       = "response.skip_it"
        KeyResponseGotReplyAuto = "response.got_reply_auto"
        KeyResponseAutoReplySkip = "response.auto_reply_skip"
        KeyResponseStillFollowUp = "response.still_follow_up"
        KeyResponseStopDetected = "response.stop_detected"
        KeyResponseAutoAddedDNC = "response.auto_added_dnc"
        KeyResponseBlockConfirm = "response.block_confirm"
        KeyResponseBlockAllNiche = "response.block_all_niche"
        KeyResponseCancelBlock  = "response.cancel_block"
        KeyResponseDealDetected = "response.deal_detected"
        KeyResponseAutoDeal     = "response.auto_deal"
        KeyResponseConfirmDeal  = "response.confirm_deal"
        KeyResponseNotDeal      = "response.not_deal"
        KeyResponseReplyFirst   = "response.reply_first"
        KeyResponseHotLead      = "response.hot_lead"
        KeyResponseAutoPrioritize = "response.auto_prioritize"
        KeyResponseOfferNow     = "response.offer_now"
        KeyResponseOfferPreview = "response.offer_preview"
        KeyResponseSend         = "response.send"
        KeyResponseChangeTemplate = "response.change_template"
        KeyResponseEditFirst    = "response.edit_first"
        KeyResponseMultiQueue   = "response.multi_queue"
        KeyResponsePositif      = "response.positif"
        KeyResponseCuriousBadge = "response.curious_badge"
        KeyResponseAutoReplyBadge = "response.auto_reply_badge"
        KeyResponseAutoOfferPos = "response.auto_offer_pos"
        KeyResponseAutoPerType  = "response.auto_per_type"
        KeyResponseConversion       = "response.conversion"
        KeyResponseConversionDeal    = "response.conversion_deal"
        KeyResponseStopNotice1       = "response.stop_notice_1"
        KeyResponseStopNotice2       = "response.stop_notice_2"
        KeyResponseStopNotice3       = "response.stop_notice_3"
        KeyResponseDealNotice1       = "response.deal_notice_1"
        KeyResponseDealNotice2       = "response.deal_notice_2"
        KeyResponseClassPositif      = "response.class_positif"
        KeyResponseClassCurious      = "response.class_curious"
        KeyResponseClassAutoReply    = "response.class_auto_reply"
        KeyResponseClosingTrigger  = "response.closing_trigger"
        KeyResponseAutoDealShort = "response.auto_deal_short"
        KeyResponseHotLeadTrigger   = "response.hot_lead_trigger"
        KeyResponseHotLeadAutoPrior = "response.hot_lead_auto_prior"
        KeyResponseHotLeadOfferHint = "response.hot_lead_offer_hint"
)

// Conversion screen
const (
        KeyConversionMarkConverted = "conversion.mark_converted"
        KeyConversionFromPipeline  = "conversion.from_pipeline"
        KeyConversionTimeTaken     = "conversion.time_taken"
        KeyConversionOrdinal       = "conversion.ordinal"
        KeyConversionTrophyWeek    = "conversion.trophy_week"
        KeyConversionRevenueWeek   = "conversion.revenue_week"
        KeyConversionCelebrating    = "conversion.celebrating"
        KeyConversionDissolving     = "conversion.dissolving"
)

// Settings screen (legacy)
const (
        KeySettingsOverview = "settings.overview"
        KeySettingsReload   = "settings.reload"
        KeySettingsError    = "settings.error"
)

// Guardrail screen (legacy)
const (
        KeyGuardrailClean    = "guardrail.clean"
        KeyGuardrailErrors   = "guardrail.errors"
        KeyGuardrailWarnings = "guardrail.warnings"
)

// Workers screen (infra)
const (
        KeyWorkersTitle         = "workers.title"
        KeyWorkersActive        = "workers.active"
        KeyWorkersIdle          = "workers.idle"
        KeyWorkersPausedStatus  = "workers.paused_status"
        KeyWorkersDetail        = "workers.detail"
        KeyWorkersAddNiche      = "workers.add_niche"
        KeyWorkersScrape        = "workers.scrape"
        KeyWorkersReview        = "workers.review"
        KeyWorkersQueue         = "workers.queue"
        KeyWorkersSend          = "workers.send"
        KeyWorkersArea          = "workers.area"
        KeyWorkersTemplate      = "workers.template"
        KeyWorkersTotalPipeline = "workers.total_pipeline"
        KeyWorkersParallel      = "workers.parallel"
        KeyWorkersPauseWorker   = "workers.pause_worker"
        KeyWorkersForceScrape   = "workers.force_scrape"
        KeyWorkersViewLeads     = "workers.view_leads"
        KeyWorkersResume        = "workers.resume"
        KeyWorkersDelete        = "workers.delete"
        KeyWorkersCollected     = "workers.collected"
        KeyWorkersAutoResume    = "workers.auto_resume"
        KeyWorkersMoreNiches    = "workers.more_niches"
        KeyWorkersIndependent   = "workers.independent"
        KeyWorkersQuery         = "workers.query"
        KeyWorkersDone          = "workers.done"
        KeyWorkersPassed        = "workers.passed"
        KeyWorkersSkipped       = "workers.skipped"
        KeyWorkersDuplicates    = "workers.duplicates"
        KeyWorkersLowRating     = "workers.low_rating"
        KeyWorkersNext          = "workers.next"
        KeyWorkersToday         = "workers.today"
        KeyWorkersResponseRate  = "workers.response_rate"
        KeyWorkersConversionRate = "workers.conversion_rate"
        KeyWorkersAvgRespond    = "workers.avg_respond"
        KeyWorkersPerforma      = "workers.performa"

        // 3F fix additions
        KeyWorkersChoose         = "workers.choose"
        KeyWorkersWorkerLabel    = "workers.worker_label"
        KeyWorkersYouPaused      = "workers.you_paused"
        KeyWorkersLeadsInDB      = "workers.leads_in_db"
        KeyWorkersDuplicateCount = "workers.duplicate_count"
        KeyWorkersLowRatingCount = "workers.low_rating_count"
        KeyWorkersIceBreakerCount = "workers.ice_breaker_count"
        KeyWorkersAutoOfferCount = "workers.auto_offer_count"
        KeyWorkersNextIn         = "workers.next_in"
        KeyWorkersTodayCount     = "workers.today_count"

        // 3F-2 fix additions — pipeline labels + qualify
        KeyWorkersQualify        = "workers.qualify"
        KeyWorkersFound          = "workers.found"
        KeyWorkersQueued         = "workers.queued"
        KeyWorkersSentLabel      = "workers.sent_label"
        KeyWorkersQueryDone      = "workers.query_done"
        KeyWorkersQueryScanning  = "workers.query_scanning"
        KeyWorkersQueryWaiting   = "workers.query_waiting"
        KeyWorkersActiveLabel    = "workers.active_label"

        // 3F-3 fix additions — pipeline label
        KeyWorkersPipelineLabel   = "workers.pipeline_label"

        // TUI-vs-doc misalignment fix additions
        KeyWorkersPickPool       = "workers.pick_pool"
        KeyWorkersAutoStart      = "workers.auto_start"
        KeyWorkersAddLabel       = "workers.add_label"
        KeyWorkersQuitCancel     = "workers.quit_cancel"
        KeyWorkersViewLeadsLong  = "workers.view_leads_long"

        // 3F-5 fix additions — paused view + found/passed summary
        KeyWorkersFoundPassed    = "workers.found_passed"
        KeyWorkersStatusPaused   = "workers.status_paused"
        KeyWorkersOptionsResume  = "workers.options_resume"
        KeyWorkersOptionsDelete  = "workers.options_delete"
        KeyWorkersOptionsViewLeads = "workers.options_view_leads"
        KeyWorkersLeadCount      = "workers.lead_count"
        KeyWorkersAutoResumeMsg  = "workers.auto_resume_msg"

        // 3F-6 fix additions — i18n-ize add-niche list items
        KeyWorkersNicheFotografer     = "workers.niche_fotografer"
        KeyWorkersNicheFotograferDesc = "workers.niche_fotografer_desc"
        KeyWorkersNicheAkuntan        = "workers.niche_akuntan"
        KeyWorkersNicheAkuntanDesc    = "workers.niche_akuntan_desc"
        KeyWorkersNicheCustom         = "workers.niche_custom"
        KeyWorkersNicheCustomDesc     = "workers.niche_custom_desc"
)

// Anti-ban shield screen (infra)
const (
        KeyShieldTitle            = "shield.title"
        KeyShieldAllSafe          = "shield.all_safe"
        KeyShieldHasWarning       = "shield.has_warning"
        KeyShieldDangerLabel      = "shield.danger_label"
        KeyShieldHealthScore      = "shield.health_score"
        KeyShieldWARotator        = "shield.wa_rotator"
        KeyShieldNumbers          = "shield.numbers"
        KeyShieldSlotActive       = "shield.slot_active"
        KeyShieldSlotCooldown     = "shield.slot_cooldown"
        KeyShieldSlotFlagged      = "shield.slot_flagged"
        KeyShieldRateLimiting     = "shield.rate_limiting"
        KeyShieldPerSlot          = "shield.per_slot"
        KeyShieldPerDay           = "shield.per_day"
        KeyShieldPerLead          = "shield.per_lead"
        KeyShieldDailyBudget      = "shield.daily_budget"
        KeyShieldWorkHoursGuard   = "shield.work_hours_guard"
        KeyShieldWorkHours        = "shield.work_hours"
        KeyShieldTimezone         = "shield.timezone"
        KeyShieldSendHours        = "shield.send_hours"
        KeyShieldScrapeHours      = "shield.scrape_hours"
        KeyShieldNow              = "shield.now"
        KeyShieldInWorkHours      = "shield.in_work_hours"
        KeyShieldPatternGuard     = "shield.pattern_guard"
        KeyShieldTemplateRotation = "shield.template_rotation"
        KeyShieldSpamGuard        = "shield.spam_guard"
        KeyShieldBanRiskScore     = "shield.ban_risk_score"
        KeyShieldRiskLow          = "shield.risk_low"
        KeyShieldRiskMedium       = "shield.risk_medium"
        KeyShieldRiskHigh         = "shield.risk_high"
        KeyShieldIndicators       = "shield.indicators"
        KeyShieldEvenDist         = "shield.even_distribution"
        KeyShieldCooldownOk       = "shield.cooldown_sufficient"
        KeyShieldTemplateVaried   = "shield.template_varied"
        KeyShieldNoOverload       = "shield.no_overload"
        KeyShieldWorkHoursOk      = "shield.work_hours_followed"
        KeyShieldSpamGuardActive  = "shield.spam_guard_active"
        KeyShieldDNCRespected     = "shield.dnc_respected"
        KeyShieldSlotDetail       = "shield.slot_detail"
        KeyShieldSlotNumber       = "shield.slot_number"
        KeyShieldStatistics       = "shield.statistics"
        KeyShield7Days            = "shield.7days"
        KeyShieldConfig           = "shield.config"
        KeyShieldEditConfig       = "shield.edit_config"
        KeyShieldAutoPause        = "shield.auto_pause"
        KeyShieldAutoRecover      = "shield.auto_recover"
        KeyShieldAddNumber        = "shield.add_number"
        KeyShieldPauseSend        = "shield.pause_send"
        KeyShieldAllSafeEmoji     = "shield.all_safe_emoji"
        KeyShieldWarningEmoji     = "shield.warning_emoji"
        KeyShieldDangerEmoji      = "shield.danger_emoji"
        KeyShieldLetItBe          = "shield.let_it_be"
        KeyShieldAlreadyMoved     = "shield.already_moved"
        KeyShieldHour             = "shield.hour"
        KeyShieldCooldown         = "shield.cooldown"
        KeyShieldTotal            = "shield.total"
        KeyShieldStatus           = "shield.status"
        KeyShieldHealthyStatus    = "shield.healthy_status"
        KeyShieldReady            = "shield.ready"

        // Shield settings detail labels (anti-ban config display)
        KeyShieldAntiBan          = "shield.anti_ban"
        KeyShieldMinDelay         = "shield.min_delay"
        KeyShieldMaxDelay         = "shield.max_delay"
        KeyShieldDelayVariance    = "shield.delay_variance"
        KeyShieldCooldownLimit    = "shield.cooldown_limit"
        KeyShieldHealthThreshold  = "shield.health_threshold"
        KeyShieldRotatorMode      = "shield.rotator_mode_setting"
        KeyShieldRotationMode     = "shield.rotation_mode"
        KeyShieldEmojiVariation   = "shield.emoji_variation"
        KeyShieldParagraphShuffle = "shield.paragraph_shuffle"
        KeyShieldConfigSection    = "shield.config_section"

        // Shield settings — spam_guard section labels
        KeyShieldMsgInterval      = "shield.msg_interval"
        KeyShieldFollowupDelay    = "shield.followup_delay"
        KeyShieldFollowupVariant  = "shield.followup_variant"
        KeyShieldColdAfter        = "shield.cold_after"
        KeyShieldRecontactDelay   = "shield.recontact_delay"
        KeyShieldAutoBlock        = "shield.auto_block"
        KeyShieldDupCrossNiche    = "shield.dup_cross_niche"
        KeyShieldWAPreValidation  = "shield.wa_pre_validation"
        KeyShieldWAValidationMethod = "shield.wa_validation_method"

        // 3F fix additions — pattern guard
        KeyShieldPatternTemplate  = "shield.pattern_template"
        KeyShieldPatternTimeVar   = "shield.pattern_time_variance"
        KeyShieldPatternMsgVar    = "shield.pattern_msg_variance"
        KeyShieldPatternEmoji     = "shield.pattern_emoji"
        KeyShieldPatternShuffle   = "shield.pattern_shuffle"

        // 3F fix additions — spam guard
        KeyShieldSpamPerLead      = "shield.spam_per_lead"
        KeyShieldSpamLifetime     = "shield.spam_lifetime"
        KeyShieldDNCCount         = "shield.dnc_count"
        KeyShieldStopDet          = "shield.stop_detection_detail"
        KeyShieldDupGuard         = "shield.duplicate_guard"
        KeyShieldRecontact        = "shield.recontact_delay_full"

        // 3F fix additions — recommendations + other
        KeyShieldRec1             = "shield.recommend_1"
        KeyShieldRec2             = "shield.recommend_2"
        KeyShieldRec3             = "shield.recommend_3"
        KeyShieldClosingTriggers  = "shield.closing_triggers"
        KeyShieldConfigPerNiche   = "shield.config_per_niche"
        KeyShieldAutoMarkDeal     = "shield.auto_mark_deal"
        KeyShieldAutoMarkHot      = "shield.auto_mark_hot"
        KeyShieldAutoBlockStop    = "shield.auto_block_stop"
        KeyShieldManualOverride   = "shield.manual_override"
        KeyShieldHealthyCheck     = "shield.healthy_check"
        KeyShieldReplacing        = "shield.replacing_load"
        KeyShieldRecommendation   = "shield.recommendation"
        KeyShieldSlotAutoPaused   = "shield.slot_auto_paused"
        KeyShieldAllMsgsMoved     = "shield.all_msgs_moved"
        KeyShieldNothingToDo      = "shield.nothing_to_do"
        KeyShieldPossibility      = "shield.possibility"
        KeyShieldAction           = "shield.action"

        // 3F-2 fix additions — hardcoded string replacements
        KeyShieldPerDayDetail     = "shield.per_day_detail"
        KeyShieldReadyFmt         = "shield.ready_fmt"
        KeyShieldWarningZero      = "shield.warning_zero"
        KeyShieldWarningCountFmt  = "shield.warning_count_fmt"
        KeyShieldSentStat         = "shield.sent_stat"
        KeyShieldTooManyHour      = "shield.too_many_hour"
        KeyShieldAutoReduce       = "shield.auto_reduce"
        KeyShieldFlaggedMsg       = "shield.flagged_msg"
        KeyShieldStats7Day        = "shield.stats_7day"
        KeyShieldStatSent         = "shield.stat_sent"
        KeyShieldStatRespond      = "shield.stat_respond"
        KeyShieldStatFailed       = "shield.stat_failed"
        KeyShieldStatWarning      = "shield.stat_warning"
        KeyShieldHistoryLabel     = "shield.history_label"
        KeyShieldHealthScoreFmt   = "shield.health_score_fmt"
        KeyShieldHealthBelow50    = "shield.health_below_50"
        KeyShieldHealthUp5        = "shield.health_up_5"
        KeyShieldSectionAntiBan   = "shield.section_anti_ban"
        KeyShieldSectionSpamGuard = "shield.section_spam_guard"

        // 3F-3 fix additions — work hours & rate limiting demo strings
        KeyShieldTimezoneDemo     = "shield.timezone_demo"
        KeyShieldSendHoursDemo    = "shield.send_hours_demo"
        KeyShieldScrapeHoursDemo  = "shield.scrape_hours_demo"
        KeyShieldNowInWorkDemo    = "shield.now_in_work_demo"
        KeyShieldPerSlotDemo      = "shield.per_slot_demo"
        KeyShieldPerNumber        = "shield.per_number"
        KeyShieldPerNumberLabel   = "shield.per_number_label"
        KeyShieldPerLeadDetail    = "shield.per_lead_detail"
        KeyShieldStatFailedDemo   = "shield.stat_failed_demo"
        KeyShieldWaitingData     = "shield.waiting_data"

        // 3F-4 fix additions — pattern guard key-value labels
        KeyShieldPatternTemplateLabel = "shield.pattern_template_label"
        KeyShieldPatternTimeVarLabel  = "shield.pattern_time_var_label"
        KeyShieldPatternMsgVarLabel   = "shield.pattern_msg_var_label"
        KeyShieldPatternEmojiLabel    = "shield.pattern_emoji_label"
        KeyShieldPatternShuffleLabel  = "shield.pattern_shuffle_label"

        // 3F-4 fix additions — spam guard key-value labels
        KeyShieldSpamPerLeadLabel     = "shield.spam_per_lead_label"
        KeyShieldSpamLifetimeLabel    = "shield.spam_lifetime_label"
        KeyShieldDNCCountLabel        = "shield.dnc_count_label"
        KeyShieldStopDetLabel         = "shield.stop_det_label"
        KeyShieldDupGuardLabel        = "shield.dup_guard_label"
        KeyShieldRecontactLabel       = "shield.recontact_label"
        KeyShieldRecontactDelayFull   = "shield.recontact_delay_full"
        KeyShieldRecontactDelayDetail = "shield.recontact_delay_detail"

        // 3F-4 fix additions — slot detail
        KeyShieldPauseNumber          = "shield.pause_number"
)

// Settings screen (infra)
const (
        KeySettingsTitle        = "settings.title"
        KeySettingsAllInFiles   = "settings.all_in_files"
        KeySettingsConfigMain   = "settings.config_main"
        KeySettingsTheme        = "settings.theme"
        KeySettingsQueries      = "settings.queries"
        KeySettingsNicheFolder  = "settings.niche_folder"
        KeySettingsActiveConfig = "settings.active_config"
        KeySettingsActiveNiches = "settings.active_niches"
        KeySettingsWASlots      = "settings.wa_slots"
        KeySettingsWorkerPool   = "settings.worker_pool"
        KeySettingsArea         = "settings.area"
        KeySettingsWorkHours    = "settings.work_hours"
        KeySettingsRateLimit    = "settings.rate_limit"
        KeySettingsRotatorMode  = "settings.rotator_mode"
        KeySettingsAutopilot    = "settings.autopilot"
        KeySettingsEditConfig   = "settings.edit_config"
        KeySettingsReloadOK     = "settings.reload_success"
        KeySettingsChanges      = "settings.changes"
        KeySettingsApplied      = "settings.applied"
        KeySettingsOpenFile     = "settings.open_file"
        KeySettingsRevertBackup = "settings.revert_backup"
        KeySettingsOldConfig    = "settings.old_config"
        KeySettingsFixFirst     = "settings.fix_first"
        KeySettingsBackupNote   = "settings.backup_note"
        KeySettingsErrorCount   = "settings.error_count"
        KeySettingsReloadErr    = "settings.reload_error"
        KeySettingsAfterSave    = "settings.after_save"
        KeySettingsOpenEditor   = "settings.open_in_editor"
        KeySettingsEdit          = "settings.edit"

        // 3F fix addition
        KeySettingsBackDashboard = "settings.back_dashboard"

        // 3F-2 fix additions — hardcoded string replacements
        KeySettingsConfigPath    = "settings.config_path"
        KeySettingsBackupPath    = "settings.backup_path"

        // 3F-3 fix additions — file path strings
        KeySettingsThemePath     = "settings.theme_path"
        KeySettingsQueriesPath   = "settings.queries_path"
        KeySettingsNicheFolderPath = "settings.niche_folder_path"

        // 3F-5 fix additions — reload error view
        KeySettingsFilePath      = "settings.file_path"
)

// Guardrail screen (infra)
const (
        KeyGuardrailTitle          = "guardrail.title"
        KeyGuardrailCleanEmoji     = "guardrail.clean_emoji"
        KeyGuardrailAllClean       = "guardrail.all_clean"
        KeyGuardrailErrorEmoji     = "guardrail.error_emoji"
        KeyGuardrailHasErrors      = "guardrail.has_errors"
        KeyGuardrailWarningEmoji   = "guardrail.warning_emoji"
        KeyGuardrailHasWarnings    = "guardrail.has_warnings"
        KeyGuardrailRevalidating   = "guardrail.revalidating"
        KeyGuardrailFirstTime      = "guardrail.first_time"
        KeyGuardrailNoConfig       = "guardrail.no_config"
        KeyGuardrailNeed           = "guardrail.need"
        KeyGuardrailRelax          = "guardrail.relax"
        KeyGuardrailAutoGenerate   = "guardrail.auto_generate"
        KeyGuardrailSeeExample     = "guardrail.see_example"
        KeyGuardrailTemplatePath   = "guardrail.template_path"
        KeyGuardrailNiches         = "guardrail.niches"
        KeyGuardrailTemplates      = "guardrail.templates"
        KeyGuardrailAreas          = "guardrail.areas"
        KeyGuardrailReady          = "guardrail.ready"
        KeyGuardrailPaused         = "guardrail.paused"
        KeyGuardrailCanGo          = "guardrail.can_go"
        KeyGuardrailArmyReady      = "guardrail.army_ready"
        KeyGuardrailGo             = "guardrail.go"
        KeyGuardrailContinue       = "guardrail.continue"
        KeyGuardrailChecking       = "guardrail.checking"
        KeyGuardrailWaiting        = "guardrail.waiting"
        KeyGuardrailFixed          = "guardrail.fixed"
        KeyGuardrailPausedNiches   = "guardrail.paused_niches"
        KeyGuardrailDetailErrors   = "guardrail.detail_errors"
        KeyGuardrailErrorCountFmt  = "guardrail.error_count_fmt"
        KeyGuardrailLine           = "guardrail.line"
        KeyGuardrailOpenFirst      = "guardrail.open_first"
        KeyGuardrailOpenNext       = "guardrail.open_next"
        KeyGuardrailAllPassed      = "guardrail.all_passed"
        KeyGuardrailAutoResume     = "guardrail.auto_resume"
        KeyGuardrailArmyReadyAgain = "guardrail.army_ready_again"
        KeyGuardrailNotBlocking    = "guardrail.not_blocking"
        KeyGuardrailDeprecated     = "guardrail.deprecated"
        KeyGuardrailUseInstead     = "guardrail.use_instead"
        KeyGuardrailStillWorks     = "guardrail.still_works"
        KeyGuardrailSuggestion     = "guardrail.suggestion"
        KeyGuardrail1PerNiche      = "guardrail.1_per_niche"

        // 3F fix additions
        KeyGuardrailNichesReady     = "guardrail.niches_ready"
        KeyGuardrailTemplatesLoaded = "guardrail.templates_loaded"
        KeyGuardrailWorkersOK       = "guardrail.workers_ok"
        KeyGuardrailCanGas          = "guardrail.can_gas"
        KeyGuardrailNeedPrepare     = "guardrail.need_prepare"
        KeyGuardrailJustFill        = "guardrail.just_fill"
        KeyGuardrailPausedNichesFmt = "guardrail.paused_niches_fmt"
        KeyGuardrailStillRunning    = "guardrail.still_running"
        KeyGuardrailDetailWarning   = "guardrail.detail_warning"
        KeyGuardrailCheckingFile    = "guardrail.checking_file"
        KeyGuardrailOfDone          = "guardrail.of_done"
        KeyGuardrailPress1Next      = "guardrail.press_1_next"

        // 3F-2 fix additions — hardcoded string replacements
        KeyGuardrailConfigPathMain  = "guardrail.config_path_main"
        KeyGuardrailConfigPathDesc  = "guardrail.config_path_desc"
        KeyGuardrailNicheFolderPath = "guardrail.niche_folder_path"
        KeyGuardrailNicheFolderDesc = "guardrail.niche_folder_desc"
        KeyGuardrailTemplatePathLabel = "guardrail.template_path_label"
        KeyGuardrailWaitingLabel = "guardrail.waiting_label"
        KeyGuardrailFileSectionFmt = "guardrail.file_section_fmt"

        // Doc-alignment fix additions — key hint labels
        KeyGuardrailContinueAnyway = "guardrail.continue_anyway"
        KeyGuardrailQuitExit       = "guardrail.quit_exit"
        KeyGuardrailOpenFile       = "guardrail.open_file"

        // 3F-5 fix additions — still running reassurance
        KeyGuardrailStillRunningMsg = "guardrail.still_running_msg"

        // 3F-6 fix additions — i18n pluralization for error/warning counts
        KeyGuardrailErrorCountOne   = "guardrail.error_count_one"
        KeyGuardrailErrorCountMany  = "guardrail.error_count_many"
        KeyGuardrailWarningCountOne = "guardrail.warning_count_one"
        KeyGuardrailWarningCountMany = "guardrail.warning_count_many"
)

// License screen
const (
        KeyLicenseValidating = "license.validating"
        KeyLicenseValid      = "license.valid"
        KeyLicenseInvalid    = "license.invalid"

        // License screen — extended labels
        KeyLicenseTitle          = "license.title"
        KeyLicenseNeedsLicense   = "license.needs_license"
        KeyLicenseEnterKeyBelow  = "license.enter_key_below"
        KeyLicenseActionValidate = "license.action_validate"
        KeyLicenseActionBuyLicense = "license.action_buy"
        KeyLicenseActionExit     = "license.action_exit"
        KeyLicenseConnecting     = "license.connecting"
        KeyLicenseStoredAt       = "license.stored_at"
        KeyLicenseOneDevice      = "license.one_device"

        // License screen — validation step labels
        KeyLicenseCheckValidity   = "license.check_validity"
        KeyLicenseCheckDevice     = "license.check_device"
        KeyLicenseValidShort      = "license.valid_short"
        KeyLicenseConnected       = "license.connected"
        KeyLicenseValidDevice     = "license.valid_device"
        KeyLicenseKeySaved        = "license.key_saved"
        KeyLicenseReadyToRun      = "license.ready_to_run"
        KeyLicenseActionContinue  = "license.action_continue"
        KeyLicenseInvalidCheck    = "license.invalid_check"

        // License screen — invalid state labels
        KeyLicenseKeyNotMatch       = "license.key_not_match"
        KeyLicenseNoTypo            = "license.no_typo"
        KeyLicenseActionTryAgain    = "license.action_try_again"

        // License screen — expired state labels
        KeyLicenseExpiredLong       = "license.expired_long"
        KeyLicenseWorkersPaused     = "license.workers_paused"
        KeyLicenseRenewToContinue   = "license.renew_to_continue"
        KeyLicenseActionNewLicense  = "license.action_new_license"
        KeyLicenseActionBuyRenewal  = "license.action_buy_renewal"

        // License screen — conflict state labels
        KeyLicenseConflictLong       = "license.conflict_long"
        KeyLicenseOneDeviceExplain   = "license.one_device_explain"
        KeyLicenseMoveDevice         = "license.move_device"
        KeyLicenseActionDisconnect   = "license.action_disconnect"
        KeyLicenseForceDisconnectExplain = "license.force_disconnect_explain"

        // License screen — server error state labels
        KeyLicenseServerFail           = "license.server_fail"
        KeyLicenseServerReachable      = "license.server_reachable"
        KeyLicenseConnectionFailed     = "license.connection_failed"
        KeyLicenseHadValidBefore       = "license.had_valid_before"
        KeyLicenseOfflineGrace         = "license.offline_grace"
        KeyLicenseMustOnline           = "license.must_online"
        KeyLicenseGraceRemaining       = "license.grace_remaining"
        KeyLicenseActionOffline        = "license.action_offline"
        KeyLicenseActionRetry          = "license.action_retry"

        // License screen — field labels (i18n, not hardcoded)
        KeyLicenseLabelKey         = "license.label_key"
        KeyLicenseLabelDevice      = "license.label_device"
        KeyLicenseLabelExpires     = "license.label_expires"
        KeyLicenseLabelExpired     = "license.label_expired"
        KeyLicenseLabelThisDevice  = "license.label_this_device"
        KeyLicenseLabelOtherDevice = "license.label_other_device"
        KeyLicenseLabelLastActive  = "license.label_last_active"
        KeyLicenseUnknownState     = "license.unknown_state"
)

// Update screen
const (
        KeyUpdateAvailable   = "update.available"
        KeyUpdateDownloading = "update.downloading"
        KeyUpdateReady       = "update.ready"
        KeyUpgradeAvailable  = "upgrade.available"
        KeyUpgradeLicense    = "upgrade.license"

        // Update screen — detailed labels
        KeyUpdateNow                = "update.now"
        KeyUpdateLater              = "update.later"
        KeyUpdateRestartNote        = "update.restart_note"
        KeyUpdateDownloadingLabel   = "update.downloading_label"
        KeyUpdateETA                = "update.eta"
        KeyUpdateSource             = "update.source"
        KeyUpdateSpeed              = "update.speed"
        KeyUpdateWorkerNote         = "update.worker_note"
        KeyUpdateRestartApproval    = "update.restart_approval"
        KeyUpdateCancelDownload     = "update.cancel_download"
        KeyUpdateReadyInstall       = "update.ready_install"
        KeyUpdateDownloadComplete   = "update.download_complete"
        KeyUpdateChecksum           = "update.checksum"
        KeyUpdateVerified           = "update.verified"
        KeyUpdateBackup             = "update.backup"
        KeyUpdateChecksumPending    = "update.checksum_pending"
        KeyUpdateRestartPrompt      = "update.restart_prompt"
        KeyUpdateRestartDuration    = "update.restart_duration"
        KeyUpdateRestartNow         = "update.restart_now"
        KeyUpdateLaterKeep          = "update.later_keep"
        KeyUpdateSkipReminder       = "update.skip_reminder"
        KeyUpdateCurrentVersion     = "update.current_version"
        KeyUpdateNewVersion         = "update.new_version"
        KeyUpdateWhatsNew           = "update.whats_new"
        KeyUpdateMinorFree          = "update.minor_free"
        KeyUpdateLicenseValid       = "update.license_valid"

        // Upgrade screen — detailed labels
        KeyUpgradeMajorWarning      = "upgrade.major_warning"
        KeyUpgradeNeedsLicense      = "upgrade.needs_license"
        KeyUpgradeV1NotValid        = "upgrade.v1_not_valid"
        KeyUpgradeNotGreedy         = "upgrade.not_greedy"
        KeyUpgradeV1StillWorks      = "upgrade.v1_still_works"
        KeyUpgradeNoForced          = "upgrade.no_forced"
        KeyUpgradeV1License         = "upgrade.v1_license"
        KeyUpgradeV1StillRuns       = "upgrade.v1_still_runs"
        KeyUpgradeBuyV2             = "upgrade.buy_v2"
        KeyUpgradeViewDetails       = "upgrade.view_details"
        KeyUpgradeStayV1            = "upgrade.stay_v1"
        KeyUpgradeV2Title           = "upgrade.v2_title"
        KeyUpgradeV2EnterLicense    = "upgrade.v2_enter_license"
        KeyUpgradeV1StillActive     = "upgrade.v1_still_active"
        KeyUpgradeV2InvalidFallback = "upgrade.v2_invalid_fallback"
        KeyUpgradeValidateV2        = "upgrade.validate_v2"
        KeyUpgradeLicensePlaceholder = "upgrade.license_placeholder"
        KeyUpgradeLicenseKey        = "upgrade.license_key"
        KeyUpgradeLicenseExpires    = "upgrade.license_expires"
        KeyUpgradeLicenseStatus     = "upgrade.license_status"

        // License expired with upgrade
        KeyLicenseExpiredTitle      = "license_expired.title"
        KeyLicenseExpiredV1         = "license_expired.v1"
        KeyLicenseExpiredV1Stopped  = "license_expired.v1_stopped"
        KeyLicenseExpiredButV2      = "license_expired.but_v2"
        KeyLicenseExpiredV2Available = "license_expired.v2_available"
        KeyLicenseExpiredV2Features = "license_expired.v2_features"
        KeyLicenseExpiredV2Migrate  = "license_expired.v2_migrate"
        KeyLicenseExpiredRenewV1    = "license_expired.renew_v1"
        KeyLicenseExpiredUpgradeV2  = "license_expired.upgrade_v2"
        KeyLicenseExpiredNewLicense = "license_expired.new_license"
        KeyLicenseExpiredExit       = "license_expired.exit"

        // Startup check variant (background update check at boot)
        KeyStartupCheckTitle    = "startup_check.title"
        KeyStartupCheckLatest   = "startup_check.latest"
        KeyStartupCheckUpdate   = "startup_check.update_available"
        KeyStartupCheckUpgrade  = "startup_check.upgrade_available"
)

// Command palette overlay
const (
        KeyPaletteTitle  = "palette.title"
        KeyPaletteSearch = "palette.search"
        KeyPaletteEmpty  = "palette.empty"
        KeyPaletteFooter            = "palette.footer"
        KeyPaletteRecentHeader      = "palette.recent_header"
        KeyPaletteAllCommandsHeader = "palette.all_commands_header"
        KeyPaletteTryPrefix         = "palette.try_prefix"

        // Command palette — command names
        KeyPaletteCmdDashboard    = "palette.cmd_dashboard"
        KeyPaletteCmdLeads        = "palette.cmd_leads"
        KeyPaletteCmdSend         = "palette.cmd_send"
        KeyPaletteCmdWorkers      = "palette.cmd_workers"
        KeyPaletteCmdTemplates    = "palette.cmd_templates"
        KeyPaletteCmdShield       = "palette.cmd_shield"
        KeyPaletteCmdFollowup     = "palette.cmd_followup"
        KeyPaletteCmdSettings     = "palette.cmd_settings"
        KeyPaletteCmdHistory      = "palette.cmd_history"
        KeyPaletteCmdExplorer     = "palette.cmd_explorer"
        KeyPaletteCmdLicense      = "palette.cmd_license"
        KeyPaletteCmdUpdate       = "palette.cmd_update"
        KeyPaletteCmdScrapeAll    = "palette.cmd_scrape_all"
        KeyPaletteCmdPauseAll     = "palette.cmd_pause_all"
        KeyPaletteCmdFollowupAll  = "palette.cmd_followup_all"
        KeyPaletteCmdEditConfig   = "palette.cmd_edit_config"
        KeyPaletteCmdReloadConfig = "palette.cmd_reload_config"
        KeyPaletteCmdValidate     = "palette.cmd_validate"
        KeyPaletteCmdNerdStats    = "palette.cmd_nerd_stats"
        KeyPaletteCmdCmdPalette   = "palette.cmd_cmd_palette"
        KeyPaletteCmdShortcuts    = "palette.cmd_shortcuts"
        KeyPaletteCmdCompose      = "palette.cmd_compose"
        KeyPaletteCmdSearchLeads  = "palette.cmd_search_leads"
        KeyPaletteCmdExportCSV    = "palette.cmd_export_csv"
        KeyPaletteCmdMarkConverted = "palette.cmd_mark_converted"
        KeyPaletteCmdBlockLead    = "palette.cmd_block_lead"
        KeyPaletteCmdRecontact    = "palette.cmd_recontact"
        KeyPaletteCmdScrapeOne   = "palette.cmd_scrape_one"
        KeyPaletteCmdResumeAll   = "palette.cmd_resume_all"
        KeyPaletteCmdLogoutWA    = "palette.cmd_logout_wa"
        KeyPaletteCmdForceRetry  = "palette.cmd_force_retry"
)

// Notification severities
const (
        KeyNotifCritical    = "notif.critical"
        KeyNotifPositive    = "notif.positive"
        KeyNotifNeutral     = "notif.neutral"
        KeyNotifInformative = "notif.informative"

        // Notification type display titles
        KeyNotifResponseReceived       = "notif.response_received"
        KeyNotifResponseReceivedSuffix = "notif.response_received_suffix"
        KeyNotifMultiResponse          = "notif.multi_response"
        KeyNotifMultiResponseSuffix    = "notif.multi_response_suffix"
        KeyNotifScrapeComplete         = "notif.scrape_complete"
        KeyNotifBatchSendComplete      = "notif.batch_send_complete"
        KeyNotifWADisconnect           = "notif.wa_disconnect"
        KeyNotifWAFlag                 = "notif.wa_flag"
        KeyNotifHealthScoreDrop        = "notif.health_score_drop"
        KeyNotifDailyLimit             = "notif.daily_limit"
        KeyNotifStreakMilestone        = "notif.streak_milestone"
        KeyNotifConfigError            = "notif.config_error"
        KeyNotifValidationError        = "notif.validation_error"
        KeyNotifLicenseExpired         = "notif.license_expired"
        KeyNotifDeviceConflict         = "notif.device_conflict"
        KeyNotifFollowUpScheduled      = "notif.follow_up_scheduled"
        KeyNotifLeadCold               = "notif.lead_cold"
        KeyNotifUpdateAvailable        = "notif.update_available"
        KeyNotifUpgradeAvailable       = "notif.upgrade_available"

        // Notification action labels — per-type (doc-aligned)
        KeyNotifActionReply          = "notif.action_reply"
        KeyNotifActionLaterShort     = "notif.action_later_short"
        KeyNotifActionDismissShort   = "notif.action_dismiss_short"
        KeyNotifActionViewStats      = "notif.action_view_stats"
        KeyNotifActionNice           = "notif.action_nice"
        KeyNotifActionUpdateNow      = "notif.action_update_now"
        KeyNotifActionSkip           = "notif.action_skip"
        KeyNotifActionULater         = "notif.action_u_later"
        KeyNotifActionViewInfo       = "notif.action_view_info"
        KeyNotifActionUUpgrade       = "notif.action_u_upgrade"
        KeyNotifActionProcess        = "notif.action_process"
        KeyNotifAction1AutoOffer     = "notif.action_1_auto_offer"
        KeyNotifAction1Relogin       = "notif.action_1_relogin"
        KeyNotifActionOkProceed      = "notif.action_ok_proceed"

        // Notification action labels — shared
        KeyNotifActionEnterLicense    = "notif.action_enter_license"
        KeyNotifActionExit            = "notif.action_exit"
        KeyNotifActionViewLicense     = "notif.action_view_license"
        KeyNotifActionDisconnectOther = "notif.action_disconnect_other"
        KeyNotifActionViewError       = "notif.action_view_error"
        KeyNotifActionValidateAll     = "notif.action_validate_all"
        KeyNotifActionRelogin         = "notif.action_relogin"
        KeyNotifActionDismiss         = "notif.action_dismiss"
        KeyNotifActionViewShield      = "notif.action_view_shield"
        KeyNotifActionLetWaclaw       = "notif.action_let_waclaw"
        KeyNotifActionLater     = "notif.action_later"
        KeyNotifActionAutoOffer = "notif.action_auto_offer"
        KeyNotifActionUpdate    = "notif.action_update"
        KeyNotifActionUpgrade   = "notif.action_upgrade"
        KeyNotifActionView            = "notif.action_view"
        KeyNotifActionOk              = "notif.action_ok"

        // Notification body strings
        KeyNotifBodyScrapeComplete    = "notif.body_scrape_complete"
        KeyNotifBodyBatchSendComplete = "notif.body_batch_send_complete"
        KeyNotifBodyConfigError       = "notif.body_config_error"
        KeyNotifBodyDeviceConflict    = "notif.body_device_conflict"
        KeyNotifBodyWADisconnect      = "notif.body_wa_disconnect"
        KeyNotifBodyWAFlag            = "notif.body_wa_flag"
        KeyNotifBodyMultiResponse     = "notif.body_multi_response"
        KeyNotifBodyLicenseExpired    = "notif.body_license_expired"
        KeyNotifBodyHealthScoreDrop   = "notif.body_health_score_drop"
)

// Confirmation overlay
const (
        KeyConfirmProceed = "confirm.proceed"
        KeyConfirmCancel  = "confirm.cancel"

        // Confirmation type titles and details
        KeyConfirmBulkOfferTitle          = "confirm.bulk_offer_title"
        KeyConfirmBulkOfferDetail         = "confirm.bulk_offer_detail"
        KeyConfirmBulkDeleteTitle         = "confirm.bulk_delete_title"
        KeyConfirmBulkDeleteDetail        = "confirm.bulk_delete_detail"
        KeyConfirmBulkArchiveTitle        = "confirm.bulk_archive_title"
        KeyConfirmBulkArchiveDetail       = "confirm.bulk_archive_detail"
        KeyConfirmForceDisconnectTitle    = "confirm.force_disconnect_title"
        KeyConfirmForceDisconnectDetail   = "confirm.force_disconnect_detail"
)

// Shortcuts overlay
const (
        KeyShortcutsTitle        = "shortcuts.title"
        KeyShortcutsMove         = "shortcuts.move"
        KeyShortcutsPrimaryAction = "shortcuts.primary_action"
        KeyShortcutsPickOption   = "shortcuts.pick_option"
        KeyShortcutsSkip         = "shortcuts.skip"
        KeyShortcutsBackQuit     = "shortcuts.back_quit"
        KeyShortcutsPause        = "shortcuts.pause"
        KeyShortcutsSearch       = "shortcuts.search"
        KeyShortcutsValidateConfig = "shortcuts.validate_config"
        KeyShortcutsLicense      = "shortcuts.license"
        KeyShortcutsHistory      = "shortcuts.history"
        KeyShortcutsReload       = "shortcuts.reload"
        KeyShortcutsNerdStats    = "shortcuts.nerd_stats"
        KeyShortcutsCheckUpdate  = "shortcuts.check_update"
        KeyShortcutsCmdPalette   = "shortcuts.cmd_palette"
        KeyShortcutsPressAnyKey  = "shortcuts.press_any_key"
)

// Nerd stats overlay
const (
        KeyNerdStatsHeader    = "nerd_stats.header"
        KeyNerdStatsLogsHeader = "nerd_stats.logs_header"
)

// Validation overlay
const (
        KeyValidationErrorCount  = "validation.error_count"
        KeyValidationWarningCount = "validation.warning_count"
        KeyValidationAllValid     = "validation.all_valid"
)

// App
const (
        KeyAppNoScreen = "app.no_screen"
)

// Compose screen
const (
        KeyComposeDraftTitle    = "compose.draft_title"
        KeyComposePreviewTitle  = "compose.preview_title"
        KeyComposePickTitle     = "compose.pick_title"
        KeyComposeDraftHint     = "compose.draft_hint"
        KeyComposeSend          = "compose.send"
        KeyComposeSendSingle    = "compose.send_single"
        KeyComposeTab           = "compose.tab"
        KeyComposeCancel        = "compose.cancel"
        KeyComposeEdit          = "compose.edit"
        KeyComposeEditFirst     = "compose.edit_first"
        KeyComposeUse           = "compose.use"
        KeyComposeSelect        = "compose.select"
        KeyComposeTip           = "compose.tip"
        KeyComposeSendTo        = "compose.send_to"
        KeyComposeSnippetPath   = "compose.snippet_path"
        KeyComposeCategory      = "compose.category"
        KeyComposeSnippetSubtitle = "compose.snippet_subtitle"

        // Snippet category i18n keys (display-time lookup)
        KeyComposeCatSoftPitch  = "compose.cat_soft_pitch"
        KeyComposeCatFreeOffer  = "compose.cat_free_offer"
        KeyComposeCatMoveToCall = "compose.cat_move_to_call"
        KeyComposeCatDirectPrice = "compose.cat_direct_price"
        KeyComposeCatSendSample = "compose.cat_send_sample"
)

// History screen
const (
        KeyHistoryToday       = "history.today"
        KeyHistoryWeek        = "history.week"
        KeyHistoryDayDetail   = "history.day_detail"
        KeyHistoryTimeline    = "history.timeline"
        KeyHistorySummary     = "history.summary"
        KeyHistorySent        = "history.sent"
        KeyHistoryRespond     = "history.respond"
        KeyHistoryConvert     = "history.convert"
        KeyHistoryNewLeads    = "history.new_leads"
        KeyHistoryScrape      = "history.scrape"
        KeyHistoryBestDay     = "history.best_day"
        KeyHistoryBestLabel   = "history.best_label"
        KeyHistoryPrimeTime   = "history.prime_time"
        KeyHistoryConvRate    = "history.conv_rate"
        KeyHistoryMsgsPerDay  = "history.msgs_per_day"
        KeyHistoryRespPerDay  = "history.resp_per_day"
        KeyHistoryViewDetail  = "history.view_detail"
        KeyHistoryViewEvent   = "history.view_event"
        KeyHistoryPrevDay     = "history.prev_day"
        KeyHistoryNextDay     = "history.next_day"
        KeyHistoryWeekTotal   = "history.week_total"
        KeyHistoryAvgRespTime = "history.avg_resp_time"
        KeyHistoryWeekend     = "history.weekend"
        KeyHistorySelectDay    = "history.select_day"

        // History — w/t key binding labels
        KeyHistorySwitchWeek  = "history.switch_week"
        KeyHistorySwitchToday = "history.switch_today"
)

// Follow-up screen
const (
        KeyFollowUpTitle       = "followup.title"
        KeyFollowUpAutoRunning = "followup.auto_running"
        KeyFollowUpSending     = "followup.sending"
        KeyFollowUpQueueToday  = "followup.queue_today"
        KeyFollowUpTotal       = "followup.total"
        KeyFollowUpCold        = "followup.cold"
        KeyFollowUpAutoAll     = "followup.auto_all"
        KeyFollowUpAutoNote    = "followup.auto_note"
        KeyFollowUpEmpty       = "followup.empty"
        KeyFollowUpEmptyDesc   = "followup.empty_desc"
        KeyFollowUpColdList    = "followup.cold_list"
        KeyFollowUpColdDesc    = "followup.cold_desc"
        KeyFollowUpColdWarning = "followup.cold_warning"
        KeyFollowUpRecontact   = "followup.recontact"
        KeyFollowUpRecontactDesc = "followup.recontact_desc"
        KeyFollowUpRateInfo    = "followup.rate_info"
        KeyFollowUpVariantNote = "followup.variant_note"
        KeyFollowUpSkipWait       = "followup.skip_wait"
        KeyFollowUpSendingNow     = "followup.sending_now"
        KeyFollowUpWaitNext       = "followup.wait_next"
        KeyFollowUpAutoSchedule   = "followup.auto_schedule"
        KeyFollowUpNoDuplicate    = "followup.no_duplicate"
        KeyFollowUpDiffVariant   = "followup.diff_variant"
        KeyFollowUpRotation      = "followup.rotation"
        KeyFollowUpMessagesLabel = "followup.messages_label"
        KeyFollowUpIceBreakerUnanswered = "followup.ice_breaker_unanswered"
        KeyFollowUpColdDetail    = "followup.cold_detail"
        KeyFollowUpRecontactTemplate = "followup.recontact_template"
        KeyFollowUpRecontactTone = "followup.recontact_tone"
        KeyFollowUpRecontactVibes = "followup.recontact_vibes"
        KeyFollowUpCanRecontact  = "followup.can_recontact"
        KeyFollowUpAfterThat     = "followup.after_that"
        KeyFollowUpAfterThatSilent = "followup.after_that_silent"
        KeyFollowUpEditFirst     = "followup.edit_first"
        KeyFollowUpTabNiche      = "followup.tab_niche"
        KeyFollowUpSendFinal     = "followup.send_final"

        // NEW keys (3G fix additions)
        KeyFollowUpFU1Detail         = "followup.fu1_detail"
        KeyFollowUpFU2Detail         = "followup.fu2_detail"
        KeyFollowUpColdDashboard     = "followup.cold_dashboard"
        KeyFollowUpTotalCold         = "followup.total_cold"
        KeyFollowUpVariantPreview1   = "followup.variant_preview_1"
        KeyFollowUpVariantPreview2   = "followup.variant_preview_2"
        KeyFollowUpVariantPreview3   = "followup.variant_preview_3"
        KeyFollowUpVariantManualOnly = "followup.variant_manual_only"
        KeyFollowUpRatePerHour       = "followup.rate_per_hour"
        KeyFollowUpEmptyUnanswered   = "followup.empty_unanswered"
        KeyFollowUpDaysAgo           = "followup.days_ago"
        KeyFollowUpOfferSent         = "followup.offer_sent"
        KeyFollowUpNextLabel         = "followup.next_label"
        KeyFollowUpVariantPerFu      = "followup.variant_per_fu"
        KeyFollowUpVariantLabel      = "followup.variant_label"
        KeyFollowUpSnippetSubtitle   = "followup.snippet_subtitle"
        KeyFollowUpViewCold          = "followup.view_cold"

        // Follow-up — label and noun keys (i18n fix)
        KeyFollowUpFU1Label = "followup.fu1_label"
        KeyFollowUpFU2Label = "followup.fu2_label"
        KeyFollowUpLeadNoun = "followup.lead_noun"
        KeyFollowUpSlot     = "followup.slot"
        KeyFollowUpOthers        = "followup.others"
        KeyFollowUpViewDetail     = "followup.view_detail"
        KeyFollowUpSendRecontact  = "followup.send_recontact"

        // Follow-up — 3G fix additions
        KeyFollowUpEmptyUnansweredIB = "followup.empty_unanswered_ib"
        KeyFollowUpColdCountLabel    = "followup.cold_count_label"
)

// Action labels
const (
        KeyActionArchiveAll = "action.archive_all"
)

// Data screens — Leads DB
const (
        KeyDataLeadsTitle        = "data.leads_title"
        KeyDataFilterAll         = "data.filter_all"
        KeyDataFilterSearch      = "data.filter_search"
        KeyDataRecent            = "data.recent"
        KeyDataMove              = "data.move"
        KeyDataViewDetail        = "data.view_detail"
        KeyDataSendOfferAll      = "data.send_offer_all"
        KeyDataScore             = "data.score"
        KeyDataTimeline          = "data.timeline"
        KeyDataStatus            = "data.status"
        KeyDataNotContacted      = "data.not_contacted"
        KeyDataFollowUpNext      = "data.follow_up_next"
        KeyDataAutoSchedule      = "data.auto_schedule"
        KeyDataColdDesc          = "data.cold_desc"
        KeyDataColdOption        = "data.cold_option"
        KeyDataColdAuto          = "data.cold_auto"
        KeyDataNoResponseYet     = "data.no_response_yet"
        KeyDataResponse          = "data.response"
        KeyDataOfferNotSent      = "data.offer_not_sent"
        KeyDataNotReplied        = "data.not_replied"
        KeyDataSendOffer         = "data.send_offer"
        KeyDataCustomReply       = "data.custom_reply"
        KeyDataLater             = "data.later"
        KeyDataMarkConvert       = "data.mark_convert"
        KeyDataArchive           = "data.archive_lead"
        KeyDataBlock             = "data.block"
        KeyDataSendFollowUp      = "data.send_followup"
        KeyDataMarkCold          = "data.mark_cold"
        KeyDataSendIceBreaker    = "data.send_icebreaker"
        KeyDataLastFollowUp      = "data.last_followup"
        KeyDataConversionRevenue = "data.conversion_revenue"
        KeyDataDuration          = "data.duration"
        KeyDataTemplate          = "data.template"
        KeyDataWorker            = "data.worker"
        KeyDataIceBreakerSent    = "data.ice_breaker_sent"
        KeyDataResponseReceived  = "data.response_received"
        KeyDataOfferSent         = "data.offer_sent"
        KeyDataMarkConvertAction = "data.mark_convert_action"
        KeyDataWaitingOffer      = "data.waiting_offer"
        KeyDataContactCount      = "data.contact_count"
        KeyDataFollowUpCount     = "data.followup_count"
        KeyDataTotal             = "data.total"
        KeyDataFiltered          = "data.filtered"
        KeyDataView              = "data.view"

        // Phase display labels — used in filter headers and badge text.
        KeyDataIceBreakerLabel = "data.ice_breaker_label"
        KeyDataFollowUp1Label  = "data.follow_up_1_label"
        KeyDataFollowUp2Label  = "data.follow_up_2_label"
        KeyDataNegativeLabel   = "data.negative_label"
        KeyDataArchivedLabel   = "data.archived_label"
        KeyDataAutoReplyLabel  = "data.auto_reply_label"
        KeyDataColdLabel       = "data.cold_label"
        KeyDataOfferSentLabel  = "data.offer_sent_label"
        KeyDataConvertLabel    = "data.convert_label"
        KeyDataFailedLabel     = "data.failed_label"
        KeyDataBlockedLabel    = "data.blocked_label"

        // Context-line labels.
        KeyDataIceBreakerTime = "data.ice_breaker_time"
        KeyDataOfferTime      = "data.offer_time"

        // Status indicators.
        KeyDataHasWebsite   = "data.has_website"
        KeyDataNoWebsite    = "data.no_website"
        KeyDataHasInstagram = "data.has_instagram"
        KeyDataNoInstagram  = "data.no_instagram"
        KeyDataReviewsLabel = "data.reviews_label"
        KeyDataPhotosCount  = "data.photos_count"
        KeyDataStillQuiet    = "data.still_quiet"
        KeyDataDueToday      = "data.due_today"
        KeyDataCategoryNote  = "data.category_note"

        // Inline labels — formerly hardcoded in view functions.
        KeyDataFilterLabel     = "data.filter_label"       // "filter:" — list view filter line
        KeyDataIceBreakerColon = "data.ice_breaker_colon"  // "ice breaker:" — detail views
        KeyDataFollowUpNum     = "data.follow_up_num"      // "follow-up %d:" — cold view
        KeyDataNiche           = "data.niche"              // "niche:" — converted lead detail
        KeyDataTrophy          = "data.trophy"             // "🏆" — conversion revenue prefix
        KeyDataErrorMark       = "data.error_mark"         // "✗" — error marker in validation
)

// Data screens — Template Manager
const (
        KeyTemplateTitle            = "template.title"
        KeyTemplateNiche            = "template.niche"
        KeyTemplateIceBreaker       = "template.ice_breaker"
        KeyTemplateOffer            = "template.offer"
        KeyTemplateRecommended      = "template.recommended"
        KeyTemplateSelect           = "template.select"
        KeyTemplatePreviewAction    = "template.preview_action"
        KeyTemplateNew              = "template.new"
        KeyTemplateEdit             = "template.edit"
        KeyTemplateEditTitle        = "template.edit_title"
        KeyTemplateUseThis          = "template.use_this"
        KeyTemplateEditFile         = "template.edit_file"
        KeyTemplateSavedAsFile      = "template.saved_as_file"
        KeyTemplateOpenInEditor     = "template.open_in_editor"
        KeyTemplateAfterSave        = "template.after_save"
        KeyTemplateErrorTitle       = "template.error_title"
        KeyTemplateFileEmpty        = "template.file_empty"
        KeyTemplateMinPlaceholder   = "template.min_placeholder"
        KeyTemplateUnknownPlaceholder   = "template.unknown_placeholder"
        KeyTemplateAvailablePlaceholders = "template.available_placeholders"
        KeyTemplateEncodingError    = "template.encoding_error"
        KeyTemplateEncodingHint     = "template.encoding_hint"
        KeyTemplateWorkerPaused     = "template.worker_paused"
        KeyTemplateIceBreakerReq    = "template.ice_breaker_required"
        KeyTemplateOpenFile         = "template.open_file"
        KeyTemplateViewPlaceholders = "template.view_placeholders"
        KeyTemplatePlaceholderLabel = "template.placeholder_label"
        KeyTemplatePreviewLabel     = "template.preview_label"
        KeyTemplateReload           = "template.reload"
        KeyTemplateLine             = "template.line"              // "baris" — template error line label
        KeyTemplateWorkerPausedFmt  = "template.worker_paused_fmt" // full "worker %s paused..." sentence
)

// Review screen
const (
        KeyReviewWaiting     = "review.waiting"
        KeyReviewTitle       = "review.title"
        KeyReviewNoWeb       = "review.no_web"
        KeyReviewNoIG        = "review.no_ig"
        KeyReviewWAReg       = "review.wa_reg"
        KeyReviewQA          = "review.qa"
        KeyReviewSkip        = "review.skip"
        KeyReviewDetail      = "review.detail"
        KeyReviewBlock       = "review.block"
        KeyReviewQueued      = "review.queued"
        KeyReviewSkipped     = "review.skipped"
        KeyReviewRemaining   = "review.remaining"
        KeyReviewBlocked     = "review.blocked"
        KeyReviewMove        = "review.move"
        KeyReviewDone        = "review.done"
        KeyReviewPhotos      = "review.photos"
        KeyReviewGoogleSrch  = "review.google_search"
        KeyReviewScore       = "review.score"
        KeyReviewHistory     = "review.history"
        KeyReviewPickVar     = "review.pick_variant"
        KeyReviewPreviewVar  = "review.preview_variant"
        KeyReviewUseThis     = "review.use_this"
        KeyReviewChangeVar   = "review.change_variant"
        KeyReviewCancel      = "review.cancel"
        KeyReviewComplete    = "review.complete"
        KeyReviewAutoReview  = "review.auto_review"
        KeyReviewDashLink    = "review.dash_link"
        KeyReviewOtherVar    = "review.other_variant"
        KeyReviewApproveAction = "review.approve_action"
        KeyReviewApproveShort  = "review.approve_short"
        KeyReviewFollowUp      = "review.follow_up"
        KeyReviewExit          = "review.exit"

        KeyReviewGasQueue        = "review.gas_queue"
        KeyReviewAutoSendTiming  = "review.auto_send_timing"
        KeyReviewAutoPilotRelax  = "review.auto_pilot_relax"
        KeyReviewSkipBlock       = "review.skip_block"
        KeyReviewBack            = "review.back"
        KeyReviewPickTemplate    = "review.pick_template"
)

// Shared word tokens (display units)
const (
        KeyWordReviews      = "word.reviews"
        KeyWordActiveIG     = "word.active_ig"
        KeyWordThisHour     = "word.this_hour"
        KeyWordFrom         = "word.from"
        KeyWordLater        = "word.later"
        KeyWordDisconnected = "word.disconnected"
)

// General
const (
        KeyGeneralLoading = "general.loading"
        KeyGeneralRetry   = "general.retry"
        KeyGeneralSkip    = "general.skip"
        KeyGeneralArchive = "general.archive"
        KeyGeneralDelete  = "general.delete"
        KeyGeneralSave    = "general.save"
        KeyGeneralClose   = "general.close"
        KeyGeneralOff     = "general.off"
        KeyGeneralOn      = "general.on"
)

// AllKeys returns every defined i18n key constant for validation purposes.
func AllKeys() []string {
        return []string{
                // Navigation
                KeyLabelProceed, KeyLabelCancel, KeyLabelBack, KeyLabelPause,
                KeyLabelHelp, KeyLabelRefresh, KeyLabelSearch, KeyLabelValidate,
                KeyLabelLicense, KeyLabelHistory, KeyLabelUpdate, KeyLabelNerd,
                KeyLabelEdit, KeyLabelOpen, KeyLabelReload,
                // Status
                KeyStatusActive, KeyStatusError, KeyStatusPaused, KeyStatusSent,
                KeyStatusDelivered, KeyStatusRead, KeyStatusIdle, KeyStatusScraping,
                KeyStatusSending, KeyStatusQualifying, KeyStatusQueuing,
                KeyStatusHealthy, KeyStatusFlagged, KeyStatusCooldown, KeyStatusDone,
                KeyStatusChecking, KeyStatusWaiting, KeyStatusWarning, KeyStatusStarting,
                // Celebration
                KeyLabelJackpot, KeyLabelDeal,
                // Shield levels
                KeyShieldHealthy, KeyShieldWarning, KeyShieldDanger,
                // Badges
                KeyBadgeNew, KeyBadgeCold, KeyBadgeConverted, KeyBadgeResponded, KeyBadgeFailed,
                // Session
                KeySessionEnd,
                // Boot
                KeyBootTagline, KeyBootReturning, KeyBootFirstTime,
                // Login
                KeyLoginQRWait, KeyLoginScanned, KeyLoginSuccess, KeyLoginExpired, KeyLoginFailed,
                // Niche
                KeyNicheSelect, KeyNicheSelected, KeyNicheCustom, KeyNicheFilters, KeyNicheConfigErr,
                KeyNicheExplorerTitle, KeyNicheExplorerSubtitle, KeyNicheExplorerPopular,
                KeyNicheExplorerSearching, KeyNicheExplorerResults,
                KeyNicheExplorerSource, KeyNicheExplorerGenConfig, KeyNicheExplorerGenSuccess,
                KeyNicheExplorerGenProgress, KeyNicheExplorerAreaAuto, KeyNicheExplorerEditFile,
                KeyNicheExplorerReload, KeyNicheExplorerParallel,
                KeyNicheNicheIs, KeyNicheCustomDir, KeyNicheCustomMin, KeyNicheCustomExample,
                KeyNicheCustomReady, KeyNicheErrPaused, KeyNicheErrOtherOK, KeyNicheProblems,
                KeyNicheTargets, KeyNicheAreaCount, KeyNicheFilterDefault, KeyNicheTemplateGen,
                KeyNicheJustRight, KeyNicheMoreArea, KeyNicheCanEdit, KeyNicheWorkerParallel,
                KeyNicheScrapeOwn,
                // Scrape
                KeyScrapeActive, KeyScrapeComplete, KeyScrapeEmpty, KeyScrapeError,
                // Send
                KeySendActive, KeySendPaused, KeySendOffHours, KeySendRateLimit, KeySendDailyLimit, KeySendFailed,
                // Monitor
                KeyMonitorLive, KeyMonitorIdle, KeyMonitorNight, KeyMonitorError,
                KeyMonitorWAConnected, KeyMonitorNicheActive, KeyMonitorWARotator,
                KeyMonitorWorkerPool, KeyMonitorRecentActivity, KeyMonitorAutoPilot,
                KeyMonitorIdleBody,
                KeyMonitorQueued, KeyMonitorMsgsSent, KeyMonitorResponses,
                KeyMonitorCanMinimize, KeyMonitorWorkersDesc, KeyMonitorWillNotify,
                KeyMonitorViewDetail, KeyMonitorNightMode, KeyMonitorWorkHours,
                KeyMonitorNow,
                KeyMonitorSenderPaused, KeyMonitorScraperRunning, KeyMonitorAutoResume,
                KeyMonitorDaySummary, KeyMonitorConverts, KeyMonitorPerNiche,
                KeyMonitorBestTime, KeyMonitorWADisconnected, KeyMonitorSlotDisconnected,
                KeyMonitorScraperOK, KeyMonitorDatabaseOK, KeyMonitorSlotsActive,
                KeyMonitorSlotPending, KeyMonitorAutoReconnect, KeyMonitorRelogin,
                KeyMonitorNoData, KeyMonitorStartScrape, KeyMonitorSetupNiche,
                KeyMonitorStartSend, KeyMonitorTipStartScrape, KeyMonitorPending,
                KeyMonitorAutoOfferAll, KeyMonitorAutoPerNiche, KeyMonitorOfferPerNiche,
                KeyMonitorToday, KeyMonitorCooldown, KeyMonitorResponded,
                KeyMonitorLeadsFound, KeyMonitorThisWeek, KeyMonitorConversionRate,
                KeyMonitorBestDay,
                KeyMonitorNavLeads, KeyMonitorNavMessages, KeyMonitorNavWorkers,
                KeyMonitorNavTemplate, KeyMonitorNavAntiban, KeyMonitorNavFollowup,
                KeyMonitorNavSettings, KeyMonitorScrapeAll,
                // Response
                KeyResponsePositive, KeyResponseCurious, KeyResponseNegative, KeyResponseMaybe, KeyResponseAuto,
                KeyResponseProcessOne,
                KeyResponseGotReply, KeyResponseSendOffer, KeyResponseReplyCustom, KeyResponseLater,
                KeyResponseSendInfo, KeyResponseGotReplyPlain, KeyResponseMarkInvalid,
                KeyResponseFollowUpLater, KeyResponseSkipIt,
                KeyResponseGotReplyAuto, KeyResponseAutoReplySkip, KeyResponseStillFollowUp,
                KeyResponseStopDetected, KeyResponseAutoAddedDNC,
                KeyResponseBlockConfirm, KeyResponseBlockAllNiche, KeyResponseCancelBlock,
                KeyResponseDealDetected, KeyResponseAutoDeal, KeyResponseConfirmDeal,
                KeyResponseNotDeal, KeyResponseReplyFirst,
                KeyResponseHotLead, KeyResponseAutoPrioritize, KeyResponseOfferNow,
                KeyResponseOfferPreview, KeyResponseSend, KeyResponseChangeTemplate,
                KeyResponseEditFirst, KeyResponseMultiQueue,
                KeyResponsePositif, KeyResponseCuriousBadge, KeyResponseAutoReplyBadge,
                KeyResponseAutoOfferPos, KeyResponseAutoPerType,
                KeyResponseConversion, KeyResponseClosingTrigger, KeyResponseAutoDealShort,
                // Conversion
                KeyConversionMarkConverted, KeyConversionFromPipeline,
                KeyConversionTimeTaken, KeyConversionOrdinal,
                KeyConversionTrophyWeek, KeyConversionRevenueWeek,
                // Settings (legacy)
                KeySettingsOverview, KeySettingsReload, KeySettingsError,
                // Guardrail (legacy)
                KeyGuardrailClean, KeyGuardrailErrors, KeyGuardrailWarnings,
                // Workers (infra)
                KeyWorkersTitle, KeyWorkersActive, KeyWorkersIdle, KeyWorkersPausedStatus,
                KeyWorkersDetail, KeyWorkersAddNiche, KeyWorkersScrape, KeyWorkersReview,
                KeyWorkersQueue, KeyWorkersSend, KeyWorkersArea, KeyWorkersTemplate,
                KeyWorkersTotalPipeline, KeyWorkersParallel, KeyWorkersPauseWorker,
                KeyWorkersForceScrape, KeyWorkersViewLeads, KeyWorkersResume, KeyWorkersDelete,
                KeyWorkersCollected, KeyWorkersAutoResume, KeyWorkersMoreNiches, KeyWorkersIndependent,
                KeyWorkersQuery, KeyWorkersDone, KeyWorkersPassed, KeyWorkersSkipped,
                KeyWorkersDuplicates, KeyWorkersLowRating, KeyWorkersNext, KeyWorkersToday,
                KeyWorkersResponseRate, KeyWorkersConversionRate, KeyWorkersAvgRespond,
                KeyWorkersPerforma,
                // Anti-ban shield (infra)
                KeyShieldTitle, KeyShieldAllSafe, KeyShieldHasWarning, KeyShieldDangerLabel,
                KeyShieldHealthScore, KeyShieldWARotator, KeyShieldNumbers,
                KeyShieldSlotActive, KeyShieldSlotCooldown, KeyShieldSlotFlagged,
                KeyShieldRateLimiting, KeyShieldPerSlot, KeyShieldPerDay, KeyShieldPerLead,
                KeyShieldDailyBudget, KeyShieldWorkHoursGuard, KeyShieldTimezone,
                KeyShieldSendHours, KeyShieldScrapeHours, KeyShieldNow, KeyShieldInWorkHours,
                KeyShieldPatternGuard, KeyShieldTemplateRotation, KeyShieldSpamGuard,
                KeyShieldBanRiskScore, KeyShieldRiskLow, KeyShieldRiskMedium, KeyShieldRiskHigh,
                KeyShieldIndicators, KeyShieldEvenDist, KeyShieldCooldownOk,
                KeyShieldTemplateVaried, KeyShieldNoOverload, KeyShieldWorkHoursOk,
                KeyShieldSpamGuardActive, KeyShieldDNCRespected, KeyShieldSlotDetail,
                KeyShieldStatistics, KeyShield7Days, KeyShieldConfig, KeyShieldEditConfig,
                KeyShieldAutoPause, KeyShieldAutoRecover, KeyShieldAddNumber, KeyShieldPauseSend,
                KeyShieldAllSafeEmoji, KeyShieldWarningEmoji, KeyShieldDangerEmoji,
                KeyShieldLetItBe, KeyShieldAlreadyMoved,
                KeyShieldHour, KeyShieldCooldown, KeyShieldTotal,
                KeyShieldStatus, KeyShieldHealthyStatus, KeyShieldReady,
                KeyShieldWaitingData,
                // Settings (infra)
                KeySettingsTitle, KeySettingsAllInFiles, KeySettingsConfigMain,
                KeySettingsTheme, KeySettingsQueries, KeySettingsNicheFolder,
                KeySettingsActiveConfig, KeySettingsActiveNiches, KeySettingsWASlots,
                KeySettingsWorkerPool, KeySettingsArea, KeySettingsWorkHours,
                KeySettingsRateLimit, KeySettingsRotatorMode, KeySettingsAutopilot,
                KeySettingsEditConfig, KeySettingsReloadOK, KeySettingsChanges,
                KeySettingsApplied, KeySettingsOpenFile, KeySettingsRevertBackup,
                KeySettingsOldConfig, KeySettingsFixFirst, KeySettingsBackupNote,
                KeySettingsErrorCount, KeySettingsReloadErr, KeySettingsAfterSave,
                KeySettingsOpenEditor, KeySettingsEdit,
                // Guardrail (infra)
                KeyGuardrailTitle, KeyGuardrailCleanEmoji, KeyGuardrailAllClean,
                KeyGuardrailErrorEmoji, KeyGuardrailHasErrors, KeyGuardrailWarningEmoji,
                KeyGuardrailHasWarnings, KeyGuardrailRevalidating, KeyGuardrailFirstTime,
                KeyGuardrailNoConfig, KeyGuardrailNeed, KeyGuardrailRelax,
                KeyGuardrailAutoGenerate, KeyGuardrailSeeExample, KeyGuardrailTemplatePath,
                KeyGuardrailNiches, KeyGuardrailTemplates, KeyGuardrailAreas,
                KeyGuardrailReady, KeyGuardrailPaused, KeyGuardrailCanGo,
                KeyGuardrailArmyReady, KeyGuardrailGo, KeyGuardrailContinue,
                KeyGuardrailChecking, KeyGuardrailWaiting, KeyGuardrailFixed,
                KeyGuardrailPausedNiches, KeyGuardrailDetailErrors, KeyGuardrailLine,
                KeyGuardrailOpenFirst, KeyGuardrailOpenNext, KeyGuardrailAllPassed,
                KeyGuardrailAutoResume, KeyGuardrailArmyReadyAgain, KeyGuardrailNotBlocking,
                KeyGuardrailDeprecated, KeyGuardrailUseInstead, KeyGuardrailStillWorks,
                KeyGuardrailSuggestion, KeyGuardrail1PerNiche,
                // License
                KeyLicenseValidating, KeyLicenseValid, KeyLicenseInvalid,
                KeyLicenseTitle, KeyLicenseNeedsLicense, KeyLicenseEnterKeyBelow,
                KeyLicenseActionValidate, KeyLicenseActionBuyLicense, KeyLicenseActionExit,
                KeyLicenseConnecting, KeyLicenseStoredAt, KeyLicenseOneDevice,
                KeyLicenseCheckValidity, KeyLicenseCheckDevice, KeyLicenseValidShort,
                KeyLicenseConnected, KeyLicenseValidDevice, KeyLicenseKeySaved,
                KeyLicenseReadyToRun, KeyLicenseActionContinue, KeyLicenseInvalidCheck,
                KeyLicenseKeyNotMatch, KeyLicenseNoTypo, KeyLicenseActionTryAgain,
                KeyLicenseExpiredLong, KeyLicenseWorkersPaused, KeyLicenseRenewToContinue,
                KeyLicenseActionNewLicense, KeyLicenseActionBuyRenewal,
                KeyLicenseConflictLong, KeyLicenseOneDeviceExplain, KeyLicenseMoveDevice,
                KeyLicenseActionDisconnect, KeyLicenseForceDisconnectExplain,
                KeyLicenseServerFail, KeyLicenseServerReachable, KeyLicenseConnectionFailed,
                KeyLicenseHadValidBefore, KeyLicenseOfflineGrace, KeyLicenseMustOnline,
                KeyLicenseGraceRemaining, KeyLicenseActionOffline, KeyLicenseActionRetry,
                KeyLicenseLabelKey, KeyLicenseLabelDevice, KeyLicenseLabelExpires,
                KeyLicenseLabelExpired, KeyLicenseLabelThisDevice, KeyLicenseLabelOtherDevice,
                KeyLicenseLabelLastActive,
                KeyLicenseUnknownState,
                // Update
                KeyUpdateAvailable, KeyUpdateDownloading, KeyUpdateReady,
                KeyUpgradeAvailable, KeyUpgradeLicense,
                KeyUpdateNow, KeyUpdateLater, KeyUpdateRestartNote,
                KeyUpdateDownloadingLabel, KeyUpdateETA, KeyUpdateSource,
                KeyUpdateWorkerNote, KeyUpdateRestartApproval, KeyUpdateCancelDownload,
                KeyUpdateReadyInstall, KeyUpdateDownloadComplete, KeyUpdateChecksum,
                KeyUpdateVerified, KeyUpdateBackup, KeyUpdateChecksumPending, KeyUpdateRestartPrompt,
                KeyUpdateRestartDuration, KeyUpdateRestartNow, KeyUpdateLaterKeep,
                KeyUpdateSkipReminder, KeyUpdateCurrentVersion, KeyUpdateNewVersion,
                KeyUpdateWhatsNew, KeyUpdateMinorFree, KeyUpdateLicenseValid,
                KeyUpgradeMajorWarning, KeyUpgradeNeedsLicense, KeyUpgradeV1NotValid,
                KeyUpgradeNotGreedy, KeyUpgradeV1StillWorks, KeyUpgradeNoForced,
                KeyUpgradeV1License, KeyUpgradeV1StillRuns, KeyUpgradeBuyV2,
                KeyUpgradeViewDetails, KeyUpgradeStayV1, KeyUpgradeV2Title,
                KeyUpgradeV2EnterLicense, KeyUpgradeV1StillActive, KeyUpgradeV2InvalidFallback,
                KeyUpgradeValidateV2,
                KeyUpgradeLicensePlaceholder, KeyUpgradeLicenseKey, KeyUpgradeLicenseExpires, KeyUpgradeLicenseStatus,
                KeyLicenseExpiredTitle, KeyLicenseExpiredV1, KeyLicenseExpiredV1Stopped,
                KeyLicenseExpiredButV2, KeyLicenseExpiredV2Available, KeyLicenseExpiredV2Features,
                KeyLicenseExpiredV2Migrate, KeyLicenseExpiredRenewV1, KeyLicenseExpiredUpgradeV2,
                KeyLicenseExpiredNewLicense, KeyLicenseExpiredExit,
                KeyStartupCheckTitle, KeyStartupCheckLatest, KeyStartupCheckUpdate, KeyStartupCheckUpgrade,
                KeyUpdateSpeed,
                // Command palette
                KeyPaletteTitle, KeyPaletteSearch, KeyPaletteEmpty,
                // Notification severities
                KeyNotifCritical, KeyNotifPositive, KeyNotifNeutral, KeyNotifInformative,
                // Confirmation
                KeyConfirmProceed, KeyConfirmCancel,
                // Shortcuts
                KeyShortcutsTitle,
                // Compose
                KeyComposeDraftTitle, KeyComposePreviewTitle, KeyComposePickTitle,
                KeyComposeDraftHint, KeyComposeSend, KeyComposeSendSingle, KeyComposeTab,
                KeyComposeCancel, KeyComposeEdit, KeyComposeUse,
                KeyComposeSelect, KeyComposeTip, KeyComposeSendTo,
                KeyComposeSnippetPath, KeyComposeCategory,
                // History
                KeyHistoryToday, KeyHistoryWeek, KeyHistoryDayDetail,
                KeyHistoryTimeline, KeyHistorySummary,
                KeyHistorySent, KeyHistoryRespond, KeyHistoryConvert,
                KeyHistoryNewLeads, KeyHistoryScrape,
                KeyHistoryBestDay, KeyHistoryBestLabel, KeyHistoryPrimeTime, KeyHistoryConvRate,
                KeyHistoryMsgsPerDay, KeyHistoryRespPerDay,
                KeyHistoryViewDetail, KeyHistoryPrevDay, KeyHistoryNextDay,
                KeyHistoryWeekTotal, KeyHistoryAvgRespTime, KeyHistoryWeekend,
                KeyHistorySelectDay,
                // Follow-up
                KeyFollowUpTitle, KeyFollowUpAutoRunning, KeyFollowUpSending,
                KeyFollowUpQueueToday, KeyFollowUpTotal, KeyFollowUpCold,
                KeyFollowUpAutoAll, KeyFollowUpAutoNote,
                KeyFollowUpEmpty, KeyFollowUpEmptyDesc,
                KeyFollowUpColdList, KeyFollowUpColdDesc, KeyFollowUpColdWarning,
                KeyFollowUpRecontact, KeyFollowUpRecontactDesc,
                KeyFollowUpRateInfo, KeyFollowUpVariantNote,
                KeyFollowUpSkipWait, KeyFollowUpSendingNow, KeyFollowUpWaitNext,
                KeyFollowUpAutoSchedule, KeyFollowUpNoDuplicate, KeyFollowUpDiffVariant,
                KeyFollowUpRotation, KeyFollowUpMessagesLabel,
                KeyFollowUpIceBreakerUnanswered, KeyFollowUpColdDetail,
                KeyFollowUpRecontactTemplate, KeyFollowUpRecontactTone,
                KeyFollowUpRecontactVibes, KeyFollowUpCanRecontact,
                KeyFollowUpAfterThat, KeyFollowUpAfterThatSilent,
                KeyFollowUpEditFirst, KeyFollowUpTabNiche,
                KeyFollowUpSendFinal,
                KeyFollowUpFU1Detail, KeyFollowUpFU2Detail,
                KeyFollowUpColdDashboard, KeyFollowUpTotalCold,
                KeyFollowUpVariantPreview1, KeyFollowUpVariantPreview2, KeyFollowUpVariantPreview3,
                KeyFollowUpVariantManualOnly, KeyFollowUpRatePerHour,
                KeyFollowUpVariantLabel,
                KeyFollowUpEmptyUnanswered,
                KeyFollowUpFU1Label, KeyFollowUpFU2Label,
                KeyFollowUpLeadNoun, KeyFollowUpSlot, KeyFollowUpOthers,
                KeyFollowUpViewDetail, KeyFollowUpSendRecontact,
                // Data screens — Leads DB
                KeyDataLeadsTitle, KeyDataFilterAll, KeyDataFilterSearch,
                KeyDataRecent, KeyDataMove, KeyDataViewDetail, KeyDataSendOfferAll,
                KeyDataScore, KeyDataTimeline, KeyDataStatus, KeyDataNotContacted,
                KeyDataFollowUpNext, KeyDataAutoSchedule, KeyDataColdDesc,
                KeyDataColdOption, KeyDataColdAuto, KeyDataNoResponseYet,
                KeyDataResponse, KeyDataOfferNotSent, KeyDataNotReplied,
                KeyDataSendOffer, KeyDataCustomReply, KeyDataLater,
                KeyDataMarkConvert, KeyDataArchive, KeyDataBlock,
                KeyDataSendFollowUp, KeyDataMarkCold, KeyDataSendIceBreaker,
                KeyDataLastFollowUp, KeyDataConversionRevenue, KeyDataDuration,
                KeyDataTemplate, KeyDataWorker, KeyDataIceBreakerSent,
                KeyDataResponseReceived, KeyDataOfferSent, KeyDataMarkConvertAction,
                KeyDataWaitingOffer, KeyDataContactCount, KeyDataFollowUpCount,
                KeyDataTotal, KeyDataFiltered, KeyDataView,
                // Data screens — Template Manager
                KeyTemplateTitle, KeyTemplateNiche, KeyTemplateIceBreaker,
                KeyTemplateOffer, KeyTemplateRecommended, KeyTemplateSelect,
                KeyTemplatePreviewAction, KeyTemplateNew, KeyTemplateEdit, KeyTemplateEditTitle,
                KeyTemplateUseThis, KeyTemplateEditFile, KeyTemplateSavedAsFile,
                KeyTemplateOpenInEditor, KeyTemplateAfterSave, KeyTemplateErrorTitle,
                KeyTemplateFileEmpty, KeyTemplateMinPlaceholder,
                KeyTemplateUnknownPlaceholder, KeyTemplateAvailablePlaceholders,
                KeyTemplateEncodingError, KeyTemplateEncodingHint,
                KeyTemplateWorkerPaused, KeyTemplateIceBreakerReq,
                KeyTemplateOpenFile, KeyTemplateViewPlaceholders,
                KeyTemplatePlaceholderLabel, KeyTemplatePreviewLabel,
                KeyTemplateReload,
                // General
                KeyGeneralLoading, KeyGeneralRetry, KeyGeneralSkip,
                KeyGeneralArchive, KeyGeneralDelete, KeyGeneralSave,
                KeyGeneralClose, KeyGeneralOff, KeyGeneralOn,
                // Action
                KeyActionArchiveAll,
                // Niche explorer — extended
                KeyNicheExplorerGenBarLabel,
                // Scrape — batch
                KeyScrapeBatchDone,
                KeyScrapeScraperError, KeyScrapeIdleNewLeads, KeyScrapeNicheLabel,
                KeyScrapeEmptyAreaHint, KeyScrapeEmptyFilterHint, KeyScrapeEmptyQueryHint,
        }
}
