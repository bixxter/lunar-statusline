package views

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"
	"statusline-config/config"
)

// NotificationCategory represents a notification type category
type NotificationCategory struct {
	Key         string
	Label       string
	Description string
	Enabled     *bool
}

// SoundOption represents a sound file option
type SoundOption struct {
	Name string
	Path string
}

// VolumeOption represents a volume level option
type VolumeOption struct {
	Name  string
	Value float64
}

// NotificationsView handles the notifications settings screen
type NotificationsView struct {
	Config           *config.Config
	Categories       []NotificationCategory
	Selected         int
	InCategory       bool
	SubSelected      int
	SoundOptions     []SoundOption
	SoundSelected    int
	SelectingSound   bool
	EditingThreshold bool
	ThresholdInput   textinput.Model
	EditingTitle     bool
	TitleInput       textinput.Model
	SelectingVolume  bool
	VolumeOptions    []VolumeOption
	VolumeSelected   int
}

// NewNotificationsView creates a new notifications view
func NewNotificationsView(cfg *config.Config) *NotificationsView {
	ti := textinput.New()
	ti.Placeholder = "70"
	ti.CharLimit = 3
	ti.Width = 10

	titleInput := textinput.New()
	titleInput.Placeholder = "Notification title"
	titleInput.CharLimit = 100
	titleInput.Width = 40

	v := &NotificationsView{
		Config: cfg,
		VolumeOptions: []VolumeOption{
			{Name: "Normal", Value: 1.0},
			{Name: "Loud", Value: 2.0},
			{Name: "Max", Value: 4.0},
		},
		Categories: []NotificationCategory{
			{Key: "desktop", Label: "Desktop Notifications", Description: "Show system notification popups", Enabled: &cfg.Notifications.Desktop.Enabled},
			{Key: "terminal_bell", Label: "Terminal Bell", Description: "Ring terminal bell on alerts", Enabled: &cfg.Notifications.TerminalBell.Enabled},
			{Key: "blinking_text", Label: "Blinking Text", Description: "Blink statusline text on alerts", Enabled: &cfg.Notifications.BlinkingText.Enabled},
			{Key: "terminal_title", Label: "Terminal Title", Description: "Update terminal title bar", Enabled: &cfg.Notifications.TerminalTitle.Enabled},
			{Key: "tmux", Label: "Tmux Alerts", Description: "Send tmux notifications", Enabled: &cfg.Notifications.Tmux.Enabled},
		},
		Selected:       0,
		ThresholdInput: ti,
		TitleInput:     titleInput,
	}
	v.loadSoundOptions()
	v.selectCurrentSound()
	return v
}

// loadSoundOptions loads available sounds from system and user directories
func (n *NotificationsView) loadSoundOptions() {
	n.SoundOptions = []SoundOption{
		{Name: "(None)", Path: ""},
	}

	// Add macOS system sounds
	if runtime.GOOS == "darwin" {
		systemSoundsDir := "/System/Library/Sounds"
		if entries, err := os.ReadDir(systemSoundsDir); err == nil {
			for _, entry := range entries {
				if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".aiff") {
					name := strings.TrimSuffix(entry.Name(), ".aiff")
					n.SoundOptions = append(n.SoundOptions, SoundOption{
						Name: name + " (System)",
						Path: filepath.Join(systemSoundsDir, entry.Name()),
					})
				}
			}
		}
	}

	// Add user sounds from ~/.claude/sounds/
	homeDir, err := os.UserHomeDir()
	if err == nil {
		userSoundsDir := filepath.Join(homeDir, ".claude", "sounds")
		if entries, err := os.ReadDir(userSoundsDir); err == nil {
			for _, entry := range entries {
				if !entry.IsDir() {
					ext := strings.ToLower(filepath.Ext(entry.Name()))
					if ext == ".aiff" || ext == ".mp3" || ext == ".wav" || ext == ".m4a" {
						name := strings.TrimSuffix(entry.Name(), ext)
						n.SoundOptions = append(n.SoundOptions, SoundOption{
							Name: name + " (Custom)",
							Path: filepath.Join(userSoundsDir, entry.Name()),
						})
					}
				}
			}
		}
	}
}

