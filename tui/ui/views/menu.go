package views

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// MenuItem represents a menu option
type MenuItem struct {
	Label       string
	Description string
	IsSeparator bool
}

// MenuView handles the main menu rendering
type MenuView struct {
	Items    []MenuItem
	Selected int
	Width    int
}

// NewMenuView creates a new menu with default items
func NewMenuView() *MenuView {
	return &MenuView{
		Items: []MenuItem{
			{Label: "Sections", Description: "Toggle which sections are displayed"},
			{Label: "Icons & Emojis", Description: "Customize icons and emojis"},
			{Label: "Mascot Settings", Description: "Configure mascot moods and triggers"},
			{Label: "Display Options", Description: "Separator and formatting settings"},
			{Label: "Notifications", Description: "Configure alerts, sounds, and notification triggers"},
			{IsSeparator: true},
			{Label: "Save & Apply", Description: "Save config and install statusline to ~/.claude/"},
			{Label: "Save Config Only", Description: "Save config without installing globally"},
		},
		Selected: 0,
		Width:    50,
	}
}

// Up moves selection up
func (m *MenuView) Up() {
	m.Selected--
	if m.Selected < 0 {
		m.Selected = len(m.Items) - 1
	}
	// Skip separators
	if m.Items[m.Selected].IsSeparator {
		m.Up()
	}
}

// Down moves selection down
func (m *MenuView) Down() {
	m.Selected++
	if m.Selected >= len(m.Items) {
		m.Selected = 0
	}
	// Skip separators
	if m.Items[m.Selected].IsSeparator {
		m.Down()
	}
}

// SelectedItem returns the currently selected item
func (m *MenuView) SelectedItem() MenuItem {
	return m.Items[m.Selected]
}

// Render returns the menu view string
func (m *MenuView) Render() string {
	var b strings.Builder

	selectedStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#10B981")).
		Bold(true)

	normalStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#9CA3AF"))

	separatorStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#4B5563"))

	descStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#6B7280")).
		Italic(true)

	for i, item := range m.Items {
		if item.IsSeparator {
			b.WriteString(separatorStyle.Render("  " + strings.Repeat("─", 30)))
			b.WriteString("\n")
			continue
		}

		var line string
		if i == m.Selected {
			line = selectedStyle.Render("  > " + item.Label)
		} else {
			line = normalStyle.Render("    " + item.Label)
		}

		b.WriteString(line)
		if i == m.Selected && item.Description != "" {
			b.WriteString("\n")
			b.WriteString(descStyle.Render("      " + item.Description))
		}
		b.WriteString("\n")
	}

	return b.String()
}
