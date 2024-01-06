package command

import (
	"context"
	"fmt"
	pconfigs "github.com/raffaele-pilloni/axxon-test/configs"
	pexecutor "github.com/raffaele-pilloni/axxon-test/internal/app/command/executor"
	"github.com/raffaele-pilloni/axxon-test/internal/client"
	"github.com/raffaele-pilloni/axxon-test/internal/db"
	"github.com/raffaele-pilloni/axxon-test/internal/repository"
	"github.com/raffaele-pilloni/axxon-test/internal/service"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"net/http"
	"time"
)

type executors []pexecutor.Interface

type Dispatcher struct {
	configs   *pconfigs.Configs
	gormDB    *gorm.DB
	executors executors
}

func NewDispatcher(
	configs *pconfigs.Configs,
) (*Dispatcher, error) {
	/****************
	 * Init Clients *
	 ****************/

	//Gorm Sql Client
	gormDB, err := gorm.Open(
		mysql.Open(configs.DB.ConnectionString),
		&gorm.Config{
			PrepareStmt: true,
			NamingStrategy: schema.NamingStrategy{
				SingularTable: true,
			},
		})
	if err != nil {
		return nil, err
	}

	// Http Client
	httpClient := client.NewHttpClient(&http.Client{
		Timeout: time.Second * configs.HTTPClient.RequestTimeout,
	})

	/****************************
	 * Init and inject services *
	 ****************************/

	// Data Access Layer
	dal := db.NewDAL(
		gormDB,
		configs.DB.QueryTimeout,
	)

	// Task Repository
	taskRepository := repository.NewTaskRepository(dal)

	// Task Service
	taskService := service.NewTaskService(dal, httpClient)

	/******************
	 * Init Executors *
	 ******************/

	// Process Task Executor
	processTaskExecutor := pexecutor.NewProcessTaskExecutor(
		taskRepository,
		taskService,
		configs.App.ProcessTaskConcurrency,
	)

	return &Dispatcher{
		configs: configs,
		gormDB:  gormDB,
		executors: []pexecutor.Interface{
			processTaskExecutor,
		},
	}, nil
}

func (a *Dispatcher) Run(ctx context.Context, commandName string, args []string) error {
	sqlDB, err := a.gormDB.DB()
	if err != nil {
		return err
	}
	defer sqlDB.Close()

	/*************************
	 * Find And Run Executor *
	 *************************/
	for _, executor := range a.executors {
		if pexecutor.Name(commandName) == executor.GetName() {
			return executor.Run(ctx, args)
		}
	}

	return fmt.Errorf("executor %s not found", commandName)
}
