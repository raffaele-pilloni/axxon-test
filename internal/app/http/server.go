package http

import (
	"context"
	"github.com/gorilla/mux"
	pconfigs "github.com/raffaele-pilloni/axxon-test/configs"
	"github.com/raffaele-pilloni/axxon-test/internal/app/http/controller"
	pmiddleware "github.com/raffaele-pilloni/axxon-test/internal/app/http/middleware"
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

type Server struct {
	configs *pconfigs.Configs
	gormDB  *gorm.DB
	server  *http.Server
}

func NewServer(
	configs *pconfigs.Configs,
) (*Server, error) {
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
	httpClient := client.NewHTTPClient(&http.Client{
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

	// Task Controller
	taskController := controller.NewTaskController(
		taskRepository,
		taskService,
	)

	//Middleware
	middleware := pmiddleware.NewMiddleware()

	/**************************
	 * Init router and server *
	 **************************/
	router := mux.NewRouter()

	router.HandleFunc("/task/{taskId:[0-9]+}", taskController.GetTask).Methods("GET")
	router.HandleFunc("/task", taskController.CreateTask).Methods("POST")

	router.Use(middleware.Handle)

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

func (a *Server) Stop() error {
	ctx, cancelCtx := context.WithTimeout(context.Background(), a.configs.Server.ShutdownTimeout*time.Second)
	defer cancelCtx()

	/************************
	 * Shutdown http server *
	 ************************/
	if err := a.server.Shutdown(ctx); err != nil {
		return err
	}

	return nil
}
