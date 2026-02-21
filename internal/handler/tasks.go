package handler

import (
	"encoding/json"
	"net/http"
	service "taskmanager/internal/service/taskservice"
	"taskmanager/internal/utils"
)

type TaskHandler struct {
	service *service.TaskService
}

type CreateTaskRequest struct {
	Title string `json:"title"`
}
type UpdateTaskRequest struct {
	Done bool `json:"done"`
}
type UpdateResponse struct {
	Updated bool `json:"updated"`
}
type ErrorResponse struct {
	Error string `json:"error"`
}

func NewTaskHandler(service *service.TaskService) *TaskHandler {
	return &TaskHandler{service: service}
}

func (h *TaskHandler) GetTasks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id := r.URL.Query().Get("id")

	if id != "" {
		task, err := h.service.GetTaskByID(id)

		if err != nil {
			switch err {
			case utils.ErrObjectNotFound:
				w.WriteHeader(http.StatusNotFound)
				json.NewEncoder(w).Encode(ErrorResponse{Error: err.Error()})
			case utils.ErrIDEmpty, utils.ErrIDNotNumber, utils.ErrIDNotPositive:
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(ErrorResponse{Error: err.Error()})
			}
			return
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(task)
		return
	}

	tasks, err := h.service.GetAllTasks()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "internal server error"})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(tasks)
}

func (h *TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var input CreateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "invalid request body"})
		return
	}

	defer r.Body.Close()

	task, err := h.service.CreateTask(input.Title)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{err.Error()})
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(task)
}

func (h *TaskHandler) UpdateTaskStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id := r.URL.Query().Get("id")

	var input UpdateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "invalid request body"})
		return
	}
	defer r.Body.Close()

	err := h.service.UpdateTaskStatus(id, input.Done)

	if err != nil {
		switch err {
		case utils.ErrObjectNotFound:
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(ErrorResponse{Error: err.Error()})
		case utils.ErrIDEmpty, utils.ErrIDNotNumber, utils.ErrIDNotPositive:
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrorResponse{Error: err.Error()})
		default:
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "invalid id"})

		}
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(UpdateResponse{Updated: true})
}
