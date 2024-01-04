package http

import (
	"context"
	"github.com/gorilla/mux"
	pconfigs "github.com/raffaele-pilloni/axxon-test/configs"
	"github.com/raffaele-pilloni/axxon-test/internal/app/http/controller"
	pdal "github.com/raffaele-pilloni/axxon-test/internal/dal"
	"github.com/raffaele-pilloni/axxon-test/internal/repository"
	"github.com/raffaele-pilloni/axxon-test/internal/service"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"net/http"
	"time"
)

type Server struct {
	configs *pconfigs.Configs
	gormDB  *gorm.DB
	server  *http.Server
}

// NewServer application
func NewServer(
	configs *pconfigs.Configs,
) (*Server, error) {
	/************************
	 * Init SQL DB Client *
	 ************************/
	gormDB, err := gorm.Open(mysql.Open(configs.DB.ConnectionString), &gorm.Config{PrepareStmt: true})
	if err != nil {
		return nil, err
	}

	/****************************
	 * Init and inject services *
	 ****************************/

	// Data Access Layer
	dal := pdal.NewDAL(
		gormDB,
		configs.DB.QueryTimeout,
	)

	// Task Repository
	taskRepository := repository.NewTaskRepository(dal)

	// Task Service
	taskService := service.NewTaskService(dal)

	// Task Controller
	taskController := controller.NewTaskController(
		taskRepository,
		taskService,
	)

	/**************************
	 * Init router and server *
	 **************************/
	router := mux.NewRouter()

	router.HandleFunc("/task/{id:[0-9]+}", taskController.GetTask).Methods("GET")
	router.HandleFunc("/task", taskController.CreateTask).Methods("POST")

	server := &http.Server{
		Addr:              configs.Server.Addr,
		Handler:           http.TimeoutHandler(router, configs.Server.HandlerTimeout*time.Second, "request timeout"),
		ReadHeaderTimeout: configs.Server.ReadHeaderTimeout * time.Second,
		ReadTimeout:       configs.Server.ReadTimeout * time.Second,
		WriteTimeout:      configs.Server.WriteTimeout * time.Second,
	}

	return &Server{
		configs: configs,
		gormDB:  gormDB,
		server:  server,
	}, nil
}

// Run application
func (a *Server) Run() error {
	sqlDB, err := a.gormDB.DB()
	if err != nil {
		return err
	}
	defer sqlDB.Close()

	/*********************
	 * Start http server *
	 *********************/
	if err := a.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}

	return nil
}

// Stop application
func (a *Server) Stop() error {
	ctx, cancelCtx := context.WithTimeout(context.Background(), a.configs.Server.ShutdownTimeout*time.Second)
	defer cancelCtx()

	if err := a.server.Shutdown(ctx); err != nil {
		return err
	}

	return nil
}
