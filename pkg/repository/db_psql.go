package repository

import (
	"database/sql"
	"fmt"
)

const (
	usersTable      = "users"
	todosTable      = "todos"
	usersTodosTable = "users_todos"
)

type ConfigPsql struct {
	Username string
	Password string
	Host     string
	Port     string
	DbName   string
	SSLMode  string
}

func NewDbPsql(cfg ConfigPsql) (*sql.DB, error) {
	db, err := sql.Open("postgres",
		fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=%s",
			cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.DbName, cfg.SSLMode))
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
