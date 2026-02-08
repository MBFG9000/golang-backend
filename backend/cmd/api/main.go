package main

import (
	"log"
	"net/http"
	"taskmanager/backend/internal/handler"
	"taskmanager/backend/internal/middleware"
	"taskmanager/backend/internal/repository"
	"taskmanager/backend/internal/service"
)

func main() {
	repo := repository.NewMemoryTaskRepository()
	service := service.NewTaskService(repo)
	handler := handler.NewTaskHandler(service)

	mux := http.NewServeMux()

	mux.HandleFunc("/tasks", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handler.GetTasks(w, r)
		case http.MethodPost:
			handler.CreateTask(w, r)
		case http.MethodPatch:
			handler.UpdateTaskStatus(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	wrappedHandler := middleware.LoggingMiddleware(
		middleware.APIKeyMiddleware(mux),
	)

	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", wrappedHandler))
}
