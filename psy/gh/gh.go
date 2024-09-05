/*
Package gh is a simple package to interface with github

It's not supposed to do everything. My specific use is to set github
labels, according to a configuration that I think is OK.

usage:

	psy gh test-labels <config.yaml> <owner> <repo>
	  will do a dry run of label nomenclature application

	psy gh poison <config.yaml> <owner> <repo>
	  will edit the labels to respect the given specs.

The config file should look something like this:

	---
	# Github label configuration.

	# T = Type
	# L = Lifecycle

	rename:
	  bug:         ["T-bug",       "a00000", "A software fault or failure"]
	  duplicate:   ["L-duplicate", "16336d", "dOops!"]
	  enhancement: ["T-upkeep",    "009a00", "Anything that rejuvenates"]
	  invalid:     ["L-invalid",   "ffbb0e", "Anything that is invalid"]
	  question:    ["T-question",  "ffaa99", "Questions are good!"]
	  wontfix:     ["L-wontfix",   "231f20", "Things not to fix"]

	create:
	  - ["T-perf",    "f51919", "Anything concerning perf"]
	  - ["T-feature", "078a00", "Anything that is adds new behavior"]
	  - ["T-doc",     "001c8a", "Anything that adds documentation"]
	  - ["T-hotfix",  "ff5100", "Critical things to deploy"]
	  - ["L-ready",   "ffe100", "Anything ready for merging and releasing"]

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
package gh

import (
	"flag"
	"fmt"
	"io"
	"os"

	"git.sr.ht/~psyomn/ecophagy/psy/common"

	"gopkg.in/yaml.v3"
)

func usage(fs *flag.FlagSet) error {
	fs.Usage()
	return ErrWrongUsage
}

type session struct {
	help           bool
	generateConfig bool
	listLabels     bool
	dryRun         bool
	poison         bool
}

// Run will run the subcommand with provided parameters
func Run(args common.RunParams) common.RunReturn {
	sess := session{}
	ghCmd := flag.NewFlagSet("gh", flag.ExitOnError)
	ghCmd.BoolVar(&sess.help, "help", false, "print help info")
	ghCmd.BoolVar(&sess.generateConfig, "generate-config", false, "print generic config into stdout")
	ghCmd.BoolVar(&sess.listLabels, "list-labels", false, "<owner> <repo> - list the labels on a github repository")
	ghCmd.BoolVar(&sess.dryRun, "dryrun", false, "set to dryrun poison")
	ghCmd.BoolVar(&sess.poison, "poison", false, "<config.yml> <owner> <repo> - run actual changes to labels")

	if err := ghCmd.Parse(args); err != nil {
		return err
	}

	ghToken, err := getGithubToken()
	if err != nil {
		return err
	}

	ghArgs := ghCmd.Args()
	switch {
	case sess.generateConfig:
		return generateConfig()
	case sess.listLabels:
		if len(ghCmd.Args()) < 2 {
			return usage(ghCmd)
		}
		owner := ghArgs[0]
		repo := ghArgs[1]
		return listLabels(owner, repo, ghToken)
	case sess.poison:
		if len(ghCmd.Args()) < 3 {
			return usage(ghCmd)
		}
		config := ghArgs[0]
		owner := ghArgs[1]
		repoName := ghArgs[2]
		return poison(sess.dryRun, config, owner, repoName, ghToken)
	case sess.help:
		return usage(ghCmd)
	default:
	}

	return usage(ghCmd)
}

func getGithubToken() (string, error) {
	// Maybe sourced as env variable for now, maybe something else
	// in the future

	if token := os.Getenv("GITHUB_TOKEN"); token != "" {
		return token, nil
	}

	return "", ErrNoGithubToken
}

func parseLabelConfig(path string) (*labelActions, error) {
	fh, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	contents, err := io.ReadAll(fh)
	if err != nil {
		return nil, err
	}

	var actions labelActions
	err = yaml.Unmarshal(contents, &actions)

	return &actions, err
}

func generateConfig() error {
	fmt.Println(`---
# Github label configuration.

# T = Type
# L = Lifecycle

rename:
  bug:         ["T-bug",       "a00000", "A software fault or failure"]
  duplicate:   ["L-duplicate", "16336d", "dOops!"]
  enhancement: ["T-upkeep",    "009a00", "Anything that rejuvenates"]
  invalid:     ["L-invalid",   "ffbb0e", "Anything that is invalid"]
  question:    ["T-question",  "ffaa99", "Questions are good!"]
  wontfix:     ["L-wontfix",   "231f20", "Things not to fix"]

create:
  - ["T-perf",    "f51919", "Anything concerning perf"]
  - ["T-feature", "078a00", "Anything that is adds new behavior"]
  - ["T-doc",     "001c8a", "Anything that adds documentation"]
  - ["T-hotfix",  "ff5100", "Critical things to deploy"]
  - ["L-ready",   "ffe100", "Anything ready for merging and releasing"]`)

	return nil
}

func listLabels(owner, repoName, token string) error {
	repo := repo{
		owner: owner,
		repo:  repoName,
		token: token,
	}

	return repo.getLabels()
}

func poison(dryRun bool, configFileName, owner, repoName, ghToken string) error {
	config, err := parseLabelConfig(configFileName)
	if err != nil {
		return err
	}

	repo := repo{
		owner:   owner,
		repo:    repoName,
		token:   ghToken,
		actions: config,
	}

	if err := repo.dryRun(dryRun).poison(); err != nil {
		return err
	}

	fmt.Println(repo.String())

	return nil
}
