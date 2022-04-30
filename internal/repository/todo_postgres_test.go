package repository

import (
	"context"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/grimerssy/todo-service/internal/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTodoPostgres_Create(t *testing.T) {
	const (
		id          = 1
		title       = "t"
		description = "d"
		completed   = true
	)

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	r := NewTodoPostgres(db)

	tests := []struct {
		name      string
		mock      func(m sqlmock.Sqlmock)
		input     core.Todo
		errAssert assert.ErrorAssertionFunc
	}{
		{
			name: "ok",
			mock: func(m sqlmock.Sqlmock) {
				m.ExpectBegin()

				rows := sqlmock.NewRows([]string{"id"}).
					AddRow(id)
				m.ExpectQuery("INSERT INTO "+todosTable).
					WithArgs(title, description, completed).
					WillReturnRows(rows)

				m.ExpectExec("INSERT INTO "+usersTodosTable).
					WithArgs(id, id).
					WillReturnResult(sqlmock.NewResult(id, 1))

				m.ExpectCommit()
			},
			input: core.Todo{
				Title:       title,
				Description: description,
				Completed:   completed,
			},
			errAssert: assert.NoError,
		},
		{
			name: "fail to insert",
			mock: func(m sqlmock.Sqlmock) {
				m.ExpectBegin()

				rows := sqlmock.NewRows([]string{"id"})
				m.ExpectQuery("INSERT INTO "+todosTable).
					WithArgs(title, description, completed).
					WillReturnRows(rows)

				m.ExpectRollback()
			},
			input: core.Todo{
				Title:       title,
				Description: description,
				Completed:   completed,
			},
			errAssert: assert.Error,
		},
	}
	for _, tt := range tests {
		tt.mock(mock)
		err := r.Create(context.Background(), id, tt.input)
		tt.errAssert(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	}
}

func TestTodoPostgres_GetByID(t *testing.T) {
	const (
		id          = 1
		title       = "t"
		description = "d"
		completed   = true
	)
	now := time.Now()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	r := NewTodoPostgres(db)

	tests := []struct {
		name      string
		mock      func(m sqlmock.Sqlmock)
		userID    uint
		todoID    uint
		want      core.Todo
		errAssert assert.ErrorAssertionFunc
	}{
		{
			name: "ok",
			mock: func(m sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					"id", "title", "description", "completed", "created_at", "updated_at"}).
					AddRow(id, title, description, completed, now, now)
				m.ExpectQuery("SELECT td.id, td.title, td.description, td.completed, td.created_at, td.updated_at FROM "+todosTable).
					WithArgs(id, id).
					WillReturnRows(rows)
			},
			userID: id,
			todoID: id,
			want: core.Todo{
				ID:          id,
				Title:       title,
				Description: description,
				Completed:   completed,
				CreatedAt:   now,
				UpdatedAt:   now,
			},
			errAssert: assert.NoError,
		},
		{
			name: "no user",
			mock: func(m sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					"id", "title", "description", "completed", "created_at", "updated_at"})
				m.ExpectQuery("SELECT td.id, td.title, td.description, td.completed, td.created_at, td.updated_at FROM "+todosTable).
					WithArgs(0, id).
					WillReturnRows(rows)
			},
			userID:    0,
			todoID:    id,
			want:      core.Todo{},
			errAssert: assert.Error,
		},
		{
			name: "no todo",
			mock: func(m sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					"id", "title", "description", "completed", "created_at", "updated_at"})
				m.ExpectQuery("SELECT td.id, td.title, td.description, td.completed, td.created_at, td.updated_at FROM "+todosTable).
					WithArgs(id, 0).
					WillReturnRows(rows)
			},
			userID:    id,
			todoID:    0,
			want:      core.Todo{},
			errAssert: assert.Error,
		},
	}
	for _, tt := range tests {
		tt.mock(mock)
		got, err := r.GetByID(context.Background(), tt.userID, tt.todoID)
		tt.errAssert(t, err)
		assert.Equal(t, tt.want, got)
		assert.NoError(t, mock.ExpectationsWereMet())
	}
}

