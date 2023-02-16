/*
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
package main

import (
	"encoding/json"
)

// userComment is the way we want to package the json object inside
// the user comment section of the photograph. It should look
// something like this:
//
//	{ "phi": {
//	    "username": "someone",
//	    "timestamp": 12312312312,
//	    "tags": ["me", "mom", "woods", vacation"]
//	  }
//	}
//
// TODO: this might be a good common utility. It might make sense to
//
//	extract it into the 'common' package
type phi struct {
	Username  string   `json:"username"`
	Timestamp int64    `json:"timestamp"`
	Tags      []string `json:"tags"`
}

type userComment struct {
	Phi phi `json:"phi"`
}

func (s *userComment) toJSON() []byte {
	bytes, err := json.Marshal(s)
	if err != nil {
		panic("unsupported json marshal field: " + err.Error())
	}

	return bytes
}
