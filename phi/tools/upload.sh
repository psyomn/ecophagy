#!/usr/bin/env bash
TOKEN="$(cat current-token)"
AUTH_HEADER="Authorization: token $TOKEN"
URL="http://127.0.0.1:9876/upload"

echo "using auth token: $TOKEN"
echo "auth header: $AUTH_HEADER"

if [ ! -f "balls.jpg" ]; then
    echo "no balls.jpg found -- downloading"
    wget -O balls.jpg 'https://upload.wikimedia.org/wikipedia/commons/1/16/HDRI_Sample_Scene_Balls_%28JPEG-HDR%29.jpg'
fi

curl -v -s -i \
     -H "$AUTH_HEADER" \
     -H "Content-Type: application/octet-stream" \
     --data-binary @balls.jpg \
     -X POST \
     "$URL"