func TestTodoPostgres_GetByCompletion(t *testing.T) {
	const (
		id          = 1
		title       = "t"
		description = "d"
		completed   = true
	)
	now := time.Now()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	r := NewTodoPostgres(db)

	tests := []struct {
		name      string
		mock      func(m sqlmock.Sqlmock)
		userID    uint
		input     bool
		want      []core.Todo
		errAssert assert.ErrorAssertionFunc
	}{
		{
			name: "ok",
			mock: func(m sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					"id", "title", "description", "completed", "created_at", "updated_at"}).
					AddRow(id, title, description, completed, now, now)
				m.ExpectQuery("SELECT td.id, td.title, td.description, td.completed, td.created_at, td.updated_at FROM "+todosTable).
					WithArgs(id, completed).
					WillReturnRows(rows)
			},
			userID: id,
			input:  completed,
			want: []core.Todo{
				{
					ID:          id,
					Title:       title,
					Description: description,
					Completed:   completed,
					CreatedAt:   now,
					UpdatedAt:   now,
				},
			},
			errAssert: assert.NoError,
		},
		{
			name: "no user",
			mock: func(m sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					"id", "title", "description", "completed", "created_at", "updated_at"})
				m.ExpectQuery("SELECT td.id, td.title, td.description, td.completed, td.created_at, td.updated_at FROM "+todosTable).
					WithArgs(0, completed).
					WillReturnRows(rows)
			},
			userID:    0,
			input:     completed,
			want:      nil,
			errAssert: assert.NoError,
		},
		{
			name: "no results",
			mock: func(m sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					"id", "title", "description", "completed", "created_at", "updated_at"})
				m.ExpectQuery("SELECT td.id, td.title, td.description, td.completed, td.created_at, td.updated_at FROM "+todosTable).
					WithArgs(id, completed).
					WillReturnRows(rows)
			},
			userID:    id,
			input:     completed,
			want:      nil,
			errAssert: assert.NoError,
		},
	}
	for _, tt := range tests {
		tt.mock(mock)
		got, err := r.GetByCompletion(context.Background(), tt.userID, tt.input)
		tt.errAssert(t, err)
		assert.Equal(t, tt.want, got)
		assert.NoError(t, mock.ExpectationsWereMet())
	}
}

func TestTodoPostgres_GetAll(t *testing.T) {
	const (
		id          = 1
		title       = "t"
		description = "d"
		completed   = true
	)
	now := time.Now()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	r := NewTodoPostgres(db)

	tests := []struct {
		name      string
		mock      func(m sqlmock.Sqlmock)
		userID    uint
		want      []core.Todo
		errAssert assert.ErrorAssertionFunc
	}{
		{
			name: "ok",
			mock: func(m sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					"id", "title", "description", "completed", "created_at", "updated_at"}).
					AddRow(id, title, description, completed, now, now)
				m.ExpectQuery("SELECT td.id, td.title, td.description, td.completed, td.created_at, td.updated_at FROM " + todosTable).
					WithArgs(id).
					WillReturnRows(rows)
			},
			userID: id,
			want: []core.Todo{
				{
					ID:          id,
					Title:       title,
					Description: description,
					Completed:   completed,
					CreatedAt:   now,
					UpdatedAt:   now,
				},
			},
			errAssert: assert.NoError,
		},
		{
			name: "no user",
			mock: func(m sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					"id", "title", "description", "completed", "created_at", "updated_at"})
				m.ExpectQuery("SELECT td.id, td.title, td.description, td.completed, td.created_at, td.updated_at FROM " + todosTable).
					WithArgs(0).
					WillReturnRows(rows)
			},
			userID:    0,
			want:      nil,
			errAssert: assert.NoError,
		},
		{
			name: "no results",
			mock: func(m sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					"id", "title", "description", "completed", "created_at", "updated_at"})
				m.ExpectQuery("SELECT td.id, td.title, td.description, td.completed, td.created_at, td.updated_at FROM " + todosTable).
					WithArgs(id).
					WillReturnRows(rows)
			},
			userID:    id,
			want:      nil,
			errAssert: assert.NoError,
		},
	}
	for _, tt := range tests {
		tt.mock(mock)
		got, err := r.GetAll(context.Background(), tt.userID)
		tt.errAssert(t, err)
		assert.Equal(t, tt.want, got)
		assert.NoError(t, mock.ExpectationsWereMet())
	}
}

