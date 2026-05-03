package overlay

import (
	"fmt"
	"strings"

	"github.com/WaClaw-App/waclaw/internal/tui/anim"
	"github.com/WaClaw-App/waclaw/internal/tui/i18n"
	"github.com/WaClaw-App/waclaw/internal/tui/style"
	"github.com/WaClaw-App/waclaw/pkg/protocol"
	"github.com/charmbracelet/lipgloss"
)

// ConfirmationOverlay holds the state for the confirmation dialog overlay.
//
// Spec (doc/14-notification-system.md):
//   - Appears before destructive/high-impact actions
//   - Always "↵ gas" / "↵ go" + "s batal" / "s cancel" — never more than 2 options
//   - 4 types: bulk_offer, bulk_delete, bulk_archive, force_disconnect
//   - Always requires explicit confirmation
type ConfirmationOverlay struct {
	// Active is the currently displayed confirmation, or nil.
	Active *ConfirmationData

	// Width is the available terminal width.
	Width int

	// Anim tracks the fade-in animation.
	Anim anim.AnimationState

	// OnConfirm is called when the user confirms the action.
	OnConfirm func(confirmationType protocol.ConfirmationType, data map[string]any)

	// OnCancel is called when the user cancels the action.
	OnCancel func()
}

// ConfirmationData carries all information for a confirmation dialog.
type ConfirmationData struct {
	// Type is the confirmation category.
	Type protocol.ConfirmationType

	// Title is the primary question text.
	Title string

	// Detail is the supporting information (e.g. breakdown of affected items).
	Detail string

	// Data carries arbitrary context for the confirm/cancel callbacks.
	Data map[string]any
}

// NewConfirmationOverlay creates a ConfirmationOverlay with default settings.
func NewConfirmationOverlay() ConfirmationOverlay {
	return ConfirmationOverlay{}
}

// Show displays a confirmation dialog.
func (co *ConfirmationOverlay) Show(data ConfirmationData) {
	co.Active = &data
	co.Anim = anim.NewAnimationState(anim.AnimFade, anim.NotifSlideIn)
}

// Dismiss removes the confirmation dialog without confirming.
func (co *ConfirmationOverlay) Dismiss() {
	co.Active = nil
	if co.OnCancel != nil {
		co.OnCancel()
	}
}

// Confirm confirms the action and dismisses the dialog.
func (co *ConfirmationOverlay) Confirm() {
	if co.Active == nil {
		return
	}
	confirmType := co.Active.Type
	data := co.Active.Data
	co.Active = nil
	if co.OnConfirm != nil {
		co.OnConfirm(confirmType, data)
	}
}

// IsVisible returns true if a confirmation is currently displayed.
func (co ConfirmationOverlay) IsVisible() bool {
	return co.Active != nil
}

// HandleKey processes a key event while the confirmation is visible.
// Returns true if the key was consumed.
func (co *ConfirmationOverlay) HandleKey(keyStr string) bool {
	if co.Active == nil {
		return false
	}

	switch keyStr {
	case "enter":
		co.Confirm()
		return true
	case "s":
		co.Dismiss()
		return true
	case "esc":
		co.Dismiss()
		return true
	}
	return true // consume all keys while confirmation is visible
}

// Tick advances the animation state.
func (co *ConfirmationOverlay) Tick() {
	if co.Active == nil {
		return
	}
	co.Anim.UpdateProgress()
}

// View renders the confirmation overlay.
func (co ConfirmationOverlay) View() string {
	if co.Active == nil {
		return ""
	}

	c := co.Active
	width := co.Width
	if width < 40 {
		width = 40
	}

	var lines []string

	// Warning icon + title.
	titleStyle := lipgloss.NewStyle().Foreground(style.Warning).Bold(true)
	lines = append(lines, titleStyle.Render(fmt.Sprintf("⚠️  %s", c.Title)))

	// Detail.
	if c.Detail != "" {
		detailStyle := lipgloss.NewStyle().Foreground(style.TextMuted)
		for _, line := range strings.Split(c.Detail, "\n") {
			lines = append(lines, detailStyle.Render(line))
		}
	}

	// Actions — always exactly 2 options.
	lines = append(lines, "")
	proceed := i18n.T("confirm.proceed")
	cancel := i18n.T("confirm.cancel")
	actionStyle := lipgloss.NewStyle().Foreground(style.TextDim)
	lines = append(lines, actionStyle.Render(fmt.Sprintf("%s    %s", proceed, cancel)))

	// Wrap in panel.
	content := strings.Join(lines, "\n")
	panel := lipgloss.NewStyle().
		Background(style.BgRaised).
		Width(width).
		Padding(0, 2).
		Render(content)

	return panel
}

// ConfirmationDataFromType creates a ConfirmationData from a confirmation type
// with default text based on the type.
func ConfirmationDataFromType(confirmType protocol.ConfirmationType, data map[string]any) ConfirmationData {
	title := confirmationTitle(confirmType, data)
	detail := confirmationDetail(confirmType, data)

	return ConfirmationData{
		Type:   confirmType,
		Title:  title,
		Detail: detail,
		Data:   data,
	}
}

// confirmationTitle returns the default title for each confirmation type.
func confirmationTitle(confirmType protocol.ConfirmationType, data map[string]any) string {
	if msg, ok := data["message"].(string); ok {
		return msg
	}

	switch confirmType {
	case protocol.ConfirmBulkOffer:
		return i18n.T("confirm.bulk_offer_title")
	case protocol.ConfirmBulkDelete:
		return i18n.T("confirm.bulk_delete_title")
	case protocol.ConfirmBulkArchive:
		return i18n.T("confirm.bulk_archive_title")
	case protocol.ConfirmForceDisconnect:
		return i18n.T("confirm.force_disconnect_title")
	default:
		return string(confirmType)
	}
}

// confirmationDetail returns the supporting detail for each confirmation type.
func confirmationDetail(confirmType protocol.ConfirmationType, data map[string]any) string {
	if detail, ok := data["detail"].(string); ok {
		return detail
	}

	switch confirmType {
	case protocol.ConfirmBulkOffer:
		if count, ok := data["count"].(float64); ok {
			return fmt.Sprintf(i18n.T("confirm.bulk_offer_detail"), int(count))
		}
	case protocol.ConfirmBulkDelete:
		if count, ok := data["count"].(float64); ok {
			return fmt.Sprintf(i18n.T("confirm.bulk_delete_detail"), int(count))
		}
	case protocol.ConfirmBulkArchive:
		if count, ok := data["count"].(float64); ok {
			return fmt.Sprintf(i18n.T("confirm.bulk_archive_detail"), int(count))
		}
	case protocol.ConfirmForceDisconnect:
		if device, ok := data["device"].(string); ok {
			return fmt.Sprintf(i18n.T("confirm.force_disconnect_detail"), device)
		}
	}
	return ""
}
