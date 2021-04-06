package memo

import (
	"errors"
)

var (
	ErrNonExistFile  = errors.New("fool! you can't memo what does not exist")
	ErrCantFindEntry = errors.New("could not find entry")
)
