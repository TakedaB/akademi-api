package repository

import (
	"database/sql"

	"github.com/TakedaB/akademi-api/internal/model"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(tx *sql.Tx, u *model.User) error {
	query := `
		INSERT INTO users (name, email, password_hash, role)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, updated_at`

	return tx.QueryRow(
		query,
		u.Name, u.Email, u.PasswordHash, u.Role,
	).Scan(&u.ID, &u.CreatedAt, &u.UpdatedAt)
}

func (r *UserRepository) FindByEmail(email string) (*model.User, error) {
	query := `SELECT id, name, email, password_hash, role, created_at, updated_at FROM users WHERE email = $1`

	var u model.User
	err := r.db.QueryRow(query, email).Scan(&u.ID, &u.Name, &u.Email, &u.PasswordHash, &u.Role, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &u, nil
}

func (r *UserRepository) BeginTx() (*sql.Tx, error) {
	return r.db.Begin()
}

func (r *UserRepository) FindByID(id string) (*model.User, error) {
	query := `SELECT id, name, email, password_hash, role, created_at, updated_at FROM users WHERE id = $1`

	var u model.User
	err := r.db.QueryRow(query, id).Scan(&u.ID, &u.Name, &u.Email, &u.PasswordHash, &u.Role, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &u, nil
}
