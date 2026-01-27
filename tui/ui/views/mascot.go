package views

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"
	"statusline-config/config"
)

// MascotCategory represents a mascot mood category
type MascotCategory struct {
	Key         string
	Label       string
	Description string
	Enabled     *bool
	Threshold   *int // nil for time-based (no threshold)
	Emojis      *[]string
}

// MascotView handles the mascot settings screen
type MascotView struct {
	Categories      []MascotCategory
	Selected        int
	SubSelected     int // For editing within a category
	InCategory      bool
	EditingEmoji    bool
	EditingThreshold bool
	EmojiInput      textinput.Model
	ThresholdInput  textinput.Model
	Config          *config.Config
}

// NewMascotView creates a new mascot view
func NewMascotView(cfg *config.Config) *MascotView {
	emojiInput := textinput.New()
	emojiInput.CharLimit = 30
	emojiInput.Width = 25

	thresholdInput := textinput.New()
	thresholdInput.CharLimit = 4
	thresholdInput.Width = 5

	view := &MascotView{
		Config: cfg,
		Categories: []MascotCategory{
			{
				Key:         "context_panic",
				Label:       "Context Panic Mode",
				Description: "When context usage exceeds threshold",
				Enabled:     &cfg.Mascot.ContextPanic.Enabled,
				Threshold:   &cfg.Mascot.ContextPanic.Threshold,
				Emojis:      &cfg.Mascot.ContextPanic.Emojis,
			},
			{
				Key:         "productive",
				Label:       "Productive Mode",
				Description: "When many lines have been added",
				Enabled:     &cfg.Mascot.Productive.Enabled,
				Threshold:   &cfg.Mascot.Productive.Threshold,
				Emojis:      &cfg.Mascot.Productive.Emojis,
			},
			{
				Key:         "deletion",
				Label:       "Deletion Mode",
				Description: "When more lines removed than added",
				Enabled:     &cfg.Mascot.Deletion.Enabled,
				Threshold:   &cfg.Mascot.Deletion.Threshold,
				Emojis:      &cfg.Mascot.Deletion.Emojis,
			},
			{
				Key:         "time_night",
				Label:       "Night Mood (12am-6am)",
				Description: "Late night/early morning moods",
				Enabled:     &cfg.Mascot.TimeBased.Enabled,
				Threshold:   nil,
				Emojis:      &cfg.Mascot.TimeBased.Night,
			},
			{
				Key:         "time_morning",
				Label:       "Morning Mood (6am-12pm)",
				Description: "Morning time moods",
				Enabled:     &cfg.Mascot.TimeBased.Enabled,
				Threshold:   nil,
				Emojis:      &cfg.Mascot.TimeBased.Morning,
			},
			{
				Key:         "time_afternoon",
				Label:       "Afternoon Mood (12pm-6pm)",
				Description: "Afternoon moods",
				Enabled:     &cfg.Mascot.TimeBased.Enabled,
				Threshold:   nil,
				Emojis:      &cfg.Mascot.TimeBased.Afternoon,
			},
			{
				Key:         "time_evening",
				Label:       "Evening Mood (6pm-12am)",
				Description: "Evening moods",
				Enabled:     &cfg.Mascot.TimeBased.Enabled,
				Threshold:   nil,
				Emojis:      &cfg.Mascot.TimeBased.Evening,
			},
		},
		EmojiInput:     emojiInput,
		ThresholdInput: thresholdInput,
	}

	return view
}

// Up moves selection up
func (v *MascotView) Up() {
	if v.EditingEmoji || v.EditingThreshold {
		return
	}
	if v.InCategory {
		v.SubSelected--
		cat := v.Categories[v.Selected]
		maxItems := len(*cat.Emojis) + 1 // emojis + toggle enabled + threshold (if applicable)
		if cat.Threshold != nil {
			maxItems++
		}
		if v.SubSelected < 0 {
			v.SubSelected = maxItems - 1
		}
	} else {
		v.Selected--
		if v.Selected < 0 {
			v.Selected = len(v.Categories) - 1
		}
	}
}

// Down moves selection down
func (v *MascotView) Down() {
	if v.EditingEmoji || v.EditingThreshold {
		return
	}
	if v.InCategory {
		v.SubSelected++
		cat := v.Categories[v.Selected]
		maxItems := len(*cat.Emojis) + 1 // emojis + toggle enabled
		if cat.Threshold != nil {
			maxItems++
		}
		if v.SubSelected >= maxItems {
			v.SubSelected = 0
		}
	} else {
		v.Selected++
		if v.Selected >= len(v.Categories) {
			v.Selected = 0
		}
	}
}