// selectCurrentSound finds and selects the current sound in the list
func (n *NotificationsView) selectCurrentSound() {
	currentPath := n.Config.Notifications.Desktop.SoundPath
	for i, opt := range n.SoundOptions {
		if opt.Path == currentPath {
			n.SoundSelected = i
			return
		}
	}
	n.SoundSelected = 0
}

// selectCurrentVolume finds and selects the current volume in the list
func (n *NotificationsView) selectCurrentVolume() {
	currentVol := n.Config.Notifications.Desktop.SoundVolume
	for i, opt := range n.VolumeOptions {
		if opt.Value == currentVol {
			n.VolumeSelected = i
			return
		}
	}
	n.VolumeSelected = 0
}

// Up moves selection up
func (n *NotificationsView) Up() {
	if n.SelectingSound {
		n.SoundSelected--
		if n.SoundSelected < 0 {
			n.SoundSelected = len(n.SoundOptions) - 1
		}
		return
	}
	if n.SelectingVolume {
		n.VolumeSelected--
		if n.VolumeSelected < 0 {
			n.VolumeSelected = len(n.VolumeOptions) - 1
		}
		return
	}
	if n.InCategory {
		n.SubSelected--
		maxSub := n.getMaxSubItems()
		if n.SubSelected < 0 {
			n.SubSelected = maxSub - 1
		}
		return
	}
	n.Selected--
	if n.Selected < 0 {
		n.Selected = len(n.Categories) - 1
	}
}

// Down moves selection down
func (n *NotificationsView) Down() {
	if n.SelectingSound {
		n.SoundSelected++
		if n.SoundSelected >= len(n.SoundOptions) {
			n.SoundSelected = 0
		}
		return
	}
	if n.SelectingVolume {
		n.VolumeSelected++
		if n.VolumeSelected >= len(n.VolumeOptions) {
			n.VolumeSelected = 0
		}
		return
	}
	if n.InCategory {
		n.SubSelected++
		maxSub := n.getMaxSubItems()
		if n.SubSelected >= maxSub {
			n.SubSelected = 0
		}
		return
	}
	n.Selected++
	if n.Selected >= len(n.Categories) {
		n.Selected = 0
	}
}

// getMaxSubItems returns the number of sub-items for the current category
func (n *NotificationsView) getMaxSubItems() int {
	cat := n.Categories[n.Selected]
	switch cat.Key {
	case "desktop":
		return 7 // enabled, on_context_panic, threshold, title, sound, sound_path, sound_volume
	case "terminal_title":
		return 5 // enabled, show_model, show_context, alert_on_panic, threshold
	default:
		return 3 // enabled, on_context_panic, threshold
	}
}

// Enter handles enter key
func (n *NotificationsView) Enter() {
	if n.SelectingSound {
		// Confirm sound selection
		n.Config.Notifications.Desktop.SoundPath = n.SoundOptions[n.SoundSelected].Path
		n.SelectingSound = false
		return
	}

	if n.SelectingVolume {
		// Confirm volume selection
		n.Config.Notifications.Desktop.SoundVolume = n.VolumeOptions[n.VolumeSelected].Value
		n.SelectingVolume = false
		return
	}

	if n.EditingThreshold {
		// Save threshold
		var threshold int
		if _, err := parseThreshold(n.ThresholdInput.Value()); err == nil {
			threshold, _ = parseThreshold(n.ThresholdInput.Value())
		}
		n.setThreshold(threshold)
		n.EditingThreshold = false
		return
	}

	if n.EditingTitle {
		n.Config.Notifications.Desktop.Title = n.TitleInput.Value()
		n.EditingTitle = false
		return
	}


	if !n.InCategory {
		n.InCategory = true
		n.SubSelected = 0
		return
	}

	// Handle sub-item actions
	cat := n.Categories[n.Selected]
	switch cat.Key {
	case "desktop":
		n.handleDesktopAction()
	case "terminal_bell", "blinking_text", "tmux":
		n.handleBasicAction()
	case "terminal_title":
		n.handleTerminalTitleAction()
	}
}

