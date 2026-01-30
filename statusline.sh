#!/bin/bash
# Claude Code Status Line
# Displays: git branch + status | directory | model | context moons | reactive mascot
# Configuration is read from ~/.claude/.statusline.config

input=$(cat)
echo "$input" >> /tmp/statusline-debug.json

# === Load Config ===
CONFIG_FILE="$HOME/.claude/.statusline.config"

# Helper function to read config values with defaults
cfg() {
    local path="$1"
    local default="$2"
    if [ -f "$CONFIG_FILE" ]; then
        local value
        value=$(jq -r "$path // empty" "$CONFIG_FILE" 2>/dev/null)
        if [ -n "$value" ] && [ "$value" != "null" ]; then
            echo "$value"
            return
        fi
    fi
    echo "$default"
}

cfg_bool() {
    local path="$1"
    local default="$2"
    if [ -f "$CONFIG_FILE" ]; then
        local value
        value=$(jq -r "$path // empty" "$CONFIG_FILE" 2>/dev/null)
        if [ "$value" = "true" ]; then
            echo "true"
            return
        elif [ "$value" = "false" ]; then
            echo "false"
            return
        fi
    fi
    echo "$default"
}

cfg_array() {
    local path="$1"
    local index="$2"
    local default="$3"
    if [ -f "$CONFIG_FILE" ]; then
        local value
        value=$(jq -r "${path}[${index}] // empty" "$CONFIG_FILE" 2>/dev/null)
        if [ -n "$value" ] && [ "$value" != "null" ]; then
            echo "$value"
            return
        fi
    fi
    echo "$default"
}

# === Read enabled sections ===
SHOW_GIT=$(cfg_bool '.enabled_sections.git' 'true')
SHOW_DIR=$(cfg_bool '.enabled_sections.directory' 'true')
SHOW_MODEL=$(cfg_bool '.enabled_sections.model' 'true')
SHOW_MOONS=$(cfg_bool '.enabled_sections.context_moons' 'true')
SHOW_TOKENS=$(cfg_bool '.enabled_sections.token_count' 'true')
SHOW_PERCENT=$(cfg_bool '.enabled_sections.percentage' 'true')
SHOW_MASCOT=$(cfg_bool '.enabled_sections.mascot' 'true')

# === Read icons ===
ICON_GIT_CLEAN=$(cfg '.icons.git_clean' '✅')
ICON_GIT_DIRTY=$(cfg '.icons.git_dirty' '⚠️')
ICON_DIR=$(cfg '.icons.directory' '🗂️')

# Moon phases from config
MOON_1=$(cfg_array '.icons.moons' 0 '●')
MOON_2=$(cfg_array '.icons.moons' 1 '◐')
MOON_3=$(cfg_array '.icons.moons' 2 '◑')
MOON_4=$(cfg_array '.icons.moons' 3 '◕')
MOON_5=$(cfg_array '.icons.moons' 4 '○')

# === Read display settings ===
SEPARATOR=$(cfg '.display.separator' ' • ')
DIR_MAX_LEN=$(cfg '.thresholds.directory_max_length' '15')
DIR_TRUNCATE=$(cfg '.thresholds.directory_truncate_to' '12')

# === Read mascot settings ===
MASCOT_PANIC_ENABLED=$(cfg_bool '.mascot.context_panic.enabled' 'true')
MASCOT_PANIC_THRESHOLD=$(cfg '.mascot.context_panic.threshold' '90')
MASCOT_PROD_ENABLED=$(cfg_bool '.mascot.productive.enabled' 'true')
MASCOT_PROD_THRESHOLD=$(cfg '.mascot.productive.threshold' '100')
MASCOT_DEL_ENABLED=$(cfg_bool '.mascot.deletion.enabled' 'true')
MASCOT_DEL_THRESHOLD=$(cfg '.mascot.deletion.threshold' '30')
MASCOT_TIME_ENABLED=$(cfg_bool '.mascot.time_based.enabled' 'true')

