/*
Package barf contains code relevant to barfing.

Copyright 2019-2021 Simon Symeonidis (psyomn)

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
package barf

import (
	"fmt"

	"github.com/psyomn/ecophagy/psy/common"
)

var runCommands = map[string]func(common.RunParams) common.RunReturn{
	"cmake":    cmake,
	"ada":      ada,
	"lilypond": lilypond,
}

func printUsage() {
	fmt.Println("usage:")
	fmt.Println("  barf <target>")
	fmt.Println("current targets: ")
	for k := range runCommands {
		fmt.Println(" ", k)
	}
}

// Run will run the barf command, that should barf out specific
// configurations on the fly. I always wanted something like this so
// that I could bootstrap new projects, and get rid of boilerplate.
func Run(args common.RunParams) common.RunReturn {
	if len(args) == 0 {
		printUsage()
		return ErrNeedProjectName
	}

	cmdFn, ok := runCommands[args[0]]
	if !ok {
		printUsage()
		return ErrNoSuchCommand
	}

	return cmdFn(args[1:])
}
