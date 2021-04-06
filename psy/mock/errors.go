package mock

import (
	"errors"
)

var (
	ErrWrongCmdUsage  = errors.New("wrong command usage for mock")
	ErrUnknownService = errors.New("unknown type of service to create")
)
