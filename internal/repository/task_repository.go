package repository

import (
	"context"
	"github.com/raffaele-pilloni/axxon-test/internal/dal"
	"github.com/raffaele-pilloni/axxon-test/internal/entity"
)

type TaskRepository struct {
	dal *dal.DAL
}

func NewTaskRepository(
	dal *dal.DAL,
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
