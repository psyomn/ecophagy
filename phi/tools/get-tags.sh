#!/usr/bin/env bash
TOKEN="$(cat current-token)"
AUTH_HEADER="Authorization: token $TOKEN"

URL="http://127.0.0.1:9876/tag/2014-10-08/balls.jpg"

curl -H "$AUTH_HEADER" -v -i "$URL"