func TestTodoPostgres_UpdateByID(t *testing.T) {
	const (
		id          = 1
		title       = "t"
		description = "d"
		completed   = true
	)

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	r := NewTodoPostgres(db)

	tests := []struct {
		name      string
		mock      func(m sqlmock.Sqlmock)
		userID    uint
		todoID    uint
		input     core.Todo
		want      uint
		errAssert assert.ErrorAssertionFunc
	}{
		{
			name: "ok",
			mock: func(m sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id"}).
					AddRow(id)
				m.ExpectQuery("UPDATE "+todosTable).
					WithArgs(title, description, completed, id, id).
					WillReturnRows(rows)
			},
			userID: id,
			todoID: id,
			input: core.Todo{
				Title:       title,
				Description: description,
				Completed:   completed,
			},
			want:      id,
			errAssert: assert.NoError,
		},
		{
			name: "no user",
			mock: func(m sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id"})
				m.ExpectQuery("UPDATE "+todosTable).
					WithArgs(title, description, completed, 0, id).
					WillReturnRows(rows)
			},
			userID: 0,
			todoID: id,
			input: core.Todo{
				Title:       title,
				Description: description,
				Completed:   completed,
			},
			want:      0,
			errAssert: assert.Error,
		},
		{
			name: "no todo",
			mock: func(m sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id"})
				m.ExpectQuery("UPDATE "+todosTable).
					WithArgs(title, description, completed, id, 0).
					WillReturnRows(rows)
			},
			userID: id,
			todoID: 0,
			input: core.Todo{
				Title:       title,
				Description: description,
				Completed:   completed,
			},
			want:      0,
			errAssert: assert.Error,
		},
	}
	for _, tt := range tests {
		tt.mock(mock)
		got, err := r.UpdateByID(context.Background(), tt.userID, tt.todoID, tt.input)
		tt.errAssert(t, err)
		assert.Equal(t, tt.want, got)
		assert.NoError(t, mock.ExpectationsWereMet())
	}
}

func TestTodoPostgres_PatchID(t *testing.T) {
	const (
		id          = 1
		title       = "t"
		description = "d"
		completed   = true
	)

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	r := NewTodoPostgres(db)

	tests := []struct {
		name      string
		mock      func(m sqlmock.Sqlmock)
		userID    uint
		todoID    uint
		input     core.Todo
		want      uint
		errAssert assert.ErrorAssertionFunc
	}{
		{
			name: "ok all fields",
			mock: func(m sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id"}).
					AddRow(id)
				m.ExpectQuery("UPDATE "+todosTable).
					WithArgs(title, description, completed, id, id).
					WillReturnRows(rows)
			},
			userID: id,
			todoID: id,
			input: core.Todo{
				Title:       title,
				Description: description,
				Completed:   completed,
			},
			want:      id,
			errAssert: assert.NoError,
		},
		{
			name: "ok no title",
			mock: func(m sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id"}).
					AddRow(id)
				m.ExpectQuery("UPDATE "+todosTable).
					WithArgs(description, completed, id, id).
					WillReturnRows(rows)
			},
			userID: id,
			todoID: id,
			input: core.Todo{
				Description: description,
				Completed:   completed,
			},
			want:      id,
			errAssert: assert.NoError,
		},
		{
			name: "ok no description",
			mock: func(m sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id"}).
					AddRow(id)
				m.ExpectQuery("UPDATE "+todosTable).
					WithArgs(title, completed, id, id).
					WillReturnRows(rows)
			},
			userID: id,
			todoID: id,
			input: core.Todo{
				Title:     title,
				Completed: completed,
			},
			want:      id,
			errAssert: assert.NoError,
		},
		{
			name: "ok no completed",
			mock: func(m sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id"}).
					AddRow(id)
				m.ExpectQuery("UPDATE "+todosTable).
					WithArgs(title, description, false, id, id).
					WillReturnRows(rows)
			},
			userID: id,
			todoID: id,
			input: core.Todo{
				Title:       title,
				Description: description,
			},
			want:      id,
			errAssert: assert.NoError,
		},
		{
			name: "ok only title",
			mock: func(m sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id"}).
					AddRow(id)
				m.ExpectQuery("UPDATE "+todosTable).
					WithArgs(title, false, id, id).
					WillReturnRows(rows)
			},
			userID: id,
			todoID: id,
			input: core.Todo{
				Title: title,
			},
			want:      id,
			errAssert: assert.NoError,
		},
		{
			name: "ok only description",
			mock: func(m sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id"}).
					AddRow(id)
				m.ExpectQuery("UPDATE "+todosTable).
					WithArgs(description, false, id, id).
					WillReturnRows(rows)
			},
			userID: id,
			todoID: id,
			input: core.Todo{
				Description: description,
			},
			want:      id,
			errAssert: assert.NoError,
		},
		{
			name: "ok only completed",
			mock: func(m sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id"}).
					AddRow(id)
				m.ExpectQuery("UPDATE "+todosTable).
					WithArgs(completed, id, id).
					WillReturnRows(rows)
			},
			userID: id,
			todoID: id,
			input: core.Todo{
				Completed: completed,
			},
			want:      id,
			errAssert: assert.NoError,
		},
		{
			name: "ok empty",
			mock: func(m sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id"}).
					AddRow(id)
				m.ExpectQuery("UPDATE "+todosTable).
					WithArgs(false, id, id).
					WillReturnRows(rows)
			},
			userID:    id,
			todoID:    id,
			input:     core.Todo{},
			want:      id,
			errAssert: assert.NoError,
		},
		{
			name: "no user",
			mock: func(m sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id"})
				m.ExpectQuery("UPDATE "+todosTable).
					WithArgs(title, description, completed, 0, id).
					WillReturnRows(rows)
			},
			userID: 0,
			todoID: id,
			input: core.Todo{
				Title:       title,
				Description: description,
				Completed:   completed,
			},
			want:      0,
			errAssert: assert.Error,
		},
		{
			name: "no todo",
			mock: func(m sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id"})
				m.ExpectQuery("UPDATE "+todosTable).
					WithArgs(title, description, completed, id, 0).
					WillReturnRows(rows)
			},
			userID: id,
			todoID: 0,
			input: core.Todo{
				Title:       title,
				Description: description,
				Completed:   completed,
			},
			want:      0,
			errAssert: assert.Error,
		},
	}
	for _, tt := range tests {
		tt.mock(mock)
		got, err := r.PatchByID(context.Background(), tt.userID, tt.todoID, tt.input)
		tt.errAssert(t, err)
		assert.Equal(t, tt.want, got)
		assert.NoError(t, mock.ExpectationsWereMet())
	}
}

