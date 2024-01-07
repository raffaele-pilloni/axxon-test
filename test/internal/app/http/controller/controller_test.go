package controller_test

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/raffaele-pilloni/axxon-test/internal/app/http/controller"
	"github.com/raffaele-pilloni/axxon-test/internal/app/http/handler"
	"github.com/raffaele-pilloni/axxon-test/internal/app/http/model/request"
	"github.com/raffaele-pilloni/axxon-test/internal/app/http/model/response"
	"github.com/raffaele-pilloni/axxon-test/internal/entity"
	applicationerror "github.com/raffaele-pilloni/axxon-test/internal/error"
	"github.com/raffaele-pilloni/axxon-test/internal/service/dto"
	dmock "github.com/raffaele-pilloni/axxon-test/mock"
	"github.com/stretchr/testify/mock"
	"gorm.io/datatypes"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strconv"
)

var _ = Describe("Task Controller Tests", func() {
	var (
		mockTaskRepository *dmock.TaskRepositoryInterface
		mockTaskService    *dmock.TaskServiceInterface

		taskController *controller.TaskController
	)

	BeforeEach(func() {
		mockTaskRepository = dmock.NewTaskRepositoryInterface(GinkgoT())
		mockTaskService = dmock.NewTaskServiceInterface(GinkgoT())

		taskController = controller.NewTaskController(
			mockTaskRepository,
			mockTaskService,
		)
	})

	It("should get task properly (200)", func() {
		taskID := 1

		request := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/task/%d", taskID), nil)
		w := httptest.NewRecorder()

		request = mux.SetURLVars(request, map[string]string{
			"taskId": strconv.Itoa(taskID),
		})

		task := entity.Task{
			ID:     0,
			Status: entity.StatusDone,
			ResponseHeaders: datatypes.NewJSONType(map[string][]string{
				"test": {"test"},
			}),
			ResponseStatusCode:    sql.NullInt64{Int64: 200},
			ResponseContentLength: sql.NullInt64{Int64: 30},
		}

		mockTaskRepository.On(
			"FindTaskByID",
			request.Context(),
			taskID,
		).Once().Return(&task, nil)

		taskController.GetTask(w, request)

		expectedResponseBody, err := json.Marshal(&response.GetTaskModelResponse{
			ID:             task.ID,
			Status:         task.StatusToString(),
			Headers:        task.ResponseHeadersToMap(),
			HTTPStatusCode: task.ResponseStatusCodeToInt(),
			Length:         task.ResponseContentLengthToInt(),
		})

		Ω(err).To(BeNil())

		Ω(w.Code).To(Equal(http.StatusOK))
		Ω(w.Body.Bytes()).To(Equal(expectedResponseBody))
	})

	It("should create task properly (200)", func() {
		createTaskModelRequest := &request.CreateTaskModelRequest{
			Method: "POST",
			URL:    "http://test.com",
			Headers: map[string][]string{
				"test": {"test"},
			},
			Body: map[string]interface{}{
				"test": "test",
			},
		}

		requestBody, err := json.Marshal(createTaskModelRequest)

		Ω(err).To(BeNil())

		request := httptest.NewRequest(http.MethodPost, "/task", bytes.NewReader(requestBody))
		w := httptest.NewRecorder()

		task := entity.Task{
			ID: 1,
		}

		mockTaskService.On(
			"CreateTask",
			request.Context(),
			mock.MatchedBy(func(actualCreateTaskDTO *dto.CreateTaskDTO) bool {
				return actualCreateTaskDTO.Method == createTaskModelRequest.Method &&
					actualCreateTaskDTO.URL == createTaskModelRequest.URL &&
					reflect.DeepEqual(actualCreateTaskDTO.Headers, createTaskModelRequest.Headers) &&
					reflect.DeepEqual(actualCreateTaskDTO.Body, createTaskModelRequest.Body)
			}),
		).Once().Return(&task, nil)

		taskController.CreateTask(w, request)

		expectedResponseBody, err := json.Marshal(&response.CreateTaskModelResponse{
			ID: task.ID,
		})

		Ω(err).To(BeNil())

		Ω(w.Code).To(Equal(http.StatusOK))
		Ω(w.Body.Bytes()).To(Equal(expectedResponseBody))
	})

	It("should fail get task when task id is not valid (500)", func() {
		taskID := "invalid"

		request := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/task/%s", taskID), nil)
		w := httptest.NewRecorder()

		request = mux.SetURLVars(request, map[string]string{
			"taskId": taskID,
		})

		taskController.GetTask(w, request)

		expectedResponseBody, err := json.Marshal(&response.ErrorModelResponse{
			Code:    handler.HTTPInternalServerErrorCode,
			Message: handler.InternalServerErrorMessage,
		})

		Ω(err).To(BeNil())

		Ω(w.Code).To(Equal(http.StatusInternalServerError))
		Ω(w.Body.Bytes()).To(Equal(expectedResponseBody))
	})

	It("should fail get task when entity not found (404)", func() {
		taskID := 1

		request := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/task/%d", taskID), nil)
		w := httptest.NewRecorder()

		request = mux.SetURLVars(request, map[string]string{
			"taskId": strconv.Itoa(taskID),
		})

		mockTaskRepository.On(
			"FindTaskByID",
			request.Context(),
			taskID,
		).Once().Return(nil, &applicationerror.EntityNotFoundError{})

		taskController.GetTask(w, request)

		expectedResponseBody, err := json.Marshal(&response.ErrorModelResponse{
			Code:    handler.HTTPEntityNotFoundErrorCode,
			Message: "There is no Task with id 1.",
		})

		Ω(err).To(BeNil())

		Ω(w.Code).To(Equal(http.StatusNotFound))
		Ω(w.Body.Bytes()).To(Equal(expectedResponseBody))
	})

	It("should fail get task when find task by id fail (500)", func() {
		taskID := 1

		request := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/task/%d", taskID), nil)
		w := httptest.NewRecorder()

		request = mux.SetURLVars(request, map[string]string{
			"taskId": strconv.Itoa(taskID),
		})

		mockTaskRepository.On(
			"FindTaskByID",
			request.Context(),
			taskID,
		).Once().Return(nil, errors.New("error test"))

		taskController.GetTask(w, request)

		expectedResponseBody, err := json.Marshal(&response.ErrorModelResponse{
			Code:    handler.HTTPInternalServerErrorCode,
			Message: handler.InternalServerErrorMessage,
		})

		Ω(err).To(BeNil())

		Ω(w.Code).To(Equal(http.StatusInternalServerError))
		Ω(w.Body.Bytes()).To(Equal(expectedResponseBody))
	})

	It("should fail create task when body is not json (400)", func() {
		requestBody := []byte("invalid json")

		request := httptest.NewRequest(http.MethodPost, "/task", bytes.NewReader(requestBody))
		w := httptest.NewRecorder()

		taskController.CreateTask(w, request)

		expectedResponseBody, err := json.Marshal(&response.ErrorModelResponse{
			Code:    handler.HTTPBadRequestErrorCode,
			Message: handler.InvalidJSONBodyErrorMessage,
		})

		Ω(err).To(BeNil())

		Ω(w.Code).To(Equal(http.StatusBadRequest))
		Ω(w.Body.Bytes()).To(Equal(expectedResponseBody))
	})

	It("should fail create task when body is invalid (200)", func() {
		createTaskModelRequest := &request.CreateTaskModelRequest{
			Method: "",
			URL:    "http://test.com",
			Headers: map[string][]string{
				"test": {"test"},
			},
			Body: map[string]interface{}{
				"test": "test",
			},
		}

		requestBody, err := json.Marshal(createTaskModelRequest)

		Ω(err).To(BeNil())

		request := httptest.NewRequest(http.MethodPost, "/task", bytes.NewReader(requestBody))
		w := httptest.NewRecorder()

		taskController.CreateTask(w, request)

		expectedResponseBody, err := json.Marshal(&response.ErrorModelResponse{
			Code:    handler.HTTPBadRequestErrorCode,
			Message: "Key: 'CreateTaskModelRequest.Method' Error:Field validation for 'Method' failed on the 'required' tag",
		})

		Ω(err).To(BeNil())

		Ω(w.Code).To(Equal(http.StatusBadRequest))
		Ω(w.Body.Bytes()).To(Equal(expectedResponseBody))
	})

	It("should fail create task and return correctly error when service create task fail with application error (422)", func() {
		createTaskModelRequest := &request.CreateTaskModelRequest{
			Method: "POST",
			URL:    "http://test.com",
			Headers: map[string][]string{
				"test": {"test"},
			},
			Body: map[string]interface{}{
				"test": "test",
			},
		}

		requestBody, err := json.Marshal(createTaskModelRequest)

		Ω(err).To(BeNil())

		request := httptest.NewRequest(http.MethodPost, "/task", bytes.NewReader(requestBody))
		w := httptest.NewRecorder()

		mockApplicationError := dmock.NewApplicationErrorInterface(GinkgoT())

		errorCode := 1000
		errorMessage := "errorTest"

		mockApplicationError.On("Code").Once().Return(errorCode)
		mockApplicationError.On("Message").Once().Return(errorMessage)
		mockApplicationError.On("Error").Once().Return(errorMessage)

		mockTaskService.On(
			"CreateTask",
			request.Context(),
			mock.MatchedBy(func(actualCreateTaskDTO *dto.CreateTaskDTO) bool {
				return actualCreateTaskDTO.Method == createTaskModelRequest.Method &&
					actualCreateTaskDTO.URL == createTaskModelRequest.URL &&
					reflect.DeepEqual(actualCreateTaskDTO.Headers, createTaskModelRequest.Headers) &&
					reflect.DeepEqual(actualCreateTaskDTO.Body, createTaskModelRequest.Body)
			}),
		).Once().Return(nil, mockApplicationError)

		taskController.CreateTask(w, request)

		expectedResponseBody, err := json.Marshal(&response.ErrorModelResponse{
			Code:    errorCode,
			Message: errorMessage,
		})

		Ω(err).To(BeNil())

		Ω(w.Code).To(Equal(http.StatusUnprocessableEntity))
		Ω(w.Body.Bytes()).To(Equal(expectedResponseBody))
	})

	It("should fail create task when service create task fail (422)", func() {
		createTaskModelRequest := &request.CreateTaskModelRequest{
			Method: "POST",
			URL:    "http://test.com",
			Headers: map[string][]string{
				"test": {"test"},
			},
			Body: map[string]interface{}{
				"test": "test",
			},
		}

		requestBody, err := json.Marshal(createTaskModelRequest)

		Ω(err).To(BeNil())

		request := httptest.NewRequest(http.MethodPost, "/task", bytes.NewReader(requestBody))
		w := httptest.NewRecorder()

		mockTaskService.On(
			"CreateTask",
			request.Context(),
			mock.MatchedBy(func(actualCreateTaskDTO *dto.CreateTaskDTO) bool {
				return actualCreateTaskDTO.Method == createTaskModelRequest.Method &&
					actualCreateTaskDTO.URL == createTaskModelRequest.URL &&
					reflect.DeepEqual(actualCreateTaskDTO.Headers, createTaskModelRequest.Headers) &&
					reflect.DeepEqual(actualCreateTaskDTO.Body, createTaskModelRequest.Body)
			}),
		).Once().Return(nil, errors.New("error test"))

		taskController.CreateTask(w, request)

		expectedResponseBody, err := json.Marshal(&response.ErrorModelResponse{
			Code:    handler.HTTPInternalServerErrorCode,
			Message: handler.InternalServerErrorMessage,
		})

		Ω(err).To(BeNil())

		Ω(w.Code).To(Equal(http.StatusInternalServerError))
		Ω(w.Body.Bytes()).To(Equal(expectedResponseBody))
	})
})
