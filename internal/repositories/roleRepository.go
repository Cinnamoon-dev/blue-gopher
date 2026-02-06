package repositories

import (
	"database/sql"
	"fmt"

	"github.com/Cinnamoon-dev/blue-gopher/internal/customerrors"
	"github.com/Cinnamoon-dev/blue-gopher/internal/domain"
)

type RoleRepository struct {
	DB *sql.DB
}

func NewRoleRepository(db *sql.DB) RoleRepository {
	return RoleRepository{DB: db}
}

func (r *RoleRepository) Get(id int64) (*domain.Role, error) {
	var role domain.Role
	row := r.DB.QueryRow("SELECT id, name FROM roles WHERE id = ?", id)

	err := row.Scan(&role.ID, &role.Name)
	if err != nil {
		return nil, &customerrors.HTTPError{Status: 404, Message: fmt.Sprintf("Role with id %d not found", id)}
	}

	return &role, nil
}
