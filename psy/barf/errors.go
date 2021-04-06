package barf

import (
	"errors"
)

var (
	ErrNeedProjectName = errors.New("you need to provide a project name")
	ErrNeedOneArg      = errors.New("need to provide at least one argument")
	ErrNoSuchCommand   = errors.New("no such command")
)
