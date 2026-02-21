package app

import (
	"context"
	"log"
	"net/http"
	"taskmanager/internal/config"
	"taskmanager/internal/handler"
	"taskmanager/internal/middleware"
	"taskmanager/internal/repository"
	postgres "taskmanager/internal/repository/postgresql"
	"taskmanager/internal/service/userservice"
)

func Run() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	config.LoadEnviroment()

	dbConfig := config.GetConfig()

	postgres := postgres.NewPostgres(ctx, dbConfig)

	repositories := repository.NewRepositories(postgres)

	userUseCase := userservice.NewUserUseCase(repositories)

	userHandler := handler.NewUserHandler(userUseCase)

	mux := http.NewServeMux()

	handler.SetupRoutes(mux, userHandler)

	wrappedHandler := middleware.LoggingMiddleware(
		middleware.APIKeyMiddleware(mux),
	)

	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", wrappedHandler))
}
