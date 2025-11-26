#!/bin/bash

# Enhanced CI/CD Monitor and Analyzer for EchoMind Project
# Usage: ./scripts/check_ci_enhanced.sh [options]
# Options: --watch, --history N, --analyze, --compare RUN1 RUN2

# Colors & Formatting
RED='\033[0;31m'
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
PURPLE='\033[0;35m'
CYAN='\033[0;36m'
WHITE='\033[1;37m'
GRAY='\033[0;37m'
BOLD='\033[1m'
DIM='\033[2m'
NC='\033[0m' # No Color

# Symbols
SUCCESS="✅"
FAILURE="❌"
WARNING="⚠️"
INFO="ℹ️"
PROGRESS="⏳"
GEAR="⚙️"
CHART="📊"

# Config
DEFAULT_HISTORY=5
TEMP_DIR="/tmp/echomind_ci_$$"
HISTORY_FILE="$HOME/.echomind_ci_history"

# Create temp directory
mkdir -p "$TEMP_DIR"
trap "rm -rf $TEMP_DIR" EXIT

# Helper functions
print_header() {
    echo -e "${CYAN}${BOLD}🔍 EchoMind CI/CD Enhanced Monitor${NC}"
    echo -e "${GRAY}========================================${NC}"
}

print_separator() {
    echo -e "${GRAY}----------------------------------------${NC}"
}

print_status() {
    local status=$1
    local message=$2
    case $status in
        "success") echo -e "${GREEN}${SUCCESS} ${message}${NC}" ;;
        "failure") echo -e "${RED}${FAILURE} ${message}${NC}" ;;
        "warning") echo -e "${YELLOW}${WARNING} ${message}${NC}" ;;
        "info") echo -e "${BLUE}${INFO} ${message}${NC}" ;;
        "progress") echo -e "${YELLOW}${PROGRESS} ${message}${NC}" ;;
    esac
}

show_progress() {
    local duration=$1
    local steps=$2
    local current=0

    while [ $current -le $steps ]; do
        local percent=$((current * 100 / steps))
        local filled=$((percent / 2))
        local empty=$((50 - filled))

        printf "\r${YELLOW}[${PROGRESS}]${NC} "
        printf "${GREEN}%*s${NC}" $filled | tr ' ' '█'
        printf "${GRAY}%*s${NC}" $empty | tr ' ' '░'
        printf " ${percent}%%"

        sleep $(($duration / $steps))
        ((current++))
    done
    echo
}

