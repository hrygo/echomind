#!/bin/bash

# EchoMind CI/CD Monitor - Simplified & Robust
# Usage: ./scripts/ci.sh [watch|history [N]|analyze|interactive|help]

set -euo pipefail

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
GRAY='\033[0;37m'
BOLD='\033[1m'
NC='\033[0m'

# Symbols
SUCCESS='✅'
FAILURE='❌'
WARNING='⚠️'
INFO='ℹ️'
PROGRESS='⏳'
CHART='📊'

# Config with defaults
HISTORY_COUNT=${HISTORY_COUNT:-5}
BRANCH=${BRANCH:-main}
TEMP_DIR=$(mktemp -d)
trap "rm -rf $TEMP_DIR" EXIT

# Logging functions
log_info() { echo -e "${BLUE}${INFO} $1${NC}"; }
log_success() { echo -e "${GREEN}${SUCCESS} $1${NC}"; }
log_warning() { echo -e "${YELLOW}${WARNING} $1${NC}"; }
log_error() { echo -e "${RED}${FAILURE} $1${NC}"; }
log_progress() {
    local msg=$1
    local duration=${2:-2}
    local steps=${3:-20}
    echo -ne "${YELLOW}${PROGRESS} ${msg} ${NC}"
    for i in $(seq 1 $steps); do
        sleep $(echo "scale=2; $duration/$steps" | bc)
        local percent=$((i * 100 / steps))
        local filled=$((percent / 2))
        printf "\r${YELLOW}${PROGRESS} ${msg} ${NC}[${GREEN}%*s${NC}%*s] ${percent}%%" \
            $filled '' $((50 - filled)) ''
    done
    echo
}

# Validate dependencies
check_deps() {
    local missing=()
    command -v gh >/dev/null 2>&1 || missing+=("gh - GitHub CLI (brew install gh)")
    command -v jq >/dev/null 2>&1 || missing+=("jq - JSON processor (brew install jq)")

    if [[ ${#missing[@]} -gt 0 ]]; then
        log_error "Missing required dependencies:"
        printf '  %s\n' "${missing[@]}"
        exit 1
    fi

    if ! gh auth status >/dev/null 2>&1; then
        log_error "GitHub CLI not authenticated. Run: gh auth login"
        exit 1
    fi
}

# Fetch CI data with robust error handling
fetch_data() {
    local limit=${1:-1}
    local branch_filter=${2:-$BRANCH}

    local data
    data=$(gh run list --limit "$limit" --branch "$branch_filter" \
        --json status,conclusion,workflowName,headBranch,databaseId,createdAt,displayTitle \
        --jq '. | sort_by(.databaseId) | reverse' 2>/dev/null) || {
        log_error "Failed to fetch CI data"
        echo "Possible causes:"
        echo "  • No internet connection"
        echo "  • GitHub API rate limit exceeded"
        echo "  • No workflow runs found"
        echo "  • Invalid repository access"
        return 1
    }

    [[ -n "$data" && "$data" != "null" ]] || {
        log_warning "No CI/CD runs found"
        return 1
    }

    echo "$data"
}

# Parse and display run info with robust handling
show_run() {
    local run_data="$1"
    local show_details="${2:-true}"

    # Safe extraction with defaults
    local status=$(echo "$run_data" | jq -r '.status // "unknown"')
    local conclusion=$(echo "$run_data" | jq -r '.conclusion // "unknown"')
    local name=$(echo "$run_data" | jq -r '.workflowName // "Unknown Workflow"')
    local branch=$(echo "$run_data" | jq -r '.headBranch // "unknown"')
    local run_id=$(echo "$run_data" | jq -r '.databaseId // "unknown"')
    local url=$(echo "$run_data" | jq -r '.url // "#"')
    local created=$(echo "$run_data" | jq -r '.createdAt // "unknown"')
    local title=$(echo "$run_data" | jq -r '.displayTitle // "CI/CD Run"')

    # Status determination
    local status_symbol="❓" status_color=$YELLOW
    case "$status" in
        "completed")
            case "$conclusion" in
                "success") status_symbol=$SUCCESS; status_color=$GREEN ;;
                "failure") status_symbol=$FAILURE; status_color=$RED ;;
                "cancelled") status_symbol="🚫"; status_color=$YELLOW ;;
            esac
            ;;
        "in_progress") status_symbol=$PROGRESS ;;
        "queued") status_symbol="🕒" ;;
    esac

    # Display
    echo
    echo -e "${GRAY}----------------------------------------${NC}"
    echo -e "${BOLD}${status_symbol} ${title}${NC}"
    echo -e "${GRAY}Workflow:${NC} ${CYAN:-}${name}${NC}"
    echo -e "${GRAY}Branch:${NC} ${branch}"
    echo -e "${GRAY}Run ID:${NC} ${run_id}"
    echo -e "${GRAY}Created:${NC} ${created}"
    echo -e "${GRAY}Status:${NC} ${status_color}${status:-$status}${NC}"
    echo -e "${GRAY}Link:${NC} ${url}"

    [[ "$show_details" == "true" ]] && analyze_details "$run_id" "$status" "$conclusion"
}

