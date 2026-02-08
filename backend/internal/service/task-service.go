package service

import (
	"fmt"
	"strconv"
	"strings"
	"taskmanager/backend/internal/model"
	"taskmanager/backend/internal/repository"
	"taskmanager/backend/internal/utils"
)

//***********************
// Service functions
//***********************

type TaskService struct {
	repo repository.TaskRepository
}

func NewTaskService(repo repository.TaskRepository) *TaskService {
	return &TaskService{
		repo: repo,
	}
}

func (s *TaskService) GetAllTasks() ([]model.Task, error) {
	return s.repo.GetAllTasks()
}

func (s *TaskService) GetTaskByID(id string) (*model.Task, error) {
	if err := validateID(id); err != nil {
		return nil, err
	}

	return s.repo.GetTaskById(id)
}

func (s *TaskService) CreateTask(title string) (*model.Task, error) {
	if err := validateTitle(title); err != nil {
		return nil, err
	}
	task := model.Task{
		Title: title,
		Done:  false,
	}

	if err := s.repo.CreateTask(task); err != nil {
		return nil, err
	}

	tasks, _ := s.repo.GetAllTasks()
	return &tasks[len(tasks)-1], nil
}

func (s *TaskService) UpdateTaskStatus(id string, done bool) error {
	if err := validateID(id); err != nil {
		return err
	}

	task, err := s.repo.GetTaskById(id)
	if err != nil {
		return err
	}

	task.Done = done

	return s.repo.UpdateTaskByID(id, *task)
}

func (s *TaskService) DeleteTask(id string) error {
	if err := validateID(id); err != nil {
		return err
	}

	if err := s.repo.DeleteTask(id); err != nil {
		return err
	}

	return nil
}

//***********************
// Additional functions
//***********************

func validateID(id string) error {
	id = strings.TrimSpace(id)

	if id == "" {
		return utils.ErrIDEmpty
	}

	numID, err := strconv.Atoi(id)

	if err != nil {
		return utils.ErrIDNotNumber
	}

	if numID <= 0 {
		return utils.ErrIDNotPositive
	}
	return nil
}

func validateTitle(title string) error {
	title = strings.TrimSpace(title)

	if title == "" {
		return fmt.Errorf("invalid title")
	}

	return nil
}
