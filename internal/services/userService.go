package services

import (
	"fmt"

	"github.com/Cinnamoon-dev/blue-gopher/internal/domain"
	"github.com/Cinnamoon-dev/blue-gopher/internal/errors"
	"github.com/Cinnamoon-dev/blue-gopher/internal/repositories"
)

type UserService struct {
	UserRepo repositories.UserRepository
	RoleRepo repositories.RoleRepository
}

func NewUserService(userRepo repositories.UserRepository, roleRepo repositories.RoleRepository) UserService {
	return UserService{UserRepo: userRepo, RoleRepo: roleRepo}
}

func (s *UserService) GetAll() ([]domain.User, error) {
	return s.UserRepo.GetAll()
}

func (s *UserService) Get(id int) (*domain.User, error) {
	return s.UserRepo.Get(id)
}

func (s *UserService) GetByEmail(email string) (*domain.User, error) {
	return s.UserRepo.GetByEmail(email)
}

func (s *UserService) Create(newUser domain.User) (int64, error) {
	_, err := s.UserRepo.GetByEmail(newUser.Email)
	if err == nil {
		return 0, &errors.HTTPError{Status: 422, Message: fmt.Sprintf("Email %s already in use", newUser.Email)}
	}

	_, err = s.RoleRepo.Get(newUser.RoleID)
	if err == nil {
		return 0, &errors.HTTPError{Status: 422, Message: fmt.Sprintf("Role with id %d does not exists", newUser.RoleID)}
	}

	return s.UserRepo.Create(newUser)
}

func (s *UserService) Update(id int, fields domain.User) error {
	if _, err := s.UserRepo.Get(id); err != nil {
		return &errors.HTTPError{Status: 404, Message: fmt.Sprintf("User %d not found", id)}
	}

	return s.UserRepo.Update(id, fields)
}

func (s *UserService) Delete(id int) error {
	return s.UserRepo.Delete(id)
}
