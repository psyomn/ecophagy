package mock

const example = `---
some-service:
  type: tcp
  port: 9999
  return: nil # for no returns

some-service-2:
  type: tcp
  port: 9998
  return: "the text"

byte-tcp-service-3:
  type: tcp
  port: 9997
  return: 13,14,15

udp-service:
  type: udp
  port: 9996
  return: "blah"

udp-service-bytes:
  type: udp
  port: 9995
  return: 12,13,14

http-service:
  type: http
  port: 9994
  return: "<body> hello </body>"
  root: "/"
`
