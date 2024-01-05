package repository

import (
	"context"
	"github.com/raffaele-pilloni/axxon-test/internal/db"
	"github.com/raffaele-pilloni/axxon-test/internal/entity"
)

type TaskRepositoryInterface interface {
	FindTaskByID(ctx context.Context, taskID int) (*entity.Task, error)
}

type TaskRepository struct {
	dal db.DALInterface
}

func NewTaskRepository(
	dal db.DALInterface,
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
