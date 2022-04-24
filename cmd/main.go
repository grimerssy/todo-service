package main

import (
	"context"
	"database/sql"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/grimerssy/todo-service/internal/server"
	"github.com/grimerssy/todo-service/pkg/handler"
	"github.com/grimerssy/todo-service/pkg/repository"
	"github.com/grimerssy/todo-service/pkg/service"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
)

const (
	configPath = "configs"
	devConfig  = "dev"
)

type config struct {
	Http     server.ConfigHttp
	Postgres repository.ConfigPsql
	JWT      service.ConfigJWT
	Hashids  service.ConfigHashids
	Bcrypt   service.ConfigBcrypt
}

func main() {
	cfg := getConfig(devConfig)

	db, repositories := getDbAndRepositories(cfg)
	services := getServices(cfg, repositories)
	handlers := getGinHandlers(services)

	srv := new(server.HttpServer)
	if err := srv.Run(cfg.Http, handlers.InitRoutes()); err != nil {
		log.Fatalf("an error occured while running http server: %s", err.Error())
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	if err := srv.Shutdown(context.Background()); err != nil {
		log.Fatalf("an error occured while shutting down the server: %s", err.Error())
	}

	if err := db.Close(); err != nil {
		log.Fatalf("an error occured while closing db connection: %s", err.Error())
	}
}

func getConfig(configName string) *config {
	const (
		postgresPasswordKey = "PSQL_TODO"
		jwtSecretKey        = "JWT_SECRET"
	)

	var cfg *config
	viper.AddConfigPath(configPath)
	viper.SetConfigName(configName)
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

func getDbAndRepositories(cfg *config) (*sql.DB, *repository.Repositories) {
	dbPsql, err := repository.NewDbPsql(cfg.Postgres)
	if err != nil {
		log.Fatalf("an error occured while connecting to postgres: %s", err.Error())
	}
	userRepository := repository.NewUserPsql(dbPsql)
	todoRepository := repository.NewTodoPsql(dbPsql)

	return dbPsql, &repository.Repositories{
		UserRepository: userRepository,
		TodoRepository: todoRepository,
	}
}

func getServices(cfg *config, repositories *repository.Repositories) *service.Services {
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

func getGinHandlers(services *service.Services) *handler.HandlersGin {
	authGin := handler.NewAuthGin(services.AuthService, services.UserService)
	middlewareGin := handler.NewMiddlewareGin(services.AuthService)
	todoGin := handler.NewTodoGin(services.TodoService)

	return &handler.HandlersGin{
		Auth:       authGin,
		Middleware: middlewareGin,
		Todo:       todoGin,
	}
}
