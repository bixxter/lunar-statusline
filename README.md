```
██╗     ██╗   ██╗███╗   ██╗ █████╗ ██████╗
██║     ██║   ██║████╗  ██║██╔══██╗██╔══██╗
██║     ██║   ██║██╔██╗ ██║███████║██████╔╝
██║     ██║   ██║██║╚██╗██║██╔══██║██╔══██╗
███████╗╚██████╔╝██║ ╚████║██║  ██║██║  ██║
╚══════╝ ╚═════╝ ╚═╝  ╚═══╝╚═╝  ╚═╝╚═╝  ╚═╝
```

A reactive, visual status line for Claude Code that shows what matters.

## What it shows

```
🔔 WAITING (2m) │ 🌱 main │ 📁 coolest-project │ Opus 4.5 │ 🌕🌑🌑 195k (37%) │ 🚀 zooming!
```

- **Waiting indicator**: Alert when Claude needs your input (permission, question, etc.)
- **Git status**: 🌱 clean / 🥀 uncommitted changes
- **Current directory**: Compact folder name
- **Model**: Which Claude you're talking to
- **Context usage**: Moon phases 🌑→🌕 showing how full your context window is
- **Reactive mascot**: Changes based on activity, time of day, and context pressure

## Install

### Quick Install

```bash
# Install statusline + hooks for waiting indicator
./install-hooks.sh
```

This installs:
- The statusline script to `~/.claude/statusline.sh`
- Hooks that detect when Claude is waiting for your input
- Configuration to `~/.claude/settings.json`

### Manual Install

1. Copy `statusline.sh` to your Claude config directory:
```bash
cp statusline.sh ~/.claude/statusline.sh
chmod +x ~/.claude/statusline.sh
```

2. Add to your `~/.claude/settings.json`:
```json
{
  "statusLine": {
    "type": "command",
    "command": "~/.claude/statusline.sh",
    "padding": 0
  }
}
```

3. (Optional) Install hooks for waiting indicator:
```bash
./install-hooks.sh
```

4. Restart Claude Code

## Requirements

- `jq` for JSON parsing: `brew install jq` (macOS) or `apt install jq` (Linux)
- `bc` for math (usually pre-installed)
- Git (optional, for branch display)

## Waiting Indicator

Never miss when Claude needs your attention. The statusline shows a prominent alert when:

- Claude asks a question (AskUserQuestion)
- A permission prompt appears ("Allow Claude to run...?")
- Any dialog requiring user input

The indicator shows how long Claude has been waiting:
```
🔔 WAITING (2m) │ ...rest of status...
```

**How it works**: Uses Claude Code hooks to detect waiting states and writes to a state file that the statusline reads.

**Customize** in `~/.claude/.statusline.config`:
```json
{
  "waiting_indicator": {
    "icon": "🔔",
    "text": "WAITING",
    "blink": true
  },
  "notifications": {
    "enabled": true,
    "terminal_bell": true,
    "system_notification": true,
    "sound": true
  }
}
```

**Notification options:**
- `terminal_bell` - Classic `\a` bell (works in most terminals)
- `system_notification` - Native OS notification (macOS/Linux)
- `sound` - Play a sound (macOS only, uses system Ping sound)

## Mascot moods

The mascot adapts to your session:

- **Context panic** (>80%): 🫠 melting, 😰 tight fit, 🔥 toasty
- **Productive** (>200 lines added): 🚀 zooming, ⚡ on fire, 💪 crushing it, 🎯 locked in
- **Cleanup mode** (more deletions): 🧹 cleaning, ✂️ snip snip, 🗑️ declutter
- **Chill vibes**: Time-of-day themed (🦉 night owl, ☀️ morning, 🎧 in the zone, 🌆 evening)

Rotates every ~10 seconds to stay fresh without being distracting.

## Configuration Editor

A TUI for customizing your statusline without editing files.

![Editor main screen](demo_1.png)

![Customization options](demo_2.png)

```bash
./lunar-editor-macos   # macOS
./lunar-editor-linux   # Linux
./lunar-editor.exe     # Windows
```

Configure sections, icons, mascot moods, and display settings.

---

Built for context awareness and vibes.
