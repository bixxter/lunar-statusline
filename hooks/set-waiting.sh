#!/bin/bash
# Hook script: Sets waiting state when Claude needs user input
# Triggered by: Notification (idle_prompt, permission_prompt, elicitation_dialog), PermissionRequest

STATE_FILE="$HOME/.claude/.statusline-state.json"
CONFIG_FILE="$HOME/.claude/.statusline.config"

# Read hook input
INPUT=$(cat)

# Determine the type from the notification or default to "input"
NOTIFICATION_TYPE=$(echo "$INPUT" | jq -r '.notification_type // "input"')
TOOL_NAME=$(echo "$INPUT" | jq -r '.tool_name // empty')
MESSAGE=$(echo "$INPUT" | jq -r '.message // "Claude needs your input"')

# Set the waiting type based on what triggered this
if [ -n "$TOOL_NAME" ]; then
    TYPE="permission:$TOOL_NAME"
    TITLE="Permission Required"
elif [ "$NOTIFICATION_TYPE" = "permission_prompt" ]; then
    TYPE="permission"
    TITLE="Permission Required"
elif [ "$NOTIFICATION_TYPE" = "elicitation_dialog" ]; then
    TYPE="question"
    TITLE="Claude has a question"
elif [ "$NOTIFICATION_TYPE" = "idle_prompt" ]; then
    TYPE="idle"
    TITLE="Claude is waiting"
else
    TYPE="input"
    TITLE="Input Required"
fi

# Check if we already have a waiting state (avoid duplicate notifications)
ALREADY_WAITING=false
if [ -f "$STATE_FILE" ]; then
    EXISTING=$(jq -r '.waiting // false' "$STATE_FILE" 2>/dev/null)
    if [ "$EXISTING" = "true" ]; then
        ALREADY_WAITING=true
    fi
fi

# Write state file
cat > "$STATE_FILE" << EOF
{
  "waiting": true,
  "type": "$TYPE",
  "timestamp": $(date +%s),
  "message": "$MESSAGE"
}
EOF

# Only notify if this is a new waiting state
if [ "$ALREADY_WAITING" = "false" ]; then
    # Read notification settings from config (correct paths)
    NOTIFY_BELL=$(jq -r '.notifications.terminal_bell.enabled // false' "$CONFIG_FILE" 2>/dev/null)
    NOTIFY_DESKTOP=$(jq -r '.notifications.desktop.enabled // false' "$CONFIG_FILE" 2>/dev/null)
    NOTIFY_SOUND=$(jq -r '.notifications.desktop.sound // false' "$CONFIG_FILE" 2>/dev/null)
    NOTIFY_SOUND_PATH=$(jq -r '.notifications.desktop.sound_path // ""' "$CONFIG_FILE" 2>/dev/null)
    NOTIFY_SOUND_VOLUME=$(jq -r '.notifications.desktop.sound_volume // 1' "$CONFIG_FILE" 2>/dev/null)
    NOTIFY_TITLE=$(jq -r '.notifications.desktop.title // "Claude needs attention"' "$CONFIG_FILE" 2>/dev/null)

    # Terminal bell
    if [ "$NOTIFY_BELL" = "true" ]; then
        printf '\a'
    fi

    # System notification (macOS or Linux)
    if [ "$NOTIFY_DESKTOP" = "true" ]; then
        if command -v osascript &> /dev/null; then
            # macOS
            osascript -e "display notification \"$MESSAGE\" with title \"$NOTIFY_TITLE\" sound name \"\"" 2>/dev/null &
        elif command -v notify-send &> /dev/null; then
            # Linux
            notify-send -u critical "$NOTIFY_TITLE" "$MESSAGE" 2>/dev/null &
        fi
    fi

    # Sound (macOS only)
    if [ "$NOTIFY_SOUND" = "true" ]; then
        if command -v afplay &> /dev/null; then
            if [ -n "$NOTIFY_SOUND_PATH" ] && [ -f "$NOTIFY_SOUND_PATH" ]; then
                # Custom sound file with volume
                afplay -v "$NOTIFY_SOUND_VOLUME" "$NOTIFY_SOUND_PATH" 2>/dev/null &
            else
                # Fallback to system sound with volume
                afplay -v "$NOTIFY_SOUND_VOLUME" /System/Library/Sounds/Tink.aiff 2>/dev/null &
            fi
        fi
    fi
fi

exit 0
