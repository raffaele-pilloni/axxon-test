package service

import (
	"context"
	"github.com/raffaele-pilloni/axxon-test/internal/dal"
	"github.com/raffaele-pilloni/axxon-test/internal/entity"
	"github.com/raffaele-pilloni/axxon-test/internal/service/dto"
	"maps"
)

type TaskService struct {
	dal *dal.DAL
}

func NewTaskService(
	dal *dal.DAL,
) *TaskService {
	return &TaskService{
		dal: dal,
	}
}

func (t TaskService) CreateTask(ctx context.Context, createTaskDTO *dto.CreateTaskDTO) (*entity.Task, error) {
	task, err := entity.NewTask(
		createTaskDTO.Method,
		createTaskDTO.URL,
		maps.Clone(createTaskDTO.Headers),
		maps.Clone(createTaskDTO.Body),
	)
	if err != nil {
		return nil, err
	}

	if err := t.dal.Save(ctx, task); err != nil {
		return nil, err
	}

	return task, nil
}
