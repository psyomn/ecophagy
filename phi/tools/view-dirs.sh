#!/usr/bin/env bash
TOKEN="$(cat current-token)"
AUTH_HEADER="Authorization: token $TOKEN"

URL="http://127.0.0.1:9876/view/"

curl -v -s -i \
     -H "$AUTH_HEADER" \
     "$URL"
