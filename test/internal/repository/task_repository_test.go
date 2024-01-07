package repository_test

import (
	"context"
	"errors"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/raffaele-pilloni/axxon-test/internal/db"
	"github.com/raffaele-pilloni/axxon-test/internal/entity"
	"github.com/raffaele-pilloni/axxon-test/internal/repository"
	dmock "github.com/raffaele-pilloni/axxon-test/mock"
	"github.com/stretchr/testify/mock"
)

var _ = Describe("Task Repository Tests", func() {
	var (
		mockDALInterface *dmock.DALInterface

		taskRepository *repository.TaskRepository
	)

	BeforeEach(func() {
		mockDALInterface = dmock.NewDALInterface(GinkgoT())

		taskRepository = repository.NewTaskRepository(
			mockDALInterface,
		)
	})

	It("should find task by id properly", func() {
		ctx := context.Background()

		taskID := 1

		mockDALInterface.On(
			"FindByID",
			ctx,
			mock.AnythingOfType("*entity.Task"),
			taskID,
		).Once().Return(nil)

		task, err := taskRepository.FindTaskByID(ctx, taskID)
		Ω(err).To(BeNil())
		Ω(task).ToNot(BeNil())
	})

	It("should find task to process properly", func() {
		ctx := context.Background()

		limit := 10

		mockDALInterface.On(
			"FindBy",
			ctx,
			mock.AnythingOfType("*[]*entity.Task"),
			mock.MatchedBy(func(actualCriteria db.Criteria) bool {
				field, ok := actualCriteria["status"]

				return ok && field == entity.StatusNew
			}),
			"created_at desc",
			limit,
		).Once().Return(nil)

		tasks, err := taskRepository.FindTasksToProcess(ctx, limit)
		Ω(err).To(BeNil())
		Ω(tasks).ToNot(Equal(nil))
	})

	It("should fail find task by id when find by id fail", func() {
		ctx := context.Background()

		taskID := 1

		mockDALInterface.On(
			"FindByID",
			ctx,
			mock.AnythingOfType("*entity.Task"),
			taskID,
		).Once().Return(errors.New("error test"))

		task, err := taskRepository.FindTaskByID(ctx, taskID)
		Ω(err).ToNot(BeNil())
		Ω(task).To(BeNil())

		Ω(err.Error()).To(Equal("error test"))
	})

	It("should fail find task to process when find by fail", func() {
		ctx := context.Background()

		limit := 10

		mockDALInterface.On(
			"FindBy",
			ctx,
			mock.AnythingOfType("*[]*entity.Task"),
			mock.MatchedBy(func(actualCriteria db.Criteria) bool {
				field, ok := actualCriteria["status"]

				return ok && field == entity.StatusNew
			}),
			"created_at desc",
			limit,
		).Once().Return(errors.New("error test"))

		tasks, err := taskRepository.FindTasksToProcess(ctx, limit)
		Ω(err).ToNot(BeNil())
		Ω(tasks).To(BeNil())

		Ω(err.Error()).To(Equal("error test"))
	})
})