// Enter enters a category or edits an item
func (v *MascotView) Enter() {
	if v.EditingEmoji {
		v.StopEditEmoji()
		return
	}
	if v.EditingThreshold {
		v.StopEditThreshold()
		return
	}
	if !v.InCategory {
		v.InCategory = true
		v.SubSelected = 0
		return
	}

	cat := v.Categories[v.Selected]
	// SubSelected 0 = toggle enabled
	if v.SubSelected == 0 {
		*cat.Enabled = !*cat.Enabled
		return
	}

	// If has threshold, SubSelected 1 = threshold
	if cat.Threshold != nil {
		if v.SubSelected == 1 {
			v.StartEditThreshold()
			return
		}
		// Otherwise it's an emoji (SubSelected - 2)
		emojiIdx := v.SubSelected - 2
		if emojiIdx >= 0 && emojiIdx < len(*cat.Emojis) {
			v.StartEditEmoji(emojiIdx)
		}
	} else {
		// No threshold, SubSelected - 1 = emoji index
		emojiIdx := v.SubSelected - 1
		if emojiIdx >= 0 && emojiIdx < len(*cat.Emojis) {
			v.StartEditEmoji(emojiIdx)
		}
	}
}

// Back goes back from category view
func (v *MascotView) Back() bool {
	if v.EditingEmoji {
		v.CancelEditEmoji()
		return false
	}
	if v.EditingThreshold {
		v.CancelEditThreshold()
		return false
	}
	if v.InCategory {
		v.InCategory = false
		v.SubSelected = 0
		return false
	}
	return true // Signal to go back to main menu
}

// StartEditEmoji begins editing an emoji
func (v *MascotView) StartEditEmoji(idx int) {
	cat := v.Categories[v.Selected]
	if idx < len(*cat.Emojis) {
		v.EmojiInput.SetValue((*cat.Emojis)[idx])
		v.EmojiInput.Focus()
		v.EditingEmoji = true
	}
}

// StopEditEmoji saves the emoji edit
func (v *MascotView) StopEditEmoji() {
	cat := v.Categories[v.Selected]
	var emojiIdx int
	if cat.Threshold != nil {
		emojiIdx = v.SubSelected - 2
	} else {
		emojiIdx = v.SubSelected - 1
	}
	if emojiIdx >= 0 && emojiIdx < len(*cat.Emojis) {
		(*cat.Emojis)[emojiIdx] = v.EmojiInput.Value()
	}
	v.EmojiInput.Blur()
	v.EditingEmoji = false
}

// CancelEditEmoji cancels the emoji edit
func (v *MascotView) CancelEditEmoji() {
	v.EmojiInput.Blur()
	v.EditingEmoji = false
}

// StartEditThreshold begins editing the threshold
func (v *MascotView) StartEditThreshold() {
	cat := v.Categories[v.Selected]
	if cat.Threshold != nil {
		v.ThresholdInput.SetValue(strconv.Itoa(*cat.Threshold))
		v.ThresholdInput.Focus()
		v.EditingThreshold = true
	}
}

// StopEditThreshold saves the threshold edit
func (v *MascotView) StopEditThreshold() {
	cat := v.Categories[v.Selected]
	if cat.Threshold != nil {
		if val, err := strconv.Atoi(v.ThresholdInput.Value()); err == nil {
			*cat.Threshold = val
		}
	}
	v.ThresholdInput.Blur()
	v.EditingThreshold = false
}

// CancelEditThreshold cancels the threshold edit
func (v *MascotView) CancelEditThreshold() {
	v.ThresholdInput.Blur()
	v.EditingThreshold = false
}

// AddEmoji adds a new emoji to the current category
func (v *MascotView) AddEmoji() {
	if v.InCategory && !v.EditingEmoji && !v.EditingThreshold {
		cat := v.Categories[v.Selected]
		*cat.Emojis = append(*cat.Emojis, "New Emoji")
	}
}