# Validate prerequisites
check_prerequisites() {
    local missing_tools=()

    if ! command -v gh &> /dev/null; then
        missing_tools+=("gh - GitHub CLI (brew install gh)")
    fi

    if ! command -v jq &> /dev/null; then
        missing_tools+=("jq - JSON processor (brew install jq)")
    fi

    if [ ${#missing_tools[@]} -gt 0 ]; then
        print_status "failure" "Missing required tools:"
        for tool in "${missing_tools[@]}"; do
            echo -e "  ${RED}•${NC} $tool"
        done
        exit 1
    fi

    if ! gh auth status &> /dev/null; then
        print_status "failure" "Not logged into GitHub CLI"
        echo -e "${BLUE}Fix:${NC} Run 'gh auth login'"
        exit 1
    fi
}

# Fetch CI/CD data with enhanced error handling
fetch_ci_data() {
    local limit=${1:-1}
    local branch=${2:-"main"}

    print_status "progress" "Fetching CI/CD data..."

    local data
    data=$(gh run list --limit "$limit" --branch "$branch" --json status,conclusion,workflowName,url,headBranch,databaseId,createdAt,updatedAt,displayTitle --jq '. | sort_by(.databaseId) | reverse' 2>/dev/null)

    if [ $? -ne 0 ] || [ -z "$data" ] || [ "$data" == "null" ]; then
        print_status "failure" "Failed to fetch CI/CD data"
        echo -e "${BLUE}Possible causes:${NC}"
        echo -e "  • No internet connection"
        echo -e "  • GitHub API rate limit exceeded"
        echo -e "  • No workflow runs found"
        echo -e "  • Invalid repository access"
        return 1
    fi

    echo "$data"
}

# Parse and display run information
display_run_info() {
    local run_data="$1"
    local show_details=${2:-true}

    local status=$(echo "$run_data" | jq -r '.status // "unknown"')
    local conclusion=$(echo "$run_data" | jq -r '.conclusion // "unknown"')
    local name=$(echo "$run_data" | jq -r '.workflowName // "Unknown"')
    local branch=$(echo "$run_data" | jq -r '.headBranch // "unknown"')
    local run_id=$(echo "$run_data" | jq -r '.databaseId // "unknown"')
    local url=$(echo "$run_data" | jq -r '.url // "#"')
    local created=$(echo "$run_data" | jq -r '.createdAt // "unknown"')
    local title=$(echo "$run_data" | jq -r '.displayTitle // "CI/CD Run"')

    # Status emoji and color
    local status_emoji="❓"
    local status_color=$YELLOW
    case $status in
        "in_progress")
            status_emoji=$PROGRESS
            status_color=$YELLOW
            ;;
        "queued")
            status_emoji="🕒"
            status_color=$YELLOW
            ;;
        "completed")
            case $conclusion in
                "success")
                    status_emoji=$SUCCESS
                    status_color=$GREEN
                    ;;
                "failure")
                    status_emoji=$FAILURE
                    status_color=$RED
                    ;;
                "cancelled")
                    status_emoji="🚫"
                    status_color=$YELLOW
                    ;;
            esac
            ;;
    esac

    # Header
    echo
    print_separator
    echo -e "${WHITE}${BOLD}${status_emoji} ${title}${NC}"
    echo -e "${GRAY}Workflow:${NC} ${CYAN}${name}${NC}"
    echo -e "${GRAY}Branch:${NC} ${branch}"
    echo -e "${GRAY}Run ID:${NC} ${run_id}"
    echo -e "${GRAY}Created:${NC} ${created}"

    # Status
    if [ "$status" == "completed" ]; then
        echo -e "${GRAY}Status:${NC} ${status_color}${conclusion^^}${NC}"
    else
        echo -e "${GRAY}Status:${NC} ${status_color}${status^^}${NC}"
    fi

    echo -e "${GRAY}Link:${NC} ${url}"

    if [ "$show_details" = true ]; then
        analyze_run_details "$run_id" "$status" "$conclusion"
    fi
}

# Analyze run details with smart filtering
analyze_run_details() {
    local run_id="$1"
    local status="$2"
    local conclusion="$3"

    if [ "$conclusion" = "failure" ]; then
        echo
        print_status "warning" "Failure Analysis:"

        # Get failed jobs
        local failed_jobs
        failed_jobs=$(gh run view "$run_id" --json jobs --jq '.jobs[] | select(.conclusion=="failure") | {name: .name, steps: [.steps[] | select(.conclusion=="failure") | .name]}' 2>/dev/null)

        if [ -n "$failed_jobs" ]; then
            echo -e "${RED}Failed Jobs:${NC}"
            echo "$failed_jobs" | jq -r '"  • " + .name + " ❌" + (if .steps | length > 0 then "\n    Failed steps: " + (.steps | join(", ")) else "" end)'
        fi

        # Get error logs (smarter)
        echo
        print_status "info" "Error Analysis:"
        if gh run view "$run_id" --log-failed > "$TEMP_DIR/failed_logs.txt" 2>/dev/null; then
            if [ -s "$TEMP_DIR/failed_logs.txt" ]; then
                # Extract key error patterns
                echo -e "${RED}Key Error Patterns:${NC}"

                # TypeScript errors
                if grep -q "error TS" "$TEMP_DIR/failed_logs.txt"; then
                    echo -e "  ${RED}•${NC} TypeScript compilation errors"
                    grep "error TS" "$TEMP_DIR/failed_logs.txt" | head -3 | sed 's/^/    /'
                fi

                # Test failures
                if grep -q "FAIL\|Test failed" "$TEMP_DIR/failed_logs.txt"; then
                    echo -e "  ${RED}•${NC} Test failures"
                    grep -E "FAIL|Test failed" "$TEMP_DIR/failed_logs.txt" | head -3 | sed 's/^/    /'
                fi

                # Build errors
                if grep -q "npm ERR\|ELIFECYCLE" "$TEMP_DIR/failed_logs.txt"; then
                    echo -e "  ${RED}•${NC} Build/Dependency errors"
                    grep -E "npm ERR|ELIFECYCLE" "$TEMP_DIR/failed_logs.txt" | head -2 | sed 's/^/    /'
                fi

                # Permission/Docker errors
                if grep -q "permission denied\|docker" "$TEMP_DIR/failed_logs.txt"; then
                    echo -e "  ${RED}•${NC} Permission or Docker issues"
                fi
            fi
        fi

        # Suggest fixes
        suggest_fixes "$run_id"
    fi
}

