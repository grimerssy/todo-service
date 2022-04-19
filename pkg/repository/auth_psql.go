package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/grimerssy/todo-service/internal/core"
)

type AuthPsql struct {
	db *sql.DB
}

func NewAuthPsql(db *sql.DB) *AuthPsql {
	return &AuthPsql{db: db}
}

func (r *AuthPsql) CreateUser(ctx context.Context, user core.User) (uint, error) {
	query := fmt.Sprintf(`
INSERT INTO %s (first_name, last_name, email, username, password) 
VALUES ($1, $2, $3, $4, $5) 
RETURNING id;
`, usersTable)

	var id uint
	row := r.db.QueryRowContext(ctx, query,
		user.FirstName, user.LastName, user.Email, user.Username, user.Password)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func (r *AuthPsql) GetUserId(ctx context.Context, username string, password string) (uint, error) {
	query := fmt.Sprintf(`
SELECT id FROM %s 
WHERE username = $1 
	AND password = $2;
`, usersTable)

	var id uint
	row := r.db.QueryRowContext(ctx, query, username, password)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}
