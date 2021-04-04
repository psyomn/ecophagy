package main

import "errors"

type errorResponse struct {
	Error string `json:"error"`
}

var (
	ErrMethodNotSupported = errors.New("http method not supported")
	ErrMalformedURL       = errors.New("url is malformed")
	ErrBadAuthHeader      = errors.New("badauth: expected 'Authorization: token XXX' format")
	ErrBadBody            = errors.New("could not read body")
	ErrNeedLogin          = errors.New("login needed")
	ErrSmallPassword      = errors.New("password too small")
	ErrSmallUsername      = errors.New("username too small")
	ErrGenericLoginError  = errors.New("username/password don't match or don't exist")
	ErrUsernameTaken      = errors.New("username has been taken")
)
