package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"

	"github.com/Cinnamoon-dev/blue-gopher/internal/customerrors"
	"github.com/Cinnamoon-dev/blue-gopher/internal/domain"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return UserRepository{db: db}
}

func (r *UserRepository) GetPermission(ctx context.Context, id int64, action string, controller string) (bool, error) {
	rows, err := r.db.Query(`
		SELECT usuarios.id as u_id, rules.action as action, rules.permission as permission, controllers.name as controller
		FROM usuarios
		JOIN rules ON usuarios.role_id = rules.role_id
		JOIN controllers ON rules.controller_id = controllers.id
		WHERE
		u_id = ? AND
		action = ? AND
		controller = ?
		`,
		id,
		action,
		controller,
	)
	if err != nil {
		return false, &customerrors.HTTPError{Status: http.StatusInternalServerError, Message: err.Error()}
	}

	var output struct {
		UserID     int64
		Action     string
		Permission bool
		Controller string
	}

	if rows.Next() {
		err = rows.Scan(&output.UserID, &output.Action, &output.Permission, &output.Controller)
		if err != nil {
			return false, &customerrors.HTTPError{Status: http.StatusInternalServerError, Message: err.Error()}
		}
	} else {
		return false, &customerrors.HTTPError{Message: "Rule not found", Status: http.StatusInternalServerError}
	}

	return output.Permission, nil
}

func (r *UserRepository) GetAll(ctx context.Context) ([]domain.User, error) {
	rows, err := r.db.Query("SELECT id, email, password, is_verified, role_id FROM usuarios ORDER BY id;")
	if err != nil {
		return nil, &customerrors.HTTPError{Message: "get rows: " + err.Error(), Status: http.StatusInternalServerError}
	}
	defer rows.Close()

	var data []domain.User
	var user domain.User

	for rows.Next() {
		err := rows.Scan(&user.ID, &user.Email, &user.Password, &user.IsVerified, &user.RoleID)
		if err != nil {
			return nil, &customerrors.HTTPError{Message: "write user: " + err.Error(), Status: http.StatusInternalServerError}
		}

		data = append(data, user)
	}

	return data, nil
}

func (r *UserRepository) Get(ctx context.Context, id int64) (*domain.User, error) {
	var user domain.User

	row := r.db.QueryRow("SELECT id, email, password, is_verified, role_id FROM usuarios WHERE id = ? ORDER BY id", id)
	if err := row.Scan(&user.ID, &user.Email, &user.Password, &user.IsVerified, &user.RoleID); err != nil {
		return nil, &customerrors.HTTPError{Message: fmt.Sprintf("User %d not found", id), Status: http.StatusNotFound}
	}

	return &user, nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	var user domain.User
	row := r.db.QueryRow("SELECT id, email, password, is_verified, role_id FROM usuarios WHERE email = ? ORDER BY id LIMIT 1", email)
	if err := row.Scan(&user.ID, &user.Email, &user.Password, &user.IsVerified, &user.RoleID); err != nil {
		return nil, &customerrors.HTTPError{Status: http.StatusNotFound, Message: fmt.Sprintf("User %s not found", email)}
	}

	return &user, nil
}

func (r *UserRepository) Create(ctx context.Context, user domain.User) (int64, error) {
	result, err := r.db.Exec("INSERT INTO usuarios(email, password, is_verified, role_id) VALUES (?, ?, ?, ?)", user.Email, user.Password, user.IsVerified, user.RoleID)
	if err != nil {
		return 0, &customerrors.HTTPError{Message: "Database error", Status: http.StatusInternalServerError}
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, &customerrors.HTTPError{Message: "Database error", Status: http.StatusInternalServerError}
	}

	return id, nil
}

func (r *UserRepository) Update(ctx context.Context, id int64, user domain.User) error {
	_, err := r.db.Exec("UPDATE usuarios SET email = ?, password = ?, is_verified = ?, role_id = ? WHERE id = ?", user.Email, user.Password, user.IsVerified, user.RoleID, id)
	if err != nil {
		return &customerrors.HTTPError{Message: "Database error", Status: http.StatusInternalServerError}
	}

	return nil
}

func (r *UserRepository) Delete(ctx context.Context, id int64) error {
	_, err := r.db.Exec("DELETE FROM usuarios WHERE id = ?", id)
	if err != nil {
		return &customerrors.HTTPError{Message: "Database error", Status: http.StatusInternalServerError}
	}

	return nil
}
