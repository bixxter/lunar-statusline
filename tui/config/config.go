package config

// Config represents the complete statusline configuration
type Config struct {
	Version         string          `json:"version"`
	EnabledSections EnabledSections `json:"enabled_sections"`
	Colors          Colors          `json:"colors"`
	Icons           Icons           `json:"icons"`
	Mascot          Mascot          `json:"mascot"`
	Thresholds      Thresholds      `json:"thresholds"`
	Display         Display         `json:"display"`
}

// EnabledSections controls which sections are displayed
type EnabledSections struct {
	Git          bool `json:"git"`
	Directory    bool `json:"directory"`
	Model        bool `json:"model"`
	ContextMoons bool `json:"context_moons"`
	TokenCount   bool `json:"token_count"`
	Percentage   bool `json:"percentage"`
	Mascot       bool `json:"mascot"`
}

// Colors defines the color scheme
type Colors struct {
	Directory string `json:"directory"`
	GitClean  string `json:"git_clean"`
	GitDirty  string `json:"git_dirty"`
	Model     string `json:"model"`
	Text      string `json:"text"`
}

// Icons defines the emoji/icon set
type Icons struct {
	GitClean  string   `json:"git_clean"`
	GitDirty  string   `json:"git_dirty"`
	Directory string   `json:"directory"`
	Moons     []string `json:"moons"`
}

// Mascot defines the mascot behavior settings
type Mascot struct {
	ContextPanic MascotState   `json:"context_panic"`
	Productive   MascotState   `json:"productive"`
	Deletion     MascotState   `json:"deletion"`
	TimeBased    TimeBasedMood `json:"time_based"`
}

// MascotState represents a single mascot mood state
type MascotState struct {
	Enabled   bool     `json:"enabled"`
	Threshold int      `json:"threshold"`
	Emojis    []string `json:"emojis"`
}

// TimeBasedMood represents time-of-day moods
type TimeBasedMood struct {
	Enabled   bool     `json:"enabled"`
	Night     []string `json:"night"`
	Morning   []string `json:"morning"`
	Afternoon []string `json:"afternoon"`
	Evening   []string `json:"evening"`
}

// Thresholds defines various threshold values
type Thresholds struct {
	MoonPhases          []int `json:"moon_phases"`
	DirectoryMaxLength  int   `json:"directory_max_length"`
	DirectoryTruncateTo int   `json:"directory_truncate_to"`
	TokenKFormat        int   `json:"token_k_format"`
}

// Display defines display formatting options
type Display struct {
	Separator string `json:"separator"`
}

// DefaultConfig returns a config with sensible defaults
func DefaultConfig() *Config {
	return &Config{
		Version: "1.0",
		EnabledSections: EnabledSections{
			Git:          true,
			Directory:    true,
			Model:        true,
			ContextMoons: true,
			TokenCount:   true,
			Percentage:   true,
			Mascot:       true,
		},
		Colors: Colors{
			Directory: "bright_blue",
			GitClean:  "bright_green",
			GitDirty:  "bright_red",
			Model:     "bright_cyan",
			Text:      "default",
		},
		Icons: Icons{
			GitClean:  "✅",
			GitDirty:  "⚠️",
			Directory: "🗂️",
			Moons:     []string{"🌑", "🌘", "🌗", "🌖", "🌕"},
		},
		Mascot: Mascot{
			ContextPanic: MascotState{
				Enabled:   true,
				Threshold: 90,
				Emojis:    []string{"😱 HELP!", "🆘 SOS!"},
			},
			Productive: MascotState{
				Enabled:   true,
				Threshold: 100,
				Emojis:    []string{"⚡ POWER!", "🔥 HOT!"},
			},
			Deletion: MascotState{
				Enabled:   true,
				Threshold: 30,
				Emojis:    []string{"🗑️ TRASH!", "💀 DELETE!"},
			},
			TimeBased: TimeBasedMood{
				Enabled:   true,
				Night:     []string{"🌃 NIGHT"},
				Morning:   []string{"🌄 DAWN"},
				Afternoon: []string{"☀️ DAY"},
				Evening:   []string{"🌇 DUSK"},
			},
		},
		Thresholds: Thresholds{
			MoonPhases:          []int{20, 40, 60, 80},
			DirectoryMaxLength:  15,
			DirectoryTruncateTo: 12,
			TokenKFormat:        1000,
		},
		Display: Display{
			Separator: " • ",
		},
	}
}
