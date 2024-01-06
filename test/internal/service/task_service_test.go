package service_test

import (
	"context"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	clientdto "github.com/raffaele-pilloni/axxon-test/internal/client/dto"
	"github.com/raffaele-pilloni/axxon-test/internal/entity"
	"github.com/raffaele-pilloni/axxon-test/internal/service"
	servicedto "github.com/raffaele-pilloni/axxon-test/internal/service/dto"
	dmock "github.com/raffaele-pilloni/axxon-test/mock"
	"github.com/stretchr/testify/mock"
	"gorm.io/datatypes"
	"reflect"
	"time"
)

var _ = Describe("Task Service Tests", func() {
	var (
		mockDALInterface *dmock.DALInterface
		mockHTTPClient   *dmock.HTTPClientInterface

		taskService *service.TaskService
	)

	BeforeEach(func() {
		mockDALInterface = dmock.NewDALInterface(GinkgoT())
		mockHTTPClient = dmock.NewHTTPClientInterface(GinkgoT())

		taskService = service.NewTaskService(
			mockDALInterface,
			mockHTTPClient,
		)
	})

	It("should create task properly", func() {
		context := context.Background()

		createRequestDTO := servicedto.CreateTaskDTO{
			Method: "POST",
			URL:    "http://test.com",
			Headers: map[string][]string{
				"test": {"test"},
			},
			Body: map[string]interface{}{
				"test": "test",
			},
		}

		mockDALInterface.On(
			"Save",
			context,
			mock.MatchedBy(func(actualTask *entity.Task) bool {
				return actualTask.MethodToString() == createRequestDTO.Method &&
					actualTask.URL == createRequestDTO.URL &&
					reflect.DeepEqual(actualTask.RequestHeadersToMap(), createRequestDTO.Headers) &&
					reflect.DeepEqual(actualTask.RequestBodyToMap(), createRequestDTO.Body) &&
					actualTask.ID == 0 &&
					actualTask.Status == entity.StatusNew &&
					actualTask.ResponseStatusCodeToInt() == 0 &&
					actualTask.ResponseContentLengthToInt() == 0 &&
					actualTask.CreatedAt != time.Time{} &&
					actualTask.UpdatedAt != time.Time{}
			}),
		).Return(nil)

		task, err := taskService.CreateTask(context, &createRequestDTO)
		Ω(err).To(BeNil())

		Ω(task.MethodToString()).To(Equal(createRequestDTO.Method))
		Ω(task.URL).To(Equal(createRequestDTO.URL))
		Ω(task.RequestHeadersToMap()).To(Equal(createRequestDTO.Headers))
		Ω(task.RequestBodyToMap()).To(Equal(createRequestDTO.Body))

		Ω(task.ID).To(BeZero())
		Ω(task.StatusToString()).To(Equal(string(entity.StatusNew)))
		Ω(task.ResponseHeadersToMap()).To(BeEmpty())
		Ω(task.ResponseStatusCodeToInt()).To(BeZero())
		Ω(task.ResponseContentLengthToInt()).To(BeZero())
		Ω(task.CreatedAt).ToNot(Equal(time.Time{}))
		Ω(task.UpdatedAt).ToNot(Equal(time.Time{}))
	})

	It("should start task processing properly", func() {
		context := context.Background()

		task := &entity.Task{}

		mockDALInterface.On(
			"Save",
			context,
			mock.MatchedBy(func(actualTask *entity.Task) bool {
				return actualTask.Status == entity.StatusInProcess
			}),
		).Return(nil)

		task, err := taskService.StartTaskProcessing(context, task)
		Ω(err).To(BeNil())

		Ω(task.StatusToString()).To(Equal(string(entity.StatusInProcess)))
	})

	It("should process task properly", func() {
		context := context.Background()

		task := &entity.Task{
			Method: "POST",
			URL:    "http://test.com",
			RequestHeaders: datatypes.NewJSONType(map[string][]string{
				"test": {"test"},
			}),
			RequestBody: datatypes.NewJSONType(map[string]interface{}{
				"test": "test",
			}),
		}

		responseDTO := clientdto.ResponseDTO{
			Headers: map[string][]string{
				"test": {"test"},
			},
			StatusCode: 200,
			Body:       []byte("{\"test\":\"test\"}"),
		}

		mockHTTPClient.On(
			"Do",
			context,
			mock.MatchedBy(func(actualRequestDTO *clientdto.RequestDTO) bool {
				return actualRequestDTO.Method == task.MethodToString() &&
					actualRequestDTO.URL == task.URL &&
					reflect.DeepEqual(actualRequestDTO.Body, task.RequestBodyToJSON()) &&
					reflect.DeepEqual(actualRequestDTO.Headers, task.RequestHeadersToMap())
			}),
		).Return(&responseDTO, nil)

		mockDALInterface.On(
			"Save",
			context,
			mock.MatchedBy(func(actualTask *entity.Task) bool {
				return actualTask.Status == entity.StatusInProcess &&
					reflect.DeepEqual(actualTask.ResponseHeadersToMap(), responseDTO.Headers) &&
					actualTask.ResponseStatusCodeToInt() == responseDTO.StatusCode &&
					actualTask.ResponseContentLengthToInt() == len(responseDTO.Body)
			}),
		).Return(nil)

		task, err := taskService.ProcessTask(context, task)
		Ω(err).To(BeNil())

		Ω(task.StatusToString()).To(Equal(string(entity.StatusInProcess)))
		Ω(task.ResponseHeadersToMap()).To(Equal(responseDTO.Headers))
		Ω(task.ResponseStatusCodeToInt()).To(Equal(responseDTO.StatusCode))
		Ω(task.ResponseContentLengthToInt()).To(Equal(len(responseDTO.Body)))
	})
})
