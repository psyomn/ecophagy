package main

func validatePassword(pass string) error {
	if len(pass) < minPasswordLength {
		return ErrSmallPassword
	}
	return nil
}

func validateUsername(user string) error {
	if len(user) < minUsernameLength {
		return ErrSmallUsername
	}
	return nil
}
