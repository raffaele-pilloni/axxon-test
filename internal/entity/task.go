package entity

import (
	applicationerror "github.com/raffaele-pilloni/axxon-test/internal/error"
	urlparser "net/url"
	"slices"
	"time"
)

type statusTask string
type methodTask string

const (
	//status task
	statusNew       statusTask = "new"
	statusInProcess statusTask = "in_process"
	statusDone      statusTask = "done"
	statusError     statusTask = "error"
	//methods task
	methodGet  methodTask = "GET"
	methodPost methodTask = "POST"
	methodPut  methodTask = "PUT"
)

var allowedMethods = []methodTask{methodGet, methodPost, methodPut}

type Task struct {
	ID                    int
	Status                statusTask
	Method                string
	URL                   string
	RequestHeaders        map[string]string
	RequestBody           map[string]interface{}
	ResponseStatusCode    string
	ResponseHeaders       map[string]string
	ResponseContentLength int
	CreatedAt             time.Time
	UpdatedAt             time.Time `gorm:"autoUpdateTime"`
}

func NewTask(
	method string,
	url string,
	requestHeaders map[string]string,
	requestBody map[string]interface{},
) (*Task, error) {
	if !slices.Contains(allowedMethods, methodTask(method)) {
		return nil, applicationerror.NewInvalidMethodForTaskError(method)
	}

	urlParsed, err := urlparser.Parse(url)
	if err != nil || urlParsed.Scheme == "" || urlParsed.Host == "" {
		return nil, applicationerror.NewInvalidURLForTaskError(url)
	}

	return &Task{
		Status:         statusNew,
		Method:         method,
		URL:            url,
		RequestHeaders: requestHeaders,
		RequestBody:    requestBody,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}, nil
}
