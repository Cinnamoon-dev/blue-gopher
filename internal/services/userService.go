package services

import (
	"fmt"

	"github.com/Cinnamoon-dev/blue-gopher/internal/domain"
	"github.com/Cinnamoon-dev/blue-gopher/internal/errors"
	"github.com/Cinnamoon-dev/blue-gopher/internal/repositories"
)

type UserService struct {
	Repo repositories.UserRepository
}

func NewUserService(Repo repositories.UserRepository) UserService {
	return UserService{Repo: Repo}
}

func (s *UserService) GetAll() ([]domain.User, error) {
	return s.Repo.GetAll()
}

func (s *UserService) Get(id int) (*domain.User, error) {
	return s.Repo.Get(id)
}

func (s *UserService) GetByEmail(email string) (*domain.User, error) {
	return s.Repo.GetByEmail(email)
}

func (s *UserService) Create(newUser domain.User) (int64, error) {
	_, err := s.Repo.GetByEmail(newUser.Email)
	if err == nil {
		return 0, &errors.HTTPError{Status: 422, Message: fmt.Sprintf("Email %s already in use", newUser.Email)}
	}

	return s.Repo.Create(newUser)
}

func (s *UserService) Update(id int, fields domain.User) error {
	if _, err := s.Repo.Get(id); err != nil {
		return &errors.HTTPError{Status: 404, Message: fmt.Sprintf("User %d not found", id)}
	}

	return s.Repo.Update(id, fields)
}

func (s *UserService) Delete(id int) error {
	return s.Repo.Delete(id)
}
