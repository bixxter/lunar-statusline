package ui

import (
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"statusline-config/config"
	"statusline-config/ui/views"
)

// Screen represents the current screen
type Screen int

const (
	ScreenMenu Screen = iota
	ScreenSections
	ScreenIcons
	ScreenMascot
	ScreenDisplay
	ScreenConfirmQuit
)

// Model is the main Bubble Tea model
type Model struct {
	Config       *config.Config
	OrigConfig   *config.Config // For dirty tracking
	Screen       Screen
	Keys         KeyMap
	Dirty        bool
	Width        int
	Height       int
	Error        string
	ShowHelp     bool
	ConfirmQuit  bool

	// Views
	MenuView     *views.MenuView
	SectionsView *views.SectionsView
	IconsView    *views.IconsView
	MascotView   *views.MascotView
	DisplayView  *views.DisplayView
	PreviewView  *views.PreviewView
}

// NewModel creates a new model
func NewModel(cfg *config.Config) Model {
	// Deep copy for dirty tracking
	origCfg := *cfg

	return Model{
		Config:       cfg,
		OrigConfig:   &origCfg,
		Screen:       ScreenMenu,
		Keys:         DefaultKeyMap(),
		Width:        80,
		Height:       24,
		MenuView:     views.NewMenuView(),
		SectionsView: views.NewSectionsView(cfg),
		IconsView:    views.NewIconsView(cfg),
		MascotView:   views.NewMascotView(cfg),
		DisplayView:  views.NewDisplayView(cfg),
		PreviewView:  views.NewPreviewView(cfg),
	}
}

// Init initializes the model
func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

// Update handles messages
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height
		return m, nil

	case tea.KeyMsg:
		// Handle confirm quit dialog
		if m.ConfirmQuit {
			switch msg.String() {
			case "y", "Y":
				return m, tea.Quit
			case "n", "N", "esc":
				m.ConfirmQuit = false
				return m, nil
			case "s", "S":
				// Save and quit
				if err := config.Save(m.Config); err != nil {
					m.Error = err.Error()
				} else {
					m.Dirty = false
				}
				return m, tea.Quit
			}
			return m, nil
		}

		// Global keys
		switch msg.String() {
		case "ctrl+c":
			if m.Dirty {
				m.ConfirmQuit = true
				return m, nil
			}
			return m, tea.Quit
		case "ctrl+s":
			if err := config.Save(m.Config); err != nil {
				m.Error = err.Error()
			} else {
				m.Dirty = false
				m.Error = ""
			}
			return m, nil
		}

		// Screen-specific handling
		switch m.Screen {
		case ScreenMenu:
			return m.updateMenu(msg)
		case ScreenSections:
			return m.updateSections(msg)
		case ScreenIcons:
			return m.updateIcons(msg)
		case ScreenMascot:
			return m.updateMascot(msg)
		case ScreenDisplay:
			return m.updateDisplay(msg)
		}
	}

	return m, nil
}

func (m Model) updateMenu(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "up", "k":
		m.MenuView.Up()
	case "down", "j":
		m.MenuView.Down()
	case "enter":
		selected := m.MenuView.Selected
		switch selected {
		case 0:
			m.Screen = ScreenSections
		case 1:
			m.Screen = ScreenIcons
		case 2:
			m.Screen = ScreenMascot
		case 3:
			m.Screen = ScreenDisplay
		case 5: // Save & Exit (index 5 because of separator)
			if err := config.Save(m.Config); err != nil {
				m.Error = err.Error()
			} else {
				m.Dirty = false
			}
			return m, tea.Quit
		}
	case "q":
		if m.Dirty {
			m.ConfirmQuit = true
			return m, nil
		}
		return m, tea.Quit
	}
	return m, nil
}

func (m Model) updateSections(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "up", "k":
		m.SectionsView.Up()
	case "down", "j":
		m.SectionsView.Down()
	case " ", "x", "enter":
		m.SectionsView.Toggle()
		m.Dirty = true
	case "esc":
		m.Screen = ScreenMenu
	case "q":
		m.Screen = ScreenMenu
	}
	return m, nil
}

func (m Model) updateIcons(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if m.IconsView.Editing {
		switch msg.String() {
		case "enter":
			m.IconsView.StopEdit()
			m.Dirty = true
			return m, nil
		case "esc":
			m.IconsView.StopEdit()
			return m, nil
		default:
			// Forward to text input
			var cmd tea.Cmd
			input := m.IconsView.CurrentInput()
			*input, cmd = input.Update(msg)
			return m, cmd
		}
	}

	switch msg.String() {
	case "up", "k":
		m.IconsView.Up()
	case "down", "j":
		m.IconsView.Down()
	case "enter", "e":
		m.IconsView.StartEdit()
	case "esc", "q":
		m.Screen = ScreenMenu
	}
	return m, nil
}

