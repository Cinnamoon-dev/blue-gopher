package services

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Cinnamoon-dev/blue-gopher/internal/customerrors"
	"github.com/Cinnamoon-dev/blue-gopher/internal/domain"
	"github.com/Cinnamoon-dev/blue-gopher/internal/repositories"
	"github.com/Cinnamoon-dev/blue-gopher/pkg/config"
)

type UserService struct {
	UserRepo repositories.UserRepository
	RoleRepo repositories.RoleRepository
}

func NewUserService(userRepo repositories.UserRepository, roleRepo repositories.RoleRepository) UserService {
	return UserService{UserRepo: userRepo, RoleRepo: roleRepo}
}

func (s *UserService) GetAll(ctx context.Context) ([]domain.User, error) {
	return s.UserRepo.GetAll(ctx)
}

func (s *UserService) Get(ctx context.Context, id int64) (*domain.User, error) {
	user, ok := ctx.Value(config.UserContextKey).(*domain.User)
	if ok && user != nil && id == user.ID {
		return user, nil
	}

	return s.UserRepo.Get(ctx, id)
}

func (s *UserService) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	return s.UserRepo.GetByEmail(ctx, email)
}

func (s *UserService) Create(ctx context.Context, newUser domain.User) (int64, error) {
	_, err := s.UserRepo.GetByEmail(ctx, newUser.Email)
	if err == nil {
		return 0, &customerrors.HTTPError{Status: http.StatusUnprocessableEntity, Message: fmt.Sprintf("Email %s already in use", newUser.Email)}
	}

	_, err = s.RoleRepo.Get(ctx, newUser.RoleID)
	if err != nil {
		return 0, &customerrors.HTTPError{Status: http.StatusUnprocessableEntity, Message: fmt.Sprintf("Role with id %d does not exists", newUser.RoleID)}
	}

	authService := NewAuthService()
	newUser.Password, err = authService.HashPassword(newUser.Password)
	if err != nil {
		return 0, &customerrors.HTTPError{Status: http.StatusInternalServerError, Message: fmt.Sprintf("Password hash: %s", err.Error())}
	}

	return s.UserRepo.Create(ctx, newUser)
}

func (s *UserService) Update(ctx context.Context, id int64, fields domain.User) error {
	if _, err := s.UserRepo.Get(ctx, id); err != nil {
		return &customerrors.HTTPError{Status: http.StatusNotFound, Message: fmt.Sprintf("User %d not found", id)}
	}

	authService := NewAuthService()
	hashPassword, err := authService.HashPassword(fields.Password)
	if err != nil {
		return &customerrors.HTTPError{Status: http.StatusInternalServerError, Message: fmt.Sprintf("Password hash: %s", err.Error())}
	}

	fields.Password = hashPassword
	return s.UserRepo.Update(ctx, id, fields)
}

func (s *UserService) Delete(ctx context.Context, id int64) error {
	return s.UserRepo.Delete(ctx, id)
}
