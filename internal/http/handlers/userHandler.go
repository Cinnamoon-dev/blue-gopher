package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/Cinnamoon-dev/blue-gopher/internal/customerrors"
	"github.com/Cinnamoon-dev/blue-gopher/internal/domain"
	"github.com/Cinnamoon-dev/blue-gopher/internal/services"
)

type UserHandler struct {
	Svc services.UserService
}

type UserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	RoleID   int64  `json:"role_id"`
}

func NewUserHandler(svc services.UserService) UserHandler {
	return UserHandler{Svc: svc}
}

// The idea is simple: each handler is going to have its own parseID
// So each handler can parse the URL the way they want
func parseID(path string) (int64, error) {
	// path = /user/{id}
	parts := strings.Split(path, "/")

	if len(parts) < 3 {
		return 0, http.ErrNotSupported
	}

	id, err := strconv.ParseInt(parts[2], 10, 64)
	if err != nil {
		return 0, &customerrors.HTTPError{Status: http.StatusUnprocessableEntity, Message: fmt.Sprintf("Invalid id %s", parts[2])}
	}

	return id, nil
}

func (h *UserHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	users, err := h.Svc.GetAll(ctx)
	if err != nil {
		RespondError(w, err)
		return
	}

	RespondJSON(w, http.StatusOK, users)
}

func (h *UserHandler) GetOneUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id, err := parseID(r.URL.Path)
	if err != nil {
		RespondError(w, err)
		return
	}

	user, err := h.Svc.Get(ctx, id)
	if err != nil {
		RespondError(w, err)
		return
	}

	RespondJSON(w, http.StatusOK, user)
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var newUser UserRequest
	json.NewDecoder(r.Body).Decode(&newUser)
	newUser.Email = strings.TrimSpace(newUser.Email)

	user := domain.User{
		ID:       0,
		Email:    newUser.Email,
		Password: newUser.Password,
		RoleID:   newUser.RoleID,
	}

	if err := user.ValidateEmail(); err != nil {
		RespondError(w, err)
		return
	}

	if err := user.ValidatePassword(); err != nil {
		RespondError(w, err)
		return
	}

	id, err := h.Svc.Create(ctx, user)
	if err != nil {
		RespondError(w, err)
		return
	}

	RespondJSON(w, http.StatusOK, map[string]string{"message": fmt.Sprintf("User %d created successfully", id)})
}

func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id, err := parseID(r.URL.Path)
	if err != nil {
		RespondError(w, err)
		return
	}

	var fields UserRequest
	if err := json.NewDecoder(r.Body).Decode(&fields); err != nil {
		RespondError(w, err)
		return
	}

	user := domain.User{
		ID:       0,
		Email:    fields.Email,
		Password: fields.Password,
		RoleID:   fields.RoleID,
	}

	if err := user.ValidateEmail(); err != nil {
		RespondError(w, err)
		return
	}

	if err := user.ValidatePassword(); err != nil {
		RespondError(w, err)
		return
	}

	if err := h.Svc.Update(ctx, id, user); err != nil {
		RespondError(w, err)
		return
	}

	RespondJSON(w, http.StatusOK, fields)
}

func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id, err := parseID(r.URL.Path)
	if err != nil {
		RespondError(w, err)
		return
	}

	if err := h.Svc.Delete(ctx, id); err != nil {
		RespondError(w, err)
		return
	}

	RespondJSON(w, http.StatusOK, map[string]string{"message": fmt.Sprintf("User %d deleted successfully", id)})
}
