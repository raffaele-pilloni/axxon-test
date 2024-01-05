package repository

import (
	"context"
	"github.com/raffaele-pilloni/axxon-test/internal/db"
	"github.com/raffaele-pilloni/axxon-test/internal/entity"
)

type TaskRepositoryInterface interface {
	FindTaskByID(ctx context.Context, taskID int) (*entity.Task, error)
	FindTasksToProcess(ctx context.Context, limit int) ([]*entity.Task, error)
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

func (t TaskRepository) FindTasksToProcess(ctx context.Context, limit int) ([]*entity.Task, error) {
	var tasks []*entity.Task

	if err := t.dal.FindBy(
		ctx,
		&tasks,
		db.Criteria{"Status": entity.StatusNew},
		"CreatedAt desc",
		limit); err != nil {
		return nil, err
	}

	return tasks, nil
}