# Smart fix suggestions based on error patterns
suggest_fixes() {
    local run_id="$1"

    echo
    print_status "info" "💡 Suggested Fixes:"

    # Check logs for common issues
    if gh run view "$run_id" --log-failed > "$TEMP_DIR/fix_logs.txt" 2>/dev/null; then
        if grep -q "next-themes" "$TEMP_DIR/fix_logs.txt"; then
            echo -e "  ${GREEN}•${NC} TypeScript import issues detected"
            echo -e "    Fix: Check import paths in theme components"
        fi

        if grep -q "permission denied" "$TEMP_DIR/fix_logs.txt"; then
            echo -e "  ${GREEN}•${NC} Permission issues detected"
            echo -e "    Fix: Check file permissions and script execution rights"
        fi

        if grep -q "docker" "$TEMP_DIR/fix_logs.txt"; then
            echo -e "  ${GREEN}•${NC} Docker build issues"
            echo -e "    Fix: Check Dockerfile and build context"
        fi

        if grep -q "npm ERR" "$TEMP_DIR/fix_logs.txt"; then
            echo -e "  ${GREEN}•${NC} Node.js dependency issues"
            echo -e "    Fix: Run 'npm install' or 'pnpm install'"
        fi
    fi

    # General suggestions
    echo -e "  ${GREEN}•${NC} Check the full logs: ${BLUE}gh run view $run_id --log-failed${NC}"
    echo -e "  ${GREEN}•${NC} Re-run failed jobs: ${BLUE}gh run rerun $run_id${NC}"
    echo -e "  ${GREEN}•${NC} Run locally: ${BLUE}make test && make build${NC}"
}

# Show historical data
show_history() {
    local count=${1:-$DEFAULT_HISTORY}

    print_status "info" "Fetching last $count CI/CD runs..."

    local history_data
    history_data=$(fetch_ci_data "$count")

    if [ $? -eq 0 ]; then
        echo
        print_status "info" "${CHART} CI/CD History (Last $count runs):"
        print_separator

        echo "$history_data" | jq -r '.[] |
        [
            (.databaseId | tostring),
            (.workflowName // "Unknown"),
            (.headBranch // "unknown"),
            (.status // "unknown"),
            (.conclusion // "unknown"),
            (.createdAt // "unknown")
        ] | @tsv' | while IFS=$'\t' read -r run_id name branch status conclusion created; do
            # Handle empty values
            run_id=${run_id:-"unknown"}
            name=${name:-"Unknown"}
            branch=${branch:-"unknown"}
            status=${status:-"unknown"}
            conclusion=${conclusion:-"unknown"}
            created=${created:-"unknown"}

            local status_symbol="❓"
            case $status in
                "completed")
                    case $conclusion in
                        "success") status_symbol=$SUCCESS ;;
                        "failure") status_symbol=$FAILURE ;;
                        "cancelled") status_symbol="🚫" ;;
                        *) status_symbol="❓" ;;
                    esac
                    ;;
                "in_progress") status_symbol=$PROGRESS ;;
                "queued") status_symbol="🕒" ;;
            esac

            printf "%-8s %s %-20s %-12s %s\n" \
                "${GRAY}#${run_id}${NC}" \
                "${status_symbol}" \
                "${name:0:25}" \
                "${branch}" \
                "${created:0:16}"
        done
    fi
}

