package controller_test

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/raffaele-pilloni/axxon-test/internal/app/http/controller"
	"github.com/raffaele-pilloni/axxon-test/internal/app/http/model/request"
	"github.com/raffaele-pilloni/axxon-test/internal/app/http/model/response"
	"github.com/raffaele-pilloni/axxon-test/internal/entity"
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

	It("should get task properly", func() {
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

	It("should create task properly", func() {
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
})
