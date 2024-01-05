package entity

import (
	"database/sql"
	applicationerror "github.com/raffaele-pilloni/axxon-test/internal/error"
	"gorm.io/datatypes"
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
	RequestHeaders        datatypes.JSONType[map[string]string]
	RequestBody           datatypes.JSONType[map[string]interface{}]
	ResponseStatusCode    sql.NullString
	ResponseHeaders       datatypes.JSONType[map[string]string]
	ResponseContentLength sql.NullInt64
	CreatedAt             time.Time
	UpdatedAt             time.Time
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
		RequestHeaders: datatypes.NewJSONType(requestHeaders),
		RequestBody:    datatypes.NewJSONType(requestBody),
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}, nil
}

func (t *Task) StatusToString() string {
	return string(t.Status)
}

func (t *Task) RequestHeadersToMap() map[string]string {
	return t.RequestHeaders.Data()
}

func (t *Task) RequestBodyToMap() map[string]interface{} {
	return t.RequestBody.Data()
}

func (t *Task) ResponseStatusCodeToString() string {
	return t.ResponseStatusCode.String
}

func (t *Task) ResponseHeadersToMap() map[string]string {
	return t.ResponseHeaders.Data()
}

func (t *Task) ResponseContentLengthToInt() int {
	return int(t.ResponseContentLength.Int64)
}
