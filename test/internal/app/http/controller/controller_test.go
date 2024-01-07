package controller_test

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/raffaele-pilloni/axxon-test/internal/app/http/controller"
	"github.com/raffaele-pilloni/axxon-test/internal/app/http/model/response"
	"github.com/raffaele-pilloni/axxon-test/internal/entity"
	dmock "github.com/raffaele-pilloni/axxon-test/mock"
	"gorm.io/datatypes"
	"net/http"
	"net/http/httptest"
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

		expectedBody, err := json.Marshal(&response.GetTaskModelResponse{
			ID:             task.ID,
			Status:         task.StatusToString(),
			Headers:        task.ResponseHeadersToMap(),
			HTTPStatusCode: task.ResponseStatusCodeToInt(),
			Length:         task.ResponseContentLengthToInt(),
		})

		Ω(err).To(BeNil())

		Ω(w.Code).To(Equal(http.StatusOK))
		Ω(w.Body.Bytes()).To(Equal(expectedBody))
	})
})
