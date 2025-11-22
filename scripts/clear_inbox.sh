#!/bin/bash

# EchoMind Inbox Cleaner Script
# This script allows you to clear all emails for a specific user.

# Configuration
API_URL="http://localhost:8080/api/v1"
TOKEN=$1

# Colors
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo -e "${YELLOW}🧹 EchoMind Inbox Cleaner${NC}"

# Display usage instructions if requested or if no arguments
if [[ "$1" == "--help" || "$1" == "-h" ]]; then
    echo -e "\n${GREEN}Usage:${NC}"
    echo -e "  ./scripts/clear_inbox.sh                  - ${YELLOW}Test Mode:${NC} Registers/Logs in 'cleaner@test.com' and clears its inbox."
    echo -e "                                          This is useful for development and testing."
    echo -e "  ./scripts/clear_inbox.sh \"YOUR_JWT_TOKEN\" - ${YELLOW}Live Mode:${NC} Clears the inbox for the user associated with the provided JWT token."
    echo -e "                                          To get your JWT token, log in to the frontend, open browser developer tools (F12),\"
    echo -e "                                          go to the 'Network' tab, find any API request (e.g., '/api/v1/emails'),\"
    echo -e "                                          and copy the 'Authorization' header value (excluding 'Bearer ')"
    echo -e "\n${YELLOW}Warning:${NC} This action is irreversible and will delete ALL emails for the target user."
    exit 0
fi

if [ -z "$TOKEN" ]; then
    echo -e "${YELLOW}⚠️  No token provided. Switching to Test Account Mode.${NC}"
    echo "🔄 Registering/Logging in temporary user 'cleaner@test.com'..."
    
    # Try login first to get token if user exists
    LOGIN_RESP=$(curl -s -X POST "$API_URL/auth/login" \
        -H "Content-Type: application/json" \
        -d '{"email": "cleaner@test.com", "password": "password123"}')
    
    # Extract token using simple grep/cut (avoiding jq dependency for portability)
    TOKEN=$(echo $LOGIN_RESP | grep -o '"token":"[^"']*' | cut -d'"' -f4)
    
    # If login failed or returned error (no token), try register
    if [ -z "$TOKEN" ] || [ "$TOKEN" == "null" ]; then
        echo "Login failed or user doesn't exist, trying to register..."
        REGISTER_RESP=$(curl -s -X POST "$API_URL/auth/register" \
            -H "Content-Type: application/json" \
            -d '{"name": "Cleaner", "email": "cleaner@test.com", "password": "password123"}')
        TOKEN=$(echo $REGISTER_RESP | grep -o '"token":"[^"']*' | cut -d'"' -f4)
    fi
    
    if [ -z "$TOKEN" ] || [ "$TOKEN" == "null" ]; then
        echo -e "${RED}❌ Failed to get token. Is the backend running?${NC}"
        echo "Response: $LOGIN_RESP $REGISTER_RESP"
        exit 1
    fi
    
    echo -e "${GREEN}✅ Authenticated as 'cleaner@test.com'${NC}"
else
    echo -e "${GREEN}🔑 Using provided token...${NC}"
fi

echo "🗑️  Sending delete request..."
RESPONSE=$(curl -s -X DELETE "$API_URL/emails/all" \
     -H "Authorization: Bearer $TOKEN")

echo -e "${GREEN}📩 Server Response:${NC} $RESPONSE"