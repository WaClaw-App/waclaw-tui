
## 13. State Machine Summary

```
LEAD STATES:
  baru → ice_breaker_sent → responded → offer_sent → converted
    │         │                 │            │
    │         │                 └→ negative → archived
    │         │                 └→ auto_reply → skipped
    │         └→ no_response → follow_up_1 → follow_up_2 → cold
    │         │                                    │            │
    │         │                                    └→ responded └→ offer_sent
    │         └→ failed → retry (max 3x) → dead
    └→ blocked (manual)

  RE-CONTACT: responded + dingin → 7 hari → re_contact → responded
  FOLLOW-UP LIMIT: ice_breaker + follow_up_1 + follow_up_2 = max 3 per lead
  COLD: 2x follow-up tanpa response → auto-tandai dingin

WORKER STATES (per niche):
  spawning → scraping → qualifying → queuing → sending → idle (loop)
    │            │           │           │          │
    │            └→ error → retry       │          └→ rate_limited
    │                                    └→ paused (manual)
    └→ config_error → paused → fixed → spawning
    └→ stopped (manual)

BATCH PIPELINE (per worker):
  queries.md → scrape → qualify → auto-review → queue → batch send → wait → loop
                                                                  └→ response → offer → deal

CONFIG VALIDATION STATES:
  unchecked → validating → clean | errors | warnings
    │                         │        │         │
    │                         │        └→ fix → revalidate → clean
    │                         └→ auto-dismiss (from boot)
    └→ first_time → generate template → revalidate → clean

SCREEN STATES:
  boot: boot_first_time | boot_returning | boot_returning_response | boot_returning_error | boot_returning_config_error | boot_returning_license_expired | boot_returning_device_conflict
  login: login_qr_waiting | login_qr_scanned | login_success | login_expired | login_failed
  niche: niche_list | niche_multi_selected | niche_custom | niche_edit_filters | niche_config_error
  scrape: scrape_active | scrape_multi_active | scrape_multi_staggered | scrape_idle | scrape_empty | scrape_error | scrape_gmaps_limited | scrape_auto_approved | scrape_high_value_reveal | scrape_batch_complete | scrape_wa_validation | scrape_wa_validation_progress
  review: review_reviewing | review_lead_detail | review_template_preview | review_queue_complete
  send: send_active | send_paused | send_off_hours | send_rate_limited | send_daily_limit | send_failed | send_all_slots_down | send_with_response
  monitor: monitor_live_dashboard | monitor_idle_background | monitor_night | monitor_error | monitor_empty | monitor_pending_responses
  response: response_positive | response_curious | response_negative | response_maybe | response_auto_reply | response_offer_preview | response_multi_queue | response_conversion | response_hot_lead | response_stop_detected | response_deal_detected
  leads: leads_list | leads_filtered | leads_full_detail | leads_follow_up_due | leads_cold | leads_never_contacted | leads_converted
  template: template_list | template_preview | template_edit_hint | template_validation_error
  workers: workers_overview | worker_detail | worker_add_niche | workers_paused
  shield: shield_overview | shield_warning | shield_danger | shield_slot_detail | shield_settings
  settings: settings_overview | settings_edit | settings_reload | settings_reload_error
  validation: validation_clean | validation_errors | validation_warnings | validation_fix | validation_first_time | validation_reload_error
  compose: compose_draft | compose_preview | compose_template_pick
  history: history_today | history_week | history_day_detail
  followup: followup_dashboard | followup_niche_detail | followup_sending | followup_empty | followup_cold_list | followup_recontact
  explorer: explorer_browse | explorer_search | explorer_category_detail | explorer_generating | explorer_generated
  update: update_available | update_downloading | update_ready | upgrade_available | upgrade_license_input | license_expired_with_upgrade | startup_check
  license: license_input | license_validating | license_valid | license_invalid | license_expired | license_device_conflict | license_server_error
  nerd_stats: nerd_stats_hidden | nerd_stats_minimal | nerd_stats_expanded (global overlay, not a screen)
  cmd_palette: cmd_palette_closed | cmd_palette_open | cmd_palette_executing | cmd_palette_empty | cmd_palette_with_recent | cmd_palette_quick_action (global overlay, not a screen)

NOTIFICATION TYPES (code identifiers):
  notif_response_received | notif_multi_response | notif_scrape_complete | notif_batch_send_complete | notif_wa_disconnect |
  notif_wa_flag | notif_health_score_drop | notif_daily_limit | notif_streak_milestone | notif_config_error | notif_validation_error |
  notif_license_expired | notif_device_conflict | notif_follow_up_scheduled | notif_lead_cold | notif_update_available | notif_upgrade_available

CONFIRMATION OVERLAYS (code identifiers):
  confirm_bulk_offer | confirm_bulk_delete | confirm_bulk_archive | confirm_force_disconnect
```

---

*WaClaw — army lu cerdas, keras, dan aman.*
*Broken config? Ketahuan. Multiple responses? Ditangani. History? Ada.*
*Follow-up? Persistent. Lisensi? Satu kunci satu device. Niche bingung? Explorer.*
*Update kecil? Gratis. Upgrade besar? Pilihan lu. Nerd stats? Toggle kapan aja.*
*Lupa shortcut? Ctrl+K. Semua command, dari mana aja, 3 detik.*
*Lu cuma nonton. WaClaw yang kerja. Tapi kalau lu mau intervene, 1 tombol cukup.*
