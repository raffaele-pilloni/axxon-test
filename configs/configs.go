package configs

import (
	"fmt"
	"github.com/subosito/gotenv"
	"os"
	"strconv"
	"time"
)

type App struct {
	Env                    string
	AppName                string
	ProcessTaskConcurrency int
}

type Server struct {
	Addr              string
	HandlerTimeout    time.Duration
	ReadHeaderTimeout time.Duration
	ReadTimeout       time.Duration
	WriteTimeout      time.Duration
	ShutdownTimeout   time.Duration
}

type DB struct {
	Host             string
	Port             string
	User             string
	Password         string
	Name             string
	ConnectionString string
	QueryTimeout     time.Duration
}

type HTTPClient struct {
	RequestTimeout time.Duration
}

type Configs struct {
	App        *App
	HTTPClient *HTTPClient
	Server     *Server
	DB         *DB
}

func LoadConfigs() (*Configs, error) {
	if err := gotenv.Load(); err != nil {
		return nil, err
	}

	appConfigs, err := loadAppConfigs()
	if err != nil {
		return nil, err
	}

	serverConfigs, err := loadServerConfigs()
	if err != nil {
		return nil, err
	}

	dbConfigs, err := loadDBConfigs()
	if err != nil {
		return nil, err
	}

	httpClientConfigurations, err := loadHTTPClientConfigs()
	if err != nil {
		return nil, err
	}

	return &Configs{
		App:        appConfigs,
		Server:     serverConfigs,
		DB:         dbConfigs,
		HTTPClient: httpClientConfigurations,
	}, nil
}

func loadAppConfigs() (*App, error) {
	processTaskConcurrency, err := strconv.ParseInt(os.Getenv("PROCESS_TASK_CONCURRENCY"), 10, 64)
	if err != nil {
		return nil, err
	}

	return &App{
		Env:                    os.Getenv("ENV"),
		AppName:                os.Getenv("APP_NAME"),
		ProcessTaskConcurrency: int(processTaskConcurrency),
	}, nil
}

func loadServerConfigs() (*Server, error) {
	handlerTimeout, err := strconv.ParseInt(os.Getenv("SERVER_HANDLER_TIMEOUT"), 10, 64)
	if err != nil {
		return nil, err
	}

	readHeaderTimeout, err := strconv.ParseInt(os.Getenv("SERVER_READ_HEADER_TIMEOUT"), 10, 64)
	if err != nil {
		return nil, err
	}

	readTimeout, err := strconv.ParseInt(os.Getenv("SERVER_READ_TIMEOUT"), 10, 64)
	if err != nil {
		return nil, err
	}

	writeTimeout, err := strconv.ParseInt(os.Getenv("SERVER_WRITE_TIMEOUT"), 10, 64)
	if err != nil {
		return nil, err
	}

	shutdownTimeout, err := strconv.ParseInt(os.Getenv("SERVER_SHUTDOWN_TIMEOUT"), 10, 64)
	if err != nil {
		return nil, err
	}

	return &Server{
		Addr:              os.Getenv("SERVER_ADDR"),
		HandlerTimeout:    time.Duration(handlerTimeout),
		ReadHeaderTimeout: time.Duration(readHeaderTimeout),
		ReadTimeout:       time.Duration(readTimeout),
		WriteTimeout:      time.Duration(writeTimeout),
		ShutdownTimeout:   time.Duration(shutdownTimeout),
	}, nil
}

func loadDBConfigs() (*DB, error) {
	queryTimeout, err := strconv.ParseInt(os.Getenv("DB_QUERY_TIMEOUT"), 10, 64)
	if err != nil {
		return nil, err
	}

	return &DB{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		Name:     os.Getenv("DB_NAME"),
		ConnectionString: fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			os.Getenv("DB_USER"),
			os.Getenv("DB_PASSWORD"),
			os.Getenv("DB_HOST"),
			os.Getenv("DB_PORT"),
			os.Getenv("DB_NAME")),
		QueryTimeout: time.Duration(queryTimeout),
	}, nil
}

func loadHTTPClientConfigs() (*HTTPClient, error) {
	requestTimeout, err := strconv.ParseInt(os.Getenv("HTTP_CLIENT_REQUEST_TIMEOUT"), 10, 64)
	if err != nil {
		return nil, err
	}

	return &HTTPClient{
		RequestTimeout: time.Duration(requestTimeout),
	}, nil
}
