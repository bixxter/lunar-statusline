package config

// Config represents the complete statusline configuration
type Config struct {
	Version          string           `json:"version"`
	EnabledSections  EnabledSections  `json:"enabled_sections"`
	Colors           Colors           `json:"colors"`
	Icons            Icons            `json:"icons"`
	Mascot           Mascot           `json:"mascot"`
	Thresholds       Thresholds       `json:"thresholds"`
	Display          Display          `json:"display"`
	WaitingIndicator WaitingIndicator `json:"waiting_indicator"`
	Notifications    Notifications    `json:"notifications"`
}

// WaitingIndicator settings for when Claude is waiting for user input
type WaitingIndicator struct {
	Enabled bool   `json:"enabled"`
	Icon    string `json:"icon"`
	Text    string `json:"text"`
	Blink   bool   `json:"blink"`
}

// Notifications settings for alerts when Claude needs input
type Notifications struct {
	TerminalBell  NotificationConfig `json:"terminal_bell"`
	Desktop       DesktopNotification `json:"desktop"`
	BlinkingText  NotificationConfig `json:"blinking_text"`
	TerminalTitle TerminalTitleConfig `json:"terminal_title"`
	Tmux          TmuxNotification `json:"tmux"`
}

// NotificationConfig represents a basic notification type
type NotificationConfig struct {
	Enabled          bool `json:"enabled"`
	OnContextPanic   bool `json:"on_context_panic"`
	OnSessionLimit   bool `json:"on_session_limit"`
	ContextThreshold int  `json:"context_threshold"`
	SessionThreshold int  `json:"session_threshold"`
}

// DesktopNotification extends NotificationConfig with desktop-specific options
type DesktopNotification struct {
	Enabled          bool    `json:"enabled"`
	OnContextPanic   bool    `json:"on_context_panic"`
	OnSessionLimit   bool    `json:"on_session_limit"`
	ContextThreshold int     `json:"context_threshold"`
	SessionThreshold int     `json:"session_threshold"`
	Title            string  `json:"title"`
	Sound            bool    `json:"sound"`
	SoundPath        string  `json:"sound_path,omitempty"`
	SoundVolume      float64 `json:"sound_volume,omitempty"`
}

// TerminalTitleConfig for terminal title bar notifications
type TerminalTitleConfig struct {
	Enabled          bool   `json:"enabled"`
	ShowModel        bool   `json:"show_model"`
	ShowContext      bool   `json:"show_context"`
	ShowBranch       bool   `json:"show_branch"`
	AlertOnPanic     bool   `json:"alert_on_panic"`
	PanicPrefix      string `json:"panic_prefix"`
	ContextThreshold int    `json:"context_threshold"`
}

// TmuxNotification for tmux-specific notifications
type TmuxNotification struct {
	Enabled          bool   `json:"enabled"`
	OnContextPanic   bool   `json:"on_context_panic"`
	OnSessionLimit   bool   `json:"on_session_limit"`
	ContextThreshold int    `json:"context_threshold"`
	SessionThreshold int    `json:"session_threshold"`
	DisplayMessage   bool   `json:"display_message"`
	SetWindowStyle   bool   `json:"set_window_style"`
	AlertStyle       string `json:"alert_style"`
}

// EnabledSections controls which sections are displayed
type EnabledSections struct {
	Git              bool `json:"git"`
	Directory        bool `json:"directory"`
	Model            bool `json:"model"`
	ContextMoons     bool `json:"context_moons"`
	TokenCount       bool `json:"token_count"`
	Percentage       bool `json:"percentage"`
	Mascot           bool `json:"mascot"`
	WaitingIndicator bool `json:"waiting_indicator"`
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
	Animate   bool     `json:"animate"`   // If true, cycle through emojis as animation frames
	Speed     int      `json:"speed"`     // Animation speed in milliseconds (default 500)
}

// TimeBasedMood represents time-of-day moods
type TimeBasedMood struct {
	Enabled   bool     `json:"enabled"`
	Night     []string `json:"night"`
	Morning   []string `json:"morning"`
	Afternoon []string `json:"afternoon"`
	Evening   []string `json:"evening"`
	Animate   bool     `json:"animate"`   // If true, cycle through emojis as animation frames
	Speed     int      `json:"speed"`     // Animation speed in milliseconds (default 500)
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
			Git:              true,
			Directory:        true,
			Model:            true,
			ContextMoons:     true,
			TokenCount:       true,
			Percentage:       true,
			Mascot:           true,
			WaitingIndicator: true,
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
				Emojis:    []string{"😰", "😱", "🆘", "😱"},
				Animate:   true,
				Speed:     300,
			},
			Productive: MascotState{
				Enabled:   true,
				Threshold: 100,
				Emojis:    []string{"🔨", "⚒️", "🛠️", "⚒️"},
				Animate:   true,
				Speed:     400,
			},
			Deletion: MascotState{
				Enabled:   true,
				Threshold: 30,
				Emojis:    []string{"🧹", "✨", "🗑️", "✨"},
				Animate:   true,
				Speed:     350,
			},
			TimeBased: TimeBasedMood{
				Enabled:   true,
				Night:     []string{"🦉", "💤", "🌙", "💤"},
				Morning:   []string{"☀️", "🌅", "☕", "🌅"},
				Afternoon: []string{"💻", "⌨️", "🖱️", "⌨️"},
				Evening:   []string{"🌆", "🌇", "🌃", "🌇"},
				Animate:   true,
				Speed:     600,
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
		WaitingIndicator: WaitingIndicator{
			Enabled: true,
			Icon:    "🔔",
			Text:    "WAITING",
			Blink:   true,
		},
		Notifications: Notifications{
			TerminalBell: NotificationConfig{
				Enabled:          true,
				OnContextPanic:   false,
				OnSessionLimit:   false,
				ContextThreshold: 30,
				SessionThreshold: 0,
			},
			Desktop: DesktopNotification{
				Enabled:          true,
				OnContextPanic:   true,
				OnSessionLimit:   false,
				ContextThreshold: 70,
				SessionThreshold: 0,
				Title:            "Context over 70% use /clear or /compact",
				Sound:            true,
				SoundVolume:      1.0,
			},
			BlinkingText: NotificationConfig{
				Enabled:          false,
				OnContextPanic:   false,
				OnSessionLimit:   false,
				ContextThreshold: 0,
				SessionThreshold: 0,
			},
			TerminalTitle: TerminalTitleConfig{
				Enabled:          false,
				ShowModel:        true,
				ShowContext:      false,
				ShowBranch:       false,
				AlertOnPanic:     true,
				PanicPrefix:      "",
				ContextThreshold: 30,
			},
			Tmux: TmuxNotification{
				Enabled:          false,
				OnContextPanic:   false,
				OnSessionLimit:   false,
				ContextThreshold: 0,
				SessionThreshold: 0,
				DisplayMessage:   false,
				SetWindowStyle:   false,
				AlertStyle:       "",
			},
		},
	}
}