# Interactive menu
show_interactive_menu() {
    while true; do
        echo
        print_separator
        echo -e "${WHITE}${BOLD}🎛️  CI/CD Monitor Menu:${NC}"
        echo -e "  ${1} ${BLUE}Current Status${NC}"
        echo -e "  ${2} ${BLUE}Watch Live${NC}"
        echo -e "  ${3} ${BLUE}History (Last 5)${NC}"
        echo -e "  ${4} ${BLUE}Deep Analysis${NC}"
        echo -e "  ${5} ${BLUE}Rerun Failed${NC}"
        echo -e "  ${6} ${BLUE}Exit${NC}"
        print_separator

        read -p "Choose an option [1-6]: " choice

        case $choice in
            1) main ;;
            2) watch_live ;;
            3) show_history ;;
            4) deep_analysis ;;
            5) rerun_failed ;;
            6) break ;;
            *) print_status "warning" "Invalid option" ;;
        esac
    done
}

# Watch live CI/CD run
watch_live() {
    print_status "progress" "Checking for running workflows..."

    local running_data
    running_data=$(fetch_ci_data 5 | jq '.[] | select(.status == "in_progress" or .status == "queued")')

    if [ -n "$running_data" ]; then
        local run_id=$(echo "$running_data" | jq -r '.databaseId')
        print_status "info" "Found running workflow: #$run_id"
        echo -e "${BLUE}Starting live log viewer...${NC}"
        gh run watch "$run_id" --exit-status
    else
        print_status "warning" "No running workflows found"
        echo -e "${BLUE}Tip:${NC} Push a commit to trigger a new run"
    fi
}

# Deep analysis mode
deep_analysis() {
    local history_data
    history_data=$(fetch_ci_data 10)

    if [ $? -eq 0 ]; then
        echo
        print_status "info" "${GEAR} Deep CI/CD Analysis:"
        print_separator

        # Success rate
        local total_runs=$(echo "$history_data" | jq 'length')
        local successful_runs=$(echo "$history_data" | jq '[.[] | select(.conclusion == "success")] | length')
        local failed_runs=$(echo "$history_data" | jq '[.[] | select(.conclusion == "failure")] | length')
        local success_rate=$((successful_runs * 100 / total_runs))

        echo -e "${WHITE}📊 Statistics:${NC}"
        echo -e "  Total Runs: ${total_runs}"
        echo -e "  Successful: ${GREEN}${successful_runs}${NC}"
        echo -e "  Failed: ${RED}${failed_runs}${NC}"
        echo -e "  Success Rate: ${success_rate}%"

        # Common failure patterns
        echo
        echo -e "${WHITE}🔍 Failure Patterns:${NC}"

        # Analyze failed runs
        echo "$history_data" | jq -r '.[] | select(.conclusion == "failure") | .databaseId' | while read -r run_id; do
            echo -e "  ${RED}•${NC} Run #${run_id}:"
            if gh run view "$run_id" --log-failed > "$TEMP_DIR/deep_analysis.txt" 2>/dev/null; then
                if grep -q "TypeScript\|TS" "$TEMP_DIR/deep_analysis.txt"; then
                    echo -e "    ${YELLOW}→${NC} TypeScript errors"
                fi
                if grep -q "test.*fail" "$TEMP_DIR/deep_analysis.txt"; then
                    echo -e "    ${YELLOW}→${NC} Test failures"
                fi
                if grep -q "npm.*ERR\|ELIFECYCLE" "$TEMP_DIR/deep_analysis.txt"; then
                    echo -e "    ${YELLOW}→${NC} Build/dependency issues"
                fi
            fi
        done

        # Recommendations
        echo
        echo -e "${WHITE}💡 Recommendations:${NC}"
        if [ $success_rate -lt 80 ]; then
            echo -e "  ${WARNING} Low success rate detected${NC}"
            echo -e "    Consider reviewing build stability"
        fi
        echo -e "  ${GREEN}•${NC} Set up pre-commit hooks to catch issues early"
        echo -e "  ${GREEN}•${NC} Implement better error handling in workflows"
        echo -e "  ${GREEN}•${NC} Consider parallel test execution for speed"
    fi
}

