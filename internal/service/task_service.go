package service

import (
	"bytes"
	"context"
	"github.com/raffaele-pilloni/axxon-test/internal/db"
	"github.com/raffaele-pilloni/axxon-test/internal/entity"
	"github.com/raffaele-pilloni/axxon-test/internal/service/dto"
	"maps"
	"net/http"
)

type TaskServiceInterface interface {
	CreateTask(ctx context.Context, createTaskDTO *dto.CreateTaskDTO) (*entity.Task, error)
	ProcessTask(ctx context.Context, task *entity.Task) (*entity.Task, error)
}

type TaskService struct {
	dal        db.DALInterface
	httpClient *http.Client
}

func NewTaskService(
	dal db.DALInterface,
	httpClient *http.Client,
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

func (t TaskService) ProcessTask(ctx context.Context, task *entity.Task) (*entity.Task, error) {
	if err := t.dal.Save(ctx, task.StartProcessing()); err != nil {
		return nil, err
	}

	request, err := http.NewRequestWithContext(ctx, task.MethodToString(), task.URL, bytes.NewReader(task.RequestBodyToJSON()))
	if err != nil {
		return t.errorTaskProcessing(ctx, task)
	}

	response, err := t.httpClient.Do(request)
	if err != nil {
		return t.errorTaskProcessing(ctx, task)
	}
	defer response.Body.Close()

	if _, err := t.doneTaskProcessing(ctx, task, response); err != nil {
		return t.errorTaskProcessing(ctx, task)
	}

	return task, nil
}

func (t TaskService) doneTaskProcessing(ctx context.Context, task *entity.Task, response *http.Response) (*entity.Task, error) {
	if err := t.dal.Save(ctx, task.DoneProcessing(
		maps.Clone(response.Header),
		response.StatusCode,
		int(response.ContentLength),
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
