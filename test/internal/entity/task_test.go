package document

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/raffaele-pilloni/axxon-test/internal/entity"
	applicationerror "github.com/raffaele-pilloni/axxon-test/internal/error"
	"time"
)

var _ = Describe("Task Entity Tests", func() {
	It("should create Task entity correctly", func() {
		method := "GET"
		url := "http://test.com"
		requestHeaders := map[string][]string{
			"test": {"test"},
		}

		requestBody := map[string]interface{}{
			"test": "test",
		}

		task, err := entity.NewTask(
			method,
			url,
			requestHeaders,
			requestBody,
		)

		Ω(err).To(BeNil())

		Ω(task.StatusToString()).To(Equal(string(entity.StatusNew)))
		Ω(task.MethodToString()).To(Equal(method))
		Ω(task.URL).To(Equal(url))
		Ω(task.RequestHeadersToMap()).To(Equal(requestHeaders))
		Ω(task.RequestBodyToMap()).To(Equal(requestBody))

		Ω(task.ID).To(BeZero())
		Ω(task.ResponseHeadersToMap()).To(BeEmpty())
		Ω(task.ResponseStatusCodeToInt()).To(BeZero())
		Ω(task.ResponseContentLengthToInt()).To(BeZero())
		Ω(task.CreatedAt).ToNot(Equal(time.Time{}))
		Ω(task.UpdatedAt).ToNot(Equal(time.Time{}))
	})

	It("should start task processing correctly", func() {
		task := &entity.Task{}

		task = task.StartProcessing()

		Ω(task.StatusToString()).To(Equal(string(entity.StatusInProcess)))
	})

	It("should done task processing correctly", func() {
		requestHeaders := map[string][]string{
			"test": {"test"},
		}

		responseStatusCode := 200
		responseContentLength := 40

		task := &entity.Task{}

		task = task.DoneProcessing(
			requestHeaders,
			responseStatusCode,
			responseContentLength,
		)

		Ω(task.StatusToString()).To(Equal(string(entity.StatusDone)))
		Ω(task.ResponseHeadersToMap()).To(Equal(requestHeaders))
		Ω(task.ResponseStatusCodeToInt()).To(Equal(responseStatusCode))
		Ω(task.ResponseContentLengthToInt()).To(Equal(responseContentLength))
	})

	It("should error task processing correctly", func() {
		task := &entity.Task{}

		task = task.ErrorProcessing()

		Ω(task.StatusToString()).To(Equal(string(entity.StatusError)))
	})

	It("should fail create task entity when method is not allowed", func() {
		method := "INVALID"
		url := "http://test.com"
		requestHeaders := map[string][]string{
			"test": {"test"},
		}

		requestBody := map[string]interface{}{
			"test": "test",
		}

		task, err := entity.NewTask(
			method,
			url,
			requestHeaders,
			requestBody,
		)

		Ω(err).ToNot(BeNil())
		Ω(task).To(BeNil())

		invalidMethodForTaskError, ok := err.(*applicationerror.InvalidMethodForTaskError)
		Ω(ok).To(BeTrue())

		Ω(invalidMethodForTaskError.Message()).To(Equal("The method INVALID is not valid for task."))
		Ω(invalidMethodForTaskError.Code()).To(Equal(applicationerror.InvalidMethodForTaskErrorCode))
	})

	It("should fail create task entity when url is invalid", func() {
		method := "POST"
		url := "invalid-test.com"
		requestHeaders := map[string][]string{
			"test": {"test"},
		}

		requestBody := map[string]interface{}{
			"test": "test",
		}

		task, err := entity.NewTask(
			method,
			url,
			requestHeaders,
			requestBody,
		)

		Ω(err).ToNot(BeNil())
		Ω(task).To(BeNil())

		invalidMethodForTaskError, ok := err.(*applicationerror.InvalidURLForTaskError)
		Ω(ok).To(BeTrue())

		Ω(invalidMethodForTaskError.Message()).To(Equal("The url invalid-test.com is not valid for task."))
		Ω(invalidMethodForTaskError.Code()).To(Equal(applicationerror.InvalidURLForTaskErrorCode))
	})
})
