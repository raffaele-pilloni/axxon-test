package service

import (
	"context"
	"github.com/raffaele-pilloni/axxon-test/internal/client"
	clientdto "github.com/raffaele-pilloni/axxon-test/internal/client/dto"
	"github.com/raffaele-pilloni/axxon-test/internal/db"
	"github.com/raffaele-pilloni/axxon-test/internal/entity"
	"github.com/raffaele-pilloni/axxon-test/internal/service/dto"
	"log"
)

type TaskServiceInterface interface {
	CreateTask(ctx context.Context, createTaskDTO *dto.CreateTaskDTO) (*entity.Task, error)
	ProcessTask(ctx context.Context, task *entity.Task) (*entity.Task, error)
	StartTaskProcessing(ctx context.Context, task *entity.Task) (*entity.Task, error)
}

type TaskService struct {
	dal        db.DALInterface
	httpClient client.HTTPClientInterface
}

func NewTaskService(
	dal db.DALInterface,
	httpClient client.HTTPClientInterface,
) *TaskService {
	return &TaskService{
		dal:        dal,
		httpClient: httpClient,
	}
}

func (t TaskService) CreateTask(ctx context.Context, createTaskDTO *dto.CreateTaskDTO) (*entity.Task, error) {
	task, err := entity.NewTask(
		createTaskDTO.Method,
		createTaskDTO.URL,
		createTaskDTO.Headers,
		createTaskDTO.Body,
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

func (t TaskService) ProcessTask(ctx context.Context, task *entity.Task) (*entity.Task, error) {
	responseDTO, err := t.httpClient.Do(ctx, &clientdto.RequestDTO{
		Method:  task.MethodToString(),
		URL:     task.URL,
		Headers: task.RequestHeadersToMap(),
		Body:    task.RequestBodyToJSON(),
	})
	if err != nil {
		log.Printf("task request fail. Error: %v", err)
		return t.errorTaskProcessing(ctx, task)
	}

	if _, err := t.doneTaskProcessing(ctx, task, responseDTO); err != nil {
		log.Printf("task request fail. Error: %v", err)
		return t.errorTaskProcessing(ctx, task)
	}

	return task, nil
}

func (t TaskService) doneTaskProcessing(ctx context.Context, task *entity.Task, responseDto *clientdto.ResponseDTO) (*entity.Task, error) {
	if err := t.dal.Save(ctx, task.DoneProcessing(
		responseDto.Headers,
		responseDto.StatusCode,
		len(responseDto.Body),
	)); err != nil {
		return nil, err
	}

	return task, nil
}

func (t TaskService) errorTaskProcessing(ctx context.Context, task *entity.Task) (*entity.Task, error) {
	if err := t.dal.Save(ctx, task.ErrorProcessing()); err != nil {
		return nil, err
	}

	return task, nil
}
