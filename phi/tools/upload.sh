#!/usr/bin/env bash
TOKEN="$(cat current-token)"
AUTH_HEADER="Authorization: token $TOKEN"
FILENAME="balls.jpg"
SAMPLE_LOCATION='https://upload.wikimedia.org/wikipedia/commons/1/16/HDRI_Sample_Scene_Balls_%28JPEG-HDR%29.jpg'

echo "using auth token: $TOKEN"
echo "auth header: $AUTH_HEADER"

if [ ! -f "$FILENAME" ]; then
    echo "no balls.jpg found -- downloading"
    wget -O "$FILENAME" "$SAMPLE_LOCATION"
fi

TIMESTAMP=$(stat --printf="%Y" "$FILENAME")
URL="http://127.0.0.1:9876/upload/$FILENAME/$TIMESTAMP"
curl -v -s -i \
     -H "$AUTH_HEADER" \
     -H "Content-Type: application/octet-stream" \
     --data-binary @"$FILENAME" \
     -X POST \
     "$URL"