# Rerun failed jobs
rerun_failed() {
    print_status "progress" "Finding failed runs..."

    local failed_data
    failed_data=$(fetch_ci_data 5 | jq '.[] | select(.conclusion == "failure")')

    if [ -n "$failed_data" ]; then
        local latest_failed=$(echo "$failed_data" | jq -r '.databaseId' | head -1)
        print_status "info" "Found failed run: #$latest_failed"
        read -p "Rerun this failed workflow? [y/N]: " confirm

        if [[ $confirm =~ ^[Yy]$ ]]; then
            print_status "progress" "Rerunning failed workflow..."
            gh run rerun "$latest_failed"
            print_status "success" "Rerun initiated"
        fi
    else
        print_status "info" "No failed runs found in recent history"
    fi
}

# Main function
main() {
    print_header
    check_prerequisites

    # Parse command line arguments
    case "${1:-}" in
        --watch)
            watch_live
            return
            ;;
        --history)
            show_history "${2:-$DEFAULT_HISTORY}"
            return
            ;;
        --analyze)
            deep_analysis
            return
            ;;
        --interactive)
            show_interactive_menu
            return
            ;;
        --help|-h)
            echo "Enhanced CI/CD Monitor for EchoMind"
            echo ""
            echo "Usage: $0 [options]"
            echo ""
            echo "Options:"
            echo "  --watch         Watch live CI/CD run"
            echo "  --history [N]   Show last N runs (default: 5)"
            echo "  --analyze       Deep analysis of recent runs"
            echo "  --interactive   Interactive menu mode"
            echo "  --help          Show this help message"
            echo ""
            echo "Examples:"
            echo "  $0              Show current status"
            echo "  $0 --watch      Watch running workflow"
            echo "  $0 --history 10 Show last 10 runs"
            echo "  $0 --analyze    Deep analysis"
            exit 0
            ;;
    esac

    # Default behavior - show current status
    print_status "progress" "Analyzing latest CI/CD run..."
    show_progress 2 20

    local latest_data
    latest_data=$(fetch_ci_data 1)

    if [ $? -eq 0 ]; then
        local run_count=$(echo "$latest_data" | jq 'length')

        if [ $run_count -eq 0 ]; then
            print_status "warning" "No CI/CD runs found"
            echo -e "${BLUE}Tip:${NC} Push a commit to trigger the first run"
        else
            echo "$latest_data" | jq -c '.[0]' | while read -r run; do
                display_run_info "$run" true
            done

            # Quick actions
            echo
            print_status "info" "Quick Actions:"
            echo -e "  ${BLUE}•${NC} Watch live: $0 --watch"
            echo -e "  ${BLUE}•${NC} View history: $0 --history"
            echo -e "  ${BLUE}•${NC} Deep analysis: $0 --analyze"
            echo -e "  ${BLUE}•${NC} Interactive mode: $0 --interactive"
        fi
    fi
}

# Run main function with all arguments
main "$@"