# Analyze run details with robust error handling
analyze_details() {
    local run_id="$1" status="$2" conclusion="$3"

    [[ "$conclusion" != "failure" ]] && return 0

    echo
    log_warning "Failure Analysis:"

    # Get failed jobs safely
    if failed_jobs=$(gh run view "$run_id" --json jobs \
        --jq '.jobs[] | select(.conclusion=="failure") | {name, steps: [.steps[] | select(.conclusion=="failure") | .name]}' 2>/dev/null); then
        [[ -n "$failed_jobs" ]] && {
            echo -e "${RED}Failed Jobs:${NC}"
            echo "$failed_jobs" | jq -r '"  • " + .name + " ❌" + (if .steps | length > 0 then "\n    Failed: " + (.steps | join(", ")) else "" end)'
        }
    fi

    # Smart error pattern detection
    if gh run view "$run_id" --log-failed > "$TEMP_DIR/errors.log" 2>/dev/null && [[ -s "$TEMP_DIR/errors.log" ]]; then
        echo
        echo -e "${RED}Key Error Patterns:${NC}"

        # TypeScript errors
        if grep -q "error TS" "$TEMP_DIR/errors.log"; then
            echo -e "  ${RED}•${NC} TypeScript compilation errors"
            grep "error TS" "$TEMP_DIR/errors.log" | head -2 | sed 's/^/    /'
        fi

        # Test failures
        if grep -q "FAIL\|Test.*fail" "$TEMP_DIR/errors.log"; then
            echo -e "  ${RED}•${NC} Test failures"
            grep -E "FAIL|Test.*fail" "$TEMP_DIR/errors.log" | head -2 | sed 's/^/    /'
        fi

        # Build/dependency issues
        if grep -q "npm.*ERR\|ELIFECYCLE" "$TEMP_DIR/errors.log"; then
            echo -e "  ${RED}•${NC} Build/dependency issues"
            grep -E "npm.*ERR|ELIFECYCLE" "$TEMP_DIR/errors.log" | head -1 | sed 's/^/    /'
        fi
    fi

    suggest_fixes "$run_id"
}

# Provide smart fix suggestions
suggest_fixes() {
    local run_id="$1"
    echo
    log_info "💡 Quick Fixes:"

    if [[ -f "$TEMP_DIR/errors.log" ]]; then
        if grep -q "next-themes\|TypeScript" "$TEMP_DIR/errors.log"; then
            echo -e "  ${GREEN}•${NC} Check TypeScript imports and types"
        fi
        if grep -q "permission" "$TEMP_DIR/errors.log"; then
            echo -e "  ${GREEN}•${NC} Verify file permissions"
        fi
        if grep -q "npm.*ERR" "$TEMP_DIR/errors.log"; then
            echo -e "  ${GREEN}•${NC} Run: npm install or pnpm install"
        fi
    fi

    echo -e "  ${GREEN}•${NC} Full logs: ${BLUE}gh run view $run_id --log-failed${NC}"
    echo -e "  ${GREEN}•${NC} Rerun: ${BLUE}gh run rerun $run_id${NC}"
    echo -e "  ${GREEN}•${NC} Test locally: ${BLUE}make test && make build${NC}"
}

