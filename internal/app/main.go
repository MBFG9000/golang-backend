package app

import (
	"context"
	"github.com/MBFG9000/golang-backend/internal/config"
	"github.com/MBFG9000/golang-backend/internal/handler"
	"github.com/MBFG9000/golang-backend/internal/middleware"
	"github.com/MBFG9000/golang-backend/internal/repository"
	postgres "github.com/MBFG9000/golang-backend/internal/repository/postgresql"
	userservice "github.com/MBFG9000/golang-backend/internal/service/users"
	"log"
	"net/http"
)

func Run() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	//config.LoadEnviroment()

	dbConfig := config.GetConfig()

	postgres := postgres.NewPostgres(ctx, dbConfig)

	repositories := repository.NewRepositories(postgres)

	userService := userservice.NewUserService(repositories)

	userHandler := handler.NewUserHandler(userService)

	mux := http.NewServeMux()

	handler.SetupRoutes(mux, userHandler)

	wrappedHandler := middleware.LoggingMiddleware(
		middleware.APIKeyMiddleware(mux),
	)

	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", wrappedHandler))
}
