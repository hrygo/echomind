#!/bin/bash

# scripts/check_ci.sh - Check the latest GitHub Actions CI/CD status
# Usage: ./scripts/check_ci.sh

# Colors
GREEN='\033[0;32m'
BLUE='\033[0;34m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 1. Check for gh CLI
if ! command -v gh &> /dev/null; then
    echo -e "${RED}Error: GitHub CLI (gh) is not installed.${NC}"
    echo "Please install it using Homebrew: brew install gh"
    exit 1
fi

# 2. Check Login Status
if ! gh auth status &> /dev/null; then
     echo -e "${RED}Error: You are not logged into GitHub CLI.${NC}"
     echo "Please run: gh auth login"
     exit 1
fi

echo -e "${BLUE}Fetching latest CI/CD pipeline status...${NC}"

# 3. Fetch Data
# We use a custom delimiter '|' to safely separate fields even if they contain spaces
# Fields: status, conclusion, workflowName, url, headBranch, databaseId
DATA=$(gh run list --limit 1 --json status,conclusion,workflowName,url,headBranch,databaseId --jq '.[0] | "\(.status)|\(.conclusion)|\(.workflowName)|\(.url)|\(.headBranch)|\(.databaseId)"')

if [ -z "$DATA" ] || [ "$DATA" == "null" ]; then
    echo -e "${YELLOW}No workflow runs found for this repository.${NC}"
    exit 0
fi

# 4. Parse Data
IFS='|' read -r STATUS CONCLUSION NAME URL BRANCH RUN_ID <<< "$DATA"

# 5. Display Status
echo -e "----------------------------------------"
echo -e "Workflow: ${BLUE}${NAME}${NC}"
echo -e "Branch:   ${BRANCH}"
echo -e "Run ID:   ${RUN_ID}"

if [ "$STATUS" == "in_progress" ]; then
    echo -e "Status:   ${YELLOW}IN PROGRESS${NC} ⏳"
elif [ "$STATUS" == "queued" ]; then
    echo -e "Status:   ${YELLOW}QUEUED${NC} 🕒"
elif [ "$CONCLUSION" == "success" ]; then
    echo -e "Status:   ${GREEN}SUCCESS${NC} ✅"
elif [ "$CONCLUSION" == "failure" ]; then
    echo -e "Status:   ${RED}FAILURE${NC} ❌"
    
    echo -e "\n${RED}Failure Details:${NC}"
    # List failed jobs and their failed steps
    gh run view "$RUN_ID" --json jobs --jq '.jobs[] | select(.conclusion=="failure") | "  • Job: \(.name)\n    ❌ Step: \(.steps[] | select(.conclusion=="failure") | .name)"'

    echo -e "\n${RED}--- Error Logs (Tail) ---${NC}"
    # Fetch logs for failed steps, save to temp file to process
    if gh run view "$RUN_ID" --log-failed > /tmp/gh_failed_log.txt 2>/dev/null; then
        if [ -s /tmp/gh_failed_log.txt ]; then
             # Filter out ANSI color codes for cleaner length check if needed, but tail works fine.
             # We want to show enough context.
             tail -n 20 /tmp/gh_failed_log.txt
             echo -e "\n${BLUE}Tip: Run 'gh run view $RUN_ID --log-failed' for full logs.${NC}"
        else
             echo -e "${YELLOW}(No detailed logs returned by GitHub)${NC}"
        fi
        rm -f /tmp/gh_failed_log.txt
    else
        echo -e "${YELLOW}Could not fetch logs automatically.${NC}"
    fi

elif [ "$CONCLUSION" == "cancelled" ]; then
    echo -e "Status:   ${YELLOW}CANCELLED${NC} 🚫"
else
    echo -e "Status:   ${STATUS} (${CONCLUSION})"
fi

echo -e "Link:     ${URL}"
echo -e "----------------------------------------"

# 6. Interactive Watch Option
# If the workflow is running and we are in an interactive terminal, offer to watch it.
if [[ "$STATUS" == "in_progress" || "$STATUS" == "queued" ]]; then
    if [ -t 1 ]; then
        echo -e "\n${BLUE}Would you like to watch the live logs?${NC}"
        read -p "Press [Enter] to watch, or [n] to exit: " choice
        if [[ -z "$choice" || "$choice" =~ ^[Yy]$ ]]; then
            gh run watch "$RUN_ID" --exit-status
        fi
    fi
fi
