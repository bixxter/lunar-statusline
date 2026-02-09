#!/bin/bash
# Install Lunar statusline for Claude Code
# Copies statusline + hook scripts and merges configuration into settings.json

set -e

# --- Quiet mode (used by auto-updater) ---
QUIET=false
for arg in "$@"; do
    case "$arg" in
        --quiet|-q) QUIET=true ;;
    esac
done

say() { [ "$QUIET" = false ] && echo "$@" || true; }

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
CLAUDE_DIR="$HOME/.claude"
HOOKS_DIR="$CLAUDE_DIR/hooks"
SETTINGS_FILE="$CLAUDE_DIR/settings.json"

# --- Dependency check ---
if ! command -v jq &>/dev/null; then
    say "Error: jq is required but not installed."
    say ""
    say "Install it with:"
    say "  macOS:   brew install jq"
    say "  Ubuntu:  sudo apt install jq"
    say "  Arch:    sudo pacman -S jq"
    exit 1
fi

say "Installing Lunar statusline for Claude Code..."
say ""

# --- Create directories ---
mkdir -p "$HOOKS_DIR"

# --- Copy statusline script ---
cp "$SCRIPT_DIR/statusline.sh" "$CLAUDE_DIR/statusline.sh"
chmod +x "$CLAUDE_DIR/statusline.sh"
say "  Installed statusline.sh to $CLAUDE_DIR/"

# --- Copy hook scripts ---
cp "$SCRIPT_DIR/hooks/set-waiting.sh" "$HOOKS_DIR/"
cp "$SCRIPT_DIR/hooks/clear-waiting.sh" "$HOOKS_DIR/"
cp "$SCRIPT_DIR/hooks/update.sh" "$HOOKS_DIR/"
chmod +x "$HOOKS_DIR/set-waiting.sh" "$HOOKS_DIR/clear-waiting.sh" "$HOOKS_DIR/update.sh"
say "  Installed hook scripts to $HOOKS_DIR/"

# --- Read existing settings or start fresh ---
if [ -f "$SETTINGS_FILE" ]; then
    if ! jq empty "$SETTINGS_FILE" 2>/dev/null; then
        say ""
        say "Error: $SETTINGS_FILE contains invalid JSON."
        say "Please fix it manually and re-run this script."
        exit 1
    fi
    SETTINGS=$(cat "$SETTINGS_FILE")
else
    SETTINGS='{}'
fi

# --- Read our hook definitions ---
NEW_HOOKS=$(cat "$SCRIPT_DIR/hooks/hooks.json")

# --- Merge settings: add statusLine config + merge hooks ---
# Uses --argjson to pass hooks.json as a variable (avoids the -s slurp bug).
# For each hook event key, removes any pre-existing entries that reference
# our hook commands (safe re-install), then appends our entries.
SETTINGS=$(echo "$SETTINGS" | jq --argjson new_hooks "$NEW_HOOKS" '
  # Add statusLine if not already configured
  (if .statusLine then .statusLine else {type: "command", command: "~/.claude/statusline.sh", padding: 0} end) as $sl |

  # Merge hooks
  ((.hooks // {}) as $existing |
    reduce ($new_hooks.hooks | keys[]) as $key (
      $existing;
      # Remove our old entries (by command path), then append fresh ones
      .[$key] = ([(.[$key] // [])[] | select(
        (.hooks // []) | map(.command) | any(test("(set-waiting|clear-waiting|update)\\.sh")) | not
      )] + $new_hooks.hooks[$key])
    )
  ) as $merged_hooks |

  . + {statusLine: $sl, hooks: $merged_hooks}
')

echo "$SETTINGS" | jq '.' > "${SETTINGS_FILE}.tmp"
mv "${SETTINGS_FILE}.tmp" "$SETTINGS_FILE"
say "  Updated $SETTINGS_FILE"

say ""
say "Done! Restart Claude Code to activate the statusline."
say ""
say "To customize notifications and display, edit:"
say "  ~/.claude/.statusline.config"
