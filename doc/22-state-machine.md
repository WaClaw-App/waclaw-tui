
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
  boot: first_time | returning | returning+response | returning+error | returning+config_error | returning+license_expired | returning+device_conflict
  login: qr_waiting | qr_scanned | login_success | login_expired | login_failed
  niche: niche_list | niche_multi_selected | niche_custom | niche_edit_filters | niche_config_error
  scrape: scraping_active | scraping_multi_active | scraping_multi_staggered | scrape_idle | scrape_empty | scrape_error | scrape_gmaps_limited | scrape_auto_approved | scrape_high_value_reveal | scrape_batch_complete
  review: reviewing | lead_detail | template_preview | queue_complete
  send: sending_active | sending_paused | sending_off_hours | sending_rate_limited | sending_daily_limit | sending_failed | sending_all_slots_down | sending_with_response
  monitor: live_dashboard | idle_background | dashboard_night | dashboard_error | dashboard_empty | dashboard_with_pending_responses
  response: response_positive | response_curious | response_negative | response_maybe | response_auto_reply | offer_preview | response_multi_queue | conversion
  leads: leads_list | leads_filtered | lead_full_detail | lead_follow_up_due | lead_cold | lead_never_contacted | lead_converted
  template: template_list | template_preview | template_edit_hint | template_validation_error
  workers: workers_overview | worker_detail | worker_add_niche | worker_paused
  shield: shield_overview | shield_warning | shield_danger | shield_slot_detail | shield_settings
  settings: settings_overview | settings_edit | settings_reload | settings_reload_error
  validation: validation_clean | validation_errors | validation_warnings | validation_fix | validation_first_time
  compose: compose_draft | compose_preview | compose_template_pick
  history: history_today | history_week | history_day_detail
  followup: followup_dashboard | followup_niche_detail | followup_sending | followup_empty | followup_cold_list | followup_recontact
  explorer: explorer_browse | explorer_search | explorer_category_detail | explorer_generating | explorer_generated
  update: update_available | update_downloading | update_ready | upgrade_available | upgrade_license_input | license_expired_with_upgrade
  license: license_input | license_validating | license_valid | license_invalid | license_expired | license_device_conflict | license_server_error
  nerd_stats: hidden | minimal | expanded (global overlay, not a screen)
  cmd_palette: cmd_closed | cmd_open | cmd_executing | cmd_empty | cmd_with_recent | cmd_quick_action (global overlay, not a screen)

NOTIFICATION TYPES:
  response_masuk | multi_response | scrape_selesai | batch_selesai | wa_disconnect |
  wa_flag | health_drop | limit_harian | streak_milestone | config_error | validation_error |
  license_expired | device_conflict | followup_terjadwal | lead_dingin | update_available | upgrade_available

CONFIRMATION OVERLAYS:
  bulk_offer | bulk_delete | bulk_archive | force_device_disconnect
```

---

*WaClaw — army lu cerdas, keras, dan aman.*
*Broken config? Ketahuan. Multiple responses? Ditangani. History? Ada.*
*Follow-up? Persistent. Lisensi? Satu kunci satu device. Niche bingung? Explorer.*
*Update kecil? Gratis. Upgrade besar? Pilihan lu. Nerd stats? Toggle kapan aja.*
*Lupa shortcut? Ctrl+K. Semua command, dari mana aja, 3 detik.*
*Lu cuma nonton. WaClaw yang kerja. Tapi kalau lu mau intervene, 1 tombol cukup.*
