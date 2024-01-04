package repository

import (
	"errors"
	"github.com/raffaele-pilloni/axxon-test/internal/entity"
	"gorm.io/gorm"
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

func (t TaskRepository) GetTaskById(taskID int) (*entity.Task, error) {
	var task entity.Task

	query := t.gormDB.Find(&task, taskID)
	if errors.Is(query.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	if query.Error != nil {
		return nil, query.Error
	}

	return &task, nil
}
