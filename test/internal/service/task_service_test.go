package service_test

import (
	"context"
	"errors"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	clientdto "github.com/raffaele-pilloni/axxon-test/internal/client/dto"
	"github.com/raffaele-pilloni/axxon-test/internal/entity"
	applicationerror "github.com/raffaele-pilloni/axxon-test/internal/error"
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
		ctx := context.Background()

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
			ctx,
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
		).Once().Return(nil)

		task, err := taskService.CreateTask(ctx, &createRequestDTO)
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
		ctx := context.Background()

		task := &entity.Task{}

		mockDALInterface.On(
			"Save",
			ctx,
			mock.MatchedBy(func(actualTask *entity.Task) bool {
				return actualTask.Status == entity.StatusInProcess
			}),
		).Once().Return(nil)

		task, err := taskService.StartTaskProcessing(ctx, task)
		Ω(err).To(BeNil())

		Ω(task.StatusToString()).To(Equal(string(entity.StatusInProcess)))
	})

	It("should process task properly", func() {
		ctx := context.Background()

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
			ctx,
			mock.MatchedBy(func(actualRequestDTO *clientdto.RequestDTO) bool {
				return actualRequestDTO.Method == task.MethodToString() &&
					actualRequestDTO.URL == task.URL &&
					reflect.DeepEqual(actualRequestDTO.Headers, task.RequestHeadersToMap()) &&
					reflect.DeepEqual(actualRequestDTO.Body, task.RequestBodyToJSON())
			}),
		).Once().Return(&responseDTO, nil)

		mockDALInterface.On(
			"Save",
			ctx,
			mock.MatchedBy(func(actualTask *entity.Task) bool {
				return actualTask.Status == entity.StatusDone &&
					reflect.DeepEqual(actualTask.ResponseHeadersToMap(), responseDTO.Headers) &&
					actualTask.ResponseStatusCodeToInt() == responseDTO.StatusCode &&
					actualTask.ResponseContentLengthToInt() == len(responseDTO.Body)
			}),
		).Once().Return(nil)

		task, err := taskService.ProcessTask(ctx, task)
		Ω(err).To(BeNil())

		Ω(task.StatusToString()).To(Equal(string(entity.StatusDone)))
		Ω(task.ResponseHeadersToMap()).To(Equal(responseDTO.Headers))
		Ω(task.ResponseStatusCodeToInt()).To(Equal(responseDTO.StatusCode))
		Ω(task.ResponseContentLengthToInt()).To(Equal(len(responseDTO.Body)))
	})

	It("should fail create task when method is not allowed", func() {
		ctx := context.Background()

		createRequestDTO := servicedto.CreateTaskDTO{
			Method: "INVALID",
			URL:    "http://test.com",
			Headers: map[string][]string{
				"test": {"test"},
			},
			Body: map[string]interface{}{
				"test": "test",
			},
		}

		task, err := taskService.CreateTask(ctx, &createRequestDTO)
		Ω(err).ToNot(BeNil())
		Ω(task).To(BeNil())

		invalidMethodForTaskError, ok := err.(*applicationerror.InvalidMethodForTaskError)
		Ω(ok).To(BeTrue())

		Ω(invalidMethodForTaskError.Message()).To(Equal("The method INVALID is not valid for task."))
		Ω(invalidMethodForTaskError.Code()).To(Equal(applicationerror.InvalidMethodForTaskErrorCode))
	})

	It("should fail create task when url is invalid", func() {
		ctx := context.Background()

		createRequestDTO := servicedto.CreateTaskDTO{
			Method: "POST",
			URL:    "invalid-test.com",
			Headers: map[string][]string{
				"test": {"test"},
			},
			Body: map[string]interface{}{
				"test": "test",
			},
		}

		task, err := taskService.CreateTask(ctx, &createRequestDTO)
		Ω(err).ToNot(BeNil())
		Ω(task).To(BeNil())

		invalidMethodForTaskError, ok := err.(*applicationerror.InvalidURLForTaskError)
		Ω(ok).To(BeTrue())

		Ω(invalidMethodForTaskError.Message()).To(Equal("The url invalid-test.com is not valid for task."))
		Ω(invalidMethodForTaskError.Code()).To(Equal(applicationerror.InvalidURLForTaskErrorCode))
	})

	It("should fail create task when save fail", func() {
		ctx := context.Background()

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
			ctx,
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
		).Once().Return(errors.New("error test"))

		task, err := taskService.CreateTask(ctx, &createRequestDTO)
		Ω(err).ToNot(BeNil())
		Ω(task).To(BeNil())

		Ω(err.Error()).To(Equal("error test"))
	})

	It("should fail start task processing when save fail", func() {
		ctx := context.Background()

		task := &entity.Task{}

		mockDALInterface.On(
			"Save",
			ctx,
			mock.MatchedBy(func(actualTask *entity.Task) bool {
				return actualTask.Status == entity.StatusInProcess
			}),
		).Once().Return(errors.New("error test"))

		task, err := taskService.StartTaskProcessing(ctx, task)
		Ω(err).ToNot(BeNil())
		Ω(task).To(BeNil())

		Ω(err.Error()).To(Equal("error test"))
	})

	It("should fail start task processing when save fail", func() {
		ctx := context.Background()

		task := &entity.Task{}

		mockDALInterface.On(
			"Save",
			ctx,
			mock.MatchedBy(func(actualTask *entity.Task) bool {
				return actualTask.Status == entity.StatusInProcess
			}),
		).Once().Return(errors.New("error test"))

		task, err := taskService.StartTaskProcessing(ctx, task)
		Ω(err).ToNot(BeNil())
		Ω(task).To(BeNil())

		Ω(err.Error()).To(Equal("error test"))
	})

	It("should make error task processing in process task when request fail", func() {
		ctx := context.Background()

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

		mockHTTPClient.On(
			"Do",
			ctx,
			mock.MatchedBy(func(actualRequestDTO *clientdto.RequestDTO) bool {
				return actualRequestDTO.Method == task.MethodToString() &&
					actualRequestDTO.URL == task.URL &&
					reflect.DeepEqual(actualRequestDTO.Headers, task.RequestHeadersToMap()) &&
					reflect.DeepEqual(actualRequestDTO.Body, task.RequestBodyToJSON())
			}),
		).Once().Return(nil, errors.New("error test"))

		mockDALInterface.On(
			"Save",
			ctx,
			mock.MatchedBy(func(actualTask *entity.Task) bool {
				return actualTask.Status == entity.StatusError
			}),
		).Once().Return(nil)

		task, err := taskService.ProcessTask(ctx, task)
		Ω(err).To(BeNil())

		Ω(task.StatusToString()).To(Equal(string(entity.StatusError)))
	})

	It("should make error task processing in process task when save fail", func() {
		ctx := context.Background()

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
			ctx,
			mock.MatchedBy(func(actualRequestDTO *clientdto.RequestDTO) bool {
				return actualRequestDTO.Method == task.MethodToString() &&
					actualRequestDTO.URL == task.URL &&
					reflect.DeepEqual(actualRequestDTO.Headers, task.RequestHeadersToMap()) &&
					reflect.DeepEqual(actualRequestDTO.Body, task.RequestBodyToJSON())
			}),
		).Once().Return(&responseDTO, nil)

		mockDALInterface.On(
			"Save",
			ctx,
			mock.MatchedBy(func(actualTask *entity.Task) bool {
				return actualTask.Status == entity.StatusDone &&
					reflect.DeepEqual(actualTask.ResponseHeadersToMap(), responseDTO.Headers) &&
					actualTask.ResponseStatusCodeToInt() == responseDTO.StatusCode &&
					actualTask.ResponseContentLengthToInt() == len(responseDTO.Body)
			}),
		).Once().Return(errors.New("error test"))

		mockDALInterface.On(
			"Save",
			ctx,
			mock.MatchedBy(func(actualTask *entity.Task) bool {
				return actualTask.Status == entity.StatusError
			}),
		).Once().Return(nil)

		task, err := taskService.ProcessTask(ctx, task)
		Ω(err).To(BeNil())

		Ω(task.StatusToString()).To(Equal(string(entity.StatusError)))
	})

	It("should fail make error task processing in process task when save fail", func() {
		ctx := context.Background()

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

		mockHTTPClient.On(
			"Do",
			ctx,
			mock.MatchedBy(func(actualRequestDTO *clientdto.RequestDTO) bool {
				return actualRequestDTO.Method == task.MethodToString() &&
					actualRequestDTO.URL == task.URL &&
					reflect.DeepEqual(actualRequestDTO.Headers, task.RequestHeadersToMap()) &&
					reflect.DeepEqual(actualRequestDTO.Body, task.RequestBodyToJSON())
			}),
		).Once().Return(nil, errors.New("error test"))

		mockDALInterface.On(
			"Save",
			ctx,
			mock.MatchedBy(func(actualTask *entity.Task) bool {
				return actualTask.Status == entity.StatusError
			}),
		).Once().Return(errors.New("error test"))

		task, err := taskService.ProcessTask(ctx, task)
		Ω(err).ToNot(BeNil())
		Ω(task).To(BeNil())

		Ω(err.Error()).To(Equal("error test"))
	})

})
