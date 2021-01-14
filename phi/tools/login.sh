#!/usr/bin/env bash
# set -x -e 
source "credentials.sh"

LOGIN_OUTPUT=$(curl -s -i \
                    -X POST \
                    -d "{\"username\":\"$PHIUSER\",\"password\":\"$PHIPASS\"}" \
                    "http://127.0.0.1:9876/login")

echo "$LOGIN_OUTPUT"
echo "$LOGIN_OUTPUT" | jq -j -r -R 'fromjson? | .token' > current-token
