package notes

import "errors"

var (
	ErrBadRegisterCLI = errors.New("not enough segments for USERNAME:PASSWORD format")
	ErrBadSize        = errors.New("argument had bad size")
	ErrGenerateToken  = errors.New("error generating token")
	ErrPublicBadPath  = errors.New("bad public resource path")
	ErrNotesBadPath   = errors.New("bad note resource path")
)
