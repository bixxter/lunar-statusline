#!/bin/bash
# Hook script: Clears waiting state when Claude resumes work
# Triggered by: UserPromptSubmit, PostToolUse, SessionStart

STATE_FILE="$HOME/.claude/.statusline-state.json"

# Remove the state file if it exists
rm -f "$STATE_FILE"

exit 0
