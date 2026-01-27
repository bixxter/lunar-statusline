package views

import (
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"
	"statusline-config/config"
)

// IconItem represents an editable icon
type IconItem struct {
	Key         string
	Label       string
	Description string
	Input       textinput.Model
}

// IconsView handles the icons editing screen
type IconsView struct {
	Items    []IconItem
	Selected int
	Editing  bool
	Config   *config.Config
}

// NewIconsView creates a new icons view
func NewIconsView(cfg *config.Config) *IconsView {
	items := []IconItem{
		{Key: "git_clean", Label: "Git Clean", Description: "Icon when git status is clean"},
		{Key: "git_dirty", Label: "Git Dirty", Description: "Icon when there are uncommitted changes"},
		{Key: "directory", Label: "Directory", Description: "Icon for directory name"},
		{Key: "moon_1", Label: "Moon Phase 1", Description: "First moon phase (0-20%)"},
		{Key: "moon_2", Label: "Moon Phase 2", Description: "Second moon phase (20-40%)"},
		{Key: "moon_3", Label: "Moon Phase 3", Description: "Third moon phase (40-60%)"},
		{Key: "moon_4", Label: "Moon Phase 4", Description: "Fourth moon phase (60-80%)"},
		{Key: "moon_5", Label: "Moon Phase 5", Description: "Fifth moon phase (80-100%)"},
	}

	// Initialize text inputs
	for i := range items {
		ti := textinput.New()
		ti.CharLimit = 10
		ti.Width = 10
		items[i].Input = ti
	}

	view := &IconsView{
		Items:    items,
		Selected: 0,
		Editing:  false,
		Config:   cfg,
	}

	view.LoadFromConfig()
	return view
}

// LoadFromConfig populates inputs from config
func (v *IconsView) LoadFromConfig() {
	v.Items[0].Input.SetValue(v.Config.Icons.GitClean)
	v.Items[1].Input.SetValue(v.Config.Icons.GitDirty)
	v.Items[2].Input.SetValue(v.Config.Icons.Directory)

	for i := 0; i < 5 && i < len(v.Config.Icons.Moons); i++ {
		v.Items[3+i].Input.SetValue(v.Config.Icons.Moons[i])
	}
}

// SaveToConfig writes inputs back to config
func (v *IconsView) SaveToConfig() {
	v.Config.Icons.GitClean = v.Items[0].Input.Value()
	v.Config.Icons.GitDirty = v.Items[1].Input.Value()
	v.Config.Icons.Directory = v.Items[2].Input.Value()

	moons := make([]string, 5)
	for i := 0; i < 5; i++ {
		moons[i] = v.Items[3+i].Input.Value()
	}
	v.Config.Icons.Moons = moons
}

// Up moves selection up
func (v *IconsView) Up() {
	if !v.Editing {
		v.Selected--
		if v.Selected < 0 {
			v.Selected = len(v.Items) - 1
		}
	}
}

// Down moves selection down
func (v *IconsView) Down() {
	if !v.Editing {
		v.Selected++
		if v.Selected >= len(v.Items) {
			v.Selected = 0
		}
	}
}

// StartEdit begins editing the selected item
func (v *IconsView) StartEdit() {
	v.Editing = true
	v.Items[v.Selected].Input.Focus()
}

// StopEdit finishes editing
func (v *IconsView) StopEdit() {
	v.Editing = false
	v.Items[v.Selected].Input.Blur()
	v.SaveToConfig()
}

// CurrentInput returns the currently active input
func (v *IconsView) CurrentInput() *textinput.Model {
	return &v.Items[v.Selected].Input
}

// Render returns the icons view string
func (v *IconsView) Render() string {
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

	b.WriteString(titleStyle.Render("Icons & Emojis"))
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
			value = editingStyle.Render(item.Input.View())
		} else {
			value = valueStyle.Render(item.Input.Value())
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
		b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("#6B7280")).Render(
			"  [enter] Save  [esc] Cancel"))
	} else {
		b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("#6B7280")).Render(
			"  [enter/e] Edit  [esc] Back"))
	}

	return b.String()
}
