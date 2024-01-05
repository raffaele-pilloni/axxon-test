package response

type GetTaskModelResponse struct {
	ID             int               `json:"id"`
	Status         string            `json:"status"`
	HTTPStatusCode string            `json:"httpStatusCode,omitempty"`
	Headers        map[string]string `json:"headers,omitempty"`
	Length         int               `json:"length,omitempty"`
}
