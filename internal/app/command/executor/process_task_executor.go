package executor

import (
	"context"
	"github.com/raffaele-pilloni/axxon-test/internal/entity"
	"github.com/raffaele-pilloni/axxon-test/internal/repository"
	"github.com/raffaele-pilloni/axxon-test/internal/service"
	"log"
	"sync"
	"time"
)

const (
	delayForError           time.Duration = 3
	processTaskExecutorName Name          = "process-task"
)

type ProcessTaskExecutor struct {
	taskRepository          repository.TaskRepositoryInterface
	taskService             service.TaskServiceInterface
	tasksProcessConcurrency int
	executorName            Name
}

func NewProcessTaskExecutor(
	taskRepository repository.TaskRepositoryInterface,
	taskService service.TaskServiceInterface,
	tasksProcessConcurrency int,
) *ProcessTaskExecutor {
	return &ProcessTaskExecutor{
		taskRepository:          taskRepository,
		taskService:             taskService,
		tasksProcessConcurrency: tasksProcessConcurrency,
		executorName:            processTaskExecutorName,
	}
}

func (p *ProcessTaskExecutor) GetName() Name {
	return p.executorName
}

func (p *ProcessTaskExecutor) Run(ctx context.Context, _ []string) error {
	var wg sync.WaitGroup
	defer wg.Wait()

	taskChan := p.readTasksToProcessAsync(ctx, &wg)

	for task := range taskChan {
		select {
		case <-ctx.Done():
			log.Print("Context closed")
			return nil
		default:
		}

		wg.Add(1)
		go func(task *entity.Task) {
			defer wg.Done()

			if _, err := p.taskService.ProcessTask(ctx, task); err != nil {
				log.Printf("Process task failed %v", err)
				time.Sleep(delayForError * time.Second)
				return
			}
		}(task)

	}

	wg.Wait()

	return nil
}

func (p *ProcessTaskExecutor) readTasksToProcessAsync(ctx context.Context, wg *sync.WaitGroup) <-chan *entity.Task {
	tasksChan := make(chan *entity.Task, p.tasksProcessConcurrency)

	go func(tasksChan chan<- *entity.Task) {
		defer wg.Done()

		for {
			select {
			case <-ctx.Done():
				log.Printf("Context closed")
				return
			default:
			}

			tasks, err := p.taskRepository.FindTasksToProcess(ctx, p.tasksProcessConcurrency)
			if err != nil {
				log.Printf("Find tasks to process failed %v", err)
				time.Sleep(delayForError * time.Second)
				continue
			}

			for _, task := range tasks {
				tasksChan <- task
			}
		}
	}(tasksChan)

	return tasksChan
}
