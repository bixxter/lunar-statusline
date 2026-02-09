#!/bin/bash
# Auto-update for Lunar statusline
# Triggered by SessionStart hook. Checks at most once per interval,
# then forks the actual update to background so the hook returns instantly.

CLAUDE_DIR="$HOME/.claude"
REPO_DIR="$CLAUDE_DIR/claude-statusline"
CONFIG_FILE="$CLAUDE_DIR/.statusline.config"
THROTTLE_FILE="$CLAUDE_DIR/.statusline-last-update-check"
LOG_FILE="$CLAUDE_DIR/.statusline-update.log"

# --- Read config helpers ---
cfg_val() {
    local path="$1" default="$2"
    if [ -f "$CONFIG_FILE" ]; then
        local v
        v=$(jq -r "$path // empty" "$CONFIG_FILE" 2>/dev/null)
        [ -n "$v" ] && [ "$v" != "null" ] && { echo "$v"; return; }
    fi
    echo "$default"
}

# --- Check if auto-update is enabled (default: true) ---
ENABLED=$(cfg_val '.auto_update.enabled' 'true')
[ "$ENABLED" != "true" ] && exit 0

# --- Throttle: skip if checked recently ---
INTERVAL_HOURS=$(cfg_val '.auto_update.check_interval_hours' '24')
INTERVAL_SECS=$((INTERVAL_HOURS * 3600))

if [ -f "$THROTTLE_FILE" ]; then
    LAST_CHECK=$(cat "$THROTTLE_FILE" 2>/dev/null)
    NOW=$(date +%s)
    ELAPSED=$((NOW - LAST_CHECK))
    [ "$ELAPSED" -lt "$INTERVAL_SECS" ] && exit 0
fi

# Update throttle timestamp immediately (before forking)
date +%s > "$THROTTLE_FILE"

# --- Fork the actual update to background ---
(
    log() { echo "[$(date '+%Y-%m-%d %H:%M:%S')] $*" >> "$LOG_FILE"; }

    # Sanity checks
    if [ ! -d "$REPO_DIR/.git" ]; then
        log "ERROR: $REPO_DIR is not a git repository. Skipping update."
        exit 1
    fi

    cd "$REPO_DIR" || exit 1

    # Determine default branch
    BRANCH=$(git symbolic-ref refs/remotes/origin/HEAD 2>/dev/null | sed 's@^refs/remotes/origin/@@')
    [ -z "$BRANCH" ] && BRANCH="main"

    # Fetch latest
    if ! git fetch --quiet origin "$BRANCH" 2>/dev/null; then
        log "WARN: git fetch failed (network issue?). Will retry next interval."
        exit 1
    fi

    LOCAL=$(git rev-parse HEAD 2>/dev/null)
    REMOTE=$(git rev-parse "origin/$BRANCH" 2>/dev/null)

    if [ "$LOCAL" = "$REMOTE" ]; then
        log "OK: Already up to date (${LOCAL:0:7})."
        exit 0
    fi

    log "Updating: ${LOCAL:0:7} -> ${REMOTE:0:7}"

    # Pull with --ff-only: safe, won't create merges or clobber local changes
    if ! git pull --ff-only --quiet origin "$BRANCH" 2>/dev/null; then
        log "WARN: git pull --ff-only failed. Local changes may be blocking the update."
        log "WARN: Run 'cd $REPO_DIR && git pull' manually to resolve."
        exit 1
    fi

    # Re-run installer quietly to copy updated files
    if [ -x "$REPO_DIR/install-hooks.sh" ]; then
        "$REPO_DIR/install-hooks.sh" --quiet 2>/dev/null
        log "OK: Updated and re-installed successfully (now at ${REMOTE:0:7})."
    else
        log "WARN: install-hooks.sh not found or not executable after pull."
    fi
) &

exit 0
