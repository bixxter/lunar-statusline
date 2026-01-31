package views

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"statusline-config/config"
)

// PreviewView renders a live preview of the statusline
type PreviewView struct {
	Config *config.Config
}

// NewPreviewView creates a new preview view
func NewPreviewView(cfg *config.Config) *PreviewView {
	return &PreviewView{Config: cfg}
}

// Render returns the preview string
func (v *PreviewView) Render() string {
	var parts []string

	cfg := v.Config

	// Git section
	if cfg.EnabledSections.Git {
		parts = append(parts, cfg.Icons.GitClean+" main")
	}

	// Directory section
	if cfg.EnabledSections.Directory {
		parts = append(parts, cfg.Icons.Directory+" project")
	}

	// Model section
	if cfg.EnabledSections.Model {
		parts = append(parts, "Sonnet")
	}

	// Context moons
	if cfg.EnabledSections.ContextMoons {
		moons := ""
		if len(cfg.Icons.Moons) >= 3 {
			moons = cfg.Icons.Moons[0] + cfg.Icons.Moons[1] + cfg.Icons.Moons[2]
		} else {
			moons = "🌑🌘🌗"
		}
		moonPart := moons
		if cfg.EnabledSections.TokenCount {
			moonPart += " 12k"
		}
		if cfg.EnabledSections.Percentage {
			moonPart += " (45%)"
		}
		parts = append(parts, moonPart)
	} else {
		// Show tokens/percentage without moons
		if cfg.EnabledSections.TokenCount || cfg.EnabledSections.Percentage {
			tokenPart := ""
			if cfg.EnabledSections.TokenCount {
				tokenPart = "12k"
			}
			if cfg.EnabledSections.Percentage {
				if tokenPart != "" {
					tokenPart += " "
				}
				tokenPart += "(45%)"
			}
			parts = append(parts, tokenPart)
		}
	}

	// Mascot
	if cfg.EnabledSections.Mascot {
		// Show animation frames if enabled, otherwise first emoji
		var mascotDisplay string
		if cfg.Mascot.TimeBased.Enabled && len(cfg.Mascot.TimeBased.Afternoon) > 0 {
			if cfg.Mascot.TimeBased.Animate && len(cfg.Mascot.TimeBased.Afternoon) > 1 {
				// Show animation sequence preview
				mascotDisplay = strings.Join(cfg.Mascot.TimeBased.Afternoon, "→")
			} else {
				mascotDisplay = cfg.Mascot.TimeBased.Afternoon[0]
			}
		} else {
			mascotDisplay = "💻→⌨️→🖱️→⌨️"
		}
		parts = append(parts, mascotDisplay)
	}

	separator := cfg.Display.Separator
	if separator == "" {
		separator = " | "
	}

	preview := strings.Join(parts, separator)

	// Style the preview without a border box for cleaner full-width display
	previewStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFFFFF")).
		Padding(0, 1)

	labelStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#6B7280")).
		Italic(true)

	return labelStyle.Render("Preview:") + "\n" + previewStyle.Render(preview)
}
