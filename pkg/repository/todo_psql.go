package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/grimerssy/todo-service/internal/core"
)

type TodoPsql struct {
	db *sql.DB
}

func NewTodoPsql(db *sql.DB) *TodoPsql {
	return &TodoPsql{db: db}
}

func (r *TodoPsql) Create(ctx context.Context, userID uint, todo core.Todo) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
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
		return err
	}

	query = fmt.Sprintf(`
INSERT INTO %s (user_id, todo_id)
VALUES ($1, $2);
`, usersTodosTable)

	if _, err := tx.ExecContext(ctx, query, userID, todoId); err != nil {
		return err
	}

	return tx.Commit()
}

func (r *TodoPsql) GetByID(ctx context.Context, userID uint, todoID uint) (core.Todo, error) {
	query := fmt.Sprintf(`
SELECT td.id, td.title, td.description, td.completed, td.created_at, td.updated_at
FROM %s td
INNER JOIN (SELECT todo_id FROM %s WHERE user_id = $1) AS ut
ON ut.todo_id = td.id
WHERE td.id = $2
LIMIT 1;
`, todosTable, usersTodosTable)

	var todo core.Todo
	row := r.db.QueryRowContext(ctx, query, userID, todoID)
	err := row.Scan(
		&todo.ID, &todo.Title, &todo.Description, &todo.Completed, &todo.CreatedAt, &todo.UpdatedAt)

	return todo, err
}

func (r *TodoPsql) GetByCompletion(ctx context.Context, userID uint, completed bool) ([]core.Todo, error) {
	query := fmt.Sprintf(`
SELECT td.id, td.title, td.description, td.completed, td.created_at, td.updated_at
FROM %s td
INNER JOIN (SELECT todo_id FROM %s WHERE user_id = $1) AS ut
ON ut.todo_id = td.id
WHERE td.completed = $2
`, todosTable, usersTodosTable)

	todos := []core.Todo{}
	rows, err := r.db.QueryContext(ctx, query, userID, completed)
	if err != nil {
		return todos, err
	}

	for rows.Next() {
		var todo core.Todo
		err := rows.Scan(
			&todo.ID, &todo.Title, &todo.Description, &todo.Completed, &todo.CreatedAt, &todo.UpdatedAt)
		if err != nil {
			return todos, err
		}
		todos = append(todos, todo)
	}

	return todos, rows.Err()
}

func (r *TodoPsql) GetAll(ctx context.Context, userID uint) ([]core.Todo, error) {
	query := fmt.Sprintf(`
SELECT td.id, td.title, td.description, td.completed, td.created_at, td.updated_at
FROM %s td
INNER JOIN (SELECT todo_id FROM %s WHERE user_id = $1) AS ut
ON ut.todo_id = td.id;
`, todosTable, usersTodosTable)

	todos := []core.Todo{}
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return todos, err
	}

	for rows.Next() {
		var todo core.Todo
		err := rows.Scan(
			&todo.ID, &todo.Title, &todo.Description, &todo.Completed, &todo.CreatedAt, &todo.UpdatedAt)
		if err != nil {
			return todos, err
		}
		todos = append(todos, todo)
	}

	return todos, rows.Err()
}

func (r *TodoPsql) UpdateByID(ctx context.Context, userID uint, todoID uint, todo core.Todo) error {
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

	_, err := r.db.ExecContext(ctx, query, todo.Title, todo.Description, todo.Completed, userID, todoID)
	return err
}

func (r *TodoPsql) PatchByID(ctx context.Context, userID uint, todoID uint, todo core.Todo) error {
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
SET %s
FROM %s ut
WHERE ut.user_id = $%d
    AND ut.todo_id = td.id
    AND td.id = $%d
`, todosTable, setQuery, usersTodosTable, argId, argId+1)

	args = append(args, userID, todoID)

	_, err := r.db.ExecContext(ctx, query, args...)
	return err
}

func (r *TodoPsql) DeleteByID(ctx context.Context, userID uint, todoID uint) error {
	query := fmt.Sprintf(`
DELETE FROM %s td
USING %s ut
WHERE ut.user_id = $1
    AND ut.todo_id = td.id
    AND td.id = $2;
`, todosTable, usersTodosTable)

	_, err := r.db.ExecContext(ctx, query, userID, todoID)
	return err
}

func (r *TodoPsql) DeleteByCompletion(ctx context.Context, userID uint, completed bool) error {
	query := fmt.Sprintf(`
DELETE FROM %s td
USING %s ut
WHERE ut.user_id = $1
    AND ut.todo_id = td.id
    AND td.completed = $2;
`, todosTable, usersTodosTable)

	_, err := r.db.ExecContext(ctx, query, userID, completed)
	return err
}
