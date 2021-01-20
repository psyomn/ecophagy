# PROTOCOL

Slight specifications for the protocol. HTTPS seems an okay enough
protocol for sharing bigger files accross the network, when
considering complexity to implement.

All the bellow actions assume that you're logged in. A login token is
passed through headers of the requests.

## API

```nocode
POST /register/
    Register user to serice.
    {"username": username,
     "password": password}
    RETURN 400, if errors (eg, user exists, etc)
    RETURN 200, with empty body

POST /login/
    {"username": username,
     "password": password}
    RETURN 200
        {"token": token}

Login user to service.

POST /upload/<filename>/<timestamp>
    Header/Authorization: token <token>
    Header/Content-Type: application/octet-stream
    BODY: <IMAGE-BINARY-DATA>
    RETURN 400, on bad credentials
    RETURN 200, on success

View current dated directories:

GET /view/
    Header/Authorization: token <token>
    Header/Content-Type: application/json
    RETURN 400, on bad credentials
    RETURN 200, text, on success
    {"directories": ["2021-01-01", "2021-01-02"]}

View files in given directory:

GET /view/yyyy-mm-dd/
    Header/Authorization: token <token>
    Header/Content-Type: application/json
    RETURN 400, on bad credentials
    RETURN 200, text, on success
    {"files": ["one.jpg", "two.jpg"]}

Get and view a picture:

GET /view/yyyy-mm-dd/picture.jpg
    Header/Authorization: token <token>
    Header/Content-Type: application/json
    RETURN 400, on bad credentials
    RETURN 200, data, on success

Get current tags of picture:

GET /tag/yyy-mm-dd/picture.jpg
    Header/Authorization: token <token>
    Header/Content-Type: application/json
    RETURN 500, on missing exif lib
    RETURN 400, on bad credentials
    RETURN 200, data, on success
    {"tags": ["one", "two", "three"]}

Update tags of a picture:

PATCH /tag/yyy-mm-dd/picture.jpg
    Header/Authorization: token <token>
    Header/Content-Type: application/json
    RETURN 500, on missing exif lib
    RETURN 400, on bad credentials
    RETURN 200, data, on success
    {"tags": ["set", "your", "new", "tags"]}

Get the status of the server:

GET /status/
    Header/Content-Type: application/json
    RETURN 200, data, on success
    {"status": "ok", "version": "1.0.0"}
```

Assuming that the device store on itself the date that we are
interested in. The date can't be preserved via other means, so we pass
it in the url. The server can then modify the modtime accordingly.

## TODO
- How to deal with other image formats that don't support exif? (how
  easy is it to convert pictures to jpg?)
