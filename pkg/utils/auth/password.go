package auth

import (
	"unicode"

	"github.com/pkg/errors"
)

func ValidatePassword(password string) error {
	length, special, number, upper, lower := false, false, false, false, false

	if len(password) >= 8 && len(password) <= 32 {
		length = true
	}

	for _, c := range password {
		if special && number && upper && lower && length {
			return nil
		}
		switch {
		case unicode.IsUpper(c):
			upper = true
			continue
		case unicode.IsLower(c):
			lower = true
			continue
		case unicode.IsNumber(c):
			number = true
			continue
		case unicode.IsPunct(c) || unicode.IsSymbol(c):
			special = true
			continue
		}
	}

	err := errors.New("password does not meet criteria")

	if !length {
		err = errors.Wrap(err, "[incorrect password length]")
	}
	if !special {
		err = errors.Wrap(err, "[must contain a special character]")
	}
	if !number {
		err = errors.Wrap(err, "[must contain a number]")
	}
	if !lower {
		err = errors.Wrap(err, "[must contain a lowercase letter]")
	}
	if !upper {
		err = errors.Wrap(err, "[must contain an uppercase letter]")
	}

	return err
}
