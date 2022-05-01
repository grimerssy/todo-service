package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/grimerssy/todo-service/internal/core"
)

type UserPostgres struct {
	db *sql.DB
}

func NewUserPostgres(db *sql.DB) *UserPostgres {
	return &UserPostgres{
		db: db,
	}
}

func (r *UserPostgres) Create(ctx context.Context, user core.User) error {
	query := fmt.Sprintf(`
INSERT INTO %s (first_name, last_name, email, username, password) 
VALUES ($1, $2, $3, $4, $5); 
`, usersTable)

	_, err := r.db.ExecContext(ctx, query,
		user.FirstName, user.LastName, user.Email, user.Username, user.Password)

	if err != nil {
		return fmt.Errorf("could not execute query: %s", err.Error())
	}

	return nil
}

func (r *UserPostgres) GetCredentialsByUsername(ctx context.Context, username string) (core.User, error) {
	query := fmt.Sprintf(`
SELECT id, password FROM %s 
WHERE username = $1
LIMIT 1;
`, usersTable)

	var user core.User
	row := r.db.QueryRowContext(ctx, query, username)
	if err := row.Scan(&user.ID, &user.Password); err != nil {
		return user, fmt.Errorf("could not scan row: %s", err.Error())
	}

	user.Username = username

	return user, nil
}
