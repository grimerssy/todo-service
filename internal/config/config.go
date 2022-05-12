package config

import (
	"github.com/grimerssy/todo-service/internal/handler"
	"github.com/grimerssy/todo-service/internal/server"
	"github.com/grimerssy/todo-service/pkg/auth"
	"github.com/grimerssy/todo-service/pkg/cache"
	"github.com/grimerssy/todo-service/pkg/database"
	"github.com/grimerssy/todo-service/pkg/encoding"
	"github.com/grimerssy/todo-service/pkg/hashing"
	"github.com/grimerssy/todo-service/pkg/logging"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

const (
	configPath = "configs"
)

type Config struct {
	Gin      handler.ConfigGin
	Server   server.ConfigServer
	Postgres database.ConfigPostgres
	LFU      cache.ConfigLFU
	JWT      auth.ConfigJWT
	Hashids  encoding.ConfigHashids
	Bcrypt   hashing.ConfigBcrypt
	Logrus   logging.ConfigLogrus
}

func NewConfig(environment string, logger logging.Logger) *Config {
	var cfg *Config

	viper.AddConfigPath(configPath)
	viper.SetConfigName(environment)

	if err := viper.ReadInConfig(); err != nil {
		logger.Logf(logging.FatalLevel, "could not read configuration: %s", err.Error())
	}
	if err := viper.Unmarshal(&cfg); err != nil {
		logger.Logf(logging.FatalLevel, "could not unmarshal configuration: %s", err.Error())
	}

	if err := cfg.ApplyEnvVariables(); err != nil {
		logger.Logf(logging.FatalLevel, "could not read .env file: %s", err.Error())
	}

	return cfg
}

func (c *Config) ApplyEnvVariables() error {
	const (
		jwtSigningString = "JWT_SIGNING_STRING"
		postgresUser     = "POSTGRES_USER"
		postgresPassword = "POSTGRES_PASSWORD"
	)

	env, err := godotenv.Read()
	if err != nil {
		return err
	}

	c.JWT.SigningString = env[jwtSigningString]

	c.Postgres.Username = env[postgresUser]
	c.Postgres.Password = env[postgresPassword]

	return nil
}