// DeleteEmoji removes the selected emoji
func (v *MascotView) DeleteEmoji() {
	if !v.InCategory || v.EditingEmoji || v.EditingThreshold {
		return
	}
	cat := v.Categories[v.Selected]
	var emojiIdx int
	if cat.Threshold != nil {
		emojiIdx = v.SubSelected - 2
	} else {
		emojiIdx = v.SubSelected - 1
	}
	if emojiIdx >= 0 && emojiIdx < len(*cat.Emojis) && len(*cat.Emojis) > 1 {
		emojis := *cat.Emojis
		*cat.Emojis = append(emojis[:emojiIdx], emojis[emojiIdx+1:]...)
		if v.SubSelected >= len(*cat.Emojis)+2 {
			v.SubSelected--
		}
	}
}

// CurrentInput returns the currently active input
func (v *MascotView) CurrentInput() *textinput.Model {
	if v.EditingEmoji {
		return &v.EmojiInput
	}
	if v.EditingThreshold {
		return &v.ThresholdInput
	}
	return nil
}

// Render returns the mascot view string
func (v *MascotView) Render() string {
	var b strings.Builder

	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#7C3AED"))

	selectedStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#10B981")).
		Bold(true)

	normalStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#9CA3AF"))

	valueStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#F59E0B"))

	descStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#6B7280")).
		Italic(true)

	checkStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#10B981"))

	uncheckStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#6B7280"))

	editingStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#3B82F6")).
		Bold(true)

	b.WriteString(titleStyle.Render("Mascot Settings"))
	b.WriteString("\n\n")

	if !v.InCategory {
		// Show category list
		for i, cat := range v.Categories {
			var checkbox string
			if *cat.Enabled {
				checkbox = checkStyle.Render("[x]")
			} else {
				checkbox = uncheckStyle.Render("[ ]")
			}

			var label string
			if i == v.Selected {
				label = selectedStyle.Render(cat.Label)
			} else {
				label = normalStyle.Render(cat.Label)
			}

			b.WriteString("  " + checkbox + " " + label)
			if i == v.Selected {
				b.WriteString("\n")
				b.WriteString(descStyle.Render("      " + cat.Description))
			}
			b.WriteString("\n")
		}

		b.WriteString("\n")
		b.WriteString(descStyle.Render("  [enter] Edit category  [space/x] Toggle  [esc] Back"))
	} else {
		// Show category details
		cat := v.Categories[v.Selected]
		b.WriteString(titleStyle.Render("  " + cat.Label))
		b.WriteString("\n\n")

		// Enabled toggle
		var checkbox string
		if *cat.Enabled {
			checkbox = checkStyle.Render("[x]")
		} else {
			checkbox = uncheckStyle.Render("[ ]")
		}
		enabledLabel := "Enabled"
		if v.SubSelected == 0 {
			enabledLabel = selectedStyle.Render(enabledLabel)
		} else {
			enabledLabel = normalStyle.Render(enabledLabel)
		}
		b.WriteString("    " + checkbox + " " + enabledLabel + "\n")

		// Threshold (if applicable)
		itemOffset := 1
		if cat.Threshold != nil {
			thresholdLabel := "Threshold"
			var thresholdValue string
			if v.EditingThreshold && v.SubSelected == 1 {
				thresholdValue = editingStyle.Render(v.ThresholdInput.View())
			} else {
				thresholdValue = valueStyle.Render(fmt.Sprintf("%d", *cat.Threshold))
			}
			if v.SubSelected == 1 {
				thresholdLabel = selectedStyle.Render(thresholdLabel)
			} else {
				thresholdLabel = normalStyle.Render(thresholdLabel)
			}
			b.WriteString("    " + thresholdLabel + ": " + thresholdValue + "\n")
			itemOffset = 2
		}

		// Emojis
		b.WriteString("\n")
		b.WriteString(normalStyle.Render("    Emojis:") + "\n")
		for i, emoji := range *cat.Emojis {
			var emojiDisplay string
			if v.EditingEmoji && v.SubSelected == i+itemOffset {
				emojiDisplay = editingStyle.Render(v.EmojiInput.View())
			} else {
				emojiDisplay = valueStyle.Render(emoji)
			}

			if v.SubSelected == i+itemOffset {
				b.WriteString(selectedStyle.Render("      > ") + emojiDisplay + "\n")
			} else {
				b.WriteString("        " + emojiDisplay + "\n")
			}
		}

		b.WriteString("\n")
		if v.EditingEmoji || v.EditingThreshold {
			b.WriteString(descStyle.Render("  [enter] Save  [esc] Cancel"))
		} else {
			b.WriteString(descStyle.Render("  [enter] Edit  [a] Add emoji  [d] Delete  [esc] Back"))
		}
	}

	return b.String()
}
