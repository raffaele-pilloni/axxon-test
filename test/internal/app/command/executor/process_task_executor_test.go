package executor_test

import (
	"context"
	"errors"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/raffaele-pilloni/axxon-test/internal/app/command/executor"
	"github.com/raffaele-pilloni/axxon-test/internal/entity"
	dmock "github.com/raffaele-pilloni/axxon-test/mock"
	"github.com/stretchr/testify/mock"
)

var _ = Describe("Process Task Executor Tests", func() {
	var (
		mockTaskRepository      *dmock.TaskRepositoryInterface
		mockTaskService         *dmock.TaskServiceInterface
		tasksProcessConcurrency int

		taskExecutor *executor.ProcessTaskExecutor
	)

	BeforeEach(func() {
		mockTaskRepository = dmock.NewTaskRepositoryInterface(GinkgoT())
		mockTaskService = dmock.NewTaskServiceInterface(GinkgoT())
		tasksProcessConcurrency = 10

		taskExecutor = executor.NewProcessTaskExecutor(
			mockTaskRepository,
			mockTaskService,
			tasksProcessConcurrency,
		)
	})

	It("should get name correctly", func() {
		executorName := taskExecutor.GetName()

		Ω(executorName).To(Equal(executor.Name("process-task")))
	})

	It("should run correctly", func() {
		ctx, cancelCtx := context.WithCancel(context.Background())
		defer cancelCtx()

		firstTasks := []*entity.Task{
			{ID: 1},
		}

		otherTasks := []*entity.Task{
			{ID: 2},
			{ID: 3},
		}

		mockTaskRepository.On(
			"FindTasksToProcess",
			ctx,
			tasksProcessConcurrency,
		).Once().Return(firstTasks, nil)

		mockTaskRepository.On(
			"FindTasksToProcess",
			ctx,
			tasksProcessConcurrency,
		).Once().Return(nil, errors.New("error test"))

		mockTaskRepository.On(
			"FindTasksToProcess",
			ctx,
			tasksProcessConcurrency,
		).Once().Return(otherTasks, nil)

		mockTaskRepository.On(
			"FindTasksToProcess",
			ctx,
			tasksProcessConcurrency,
		).Maybe().Return([]*entity.Task{}, nil)

		mockFirstStartProcessing := mockTaskService.On(
			"StartTaskProcessing",
			ctx,
			firstTasks[0],
		).Once().Return(firstTasks[0], nil)

		mockTaskService.On(
			"StartTaskProcessing",
			ctx,
			otherTasks[0],
		).Once().Return(nil, errors.New("error test"))

		mockThirdStartProcessing := mockTaskService.On(
			"StartTaskProcessing",
			ctx,
			otherTasks[1],
		).Once().Return(otherTasks[1], nil)

		mockTaskService.On(
			"ProcessTask",
			ctx,
			firstTasks[0],
		).Once().Return(nil, errors.New("error test")).NotBefore(mockFirstStartProcessing)

		mockTaskService.On(
			"ProcessTask",
			ctx,
			otherTasks[1],
		).Once().Return(otherTasks[1], nil).Run(func(args mock.Arguments) {
			cancelCtx()
		}).NotBefore(mockThirdStartProcessing)

		_ = taskExecutor.Run(ctx, []string{})
	})
})
