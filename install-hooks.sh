#!/bin/bash
# Install Claude Code hooks for waiting indicator
# This script copies hook scripts and merges hook configuration

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
CLAUDE_DIR="$HOME/.claude"
HOOKS_DIR="$CLAUDE_DIR/hooks"
SETTINGS_FILE="$CLAUDE_DIR/settings.json"

echo "Installing Claude Code statusline hooks..."

# Create hooks directory
mkdir -p "$HOOKS_DIR"

# Copy hook scripts
cp "$SCRIPT_DIR/hooks/set-waiting.sh" "$HOOKS_DIR/"
cp "$SCRIPT_DIR/hooks/clear-waiting.sh" "$HOOKS_DIR/"
chmod +x "$HOOKS_DIR/set-waiting.sh" "$HOOKS_DIR/clear-waiting.sh"

echo "  Copied hook scripts to $HOOKS_DIR"

# Merge hooks into settings.json
if [ -f "$SETTINGS_FILE" ]; then
    # Check if hooks already exist
    EXISTING_HOOKS=$(jq -r '.hooks // empty' "$SETTINGS_FILE" 2>/dev/null)
    if [ -n "$EXISTING_HOOKS" ] && [ "$EXISTING_HOOKS" != "null" ]; then
        echo "  Found existing hooks in settings.json"
        # Merge hooks (new hooks take precedence)
        MERGED=$(jq -s '.[0] * {hooks: (.[0].hooks // {} | . * .[1].hooks)}' "$SETTINGS_FILE" "$SCRIPT_DIR/hooks/hooks.json")
        echo "$MERGED" > "$SETTINGS_FILE"
        echo "  Merged hooks configuration"
    else
        # No existing hooks, just add them
        jq -s '.[0] * .[1]' "$SETTINGS_FILE" "$SCRIPT_DIR/hooks/hooks.json" > "${SETTINGS_FILE}.tmp"
        mv "${SETTINGS_FILE}.tmp" "$SETTINGS_FILE"
        echo "  Added hooks to settings.json"
    fi
else
    # No settings file, copy hooks.json as base
    echo '{}' | jq -s '.[0] * .[1]' - "$SCRIPT_DIR/hooks/hooks.json" > "$SETTINGS_FILE"
    echo "  Created settings.json with hooks"
fi

echo ""
echo "Installation complete!"
echo ""
echo "The statusline will now show a waiting indicator when Claude needs your input."
echo "You can customize the indicator in your statusline config:"
echo "  ~/.claude/.statusline.config"
echo ""
echo "Settings:"
echo '  "waiting_indicator": {'
echo '    "enabled": true,'
echo '    "icon": "🔔",'
echo '    "text": "WAITING",'
echo '    "blink": true'
echo '  }'
