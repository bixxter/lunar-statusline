#!/bin/bash
# Claude Code Status Line
# Displays: git branch + status | directory | model | context moons | reactive mascot

input=$(cat)

# === Directory Info ===
DIR_NAME=$(basename "$PWD")
# Truncate if longer than 20 chars
if [ ${#DIR_NAME} -gt 20 ]; then
    DIR_INFO="${DIR_NAME:0:17}..."
else
    DIR_INFO="$DIR_NAME"
fi
DIR_INFO="\033[35m📁 $DIR_INFO\033[0m"  # Magenta

# === Git Info ===
GIT_INFO=""
if git rev-parse --git-dir > /dev/null 2>&1; then
    BRANCH=$(git branch --show-current 2>/dev/null)
    if [ -n "$BRANCH" ]; then
        # Check for uncommitted changes
        if git diff --quiet 2>/dev/null && git diff --cached --quiet 2>/dev/null; then
            GIT_INFO="\033[32m🌱 $BRANCH\033[0m"  # Green - all committed
        else
            GIT_INFO="\033[31m🥀 $BRANCH\033[0m"  # Red - uncommitted changes
        fi
    fi
fi

# === Model ===
MODEL=$(echo "$input" | jq -r '.model.display_name // "?"')

# === Context Moons ===
PERCENT=$(echo "$input" | jq -r '.context_window.used_percentage // 0' | cut -d. -f1)

# Three moons representing 0-33%, 34-66%, 67-100%
get_moon() {
    local pct=$1
    if [ "$pct" -lt 15 ]; then
        echo "🌑"
    elif [ "$pct" -lt 40 ]; then
        echo "🌘"
    elif [ "$pct" -lt 60 ]; then
        echo "🌗"
    elif [ "$pct" -lt 85 ]; then
        echo "🌖"
    else
        echo "🌕"
    fi
}

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

# Token count (k format)
TOKENS=$(echo "$input" | jq -r '.context_window.total_input_tokens // 0')
if [ "$TOKENS" -gt 1000 ]; then
    TOKENS_DISPLAY="$((TOKENS / 1000))k"
else
    TOKENS_DISPLAY="$TOKENS"
fi

CONTEXT_INFO="${MOON1}${MOON2}${MOON3} ${TOKENS_DISPLAY} (${PERCENT}%)"

# === Reactive Mascot ===
LINES_ADDED=$(echo "$input" | jq -r '.cost.total_lines_added // 0')
LINES_REMOVED=$(echo "$input" | jq -r '.cost.total_lines_removed // 0')

# Mascot states based on session activity
get_mascot() {
    # Random factor for variety (changes every ~10 seconds)
    RANDOM_SEED=$(($(date +%s) / 10))
    
    # Context panic mode
    if [ "$PERCENT" -gt 80 ]; then
        case $((RANDOM_SEED % 3)) in
            0) echo "🫠 melting..." ;;
            1) echo "😰 tight fit!" ;;
            2) echo "🔥 toasty!" ;;
        esac
        return
    fi
    
    # Productive mode (lots of lines added)
    if [ "$LINES_ADDED" -gt 200 ]; then
        case $((RANDOM_SEED % 4)) in
            0) echo "🚀 zooming!" ;;
            1) echo "⚡ on fire!" ;;
            2) echo "💪 crushing it" ;;
            3) echo "🎯 locked in" ;;
        esac
        return
    fi
    
    # Deletion mode
    if [ "$LINES_REMOVED" -gt "$LINES_ADDED" ] && [ "$LINES_REMOVED" -gt 50 ]; then
        case $((RANDOM_SEED % 3)) in
            0) echo "🧹 cleaning!" ;;
            1) echo "✂️ snip snip" ;;
            2) echo "🗑️ declutter" ;;
        esac
        return
    fi
    
    # Default chill vibes
    HOUR=$(date +%H)
    if [ "$HOUR" -lt 6 ]; then
        MOODS=("🦉 night owl" "🌙 late grind" "☕ need coffee")
    elif [ "$HOUR" -lt 12 ]; then
        MOODS=("🌅 fresh start" "☀️ morning!" "🥐 coding time")
    elif [ "$HOUR" -lt 18 ]; then
        MOODS=("🎧 in the zone" "🧠 thinking..." "💭 hmm...")
    else
        MOODS=("🌆 evening mode" "🍕 dinner code" "✨ wrapping up")
    fi
    
    echo "${MOODS[$((RANDOM_SEED % 3))]}"
}

MASCOT=$(get_mascot)

# === Compose Status Line ===
# Format: git | directory | model | context | mascot
if [ -n "$GIT_INFO" ]; then
    echo -e "$GIT_INFO │ $DIR_INFO │ \033[36m$MODEL\033[0m │ $CONTEXT_INFO │ $MASCOT"
else
    echo -e "$DIR_INFO │ \033[36m$MODEL\033[0m │ $CONTEXT_INFO │ $MASCOT"
fi
