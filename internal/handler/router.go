package handler

import (
	"net/http"
)

func SetupRoutes(mux *http.ServeMux, userHandler *UserHandler) {
	mux.HandleFunc("GET /users", userHandler.GetUsers)
	mux.HandleFunc("GET /users/{id}", userHandler.GetUserByID)
	mux.HandleFunc("POST /users", userHandler.CreateUser)
	mux.HandleFunc("PUT /users/{id}", userHandler.UpdateUserByID)
	mux.HandleFunc("DELETE /users/{id}", userHandler.DeleteUserByID)
}
