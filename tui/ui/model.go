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

// ASCII Art variants for different terminal widths
const asciiArtLarge = `
██╗     ██╗   ██╗███╗   ██╗ █████╗ ██████╗     ███████╗██████╗ ██╗████████╗ ██████╗ ██████╗
██║     ██║   ██║████╗  ██║██╔══██╗██╔══██╗    ██╔════╝██╔══██╗██║╚══██╔══╝██╔═══██╗██╔══██╗
██║     ██║   ██║██╔██╗ ██║███████║██████╔╝    █████╗  ██║  ██║██║   ██║   ██║   ██║██████╔╝
██║     ██║   ██║██║╚██╗██║██╔══██║██╔══██╗    ██╔══╝  ██║  ██║██║   ██║   ██║   ██║██╔══██╗
███████╗╚██████╔╝██║ ╚████║██║  ██║██║  ██║    ███████╗██████╔╝██║   ██║   ╚██████╔╝██║  ██║
╚══════╝ ╚═════╝ ╚═╝  ╚═══╝╚═╝  ╚═╝╚═╝  ╚═╝    ╚══════╝╚═════╝ ╚═╝   ╚═╝    ╚═════╝ ╚═╝  ╚═╝`

const asciiArtMedium = `
█░░ █░█ █▄░█ ▄▀█ █▀█   █▀▀ █▀▄ █ ▀█▀ █▀█ █▀█
█▄▄ █▄█ █░▀█ █▀█ █▀▄   ██▄ █▄▀ █ ░█░ █▄█ █▀▄`

const asciiArtSmall = `
╦  ╦ ╦╔╗╔╔═╗╦═╗  ╔═╗╔╦╗╦╔╦╗╔═╗╦═╗
║  ║ ║║║║╠═╣╠╦╝  ║╣  ║║║ ║ ║ ║╠╦╝
╩═╝╚═╝╝╚╝╩ ╩╩╚═  ╚═╝═╩╝╩ ╩ ╚═╝╩╚═`

// Rainbow sparkle characters
var sparkles = []string{"✦", "✧", "★", "✶", "✴", "✵", "❋", "✺", "·", "•"}

// Rainbow colors for gay pride sparkles
var rainbowColors = []string{
	"#FF0000", // Red
	"#FF8C00", // Orange
	"#FFD700", // Yellow
	"#00FF00", // Green
	"#0000FF", // Blue
	"#8B00FF", // Violet
	"#FF69B4", // Pink
	"#00FFFF", // Cyan
}

