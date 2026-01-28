package repository

import (
	"context"
	"database/sql"
	"errors"

	"log-management-backend/internal/models"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

// GetUserByUsername - ค้นหา user จาก username
func (r *UserRepository) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	query := `
		SELECT id, username, password_hash, role, tenant_id, created_at
		FROM users
		WHERE username = $1
	`

	var user models.User
	err := r.db.QueryRowContext(ctx, query, username).Scan(
		&user.ID,
		&user.Username,
		&user.PasswordHash,
		&user.Role,
		&user.TenantID,
		&user.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New("user not found")
	}
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// CreateUser - สร้าง user ใหม่
func (r *UserRepository) CreateUser(ctx context.Context, user models.User) error {
	query := `
		INSERT INTO users (username, password_hash, role, tenant_id)
		VALUES ($1, $2, $3, $4)
	`

	_, err := r.db.ExecContext(ctx, query, user.Username, user.PasswordHash, user.Role, user.TenantID)
	return err
}
