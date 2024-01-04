package response

type GetTaskModelResponse struct {
	ID             int               `json:"id"`
	Status         string            `json:"status"`
	HttpStatusCode string            `json:"httpStatusCode"`
	Headers        map[string]string `json:"headers"`
	Length         int               `json:"length"`
}
