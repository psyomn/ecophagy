/*
more experimental than anything.

Copyright 2019-2024 Simon Symeonidis (psyomn)

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
	"fmt"
	"log"
	"os"

	"git.sr.ht/~psyomn/ecophagy/psy/barf"
	"git.sr.ht/~psyomn/ecophagy/psy/common"
	"git.sr.ht/~psyomn/ecophagy/psy/filebucket"
	"git.sr.ht/~psyomn/ecophagy/psy/gh"
	"git.sr.ht/~psyomn/ecophagy/psy/git"
	"git.sr.ht/~psyomn/ecophagy/psy/memo"
	"git.sr.ht/~psyomn/ecophagy/psy/mock"
	"git.sr.ht/~psyomn/ecophagy/psy/notes"
	"git.sr.ht/~psyomn/ecophagy/psy/uploader"
)

type command struct {
	name string
	fn   func(common.RunParams) common.RunReturn
	desc string
}

func makeCommands() []command {
	return []command{
		{"barf", barf.Run, "run code barfer"},
		{"filebucket", filebucket.Run, "bucket duplicate files"},
		{"gh", gh.Run, "personal github utils"},
		{"git", git.Run, "run git helper"},
		{"memo", memo.Run, "description on files in the system"},
		{"mock", mock.Run, "run tcp/udp mocker"},
		{"notes", notes.Run, "run notes webserver"},
		{"upload", uploader.Run, "run the uploader tool"},
		{"help", help, "print help"},
	}
}

func help(_ common.RunParams) common.RunReturn {
	fmt.Println("usage:")
	commands := makeCommands()
	for _, c := range commands {
		fmt.Println("\t", c.name, "\t", c.desc)
	}
	return nil
}

func main() {
	args := os.Args
	commands := makeCommands()

	if len(args) < 2 {
		_ = help(nil)
		os.Exit(1)
	}

	cmd := args[1]
	rest := args[2:]
	var callfn func(common.RunParams) common.RunReturn

	for _, c := range commands {
		if cmd == c.name {
			callfn = c.fn
			break
		}
	}

	if callfn == nil {
		log.Println("no such command: ", cmd)
		os.Exit(1)
	}

	err := callfn(rest)
	if err != nil {
		log.Println("error:", err)
		os.Exit(1)
	}
}