func (n *NotificationsView) handleDesktopAction() {
	switch n.SubSelected {
	case 0: // enabled
		n.Config.Notifications.Desktop.Enabled = !n.Config.Notifications.Desktop.Enabled
	case 1: // on_context_panic
		n.Config.Notifications.Desktop.OnContextPanic = !n.Config.Notifications.Desktop.OnContextPanic
	case 2: // threshold
		n.ThresholdInput.SetValue(intToStr(n.Config.Notifications.Desktop.ContextThreshold))
		n.ThresholdInput.Focus()
		n.EditingThreshold = true
	case 3: // title
		n.TitleInput.SetValue(n.Config.Notifications.Desktop.Title)
		n.TitleInput.Focus()
		n.EditingTitle = true
	case 4: // sound enabled
		n.Config.Notifications.Desktop.Sound = !n.Config.Notifications.Desktop.Sound
	case 5: // sound selection
		n.SelectingSound = true
		n.selectCurrentSound()
	case 6: // sound volume
		n.SelectingVolume = true
		n.selectCurrentVolume()
	}
}

func (n *NotificationsView) handleBasicAction() {
	cat := n.Categories[n.Selected]
	switch n.SubSelected {
	case 0: // enabled
		*cat.Enabled = !*cat.Enabled
	case 1: // on_context_panic
		switch cat.Key {
		case "terminal_bell":
			n.Config.Notifications.TerminalBell.OnContextPanic = !n.Config.Notifications.TerminalBell.OnContextPanic
		case "blinking_text":
			n.Config.Notifications.BlinkingText.OnContextPanic = !n.Config.Notifications.BlinkingText.OnContextPanic
		case "tmux":
			n.Config.Notifications.Tmux.OnContextPanic = !n.Config.Notifications.Tmux.OnContextPanic
		}
	case 2: // threshold
		n.ThresholdInput.SetValue(intToStr(n.getThreshold()))
		n.ThresholdInput.Focus()
		n.EditingThreshold = true
	}
}

func (n *NotificationsView) handleTerminalTitleAction() {
	switch n.SubSelected {
	case 0: // enabled
		n.Config.Notifications.TerminalTitle.Enabled = !n.Config.Notifications.TerminalTitle.Enabled
	case 1: // show_model
		n.Config.Notifications.TerminalTitle.ShowModel = !n.Config.Notifications.TerminalTitle.ShowModel
	case 2: // show_context
		n.Config.Notifications.TerminalTitle.ShowContext = !n.Config.Notifications.TerminalTitle.ShowContext
	case 3: // alert_on_panic
		n.Config.Notifications.TerminalTitle.AlertOnPanic = !n.Config.Notifications.TerminalTitle.AlertOnPanic
	case 4: // threshold
		n.ThresholdInput.SetValue(intToStr(n.Config.Notifications.TerminalTitle.ContextThreshold))
		n.ThresholdInput.Focus()
		n.EditingThreshold = true
	}
}

func (n *NotificationsView) getThreshold() int {
	cat := n.Categories[n.Selected]
	switch cat.Key {
	case "desktop":
		return n.Config.Notifications.Desktop.ContextThreshold
	case "terminal_bell":
		return n.Config.Notifications.TerminalBell.ContextThreshold
	case "blinking_text":
		return n.Config.Notifications.BlinkingText.ContextThreshold
	case "terminal_title":
		return n.Config.Notifications.TerminalTitle.ContextThreshold
	case "tmux":
		return n.Config.Notifications.Tmux.ContextThreshold
	}
	return 0
}

func (n *NotificationsView) setThreshold(val int) {
	cat := n.Categories[n.Selected]
	switch cat.Key {
	case "desktop":
		n.Config.Notifications.Desktop.ContextThreshold = val
	case "terminal_bell":
		n.Config.Notifications.TerminalBell.ContextThreshold = val
	case "blinking_text":
		n.Config.Notifications.BlinkingText.ContextThreshold = val
	case "terminal_title":
		n.Config.Notifications.TerminalTitle.ContextThreshold = val
	case "tmux":
		n.Config.Notifications.Tmux.ContextThreshold = val
	}
}

// Back returns true if we should go back to menu
func (n *NotificationsView) Back() bool {
	if n.SelectingSound {
		n.SelectingSound = false
		return false
	}
	if n.EditingThreshold {
		n.EditingThreshold = false
		return false
	}
	if n.EditingTitle {
		n.EditingTitle = false
		return false
	}
	if n.SelectingVolume {
		n.SelectingVolume = false
		return false
	}
	if n.InCategory {
		n.InCategory = false
		return false
	}
	return true
}

