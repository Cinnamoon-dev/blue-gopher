package domain

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

type User struct {
	ID       int64
	Email    string
	Password string
	RoleID   int64
}

func (u *User) ValidateEmail() error {
	email := u.Email
	email = strings.TrimSpace(email)

	if email == "" {
		return errors.New("Validate email: email is required")
	}

	emailRegex := `^[a-zA-Z0-9.!#$%&'*+/=?^_` + "`" + `{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$`
	re := regexp.MustCompile(emailRegex)

	if re.MatchString(email) {
		return nil
	}

	return fmt.Errorf("Validate email: invalid email %s", email)
}

func (u *User) ValidatePassword() error {
	if len(u.Password) < 4 {
		return errors.New("Validate password: password is too short (minimum 4 characters)")
	}

	return nil
}
