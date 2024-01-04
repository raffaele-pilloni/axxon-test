package repository

import (
	"context"
	"errors"
	"github.com/raffaele-pilloni/axxon-test/internal/entity"
	"gorm.io/gorm"
	"time"
)

type TaskRepository struct {
	gormDB *gorm.DB
}

func NewTaskRepository(
	gormDB *gorm.DB,
) *TaskRepository {
	return &TaskRepository{
		gormDB: gormDB,
	}
}

func (t TaskRepository) GetTaskById(ctx context.Context, taskID int) (*entity.Task, error) {
	ctx, cancelCtx := context.WithTimeout(ctx, 10*time.Second)
	defer cancelCtx()

	var task entity.Task

	query := t.gormDB.WithContext(ctx).Find(&task, taskID)
	if errors.Is(query.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	if query.Error != nil {
		return nil, query.Error
	}

	return &task, nil
}
