package repository

import (
	"context"
	"github.com/raffaele-pilloni/axxon-test/internal/dal"
	"github.com/raffaele-pilloni/axxon-test/internal/entity"
)

type TaskRepositoryInterface interface {
	FindTaskByID(ctx context.Context, taskID int) (*entity.Task, error)
}

type TaskRepository struct {
	dal dal.DALInterface
}

func NewTaskRepository(
	dal dal.DALInterface,
) *TaskRepository {
	return &TaskRepository{
		dal: dal,
	}
}

func (t TaskRepository) FindTaskByID(ctx context.Context, taskID int) (*entity.Task, error) {
	var task entity.Task

	if err := t.dal.FindByID(ctx, &task, taskID); err != nil {
		return nil, err
	}

	return &task, nil
}
