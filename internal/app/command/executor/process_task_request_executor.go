package executor

import (
	"context"
	"github.com/raffaele-pilloni/axxon-test/internal/repository"
	"github.com/raffaele-pilloni/axxon-test/internal/service"
)

const (
	processTaskRequestExecutorName Name = "process-task-request"
)

type ProcessTaskRequestExecutor struct {
	taskRepository repository.TaskRepositoryInterface
	taskService    service.TaskServiceInterface
	executorName   Name
}

func NewProcessTaskRequestExecutor(
	taskRepository repository.TaskRepositoryInterface,
	taskService service.TaskServiceInterface,
) *ProcessTaskRequestExecutor {
	return &ProcessTaskRequestExecutor{
		taskRepository: taskRepository,
		taskService:    taskService,
		executorName:   processTaskRequestExecutorName,
	}
}

func (p *ProcessTaskRequestExecutor) GetName() Name {
	return p.executorName
}

func (p *ProcessTaskRequestExecutor) Run(ctx context.Context, args []string) error {
	return nil
}
