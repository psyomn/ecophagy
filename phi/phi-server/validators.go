package main

import "errors"

func validatePassword(pass string) error {
	if len(pass) < minPasswordLength {
		return errors.New("problem registering user with small password")
	}
	return nil
}

func validateUsername(user string) error {
	if len(user) < minUsernameLength {
		return errors.New("problem registering user with small username")
	}
	return nil
}
