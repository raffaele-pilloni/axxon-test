package dto

type CreateTaskDTO struct {
	Method  string
	URL     string
	Headers map[string]string
	Body    map[string]interface{}
}
