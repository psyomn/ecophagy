/*
Package mock will mock tcp/udp endpoints

Copyright 2020 Simon Symeonidis (psyomn)

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
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
