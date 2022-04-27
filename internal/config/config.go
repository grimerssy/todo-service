package config

import (
	"os"

	"github.com/grimerssy/todo-service/internal/handler"
	"github.com/grimerssy/todo-service/internal/server"
	"github.com/grimerssy/todo-service/pkg/auth"
	"github.com/grimerssy/todo-service/pkg/database"
	"github.com/grimerssy/todo-service/pkg/encoding"
	"github.com/grimerssy/todo-service/pkg/hashing"
	"github.com/grimerssy/todo-service/pkg/logging"
	"github.com/spf13/viper"
)

const (
	configPath = "configs"
)

type Config struct {
	Gin      handler.ConfigGin
	Server   server.ConfigServer
	Postgres database.ConfigPostgres
	JWT      auth.ConfigJWT
	Hashids  encoding.ConfigHashids
	Bcrypt   hashing.ConfigBcrypt
	Logrus   logging.ConfigLogrus
}

func NewConfig(environment string, logger logging.Logger) *Config {
	const (
		jwtSigningStringVar = "JWT"
		postgresPasswordVar = "POSTGRES"
	)

	var cfg *Config
	viper.AddConfigPath(configPath)
	viper.SetConfigName(environment)
	if err := viper.ReadInConfig(); err != nil {
		logger.Logf(logging.FatalLevel, "could not read configuration: %s", err.Error())
	}
	if err := viper.Unmarshal(&cfg); err != nil {
		logger.Logf(logging.FatalLevel, "could not unmarshal configuration: %s", err.Error())
	}

	cfg.Postgres.Password = os.Getenv(postgresPasswordVar)
	cfg.JWT.SigningString = os.Getenv(jwtSigningStringVar)

	return cfg
}
