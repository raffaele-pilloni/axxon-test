package http

import (
	"context"
	"github.com/gorilla/mux"
	pconfig "github.com/raffaele-pilloni/axxon-test/config"
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
	config *pconfig.Config
	gormDB *gorm.DB
	server *http.Server
}

func NewServer(
	config *pconfig.Config,
) (*Server, error) {
	/****************
	 * Init Clients *
	 ****************/

	//Gorm Sql Client
	gormDB, err := gorm.Open(
		mysql.Open(config.DB.ConnectionString),
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
		Timeout: time.Second * config.HTTPClient.RequestTimeout,
	})

	/****************************
	 * Init and inject services *
	 ****************************/

	// Data Access Layer
	dal := db.NewDAL(
		gormDB,
		config.DB.QueryTimeout,
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
		Addr:              config.Server.Addr,
		Handler:           http.TimeoutHandler(router, config.Server.HandlerTimeout*time.Second, "request timeout"),
		ReadHeaderTimeout: config.Server.ReadHeaderTimeout * time.Second,
		ReadTimeout:       config.Server.ReadTimeout * time.Second,
		WriteTimeout:      config.Server.WriteTimeout * time.Second,
	}

	return &Server{
		config: config,
		gormDB: gormDB,
		server: server,
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
	ctx, cancelCtx := context.WithTimeout(context.Background(), a.config.Server.ShutdownTimeout*time.Second)
	defer cancelCtx()

	/************************
	 * Shutdown http server *
	 ************************/
	if err := a.server.Shutdown(ctx); err != nil {
		return err
	}

	return nil
}
