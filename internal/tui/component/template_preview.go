package component

import (
        "fmt"
        "strings"
        "time"

        "github.com/WaClaw-App/waclaw/internal/tui/anim"
        "github.com/WaClaw-App/waclaw/internal/tui/i18n"
        "github.com/WaClaw-App/waclaw/internal/tui/style"
        "github.com/charmbracelet/lipgloss"
)

// templatePlaceholders is the single source of truth for supported template
// placeholder names. DRY — used by Substituted, ViewWithHighlight, Validate,
// and exposed via TemplatePlaceholderNames for consumers (e.g. template_mgr.go).
var templatePlaceholders = []string{
        "Title", "Category", "Address", "City", "Rating", "Reviews", "Area",
}

// TemplatePlaceholderNames returns the supported template placeholder names.
// Consumers should use this instead of maintaining their own duplicate list.
func TemplatePlaceholderNames() []string {
        return templatePlaceholders
}

// TemplatePreview renders a template text with placeholder substitution
// and optional char-by-char type-out animation.
//
// Visual spec (from doc/04-screens-lead-review-send.md, doc/06-screens-database-templates.md):
//   - Placeholder substitution: {{.Title}}, {{.Category}}, {{.Address}}, {{.City}}, {{.Rating}}, {{.Reviews}}, {{.Area}}
//   - Char-by-char type-out preview (for lead review)
//   - Validation: ✗/✓ for missing placeholders
//   - Used in: lead review, template manager, compose
type TemplatePreview struct {
        // Template is the raw template text with placeholders.
        Template string

        // Vars holds the substitution values keyed by placeholder name.
        // Supported keys: Title, Category, Address, City, Rating, Reviews, Area.
        Vars map[string]string

        // TypeOutEnabled enables the char-by-char typing animation.
        TypeOutEnabled bool

        // TypeOutProgress is the fraction of characters revealed [0.0, 1.0].
        TypeOutProgress float64

        // TypeOutStart is when the type-out animation began.
        TypeOutStart time.Time

        // Width limits the preview width.
        Width int
}

// NewTemplatePreview creates a TemplatePreview with the given template.
func NewTemplatePreview(template string) TemplatePreview {
        return TemplatePreview{
                Template: template,
                Vars:     make(map[string]string),
                Width:    60,
        }
}

// SetVar sets a substitution variable.
func (tp *TemplatePreview) SetVar(key, value string) {
        tp.Vars[key] = value
}

// Substituted returns the template text with all placeholders replaced.
func (tp TemplatePreview) Substituted() string {
        result := tp.Template
        for _, ph := range templatePlaceholders {
                placeholder := fmt.Sprintf("{{.%s}}", ph)
                value, ok := tp.Vars[ph]
                if !ok || value == "" {
                        value = placeholder
                }
                result = strings.ReplaceAll(result, placeholder, value)
        }
        return result
}

// StartTypeOut begins the char-by-char typing animation.
func (tp *TemplatePreview) StartTypeOut() {
        tp.TypeOutEnabled = true
        tp.TypeOutStart = time.Now()
        tp.TypeOutProgress = 0
}

// Tick advances the type-out animation.
func (tp *TemplatePreview) Tick(now time.Time) {
        if !tp.TypeOutEnabled {
                return
        }

        text := tp.Substituted()
        totalChars := len(text)
        if totalChars == 0 {
                tp.TypeOutProgress = 1.0
                return
        }

        // Type out at anim.TypeOutCharDelay per character.
        elapsed := now.Sub(tp.TypeOutStart)
        typedChars := int(elapsed / anim.TypeOutCharDelay)
        tp.TypeOutProgress = float64(typedChars) / float64(totalChars)

        if tp.TypeOutProgress >= 1.0 {
                tp.TypeOutProgress = 1.0
                tp.TypeOutEnabled = false
        }
}

// View renders the template preview.
func (tp TemplatePreview) View() string {
        text := tp.Substituted()

        // Apply type-out animation if active.
        if tp.TypeOutEnabled && tp.TypeOutProgress < 1.0 {
                totalChars := len(text)
                visibleChars := int(float64(totalChars) * tp.TypeOutProgress)
                if visibleChars > totalChars {
                        visibleChars = totalChars
                }
                text = text[:visibleChars]
        }

        rendered := lipgloss.NewStyle().Foreground(style.Text).Render(text)

        if tp.TypeOutEnabled && tp.TypeOutProgress < 1.0 {
                rendered += lipgloss.NewStyle().Foreground(style.Accent).Render("▎")
        }

        return rendered
}

// ViewWithHighlight renders the template with unsubstituted placeholders
// highlighted in accent and substituted values in normal text.
func (tp TemplatePreview) ViewWithHighlight() string {
        result := tp.Template
        for _, ph := range templatePlaceholders {
                placeholder := fmt.Sprintf("{{.%s}}", ph)
                value, ok := tp.Vars[ph]
                if ok && value != "" {
                        styled := lipgloss.NewStyle().Foreground(style.Text).Render(value)
                        result = strings.ReplaceAll(result, placeholder, styled)
                } else {
                        styled := lipgloss.NewStyle().Foreground(style.Accent).Render(placeholder)
                        result = strings.ReplaceAll(result, placeholder, styled)
                }
        }
        return result
}

// ValidationResult holds the result of template placeholder validation.
type ValidationResult struct {
        // Missing lists placeholder names that are present in the template
        // but have no corresponding value in Vars.
        Missing []string

        // Extra lists placeholder names in Vars that are not used in the template.
        Extra []string

        // Valid is true if all placeholders have values.
        Valid bool
}

// Validate checks that all placeholders in the template have substitution values.
func (tp TemplatePreview) Validate() ValidationResult {
        var missing []string
        var used []string

        for _, ph := range templatePlaceholders {
                placeholder := fmt.Sprintf("{{.%s}}", ph)
                if strings.Contains(tp.Template, placeholder) {
                        used = append(used, ph)
                        if _, ok := tp.Vars[ph]; !ok {
                                missing = append(missing, ph)
                        }
                }
        }

        // Check for extra vars not used in template.
        var extra []string
        for k := range tp.Vars {
                found := false
                for _, u := range used {
                        if k == u {
                                found = true
                                break
                        }
                }
                if !found {
                        extra = append(extra, k)
                }
        }

        return ValidationResult{
                Missing: missing,
                Extra:   extra,
                Valid:   len(missing) == 0,
        }
}

// ViewValidation renders a validation status line.
// Shows ✓ if valid, ✗ with missing placeholders if invalid.
func (tp TemplatePreview) ViewValidation() string {
        result := tp.Validate()
        if result.Valid {
                return lipgloss.NewStyle().Foreground(style.Success).Render(
                        i18n.T(i18n.KeyGuardrailClean),
                )
        }

        var lines []string
        lines = append(lines, lipgloss.NewStyle().Foreground(style.Danger).Render(
                fmt.Sprintf("✗ %d %s", len(result.Missing), i18n.T(i18n.KeyGuardrailErrors)),
        ))
        for _, m := range result.Missing {
                placeholder := fmt.Sprintf("{{.%s}}", m)
                lines = append(lines, lipgloss.NewStyle().Foreground(style.Danger).Render(
                        fmt.Sprintf("  %s", placeholder),
                ))
        }

        return strings.Join(lines, "\n")
}

// String implements fmt.Stringer.
func (tp TemplatePreview) String() string {
        return fmt.Sprintf("TemplatePreview{template=%d chars, vars=%d, typeOut=%.0f%%}",
                len(tp.Template), len(tp.Vars), tp.TypeOutProgress*100)
}
