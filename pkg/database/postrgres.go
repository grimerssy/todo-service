package database

import (
	"database/sql"
	"fmt"
)

type ConfigPostgres struct {
	Username string
	Password string
	Host     string
	Port     string
	DbName   string
	SSLMode  string
}

func NewPostgres(cfg ConfigPostgres) (*sql.DB, error) {
	db, err := sql.Open("postgres",
		fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=%s",
			cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.DbName, cfg.SSLMode))
	if err != nil {
		return nil, err
	}

	return db, db.Ping()
}
