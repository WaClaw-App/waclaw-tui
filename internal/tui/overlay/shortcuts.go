package overlay

import (
        "strings"

        "github.com/WaClaw-App/waclaw/internal/tui/anim"
        "github.com/WaClaw-App/waclaw/internal/tui/i18n"
        "github.com/WaClaw-App/waclaw/internal/tui/style"
        "github.com/charmbracelet/lipgloss"
)

// ShortcutsOverlay holds the state for the keyboard shortcuts overlay.
//
// Spec (doc/19-keyboard-grammar.md):
//   - "?" key shows shortcuts on current screen
//   - Fade in on current screen, fade out on any keypress
//   - Never lose context — it's an overlay, not a screen
//   - Compact: all shortcuts visible at a glance
type ShortcutsOverlay struct {
        // Visible indicates whether the shortcuts are currently displayed.
        Visible bool

        // Width is the available terminal width.
        Width int

        // Anim tracks the fade-in/out animation.
        Anim anim.AnimationState
}

// NewShortcutsOverlay creates a ShortcutsOverlay in hidden state.
func NewShortcutsOverlay() ShortcutsOverlay {
        return ShortcutsOverlay{}
}

// Show displays the shortcuts overlay with a fade-in animation.
func (so *ShortcutsOverlay) Show() {
        so.Visible = true
        so.Anim = anim.NewAnimationState(anim.AnimFade, anim.ShortcutsFadeIn)
}

// Hide hides the shortcuts overlay with a fade-out animation.
func (so *ShortcutsOverlay) Hide() {
        so.Visible = false
        so.Anim = anim.NewAnimationState(anim.AnimFade, anim.ShortcutsFadeOut)
}

// Toggle shows or hides the shortcuts overlay.
func (so *ShortcutsOverlay) Toggle() {
        if so.Visible {
                so.Hide()
        } else {
                so.Show()
        }
}

// IsVisible returns true if the shortcuts are displayed.
func (so ShortcutsOverlay) IsVisible() bool {
        return so.Visible
}

// Tick advances the animation state.
func (so *ShortcutsOverlay) Tick() {
        so.Anim.UpdateProgress()
}

// View renders the shortcuts overlay.
//
// Spec (doc/19-keyboard-grammar.md):
//
//        ── tombol ─────────────
//
//        ↑↓   pindah
//        ↵    aksi utama
//        1-3  pilih opsi
//        s    skip
//        q    balik/keluar
//        /    cari
//        v    validasi config
//        l    lisensi
//        h    riwayat
//        r    reload
//        `    nerd stats
//        u    cek update
//        ^K   command palette
//
//        tekan apa aja buat tutup
func (so ShortcutsOverlay) View() string {
        if !so.Visible {
                return ""
        }

        width := so.Width
        if width < 36 {
                width = 36
        }

        title := i18n.T("shortcuts.title")
        headerStyle := lipgloss.NewStyle().Foreground(style.TextDim)
        keyStyle := lipgloss.NewStyle().Foreground(style.Accent).Bold(true)
        descStyle := lipgloss.NewStyle().Foreground(style.TextMuted)
        footerStyle := lipgloss.NewStyle().Foreground(style.TextDim)

        var lines []string

        // Header.
        lines = append(lines, headerStyle.Render(
                "── "+title+" "+"──────────────────────"[:max(0, width-4-len(title))]))

        lines = append(lines, "")

        // Shortcut entries. Each entry: key (left) + description (right).
        entries := []struct {
                key  string
                desc string
        }{
                {"↑↓", i18n.T("shortcuts.move")},
                {"↵", i18n.T("shortcuts.primary_action")},
                {"1-3", i18n.T("shortcuts.pick_option")},
                {"s", i18n.T("shortcuts.skip")},
                {"q", i18n.T("shortcuts.back_quit")},
                {"p", i18n.T("shortcuts.pause")},
                {"/", i18n.T("shortcuts.search")},
                {"v", i18n.T("shortcuts.validate_config")},
                {"l", i18n.T("shortcuts.license")},
                {"h", i18n.T("shortcuts.history")},
                {"r", i18n.T("shortcuts.reload")},
                {"`", i18n.T("shortcuts.nerd_stats")},
                {"u", i18n.T("shortcuts.check_update")},
                {"^K", i18n.T("shortcuts.cmd_palette")},
        }

        for _, e := range entries {
                keyPart := keyStyle.Render(e.key)
                // Pad the key column to 5 visible characters.
                padLen := 5 - len(e.key)
                if padLen < 1 {
                        padLen = 1
                }
                pad := strings.Repeat(" ", padLen)
                descPart := descStyle.Render(e.desc)
                lines = append(lines, keyPart+pad+descPart)
        }

        lines = append(lines, "")
        lines = append(lines, footerStyle.Render(i18n.T("shortcuts.press_any_key")))

        content := strings.Join(lines, "\n")

        // Render as a floating panel with raised background.
        panel := lipgloss.NewStyle().
                Background(style.BgRaised).
                Width(min(width, 40)).
                Padding(0, 2).
                Render(content)

        return panel
}