// CurrentInput returns the current text input being edited
func (n *NotificationsView) CurrentInput() *textinput.Model {
	if n.EditingThreshold {
		return &n.ThresholdInput
	}
	if n.EditingTitle {
		return &n.TitleInput
	}
	return nil
}

// PlaySelectedSound plays the currently selected sound for preview
func (n *NotificationsView) PlaySelectedSound() {
	if n.SoundSelected > 0 && n.SoundSelected < len(n.SoundOptions) {
		path := n.SoundOptions[n.SoundSelected].Path
		if path != "" && runtime.GOOS == "darwin" {
			exec.Command("afplay", path).Start()
		}
	}
}

// Render returns the notifications view string
func (n *NotificationsView) Render() string {
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

	highlightStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#F59E0B"))

	b.WriteString(titleStyle.Render("Notification Settings"))
	b.WriteString("\n\n")

	if n.SelectingSound {
		return n.renderSoundSelector(&b, selectedStyle, normalStyle, highlightStyle)
	}

	if n.SelectingVolume {
		return n.renderVolumeSelector(&b, selectedStyle, normalStyle, highlightStyle)
	}

	if !n.InCategory {
		return n.renderCategoryList(&b, selectedStyle, normalStyle, checkStyle, uncheckStyle, descStyle)
	}

	return n.renderCategoryDetail(&b, selectedStyle, normalStyle, checkStyle, uncheckStyle, descStyle, highlightStyle)
}

func (n *NotificationsView) renderCategoryList(b *strings.Builder, selectedStyle, normalStyle, checkStyle, uncheckStyle, descStyle lipgloss.Style) string {
	for i, cat := range n.Categories {
		var checkbox string
		if *cat.Enabled {
			checkbox = checkStyle.Render("[x]")
		} else {
			checkbox = uncheckStyle.Render("[ ]")
		}

		var label string
		if i == n.Selected {
			label = selectedStyle.Render(cat.Label)
		} else {
			label = normalStyle.Render(cat.Label)
		}

		b.WriteString("  " + checkbox + " " + label)
		if i == n.Selected {
			b.WriteString("\n")
			b.WriteString(descStyle.Render("      " + cat.Description))
		}
		b.WriteString("\n")
	}

	b.WriteString("\n")
	b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("#6B7280")).Render(
		"  [enter] Configure  [space] Toggle  [esc] Back"))

	return b.String()
}

func (n *NotificationsView) renderCategoryDetail(b *strings.Builder, selectedStyle, normalStyle, checkStyle, uncheckStyle, descStyle, highlightStyle lipgloss.Style) string {
	cat := n.Categories[n.Selected]

	subTitleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#3B82F6")).
		Bold(true)

	b.WriteString(subTitleStyle.Render(cat.Label + " Settings"))
	b.WriteString("\n\n")

	switch cat.Key {
	case "desktop":
		n.renderDesktopSettings(b, selectedStyle, normalStyle, checkStyle, uncheckStyle, highlightStyle)
	case "terminal_title":
		n.renderTerminalTitleSettings(b, selectedStyle, normalStyle, checkStyle, uncheckStyle, highlightStyle)
	default:
		n.renderBasicSettings(b, selectedStyle, normalStyle, checkStyle, uncheckStyle, highlightStyle)
	}

	b.WriteString("\n")
	if n.EditingThreshold {
		b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("#6B7280")).Render(
			"  [enter] Save  [esc] Cancel"))
	} else if n.EditingTitle {
		b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("#6B7280")).Render(
			"  [enter] Save  [esc] Cancel"))
	} else {
		b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("#6B7280")).Render(
			"  [enter/space] Toggle  [esc] Back"))
	}

	return b.String()
}

