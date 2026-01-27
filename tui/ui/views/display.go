package views

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"
	"statusline-config/config"
)

// DisplayItem represents a display option
type DisplayItem struct {
	Key         string
	Label       string
	Description string
	IsString    bool // true for string values, false for int
}

// DisplayView handles the display options screen
type DisplayView struct {
	Items    []DisplayItem
	Selected int
	Editing  bool
	Input    textinput.Model
	Config   *config.Config
}

// NewDisplayView creates a new display view
func NewDisplayView(cfg *config.Config) *DisplayView {
	ti := textinput.New()
	ti.CharLimit = 20
	ti.Width = 15

	return &DisplayView{
		Items: []DisplayItem{
			{Key: "separator", Label: "Separator", Description: "Text between sections", IsString: true},
			{Key: "dir_max_len", Label: "Directory Max Length", Description: "Maximum directory name length", IsString: false},
			{Key: "dir_truncate", Label: "Directory Truncate To", Description: "Length to truncate directory to", IsString: false},
			{Key: "token_k_format", Label: "Token K Format", Description: "Threshold for showing as 'k' format", IsString: false},
		},
		Selected: 0,
		Editing:  false,
		Input:    ti,
		Config:   cfg,
	}
}

// GetValue returns the current value for an item
func (v *DisplayView) GetValue(item DisplayItem) string {
	switch item.Key {
	case "separator":
		return v.Config.Display.Separator
	case "dir_max_len":
		return strconv.Itoa(v.Config.Thresholds.DirectoryMaxLength)
	case "dir_truncate":
		return strconv.Itoa(v.Config.Thresholds.DirectoryTruncateTo)
	case "token_k_format":
		return strconv.Itoa(v.Config.Thresholds.TokenKFormat)
	}
	return ""
}

// SetValue sets the value for an item
func (v *DisplayView) SetValue(item DisplayItem, value string) {
	switch item.Key {
	case "separator":
		v.Config.Display.Separator = value
	case "dir_max_len":
		if val, err := strconv.Atoi(value); err == nil {
			v.Config.Thresholds.DirectoryMaxLength = val
		}
	case "dir_truncate":
		if val, err := strconv.Atoi(value); err == nil {
			v.Config.Thresholds.DirectoryTruncateTo = val
		}
	case "token_k_format":
		if val, err := strconv.Atoi(value); err == nil {
			v.Config.Thresholds.TokenKFormat = val
		}
	}
}

// Up moves selection up
func (v *DisplayView) Up() {
	if !v.Editing {
		v.Selected--
		if v.Selected < 0 {
			v.Selected = len(v.Items) - 1
		}
	}
}

// Down moves selection down
func (v *DisplayView) Down() {
	if !v.Editing {
		v.Selected++
		if v.Selected >= len(v.Items) {
			v.Selected = 0
		}
	}
}

// StartEdit begins editing the selected item
func (v *DisplayView) StartEdit() {
	item := v.Items[v.Selected]
	v.Input.SetValue(v.GetValue(item))
	v.Input.Focus()
	v.Editing = true
}

// StopEdit finishes editing and saves
func (v *DisplayView) StopEdit() {
	item := v.Items[v.Selected]
	v.SetValue(item, v.Input.Value())
	v.Input.Blur()
	v.Editing = false
}

// CancelEdit cancels editing
func (v *DisplayView) CancelEdit() {
	v.Input.Blur()
	v.Editing = false
}

// CurrentInput returns the input model
func (v *DisplayView) CurrentInput() *textinput.Model {
	return &v.Input
}

// Render returns the display view string
func (v *DisplayView) Render() string {
	var b strings.Builder

	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#7C3AED")).
		MarginBottom(1)

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

	editingStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#3B82F6")).
		Bold(true)

	b.WriteString(titleStyle.Render("Display Options"))
	b.WriteString("\n\n")

	for i, item := range v.Items {
		var label string
		if i == v.Selected {
			label = selectedStyle.Render(item.Label)
		} else {
			label = normalStyle.Render(item.Label)
		}

		var value string
		if v.Editing && i == v.Selected {
			value = editingStyle.Render(v.Input.View())
		} else {
			val := v.GetValue(item)
			if item.IsString {
				value = valueStyle.Render(fmt.Sprintf("%q", val))
			} else {
				value = valueStyle.Render(val)
			}
		}

		b.WriteString("  " + label + ": " + value)
		if i == v.Selected {
			b.WriteString("\n")
			b.WriteString(descStyle.Render("      " + item.Description))
		}
		b.WriteString("\n")
	}

	b.WriteString("\n")
	if v.Editing {
		b.WriteString(descStyle.Render("  [enter] Save  [esc] Cancel"))
	} else {
		b.WriteString(descStyle.Render("  [enter/e] Edit  [esc] Back"))
	}

	return b.String()
}
