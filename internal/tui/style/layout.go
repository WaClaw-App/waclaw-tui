package style

import (
        "strings"
        "unicode/utf8"
)

// Layout constants — vertical borderless. No borders, no boxes, only space.
const (
        // SpacingUnit is 1 line = 8px equivalent.
        SpacingUnit = 1

        // SectionGap is the gap between major sections: 2 lines (16px).
        SectionGap = 2

        // SubSectionGap is the gap between sub-sections: 1 line (8px).
        SubSectionGap = 1

        // ItemGap is the gap between items in the same group: 0 lines (touching = grouped).
        ItemGap = 0

        // IndentPerLevel is 2 spaces per nesting level.
        IndentPerLevel = 2
)

// Section returns a string of newlines for section spacing.
func Section(gap int) string {
        if gap <= 0 {
                return ""
        }
        return strings.Repeat("\n", gap)
}

// Indent returns indentation for the given nesting level.
func Indent(level int) string {
        if level <= 0 {
                return ""
        }
        return strings.Repeat(" ", level*IndentPerLevel)
}

// SeparatorWidth is the default width for horizontal separator lines.
const SeparatorWidth = 50

// Input field layout constants — DRY, no magic numbers.
const (
        // LicenseInputWidth is the visible content width of the license key input field.
        LicenseInputWidth = 40

        // LicenseInputFrameWidth includes the 2-char padding (1 on each side).
        LicenseInputFrameWidth = LicenseInputWidth + 2

        // LicenseGlowWidth is the width of the red glow accent line under the input.
        // Matches the input frame width so the glow aligns visually.
        LicenseGlowWidth = LicenseInputFrameWidth

        // LicenseProgressBarWidth is the default progress bar width for validation.
        LicenseProgressBarWidth = 30

        // UpgradeInputWidth is the width of the v2 license input field (longer prefix).
        UpgradeInputWidth = 50

        // UpgradeProgressBarWidth is the progress bar width for the update screen.
        UpgradeProgressBarWidth = 40

        // DownloadBarMaxWidth is the maximum width for the download progress bar.
        // Caps the bar so it doesn't stretch across very wide terminals.
        DownloadBarMaxWidth = 60

        // DownloadBarPadding is the horizontal padding subtracted from the terminal
        // width when sizing the download progress bar (accounts for indent + margins).
        DownloadBarPadding = 4

        // ActionHintSpacing is the number of spaces between action hints in footers.
        ActionHintSpacing = 4

        // DefaultProgressBarWidth is the default width for progress bars across screens.
        // Used by antiban slot detail bars, send rate bars, and other generic bars.
        DefaultProgressBarWidth = 40

        // CompactProgressBarWidth is the narrower progress bar for inline contexts
        // (e.g., antiban slot cards, where the bar sits next to a label).
        CompactProgressBarWidth = 30
)

// Separator returns a horizontal visual break using dim dots instead of
// box-drawing characters (─━). The vertical borderless design system (P3)
// forbids ─ ━ characters; hierarchy is conveyed through brightness, not
// borders. The dot pattern preserves visual rhythm while complying with
// the "no horizontal separators" rule.
func Separator() string {
        return DimStyle.Render(strings.Repeat("·", SeparatorWidth))
}

// SectionLabel renders a section heading using the ·· label ······ pattern.
// Uses dots instead of ── to comply with the vertical borderless design system (P3).
// The pattern preserves the wireframe intent of "── label ────────"
// without violating the "no ─ ━" rule. The label is rendered in SubHeadingStyle.
func SectionLabel(label string) string {
        const dotPad = 2 // "·· " before label
        labelLen := utf8.RuneCountInString(label)
        afterLen := SeparatorWidth - labelLen - dotPad
        if afterLen < 3 {
                afterLen = 3
        }
        return SubHeadingStyle.Render("·· "+label+" "+strings.Repeat("·", afterLen)) + "\n"
}