func (n *NotificationsView) renderDesktopSettings(b *strings.Builder, selectedStyle, normalStyle, checkStyle, uncheckStyle, highlightStyle lipgloss.Style) {
	items := []struct {
		label   string
		enabled bool
		value   string
	}{
		{"Enabled", n.Config.Notifications.Desktop.Enabled, ""},
		{"Trigger on Context Panic", n.Config.Notifications.Desktop.OnContextPanic, ""},
		{"Context Threshold", false, intToStr(n.Config.Notifications.Desktop.ContextThreshold) + "%"},
		{"Notification Title", false, truncateStr(n.Config.Notifications.Desktop.Title, 30)},
		{"Sound Enabled", n.Config.Notifications.Desktop.Sound, ""},
		{"Sound File", false, n.getCurrentSoundName()},
		{"Sound Volume", false, n.getCurrentVolumeName()},
	}

	for i, item := range items {
		var line string
		style := normalStyle
		if i == n.SubSelected {
			style = selectedStyle
		}

		if item.value != "" {
			if i == 2 && n.SubSelected == 2 && n.EditingThreshold {
				line = style.Render("    " + item.label + ": ") + n.ThresholdInput.View()
			} else if i == 3 && n.SubSelected == 3 && n.EditingTitle {
				line = style.Render("    " + item.label + ": ") + n.TitleInput.View()
			} else {
				line = style.Render("    "+item.label+": ") + highlightStyle.Render(item.value)
			}
		} else {
			var checkbox string
			if item.enabled {
				checkbox = checkStyle.Render("[x]")
			} else {
				checkbox = uncheckStyle.Render("[ ]")
			}
			line = "  " + checkbox + " " + style.Render(item.label)
		}
		b.WriteString(line + "\n")
	}
}

func (n *NotificationsView) renderBasicSettings(b *strings.Builder, selectedStyle, normalStyle, checkStyle, uncheckStyle, highlightStyle lipgloss.Style) {
	cat := n.Categories[n.Selected]
	var enabled, onPanic bool
	var threshold int

	switch cat.Key {
	case "terminal_bell":
		enabled = n.Config.Notifications.TerminalBell.Enabled
		onPanic = n.Config.Notifications.TerminalBell.OnContextPanic
		threshold = n.Config.Notifications.TerminalBell.ContextThreshold
	case "blinking_text":
		enabled = n.Config.Notifications.BlinkingText.Enabled
		onPanic = n.Config.Notifications.BlinkingText.OnContextPanic
		threshold = n.Config.Notifications.BlinkingText.ContextThreshold
	case "tmux":
		enabled = n.Config.Notifications.Tmux.Enabled
		onPanic = n.Config.Notifications.Tmux.OnContextPanic
		threshold = n.Config.Notifications.Tmux.ContextThreshold
	}

	items := []struct {
		label   string
		enabled bool
		value   string
	}{
		{"Enabled", enabled, ""},
		{"Trigger on Context Panic", onPanic, ""},
		{"Context Threshold", false, intToStr(threshold) + "%"},
	}

	for i, item := range items {
		var line string
		style := normalStyle
		if i == n.SubSelected {
			style = selectedStyle
		}

		if item.value != "" {
			if i == 2 && n.SubSelected == 2 && n.EditingThreshold {
				line = style.Render("    " + item.label + ": ") + n.ThresholdInput.View()
			} else {
				line = style.Render("    "+item.label+": ") + highlightStyle.Render(item.value)
			}
		} else {
			var checkbox string
			if item.enabled {
				checkbox = checkStyle.Render("[x]")
			} else {
				checkbox = uncheckStyle.Render("[ ]")
			}
			line = "  " + checkbox + " " + style.Render(item.label)
		}
		b.WriteString(line + "\n")
	}
}

func (n *NotificationsView) renderTerminalTitleSettings(b *strings.Builder, selectedStyle, normalStyle, checkStyle, uncheckStyle, highlightStyle lipgloss.Style) {
	cfg := n.Config.Notifications.TerminalTitle
	items := []struct {
		label   string
		enabled bool
		value   string
	}{
		{"Enabled", cfg.Enabled, ""},
		{"Show Model", cfg.ShowModel, ""},
		{"Show Context", cfg.ShowContext, ""},
		{"Alert on Panic", cfg.AlertOnPanic, ""},
		{"Context Threshold", false, intToStr(cfg.ContextThreshold) + "%"},
	}

	for i, item := range items {
		var line string
		style := normalStyle
		if i == n.SubSelected {
			style = selectedStyle
		}

		if item.value != "" {
			if i == 4 && n.SubSelected == 4 && n.EditingThreshold {
				line = style.Render("    " + item.label + ": ") + n.ThresholdInput.View()
			} else {
				line = style.Render("    "+item.label+": ") + highlightStyle.Render(item.value)
			}
		} else {
			var checkbox string
			if item.enabled {
				checkbox = checkStyle.Render("[x]")
			} else {
				checkbox = uncheckStyle.Render("[ ]")
			}
			line = "  " + checkbox + " " + style.Render(item.label)
		}
		b.WriteString(line + "\n")
	}
}

