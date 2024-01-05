package service

import (
	"context"
	"github.com/raffaele-pilloni/axxon-test/internal/db"
	"github.com/raffaele-pilloni/axxon-test/internal/entity"
	"github.com/raffaele-pilloni/axxon-test/internal/service/dto"
	"maps"
)

type TaskServiceInterface interface {
	CreateTask(ctx context.Context, createTaskDTO *dto.CreateTaskDTO) (*entity.Task, error)
	StartTaskProcessing(ctx context.Context, task *entity.Task) (*entity.Task, error)
	DoneTaskProcessing(ctx context.Context, task *entity.Task, doneTaskProcessingDTO *dto.DoneTaskProcessingDTO) (*entity.Task, error)
	ErrorTaskProcessing(ctx context.Context, task *entity.Task) (*entity.Task, error)
}

type TaskService struct {
	dal db.DALInterface
}

func NewTaskService(
	dal db.DALInterface,
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

func (t TaskService) StartTaskProcessing(ctx context.Context, task *entity.Task) (*entity.Task, error) {
	if err := t.dal.Save(ctx, task.StartProcessing()); err != nil {
		return nil, err
	}

	return task, nil
}

func (t TaskService) DoneTaskProcessing(ctx context.Context, task *entity.Task, doneTaskProcessingDTO *dto.DoneTaskProcessingDTO) (*entity.Task, error) {
	if err := t.dal.Save(ctx, task.DoneProcessing(
		maps.Clone(doneTaskProcessingDTO.ResponseHeaders),
		doneTaskProcessingDTO.ResponseStatusCode,
		doneTaskProcessingDTO.ResponseContentLength,
	)); err != nil {
		return nil, err
	}

	return task, nil
}

func (t TaskService) ErrorTaskProcessing(ctx context.Context, task *entity.Task) (*entity.Task, error) {
	if err := t.dal.Save(ctx, task.ErrorProcessing()); err != nil {
		return nil, err
	}

	return task, nil
}
