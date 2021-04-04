package main

import (
	"errors"
)

var (
	ErrNewFileDetected = errors.New("mismatched expected csv hash")
	ErrScoreThreshold  = errors.New("rating threshold violated")
)
