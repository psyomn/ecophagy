#!/usr/bin/env bash
TOKEN="$(cat current-token)"
AUTH_HEADER="Authorization: token $TOKEN"

URL="http://127.0.0.1:9876/view/2020-09-03/balls.jpg"

curl -s \
     -H "$AUTH_HEADER" \
     "$URL" > fetch-balls.jpg
