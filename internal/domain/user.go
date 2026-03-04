package domain

import (
	"errors"
	"fmt"
	"net/mail"
	"strings"
)

type User struct {
	ID       int64
	Email    string
	Password string
	RoleID   int64
}

func NewUser(id int64, email string, password string, roleID int64) (*User, error) {
	user := &User{
		ID:       id,
		Email:    email,
		Password: password,
		RoleID:   roleID,
	}

	err := user.ValidateEmail()
	if err != nil {
		return nil, err
	}

	err = user.ValidatePassword()
	if err != nil {
		return nil, err
	}

	return user, err
}

func (u *User) ValidateEmail() error {
	email := u.Email
	email = strings.TrimSpace(email)

	if email == "" {
		return errors.New("Validate email: email is required")
	}

	if _, err := mail.ParseAddress(email); err != nil {
		return fmt.Errorf("Validate email: invalid email %s", email)
	}

	return nil
}

func (u *User) ValidatePassword() error {
	if len(u.Password) < 4 {
		return errors.New("Validate password: password is too short (minimum 4 characters)")
	}

	if len(u.Password) > 72 { // bcrypt hash only uses until 72 characters
		return errors.New("Validate password: password too long (maximum 72 characters)")
	}

	return nil
}
