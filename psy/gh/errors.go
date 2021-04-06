package gh

import "errors"

var (
	ErrNoGithubToken = errors.New("no github token set in environment variable GITHUB_TOKEN")
	ErrWrongUsage    = errors.New("wrong gh usage")
)