func TestTodoPostgres_DeleteByID(t *testing.T) {
	const id = 1

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	r := NewTodoPostgres(db)

	tests := []struct {
		name      string
		mock      func(m sqlmock.Sqlmock)
		userID    uint
		todoID    uint
		want      uint
		errAssert assert.ErrorAssertionFunc
	}{
		{
			name: "ok",
			mock: func(m sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id"}).
					AddRow(id)
				m.ExpectQuery("DELETE FROM "+todosTable).
					WithArgs(id, id).
					WillReturnRows(rows)
			},
			userID:    id,
			todoID:    id,
			want:      id,
			errAssert: assert.NoError,
		},
		{
			name: "no user",
			mock: func(m sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id"})
				m.ExpectQuery("DELETE FROM "+todosTable).
					WithArgs(0, id).
					WillReturnRows(rows)
			},
			userID:    0,
			todoID:    id,
			want:      0,
			errAssert: assert.Error,
		},
		{
			name: "no todo",
			mock: func(m sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id"})
				m.ExpectQuery("DELETE FROM "+todosTable).
					WithArgs(id, 0).
					WillReturnRows(rows)
			},
			userID:    id,
			todoID:    0,
			want:      0,
			errAssert: assert.Error,
		},
	}
	for _, tt := range tests {
		tt.mock(mock)
		got, err := r.DeleteByID(context.Background(), tt.userID, tt.todoID)
		tt.errAssert(t, err)
		assert.Equal(t, tt.want, got)
		assert.NoError(t, mock.ExpectationsWereMet())
	}
}

func TestTodoPostgres_DeleteByCompletion(t *testing.T) {
	const (
		id        = 1
		completed = true
	)

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	r := NewTodoPostgres(db)

	tests := []struct {
		name      string
		mock      func(m sqlmock.Sqlmock)
		userID    uint
		input     bool
		want      []uint
		errAssert assert.ErrorAssertionFunc
	}{
		{
			name: "ok",
			mock: func(m sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id"}).
					AddRow(id)
				m.ExpectQuery("DELETE FROM "+todosTable).
					WithArgs(id, completed).
					WillReturnRows(rows)
			},
			userID:    id,
			input:     completed,
			want:      []uint{id},
			errAssert: assert.NoError,
		},
		{
			name: "no user",
			mock: func(m sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id"})
				m.ExpectQuery("DELETE FROM "+todosTable).
					WithArgs(0, completed).
					WillReturnRows(rows)
			},
			userID:    0,
			input:     completed,
			want:      nil,
			errAssert: assert.NoError,
		},
		{
			name: "no matches",
			mock: func(m sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id"})
				m.ExpectQuery("DELETE FROM "+todosTable).
					WithArgs(id, completed).
					WillReturnRows(rows)
			},
			userID:    id,
			input:     completed,
			want:      nil,
			errAssert: assert.NoError,
		},
	}
	for _, tt := range tests {
		tt.mock(mock)
		got, err := r.DeleteByCompletion(context.Background(), tt.userID, tt.input)
		tt.errAssert(t, err)
		assert.Equal(t, tt.want, got)
		assert.NoError(t, mock.ExpectationsWereMet())
	}
}
