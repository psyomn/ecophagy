/*
Package git just has some dumb helpers. I mainly want to use this so
that I can sort lists of semvers real fast.

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
package git

import (
	"fmt"
	"os/exec"
	"sort"
	"strings"

	"github.com/psyomn/ecophagy/psy/common"

	"github.com/Masterminds/semver"
)

func Run(_ common.RunParams) common.RunReturn {
	// Right now, the only thing that I'm looking for is for
	// something that can quickly list me the 5 latest tags on the
	// local repository. This is a crappy script that should help
	// and do that.

	cmd := exec.Command("git", "tag")
	output, err := cmd.Output()
	if err != nil {
		return err
	}

	outputStr := string(output)
	lines := strings.Split(outputStr, "\n")

	semvers := make([]*semver.Version, 0)
	for index := range lines {
		if lines[index] == "" {
			// last line might be blank; ignore
			continue
		}

		v, err := semver.NewVersion(lines[index])
		if err != nil {
			fmt.Println("warn: could not parse: ", err, " version:", lines[index])
			continue
		}

		semvers = append(semvers, v)
	}

	sort.Sort(semver.Collection(semvers))

	numSemvers := len(semvers)
	lastFive := numSemvers - 5
	if lastFive < 0 {
		lastFive = 0
	}

	lastFiveSemvers := semvers[lastFive:]

	fmt.Println("last 5 tags: ")
	for _, el := range lastFiveSemvers {
		fmt.Println("  ", el)
	}

	return nil
}
