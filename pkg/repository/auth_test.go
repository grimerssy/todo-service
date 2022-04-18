package repository

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/grimerssy/todo-service/internal/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAuthPostgres_CreateUser(t *testing.T) {
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

	r := NewAuthPostgres(db)

	tests := []struct {
		name      string
		mock      func(m sqlmock.Sqlmock)
		input     core.User
		want      uint
		errAssert assert.ErrorAssertionFunc
	}{
		{
			name: "ok",
			mock: func(m sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id"}).
					AddRow(id)
				m.ExpectQuery("INSERT INTO "+usersTable).
					WithArgs(firstName, lastName, email, username, password).
					WillReturnRows(rows)
			},
			input: core.User{
				FirstName: firstName,
				LastName:  lastName,
				Email:     email,
				Username:  username,
				Password:  password,
			},
			want:      id,
			errAssert: assert.NoError,
		},
		{
			name: "fail to insert",
			mock: func(m sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id"})
				m.ExpectQuery("INSERT INTO "+usersTable).
					WithArgs(firstName, lastName, email, username, password).
					WillReturnRows(rows)
			},
			input: core.User{
				FirstName: firstName,
				LastName:  lastName,
				Email:     email,
				Username:  username,
				Password:  password,
			},
			want:      0,
			errAssert: assert.Error,
		},
	}
	for _, tt := range tests {
		tt.mock(mock)
		got, err := r.CreateUser(context.Background(), tt.input)
		tt.errAssert(t, err)
		assert.Equal(t, tt.want, got)
		assert.NoError(t, mock.ExpectationsWereMet())
	}
}

func TestAuthPostgres_GetUserId(t *testing.T) {
	type input struct {
		username string
		password string
	}

	const (
		id       = 1
		username = "un"
		password = "pw"
		invalid  = ""
	)

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	r := NewAuthPostgres(db)

	tests := []struct {
		name      string
		mock      func(m sqlmock.Sqlmock)
		input     input
		want      uint
		errAssert assert.ErrorAssertionFunc
	}{
		{
			name: "ok",
			mock: func(m sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id"}).
					AddRow(id)
				m.ExpectQuery("SELECT id FROM "+usersTable).
					WithArgs(username, password).
					WillReturnRows(rows)
			},
			input: input{
				username: username,
				password: password,
			},
			want:      id,
			errAssert: assert.NoError,
		},
		{
			name: "invalid password",
			mock: func(m sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id"})
				m.ExpectQuery("SELECT id FROM "+usersTable).
					WithArgs(username, invalid).
					WillReturnRows(rows)
			},
			input: input{
				username: username,
				password: invalid,
			},
			want:      0,
			errAssert: assert.Error,
		},
	}
	for _, tt := range tests {
		tt.mock(mock)
		got, err := r.GetUserId(context.Background(), tt.input.username, tt.input.password)
		tt.errAssert(t, err)
		assert.Equal(t, tt.want, got)
		assert.NoError(t, mock.ExpectationsWereMet())
	}
}