# Show history with robust formatting
show_history() {
    local count=${1:-$HISTORY_COUNT}
    log_info "Fetching last $count runs..."

    local data
    data=$(fetch_data "$count") || return 1

    echo
    log_info "${CHART} CI/CD History:"
    echo -e "${GRAY}----------------------------------------${NC}"

    # Safe parsing with error handling
    echo "$data" | jq -r '.[] |
        [
            (.databaseId // "unknown" | tostring),
            (.workflowName // "Unknown"),
            (.headBranch // "unknown"),
            (.status // "unknown"),
            (.conclusion // "unknown"),
            (.createdAt // "unknown")
        ] | @tsv' 2>/dev/null | while IFS=$'\t' read -r run_id name branch status conclusion created; do

        # Handle empty values
        : "${run_id:=unknown}" "${name:=Unknown}" "${branch:=unknown}" \
          "${status:=unknown}" "${conclusion:=unknown}" "${created:=unknown}"

        # Status symbol
        local symbol="❓"
        case "$status" in
            "completed")
                case "$conclusion" in
                    "success") symbol=$SUCCESS ;;
                    "failure") symbol=$FAILURE ;;
                    "cancelled") symbol="🚫" ;;
                esac
                ;;
            "in_progress") symbol=$PROGRESS ;;
            "queued") symbol="🕒" ;;
        esac

        printf "%-8s %s %-22s %-10s %s\n" \
            "${GRAY}#${run_id}${NC}" \
            "$symbol" \
            "${name:0:25}" \
            "$branch" \
            "${created:0:16}" | sed 's/\\033\[[0-9;]*m//g'
    done
}

# Deep analysis with statistics
analyze_deep() {
    log_info "Performing deep analysis..."
    log_progress "Analyzing patterns" 3 30

    local data
    data=$(fetch_data 10) || return 1

    local total=$(echo "$data" | jq 'length')
    local successful=$(echo "$data" | jq '[.[] | select(.conclusion == "success")] | length')
    local failed=$(echo "$data" | jq '[.[] | select(.conclusion == "failure")] | length')
    local success_rate=$((successful * 100 / total))

    echo
    log_info "${CHART} CI/CD Statistics:"
    echo -e "${GRAY}----------------------------------------${NC}"
    echo -e "Total Runs: ${total}"
    echo -e "Successful: ${GREEN}${successful}${NC}"
    echo -e "Failed: ${RED}${failed}${NC}"
    echo -e "Success Rate: ${success_rate}%"

    # Recommendations
    echo
    log_info "💡 Recommendations:"
    [[ $success_rate -lt 80 ]] && echo -e "  ${WARNING} Low success rate - review build stability${NC}"
    echo -e "  ${GREEN}•${NC} Use pre-commit hooks to catch issues early"
    echo -e "  ${GREEN}•${NC} Monitor build times and optimize dependencies"
    echo -e "  ${GREEN}•${NC} Consider parallel test execution"
}

# Watch live run
watch_live() {
    log_info "Finding running workflows..."

    local running_data
    running_data=$(fetch_data 5 | jq '.[] | select(.status == "in_progress" or .status == "queued")')

    if [[ -n "$running_data" ]]; then
        local run_id=$(echo "$running_data" | jq -r '.databaseId')
        log_info "Found running workflow: #$run_id"
        echo -e "${BLUE}Starting live monitor...${NC}"
        gh run watch "$run_id" --exit-status
    else
        log_warning "No running workflows found"
        echo -e "${BLUE}Tip:${NC} Push a commit to trigger a new run"
    fi
}

