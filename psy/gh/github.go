/*
Package gh is a simple package to interface with github

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
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

const githubBaseURL = "https://api.github.com/"

type githubLabelPatchBody struct {
	NewName     string `json:"new_name"`
	Color       string `json:"color"`
	Description string `json:"description"`
}

type githubCreateLabelBody struct {
	Name        string `json:"name"`
	Color       string `json:"color"`
	Description string `json:"description"`
}

type label struct {
	NodeID      string `json:"node_id"`
	URL         string `json:"url"`
	Name        string `json:"name"`
	Color       string `json:"color"`
	Default     bool   `json:"default"`
	Description string `json:"description"`
}

func (s *label) String() string {
	return fmt.Sprintf("[#%s %s]: %s (default: %v)",
		s.Color,
		s.Name,
		s.Description,
		s.Default)
}

type labelActions struct {
	Rename map[string][]string
	Create [][]string
}

type repo struct {
	owner  string
	repo   string
	labels []label

	token    string
	actions  *labelActions
	isDryRun bool
}

func (s *repo) create() {
	client := http.Client{}

	for _, arr := range s.actions.Create {
		createBodyStruct := githubCreateLabelBody{
			Name:        arr[0],
			Color:       arr[1],
			Description: arr[2],
		}

		fmt.Print("create ", createBodyStruct.Name, "... ")

		createLabelURL := fmt.Sprintf(
			"%srepos/%s/%s/labels",
			githubBaseURL, s.owner, s.repo,
		)

		bodyBytes, err := json.Marshal(&createBodyStruct)
		if err != nil {
			fmt.Println(err)
			continue
		}

		ctx := context.Background()
		req, err := http.NewRequestWithContext(
			ctx, http.MethodPost, createLabelURL, bytes.NewReader(bodyBytes))
		if err != nil {
			// TODO better error handling here
			fmt.Println(err)
			continue
		}

		req.Header["Accept"] = []string{"application/vnd.github.v3+json"}
		req.Header["Authorization"] = []string{fmt.Sprintf("token %s", s.token)}

		if s.isDryRun {
			fmt.Println("[SKIP]")
			continue
		}

		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("[ERROR]: ", err)
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusCreated {
			fmt.Println("[DONE]")
			continue
		}

		if resp.StatusCode == http.StatusUnprocessableEntity {
			fmt.Println("[EXISTS] (got a 422 http code, so the label probably already exists)")
			continue
		}

		fmt.Println("[ERROR]: ", resp)
	}
}

func (s *repo) getLabels() error {
	url := githubBaseURL + "repos/" + s.owner + "/" + s.repo + "/labels"
	client := &http.Client{}

	ctx := context.Background()
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		url,
		nil,
	)
	if err != nil {
		return err
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	return json.Unmarshal(bytes, &s.labels)
}

func (s *repo) dryRun(run bool) *repo {
	s.isDryRun = run
	return s
}

func (s *repo) poison() error {
	s.rename()
	s.create()

	err := s.getLabels()
	if err != nil {
		return err
	}

	return nil
}

// https://docs.github.com/en/rest/reference/issues#update-a-label
func (s *repo) rename() {
	client := http.Client{}

	for k, v := range s.actions.Rename {
		patchStruct := githubLabelPatchBody{
			NewName:     v[0],
			Color:       v[1],
			Description: v[2],
		}

		patchURL := fmt.Sprintf("%srepos/%s/%s/labels/%s",
			githubBaseURL,
			s.owner,
			s.repo,
			k,
		)

		bodyBytes, err := json.Marshal(&patchStruct)
		if err != nil {
			// TODO better error handling here
			fmt.Println(err)
			continue
		}

		fmt.Print("update ", k, " -> ", patchStruct.NewName, "... ")
		if s.isDryRun {
			fmt.Println(" [SKIP]")
			continue
		}

		ctx := context.Background()
		req, err := http.NewRequestWithContext(
			ctx,
			http.MethodPatch,
			patchURL,
			bytes.NewReader(bodyBytes),
		)

		if err != nil {
			// TODO better error handling here
			fmt.Println(err)
			continue
		}

		req.Header["Accept"] = []string{"application/vnd.github.v3+json"}
		req.Header["Authorization"] = []string{fmt.Sprintf("token %s", s.token)}

		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("[ERROR]: ", err, resp)
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusNotFound {
			fmt.Println("[NOT FOUND] (can't update a label if it doesn't exist)")
			continue
		}

		if resp.StatusCode != http.StatusOK {
			fmt.Println("[ERROR]:", resp)
			continue
		}

		fmt.Println("[OK]")
	}
}

func (s *repo) String() string {
	var build strings.Builder
	build.Grow(256)

	fmt.Fprintf(&build, "owner  : %s\n", s.owner)
	fmt.Fprintf(&build, "repo   : %s\n", s.repo)
	fmt.Fprintf(&build, "dryrun : %v\n", s.isDryRun)

	fmt.Fprintf(&build, "labels :\n")
	for _, label := range s.labels {
		fmt.Fprintf(&build, "  - %v\n", label.String())
	}

	return build.String()
}
