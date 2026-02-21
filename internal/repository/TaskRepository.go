package repository

import (
	"taskmanager/internal/model"
)

type TaskRepository interface {
	GetAllTasks() ([]model.Task, error)
	GetTaskById(id string) (*model.Task, error)
	CreateTask(task model.Task) error
	UpdateTaskByID(id string, task model.Task) error
	DeleteTask(id string) error
}