# Rerun failed workflow
rerun_failed() {
    log_progress "Finding failed runs" 2 15

    local failed_data
    failed_data=$(fetch_data 5 | jq '.[] | select(.conclusion == "failure")')

    if [[ -n "$failed_data" ]]; then
        local run_id=$(echo "$failed_data" | jq -r '.databaseId' | head -1)
        log_info "Found failed run: #$run_id"
        read -p "Rerun this workflow? [y/N]: " confirm

        if [[ $confirm =~ ^[Yy]$ ]]; then
            log_progress "Rerunning workflow" 2 20
            gh run rerun "$run_id" && log_success "Rerun initiated"
        fi
    else
        log_info "No failed runs found in recent history"
    fi
}

# Interactive menu
show_menu() {
    while true; do
        echo
        echo -e "${BOLD}🎛️  CI/CD Menu:${NC}"
        echo "  1) Current Status"
        echo "  2) Watch Live"
        echo "  3) History"
        echo "  4) Deep Analysis"
        echo "  5) Rerun Failed"
        echo "  6) Exit"
        echo

        read -p "Choose [1-6]: " choice

        case $choice in
            1) main ;;
            2) watch_live ;;
            3) show_history ;;
            4) analyze_deep ;;
            5) rerun_failed ;;
            6) break ;;
            *) log_warning "Invalid option" ;;
        esac
    done
}

# Help information
show_help() {
    cat << 'EOF'
EchoMind CI/CD Monitor

USAGE:
    ./scripts/ci.sh [COMMAND] [OPTIONS]

COMMANDS:
    (none)          Show current status
    watch           Watch live running workflow
    history [N]     Show last N runs (default: 5)
    analyze         Deep analysis of recent runs
    interactive     Interactive menu mode
    help            Show this help

ENVIRONMENT:
    HISTORY_COUNT    Default history count (default: 5)
    BRANCH          Target branch (default: main)

EXAMPLES:
    ./scripts/ci.sh                 # Current status
    ./scripts/ci.sh watch           # Watch live
    ./scripts/ci.sh history 10      # Last 10 runs
    ./scripts/ci.sh analyze         # Deep analysis
    ./scripts/ci.sh interactive     # Interactive menu

ALIASES:
    alias ci='./scripts/ci.sh'

REQUIREMENTS:
    - GitHub CLI (gh)
    - jq (JSON processor)
    - GitHub authentication: gh auth login
EOF
}

# Main function
main() {
    local cmd=${1:-"status"}
    shift || true

    # Header
    echo -e "${BLUE}${BOLD}🔍 EchoMind CI/CD Monitor${NC}"
    echo -e "${GRAY}========================================${NC}"

    # Check dependencies
    check_deps

    # Route command
    case $cmd in
        "status"|"")
            log_progress "Checking latest run" 2 20
            local data
            data=$(fetch_data 1) || {
                log_warning "No recent runs found"
                echo -e "${BLUE}Quick actions:${NC}"
                echo "  ./scripts/ci.sh watch      # Watch live"
                echo "  ./scripts/ci.sh history    # View history"
                echo "  ./scripts/ci.sh interactive # Interactive mode"
                return 1
            }
            echo "$data" | jq -c '.[0]' | while read -r run; do
                show_run "$run" false
            done
            ;;
        "watch")
            watch_live
            ;;
        "history")
            show_history "${1:-$HISTORY_COUNT}"
            ;;
        "analyze")
            analyze_deep
            ;;
        "interactive"|"menu")
            show_menu
            ;;
        "help"|"-h"|"--help")
            show_help
            ;;
        *)
            log_error "Unknown command: $cmd"
            echo "Run './scripts/ci.sh help' for usage information"
            return 1
            ;;
    esac
}

# Execute main with all arguments
main "$@"