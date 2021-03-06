package repository

import (
	"context"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/grimerssy/todo-service/internal/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserPostgres_Create(t *testing.T) {
	const (
		id        = 1
		firstName = "fn"
		lastName  = "ln"
		email     = "em"
		username  = "un"
		password  = "pw"
	)

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	r := NewUserPostgres(db)

	tests := []struct {
		name      string
		mock      func(m sqlmock.Sqlmock)
		input     core.User
		errAssert assert.ErrorAssertionFunc
	}{
		{
			name: "ok",
			mock: func(m sqlmock.Sqlmock) {
				m.ExpectExec("INSERT INTO "+usersTable).
					WithArgs(firstName, lastName, email, username, password).
					WillReturnResult(sqlmock.NewResult(id, 1))
			},
			input: core.User{
				FirstName: firstName,
				LastName:  lastName,
				Email:     email,
				Username:  username,
				Password:  password,
			},
			errAssert: assert.NoError,
		},
		{
			name: "fail to insert",
			mock: func(m sqlmock.Sqlmock) {
				m.ExpectExec("INSERT INTO "+usersTable).
					WithArgs(firstName, lastName, email, username, password).
					WillReturnResult(sqlmock.NewResult(0, 0)).
					WillReturnError(errors.New(""))
			},
			input: core.User{
				FirstName: firstName,
				LastName:  lastName,
				Email:     email,
				Username:  username,
				Password:  password,
			},
			errAssert: assert.Error,
		},
	}
	for _, tt := range tests {
		tt.mock(mock)
		err := r.Create(context.Background(), tt.input)
		tt.errAssert(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	}
}

func TestUserPostgres_GetCredentialsByUsername(t *testing.T) {
	const (
		id       = 1
		username = "un"
		password = "pw"
		invalid  = ""
	)

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	r := NewUserPostgres(db)

	tests := []struct {
		name      string
		mock      func(m sqlmock.Sqlmock)
		username  string
		want      core.User
		errAssert assert.ErrorAssertionFunc
	}{
		{
			name: "ok",
			mock: func(m sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "password"}).
					AddRow(id, password)
				m.ExpectQuery("SELECT id, password FROM " + usersTable).
					WithArgs(username).
					WillReturnRows(rows)
			},
			username: username,
			want: core.User{
				ID:       id,
				Username: username,
				Password: password,
			},
			errAssert: assert.NoError,
		},
		{
			name: "invalid username",
			mock: func(m sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "password"})
				m.ExpectQuery("SELECT id, password FROM " + usersTable).
					WithArgs(invalid).
					WillReturnRows(rows)
			},
			username:  invalid,
			want:      core.User{},
			errAssert: assert.Error,
		},
	}
	for _, tt := range tests {
		tt.mock(mock)
		got, err := r.GetCredentialsByUsername(context.Background(), tt.username)
		tt.errAssert(t, err)
		assert.Equal(t, tt.want, got)
		assert.NoError(t, mock.ExpectationsWereMet())
	}
}
