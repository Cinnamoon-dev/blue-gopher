package services

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Cinnamoon-dev/blue-gopher/internal/customerrors"
	"github.com/Cinnamoon-dev/blue-gopher/internal/domain"
	"github.com/Cinnamoon-dev/blue-gopher/internal/messaging/events"
	"github.com/Cinnamoon-dev/blue-gopher/internal/messaging/rabbitmq"
	"github.com/Cinnamoon-dev/blue-gopher/internal/repositories"
	"github.com/Cinnamoon-dev/blue-gopher/pkg/config"
)

type UserService struct {
	UserRepo  repositories.UserRepository
	RoleRepo  repositories.RoleRepository
	Publisher rabbitmq.RabbitPublisher
}

func NewUserService(
	userRepo repositories.UserRepository,
	roleRepo repositories.RoleRepository,
	publisher rabbitmq.RabbitPublisher,
) UserService {
	return UserService{
		UserRepo:  userRepo,
		RoleRepo:  roleRepo,
		Publisher: publisher,
	}
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
	newUser.Password, err = authService.Hash(newUser.Password)
	if err != nil {
		return 0, &customerrors.HTTPError{Status: http.StatusInternalServerError, Message: fmt.Sprintf("Password hash: %s", err.Error())}
	}

	userID, err := s.UserRepo.Create(ctx, newUser)
	if err != nil {
		return userID, err
	}

	event := events.EmailVerificationRequested{
		Event: events.Event{
			ID:        "1",
			Type:      "email.verification_requested",
			CreatedAt: time.Now(),
		},
		Email: newUser.Email,
	}

	err = s.Publisher.PublishEvent("", "email", event)
	if err != nil {
		log.Printf("[UserService.Create] Failed at sending email.verification_requested event: %s", err)
	}

	return userID, err
}

func (s *UserService) Update(ctx context.Context, id int64, fields domain.User) error {
	if _, err := s.UserRepo.Get(ctx, id); err != nil {
		return &customerrors.HTTPError{Status: http.StatusNotFound, Message: fmt.Sprintf("User %d not found", id)}
	}

	authService := NewAuthService()
	hashPassword, err := authService.Hash(fields.Password)
	if err != nil {
		return &customerrors.HTTPError{Status: http.StatusInternalServerError, Message: fmt.Sprintf("Password hash: %s", err.Error())}
	}

	fields.Password = hashPassword
	return s.UserRepo.Update(ctx, id, fields)
}

func (s *UserService) Delete(ctx context.Context, id int64) error {
	return s.UserRepo.Delete(ctx, id)
}
