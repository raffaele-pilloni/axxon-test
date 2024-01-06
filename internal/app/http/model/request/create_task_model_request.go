package request

type CreateTaskModelRequest struct {
	Method  string                 `json:"method" validate:"required"`
	URL     string                 `json:"url" validate:"required"`
	Headers map[string][]string    `json:"headers" validate:"dive,required"`
	Body    map[string]interface{} `json:"body"`
}
