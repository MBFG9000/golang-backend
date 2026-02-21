package slicerepo

import (
	"fmt"
	"taskmanager/internal/model"
	"taskmanager/internal/utils"
)

type MemoryTaskRepository struct {
	tasks  []model.Task
	lastID int
}

func NewMemoryTaskRepository() *MemoryTaskRepository {
	return &MemoryTaskRepository{
		tasks: []model.Task{
			{ID: "1", Title: "Изучить Go", Done: true},
			{ID: "2", Title: "Написать программу", Done: false},
			{ID: "3", Title: "Протестировать код", Done: false},
			{ID: "4", Title: "Оптимизировать производительность", Done: false},
		},
		lastID: 4,
	}
}

//***********************
// Repository functions
//***********************

func (r *MemoryTaskRepository) GetAllTasks() ([]model.Task, error) {
	return r.tasks, nil
}

func (r *MemoryTaskRepository) GetTaskById(id string) (*model.Task, error) {
	for i := range r.tasks {
		if r.tasks[i].ID == id {
			return &r.tasks[i], nil
		}

	}

	return nil, utils.ErrObjectNotFound
}

func (r *MemoryTaskRepository) CreateTask(task model.Task) error {
	r.lastID++
	task.ID = fmt.Sprintf("%d", r.lastID)

	r.tasks = append(r.tasks, task)
	return nil
}

func (r *MemoryTaskRepository) UpdateTaskByID(id string, task model.Task) error {
	index, err := getIndexById(r, id)

	if err != nil {
		return err
	}

	r.tasks[index] = task
	return nil
}

func (r *MemoryTaskRepository) DeleteTask(id string) error {
	index, err := getIndexById(r, id)

	if err != nil {
		return err
	}

	newtasks, err := utils.RemoveOrdered(r.tasks, index)

	if err != nil {
		return err
	}

	r.tasks = newtasks
	return nil
}

func (r *MemoryTaskRepository) MarkTaskDone(id string) error {
	// Faster but breaks the business logic

	index, err := getIndexById(r, id)

	if err != nil {
		return err
	}

	r.tasks[index].Done = true
	return nil
}

//***********************
// Additional functions
//***********************

func getIndexById(r *MemoryTaskRepository, id string) (int, error) {
	for i, t := range r.tasks {
		if t.ID == id {
			return i, nil
		}
	}

	return -1, fmt.Errorf("Task with ID %s does not exist", id)
}
