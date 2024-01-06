package executor

import (
	"context"
	"github.com/raffaele-pilloni/axxon-test/internal/entity"
	"github.com/raffaele-pilloni/axxon-test/internal/repository"
	"github.com/raffaele-pilloni/axxon-test/internal/service"
	"github.com/raffaele-pilloni/axxon-test/internal/service/dto"
	"log"
	"net/http"
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
	httpClient              *http.Client
	tasksProcessConcurrency int
	executorName            Name
}

func NewProcessTaskExecutor(
	taskRepository repository.TaskRepositoryInterface,
	taskService service.TaskServiceInterface,
	httpClient *http.Client,
	tasksProcessConcurrency int,
) *ProcessTaskExecutor {
	return &ProcessTaskExecutor{
		taskRepository:          taskRepository,
		taskService:             taskService,
		httpClient:              httpClient,
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
			log.Printf("Context closed")
			return nil
		default:
		}

		wg.Add(1)
		go func(task *entity.Task) {
			defer wg.Done()

			p.processTask(ctx, task)
		}(task)

	}

	wg.Wait()

	return nil
}

func (p *ProcessTaskExecutor) processTask(ctx context.Context, task *entity.Task) {
	if _, err := p.taskService.StartTaskProcessing(ctx, task); err != nil {
		log.Printf("Start task processing failed %v", err)
		time.Sleep(delayForError * time.Second)
		return
	}

	if _, err := p.taskService.ErrorTaskProcessing(ctx, task); err != nil {
		log.Printf("Error task processing failed %v", err)
		time.Sleep(delayForError * time.Second)
		return
	}

	if _, err := p.taskService.DoneTaskProcessing(ctx, task, &dto.DoneTaskProcessingDTO{}); err != nil {
		log.Printf("Done task processing failed %v", err)
		time.Sleep(delayForError * time.Second)
		return
	}
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