// Model is the main Bubble Tea model
type Model struct {
	Config      *config.Config
	OrigConfig  *config.Config // For dirty tracking
	Screen      Screen
	Keys        KeyMap
	Dirty       bool
	Width       int
	Height      int
	Error       string
	ShowHelp    bool
	ConfirmQuit bool

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
				// Save, install, and quit
				if err := config.SaveAndInstall(m.Config); err != nil {
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
		case 5: // Save & Apply (index 5 because of separator)
			if err := config.SaveAndInstall(m.Config); err != nil {
				m.Error = "Save failed: " + err.Error()
				return m, nil
			}
			m.Dirty = false
			m.Error = ""
			return m, tea.Quit
		case 6: // Save Config Only
			if err := config.Save(m.Config); err != nil {
				m.Error = "Save failed: " + err.Error()
				return m, nil
			}
			m.Dirty = false
			m.Error = ""
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

// generateSparkles creates a line of rainbow sparkles
func generateSparkles(width int, seed int) string {
	if width <= 0 {
		return ""
	}
	var result strings.Builder
	for i := 0; i < width; i++ {
		// Sparse sparkles - only ~20% of positions have sparkles
		if (i+seed)%5 == 0 {
			sparkle := sparkles[(i+seed)%len(sparkles)]
			color := rainbowColors[(i+seed)%len(rainbowColors)]
			style := lipgloss.NewStyle().Foreground(lipgloss.Color(color)).Bold(true)
			result.WriteString(style.Render(sparkle))
		} else {
			result.WriteString(" ")
		}
	}
	return result.String()
}

// renderHeader renders the ASCII art header with gradient colors and sparkles
func (m Model) renderHeader() string {
	// Choose ASCII art based on terminal width
	var asciiArt string
	var artWidth int

	if m.Width >= 95 {
		asciiArt = asciiArtLarge
		artWidth = 90
	} else if m.Width >= 50 {
		asciiArt = asciiArtMedium
		artWidth = 47
	} else {
		asciiArt = asciiArtSmall
		artWidth = 33
	}

	// Gradient colors from purple to cyan (moon/lunar theme)
	gradientColors := []string{
		"#9333EA", // Purple
		"#7C3AED", // Violet
		"#6366F1", // Indigo
		"#3B82F6", // Blue
		"#0EA5E9", // Sky
		"#06B6D4", // Cyan
	}

	lines := strings.Split(strings.TrimPrefix(asciiArt, "\n"), "\n")
	var coloredLines []string

	// Add top sparkle border
	coloredLines = append(coloredLines, generateSparkles(artWidth, 0))

	for i, line := range lines {
		colorIdx := i % len(gradientColors)
		style := lipgloss.NewStyle().
			Foreground(lipgloss.Color(gradientColors[colorIdx])).
			Bold(true)

		// Add sparkles on the sides
		leftSparkle := sparkles[(i*3)%len(sparkles)]
		rightSparkle := sparkles[(i*3+2)%len(sparkles)]
		leftColor := rainbowColors[(i*2)%len(rainbowColors)]
		rightColor := rainbowColors[(i*2+3)%len(rainbowColors)]

		leftStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(leftColor)).Bold(true)
		rightStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(rightColor)).Bold(true)

		coloredLine := leftStyle.Render(leftSparkle+" ") + style.Render(line) + rightStyle.Render(" "+rightSparkle)
		coloredLines = append(coloredLines, coloredLine)
	}

	// Add bottom sparkle border
	coloredLines = append(coloredLines, generateSparkles(artWidth, 7))

	return strings.Join(coloredLines, "\n")
}

