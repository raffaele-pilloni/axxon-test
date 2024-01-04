package controller

import (
	"github.com/gorilla/mux"
	"github.com/raffaele-pilloni/axxon-test/internal/app/http/handler"
	"github.com/raffaele-pilloni/axxon-test/internal/app/http/model/request"
	"github.com/raffaele-pilloni/axxon-test/internal/app/http/model/response"
	"github.com/raffaele-pilloni/axxon-test/internal/repository"
	"github.com/raffaele-pilloni/axxon-test/internal/service"
	"github.com/raffaele-pilloni/axxon-test/internal/service/dto"
	"maps"
	"net/http"
	"strconv"
)

type TaskController struct {
	taskRepository *repository.TaskRepository
	taskService    *service.TaskService
}

func NewTaskController(
	taskRepository *repository.TaskRepository,
	taskService *service.TaskService,
) *TaskController {
	return &TaskController{
		taskRepository: taskRepository,
		taskService:    taskService,
	}
}

func (t *TaskController) GetTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	taskID, err := strconv.Atoi(vars["taskId"])
	if err != nil {
		handler.HandleError(w, err)
		return
	}

	task, err := t.taskRepository.FindTaskById(r.Context(), taskID)
	if err != nil {
		handler.HandleError(w, err)
		return
	}

	handler.HandleSuccess(w, &response.GetTaskModelResponse{
		ID:             task.ID,
		Status:         string(task.Status),
		HttpStatusCode: task.ResponseStatusCode,
		Headers:        maps.Clone(task.ResponseHeaders),
		Length:         task.ResponseContentLength,
	})
}

func (t *TaskController) CreateTask(w http.ResponseWriter, r *http.Request) {
	createTaskModelRequest, err := handler.HandleRequest[request.CreateTaskModelRequest](r)
	if err != nil {
		handler.HandleError(w, err)
		return
	}

	task, err := t.taskService.CreateTask(
		r.Context(),
		&dto.CreateTaskDTO{
			Method:  createTaskModelRequest.Method,
			URL:     createTaskModelRequest.URL,
			Headers: maps.Clone(createTaskModelRequest.Headers),
			Body:    maps.Clone(createTaskModelRequest.Body),
		})
	if err != nil {
		handler.HandleError(w, err)
		return
	}

	handler.HandleSuccess(w, &response.CreateTaskModelResponse{
		ID: task.ID,
	})
}
