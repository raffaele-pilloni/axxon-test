package dto

type CreateTaskDTO struct {
	Method  string                 `json:"method" validate:"required"`
	URL     string                 `json:"url" validate:"required,http_url"`
	Headers map[string]string      `json:"headers"`
	Body    map[string]interface{} `json:"body"`
}
