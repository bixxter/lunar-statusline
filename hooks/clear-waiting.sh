#!/bin/bash
# Hook script: Clears waiting state when Claude resumes work
# Triggered by: UserPromptSubmit, PostToolUse, SessionStart

STATE_FILE="$HOME/.claude/.statusline-state.json"
CONFIG_FILE="$HOME/.claude/.statusline.config"

# Check if we were in waiting state before clearing
WAS_WAITING=false
if [ -f "$STATE_FILE" ]; then
    EXISTING=$(jq -r '.waiting // false' "$STATE_FILE" 2>/dev/null)
    if [ "$EXISTING" = "true" ]; then
        WAS_WAITING=true
    fi
fi

# Remove the state file if it exists
rm -f "$STATE_FILE"

# Reset terminal title if it was set
if [ "$WAS_WAITING" = "true" ]; then
    NOTIFY_TERMINAL_TITLE=$(jq -r '.notifications.terminal_title.enabled // false' "$CONFIG_FILE" 2>/dev/null)
    NOTIFY_TERMINAL_TITLE_NORMAL=$(jq -r '.notifications.terminal_title.normal_text // "Claude Code"' "$CONFIG_FILE" 2>/dev/null)

    if [ "$NOTIFY_TERMINAL_TITLE" = "true" ]; then
        printf '\033]0;%s\007' "$NOTIFY_TERMINAL_TITLE_NORMAL" > /dev/tty 2>/dev/null || true
    fi
fi

exit 0
