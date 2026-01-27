package ui

import "github.com/charmbracelet/lipgloss"

var (
	// Colors
	primaryColor   = lipgloss.Color("#7C3AED") // Purple
	secondaryColor = lipgloss.Color("#10B981") // Green
	accentColor    = lipgloss.Color("#F59E0B") // Amber
	mutedColor     = lipgloss.Color("#6B7280") // Gray
	errorColor     = lipgloss.Color("#EF4444") // Red

	// Base styles
	TitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(primaryColor).
			MarginBottom(1)

	SubtitleStyle = lipgloss.NewStyle().
			Foreground(mutedColor).
			Italic(true)

	// Menu styles
	MenuItemStyle = lipgloss.NewStyle().
			PaddingLeft(2)

	SelectedMenuItemStyle = lipgloss.NewStyle().
				PaddingLeft(2).
				Foreground(secondaryColor).
				Bold(true)

	// Section styles
	SectionTitleStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(primaryColor).
				MarginBottom(1).
				BorderBottom(true).
				BorderStyle(lipgloss.NormalBorder()).
				BorderForeground(mutedColor)

	// Checkbox styles
	CheckedStyle = lipgloss.NewStyle().
			Foreground(secondaryColor)

	UncheckedStyle = lipgloss.NewStyle().
			Foreground(mutedColor)

	// Input styles
	InputLabelStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Bold(true)

	InputValueStyle = lipgloss.NewStyle().
			Foreground(accentColor)

	// Preview styles
	PreviewBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(mutedColor).
			Padding(0, 1).
			MarginTop(1)

	PreviewLabelStyle = lipgloss.NewStyle().
				Foreground(mutedColor).
				Italic(true)

	// Help styles
	HelpKeyStyle = lipgloss.NewStyle().
			Foreground(accentColor).
			Bold(true)

	HelpDescStyle = lipgloss.NewStyle().
			Foreground(mutedColor)

	// Status styles
	DirtyStyle = lipgloss.NewStyle().
			Foreground(errorColor).
			Bold(true)

	SavedStyle = lipgloss.NewStyle().
			Foreground(secondaryColor)

	// Box styles
	MainBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(primaryColor).
			Padding(1, 2)

	// Footer style
	FooterStyle = lipgloss.NewStyle().
			Foreground(mutedColor).
			MarginTop(1)

	// Cursor for editing
	CursorStyle = lipgloss.NewStyle().
			Foreground(accentColor).
			Background(lipgloss.Color("#1F2937"))
)

// Helper functions for styling
func RenderCheckbox(checked bool, label string) string {
	if checked {
		return CheckedStyle.Render("[x] ") + label
	}
	return UncheckedStyle.Render("[ ] ") + label
}

func RenderKeyHelp(key, desc string) string {
	return HelpKeyStyle.Render(key) + " " + HelpDescStyle.Render(desc)
}
