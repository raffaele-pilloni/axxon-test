package executor

import (
	"context"
	"github.com/raffaele-pilloni/axxon-test/internal/repository"
	"github.com/raffaele-pilloni/axxon-test/internal/service"
)

const (
	processTaskExecutorName Name = "process-task"
)

type ProcessTaskExecutor struct {
	taskRepository repository.TaskRepositoryInterface
	taskService    service.TaskServiceInterface
	executorName   Name
}

func NewProcessTaskExecutor(
	taskRepository repository.TaskRepositoryInterface,
	taskService service.TaskServiceInterface,
) *ProcessTaskExecutor {
	return &ProcessTaskExecutor{
		taskRepository: taskRepository,
		taskService:    taskService,
		executorName:   processTaskExecutorName,
	}
}

func (p *ProcessTaskExecutor) GetName() Name {
	return p.executorName
}

func (p *ProcessTaskExecutor) Run(ctx context.Context, args []string) error {
	return nil
}
