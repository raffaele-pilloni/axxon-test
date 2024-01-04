package service

import (
	"github.com/raffaele-pilloni/axxon-test/internal/entity"
	"github.com/raffaele-pilloni/axxon-test/internal/service/dto"
	"gorm.io/gorm"
	"maps"
)

type TaskService struct {
	gormDB *gorm.DB
}

func NewTaskService(
	gormDB *gorm.DB,
) *TaskService {
	return &TaskService{
		gormDB: gormDB,
	}
}

func (t TaskService) CreateTask(createTaskDTO *dto.CreateTaskDTO) (*entity.Task, error) {
	task, err := entity.NewTask(
		createTaskDTO.Method,
		createTaskDTO.URL,
		maps.Clone(createTaskDTO.Headers),
		maps.Clone(createTaskDTO.Body),
	)
	if err != nil {
		return nil, err
	}

	if insert := t.gormDB.Save(task); insert.Error != nil {
		return nil, insert.Error
	}

	return task, nil
}
