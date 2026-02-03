package repositories

import (
	"database/sql"
	"net/http"

	"github.com/Cinnamoon-dev/blue-gopher/errors"
)

type UserRepository struct {
	db *sql.DB
}

type User struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func NewUserRepository(db *sql.DB) UserRepository {
	return UserRepository{db: db}
}

func (r *UserRepository) GetAll() ([]User, error) {
	rows, err := r.db.Query("SELECT id, email, password FROM usuarios ORDER BY id;")
	if err != nil {
		return nil, &errors.HTTPError{Message: "Database error", Status: http.StatusInternalServerError}
	}
	defer rows.Close()

	var data []User
	var user User

	for rows.Next() {
		err := rows.Scan(&user.ID, &user.Email, &user.Password)
		if err != nil {
			return nil, &errors.HTTPError{Message: "Database error", Status: http.StatusInternalServerError}
		}

		data = append(data, user)
	}

	return data, nil
}

func (r *UserRepository) Get(id int) (*User, error) {
	var user User

	row := r.db.QueryRow("SELECT id, email, password FROM usuarios WHERE id = ? ORDER BY id", id)
	if err := row.Scan(&user.ID, &user.Email, &user.Password); err != nil {
		return nil, &errors.HTTPError{Message: "Not found", Status: http.StatusNotFound}
	}

	return &user, nil
}

func (r *UserRepository) GetByEmail(name string) (*User, error) {
	var user User
	row := r.db.QueryRow("SELECT id, email, password FROM usuarios WHERE email = ? ORDER BY id LIMIT 1", name)
	if err := row.Scan(&user.ID, &user.Email, &user.Password); err != nil {
		return nil, &errors.HTTPError{Status: http.StatusNotFound, Message: "Not found"}
	}

	return &user, nil
}

func (r *UserRepository) Create(user User) (int64, error) {
	result, err := r.db.Exec("INSERT INTO usuarios(email, password) VALUES (?, ?)", user.Email, user.Password)
	if err != nil {
		return 0, &errors.HTTPError{Message: "Database error", Status: http.StatusInternalServerError}
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, &errors.HTTPError{Message: "Database error", Status: http.StatusInternalServerError}
	}

	return id, nil
}

func (r *UserRepository) Update(id int, user User) error {
	_, err := r.db.Exec("UPDATE usuarios SET email = ?, password = ? WHERE id = ?", user.Email, user.Password, id)
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
