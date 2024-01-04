package response

type ErrorModelResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
