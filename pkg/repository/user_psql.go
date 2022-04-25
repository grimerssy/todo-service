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
	return &UserPsql{
		db: db,
	}
}

func (r *UserPsql) Create(ctx context.Context, user core.User) error {
	res := make(chan error, 1)

	go func() {
		query := fmt.Sprintf(`
INSERT INTO %s (first_name, last_name, email, username, password) 
VALUES ($1, $2, $3, $4, $5); 
`, usersTable)

		_, err := r.db.ExecContext(ctx, query,
			user.FirstName, user.LastName, user.Email, user.Username, user.Password)

		if err != nil {
			res <- fmt.Errorf("could not execute query: %s", err.Error())
			return
		}

		res <- nil
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-res:
		return err

	}
}

func (r *UserPsql) GetCredentialsByUsername(ctx context.Context, username string) (core.UserCredentials, error) {
	res := make(chan func() (core.UserCredentials, error), 1)

	go func() {
		query := fmt.Sprintf(`
SELECT id, password FROM %s 
WHERE username = $1
LIMIT 1;
`, usersTable)

		var cred core.UserCredentials
		row := r.db.QueryRowContext(ctx, query, username)
		if err := row.Scan(&cred.ID, &cred.Password); err != nil {
			res <- func() (core.UserCredentials, error) {
				return cred, fmt.Errorf("could not scan row: %s", err.Error())
			}
			return
		}

		cred.Username = username

		res <- func() (core.UserCredentials, error) {
			return cred, nil
		}
	}()

	select {
	case <-ctx.Done():
		return core.UserCredentials{}, ctx.Err()
	case f := <-res:
		return f()
	}
}
