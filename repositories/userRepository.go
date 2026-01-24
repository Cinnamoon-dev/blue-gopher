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
	ID    int    `json:"id"`
	Nome  string `json:"nome"`
	Idade int    `json:"idade"`
}

func NewUserRepository(db *sql.DB) UserRepository {
	return UserRepository{db: db}
}

// TODO
// Create custom errors to return a message and a http status code

func (r *UserRepository) GetAll() ([]User, error) {
	rows, err := r.db.Query("SELECT id, nome, idade FROM usuarios;")
	if err != nil {
		return nil, &errors.HTTPError{Message: "Database error", Status: http.StatusInternalServerError}
	}
	defer rows.Close()

	var data []User
	var user User

	for rows.Next() {
		err := rows.Scan(&user.ID, &user.Nome, &user.Idade)
		if err != nil {
			return nil, &errors.HTTPError{Message: "Database error", Status: http.StatusInternalServerError}
		}

		data = append(data, user)
	}

	return data, nil
}

func (r *UserRepository) Get(id int) (*User, error) {
	var user User

	row := r.db.QueryRow("SELECT id, nome, idade FROM usuarios WHERE id = ?", id)
	if err := row.Scan(&user.ID, &user.Nome, &user.Idade); err != nil {
		return nil, &errors.HTTPError{Message: "Not found", Status: http.StatusNotFound}
	}

	return &user, nil
}

func (r *UserRepository) GetByName(name string) (*User, error) {
	var user User
	row := r.db.QueryRow("SELECT id, nome, idade FROM usuarios WHERE nome = ?", name)
	if err := row.Scan(&user.ID, &user.Nome, &user.Idade); err != nil {
		return nil, &errors.HTTPError{Status: http.StatusNotFound, Message: "Not found"}
	}

	return &user, nil
}

func (r *UserRepository) Create(user User) (int64, error) {
	result, err := r.db.Exec("INSERT INTO usuarios(nome, idade) VALUES (?, ?)", user.Nome, user.Idade)
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
	_, err := r.db.Exec("UPDATE usuarios SET nome = ?, idade = ? WHERE id = ?", user.Nome, user.Idade, id)
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

	return err
}
