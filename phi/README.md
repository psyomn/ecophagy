# Phi

Originally taken from [phi](https://github.com/psyomn/phi), but decided
to leave it here as it would probably have a better home here.

This is very barebones. Use at your own risk. Please only use at home.

Is a family photo backup tool. Originally is supposed to be an
alternate way to sync photographs from different devices.

Use cases:
- be mainly for personal use, on self hosted server/home server
- be secure-ish
- support multiple users (no sharing with each other though)
- not supposed to be for more than 10 users
- very little concurrency concerns

Soft Dependencies:
- exiftool: will add metadata through embedding json inside the
  comment exif field. If you don't have exiftool, then images will
  simply not be tagged.

# start simple server:

```bash
phi-server -config phi/phi-server/config.sample.json
```

Make sure to read the configuration.
