#!/usr/bin/env bash
source "credentials.sh"
curl -s -i -X POST "http://127.0.0.1:9876/register" -d "{\"username\":\"$PHIUSER\",\"password\":\"$PHIPASS\"}"
