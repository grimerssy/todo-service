package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/grimerssy/todo-service/internal/core"
)

type TodoPostgres struct {
	db *sql.DB
}

func NewTodoPostgres(db *sql.DB) *TodoPostgres {
	return &TodoPostgres{db: db}
}

func (r *TodoPostgres) Create(ctx context.Context, userId uint, todo core.Todo) (uint, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	query := fmt.Sprintf(`
INSERT INTO %s (title, description, completed) 
VALUES ($1, $2, $3) 
RETURNING id;
`, todosTable)

	var todoId uint
	row := tx.QueryRowContext(ctx, query, todo.Title, todo.Description, todo.Completed)
	if err := row.Scan(&todoId); err != nil {
		return 0, err
	}

	query = fmt.Sprintf(`
INSERT INTO %s (user_id, todo_id)
VALUES ($1, $2);
`, usersTodosTable)

	if _, err := tx.ExecContext(ctx, query, userId, todoId); err != nil {
		return 0, err
	}

	return todoId, tx.Commit()
}

func (r *TodoPostgres) GetById(ctx context.Context, userId uint, todoId uint) (core.Todo, error) {
	query := fmt.Sprintf(`
SELECT td.id, td.title, td.description, td.completed, td.created_at, td.updated_at
FROM %s td
INNER JOIN (SELECT todo_id FROM %s WHERE user_id = $1) AS ut
ON ut.todo_id = td.id
WHERE td.id = $2
LIMIT 1;
`, todosTable, usersTodosTable)

	var todo core.Todo
	row := r.db.QueryRowContext(ctx, query, userId, todoId)
	err := row.Scan(
		&todo.Id, &todo.Title, &todo.Description, &todo.Completed, &todo.CreatedAt, &todo.UpdatedAt)

	return todo, err
}

func (r *TodoPostgres) GetByCompletion(ctx context.Context, userId uint, completed bool) ([]core.Todo, error) {
	query := fmt.Sprintf(`
SELECT td.id, td.title, td.description, td.completed, td.created_at, td.updated_at
FROM %s td
INNER JOIN (SELECT todo_id FROM %s WHERE user_id = $1) AS ut
ON ut.todo_id = td.id
WHERE td.completed = $2
`, todosTable, usersTodosTable)

	todos := []core.Todo{}
	rows, err := r.db.QueryContext(ctx, query, userId, completed)
	if err != nil {
		return todos, err
	}

	for rows.Next() {
		var todo core.Todo
		err := rows.Scan(
			&todo.Id, &todo.Title, &todo.Description, &todo.Completed, &todo.CreatedAt, &todo.UpdatedAt)
		if err != nil {
			return todos, err
		}
		todos = append(todos, todo)
	}

	return todos, rows.Err()
}

func (r *TodoPostgres) GetAll(ctx context.Context, userId uint) ([]core.Todo, error) {
	query := fmt.Sprintf(`
SELECT td.id, td.title, td.description, td.completed, td.created_at, td.updated_at
FROM %s td
INNER JOIN (SELECT todo_id FROM %s WHERE user_id = $1) AS ut
ON ut.todo_id = td.id;
`, todosTable, usersTodosTable)

	todos := []core.Todo{}
	rows, err := r.db.QueryContext(ctx, query, userId)
	if err != nil {
		return todos, err
	}

	for rows.Next() {
		var todo core.Todo
		err := rows.Scan(
			&todo.Id, &todo.Title, &todo.Description, &todo.Completed, &todo.CreatedAt, &todo.UpdatedAt)
		if err != nil {
			return todos, err
		}
		todos = append(todos, todo)
	}

	return todos, rows.Err()
}

func (r *TodoPostgres) Update(ctx context.Context, userId uint, todoId uint, todo core.Todo) error {
	query := fmt.Sprintf(`
UPDATE %s td
SET title = $1,
    description = $2,
    completed = $3
FROM %s ut
WHERE ut.user_id = $4
    AND ut.todo_id = td.id
    AND td.id = $5
`, todosTable, usersTodosTable)

	_, err := r.db.ExecContext(ctx, query, todo.Title, todo.Description, todo.Completed, userId, todoId)
	return err
}

func (r *TodoPostgres) Patch(ctx context.Context, userId uint, todoId uint, todo core.Todo) error {
	setStatements := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1

	if len(todo.Title) != 0 {
		setStatements = append(setStatements, fmt.Sprintf("title = $%d", argId))
		args = append(args, todo.Title)
		argId++
	}
	if len(todo.Description) != 0 {
		setStatements = append(setStatements, fmt.Sprintf("description = $%d", argId))
		args = append(args, todo.Description)
		argId++
	}
	setStatements = append(setStatements, fmt.Sprintf("completed = $%d", argId))
	args = append(args, todo.Completed)
	argId++

	setQuery := strings.Join(setStatements, ", ")

	query := fmt.Sprintf(`
UPDATE %s td
%s
FROM %s ut
WHERE ut.user_id = $%d
    AND ut.todo_id = td.id
    AND td.id = $%d
`, todosTable, setQuery, usersTodosTable, argId, argId+1)

	args = append(args, userId, todoId)

	_, err := r.db.ExecContext(ctx, query, args...)
	return err
}

func (r *TodoPostgres) DeleteById(ctx context.Context, userId uint, todoId uint) error {
	query := fmt.Sprintf(`
DELETE FROM %s td
USING %s ut
WHERE ut.user_id = $1
    AND ut.todo_id = td.id
    AND td.id = $2;
`, todosTable, usersTodosTable)

	_, err := r.db.ExecContext(ctx, query, userId, todoId)
	return err
}

func (r *TodoPostgres) DeleteByCompletion(ctx context.Context, userId uint, completed bool) error {
	query := fmt.Sprintf(`
DELETE FROM %s td
USING %s ut
WHERE ut.user_id = $1
    AND ut.todo_id = td.id
    AND td.completed = $2;
`, todosTable, usersTodosTable)

	_, err := r.db.ExecContext(ctx, query, userId, completed)
	return err
}
