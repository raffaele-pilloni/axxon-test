package entity

import (
	"database/sql"
	"encoding/json"
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
	StatusNew       statusTask = "new"
	StatusInProcess statusTask = "in_process"
	StatusDone      statusTask = "done"
	StatusError     statusTask = "error"
	//methods task
	methodGet  methodTask = "GET"
	methodPost methodTask = "POST"
	methodPut  methodTask = "PUT"
)

var allowedMethods = []methodTask{methodGet, methodPost, methodPut}

type Task struct {
	ID                    int
	Status                statusTask
	Method                methodTask
	URL                   string
	RequestHeaders        datatypes.JSONType[map[string][]string]
	RequestBody           datatypes.JSONType[map[string]interface{}]
	ResponseHeaders       datatypes.JSONType[map[string][]string]
	ResponseStatusCode    sql.NullInt64
	ResponseContentLength sql.NullInt64
	CreatedAt             time.Time
	UpdatedAt             time.Time
}

func NewTask(
	method string,
	url string,
	requestHeaders map[string][]string,
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
		Status:         StatusNew,
		Method:         methodTask(method),
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

func (t *Task) MethodToString() string {
	return string(t.Method)
}

func (t *Task) RequestHeadersToMap() map[string][]string {
	return t.RequestHeaders.Data()
}

func (t *Task) RequestBodyToMap() map[string]interface{} {
	return t.RequestBody.Data()
}

func (t *Task) RequestBodyToJSON() []byte {
	if t.RequestBodyToMap() == nil {
		return nil
	}

	bodyMarshalled, _ := json.Marshal(t.RequestBodyToMap())

	return bodyMarshalled
}

func (t *Task) ResponseStatusCodeToInt() int {
	return int(t.ResponseStatusCode.Int64)
}

func (t *Task) ResponseHeadersToMap() map[string][]string {
	return t.ResponseHeaders.Data()
}

func (t *Task) ResponseContentLengthToInt() int {
	return int(t.ResponseContentLength.Int64)
}

func (t *Task) StartProcessing() *Task {
	t.Status = StatusInProcess

	return t
}

func (t *Task) DoneProcessing(
	responseHeaders map[string][]string,
	responseStatusCode int,
	responseContentLength int,
) *Task {
	t.Status = StatusDone

	t.ResponseHeaders = datatypes.NewJSONType(responseHeaders)

	_ = t.ResponseStatusCode.Scan(responseStatusCode)
	_ = t.ResponseContentLength.Scan(responseContentLength)

	return t
}

func (t *Task) ErrorProcessing() *Task {
	t.Status = StatusError

	return t
}
