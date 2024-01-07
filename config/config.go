package config

import (
	"fmt"
	"github.com/subosito/gotenv"
	"os"
	"strconv"
	"time"
)

const (
	defaultProjectDir string = "./"
	dotEnvFile        string = "%s/.env"
	dotEnvFileTest    string = "%s/.env.test"
)

type App struct {
	ProjectDir             string
	Env                    string
	AppName                string
	ServiceName            string
	LogOutputEnabled       bool
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

type Config struct {
	App        *App
	HTTPClient *HTTPClient
	Server     *Server
	DB         *DB
}

func LoadConfig(isTest bool) (*Config, error) {
	projectDir := defaultProjectDir
	if os.Getenv("PROJECT_DIR") != "" {
		projectDir = os.Getenv("PROJECT_DIR")
	}

	dotEnvFiles := []string{fmt.Sprintf(dotEnvFile, projectDir)}
	if isTest {
		dotEnvFiles = []string{
			fmt.Sprintf(dotEnvFileTest, projectDir),
			fmt.Sprintf(dotEnvFile, projectDir),
		}
	}

	if err := gotenv.Load(dotEnvFiles...); err != nil {
		return nil, fmt.Errorf("define PROJECT_DIR env variable or create .env file in current directory: %v", err)
	}

	appConfig, err := loadAppConfig(projectDir)
	if err != nil {
		return nil, err
	}

	serverConfig, err := loadServerConfig()
	if err != nil {
		return nil, err
	}

	dbConfig, err := loadDBConfig()
	if err != nil {
		return nil, err
	}

	httpClientConfiguration, err := loadHTTPClientConfig()
	if err != nil {
		return nil, err
	}

	return &Config{
		App:        appConfig,
		Server:     serverConfig,
		DB:         dbConfig,
		HTTPClient: httpClientConfiguration,
	}, nil
}

func loadAppConfig(projectDir string) (*App, error) {
	processTaskConcurrency, err := strconv.ParseInt(os.Getenv("PROCESS_TASK_CONCURRENCY"), 10, 64)
	if err != nil {
		return nil, err
	}

	logOutputEnabled, err := strconv.ParseBool(os.Getenv("LOG_OUTPUT_ENABLED"))
	if err != nil {
		return nil, err
	}

	return &App{
		ProjectDir:             projectDir,
		Env:                    os.Getenv("ENV"),
		AppName:                os.Getenv("APP_NAME"),
		ServiceName:            os.Getenv("SERVICE_NAME"),
		LogOutputEnabled:       logOutputEnabled,
		ProcessTaskConcurrency: int(processTaskConcurrency),
	}, nil
}

func loadServerConfig() (*Server, error) {
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

func loadDBConfig() (*DB, error) {
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

func loadHTTPClientConfig() (*HTTPClient, error) {
	requestTimeout, err := strconv.ParseInt(os.Getenv("HTTP_CLIENT_REQUEST_TIMEOUT"), 10, 64)
	if err != nil {
		return nil, err
	}

	return &HTTPClient{
		RequestTimeout: time.Duration(requestTimeout),
	}, nil
}
