package views

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"statusline-config/config"
)

// SectionItem represents a toggleable section
type SectionItem struct {
	Key         string
	Label       string
	Description string
	Enabled     *bool
}

// SectionsView handles the sections toggle screen
type SectionsView struct {
	Items    []SectionItem
	Selected int
	Config   *config.Config
}

// NewSectionsView creates a new sections view
func NewSectionsView(cfg *config.Config) *SectionsView {
	return &SectionsView{
		Config: cfg,
		Items: []SectionItem{
			{Key: "waiting_indicator", Label: "Waiting Indicator", Description: "Show alert when Claude needs your input", Enabled: &cfg.EnabledSections.WaitingIndicator},
			{Key: "git", Label: "Git Branch", Description: "Show current git branch and status", Enabled: &cfg.EnabledSections.Git},
			{Key: "directory", Label: "Directory", Description: "Show current directory name", Enabled: &cfg.EnabledSections.Directory},
			{Key: "model", Label: "Model Name", Description: "Show Claude model in use", Enabled: &cfg.EnabledSections.Model},
			{Key: "context_moons", Label: "Context Moons", Description: "Visual moon phases for context usage", Enabled: &cfg.EnabledSections.ContextMoons},
			{Key: "token_count", Label: "Token Count", Description: "Show token count (e.g., 12k)", Enabled: &cfg.EnabledSections.TokenCount},
			{Key: "percentage", Label: "Percentage", Description: "Show context usage percentage", Enabled: &cfg.EnabledSections.Percentage},
			{Key: "mascot", Label: "Mascot", Description: "Show reactive mascot emoji", Enabled: &cfg.EnabledSections.Mascot},
		},
		Selected: 0,
	}
}

// UpdateConfig refreshes the view with new config
func (s *SectionsView) UpdateConfig(cfg *config.Config) {
	s.Config = cfg
	s.Items[0].Enabled = &cfg.EnabledSections.WaitingIndicator
	s.Items[1].Enabled = &cfg.EnabledSections.Git
	s.Items[2].Enabled = &cfg.EnabledSections.Directory
	s.Items[3].Enabled = &cfg.EnabledSections.Model
	s.Items[4].Enabled = &cfg.EnabledSections.ContextMoons
	s.Items[5].Enabled = &cfg.EnabledSections.TokenCount
	s.Items[6].Enabled = &cfg.EnabledSections.Percentage
	s.Items[7].Enabled = &cfg.EnabledSections.Mascot
}

// Up moves selection up
func (s *SectionsView) Up() {
	s.Selected--
	if s.Selected < 0 {
		s.Selected = len(s.Items) - 1
	}
}

// Down moves selection down
func (s *SectionsView) Down() {
	s.Selected++
	if s.Selected >= len(s.Items) {
		s.Selected = 0
	}
}

// Toggle toggles the selected item
func (s *SectionsView) Toggle() {
	if s.Selected >= 0 && s.Selected < len(s.Items) {
		*s.Items[s.Selected].Enabled = !*s.Items[s.Selected].Enabled
	}
}

// Render returns the sections view string
func (s *SectionsView) Render() string {
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

	checkStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#10B981"))

	uncheckStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#6B7280"))

	descStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#6B7280")).
		Italic(true)

	b.WriteString(titleStyle.Render("Toggle Sections"))
	b.WriteString("\n\n")

	for i, item := range s.Items {
		var checkbox string
		if *item.Enabled {
			checkbox = checkStyle.Render("[x]")
		} else {
			checkbox = uncheckStyle.Render("[ ]")
		}

		var label string
		if i == s.Selected {
			label = selectedStyle.Render(item.Label)
		} else {
			label = normalStyle.Render(item.Label)
		}

		b.WriteString("  " + checkbox + " " + label)
		if i == s.Selected {
			b.WriteString("\n")
			b.WriteString(descStyle.Render("      " + item.Description))
		}
		b.WriteString("\n")
	}

	b.WriteString("\n")
	b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("#6B7280")).Render(
		"  [space/x] Toggle  [esc] Back"))

	return b.String()
}
