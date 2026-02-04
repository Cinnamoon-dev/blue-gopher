package repositories

import (
	"database/sql"
	"net/http"

	"github.com/Cinnamoon-dev/blue-gopher/internal/domain"
	"github.com/Cinnamoon-dev/blue-gopher/internal/errors"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return UserRepository{db: db}
}

func (r *UserRepository) GetPermission(id int, action string, controller string) (bool, error) {
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
		return false, err
	}

	var output struct {
		UserID     int
		Action     string
		Permission bool
		Controller string
	}

	if rows.Next() {
		err = rows.Scan(&output.UserID, &output.Action, &output.Permission, &output.Controller)
		if err != nil {
			return false, err
		}
	} else {
		return false, &errors.HTTPError{Message: "Rule not found", Status: http.StatusInternalServerError}
	}

	if output.Permission == false {
		return false, nil
	}

	return true, nil
}

func (r *UserRepository) GetAll() ([]domain.User, error) {
	rows, err := r.db.Query("SELECT id, email, password, role_id FROM usuarios ORDER BY id;")
	if err != nil {
		return nil, &errors.HTTPError{Message: "Database error", Status: http.StatusInternalServerError}
	}
	defer rows.Close()

	var data []domain.User
	var user domain.User

	for rows.Next() {
		err := rows.Scan(&user.ID, &user.Email, &user.Password, &user.RoleID)
		if err != nil {
			return nil, &errors.HTTPError{Message: "Database error", Status: http.StatusInternalServerError}
		}

		data = append(data, user)
	}

	return data, nil
}

func (r *UserRepository) Get(id int) (*domain.User, error) {
	var user domain.User

	row := r.db.QueryRow("SELECT id, email, password, role_id FROM usuarios WHERE id = ? ORDER BY id", id)
	if err := row.Scan(&user.ID, &user.Email, &user.Password, &user.RoleID); err != nil {
		return nil, &errors.HTTPError{Message: "Not found", Status: http.StatusNotFound}
	}

	return &user, nil
}

func (r *UserRepository) GetByEmail(name string) (*domain.User, error) {
	var user domain.User
	row := r.db.QueryRow("SELECT id, email, password, role_id FROM usuarios WHERE email = ? ORDER BY id LIMIT 1", name)
	if err := row.Scan(&user.ID, &user.Email, &user.Password, &user.RoleID); err != nil {
		return nil, &errors.HTTPError{Status: http.StatusNotFound, Message: "Not found"}
	}

	return &user, nil
}

func (r *UserRepository) Create(user domain.User) (int64, error) {
	result, err := r.db.Exec("INSERT INTO usuarios(email, password, role_id) VALUES (?, ?, ?)", user.Email, user.Password, user.RoleID)
	if err != nil {
		return 0, &errors.HTTPError{Message: "Database error", Status: http.StatusInternalServerError}
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, &errors.HTTPError{Message: "Database error", Status: http.StatusInternalServerError}
	}

	return id, nil
}

func (r *UserRepository) Update(id int, user domain.User) error {
	_, err := r.db.Exec("UPDATE usuarios SET email = ?, password = ?, role_id = ? WHERE id = ?", user.Email, user.Password, user.RoleID, id)
	if err != nil {
		return &errors.HTTPError{Message: "Database error", Status: http.StatusInternalServerError}
	}

	return nil
}

func (r *UserRepository) Delete(id int) error {
	_, err := r.db.Exec("DELETE FROM usuarios WHERE id = ?", id)
	if err != nil {
		return &errors.HTTPError{Message: "Database error", Status: http.StatusInternalServerError}
	}

	return nil
}
