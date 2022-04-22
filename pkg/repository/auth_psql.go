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

func (r *AuthPsql) CreateUser(ctx context.Context, user core.User) error {
	query := fmt.Sprintf(`
INSERT INTO %s (first_name, last_name, email, username, password) 
VALUES ($1, $2, $3, $4, $5); 
`, usersTable)

	_, err := r.db.ExecContext(ctx, query,
		user.FirstName, user.LastName, user.Email, user.Username, user.Password)

	return err
}

func (r *AuthPsql) GetUserAuth(ctx context.Context, username string) (core.UserAuth, error) {
	query := fmt.Sprintf(`
SELECT id, password FROM %s 
WHERE username = $1
LIMIT 1;
`, usersTable)

	var auth core.UserAuth
	row := r.db.QueryRowContext(ctx, query, username)
	err := row.Scan(&auth.ID, &auth.Password)

	auth.Username = username

	return auth, err
}
