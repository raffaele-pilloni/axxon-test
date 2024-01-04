package configs

import (
	"fmt"
	"github.com/subosito/gotenv"
	"os"
	"strconv"
	"time"
)

type Server struct {
	Addr              string
	HandlerTimeout    time.Duration
	ReadHeaderTimeout time.Duration
	ReadTimeout       time.Duration
	WriteTimeout      time.Duration
}

type DB struct {
	Host             string
	Port             string
	User             string
	Password         string
	Name             string
	ConnectionString string
}

type Configs struct {
	Env     string
	AppName string
	Server  *Server
	DB      *DB
}

func LoadConfigs() (*Configs, error) {
	if err := gotenv.Load(); err != nil {
		return nil, err
	}

	serverConfig, err := loadServerConfigs()
	if err != nil {
		return nil, err
	}

	return &Configs{
		Env:     os.Getenv("ENV"),
		AppName: os.Getenv("APP_NAME"),
		Server:  serverConfig,
		DB:      loadDBConfigs(),
	}, nil
}

func loadServerConfigs() (*Server, error) {
	handlerTimeout, err := strconv.ParseInt(os.Getenv("HANDLER_TIMEOUT"), 10, 64)
	if err != nil {
		return nil, err
	}

	readHeaderTimeout, err := strconv.ParseInt(os.Getenv("READ_HEADER_TIMEOUT"), 10, 64)
	if err != nil {
		return nil, err
	}

	readTimeout, err := strconv.ParseInt(os.Getenv("READ_TIMEOUT"), 10, 64)
	if err != nil {
		return nil, err
	}

	writeTimeout, err := strconv.ParseInt(os.Getenv("WRITE_TIMEOUT"), 10, 64)
	if err != nil {
		return nil, err
	}

	return &Server{
		Addr:              os.Getenv("SERVER_ADDR"),
		HandlerTimeout:    time.Duration(handlerTimeout),
		ReadHeaderTimeout: time.Duration(readHeaderTimeout),
		ReadTimeout:       time.Duration(readTimeout),
		WriteTimeout:      time.Duration(writeTimeout),
	}, nil
}

func loadDBConfigs() *DB {
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
			os.Getenv("DB_NAME"),
		),
	}
}