// View renders the model
func (m Model) View() string {
	// Main container that fills the screen (no background for clean look)
	mainStyle := lipgloss.NewStyle().
		Width(m.Width).
		Height(m.Height)

	var content strings.Builder

	// Top padding to ensure ASCII art is visible (extra padding for terminal tabs)
	content.WriteString("\n\n\n\n\n\n\n\n")

	// === HEADER SECTION ===
	header := m.renderHeader()
	headerBox := lipgloss.NewStyle().
		Width(m.Width).
		Align(lipgloss.Center).
		Padding(0, 0).
		Render(header)
	content.WriteString(headerBox)

	// Subtitle with version
	subtitleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#64748B")).
		Italic(true).
		Width(m.Width).
		Align(lipgloss.Center)
	content.WriteString(subtitleStyle.Render("Claude Statusline Configuration Tool  v1.0"))
	content.WriteString("\n")

	// Status bar (dirty indicator / error)
	statusStyle := lipgloss.NewStyle().
		Width(m.Width).
		Align(lipgloss.Center).
		Height(1)

	if m.Error != "" {
		errorStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#EF4444")).
			Bold(true)
		content.WriteString(statusStyle.Render(errorStyle.Render("Error: " + m.Error)))
	} else if m.Dirty {
		dirtyStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#F59E0B")).
			Bold(true)
		content.WriteString(statusStyle.Render(dirtyStyle.Render("● Unsaved Changes")))
	} else {
		savedStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#10B981"))
		content.WriteString(statusStyle.Render(savedStyle.Render("✓ Saved")))
	}
	content.WriteString("\n\n")

	// === MAIN CONTENT SECTION ===
	// Calculate available height for content based on ASCII art size
	var headerHeight int
	if m.Width >= 95 {
		headerHeight = 16 // Large ASCII (6 lines) + sparkle borders (2) + subtitle + status + extra top padding
	} else if m.Width >= 50 {
		headerHeight = 12 // Medium ASCII (2 lines) + sparkle borders (2) + subtitle + status + extra top padding
	} else {
		headerHeight = 13 // Small ASCII (3 lines) + sparkle borders (2) + subtitle + status + extra top padding
	}
	footerHeight := 6 // Preview + help
	contentHeight := m.Height - headerHeight - footerHeight
	if contentHeight < 10 {
		contentHeight = 10
	}

	contentWidth := m.Width - 8
	if contentWidth < 40 {
		contentWidth = 40
	}

	// Content box with border
	contentBoxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#3B82F6")).
		Padding(1, 2).
		Width(contentWidth).
		Height(contentHeight)

	// Confirm quit dialog
	if m.ConfirmQuit {
		dialogStyle := lipgloss.NewStyle().
			Border(lipgloss.DoubleBorder()).
			BorderForeground(lipgloss.Color("#F59E0B")).
			Padding(2, 4).
			Width(40).
			Align(lipgloss.Center)

		titleStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#F59E0B")).
			Bold(true)

		dialog := titleStyle.Render("⚠ Unsaved Changes") + "\n\n"
		dialog += "What would you like to do?\n\n"
		dialog += lipgloss.NewStyle().Foreground(lipgloss.Color("#10B981")).Render("[s]") + " Save & apply globally\n"
		dialog += lipgloss.NewStyle().Foreground(lipgloss.Color("#EF4444")).Render("[y]") + " Quit without saving\n"
		dialog += lipgloss.NewStyle().Foreground(lipgloss.Color("#6B7280")).Render("[n]") + " Cancel"

		dialogRendered := dialogStyle.Render(dialog)
		centeredDialog := lipgloss.NewStyle().
			Width(m.Width).
			Align(lipgloss.Center).
			Render(dialogRendered)
		content.WriteString(centeredDialog)
	} else {
		// Render current screen content
		var screenContent string
		switch m.Screen {
		case ScreenMenu:
			screenContent = m.MenuView.Render()
		case ScreenSections:
			screenContent = m.SectionsView.Render()
		case ScreenIcons:
			screenContent = m.IconsView.Render()
		case ScreenMascot:
			screenContent = m.MascotView.Render()
		case ScreenDisplay:
			screenContent = m.DisplayView.Render()
		}

		contentBox := contentBoxStyle.Render(screenContent)
		centeredContent := lipgloss.NewStyle().
			Width(m.Width).
			Align(lipgloss.Center).
			Render(contentBox)
		content.WriteString(centeredContent)
	}

	content.WriteString("\n")

	// === PREVIEW SECTION ===
	if !m.ConfirmQuit {
		previewStyle := lipgloss.NewStyle().
			Width(m.Width - 4).
			Padding(0, 2)
		preview := previewStyle.Render(m.PreviewView.Render())
		content.WriteString(preview)
		content.WriteString("\n")
	}

	// === FOOTER / HELP ===
	helpStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#64748B")).
		Width(m.Width).
		Align(lipgloss.Center).
		Padding(1, 0)

	keyStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#0EA5E9")).
		Bold(true)

	var helpParts []string
	if m.Screen != ScreenMenu && !m.ConfirmQuit {
		helpParts = append(helpParts, keyStyle.Render("esc")+" back")
	}
	helpParts = append(helpParts, keyStyle.Render("↑↓/jk")+" navigate")
	helpParts = append(helpParts, keyStyle.Render("enter")+" select")
	helpParts = append(helpParts, keyStyle.Render("ctrl+s")+" save")
	helpParts = append(helpParts, keyStyle.Render("q")+" quit")

	helpText := strings.Join(helpParts, "  │  ")
	content.WriteString(helpStyle.Render(helpText))

	// Apply main style and return
	return mainStyle.Render(content.String())
}
