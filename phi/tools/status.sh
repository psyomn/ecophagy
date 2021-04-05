#!/usr/bin/env bash
source "constants.sh"
curl -s -i "$PHI_SERVER_HOST:$PHI_SERVER_PORT/status"
