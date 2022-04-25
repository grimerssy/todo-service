package config

import (
	"log"
	"os"
	"time"

	"github.com/grimerssy/todo-service/internal/server"
	"github.com/grimerssy/todo-service/pkg/handler"
	"github.com/grimerssy/todo-service/pkg/repository"
	"github.com/grimerssy/todo-service/pkg/service"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

const (
	configPath = "configs"
)

type Config struct {
	LogFormatting   string
	RequestSeconds  time.Duration
	ShutdownSeconds time.Duration
	Server          server.ConfigServer
	Postgres        repository.ConfigPsql
	JWT             service.ConfigJWT
	Hashids         service.ConfigHashids
	Bcrypt          service.ConfigBcrypt
}

func GetLogger(logFormatting, environment string) *logrus.Logger {
	var logLevel logrus.Level
	var logFormatter logrus.Formatter

	switch environment {
	case "dev":
		logLevel = logrus.DebugLevel
	default:
		logLevel = logrus.ErrorLevel
	}

	switch logFormatting {
	case "JSON":
		logFormatter = &logrus.JSONFormatter{}
	default:
		logFormatter = &logrus.TextFormatter{}
	}

	logger := logrus.New()
	logger.Level = logLevel
	logger.Formatter = logFormatter

	return logger
}

func GetConfig(environment string) *Config {
	const (
		postgresPasswordKey = "PSQL_TODO"
		jwtSecretKey        = "JWT_SECRET"
	)

	var cfg *Config
	viper.AddConfigPath(configPath)
	viper.SetConfigName(environment)
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("an error occured while reading configuration: %s", err.Error())
	}
	if err := viper.Unmarshal(&cfg); err != nil {
		log.Fatalf("an error occured while unmarshalling configuration: %s", err.Error())
	}

	cfg.Postgres.Password = os.Getenv(postgresPasswordKey)
	cfg.JWT.Secret = os.Getenv(jwtSecretKey)

	return cfg
}

func GetRepositories(cfg *Config) (*repository.Repositories, func() error) {
	dbPsql, err := repository.NewDbPsql(cfg.Postgres)
	if err != nil {
		log.Fatalf("an error occured while connecting to postgres: %s", err.Error())
	}
	userRepository := repository.NewUserPsql(dbPsql)
	todoRepository := repository.NewTodoPsql(dbPsql)

	closeDB := func() error {
		return dbPsql.Close()
	}

	return &repository.Repositories{
		UserRepository: userRepository,
		TodoRepository: todoRepository,
	}, closeDB
}

func GetServices(cfg *Config, repositories *repository.Repositories) *service.Services {
	hasher := service.NewHashBcrypt(cfg.Bcrypt)

	userEncoder, err := service.NewEncoderHashids(cfg.Hashids, service.UserKey)
	if err != nil {
		log.Fatalf("an error occured while initializing user encoder: %s", err.Error())
	}
	todoEncoder, err := service.NewEncoderHashids(cfg.Hashids, service.TodoKey)
	if err != nil {
		log.Fatalf("an error occured while initializing todo encoder: %s", err.Error())
	}

	userService := service.NewUserEncoded(hasher, userEncoder, repositories.UserRepository)
	todoService := service.NewTodoEncoded(userEncoder, todoEncoder, repositories.TodoRepository)
	authService := service.NewAuthJWT(cfg.JWT, userService)

	return &service.Services{
		AuthService: authService,
		UserService: userService,
		TodoService: todoService,
	}
}

func GetGinHandlers(cfg *Config, logger logrus.FieldLogger, services *service.Services) *handler.HandlersGin {
	requestTimeout := cfg.RequestSeconds * time.Second

	authGin := handler.NewAuthGin(logger, services.AuthService, services.UserService, requestTimeout)
	middlewareGin := handler.NewMiddlewareGin(logger, services.AuthService, requestTimeout)
	todoGin := handler.NewTodoGin(logger, services.TodoService, requestTimeout)

	return &handler.HandlersGin{
		Auth:       authGin,
		Middleware: middlewareGin,
		Todo:       todoGin,
	}
}
