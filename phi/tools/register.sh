#!/usr/bin/env bash
source "credentials.sh"
source "constants.sh"
curl -s -i -X POST "$PHI_SERVER_HOST:$PHI_SERVER_PORT/register" \
     -d "{\"username\":\"$PHIUSER\",\"password\":\"$PHIPASS\"}"
