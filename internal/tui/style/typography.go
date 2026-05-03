package style

import "github.com/charmbracelet/lipgloss"

// Typography styles — weight hierarchy: bold primary, muted secondary, dim tertiary.
// Badge styles are consolidated here as the single source of truth for
// status badge rendering — no badge definitions in colors.go.
var (
        // HeadingStyle is the primary heading — bold, bright, maximum weight.
        HeadingStyle = lipgloss.NewStyle().
                        Bold(true).
                        Foreground(Text)

        // SubHeadingStyle is a secondary heading — bold, muted.
        SubHeadingStyle = lipgloss.NewStyle().
                        Bold(true).
                        Foreground(TextMuted)

        // BodyStyle is the default body text.
        BodyStyle = lipgloss.NewStyle().
                        Foreground(Text)

        // CaptionStyle is tertiary/dim text — labels, metadata.
        CaptionStyle = lipgloss.NewStyle().
                        Foreground(TextDim)

        // ActionStyle is for interactive elements (accent color = clickable).
        ActionStyle = lipgloss.NewStyle().
                        Foreground(Accent).
                        Bold(true)

        // BadgeNewStyle for "new" status badges — accent (#818CF8).
        BadgeNewStyle = lipgloss.NewStyle().Foreground(Accent)

        // BadgeRespondedStyle for "responded" status badges — warning/amber (#FBBF24).
        BadgeRespondedStyle = lipgloss.NewStyle().Foreground(Warning)

        // BadgeConvertedStyle for "converted/deal" status badges — success (#34D399).
        BadgeConvertedStyle = lipgloss.NewStyle().Foreground(Success)

        // BadgeFailedStyle for "failed" status badges — text_dim (#3D3D44).
        BadgeFailedStyle = lipgloss.NewStyle().Foreground(TextDim)

        // BadgeColdStyle for "cold" status badges — dim blue (#5B7A9D).
        // Doc says: "❄ COLD badge in dim blue = bukan mati, tapi dingin".
        BadgeColdStyle = lipgloss.NewStyle().Foreground(DimBlue)

        // TextStyle is an alias for BodyStyle — regular weight primary text.
        TextStyle = BodyStyle

        // SelectedBodyStyle is body text with bold for selected/focused items.
        SelectedBodyStyle = lipgloss.NewStyle().Foreground(Text).Bold(true)

        // WarningBoldStyle is warning amber with bold for selected recontact leads.
        WarningBoldStyle = lipgloss.NewStyle().Foreground(Warning).Bold(true)

        // BarLabelStyle is a fixed-width label for bar chart day names.
        BarLabelStyle = lipgloss.NewStyle().Foreground(TextMuted).Width(8)
)

// rebuildTypography recreates all typography styles from the current color tokens.
// Called by RebuildStyles during theme hot-reload.
func rebuildTypography() {
        HeadingStyle = lipgloss.NewStyle().
                Bold(true).
                Foreground(Text)

        SubHeadingStyle = lipgloss.NewStyle().
                Bold(true).
                Foreground(TextMuted)

        BodyStyle = lipgloss.NewStyle().
                Foreground(Text)

        CaptionStyle = lipgloss.NewStyle().
                Foreground(TextDim)

        ActionStyle = lipgloss.NewStyle().
                Foreground(Accent).
                Bold(true)

        BadgeNewStyle = lipgloss.NewStyle().Foreground(Accent)
        BadgeRespondedStyle = lipgloss.NewStyle().Foreground(Warning)
        BadgeConvertedStyle = lipgloss.NewStyle().Foreground(Success)
        BadgeFailedStyle = lipgloss.NewStyle().Foreground(TextDim)
        BadgeColdStyle = lipgloss.NewStyle().Foreground(DimBlue)
        TextStyle = BodyStyle

        SelectedBodyStyle = lipgloss.NewStyle().Foreground(Text).Bold(true)
        WarningBoldStyle = lipgloss.NewStyle().Foreground(Warning).Bold(true)
        BarLabelStyle = lipgloss.NewStyle().Foreground(TextMuted).Width(8)
}
