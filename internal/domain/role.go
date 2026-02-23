package domain

import (
	"errors"
	"strings"
)

type Role struct {
	ID   int64
	Name string
}

func (r *Role) ValidateName() error {
	r.Name = strings.TrimSpace(r.Name)

	if r.Name == "" {
		return errors.New("Validate name: name is required")
	}

	return nil
}