func (m Model) updateMascot(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if m.MascotView.EditingEmoji || m.MascotView.EditingThreshold {
		switch msg.String() {
		case "enter":
			m.MascotView.Enter()
			m.Dirty = true
			return m, nil
		case "esc":
			m.MascotView.Back()
			return m, nil
		default:
			// Forward to text input
			input := m.MascotView.CurrentInput()
			if input != nil {
				var cmd tea.Cmd
				*input, cmd = input.Update(msg)
				return m, cmd
			}
		}
		return m, nil
	}

	switch msg.String() {
	case "up", "k":
		m.MascotView.Up()
	case "down", "j":
		m.MascotView.Down()
	case "enter":
		m.MascotView.Enter()
		m.Dirty = true
	case " ", "x":
		if !m.MascotView.InCategory {
			// Toggle enabled on category
			cat := &m.MascotView.Categories[m.MascotView.Selected]
			*cat.Enabled = !*cat.Enabled
			m.Dirty = true
		} else if m.MascotView.SubSelected == 0 {
			// Toggle enabled in category view
			m.MascotView.Enter()
			m.Dirty = true
		}
	case "a":
		m.MascotView.AddEmoji()
		m.Dirty = true
	case "d", "backspace":
		m.MascotView.DeleteEmoji()
		m.Dirty = true
	case "esc", "q":
		if m.MascotView.Back() {
			m.Screen = ScreenMenu
		}
	}
	return m, nil
}

func (m Model) updateDisplay(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if m.DisplayView.Editing {
		switch msg.String() {
		case "enter":
			m.DisplayView.StopEdit()
			m.Dirty = true
			return m, nil
		case "esc":
			m.DisplayView.CancelEdit()
			return m, nil
		default:
			// Forward to text input
			var cmd tea.Cmd
			input := m.DisplayView.CurrentInput()
			*input, cmd = input.Update(msg)
			return m, cmd
		}
	}

	switch msg.String() {
	case "up", "k":
		m.DisplayView.Up()
	case "down", "j":
		m.DisplayView.Down()
	case "enter", "e":
		m.DisplayView.StartEdit()
	case "esc", "q":
		m.Screen = ScreenMenu
	}
	return m, nil
}

// View renders the model
func (m Model) View() string {
	var b strings.Builder

	// Header
	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#7C3AED")).
		Background(lipgloss.Color("#1F2937")).
		Padding(0, 2).
		Width(60)

	versionStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#6B7280")).
		Align(lipgloss.Right)

	header := headerStyle.Render("Claude Statusline Configurator")
	version := versionStyle.Render("v1.0")
	b.WriteString(lipgloss.JoinHorizontal(lipgloss.Top, header, "  ", version))
	b.WriteString("\n")

	// Dirty indicator
	if m.Dirty {
		dirtyStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#F59E0B")).
			Bold(true)
		b.WriteString(dirtyStyle.Render("  [Unsaved Changes]"))
		b.WriteString("\n")
	}

	// Error message
	if m.Error != "" {
		errorStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#EF4444")).
			Bold(true)
		b.WriteString(errorStyle.Render("  Error: " + m.Error))
		b.WriteString("\n")
	}

	b.WriteString("\n")

	// Confirm quit dialog
	if m.ConfirmQuit {
		dialogStyle := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#F59E0B")).
			Padding(1, 2)

		dialog := "You have unsaved changes!\n\n"
		dialog += "[s] Save and quit\n"
		dialog += "[y] Quit without saving\n"
		dialog += "[n] Cancel"

		b.WriteString(dialogStyle.Render(dialog))
		return b.String()
	}

	// Main content
	contentStyle := lipgloss.NewStyle().
		Padding(0, 2)

	switch m.Screen {
	case ScreenMenu:
		b.WriteString(contentStyle.Render(m.MenuView.Render()))
	case ScreenSections:
		b.WriteString(contentStyle.Render(m.SectionsView.Render()))
	case ScreenIcons:
		b.WriteString(contentStyle.Render(m.IconsView.Render()))
	case ScreenMascot:
		b.WriteString(contentStyle.Render(m.MascotView.Render()))
	case ScreenDisplay:
		b.WriteString(contentStyle.Render(m.DisplayView.Render()))
	}

	// Preview
	b.WriteString("\n")
	b.WriteString(contentStyle.Render(m.PreviewView.Render()))

	// Footer help
	b.WriteString("\n\n")
	footerStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#6B7280"))

	helpText := "[ctrl+s] Save  [q] Quit"
	if m.Screen != ScreenMenu {
		helpText = "[esc] Back  " + helpText
	}
	b.WriteString(footerStyle.Render("  " + helpText))

	return b.String()
}
