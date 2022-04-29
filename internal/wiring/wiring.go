package wiring

import (
	"github.com/grimerssy/todo-service/internal/config"
	"github.com/grimerssy/todo-service/internal/handler"
	"github.com/grimerssy/todo-service/internal/repository"
	"github.com/grimerssy/todo-service/internal/service"
	"github.com/grimerssy/todo-service/pkg/auth"
	"github.com/grimerssy/todo-service/pkg/cache"
	"github.com/grimerssy/todo-service/pkg/database"
	"github.com/grimerssy/todo-service/pkg/encoding"
	"github.com/grimerssy/todo-service/pkg/hashing"
	"github.com/grimerssy/todo-service/pkg/logging"
)

func GetRepositories(cfg *config.Config, logger logging.Logger) (*repository.Repositories, func() error) {
	dbPsql, err := database.NewPostgres(cfg.Postgres)
	if err != nil {
		logger.Logf(logging.FatalLevel, "could not connect to postgres: %s", err.Error())
	}
	userRepository := repository.NewUserPostgres(dbPsql)
	todoRepository := repository.NewTodoPostgres(dbPsql)

	closeDB := func() error {
		return dbPsql.Close()
	}

	return &repository.Repositories{
		UserRepository: userRepository,
		TodoRepository: todoRepository,
	}, closeDB
}

func GetServices(cfg *config.Config, logger logging.Logger, repositories *repository.Repositories) *service.Services {
	todoCache := cache.NewLFU(cfg.LFU)

	hash := hashing.NewBcrypt(cfg.Bcrypt)

	userEncoder, err := encoding.NewHashids(cfg.Hashids, encoding.UserKey)
	if err != nil {
		logger.Logf(logging.FatalLevel, "could not initialize user encoder: %s", err.Error())
	}
	todoEncoder, err := encoding.NewHashids(cfg.Hashids, encoding.TodoKey)
	if err != nil {
		logger.Logf(logging.FatalLevel, "could not initialize todo encoder: %s", err.Error())
	}

	authenticator := auth.NewJWT(cfg.JWT)

	userService := service.NewUserEncoded(hash, userEncoder, authenticator, repositories.UserRepository)
	todoService := service.NewTodoEncoded(todoCache, userEncoder, todoEncoder, repositories.TodoRepository)

	return &service.Services{
		UserService: userService,
		TodoService: todoService,
	}
}

func GetGinHandlers(cfg *config.Config, logger logging.Logger, services *service.Services) *handler.HandlersGin {
	authGin := handler.NewAuthGin(cfg.Gin, logger, services.UserService)
	middlewareGin := handler.NewMiddlewareGin(cfg.Gin, logger, services.UserService)
	todoGin := handler.NewTodoGin(cfg.Gin, logger, services.TodoService)

	return &handler.HandlersGin{
		Auth:       authGin,
		Middleware: middlewareGin,
		Todo:       todoGin,
	}
}
