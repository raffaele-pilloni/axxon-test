package http

import (
	"context"
	"github.com/gorilla/mux"
	"github.com/raffaele-pilloni/axxon-test/configs"
	"github.com/raffaele-pilloni/axxon-test/internal/app/http/controller"
	"github.com/raffaele-pilloni/axxon-test/internal/repository"
	"github.com/raffaele-pilloni/axxon-test/internal/service"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"net/http"
	"time"
)

type Server struct {
	gormDB *gorm.DB
	server *http.Server
}

// NewServer application
func NewServer(
	configs *configs.Configs,
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

	// Task Repository
	taskRepository := repository.NewTaskRepository(gormDB)

	// Task Service
	taskService := service.NewTaskService(gormDB)

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
		WriteTimeout:      configs.Server.WriteTimeout * time.Minute,
	}

	return &Server{
		gormDB: gormDB,
		server: server,
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
	backgroundCtx := context.Background()

	serverShutdownCtx, cancelServerShutdownCtx := context.WithTimeout(backgroundCtx, 10*time.Second)
	defer cancelServerShutdownCtx()

	if err := a.server.Shutdown(serverShutdownCtx); err != nil {
		return err
	}

	return nil
}