# === Directory Info ===
DIR_INFO=""
if [ "$SHOW_DIR" = "true" ]; then
    DIR_NAME=$(basename "$PWD")
    # Truncate if longer than max length
    if [ ${#DIR_NAME} -gt "$DIR_MAX_LEN" ]; then
        DIR_INFO="${DIR_NAME:0:$DIR_TRUNCATE}..."
    else
        DIR_INFO="$DIR_NAME"
    fi
    DIR_INFO="\033[35m$ICON_DIR $DIR_INFO\033[0m"  # Magenta
fi

# === Git Info ===
GIT_INFO=""
if [ "$SHOW_GIT" = "true" ] && git rev-parse --git-dir > /dev/null 2>&1; then
    BRANCH=$(git branch --show-current 2>/dev/null)
    if [ -n "$BRANCH" ]; then
        # Check for uncommitted changes
        if git diff --quiet 2>/dev/null && git diff --cached --quiet 2>/dev/null; then
            GIT_INFO="\033[32m$ICON_GIT_CLEAN $BRANCH\033[0m"  # Green - all committed
        else
            GIT_INFO="\033[31m$ICON_GIT_DIRTY $BRANCH\033[0m"  # Red - uncommitted changes
        fi
    fi
fi

# === Model ===
MODEL_INFO=""
if [ "$SHOW_MODEL" = "true" ]; then
    MODEL=$(echo "$input" | jq -r '.model.display_name // "?"')
    MODEL_INFO="\033[36m$MODEL\033[0m"
fi

# === Context Moons ===
PERCENT=$(echo "$input" | jq -r '.context_window.used_percentage // 0' | cut -d. -f1)

get_moon() {
    local pct=$1
    if [ "$pct" -lt 15 ]; then
        echo "$MOON_1"
    elif [ "$pct" -lt 40 ]; then
        echo "$MOON_2"
    elif [ "$pct" -lt 60 ]; then
        echo "$MOON_3"
    elif [ "$pct" -lt 85 ]; then
        echo "$MOON_4"
    else
        echo "$MOON_5"
    fi
}

CONTEXT_INFO=""
if [ "$SHOW_MOONS" = "true" ]; then
    # Split into thirds for visualization
    THIRD1=$((PERCENT * 3))
    THIRD2=$(((PERCENT - 33) * 3))
    THIRD3=$(((PERCENT - 66) * 3))
    [ "$THIRD1" -lt 0 ] && THIRD1=0
    [ "$THIRD2" -lt 0 ] && THIRD2=0
    [ "$THIRD3" -lt 0 ] && THIRD3=0
    [ "$THIRD1" -gt 100 ] && THIRD1=100
    [ "$THIRD2" -gt 100 ] && THIRD2=100
    [ "$THIRD3" -gt 100 ] && THIRD3=100

    MOON1=$(get_moon $THIRD1)
    MOON2=$(get_moon $THIRD2)
    MOON3=$(get_moon $THIRD3)

    CONTEXT_INFO="${MOON1}${MOON2}${MOON3}"
fi

# Token count (k format)
if [ "$SHOW_TOKENS" = "true" ]; then
    TOKENS=$(echo "$input" | jq -r '.context_window.total_input_tokens // 0')
    TOKEN_K_FORMAT=$(cfg '.thresholds.token_k_format' '1000')
    if [ "$TOKENS" -gt "$TOKEN_K_FORMAT" ]; then
        TOKENS_DISPLAY="$((TOKENS / 1000))k"
    else
        TOKENS_DISPLAY="$TOKENS"
    fi
    if [ -n "$CONTEXT_INFO" ]; then
        CONTEXT_INFO="$CONTEXT_INFO $TOKENS_DISPLAY"
    else
        CONTEXT_INFO="$TOKENS_DISPLAY"
    fi
fi

# Percentage
if [ "$SHOW_PERCENT" = "true" ]; then
    if [ -n "$CONTEXT_INFO" ]; then
        CONTEXT_INFO="$CONTEXT_INFO (${PERCENT}%)"
    else
        CONTEXT_INFO="(${PERCENT}%)"
    fi
fi

# === Reactive Mascot ===
MASCOT=""
if [ "$SHOW_MASCOT" = "true" ]; then
    LINES_ADDED=$(echo "$input" | jq -r '.cost.total_lines_added // 0')
    LINES_REMOVED=$(echo "$input" | jq -r '.cost.total_lines_removed // 0')

    get_mascot() {
        # Random factor for variety (changes every ~10 seconds)
        RANDOM_SEED=$(($(date +%s) / 10))

        # Context panic mode
        if [ "$MASCOT_PANIC_ENABLED" = "true" ] && [ "$PERCENT" -gt "$MASCOT_PANIC_THRESHOLD" ]; then
            # Get emojis from config
            PANIC_COUNT=$(jq -r '.mascot.context_panic.emojis | length' "$CONFIG_FILE" 2>/dev/null)
            if [ -n "$PANIC_COUNT" ] && [ "$PANIC_COUNT" -gt 0 ]; then
                IDX=$((RANDOM_SEED % PANIC_COUNT))
                jq -r ".mascot.context_panic.emojis[$IDX]" "$CONFIG_FILE" 2>/dev/null
            else
                case $((RANDOM_SEED % 3)) in
                    0) echo "🫠 melting..." ;;
                    1) echo "😰 tight fit!" ;;
                    2) echo "🔥 toasty!" ;;
                esac
            fi
            return
        fi

        # Productive mode (lots of lines added)
        if [ "$MASCOT_PROD_ENABLED" = "true" ] && [ "$LINES_ADDED" -gt "$MASCOT_PROD_THRESHOLD" ]; then
            PROD_COUNT=$(jq -r '.mascot.productive.emojis | length' "$CONFIG_FILE" 2>/dev/null)
            if [ -n "$PROD_COUNT" ] && [ "$PROD_COUNT" -gt 0 ]; then
                IDX=$((RANDOM_SEED % PROD_COUNT))
                jq -r ".mascot.productive.emojis[$IDX]" "$CONFIG_FILE" 2>/dev/null
            else
                case $((RANDOM_SEED % 4)) in
                    0) echo "🚀 zooming!" ;;
                    1) echo "⚡ on fire!" ;;
                    2) echo "💪 crushing it" ;;
                    3) echo "🎯 locked in" ;;
                esac
            fi
            return
        fi

        # Deletion mode
        if [ "$MASCOT_DEL_ENABLED" = "true" ] && [ "$LINES_REMOVED" -gt "$LINES_ADDED" ] && [ "$LINES_REMOVED" -gt "$MASCOT_DEL_THRESHOLD" ]; then
            DEL_COUNT=$(jq -r '.mascot.deletion.emojis | length' "$CONFIG_FILE" 2>/dev/null)
            if [ -n "$DEL_COUNT" ] && [ "$DEL_COUNT" -gt 0 ]; then
                IDX=$((RANDOM_SEED % DEL_COUNT))
                jq -r ".mascot.deletion.emojis[$IDX]" "$CONFIG_FILE" 2>/dev/null
            else
                case $((RANDOM_SEED % 3)) in
                    0) echo "🧹 cleaning!" ;;
                    1) echo "✂️ snip snip" ;;
                    2) echo "🗑️ declutter" ;;
                esac
            fi
            return
        fi

        # Time-based moods (default)
        if [ "$MASCOT_TIME_ENABLED" = "true" ]; then
            HOUR=$(date +%H)
            if [ "$HOUR" -lt 6 ]; then
                TIME_KEY="night"
            elif [ "$HOUR" -lt 12 ]; then
                TIME_KEY="morning"
            elif [ "$HOUR" -lt 18 ]; then
                TIME_KEY="afternoon"
            else
                TIME_KEY="evening"
            fi

            TIME_COUNT=$(jq -r ".mascot.time_based.$TIME_KEY | length" "$CONFIG_FILE" 2>/dev/null)
            if [ -n "$TIME_COUNT" ] && [ "$TIME_COUNT" -gt 0 ]; then
                IDX=$((RANDOM_SEED % TIME_COUNT))
                jq -r ".mascot.time_based.${TIME_KEY}[$IDX]" "$CONFIG_FILE" 2>/dev/null
            else
                # Fallback defaults
                case $TIME_KEY in
                    night) echo "🦉 night owl" ;;
                    morning) echo "🌅 fresh start" ;;
                    afternoon) echo "🎧 in the zone" ;;
                    evening) echo "🌆 evening mode" ;;
                esac
            fi
        else
            echo "🤖 working"
        fi
    }

    MASCOT=$(get_mascot)
fi

# === Compose Status Line ===
# Build parts array
PARTS=()
[ -n "$GIT_INFO" ] && PARTS+=("$GIT_INFO")
[ -n "$DIR_INFO" ] && PARTS+=("$DIR_INFO")
[ -n "$MODEL_INFO" ] && PARTS+=("$MODEL_INFO")
[ -n "$CONTEXT_INFO" ] && PARTS+=("$CONTEXT_INFO")
[ -n "$MASCOT" ] && PARTS+=("$MASCOT")

# Join with separator
OUTPUT=""
for i in "${!PARTS[@]}"; do
    if [ "$i" -gt 0 ]; then
        OUTPUT+=" │ "
    fi
    OUTPUT+="${PARTS[$i]}"
done

echo -e "$OUTPUT"