func (n *NotificationsView) renderSoundSelector(b *strings.Builder, selectedStyle, normalStyle, highlightStyle lipgloss.Style) string {
	subTitleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#3B82F6")).
		Bold(true)

	b.WriteString(subTitleStyle.Render("Select Notification Sound"))
	b.WriteString("\n\n")

	// Show a scrollable list (max 10 visible)
	startIdx := 0
	visibleCount := 10
	if n.SoundSelected >= visibleCount {
		startIdx = n.SoundSelected - visibleCount + 1
	}
	endIdx := startIdx + visibleCount
	if endIdx > len(n.SoundOptions) {
		endIdx = len(n.SoundOptions)
	}

	for i := startIdx; i < endIdx; i++ {
		opt := n.SoundOptions[i]
		style := normalStyle
		prefix := "  "
		if i == n.SoundSelected {
			style = selectedStyle
			prefix = "> "
		}

		// Mark current selection
		marker := ""
		if opt.Path == n.Config.Notifications.Desktop.SoundPath {
			marker = highlightStyle.Render(" (current)")
		}

		b.WriteString(prefix + style.Render(opt.Name) + marker + "\n")
	}

	if len(n.SoundOptions) > visibleCount {
		b.WriteString("\n")
		b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("#6B7280")).Italic(true).Render(
			"  Showing " + intToStr(startIdx+1) + "-" + intToStr(endIdx) + " of " + intToStr(len(n.SoundOptions))))
	}

	b.WriteString("\n\n")
	b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("#6B7280")).Render(
		"  [enter] Select  [p] Preview  [esc] Cancel"))

	return b.String()
}

func (n *NotificationsView) getCurrentSoundName() string {
	path := n.Config.Notifications.Desktop.SoundPath
	if path == "" {
		return "(None)"
	}
	for _, opt := range n.SoundOptions {
		if opt.Path == path {
			return opt.Name
		}
	}
	return filepath.Base(path)
}

func (n *NotificationsView) getCurrentVolumeName() string {
	vol := n.Config.Notifications.Desktop.SoundVolume
	for _, opt := range n.VolumeOptions {
		if opt.Value == vol {
			return opt.Name
		}
	}
	return floatToStr(vol) + "x"
}

func (n *NotificationsView) renderVolumeSelector(b *strings.Builder, selectedStyle, normalStyle, highlightStyle lipgloss.Style) string {
	subTitleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#3B82F6")).
		Bold(true)

	b.WriteString(subTitleStyle.Render("Select Volume Level"))
	b.WriteString("\n\n")

	for i, opt := range n.VolumeOptions {
		style := normalStyle
		prefix := "  "
		if i == n.VolumeSelected {
			style = selectedStyle
			prefix = "> "
		}

		// Mark current selection
		marker := ""
		if opt.Value == n.Config.Notifications.Desktop.SoundVolume {
			marker = highlightStyle.Render(" (current)")
		}

		label := opt.Name + " (" + floatToStr(opt.Value) + "x)"
		b.WriteString(prefix + style.Render(label) + marker + "\n")
	}

	b.WriteString("\n")
	b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("#6B7280")).Render(
		"  [enter] Select  [esc] Cancel"))

	return b.String()
}

func intToStr(i int) string {
	return strconv.Itoa(i)
}

func parseThreshold(s string) (int, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0, nil
	}
	val, err := strconv.Atoi(s)
	if err != nil {
		return 0, err
	}
	if val > 100 {
		val = 100
	}
	if val < 0 {
		val = 0
	}
	return val, nil
}

func truncateStr(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max-3] + "..."
}

func floatToStr(f float64) string {
	return strconv.FormatFloat(f, 'f', 1, 64)
}

func parseVolume(s string) (float64, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return 1.0, nil
	}
	val, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 1.0, err
	}
	if val < 0.1 {
		val = 0.1
	}
	if val > 10.0 {
		val = 10.0
	}
	return val, nil
}
