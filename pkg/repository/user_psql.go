package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/grimerssy/todo-service/internal/core"
)

type UserPsql struct {
	db *sql.DB
}

func NewUserPsql(db *sql.DB) *UserPsql {
	return &UserPsql{db: db}
}

func (r *UserPsql) Create(ctx context.Context, user core.User) error {
	query := fmt.Sprintf(`
INSERT INTO %s (first_name, last_name, email, username, password) 
VALUES ($1, $2, $3, $4, $5); 
`, usersTable)

	_, err := r.db.ExecContext(ctx, query,
		user.FirstName, user.LastName, user.Email, user.Username, user.Password)

	return err
}

func (r *UserPsql) GetCredentialsByUsername(ctx context.Context, username string) (core.UserCredentials, error) {
	query := fmt.Sprintf(`
SELECT id, password FROM %s 
WHERE username = $1
LIMIT 1;
`, usersTable)

	var cred core.UserCredentials
	row := r.db.QueryRowContext(ctx, query, username)
	err := row.Scan(&cred.ID, &cred.Password)

	cred.Username = username

	return cred, err
}