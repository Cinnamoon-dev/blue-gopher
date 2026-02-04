package repositories

import (
	"database/sql"

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
		return nil, err
	}

	return &role, nil
}
