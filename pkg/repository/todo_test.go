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
		want      uint
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
					WillReturnResult(sqlmock.NewResult(1, 1))

				m.ExpectCommit()
			},
			input: core.Todo{
				Title:       title,
				Description: description,
				Completed:   completed,
			},
			want:      id,
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
			want:      0,
			errAssert: assert.Error,
		},
	}
	for _, tt := range tests {
		tt.mock(mock)
		got, err := r.Create(context.Background(), id, tt.input)
		tt.errAssert(t, err)
		assert.Equal(t, tt.want, got)
		assert.NoError(t, mock.ExpectationsWereMet())
	}
}

func TestTodoPostgres_GetById(t *testing.T) {
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
		userId    uint
		todoId    uint
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
			userId: id,
			todoId: id,
			want: core.Todo{
				Id:          id,
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
			userId:    0,
			todoId:    id,
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
			userId:    id,
			todoId:    0,
			want:      core.Todo{},
			errAssert: assert.Error,
		},
	}
	for _, tt := range tests {
		tt.mock(mock)
		got, err := r.GetById(context.Background(), tt.userId, tt.todoId)
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
		userId    uint
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
			userId: id,
			input:  completed,
			want: []core.Todo{
				{
					Id:          id,
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
			userId:    0,
			input:     completed,
			want:      []core.Todo{},
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
			userId:    id,
			input:     completed,
			want:      []core.Todo{},
			errAssert: assert.NoError,
		},
	}
	for _, tt := range tests {
		tt.mock(mock)
		got, err := r.GetByCompletion(context.Background(), tt.userId, tt.input)
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
		userId    uint
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
			userId: id,
			want: []core.Todo{
				{
					Id:          id,
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
			userId:    0,
			want:      []core.Todo{},
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
			userId:    id,
			want:      []core.Todo{},
			errAssert: assert.NoError,
		},
	}
	for _, tt := range tests {
		tt.mock(mock)
		got, err := r.GetAll(context.Background(), tt.userId)
		tt.errAssert(t, err)
		assert.Equal(t, tt.want, got)
		assert.NoError(t, mock.ExpectationsWereMet())
	}
}

func TestTodoPostgres_Update(t *testing.T) {
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
		userId    uint
		todoId    uint
		input     core.Todo
		errAssert assert.ErrorAssertionFunc
	}{
		{
			name: "ok",
			mock: func(m sqlmock.Sqlmock) {
				m.ExpectExec("UPDATE "+todosTable).
					WithArgs(title, description, completed, id, id).
					WillReturnResult(sqlmock.NewResult(id, 1))
			},
			userId: id,
			todoId: id,
			input: core.Todo{
				Title:       title,
				Description: description,
				Completed:   completed,
			},
			errAssert: assert.NoError,
		},
		{
			name: "no user",
			mock: func(m sqlmock.Sqlmock) {
				m.ExpectExec("UPDATE "+todosTable).
					WithArgs(title, description, completed, 0, id).
					WillReturnResult(sqlmock.NewResult(0, 0))
			},
			userId: 0,
			todoId: id,
			input: core.Todo{
				Title:       title,
				Description: description,
				Completed:   completed,
			},
			errAssert: assert.NoError,
		},
		{
			name: "no todo",
			mock: func(m sqlmock.Sqlmock) {
				m.ExpectExec("UPDATE "+todosTable).
					WithArgs(title, description, completed, id, 0).
					WillReturnResult(sqlmock.NewResult(0, 0))
			},
			userId: id,
			todoId: 0,
			input: core.Todo{
				Title:       title,
				Description: description,
				Completed:   completed,
			},
			errAssert: assert.NoError,
		},
	}
	for _, tt := range tests {
		tt.mock(mock)
		err := r.Update(context.Background(), tt.userId, tt.todoId, tt.input)
		tt.errAssert(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	}
}

func TestTodoPostgres_Patch(t *testing.T) {
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
		userId    uint
		todoId    uint
		input     core.Todo
		errAssert assert.ErrorAssertionFunc
	}{
		{
			name: "ok all fields",
			mock: func(m sqlmock.Sqlmock) {
				m.ExpectExec("UPDATE "+todosTable).
					WithArgs(title, description, completed, id, id).
					WillReturnResult(sqlmock.NewResult(id, 1))
			},
			userId: id,
			todoId: id,
			input: core.Todo{
				Title:       title,
				Description: description,
				Completed:   completed,
			},
			errAssert: assert.NoError,
		},
		{
			name: "ok no title",
			mock: func(m sqlmock.Sqlmock) {
				m.ExpectExec("UPDATE "+todosTable).
					WithArgs(description, completed, id, id).
					WillReturnResult(sqlmock.NewResult(id, 1))
			},
			userId: id,
			todoId: id,
			input: core.Todo{
				Description: description,
				Completed:   completed,
			},
			errAssert: assert.NoError,
		},
		{
			name: "ok no description",
			mock: func(m sqlmock.Sqlmock) {
				m.ExpectExec("UPDATE "+todosTable).
					WithArgs(title, completed, id, id).
					WillReturnResult(sqlmock.NewResult(id, 1))
			},
			userId: id,
			todoId: id,
			input: core.Todo{
				Title:     title,
				Completed: completed,
			},
			errAssert: assert.NoError,
		},
		{
			name: "ok no completed",
			mock: func(m sqlmock.Sqlmock) {
				m.ExpectExec("UPDATE "+todosTable).
					WithArgs(title, description, false, id, id).
					WillReturnResult(sqlmock.NewResult(id, 1))
			},
			userId: id,
			todoId: id,
			input: core.Todo{
				Title:       title,
				Description: description,
			},
			errAssert: assert.NoError,
		},
		{
			name: "ok only title",
			mock: func(m sqlmock.Sqlmock) {
				m.ExpectExec("UPDATE "+todosTable).
					WithArgs(title, false, id, id).
					WillReturnResult(sqlmock.NewResult(id, 1))
			},
			userId: id,
			todoId: id,
			input: core.Todo{
				Title: title,
			},
			errAssert: assert.NoError,
		},
		{
			name: "ok only description",
			mock: func(m sqlmock.Sqlmock) {
				m.ExpectExec("UPDATE "+todosTable).
					WithArgs(description, false, id, id).
					WillReturnResult(sqlmock.NewResult(id, 1))
			},
			userId: id,
			todoId: id,
			input: core.Todo{
				Description: description,
			},
			errAssert: assert.NoError,
		},
		{
			name: "ok only completed",
			mock: func(m sqlmock.Sqlmock) {
				m.ExpectExec("UPDATE "+todosTable).
					WithArgs(completed, id, id).
					WillReturnResult(sqlmock.NewResult(id, 1))
			},
			userId: id,
			todoId: id,
			input: core.Todo{
				Completed: completed,
			},
			errAssert: assert.NoError,
		},
		{
			name: "ok empty",
			mock: func(m sqlmock.Sqlmock) {
				m.ExpectExec("UPDATE "+todosTable).
					WithArgs(false, id, id).
					WillReturnResult(sqlmock.NewResult(id, 1))
			},
			userId:    id,
			todoId:    id,
			input:     core.Todo{},
			errAssert: assert.NoError,
		},
		{
			name: "no user",
			mock: func(m sqlmock.Sqlmock) {
				m.ExpectExec("UPDATE "+todosTable).
					WithArgs(title, description, completed, 0, id).
					WillReturnResult(sqlmock.NewResult(0, 0))
			},
			userId: 0,
			todoId: id,
			input: core.Todo{
				Title:       title,
				Description: description,
				Completed:   completed,
			},
			errAssert: assert.NoError,
		},
		{
			name: "no todo",
			mock: func(m sqlmock.Sqlmock) {
				m.ExpectExec("UPDATE "+todosTable).
					WithArgs(title, description, completed, id, 0).
					WillReturnResult(sqlmock.NewResult(0, 0))
			},
			userId: id,
			todoId: 0,
			input: core.Todo{
				Title:       title,
				Description: description,
				Completed:   completed,
			},
			errAssert: assert.NoError,
		},
	}
	for _, tt := range tests {
		tt.mock(mock)
		err := r.Patch(context.Background(), tt.userId, tt.todoId, tt.input)
		tt.errAssert(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	}
}

func TestTodoPostgres_DeleteById(t *testing.T) {
	const id = 1

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	r := NewTodoPostgres(db)

	tests := []struct {
		name      string
		mock      func(m sqlmock.Sqlmock)
		userId    uint
		todoId    uint
		errAssert assert.ErrorAssertionFunc
	}{
		{
			name: "ok",
			mock: func(m sqlmock.Sqlmock) {
				m.ExpectExec("DELETE FROM "+todosTable).
					WithArgs(id, id).
					WillReturnResult(sqlmock.NewResult(id, 1))
			},
			userId:    id,
			todoId:    id,
			errAssert: assert.NoError,
		},
		{
			name: "no user",
			mock: func(m sqlmock.Sqlmock) {
				m.ExpectExec("DELETE FROM "+todosTable).
					WithArgs(0, id).
					WillReturnResult(sqlmock.NewResult(0, 0))
			},
			userId:    0,
			todoId:    id,
			errAssert: assert.NoError,
		},
		{
			name: "no todo",
			mock: func(m sqlmock.Sqlmock) {
				m.ExpectExec("DELETE FROM "+todosTable).
					WithArgs(id, 0).
					WillReturnResult(sqlmock.NewResult(0, 0))
			},
			userId:    id,
			todoId:    0,
			errAssert: assert.NoError,
		},
	}
	for _, tt := range tests {
		tt.mock(mock)
		err := r.DeleteById(context.Background(), tt.userId, tt.todoId)
		tt.errAssert(t, err)
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
		userId    uint
		input     bool
		errAssert assert.ErrorAssertionFunc
	}{
		{
			name: "ok",
			mock: func(m sqlmock.Sqlmock) {
				m.ExpectExec("DELETE FROM "+todosTable).
					WithArgs(id, completed).
					WillReturnResult(sqlmock.NewResult(id, 1))
			},
			userId:    id,
			input:     completed,
			errAssert: assert.NoError,
		},
		{
			name: "no user",
			mock: func(m sqlmock.Sqlmock) {
				m.ExpectExec("DELETE FROM "+todosTable).
					WithArgs(0, completed).
					WillReturnResult(sqlmock.NewResult(0, 0))
			},
			userId:    0,
			input:     completed,
			errAssert: assert.NoError,
		},
		{
			name: "no matches",
			mock: func(m sqlmock.Sqlmock) {
				m.ExpectExec("DELETE FROM "+todosTable).
					WithArgs(id, completed).
					WillReturnResult(sqlmock.NewResult(0, 0))
			},
			userId:    id,
			input:     completed,
			errAssert: assert.NoError,
		},
	}
	for _, tt := range tests {
		tt.mock(mock)
		err := r.DeleteByCompletion(context.Background(), tt.userId, tt.input)
		tt.errAssert(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	}
}
